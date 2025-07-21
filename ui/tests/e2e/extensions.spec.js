import { test, expect } from '@playwright/test';
import { ExtensionsPage } from '../pages/extensionsPage.js';

const URLS = {
  KANVAS: {
    DOCS: 'https://docs.layer5.io/kanvas/',
    DESIGNER_EMBED: 'https://docs.layer5.io/kanvas/designer/embedding-designs/',
  },
  DOCKER: {
    EXTENSION: 'https://hub.docker.com/extensions/meshery/docker-extension-meshery',
  },
  MESHERY: {
    CATALOG: 'https://meshery.io/catalog',
    ADATPER_DOCS: 'https://docs.meshery.io/concepts/architecture/adapters',
  },
};

test.describe('Extensions Section Tests', () => {
  let extensionsPage;

  test.beforeEach(async ({ page }) => {
    extensionsPage = new ExtensionsPage(page);
    await extensionsPage.navigate();
  });

  test('Verify Kanvas Snapshot using data-testid', async () => {
    const { heading, description, enableBtn, image } = extensionsPage.getKanvasSnapshotElements();
    await expect(heading).toBeVisible();
    await expect(description).toBeVisible();
    await expect(enableBtn).toBeVisible();
    await expect(enableBtn).toBeEnabled();
    await expect(image).toBeVisible();
  });

  test('Verify Performance Analysis Details', async () => {
    const { heading, enableBtn } = extensionsPage.getPerformanceAnalysisElements();
    await expect(heading).toBeVisible();
    await expect(enableBtn).toBeVisible();
    await expect(enableBtn).toBeEnabled();
  });

  test('Verify Kanvas Details', async ({ context }) => {
    const { heading, button } = extensionsPage.getKanvasSignupElements();
    await expect(heading).toBeVisible();
    await expect(button).toBeVisible();
    if (await button.isEnabled()) {
      const docsPage = await extensionsPage.openKanvasDocs(context);
      try {
        await expect(docsPage).toHaveURL(URLS.KANVAS.DOCS);
      } finally {
        await docsPage.close();
      }
    }
  });

  test('Verify Meshery Docker Extension Details', async ({ context }) => {
    const newPage = await extensionsPage.openDockerExtension(context);
    try {
      await expect(newPage).toHaveURL(URLS.DOCKER.EXTENSION);
    } finally {
      await newPage.close();
    }
  });

  test('Verify Meshery Design Embed Details', async ({ context }) => {
    const newPage = await extensionsPage.openDesignEmbed(context);
    try {
      await expect(newPage).toHaveURL(URLS.KANVAS.DESIGNER_EMBED);
    } finally {
      await newPage.close();
    }
  });

  test('Verify Meshery Catalog Section Details', async ({ context }) => {
    await extensionsPage.toggleCatalogSection();
    const newPage = await extensionsPage.openCatalogLink(context);
    try {
      await expect(newPage).toHaveURL(URLS.MESHERY.CATALOG);
    } finally {
      await newPage.close();
    }
  });

  test('Verify Meshery Adapter for Istio Section', async ({ context }) => {
    const docsPage = await extensionsPage.openAdapterIstioDocs(context);
    try {
      await expect(docsPage).toHaveURL(URLS.MESHERY.ADATPER_DOCS);
    } finally {
      await docsPage.close();
    }
  });
});
