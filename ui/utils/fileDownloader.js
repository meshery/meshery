export default function downloadFile(props) {
  const { id, type, name, oci } = props;
  let dataUri = '';

  if (type === 'pattern') {
    dataUri = oci ? `/api/pattern/download/${id}?oci=true` : `/api/pattern/download/${id}`;
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
