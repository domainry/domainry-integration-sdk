export interface IntegrationRequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  body?: unknown
  headers?: Record<string, string>
  requestId?: string
  signal?: AbortSignal
}

export type IntegrationRequest = <T>(path: string, options?: IntegrationRequestOptions) => Promise<T>

export interface IntegrationResponse<T> {
  data: T
  headers: {
    get(name: string): string | null
  }
  status: number
}

export type IntegrationRequestWithResponse = <T>(
  path: string,
  options?: IntegrationRequestOptions,
) => Promise<IntegrationResponse<T>>

export interface IntegrationClientDependencies {
  request: IntegrationRequest
  requestWithResponse?: IntegrationRequestWithResponse
}

export interface IntegrationFieldContract {
  key: string
  name?: string
  label?: string
  type: string
  required?: boolean
  unique?: boolean
  default?: unknown
  default_value?: unknown
  config?: Record<string, unknown>
  validation?: Record<string, unknown>
  i18n?: Record<string, { name?: string; label?: string; description?: string }>
}

export interface IntegrationConnectorOperation {
  key: string
  name?: string
  description?: string
  method?: string
  execution_mode?: 'sync' | 'async'
  side_effect?: string
  input?: IntegrationFieldContract[]
  output?: IntegrationFieldContract[]
  timeout_default_seconds?: number
  timeout_max_seconds?: number
  idempotency_supported?: boolean
  compensation_operation?: string
  test_supported?: boolean
  dry_run_supported?: boolean
}

export interface IntegrationConnectorProvider {
  key: string
  name?: string
  adapter_key?: string
  readiness?: 'provider_available' | 'adapter_ready' | 'connection_ready'
  config_fields?: IntegrationFieldContract[]
  secret_fields?: IntegrationFieldContract[]
  operation_keys?: string[]
}

export interface IntegrationConnector {
  key: string
  name?: string
  description?: string
  type: string
  provider: string
  adapter_key?: string
  version?: string
  source?: string
  config_fields?: string[]
  secret_refs?: string[]
  readiness?: 'catalog_only' | 'adapter_ready' | 'connection_ready'
  definition_ready?: boolean
  adapter_ready?: boolean
  connection_ready?: boolean
  capabilities?: string[]
  providers?: IntegrationConnectorProvider[]
  operations?: IntegrationConnectorOperation[]
  config?: Record<string, unknown>
}

export interface IntegrationConnection {
  key: string
  connector_key: string
  provider_key: string
  name?: string
  status: 'draft' | 'configured' | 'verified' | 'active' | 'degraded' | 'disabled'
  config?: Record<string, unknown>
  secret_refs?: Record<string, string>
  created_at?: string
  updated_at?: string
}

export interface IntegrationConnectionInput {
  connector_key: string
  provider_key: string
  name?: string
  status?: IntegrationConnection['status']
  config?: Record<string, unknown>
  secret_refs?: Record<string, string>
}

export interface IntegrationCatalog {
  connectors: IntegrationConnector[]
  connections: IntegrationConnection[]
  count: number
}

export interface IntegrationWebhookSubscription {
  key: string
  workspace_id?: string
  name?: string
  connector_key: string
  connection_key: string
  event_types?: string[]
  status: 'active' | 'disabled'
  description?: string
  created_by?: string
  created_at?: string
  updated_at?: string
  disabled_at?: string
}

export interface IntegrationWebhookSubscriptionInput {
  name?: string
  connector_key: string
  connection_key: string
  event_types?: string[]
  status?: 'active' | 'disabled'
  description?: string
}

export interface IntegrationSecret {
  key: string
  kind: string
  status: 'active' | 'disabled' | 'expired' | 'revoked'
  configured: boolean
  description?: string
  fingerprint?: string
  created_at?: string
  updated_at?: string
  disabled_at?: string
  expires_at?: string
  rotated_at?: string
  revoked_at?: string
  last_tested_at?: string
  last_test_status?: string
  last_test_error?: string
}

export interface IntegrationSecretInput {
  key?: string
  kind?: string
  description?: string
  value?: string
  expires_at?: string
}

export interface IntegrationAPIKey {
  key: string
  workspace_id?: string
  name?: string
  token_prefix?: string
  actor_id: string
  role_key: string
  scopes?: string[]
  status: string
  expires_at?: string
  last_used_at?: string
  created_by?: string
  created_at?: string
  updated_at?: string
  disabled_at?: string
}

export interface IntegrationAPIKeyInput {
  key?: string
  name?: string
  actor_id: string
  role_key: string
  scopes?: string[]
  expires_at?: string
}

export interface IntegrationAPIKeyCredential {
  api_key: IntegrationAPIKey
  token: string
}

export interface IntegrationExternalIdentity {
  key: string
  workspace_id?: string
  provider: string
  external_subject: string
  external_subject_type?: string
  external_name?: string
  external_organization?: string
  external_department?: string
  external_group?: string
  external_bot_id?: string
  actor_id: string
  role_key: string
  status: string
  last_resolved_at?: string
  created_by?: string
  created_at?: string
  updated_at?: string
  disabled_at?: string
}

export interface IntegrationExternalIdentityInput {
  provider: string
  external_subject: string
  external_subject_type?: string
  external_name?: string
  external_organization?: string
  external_department?: string
  external_group?: string
  external_bot_id?: string
  actor_id: string
  role_key: string
  status?: string
}

export interface IntegrationDeliveryReceipt {
  message_id: string
  invocation_id: string
  status: string
  result_ref?: string
  error_code?: string
}

export interface IntegrationConnectorOperationTestResult {
  connection: IntegrationConnection
  operation: IntegrationConnectorOperation | string
  response: Record<string, unknown>
  receipt?: IntegrationDeliveryReceipt
  invocation?: {
    id: string
    status: string
    request_ref?: string
    response_ref?: string
    duration_ms?: number
    error?: string
  }
}

export interface IntegrationConnectionTestRequest {
  operation: string
  input: Record<string, unknown>
  confirm?: boolean
}

export interface IntegrationWebPushSubscription {
  id: string
  user_id: string
  endpoint_hash: string
  status: string
  created_at: string
  updated_at: string
  expires_at?: string
  revoked_at?: string
}

export interface IntegrationWebPushSubscriptionUpsertRequest {
  endpoint: string
  p256dh: string
  auth: string
  expires_at?: string
}

export interface IntegrationWebPushReadiness {
  ready: boolean
  status: string
  public_key?: string
  connection_key?: string
  reason?: string
}

export interface IntegrationInvocation {
  id: string
  connector_key: string
  provider_key?: string
  connection_key?: string
  operation: string
  status: string
  duration_ms?: number
  request_ref?: string
  response_ref?: string
  error?: string
  metadata?: Record<string, unknown>
  created_at?: string
  updated_at?: string
}

export interface IntegrationEvent {
  id: string
  connector_key?: string
  connection_key?: string
  provider: string
  event_type: string
  external_id: string
  status: string
  payload?: Record<string, unknown>
  error?: string
  attempt_count: number
  next_retry_at?: string
  last_attempt_at?: string
  received_at?: string
  updated_at?: string
}
