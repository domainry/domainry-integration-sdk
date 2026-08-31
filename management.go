package integrationsdk

import (
	"context"
	"encoding/json"
)

// Connection is Integration-owned connector configuration. Secret material is
// referenced by stable keys and is never projected through this contract.
type Connection struct {
	Key          string            `json:"key"`
	WorkspaceID  string            `json:"workspace_id,omitempty"`
	ConnectorKey string            `json:"connector_key"`
	ProviderKey  string            `json:"provider_key"`
	Name         string            `json:"name,omitempty"`
	Status       string            `json:"status"`
	Config       map[string]any    `json:"config,omitempty"`
	SecretRefs   map[string]string `json:"secret_refs,omitempty"`
	CreatedBy    string            `json:"created_by,omitempty"`
	CreatedAt    string            `json:"created_at,omitempty"`
	UpdatedAt    string            `json:"updated_at,omitempty"`
}

type ConnectionInput struct {
	ConnectorKey string            `json:"connector_key"`
	ProviderKey  string            `json:"provider_key"`
	Name         string            `json:"name,omitempty"`
	Status       string            `json:"status,omitempty"`
	Config       map[string]any    `json:"config,omitempty"`
	SecretRefs   map[string]string `json:"secret_refs,omitempty"`
}

type Secret struct {
	Key            string `json:"key"`
	WorkspaceID    string `json:"workspace_id,omitempty"`
	Kind           string `json:"kind"`
	Status         string `json:"status"`
	Configured     bool   `json:"configured"`
	Description    string `json:"description,omitempty"`
	ValueRef       string `json:"value_ref,omitempty"`
	Fingerprint    string `json:"fingerprint,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	DisabledAt     string `json:"disabled_at,omitempty"`
	ExpiresAt      string `json:"expires_at,omitempty"`
	RotatedAt      string `json:"rotated_at,omitempty"`
	RevokedAt      string `json:"revoked_at,omitempty"`
	LastTestedAt   string `json:"last_tested_at,omitempty"`
	LastTestStatus string `json:"last_test_status,omitempty"`
	LastTestError  string `json:"last_test_error,omitempty"`
}

type SecretInput struct {
	Kind        string `json:"kind"`
	Description string `json:"description,omitempty"`
	Value       string `json:"value,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"`
}

type APIKey struct {
	Key         string   `json:"key"`
	WorkspaceID string   `json:"workspace_id,omitempty"`
	Name        string   `json:"name,omitempty"`
	TokenPrefix string   `json:"token_prefix,omitempty"`
	ActorID     string   `json:"actor_id"`
	RoleKey     string   `json:"role_key"`
	Scopes      []string `json:"scopes,omitempty"`
	Status      string   `json:"status"`
	ExpiresAt   string   `json:"expires_at,omitempty"`
	LastUsedAt  string   `json:"last_used_at,omitempty"`
	CreatedBy   string   `json:"created_by,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
	DisabledAt  string   `json:"disabled_at,omitempty"`
}

type APIKeyInput struct {
	Key       string   `json:"key,omitempty"`
	Name      string   `json:"name,omitempty"`
	ActorID   string   `json:"actor_id"`
	RoleKey   string   `json:"role_key"`
	Scopes    []string `json:"scopes,omitempty"`
	ExpiresAt string   `json:"expires_at,omitempty"`
}

type APIKeyCredential struct {
	APIKey APIKey `json:"api_key"`
	Token  string `json:"token"`
}

type ExternalIdentity struct {
	Key                  string `json:"key"`
	WorkspaceID          string `json:"workspace_id,omitempty"`
	Provider             string `json:"provider"`
	ExternalSubject      string `json:"external_subject"`
	ExternalSubjectType  string `json:"external_subject_type,omitempty"`
	ExternalName         string `json:"external_name,omitempty"`
	ExternalOrganization string `json:"external_organization,omitempty"`
	ExternalDepartment   string `json:"external_department,omitempty"`
	ExternalGroup        string `json:"external_group,omitempty"`
	ExternalBotID        string `json:"external_bot_id,omitempty"`
	ActorID              string `json:"actor_id"`
	RoleKey              string `json:"role_key"`
	Status               string `json:"status"`
	LastResolvedAt       string `json:"last_resolved_at,omitempty"`
	CreatedBy            string `json:"created_by,omitempty"`
	CreatedAt            string `json:"created_at,omitempty"`
	UpdatedAt            string `json:"updated_at,omitempty"`
	DisabledAt           string `json:"disabled_at,omitempty"`
}

type ExternalIdentityInput struct {
	Provider             string `json:"provider"`
	ExternalSubject      string `json:"external_subject"`
	ExternalSubjectType  string `json:"external_subject_type,omitempty"`
	ExternalName         string `json:"external_name,omitempty"`
	ExternalOrganization string `json:"external_organization,omitempty"`
	ExternalDepartment   string `json:"external_department,omitempty"`
	ExternalGroup        string `json:"external_group,omitempty"`
	ExternalBotID        string `json:"external_bot_id,omitempty"`
	ActorID              string `json:"actor_id"`
	RoleKey              string `json:"role_key"`
	Status               string `json:"status,omitempty"`
}

type WebhookSubscription struct {
	Key           string   `json:"key"`
	WorkspaceID   string   `json:"workspace_id,omitempty"`
	Name          string   `json:"name,omitempty"`
	ConnectorKey  string   `json:"connector_key"`
	ConnectionKey string   `json:"connection_key"`
	EventTypes    []string `json:"event_types,omitempty"`
	Status        string   `json:"status"`
	Description   string   `json:"description,omitempty"`
	CreatedBy     string   `json:"created_by,omitempty"`
	CreatedAt     string   `json:"created_at,omitempty"`
	UpdatedAt     string   `json:"updated_at,omitempty"`
	DisabledAt    string   `json:"disabled_at,omitempty"`
}

type WebhookSubscriptionInput struct {
	Name          string   `json:"name,omitempty"`
	ConnectorKey  string   `json:"connector_key"`
	ConnectionKey string   `json:"connection_key"`
	EventTypes    []string `json:"event_types,omitempty"`
	Status        string   `json:"status,omitempty"`
	Description   string   `json:"description,omitempty"`
}

type ConnectionTestRequest struct {
	Operation string          `json:"operation,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	// Input is the product HTTP compatibility name. Payload remains the
	// deployment-neutral service name; owners accept either during migration.
	Input   json.RawMessage `json:"input,omitempty"`
	Confirm bool            `json:"confirm,omitempty"`
}

type ConnectionTestResult struct {
	Connection Connection      `json:"connection"`
	Operation  string          `json:"operation"`
	Response   json.RawMessage `json:"response,omitempty"`
	Receipt    DeliveryReceipt `json:"receipt"`
}

// Management owns Integration configuration and security lifecycle. Actor IDs
// are evidence supplied by the authenticated product surface, not authority to
// impersonate another Runtime principal.
type Management interface {
	ListConnections(context.Context, string) ([]Connection, error)
	GetConnection(context.Context, string, string) (Connection, error)
	UpsertConnection(context.Context, string, string, string, ConnectionInput) (Connection, error)
	DeleteConnection(context.Context, string, string) error
	SetConnectionStatus(context.Context, string, string, string, string) (Connection, error)
	TestConnection(context.Context, string, string, ConnectionTestRequest) (ConnectionTestResult, error)

	ListSecrets(context.Context, string) ([]Secret, error)
	UpsertSecret(context.Context, string, string, string, SecretInput) (Secret, error)
	TransitionSecret(context.Context, string, string, string, string) (Secret, error)

	ListAPIKeys(context.Context, string) ([]APIKey, error)
	CreateAPIKey(context.Context, string, string, APIKeyInput) (APIKeyCredential, error)
	DisableAPIKey(context.Context, string, string, string) (APIKey, error)
	RotateAPIKey(context.Context, string, string, string) (APIKeyCredential, error)

	ListExternalIdentities(context.Context, string) ([]ExternalIdentity, error)
	UpsertExternalIdentity(context.Context, string, string, string, ExternalIdentityInput) (ExternalIdentity, error)
	DisableExternalIdentity(context.Context, string, string, string) (ExternalIdentity, error)
	ResolveExternalIdentity(context.Context, string, string, string) (ExternalIdentity, error)

	ListWebhookSubscriptions(context.Context, string) ([]WebhookSubscription, error)
	UpsertWebhookSubscription(context.Context, string, string, string, WebhookSubscriptionInput) (WebhookSubscription, error)
	DeleteWebhookSubscription(context.Context, string, string) error
	DisableWebhookSubscription(context.Context, string, string, string) (WebhookSubscription, error)
}

type ManagementBinding interface {
	Management() Management
}
