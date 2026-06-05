'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useCheckoutPoller } from '@/hooks/useCheckoutPoller'
import { api } from '@/lib/api'
import { CheckoutStatus } from '@/components/CheckoutStatus'

export default function CheckoutStatusPage({
  params,
}: {
  params: Promise<{ id: string }>
}) {
  const [checkoutId, setCheckoutId] = useState<string | null>(null)
  const [isConfirming, setIsConfirming] = useState(false)
  const router = useRouter()

  // Initialize checkout ID from params
  if (checkoutId === null) {
    params.then((p) => setCheckoutId(p.id))
  }

  const { session, isLoading, error } = useCheckoutPoller(checkoutId)

  const handleConfirm = async () => {
    if (!session) return

    try {
      setIsConfirming(true)
      const res = await api.confirmCheckout(session.checkout_id)
      router.replace(`/orders/${res.data.id}`)
    } catch (err) {
      setIsConfirming(false)
      alert(err instanceof Error ? err.message : 'Failed to confirm order')
    }
  }

  if (isLoading || !session) {
    return (
      <div className="min-h-screen bg-surface flex items-center justify-center px-6">
        <div className="flex flex-col items-center">
          <div className="animate-spin mb-6">
            <div className="w-12 h-12 border-4 border-mist border-t-teal rounded-full" />
          </div>
          <p className="text-lg text-graphite">Loading checkout...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-surface flex items-center justify-center px-6">
        <div className="bg-white rounded-lg shadow-modal max-w-md w-full p-8 text-center">
          <p className="text-red-900 font-semibold mb-2">Error</p>
          <p className="text-red-700 text-sm mb-6">{error.message}</p>
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
    <div className="min-h-screen bg-surface">
      <div className="max-w-3xl mx-auto px-6 py-12">
        <CheckoutStatus
          session={session}
          onConfirm={handleConfirm}
          isConfirming={isConfirming}
        />
      </div>
    </div>
  )
}
