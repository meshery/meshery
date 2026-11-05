import { expect } from '@playwright/test';

export class DeployModal {
  constructor(page) {
    this.page = page;
    this.modalContent = page.getByTestId('stepper-content');
    this.nextButton = page.getByTestId('deploy-stepper-next-btn');
    this.checkboxBypassDryRun = page.getByTestId('bypass-dry-run-checkbox');
    this.loader = page.getByTestId('loading');
  }

  getStepLabel(index) {
    return this.page.getByTestId(`step-label-${index}`);
  }

  getEnvironmentSelect(envKey) {
    return this.page.getByTestId(`env-${envKey}`);
  }

  async expectStepToBeVisible(index) {
    const step = this.getStepLabel(index);
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
    await envSelector.click();
  }

  async waitForLoader() {
    await this.loader.waitFor({ state: 'detached' });
  }

  async deployWithBypass() {
    const errorMessage = this.page.getByTestId('dry-run-summary-errors');
    if (await errorMessage.isVisible()) {
      await this.toggleBypassDryRun();
    }
    await this.goToNextStep();
  }
}
