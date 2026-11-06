import { expect, test } from '@playwright/test';
import { CommandPalettePage } from './pages/CommandPalettePage';
import { DashboardPage } from './pages/DashboardPage';

test.describe('Command Palette E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    const dashboardPage = new DashboardPage(page);
    await dashboardPage.navigateToDashboard();
    await expect(page.getByTestId('command-palette-trigger')).toBeVisible();
  });

  test('opens via header button with input focused and closes on single Escape', async ({
    page,
  }) => {
    const palettePage = new CommandPalettePage(page);

    await palettePage.trigger.click();
    await expect(palettePage.palette).toBeVisible();
    await expect(palettePage.input).toBeFocused();

    await page.keyboard.press('Escape');
    await expect(palettePage.palette).toBeHidden();
  });

  test('opens via shortcut with input focused and toggles closed on shortcut again', async ({
    page,
  }) => {
    const palettePage = new CommandPalettePage(page);

    await page.keyboard.press('ControlOrMeta+k');
    await expect(palettePage.palette).toBeVisible();
    await expect(palettePage.input).toBeFocused();

    await page.keyboard.press('ControlOrMeta+k');
    await expect(palettePage.palette).toBeHidden();
  });

  test('highlights initial active option and updates active option with ArrowDown and ArrowUp', async ({
    page,
  }) => {
    const palettePage = new CommandPalettePage(page);
    await palettePage.openViaShortcut();

    // Active row on unfiltered list: initial highlight is first option (index 0: action:create-design)
    await expect(palettePage.input).toHaveAttribute('aria-activedescendant', /-option-0$/);
    const activeId0 = await palettePage.input.getAttribute('aria-activedescendant');
    expect(activeId0).toBeTruthy();
    await expect(palettePage.getItem('action:create-design')).toHaveAttribute('id', activeId0!);
    await expect(page.locator(`[id="${activeId0}"]`)).toHaveAttribute(
      'data-testid',
      'command-palette-item-action:create-design',
    );

    // ArrowDown navigates to second option (index 1: action:import-design)
    await page.keyboard.press('ArrowDown');
    await expect(palettePage.input).toHaveAttribute('aria-activedescendant', /-option-1$/);
    const activeId1 = await palettePage.input.getAttribute('aria-activedescendant');
    expect(activeId1).toBeTruthy();
    await expect(palettePage.getItem('action:import-design')).toHaveAttribute('id', activeId1!);
    await expect(page.locator(`[id="${activeId1}"]`)).toHaveAttribute(
      'data-testid',
      'command-palette-item-action:import-design',
    );

    // ArrowUp navigates back to first option (index 0: action:create-design)
    await page.keyboard.press('ArrowUp');
    await expect(palettePage.input).toHaveAttribute('aria-activedescendant', /-option-0$/);
    const activeIdBack = await palettePage.input.getAttribute('aria-activedescendant');
    expect(activeIdBack).toBeTruthy();
    await expect(palettePage.getItem('action:create-design')).toHaveAttribute('id', activeIdBack!);
    await expect(page.locator(`[id="${activeIdBack}"]`)).toHaveAttribute(
      'data-testid',
      'command-palette-item-action:create-design',
    );
  });

  test('filters commands by query and hides non-matching items', async ({ page }) => {
    const palettePage = new CommandPalettePage(page);
    await palettePage.openViaShortcut();

    await palettePage.input.fill('des');

    // Matching items visible
    await expect(palettePage.getItem('nav:/configuration/designs')).toBeVisible();
    await expect(palettePage.getItem('action:create-design')).toBeVisible();
    await expect(palettePage.getItem('action:import-design')).toBeVisible();

    // Dashboard navigation is not present
    await expect(palettePage.getItem('nav:/')).toHaveCount(0);
  });

  test('navigates to page on Enter and verifies palette closes first', async ({ page }) => {
    const palettePage = new CommandPalettePage(page);
    await palettePage.openViaShortcut();

    await palettePage.input.fill('Designs');
    await expect(palettePage.getItem('nav:/configuration/designs')).toBeVisible();
    await page.keyboard.press('Enter');

    // Palette is gone FIRST
    await expect(palettePage.palette).toBeHidden();

    // URL and document title match expected values
    await expect(page).toHaveURL(/\/configuration\/designs/);
    await expect(page).toHaveTitle(/Designs \| Meshery/);
  });

  test('navigates to import design modal and cleans up action query param', async ({ page }) => {
    const palettePage = new CommandPalettePage(page);

    await page.keyboard.press('ControlOrMeta+k');
    await expect(palettePage.palette).toBeVisible();

    await palettePage.input.fill('Import Design');
    await expect(palettePage.getItem('action:import-design')).toBeVisible();
    await page.keyboard.press('Enter');

    // Palette is gone FIRST
    await expect(palettePage.palette).toBeHidden();

    // getByRole('dialog') containing textbox "Design file name" is visible
    const importDialog = page
      .getByRole('dialog')
      .filter({ has: page.getByRole('textbox', { name: 'Design file name' }) });
    await expect(importDialog).toBeVisible();

    // URL path is /configuration/designs and does not contain action=import
    await expect(page).toHaveURL(/\/configuration\/designs/);
    await expect(page).not.toHaveURL(/action=import/);
  });

  test('opens command palette when shortcut is pressed inside a text input on the Designs page', async ({
    page,
  }) => {
    const palettePage = new CommandPalettePage(page);

    await page.goto('/configuration/designs', { waitUntil: 'domcontentloaded' });
    await expect(palettePage.trigger).toBeVisible();

    // Locate real stable text input on the Designs page toolbar
    const searchBar = page.getByTestId('meshery-patterns-search-bar');
    const searchInput = page.getByPlaceholder('Search designs...');

    await expect(async () => {
      if (!(await searchInput.isVisible())) {
        await searchBar.getByRole('button').first().click();
      }
      await expect(searchInput).toBeVisible({ timeout: 2_000 });
    }).toPass({ timeout: 30_000 });

    await searchInput.focus();
    await expect(searchInput).toBeFocused();

    // Ctrl+K while focus is inside the text input toggles command palette
    await page.keyboard.press('ControlOrMeta+k');
    await expect(palettePage.palette).toBeVisible();
    await expect(palettePage.input).toBeFocused();

    // Ensure the page's search input value was NOT changed (no 'k' inserted)
    await expect(searchInput).toHaveValue('');
  });
});
