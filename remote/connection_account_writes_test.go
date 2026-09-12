package remote

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	sdk "github.com/domainry/domainry-integration-sdk"
)

func remoteWriteRequest() sdk.ConnectionAccountWriteRequest {
	return sdk.ConnectionAccountWriteRequest{RequestID: "execution-1", ExpectedSource: sdk.ConnectionAccountWriteSource{WorkspaceID: "workspace", ConnectionKey: "account", ConnectorKey: "connector", ProviderKey: "provider", AccountUpdatedAt: "revision", Operation: "mail_send", ContractSHA256: strings.Repeat("a", 64)}, Payload: json.RawMessage(`{"message":{}}`)}
}

func TestRemoteAccountWriteUsesExactHostRequestAndReadOnlyReceiptPath(t *testing.T) {
	input := remoteWriteRequest()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != "POST" || r.Header.Get("Authorization") != "Bearer private-service" || r.Header.Get("X-Domainry-Runtime-ID") != "runtime" || r.Header.Get("Cookie") != "" || r.URL.Query().Get("user_id") != "user" || r.URL.Query().Get("allow_personal") != "true" {
			t.Error("untrusted request envelope")
		}
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/write-access") {
			var op sdk.ConnectionAccountWriteOperation
			if json.NewDecoder(r.Body).Decode(&op) != nil || op != input.ExpectedSource.OperationContract() {
				t.Error("wrong contract")
			}
			_ = json.NewEncoder(w).Encode(sdk.ConnectionAccountWriteAccess{Source: input.ExpectedSource, ScopeAlternatives: [][]string{{"send"}}})
			return
		}
		if r.URL.Path != "/integration/v1/connection-accounts/account/write" && r.URL.Path != "/integration/v1/connection-accounts/account/write-receipt" {
			t.Error("wrong owner endpoint", r.URL.Path)
		}
		var in sdk.ConnectionAccountWriteRequest
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.RequestID != input.RequestID || in.ExpectedSource != input.ExpectedSource || string(in.Payload) != string(input.Payload) {
			t.Error("request identity changed")
		}
		_ = json.NewEncoder(w).Encode(sdk.ConnectionAccountWriteResult{Source: input.ExpectedSource, Status: sdk.AccountWriteSucceeded, InvocationID: "account-write:original", RecordedAt: "2026-09-12T00:00:00Z", Receipt: json.RawMessage(`{"status":"accepted"}`)})
	}))
	defer server.Close()
	base, _ := url.Parse(server.URL)
	client := server.Client()
	client.Jar, _ = cookiejar.New(nil)
	client.Jar.SetCookies(base, []*http.Cookie{{Name: "browser-session", Value: "private-cookie"}})
	c := &remoteClient{base: base, client: client, token: "private-service", runtimeID: "runtime"}
	subject := sdk.ConnectionAccountSubject{WorkspaceID: "workspace", UserID: "user", Access: sdk.ConnectionAccountAccess{Personal: true}}
	if _, err := c.AuthorizeConnectionAccountWrite(t.Context(), subject, "account", input.ExpectedSource.OperationContract()); err != nil {
		t.Fatal(err)
	}
	if _, err := c.WriteConnectionAccount(t.Context(), subject, "account", input); err != nil {
		t.Fatal(err)
	}
	if _, err := c.ReadConnectionAccountWriteReceipt(t.Context(), subject, "account", input); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 3 || client.Jar == nil {
		t.Fatal("unexpected calls or mutated host client")
	}
}

func TestRemoteAccountWriteRejectsRedirectOversizeLossAndMalformedResponses(t *testing.T) {
	for _, scenario := range []string{"redirect", "lost", "oversize", "trailing", "foreign-source", "private-error", "invalid-status"} {
		t.Run(scenario, func(t *testing.T) {
			var calls atomic.Int32
			input := remoteWriteRequest()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				switch scenario {
				case "redirect":
					w.Header().Set("Location", "/other")
					w.WriteHeader(307)
				case "lost":
					conn, _, err := w.(http.Hijacker).Hijack()
					if err != nil {
						t.Error(err)
					} else {
						_ = conn.Close()
					}
				case "oversize":
					_, _ = w.Write([]byte(strings.Repeat(" ", 65<<10)))
				case "trailing":
					_, _ = w.Write([]byte(`{} {}`))
				case "foreign-source":
					_ = json.NewEncoder(w).Encode(sdk.ConnectionAccountWriteResult{Status: sdk.AccountWriteNotFound})
				case "private-error":
					w.WriteHeader(503)
					_, _ = w.Write([]byte("PRIVATE-PROVIDER-CONTENT"))
				case "invalid-status":
					_ = json.NewEncoder(w).Encode(sdk.ConnectionAccountWriteResult{Source: input.ExpectedSource, Status: "delivered"})
				}
			}))
			defer server.Close()
			base, _ := url.Parse(server.URL)
			c := &remoteClient{base: base, client: server.Client(), token: "private", runtimeID: "runtime"}
			_, err := c.WriteConnectionAccount(t.Context(), sdk.ConnectionAccountSubject{WorkspaceID: "workspace", UserID: "user", Access: sdk.ConnectionAccountAccess{Personal: true}}, "account", input)
			if err == nil || calls.Load() != 1 || strings.Contains(err.Error(), "PRIVATE-PROVIDER-CONTENT") {
				t.Fatal("unsafe write recovery", err, calls.Load())
			}
		})
	}
}
