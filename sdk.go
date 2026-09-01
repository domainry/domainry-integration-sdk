// Package integrationsdk defines the deployment-neutral Integration owner boundary.
package integrationsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/domainry/domainry-foundation/modulecapability"
)

const ProtocolVersionV1 = "domainry-integration-protocol-v1"

type DeploymentMode string

const (
	DeploymentModeModule DeploymentMode = "module"
	DeploymentModeSaaS   DeploymentMode = "saas"
)

type ApplicationRef struct {
	RuntimeID string `json:"runtime_id"`
}

func (r ApplicationRef) Validate() error {
	if strings.TrimSpace(r.RuntimeID) == "" {
		return fmt.Errorf("Integration runtime identity is required")
	}
	return nil
}

type Descriptor struct {
	ProtocolVersion string         `json:"protocol_version"`
	Mode            DeploymentMode `json:"mode"`
	Audience        string         `json:"audience,omitempty"`
	Capabilities    []string       `json:"capabilities"`
}

func (d Descriptor) Validate() error {
	if d.ProtocolVersion != ProtocolVersionV1 {
		return fmt.Errorf("unsupported Integration protocol %q", d.ProtocolVersion)
	}
	if d.Mode != DeploymentModeModule && d.Mode != DeploymentModeSaaS {
		return fmt.Errorf("invalid Integration deployment mode %q", d.Mode)
	}
	required := map[string]bool{
		"catalog.read": false, "requirements.connections.sync": false,
		"delivery.accept": false, "delivery.query": false,
		"web_push_subscriptions.manage": false,
		"management.connections":        false, "management.secrets": false,
		"management.api_keys": false, "management.external_identities": false,
		"management.webhook_subscriptions": false,
		"operations.call":                  false, "operations.invocations.query": false,
		"inbound.webhooks.accept": false, "inbound.events.query": false,
	}
	for _, capability := range d.Capabilities {
		if _, ok := required[strings.TrimSpace(capability)]; ok {
			required[strings.TrimSpace(capability)] = true
		}
	}
	for capability, present := range required {
		if !present {
			return fmt.Errorf("Integration capability %q is required", capability)
		}
	}
	return nil
}

// Catalog and Delivery are deliberately deployment-neutral. Runtime retains
// only its local outbox handoff and calls these owner ports after commit.
type ConnectorDefinition struct {
	Key         string          `json:"key"`
	DisplayName string          `json:"display_name"`
	Operations  []string        `json:"operations"`
	Schema      json.RawMessage `json:"schema,omitempty"`
	// Definition preserves the complete source-owned connector contract for
	// consumers that need provider/configuration metadata beyond the summary.
	Definition json.RawMessage `json:"definition"`
}

type Catalog interface {
	ListConnectorDefinitions(context.Context) ([]ConnectorDefinition, error)
}

// ConnectionRequirement is an application-owned declaration consumed by the
// Integration owner. Runtime does not persist or manage the resulting
// connection; Module and SaaS implementations materialize it in owner state.
type ConnectionRequirement struct {
	Key          string          `json:"key"`
	WorkspaceID  string          `json:"workspace_id"`
	ConnectorKey string          `json:"connector_key"`
	ProviderKey  string          `json:"provider_key"`
	Name         string          `json:"name,omitempty"`
	Status       string          `json:"status,omitempty"`
	Config       json.RawMessage `json:"config"`
}

func (r ConnectionRequirement) Validate() error {
	for name, value := range map[string]string{"key": r.Key, "workspace_id": r.WorkspaceID, "connector_key": r.ConnectorKey, "provider_key": r.ProviderKey} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("Integration connection requirement %s is required", name)
		}
	}
	if !json.Valid(r.Config) {
		return fmt.Errorf("Integration connection requirement config must be valid JSON")
	}
	return nil
}

type Requirements interface {
	SynchronizeConnections(context.Context, []ConnectionRequirement) error
	SynchronizeEventMappings(context.Context, []EventMappingRequirement) error
}

// WebPushSubscription is an Integration-owned browser delivery target. Secret
// endpoint material is accepted on write and deliberately omitted on reads.
type WebPushSubscription struct {
	ID           string `json:"id"`
	WorkspaceID  string `json:"workspace_id,omitempty"`
	UserID       string `json:"user_id"`
	EndpointHash string `json:"endpoint_hash"`
	Status       string `json:"status"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	RevokedAt    string `json:"revoked_at,omitempty"`
}

type WebPushSubscriptionInput struct {
	Endpoint  string `json:"endpoint"`
	P256DH    string `json:"p256dh"`
	Auth      string `json:"auth"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

type WebPushReadiness struct {
	Ready         bool   `json:"ready"`
	PublicKey     string `json:"public_key,omitempty"`
	ConnectionKey string `json:"connection_key,omitempty"`
	Status        string `json:"status"`
	Reason        string `json:"reason,omitempty"`
}

type WebPushSubscriptions interface {
	Readiness(context.Context, string) (WebPushReadiness, error)
	List(context.Context, string, string) ([]WebPushSubscription, error)
	Upsert(context.Context, string, string, string, WebPushSubscriptionInput) (WebPushSubscription, error)
	Revoke(context.Context, string, string, string) (WebPushSubscription, error)
	CleanupExpired(context.Context, string) (int, error)
}

type DeliveryRequest struct {
	MessageID        string          `json:"message_id"`
	DeduplicationKey string          `json:"deduplication_key"`
	WorkspaceID      string          `json:"workspace_id"`
	ConnectorKey     string          `json:"connector_key"`
	ConnectionKey    string          `json:"connection_key,omitempty"`
	Operation        string          `json:"operation"`
	Payload          json.RawMessage `json:"payload"`
}

func (r DeliveryRequest) Validate() error {
	for name, value := range map[string]string{"message_id": r.MessageID, "deduplication_key": r.DeduplicationKey, "workspace_id": r.WorkspaceID, "connector_key": r.ConnectorKey, "operation": r.Operation} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("Integration delivery %s is required", name)
		}
	}
	if !json.Valid(r.Payload) {
		return fmt.Errorf("Integration delivery payload must be valid JSON")
	}
	return nil
}

type DeliveryStatus string

const (
	DeliveryStatusAccepted  DeliveryStatus = "accepted"
	DeliveryStatusRunning   DeliveryStatus = "running"
	DeliveryStatusSucceeded DeliveryStatus = "succeeded"
	DeliveryStatusFailed    DeliveryStatus = "failed"
)

type DeliveryReceipt struct {
	MessageID    string         `json:"message_id"`
	InvocationID string         `json:"invocation_id"`
	Status       DeliveryStatus `json:"status"`
	ResultRef    string         `json:"result_ref,omitempty"`
	ErrorCode    string         `json:"error_code,omitempty"`
}

type Delivery interface {
	Accept(context.Context, DeliveryRequest) (DeliveryReceipt, error)
	Query(context.Context, string) (DeliveryReceipt, error)
}

type Binding interface {
	modulecapability.Binding
	Descriptor() Descriptor
	Catalog() Catalog
	Requirements() Requirements
	Delivery() Delivery
	Close(context.Context) error
}

// LocalWorkers is exposed only by an embedded Integration Module. A SaaS
// deployment owns these loops in its own process and remote bindings do not
// expose them.
type LocalWorkers interface {
	ProcessDueEvents(context.Context, int) (int, error)
	ProcessDueProviderTasks(context.Context, int) (int, error)
	ProcessDueReconciliations(context.Context, int) (int, error)
	ProcessDueCredentialExpirations(context.Context, int) (int, error)
}

type LocalWorkerBinding interface {
	LocalWorkers() (LocalWorkers, bool)
}

type WebPushBinding interface {
	WebPushSubscriptions() WebPushSubscriptions
}

// Factory identifies the deployment topology selected by the application
// composition root. The matching modulehost or saashost contract performs the
// actual open because the two topologies intentionally receive different host
// capabilities.
type Factory interface {
	DeploymentMode() DeploymentMode
}
