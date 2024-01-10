// import dataFetch from "lib/data-fetch";

export default function downloadFile({ id, type, name, oci }) {
  if (type === 'pattern') {
    oci ? `/api/pattern/download/${id}?oci=true` : dataUri = `/api/pattern/download/${id}`;
  }

  if (type === 'filter') {
    dataUri = `/api/filter/download/${id}`;
  }

  let linkElement = document.createElement('a');
  linkElement.setAttribute('href', dataUri);
  linkElement.setAttribute('download', name);
  linkElement.click();
  linkElement.remove();
}
