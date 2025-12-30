export const waitForAuthRedirection = async (page) => {
  await page.waitForURL((url) => {
    const pathname = new URL(url).pathname;
    return ['/', '/dashboard'].includes(pathname);
  });
};

export class LoginPage {
  constructor(page) {
    this.page = page;
    this.emailInput = this.page.locator('input[name="identifier"]');
    this.passwordInput = this.page.locator('input[name="password"]');
    const submitButton = page.locator('button[name="method"][value="password"]');
    this.loginButton = submitButton;
  }

  async navigateToLogin() {
    await this.page.goto('/login');
  }

  async login(email, password) {
    if (!email) {
      throw new Error('Email is required for login');
    }
    if (!password) {
      throw new Error('Password is required for login');
    }

    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.loginButton.click();
  }

  async waitForRedirection() {
    await waitForAuthRedirection(this.page);
  }
}
