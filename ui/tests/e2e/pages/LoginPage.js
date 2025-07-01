export class LoginPage {
  constructor(page) {
    this.page = page;
    this.emailInput = this.page.getByLabel('E-Mail');
    this.passwordInput = this.page.getByLabel('Password');
    this.loginButton = this.page.getByRole('button', { name: 'Sign in', exact: true });
  }

  async navigateToLogin() {
    await this.page.goto('/login');
  }

  async login(email, password) {
    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.loginButton.click();
  }

  async waitForRedirection() {
    await this.page.waitForURL(url => {
      const pathname = new URL(url).pathname;
      return ['/', '/dashboard'].includes(pathname);
    });
  }
  
}
