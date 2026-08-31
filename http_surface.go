package integrationsdk

import (
	"strings"
	"unicode"
)

const IntegrationHTTPSurfaceContractVersion = "domainry-integration-http-surface-v1"

type HTTPRouteContract struct {
	Pattern             string   `json:"pattern"`
	Exposures           []string `json:"exposures"`
	Authentication      string   `json:"authentication"`
	Permission          string   `json:"permission,omitempty"`
	AnyPermissions      []string `json:"any_permissions,omitempty"`
	PrincipalOnly       bool     `json:"principal_only,omitempty"`
	EffectClass         string   `json:"effect_class"`
	HighRiskPolicy      string   `json:"high_risk_policy"`
	IdempotencyDecision string   `json:"idempotency_decision"`
	AuditClass          string   `json:"audit_class"`
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
	admin := func(pattern string) HTTPRouteContract {
		return integrationHTTPRoute(pattern, []string{"tenant_admin"}, "authenticated", "", []string{"workspace.admin"}, false)
	}
	user := func(pattern string) HTTPRouteContract {
		return integrationHTTPRoute(pattern, []string{"public"}, "authenticated", "", nil, true)
	}
	public := func(pattern string) HTTPRouteContract {
		return integrationHTTPRoute(pattern, []string{"public"}, "anonymous", "", nil, false)
	}
	return []HTTPRouteContract{
		admin("GET /tenant-admin/integrations/catalog"),
		admin("GET /tenant-admin/integrations/connectors"),
		admin("GET /tenant-admin/integrations/connections"),
		admin("GET /tenant-admin/integrations/connections/{connectionKey}"),
		admin("POST /tenant-admin/integrations/connections/{connectionKey}/validate"),
		admin("PUT /tenant-admin/integrations/connections/{connectionKey}"),
		admin("DELETE /tenant-admin/integrations/connections/{connectionKey}"),
		admin("POST /tenant-admin/integrations/connections/{connectionKey}/disable"),
		admin("POST /tenant-admin/integrations/connections/{connectionKey}/test-operation"),
		admin("GET /tenant-admin/integrations/secrets"),
		admin("PUT /tenant-admin/integrations/secrets/{secretKey}"),
		admin("POST /tenant-admin/integrations/secrets/{secretKey}/disable"),
		admin("POST /tenant-admin/integrations/secrets/{secretKey}/rotate"),
		admin("POST /tenant-admin/integrations/secrets/{secretKey}/expire"),
		admin("POST /tenant-admin/integrations/secrets/{secretKey}/revoke"),
		admin("GET /tenant-admin/integrations/api-keys"),
		admin("POST /tenant-admin/integrations/api-keys"),
		admin("POST /tenant-admin/integrations/api-keys/{apiKey}/disable"),
		admin("POST /tenant-admin/integrations/api-keys/{apiKey}/rotate"),
		admin("GET /tenant-admin/integrations/external-identities"),
		admin("PUT /tenant-admin/integrations/external-identities/{identityKey}"),
		admin("POST /tenant-admin/integrations/external-identities/{identityKey}/disable"),
		admin("POST /tenant-admin/integrations/external-identities/resolve"),
		admin("GET /tenant-admin/integrations/webhook-subscriptions"),
		admin("PUT /tenant-admin/integrations/webhook-subscriptions/{subscriptionKey}"),
		admin("DELETE /tenant-admin/integrations/webhook-subscriptions/{subscriptionKey}"),
		admin("POST /tenant-admin/integrations/webhook-subscriptions/{subscriptionKey}/disable"),
		user("GET /business/notifications/web-push/readiness"),
		user("GET /business/notifications/web-push/subscriptions"),
		user("PUT /business/notifications/web-push/subscriptions/{subscriptionID}"),
		user("POST /business/notifications/web-push/subscriptions/{subscriptionID}/revoke"),
		admin("POST /integrations/web-push/subscriptions/cleanup-expired"),
		admin("GET /tenant-admin/integrations/invocations"),
		admin("GET /tenant-admin/integrations/invocations/{invocationID}"),
		admin("GET /tenant-admin/integrations/events"),
		admin("GET /tenant-admin/integrations/events/{eventID}"),
		admin("POST /tenant-admin/integrations/events/{eventID}/replay"),
		public("POST /integrations/webhooks/{workspaceID}/{connectorKey}/{connectionKey}"),
	}
}

func integrationHTTPRoute(pattern string, exposures []string, authentication, permission string, anyPermissions []string, principalOnly bool) HTTPRouteContract {
	method, path, _ := strings.Cut(strings.TrimSpace(pattern), " ")
	effect, idempotency, auditClass := "write", "caller_key_or_natural_resource_identity", "integration_owner_mutation"
	if method == "GET" || method == "HEAD" || method == "OPTIONS" {
		effect, idempotency, auditClass = "read", "not_applicable", "integration_owner_read"
	}
	if strings.HasPrefix(path, "/integrations/webhooks/") {
		idempotency, auditClass = "provider_event_identity", "integration_webhook_ingress"
	}
	return HTTPRouteContract{
		Pattern: pattern, Exposures: append([]string(nil), exposures...), Authentication: authentication,
		Permission: permission, AnyPermissions: append([]string(nil), anyPermissions...), PrincipalOnly: principalOnly,
		EffectClass: effect, HighRiskPolicy: "none", IdempotencyDecision: idempotency, AuditClass: auditClass,
	}
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
