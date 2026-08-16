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

  const detailsId = parentRow.dataset.detailsId;
  if (!detailsId) {
    return;
  }

  const hiddenRow = document.getElementById(detailsId);
  if (!hiddenRow) {
    return;
  }

  hiddenRow.style.display = "table-row";
  hiddenRow.style.visibility = "visible";
}

document.addEventListener("DOMContentLoaded", displayErrorCode);
window.addEventListener("hashchange", displayErrorCode);
