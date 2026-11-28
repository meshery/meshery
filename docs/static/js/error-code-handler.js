function displayErrorCode() {
  const rawHash = window.location.hash.substring(1);
  if (!rawHash) {
    return;
  }

  let hash = rawHash;
  try {
    hash = decodeURIComponent(rawHash);
  } catch (e) {
    // Keep rawHash if it cannot be URI-decoded
  }

  const errorRow = document.getElementById(hash) || document.getElementById(rawHash);
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
