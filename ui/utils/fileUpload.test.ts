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

  it('prefers the selected file inside the active dialog when present', () => {
    document.body.innerHTML = `
      <div role="dialog">
        <input type="file" id="dialog-input" />
      </div>
      <input type="file" id="page-input" />
    `;

    const dialogInput = document.getElementById('dialog-input') as HTMLInputElement;
    const pageInput = document.getElementById('page-input') as HTMLInputElement;
    const dialogFile = new File(['dialog'], 'dialog.yaml', { type: 'application/yaml' });
    const pageFile = new File(['page'], 'page.yaml', { type: 'application/yaml' });

    selectFile(dialogInput, dialogFile);
    selectFile(pageInput, pageFile);

    expect(findSelectedFileInDialog()).toBe(dialogFile);
  });
});
