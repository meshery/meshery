document.addEventListener("DOMContentLoaded", function () {
  const toggleBtn = document.getElementById("mode-toggle-btn");
  const logo = document.getElementById("logo-dark-light");
  const body = document.querySelector("body");

  // Load saved theme if available
  let savedTheme = localStorage.getItem("mode");
  let isDark;

  if (savedTheme === "dark") {
    body.classList.add("dark-mode");
    logo.src = logo.dataset.logoForDark;
    isDark = true;
  } else if (savedTheme === "light") {
    body.classList.remove("dark-mode");
    logo.src = logo.dataset.logoForLight;
    isDark = false;
  } else {
    // No saved theme — fallback to OS preference
    const osPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    if (osPrefersDark) {
      body.classList.add("dark-mode");
      logo.src = logo.dataset.logoForDark;
      isDark = true;
    } else {
      body.classList.remove("dark-mode");
      logo.src = logo.dataset.logoForLight;
      isDark = false;
    }
  }

  // Toggle button
  toggleBtn.addEventListener("click", function () {
    isDark = !isDark;
    body.classList.toggle("dark-mode");
    logo.src = isDark ? logo.dataset.logoForDark : logo.dataset.logoForLight;
    localStorage.setItem("mode", isDark ? "dark" : "light");
  });
});

