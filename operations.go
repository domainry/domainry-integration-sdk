package integrationsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Invocation is Integration-owned provider execution evidence exposed to
// Runtime only through owner queries. Runtime never persists this projection.
type Invocation struct {
	ID                  string         `json:"id"`
	WorkspaceID         string         `json:"workspace_id,omitempty"`
	ConnectorKey        string         `json:"connector_key"`
	ProviderKey         string         `json:"provider_key,omitempty"`
	ConnectionKey       string         `json:"connection_key,omitempty"`
	Operation           string         `json:"operation"`
	Status              string         `json:"status"`
	DurationMS          int64          `json:"duration_ms,omitempty"`
	RequestRef          string         `json:"request_ref,omitempty"`
	ResponseRef         string         `json:"response_ref,omitempty"`
	Error               string         `json:"error,omitempty"`
	EventID             string         `json:"event_id,omitempty"`
	ObjectKey           string         `json:"object_key,omitempty"`
	RecordID            string         `json:"record_id,omitempty"`
	WorkflowExecutionID string         `json:"workflow_execution_id,omitempty"`
	Metadata            map[string]any `json:"metadata,omitempty"`
	CreatedAt           string         `json:"created_at,omitempty"`
	UpdatedAt           string         `json:"updated_at,omitempty"`
}

type InvocationQuery struct {
	WorkspaceID   string `json:"workspace_id"`
	ConnectorKey  string `json:"connector_key,omitempty"`
	ConnectionKey string `json:"connection_key,omitempty"`
	Operation     string `json:"operation,omitempty"`
	Status        string `json:"status,omitempty"`
	CreatedFrom   string `json:"created_from,omitempty"`
	Limit         int    `json:"limit,omitempty"`
}

type ProviderCallRequest struct {
	RequestID     string          `json:"request_id"`
	WorkspaceID   string          `json:"workspace_id"`
	ConnectorKey  string          `json:"connector_key"`
	ConnectionKey string          `json:"connection_key,omitempty"`
	Operation     string          `json:"operation"`
	Payload       json.RawMessage `json:"payload"`
	ActorID       string          `json:"actor_id,omitempty"`
	RoleKey       string          `json:"role_key,omitempty"`
}

func (r ProviderCallRequest) Validate() error {
	for name, value := range map[string]string{"request_id": r.RequestID, "workspace_id": r.WorkspaceID, "connector_key": r.ConnectorKey, "operation": r.Operation} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("Integration provider call %s is required", name)
		}
	}
	if !json.Valid(r.Payload) {
		return fmt.Errorf("Integration provider call payload must be valid JSON")
	}
	return nil
}

type ProviderCallResult struct {
	Invocation Invocation      `json:"invocation"`
	Response   json.RawMessage `json:"response,omitempty"`
}

type Event struct {
	ID            string                   `json:"id"`
	WorkspaceID   string                   `json:"workspace_id,omitempty"`
	ConnectorKey  string                   `json:"connector_key,omitempty"`
	ConnectionKey string                   `json:"connection_key,omitempty"`
	Provider      string                   `json:"provider"`
	EventType     string                   `json:"event_type"`
	ExternalID    string                   `json:"external_id"`
	Status        string                   `json:"status"`
	Payload       json.RawMessage          `json:"payload,omitempty"`
	Error         string                   `json:"error,omitempty"`
	AttemptCount  int                      `json:"attempt_count"`
	NextRetryAt   string                   `json:"next_retry_at,omitempty"`
	LastAttemptAt string                   `json:"last_attempt_at,omitempty"`
	ReceivedAt    string                   `json:"received_at"`
	UpdatedAt     string                   `json:"updated_at"`
	Execution     *RuntimeExecutionReceipt `json:"execution,omitempty"`
}

type EventQuery struct {
	WorkspaceID string `json:"workspace_id"`
	Provider    string `json:"provider,omitempty"`
	EventType   string `json:"event_type,omitempty"`
	Status      string `json:"status,omitempty"`
	Limit       int    `json:"limit,omitempty"`
}

// WebhookRequest is the deployment-neutral ingress envelope. Integration
// resolves the connection and secret material before calling the Provider's
// WebhookVerifier; Runtime never sees those credentials.
type WebhookRequest struct {
	WorkspaceID   string              `json:"workspace_id"`
	ConnectorKey  string              `json:"connector_key"`
	ConnectionKey string              `json:"connection_key"`
	Headers       map[string][]string `json:"headers,omitempty"`
	Query         map[string][]string `json:"query,omitempty"`
	Body          []byte              `json:"body"`
	ReceivedAt    time.Time           `json:"received_at"`
}

func (r WebhookRequest) Validate() error {
	for name, value := range map[string]string{"workspace_id": r.WorkspaceID, "connector_key": r.ConnectorKey, "connection_key": r.ConnectionKey} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("Integration webhook %s is required", name)
		}
	}
	if r.ReceivedAt.IsZero() {
		return fmt.Errorf("Integration webhook received_at is required")
	}
	return nil
}

type WebhookReceipt struct {
	Event     Event  `json:"event"`
	Challenge string `json:"challenge,omitempty"`
	Format    string `json:"challenge_format,omitempty"`
}

type TriggerTarget struct {
	Type        string         `json:"type"`
	WorkflowKey string         `json:"workflow_key,omitempty"`
	ObjectKey   string         `json:"object_key,omitempty"`
	RecordID    string         `json:"record_id,omitempty"`
	ActionKey   string         `json:"action_key,omitempty"`
	Input       map[string]any `json:"input,omitempty"`
}

type TriggerPrincipal struct {
	ActorID      string `json:"actor_id,omitempty"`
	RoleKey      string `json:"role_key,omitempty"`
	ExternalName string `json:"external_name,omitempty"`
}

type TriggerRequest struct {
	EventID        string           `json:"event_id"`
	WorkspaceID    string           `json:"workspace_id"`
	MappingKey     string           `json:"mapping_key"`
	IdempotencyKey string           `json:"idempotency_key"`
	Target         TriggerTarget    `json:"target"`
	Principal      TriggerPrincipal `json:"principal,omitempty"`
}

type RuntimeExecutionReceipt struct {
	EventID     string `json:"event_id"`
	MappingKey  string `json:"mapping_key"`
	ExecutionID string `json:"execution_id"`
	TargetType  string `json:"target_type"`
	Status      string `json:"status"`
	ErrorCode   string `json:"error_code,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
}

// TriggerSink is implemented by Runtime. Integration owns verification,
// event durability and mapping; Runtime owns Action/Workflow execution and
// returns an idempotent receipt for that local execution.
type TriggerSink interface {
	Trigger(context.Context, TriggerRequest) (RuntimeExecutionReceipt, error)
}

type EventMappingRequirement struct {
	Key              string                             `json:"key"`
	WorkspaceID      string                             `json:"workspace_id"`
	Provider         string                             `json:"provider"`
	ConnectionKey    string                             `json:"connection_key,omitempty"`
	EventType        string                             `json:"event_type,omitempty"`
	CommandPrefix    string                             `json:"command_prefix,omitempty"`
	TargetType       string                             `json:"target_type"`
	WorkflowKey      string                             `json:"workflow_key,omitempty"`
	ObjectKey        string                             `json:"object_key,omitempty"`
	ObjectKeyPath    string                             `json:"object_key_path,omitempty"`
	RecordID         string                             `json:"record_id,omitempty"`
	RecordIDPath     string                             `json:"record_id_path,omitempty"`
	ActionKey        string                             `json:"action_key,omitempty"`
	ActionKeyPath    string                             `json:"action_key_path,omitempty"`
	ActionInput      map[string]string                  `json:"action_input,omitempty"`
	WorkflowInput    map[string]string                  `json:"workflow_input,omitempty"`
	EventFields      []EventFieldRequirement            `json:"event_fields,omitempty"`
	ExternalIdentity ExternalIdentityMappingRequirement `json:"external_identity,omitempty"`
	Payload          map[string]any                     `json:"payload,omitempty"`
	Enabled          bool                               `json:"enabled"`
}

type EventFieldRequirement struct {
	Path     string   `json:"path"`
	Type     string   `json:"type"`
	Options  []string `json:"options,omitempty"`
	Required bool     `json:"required,omitempty"`
}

type ExternalIdentityMappingRequirement struct {
	Provider    string `json:"provider,omitempty"`
	SubjectPath string `json:"subject_path,omitempty"`
	SubjectType string `json:"subject_type,omitempty"`
	NamePath    string `json:"name_path,omitempty"`
	OnUnmapped  string `json:"on_unmapped,omitempty"`
}

func (r EventMappingRequirement) Validate() error {
	for name, value := range map[string]string{"key": r.Key, "workspace_id": r.WorkspaceID, "provider": r.Provider, "target_type": r.TargetType} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("Integration event mapping %s is required", name)
		}
	}
	switch strings.TrimSpace(r.TargetType) {
	case "action":
		if strings.TrimSpace(r.ObjectKeyPath) != "" || strings.TrimSpace(r.ActionKeyPath) != "" {
			return fmt.Errorf("Integration action event mapping requires finite static object_key and action_key targets")
		}
		if strings.TrimSpace(r.ActionKey) == "" || strings.TrimSpace(r.ObjectKey) == "" {
			return fmt.Errorf("Integration action event mapping requires object_key and action_key")
		}
		if strings.TrimSpace(r.RecordID) == "" && strings.TrimSpace(r.RecordIDPath) == "" {
			return fmt.Errorf("Integration action event mapping requires record_id or record_id_path")
		}
	case "workflow":
		if strings.TrimSpace(r.WorkflowKey) == "" {
			return fmt.Errorf("Integration workflow event mapping requires workflow_key")
		}
	default:
		return fmt.Errorf("Integration event mapping target_type %q is unsupported", r.TargetType)
	}
	declared := map[string]EventFieldRequirement{}
	for index, field := range r.EventFields {
		path := strings.TrimSpace(field.Path)
		if !validEventPath(path) {
			return fmt.Errorf("Integration event mapping event_fields[%d].path is invalid", index)
		}
		if _, exists := declared[path]; exists {
			return fmt.Errorf("Integration event mapping event field path %q is duplicated", path)
		}
		if !supportedEventFieldType(field.Type) {
			return fmt.Errorf("Integration event mapping event field %q has unsupported type %q", path, field.Type)
		}
		declared[path] = field
	}
	requireDeclared := func(name, path string) error {
		path = strings.TrimSpace(path)
		if path == "" {
			return nil
		}
		if _, exists := declared[path]; !exists {
			return fmt.Errorf("Integration event mapping %s path %q is not declared in event_fields", name, path)
		}
		return nil
	}
	for key, path := range r.ActionInput {
		if err := requireDeclared("action_input."+key, path); err != nil {
			return err
		}
	}
	for key, path := range r.WorkflowInput {
		if err := requireDeclared("workflow_input."+key, path); err != nil {
			return err
		}
	}
	for name, path := range map[string]string{
		"record_id_path":                 r.RecordIDPath,
		"external_identity.subject_path": r.ExternalIdentity.SubjectPath,
		"external_identity.name_path":    r.ExternalIdentity.NamePath,
	} {
		if err := requireDeclared(name, path); err != nil {
			return err
		}
	}
	if strings.TrimSpace(r.CommandPrefix) != "" {
		if err := requireDeclared("command_prefix", "command"); err != nil {
			return err
		}
	}
	return nil
}

func validEventPath(path string) bool {
	if path == "" {
		return false
	}
	for _, segment := range strings.Split(path, ".") {
		if strings.TrimSpace(segment) == "" || strings.TrimSpace(segment) != segment {
			return false
		}
	}
	return true
}

func supportedEventFieldType(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "text", "string", "long_text", "email", "url", "date", "datetime", "relation", "user", "file",
		"number", "decimal", "currency", "percent", "integer", "boolean", "bool", "object", "map", "array", "list", "json":
		return true
	default:
		return false
	}
}

type Operations interface {
	Call(context.Context, ProviderCallRequest) (ProviderCallResult, error)
	ListInvocations(context.Context, InvocationQuery) ([]Invocation, error)
	GetInvocation(context.Context, string, string) (Invocation, error)
	AcceptWebhook(context.Context, WebhookRequest) (WebhookReceipt, error)
	ListEvents(context.Context, EventQuery) ([]Event, error)
	GetEvent(context.Context, string, string) (Event, error)
	ReplayEvent(context.Context, string, string) (Event, error)
}

type OperationsBinding interface {
	Operations() Operations
}

type ProviderResourceHealth struct {
	ObservationID     string `json:"observation_id"`
	Kind              string `json:"kind"`
	State             string `json:"state"`
	PreviousState     string `json:"previous_state,omitempty"`
	EvidenceSource    string `json:"evidence_source"`
	QuotaUsedPercent  *int   `json:"quota_used_percent,omitempty"`
	BalanceBand       string `json:"balance_band,omitempty"`
	CapabilityBlocked bool   `json:"capability_blocked"`
	ObservedAt        string `json:"observed_at"`
	ErrorCode         string `json:"error_code,omitempty"`
}
