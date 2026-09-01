package integrationsdk

import (
	"strings"
	"testing"
)

func TestIntegrationHTTPSurfaceContractIsCompleteAndSourceOwned(t *testing.T) {
	contract := IntegrationHTTPSurfaceContract()
	if contract.ContractVersion != IntegrationHTTPSurfaceContractVersion || contract.Owner != "integration" || contract.Name == "" {
		t.Fatalf("incomplete Integration HTTP contract: %#v", contract)
	}
	seen := map[string]bool{}
	browserRoutes := 0
	for _, route := range contract.Routes {
		pattern := route.Pattern()
		if seen[pattern] {
			t.Fatalf("duplicate route %q", pattern)
		}
		seen[pattern] = true
		if route.Action.EffectClass == "" || route.Action.IdempotencyDecision == "" || route.Action.AuditClass == "" {
			t.Fatalf("route %q has incomplete governance: %#v", pattern, route)
		}
		operation := contract.OpenAPI[pattern]
		if strings.TrimSpace(operation["operationId"].(string)) == "" || operation["responses"] == nil {
			t.Fatalf("route %q has incomplete OpenAPI: %#v", pattern, operation)
		}
		if strings.Contains(pattern, "/operations/integrations") || strings.Contains(pattern, "/integrations/outbox") {
			t.Fatalf("retired Runtime-owned Integration route returned: %q", pattern)
		}
		if route.Action.Authorization.Strategy == "anonymous_protocol" {
			if route.BrowserClientPackage != "" || route.BrowserClientMethod != "" {
				t.Fatalf("anonymous ingress must not claim a browser client: %#v", route)
			}
		} else {
			browserRoutes++
			if route.BrowserClientPackage != "@domainry/integration-client" || strings.TrimSpace(route.BrowserClientMethod) == "" {
				t.Fatalf("browser route %q has no Integration owner client: %#v", pattern, route)
			}
			if operation["x-domainry-owner-client-package"] != route.BrowserClientPackage || operation["x-domainry-owner-client-method"] != route.BrowserClientMethod {
				t.Fatalf("route %q has incomplete owner client OpenAPI extension: %#v", pattern, operation)
			}
		}
	}
	if browserRoutes != 37 {
		t.Fatalf("expected 37 browser-owned routes, got %d", browserRoutes)
	}
	if len(seen) != len(contract.OpenAPI) {
		t.Fatalf("route/OpenAPI mismatch: routes=%d operations=%d", len(seen), len(contract.OpenAPI))
	}
}
