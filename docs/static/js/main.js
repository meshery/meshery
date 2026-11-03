(function ($) {
    'use strict';
    $(function () {
        $('[data-toggle="tooltip"]').tooltip();
        $('[data-toggle="popover"]').popover();
        $('.popover-dismiss').popover({
            trigger: 'focus'
        })
    });

    function bottomPos(element) {
        return element.offset().top + element.outerHeight();
    }
    $(function () {
        var promo = $(".js-td-cover");
        if (!promo.length) {
            return
        }
        var promoOffset = bottomPos(promo);
        var navbarOffset = $('.js-navbar-scroll').offset().top;
        var threshold = Math.ceil($('.js-navbar-scroll').outerHeight());
        if ((promoOffset - navbarOffset) < threshold) {
            $('.js-navbar-scroll').addClass('navbar-bg-onscroll');
        }
        $(window).on('scroll', function () {
            var navtop = $('.js-navbar-scroll').offset().top - $(window).scrollTop();
            var promoOffset = bottomPos($('.js-td-cover'));
            var navbarOffset = $('.js-navbar-scroll').offset().top;
            if ((promoOffset - navbarOffset) < threshold) {
                $('.js-navbar-scroll').addClass('navbar-bg-onscroll');
            } else {
                $('.js-navbar-scroll').removeClass('navbar-bg-onscroll');
                $('.js-navbar-scroll').addClass('navbar-bg-onscroll--fade');
            }
        });
    });
}(jQuery));


/*popup-hidden-div*/

function HideToggleFunction() {
    var hide = document.getElementById("hiddendiv");
    
    if (hide) {
        if (hide.style.display === "block") {
            hide.style.display = "none";
        } else {
            hide.style.display = "block";
        }
    }
}

/*clipboard*/

var getcodeelement = $('.clipboardjs');

getcodeelement.each(function () {
    
    /*trigger*/
    
    var clipButton = '<div class="btn-copy-wrap"><button type="button" class="clipbtn" aria-label="Copy code to clipboard"><i class="far fa-copy" aria-hidden="true"></i></button></div>';
    $(this).after(clipButton);
});

var clipboard = new Clipboard('.clipbtn', {
    text: function (trigger) {
        var container = trigger.closest('pre') || trigger.closest('.highlight');
        var content = container ? container.querySelector('.clipboardjs') : null;
        var text = content ? content.textContent : '';
        return text.trim().replace(/^\s*\$\s+/gm, '');
    }
});

/* Change copy icon to check icon when successfully copied*/
clipboard.on("success", (e) => {
    const button = e.trigger;
    if (button.dataset.isCopying === "true") {
        return;
    }

    const icon = button.querySelector('i');
    if (!icon) return;

    button.dataset.isCopying = "true";
    button.setAttribute('aria-label', 'Copied to clipboard');

    const originalIcon = icon.cloneNode(true);
    const text = document.createElement('span');
    const checkIcon = document.createElement('i');
    checkIcon.className = 'fas fa-check';
    checkIcon.style.color = "var(--brand-color-secondary)";
    checkIcon.setAttribute('aria-hidden', 'true');
    text.appendChild(checkIcon);
    text.appendChild(document.createTextNode(' Copied!'));
    text.style.color = "var(--brand-color-secondary)";

    button.replaceChild(text, icon);

    setTimeout(() => {
        // 1. Always reset accessible name/state first (safe even if detached)
        button.setAttribute('aria-label', 'Copy code to clipboard');
        button.removeAttribute("data-is-copying");
        // 2. Guard only the physical DOM manipulation
        if (button.isConnected === false) {
            return;
        }
        if (text.parentNode === button) {
            button.replaceChild(originalIcon, text);
        }
    }, 2000);
});
const toggleBtnSidebarNav = document.querySelector(".nav-toggle-btn--document");

if (toggleBtnSidebarNav) {
    toggleBtnSidebarNav.addEventListener("click", () => {
        const leftContainer = document.querySelector(".left-container");

        if (leftContainer) {
            const isActive = leftContainer.classList.toggle('left-container--active');

            const newState = isActive ? 'active' : 'inactive';
            localStorage.setItem('leftContainer-state', newState);
        }
    })
}

const toggleBtnMainNav = document.querySelector(".nav-toggle-btn--main");

if (toggleBtnMainNav) {
    toggleBtnMainNav.addEventListener("click", () => {
        let sidebarNav = document.getElementById("main_navbar")
        if (sidebarNav) {
            sidebarNav.classList.toggle("main-navbar--active")
        }      
    })
}

document.addEventListener("click", (event) => {
    let sidebarNav = document.getElementById("main_navbar")
    if (sidebarNav) {
        let isClickInsideSidebar = sidebarNav.contains(event.target)
        let isClickOnToggleButton = toggleBtnMainNav && toggleBtnMainNav.contains(event.target)

        if (!isClickInsideSidebar && !isClickOnToggleButton) {
            sidebarNav.classList.remove("main-navbar--active")
        }
    }
});

/* Table of Contents Scrollspy & Auto-scroll (Modern Docs) */
function initTableOfContentsScrollSpy() {
    const tocWrapper = document.querySelector(".content-table-wrapper");
    const tocNav = tocWrapper ? (tocWrapper.querySelector(".td-toc nav") || tocWrapper.querySelector("nav")) : null;
    if (!tocWrapper || !tocNav) return;

    const tocLinks = Array.from(tocNav.querySelectorAll('a[href^="#"]'));
    if (!tocLinks.length) return;

    const headings = tocLinks
        .map(link => {
            const href = link.getAttribute("href");
            if (!href || href === "#") return null;
            let id = href.slice(1);
            try {
                id = decodeURIComponent(id);
            } catch (e) {
                // Keep original if decoding fails
            }
            const target = document.getElementById(id);
            return target ? { link, target } : null;
        })
        .filter(Boolean);

    if (!headings.length) return;

    function scrollTocContainer(link) {
        if (!tocWrapper || tocWrapper.offsetParent === null) return;
        const wrapperRect = tocWrapper.getBoundingClientRect();
        const linkRect = link.getBoundingClientRect();
        const padding = 20;

        if (linkRect.top < wrapperRect.top + padding || linkRect.bottom > wrapperRect.bottom - padding) {
            const targetScrollTop = tocWrapper.scrollTop + (linkRect.top - wrapperRect.top) - (wrapperRect.height / 3);
            tocWrapper.scrollTo({
                top: Math.max(0, targetScrollTop),
                behavior: "smooth"
            });
        }
    }

    let activeLink = null;
    function setActive(link) {
        if (!link || link === activeLink) return;
        if (activeLink) activeLink.classList.remove("active");
        link.classList.add("active");
        activeLink = link;
        scrollTocContainer(link);
    }

    function onScroll() {
        const scrollBottom = window.innerHeight + window.scrollY;
        if (scrollBottom >= document.documentElement.scrollHeight - 2) {
            setActive(headings[headings.length - 1].link);
            return;
        }

        const readingOffset = 130;
        let current = headings[0];
        for (let i = 0; i < headings.length; i++) {
            if (headings[i].target.getBoundingClientRect().top <= readingOffset) {
                current = headings[i];
            } else {
                break;
            }
        }
        if (current) setActive(current.link);
    }

    let ticking = false;
    window.addEventListener("scroll", () => {
        if (!ticking) {
            window.requestAnimationFrame(() => {
                onScroll();
                ticking = false;
            });
            ticking = true;
        }
    }, { passive: true });

    headings.forEach(({ link }) => {
        link.addEventListener("click", () => setActive(link));
    });

    onScroll();
}

if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initTableOfContentsScrollSpy);
} else {
    initTableOfContentsScrollSpy();
}


