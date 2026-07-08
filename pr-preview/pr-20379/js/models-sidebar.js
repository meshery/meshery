document.addEventListener('DOMContentLoaded', function() {
  const wrapper = document.querySelector('.models-sidebar-wrapper');
  const sections = document.querySelectorAll('.models-list section[id^="models-group-"]');
  const links = document.querySelectorAll('.alphabet-link');

  if (!wrapper || sections.length === 0 || links.length === 0) return;

  function updateActiveLetter() {
    let currentActive = null;
    const viewportMid = window.innerHeight / 3;

    for (let section of sections) {
      const rect = section.getBoundingClientRect();
      if (rect.top <= viewportMid) {
        currentActive = section.id;
      }
    }

    links.forEach(link => {
      link.classList.remove('active');
      if (currentActive) {
        const targetId = link.getAttribute('href').substring(1);
        if (targetId === currentActive) {
          link.classList.add('active');
        }
      }
    });
  }

  links.forEach(link => {
    link.addEventListener('click', function(e) {
      e.preventDefault();
      const targetId = this.getAttribute('href').substring(1);
      const targetSection = document.getElementById(targetId);

      if (targetSection) {
        targetSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
        links.forEach(l => l.classList.remove('active'));
        this.classList.add('active');
      }
    });
  });

  function updateSidebarPlacement() {
    const nav = document.querySelector('.models-top-nav');
    if (!nav) return;
    wrapper.style.top = `${nav.getBoundingClientRect().top}px`;
  }

  window.addEventListener('resize', updateSidebarPlacement);
  window.addEventListener('scroll', updateActiveLetter, { passive: true });

  updateSidebarPlacement();
  updateActiveLetter();
});