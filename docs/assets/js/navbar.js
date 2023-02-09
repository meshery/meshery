"use strict";


// DARK/LIGHT MODE TOGGLE

let toggleBtn = document.getElementById("mode-toggle-btn")
toggleBtn.onclick = setMode;
if (localStorage.getItem("mode")) setMode();
function setMode() {
    document.body.classList.toggle("dark-mode")
    // let layer5Logos = document.querySelectorAll("#layer5-logo");
    // let allLogos=document.querySelectorAll("#logo-dark-light");
    // console.log(allLogos)
    if (document.body.classList.contains("dark-mode")) {
        // layer5Logos.forEach(e => e.src = '../images/company-logo/layer5-dark-mode-logo.svg')
        // allLogos.forEach(e=>e.src=e.dataset.logoForDark)
    } else {
        // layer5Logos.forEach(e => e.src = '../images/company-logo/layer5-no-trim.svg')
        // allLogos.forEach(e=>e.src=e.dataset.logoForLight)
    }
    if (document.body.classList.contains("dark-mode")) {
        localStorage.setItem("mode", "dark-mode")
    } else {
        localStorage.setItem("mode", "")
    }
}