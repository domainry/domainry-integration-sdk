package remote

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/domainry/domainry-foundation/modulecapability"
	"github.com/domainry/domainry-foundation/modulecapability/contracttest"
	integrationsdk "github.com/domainry/domainry-integration-sdk"
)

func TestRemoteBindingUsesStableRuntimeAndDeliveryIdentity(t *testing.T) {
	requirementsSeen := false
	webPushCalls := map[string]bool{}
	fixture, err := contracttest.NewFixtureBinding("integration")
	if err != nil {
		t.Fatal(err)
	}
	capabilitySummary, err := fixture.CapabilitySummary(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	capabilityHandler, err := modulecapability.NewHTTPHandler(fixture, func(*http.Request) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Domainry-Runtime-ID") != "runtime-a" || request.Header.Get("Authorization") != "Bearer secret" {
			t.Fatalf("headers = %#v", request.Header)
		}
		if strings.HasPrefix(request.URL.Path, modulecapability.HTTPPrefix) {
			capabilityHandler.ServeHTTP(response, request)
			return
		}
		if request.Method == http.MethodPost && request.URL.Path == "/v1/deliveries" {
			var delivery integrationsdk.DeliveryRequest
			if err := json.NewDecoder(request.Body).Decode(&delivery); err != nil || delivery.MessageID != "message-1" || delivery.DeduplicationKey != "record:1:sync" {
				t.Fatalf("delivery=%#v err=%v", delivery, err)
			}
			_ = json.NewEncoder(response).Encode(integrationsdk.DeliveryReceipt{MessageID: delivery.MessageID, InvocationID: "invocation-1", Status: integrationsdk.DeliveryStatusAccepted})
			return
		}
		if request.Method == http.MethodPut && request.URL.Path == "/v1/application-requirements/connections" {
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
		if request.Method == http.MethodGet && request.URL.Path == "/v1/web-push/readiness" {
			if request.URL.Query().Get("workspace_id") != "workspace-a" {
				t.Fatalf("readiness query=%s", request.URL.RawQuery)
			}
			webPushCalls["readiness"] = true
			_ = json.NewEncoder(response).Encode(integrationsdk.WebPushReadiness{Ready: true, PublicKey: "public-key", ConnectionKey: "push", Status: "verified"})
			return
		}
		if request.Method == http.MethodGet && request.URL.Path == "/v1/web-push-subscriptions" {
			if request.URL.Query().Get("workspace_id") != "workspace-a" || request.URL.Query().Get("user_id") != "user-a" {
				t.Fatalf("list query=%s", request.URL.RawQuery)
			}
			webPushCalls["list"] = true
			_ = json.NewEncoder(response).Encode(map[string]any{"items": []integrationsdk.WebPushSubscription{{ID: "browser-a", UserID: "user-a", Status: "active"}}})
			return
		}
		if request.Method == http.MethodPut && request.URL.Path == "/v1/web-push-subscriptions/browser-a" {
			var input integrationsdk.WebPushSubscriptionInput
			if err := json.NewDecoder(request.Body).Decode(&input); err != nil || input.Endpoint != "https://push.example/a" {
				t.Fatalf("upsert input=%#v err=%v", input, err)
			}
			webPushCalls["upsert"] = true
			_ = json.NewEncoder(response).Encode(integrationsdk.WebPushSubscription{ID: "browser-a", UserID: "user-a", Status: "active"})
			return
		}
		if request.Method == http.MethodPost && request.URL.Path == "/v1/web-push-subscriptions/browser-a/revoke" {
			webPushCalls["revoke"] = true
			_ = json.NewEncoder(response).Encode(integrationsdk.WebPushSubscription{ID: "browser-a", UserID: "user-a", Status: "revoked"})
			return
		}
		if request.Method == http.MethodPost && request.URL.Path == "/v1/web-push-subscriptions/cleanup-expired" {
			webPushCalls["cleanup"] = true
			_ = json.NewEncoder(response).Encode(map[string]int{"cleaned": 2})
			return
		}
		http.NotFound(response, request)
	}))
	defer server.Close()

	binding, err := NewFactory(Options{BaseURL: server.URL, Token: "secret", HTTPClient: server.Client(), CapabilityContractSHA256: capabilitySummary.Identity.ContractSHA256}).OpenSaaS(context.Background(), integrationsdk.ApplicationRef{RuntimeID: "runtime-a"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	contracttest.VerifyBinding(t, binding)
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
}
