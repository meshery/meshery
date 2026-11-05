"use strict";

// DARK/LIGHT MODE TOGGLE

const toggleBtn = document.getElementById("mode-toggle-btn");
toggleBtn.addEventListener("click", setMode);


function setMode() {
  document.body.classList.toggle("dark-mode");
  updateLogos();
  if (document.body.classList.contains("dark-mode")) {
    localStorage.setItem("mode", "dark-mode");
  } else {
    localStorage.setItem("mode", "light-mode");
  }
}

function updateLogos() {
  const allLogos = document.querySelectorAll("#logo-dark-light");
  allLogos.forEach(logo => {
    if (document.body.classList.contains("dark-mode")) {
      logo.src = logo.dataset.logoForDark;
    } else {
      logo.src = logo.dataset.logoForLight;
    }
  });
}
