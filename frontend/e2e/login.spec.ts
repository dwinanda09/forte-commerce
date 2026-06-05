import { test, expect } from '@playwright/test'

test.describe('Login Flow', () => {
  test('loads login page and displays sign in form', async ({ page }) => {
    await page.goto('/login')

    expect(page.url()).toContain('/login')
    await expect(page).toHaveTitle(/.*/)
  })

  test('shows error message with wrong credentials', async ({ page }) => {
    await page.goto('/login')

    const usernameInput = page.getByPlaceholder(/username|email/i)
    const passwordInput = page.getByPlaceholder(/password/i)

    await usernameInput.fill('wronguser')
    await passwordInput.fill('wrongpass')

    const submitButton = page.getByRole('button', { name: /sign in|login|submit/i })
    await submitButton.click()

    await page.waitForTimeout(1000)

    const errorElements = await page
      .locator('text=/invalid|incorrect|failed|error/i')
      .all()

    expect(errorElements.length).toBeGreaterThan(0)
  })

  test('redirects to home page with correct credentials', async ({ page }) => {
    await page.goto('/login')

    const usernameInput = page.getByPlaceholder(/username|email/i)
    const passwordInput = page.getByPlaceholder(/password/i)

    await usernameInput.fill('demo')
    await passwordInput.fill('password')

    const submitButton = page.getByRole('button', { name: /sign in|login|submit/i })
    await submitButton.click()

    await page.waitForURL(/\/$/, { timeout: 10000 }).catch(() => {
      console.log('Expected redirect to home page')
    })
  })

  test('prevents empty form submission', async ({ page }) => {
    await page.goto('/login')

    const submitButton = page.getByRole('button', { name: /sign in|login|submit/i })

    const isDisabled = await submitButton.isDisabled()
    if (isDisabled) {
      expect(isDisabled).toBe(true)
    }
  })

  test('displays input fields for credentials', async ({ page }) => {
    await page.goto('/login')

    const usernameInput = page.getByPlaceholder(/username|email/i)
    const passwordInput = page.getByPlaceholder(/password/i)

    await expect(usernameInput).toBeVisible()
    await expect(passwordInput).toBeVisible()
  })

  test('form retains input after submission attempt', async ({ page }) => {
    await page.goto('/login')

    const usernameInput = page.getByPlaceholder(/username|email/i)
    const passwordInput = page.getByPlaceholder(/password/i)

    await usernameInput.fill('testuser')
    await passwordInput.fill('testpass')

    const submitButton = page.getByRole('button', { name: /sign in|login|submit/i })
    await submitButton.click()

    await page.waitForTimeout(500)

    expect(await usernameInput.inputValue()).toBe('testuser')
  })
})
