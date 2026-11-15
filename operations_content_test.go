package integrationsdk

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	connector "github.com/domainry/domainry-connector-sdk"
)

func contentRequest(t *testing.T) ProviderCallRequest {
	t.Helper()
	content, err := NewProviderCallContent([]byte("# Verified solution\nNo private customer material."), "application/octet-stream")
	if err != nil {
		t.Fatal(err)
	}
	return ProviderCallRequest{RequestID: "publish-a", WorkspaceID: "workspace-a", ConnectionKey: "kb-a", ConnectorKey: "knowledge_base", Operation: "put_document", ActorID: "author-a", RoleKey: "publisher", Source: InvocationSource{ExecutionID: "publish-a", ObjectKey: "knowledge_article", RecordID: "article-a"}, PersistenceMode: ProviderCallPersistenceSensitive, MaskedDestination: "team-kb", Payload: json.RawMessage(`{"doc_id":"article-a"}`), Content: content, ExpectedConnection: &ProviderCallTarget{ConnectionUpdatedAt: "revision-a", ConnectionSHA256: strings.Repeat("b", 64), ProviderKey: "http_api", ContractSHA256: strings.Repeat("a", 64)}}
}

func TestContentCallRequiresExactBoundedBytesAndHostIdentity(t *testing.T) {
	for name, mutate := range map[string]func(*ProviderCallRequest){
		"empty": func(r *ProviderCallRequest) { r.Content.Data = nil },
		"oversize": func(r *ProviderCallRequest) {
			r.Content.Data = bytes.Repeat([]byte("x"), MaximumProviderCallContentBytes+1)
		},
		"hash":                 func(r *ProviderCallRequest) { r.Content.SHA256 = strings.Repeat("a", 64) },
		"case of hash":         func(r *ProviderCallRequest) { r.Content.SHA256 = strings.ToUpper(r.Content.SHA256) },
		"MIME":                 func(r *ProviderCallRequest) { r.Content.ContentType = "application/octet-stream\r\nX-Secret: x" },
		"missing MIME":         func(r *ProviderCallRequest) { r.Content.ContentType = "" },
		"standard persistence": func(r *ProviderCallRequest) { r.PersistenceMode = ProviderCallPersistenceStandard },
		"connection":           func(r *ProviderCallRequest) { r.ConnectionKey = "" },
		"actor":                func(r *ProviderCallRequest) { r.ActorID = "" },
		"role":                 func(r *ProviderCallRequest) { r.RoleKey = "" },
		"source":               func(r *ProviderCallRequest) { r.Source = InvocationSource{} },
		"target":               func(r *ProviderCallRequest) { r.ExpectedConnection = nil },
		"target revision":      func(r *ProviderCallRequest) { r.ExpectedConnection.ConnectionUpdatedAt = "" },
		"target contract":      func(r *ProviderCallRequest) { r.ExpectedConnection.ContractSHA256 = "wrong" },
		"unbounded metadata": func(r *ProviderCallRequest) {
			r.Payload, _ = json.Marshal(strings.Repeat("x", MaximumProviderCallContentBytes))
		},
	} {
		t.Run(name, func(t *testing.T) {
			r := contentRequest(t)
			mutate(&r)
			if r.Validate() == nil {
				t.Fatal("invalid body crossed the public boundary")
			}
		})
	}
	r := contentRequest(t)
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(r)
	var decoded ProviderCallRequest
	if err != nil || json.Unmarshal(raw, &decoded) != nil || decoded.Validate() != nil || !bytes.Equal(decoded.Content.Data, r.Content.Data) {
		t.Fatal("public content contract lost its exact bytes", err)
	}
}

func TestContentConstructorCopiesAndAcceptsItsExactMaximum(t *testing.T) {
	data := bytes.Repeat([]byte("x"), MaximumProviderCallContentBytes)
	c, err := NewProviderCallContent(data, "application/octet-stream")
	if err != nil {
		t.Fatal(err)
	}
	data[0] = 'y'
	if c.Data[0] != 'x' || c.Validate() != nil {
		t.Fatal("caller mutated the frozen content")
	}
}

func TestCurrentTargetBindsScopeEvenWhenOwnerTimestampIsUnchanged(t *testing.T) {
	operation := connector.OperationDescriptor{ConnectorKey: "knowledge_base", ProviderKey: "http_api", Key: "put_document", Mode: connector.ModeCall, ContractSHA256: strings.Repeat("c", 64), Reliability: connector.ReliabilityContract{Effect: connector.EffectWrite, Idempotency: connector.IdempotencyContract{Strategy: connector.IdempotencyNone}, Reconciliation: connector.ReconciliationNone, Compensation: connector.CompensationContract{Mode: connector.CompensationNone}}}
	c := Connection{Key: "team-kb", WorkspaceID: "workspace-a", ConnectorKey: "knowledge_base", ProviderKey: "http_api", Status: "active", UpdatedAt: "current-owner-token", Config: map[string]any{"kb_id": "kb-a"}, SecretRefs: map[string]string{"api_key": "secret:model-key"}}
	before, err := NewProviderCallTarget(c, operation)
	if err != nil {
		t.Fatal(err)
	}
	c.Config["kb_id"] = "kb-b"
	after, err := NewProviderCallTarget(c, operation)
	if err != nil || before.ConnectionUpdatedAt != after.ConnectionUpdatedAt || before.ConnectionSHA256 == after.ConnectionSHA256 {
		t.Fatal("same-millisecond scope change was invisible", err)
	}
	c.Status = "inactive"
	if _, err := NewProviderCallTarget(c, operation); err == nil {
		t.Fatal("inactive connection became a publish target")
	}
}
