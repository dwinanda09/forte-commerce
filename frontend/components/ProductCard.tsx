'use client'

import { useState } from 'react'
import Link from 'next/link'
import type { Product } from '@/lib/types'

// ─── Buyer actions ────────────────────────────────────────────────────────────

interface BuyerActionsProps {
  product: Product
  onAddToCart: (product: Product, qty: number) => void
}

function BuyerActions({ product, onAddToCart }: BuyerActionsProps) {
  const [qty, setQty] = useState(1)

  const handleAdd = () => {
    onAddToCart(product, qty)
    setQty(1)
  }

  return (
    <div className="space-y-2 pt-3 border-t border-mist/50">
      <div className="flex items-center gap-2 bg-surface rounded-md px-2 py-1.5">
        <button
          onClick={() => setQty(q => Math.max(1, q - 1))}
          className="w-7 h-7 flex items-center justify-center text-teal hover:bg-white rounded transition-colors focus-visible:outline-2 focus-visible:outline-teal"
        >
          −
        </button>
        <span className="flex-1 text-center text-sm font-medium text-graphite tabular-nums">
          {qty}
        </span>
        <button
          onClick={() => setQty(q => Math.min(product.available_qty, q + 1))}
          className="w-7 h-7 flex items-center justify-center text-teal hover:bg-white rounded transition-colors focus-visible:outline-2 focus-visible:outline-teal"
        >
          +
        </button>
      </div>
      <button
        onClick={handleAdd}
        className="w-full bg-teal text-white py-2.5 rounded-md text-sm font-medium hover:bg-teal-hover transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal"
      >
        Add to Cart
      </button>
    </div>
  )
}

// ─── Seller actions ───────────────────────────────────────────────────────────

interface SellerActionsProps {
  productId: string
  onDelete: () => void
  isDeleting?: boolean
}

function SellerActions({ productId, onDelete, isDeleting }: SellerActionsProps) {
  return (
    <div className="flex items-center justify-between pt-3 border-t border-mist/50">
      <span className="text-xs text-steel/70">
        ID <span className="font-mono">{productId.slice(0, 8)}</span>
      </span>
      <div className="flex items-center gap-3">
        <Link
          href={`/seller/products/${productId}/edit`}
          className="text-xs font-medium text-teal hover:text-teal-hover transition-colors"
        >
          Edit
        </Link>
        <button
          onClick={onDelete}
          disabled={isDeleting}
          className="text-xs text-steel hover:text-red-600 transition-colors disabled:opacity-40"
        >
          {isDeleting ? '…' : 'Delete'}
        </button>
      </div>
    </div>
  )
}

// ─── Unified ProductCard ──────────────────────────────────────────────────────

type ProductCardProps =
  | {
      product: Product
      variant: 'buyer'
      onAddToCart: (product: Product, qty: number) => void
    }
  | {
      product: Product
      variant: 'seller'
      onDelete: () => void
      isDeleting?: boolean
    }

export function ProductCard(props: ProductCardProps) {
  const { product, variant } = props
  const isOutOfStock = product.available_qty <= 0

  return (
    <div className="bg-white rounded-xl border border-mist/80 shadow-card overflow-hidden hover:shadow-modal hover:border-steel/40 transition-all group">
      {/* Image placeholder */}
      <div className="w-full aspect-[16/9] bg-gradient-to-br from-mist/60 to-steel/25 border-b border-mist/50 flex items-center justify-center">
        <span className={`text-4xl transition-opacity ${isOutOfStock ? 'opacity-30' : 'opacity-55 group-hover:opacity-75'}`}>
          📦
        </span>
      </div>

      {/* Body */}
      <div className="p-4 flex flex-col gap-1.5">
        <p className="text-[11px] font-mono text-steel/80 uppercase tracking-widest">
          {product.sku}
        </p>

        <h3 className="text-sm font-semibold text-graphite line-clamp-2 leading-snug">
          {product.name}
        </h3>

        <p className={`text-xl font-bold mt-0.5 ${isOutOfStock ? 'text-steel/40' : 'text-teal'}`}>
          ${product.price.toFixed(2)}
        </p>

        <p className="text-xs">
          {isOutOfStock ? (
            <span className="text-red-500 font-medium">Out of stock</span>
          ) : variant === 'seller' ? (
            <span className="text-steel">
              {product.inventory_qty} total
              {product.reserved_qty > 0 && (
                <span className="ml-1.5 text-amber-600">· {product.reserved_qty} reserved</span>
              )}
            </span>
          ) : (
            <span className="text-steel">{product.available_qty} in stock</span>
          )}
        </p>

        {/* Actions */}
        <div className="mt-1">
          {variant === 'buyer' && !isOutOfStock && (
            <BuyerActions product={product} onAddToCart={props.onAddToCart} />
          )}

          {variant === 'buyer' && isOutOfStock && (
            <div className="pt-3 border-t border-mist/50">
              <button
                disabled
                className="w-full bg-teal/30 text-white/70 py-2.5 rounded-md text-sm font-medium cursor-not-allowed"
              >
                Out of Stock
              </button>
            </div>
          )}

          {variant === 'seller' && (
            <SellerActions
              productId={product.id}
              onDelete={props.onDelete}
              isDeleting={props.isDeleting}
            />
          )}
        </div>
      </div>
    </div>
  )
}
