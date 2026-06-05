'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { api } from '@/lib/api'

export default function NewProductPage() {
  const router = useRouter()
  const [form, setForm] = useState({ sku: '', name: '', price: '', inventory_qty: '' })
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await api.createProduct({
        sku: form.sku,
        name: form.name,
        price: parseFloat(form.price),
        inventory_qty: parseInt(form.inventory_qty, 10),
      })
      router.push('/seller/products')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create product')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-surface">
      <div className="max-w-xl mx-auto px-6 py-12">
        <h1 className="text-3xl font-bold text-graphite mb-8">Add Product</h1>

        <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow-card p-6 space-y-5">
          {error && (
            <div className="bg-red-50 border border-red-200 rounded-md px-4 py-3 text-sm text-red-700">
              {error}
            </div>
          )}

          {[
            { label: 'SKU', key: 'sku', type: 'text', placeholder: 'e.g. ABC123' },
            { label: 'Product Name', key: 'name', type: 'text', placeholder: 'e.g. Google Home' },
            { label: 'Price ($)', key: 'price', type: 'number', placeholder: '49.99' },
            { label: 'Inventory Qty', key: 'inventory_qty', type: 'number', placeholder: '10' },
          ].map(({ label, key, type, placeholder }) => (
            <div key={key}>
              <label className="block text-sm font-medium text-graphite mb-1.5">{label}</label>
              <input
                type={type}
                value={form[key as keyof typeof form]}
                onChange={(e) => setForm({ ...form, [key]: e.target.value })}
                placeholder={placeholder}
                required
                min={type === 'number' ? '0' : undefined}
                step={key === 'price' ? '0.01' : undefined}
                className="w-full border-1.5 border-mist rounded-md px-3.5 py-2.5 text-sm text-graphite bg-white focus:outline-none focus:border-teal focus:ring-2 focus:ring-teal/15 transition"
              />
            </div>
          ))}

          <div className="flex gap-3 pt-2">
            <button
              type="button"
              onClick={() => router.back()}
              className="flex-1 py-2.5 border border-mist rounded-md text-sm font-medium text-steel hover:bg-surface transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="flex-1 py-2.5 bg-teal text-white rounded-md text-sm font-medium hover:bg-teal-hover transition-colors disabled:opacity-50"
            >
              {loading ? 'Creating…' : 'Create Product'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
