package integrationsdk

import (
	"encoding/json"
	"testing"
)

func TestAuthoringEvidenceSchemasMatchBusinessIDWireContracts(t *testing.T) {
	for _, item := range []struct {
		name, idKey string
		schema      AuthoringSchema
		value       any
	}{
		{"invocation", "invocation_id", integrationInvocationSchema(), Invocation{InvocationID: "invocation-business", WorkspaceID: "workspace", ConnectorKey: "connector", ProviderKey: "provider", ConnectionKey: "connection", Operation: "lookup", Status: "succeeded", DurationMS: 1, RequestRef: "request", ResponseRef: "response", Error: "test", EventID: "event", ObjectKey: "account", RecordID: "account-business", WorkflowExecutionID: "workflow", Metadata: map[string]any{"test": true}, CreatedAt: "2026-10-22T00:00:00Z", UpdatedAt: "2026-10-22T00:00:00Z"}},
		{"event", "event_id", integrationEventSchema(), Event{EventID: "event-business", WorkspaceID: "workspace", ConnectorKey: "connector", ConnectionKey: "connection", Provider: "provider", EventType: "changed", ExternalID: "external", Status: "received", Payload: json.RawMessage(`{}`), Error: "test", AttemptCount: 1, NextRetryAt: "2026-10-22T00:00:00Z", LastAttemptAt: "2026-10-22T00:00:00Z", ReceivedAt: "2026-10-22T00:00:00Z", UpdatedAt: "2026-10-22T00:00:00Z", Execution: &RuntimeExecutionReceipt{}}},
	} {
		t.Run(item.name, func(t *testing.T) {
			raw, err := json.Marshal(item.value)
			if err != nil {
				t.Fatal(err)
			}
			var wire map[string]json.RawMessage
			if err := json.Unmarshal(raw, &wire); err != nil {
				t.Fatal(err)
			}
			if _, found := item.schema.Properties[item.idKey]; !found || wire[item.idKey] == nil {
				t.Fatalf("schema and wire do not declare %s: %s", item.idKey, raw)
			}
			if _, found := item.schema.Properties["id"]; found || wire["id"] != nil {
				t.Fatalf("invented id alias in %s contract", item.name)
			}
			for _, required := range item.schema.Required {
				if wire[required] == nil {
					t.Errorf("required schema property %s missing from wire", required)
				}
			}
			for key := range wire {
				if _, found := item.schema.Properties[key]; !found {
					t.Errorf("wire property %s absent from closed schema", key)
				}
			}
		})
	}
	for _, capability := range IntegrationAuthoringDomain().Capabilities {
		if capability.Key != "integration.event.replay" {
			continue
		}
		for _, output := range capability.OutputVariables {
			if output.Name == "event_id" && output.JSONPointer != "/event_id" {
				t.Fatalf("event replay points to %s instead of /event_id", output.JSONPointer)
			}
		}
	}
}
