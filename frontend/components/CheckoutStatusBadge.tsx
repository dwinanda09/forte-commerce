'use client'

import type { CheckoutSession } from '@/lib/types'

interface CheckoutStatusBadgeProps {
  status: CheckoutSession['status']
}

const statusConfig: Record<
  CheckoutSession['status'],
  { bg: string; color: string; label: string }
> = {
  pending: {
    bg: '#FFF3CD',
    color: '#856404',
    label: 'Processing',
  },
  completed: {
    bg: '#D1ECE4',
    color: '#0a5e45',
    label: 'Completed',
  },
  expired: {
    bg: '#F8D7DA',
    color: '#842029',
    label: 'Expired',
  },
  failed: {
    bg: '#F8D7DA',
    color: '#842029',
    label: 'Failed',
  },
}

export function CheckoutStatusBadge({ status }: CheckoutStatusBadgeProps) {
  const config = statusConfig[status]

  return (
    <span
      className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold"
      style={{ backgroundColor: config.bg, color: config.color }}
    >
      <span
        className="w-1.5 h-1.5 rounded-full"
        style={{ backgroundColor: config.color }}
      />
      {config.label}
    </span>
  )
}
