import type {
  IntegrationAPIKey,
  IntegrationAPIKeyCredential,
  IntegrationAPIKeyInput,
  IntegrationCatalog,
  IntegrationClientDependencies,
  IntegrationConnection,
  IntegrationConnectionInput,
  IntegrationConnectionTestRequest,
  IntegrationConnectorOperationTestResult,
  IntegrationEvent,
  IntegrationExternalIdentity,
  IntegrationExternalIdentityInput,
  IntegrationInvocation,
  IntegrationRequestOptions,
  IntegrationSecret,
  IntegrationSecretInput,
  IntegrationWebhookSubscription,
  IntegrationWebhookSubscriptionInput,
  IntegrationWebPushReadiness,
  IntegrationWebPushSubscription,
  IntegrationWebPushSubscriptionUpsertRequest,
} from './types.js'

export class IntegrationClient {
  readonly #dependencies: IntegrationClientDependencies

  constructor(dependencies: IntegrationClientDependencies) {
    this.#dependencies = dependencies
  }

  catalog<T = IntegrationCatalog>(signal?: AbortSignal) {
    return this.#dependencies.request<T>('/tenant-admin/integrations/catalog', { signal })
  }

  connectors<T = IntegrationCatalog>(signal?: AbortSignal) {
    return this.#dependencies.request<T>('/tenant-admin/integrations/connectors', { signal })
  }

  listConnections<T = { connections: IntegrationConnection[]; count: number }>(query: IntegrationQuery = {}, signal?: AbortSignal) {
    return this.#dependencies.request<T>(
      `/tenant-admin/integrations/connections${queryString(query)}`,
      { signal },
    )
  }

  getConnection<T = IntegrationConnection>(connectionKey: string, signal?: AbortSignal) {
    required(connectionKey, 'connectionKey')
    return this.#dependencies.request<T>(integrationResourcePath('connections', connectionKey), { signal })
  }

  getConnectionWithResponse<T = IntegrationConnection>(connectionKey: string, signal?: AbortSignal) {
    required(connectionKey, 'connectionKey')
    if (!this.#dependencies.requestWithResponse) {
      throw new Error('IntegrationClient requestWithResponse dependency is required')
    }
    return this.#dependencies.requestWithResponse<T>(integrationResourcePath('connections', connectionKey), { signal })
  }

  validateConnection<T = { valid: boolean; connection_key: string; normalized: IntegrationConnectionInput }>(connectionKey: string, body: IntegrationConnectionInput, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('connections', connectionKey, 'validate', 'POST', body, options)
  }

  upsertConnection<T = IntegrationConnection>(connectionKey: string, body: IntegrationConnectionInput, options: IntegrationRequestOptions = {}) {
    return this.#resourceCommand<T>('connections', connectionKey, 'PUT', body, options)
  }

  deleteConnection<T = void>(connectionKey: string, options: IntegrationRequestOptions = {}) {
    return this.#resourceCommand<T>('connections', connectionKey, 'DELETE', undefined, options)
  }

  disableConnection<T = IntegrationConnection>(connectionKey: string, body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('connections', connectionKey, 'disable', 'POST', body, options)
  }

  testOperation<T = IntegrationConnectorOperationTestResult>(connectionKey: string, body: IntegrationConnectionTestRequest, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('connections', connectionKey, 'test-operation', 'POST', body, options)
  }

  listSecrets<T = { items: IntegrationSecret[]; count: number }>(query: IntegrationQuery = {}, signal?: AbortSignal) {
    return this.#dependencies.request<T>(`/tenant-admin/integrations/secrets${queryString(query)}`, { signal })
  }

  upsertSecret<T = IntegrationSecret>(secretKey: string, body: IntegrationSecretInput, options: IntegrationRequestOptions = {}) {
    return this.#resourceCommand<T>('secrets', secretKey, 'PUT', body, options)
  }

  disableSecret<T = IntegrationSecret>(secretKey: string, body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('secrets', secretKey, 'disable', 'POST', body, options)
  }

  rotateSecret<T = IntegrationSecret>(secretKey: string, body?: IntegrationSecretInput, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('secrets', secretKey, 'rotate', 'POST', body, options)
  }

  expireSecret<T = IntegrationSecret>(secretKey: string, body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('secrets', secretKey, 'expire', 'POST', body, options)
  }

  revokeSecret<T = IntegrationSecret>(secretKey: string, body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('secrets', secretKey, 'revoke', 'POST', body, options)
  }

  listAPIKeys<T = { api_keys: IntegrationAPIKey[]; count: number }>(query: IntegrationQuery = {}, signal?: AbortSignal) {
    return this.#dependencies.request<T>(`/tenant-admin/integrations/api-keys${queryString(query)}`, { signal })
  }

  createAPIKey<T = IntegrationAPIKeyCredential>(body: IntegrationAPIKeyInput, options: IntegrationRequestOptions = {}) {
    return this.#dependencies.request<T>('/tenant-admin/integrations/api-keys', fixedOptions(options, 'POST', body))
  }

  disableAPIKey<T = IntegrationAPIKey>(apiKey: string, body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('api-keys', apiKey, 'disable', 'POST', body, options)
  }

  rotateAPIKey<T = IntegrationAPIKeyCredential>(apiKey: string, body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('api-keys', apiKey, 'rotate', 'POST', body, options)
  }

  listExternalIdentities<T = { identities: IntegrationExternalIdentity[]; count: number }>(query: IntegrationQuery = {}, signal?: AbortSignal) {
    return this.#dependencies.request<T>(
      `/tenant-admin/integrations/external-identities${queryString(query)}`,
      { signal },
    )
  }

  upsertExternalIdentity<T = IntegrationExternalIdentity>(identityKey: string, body: IntegrationExternalIdentityInput, options: IntegrationRequestOptions = {}) {
    return this.#resourceCommand<T>('external-identities', identityKey, 'PUT', body, options)
  }

  disableExternalIdentity<T = IntegrationExternalIdentity>(identityKey: string, body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('external-identities', identityKey, 'disable', 'POST', body, options)
  }

  resolveExternalIdentity<T = IntegrationExternalIdentity>(body: { provider: string; external_subject: string }, options: IntegrationRequestOptions = {}) {
    return this.#dependencies.request<T>(
      '/tenant-admin/integrations/external-identities/resolve',
      fixedOptions(options, 'POST', body),
    )
  }

  listWebhookSubscriptions<T = { subscriptions: IntegrationWebhookSubscription[]; count: number }>(query: IntegrationQuery = {}, signal?: AbortSignal) {
    return this.#dependencies.request<T>(
      `/tenant-admin/integrations/webhook-subscriptions${queryString(query)}`,
      { signal },
    )
  }

  upsertWebhookSubscription<T = IntegrationWebhookSubscription>(subscriptionKey: string, body: IntegrationWebhookSubscriptionInput, options: IntegrationRequestOptions = {}) {
    return this.#resourceCommand<T>('webhook-subscriptions', subscriptionKey, 'PUT', body, options)
  }

  deleteWebhookSubscription<T = void>(subscriptionKey: string, options: IntegrationRequestOptions = {}) {
    return this.#resourceCommand<T>('webhook-subscriptions', subscriptionKey, 'DELETE', undefined, options)
  }

  disableWebhookSubscription<T = IntegrationWebhookSubscription>(subscriptionKey: string, body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('webhook-subscriptions', subscriptionKey, 'disable', 'POST', body, options)
  }

  listInvocations(query: Record<string, string> = {}, signal?: AbortSignal) {
    return this.#dependencies.request<{ invocations: IntegrationInvocation[]; count: number }>(
      `/tenant-admin/integrations/invocations${queryString(query)}`,
      { signal },
    )
  }

  getInvocation(invocationID: string, signal?: AbortSignal) {
    required(invocationID, 'invocationID')
    return this.#dependencies.request<IntegrationInvocation>(
      integrationResourcePath('invocations', invocationID),
      { signal },
    )
  }

  listEvents(query: Record<string, string> = {}, signal?: AbortSignal) {
    return this.#dependencies.request<{ events: IntegrationEvent[]; count: number }>(
      `/tenant-admin/integrations/events${queryString(query)}`,
      { signal },
    )
  }

  getEvent(eventID: string, signal?: AbortSignal) {
    required(eventID, 'eventID')
    return this.#dependencies.request<IntegrationEvent>(integrationResourcePath('events', eventID), { signal })
  }

  replayEvent(eventID: string) {
    required(eventID, 'eventID')
    return this.#dependencies.request<IntegrationEvent>(
      `/tenant-admin/integrations/events/${encodeURIComponent(eventID)}/replay`,
      { method: 'POST' },
    )
  }

  webPushReadiness() {
    return this.#dependencies.request<IntegrationWebPushReadiness>(
      '/business/notifications/web-push/readiness',
    )
  }

  webPushSubscriptions() {
    return this.#dependencies.request<{ subscriptions: IntegrationWebPushSubscription[]; count: number }>(
      '/business/notifications/web-push/subscriptions',
    )
  }

  upsertWebPushSubscription(
    subscriptionID: string,
    request: IntegrationWebPushSubscriptionUpsertRequest,
    options: { idempotencyKey: string; requestId?: string },
  ) {
    required(subscriptionID, 'subscriptionID')
    required(options.idempotencyKey, 'Idempotency-Key')
    return this.#dependencies.request<IntegrationWebPushSubscription>(
      `/business/notifications/web-push/subscriptions/${encodeURIComponent(subscriptionID)}`,
      {
        method: 'PUT',
        requestId: options.requestId,
        headers: { 'Idempotency-Key': options.idempotencyKey },
        body: request,
      },
    )
  }

  revokeWebPushSubscription(
    subscriptionID: string,
    options: { idempotencyKey: string; requestId?: string },
  ) {
    required(subscriptionID, 'subscriptionID')
    required(options.idempotencyKey, 'Idempotency-Key')
    return this.#dependencies.request<IntegrationWebPushSubscription>(
      `/business/notifications/web-push/subscriptions/${encodeURIComponent(subscriptionID)}/revoke`,
      {
        method: 'POST',
        requestId: options.requestId,
        headers: { 'Idempotency-Key': options.idempotencyKey },
      },
    )
  }

  cleanupExpiredWebPushSubscriptions<T = unknown>(body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#dependencies.request<T>(
      '/integrations/web-push/subscriptions/cleanup-expired',
      fixedOptions(options, 'POST', body),
    )
  }

  #resourceCommand<T>(
    resource: string,
    key: string,
    method: 'PUT' | 'DELETE',
    body: unknown,
    options: IntegrationRequestOptions,
  ) {
    required(key, `${resource} key`)
    return this.#dependencies.request<T>(
      integrationResourcePath(resource, key),
      fixedOptions(options, method, body),
    )
  }

  #command<T>(
    resource: string,
    key: string,
    command: string,
    method: 'POST',
    body: unknown,
    options: IntegrationRequestOptions,
  ) {
    required(key, `${resource} key`)
    return this.#dependencies.request<T>(
      `${integrationResourcePath(resource, key)}/${command}`,
      fixedOptions(options, method, body),
    )
  }
}

export type IntegrationQuery = Record<string, string | number | boolean | undefined>

function queryString(query: IntegrationQuery): string {
  const values = new URLSearchParams()
  for (const [key, value] of Object.entries(query)) {
    const normalized = value === undefined ? '' : String(value).trim()
    if (normalized) values.set(key, normalized)
  }
  return values.size ? `?${values.toString()}` : ''
}

function integrationResourcePath(resource: string, key: string): string {
  return `/tenant-admin/integrations/${resource}/${encodeURIComponent(key)}`
}

function fixedOptions(
  options: IntegrationRequestOptions,
  method: 'POST' | 'PUT' | 'DELETE',
  body: unknown,
): IntegrationRequestOptions {
  const result: IntegrationRequestOptions = { ...options, method }
  if (body !== undefined) result.body = body
  return result
}

function required(value: string, name: string): void {
  if (!value.trim()) throw new Error(`${name} is required`)
}
