import assert from 'node:assert/strict'
import test from 'node:test'

import { IntegrationClient } from './dist/index.js'

test('publishes one owner-client method for every browser Integration route', async () => {
  const calls = []
  const client = new IntegrationClient({
    request: async (path, options) => {
      calls.push({ path, options })
      if (path.includes('/invocations')) return { invocations: [], count: 0 }
      if (path.includes('/events')) return { events: [], count: 0 }
      if (path.endsWith('/readiness')) return { ready: true, status: 'verified' }
      return { subscriptions: [], count: 0 }
    },
  })

  await client.catalog()
  await client.connectors()
  await client.listConnections()
  await client.getConnection('connection/1')
  await client.validateConnection('connection/1', {})
  await client.upsertConnection('connection/1', {})
  await client.deleteConnection('connection/1')
  await client.disableConnection('connection/1')
  await client.testOperation('connection/1', {})
  await client.listSecrets()
  await client.upsertSecret('secret/1', {})
  await client.disableSecret('secret/1')
  await client.rotateSecret('secret/1')
  await client.expireSecret('secret/1')
  await client.revokeSecret('secret/1')
  await client.listAPIKeys()
  await client.createAPIKey({})
  await client.disableAPIKey('api/1')
  await client.rotateAPIKey('api/1')
  await client.listExternalIdentities()
  await client.upsertExternalIdentity('identity/1', {})
  await client.disableExternalIdentity('identity/1')
  await client.resolveExternalIdentity({})
  await client.listWebhookSubscriptions()
  await client.upsertWebhookSubscription('webhook/1', {})
  await client.deleteWebhookSubscription('webhook/1')
  await client.disableWebhookSubscription('webhook/1')
  await client.webPushReadiness()
  await client.webPushSubscriptions()
  await client.upsertWebPushSubscription(
    'sub/1',
    { endpoint: 'https://push.example/sub', p256dh: 'p256-secret', auth: 'auth-secret' },
    { idempotencyKey: 'upsert-1' },
  )
  await client.revokeWebPushSubscription('sub/1', { idempotencyKey: 'revoke-1' })
  await client.cleanupExpiredWebPushSubscriptions()
  await client.listInvocations()
  await client.getInvocation('invocation/1')
  await client.listEvents()
  await client.getEvent('event/1')
  await client.replayEvent('event/1')

  assert.deepEqual(calls.map(({ path, options }) => `${options?.method ?? 'GET'} ${path}`), [
    'GET /tenant-admin/integrations/catalog',
    'GET /tenant-admin/integrations/connectors',
    'GET /tenant-admin/integrations/connections',
    'GET /tenant-admin/integrations/connections/connection%2F1',
    'POST /tenant-admin/integrations/connections/connection%2F1/validate',
    'PUT /tenant-admin/integrations/connections/connection%2F1',
    'DELETE /tenant-admin/integrations/connections/connection%2F1',
    'POST /tenant-admin/integrations/connections/connection%2F1/disable',
    'POST /tenant-admin/integrations/connections/connection%2F1/test-operation',
    'GET /tenant-admin/integrations/secrets',
    'PUT /tenant-admin/integrations/secrets/secret%2F1',
    'POST /tenant-admin/integrations/secrets/secret%2F1/disable',
    'POST /tenant-admin/integrations/secrets/secret%2F1/rotate',
    'POST /tenant-admin/integrations/secrets/secret%2F1/expire',
    'POST /tenant-admin/integrations/secrets/secret%2F1/revoke',
    'GET /tenant-admin/integrations/api-keys',
    'POST /tenant-admin/integrations/api-keys',
    'POST /tenant-admin/integrations/api-keys/api%2F1/disable',
    'POST /tenant-admin/integrations/api-keys/api%2F1/rotate',
    'GET /tenant-admin/integrations/external-identities',
    'PUT /tenant-admin/integrations/external-identities/identity%2F1',
    'POST /tenant-admin/integrations/external-identities/identity%2F1/disable',
    'POST /tenant-admin/integrations/external-identities/resolve',
    'GET /tenant-admin/integrations/webhook-subscriptions',
    'PUT /tenant-admin/integrations/webhook-subscriptions/webhook%2F1',
    'DELETE /tenant-admin/integrations/webhook-subscriptions/webhook%2F1',
    'POST /tenant-admin/integrations/webhook-subscriptions/webhook%2F1/disable',
    'GET /business/notifications/web-push/readiness',
    'GET /business/notifications/web-push/subscriptions',
    'PUT /business/notifications/web-push/subscriptions/sub%2F1',
    'POST /business/notifications/web-push/subscriptions/sub%2F1/revoke',
    'POST /integrations/web-push/subscriptions/cleanup-expired',
    'GET /tenant-admin/integrations/invocations',
    'GET /tenant-admin/integrations/invocations/invocation%2F1',
    'GET /tenant-admin/integrations/events',
    'GET /tenant-admin/integrations/events/event%2F1',
    'POST /tenant-admin/integrations/events/event%2F1/replay',
  ])
  assert.equal(calls.length, 37)
  assert.equal(calls.some((call) => call.path.startsWith('/operations/integrations')), false)
})

test('keeps response metadata behind the Integration client boundary', async () => {
  const client = new IntegrationClient({
    request: async () => ({}),
    requestWithResponse: async (path) => ({
      data: { key: 'connection-1' },
      headers: { get: (name) => name === 'X-Resource-Hash' ? 'hash-1' : null },
      status: 200,
      path,
    }),
  })

  const response = await client.getConnectionWithResponse('connection/1')
  assert.equal(response.headers.get('X-Resource-Hash'), 'hash-1')
  assert.equal(response.data.key, 'connection-1')
})
