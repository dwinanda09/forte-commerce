'use client'

import { useState } from 'react'
import useSWR from 'swr'
import Link from 'next/link'
import { api } from '@/lib/api'
import { ProductCard } from '@/components/ProductCard'
import type { Product } from '@/lib/types'

export default function SellerProductsPage() {
  const { data: products, isLoading, mutate } = useSWR(
    'products-seller',
    async () => {
      const res = await api.getProducts()
      return res.data
    }
  )

  const [deleting, setDeleting] = useState<string | null>(null)

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this product?')) return
    setDeleting(id)
    try {
      await api.deleteProduct(id)
      mutate()
    } catch {
      alert('Failed to delete product')
    } finally {
      setDeleting(null)
    }
  }

  return (
    <div className="min-h-screen bg-surface">
      <div className="max-w-5xl mx-auto px-6 py-12">
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-2xl font-bold text-graphite">Product Management</h1>
          <Link
            href="/seller/products/new"
            className="inline-flex items-center gap-1.5 text-sm font-medium text-teal hover:text-teal-hover transition-colors"
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-4 h-4">
              <path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
            </svg>
            New product
          </Link>
        </div>

        {isLoading && (
          <div className="flex justify-center py-16">
            <div className="w-8 h-8 border-2 border-mist border-t-teal rounded-full animate-spin" />
          </div>
        )}

        {products && products.length === 0 && (
          <div className="bg-white rounded-lg shadow-card border border-mist/70 px-8 py-16 text-center">
            <p className="text-steel text-sm mb-4">No products yet.</p>
            <Link
              href="/seller/products/new"
              className="inline-flex items-center gap-1.5 text-sm text-teal font-medium hover:text-teal-hover transition-colors"
            >
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-4 h-4">
                <path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
              </svg>
              Add your first product
            </Link>
          </div>
        )}

        {products && products.length > 0 && (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {products.map((p: Product) => (
              <ProductCard
                key={p.id}
                product={p}
                variant="seller"
                onDelete={() => handleDelete(p.id)}
                isDeleting={deleting === p.id}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
