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

	integrationsdk "github.com/domainry/domainry-integration-sdk"
	"github.com/domainry/domainry-integration-sdk/saashost"
)

type Options struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type Factory struct{ options Options }

func NewFactory(options Options) *Factory { return &Factory{options: options} }

func (*Factory) DeploymentMode() integrationsdk.DeploymentMode {
	return integrationsdk.DeploymentModeSaaS
}

func (f *Factory) OpenSaaS(_ context.Context, application integrationsdk.ApplicationRef, _ saashost.Host) (integrationsdk.Binding, error) {
	if err := application.Validate(); err != nil {
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
	return &binding{client: transport}, nil
}

type binding struct{ client *remoteClient }

func (*binding) Descriptor() integrationsdk.Descriptor {
	return integrationsdk.Descriptor{ProtocolVersion: integrationsdk.ProtocolVersionV1, Mode: integrationsdk.DeploymentModeSaaS, Capabilities: []string{"catalog.read", "requirements.connections.sync", "delivery.accept", "delivery.query", "web_push_subscriptions.manage"}}
}
func (b *binding) Catalog() integrationsdk.Catalog                           { return b.client }
func (b *binding) Requirements() integrationsdk.Requirements                 { return b.client }
func (b *binding) WebPushSubscriptions() integrationsdk.WebPushSubscriptions { return b.client }
func (b *binding) Delivery() integrationsdk.Delivery                         { return b.client }
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
