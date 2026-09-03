package integrationsdk

import (
	"encoding/json"
	"strings"
	"testing"

	actioncontract "github.com/domainry/domainry-foundation/action"
)

func TestIntegrationAuthoringDomainOwnsCompleteCapabilitiesWithoutRuntimeOutbox(t *testing.T) {
	domain := IntegrationAuthoringDomain()
	if domain.Key != "integration" {
		t.Fatalf("domain key = %q", domain.Key)
	}
	want := map[string]bool{
		"integration.binding_validation": true,
		"integration.catalog":            true,
		"integration.connection":         true,
		"integration.connection.delete":  true,
		"integration.connection.disable": true,
		"integration.connection.rotate":  true,
		"integration.event.list":         true,
		"integration.event.replay":       true,
		"integration.invocation.list":    true,
		"integration.operation_test":     true,
	}
	seen := map[string]bool{}
	for _, capability := range domain.Capabilities {
		if seen[capability.Key] {
			t.Fatalf("duplicate capability %q", capability.Key)
		}
		seen[capability.Key] = true
		if strings.Contains(capability.Key, "outbox") {
			t.Fatalf("Runtime publication handoff leaked into Integration authoring as %q", capability.Key)
		}
		if len(capability.Examples) < 2 {
			t.Fatalf("capability %q does not publish minimal and representative examples", capability.Key)
		}
		examples := map[string]bool{}
		for _, example := range capability.Examples {
			examples[example.Name] = true
		}
		if !examples["minimal_valid"] || !examples["representative"] {
			t.Fatalf("capability %q examples are incomplete: %#v", capability.Key, examples)
		}
		if len(capability.Sources) == 0 {
			t.Fatalf("capability %q has no source evidence", capability.Key)
		}
		for _, source := range capability.Sources {
			if strings.Contains(source.Path, "domainry-runtime") {
				t.Fatalf("capability %q points back to Runtime source %q", capability.Key, source.Path)
			}
		}
		for _, route := range capability.ConfigurationRoutes {
			if strings.Contains(route, "/operations/integrations") || strings.Contains(route, "/integrations/outbox") {
				t.Fatalf("capability %q publishes retired route %q", capability.Key, route)
			}
		}
	}
	if len(seen) != len(want) {
		t.Fatalf("capability count = %d, want %d", len(seen), len(want))
	}
	for key := range want {
		if !seen[key] {
			t.Errorf("missing Integration capability %q", key)
		}
	}
}

func TestIntegrationAuthoringDomainProjectsOnlyExactActionPermissions(t *testing.T) {
	actions := make(map[string]bool)
	for _, route := range IntegrationHTTPSurfaceContract().Routes {
		action := route.Action
		if action.Authorization.Strategy != actioncontract.AuthorizationAuthenticated || action.Permission == nil {
			continue
		}
		if action.Permission == nil || action.Permission.Key != action.Key {
			t.Fatalf("Integration action %q is not a same-key permission", action.Key)
		}
		actions[action.Key] = true
	}

	retiredAliases := map[string]bool{
		"integration.audit.view":        true,
		"integration.catalog.view":      true,
		"integration.connection.manage": true,
		"integration.connection.test":   true,
		"integration.retry":             true,
	}
	for _, capability := range IntegrationAuthoringDomain().Capabilities {
		if capability.Execution != nil && capability.Execution.PermissionModel != exactActionSameKeyPermissionModel {
			t.Errorf("capability %q permission model = %q", capability.Key, capability.Execution.PermissionModel)
		}
		for _, permission := range capability.Permissions {
			if retiredAliases[permission] {
				t.Errorf("capability %q publishes retired permission alias %q", capability.Key, permission)
			}
			if !actions[permission] {
				t.Errorf("capability %q permission %q is not an Integration action", capability.Key, permission)
			}
		}
	}
}

func TestSpecializeIntegrationAuthoringCapabilityUsesConnectorOwnedDefinition(t *testing.T) {
	definition := json.RawMessage(`{
		"key":"email",
		"providers":[{"key":"smtp","config_fields":[{"key":"host","type":"string","required":true,"validation":{"min_length":3}}],"secret_fields":[{"key":"password","type":"string","required":true}]}],
		"operations":[{"key":"send","input":[{"key":"to","type":"string","required":true}],"output":[{"key":"message_id","type":"string","required":true}]}]
	}`)

	connection, found, err := SpecializeIntegrationAuthoringCapability("integration.connection", definition, AuthoringSelection{ProviderKey: "smtp"})
	if err != nil || !found {
		t.Fatalf("specialize connection: found=%v err=%v", found, err)
	}
	config := connection.InputSchema.Properties["config"]
	if config.Properties["host"].Type != "string" || len(config.Required) != 1 || config.Required[0] != "host" {
		t.Fatalf("Connector provider config did not close the connection schema: %#v", config)
	}
	secrets := connection.InputSchema.Properties["secret_refs"]
	if secrets.Properties["password"].Description == "" {
		t.Fatalf("Connector secret field did not close the secret reference schema: %#v", secrets)
	}

	operation, found, err := SpecializeIntegrationAuthoringCapability("integration.operation_test", definition, AuthoringSelection{OperationKey: "send"})
	if err != nil || !found {
		t.Fatalf("specialize operation: found=%v err=%v", found, err)
	}
	input := operation.InputSchema.Properties["input"]
	if input.Properties["to"].Type != "string" || len(input.Required) != 1 || input.Required[0] != "to" {
		t.Fatalf("Connector operation input did not close the operation schema: %#v", input)
	}
	if operation.OutputSchema.Properties["response"].Properties["message_id"].Type != "string" {
		t.Fatalf("Connector operation output did not close the response schema: %#v", operation.OutputSchema)
	}
}
