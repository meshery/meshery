export function fileDownloader(id, name) {
  let dataStr = JSON.stringify(id);
  let dataUri = `/api/application/download/${id}?source-type=k8s_manifest`+ encodeURIComponent(dataStr);

  let downloadFileDefaultName = name;

  let linkElement = document.createElement('a');
  linkElement.setAttribute('href', dataUri);
  linkElement.setAttribute('download', downloadFileDefaultName);
  linkElement.click();
  linkElement.remove()
}
