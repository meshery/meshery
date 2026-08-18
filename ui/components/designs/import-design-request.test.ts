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
      requestBody: {
        name: 'Imported design',
        fileName: 'imported-design.yaml',
        file: [1, 2, 3],
      },
    });
  });

  it('builds the URL import request body from the contract fields', async () => {
    await expect(
      buildImportDesignRequestBody({
        uploadType: 'URL Import',
        name: 'Imported design',
        url: 'https://example.com/design.yaml',
      }),
    ).resolves.toEqual({
      requestBody: {
        url: 'https://example.com/design.yaml',
        name: 'Imported design',
      },
    });

    expect(resolveImportedDesignFile).not.toHaveBeenCalled();
  });

  it('returns a user-facing error when the URL import has a missing or blank URL', async () => {
    await expect(
      buildImportDesignRequestBody({
        uploadType: 'URL Import',
        name: 'Imported design',
        url: '   ',
      }),
    ).resolves.toEqual({
      errorMessage: 'Please enter a design URL before continuing.',
    });

    await expect(
      buildImportDesignRequestBody({
        uploadType: 'URL Import',
        name: 'Imported design',
      }),
    ).resolves.toEqual({
      errorMessage: 'Please enter a design URL before continuing.',
    });
  });

  it('returns a user-facing error for an unrecognized upload type', async () => {
    await expect(
      buildImportDesignRequestBody({
        uploadType: 'Something Else',
        name: 'Imported design',
      }),
    ).resolves.toEqual({
      errorMessage: 'Please choose a valid design import source before continuing.',
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
