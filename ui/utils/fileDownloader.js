const DOWNLOAD_PATH = {
  [PATTERN]: ({ id, source_type, oci }) =>
    source_type
      ? `/api/pattern/download/${id}/${source_type}`
      : oci
      ? `/api/pattern/download/${id}?oci=true`
      : `/api/pattern/download/${id}`,
  [FILTER]: ({ id }) => `/api/filter/download/${id}`
};

export const downloadFileFromUrl = (url, name) => {
  let linkElement = document.createElement("a");
  linkElement.setAttribute("href", url);
  linkElement.setAttribute("download", name);
  linkElement.click();
  linkElement.remove();
};

/*
 * @param {string} type - application, pattern, filter
 * @param {string} id - id of the file to download
 * @param {string} name - name of the file to download
 * @param {string} source_type - source type of the file to download
 * @returns {void}
 */
export default function downloadContent({ id, type, name, source_type, oci=false }) {
  const uri = DOWNLOAD_PATH[type]({ id, source_type, oci });
  if (!uri) {
    throw new Error("Invalid type of content to download", type);
  }
  console.log("Download Uri", uri);
  downloadFileFromUrl(uri, name);
}

export const downloadFileFromContent = (content, fileName, type) => {
  const blob = new Blob([content], { type: type });
  const url = window.URL.createObjectURL(blob);
  downloadFileFromUrl(url, fileName);
};