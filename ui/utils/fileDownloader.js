export function fileDownloader(id, name, source_type) {
  let dataStr = JSON.stringify(id);
  let dataUri = `/api/application/download/${id}/${source_type}`+ encodeURIComponent(dataStr);

  let downloadFileDefaultName = name;

  let linkElement = document.createElement('a');
  linkElement.setAttribute('href', dataUri);
  linkElement.setAttribute('download', downloadFileDefaultName);
  linkElement.click();
  linkElement.remove()
}
