import { getUnit8ArrayDecodedFile } from '@/utils/utils';
import { findSelectedFileInDialog, readFileAsBytes } from '@/utils/fileUpload';

export const filenameFromDataUrl = (dataUrl?: string): string | undefined => {
  if (!dataUrl || typeof dataUrl !== 'string') return undefined;
  const commaIndex = dataUrl.indexOf(',');
  if (commaIndex === -1) return undefined;
  const match = dataUrl.slice(0, commaIndex).match(/;name=([^;,]+)(?:;|$)/);
  if (!match) return undefined;
  try {
    return decodeURIComponent(match[1]);
  } catch {
    return match[1];
  }
};

export const decodeDesignDataUrlToBytes = (dataUrl?: string): number[] | null => {
  if (!dataUrl || typeof dataUrl !== 'string') return null;
  const commaIndex = dataUrl.indexOf(',');
  if (commaIndex === -1) return null;
  const metadata = dataUrl.slice(0, commaIndex);
  if (!metadata.startsWith('data:') || !metadata.endsWith(';base64')) {
    return null;
  }
  const payload = dataUrl.slice(commaIndex + 1);
  const normalizedPayload = payload.replace(/\s+/g, '');
  if (
    !normalizedPayload ||
    payload.includes(';base64,') ||
    !/^[A-Za-z0-9+/]+={0,2}$/.test(normalizedPayload)
  ) {
    return null;
  }
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
  const fileName = fileData
    ? (filenameFromDataUrl(dataUrl) ?? 'design.yaml')
    : (selectedFile?.name ?? 'design.yaml');

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
