'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { api } from '@/lib/api'
import { useCart } from '@/hooks/useCart'

export default function CheckoutInitPage() {
  const router = useRouter()
  const cart = useCart()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (cart.items.length === 0) {
      router.replace('/')
      return
    }

    const submitCheckout = async () => {
      try {
        setIsSubmitting(true)
        const skuList = cart.skuList()
        const res = await api.submitCheckout(skuList)
        cart.clear()
        router.replace(`/checkout/${res.data.checkout_id}`)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to submit checkout')
        setIsSubmitting(false)
      }
    }

    submitCheckout()
  }, [])

  if (error) {
    return (
      <div className="min-h-screen bg-surface flex items-center justify-center px-6">
        <div className="bg-white rounded-lg shadow-modal max-w-md w-full p-8 text-center">
          <p className="text-red-900 font-semibold mb-2">Checkout Failed</p>
          <p className="text-red-700 text-sm mb-6">{error}</p>
          <a href="/">
            <button className="bg-teal text-white px-6 py-2 rounded-md font-medium hover:bg-teal-hover transition-colors focus-visible:outline-2 focus-visible:outline-teal">
              Return to Products
            </button>
          </a>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-surface flex items-center justify-center px-6">
      <div className="flex flex-col items-center">
        <div className="animate-spin mb-6">
          <div className="w-12 h-12 border-4 border-mist border-t-teal rounded-full" />
        </div>
        <p className="text-lg text-graphite">Submitting your checkout...</p>
      </div>
    </div>
  )
}
