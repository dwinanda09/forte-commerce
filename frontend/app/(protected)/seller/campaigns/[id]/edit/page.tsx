'use client'

import { useState } from 'react'
import { useRouter, useParams } from 'next/navigation'
import useSWR from 'swr'
import { api } from '@/lib/api'
import type { Condition, Action } from '@/lib/types'
import { CampaignForm } from '@/components/CampaignForm'

export default function EditCampaignPage() {
  const router = useRouter()
  const params = useParams()
  const id = params.id as string

  const { data: campaign, isLoading, error: fetchError } = useSWR(
    id ? `campaign-${id}` : null,
    () => api.getCampaign(id).then(r => r.data)
  )

  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (data: {
    name: string
    description: string
    is_active: boolean
    priority: number
    conditions: Condition[]
    actions: Action[]
  }) => {
    setError(null)
    setLoading(true)
    try {
      await api.updateCampaign(id, data)
      router.push('/seller/campaigns')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update campaign')
    } finally {
      setLoading(false)
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-surface flex items-center justify-center">
        <div className="w-10 h-10 border-4 border-mist border-t-teal rounded-full animate-spin" />
      </div>
    )
  }

  if (fetchError || !campaign) {
    return (
      <div className="min-h-screen bg-surface flex items-center justify-center">
        <p className="text-steel text-sm">Campaign not found.</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-surface">
      <div className="max-w-2xl mx-auto px-6 py-12">
        <h1 className="text-3xl font-bold text-graphite mb-8">Edit Campaign</h1>
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-md px-4 py-3 text-sm text-red-700 mb-6">
            {error}
          </div>
        )}
        <CampaignForm
          initialValues={{
            name: campaign.name,
            description: campaign.description,
            is_active: campaign.is_active,
            priority: campaign.priority,
            conditions: campaign.conditions ?? [],
            actions: campaign.actions ?? [],
          }}
          onSubmit={handleSubmit}
          loading={loading}
          submitLabel="Save Changes"
        />
      </div>
    </div>
  )
}
