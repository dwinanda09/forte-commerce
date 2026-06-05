'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { api } from '@/lib/api'
import type { Condition, Action, ConditionType, ActionType } from '@/lib/types'
import { CampaignForm } from '@/components/CampaignForm'

export default function NewCampaignPage() {
  const router = useRouter()
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
      await api.createCampaign(data)
      router.push('/seller/campaigns')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create campaign')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-surface">
      <div className="max-w-2xl mx-auto px-6 py-12">
        <h1 className="text-3xl font-bold text-graphite mb-8">New Campaign</h1>
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-md px-4 py-3 text-sm text-red-700 mb-6">
            {error}
          </div>
        )}
        <CampaignForm
          onSubmit={handleSubmit}
          loading={loading}
          submitLabel="Create Campaign"
        />
      </div>
    </div>
  )
}
