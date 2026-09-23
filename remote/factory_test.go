package remote

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	integrationsdk "github.com/domainry/domainry-integration-sdk"
)

func TestRemoteBindingUsesStableRuntimeAndDeliveryIdentity(t *testing.T) {
	requirementsSeen := false
	webPushCalls := map[string]bool{}
	accountCalls := map[string]bool{}
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Domainry-Runtime-ID") != "runtime-a" || request.Header.Get("Authorization") != "Bearer secret" {
			t.Fatalf("headers = %#v", request.Header)
		}
		if request.Method == http.MethodGet && request.URL.Path == "/integration/v1/descriptor" {
			_ = json.NewEncoder(response).Encode(integrationsdk.Descriptor{
				ProtocolVersion: integrationsdk.ProtocolVersionV1, Mode: integrationsdk.DeploymentModeSaaS, Audience: "runtime-a",
				Capabilities: []string{"catalog.read", "requirements.connections.sync", "delivery.accept", "delivery.query", "web_push_subscriptions.manage", "management.connections", "connection_accounts.manage", "connection_accounts.read", "connection_accounts.write", "oauth_applications.manage", "oauth_authorizations.manage", "management.secrets", "management.api_keys", "management.external_identities", "management.webhook_subscriptions", "operations.call", "operations.invocations.query", "inbound.webhooks.accept", "inbound.events.query", "subjects.lifecycle"},
			})
			return
		}
		if request.Method == http.MethodPost && request.URL.Path == "/integration/v1/deliveries" {
			var delivery integrationsdk.DeliveryRequest
			if err := json.NewDecoder(request.Body).Decode(&delivery); err != nil || delivery.MessageID != "message-1" || delivery.DeduplicationKey != "record:1:sync" {
				t.Fatalf("delivery=%#v err=%v", delivery, err)
			}
			_ = json.NewEncoder(response).Encode(integrationsdk.DeliveryReceipt{MessageID: delivery.MessageID, InvocationID: "invocation-1", Status: integrationsdk.DeliveryStatusAccepted})
			return
		}
		if request.Method == http.MethodPut && request.URL.Path == "/integration/v1/application-requirements/connections" {
			var body struct {
				Items []integrationsdk.ConnectionRequirement `json:"items"`
			}
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil || len(body.Items) != 1 || body.Items[0].Key != "primary" {
				t.Fatalf("requirements=%#v err=%v", body, err)
			}
			requirementsSeen = true
			response.WriteHeader(http.StatusNoContent)
			return
		}
		if request.Method == http.MethodGet && request.URL.Path == "/integration/v1/web-push/readiness" {
			if request.URL.Query().Get("workspace_id") != "workspace-a" {
				t.Fatalf("readiness query=%s", request.URL.RawQuery)
			}
			webPushCalls["readiness"] = true
			_ = json.NewEncoder(response).Encode(integrationsdk.WebPushReadiness{Ready: true, PublicKey: "public-key", ConnectionKey: "push", Status: "verified"})
			return
		}
		if request.Method == http.MethodGet && request.URL.Path == "/integration/v1/web-push-subscriptions" {
			if request.URL.Query().Get("workspace_id") != "workspace-a" || request.URL.Query().Get("user_id") != "user-a" {
				t.Fatalf("list query=%s", request.URL.RawQuery)
			}
			webPushCalls["list"] = true
			_ = json.NewEncoder(response).Encode(map[string]any{"items": []integrationsdk.WebPushSubscription{{ID: "browser-a", UserID: "user-a", Status: "active"}}})
			return
		}
		if request.Method == http.MethodPut && request.URL.Path == "/integration/v1/web-push-subscriptions/browser-a" {
			var input integrationsdk.WebPushSubscriptionInput
			if err := json.NewDecoder(request.Body).Decode(&input); err != nil || input.Endpoint != "https://push.example/a" {
				t.Fatalf("upsert input=%#v err=%v", input, err)
			}
			webPushCalls["upsert"] = true
			_ = json.NewEncoder(response).Encode(integrationsdk.WebPushSubscription{ID: "browser-a", UserID: "user-a", Status: "active"})
			return
		}
		if request.Method == http.MethodPost && request.URL.Path == "/integration/v1/web-push-subscriptions/browser-a/revoke" {
			webPushCalls["revoke"] = true
			_ = json.NewEncoder(response).Encode(integrationsdk.WebPushSubscription{ID: "browser-a", UserID: "user-a", Status: "revoked"})
			return
		}
		if request.Method == http.MethodPost && request.URL.Path == "/integration/v1/web-push-subscriptions/cleanup-expired" {
			webPushCalls["cleanup"] = true
			_ = json.NewEncoder(response).Encode(map[string]int{"cleaned": 2})
			return
		}
		if request.Method == http.MethodPost && request.URL.Path == "/integration/v1/management/connections/primary/account" {
			if request.URL.Query().Get("workspace_id") != "workspace-a" || request.URL.Query().Get("actor_id") != "admin" {
				t.Fatalf("account registration query=%s", request.URL.RawQuery)
			}
			var input integrationsdk.ConnectionAccountRegistration
			if err := json.NewDecoder(request.Body).Decode(&input); err != nil || input.Scope != integrationsdk.ConnectionAccountScopePersonal || input.OwnerUserID != "user-a" {
				t.Fatalf("account registration=%#v err=%v", input, err)
			}
			accountCalls["register"] = true
			_ = json.NewEncoder(response).Encode(integrationsdk.ConnectionAccount{Key: "primary", Scope: input.Scope, OwnerUserID: input.OwnerUserID, Status: "active", UpdatedAt: "revision-1"})
			return
		}
		if request.Method == http.MethodGet && request.URL.Path == "/integration/v1/connection-accounts" {
			assertConnectionAccountSubject(t, request)
			accountCalls["list"] = true
			_ = json.NewEncoder(response).Encode(map[string]any{"items": []integrationsdk.ConnectionAccount{{Key: "primary", Scope: integrationsdk.ConnectionAccountScopePersonal, OwnerUserID: "user-a", Status: "active", UpdatedAt: "revision-1"}}})
			return
		}
		if request.Method == http.MethodGet && request.URL.Path == "/integration/v1/connection-accounts/primary" {
			assertConnectionAccountSubject(t, request)
			accountCalls["get"] = true
			_ = json.NewEncoder(response).Encode(integrationsdk.ConnectionAccount{Key: "primary", Scope: integrationsdk.ConnectionAccountScopePersonal, OwnerUserID: "user-a", Status: "active", UpdatedAt: "revision-1"})
			return
		}
		if request.Method == http.MethodPost && request.URL.Path == "/integration/v1/connection-accounts/primary/test" {
			assertConnectionAccountSubject(t, request)
			accountCalls["test"] = true
			_ = json.NewEncoder(response).Encode(integrationsdk.ConnectionAccountTestResult{Account: integrationsdk.ConnectionAccount{Key: "primary", Scope: integrationsdk.ConnectionAccountScopePersonal, OwnerUserID: "user-a", Status: "active"}, Operation: "test_connection"})
			return
		}
		if request.Method == http.MethodPost && request.URL.Path == "/integration/v1/connection-accounts/primary/revoke" {
			assertConnectionAccountSubject(t, request)
			var input map[string]string
			if err := json.NewDecoder(request.Body).Decode(&input); err != nil || input["expected_updated_at"] != "revision-1" {
				t.Fatalf("account revoke=%#v err=%v", input, err)
			}
			accountCalls["revoke"] = true
			_ = json.NewEncoder(response).Encode(integrationsdk.ConnectionAccount{Key: "primary", Scope: integrationsdk.ConnectionAccountScopePersonal, OwnerUserID: "user-a", Status: "revoked", UpdatedAt: "revision-2"})
			return
		}
		http.NotFound(response, request)
	}))
	defer server.Close()

	binding, err := NewFactory(Options{BaseURL: server.URL, Token: "secret", HTTPClient: server.Client()}).OpenSaaS(context.Background(), integrationsdk.ApplicationRef{RuntimeID: "runtime-a"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := binding.Requirements().SynchronizeConnections(context.Background(), []integrationsdk.ConnectionRequirement{{Key: "primary", WorkspaceID: "workspace-a", ConnectorKey: "crm", ProviderKey: "probe", Config: json.RawMessage(`{}`)}}); err != nil {
		t.Fatal(err)
	}
	receipt, err := binding.Delivery().Accept(context.Background(), integrationsdk.DeliveryRequest{MessageID: "message-1", DeduplicationKey: "record:1:sync", WorkspaceID: "workspace-a", ConnectorKey: "crm", Operation: "upsert", Payload: json.RawMessage(`{"id":"1"}`)})
	if err != nil || receipt.InvocationID != "invocation-1" {
		t.Fatalf("receipt=%#v err=%v", receipt, err)
	}
	if !requirementsSeen {
		t.Fatal("SaaS connection requirements were not sent")
	}
	webPushBinding, ok := binding.(integrationsdk.WebPushBinding)
	if !ok {
		t.Fatal("SaaS binding has no Web Push owner port")
	}
	webPush := webPushBinding.WebPushSubscriptions()
	readiness, err := webPush.Readiness(context.Background(), "workspace-a")
	if err != nil || !readiness.Ready || readiness.PublicKey != "public-key" {
		t.Fatalf("readiness=%#v err=%v", readiness, err)
	}
	items, err := webPush.List(context.Background(), "workspace-a", "user-a")
	if err != nil || len(items) != 1 || items[0].ID != "browser-a" {
		t.Fatalf("items=%#v err=%v", items, err)
	}
	if _, err := webPush.Upsert(context.Background(), "workspace-a", "user-a", "browser-a", integrationsdk.WebPushSubscriptionInput{Endpoint: "https://push.example/a", P256DH: "key", Auth: "auth"}); err != nil {
		t.Fatal(err)
	}
	if _, err := webPush.Revoke(context.Background(), "workspace-a", "user-a", "browser-a"); err != nil {
		t.Fatal(err)
	}
	if cleaned, err := webPush.CleanupExpired(context.Background(), "workspace-a"); err != nil || cleaned != 2 {
		t.Fatalf("cleaned=%d err=%v", cleaned, err)
	}
	for _, call := range []string{"readiness", "list", "upsert", "revoke", "cleanup"} {
		if !webPushCalls[call] {
			t.Fatalf("missing Web Push SaaS call %q", call)
		}
	}
	accountAdmin := binding.(integrationsdk.ConnectionAccountAdministrationBinding).ConnectionAccountAdministration()
	registered, err := accountAdmin.RegisterConnectionAccount(context.Background(), "workspace-a", "primary", "admin", integrationsdk.ConnectionAccountRegistration{Scope: integrationsdk.ConnectionAccountScopePersonal, OwnerUserID: "user-a"})
	if err != nil || registered.OwnerUserID != "user-a" {
		t.Fatalf("registered=%#v err=%v", registered, err)
	}
	accounts := binding.(integrationsdk.ConnectionAccountsBinding).ConnectionAccounts()
	subject := integrationsdk.ConnectionAccountSubject{WorkspaceID: "workspace-a", UserID: "user-a", Access: integrationsdk.ConnectionAccountAccess{Personal: true, Workspace: true}}
	if listed, err := accounts.ListConnectionAccounts(context.Background(), subject); err != nil || len(listed) != 1 || listed[0].Key != "primary" {
		t.Fatalf("accounts=%#v err=%v", listed, err)
	}
	if value, err := accounts.GetConnectionAccount(context.Background(), subject, "primary"); err != nil || value.OwnerUserID != "user-a" {
		t.Fatalf("account=%#v err=%v", value, err)
	}
	if result, err := accounts.TestConnectionAccount(context.Background(), subject, "primary", integrationsdk.ConnectionTestRequest{Operation: "test_connection"}); err != nil || result.Operation != "test_connection" {
		t.Fatalf("test result=%#v err=%v", result, err)
	}
	if revoked, err := accounts.RevokeConnectionAccount(context.Background(), subject, "primary", "revision-1"); err != nil || revoked.Status != "revoked" {
		t.Fatalf("revoked=%#v err=%v", revoked, err)
	}
	if _, err := accounts.GetConnectionAccount(context.Background(), integrationsdk.ConnectionAccountSubject{WorkspaceID: "workspace-a"}, "primary"); err == nil {
		t.Fatal("incomplete account subject was accepted")
	}
	for _, call := range []string{"register", "list", "get", "test", "revoke"} {
		if !accountCalls[call] {
			t.Fatalf("missing connection account SaaS call %q", call)
		}
	}
}

func assertConnectionAccountSubject(t *testing.T, request *http.Request) {
	t.Helper()
	if request.URL.Query().Get("workspace_id") != "workspace-a" || request.URL.Query().Get("user_id") != "user-a" || request.URL.Query().Get("allow_personal") != "true" || request.URL.Query().Get("allow_workspace") != "true" {
		t.Fatalf("connection account query=%s", request.URL.RawQuery)
	}
}
