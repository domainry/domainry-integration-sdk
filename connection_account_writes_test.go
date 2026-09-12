package integrationsdk_test

import (
	"encoding/json"
	"strings"
	"testing"

	sdk "github.com/domainry/domainry-integration-sdk"
)

func TestAccountWriteOptionalHostContractAndBounds(t *testing.T) {
	valid := sdk.ConnectionAccountWriteRequest{RequestID: "execution", ExpectedSource: sdk.ConnectionAccountWriteSource{WorkspaceID: "workspace", ConnectionKey: "account", ConnectorKey: "connector", ProviderKey: "provider", AccountUpdatedAt: "revision", Operation: "mail_send", ContractSHA256: strings.Repeat("a", 64)}, Payload: json.RawMessage(`{}`)}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"missing-id", "control-id", "long-id", "missing-source", "wrong-hash", "missing-revision", "payload-limit", "payload-json"} {
		bad := valid
		switch name {
		case "missing-id":
			bad.RequestID = ""
		case "control-id":
			bad.RequestID = "execution\n"
		case "long-id":
			bad.RequestID = strings.Repeat("x", 257)
		case "missing-source":
			bad.ExpectedSource.ConnectionKey = ""
		case "wrong-hash":
			bad.ExpectedSource.ContractSHA256 = strings.Repeat("A", 64)
		case "missing-revision":
			bad.ExpectedSource.AccountUpdatedAt = ""
		case "payload-limit":
			bad.Payload = json.RawMessage(`"` + strings.Repeat("x", 1<<20) + `"`)
		case "payload-json":
			bad.Payload = json.RawMessage(`{`)
		}
		if bad.Validate() == nil {
			t.Fatal("invalid write request accepted", name)
		}
	}
	for _, route := range sdk.IntegrationHTTPAdapterContract().Routes {
		if strings.Contains(route.Pattern(), "/write") {
			t.Fatal("host write port leaked into public browser routes", route.Pattern())
		}
	}
}
