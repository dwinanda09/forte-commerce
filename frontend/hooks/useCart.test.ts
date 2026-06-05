import { describe, it, expect, beforeEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useCart, type CartItem } from './useCart'

describe('useCart', () => {
  beforeEach(() => {
    const { result } = renderHook(() => useCart())
    act(() => {
      result.current.clear()
    })
  })

  describe('addItem', () => {
    it('adds new item with qty 1', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
      })

      expect(result.current.items).toEqual([{ ...item, qty: 1 }])
    })

    it('increments qty when adding duplicate sku', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
        result.current.addItem(item)
      })

      expect(result.current.items).toHaveLength(1)
      expect(result.current.items[0].qty).toBe(2)
    })

    it('consolidates multiple duplicate additions', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
        result.current.addItem(item)
        result.current.addItem(item)
      })

      expect(result.current.items).toHaveLength(1)
      expect(result.current.items[0].qty).toBe(3)
    })

    it('handles multiple different items', () => {
      const { result } = renderHook(() => useCart())

      const item1 = { sku: 'SKU-001', name: 'Product 1', price: 100 }
      const item2 = { sku: 'SKU-002', name: 'Product 2', price: 200 }

      act(() => {
        result.current.addItem(item1)
        result.current.addItem(item2)
      })

      expect(result.current.items).toHaveLength(2)
      expect(result.current.items[0].sku).toBe('SKU-001')
      expect(result.current.items[1].sku).toBe('SKU-002')
    })

    it('preserves item properties when adding', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Premium Product', price: 999.99 }

      act(() => {
        result.current.addItem(item)
      })

      expect(result.current.items[0].name).toBe('Premium Product')
      expect(result.current.items[0].price).toBe(999.99)
    })
  })

  describe('removeItem', () => {
    it('removes item by sku', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
        result.current.removeItem('SKU-001')
      })

      expect(result.current.items).toHaveLength(0)
    })

    it('only removes matching sku', () => {
      const { result } = renderHook(() => useCart())

      const item1 = { sku: 'SKU-001', name: 'Product 1', price: 100 }
      const item2 = { sku: 'SKU-002', name: 'Product 2', price: 200 }

      act(() => {
        result.current.addItem(item1)
        result.current.addItem(item2)
        result.current.removeItem('SKU-001')
      })

      expect(result.current.items).toHaveLength(1)
      expect(result.current.items[0].sku).toBe('SKU-002')
    })

    it('handles removing non-existent sku gracefully', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
        result.current.removeItem('SKU-999')
      })

      expect(result.current.items).toHaveLength(1)
    })
  })

  describe('updateQty', () => {
    it('updates quantity for existing item', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
        result.current.updateQty('SKU-001', 5)
      })

      expect(result.current.items[0].qty).toBe(5)
    })

    it('removes item when qty is 0', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
        result.current.updateQty('SKU-001', 0)
      })

      expect(result.current.items).toHaveLength(0)
    })

    it('removes item when qty is negative', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
        result.current.updateQty('SKU-001', -1)
      })

      expect(result.current.items).toHaveLength(0)
    })

    it('preserves other items when updating qty', () => {
      const { result } = renderHook(() => useCart())

      const item1 = { sku: 'SKU-001', name: 'Product 1', price: 100 }
      const item2 = { sku: 'SKU-002', name: 'Product 2', price: 200 }

      act(() => {
        result.current.addItem(item1)
        result.current.addItem(item2)
        result.current.updateQty('SKU-001', 10)
      })

      expect(result.current.items).toHaveLength(2)
      expect(result.current.items[1].qty).toBe(1)
    })

    it('handles non-existent sku gracefully', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
        result.current.updateQty('SKU-999', 5)
      })

      expect(result.current.items).toHaveLength(1)
    })
  })

  describe('clear', () => {
    it('removes all items', () => {
      const { result } = renderHook(() => useCart())

      const item1 = { sku: 'SKU-001', name: 'Product 1', price: 100 }
      const item2 = { sku: 'SKU-002', name: 'Product 2', price: 200 }

      act(() => {
        result.current.addItem(item1)
        result.current.addItem(item2)
        result.current.clear()
      })

      expect(result.current.items).toHaveLength(0)
    })

    it('handles clearing empty cart', () => {
      const { result } = renderHook(() => useCart())

      act(() => {
        result.current.clear()
      })

      expect(result.current.items).toHaveLength(0)
    })
  })

  describe('totalItems', () => {
    it('returns sum of all quantities', () => {
      const { result } = renderHook(() => useCart())

      const item1 = { sku: 'SKU-001', name: 'Product 1', price: 100 }
      const item2 = { sku: 'SKU-002', name: 'Product 2', price: 200 }

      act(() => {
        result.current.addItem(item1)
        result.current.addItem(item2)
        result.current.addItem(item2)
      })

      expect(result.current.totalItems()).toBe(3)
    })

    it('returns 0 for empty cart', () => {
      const { result } = renderHook(() => useCart())

      expect(result.current.totalItems()).toBe(0)
    })

    it('returns correct total after updates', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
        result.current.addItem(item)
        result.current.updateQty('SKU-001', 5)
      })

      expect(result.current.totalItems()).toBe(5)
    })

    it('returns correct total after removal', () => {
      const { result } = renderHook(() => useCart())

      const item1 = { sku: 'SKU-001', name: 'Product 1', price: 100 }
      const item2 = { sku: 'SKU-002', name: 'Product 2', price: 200 }

      act(() => {
        result.current.addItem(item1)
        result.current.addItem(item1)
        result.current.addItem(item2)
        result.current.removeItem('SKU-001')
      })

      expect(result.current.totalItems()).toBe(1)
    })
  })

  describe('skuList', () => {
    it('returns flat array with one sku per qty', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
        result.current.addItem(item)
      })

      expect(result.current.skuList()).toEqual(['SKU-001', 'SKU-001'])
    })

    it('handles multiple items with different quantities', () => {
      const { result } = renderHook(() => useCart())

      const item1 = { sku: 'SKU-001', name: 'Product 1', price: 100 }
      const item2 = { sku: 'SKU-002', name: 'Product 2', price: 200 }

      act(() => {
        result.current.addItem(item1)
        result.current.addItem(item1)
        result.current.addItem(item2)
      })

      const skuList = result.current.skuList()
      expect(skuList).toHaveLength(3)
      expect(skuList.filter((s) => s === 'SKU-001')).toHaveLength(2)
      expect(skuList.filter((s) => s === 'SKU-002')).toHaveLength(1)
    })

    it('returns empty array for empty cart', () => {
      const { result } = renderHook(() => useCart())

      expect(result.current.skuList()).toEqual([])
    })

    it('returns correct list after manual qty update', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.addItem(item)
        result.current.updateQty('SKU-001', 3)
      })

      expect(result.current.skuList()).toEqual(['SKU-001', 'SKU-001', 'SKU-001'])
    })
  })

  describe('cart state persistence', () => {
    it('maintains state across multiple operations', () => {
      const { result } = renderHook(() => useCart())

      const item1 = { sku: 'SKU-001', name: 'Product 1', price: 100 }
      const item2 = { sku: 'SKU-002', name: 'Product 2', price: 200 }

      act(() => {
        result.current.addItem(item1)
        result.current.addItem(item2)
        result.current.addItem(item1)
      })

      expect(result.current.totalItems()).toBe(3)

      act(() => {
        result.current.updateQty('SKU-001', 2)
      })

      expect(result.current.totalItems()).toBe(3)
      expect(result.current.items).toHaveLength(2)
    })
  })

  describe('edge cases', () => {
    it('handles large quantities', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product 1', price: 100 }

      act(() => {
        result.current.updateQty('SKU-001', 1)
        result.current.addItem(item)
        result.current.updateQty('SKU-001', 10000)
      })

      expect(result.current.totalItems()).toBe(10000)
    })

    it('handles items with special characters in sku', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001-SPECIAL_123', name: 'Product', price: 100 }

      act(() => {
        result.current.addItem(item)
      })

      expect(result.current.items[0].sku).toBe('SKU-001-SPECIAL_123')
    })

    it('handles decimal prices', () => {
      const { result } = renderHook(() => useCart())

      const item = { sku: 'SKU-001', name: 'Product', price: 99.99 }

      act(() => {
        result.current.addItem(item)
      })

      expect(result.current.items[0].price).toBe(99.99)
    })
  })
})
