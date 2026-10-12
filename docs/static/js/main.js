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
})

/*
 * Version banner – injects a dismissible warning when the user is
 * viewing a pre-built static snapshot (e.g. /v0.8/, /v0.9/).
 */
;(function () {
    var LATEST_BASE = 'https://docs.meshery.io';
    var STORAGE_KEY = 'version-banner-dismissed';
    var VERSION_RE = /^\/(v\d+\.\d+)(?:\/|$)/;

    // sessionStorage can throw when storage is disabled (e.g. privacy modes).
    // Guard access so the banner still renders and dismissal never errors.
    function isDismissed() {
        try { return sessionStorage.getItem(STORAGE_KEY) === '1'; } catch (e) { return false; }
    }
    function rememberDismissed() {
        try { sessionStorage.setItem(STORAGE_KEY, '1'); } catch (e) {}
    }

    if (isDismissed()) return;

    var path = window.location.pathname;
    var match = VERSION_RE.exec(path);

    if (!match) return;

    var matchedPrefix = match[1];

    var pagePath = path.replace('/' + matchedPrefix, '') || '/';
    var latestUrl = LATEST_BASE + pagePath;

    var banner = document.createElement('div');
    banner.id = 'version-banner';
    banner.setAttribute('role', 'alert');
    banner.style.cssText = [
        'position:fixed',
        'top:0',
        'left:0',
        'right:0',
        'z-index:1100',
        'display:flex',
        'align-items:center',
        'justify-content:center',
        'gap:0.5rem',
        'padding:0.55rem 2.5rem 0.55rem 1rem',
        'background-color:#FBF3B0',
        'color:#1a1a1a',
        'font-size:0.875rem',
        'font-weight:600',
        'text-align:center',
        'line-height:1.4',
        'box-shadow:0 2px 6px rgba(0,0,0,0.15)',
        'font-family:inherit'
    ].join(';');

    var icon = document.createElement('span');
    icon.setAttribute('aria-hidden', 'true');
    icon.textContent = '\u26A0';
    icon.style.cssText = 'flex-shrink:0;font-size:1rem;';

    var text = document.createElement('span');
    text.innerHTML =
        'You are not viewing the latest version of the documentation. ' +
        '<a href="' + latestUrl + '" style="color:#0A7D6B;font-weight:700;text-decoration:underline;text-underline-offset:2px;margin-left:0.25rem;">View the latest version &rarr;</a>';

    var closeBtn = document.createElement('button');
    closeBtn.innerHTML = '&times;';
    closeBtn.setAttribute('aria-label', 'Dismiss version warning');
    closeBtn.style.cssText = [
        'position:absolute',
        'right:0.75rem',
        'top:50%',
        'transform:translateY(-50%)',
        'background:none',
        'border:none',
        'color:#1a1a1a',
        'font-size:1.25rem',
        'cursor:pointer',
        'padding:0.25rem',
        'line-height:1',
        'opacity:0.7'
    ].join(';');

    closeBtn.addEventListener('click', function () {
        banner.remove();
        var css = document.getElementById('version-banner-adjust');
        if (css) css.remove();
        rememberDismissed();
    });

    banner.appendChild(icon);
    banner.appendChild(text);
    banner.appendChild(closeBtn);

    document.body.insertBefore(banner, document.body.firstChild);

    var FALLBACK_BANNER_HEIGHT = 36;
    var bh = banner.offsetHeight || FALLBACK_BANNER_HEIGHT;
    var px = bh + 'px';

    var navbar = document.querySelector('.navbar');
    var navH = navbar ? navbar.offsetHeight : 104;
    var totalOffset = (navH + bh) + 'px';

    var adjustCSS = document.createElement('style');
    adjustCSS.id = 'version-banner-adjust';
    adjustCSS.textContent = [
        'header { top: ' + px + ' !important; }',
        '.navbar { top: ' + px + ' !important; }',
        '.main-container { padding-top: ' + totalOffset + ' !important; }',
        '.breadcrumb { margin-top: ' + totalOffset + ' !important; }',
        '.left-container { top: ' + totalOffset + ' !important; margin-top: ' + totalOffset + ' !important; height: calc(100vh - ' + totalOffset + ') !important; }',
        '.sidebar-container { height: calc(100vh - ' + totalOffset + ' - 3vh) !important; }',
        '.content-table { top: calc(' + totalOffset + ' + 5rem) !important; max-height: calc(100vh - ' + totalOffset + ' - 6rem) !important; }',
        '#main_navbar { top: ' + totalOffset + ' !important; }',
        'html { scroll-padding-top: ' + totalOffset + ' !important; }'
    ].join('\n');
    document.head.appendChild(adjustCSS);
})()
