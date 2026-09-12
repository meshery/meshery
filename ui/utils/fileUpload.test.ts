import { beforeEach, describe, expect, it } from 'vitest';
import { findSelectedFileInDialog } from './fileUpload';

const selectFile = (input: HTMLInputElement, file: File) => {
  Object.defineProperty(input, 'files', {
    value: [file],
    configurable: true,
  });
};

describe('findSelectedFileInDialog', () => {
  beforeEach(() => {
    document.body.innerHTML = '';
  });

  it('falls back to the full document when the first dialog has no selected file', () => {
    document.body.innerHTML = `
      <div role="dialog">
        <input type="file" id="dialog-input" />
      </div>
      <input type="file" id="page-input" />
    `;

    const pageInput = document.getElementById('page-input') as HTMLInputElement;
    const pageFile = new File(['design'], 'design.yaml', { type: 'application/yaml' });
    selectFile(pageInput, pageFile);

    expect(findSelectedFileInDialog()).toBe(pageFile);
  });

  it('prefers the selected file inside the top-most dialog when multiple dialogs are open', () => {
    document.body.innerHTML = `
      <div role="dialog">
        <input type="file" id="background-dialog-input" />
      </div>
      <div role="dialog">
        <input type="file" id="active-dialog-input" />
      </div>
      <input type="file" id="page-input" />
    `;

    const backgroundDialogInput = document.getElementById(
      'background-dialog-input',
    ) as HTMLInputElement;
    const activeDialogInput = document.getElementById('active-dialog-input') as HTMLInputElement;
    const pageInput = document.getElementById('page-input') as HTMLInputElement;
    const backgroundDialogFile = new File(['background'], 'background.yaml', {
      type: 'application/yaml',
    });
    const activeDialogFile = new File(['dialog'], 'dialog.yaml', { type: 'application/yaml' });
    const pageFile = new File(['page'], 'page.yaml', { type: 'application/yaml' });

    selectFile(backgroundDialogInput, backgroundDialogFile);
    selectFile(activeDialogInput, activeDialogFile);
    selectFile(pageInput, pageFile);

    expect(findSelectedFileInDialog()).toBe(activeDialogFile);
  });
});
