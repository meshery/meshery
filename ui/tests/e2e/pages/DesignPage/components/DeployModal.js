import { expect } from '@playwright/test';

export class DeployModal {
  constructor(page) {
    this.page = page;
    this.modalContent = page.locator('[data-testid="stepper-content"], [role="dialog"]').first();
    this.nextButton = page.getByTestId('deploy-stepper-next-btn');
    this.checkboxBypassDryRun = page
      .locator('[data-testid="bypass-dry-run"], [data-testid="bypass-dry-run-checkbox"]')
      .first();
    this.loader = page.getByTestId('loading').first();
  }

  getStepLabel(index, name) {
    if (name) {
      return this.page.getByText(name, { exact: false });
    }
    return this.page
      .locator(`[data-testid="step-label-${index}"], [data-testid="step-${index}"]`)
      .first();
  }

  getEnvironmentSelect(envKey) {
    return this.page
      .getByTestId(`env-${envKey}`)
      .or(this.page.locator(`[data-testid*="env-"]`).filter({ hasText: envKey }))
      .first();
  }

  async expectStepToBeVisible(index, name) {
    const step = this.getStepLabel(index, name);
    await expect(step).toBeVisible();
  }

  async goToNextStep() {
    await this.nextButton.click();
  }

  async toggleBypassDryRun() {
    await this.checkboxBypassDryRun.click();
  }

  async selectEnvironment(envKey) {
    const envSelector = this.getEnvironmentSelect(envKey);
    if (await envSelector.isVisible()) {
      await envSelector.click();
    }
  }

  async waitForLoader() {
    if (await this.loader.isVisible()) {
      await this.loader.waitFor({ state: 'detached', timeout: 15000 });
    }
  }

  async deployWithBypass() {
    const errorMessage = this.page.getByTestId('dry-run-summary-errors');
    if (await errorMessage.isVisible()) {
      await this.toggleBypassDryRun();
    }
    await this.goToNextStep();
  }
}
