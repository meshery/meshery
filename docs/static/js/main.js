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

var clipButtonHtml = '<div class="btn-copy-wrap"><button type="button" class="clipbtn" aria-label="Copy code to clipboard"><i class="far fa-copy" aria-hidden="true"></i></button></div>';

// 1. Shortcode blocks (.codeblock-pre)
$('pre.codeblock-pre').each(function () {
    var pre = $(this);
    var code = pre.find('.clipboardjs');
    if (code.length) {
        var text = code.text().trim();
        if (text.indexOf('$') !== -1) {
            text = text.replace(/^[ \t]+(?=\$ )/gm, '');
        }
        code.text(text);
    }
    if (!pre.find('.clipbtn').length) {
        pre.append(clipButtonHtml);
    }
});

// 2. Standard Markdown syntax-highlighted blocks (.highlight)
$('.highlight').each(function () {
    if (!$(this).find('.clipbtn').length) {
        $(this).append(clipButtonHtml);
    }
});

// 3. Standalone pre blocks (untagged markdown blocks not inside .highlight or .codeblock-pre)
$('pre:not(.codeblock-pre)').each(function () {
    if ($(this).closest('.highlight').length || $(this).closest('.codeblock-pre').length) return;
    if (!$(this).find('.clipbtn').length) {
        $(this).append(clipButtonHtml);
    }
});

var clipboard = new Clipboard('.clipbtn', {
    text: function (trigger) {
        var container = trigger.closest('pre') || trigger.closest('.highlight');
        if (!container) return '';
        var content = container.querySelector('.clipboardjs') || container.querySelector('code') || container;
        var clone = content ? content.cloneNode(true) : null;
        if (clone) {
            var btn = clone.querySelector('.btn-copy-wrap');
            if (btn) btn.remove();
            return (clone.textContent || '').trim().replace(/^\s*\$\s+/gm, '');
        }
        return '';
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

    const originalLabel = button.getAttribute('aria-label') || 'Copy code to clipboard';
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
        button.setAttribute('aria-label', originalLabel);
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

