import type {
 IntegrationConnectionAccountReadOperation, IntegrationConnectionAccountReadRequest, IntegrationConnectionAccountReadAccess, IntegrationConnectionAccountReadResult,
 IntegrationConnectionAccount, IntegrationOAuthApplication, IntegrationOAuthApplicationInput,
 IntegrationOAuthAuthorizationOption, IntegrationOAuthAuthorizationInput, IntegrationOAuthAuthorizationSession, IntegrationOAuthAuthorizationCallback,
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
    return this.#dependencies.request<T>('/integration/catalog', { signal })
  }

  connectors<T = IntegrationCatalog>(signal?: AbortSignal) {
    return this.#dependencies.request<T>('/integration/connectors', { signal })
  }

  listConnections<T = { connections: IntegrationConnection[]; count: number }>(query: IntegrationQuery = {}, signal?: AbortSignal) {
    return this.#dependencies.request<T>(
      `/integration/connections${queryString(query)}`,
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


  listConnectionAccounts(signal?: AbortSignal) {
    return this.#dependencies.request<{accounts: IntegrationConnectionAccount[]; count: number}>('/integration/connection-accounts', {signal})
  }
  getConnectionAccount(key: string, signal?: AbortSignal) {required(key, 'key'); return this.#dependencies.request<IntegrationConnectionAccount>(integrationResourcePath('connection-accounts', key), {signal})}
  registerConnectionAccount(key: string, body: {scope: 'personal' | 'workspace'; owner_user_id?: string}, options: IntegrationRequestOptions = {}) {return this.#command<IntegrationConnectionAccount>('connections',key,'account','POST',body,options)}
  testConnectionAccount(key: string, options: IntegrationRequestOptions = {}) {return this.#command<{account: IntegrationConnectionAccount; operation: string; connected: boolean}>('connection-accounts',key,'test','POST',{},options)}
  revokeConnectionAccount(key: string, expectedUpdatedAt: string, options: IntegrationRequestOptions = {}) {return this.#command<IntegrationConnectionAccount>('connection-accounts',key,'revoke','POST',{expected_updated_at: expectedUpdatedAt},options)}
  authorizeConnectionAccountRead(key: string, body: IntegrationConnectionAccountReadOperation, options: IntegrationRequestOptions = {}) {return this.#command<IntegrationConnectionAccountReadAccess>('connection-accounts',key,'read-access','POST',body,options)}
  readConnectionAccount(key: string, body: IntegrationConnectionAccountReadRequest, options: IntegrationRequestOptions = {}) {return this.#command<IntegrationConnectionAccountReadResult>('connection-accounts',key,'read','POST',body,options)}
  listOAuthApplications(signal?: AbortSignal) {return this.#dependencies.request<{applications: IntegrationOAuthApplication[]}>('/integration/oauth-applications',{signal})}
  upsertOAuthApplication(key: string, body: IntegrationOAuthApplicationInput, options: IntegrationRequestOptions = {}) {return this.#resourceCommand<IntegrationOAuthApplication>('oauth-applications',key,'PUT',body,options)}
  listOAuthAuthorizationOptions(signal?: AbortSignal) {return this.#dependencies.request<{options: IntegrationOAuthAuthorizationOption[]}>('/integration/oauth-authorizations/options',{signal})}
  startOAuthAuthorization(body: IntegrationOAuthAuthorizationInput, options: IntegrationRequestOptions = {}) {return this.#dependencies.request<IntegrationOAuthAuthorizationSession>('/integration/oauth-authorizations',{...options,method:'POST',body})}
  getOAuthAuthorization(id: string, signal?: AbortSignal) {required(id,'id');return this.#dependencies.request<IntegrationOAuthAuthorizationSession>(integrationResourcePath('oauth-authorizations',id),{signal})}
  completeOAuthAuthorization(body: IntegrationOAuthAuthorizationCallback, options: IntegrationRequestOptions = {}) {return this.#dependencies.request<IntegrationOAuthAuthorizationSession>('/integration/oauth-authorizations/callback',{...options,method:'POST',body})}

  listSecrets<T = { items: IntegrationSecret[]; count: number }>(query: IntegrationQuery = {}, signal?: AbortSignal) {
    return this.#dependencies.request<T>(`/integration/secrets${queryString(query)}`, { signal })
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
    return this.#dependencies.request<T>(`/integration/api-keys${queryString(query)}`, { signal })
  }

  createAPIKey<T = IntegrationAPIKeyCredential>(body: IntegrationAPIKeyInput, options: IntegrationRequestOptions = {}) {
    return this.#dependencies.request<T>('/integration/api-keys', fixedOptions(options, 'POST', body))
  }

  disableAPIKey<T = IntegrationAPIKey>(apiKey: string, body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('api-keys', apiKey, 'disable', 'POST', body, options)
  }

  rotateAPIKey<T = IntegrationAPIKeyCredential>(apiKey: string, body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#command<T>('api-keys', apiKey, 'rotate', 'POST', body, options)
  }

  listExternalIdentities<T = { identities: IntegrationExternalIdentity[]; count: number }>(query: IntegrationQuery = {}, signal?: AbortSignal) {
    return this.#dependencies.request<T>(
      `/integration/external-identities${queryString(query)}`,
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
      '/integration/external-identities/resolve',
      fixedOptions(options, 'POST', body),
    )
  }

  listWebhookSubscriptions<T = { subscriptions: IntegrationWebhookSubscription[]; count: number }>(query: IntegrationQuery = {}, signal?: AbortSignal) {
    return this.#dependencies.request<T>(
      `/integration/webhook-subscriptions${queryString(query)}`,
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
      `/integration/invocations${queryString(query)}`,
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
      `/integration/events${queryString(query)}`,
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
      `/integration/events/${encodeURIComponent(eventID)}/replay`,
      { method: 'POST' },
    )
  }

  webPushReadiness() {
    return this.#dependencies.request<IntegrationWebPushReadiness>(
      '/integration/web-push/readiness',
    )
  }

  webPushSubscriptions() {
    return this.#dependencies.request<{ subscriptions: IntegrationWebPushSubscription[]; count: number }>(
      '/integration/web-push/subscriptions',
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
      `/integration/web-push/subscriptions/${encodeURIComponent(subscriptionID)}`,
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
      `/integration/web-push/subscriptions/${encodeURIComponent(subscriptionID)}/revoke`,
      {
        method: 'POST',
        requestId: options.requestId,
        headers: { 'Idempotency-Key': options.idempotencyKey },
      },
    )
  }

  cleanupExpiredWebPushSubscriptions<T = unknown>(body?: unknown, options: IntegrationRequestOptions = {}) {
    return this.#dependencies.request<T>(
      '/integration/web-push/subscriptions/cleanup-expired',
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
  return `/integration/${resource}/${encodeURIComponent(key)}`
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
