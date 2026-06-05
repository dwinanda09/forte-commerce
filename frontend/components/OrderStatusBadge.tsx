'use client'

import type { Order } from '@/lib/types'

interface OrderStatusBadgeProps {
  status: Order['status']
}

const statusConfig: Record<
  Order['status'],
  { bg: string; color: string; label: string }
> = {
  pending: {
    bg: '#FFF3CD',
    color: '#856404',
    label: 'Pending',
  },
  paid: {
    bg: '#D1ECE4',
    color: '#0a5e45',
    label: 'Paid',
  },
  cancelled: {
    bg: '#E2E3E5',
    color: '#41464b',
    label: 'Cancelled',
  },
}

export function OrderStatusBadge({ status }: OrderStatusBadgeProps) {
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
