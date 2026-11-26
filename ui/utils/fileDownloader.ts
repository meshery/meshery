import { FILTER, PATTERN } from '@/constants/navigator';

const DOWNLOAD_PATH = {
  [PATTERN]: ({ id, source_type, params }) =>
    source_type
      ? `/api/pattern/download/${id}/${source_type}`
      : params
        ? `/api/pattern/download/${id}?${params}`
        : `/api/pattern/download/${id}`,
  [FILTER]: ({ id }) => `/api/filter/download/${id}`,
};

export const downloadFileFromUrl = (url, name) => {
  let linkElement = document.createElement('a');
  linkElement.setAttribute('href', url);
  linkElement.setAttribute('download', name);
  linkElement.click();
  linkElement.remove();
};

/*
 * @param {string} type - application, pattern, filter
 * @param {string} id - id of the file to download
 * @param {string} name - name of the file to download
 * @param {string} source_type - source type of the file to download
 * @returns {Promise<void>}
 */
export default async function downloadContent({ id, type, name, source_type, params }) {
  const uri = DOWNLOAD_PATH[type]({ id, source_type, params });
  if (!uri) {
    throw new Error('Invalid type of content to download', type);
  }
  const response = await fetch(uri, { redirect: 'error' });
  if (!response.ok) {
    throw new Error(`Download failed (HTTP ${response.status})`);
  }
  const disposition = response.headers.get('Content-Disposition');
  const filename = disposition?.match(/(?:^|;)\s*filename\s*=\s*(?:"((?:[^"\\]|\\.)*)"|([^;]*))/i);
  const downloadName = filename?.[1]?.replace(/\\(.)/g, '$1') || filename?.[2]?.trim() || name;
  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  try {
    downloadFileFromUrl(url, downloadName);
  } finally {
    window.URL.revokeObjectURL(url);
  }
}

export const downloadFileFromContent = (content, fileName, type) => {
  const blob = new Blob([content], { type: type });
  const url = window.URL.createObjectURL(blob);
  downloadFileFromUrl(url, fileName);
};
