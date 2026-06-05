'use client'

import Image from 'next/image'
import Link from 'next/link'

export function TopHeader() {
  return (
    <header className="bg-white border-b border-mist shadow-card sticky top-0 z-40">
      <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-center">
        <Link
          href="/"
          className="flex items-center gap-2.5 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal rounded"
        >
          <Image src="/fc-logo.png" alt="ForteCommerce" width={28} height={28} priority />
          <span className="font-bold text-base text-graphite">ForteCommerce</span>
        </Link>
      </div>
    </header>
  )
}
