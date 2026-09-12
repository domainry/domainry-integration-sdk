package remote

import (
	"context"
	integrationsdk "github.com/domainry/domainry-integration-sdk"
	"net/http"
	"net/url"
	"strings"
)

func (b *binding) OAuthApplications() integrationsdk.OAuthApplications     { return b.client }
func (b *binding) OAuthAuthorizations() integrationsdk.OAuthAuthorizations { return b.client }
func (c *remoteClient) ListOAuthApplications(ctx context.Context, workspace string) ([]integrationsdk.OAuthApplication, error) {
	var result struct {
		Items []integrationsdk.OAuthApplication `json:"items"`
	}
	err := c.call(ctx, http.MethodGet, managementPath("oauth-applications", workspace), nil, &result)
	return result.Items, err
}
func (c *remoteClient) UpsertOAuthApplication(ctx context.Context, workspace, key, actor string, input integrationsdk.OAuthApplicationInput) (integrationsdk.OAuthApplication, error) {
	var result integrationsdk.OAuthApplication
	err := c.call(ctx, http.MethodPut, managementPath("oauth-applications/"+url.PathEscape(key), workspace)+"&actor_id="+url.QueryEscape(actor), input, &result)
	return result, err
}
func oauthPath(suffix string, subject integrationsdk.ConnectionAccountSubject) string {
	return strings.Replace(connectionAccountPath(suffix, subject), "/connection-accounts", "/oauth-authorizations", 1)
}
func (c *remoteClient) ListOAuthAuthorizationOptions(ctx context.Context, subject integrationsdk.ConnectionAccountSubject) ([]integrationsdk.OAuthAuthorizationOption, error) {
	if err := subject.Validate(); err != nil {
		return nil, err
	}
	var result struct {
		Items []integrationsdk.OAuthAuthorizationOption `json:"items"`
	}
	err := c.call(ctx, http.MethodGet, oauthPath("/options", subject), nil, &result)
	return result.Items, err
}
func (c *remoteClient) StartOAuthAuthorization(ctx context.Context, subject integrationsdk.ConnectionAccountSubject, input integrationsdk.OAuthAuthorizationInput) (integrationsdk.OAuthAuthorizationSession, error) {
	if err := subject.Validate(); err != nil {
		return integrationsdk.OAuthAuthorizationSession{}, err
	}
	var result integrationsdk.OAuthAuthorizationSession
	err := c.call(ctx, http.MethodPost, oauthPath("", subject), input, &result)
	return result, err
}
func (c *remoteClient) GetOAuthAuthorization(ctx context.Context, subject integrationsdk.ConnectionAccountSubject, id string) (integrationsdk.OAuthAuthorizationSession, error) {
	if err := subject.Validate(); err != nil {
		return integrationsdk.OAuthAuthorizationSession{}, err
	}
	var result integrationsdk.OAuthAuthorizationSession
	err := c.call(ctx, http.MethodGet, oauthPath("/"+url.PathEscape(id), subject), nil, &result)
	return result, err
}
func (c *remoteClient) CompleteOAuthAuthorization(ctx context.Context, subject integrationsdk.ConnectionAccountSubject, input integrationsdk.OAuthAuthorizationCallback) (integrationsdk.OAuthAuthorizationSession, error) {
	if err := subject.Validate(); err != nil {
		return integrationsdk.OAuthAuthorizationSession{}, err
	}
	var result integrationsdk.OAuthAuthorizationSession
	err := c.call(ctx, http.MethodPost, oauthPath("/callback", subject), input, &result)
	return result, err
}

var _ integrationsdk.OAuthApplications = (*remoteClient)(nil)
var _ integrationsdk.OAuthAuthorizations = (*remoteClient)(nil)
