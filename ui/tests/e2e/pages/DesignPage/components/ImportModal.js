import path from 'path';

export class ImportModal {
  constructor(page) {
    this.page = page;
    this.modal = page.locator('[role="dialog"]').filter({ hasText: 'Import Design' });
    this.designFileInput = page
      .locator('#root_name, input[placeholder="Untitled Design"], input[name="name"]')
      .first();
    this.urlInput = page.locator('#root_url, input[type="url"], input[name="url"]').first();
    this.fileInput = page.locator('input[type="file"]');
    this.importBtn = page.getByRole('button', { name: /^import/i });
  }

  async clickImportType(type) {
    const radioLabel = type === 'File' ? /File Upload/i : /URL Import/i;
    const radio = this.page
      .getByRole('radio', { name: radioLabel })
      .or(
        this.page.locator(
          `input[type="radio"][value="${type === 'File' ? 'File Upload' : 'URL Import'}"]`,
        ),
      );
    if ((await radio.count()) > 0) {
      await radio.first().click();
    }
  }

  async importDesign(type, filePathOrUrl, designName) {
    if (type === 'File') {
      return this._importDesignFromFile(filePathOrUrl, designName);
    } else if (type === 'URL') {
      return this._importDesignFromURL(filePathOrUrl, designName);
    } else {
      throw new Error(`Invalid import type: ${type}`);
    }
  }

  async _importDesignFromFile(filePath, designName) {
    await this.clickImportType('File');
    if (await this.designFileInput.isVisible()) {
      await this.designFileInput.fill(designName);
    }
    const resolvedPath = path.resolve(filePath);
    await this.fileInput.setInputFiles(resolvedPath);
    await this.importBtn.click();
  }

  async _importDesignFromURL(url, designName) {
    await this.clickImportType('URL');
    if (await this.designFileInput.isVisible()) {
      await this.designFileInput.fill(designName);
    }
    await this.urlInput.fill(url);
    await this.importBtn.click();
  }
}
