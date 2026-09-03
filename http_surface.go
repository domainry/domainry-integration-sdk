package integrationsdk

import (
	"strings"
	"unicode"

	actioncontract "github.com/domainry/domainry-foundation/action"
)

const IntegrationHTTPSurfaceContractVersion = "domainry-integration-http-surface-v2"

const (
	CapabilityIntegrationConnections   = "integration.connections"
	CapabilityIntegrationCredentials   = "integration.credentials"
	CapabilityIntegrationOperations    = "integration.operations"
	CapabilityIntegrationSubscriptions = "integration.subscriptions"

	ActionIntegrationCatalogRead                        = "integration.catalog.read"
	ActionIntegrationConnectorsList                     = "integration.connectors.list"
	ActionIntegrationConnectionsList                    = "integration.connections.list"
	ActionIntegrationConnectionsGet                     = "integration.connections.get"
	ActionIntegrationConnectionsValidate                = "integration.connections.validate"
	ActionIntegrationConnectionsUpsert                  = "integration.connections.upsert"
	ActionIntegrationConnectionsDelete                  = "integration.connections.delete"
	ActionIntegrationConnectionsDisable                 = "integration.connections.disable"
	ActionIntegrationConnectionsTestOperation           = "integration.connections.test_operation"
	ActionIntegrationSecretsList                        = "integration.secrets.list"
	ActionIntegrationSecretsUpsert                      = "integration.secrets.upsert"
	ActionIntegrationSecretsDisable                     = "integration.secrets.disable"
	ActionIntegrationSecretsRotate                      = "integration.secrets.rotate"
	ActionIntegrationSecretsExpire                      = "integration.secrets.expire"
	ActionIntegrationSecretsRevoke                      = "integration.secrets.revoke"
	ActionIntegrationAPIKeysList                        = "integration.api_keys.list"
	ActionIntegrationAPIKeysCreate                      = "integration.api_keys.create"
	ActionIntegrationAPIKeysDisable                     = "integration.api_keys.disable"
	ActionIntegrationAPIKeysRotate                      = "integration.api_keys.rotate"
	ActionIntegrationExternalIdentitiesList             = "integration.external_identities.list"
	ActionIntegrationExternalIdentitiesUpsert           = "integration.external_identities.upsert"
	ActionIntegrationExternalIdentitiesDisable          = "integration.external_identities.disable"
	ActionIntegrationExternalIdentitiesResolve          = "integration.external_identities.resolve"
	ActionIntegrationWebhookSubscriptionsList           = "integration.webhook_subscriptions.list"
	ActionIntegrationWebhookSubscriptionsUpsert         = "integration.webhook_subscriptions.upsert"
	ActionIntegrationWebhookSubscriptionsDelete         = "integration.webhook_subscriptions.delete"
	ActionIntegrationWebhookSubscriptionsDisable        = "integration.webhook_subscriptions.disable"
	ActionIntegrationWebPushReadiness                   = "integration.web_push.readiness"
	ActionIntegrationWebPushSubscriptionsList           = "integration.web_push_subscriptions.list"
	ActionIntegrationWebPushSubscriptionsUpsert         = "integration.web_push_subscriptions.upsert"
	ActionIntegrationWebPushSubscriptionsRevoke         = "integration.web_push_subscriptions.revoke"
	ActionIntegrationWebPushSubscriptionsCleanupExpired = "integration.web_push_subscriptions.cleanup_expired"
	ActionIntegrationInvocationsList                    = "integration.invocations.list"
	ActionIntegrationInvocationsGet                     = "integration.invocations.get"
	ActionIntegrationEventsList                         = "integration.events.list"
	ActionIntegrationEventsGet                          = "integration.events.get"
	ActionIntegrationEventsReplay                       = "integration.events.replay"
	ActionIntegrationWebhooksIngest                     = "integration.webhooks.ingest"
)

type HTTPRouteContract struct {
	Action               actioncontract.ActionDefinition `json:"action"`
	BrowserClientPackage string                          `json:"browser_client_package,omitempty"`
	BrowserClientMethod  string                          `json:"browser_client_method,omitempty"`
}

func (route HTTPRouteContract) Pattern() string {
	if route.Action.HTTP == nil {
		return ""
	}
	return route.Action.HTTP.Method + " " + route.Action.HTTP.RouteTemplate
}

type HTTPSurfaceContract struct {
	ContractVersion string                    `json:"contract_version"`
	Owner           string                    `json:"owner"`
	Name            string                    `json:"name"`
	Routes          []HTTPRouteContract       `json:"routes"`
	OpenAPI         map[string]map[string]any `json:"openapi_operations"`
}

// IntegrationHTTPSurfaceContract is the deployment-neutral HTTP contract
// implemented by both the embedded module and the SaaS transport. Runtime and
// Control Plane may mount or generate from it, but must not recreate it.
func IntegrationHTTPSurfaceContract() HTTPSurfaceContract {
	routes := integrationHTTPRoutes()
	return HTTPSurfaceContract{
		ContractVersion: IntegrationHTTPSurfaceContractVersion,
		Owner:           "integration",
		Name:            "integration_product",
		Routes:          routes,
		OpenAPI:         integrationHTTPOperations(routes),
	}
}

func integrationHTTPRoutes() []HTTPRouteContract {
	admin := func(key, pattern, browserClientMethod string) HTTPRouteContract {
		return integrationHTTPRoute(key, pattern, browserClientMethod, []actioncontract.Exposure{actioncontract.ExposureTenantAdmin}, actioncontract.AuthorizationAuthenticated, true)
	}
	user := func(key, pattern, browserClientMethod string) HTTPRouteContract {
		return integrationHTTPRoute(key, pattern, browserClientMethod, []actioncontract.Exposure{actioncontract.ExposurePublic}, actioncontract.AuthorizationAuthenticated, false)
	}
	signed := func(key, pattern string) HTTPRouteContract {
		return integrationHTTPRoute(key, pattern, "", []actioncontract.Exposure{actioncontract.ExposurePublic}, actioncontract.AuthorizationSigned, false)
	}
	return []HTTPRouteContract{
		admin(ActionIntegrationCatalogRead, "GET /tenant-admin/integrations/catalog", "catalog"),
		admin(ActionIntegrationConnectorsList, "GET /tenant-admin/integrations/connectors", "connectors"),
		admin(ActionIntegrationConnectionsList, "GET /tenant-admin/integrations/connections", "listConnections"),
		admin(ActionIntegrationConnectionsGet, "GET /tenant-admin/integrations/connections/{connectionKey}", "getConnection"),
		admin(ActionIntegrationConnectionsValidate, "POST /tenant-admin/integrations/connections/{connectionKey}/validate", "validateConnection"),
		admin(ActionIntegrationConnectionsUpsert, "PUT /tenant-admin/integrations/connections/{connectionKey}", "upsertConnection"),
		admin(ActionIntegrationConnectionsDelete, "DELETE /tenant-admin/integrations/connections/{connectionKey}", "deleteConnection"),
		admin(ActionIntegrationConnectionsDisable, "POST /tenant-admin/integrations/connections/{connectionKey}/disable", "disableConnection"),
		admin(ActionIntegrationConnectionsTestOperation, "POST /tenant-admin/integrations/connections/{connectionKey}/test-operation", "testOperation"),
		admin(ActionIntegrationSecretsList, "GET /tenant-admin/integrations/secrets", "listSecrets"),
		admin(ActionIntegrationSecretsUpsert, "PUT /tenant-admin/integrations/secrets/{secretKey}", "upsertSecret"),
		admin(ActionIntegrationSecretsDisable, "POST /tenant-admin/integrations/secrets/{secretKey}/disable", "disableSecret"),
		admin(ActionIntegrationSecretsRotate, "POST /tenant-admin/integrations/secrets/{secretKey}/rotate", "rotateSecret"),
		admin(ActionIntegrationSecretsExpire, "POST /tenant-admin/integrations/secrets/{secretKey}/expire", "expireSecret"),
		admin(ActionIntegrationSecretsRevoke, "POST /tenant-admin/integrations/secrets/{secretKey}/revoke", "revokeSecret"),
		admin(ActionIntegrationAPIKeysList, "GET /tenant-admin/integrations/api-keys", "listAPIKeys"),
		admin(ActionIntegrationAPIKeysCreate, "POST /tenant-admin/integrations/api-keys", "createAPIKey"),
		admin(ActionIntegrationAPIKeysDisable, "POST /tenant-admin/integrations/api-keys/{apiKey}/disable", "disableAPIKey"),
		admin(ActionIntegrationAPIKeysRotate, "POST /tenant-admin/integrations/api-keys/{apiKey}/rotate", "rotateAPIKey"),
		admin(ActionIntegrationExternalIdentitiesList, "GET /tenant-admin/integrations/external-identities", "listExternalIdentities"),
		admin(ActionIntegrationExternalIdentitiesUpsert, "PUT /tenant-admin/integrations/external-identities/{identityKey}", "upsertExternalIdentity"),
		admin(ActionIntegrationExternalIdentitiesDisable, "POST /tenant-admin/integrations/external-identities/{identityKey}/disable", "disableExternalIdentity"),
		admin(ActionIntegrationExternalIdentitiesResolve, "POST /tenant-admin/integrations/external-identities/resolve", "resolveExternalIdentity"),
		admin(ActionIntegrationWebhookSubscriptionsList, "GET /tenant-admin/integrations/webhook-subscriptions", "listWebhookSubscriptions"),
		admin(ActionIntegrationWebhookSubscriptionsUpsert, "PUT /tenant-admin/integrations/webhook-subscriptions/{subscriptionKey}", "upsertWebhookSubscription"),
		admin(ActionIntegrationWebhookSubscriptionsDelete, "DELETE /tenant-admin/integrations/webhook-subscriptions/{subscriptionKey}", "deleteWebhookSubscription"),
		admin(ActionIntegrationWebhookSubscriptionsDisable, "POST /tenant-admin/integrations/webhook-subscriptions/{subscriptionKey}/disable", "disableWebhookSubscription"),
		user(ActionIntegrationWebPushReadiness, "GET /business/notifications/web-push/readiness", "webPushReadiness"),
		user(ActionIntegrationWebPushSubscriptionsList, "GET /business/notifications/web-push/subscriptions", "webPushSubscriptions"),
		user(ActionIntegrationWebPushSubscriptionsUpsert, "PUT /business/notifications/web-push/subscriptions/{subscriptionID}", "upsertWebPushSubscription"),
		user(ActionIntegrationWebPushSubscriptionsRevoke, "POST /business/notifications/web-push/subscriptions/{subscriptionID}/revoke", "revokeWebPushSubscription"),
		admin(ActionIntegrationWebPushSubscriptionsCleanupExpired, "POST /integrations/web-push/subscriptions/cleanup-expired", "cleanupExpiredWebPushSubscriptions"),
		admin(ActionIntegrationInvocationsList, "GET /tenant-admin/integrations/invocations", "listInvocations"),
		admin(ActionIntegrationInvocationsGet, "GET /tenant-admin/integrations/invocations/{invocationID}", "getInvocation"),
		admin(ActionIntegrationEventsList, "GET /tenant-admin/integrations/events", "listEvents"),
		admin(ActionIntegrationEventsGet, "GET /tenant-admin/integrations/events/{eventID}", "getEvent"),
		admin(ActionIntegrationEventsReplay, "POST /tenant-admin/integrations/events/{eventID}/replay", "replayEvent"),
		signed(ActionIntegrationWebhooksIngest, "POST /integrations/webhooks/{workspaceID}/{connectorKey}/{connectionKey}"),
	}
}

func integrationHTTPRoute(key, pattern, browserClientMethod string, exposures []actioncontract.Exposure, strategy actioncontract.AuthorizationStrategy, requirePermission bool) HTTPRouteContract {
	method, path, _ := strings.Cut(strings.TrimSpace(pattern), " ")
	effect, risk, idempotency, auditClass := actioncontract.EffectWrite, actioncontract.RiskMedium, "caller_key_or_natural_resource_identity", "integration_owner_mutation"
	if method == "GET" || method == "HEAD" || method == "OPTIONS" {
		effect, risk, idempotency, auditClass = actioncontract.EffectRead, actioncontract.RiskLow, "not_applicable", "integration_owner_read"
	}
	if strings.HasPrefix(path, "/integrations/webhooks/") {
		risk, idempotency, auditClass = actioncontract.RiskHigh, "provider_event_identity", "integration_webhook_ingress"
	}
	separator := strings.LastIndex(key, ".")
	capabilityKey, capabilityLabel := integrationHTTPCapability(key)
	label := strings.TrimSpace(browserClientMethod)
	if label == "" {
		label = key
	}
	definition := actioncontract.ActionDefinition{
		Key: key, Owner: "module:integration", SourceKind: "module_surface", CapabilityKey: capabilityKey, CapabilityLabel: capabilityLabel,
		OperationKey: key[separator+1:], OperationLabel: label, Label: label, Exposures: exposures,
		HTTP: &actioncontract.HTTPBinding{Method: method, RouteTemplate: path}, EffectClass: effect, RiskLevel: risk,
		IdempotencyDecision: idempotency, AuditClass: auditClass, LifecycleStatus: actioncontract.LifecycleActive,
		Authorization: actioncontract.Authorization{Strategy: strategy},
	}
	if requirePermission {
		definition.Permission = &actioncontract.PermissionDefinition{Key: key, Owner: definition.Owner, ResourceKey: key[:separator], OperationKey: key[separator+1:], Label: label, Category: "Integration", LifecycleStatus: actioncontract.LifecycleActive}
	}
	if strategy == actioncontract.AuthorizationSigned {
		definition.Authorization = actioncontract.Authorization{Strategy: actioncontract.AuthorizationSigned, PolicyKey: "integration.webhook.signature"}
	}
	return HTTPRouteContract{
		Action: definition, BrowserClientPackage: integrationBrowserClientPackage(browserClientMethod), BrowserClientMethod: browserClientMethod,
	}
}

func integrationHTTPCapability(key string) (string, string) {
	switch {
	case strings.HasPrefix(key, "integration.secrets."), strings.HasPrefix(key, "integration.api_keys."), strings.HasPrefix(key, "integration.external_identities."):
		return CapabilityIntegrationCredentials, "Integration credentials"
	case strings.HasPrefix(key, "integration.webhook_subscriptions."), strings.HasPrefix(key, "integration.web_push."), strings.HasPrefix(key, "integration.web_push_subscriptions."), strings.HasPrefix(key, "integration.webhooks."):
		return CapabilityIntegrationSubscriptions, "Inbound and push subscriptions"
	case strings.HasPrefix(key, "integration.invocations."), strings.HasPrefix(key, "integration.events."):
		return CapabilityIntegrationOperations, "Invocation and event operations"
	default:
		return CapabilityIntegrationConnections, "Connector connections"
	}
}

func integrationBrowserClientPackage(method string) string {
	if strings.TrimSpace(method) == "" {
		return ""
	}
	return "@domainry/integration-client"
}

func integrationHTTPOperations(routes []HTTPRouteContract) map[string]map[string]any {
	result := make(map[string]map[string]any, len(routes))
	for _, route := range routes {
		pattern := route.Pattern()
		method, _, found := strings.Cut(pattern, " ")
		if !found {
			continue
		}
		operation := map[string]any{
			"operationId": integrationHTTPOperationID(route.Action.Key),
			"tags":        []string{"Integration"},
			"summary":     route.Action.Label,
			"responses": map[string]any{
				integrationHTTPSuccessStatus(method): map[string]any{
					"description": "Integration owner response",
					"content": map[string]any{"application/json": map[string]any{
						"schema": map[string]any{"type": "object", "additionalProperties": true},
					}},
				},
			},
		}
		if route.BrowserClientMethod != "" {
			operation["x-domainry-owner-client-package"] = route.BrowserClientPackage
			operation["x-domainry-owner-client-method"] = route.BrowserClientMethod
		}
		if route.Action.Authorization.Strategy == actioncontract.AuthorizationSigned {
			operation["security"] = []any{}
		} else {
			operation["security"] = []map[string]any{{"BearerAuth": []string{}}}
		}
		if method != "GET" && method != "DELETE" {
			operation["requestBody"] = map[string]any{
				"required": false,
				"content": map[string]any{"application/json": map[string]any{
					"schema": map[string]any{"type": "object", "additionalProperties": true},
				}},
			}
		}
		result[pattern] = operation
	}
	return result
}

func integrationHTTPOperationID(actionKey string) string {
	parts := []string{}
	for _, segment := range strings.FieldsFunc(actionKey, func(value rune) bool {
		return value == '.' || value == '_' || value == '-'
	}) {
		if segment == "integration" {
			continue
		}
		runes := []rune(segment)
		if len(runes) != 0 {
			runes[0] = unicode.ToUpper(runes[0])
			parts = append(parts, string(runes))
		}
	}
	return "integration" + strings.Join(parts, "")
}

func integrationHTTPSuccessStatus(method string) string {
	if strings.EqualFold(strings.TrimSpace(method), "DELETE") {
		return "204"
	}
	return "200"
}
