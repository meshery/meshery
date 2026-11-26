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

import {
  filenameFromDataUrl,
  decodeDesignDataUrlToBytes,
  resolveImportedDesignFile,
} from './import-design-file';

describe('filenameFromDataUrl', () => {
  it('extracts and decodes filename enclosed with semicolons', () => {
    expect(filenameFromDataUrl('data:application/yaml;name=my%20design.yaml;base64,QQ==')).toBe(
      'my design.yaml',
    );
  });

  it('extracts filename preceding the comma data separator without a trailing semicolon', () => {
    expect(filenameFromDataUrl('data:application/yaml;base64,QQ==;name=design.yaml')).toBe(
      'design.yaml',
    );
    expect(filenameFromDataUrl('data:application/x-yaml;name=app.yaml,content')).toBe('app.yaml');
  });

  it('gracefully falls back to raw string when URI decoding fails on malformed percent encoding', () => {
    expect(filenameFromDataUrl('data:application/yaml;name=malformed%20%2.yaml;base64,QQ==')).toBe(
      'malformed%20%2.yaml',
    );
  });

  it('returns undefined for nullish, non-string, or dataUrl lacking name', () => {
    expect(filenameFromDataUrl(undefined)).toBeUndefined();
    expect(filenameFromDataUrl('')).toBeUndefined();
    // @ts-expect-error testing invalid argument
    expect(filenameFromDataUrl(null)).toBeUndefined();
    expect(filenameFromDataUrl('data:text/plain;base64,QQ==')).toBeUndefined();
  });
});

describe('decodeDesignDataUrlToBytes', () => {
  it('returns decoded byte array when getUnit8ArrayDecodedFile succeeds', () => {
    getUnit8ArrayDecodedFile.mockReturnValue([65, 66, 67]);
    expect(decodeDesignDataUrlToBytes('data:text/plain;base64,QUJD')).toEqual([65, 66, 67]);
  });

  it('returns null when getUnit8ArrayDecodedFile throws on invalid base64 or malformed URL', () => {
    getUnit8ArrayDecodedFile.mockImplementation(() => {
      throw new Error('Invalid base64 string');
    });
    expect(decodeDesignDataUrlToBytes('data:text/plain;base64,invalid@@')).toBeNull();
  });

  it('returns null for empty or non-string inputs without throwing', () => {
    expect(decodeDesignDataUrlToBytes(undefined)).toBeNull();
    expect(decodeDesignDataUrlToBytes('')).toBeNull();
    // @ts-expect-error testing invalid argument
    expect(decodeDesignDataUrlToBytes(null)).toBeNull();
  });
});

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

  it('uses filename extracted from dataUrl when present', async () => {
    getUnit8ArrayDecodedFile.mockReturnValue([1, 2, 3]);
    findSelectedFileInDialog.mockReturnValue(undefined);

    await expect(
      resolveImportedDesignFile('data:application/yaml;name=custom.yaml;base64,QQ=='),
    ).resolves.toEqual({
      fileData: [1, 2, 3],
      fileName: 'custom.yaml',
    });
  });

  it('falls back to selectedFile when dataUrl decoding returns null', async () => {
    getUnit8ArrayDecodedFile.mockReturnValue(null);
    const mockFile = new File(['content'], 'picked.yaml', { type: 'application/yaml' });
    findSelectedFileInDialog.mockReturnValue(mockFile);
    readFileAsBytes.mockResolvedValue([10, 20, 30]);

    await expect(resolveImportedDesignFile(undefined)).resolves.toEqual({
      fileData: [10, 20, 30],
      fileName: 'picked.yaml',
    });
  });

  it('safely returns null when readFileAsBytes rejects', async () => {
    getUnit8ArrayDecodedFile.mockReturnValue(null);
    const mockFile = new File(['content'], 'error.yaml', { type: 'application/yaml' });
    findSelectedFileInDialog.mockReturnValue(mockFile);
    readFileAsBytes.mockRejectedValue(new Error('Read failed'));

    await expect(resolveImportedDesignFile(undefined)).resolves.toBeNull();
  });

  it('returns null when no dataUrl bytes or selectedFile can be resolved', async () => {
    getUnit8ArrayDecodedFile.mockReturnValue(null);
    findSelectedFileInDialog.mockReturnValue(undefined);

    await expect(resolveImportedDesignFile(undefined)).resolves.toBeNull();
  });
});
