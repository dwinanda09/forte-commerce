'use client'

import useSWR from 'swr'
import { api } from '@/lib/api'
import type { CheckoutSession } from '@/lib/types'

export function useCheckoutPoller(id: string | null) {
  const shouldFetch = id !== null

  const { data, error, isLoading } = useSWR(
    shouldFetch ? id : null,
    (checkoutId) => api.getCheckout(checkoutId).then((res) => res.data),
    {
      refreshInterval: (data?: CheckoutSession) => {
        if (!data) return 2000
        if (data.status === 'pending') return 2000
        return 0
      },
      revalidateOnFocus: false,
      revalidateOnReconnect: false,
    }
  )

  return {
    session: data,
    isLoading,
    error,
  }
}
