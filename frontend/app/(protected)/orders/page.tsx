'use client'

import useSWR from 'swr'
import Link from 'next/link'
import { api } from '@/lib/api'
import { OrderStatusBadge } from '@/components/OrderStatusBadge'

export default function OrdersPage() {
  const { data: orders, isLoading, error } = useSWR(
    'orders',
    async () => {
      const res = await api.listOrders()
      return res.data
    }
  )

  return (
    <div className="min-h-screen bg-surface">
      <div className="max-w-6xl mx-auto px-6 py-12">
        <h1 className="text-4xl font-bold text-graphite mb-12">Orders</h1>

        {isLoading && (
          <div className="flex items-center justify-center py-16">
            <div className="animate-spin">
              <div className="w-12 h-12 border-4 border-mist border-t-teal rounded-full" />
            </div>
          </div>
        )}

        {error && (
          <div className="bg-red-50 rounded-lg shadow-card p-8 text-center">
            <p className="text-red-900 font-semibold mb-2">Failed to load orders</p>
            <p className="text-red-700 text-sm">{error.message}</p>
          </div>
        )}

        {orders && orders.length === 0 && (
          <div className="flex flex-col items-center justify-center py-24 text-center">
            <div className="w-20 h-20 rounded-2xl bg-mist/40 flex items-center justify-center mb-6">
              <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="7" y="10" width="26" height="24" rx="3" stroke="#6D8196" strokeWidth="2" fill="none"/>
                <path d="M14 10V8a6 6 0 0 1 12 0v2" stroke="#6D8196" strokeWidth="2" strokeLinecap="round"/>
                <path d="M14 20h12M14 26h8" stroke="#B0C4DE" strokeWidth="2" strokeLinecap="round"/>
              </svg>
            </div>
            <h2 className="text-xl font-semibold text-graphite mb-2">No orders yet</h2>
            <p className="text-steel text-sm mb-8 max-w-xs">
              Your completed purchases will appear here. Browse the shop to get started.
            </p>
            <Link href="/">
              <button className="bg-teal text-white px-6 py-2.5 rounded-md font-medium hover:bg-teal-hover transition-colors focus-visible:outline-2 focus-visible:outline-teal">
                Browse Shop
              </button>
            </Link>
          </div>
        )}

        {orders && orders.length > 0 && (
          <div className="bg-white rounded-lg shadow-card overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="border-b border-mist bg-surface">
                  <tr>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-graphite">
                      Order ID
                    </th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-graphite">
                      Items
                    </th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-graphite">
                      Total
                    </th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-graphite">
                      Status
                    </th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-graphite">
                      Action
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-mist">
                  {orders.map((order) => (
                    <tr key={order.id} className="hover:bg-surface transition-colors">
                      <td className="px-6 py-4 text-sm font-mono text-steel">
                        {order.id.slice(0, 12)}...
                      </td>
                      <td className="px-6 py-4 text-sm text-graphite">
                        {order.items.length} item{order.items.length !== 1 ? 's' : ''}
                      </td>
                      <td className="px-6 py-4 text-sm font-semibold text-teal">
                        ${order.total.toFixed(2)}
                      </td>
                      <td className="px-6 py-4 text-sm">
                        <OrderStatusBadge status={order.status} />
                      </td>
                      <td className="px-6 py-4 text-sm">
                        <Link href={`/orders/${order.id}`}>
                          <button className="text-teal hover:text-teal-hover font-medium focus-visible:outline-2 focus-visible:outline-teal">
                            View Details
                          </button>
                        </Link>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
