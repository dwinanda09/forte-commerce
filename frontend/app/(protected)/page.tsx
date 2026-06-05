'use client'

import { useState } from 'react'
import useSWR from 'swr'
import { api } from '@/lib/api'
import { useCart } from '@/hooks/useCart'
import { ProductCard } from '@/components/ProductCard'
import { CartDrawer } from '@/components/CartDrawer'
import { Toast } from '@/components/Toast'
import type { Product } from '@/lib/types'

export default function ProductsPage() {
  const [cartOpen, setCartOpen] = useState(false)
  const [toast, setToast] = useState<string | null>(null)
  const cart = useCart()

  const { data: products, isLoading, error } = useSWR(
    'products',
    async () => {
      const res = await api.getProducts()
      return res.data
    }
  )

  const handleAddToCart = (product: Product, qty: number) => {
    for (let i = 0; i < qty; i++) {
      cart.addItem({
        sku: product.sku,
        name: product.name,
        price: product.price,
      })
    }
    setToast(`${product.name} added to cart`)
  }

  const handleRemoveFromCart = (sku: string) => {
    cart.removeItem(sku)
  }

  const handleUpdateQty = (sku: string, qty: number) => {
    if (qty <= 0) {
      cart.removeItem(sku)
    } else {
      cart.updateQty(sku, qty)
    }
  }

  const handleClearCart = () => {
    cart.clear()
  }

  return (
    <div className="min-h-screen bg-surface">
      {/* Cart Button (Fixed) */}
      {cart.items.length > 0 && (
        <button
          onClick={() => setCartOpen(true)}
          className="fixed bottom-24 right-6 z-40 bg-teal text-white rounded-full w-16 h-16 flex items-center justify-center shadow-modal hover:bg-teal-hover transition-all focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal"
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-6 h-6">
            <circle cx="9" cy="21" r="1" />
            <circle cx="20" cy="21" r="1" />
            <path d="M1 1h4l2.68 13.39a2 2 0 002 1.61h9.72a2 2 0 002-1.61L23 6H6" />
          </svg>
          <span className="absolute -top-1.5 -right-1.5 bg-white text-teal text-xs font-bold rounded-full w-5 h-5 flex items-center justify-center shadow-card leading-none">
            {cart.totalItems()}
          </span>
        </button>
      )}

      {/* Cart Drawer */}
      <CartDrawer
        items={cart.items}
        isOpen={cartOpen}
        onClose={() => setCartOpen(false)}
        onRemove={handleRemoveFromCart}
        onUpdateQty={handleUpdateQty}
        onClear={handleClearCart}
      />

      {/* Toast Notification */}
      {toast && <Toast message={toast} onDismiss={() => setToast(null)} />}

      {/* Products Container */}
      <div className="max-w-7xl mx-auto px-6 py-12">
        <h1 className="text-4xl font-bold text-graphite mb-12">Products</h1>

        {isLoading && (
          <div className="flex items-center justify-center py-16">
            <div className="animate-spin">
              <div className="w-12 h-12 border-4 border-mist border-t-teal rounded-full" />
            </div>
          </div>
        )}

        {error && (
          <div className="bg-red-50 rounded-lg shadow-card p-8 text-center">
            <p className="text-red-900 font-semibold mb-2">Failed to load products</p>
            <p className="text-red-700 text-sm">{error.message}</p>
          </div>
        )}

        {products && products.length === 0 && (
          <div className="text-center py-16">
            <p className="text-steel text-lg">No products available</p>
          </div>
        )}

        {products && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {products.map((product) => (
              <ProductCard
                key={product.id}
                product={product}
                variant="buyer"
                onAddToCart={handleAddToCart}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
