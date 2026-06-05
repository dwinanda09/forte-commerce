'use client'

import { useEffect } from 'react'

interface ToastProps {
  message: string
  onDismiss: () => void
  duration?: number
}

export function Toast({ message, onDismiss, duration = 2500 }: ToastProps) {
  useEffect(() => {
    const t = setTimeout(onDismiss, duration)
    return () => clearTimeout(t)
  }, [onDismiss, duration])

  return (
    <div className="fixed bottom-28 right-6 z-50 flex items-center gap-3 bg-graphite text-white px-4 py-3 rounded-lg shadow-modal text-sm font-medium animate-slide-up max-w-xs">
      <span className="text-teal text-base">✓</span>
      {message}
    </div>
  )
}
