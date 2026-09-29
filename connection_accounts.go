package integrationsdk

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ConnectionAccountScope describes who may use one Integration-owned external
// account. Personal accounts belong to exactly one user. Workspace accounts are
// shared with authorized users in the same workspace.
type ConnectionAccountScope string

const (
	ConnectionAccountScopePersonal  ConnectionAccountScope = "personal"
	ConnectionAccountScopeWorkspace ConnectionAccountScope = "workspace"
)

func (s ConnectionAccountScope) Valid() bool {
	return s == ConnectionAccountScopePersonal || s == ConnectionAccountScopeWorkspace
}

// ConnectionAccountSubject is supplied by a trusted product adapter after it
// authenticates the current user. It is never accepted from model output.
type ConnectionAccountSubject struct {
	WorkspaceID string `json:"workspace_id"`
	UserID      string `json:"user_id"`
	// Access is compiled for the current action by the trusted product host.
	// Its zero value denies access. Browser and model input must never set it.
	Access ConnectionAccountAccess `json:"access"`
}

// ConnectionAccountAccess is a closed projection of current authorization, not
// a grant stored on the connection. Hosts recompute it for every action; SaaS
// carries it only over the authenticated service-to-service boundary.
type ConnectionAccountAccess struct {
	Personal  bool `json:"personal"`
	Workspace bool `json:"workspace"`
}

func (s ConnectionAccountSubject) Validate() error {
	if strings.TrimSpace(s.WorkspaceID) == "" || strings.TrimSpace(s.UserID) == "" {
		return fmt.Errorf("Integration connection account subject is incomplete")
	}
	return nil
}

// ConnectionAccount is the safe user-facing projection of a Connection. It
// deliberately excludes provider configuration, secret references, credential
// fingerprints and creator evidence.
type ConnectionAccount struct {
	Key          string                            `json:"key"`
	WorkspaceID  string                            `json:"workspace_id,omitempty"`
	ConnectorKey string                            `json:"connector_key"`
	ProviderKey  string                            `json:"provider_key"`
	Name         string                            `json:"name,omitempty"`
	Scope        ConnectionAccountScope            `json:"scope"`
	OwnerUserID  string                            `json:"owner_user_id,omitempty"`
	Status       string                            `json:"status"`
	HealthState  ConnectionAccountHealthState      `json:"health_state,omitempty"`
	Readiness    *ConnectionAccountReadiness       `json:"readiness,omitempty"`
	Background   []ConnectionAccountBackgroundTask `json:"background,omitempty"`
	CreatedAt    string                            `json:"created_at,omitempty"`
	UpdatedAt    string                            `json:"updated_at,omitempty"`
}

// ConnectionAccountBackgroundTask is an owner-supplied, secret-free snapshot of
// a provider state run. "ready" means scheduled, not that a sync succeeded.
type ConnectionAccountBackgroundTask struct {
	TaskKey       string                       `json:"task_key"`
	Status        string                       `json:"status"`
	HealthState   ConnectionAccountHealthState `json:"health_state,omitempty"`
	AttemptCount  int                          `json:"attempt_count"`
	DueAt         string                       `json:"due_at,omitempty"`
	LastSuccessAt string                       `json:"last_success_at,omitempty"`
	LastErrorCode string                       `json:"last_error_code,omitempty"`
	UpdatedAt     string                       `json:"updated_at,omitempty"`
}

// Health is a read-time interpretation of provider-run state, not a second
// persisted queue or a claim that a scheduled run has completed successfully.
type ConnectionAccountHealthState string

const (
	ConnectionAccountHealthHealthy       ConnectionAccountHealthState = "healthy"
	ConnectionAccountHealthInitializing  ConnectionAccountHealthState = "initializing"
	ConnectionAccountHealthDelayed       ConnectionAccountHealthState = "delayed"
	ConnectionAccountHealthFailing       ConnectionAccountHealthState = "failing"
	ConnectionAccountHealthInactive      ConnectionAccountHealthState = "inactive"
	ConnectionAccountHealthNotApplicable ConnectionAccountHealthState = "not_applicable"
	ConnectionAccountHealthUnknown       ConnectionAccountHealthState = "unknown"
)

// Retry is an optimistic, current-user request for one failed provider state
// run. It schedules the existing run; it never creates a product-local job.
type ConnectionAccountBackgroundRetryRequest struct {
	ExpectedUpdatedAt string `json:"expected_updated_at"`
}

func (r ConnectionAccountBackgroundRetryRequest) Validate() error {
	if _, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(r.ExpectedUpdatedAt)); err != nil {
		return fmt.Errorf("Integration background task revision is required")
	}
	return nil
}

// ConnectionAccountReadiness reports current local configuration only. It does
// not probe the upstream, refresh credentials or guarantee the next operation.
// Missing readiness from an older owner must be treated as unknown/unavailable.
// GrantedScopes come only from an actual OAuth grant, never requested scopes.
type ConnectionAccountReadiness struct {
	Available           bool                              `json:"available"`
	State               string                            `json:"state"`
	GrantedScopes       []string                          `json:"granted_scopes,omitempty"`
	CredentialExpiresAt string                            `json:"credential_expires_at,omitempty"`
	Test                *ConnectionAccountTestEligibility `json:"test,omitempty"`
}

// Probe eligibility is independent from business-tool eligibility. Missing
// probe scope does not invalidate otherwise usable calendar/mail credentials.
type ConnectionAccountTestEligibility struct {
	Allowed           bool       `json:"allowed"`
	State             string     `json:"state"`
	ScopeAlternatives [][]string `json:"scope_alternatives,omitempty"`
}

type ConnectionAccountRegistration struct {
	Scope       ConnectionAccountScope `json:"scope"`
	OwnerUserID string                 `json:"owner_user_id,omitempty"`
}

type ConnectionAccountTestResult struct {
	Account   ConnectionAccount `json:"account"`
	Operation string            `json:"operation"`
	Connected bool              `json:"connected"`
	Receipt   DeliveryReceipt   `json:"receipt"`
}

func (r ConnectionAccountRegistration) Validate() error {
	if !r.Scope.Valid() {
		return fmt.Errorf("Integration connection account scope is invalid")
	}
	owner := strings.TrimSpace(r.OwnerUserID)
	if r.Scope == ConnectionAccountScopePersonal && owner == "" {
		return fmt.Errorf("Integration personal connection owner is required")
	}
	if r.Scope == ConnectionAccountScopeWorkspace && owner != "" {
		return fmt.Errorf("Integration workspace connection must not have a personal owner")
	}
	return nil
}

// ConnectionAccounts is the current-user boundary. Administrative connection
// configuration remains on Management. Implementations must include only
// authorized workspace accounts and personal accounts owned by Subject.UserID.
// Test requests may only select the fixed connection test, never a business
// operation. Provider response bodies are not part of this safe account API.
type ConnectionAccounts interface {
	ListConnectionAccounts(context.Context, ConnectionAccountSubject) ([]ConnectionAccount, error)
	GetConnectionAccount(context.Context, ConnectionAccountSubject, string) (ConnectionAccount, error)
	TestConnectionAccount(context.Context, ConnectionAccountSubject, string, ConnectionTestRequest) (ConnectionAccountTestResult, error)
	RetryConnectionAccountBackgroundTask(context.Context, ConnectionAccountSubject, string, string, ConnectionAccountBackgroundRetryRequest) (ConnectionAccountBackgroundTask, error)
	RevokeConnectionAccount(context.Context, ConnectionAccountSubject, string, string) (ConnectionAccount, error)
}

// ConnectionAccountAdministration publishes an administratively configured
// Connection to the current-user account catalog. It is intentionally separate
// from ConnectionAccounts so ordinary users cannot change account ownership.
type ConnectionAccountAdministration interface {
	RegisterConnectionAccount(context.Context, string, string, string, ConnectionAccountRegistration) (ConnectionAccount, error)
}

type ConnectionAccountsBinding interface {
	ConnectionAccounts() ConnectionAccounts
}

type ConnectionAccountAdministrationBinding interface {
	ConnectionAccountAdministration() ConnectionAccountAdministration
}
