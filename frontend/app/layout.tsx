import type { Metadata, Viewport } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'ForteCommerce',
  description: 'Premium e-commerce checkout',
  manifest: '/manifest.json',
  appleWebApp: {
    capable: true,
    statusBarStyle: 'default',
    title: 'ForteCommerce',
  },
}

export const viewport: Viewport = {
  themeColor: '#01796F',
  width: 'device-width',
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <head>
        <link rel="apple-touch-icon" href="/fc-logo.png" />
      </head>
      <body className="bg-surface text-graphite">
        {children}
      </body>
    </html>
  )
}
