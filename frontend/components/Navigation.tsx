'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'

function HomeIcon({ active }: { active: boolean }) {
  return (
    <svg
      width="22"
      height="22"
      viewBox="0 0 24 24"
      fill={active ? 'currentColor' : 'none'}
      stroke="currentColor"
      strokeWidth={active ? 0 : 1.75}
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M3 9.5L12 3l9 6.5V20a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V9.5z" />
      <path d="M9 21V12h6v9" />
    </svg>
  )
}

function OrdersIcon({ active }: { active: boolean }) {
  return (
    <svg
      width="22"
      height="22"
      viewBox="0 0 24 24"
      fill={active ? 'currentColor' : 'none'}
      stroke="currentColor"
      strokeWidth={active ? 0 : 1.75}
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <rect x="3" y="3" width="18" height="18" rx="2" />
      {!active && (
        <>
          <line x1="8" y1="8" x2="16" y2="8" />
          <line x1="8" y1="12" x2="16" y2="12" />
          <line x1="8" y1="16" x2="12" y2="16" />
        </>
      )}
      {active && (
        <path d="M3 3h18v18H3V3zm5 5h8M8 12h8M8 16h4" stroke="white" strokeWidth="1.75" fill="none" strokeLinecap="round" />
      )}
    </svg>
  )
}

function LibraryIcon({ active }: { active: boolean }) {
  return (
    <svg
      width="22"
      height="22"
      viewBox="0 0 24 24"
      fill={active ? 'currentColor' : 'none'}
      stroke="currentColor"
      strokeWidth={active ? 0 : 1.75}
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z" />
      <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" />
      {active && (
        <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2zM22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z" stroke="white" strokeWidth="1.5" fill="none" />
      )}
    </svg>
  )
}

function StoreIcon({ active }: { active: boolean }) {
  return (
    <svg
      width="22"
      height="22"
      viewBox="0 0 24 24"
      fill={active ? 'currentColor' : 'none'}
      stroke="currentColor"
      strokeWidth={active ? 0 : 1.75}
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M3 9l1-6h16l1 6" />
      <path d="M3 9h18v12a1 1 0 01-1 1H4a1 1 0 01-1-1V9z" />
      <path d="M9 9v3a3 3 0 006 0V9" />
    </svg>
  )
}

function SignOutIcon() {
  return (
    <svg
      width="22"
      height="22"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={1.75}
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
      <polyline points="16 17 21 12 16 7" />
      <line x1="21" y1="12" x2="9" y2="12" />
    </svg>
  )
}

export function Navigation() {
  const pathname = usePathname()
  const { logout, role } = useAuth()

  const isHome = pathname === '/'
  const isOrders = pathname.startsWith('/orders')
  const isLibrary = pathname.startsWith('/library')
  const isManager = pathname.startsWith('/seller')

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-50 bg-white border-t border-mist shadow-[0_-1px_3px_rgba(0,0,0,.06)]">
      <div className="max-w-7xl mx-auto flex items-stretch h-16">
        {/* Shop */}
        <Link
          href="/"
          className={`flex-1 flex flex-col items-center justify-center gap-0.5 text-xs font-medium transition-colors focus-visible:outline-none ${
            isHome ? 'text-teal' : 'text-steel hover:text-graphite'
          }`}
        >
          <HomeIcon active={isHome} />
          <span>Shop</span>
        </Link>

        {/* Orders */}
        <Link
          href="/orders"
          className={`flex-1 flex flex-col items-center justify-center gap-0.5 text-xs font-medium transition-colors focus-visible:outline-none ${
            isOrders ? 'text-teal' : 'text-steel hover:text-graphite'
          }`}
        >
          <OrdersIcon active={isOrders} />
          <span>Orders</span>
        </Link>

        {/* Library */}
        <Link
          href="/library"
          className={`flex-1 flex flex-col items-center justify-center gap-0.5 text-xs font-medium transition-colors focus-visible:outline-none ${
            isLibrary ? 'text-teal' : 'text-steel hover:text-graphite'
          }`}
        >
          <LibraryIcon active={isLibrary} />
          <span>Library</span>
        </Link>

        {/* Manage (Seller Only) */}
        {role === 'seller' && (
          <Link
            href="/seller/products"
            className={`flex-1 flex flex-col items-center justify-center gap-0.5 text-xs font-medium transition-colors focus-visible:outline-none ${
              isManager ? 'text-teal' : 'text-steel hover:text-graphite'
            }`}
          >
            <StoreIcon active={isManager} />
            <span>Manage</span>
          </Link>
        )}

        {/* Sign Out */}
        <button
          onClick={logout}
          className="flex-1 flex flex-col items-center justify-center gap-0.5 text-xs font-medium text-steel hover:text-graphite transition-colors focus-visible:outline-none"
        >
          <SignOutIcon />
          <span>Sign Out</span>
        </button>
      </div>

      {/* Safe area bottom padding for notched devices */}
      <div className="h-safe-area-inset-bottom bg-white" style={{ height: 'env(safe-area-inset-bottom)' }} />
    </nav>
  )
}
