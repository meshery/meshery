import path from 'path';

export class ImportModal {
  constructor(page) {
    this.page = page;
    this.modal = page.locator('[role="dialog"]').filter({ hasText: 'Import Design' });
    this.designFileInput = this.modal
      .getByLabel(/design file name/i)
      .or(
        this.modal.locator('#root_name, input[placeholder="Untitled Design"], input[name="name"]'),
      )
      .first();
    this.urlInput = this.modal
      .getByLabel(/^url/i)
      .or(
        this.modal.locator(
          '#root_url, input[type="url"], input[name="url"], input[placeholder*="http"]',
        ),
      )
      .first();
    this.fileInput = this.modal.locator('input[type="file"]').first();
    this.importBtn = this.modal
      .getByRole('button', { name: /^import$/i })
      .or(this.modal.getByText(/^import$/i))
      .first();
  }

  async clickImportType(type) {
    const radioLabel = type === 'File' ? /File Upload/i : /URL Import/i;
    const radio = this.modal
      .getByRole('radio', { name: radioLabel })
      .or(this.modal.getByLabel(radioLabel))
      .or(this.modal.getByText(radioLabel))
      .or(
        this.modal.locator(
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
