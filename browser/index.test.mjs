import assert from 'node:assert/strict'
import test from 'node:test'

import { IntegrationClient } from './dist/index.js'

test('calls only Integration owner routes', async () => {
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

  await client.listInvocations({ connector_key: 'crm' })
  await client.listEvents({ status: 'failed' })
  await client.replayEvent('event/1')
  await client.webPushReadiness()
  await client.upsertWebPushSubscription(
    'sub/1',
    { endpoint: 'https://push.example/sub', p256dh: 'p256-secret', auth: 'auth-secret' },
    { idempotencyKey: 'upsert-1' },
  )

  assert.deepEqual(calls.map((call) => call.path), [
    '/tenant-admin/integrations/invocations?connector_key=crm',
    '/tenant-admin/integrations/events?status=failed',
    '/tenant-admin/integrations/events/event%2F1/replay',
    '/business/notifications/web-push/readiness',
    '/business/notifications/web-push/subscriptions/sub%2F1',
  ])
  assert.equal(calls.at(-1).options.headers['Idempotency-Key'], 'upsert-1')
  assert.equal(calls.some((call) => call.path.startsWith('/operations/integrations')), false)
})
