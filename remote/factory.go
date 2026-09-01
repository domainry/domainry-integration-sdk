package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/domainry/domainry-foundation/modulecapability"
	integrationsdk "github.com/domainry/domainry-integration-sdk"
	"github.com/domainry/domainry-integration-sdk/saashost"
)

type Options struct {
	BaseURL                  string
	Token                    string
	HTTPClient               *http.Client
	CapabilityContractSHA256 string
}

type Factory struct{ options Options }

func NewFactory(options Options) *Factory { return &Factory{options: options} }

func (*Factory) DeploymentMode() integrationsdk.DeploymentMode {
	return integrationsdk.DeploymentModeSaaS
}

func (f *Factory) OpenSaaS(ctx context.Context, application integrationsdk.ApplicationRef, _ saashost.Host) (integrationsdk.Binding, error) {
	if err := application.Validate(); err != nil {
		return nil, err
	}
	if err := modulecapability.ValidateRemoteExpectation("integration", f.options.CapabilityContractSHA256); err != nil {
		return nil, err
	}
	base, err := url.Parse(strings.TrimSpace(f.options.BaseURL))
	if err != nil || base.Scheme == "" || base.Host == "" {
		return nil, fmt.Errorf("Integration SaaS base URL is invalid")
	}
	client := f.options.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	transport := &remoteClient{base: base, token: strings.TrimSpace(f.options.Token), runtimeID: application.RuntimeID, client: client}
	capability, err := modulecapability.OpenRemote(ctx, modulecapability.RemoteConfig{
		BaseURL: strings.TrimRight(base.String(), "/"), Client: client, ExpectedModuleKey: "integration", ExpectedContractSHA256: f.options.CapabilityContractSHA256,
		Authorize: func(request *http.Request) error {
			request.Header.Set("Authorization", "Bearer "+transport.token)
			request.Header.Set("X-Domainry-Runtime-ID", transport.runtimeID)
			return nil
		},
	})
	if err != nil {
		return nil, err
	}
	return &binding{client: transport, capability: capability}, nil
}

type binding struct {
	client     *remoteClient
	capability modulecapability.Binding
}

func (b *binding) CapabilitySummary(ctx context.Context) (modulecapability.ModuleSummary, error) {
	return b.capability.CapabilitySummary(ctx)
}
func (b *binding) CapabilityCategory(ctx context.Context, key string) (modulecapability.CategoryDocument, error) {
	return b.capability.CapabilityCategory(ctx, key)
}
func (b *binding) ValidateCapabilityCandidate(ctx context.Context, request modulecapability.ValidationRequest) (modulecapability.ValidationResult, error) {
	return b.capability.ValidateCapabilityCandidate(ctx, request)
}

func (*binding) Descriptor() integrationsdk.Descriptor {
	return integrationsdk.Descriptor{ProtocolVersion: integrationsdk.ProtocolVersionV1, Mode: integrationsdk.DeploymentModeSaaS, Capabilities: []string{"catalog.read", "requirements.connections.sync", "delivery.accept", "delivery.query", "web_push_subscriptions.manage", "management.connections", "management.secrets", "management.api_keys", "management.external_identities", "management.webhook_subscriptions", "operations.call", "operations.invocations.query", "inbound.webhooks.accept", "inbound.events.query"}}
}
func (b *binding) Catalog() integrationsdk.Catalog                           { return b.client }
func (b *binding) Requirements() integrationsdk.Requirements                 { return b.client }
func (b *binding) WebPushSubscriptions() integrationsdk.WebPushSubscriptions { return b.client }
func (b *binding) Delivery() integrationsdk.Delivery                         { return b.client }
func (b *binding) Management() integrationsdk.Management                     { return b.client }
func (b *binding) Operations() integrationsdk.Operations                     { return b.client }
func (*binding) Close(context.Context) error                                 { return nil }

type remoteClient struct {
	base      *url.URL
	token     string
	runtimeID string
	client    *http.Client
}

func (c *remoteClient) ListConnectorDefinitions(ctx context.Context) ([]integrationsdk.ConnectorDefinition, error) {
	var response struct {
		Items []integrationsdk.ConnectorDefinition `json:"items"`
	}
	if err := c.call(ctx, http.MethodGet, "/v1/connector-definitions", nil, &response); err != nil {
		return nil, err
	}
	return response.Items, nil
}

func (c *remoteClient) SynchronizeConnections(ctx context.Context, requirements []integrationsdk.ConnectionRequirement) error {
	for _, requirement := range requirements {
		if err := requirement.Validate(); err != nil {
			return err
		}
	}
	return c.call(ctx, http.MethodPut, "/v1/application-requirements/connections", struct {
		Items []integrationsdk.ConnectionRequirement `json:"items"`
	}{Items: requirements}, nil)
}

func (c *remoteClient) SynchronizeEventMappings(ctx context.Context, requirements []integrationsdk.EventMappingRequirement) error {
	for _, requirement := range requirements {
		if err := requirement.Validate(); err != nil {
			return err
		}
	}
	return c.call(ctx, http.MethodPut, "/v1/application-requirements/event-mappings", struct {
		Items []integrationsdk.EventMappingRequirement `json:"items"`
	}{Items: requirements}, nil)
}

func (c *remoteClient) Accept(ctx context.Context, request integrationsdk.DeliveryRequest) (integrationsdk.DeliveryReceipt, error) {
	if err := request.Validate(); err != nil {
		return integrationsdk.DeliveryReceipt{}, err
	}
	var receipt integrationsdk.DeliveryReceipt
	if err := c.call(ctx, http.MethodPost, "/v1/deliveries", request, &receipt); err != nil {
		return integrationsdk.DeliveryReceipt{}, err
	}
	return receipt, nil
}

func (c *remoteClient) Query(ctx context.Context, messageID string) (integrationsdk.DeliveryReceipt, error) {
	if strings.TrimSpace(messageID) == "" {
		return integrationsdk.DeliveryReceipt{}, fmt.Errorf("Integration delivery message ID is required")
	}
	var receipt integrationsdk.DeliveryReceipt
	path := "/v1/deliveries/" + url.PathEscape(messageID)
	if err := c.call(ctx, http.MethodGet, path, nil, &receipt); err != nil {
		return integrationsdk.DeliveryReceipt{}, err
	}
	return receipt, nil
}

func (c *remoteClient) Readiness(ctx context.Context, workspaceID string) (integrationsdk.WebPushReadiness, error) {
	var value integrationsdk.WebPushReadiness
	path := "/v1/web-push/readiness?workspace_id=" + url.QueryEscape(strings.TrimSpace(workspaceID))
	if err := c.call(ctx, http.MethodGet, path, nil, &value); err != nil {
		return value, err
	}
	return value, nil
}

func (c *remoteClient) List(ctx context.Context, workspaceID, userID string) ([]integrationsdk.WebPushSubscription, error) {
	var response struct {
		Items []integrationsdk.WebPushSubscription `json:"items"`
	}
	path := "/v1/web-push-subscriptions?workspace_id=" + url.QueryEscape(strings.TrimSpace(workspaceID)) + "&user_id=" + url.QueryEscape(strings.TrimSpace(userID))
	if err := c.call(ctx, http.MethodGet, path, nil, &response); err != nil {
		return nil, err
	}
	return response.Items, nil
}

func (c *remoteClient) Upsert(ctx context.Context, workspaceID, userID, id string, input integrationsdk.WebPushSubscriptionInput) (integrationsdk.WebPushSubscription, error) {
	var value integrationsdk.WebPushSubscription
	path := "/v1/web-push-subscriptions/" + url.PathEscape(strings.TrimSpace(id)) + "?workspace_id=" + url.QueryEscape(strings.TrimSpace(workspaceID)) + "&user_id=" + url.QueryEscape(strings.TrimSpace(userID))
	if err := c.call(ctx, http.MethodPut, path, input, &value); err != nil {
		return value, err
	}
	return value, nil
}

func (c *remoteClient) Revoke(ctx context.Context, workspaceID, userID, id string) (integrationsdk.WebPushSubscription, error) {
	var value integrationsdk.WebPushSubscription
	path := "/v1/web-push-subscriptions/" + url.PathEscape(strings.TrimSpace(id)) + "/revoke?workspace_id=" + url.QueryEscape(strings.TrimSpace(workspaceID)) + "&user_id=" + url.QueryEscape(strings.TrimSpace(userID))
	if err := c.call(ctx, http.MethodPost, path, nil, &value); err != nil {
		return value, err
	}
	return value, nil
}

func (c *remoteClient) CleanupExpired(ctx context.Context, workspaceID string) (int, error) {
	var response struct {
		Cleaned int `json:"cleaned"`
	}
	path := "/v1/web-push-subscriptions/cleanup-expired?workspace_id=" + url.QueryEscape(strings.TrimSpace(workspaceID))
	if err := c.call(ctx, http.MethodPost, path, nil, &response); err != nil {
		return 0, err
	}
	return response.Cleaned, nil
}

func managementPath(resource, workspaceID string) string {
	return "/v1/management/" + resource + "?workspace_id=" + url.QueryEscape(strings.TrimSpace(workspaceID))
}

func (c *remoteClient) ListConnections(ctx context.Context, workspaceID string) ([]integrationsdk.Connection, error) {
	var response struct {
		Items []integrationsdk.Connection `json:"items"`
	}
	err := c.call(ctx, http.MethodGet, managementPath("connections", workspaceID), nil, &response)
	return response.Items, err
}
func (c *remoteClient) GetConnection(ctx context.Context, workspaceID, key string) (integrationsdk.Connection, error) {
	var value integrationsdk.Connection
	err := c.call(ctx, http.MethodGet, managementPath("connections/"+url.PathEscape(strings.TrimSpace(key)), workspaceID), nil, &value)
	return value, err
}
func (c *remoteClient) UpsertConnection(ctx context.Context, workspaceID, key, actorID string, input integrationsdk.ConnectionInput) (integrationsdk.Connection, error) {
	var value integrationsdk.Connection
	path := managementPath("connections/"+url.PathEscape(strings.TrimSpace(key)), workspaceID) + "&actor_id=" + url.QueryEscape(strings.TrimSpace(actorID))
	err := c.call(ctx, http.MethodPut, path, input, &value)
	return value, err
}
func (c *remoteClient) DeleteConnection(ctx context.Context, workspaceID, key string) error {
	return c.call(ctx, http.MethodDelete, managementPath("connections/"+url.PathEscape(strings.TrimSpace(key)), workspaceID), nil, nil)
}
func (c *remoteClient) SetConnectionStatus(ctx context.Context, workspaceID, key, status, actorID string) (integrationsdk.Connection, error) {
	var value integrationsdk.Connection
	path := managementPath("connections/"+url.PathEscape(strings.TrimSpace(key))+"/status", workspaceID) + "&actor_id=" + url.QueryEscape(strings.TrimSpace(actorID))
	err := c.call(ctx, http.MethodPost, path, map[string]string{"status": status}, &value)
	return value, err
}
func (c *remoteClient) TestConnection(ctx context.Context, workspaceID, key string, request integrationsdk.ConnectionTestRequest) (integrationsdk.ConnectionTestResult, error) {
	var value integrationsdk.ConnectionTestResult
	err := c.call(ctx, http.MethodPost, managementPath("connections/"+url.PathEscape(strings.TrimSpace(key))+"/test", workspaceID), request, &value)
	return value, err
}
func (c *remoteClient) ListSecrets(ctx context.Context, workspaceID string) ([]integrationsdk.Secret, error) {
	var response struct {
		Items []integrationsdk.Secret `json:"items"`
	}
	err := c.call(ctx, http.MethodGet, managementPath("secrets", workspaceID), nil, &response)
	return response.Items, err
}
func (c *remoteClient) UpsertSecret(ctx context.Context, workspaceID, key, actorID string, input integrationsdk.SecretInput) (integrationsdk.Secret, error) {
	var value integrationsdk.Secret
	path := managementPath("secrets/"+url.PathEscape(strings.TrimSpace(key)), workspaceID) + "&actor_id=" + url.QueryEscape(strings.TrimSpace(actorID))
	err := c.call(ctx, http.MethodPut, path, input, &value)
	return value, err
}
func (c *remoteClient) TransitionSecret(ctx context.Context, workspaceID, key, transition, actorID string) (integrationsdk.Secret, error) {
	var value integrationsdk.Secret
	path := managementPath("secrets/"+url.PathEscape(strings.TrimSpace(key))+"/transition", workspaceID) + "&actor_id=" + url.QueryEscape(strings.TrimSpace(actorID))
	err := c.call(ctx, http.MethodPost, path, map[string]string{"transition": transition}, &value)
	return value, err
}
func (c *remoteClient) ListAPIKeys(ctx context.Context, workspaceID string) ([]integrationsdk.APIKey, error) {
	var response struct {
		Items []integrationsdk.APIKey `json:"items"`
	}
	err := c.call(ctx, http.MethodGet, managementPath("api-keys", workspaceID), nil, &response)
	return response.Items, err
}
func (c *remoteClient) CreateAPIKey(ctx context.Context, workspaceID, actorID string, input integrationsdk.APIKeyInput) (integrationsdk.APIKeyCredential, error) {
	var value integrationsdk.APIKeyCredential
	path := managementPath("api-keys", workspaceID) + "&actor_id=" + url.QueryEscape(strings.TrimSpace(actorID))
	err := c.call(ctx, http.MethodPost, path, input, &value)
	return value, err
}
func (c *remoteClient) DisableAPIKey(ctx context.Context, workspaceID, key, actorID string) (integrationsdk.APIKey, error) {
	var value integrationsdk.APIKey
	path := managementPath("api-keys/"+url.PathEscape(strings.TrimSpace(key))+"/disable", workspaceID) + "&actor_id=" + url.QueryEscape(strings.TrimSpace(actorID))
	err := c.call(ctx, http.MethodPost, path, nil, &value)
	return value, err
}
func (c *remoteClient) RotateAPIKey(ctx context.Context, workspaceID, key, actorID string) (integrationsdk.APIKeyCredential, error) {
	var value integrationsdk.APIKeyCredential
	path := managementPath("api-keys/"+url.PathEscape(strings.TrimSpace(key))+"/rotate", workspaceID) + "&actor_id=" + url.QueryEscape(strings.TrimSpace(actorID))
	err := c.call(ctx, http.MethodPost, path, nil, &value)
	return value, err
}
func (c *remoteClient) ListExternalIdentities(ctx context.Context, workspaceID string) ([]integrationsdk.ExternalIdentity, error) {
	var response struct {
		Items []integrationsdk.ExternalIdentity `json:"items"`
	}
	err := c.call(ctx, http.MethodGet, managementPath("external-identities", workspaceID), nil, &response)
	return response.Items, err
}
func (c *remoteClient) UpsertExternalIdentity(ctx context.Context, workspaceID, key, actorID string, input integrationsdk.ExternalIdentityInput) (integrationsdk.ExternalIdentity, error) {
	var value integrationsdk.ExternalIdentity
	path := managementPath("external-identities/"+url.PathEscape(strings.TrimSpace(key)), workspaceID) + "&actor_id=" + url.QueryEscape(strings.TrimSpace(actorID))
	err := c.call(ctx, http.MethodPut, path, input, &value)
	return value, err
}
func (c *remoteClient) DisableExternalIdentity(ctx context.Context, workspaceID, key, actorID string) (integrationsdk.ExternalIdentity, error) {
	var value integrationsdk.ExternalIdentity
	path := managementPath("external-identities/"+url.PathEscape(strings.TrimSpace(key))+"/disable", workspaceID) + "&actor_id=" + url.QueryEscape(strings.TrimSpace(actorID))
	err := c.call(ctx, http.MethodPost, path, nil, &value)
	return value, err
}
func (c *remoteClient) ResolveExternalIdentity(ctx context.Context, workspaceID, provider, externalSubject string) (integrationsdk.ExternalIdentity, error) {
	var value integrationsdk.ExternalIdentity
	path := managementPath("external-identities/resolve", workspaceID) + "&provider=" + url.QueryEscape(strings.TrimSpace(provider)) + "&external_subject=" + url.QueryEscape(strings.TrimSpace(externalSubject))
	err := c.call(ctx, http.MethodPost, path, nil, &value)
	return value, err
}
func (c *remoteClient) ListWebhookSubscriptions(ctx context.Context, workspaceID string) ([]integrationsdk.WebhookSubscription, error) {
	var response struct {
		Items []integrationsdk.WebhookSubscription `json:"items"`
	}
	err := c.call(ctx, http.MethodGet, managementPath("webhook-subscriptions", workspaceID), nil, &response)
	return response.Items, err
}
func (c *remoteClient) UpsertWebhookSubscription(ctx context.Context, workspaceID, key, actorID string, input integrationsdk.WebhookSubscriptionInput) (integrationsdk.WebhookSubscription, error) {
	var value integrationsdk.WebhookSubscription
	path := managementPath("webhook-subscriptions/"+url.PathEscape(strings.TrimSpace(key)), workspaceID) + "&actor_id=" + url.QueryEscape(strings.TrimSpace(actorID))
	err := c.call(ctx, http.MethodPut, path, input, &value)
	return value, err
}
func (c *remoteClient) DeleteWebhookSubscription(ctx context.Context, workspaceID, key string) error {
	return c.call(ctx, http.MethodDelete, managementPath("webhook-subscriptions/"+url.PathEscape(strings.TrimSpace(key)), workspaceID), nil, nil)
}
func (c *remoteClient) DisableWebhookSubscription(ctx context.Context, workspaceID, key, actorID string) (integrationsdk.WebhookSubscription, error) {
	var value integrationsdk.WebhookSubscription
	path := managementPath("webhook-subscriptions/"+url.PathEscape(strings.TrimSpace(key))+"/disable", workspaceID) + "&actor_id=" + url.QueryEscape(strings.TrimSpace(actorID))
	err := c.call(ctx, http.MethodPost, path, nil, &value)
	return value, err
}

func (c *remoteClient) Call(ctx context.Context, request integrationsdk.ProviderCallRequest) (integrationsdk.ProviderCallResult, error) {
	if err := request.Validate(); err != nil {
		return integrationsdk.ProviderCallResult{}, err
	}
	var result integrationsdk.ProviderCallResult
	err := c.call(ctx, http.MethodPost, "/v1/operations/call", request, &result)
	return result, err
}

func (c *remoteClient) ListInvocations(ctx context.Context, query integrationsdk.InvocationQuery) ([]integrationsdk.Invocation, error) {
	var response struct {
		Items []integrationsdk.Invocation `json:"items"`
	}
	values := url.Values{}
	values.Set("workspace_id", strings.TrimSpace(query.WorkspaceID))
	values.Set("connector_key", strings.TrimSpace(query.ConnectorKey))
	values.Set("connection_key", strings.TrimSpace(query.ConnectionKey))
	values.Set("operation", strings.TrimSpace(query.Operation))
	values.Set("status", strings.TrimSpace(query.Status))
	values.Set("created_from", strings.TrimSpace(query.CreatedFrom))
	if query.Limit > 0 {
		values.Set("limit", fmt.Sprint(query.Limit))
	}
	err := c.call(ctx, http.MethodGet, "/v1/operations/invocations?"+values.Encode(), nil, &response)
	return response.Items, err
}

func (c *remoteClient) GetInvocation(ctx context.Context, workspaceID, id string) (integrationsdk.Invocation, error) {
	var value integrationsdk.Invocation
	path := "/v1/operations/invocations/" + url.PathEscape(strings.TrimSpace(id)) + "?workspace_id=" + url.QueryEscape(strings.TrimSpace(workspaceID))
	err := c.call(ctx, http.MethodGet, path, nil, &value)
	return value, err
}

func (c *remoteClient) AcceptWebhook(ctx context.Context, request integrationsdk.WebhookRequest) (integrationsdk.WebhookReceipt, error) {
	if err := request.Validate(); err != nil {
		return integrationsdk.WebhookReceipt{}, err
	}
	var value integrationsdk.WebhookReceipt
	err := c.call(ctx, http.MethodPost, "/v1/inbound/webhooks", request, &value)
	return value, err
}

func (c *remoteClient) ListEvents(ctx context.Context, query integrationsdk.EventQuery) ([]integrationsdk.Event, error) {
	var response struct {
		Items []integrationsdk.Event `json:"items"`
	}
	values := url.Values{}
	values.Set("workspace_id", strings.TrimSpace(query.WorkspaceID))
	values.Set("provider", strings.TrimSpace(query.Provider))
	values.Set("event_type", strings.TrimSpace(query.EventType))
	values.Set("status", strings.TrimSpace(query.Status))
	if query.Limit > 0 {
		values.Set("limit", fmt.Sprint(query.Limit))
	}
	err := c.call(ctx, http.MethodGet, "/v1/inbound/events?"+values.Encode(), nil, &response)
	return response.Items, err
}

func (c *remoteClient) GetEvent(ctx context.Context, workspaceID, id string) (integrationsdk.Event, error) {
	var value integrationsdk.Event
	path := "/v1/inbound/events/" + url.PathEscape(strings.TrimSpace(id)) + "?workspace_id=" + url.QueryEscape(strings.TrimSpace(workspaceID))
	err := c.call(ctx, http.MethodGet, path, nil, &value)
	return value, err
}

func (c *remoteClient) ReplayEvent(ctx context.Context, workspaceID, id string) (integrationsdk.Event, error) {
	var value integrationsdk.Event
	path := "/v1/inbound/events/" + url.PathEscape(strings.TrimSpace(id)) + "/replay?workspace_id=" + url.QueryEscape(strings.TrimSpace(workspaceID))
	err := c.call(ctx, http.MethodPost, path, nil, &value)
	return value, err
}

var _ integrationsdk.Operations = (*remoteClient)(nil)

func (c *remoteClient) call(ctx context.Context, method, path string, body any, target any) error {
	reference, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("parse Integration SaaS path: %w", err)
	}
	endpoint := c.base.ResolveReference(reference)
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode Integration SaaS request: %w", err)
		}
		reader = bytes.NewReader(payload)
	}
	request, err := http.NewRequestWithContext(ctx, method, endpoint.String(), reader)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Domainry-Runtime-ID", c.runtimeID)
	if c.token != "" {
		request.Header.Set("Authorization", "Bearer "+c.token)
	}
	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("call Integration SaaS: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		limited, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return fmt.Errorf("Integration SaaS returned %s: %s", response.Status, strings.TrimSpace(string(limited)))
	}
	if target == nil {
		return nil
	}
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf("decode Integration SaaS response: %w", err)
	}
	return nil
}

var _ saashost.Factory = (*Factory)(nil)
var _ integrationsdk.Catalog = (*remoteClient)(nil)
var _ integrationsdk.Requirements = (*remoteClient)(nil)
var _ integrationsdk.Delivery = (*remoteClient)(nil)
var _ integrationsdk.WebPushSubscriptions = (*remoteClient)(nil)
var _ integrationsdk.Management = (*remoteClient)(nil)
var _ integrationsdk.ManagementBinding = (*binding)(nil)
