'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

export default function SellerLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()

  const tabs = [
    { href: '/seller/products', label: 'Products' },
    { href: '/seller/campaigns', label: 'Campaigns' },
  ]

  return (
    <>
      <div className="bg-white border-b border-mist">
        <div className="max-w-5xl mx-auto px-6 flex gap-1 pt-2">
          {tabs.map(tab => {
            const active = pathname.startsWith(tab.href)
            return (
              <Link
                key={tab.href}
                href={tab.href}
                className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
                  active
                    ? 'border-teal text-teal'
                    : 'border-transparent text-steel hover:text-graphite hover:border-mist'
                }`}
              >
                {tab.label}
              </Link>
            )
          })}
        </div>
      </div>
      {children}
    </>
  )
}
