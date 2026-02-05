import { expect, Page } from '@playwright/test';

export const waitForSnackBar = async (page: Page, message: string): Promise<void> => {
  // Now TypeScript know 'page' has a .locator() method;
  const snackbar = page.locator(`text=${message}`).first();
  await expect(snackbar).toBeVisible();

  // Wait for the notification to dissappear automatically
  await snackbar.waitFor({state: 'detached', timeout: 10000});
};


