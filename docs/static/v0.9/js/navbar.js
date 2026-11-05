"use strict";

document.addEventListener("DOMContentLoaded", () => {

    const savedMode = localStorage.getItem("mode");

    if (!savedMode || savedMode === "dark-mode") {
        document.body.classList.add("dark-mode");
    } else {
        document.body.classList.remove("dark-mode");
    }

    updateLogos();

    const toggleBtn = document.getElementById("mode-toggle-btn");
    if (toggleBtn) {
        toggleBtn.addEventListener("click", setMode);
    }
});

function setMode() {
    document.body.classList.toggle("dark-mode");
    updateLogos();
    localStorage.setItem(
        "mode",
        document.body.classList.contains("dark-mode") ? "dark-mode" : "light-mode"
    );
}

function updateLogos() {
    const allLogos = document.querySelectorAll(".logo-dark-light");
    allLogos.forEach(logo => {
        if (document.body.classList.contains("dark-mode")) {
            logo.src = logo.dataset.logoForDark;
        } else {
            logo.src = logo.dataset.logoForLight;
        }
    });
}

document.addEventListener("keydown", function (event) {
  const searchInput = document.getElementById("sidebar-search-input");
  if (!searchInput) return;

  const activeEl = document.activeElement;

  const isEditable =
    activeEl &&
    (activeEl.isContentEditable ||
      ["INPUT", "TEXTAREA"].includes(activeEl.tagName));

  if (
    event.key === "/" &&
    !event.ctrlKey &&
    !event.metaKey &&
    !event.altKey &&
    !isEditable
  ) {
    event.preventDefault();
    searchInput.focus();
  }
});