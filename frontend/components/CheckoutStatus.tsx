'use client'

import { useEffect, useState } from 'react'
import type { CheckoutSession } from '@/lib/types'
import { CheckoutStatusBadge } from './CheckoutStatusBadge'

interface CheckoutStatusProps {
  session: CheckoutSession
  onConfirm: () => void
  isConfirming: boolean
}

export function CheckoutStatus({
  session,
  onConfirm,
  isConfirming,
}: CheckoutStatusProps) {
  const [timeLeft, setTimeLeft] = useState<string>('')

  useEffect(() => {
    if (session.status !== 'pending') return

    const timer = setInterval(() => {
      const expiresAt = new Date(session.expires_at).getTime()
      const now = Date.now()
      const diff = expiresAt - now

      if (diff <= 0) {
        setTimeLeft('Expired')
        clearInterval(timer)
      } else {
        const minutes = Math.floor(diff / 60000)
        const seconds = Math.floor((diff % 60000) / 1000)
        setTimeLeft(`${minutes}m ${seconds}s`)
      }
    }, 1000)

    return () => clearInterval(timer)
  }, [session.status, session.expires_at])

  if (session.status === 'pending') {
    return (
      <div className="flex flex-col items-center justify-center py-16">
        <div className="animate-spin mb-6">
          <div className="w-12 h-12 border-4 border-mist border-t-teal rounded-full" />
        </div>
        <p className="text-lg text-graphite mb-2">Processing your checkout...</p>
        <p className="text-sm text-steel">
          Expires in <span className="font-mono font-semibold">{timeLeft}</span>
        </p>
      </div>
    )
  }

  if (session.status === 'completed' && session.result) {
    const { result } = session

    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-semibold text-graphite">Order Summary</h2>
          <CheckoutStatusBadge status={session.status} />
        </div>

        {/* Items */}
        <div className="bg-white rounded-lg shadow-card overflow-hidden">
          <div className="border-b border-mist px-6 py-4">
            <h3 className="font-semibold text-graphite">Items</h3>
          </div>
          <div className="divide-y divide-mist">
            {result.items.map((item, idx) => (
              <div key={idx} className="px-6 py-4 flex justify-between items-center">
                <div>
                  <p className="font-medium text-graphite">{item.name}</p>
                  <p className="text-sm text-steel">Qty: {item.qty}</p>
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
        {result.promotions_applied.length > 0 && (
          <div className="bg-white rounded-lg shadow-card p-6">
            <h3 className="font-semibold text-graphite mb-4">Promotions Applied</h3>
            <div className="space-y-3">
              {result.promotions_applied.map((promo, idx) => (
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
        <div className="bg-white rounded-lg shadow-card p-6 space-y-3">
          <div className="flex justify-between text-graphite">
            <span>Subtotal</span>
            <span>${result.subtotal.toFixed(2)}</span>
          </div>
          {result.total_discount > 0 && (
            <div className="flex justify-between text-teal">
              <span>Total Discount</span>
              <span>-${result.total_discount.toFixed(2)}</span>
            </div>
          )}
          <div className="border-t border-mist pt-3 flex justify-between text-xl font-bold text-teal">
            <span>Total</span>
            <span>${result.total.toFixed(2)}</span>
          </div>
        </div>

        {/* Confirm Button */}
        <button
          onClick={onConfirm}
          disabled={isConfirming}
          className="w-full bg-teal text-white py-4 rounded-lg font-semibold text-lg hover:bg-teal-hover transition-colors disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal"
        >
          {isConfirming ? 'Confirming...' : 'Confirm & Continue to Payment'}
        </button>
      </div>
    )
  }

  if (session.status === 'expired' || session.status === 'failed') {
    return (
      <div className="bg-red-50 rounded-lg shadow-card p-8 text-center">
        <p className="text-2xl text-red-900 font-semibold mb-2">
          {session.status === 'expired'
            ? 'Checkout Expired'
            : 'Checkout Failed'}
        </p>
        {session.error_message && (
          <p className="text-red-700 mb-4">{session.error_message}</p>
        )}
        <a href="/">
          <button className="bg-teal text-white px-6 py-2 rounded-md font-medium hover:bg-teal-hover transition-colors focus-visible:outline-2 focus-visible:outline-teal">
            Return to Products
          </button>
        </a>
      </div>
    )
  }

  return null
}
