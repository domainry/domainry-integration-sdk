package integrationsdk

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime"
	"strings"

	connector "github.com/domainry/domainry-connector-sdk"
)

// MaximumProviderCallContentBytes bounds small, generated source documents.
// Large transfers remain owned by their streaming/file or Data Exchange port.
const MaximumProviderCallContentBytes = 256 << 10

// ProviderCallContent travels only with a sensitive synchronous Provider call.
// Integration verifies the declared hash and gives the existing Connector
// operation a fresh reader plus exact byte count, MIME and SHA-256. Neither
// these bytes nor the Provider response are retained in Invocation metadata.
// JSON transports encode Data as base64; this is not a stored blob reference.
type ProviderCallContent struct {
	Data        []byte `json:"data"`
	SHA256      string `json:"sha256"`
	ContentType string `json:"content_type"`
}

// ProviderCallTarget freezes the current owner Connection and exact operation
// before a content effect. ConnectionUpdatedAt is the existing Connection's
// opaque revision token; callers copy it unchanged, never synthesize a time.
// Integration checks it again after secret resolution immediately before I/O.
type ProviderCallTarget struct {
	ConnectionUpdatedAt string `json:"connection_updated_at"`
	ConnectionSHA256    string `json:"connection_sha256"`
	ProviderKey         string `json:"provider_key"`
	ContractSHA256      string `json:"contract_sha256"`
}

func (v ProviderCallTarget) Validate() error {
	if !boundedAccountWriteValue(v.ConnectionUpdatedAt, 1024) || !connector.ValidSHA256(v.ConnectionSHA256) || !boundedAccountWriteValue(v.ProviderKey, 255) || !connector.ValidSHA256(v.ContractSHA256) {
		return fmt.Errorf("Integration provider target expectation is invalid")
	}
	return nil
}

// GuardedCallSHA256 identifies the exact immutable host effect without storing
// its content. Integration includes it in the sensitive Invocation receipt.
func (r ProviderCallRequest) GuardedCallSHA256() string {
	if r.ExpectedConnection == nil {
		return ""
	}
	var content any
	if r.Content != nil {
		content = struct {
			Bytes               int
			SHA256, ContentType string
		}{len(r.Content.Data), r.Content.SHA256, r.Content.ContentType}
	}
	raw, err := json.Marshal(struct {
		RequestID, WorkspaceID, ActorID, RoleKey, ConnectorKey, ConnectionKey, Operation, MaskedDestination string
		Source                                                                                              InvocationSource
		Target                                                                                              *ProviderCallTarget
		Payload                                                                                             json.RawMessage
		Content                                                                                             any
	}{r.RequestID, r.WorkspaceID, r.ActorID, r.RoleKey, r.ConnectorKey, r.ConnectionKey, r.Operation, r.MaskedDestination, r.Source, r.ExpectedConnection, r.Payload, content})
	if err != nil {
		return ""
	}
	digest := sha256.Sum256(raw)
	return hex.EncodeToString(digest[:])
}

// NewProviderCallTarget uses the exact current owner projection, including its
// configuration and secret references. The digest also detects a scope edit
// made within the same UpdatedAt millisecond. It never includes secret values.
func NewProviderCallTarget(c Connection, operation connector.OperationDescriptor) (*ProviderCallTarget, error) {
	if c.WorkspaceID == "" || c.Key == "" || c.Status != "active" || c.ConnectorKey != operation.ConnectorKey || c.ProviderKey != operation.ProviderKey || operation.Mode != connector.ModeCall || operation.Validate() != nil {
		return nil, fmt.Errorf("Integration provider target does not match the current connection")
	}
	config, refs := c.Config, c.SecretRefs
	if len(config) == 0 {
		config = nil
	}
	if len(refs) == 0 {
		refs = nil
	}
	raw, err := json.Marshal([]any{c.WorkspaceID, c.Key, c.ConnectorKey, c.ProviderKey, c.Status, c.UpdatedAt, config, refs})
	if err != nil {
		return nil, fmt.Errorf("Integration provider target configuration is invalid")
	}
	digest := sha256.Sum256(raw)
	value := &ProviderCallTarget{ConnectionUpdatedAt: c.UpdatedAt, ConnectionSHA256: hex.EncodeToString(digest[:]), ProviderKey: c.ProviderKey, ContractSHA256: operation.ContractSHA256}
	return value, value.Validate()
}

func NewProviderCallContent(data []byte, contentType string) (*ProviderCallContent, error) {
	if len(data) == 0 || len(data) > MaximumProviderCallContentBytes {
		return nil, fmt.Errorf("Integration provider content is invalid or exceeds its bound")
	}
	digest := sha256.Sum256(data)
	value := &ProviderCallContent{Data: append([]byte(nil), data...), SHA256: hex.EncodeToString(digest[:]), ContentType: contentType}
	if err := value.Validate(); err != nil {
		return nil, err
	}
	return value, nil
}

func (c ProviderCallContent) Validate() error {
	if len(c.Data) == 0 || len(c.Data) > MaximumProviderCallContentBytes || len(c.ContentType) > 128 || c.ContentType != strings.TrimSpace(c.ContentType) || strings.ContainsAny(c.ContentType, "\r\n\x00") {
		return fmt.Errorf("Integration provider content is invalid or exceeds its bound")
	}
	if _, _, err := mime.ParseMediaType(c.ContentType); err != nil {
		return fmt.Errorf("Integration provider content type is invalid")
	}
	digest := sha256.Sum256(c.Data)
	if c.SHA256 != hex.EncodeToString(digest[:]) {
		return fmt.Errorf("Integration provider content SHA-256 does not match")
	}
	return nil
}
