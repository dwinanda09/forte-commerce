'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import useSWR from 'swr'
import Link from 'next/link'
import { api } from '@/lib/api'
import { OrderStatusBadge } from '@/components/OrderStatusBadge'

export default function OrderDetailPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const router = useRouter()
  const [orderId, setOrderId] = useState<string | null>(null)
  const [isActionLoading, setIsActionLoading] = useState(false)
  const [actionError, setActionError] = useState<string | null>(null)
  const [showCancelConfirm, setShowCancelConfirm] = useState(false)

  if (orderId === null) {
    params.then((p) => setOrderId(p.id))
  }

  const { data: order, isLoading, error, mutate } = useSWR(
    orderId ? `order-${orderId}` : null,
    async () => {
      if (!orderId) return null
      const res = await api.getOrder(orderId)
      return res.data
    }
  )

  const handlePay = async () => {
    if (!order) return
    setActionError(null)
    try {
      setIsActionLoading(true)
      await api.payOrder(order.id)
      mutate()
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to process payment')
    } finally {
      setIsActionLoading(false)
    }
  }

  const handleCancelConfirm = async () => {
    if (!order) return
    setShowCancelConfirm(false)
    setActionError(null)
    try {
      setIsActionLoading(true)
      await api.cancelOrder(order.id)
      mutate()
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to cancel order')
    } finally {
      setIsActionLoading(false)
    }
  }

  if (isLoading || !order) {
    return (
      <div className="min-h-screen bg-surface flex items-center justify-center px-6">
        <div className="flex flex-col items-center">
          <div className="animate-spin mb-6">
            <div className="w-12 h-12 border-4 border-mist border-t-teal rounded-full" />
          </div>
          <p className="text-lg text-graphite">Loading order...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-surface flex items-center justify-center px-6">
        <div className="bg-white rounded-lg shadow-modal max-w-md w-full p-8 text-center">
          <p className="text-red-900 font-semibold mb-2">Failed to load order</p>
          <p className="text-red-700 text-sm mb-6">{error.message}</p>
          <Link href="/orders">
            <button className="bg-teal text-white px-6 py-2 rounded-md font-medium hover:bg-teal-hover transition-colors focus-visible:outline-2 focus-visible:outline-teal">
              Back to Orders
            </button>
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-surface">
      {/* Cancel Confirmation Modal */}
      {showCancelConfirm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center px-6 bg-black/40">
          <div className="bg-white rounded-lg shadow-modal max-w-sm w-full p-8">
            <h2 className="text-xl font-bold text-graphite mb-3">Cancel Order</h2>
            <p className="text-steel mb-8">Are you sure you want to cancel this order? This cannot be undone.</p>
            <div className="flex gap-3">
              <button
                onClick={() => setShowCancelConfirm(false)}
                className="flex-1 border border-mist text-steel py-2.5 rounded-md font-medium hover:text-graphite hover:border-steel transition-colors focus-visible:outline-2 focus-visible:outline-teal"
              >
                Keep Order
              </button>
              <button
                onClick={handleCancelConfirm}
                className="flex-1 bg-red-600 text-white py-2.5 rounded-md font-semibold hover:bg-red-700 transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal"
              >
                Yes, Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="max-w-3xl mx-auto px-6 py-12">
        {/* Header */}
        <div className="mb-8">
          <Link href="/orders">
            <button className="text-teal hover:text-teal-hover font-medium mb-4 focus-visible:outline-2 focus-visible:outline-teal">
              ← Back to Orders
            </button>
          </Link>

          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-graphite mb-2">Order Details</h1>
              <p className="text-steel font-mono">ID: {order.id}</p>
            </div>
            <OrderStatusBadge status={order.status} />
          </div>
        </div>

        {/* Action Error */}
        {actionError && (
          <div className="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
            <p className="text-red-700 text-sm">{actionError}</p>
          </div>
        )}

        {/* Order Info */}
        <div className="bg-white rounded-lg shadow-card p-6 mb-6">
          <div className="grid grid-cols-2 gap-6 mb-6 pb-6 border-b border-mist">
            <div>
              <p className="text-sm text-steel mb-1">Created</p>
              <p className="font-semibold text-graphite">
                {new Date(order.created_at).toLocaleDateString()}
              </p>
            </div>
            <div>
              <p className="text-sm text-steel mb-1">Updated</p>
              <p className="font-semibold text-graphite">
                {new Date(order.updated_at).toLocaleDateString()}
              </p>
            </div>
          </div>
        </div>

        {/* Items */}
        <div className="bg-white rounded-lg shadow-card overflow-hidden mb-6">
          <div className="border-b border-mist px-6 py-4">
            <h2 className="font-semibold text-graphite">Items</h2>
          </div>
          <div className="divide-y divide-mist">
            {order.items.map((item, idx) => (
              <div
                key={idx}
                className="px-6 py-4 flex justify-between items-center hover:bg-surface transition-colors"
              >
                <div>
                  <p className="font-medium text-graphite">{item.name}</p>
                  <p className="text-sm text-steel font-mono">{item.sku}</p>
                </div>
                <div className="text-right">
                  <p className="text-sm text-steel">
                    ${item.price.toFixed(2)} × {item.qty}
                  </p>
                  <p className="font-semibold text-teal">
                    ${item.total.toFixed(2)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Promotions */}
        {order.promotions_applied.length > 0 && (
          <div className="bg-white rounded-lg shadow-card p-6 mb-6">
            <h2 className="font-semibold text-graphite mb-4">Promotions Applied</h2>
            <div className="space-y-3">
              {order.promotions_applied.map((promo, idx) => (
                <div
                  key={idx}
                  className="flex justify-between items-start pb-3 border-b border-mist last:border-b-0"
                >
                  <div>
                    <p className="font-medium text-graphite">{promo.name}</p>
                    <p className="text-sm text-steel">{promo.description}</p>
                  </div>
                  <p className="font-semibold text-teal">
                    -${promo.discount.toFixed(2)}
                  </p>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Totals */}
        <div className="bg-white rounded-lg shadow-card p-6 mb-6 space-y-3">
          <div className="flex justify-between text-graphite">
            <span>Subtotal</span>
            <span>${order.subtotal.toFixed(2)}</span>
          </div>
          {order.total_discount > 0 && (
            <div className="flex justify-between text-teal">
              <span>Total Discount</span>
              <span>-${order.total_discount.toFixed(2)}</span>
            </div>
          )}
          <div className="border-t border-mist pt-3 flex justify-between text-xl font-bold text-teal">
            <span>Total</span>
            <span>${order.total.toFixed(2)}</span>
          </div>
        </div>

        {/* Actions */}
        <div className="flex gap-4">
          {order.status === 'pending' && (
            <>
              <button
                onClick={handlePay}
                disabled={isActionLoading}
                className="flex-1 bg-teal text-white py-3 rounded-lg font-semibold hover:bg-teal-hover transition-colors disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal"
              >
                {isActionLoading ? 'Processing...' : 'Pay Now'}
              </button>
              <button
                onClick={() => setShowCancelConfirm(true)}
                disabled={isActionLoading}
                className="flex-1 bg-red-600 text-white py-3 rounded-lg font-semibold hover:bg-red-700 transition-colors disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal"
              >
                Cancel Order
              </button>
            </>
          )}

          {order.status !== 'pending' && (
            <Link href="/orders" className="w-full">
              <button className="w-full bg-teal text-white py-3 rounded-lg font-semibold hover:bg-teal-hover transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal">
                Back to Orders
              </button>
            </Link>
          )}
        </div>
      </div>
    </div>
  )
}
