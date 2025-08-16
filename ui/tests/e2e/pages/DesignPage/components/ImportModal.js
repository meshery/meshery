import path from 'path';

export class ImportModal {
  constructor(page) {
    this.page = page;
    this.modal = page.getByTestId('import-design-modal');
    this.designFileInput = page.getByLabel('Design file name');
    this.urlInput = page.getByLabel('URL *');
    this.fileChoose = page.getByRole('textbox', { name: 'file', exact: true });
    this.importBtn = page.getByRole('button', { name: 'Import' });
  }

  async clickImportType(type) {
    await this.page.getByRole('radio', { name: type }).click();
  }

  importDesign(type, path, designName) {
    if (type === 'File') {
      return this._importDesignFromFile(path, designName);
    } else if (type === 'URL') {
      return this._importDesignFromURL(path, designName);
    } else {
      throw new Error('Invalid import type');
    }
  }

  async _importDesignFromFile(filePath, designName) {
    await this.designFileInput.fill(designName);
    const fileChooserPromoise = this.page.waitForEvent('filechooser');
    await this.fileChoose.click();
    const fileChooser = await fileChooserPromoise;
    await fileChooser.setFiles(path.resolve(filePath));
    await this.importBtn.click();
  }

  async _importDesignFromURL(url, designName) {
    await this.urlInput.fill(url);
    await this.designFileInput.fill(designName);
    await this.importBtn.click();
  }
}
