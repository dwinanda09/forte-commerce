'use client'

import useSWR from 'swr'
import Link from 'next/link'
import { api } from '@/lib/api'
import type { CheckoutItem, Order } from '@/lib/types'

interface OwnedProduct {
  sku: string
  name: string
  price: number
  totalQty: number
  totalPaid: number
  lastPurchased: string
  orderIds: string[]
}

function buildLibrary(orders: Order[]): OwnedProduct[] {
  const map = new Map<string, OwnedProduct>()

  for (const order of orders) {
    if (order.status !== 'paid') continue
    for (const item of order.items) {
      const existing = map.get(item.sku)
      if (existing) {
        existing.totalQty += item.qty
        existing.totalPaid += item.total
        if (order.created_at > existing.lastPurchased) {
          existing.lastPurchased = order.created_at
        }
        existing.orderIds.push(order.id)
      } else {
        map.set(item.sku, {
          sku: item.sku,
          name: item.name,
          price: item.price,
          totalQty: item.qty,
          totalPaid: item.total,
          lastPurchased: order.created_at,
          orderIds: [order.id],
        })
      }
    }
  }

  return Array.from(map.values()).sort(
    (a, b) => new Date(b.lastPurchased).getTime() - new Date(a.lastPurchased).getTime()
  )
}

function ProductIcon({ name }: { name: string }) {
  const lower = name.toLowerCase()
  if (lower.includes('macbook') || lower.includes('laptop')) return <>💻</>
  if (lower.includes('raspberry')) return <>🍓</>
  if (lower.includes('alexa') || lower.includes('speaker')) return <>🔊</>
  if (lower.includes('google') || lower.includes('home')) return <>🏠</>
  if (lower.includes('iphone') || lower.includes('phone')) return <>📱</>
  return <>📦</>
}

function OwnedCard({ product }: { product: OwnedProduct }) {
  const date = new Date(product.lastPurchased).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })

  return (
    <div className="bg-white rounded-lg shadow-card overflow-hidden flex flex-col group hover:shadow-modal transition-shadow duration-200">
      {/* Artwork area */}
      <div className="h-32 bg-gradient-to-br from-mist to-steel/40 flex items-center justify-center text-4xl select-none">
        <ProductIcon name={product.name} />
      </div>

      {/* Body */}
      <div className="p-4 flex flex-col flex-1">
        <p className="text-xs font-mono text-steel mb-1">{product.sku}</p>
        <p className="font-semibold text-graphite leading-snug mb-3 flex-1">{product.name}</p>

        <div className="flex items-end justify-between mt-auto pt-3 border-t border-mist">
          <div>
            <p className="text-xs text-steel">Owned</p>
            <p className="text-lg font-bold text-teal">
              {product.totalQty} <span className="text-sm font-medium">unit{product.totalQty !== 1 ? 's' : ''}</span>
            </p>
          </div>
          <div className="text-right">
            <p className="text-xs text-steel">Total paid</p>
            <p className="font-semibold text-graphite">${product.totalPaid.toFixed(2)}</p>
          </div>
        </div>

        <p className="text-xs text-steel mt-2">Last purchased {date}</p>
      </div>
    </div>
  )
}

export default function LibraryPage() {
  const { data: orders, isLoading, error } = useSWR('orders-library', async () => {
    const res = await api.listOrders()
    return res.data
  })

  const library = orders ? buildLibrary(orders) : []
  const totalUnits = library.reduce((s, p) => s + p.totalQty, 0)
  const totalSpent = orders
    ? orders.filter((o) => o.status === 'paid').reduce((s, o) => s + o.total, 0)
    : 0

  return (
    <div className="min-h-screen bg-surface">
      <div className="max-w-6xl mx-auto px-6 py-12">
        {/* Header */}
        <div className="mb-10">
          <h1 className="text-4xl font-bold text-graphite mb-1">My Library</h1>
          <p className="text-steel">Products you&apos;ve purchased and own</p>
        </div>

        {/* Loading */}
        {isLoading && (
          <div className="flex items-center justify-center py-20">
            <div className="w-12 h-12 border-4 border-mist border-t-teal rounded-full animate-spin" />
          </div>
        )}

        {/* Error */}
        {error && (
          <div className="bg-red-50 rounded-lg shadow-card p-8 text-center">
            <p className="text-red-900 font-semibold mb-1">Failed to load library</p>
            <p className="text-red-700 text-sm">{error.message}</p>
          </div>
        )}

        {/* Empty — no paid orders */}
        {!isLoading && !error && library.length === 0 && (
          <div className="bg-white rounded-lg shadow-card p-16 text-center">
            <div className="text-5xl mb-5 select-none">📭</div>
            <p className="text-xl font-semibold text-graphite mb-2">Your library is empty</p>
            <p className="text-steel mb-8">Items appear here once an order is paid.</p>
            <Link href="/">
              <button className="bg-teal text-white px-8 py-2.5 rounded-md font-medium hover:bg-teal-hover transition-colors focus-visible:outline-2 focus-visible:outline-teal">
                Browse Products
              </button>
            </Link>
          </div>
        )}

        {/* Collection */}
        {library.length > 0 && (
          <>
            {/* Summary strip */}
            <div className="flex gap-6 mb-8">
              <div className="bg-white rounded-lg shadow-card px-6 py-4">
                <p className="text-xs text-steel uppercase tracking-wide mb-0.5">Products</p>
                <p className="text-2xl font-bold text-graphite">{library.length}</p>
              </div>
              <div className="bg-white rounded-lg shadow-card px-6 py-4">
                <p className="text-xs text-steel uppercase tracking-wide mb-0.5">Units owned</p>
                <p className="text-2xl font-bold text-graphite">{totalUnits}</p>
              </div>
              <div className="bg-white rounded-lg shadow-card px-6 py-4">
                <p className="text-xs text-steel uppercase tracking-wide mb-0.5">Total spent</p>
                <p className="text-2xl font-bold text-teal">${totalSpent.toFixed(2)}</p>
              </div>
            </div>

            {/* Grid */}
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
              {library.map((product) => (
                <OwnedCard key={product.sku} product={product} />
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  )
}
