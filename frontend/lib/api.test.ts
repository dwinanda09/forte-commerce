import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { api } from './api'

describe('api module', () => {
  let fetchMock: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchMock = vi.fn()
    global.fetch = fetchMock
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('login', () => {
    it('posts credentials and returns token', async () => {
      const responseData = { data: { token: 'test-token' }, success: true, meta: { request_id: 'req-1', timestamp: '2026-01-01T00:00:00Z' } }
      fetchMock.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      })

      const result = await api.login('demo', 'password')

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/auth/login', {
        method: 'POST',
        body: JSON.stringify({ username: 'demo', password: 'password' }),
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      })
      expect(result.data.token).toBe('test-token')
    })

    it('throws error on failed login', async () => {
      fetchMock.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => ({ message: 'Invalid credentials' }),
      })

      await expect(api.login('demo', 'wrong')).rejects.toThrow('Invalid credentials')
    })

    it('handles unknown error response', async () => {
      fetchMock.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => {
          throw new Error()
        },
      })

      await expect(api.login('demo', 'pass')).rejects.toThrow()
    })
  })

  describe('getProducts', () => {
    it('fetches products list', async () => {
      const products = [
        { id: '1', sku: 'SKU-001', name: 'Product 1', price: 100, inventory_qty: 10, reserved_qty: 2, available_qty: 8 },
      ]
      const responseData = { data: products, success: true, meta: { request_id: 'req-1', timestamp: '2026-01-01T00:00:00Z' } }
      fetchMock.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      })

      const result = await api.getProducts()

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/products', {
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      })
      expect(result.data).toEqual(products)
    })
  })

  describe('submitCheckout', () => {
    it('posts items to checkout endpoint', async () => {
      const responseData = { data: { checkout_id: 'chk-123' }, success: true, meta: { request_id: 'req-1', timestamp: '2026-01-01T00:00:00Z' } }
      fetchMock.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      })

      const items = ['SKU-001', 'SKU-002']
      const result = await api.submitCheckout(items)

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/checkout', {
        method: 'POST',
        body: JSON.stringify({ items }),
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      })
      expect(result.data.checkout_id).toBe('chk-123')
    })
  })

  describe('getCheckout', () => {
    it('fetches checkout session by id', async () => {
      const session = {
        checkout_id: 'chk-123',
        status: 'pending' as const,
        expires_at: '2026-01-02T00:00:00Z',
      }
      const responseData = { data: session, success: true, meta: { request_id: 'req-1', timestamp: '2026-01-01T00:00:00Z' } }
      fetchMock.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      })

      const result = await api.getCheckout('chk-123')

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/checkout/chk-123', {
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      })
      expect(result.data.checkout_id).toBe('chk-123')
    })
  })

  describe('confirmCheckout', () => {
    it('confirms checkout and returns order', async () => {
      const order = {
        id: 'order-123',
        checkout_session_id: 'chk-123',
        status: 'pending' as const,
        items: [],
        promotions_applied: [],
        subtotal: 100,
        total_discount: 0,
        total: 100,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      }
      const responseData = { data: order, success: true, meta: { request_id: 'req-1', timestamp: '2026-01-01T00:00:00Z' } }
      fetchMock.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      })

      const result = await api.confirmCheckout('chk-123')

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/checkout/chk-123/confirm', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      })
      expect(result.data.id).toBe('order-123')
    })
  })

  describe('payOrder', () => {
    it('posts payment for order', async () => {
      const order = {
        id: 'order-123',
        checkout_session_id: 'chk-123',
        status: 'paid' as const,
        items: [],
        promotions_applied: [],
        subtotal: 100,
        total_discount: 0,
        total: 100,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      }
      const responseData = { data: order, success: true, meta: { request_id: 'req-1', timestamp: '2026-01-01T00:00:00Z' } }
      fetchMock.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      })

      const result = await api.payOrder('order-123')

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/orders/order-123/pay', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      })
      expect(result.data.status).toBe('paid')
    })
  })

  describe('cancelOrder', () => {
    it('cancels order', async () => {
      const order = {
        id: 'order-123',
        checkout_session_id: 'chk-123',
        status: 'cancelled' as const,
        items: [],
        promotions_applied: [],
        subtotal: 100,
        total_discount: 0,
        total: 100,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      }
      const responseData = { data: order, success: true, meta: { request_id: 'req-1', timestamp: '2026-01-01T00:00:00Z' } }
      fetchMock.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      })

      const result = await api.cancelOrder('order-123')

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/orders/order-123/cancel', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      })
      expect(result.data.status).toBe('cancelled')
    })
  })

  describe('getOrder', () => {
    it('fetches single order by id', async () => {
      const order = {
        id: 'order-123',
        checkout_session_id: 'chk-123',
        status: 'paid' as const,
        items: [],
        promotions_applied: [],
        subtotal: 100,
        total_discount: 0,
        total: 100,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      }
      const responseData = { data: order, success: true, meta: { request_id: 'req-1', timestamp: '2026-01-01T00:00:00Z' } }
      fetchMock.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      })

      const result = await api.getOrder('order-123')

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/orders/order-123', {
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      })
      expect(result.data.id).toBe('order-123')
    })
  })

  describe('listOrders', () => {
    it('fetches all orders', async () => {
      const orders = [
        {
          id: 'order-123',
          checkout_session_id: 'chk-123',
          status: 'paid' as const,
          items: [],
          promotions_applied: [],
          subtotal: 100,
          total_discount: 0,
          total: 100,
          created_at: '2026-01-01T00:00:00Z',
          updated_at: '2026-01-01T00:00:00Z',
        },
      ]
      const responseData = { data: orders, success: true, meta: { request_id: 'req-1', timestamp: '2026-01-01T00:00:00Z' } }
      fetchMock.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      })

      const result = await api.listOrders()

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/orders', {
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      })
      expect(result.data).toEqual(orders)
    })
  })

  it('includes credentials and content-type header on all requests', async () => {
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ data: [], success: true, meta: { request_id: 'req-1', timestamp: '2026-01-01T00:00:00Z' } }),
    })

    await api.getProducts()

    const call = fetchMock.mock.calls[0]
    expect(call[1].credentials).toBe('include')
    expect(call[1].headers['Content-Type']).toBe('application/json')
  })
})
