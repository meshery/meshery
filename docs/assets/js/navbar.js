"use strict";

let toggleBtn = document.getElementById("mode-toggle-btn")
toggleBtn.onclick = setMode;
document.body.classList.add("dark-mode");
if (localStorage.getItem("mode")) setMode();
function setMode() {
  document.body.classList.toggle("dark-mode")

  let allLogos = document.querySelectorAll("#logo-dark-light");

  if (document.body.classList.contains("dark-mode")) {
    allLogos.forEach(e => e.src = e.dataset.logoForDark)
  } else {

    allLogos.forEach(e => e.src = e.dataset.logoForLight)
  }
  if (document.body.classList.contains("dark-mode")) {
    localStorage.setItem("mode", "dark-mode")
  } else {
    localStorage.setItem("mode", "")
  }
}
