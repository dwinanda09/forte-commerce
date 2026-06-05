import { test, expect, APIRequestContext } from '@playwright/test'
import { execSync } from 'child_process'
import {
  fetchPriceMap,
  fetchActiveCampaigns,
  computeTC1FreeItem,
  computeTC2BuyNGetM,
  computeTC3PctDiscountOnSKU,
  type ExpectedCheckout,
} from './helpers/compute-expected'

// Serial mode: inventory state must be predictable across tests.
test.describe.configure({ mode: 'serial' })

const BACKEND = 'http://localhost:8080'
const COMPOSE_DIR = '/Users/dwinanda/saas/fortecommerce'

function resetInventory() {
  execSync(
    `docker compose exec -T postgres psql -U forte -d forte_commerce -c "` +
      `UPDATE products SET inventory_qty=5, reserved_qty=0 WHERE sku='43N23P';` +
      `UPDATE products SET inventory_qty=2, reserved_qty=0 WHERE sku='234234';` +
      `UPDATE products SET inventory_qty=10, reserved_qty=0 WHERE sku='120P90';` +
      `UPDATE products SET inventory_qty=10, reserved_qty=0 WHERE sku='A304SD';` +
      `DELETE FROM orders;` +
      `DELETE FROM checkout_sessions;"`,
    { cwd: COMPOSE_DIR }
  )
}

async function login(request: APIRequestContext): Promise<string> {
  const res = await request.post(`${BACKEND}/api/v1/auth/login`, {
    data: { username: 'demo1', password: 'demo1' },
  })
  expect(res.status()).toBe(200)
  const body = await res.json()
  return body.data.token as string
}

async function submitCheckout(
  request: APIRequestContext,
  token: string,
  items: string[]
): Promise<string> {
  const res = await request.post(`${BACKEND}/api/v1/checkout`, {
    headers: { Authorization: `Bearer ${token}` },
    data: { items },
  })
  expect(res.status()).toBe(202)
  const body = await res.json()
  return body.data.checkout_id as string
}

async function pollCheckout(
  request: APIRequestContext,
  token: string,
  id: string,
  timeoutMs = 15000
): Promise<Record<string, unknown>> {
  const deadline = Date.now() + timeoutMs
  while (Date.now() < deadline) {
    const res = await request.get(`${BACKEND}/api/v1/checkout/${id}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    expect(res.status()).toBe(200)
    const body = await res.json()
    const session = body.data as Record<string, unknown>
    if (session.status !== 'pending') return session
    await new Promise((r) => setTimeout(r, 500))
  }
  throw new Error(`Checkout ${id} did not reach terminal state within ${timeoutMs}ms`)
}

test.describe('Checkout Promotions (PDF Core Cases)', () => {
  let tc1: ExpectedCheckout
  let tc2: ExpectedCheckout
  let tc3: ExpectedCheckout

  // Fetch live product prices and active campaigns once, compute expected values from DB state.
  // If a campaign no longer exists or no longer applies, setup throws — tests fail loudly.
  test.beforeAll(async ({ request }) => {
    const token = await login(request)
    const [priceMap, campaigns] = await Promise.all([
      fetchPriceMap(request, token),
      fetchActiveCampaigns(request, token),
    ])
    tc1 = computeTC1FreeItem(priceMap, campaigns)
    tc2 = computeTC2BuyNGetM(priceMap, campaigns)
    tc3 = computeTC3PctDiscountOnSKU(priceMap, campaigns)
  })

  test.beforeEach(() => {
    resetInventory()
  })

  test('TC1 — MacBook Pro + Raspberry Pi: free item promotion', async ({ request }) => {
    const token = await login(request)
    const id = await submitCheckout(request, token, ['43N23P', '234234'])
    const session = await pollCheckout(request, token, id)

    expect(session.status).toBe('completed')

    const result = session.result as Record<string, unknown>
    expect(parseFloat(result.subtotal as string)).toBeCloseTo(tc1.subtotal, 2)
    expect(parseFloat(result.total_discount as string)).toBeCloseTo(tc1.total_discount, 2)
    expect(parseFloat(result.total as string)).toBeCloseTo(tc1.total, 2)

    const promos = result.promotions_applied as Array<Record<string, unknown>>
    expect(promos.length).toBeGreaterThan(0)
    expect(promos[0].name).toBe(tc1.campaignName)
    expect(parseFloat(promos[0].discount as string)).toBeCloseTo(tc1.total_discount, 2)
  })

  test('TC2 — 3× Google Home: bundle 3-for-2 promotion', async ({ request }) => {
    const token = await login(request)
    const id = await submitCheckout(request, token, ['120P90', '120P90', '120P90'])
    const session = await pollCheckout(request, token, id)

    expect(session.status).toBe('completed')

    const result = session.result as Record<string, unknown>
    expect(parseFloat(result.subtotal as string)).toBeCloseTo(tc2.subtotal, 2)
    expect(parseFloat(result.total_discount as string)).toBeCloseTo(tc2.total_discount, 2)
    expect(parseFloat(result.total as string)).toBeCloseTo(tc2.total, 2)

    const promos = result.promotions_applied as Array<Record<string, unknown>>
    expect(promos.length).toBeGreaterThan(0)
    expect(promos[0].name).toBe(tc2.campaignName)
  })

  test('TC3 — 3× Alexa Speaker: 10% quantity discount', async ({ request }) => {
    const token = await login(request)
    const id = await submitCheckout(request, token, ['A304SD', 'A304SD', 'A304SD'])
    const session = await pollCheckout(request, token, id)

    expect(session.status).toBe('completed')

    const result = session.result as Record<string, unknown>
    expect(parseFloat(result.subtotal as string)).toBeCloseTo(tc3.subtotal, 2)
    expect(parseFloat(result.total_discount as string)).toBeCloseTo(tc3.total_discount, 2)
    expect(parseFloat(result.total as string)).toBeCloseTo(tc3.total, 2)

    const promos = result.promotions_applied as Array<Record<string, unknown>>
    expect(promos.length).toBeGreaterThan(0)
    expect(promos[0].name).toBe(tc3.campaignName)
  })

  test('TC4 — Concurrency: only 1 of 5 concurrent checkouts succeeds for stock=5 MacBook', async ({
    request,
  }) => {
    const token = await login(request)

    // 5 concurrent checkouts each requesting all 5 MacBook Pro in stock
    const ids = await Promise.all(
      Array.from({ length: 5 }, () =>
        submitCheckout(request, token, ['43N23P', '43N23P', '43N23P', '43N23P', '43N23P'])
      )
    )

    const sessions = await Promise.all(ids.map((id) => pollCheckout(request, token, id)))

    const completed = sessions.filter((s) => s.status === 'completed')
    const failed = sessions.filter((s) => s.status === 'failed')

    expect(completed).toHaveLength(1)
    expect(failed).toHaveLength(4)

    for (const s of failed) {
      expect((s.error_message as string).toLowerCase()).toMatch(/insufficient/)
    }
  })
})
