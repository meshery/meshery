"use strict";

// DARK/LIGHT MODE TOGGLE

let toggleBtn = document.getElementById("mode-toggle-btn");
toggleBtn.onclick = setMode;

// Set dark mode by default
document.body.classList.add("dark-mode");

if (localStorage.getItem("mode") !== null && localStorage.getItem("mode") !== "dark-mode") {
  document.body.classList.remove("dark-mode");
}
 
// Set logo for dark mode
if (document.body.classList.contains("dark-mode")) {
  let allLogos = document.querySelectorAll("#logo-dark-light");
  allLogos.forEach(e => e.src = e.dataset.logoForDark);
  localStorage.setItem("mode", "dark-mode");
}

// Set logo for light mode
function setMode() {
  document.body.classList.toggle("dark-mode");
  let allLogos = document.querySelectorAll("#logo-dark-light");

  if (document.body.classList.contains("dark-mode")) {
    allLogos.forEach(e => e.src = e.dataset.logoForDark);
    localStorage.setItem("mode", "dark-mode");
  } 
  
  else {
    allLogos.forEach(e => e.src = e.dataset.logoForLight);
    localStorage.setItem("mode", "");
  }
}
