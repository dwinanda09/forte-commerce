import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ProductCard } from './ProductCard'
import type { Product } from '@/lib/types'

describe('ProductCard', () => {
  const mockProduct: Product = {
    id: '1',
    sku: 'SKU-001',
    name: 'Test Product',
    price: 99.99,
    inventory_qty: 10,
    reserved_qty: 2,
    available_qty: 8,
  }

  it('renders product sku', () => {
    const mockOnAdd = vi.fn()
    render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    expect(screen.getByText('SKU-001')).toBeInTheDocument()
  })

  it('renders product name', () => {
    const mockOnAdd = vi.fn()
    render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    expect(screen.getByText('Test Product')).toBeInTheDocument()
  })

  it('renders product price with 2 decimals', () => {
    const mockOnAdd = vi.fn()
    render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    expect(screen.getByText('$99.99')).toBeInTheDocument()
  })

  it('renders available quantity', () => {
    const mockOnAdd = vi.fn()
    render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    expect(screen.getByText('8 in stock')).toBeInTheDocument()
  })

  it('shows out of stock message when available qty is 0', () => {
    const mockOnAdd = vi.fn()
    const outOfStockProduct = { ...mockProduct, available_qty: 0 }
    render(<ProductCard product={outOfStockProduct} onAddToCart={mockOnAdd} />)

    expect(screen.getByText('Out of stock')).toBeInTheDocument()
  })

  it('disables add button when out of stock', () => {
    const mockOnAdd = vi.fn()
    const outOfStockProduct = { ...mockProduct, available_qty: 0 }
    render(<ProductCard product={outOfStockProduct} onAddToCart={mockOnAdd} />)

    const button = screen.getByRole('button', { name: 'Out of Stock' })
    expect(button).toBeDisabled()
  })

  it('shows quantity selector when in stock', () => {
    const mockOnAdd = vi.fn()
    render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    expect(screen.getByText('1')).toBeInTheDocument()
  })

  it('hides quantity selector when out of stock', () => {
    const mockOnAdd = vi.fn()
    const outOfStockProduct = { ...mockProduct, available_qty: 0 }
    render(<ProductCard product={outOfStockProduct} onAddToCart={mockOnAdd} />)

    const qtyDisplayElements = screen.queryAllByText(/^\d+$/)
    expect(qtyDisplayElements).toHaveLength(0)
  })

  it('increases quantity when plus button is clicked', async () => {
    const mockOnAdd = vi.fn()
    const user = userEvent.setup()
    render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    const plusButton = screen.getAllByRole('button').find((btn) => btn.textContent === '+')
    await user.click(plusButton!)

    expect(screen.getByText('2')).toBeInTheDocument()
  })

  it('decreases quantity when minus button is clicked', async () => {
    const mockOnAdd = vi.fn()
    const user = userEvent.setup()
    render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    const plusButton = screen.getAllByRole('button').find((btn) => btn.textContent === '+')
    await user.click(plusButton!)
    expect(screen.getByText('2')).toBeInTheDocument()

    const minusButton = screen.getAllByRole('button').find((btn) => btn.textContent === '−')
    await user.click(minusButton!)

    expect(screen.getByText('1')).toBeInTheDocument()
  })

  it('does not decrease qty below 1', async () => {
    const mockOnAdd = vi.fn()
    const user = userEvent.setup()
    render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    const minusButton = screen.getAllByRole('button').find((btn) => btn.textContent === '−')
    await user.click(minusButton!)

    expect(screen.getByText('1')).toBeInTheDocument()
  })

  it('does not increase qty beyond available quantity', async () => {
    const mockOnAdd = vi.fn()
    const user = userEvent.setup()
    const limitedProduct = { ...mockProduct, available_qty: 2 }
    render(<ProductCard product={limitedProduct} onAddToCart={mockOnAdd} />)

    const plusButton = screen.getAllByRole('button').find((btn) => btn.textContent === '+')

    await user.click(plusButton!)
    expect(screen.getByText('2')).toBeInTheDocument()

    await user.click(plusButton!)
    expect(screen.getByText('2')).toBeInTheDocument()
  })

  it('calls onAddToCart with product and qty when add button is clicked', async () => {
    const mockOnAdd = vi.fn()
    const user = userEvent.setup()
    render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    const addButton = screen.getByRole('button', { name: 'Add to Cart' })
    await user.click(addButton)

    expect(mockOnAdd).toHaveBeenCalledWith(mockProduct, 1)
  })

  it('resets quantity to 1 after adding to cart', async () => {
    const mockOnAdd = vi.fn()
    const user = userEvent.setup()
    render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    const plusButton = screen.getAllByRole('button').find((btn) => btn.textContent === '+')
    await user.click(plusButton!)
    expect(screen.getByText('2')).toBeInTheDocument()

    const addButton = screen.getByRole('button', { name: 'Add to Cart' })
    await user.click(addButton)

    expect(screen.getByText('1')).toBeInTheDocument()
  })

  it('handles multiple add to cart calls with different quantities', async () => {
    const mockOnAdd = vi.fn()
    const user = userEvent.setup()
    render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    const plusButton = screen.getAllByRole('button').find((btn) => btn.textContent === '+')
    const addButton = screen.getByRole('button', { name: 'Add to Cart' })

    await user.click(plusButton!)
    await user.click(plusButton!)
    await user.click(addButton)

    expect(mockOnAdd).toHaveBeenCalledWith(mockProduct, 3)

    await user.click(plusButton!)
    await user.click(addButton)

    expect(mockOnAdd).toHaveBeenCalledWith(mockProduct, 2)
    expect(mockOnAdd).toHaveBeenCalledTimes(2)
  })

  it('renders with high available quantity', () => {
    const mockOnAdd = vi.fn()
    const highStockProduct = { ...mockProduct, available_qty: 1000 }
    render(<ProductCard product={highStockProduct} onAddToCart={mockOnAdd} />)

    expect(screen.getByText('1000 in stock')).toBeInTheDocument()
  })

  it('handles product name truncation in UI', () => {
    const mockOnAdd = vi.fn()
    const longNameProduct = {
      ...mockProduct,
      name: 'This is a very long product name that should be truncated in the UI to prevent layout issues',
    }
    render(<ProductCard product={longNameProduct} onAddToCart={mockOnAdd} />)

    expect(screen.getByText(longNameProduct.name)).toBeInTheDocument()
  })

  it('displays teal color for price when in stock', () => {
    const mockOnAdd = vi.fn()
    const { container } = render(<ProductCard product={mockProduct} onAddToCart={mockOnAdd} />)

    const priceElement = container.querySelector('.text-teal')
    expect(priceElement).toBeInTheDocument()
  })

  it('displays dimmed price when out of stock', () => {
    const mockOnAdd = vi.fn()
    const outOfStockProduct = { ...mockProduct, available_qty: 0 }
    const { container } = render(<ProductCard product={outOfStockProduct} onAddToCart={mockOnAdd} />)

    const priceElement = container.querySelector('.opacity-40')
    expect(priceElement).toBeInTheDocument()
  })

  it('handles decimal prices correctly', () => {
    const mockOnAdd = vi.fn()
    const decimalPriceProduct = { ...mockProduct, price: 199.99 }
    render(<ProductCard product={decimalPriceProduct} onAddToCart={mockOnAdd} />)

    expect(screen.getByText('$199.99')).toBeInTheDocument()
  })

  it('handles single unit available qty', () => {
    const mockOnAdd = vi.fn()
    const singleUnitProduct = { ...mockProduct, available_qty: 1 }
    render(<ProductCard product={singleUnitProduct} onAddToCart={mockOnAdd} />)

    expect(screen.getByText('1 in stock')).toBeInTheDocument()
  })

  it('correctly identifies out of stock with negative available qty edge case', () => {
    const mockOnAdd = vi.fn()
    const negativeProduct = { ...mockProduct, available_qty: -1 }
    render(<ProductCard product={negativeProduct} onAddToCart={mockOnAdd} />)

    expect(screen.getByText('Out of stock')).toBeInTheDocument()
  })
})
