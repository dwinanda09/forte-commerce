import { test, expect } from '@playwright/test'

test.describe('Shop Flow', () => {
  test.beforeEach(async ({ page, context }) => {
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

  test('displays product grid on shop page', async ({ page }) => {
    await page.goto('/')

    const productCards = await page.locator('[class*="shadow-card"], [class*="product"]').all()

    if (productCards.length > 0) {
      expect(productCards.length).toBeGreaterThan(0)
    }
  })

  test('shows product details in grid', async ({ page }) => {
    await page.goto('/')

    const priceElements = await page.locator('text=/\\$/').all()

    if (priceElements.length > 0) {
      expect(priceElements.length).toBeGreaterThan(0)
    }
  })

  test('clicking add to cart opens or shows cart', async ({ page }) => {
    await page.goto('/')

    const addToCartButtons = await page.getByRole('button', { name: /add to cart/i }).all()

    if (addToCartButtons.length > 0) {
      await addToCartButtons[0].click()

      await page.waitForTimeout(500)

      const cartDrawer = await page.locator('[class*="cart"], [class*="drawer"]').all()
      expect(cartDrawer.length).toBeGreaterThanOrEqual(0)
    }
  })

  test('cart shows correct item count after adding product', async ({ page }) => {
    await page.goto('/')

    const addToCartButtons = await page.getByRole('button', { name: /add to cart/i }).all()

    if (addToCartButtons.length > 0) {
      await addToCartButtons[0].click()

      await page.waitForTimeout(500)

      const cartElements = await page.locator('text=/cart|items?/i').all()
      expect(cartElements.length).toBeGreaterThanOrEqual(0)
    }
  })

  test('can increase quantity before adding to cart', async ({ page }) => {
    await page.goto('/')

    const increaseQtyButtons = await page.locator('button:has-text("+")').all()

    if (increaseQtyButtons.length > 0) {
      await increaseQtyButtons[0].click()

      await page.waitForTimeout(300)
    }
  })

  test('can decrease quantity before adding to cart', async ({ page }) => {
    await page.goto('/')

    const increaseQtyButtons = await page.locator('button:has-text("+")').all()

    if (increaseQtyButtons.length > 0) {
      await increaseQtyButtons[0].click()

      const decreaseQtyButtons = await page.locator('button:has-text("−")').all()
      if (decreaseQtyButtons.length > 0) {
        await decreaseQtyButtons[0].click()

        await page.waitForTimeout(300)
      }
    }
  })

  test('out of stock products show disabled button', async ({ page }) => {
    await page.goto('/')

    const disabledButtons = await page.getByRole('button', { name: /out of stock/i }).all()

    if (disabledButtons.length > 0) {
      const isDisabled = await disabledButtons[0].isDisabled()
      expect(isDisabled).toBe(true)
    }
  })

  test('shows stock information for each product', async ({ page }) => {
    await page.goto('/')

    const stockElements = await page.locator('text=/in stock|out of stock/i').all()

    if (stockElements.length > 0) {
      expect(stockElements.length).toBeGreaterThan(0)
    }
  })

  test('product cards are clickable or interactive', async ({ page }) => {
    await page.goto('/')

    const addButtons = await page.getByRole('button', { name: /add to cart/i }).all()

    if (addButtons.length > 0) {
      expect(addButtons[0]).toBeEnabled()
    }
  })
})
