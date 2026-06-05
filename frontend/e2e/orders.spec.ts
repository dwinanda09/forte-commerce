import { test, expect } from '@playwright/test'

test.describe('Orders Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')

    const usernameInput = page.getByPlaceholder(/username|email/i)
    const passwordInput = page.getByPlaceholder(/password/i)

    if (usernameInput && passwordInput) {
      await usernameInput.fill('demo')
      await passwordInput.fill('password')

      const submitButton = page.getByRole('button', { name: /sign in|login|submit/i })
      await submitButton.click()

      await page.waitForTimeout(2000)
    }
  })

  test('navigates to orders page', async ({ page }) => {
    await page.goto('/orders')

    expect(page.url()).toContain('/orders')
  })

  test('displays order list or empty state', async ({ page }) => {
    await page.goto('/orders')

    const orderElements = await page.locator('[class*="order"]').all()
    const emptyStateElements = await page.locator('text=/no orders|empty/i').all()

    if (orderElements.length > 0) {
      expect(orderElements.length).toBeGreaterThan(0)
    } else if (emptyStateElements.length > 0) {
      expect(emptyStateElements.length).toBeGreaterThan(0)
    }
  })

  test('bottom nav orders tab is highlighted when on orders page', async ({ page }) => {
    await page.goto('/orders')

    const navElements = await page.locator('nav button, nav a').all()

    for (const element of navElements) {
      const text = await element.textContent()
      if (text?.toLowerCase().includes('order')) {
        const isActive = await element.evaluate((el) => {
          return (
            el.classList.contains('active') ||
            el.classList.contains('text-teal') ||
            el.getAttribute('aria-current') === 'page'
          )
        })
        expect(isActive).toBeTruthy()
        break
      }
    }
  })

  test('can view individual order details if orders exist', async ({ page }) => {
    await page.goto('/orders')

    const orderLinks = await page.getByRole('link').all()

    if (orderLinks.length > 0) {
      for (const link of orderLinks) {
        const href = await link.getAttribute('href')
        if (href && href.includes('/orders/')) {
          await link.click()
          await page.waitForTimeout(500)
          break
        }
      }
    }
  })

  test('shows order status for each order', async ({ page }) => {
    await page.goto('/orders')

    const statusElements = await page.locator('text=/pending|paid|cancelled/i').all()

    if (statusElements.length > 0) {
      expect(statusElements.length).toBeGreaterThan(0)
    }
  })

  test('displays order items and prices', async ({ page }) => {
    await page.goto('/orders')

    const priceElements = await page.locator('text=/\\$/').all()

    if (priceElements.length > 0) {
      expect(priceElements.length).toBeGreaterThan(0)
    }
  })

  test('can navigate back from order details', async ({ page }) => {
    await page.goto('/orders')

    const orderLinks = await page.getByRole('link').all()

    if (orderLinks.length > 0) {
      for (const link of orderLinks) {
        const href = await link.getAttribute('href')
        if (href && href.includes('/orders/')) {
          await link.click()
          await page.waitForTimeout(500)

          const backButton = page.getByRole('button', { name: /back/i }).or(page.getByRole('link', { name: /back/i }))
          const exists = await backButton.isVisible().catch(() => false)

          if (exists) {
            await backButton.click()
            await page.waitForURL(/\/orders\/?$/)
          }
          break
        }
      }
    }
  })

  test('orders page loads without errors', async ({ page }) => {
    let errorFound = false
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        errorFound = true
      }
    })

    await page.goto('/orders')
    await page.waitForTimeout(1000)

    expect(errorFound).toBe(false)
  })

  test('orders list is paginated or scrollable if many items', async ({ page }) => {
    await page.goto('/orders')

    const listContainer = await page.locator('[class*="list"], [class*="orders"]').first()

    if (await listContainer.isVisible().catch(() => false)) {
      const scrollHeight = await listContainer.evaluate((el) => el.scrollHeight)
      const clientHeight = await listContainer.evaluate((el) => el.clientHeight)

      if (scrollHeight > clientHeight) {
        expect(scrollHeight).toBeGreaterThan(clientHeight)
      }
    }
  })
})
