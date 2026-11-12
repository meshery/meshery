"use strict";

document.addEventListener("DOMContentLoaded", () => {
    const toggleBtn = document.getElementById("mode-toggle-btn");

    if (toggleBtn) {
        toggleBtn.addEventListener("click", setMode);
    }

    // Keep the early-applied theme and update the logos.
    updateLogos();

    // Follow system preference changes only when no explicit preference exists.
    if (window.matchMedia) {
        const systemTheme = window.matchMedia("(prefers-color-scheme: dark)");

        const handleSystemThemeChange = event => {
            if (localStorage.getItem("mode") === null) {
                applyTheme(event.matches);
            }
        };

        if (systemTheme.addEventListener) {
            systemTheme.addEventListener("change", handleSystemThemeChange);
        } else {
            systemTheme.addListener(handleSystemThemeChange);
        }
    }

    // Synchronize theme changes made in other tabs.
    window.addEventListener("storage", event => {
        if (event.storageArea !== localStorage) return;

        if (event.key === "mode" || event.key === null) {
            applyCurrentPreference();
        }
    });
    // Reconcile the theme in case preferences changed before listeners were registered.
    applyCurrentPreference();
});

function applyTheme(isDark) {
    if (window.applyDocsTheme) {
        window.applyDocsTheme(isDark);
    } else {
        document.body.classList.toggle("dark-mode", isDark);
        updateLogos();
    }
}

function applyCurrentPreference() {
    const savedMode = localStorage.getItem("mode");

    if (savedMode === "light-mode") {
        applyTheme(false);
    } else if (savedMode === "dark-mode") {
        applyTheme(true);
    } else {
        const systemPreference = window.matchMedia
            ? window.matchMedia("(prefers-color-scheme: dark)").matches
            : true;

        applyTheme(systemPreference);
    }
}

function setMode() {
    const isDark = !document.body.classList.contains("dark-mode");

    localStorage.setItem(
        "mode",
        isDark ? "dark-mode" : "light-mode"
    );

    applyTheme(isDark);
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