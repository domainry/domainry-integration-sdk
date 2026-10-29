package integrationsdk

import (
	"strings"

	actioncontract "github.com/domainry/domainry-foundation/action"
)

const IntegrationHTTPAdapterContractVersion = "domainry-integration-http-adapter-v1"

const (
	ActionIntegrationConnectionAccountsReadAccess = "integration.connection_accounts.read_access"
	ActionIntegrationConnectionAccountsRead       = "integration.connection_accounts.read"
	ActionIntegrationOAuthApplicationsList        = "integration.oauth_applications.list"
	ActionIntegrationOAuthApplicationsUpsert      = "integration.oauth_applications.upsert"
	ActionIntegrationOAuthAuthorizationsOptions   = "integration.oauth_authorizations.options"
	ActionIntegrationOAuthAuthorizationsStart     = "integration.oauth_authorizations.start"
	ActionIntegrationOAuthAuthorizationsGet       = "integration.oauth_authorizations.get"
	ActionIntegrationOAuthAuthorizationsComplete  = "integration.oauth_authorizations.complete"

	CapabilityIntegrationConnections   = "integration.connections"
	CapabilityIntegrationAccounts      = "integration.connection_accounts"
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
	ActionIntegrationConnectionsRegisterAccount         = "integration.connections.register_account"
	ActionIntegrationConnectionAccountsList             = "integration.connection_accounts.list"
	ActionIntegrationConnectionAccountsGet              = "integration.connection_accounts.get"
	ActionIntegrationConnectionAccountsTest             = "integration.connection_accounts.test"
	ActionIntegrationConnectionAccountsRetryBackground  = "integration.connection_accounts.retry_background"
	ActionIntegrationConnectionAccountsRevoke           = "integration.connection_accounts.revoke"
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
	ActionIntegrationRoutedWebhooksIngest               = "integration.routed_webhooks.ingest"
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

type HTTPAdapterContract struct {
	ContractVersion string              `json:"contract_version"`
	Owner           string              `json:"owner"`
	Name            string              `json:"name"`
	Routes          []HTTPRouteContract `json:"routes"`
}

// IntegrationHTTPAdapterContract is the deployment-neutral HTTP contract
// implemented by both the embedded module and the SaaS transport. Runtime and
// Control Plane may mount or generate from it, but must not recreate it.
func IntegrationHTTPAdapterContract() HTTPAdapterContract {
	routes := integrationHTTPRoutes()
	return HTTPAdapterContract{
		ContractVersion: IntegrationHTTPAdapterContractVersion,
		Owner:           "integration",
		Name:            "integration_product",
		Routes:          routes,
	}
}

// IntegrationAuthorizationActions is the complete Integration-owned action
// manifest. Account writes are deliberately SDK-only: callers receive no HTTP
// route that could bypass the Agent tool confirmation and execution ledger.
func IntegrationAuthorizationActions() ([]actioncontract.ActionDefinition, error) {
	routes := integrationHTTPRoutes()
	actions := make([]actioncontract.ActionDefinition, 0, len(routes)+1)
	for _, route := range routes {
		actions = append(actions, route.Action)
	}
	permission := ConnectionAccountWritePermission()
	actions = append(actions, actioncontract.ActionDefinition{
		Key: permission.Key, Owner: permission.Owner, SourceKind: "integration_account_write",
		CapabilityKey: CapabilityIntegrationAccounts, CapabilityLabel: "Current-user connection accounts",
		OperationKey: permission.OperationKey, OperationLabel: permission.Label, Label: permission.Label,
		Exposures:     []actioncontract.Exposure{actioncontract.ExposurePublic},
		Authorization: actioncontract.Authorization{Strategy: actioncontract.AuthorizationAuthenticated},
		NonHTTP:       []actioncontract.NonHTTPBinding{{Kind: "sdk", InvocationKey: permission.Key}},
		EffectClass:   actioncontract.EffectWrite, RiskLevel: actioncontract.RiskMedium,
		IdempotencyDecision: "owner_receipt_reconcile", AuditClass: "integration_account_write",
		LifecycleStatus: actioncontract.LifecycleActive, Permission: &permission,
	})
	return actions, nil
}

func integrationHTTPRoutes() []HTTPRouteContract {
	admin := func(key, pattern, browserClientMethod string) HTTPRouteContract {
		return integrationHTTPRoute(key, pattern, browserClientMethod, []actioncontract.Exposure{actioncontract.ExposureManagement}, actioncontract.AuthorizationAuthenticated, true)
	}
	member := func(key, pattern, browserClientMethod string) HTTPRouteContract {
		return integrationHTTPRoute(key, pattern, browserClientMethod, []actioncontract.Exposure{actioncontract.ExposurePublic}, actioncontract.AuthorizationAuthenticated, true)
	}
	user := func(key, pattern, browserClientMethod string) HTTPRouteContract {
		return integrationHTTPRoute(key, pattern, browserClientMethod, []actioncontract.Exposure{actioncontract.ExposurePublic}, actioncontract.AuthorizationAuthenticated, false)
	}
	signed := func(key, pattern string) HTTPRouteContract {
		return integrationHTTPRoute(key, pattern, "", []actioncontract.Exposure{actioncontract.ExposurePublic}, actioncontract.AuthorizationSigned, false)
	}
	return []HTTPRouteContract{
		admin(ActionIntegrationCatalogRead, "GET /integration/catalog", "catalog"),
		admin(ActionIntegrationConnectorsList, "GET /integration/connectors", "connectors"),
		admin(ActionIntegrationConnectionsList, "GET /integration/connections", "listConnections"),
		admin(ActionIntegrationConnectionsGet, "GET /integration/connections/{connectionKey}", "getConnection"),
		admin(ActionIntegrationConnectionsValidate, "POST /integration/connections/{connectionKey}/validate", "validateConnection"),
		admin(ActionIntegrationConnectionsUpsert, "PUT /integration/connections/{connectionKey}", "upsertConnection"),
		admin(ActionIntegrationConnectionsDelete, "DELETE /integration/connections/{connectionKey}", "deleteConnection"),
		admin(ActionIntegrationConnectionsDisable, "POST /integration/connections/{connectionKey}/disable", "disableConnection"),
		admin(ActionIntegrationConnectionsTestOperation, "POST /integration/connections/{connectionKey}/test-operation", "testOperation"),
		admin(ActionIntegrationConnectionsRegisterAccount, "POST /integration/connections/{connectionKey}/account", "registerConnectionAccount"),
		member(ActionIntegrationConnectionAccountsList, "GET /integration/connection-accounts", "listConnectionAccounts"),
		member(ActionIntegrationConnectionAccountsGet, "GET /integration/connection-accounts/{connectionKey}", "getConnectionAccount"),
		member(ActionIntegrationConnectionAccountsTest, "POST /integration/connection-accounts/{connectionKey}/test", "testConnectionAccount"),
		member(ActionIntegrationConnectionAccountsRetryBackground, "POST /integration/connection-accounts/{connectionKey}/background/{taskKey}/retry", "retryConnectionAccountBackgroundTask"),
		member(ActionIntegrationConnectionAccountsRevoke, "POST /integration/connection-accounts/{connectionKey}/revoke", "revokeConnectionAccount"),
		member(ActionIntegrationConnectionAccountsReadAccess, "POST /integration/connection-accounts/{connectionKey}/read-access", "authorizeConnectionAccountRead"),
		member(ActionIntegrationConnectionAccountsRead, "POST /integration/connection-accounts/{connectionKey}/read", "readConnectionAccount"),
		admin(ActionIntegrationOAuthApplicationsList, "GET /integration/oauth-applications", "listOAuthApplications"),
		admin(ActionIntegrationOAuthApplicationsUpsert, "PUT /integration/oauth-applications/{applicationKey}", "upsertOAuthApplication"),
		member(ActionIntegrationOAuthAuthorizationsOptions, "GET /integration/oauth-authorizations/options", "listOAuthAuthorizationOptions"),
		member(ActionIntegrationOAuthAuthorizationsStart, "POST /integration/oauth-authorizations", "startOAuthAuthorization"),
		member(ActionIntegrationOAuthAuthorizationsGet, "GET /integration/oauth-authorizations/{sessionID}", "getOAuthAuthorization"),
		member(ActionIntegrationOAuthAuthorizationsComplete, "POST /integration/oauth-authorizations/callback", "completeOAuthAuthorization"),
		admin(ActionIntegrationSecretsList, "GET /integration/secrets", "listSecrets"),
		admin(ActionIntegrationSecretsUpsert, "PUT /integration/secrets/{secretKey}", "upsertSecret"),
		admin(ActionIntegrationSecretsDisable, "POST /integration/secrets/{secretKey}/disable", "disableSecret"),
		admin(ActionIntegrationSecretsRotate, "POST /integration/secrets/{secretKey}/rotate", "rotateSecret"),
		admin(ActionIntegrationSecretsExpire, "POST /integration/secrets/{secretKey}/expire", "expireSecret"),
		admin(ActionIntegrationSecretsRevoke, "POST /integration/secrets/{secretKey}/revoke", "revokeSecret"),
		admin(ActionIntegrationAPIKeysList, "GET /integration/api-keys", "listAPIKeys"),
		admin(ActionIntegrationAPIKeysCreate, "POST /integration/api-keys", "createAPIKey"),
		admin(ActionIntegrationAPIKeysDisable, "POST /integration/api-keys/{apiKey}/disable", "disableAPIKey"),
		admin(ActionIntegrationAPIKeysRotate, "POST /integration/api-keys/{apiKey}/rotate", "rotateAPIKey"),
		admin(ActionIntegrationExternalIdentitiesList, "GET /integration/external-identities", "listExternalIdentities"),
		admin(ActionIntegrationExternalIdentitiesUpsert, "PUT /integration/external-identities/{identityKey}", "upsertExternalIdentity"),
		admin(ActionIntegrationExternalIdentitiesDisable, "POST /integration/external-identities/{identityKey}/disable", "disableExternalIdentity"),
		admin(ActionIntegrationExternalIdentitiesResolve, "POST /integration/external-identities/resolve", "resolveExternalIdentity"),
		admin(ActionIntegrationWebhookSubscriptionsList, "GET /integration/webhook-subscriptions", "listWebhookSubscriptions"),
		admin(ActionIntegrationWebhookSubscriptionsUpsert, "PUT /integration/webhook-subscriptions/{subscriptionKey}", "upsertWebhookSubscription"),
		admin(ActionIntegrationWebhookSubscriptionsDelete, "DELETE /integration/webhook-subscriptions/{subscriptionKey}", "deleteWebhookSubscription"),
		admin(ActionIntegrationWebhookSubscriptionsDisable, "POST /integration/webhook-subscriptions/{subscriptionKey}/disable", "disableWebhookSubscription"),
		user(ActionIntegrationWebPushReadiness, "GET /integration/web-push/readiness", "webPushReadiness"),
		user(ActionIntegrationWebPushSubscriptionsList, "GET /integration/web-push/subscriptions", "webPushSubscriptions"),
		user(ActionIntegrationWebPushSubscriptionsUpsert, "PUT /integration/web-push/subscriptions/{subscriptionID}", "upsertWebPushSubscription"),
		user(ActionIntegrationWebPushSubscriptionsRevoke, "POST /integration/web-push/subscriptions/{subscriptionID}/revoke", "revokeWebPushSubscription"),
		admin(ActionIntegrationWebPushSubscriptionsCleanupExpired, "POST /integration/web-push/subscriptions/cleanup-expired", "cleanupExpiredWebPushSubscriptions"),
		admin(ActionIntegrationInvocationsList, "GET /integration/invocations", "listInvocations"),
		admin(ActionIntegrationInvocationsGet, "GET /integration/invocations/{invocationID}", "getInvocation"),
		admin(ActionIntegrationEventsList, "GET /integration/events", "listEvents"),
		admin(ActionIntegrationEventsGet, "GET /integration/events/{eventID}", "getEvent"),
		admin(ActionIntegrationEventsReplay, "POST /integration/events/{eventID}/replay", "replayEvent"),
		signed(ActionIntegrationWebhooksIngest, "POST /integration/webhooks/{workspaceID}/{connectorKey}/{connectionKey}"),
		signed(ActionIntegrationRoutedWebhooksIngest, "POST /integration/routed-webhooks/{workspaceID}/{connectorKey}/{connectionKey}"),
	}
}

func integrationHTTPRoute(key, pattern, browserClientMethod string, exposures []actioncontract.Exposure, strategy actioncontract.AuthorizationStrategy, requirePermission bool) HTTPRouteContract {
	method, path, _ := strings.Cut(strings.TrimSpace(pattern), " ")
	effect, risk, idempotency, auditClass := actioncontract.EffectWrite, actioncontract.RiskMedium, "caller_key_or_natural_resource_identity", "integration_owner_mutation"
	if method == "GET" || method == "HEAD" || method == "OPTIONS" || key == ActionIntegrationConnectionAccountsReadAccess || key == ActionIntegrationConnectionAccountsRead {
		effect, risk, idempotency, auditClass = actioncontract.EffectRead, actioncontract.RiskLow, "not_applicable", "integration_owner_read"
	}
	if strings.HasPrefix(path, "/integration/webhooks/") || strings.HasPrefix(path, "/integration/routed-webhooks/") {
		risk, idempotency, auditClass = actioncontract.RiskHigh, "provider_event_identity", "integration_webhook_ingress"
	}
	separator := strings.LastIndex(key, ".")
	capabilityKey, capabilityLabel := integrationHTTPCapability(key)
	label := strings.TrimSpace(browserClientMethod)
	if label == "" {
		label = key
	}
	definition := actioncontract.ActionDefinition{
		Key: key, Owner: "module:integration", SourceKind: "module_http", CapabilityKey: capabilityKey, CapabilityLabel: capabilityLabel,
		OperationKey: key[separator+1:], OperationLabel: label, Label: label, Exposures: exposures,
		HTTP: &actioncontract.HTTPBinding{Method: method, RouteTemplate: path}, EffectClass: effect, RiskLevel: risk,
		IdempotencyDecision: idempotency, AuditClass: auditClass, LifecycleStatus: actioncontract.LifecycleActive,
		Authorization: actioncontract.Authorization{Strategy: strategy},
	}
	if requirePermission {
		resourceKey, operationKey := key[:separator], key[separator+1:]
		localized := integrationPermissionText(capabilityKey, resourceKey, operationKey)
		definition.OperationLabel, definition.Label = localized["en"].OperationLabel, localized["en"].Label
		definition.Permission = &actioncontract.PermissionDefinition{
			Key: key, Owner: definition.Owner, ModuleKey: "integration", GroupKey: capabilityKey,
			ResourceKey: resourceKey, OperationKey: operationKey, Label: localized["en"].Label, Category: "Integration",
			GrantMode: actioncontract.PermissionGrantInternal, ScopeKind: actioncontract.PermissionScopeNone,
			DefaultLocale: "en", I18n: localized, LifecycleStatus: actioncontract.LifecycleActive,
		}
	}
	if strategy == actioncontract.AuthorizationSigned {
		definition.Authorization = actioncontract.Authorization{Strategy: actioncontract.AuthorizationSigned, PolicyKey: "integration.webhook.signature"}
	}
	return HTTPRouteContract{
		Action: definition, BrowserClientPackage: integrationBrowserClientPackage(browserClientMethod), BrowserClientMethod: browserClientMethod,
	}
}

func integrationPermissionText(capabilityKey, resourceKey, operationKey string) actioncontract.PermissionLocalizedText {
	groupEnglish := map[string]string{
		CapabilityIntegrationConnections: "Connector connections", CapabilityIntegrationAccounts: "Personal connection accounts",
		CapabilityIntegrationCredentials: "Credentials and identities", CapabilityIntegrationOperations: "Invocations and events",
		CapabilityIntegrationSubscriptions: "Webhook subscriptions",
	}[capabilityKey]
	groupChinese := map[string]string{
		CapabilityIntegrationConnections: "连接器与连接", CapabilityIntegrationAccounts: "个人连接账号",
		CapabilityIntegrationCredentials: "凭据与外部身份", CapabilityIntegrationOperations: "调用与事件",
		CapabilityIntegrationSubscriptions: "Webhook 订阅",
	}[capabilityKey]
	resourceEnglish := map[string]string{
		"integration.catalog": "Integration catalog", "integration.connectors": "Connectors", "integration.connections": "Connections",
		"integration.connection_accounts": "Connection accounts", "integration.oauth_applications": "OAuth applications",
		"integration.oauth_authorizations": "OAuth authorizations", "integration.secrets": "Integration secrets",
		"integration.api_keys": "API keys", "integration.external_identities": "External identities",
		"integration.webhook_subscriptions": "Webhook subscriptions", "integration.web_push_subscriptions": "Web push subscriptions",
		"integration.invocations": "Integration invocations", "integration.events": "Integration events",
	}[resourceKey]
	resourceChinese := map[string]string{
		"integration.catalog": "集成目录", "integration.connectors": "连接器", "integration.connections": "集成连接",
		"integration.connection_accounts": "连接账号", "integration.oauth_applications": "OAuth 应用",
		"integration.oauth_authorizations": "OAuth 授权", "integration.secrets": "集成密钥",
		"integration.api_keys": "API 密钥", "integration.external_identities": "外部身份",
		"integration.webhook_subscriptions": "Webhook 订阅", "integration.web_push_subscriptions": "Web Push 订阅",
		"integration.invocations": "集成调用", "integration.events": "集成事件",
	}[resourceKey]
	operationEnglish := map[string]string{
		"read": "View", "list": "List", "get": "View details", "validate": "Validate", "upsert": "Create or update",
		"delete": "Delete", "disable": "Disable", "test_operation": "Test operation", "register_account": "Register account",
		"test": "Test", "retry_background": "Retry background task", "revoke": "Revoke", "read_access": "Authorize read access",
		"write": "Write authorized data", "options": "List options", "start": "Start", "complete": "Complete",
		"rotate": "Rotate", "expire": "Expire", "create": "Create", "resolve": "Resolve", "cleanup_expired": "Clean up expired",
		"replay": "Replay",
	}[operationKey]
	operationChinese := map[string]string{
		"read": "查看", "list": "查看列表", "get": "查看详情", "validate": "校验", "upsert": "新建或更新",
		"delete": "删除", "disable": "停用", "test_operation": "测试操作", "register_account": "登记账号",
		"test": "测试", "retry_background": "重试后台任务", "revoke": "撤销", "read_access": "授权读取",
		"write": "写入已授权数据", "options": "查看可选项", "start": "开始", "complete": "完成",
		"rotate": "轮换", "expire": "设为过期", "create": "创建", "resolve": "解析", "cleanup_expired": "清理过期项",
		"replay": "重放",
	}[operationKey]
	return actioncontract.PermissionLocalizedText{
		"en":    {ModuleLabel: "Integrations", GroupLabel: groupEnglish, ResourceLabel: resourceEnglish, OperationLabel: operationEnglish, Label: resourceEnglish + " · " + operationEnglish, Description: "Provides access to " + resourceEnglish + " for the “" + operationEnglish + "” operation."},
		"zh-CN": {ModuleLabel: "集成", GroupLabel: groupChinese, ResourceLabel: resourceChinese, OperationLabel: operationChinese, Label: resourceChinese + " · " + operationChinese, Description: "允许对“" + resourceChinese + "”执行“" + operationChinese + "”操作。"},
	}
}

func integrationHTTPCapability(key string) (string, string) {
	switch {
	case strings.HasPrefix(key, "integration.connection_accounts."), strings.HasPrefix(key, "integration.oauth_authorizations."), key == ActionIntegrationConnectionsRegisterAccount:
		return CapabilityIntegrationAccounts, "Current-user connection accounts"
	case strings.HasPrefix(key, "integration.secrets."), strings.HasPrefix(key, "integration.api_keys."), strings.HasPrefix(key, "integration.external_identities."):
		return CapabilityIntegrationCredentials, "Integration credentials"
	case strings.HasPrefix(key, "integration.webhook_subscriptions."), strings.HasPrefix(key, "integration.web_push."), strings.HasPrefix(key, "integration.web_push_subscriptions."), strings.HasPrefix(key, "integration.webhooks."), strings.HasPrefix(key, "integration.routed_webhooks."):
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
