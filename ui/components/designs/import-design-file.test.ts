import { beforeEach, describe, expect, it, vi } from 'vitest';

const getUnit8ArrayDecodedFile = vi.fn();
const findSelectedFileInDialog = vi.fn();
const readFileAsBytes = vi.fn();

vi.mock('@/utils/utils', () => ({
  getUnit8ArrayDecodedFile: (dataUrl: string) => getUnit8ArrayDecodedFile(dataUrl),
}));

vi.mock('@/utils/fileUpload', () => ({
  findSelectedFileInDialog: (selector: string) => findSelectedFileInDialog(selector),
  readFileAsBytes: (file: File) => readFileAsBytes(file),
}));

import { resolveImportedDesignFile } from './import-design-file';

describe('resolveImportedDesignFile', () => {
  beforeEach(() => {
    getUnit8ArrayDecodedFile.mockReset();
    findSelectedFileInDialog.mockReset();
    readFileAsBytes.mockReset();
  });

  it('falls back to a default filename when bytes are present but the original name is missing', async () => {
    getUnit8ArrayDecodedFile.mockReturnValue([1, 2, 3]);
    findSelectedFileInDialog.mockReturnValue(undefined);

    await expect(resolveImportedDesignFile('data:text/plain;base64,QQ==')).resolves.toEqual({
      fileData: [1, 2, 3],
      fileName: 'design.yaml',
    });
  });
});
