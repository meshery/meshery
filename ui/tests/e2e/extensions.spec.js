const { test, expect } = require('@playwright/test');
const { ExtensionsPage } = require('../pages/ExtensionsPage');

test('Sapna should be able to add an extension', async ({ page }) => {
  const extensionsPage = new ExtensionsPage(page);

  await extensionsPage.goto();
  await extensionsPage.addExtension('TestExtension');

  const extensions = await extensionsPage.getExtensionList();
  expect(extensions).toContain('TestExtension');
});
