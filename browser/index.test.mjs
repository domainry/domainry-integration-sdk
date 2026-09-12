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
  await client.listConnectionAccounts()
  await client.getConnectionAccount('account/1')
  await client.registerConnectionAccount('account/1', { scope: 'personal' })
  await client.testConnectionAccount('account/1')
  await client.revokeConnectionAccount('account/1', 'revision')
  await client.authorizeConnectionAccountRead('account/1', {operation: 'calendar_list', contract_sha256: 'a'.repeat(64)})
  await client.readConnectionAccount('account/1', {request_id: 'read-1', operation: 'calendar_list', contract_sha256: 'a'.repeat(64), payload: {limit: 2}})
  await client.listOAuthApplications()
  await client.upsertOAuthApplication('app/1', {})
  await client.listOAuthAuthorizationOptions()
  await client.startOAuthAuthorization({ application_key: 'app/1', scope: 'personal' })
  await client.getOAuthAuthorization('session/1')
  await client.completeOAuthAuthorization({ state: 'state', code: 'code' })
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
    'GET /integration/catalog',
    'GET /integration/connectors',
    'GET /integration/connections',
    'GET /integration/connections/connection%2F1',
    'POST /integration/connections/connection%2F1/validate',
    'PUT /integration/connections/connection%2F1',
    'DELETE /integration/connections/connection%2F1',
    'POST /integration/connections/connection%2F1/disable',
    'POST /integration/connections/connection%2F1/test-operation',
    'GET /integration/connection-accounts',
    'GET /integration/connection-accounts/account%2F1',
    'POST /integration/connections/account%2F1/account',
    'POST /integration/connection-accounts/account%2F1/test',
    'POST /integration/connection-accounts/account%2F1/revoke',
    'POST /integration/connection-accounts/account%2F1/read-access',
    'POST /integration/connection-accounts/account%2F1/read',
    'GET /integration/oauth-applications',
    'PUT /integration/oauth-applications/app%2F1',
    'GET /integration/oauth-authorizations/options',
    'POST /integration/oauth-authorizations',
    'GET /integration/oauth-authorizations/session%2F1',
    'POST /integration/oauth-authorizations/callback',
    'GET /integration/secrets',
    'PUT /integration/secrets/secret%2F1',
    'POST /integration/secrets/secret%2F1/disable',
    'POST /integration/secrets/secret%2F1/rotate',
    'POST /integration/secrets/secret%2F1/expire',
    'POST /integration/secrets/secret%2F1/revoke',
    'GET /integration/api-keys',
    'POST /integration/api-keys',
    'POST /integration/api-keys/api%2F1/disable',
    'POST /integration/api-keys/api%2F1/rotate',
    'GET /integration/external-identities',
    'PUT /integration/external-identities/identity%2F1',
    'POST /integration/external-identities/identity%2F1/disable',
    'POST /integration/external-identities/resolve',
    'GET /integration/webhook-subscriptions',
    'PUT /integration/webhook-subscriptions/webhook%2F1',
    'DELETE /integration/webhook-subscriptions/webhook%2F1',
    'POST /integration/webhook-subscriptions/webhook%2F1/disable',
    'GET /integration/web-push/readiness',
    'GET /integration/web-push/subscriptions',
    'PUT /integration/web-push/subscriptions/sub%2F1',
    'POST /integration/web-push/subscriptions/sub%2F1/revoke',
    'POST /integration/web-push/subscriptions/cleanup-expired',
    'GET /integration/invocations',
    'GET /integration/invocations/invocation%2F1',
    'GET /integration/events',
    'GET /integration/events/event%2F1',
    'POST /integration/events/event%2F1/replay',
  ])
  assert.equal(calls.length, 50)
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


test('account and OAuth commands keep fixed payloads and forward cancellation', async () => {
  const calls = []
  const controller = new AbortController()
  const client = new IntegrationClient({ request: async (path, options) => { calls.push({ path, options }); return {} } })
  const options = { method: 'DELETE', body: { owner_user_id: 'forged' }, signal: controller.signal }
  await client.testConnectionAccount('a', options)
  await client.revokeConnectionAccount('a', 'revision', options)
  await client.startOAuthAuthorization({ application_key: 'app', scope: 'personal' }, options)
  await client.completeOAuthAuthorization({ state: 'private', code: 'private' }, options)
  for (const { options } of calls) {
    assert.equal(options.method, 'POST')
    assert.equal(options.signal, controller.signal)
    assert.equal(options.body.owner_user_id, undefined)
  }
  assert.deepEqual(calls[0].options.body, {})
  assert.deepEqual(calls[1].options.body, { expected_updated_at: 'revision' })
})
