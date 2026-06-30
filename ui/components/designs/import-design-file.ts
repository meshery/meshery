import { getUnit8ArrayDecodedFile } from '@/utils/utils';
import { findSelectedFileInDialog, readFileAsBytes } from '@/utils/fileUpload';

export const filenameFromDataUrl = (dataUrl?: string): string | undefined => {
  if (!dataUrl) return undefined;
  const match = dataUrl.match(/;name=([^;]+);/);
  return match ? decodeURIComponent(match[1]) : undefined;
};

export const decodeDesignDataUrlToBytes = (dataUrl?: string): number[] | null => {
  if (!dataUrl) return null;
  return getUnit8ArrayDecodedFile(dataUrl);
};

export const findSelectedDesignImportFile = (): File | undefined =>
  findSelectedFileInDialog('input[type="file"]');

export const resolveImportedDesignFile = async (dataUrl?: string) => {
  let fileData = decodeDesignDataUrlToBytes(dataUrl);
  const selectedFile = findSelectedDesignImportFile();
  const fileName = filenameFromDataUrl(dataUrl) ?? selectedFile?.name ?? 'design.yaml';

  if (!fileData && selectedFile) {
    fileData = await readFileAsBytes(selectedFile);
  }

  if (!fileData) {
    return null;
  }

  return { fileData, fileName };
};
