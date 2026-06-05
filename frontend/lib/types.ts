export interface Product {
  id: string
  sku: string
  name: string
  price: number
  inventory_qty: number
  reserved_qty: number
  available_qty: number
}

export interface CheckoutItem {
  sku: string
  name: string
  qty: number
  price: number
  total: number
}

export interface AppliedPromotion {
  name: string
  description: string
  discount: number
}

export interface CheckoutResult {
  items: CheckoutItem[]
  promotions_applied: AppliedPromotion[]
  subtotal: number
  total_discount: number
  total: number
}

export interface CheckoutSession {
  checkout_id: string
  status: 'pending' | 'completed' | 'expired' | 'failed'
  expires_at: string
  error_message?: string
  result?: CheckoutResult
}

export interface Order {
  id: string
  checkout_session_id: string
  status: 'pending' | 'paid' | 'cancelled'
  items: CheckoutItem[]
  promotions_applied: AppliedPromotion[]
  subtotal: number
  total_discount: number
  total: number
  created_at: string
  updated_at: string
}

export type ConditionType = 'cart_has_sku' | 'item_qty_gte' | 'cart_total_gte' | 'cart_item_count_gte'
export type ActionType = 'free_item' | 'buy_n_get_m' | 'pct_discount_on_sku' | 'pct_discount_on_cart' | 'fixed_discount'

export interface Condition {
  type: ConditionType
  sku?: string
  min_qty?: number
  qty?: number
  amount?: number
  count?: number
}

export interface Action {
  type: ActionType
  sku?: string
  trigger_sku?: string
  buy_n?: number
  pay_m?: number
  pct?: number
  amount?: number
}

export interface Campaign {
  id: string
  name: string
  description: string
  is_active: boolean
  priority: number
  conditions: Condition[]
  actions: Action[]
  created_at: string
  updated_at: string
}

export interface CampaignRequest {
  name: string
  description?: string
  is_active?: boolean
  priority?: number
  conditions?: Condition[]
  actions?: Action[]
}

export interface ApiResponse<T> {
  success: boolean
  data: T
  message?: string
  meta: { request_id: string; timestamp: string }
}
