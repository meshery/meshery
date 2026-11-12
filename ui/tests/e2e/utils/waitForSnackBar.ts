import { expect, Page } from '@playwright/test';

export const waitForSnackBar = async (page: Page, message: string): Promise<void> => {
  const snackbar = page.locator(`text=${message}`).first();
  await expect(snackbar).toBeVisible();
  await snackbar.waitFor({ state: 'detached', timeout: 10000 });
};
