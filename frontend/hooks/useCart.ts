'use client'

import { create } from 'zustand'

export interface CartItem {
  sku: string
  name: string
  price: number
  qty: number
}

interface CartStore {
  items: CartItem[]
  addItem: (item: Omit<CartItem, 'qty'>) => void
  removeItem: (sku: string) => void
  updateQty: (sku: string, qty: number) => void
  clear: () => void
  totalItems: () => number
  skuList: () => string[]
}

const useCart = create<CartStore>((set, get) => ({
  items: [],

  addItem: (item) =>
    set((state) => {
      const existing = state.items.find((i) => i.sku === item.sku)
      if (existing) {
        return {
          items: state.items.map((i) =>
            i.sku === item.sku ? { ...i, qty: i.qty + 1 } : i
          ),
        }
      }
      return {
        items: [...state.items, { ...item, qty: 1 }],
      }
    }),

  removeItem: (sku) =>
    set((state) => ({
      items: state.items.filter((i) => i.sku !== sku),
    })),

  updateQty: (sku, qty) =>
    set((state) => {
      if (qty <= 0) {
        return { items: state.items.filter((i) => i.sku !== sku) }
      }
      return {
        items: state.items.map((i) =>
          i.sku === sku ? { ...i, qty } : i
        ),
      }
    }),

  clear: () => set({ items: [] }),

  totalItems: () => {
    return get().items.reduce((sum, item) => sum + item.qty, 0)
  },

  skuList: () => {
    const list: string[] = []
    get().items.forEach((item) => {
      for (let i = 0; i < item.qty; i++) {
        list.push(item.sku)
      }
    })
    return list
  },
}))

export { useCart }
