'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import useSWR from 'swr'
import { api } from '@/lib/api'
import type { Product } from '@/lib/types'

export default function EditProductPage({ params }: { params: Promise<{ id: string }> }) {
  const router = useRouter()
  const [id, setId] = useState<string | null>(null)

  useEffect(() => {
    params.then((p) => setId(p.id))
  }, [params])

  const { data: products } = useSWR(
    id ? `products-seller-${id}` : null,
    async () => {
      const res = await api.getProducts()
      return res.data
    }
  )

  const product = products?.find((p: Product) => p.id === id)

  const [form, setForm] = useState({ sku: '', name: '', price: '', inventory_qty: '' })
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (product) {
      setForm({
        sku: product.sku,
        name: product.name,
        price: String(product.price),
        inventory_qty: String(product.inventory_qty),
      })
    }
  }, [product])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!id) return
    setError(null)
    setLoading(true)
    try {
      await api.updateProduct(id, {
        sku: form.sku,
        name: form.name,
        price: parseFloat(form.price),
        inventory_qty: parseInt(form.inventory_qty, 10),
      })
      router.push('/seller/products')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update product')
    } finally {
      setLoading(false)
    }
  }

  if (!product) {
    return (
      <div className="min-h-screen bg-surface flex items-center justify-center">
        <div className="w-10 h-10 border-4 border-mist border-t-teal rounded-full animate-spin" />
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-surface">
      <div className="max-w-xl mx-auto px-6 py-12">
        <h1 className="text-3xl font-bold text-graphite mb-8">Edit Product</h1>

        <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow-card p-6 space-y-5">
          {error && (
            <div className="bg-red-50 border border-red-200 rounded-md px-4 py-3 text-sm text-red-700">
              {error}
            </div>
          )}

          {[
            { label: 'SKU', key: 'sku', type: 'text' },
            { label: 'Product Name', key: 'name', type: 'text' },
            { label: 'Price ($)', key: 'price', type: 'number' },
            { label: 'Inventory Qty', key: 'inventory_qty', type: 'number' },
          ].map(({ label, key, type }) => (
            <div key={key}>
              <label className="block text-sm font-medium text-graphite mb-1.5">{label}</label>
              <input
                type={type}
                value={form[key as keyof typeof form]}
                onChange={(e) => setForm({ ...form, [key]: e.target.value })}
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
              {loading ? 'Saving…' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
