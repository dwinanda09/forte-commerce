import { TopHeader } from '@/components/TopHeader'
import { Navigation } from '@/components/Navigation'

export default function ProtectedLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="min-h-screen flex flex-col">
      <TopHeader />
      <main className="flex-1 pb-20">
        {children}
      </main>
      <Navigation />
    </div>
  )
}
