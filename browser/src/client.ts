import type {
  IntegrationClientDependencies,
  IntegrationEvent,
  IntegrationInvocation,
  IntegrationWebPushReadiness,
  IntegrationWebPushSubscription,
  IntegrationWebPushSubscriptionUpsertRequest,
} from './types.js'

export class IntegrationClient {
  readonly #dependencies: IntegrationClientDependencies

  constructor(dependencies: IntegrationClientDependencies) {
    this.#dependencies = dependencies
  }

  catalog(signal?: AbortSignal) {
    return this.#dependencies.request('/tenant-admin/integrations/connectors', { signal })
  }

  listInvocations(query: Record<string, string> = {}, signal?: AbortSignal) {
    return this.#dependencies.request<{ invocations: IntegrationInvocation[]; count: number }>(
      `/tenant-admin/integrations/invocations${queryString(query)}`,
      { signal },
    )
  }

  listEvents(query: Record<string, string> = {}, signal?: AbortSignal) {
    return this.#dependencies.request<{ events: IntegrationEvent[]; count: number }>(
      `/tenant-admin/integrations/events${queryString(query)}`,
      { signal },
    )
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
}

function queryString(query: Record<string, string>): string {
  const values = new URLSearchParams()
  for (const [key, value] of Object.entries(query)) {
    if (value.trim()) values.set(key, value.trim())
  }
  return values.size ? `?${values.toString()}` : ''
}

function required(value: string, name: string): void {
  if (!value.trim()) throw new Error(`${name} is required`)
}
