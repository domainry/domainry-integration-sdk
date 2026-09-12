package integrationsdk

import (
	"encoding/json"
	"testing"
)

func TestConnectionAccountContractSeparatesOwnershipFromCreatorEvidence(t *testing.T) {
	if (ConnectionAccountSubject{}).Validate() == nil || (ConnectionAccountSubject{WorkspaceID: "workspace", UserID: "user"}).Validate() != nil {
		t.Fatal("connection account subject validation is incomplete")
	}
	if !ConnectionAccountScopePersonal.Valid() || !ConnectionAccountScopeWorkspace.Valid() || ConnectionAccountScope("organization").Valid() {
		t.Fatal("connection account scopes are not closed")
	}
	for name, registration := range map[string]ConnectionAccountRegistration{
		"unknown":                {Scope: "organization"},
		"personal without owner": {Scope: ConnectionAccountScopePersonal},
		"workspace with owner":   {Scope: ConnectionAccountScopeWorkspace, OwnerUserID: "user"},
	} {
		if registration.Validate() == nil {
			t.Fatalf("%s registration passed", name)
		}
	}
	if (ConnectionAccountRegistration{Scope: ConnectionAccountScopePersonal, OwnerUserID: "user"}).Validate() != nil || (ConnectionAccountRegistration{Scope: ConnectionAccountScopeWorkspace}).Validate() != nil {
		t.Fatal("valid registrations were rejected")
	}
	value := ConnectionAccount{Key: "google-personal", WorkspaceID: "workspace", ConnectorKey: "google_workspace", ProviderKey: "google", Name: "Work mail", Scope: ConnectionAccountScopePersonal, OwnerUserID: "user", Status: "active"}
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"secret_refs", "config", "fingerprint", "created_by"} {
		var fields map[string]any
		if json.Unmarshal(raw, &fields) != nil {
			t.Fatal("invalid account JSON")
		}
		if _, exposed := fields[forbidden]; exposed {
			t.Fatalf("safe account projection exposes %q", forbidden)
		}
	}
}
