import { getUnit8ArrayDecodedFile } from '@/utils/utils';
import { findSelectedFileInDialog, readFileAsBytes } from '@/utils/fileUpload';

export const filenameFromDataUrl = (dataUrl?: string): string | undefined => {
  if (!dataUrl || typeof dataUrl !== 'string') return undefined;
  const match = dataUrl.match(/;name=([^;,]+)(?:;|,|$)/);
  if (!match) return undefined;
  try {
    return decodeURIComponent(match[1]);
  } catch {
    return match[1];
  }
};

export const decodeDesignDataUrlToBytes = (dataUrl?: string): number[] | null => {
  if (!dataUrl || typeof dataUrl !== 'string') return null;
  try {
    return getUnit8ArrayDecodedFile(dataUrl);
  } catch {
    return null;
  }
};

export const findSelectedDesignImportFile = (): File | undefined =>
  findSelectedFileInDialog('input[type="file"]');

export const resolveImportedDesignFile = async (
  dataUrl?: string,
): Promise<{ fileData: number[]; fileName: string } | null> => {
  let fileData = decodeDesignDataUrlToBytes(dataUrl);
  const selectedFile = findSelectedDesignImportFile();
  const fileName = filenameFromDataUrl(dataUrl) ?? selectedFile?.name ?? 'design.yaml';

  if (!fileData && selectedFile) {
    try {
      fileData = await readFileAsBytes(selectedFile);
    } catch {
      fileData = null;
    }
  }

  if (!fileData) {
    return null;
  }

  return { fileData, fileName };
};
