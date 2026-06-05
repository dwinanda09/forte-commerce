import type { ApiResponse, Product, CheckoutSession, Order, Campaign, CampaignRequest } from './types'

async function apiFetch<T>(path: string, options?: RequestInit): Promise<ApiResponse<T>> {
  const url = path.startsWith('http') ? path : path

  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    credentials: 'include',
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: 'Unknown error' }))
    throw new Error(error.message || `API error: ${response.status}`)
  }

  return response.json()
}

export const api = {
  login: (username: string, password: string) =>
    apiFetch<{ token: string }>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    }),

  getProducts: () => apiFetch<Product[]>('/api/v1/products'),

  submitCheckout: (items: string[]) =>
    apiFetch<{ checkout_id: string }>('/api/v1/checkout', {
      method: 'POST',
      body: JSON.stringify({ items }),
    }),

  getCheckout: (id: string) => apiFetch<CheckoutSession>(`/api/v1/checkout/${id}`),

  confirmCheckout: (id: string) =>
    apiFetch<Order>(`/api/v1/checkout/${id}/confirm`, {
      method: 'POST',
    }),

  payOrder: (id: string) =>
    apiFetch<Order>(`/api/v1/orders/${id}/pay`, {
      method: 'POST',
    }),

  cancelOrder: (id: string) =>
    apiFetch<Order>(`/api/v1/orders/${id}/cancel`, {
      method: 'POST',
    }),

  getOrder: (id: string) => apiFetch<Order>(`/api/v1/orders/${id}`),

  listOrders: () => apiFetch<Order[]>('/api/v1/orders'),

  createProduct: (data: { sku: string; name: string; price: number; inventory_qty: number }) =>
    apiFetch<Product>('/api/v1/seller/products', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  updateProduct: (id: string, data: { sku?: string; name?: string; price?: number; inventory_qty?: number }) =>
    apiFetch<Product>(`/api/v1/seller/products/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  deleteProduct: (id: string) =>
    apiFetch<{ status: string }>(`/api/v1/seller/products/${id}`, {
      method: 'DELETE',
    }),

  listCampaigns: () => apiFetch<Campaign[]>('/api/v1/campaigns'),

  getCampaign: (id: string) => apiFetch<Campaign>(`/api/v1/seller/campaigns/${id}`),

  createCampaign: (data: CampaignRequest) =>
    apiFetch<Campaign>('/api/v1/seller/campaigns', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  updateCampaign: (id: string, data: CampaignRequest) =>
    apiFetch<Campaign>(`/api/v1/seller/campaigns/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  deleteCampaign: (id: string) =>
    apiFetch<{ status: string }>(`/api/v1/seller/campaigns/${id}`, {
      method: 'DELETE',
    }),

  toggleCampaign: (id: string, active: boolean) =>
    apiFetch<{ id: string; active: boolean }>(`/api/v1/seller/campaigns/${id}/toggle`, {
      method: 'PATCH',
      body: JSON.stringify({ active }),
    }),
}
