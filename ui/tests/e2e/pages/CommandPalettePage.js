import { expect } from '@playwright/test';

export class CommandPalettePage {
  constructor(page) {
    this.page = page;
    this.trigger = page.getByTestId('command-palette-trigger');
    this.palette = page.getByTestId('command-palette');
    this.input = page.getByTestId('command-palette-input');
  }

  async openViaTrigger() {
    await this.trigger.click();
    await expect(this.palette).toBeVisible();
    await expect(this.input).toBeFocused();
  }

  async openViaShortcut() {
    await this.page.keyboard.press('ControlOrMeta+k');
    await expect(this.palette).toBeVisible();
    await expect(this.input).toBeFocused();
  }

  async closeViaEscape() {
    await this.page.keyboard.press('Escape');
    await expect(this.palette).toBeHidden();
  }

  async closeViaShortcut() {
    await this.page.keyboard.press('ControlOrMeta+k');
    await expect(this.palette).toBeHidden();
  }

  getItem(id) {
    return this.page.getByTestId(`command-palette-item-${id}`);
  }
}
