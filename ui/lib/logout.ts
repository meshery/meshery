export function submitLogout() {
  const form = document.createElement('form');
  form.method = 'POST';
  form.action = '/user/logout';
  document.body.appendChild(form);
  form.submit();
}
