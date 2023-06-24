export default function downloadFile({ id, type,name, source_type }) {
  let dataStr = JSON.stringify(id);
  let dataUri = "";

  if (source_type) {
    dataUri = `/api/application/download/${id}/${source_type}`+ encodeURIComponent(dataStr);
  } else {
    dataUri = `/api/application/download/${id}`;
  }

  if (type === "pattern") {
    dataUri = `/api/pattern/download/${id}`;
  }

  let linkElement = document.createElement('a');
  linkElement.setAttribute('href', dataUri);
  linkElement.setAttribute('download', name);
  linkElement.click();
  linkElement.remove();
}
