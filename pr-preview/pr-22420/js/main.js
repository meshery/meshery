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

/* Table of Contents Scrollspy & Auto-scroll */
function initTableOfContentsScrollSpy() {
    const tocWrapper = document.querySelector(".content-table-wrapper");
    const toc = document.getElementById("TableOfContents");
    if (!toc || !tocWrapper) return;

    const tocLinks = Array.from(toc.querySelectorAll('a[href^="#"]'));
    if (!tocLinks.length) return;

    const linksByTargetId = new Map();
    const targets = [];

    tocLinks.forEach((link) => {
        const href = link.getAttribute("href");
        if (!href || href === "#") return;
        let id = href.slice(1);
        try {
            id = decodeURIComponent(id);
        } catch (e) {
            // Keep original if decoding fails
        }
        const target = document.getElementById(id);
        if (target) {
            targets.push(target);
            linksByTargetId.set(id, link);
        }
    });

    if (!targets.length) return;

    let currentActiveLink = null;
    let isUserClicking = false;
    let clickIdleTimer = null;
    let clickFallbackTimer = null;

    function resetClickLock() {
        if (!isUserClicking) return;
        isUserClicking = false;
        clearTimeout(clickIdleTimer);
        clearTimeout(clickFallbackTimer);
        determineActiveHeading();
    }

    function scrollTocLinkIntoView(link) {
        if (!tocWrapper || tocWrapper.offsetParent === null) return;

        const wrapperRect = tocWrapper.getBoundingClientRect();
        const linkRect = link.getBoundingClientRect();
        const padding = 24;

        const isAbove = linkRect.top < wrapperRect.top + padding;
        const isBelow = linkRect.bottom > wrapperRect.bottom - padding;

        if (isAbove || isBelow) {
            const targetScrollTop = tocWrapper.scrollTop + (linkRect.top - wrapperRect.top) - (wrapperRect.height / 3);
            tocWrapper.scrollTo({
                top: Math.max(0, targetScrollTop),
                behavior: "smooth"
            });
        }
    }

    function setActiveLink(link, shouldScrollContainer) {
        if (!link || link === currentActiveLink) return;
        if (currentActiveLink) {
            currentActiveLink.classList.remove("active");
        }
        link.classList.add("active");
        currentActiveLink = link;

        if (shouldScrollContainer) {
            scrollTocLinkIntoView(link);
        }
    }

    function determineActiveHeading() {
        if (isUserClicking) return;

        // 1. If at bottom of page, activate the last heading
        const scrollBottom = window.innerHeight + window.scrollY;
        const documentHeight = document.documentElement.scrollHeight;
        if (scrollBottom >= documentHeight - 2) {
            const lastTarget = targets[targets.length - 1];
            const link = linksByTargetId.get(lastTarget.id);
            if (link) setActiveLink(link, true);
            return;
        }

        // 2. If at top of document (before first heading)
        if (window.scrollY < 80 && targets.length > 0) {
            const firstTarget = targets[0];
            const link = linksByTargetId.get(firstTarget.id);
            if (link) setActiveLink(link, true);
            return;
        }

        // 3. Reading focus line (standard in Docusaurus, VitePress, Nextra)
        // Highlights the last heading that has passed above the reading threshold
        const readingOffset = 130;
        let activeTarget = targets[0];

        for (let i = 0; i < targets.length; i++) {
            const top = targets[i].getBoundingClientRect().top;
            if (top <= readingOffset) {
                activeTarget = targets[i];
            } else {
                break;
            }
        }

        if (activeTarget) {
            const link = linksByTargetId.get(activeTarget.id);
            if (link) setActiveLink(link, true);
        }
    }

    // Handle smooth clicks and avoid intermediate flickering during jump
    tocLinks.forEach((link) => {
        link.addEventListener("click", () => {
            const href = link.getAttribute("href");
            if (!href || href === "#") return;
            isUserClicking = true;
            setActiveLink(link, true);
            clearTimeout(clickIdleTimer);
            clearTimeout(clickFallbackTimer);
            // Fallback safety timeout in case no scroll occurs (e.g. already at target)
            clickFallbackTimer = setTimeout(resetClickLock, 1500);
        });
    });

    // Reset click lock immediately if user manually interacts
    window.addEventListener("wheel", resetClickLock, { passive: true });
    window.addEventListener("touchmove", resetClickLock, { passive: true });
    window.addEventListener("keydown", (e) => {
        if (["ArrowUp", "ArrowDown", "PageUp", "PageDown", "Home", "End", " "].includes(e.key)) {
            resetClickLock();
        }
    }, { passive: true });

    // Observe heading intersections
    if ("IntersectionObserver" in window) {
        const observer = new IntersectionObserver(() => {
            determineActiveHeading();
        }, {
            rootMargin: "-80px 0px -60% 0px",
            threshold: [0, 1.0]
        });

        targets.forEach((target) => observer.observe(target));
    }

    // Throttled scroll listener via requestAnimationFrame
    let ticking = false;
    window.addEventListener("scroll", () => {
        if (isUserClicking) {
            // Keep extending the lock until smooth scrolling settles (idle for 150ms)
            clearTimeout(clickIdleTimer);
            clickIdleTimer = setTimeout(resetClickLock, 150);
        }

        if (!ticking) {
            window.requestAnimationFrame(() => {
                determineActiveHeading();
                ticking = false;
            });
            ticking = true;
        }
    }, { passive: true });

    // Initial check on load (check hash or current scroll position)
    if (window.location.hash) {
        let hashId = window.location.hash.slice(1);
        try {
            hashId = decodeURIComponent(hashId);
        } catch (e) {
            // Keep original
        }
        const initialLink = linksByTargetId.get(hashId);
        if (initialLink) {
            setActiveLink(initialLink, true);
        } else {
            determineActiveHeading();
        }
    } else {
        determineActiveHeading();
    }
}

if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initTableOfContentsScrollSpy);
} else {
    initTableOfContentsScrollSpy();
}


