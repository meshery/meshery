function displayErrorCode() {
  const hash = window.location.hash.substring(1);
  if (!hash) {
    return;
  }

  const errorRow = document.getElementById(hash);
  if (!errorRow) {
    return;
  }

  const parentRow = errorRow.closest("tr");
  if (!parentRow) {
    return;
  }

  const onclickAttr = parentRow.getAttribute("onclick");
  if (!onclickAttr) {
    return;
  }

  const matches = onclickAttr.match(/toggle_visibility\('([^']+)'\)/);
  if (!matches || !matches[1]) {
    return;
  }

  const hiddenRow = document.getElementById(matches[1]);
  if (!hiddenRow) {
    return;
  }

  hiddenRow.style.display = "table-row";
  hiddenRow.style.visibility = "visible";
}

document.addEventListener("DOMContentLoaded", displayErrorCode);
window.addEventListener("hashchange", displayErrorCode);