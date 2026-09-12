package integrationsdk

import "context"

// OAuthApplication is Integration-owned application configuration. Client secrets
// are write-only and are stored with the host cipher, never returned to products.
type OAuthApplication struct {
	Key              string         `json:"key"`
	WorkspaceID      string         `json:"workspace_id"`
	ConnectorKey     string         `json:"connector_key"`
	ProviderKey      string         `json:"provider_key"`
	Name             string         `json:"name"`
	ClientID         string         `json:"client_id"`
	RedirectURI      string         `json:"redirect_uri"`
	Scopes           []string       `json:"scopes"`
	ConnectionConfig map[string]any `json:"connection_config,omitempty"`
	Configured       bool           `json:"configured"`
	Enabled          bool           `json:"enabled"`
	UpdatedAt        string         `json:"updated_at"`
}
type OAuthApplicationInput struct {
	ConnectorKey      string         `json:"connector_key"`
	ProviderKey       string         `json:"provider_key"`
	Name              string         `json:"name"`
	ClientID          string         `json:"client_id"`
	ClientSecret      string         `json:"client_secret,omitempty"`
	RedirectURI       string         `json:"redirect_uri"`
	Scopes            []string       `json:"scopes"`
	ConnectionConfig  map[string]any `json:"connection_config,omitempty"`
	Enabled           bool           `json:"enabled"`
	ExpectedUpdatedAt string         `json:"expected_updated_at,omitempty"`
}
type OAuthApplications interface {
	ListOAuthApplications(context.Context, string) ([]OAuthApplication, error)
	UpsertOAuthApplication(context.Context, string, string, string, OAuthApplicationInput) (OAuthApplication, error)
}
type OAuthApplicationsBinding interface{ OAuthApplications() OAuthApplications }

// OAuthAuthorizationOption is the safe current-user catalog. Products get only
// explicit scopes and labels; all protocol/client configuration stays in Integration.
type OAuthAuthorizationOption struct {
	Key          string   `json:"key"`
	ConnectorKey string   `json:"connector_key"`
	ProviderKey  string   `json:"provider_key"`
	Name         string   `json:"name"`
	Scopes       []string `json:"scopes"`
}
type OAuthAuthorizationInput struct {
	ApplicationKey string                 `json:"application_key"`
	Name           string                 `json:"name"`
	Scope          ConnectionAccountScope `json:"scope"`
	Scopes         []string               `json:"scopes"`
}
type OAuthAuthorizationSession struct {
	ID              string                 `json:"id"`
	Status          string                 `json:"status"`
	ApplicationKey  string                 `json:"application_key"`
	Scope           ConnectionAccountScope `json:"scope"`
	RequestedScopes []string               `json:"requested_scopes"`
	GrantedScopes   []string               `json:"granted_scopes,omitempty"`
	Account         *ConnectionAccount     `json:"account,omitempty"`
	ExpiresAt       string                 `json:"expires_at"`
	// This one-time navigation target contains state. Hosts must use no-store and
	// must never put it in conversation messages, tools, telemetry or shared logs.
	AuthorizationURL string `json:"authorization_url,omitempty"`
}
type OAuthAuthorizationCallback struct {
	State string `json:"state"`
	Code  string `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}
type OAuthAuthorizations interface {
	ListOAuthAuthorizationOptions(context.Context, ConnectionAccountSubject) ([]OAuthAuthorizationOption, error)
	StartOAuthAuthorization(context.Context, ConnectionAccountSubject, OAuthAuthorizationInput) (OAuthAuthorizationSession, error)
	GetOAuthAuthorization(context.Context, ConnectionAccountSubject, string) (OAuthAuthorizationSession, error)
	CompleteOAuthAuthorization(context.Context, ConnectionAccountSubject, OAuthAuthorizationCallback) (OAuthAuthorizationSession, error)
}
type OAuthAuthorizationsBinding interface{ OAuthAuthorizations() OAuthAuthorizations }
