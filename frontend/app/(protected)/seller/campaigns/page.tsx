'use client'

import { useState } from 'react'
import useSWR from 'swr'
import Link from 'next/link'
import { api } from '@/lib/api'
import { AdminCard } from '@/components/AdminCard'
import type { Campaign } from '@/lib/types'

export default function SellerCampaignsPage() {
  const { data: campaigns, isLoading, mutate } = useSWR(
    'campaigns-seller',
    async () => {
      const res = await api.listCampaigns()
      return res.data
    }
  )

  const [toggling, setToggling] = useState<string | null>(null)
  const [deleting, setDeleting] = useState<string | null>(null)

  const handleToggle = async (c: Campaign) => {
    setToggling(c.id)
    try {
      await api.toggleCampaign(c.id, !c.is_active)
      mutate()
    } catch {
      alert('Failed to toggle campaign')
    } finally {
      setToggling(null)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this campaign?')) return
    setDeleting(id)
    try {
      await api.deleteCampaign(id)
      mutate()
    } catch {
      alert('Failed to delete campaign')
    } finally {
      setDeleting(null)
    }
  }

  return (
    <div className="min-h-screen bg-surface">
      <div className="max-w-5xl mx-auto px-6 py-12">
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-2xl font-bold text-graphite">Campaign Engine</h1>
          <Link
            href="/seller/campaigns/new"
            className="inline-flex items-center gap-1.5 text-sm font-medium text-teal hover:text-teal-hover transition-colors"
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-4 h-4">
              <path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
            </svg>
            New campaign
          </Link>
        </div>

        {isLoading && (
          <div className="flex justify-center py-16">
            <div className="w-8 h-8 border-2 border-mist border-t-teal rounded-full animate-spin" />
          </div>
        )}

        {campaigns && campaigns.length === 0 && (
          <div className="bg-white rounded-lg shadow-card border border-mist/70 px-8 py-16 text-center">
            <p className="text-steel text-sm mb-4">No campaigns yet.</p>
            <Link
              href="/seller/campaigns/new"
              className="inline-flex items-center gap-1.5 text-sm text-teal font-medium hover:text-teal-hover transition-colors"
            >
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-4 h-4">
                <path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
              </svg>
              Add your first campaign
            </Link>
          </div>
        )}

        {campaigns && campaigns.length > 0 && (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {campaigns.map((c: Campaign) => (
              <AdminCard
                key={c.id}
                title={c.name}
                badge={
                  <button
                    onClick={() => handleToggle(c)}
                    disabled={toggling === c.id}
                    title={c.is_active ? 'Click to deactivate' : 'Click to activate'}
                    className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-semibold transition-opacity disabled:opacity-40 ${
                      c.is_active
                        ? 'bg-emerald-50 text-emerald-700 hover:bg-emerald-100'
                        : 'bg-mist/40 text-steel hover:bg-mist/70'
                    }`}
                  >
                    <span className={`w-1.5 h-1.5 rounded-full shrink-0 ${c.is_active ? 'bg-emerald-500' : 'bg-steel/50'}`} />
                    {toggling === c.id ? '…' : c.is_active ? 'Active' : 'Off'}
                  </button>
                }
                description={c.description}
                stats={[{ key: 'Priority', value: c.priority }]}
                editHref={`/seller/campaigns/${c.id}/edit`}
                onDelete={() => handleDelete(c.id)}
                isDeleting={deleting === c.id}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
