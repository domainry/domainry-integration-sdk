package integrationsdk

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

type ConnectionAccountReadOperation struct {
	Operation      string `json:"operation"`
	ContractSHA256 string `json:"contract_sha256"`
}

func (r ConnectionAccountReadOperation) Validate() error {
	if !boundedAccountReadValue(r.Operation, 128) {
		return fmt.Errorf("Integration account read operation is invalid")
	}
	hash, err := hex.DecodeString(r.ContractSHA256)
	if err != nil || len(hash) != 32 || strings.ToLower(r.ContractSHA256) != r.ContractSHA256 {
		return fmt.Errorf("Integration account read contract is invalid")
	}
	return nil
}

type ConnectionAccountReadRequest struct {
	RequestID      string          `json:"request_id"`
	Operation      string          `json:"operation"`
	ContractSHA256 string          `json:"contract_sha256"`
	Payload        json.RawMessage `json:"payload"`
}

func (r ConnectionAccountReadRequest) OperationContract() ConnectionAccountReadOperation {
	return ConnectionAccountReadOperation{Operation: r.Operation, ContractSHA256: r.ContractSHA256}
}
func (r ConnectionAccountReadRequest) Validate() error {
	if err := r.OperationContract().Validate(); err != nil {
		return err
	}
	if !boundedAccountReadValue(r.RequestID, 256) || len(r.Payload) > 1<<20 || !json.Valid(r.Payload) {
		return fmt.Errorf("Integration account read request is invalid")
	}
	return nil
}

func boundedAccountReadValue(value string, limit int) bool {
	return value != "" && len(value) <= limit && strings.TrimSpace(value) == value && strings.IndexFunc(value, unicode.IsControl) < 0
}

// Source binds retained data to the exact authorized account revision and
// operation contract. It contains no credentials. Consumers must reauthorize
// with a fresh host-derived Subject before exposing a saved result.
type ConnectionAccountReadSource struct {
	WorkspaceID      string `json:"workspace_id"`
	ConnectionKey    string `json:"connection_key"`
	ConnectorKey     string `json:"connector_key"`
	ProviderKey      string `json:"provider_key"`
	AccountUpdatedAt string `json:"account_updated_at"`
	Operation        string `json:"operation"`
	ContractSHA256   string `json:"contract_sha256"`
}

type ConnectionAccountReadAccess struct {
	Source            ConnectionAccountReadSource `json:"source"`
	ScopeAlternatives [][]string                  `json:"scope_alternatives"`
}

type ConnectionAccountReadResult struct {
	Source       ConnectionAccountReadSource `json:"source"`
	InvocationID string                      `json:"invocation_id"`
	ReadAt       string                      `json:"read_at"`
	// Sensitive responses are never persisted by Integration. Replaying a
	// succeeded request returns evidence with PayloadAvailable=false.
	PayloadAvailable bool            `json:"payload_available"`
	Payload          json.RawMessage `json:"payload,omitempty"`
}

// Account reads are separate from Management and the fixed account probe.
// Implementations accept only registered call/read operations with declared
// scope requirements, use owner-resolved config/secrets, and recheck account
// ownership/state/grants before returning. Subject is trusted host input;
// browsers and models may not supply identity, access, scopes or credentials.
type ConnectionAccountReads interface {
	AuthorizeConnectionAccountRead(context.Context, ConnectionAccountSubject, string, ConnectionAccountReadOperation) (ConnectionAccountReadAccess, error)
	ReadConnectionAccount(context.Context, ConnectionAccountSubject, string, ConnectionAccountReadRequest) (ConnectionAccountReadResult, error)
}

type ConnectionAccountReadsBinding interface{ ConnectionAccountReads() ConnectionAccountReads }
