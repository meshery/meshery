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
    this.loginButton = this.page.locator('button[name="method"][value="password"]');
  }

  async navigateToLogin() {
    await this.page.goto('/login');
  }

  async loginWithToken(token, baseURL, provider = 'Layer5') {
    if (!token) {
      throw new Error('Token is required for token-based authentication');
    }
    if (!baseURL) {
      throw new Error('Base URL is required for token-based authentication');
    }

    await this.page.context().addCookies([
      {
        name: 'token',
        value: token,
        url: baseURL,
        expires: Math.floor(Date.now() / 1000) + 60 * 60,
      },
      {
        name: 'meshery-provider',
        value: provider,
        url: baseURL,
        expires: -1,
      },
    ]);

    // Navigate to the baseURL after setting cookies
    await this.page.goto(baseURL);
  }

  async loginWithEmail(email, password) {
    if (!email) {
      throw new Error('Email is required for login');
    }
    if (!password) {
      throw new Error('Password is required for login');
    }

    await this.emailInput.evaluate((el, value) => {
      el.value = value;
    }, email);
    await this.passwordInput.evaluate((el, value) => {
      el.value = value;
    }, password);
    await this.loginButton.click();
  }

  async waitForRedirection() {
    await waitForAuthRedirection(this.page);
  }
}
