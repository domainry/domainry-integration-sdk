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
	for _, route := range contract.Routes {
		if seen[route.Pattern] {
			t.Fatalf("duplicate route %q", route.Pattern)
		}
		seen[route.Pattern] = true
		if route.EffectClass == "" || route.IdempotencyDecision == "" || route.AuditClass == "" {
			t.Fatalf("route %q has incomplete governance: %#v", route.Pattern, route)
		}
		operation := contract.OpenAPI[route.Pattern]
		if strings.TrimSpace(operation["operationId"].(string)) == "" || operation["responses"] == nil {
			t.Fatalf("route %q has incomplete OpenAPI: %#v", route.Pattern, operation)
		}
		if strings.Contains(route.Pattern, "/operations/integrations") || strings.Contains(route.Pattern, "/integrations/outbox") {
			t.Fatalf("retired Runtime-owned Integration route returned: %q", route.Pattern)
		}
	}
	if len(seen) != len(contract.OpenAPI) {
		t.Fatalf("route/OpenAPI mismatch: routes=%d operations=%d", len(seen), len(contract.OpenAPI))
	}
}
