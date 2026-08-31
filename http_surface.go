package integrationsdk

import (
	"strings"
	"unicode"
)

const IntegrationHTTPSurfaceContractVersion = "domainry-integration-http-surface-v1"

type HTTPRouteContract struct {
	Pattern              string   `json:"pattern"`
	BrowserClientPackage string   `json:"browser_client_package,omitempty"`
	BrowserClientMethod  string   `json:"browser_client_method,omitempty"`
	Exposures            []string `json:"exposures"`
	Authentication       string   `json:"authentication"`
	Permission           string   `json:"permission,omitempty"`
	AnyPermissions       []string `json:"any_permissions,omitempty"`
	PrincipalOnly        bool     `json:"principal_only,omitempty"`
	EffectClass          string   `json:"effect_class"`
	HighRiskPolicy       string   `json:"high_risk_policy"`
	IdempotencyDecision  string   `json:"idempotency_decision"`
	AuditClass           string   `json:"audit_class"`
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
	admin := func(pattern, browserClientMethod string) HTTPRouteContract {
		return integrationHTTPRoute(pattern, browserClientMethod, []string{"tenant_admin"}, "authenticated", "", []string{"workspace.admin"}, false)
	}
	user := func(pattern, browserClientMethod string) HTTPRouteContract {
		return integrationHTTPRoute(pattern, browserClientMethod, []string{"public"}, "authenticated", "", nil, true)
	}
	public := func(pattern string) HTTPRouteContract {
		return integrationHTTPRoute(pattern, "", []string{"public"}, "anonymous", "", nil, false)
	}
	return []HTTPRouteContract{
		admin("GET /tenant-admin/integrations/catalog", "catalog"),
		admin("GET /tenant-admin/integrations/connectors", "connectors"),
		admin("GET /tenant-admin/integrations/connections", "listConnections"),
		admin("GET /tenant-admin/integrations/connections/{connectionKey}", "getConnection"),
		admin("POST /tenant-admin/integrations/connections/{connectionKey}/validate", "validateConnection"),
		admin("PUT /tenant-admin/integrations/connections/{connectionKey}", "upsertConnection"),
		admin("DELETE /tenant-admin/integrations/connections/{connectionKey}", "deleteConnection"),
		admin("POST /tenant-admin/integrations/connections/{connectionKey}/disable", "disableConnection"),
		admin("POST /tenant-admin/integrations/connections/{connectionKey}/test-operation", "testOperation"),
		admin("GET /tenant-admin/integrations/secrets", "listSecrets"),
		admin("PUT /tenant-admin/integrations/secrets/{secretKey}", "upsertSecret"),
		admin("POST /tenant-admin/integrations/secrets/{secretKey}/disable", "disableSecret"),
		admin("POST /tenant-admin/integrations/secrets/{secretKey}/rotate", "rotateSecret"),
		admin("POST /tenant-admin/integrations/secrets/{secretKey}/expire", "expireSecret"),
		admin("POST /tenant-admin/integrations/secrets/{secretKey}/revoke", "revokeSecret"),
		admin("GET /tenant-admin/integrations/api-keys", "listAPIKeys"),
		admin("POST /tenant-admin/integrations/api-keys", "createAPIKey"),
		admin("POST /tenant-admin/integrations/api-keys/{apiKey}/disable", "disableAPIKey"),
		admin("POST /tenant-admin/integrations/api-keys/{apiKey}/rotate", "rotateAPIKey"),
		admin("GET /tenant-admin/integrations/external-identities", "listExternalIdentities"),
		admin("PUT /tenant-admin/integrations/external-identities/{identityKey}", "upsertExternalIdentity"),
		admin("POST /tenant-admin/integrations/external-identities/{identityKey}/disable", "disableExternalIdentity"),
		admin("POST /tenant-admin/integrations/external-identities/resolve", "resolveExternalIdentity"),
		admin("GET /tenant-admin/integrations/webhook-subscriptions", "listWebhookSubscriptions"),
		admin("PUT /tenant-admin/integrations/webhook-subscriptions/{subscriptionKey}", "upsertWebhookSubscription"),
		admin("DELETE /tenant-admin/integrations/webhook-subscriptions/{subscriptionKey}", "deleteWebhookSubscription"),
		admin("POST /tenant-admin/integrations/webhook-subscriptions/{subscriptionKey}/disable", "disableWebhookSubscription"),
		user("GET /business/notifications/web-push/readiness", "webPushReadiness"),
		user("GET /business/notifications/web-push/subscriptions", "webPushSubscriptions"),
		user("PUT /business/notifications/web-push/subscriptions/{subscriptionID}", "upsertWebPushSubscription"),
		user("POST /business/notifications/web-push/subscriptions/{subscriptionID}/revoke", "revokeWebPushSubscription"),
		admin("POST /integrations/web-push/subscriptions/cleanup-expired", "cleanupExpiredWebPushSubscriptions"),
		admin("GET /tenant-admin/integrations/invocations", "listInvocations"),
		admin("GET /tenant-admin/integrations/invocations/{invocationID}", "getInvocation"),
		admin("GET /tenant-admin/integrations/events", "listEvents"),
		admin("GET /tenant-admin/integrations/events/{eventID}", "getEvent"),
		admin("POST /tenant-admin/integrations/events/{eventID}/replay", "replayEvent"),
		public("POST /integrations/webhooks/{workspaceID}/{connectorKey}/{connectionKey}"),
	}
}

func integrationHTTPRoute(pattern, browserClientMethod string, exposures []string, authentication, permission string, anyPermissions []string, principalOnly bool) HTTPRouteContract {
	method, path, _ := strings.Cut(strings.TrimSpace(pattern), " ")
	effect, idempotency, auditClass := "write", "caller_key_or_natural_resource_identity", "integration_owner_mutation"
	if method == "GET" || method == "HEAD" || method == "OPTIONS" {
		effect, idempotency, auditClass = "read", "not_applicable", "integration_owner_read"
	}
	if strings.HasPrefix(path, "/integrations/webhooks/") {
		idempotency, auditClass = "provider_event_identity", "integration_webhook_ingress"
	}
	return HTTPRouteContract{
		Pattern: pattern, BrowserClientPackage: integrationBrowserClientPackage(browserClientMethod), BrowserClientMethod: browserClientMethod,
		Exposures: append([]string(nil), exposures...), Authentication: authentication,
		Permission: permission, AnyPermissions: append([]string(nil), anyPermissions...), PrincipalOnly: principalOnly,
		EffectClass: effect, HighRiskPolicy: "none", IdempotencyDecision: idempotency, AuditClass: auditClass,
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
		method, path, found := strings.Cut(strings.TrimSpace(route.Pattern), " ")
		if !found {
			continue
		}
		operation := map[string]any{
			"operationId": integrationHTTPOperationID(method, path),
			"tags":        []string{"Integration"},
			"summary":     "Integration owner " + strings.ToLower(method) + " " + path,
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
		if route.Authentication == "anonymous" {
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
		result[route.Pattern] = operation
	}
	return result
}

func integrationHTTPOperationID(method, path string) string {
	parts := []string{strings.ToLower(strings.TrimSpace(method))}
	for _, segment := range strings.FieldsFunc(path, func(value rune) bool {
		return value == '/' || value == '-' || value == '{' || value == '}' || value == ':'
	}) {
		if segment == "tenant" || segment == "admin" || segment == "integrations" || segment == "integration" {
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
