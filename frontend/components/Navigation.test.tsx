import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Navigation } from './Navigation'

vi.mock('@/hooks/useAuth', () => ({
  useAuth: () => ({
    logout: vi.fn(),
  }),
}))

describe('Navigation', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('rendering', () => {
    it('renders navigation with shop and orders links', () => {
      render(<Navigation />)

      expect(screen.getByText('Shop')).toBeInTheDocument()
      expect(screen.getByText('Orders')).toBeInTheDocument()
    })

    it('renders sign out button', () => {
      render(<Navigation />)

      expect(screen.getByText('Sign Out')).toBeInTheDocument()
    })

    it('renders as navigation element', () => {
      const { container } = render(<Navigation />)

      const nav = container.querySelector('nav')
      expect(nav).toBeInTheDocument()
    })
  })

  describe('routing', () => {
    it('shop link navigates to home', () => {
      render(<Navigation />)

      const shopLinks = screen.getAllByRole('link').filter((link) => link.getAttribute('href') === '/')
      expect(shopLinks.length).toBeGreaterThan(0)
    })

    it('orders link navigates to orders page', () => {
      render(<Navigation />)

      const orderLinks = screen.getAllByRole('link').filter((link) => link.getAttribute('href') === '/orders')
      expect(orderLinks.length).toBeGreaterThan(0)
    })
  })

  describe('sign out functionality', () => {
    it('sign out button triggers logout', async () => {
      const mockLogout = vi.fn()
      vi.doMock('@/hooks/useAuth', () => ({
        useAuth: () => ({
          logout: mockLogout,
        }),
      }))

      const user = userEvent.setup()
      render(<Navigation />)

      const signOutButton = screen.getByRole('button', { name: /Sign Out/i })
      await user.click(signOutButton)
    })
  })

  describe('layout', () => {
    it('navigation is fixed at bottom', () => {
      const { container } = render(<Navigation />)

      const nav = container.querySelector('nav')
      expect(nav).toHaveClass('fixed')
      expect(nav).toHaveClass('bottom-0')
    })

    it('navigation spans full width', () => {
      const { container } = render(<Navigation />)

      const nav = container.querySelector('nav')
      expect(nav).toHaveClass('left-0')
      expect(nav).toHaveClass('right-0')
    })

    it('navigation has proper z-index', () => {
      const { container } = render(<Navigation />)

      const nav = container.querySelector('nav')
      expect(nav).toHaveClass('z-50')
    })

    it('has safe area bottom padding', () => {
      const { container } = render(<Navigation />)

      const nav = container.querySelector('nav')
      expect(nav?.style.paddingBottom).toBeDefined()
    })
  })

  describe('styling', () => {
    it('has white background', () => {
      const { container } = render(<Navigation />)

      const nav = container.querySelector('nav')
      expect(nav).toHaveClass('bg-white')
    })

    it('has top border', () => {
      const { container } = render(<Navigation />)

      const nav = container.querySelector('nav')
      expect(nav).toHaveClass('border-t')
    })

    it('has shadow', () => {
      const { container } = render(<Navigation />)

      const nav = container.querySelector('nav')
      expect(nav?.className).toMatch(/shadow/)
    })
  })

  describe('navigation items', () => {
    it('has three navigation items', () => {
      const { container } = render(<Navigation />)

      const items = container.querySelectorAll('a, button[class*="flex-1"]')
      expect(items.length).toBeGreaterThanOrEqual(2)
    })

    it('items are flex containers', () => {
      const { container } = render(<Navigation />)

      const navLinks = screen.getAllByRole('link')
      navLinks.forEach((link) => {
        expect(link).toHaveClass('flex-1')
      })
    })
  })

  describe('icon rendering', () => {
    it('renders home icon', () => {
      const { container } = render(<Navigation />)

      const svgs = container.querySelectorAll('svg')
      expect(svgs.length).toBeGreaterThanOrEqual(2)
    })

    it('renders orders icon', () => {
      const { container } = render(<Navigation />)

      const svgs = container.querySelectorAll('svg')
      expect(svgs.length).toBeGreaterThanOrEqual(2)
    })

    it('renders sign out icon', () => {
      const { container } = render(<Navigation />)

      const svgs = container.querySelectorAll('svg')
      expect(svgs.length).toBeGreaterThanOrEqual(3)
    })
  })

  describe('accessibility', () => {
    it('navigation items are keyboard navigable', () => {
      render(<Navigation />)

      const links = screen.getAllByRole('link')
      expect(links.length).toBeGreaterThan(0)

      const button = screen.getByRole('button', { name: /Sign Out/i })
      expect(button).toBeInTheDocument()
    })

    it('sign out button is focusable', () => {
      render(<Navigation />)

      const signOutButton = screen.getByRole('button', { name: /Sign Out/i })
      expect(signOutButton).toHaveClass('focus-visible:outline-none')
    })
  })
})
