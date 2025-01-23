function addAnchorListeners() {
  const anchors = document.querySelectorAll('a[href^="#"]');
  anchors.forEach((anchor) => {
    if (!anchor.dataset.listenerAdded) {
      anchor.dataset.listenerAdded = "true"; // Prevent adding the same listener multiple times
      anchor.addEventListener("click", function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute("href"));
        console.log(target);
        if (target) {
          console.log(target.offsetTop);
          const offsetPosition =
            target.getBoundingClientRect().top + window.scrollY - 400; // Adjust '100' for the offset
          window.scrollTo({
            top: offsetPosition, // Adjust for fixed headers
            behavior: "smooth",
          });
        }
      });
    }
  });
}

// Initial run
document.addEventListener("DOMContentLoaded", () => {
  addAnchorListeners();

  // Observe dynamic changes
  const observer = new MutationObserver(() => {
    addAnchorListeners();
  });
  observer.observe(document.body, { childList: true, subtree: true });
});
