export interface IntegrationRequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  body?: unknown
  headers?: Record<string, string>
  requestId?: string
  signal?: AbortSignal
}

export type IntegrationRequest = <T>(path: string, options?: IntegrationRequestOptions) => Promise<T>

export interface IntegrationClientDependencies {
  request: IntegrationRequest
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
