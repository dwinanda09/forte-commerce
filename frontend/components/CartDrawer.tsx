'use client'

import { useState } from 'react'
import Link from 'next/link'
import type { CartItem } from '@/hooks/useCart'

interface CartDrawerProps {
  items: CartItem[]
  isOpen: boolean
  onClose: () => void
  onRemove: (sku: string) => void
  onUpdateQty: (sku: string, qty: number) => void
  onClear: () => void
}

export function CartDrawer({
  items,
  isOpen,
  onClose,
  onRemove,
  onUpdateQty,
  onClear,
}: CartDrawerProps) {
  const [isCheckingOut, setIsCheckingOut] = useState(false)

  const subtotal = items.reduce((sum, item) => sum + item.price * item.qty, 0)
  const itemCount = items.reduce((sum, item) => sum + item.qty, 0)

  const handleCheckout = async () => {
    setIsCheckingOut(true)
    // Navigation happens in the parent page
  }

  return (
    <>
      {/* Overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black/20 z-40"
          onClick={onClose}
        />
      )}

      {/* Drawer */}
      <div
        className={`fixed right-0 top-0 h-screen w-full max-w-sm bg-white shadow-modal z-50 transform transition-transform duration-300 flex flex-col ${
          isOpen ? 'translate-x-0' : 'translate-x-full'
        }`}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-5 border-b border-mist">
          <h2 className="text-xl font-semibold text-graphite">
            Cart {itemCount > 0 && <span className="text-steel">({itemCount})</span>}
          </h2>
          <button
            onClick={onClose}
            className="text-2xl text-steel hover:text-graphite transition-colors focus-visible:outline-2 focus-visible:outline-teal"
          >
            ✕
          </button>
        </div>

        {/* Items */}
        <div className="flex-1 overflow-y-auto p-5 space-y-4">
          {items.length === 0 ? (
            <p className="text-center text-steel py-8">Your cart is empty</p>
          ) : (
            items.map((item) => (
              <div
                key={item.sku}
                className="pb-4 border-b border-mist last:border-b-0"
              >
                <div className="flex justify-between items-start mb-2">
                  <div>
                    <p className="font-medium text-graphite">{item.name}</p>
                    <p className="text-sm text-steel font-mono">{item.sku}</p>
                  </div>
                  <p className="font-semibold text-teal">
                    ${(item.price * item.qty).toFixed(2)}
                  </p>
                </div>

                {/* Qty Controls */}
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => onUpdateQty(item.sku, item.qty - 1)}
                    className="px-2 py-1 text-steel hover:bg-surface rounded transition-colors focus-visible:outline-2 focus-visible:outline-teal"
                  >
                    −
                  </button>
                  <span className="flex-1 text-center text-sm">{item.qty}</span>
                  <button
                    onClick={() => onUpdateQty(item.sku, item.qty + 1)}
                    className="px-2 py-1 text-steel hover:bg-surface rounded transition-colors focus-visible:outline-2 focus-visible:outline-teal"
                  >
                    +
                  </button>
                  <button
                    onClick={() => onRemove(item.sku)}
                    className="ml-2 text-red-600 hover:text-red-700 text-sm font-medium focus-visible:outline-2 focus-visible:outline-red-500"
                  >
                    Remove
                  </button>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Footer */}
        {items.length > 0 && (
          <div className="border-t border-mist p-5 space-y-4">
            <div className="flex justify-between items-center">
              <span className="text-graphite font-medium">Subtotal</span>
              <span className="text-lg font-semibold text-teal">
                ${subtotal.toFixed(2)}
              </span>
            </div>

            <Link href={`/checkout`} onClick={onClose}>
              <button
                onClick={handleCheckout}
                disabled={isCheckingOut}
                className="w-full bg-teal text-white py-3 rounded-md font-semibold hover:bg-teal-hover transition-colors disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal"
              >
                Proceed to Checkout
              </button>
            </Link>

            <button
              onClick={onClear}
              className="w-full bg-surface text-steel py-2 rounded-md font-medium hover:bg-mist transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal"
            >
              Clear Cart
            </button>
          </div>
        )}
      </div>
    </>
  )
}
