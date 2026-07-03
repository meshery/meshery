import { beforeEach, describe, expect, it, vi } from 'vitest';

const resolveImportedDesignFile = vi.fn();

vi.mock('./import-design-file', () => ({
  resolveImportedDesignFile: (dataUrl: string | undefined) => resolveImportedDesignFile(dataUrl),
}));

import { buildImportDesignRequestBody } from './import-design-request';

describe('buildImportDesignRequestBody', () => {
  beforeEach(() => {
    resolveImportedDesignFile.mockReset();
  });

  it('builds the file import request body from resolved file metadata', async () => {
    resolveImportedDesignFile.mockResolvedValue({
      fileData: [1, 2, 3],
      fileName: 'imported-design.yaml',
    });

    await expect(
      buildImportDesignRequestBody({
        uploadType: 'File Upload',
        name: 'Imported design',
        file: 'data:text/plain;base64,QQ==',
      }),
    ).resolves.toEqual({
      requestBody: JSON.stringify({
        name: 'Imported design',
        file_name: 'imported-design.yaml',
        file: [1, 2, 3],
      }),
    });
  });

  it('returns a user-facing error when no design file is available', async () => {
    resolveImportedDesignFile.mockResolvedValue(null);

    await expect(
      buildImportDesignRequestBody({
        uploadType: 'File Upload',
        name: 'Imported design',
      }),
    ).resolves.toEqual({
      errorMessage: 'Please choose a design file before continuing.',
    });
  });

  it('logs and surfaces file resolution failures', async () => {
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});
    const failure = new Error('boom');
    resolveImportedDesignFile.mockRejectedValue(failure);

    await expect(
      buildImportDesignRequestBody({
        uploadType: 'File Upload',
        name: 'Imported design',
      }),
    ).resolves.toEqual({
      errorMessage: 'Unable to read the selected design file. Please try again.',
    });

    expect(consoleError).toHaveBeenCalledWith('Error resolving design import file:', failure);
    consoleError.mockRestore();
  });
});
