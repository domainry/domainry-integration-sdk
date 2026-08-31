package integrationsdk

import "testing"

func TestDescriptorRequiresDeploymentNeutralOwnerCapabilities(t *testing.T) {
	valid := Descriptor{ProtocolVersion: ProtocolVersionV1, Mode: DeploymentModeModule, Capabilities: []string{
		"catalog.read", "requirements.connections.sync", "delivery.accept", "delivery.query", "web_push_subscriptions.manage",
		"management.connections", "management.secrets", "management.api_keys", "management.external_identities", "management.webhook_subscriptions",
		"operations.call", "operations.invocations.query", "inbound.webhooks.accept", "inbound.events.query",
	}}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []Descriptor{
		{},
		{ProtocolVersion: ProtocolVersionV1, Mode: "runtime", Capabilities: valid.Capabilities},
		{ProtocolVersion: ProtocolVersionV1, Mode: DeploymentModeSaaS, Capabilities: []string{"catalog.read"}},
	} {
		if err := invalid.Validate(); err == nil {
			t.Fatalf("invalid descriptor accepted: %#v", invalid)
		}
	}
}

func TestDeliveryRequestRequiresStableHandoffIdentity(t *testing.T) {
	valid := DeliveryRequest{MessageID: "message-1", DeduplicationKey: "record:1:sync", WorkspaceID: "workspace-1", ConnectorKey: "crm", Operation: "upsert", Payload: []byte(`{"id":"1"}`)}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	invalid := valid
	invalid.DeduplicationKey = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("missing deduplication key accepted")
	}
	invalid = valid
	invalid.Payload = []byte(`{`)
	if err := invalid.Validate(); err == nil {
		t.Fatal("invalid JSON payload accepted")
	}
}

func TestApplicationRefRequiresRuntimeIdentity(t *testing.T) {
	if err := (ApplicationRef{}).Validate(); err == nil {
		t.Fatal("empty runtime identity accepted")
	}
	if err := (ApplicationRef{RuntimeID: "runtime-a"}).Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestEventMappingRequirementClosesRuntimeTriggerTargets(t *testing.T) {
	valid := EventMappingRequirement{
		Key: "contact-change", WorkspaceID: "workspace-a", Provider: "crm", TargetType: "action",
		ObjectKey: "contact", ActionKey: "sync", RecordIDPath: "contact.id",
		ActionInput: map[string]string{"name": "contact.name"},
		EventFields: []EventFieldRequirement{{Path: "contact.id", Type: "text"}, {Path: "contact.name", Type: "text"}},
	}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	invalid := valid
	invalid.ActionKey, invalid.ActionKeyPath = "", "target.action"
	if err := invalid.Validate(); err == nil {
		t.Fatal("dynamic action target accepted")
	}
	invalid = valid
	invalid.ActionInput = map[string]string{"name": "undeclared.name"}
	if err := invalid.Validate(); err == nil {
		t.Fatal("undeclared event path accepted")
	}
}
