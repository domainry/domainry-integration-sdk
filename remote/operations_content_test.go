package remote

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	sdk "github.com/domainry/domainry-integration-sdk"
)

func TestRemoteProviderContentUsesTheExistingAuthenticatedCallWithoutRetry(t *testing.T) {
	body := []byte("# A verified solution\nA small generated document.")
	content, _ := sdk.NewProviderCallContent(body, "application/octet-stream")
	input := sdk.ProviderCallRequest{RequestID: "publication-a", WorkspaceID: "workspace-a", ConnectorKey: "knowledge_base", ConnectionKey: "shared-kb", Operation: "put_document", ActorID: "publisher-a", RoleKey: "publisher", Source: sdk.InvocationSource{ExecutionID: "publication-a"}, PersistenceMode: sdk.ProviderCallPersistenceSensitive, MaskedDestination: "shared-kb", Payload: json.RawMessage(`{"doc_id":"solution-a"}`), Content: content, ExpectedConnection: &sdk.ProviderCallTarget{ConnectionUpdatedAt: "revision-a", ConnectionSHA256: strings.Repeat("b", 64), ProviderKey: "http_api", ContractSHA256: strings.Repeat("a", 64)}}
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != "POST" || r.URL.Path != "/integration/v1/operations/call" || r.Header.Get("Authorization") != "Bearer private-service" || r.Header.Get("X-Domainry-Runtime-ID") != "runtime-a" {
			t.Error("content escaped the existing authenticated owner endpoint")
		}
		var received sdk.ProviderCallRequest
		if json.NewDecoder(r.Body).Decode(&received) != nil || received.Validate() != nil || !bytes.Equal(received.Content.Data, body) || received.Content.SHA256 != content.SHA256 || *received.ExpectedConnection != *input.ExpectedConnection {
			t.Error("content or current connection identity changed over JSON transport")
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("PRIVATE-RESPONSE-BODY"))
	}))
	defer server.Close()
	base, _ := url.Parse(server.URL)
	client := &remoteClient{base: base, client: server.Client(), token: "private-service", runtimeID: "runtime-a"}
	if _, err := client.Call(t.Context(), input); err == nil || strings.Contains(err.Error(), "PRIVATE-RESPONSE-BODY") || calls.Load() != 1 {
		t.Fatal("content failure leaked raw data or replayed the write", err, calls.Load())
	}
	input.Content.SHA256 = "invalid"
	if _, err := client.Call(t.Context(), input); err == nil || calls.Load() != 1 {
		t.Fatal("invalid content reached the remote server")
	}
}

func TestRemoteGuardedCallChecksTheExactReturnedInvocation(t *testing.T) {
	for _, kind := range []string{"valid", "workspace", "connection", "provider", "operation", "fingerprint"} {
		t.Run(kind, func(t *testing.T) {
			content, _ := sdk.NewProviderCallContent([]byte("verified content"), "application/octet-stream")
			input := sdk.ProviderCallRequest{RequestID: "publish-a", WorkspaceID: "workspace-a", ConnectorKey: "knowledge_base", ConnectionKey: "team-kb", Operation: "put_document", ActorID: "author", RoleKey: "publisher", Source: sdk.InvocationSource{ExecutionID: "publish-a"}, PersistenceMode: sdk.ProviderCallPersistenceSensitive, MaskedDestination: "team-kb", Payload: json.RawMessage(`{}`), Content: content, ExpectedConnection: &sdk.ProviderCallTarget{ConnectionUpdatedAt: "revision", ConnectionSHA256: strings.Repeat("b", 64), ProviderKey: "http_api", ContractSHA256: strings.Repeat("a", 64)}}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				v := sdk.Invocation{InvocationID: "invocation-a", WorkspaceID: input.WorkspaceID, ConnectorKey: input.ConnectorKey, ConnectionKey: input.ConnectionKey, ProviderKey: "http_api", Operation: input.Operation, Status: "succeeded", Metadata: map[string]any{"guarded_call_sha256": input.GuardedCallSHA256()}}
				switch kind {
				case "workspace":
					v.WorkspaceID = "other"
				case "connection":
					v.ConnectionKey = "other"
				case "provider":
					v.ProviderKey = "other"
				case "operation":
					v.Operation = "other"
				case "fingerprint":
					v.Metadata["guarded_call_sha256"] = strings.Repeat("c", 64)
				}
				_ = json.NewEncoder(w).Encode(sdk.ProviderCallResult{Invocation: v})
			}))
			defer server.Close()
			base, _ := url.Parse(server.URL)
			c := &remoteClient{base: base, client: server.Client(), token: "service", runtimeID: "runtime"}
			result, err := c.Call(t.Context(), input)
			if kind == "valid" {
				if err != nil || result.Invocation.InvocationID != "invocation-a" {
					t.Fatal(result, err)
				}
			} else if err == nil || result.Invocation.InvocationID != "" {
				t.Fatal("foreign receipt was accepted", result, err)
			}
		})
	}
}
