function displayErrorCode() {
  const hash = window.location.hash.substring(1);
  if (hash) {
    const errorRow = document.getElementById(hash);

    if (errorRow) {
      const parentRow = errorRow.closest("tr");

      if (parentRow) {
        const onclickAttr = parentRow.getAttribute("onclick");
        if (onclickAttr) {
          const matches = onclickAttr.match(/toggle_visibility\('([^']+)'\)/);

          if (matches && matches[1]) {
            const hiddenRowID = matches[1];
            const hiddenRow = document.getElementById(hiddenRowID);
            if (hiddenRow) {
              hiddenRow.style.display = "table-row";
              hiddenRow.style.visibility = "visible";
            }
          }
        }
      }
    }
  }
}
document.addEventListener("DOMContentLoaded", displayErrorCode);
window.addEventListener("hashchange", displayErrorCode);
