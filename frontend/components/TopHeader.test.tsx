import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { TopHeader } from './TopHeader'

describe('TopHeader', () => {
  describe('rendering', () => {
    it('renders header element', () => {
      const { container } = render(<TopHeader />)

      const header = container.querySelector('header')
      expect(header).toBeInTheDocument()
    })

    it('renders logo image', () => {
      const { container } = render(<TopHeader />)

      const images = container.querySelectorAll('img')
      expect(images.length).toBeGreaterThan(0)
    })

    it('renders brand text', () => {
      render(<TopHeader />)

      expect(screen.getByText('ForteCommerce')).toBeInTheDocument()
    })

    it('renders link to home', () => {
      render(<TopHeader />)

      const link = screen.getByRole('link')
      expect(link).toHaveAttribute('href', '/')
    })
  })

  describe('styling', () => {
    it('header has white background', () => {
      const { container } = render(<TopHeader />)

      const header = container.querySelector('header')
      expect(header).toHaveClass('bg-white')
    })

    it('header has bottom border', () => {
      const { container } = render(<TopHeader />)

      const header = container.querySelector('header')
      expect(header).toHaveClass('border-b')
      expect(header).toHaveClass('border-mist')
    })

    it('header has shadow', () => {
      const { container } = render(<TopHeader />)

      const header = container.querySelector('header')
      expect(header).toHaveClass('shadow-card')
    })

    it('header is sticky', () => {
      const { container } = render(<TopHeader />)

      const header = container.querySelector('header')
      expect(header).toHaveClass('sticky')
      expect(header).toHaveClass('top-0')
    })

    it('header has high z-index', () => {
      const { container } = render(<TopHeader />)

      const header = container.querySelector('header')
      expect(header).toHaveClass('z-40')
    })
  })

  describe('layout', () => {
    it('has centered max-width container', () => {
      const { container } = render(<TopHeader />)

      const div = container.querySelector('.max-w-7xl')
      expect(div).toBeInTheDocument()
    })

    it('has centered content', () => {
      const { container } = render(<TopHeader />)

      const contentDiv = container.querySelector('.flex.items-center.justify-center')
      expect(contentDiv).toBeInTheDocument()
    })

    it('link has flex layout', () => {
      render(<TopHeader />)

      const link = screen.getByRole('link')
      expect(link).toHaveClass('flex')
      expect(link).toHaveClass('items-center')
      expect(link).toHaveClass('gap-2.5')
    })
  })

  describe('logo image', () => {
    it('logo has correct alt text', () => {
      const { container } = render(<TopHeader />)

      const img = container.querySelector('img')
      expect(img).toHaveAttribute('alt', 'ForteCommerce')
    })

    it('logo has correct dimensions', () => {
      const { container } = render(<TopHeader />)

      const img = container.querySelector('img')
      expect(img).toHaveAttribute('width')
      expect(img).toHaveAttribute('height')
    })
  })

  describe('accessibility', () => {
    it('link is keyboard navigable', () => {
      render(<TopHeader />)

      const link = screen.getByRole('link')
      expect(link).toBeInTheDocument()
    })

    it('link has focus styling classes', () => {
      render(<TopHeader />)

      const link = screen.getByRole('link')
      expect(link).toHaveClass('focus-visible:outline-2')
      expect(link).toHaveClass('focus-visible:outline-offset-2')
      expect(link).toHaveClass('focus-visible:outline-teal')
    })

    it('link has rounded corners for focus style', () => {
      render(<TopHeader />)

      const link = screen.getByRole('link')
      expect(link).toHaveClass('rounded')
    })

    it('brand text has semantic meaning', () => {
      render(<TopHeader />)

      const text = screen.getByText('ForteCommerce')
      expect(text).toBeInTheDocument()
    })
  })

  describe('brand consistency', () => {
    it('renders both logo and text brand elements', () => {
      const { container } = render(<TopHeader />)

      const images = container.querySelectorAll('img')
      const brandText = screen.getByText('ForteCommerce')

      expect(images.length).toBeGreaterThan(0)
      expect(brandText).toBeInTheDocument()
    })

    it('brand elements are within the same link', () => {
      render(<TopHeader />)

      const link = screen.getByRole('link')
      const brandText = screen.getByText('ForteCommerce')

      expect(link.contains(brandText)).toBe(true)
    })

    it('brand text is bold', () => {
      render(<TopHeader />)

      const text = screen.getByText('ForteCommerce')
      expect(text).toHaveClass('font-bold')
    })
  })
})
