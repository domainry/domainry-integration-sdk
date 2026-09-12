package integrationsdk

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAccountReadContractBoundsAndSeparatesTrustedSubject(t *testing.T) {
	r := ConnectionAccountReadRequest{RequestID: "r", Operation: "calendar_list", ContractSHA256: strings.Repeat("a", 64), Payload: json.RawMessage(`{}`)}
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"id", "control", "operation", "hash", "upper-hash", "payload", "too-large"} {
		bad := r
		switch name {
		case "id":
			bad.RequestID = " "
		case "control":
			bad.RequestID = "r\n"
		case "operation":
			bad.Operation = ""
		case "hash":
			bad.ContractSHA256 = "a"
		case "upper-hash":
			bad.ContractSHA256 = strings.Repeat("A", 64)
		case "payload":
			bad.Payload = json.RawMessage(`{`)
		case "too-large":
			bad.Payload = json.RawMessage(`"` + strings.Repeat("a", 1<<20) + `"`)
		}
		if bad.Validate() == nil {
			t.Fatal("unbounded account input", name)
		}
	}
	b, _ := json.Marshal(r)
	var fields map[string]any
	_ = json.Unmarshal(b, &fields)
	for _, key := range []string{"workspace_id", "user_id", "access", "scopes", "secrets", "config", "provider_key"} {
		if _, ok := fields[key]; ok {
			t.Fatal("model request contains owner authority", key)
		}
	}
	for _, route := range IntegrationHTTPAdapterContract().Routes {
		if route.Action.Key == ActionIntegrationConnectionAccountsRead || route.Action.Key == ActionIntegrationConnectionAccountsReadAccess {
			if route.Action.EffectClass != "read" || route.Action.Permission == nil || route.Action.CapabilityKey != CapabilityIntegrationAccounts {
				t.Fatal("account read HTTP governance is incorrect", route.Action)
			}
		}
	}
}
