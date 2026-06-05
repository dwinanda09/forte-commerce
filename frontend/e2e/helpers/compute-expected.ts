import { APIRequestContext } from '@playwright/test'

const BACKEND = 'http://localhost:8080'

interface Condition {
  type: string
  sku?: string
  min_qty?: number
  qty?: number
  amount?: number
  count?: number
}

interface Action {
  type: string
  sku?: string
  trigger_sku?: string
  buy_n?: number
  pay_m?: number
  pct?: number
  amount?: number
}

interface Campaign {
  id: string
  name: string
  is_active: boolean
  conditions: Condition[]
  actions: Action[]
}

export interface ExpectedCheckout {
  subtotal: number
  total_discount: number
  total: number
  campaignName: string
}

export async function fetchPriceMap(
  request: APIRequestContext,
  token: string
): Promise<Record<string, number>> {
  const res = await request.get(`${BACKEND}/api/v1/products`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  const body = await res.json()
  const products = body.data as Array<{ sku: string; price: number }>
  return Object.fromEntries(products.map((p) => [p.sku, p.price]))
}

export async function fetchActiveCampaigns(
  request: APIRequestContext,
  token: string
): Promise<Campaign[]> {
  const res = await request.get(`${BACKEND}/api/v1/campaigns`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  const body = await res.json()
  const raw = body.data as Array<{
    id: string
    name: string
    is_active: boolean
    conditions: Condition[]
    actions: Action[]
  }>
  return raw.filter((c) => c.is_active)
}

function requireCampaign(
  campaigns: Campaign[],
  predicate: (c: Campaign) => boolean,
  scenario: string
): Campaign {
  const found = campaigns.find(predicate)
  if (!found) {
    throw new Error(
      `[E2E setup] No active campaign found for scenario "${scenario}". ` +
        `Campaign config in DB has changed — either restore the campaign or update the test scenario.`
    )
  }
  return found
}

function round2(n: number): number {
  return Math.round(n * 100) / 100
}

export function computeTC1FreeItem(
  priceMap: Record<string, number>,
  campaigns: Campaign[]
): ExpectedCheckout {
  const triggerSKU = '43N23P'

  const campaign = requireCampaign(
    campaigns,
    (c) => c.actions.some((a) => a.type === 'free_item' && a.trigger_sku === triggerSKU),
    `free_item triggered by ${triggerSKU}`
  )

  const action = campaign.actions.find(
    (a) => a.type === 'free_item' && a.trigger_sku === triggerSKU
  )!

  const freeSKU = action.sku!
  const subtotal = round2(priceMap[triggerSKU] + priceMap[freeSKU])
  const discount = round2(priceMap[freeSKU])
  const total = round2(subtotal - discount)

  return { subtotal, total_discount: discount, total, campaignName: campaign.name }
}

export function computeTC2BuyNGetM(
  priceMap: Record<string, number>,
  campaigns: Campaign[]
): ExpectedCheckout {
  const sku = '120P90'
  const qty = 3

  const campaign = requireCampaign(
    campaigns,
    (c) => c.actions.some((a) => a.type === 'buy_n_get_m' && a.sku === sku),
    `buy_n_get_m for ${sku}`
  )

  const action = campaign.actions.find((a) => a.type === 'buy_n_get_m' && a.sku === sku)!
  const buyN = action.buy_n!
  const payM = action.pay_m!

  const subtotal = round2(qty * priceMap[sku])
  const freeItems = Math.floor(qty / buyN) * (buyN - payM)
  const discount = round2(freeItems * priceMap[sku])
  const total = round2(subtotal - discount)

  return { subtotal, total_discount: discount, total, campaignName: campaign.name }
}

export function computeTC3PctDiscountOnSKU(
  priceMap: Record<string, number>,
  campaigns: Campaign[]
): ExpectedCheckout {
  const sku = 'A304SD'
  const qty = 3

  const campaign = requireCampaign(
    campaigns,
    (c) => c.actions.some((a) => a.type === 'pct_discount_on_sku' && a.sku === sku),
    `pct_discount_on_sku for ${sku}`
  )

  const action = campaign.actions.find((a) => a.type === 'pct_discount_on_sku' && a.sku === sku)!
  const pct = action.pct!

  const subtotal = round2(qty * priceMap[sku])
  const discount = round2(subtotal * (pct / 100))
  const total = round2(subtotal - discount)

  return { subtotal, total_discount: discount, total, campaignName: campaign.name }
}
