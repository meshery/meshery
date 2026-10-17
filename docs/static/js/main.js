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

$(document).ready(function () {
    // 1. All codeblock-pre blocks (shortcodes, mesheryctl reference, raw HTML)
    $('.codeblock-pre').each(function () {
        if (!$(this).find('.clipbtn').length) {
            var cb = $(this).find('.codeblock');
            if (cb.length && !cb.find('.clipboardjs').length) {
                cb.contents().wrapAll('<div class="clipboardjs"></div>');
            }
            var clipButton = '<div class="btn-copy-wrap"><button type="button" class="clipbtn" aria-label="Copy code to clipboard"><i class="far fa-copy" aria-hidden="true"></i></button></div>';
            $(this).append(clipButton);
        }
    });

    // Fallback for standalone .clipboardjs not inside .codeblock-pre
    $('.clipboardjs').each(function () {
        if (!$(this).closest('.codeblock-pre').length && !$(this).parent().find('.clipbtn').length) {
            var clipButton = '<div class="btn-copy-wrap"><button type="button" class="clipbtn" aria-label="Copy code to clipboard"><i class="far fa-copy" aria-hidden="true"></i></button></div>';
            $(this).after(clipButton);
        }
    });

    // 2. Standard Markdown fenced blocks (.highlight)
    $('.highlight').each(function () {
        if (!$(this).find('.clipbtn').length) {
            var clipButton = '<div class="btn-copy-wrap highlight-copy-wrap"><button type="button" class="clipbtn" aria-label="Copy code to clipboard"><i class="far fa-copy" aria-hidden="true"></i></button></div>';
            $(this).append(clipButton);
        }
    });

    // 3. Standalone pre blocks (e.g. untagged markdown code blocks, REST API explorer, bare HTML pre)
    $('pre:not(.codeblock-pre)').each(function () {
        if ($(this).closest('.highlight').length) return;
        if ($(this).closest('.codeblock-pre').length) return;
        if ($(this).find('.clipbtn').length) return;
        var clipButton = '<div class="btn-copy-wrap highlight-copy-wrap"><button type="button" class="clipbtn" aria-label="Copy code to clipboard"><i class="far fa-copy" aria-hidden="true"></i></button></div>';
        $(this).css('position', 'relative').append(clipButton);
    });

    // 4. REST API Try-it endpoint URLs
    $('.rest-api-try-item__url').each(function () {
        if (!$(this).parent().find('.clipbtn').length) {
            $(this).wrap('<div class="rest-api-try-item__box"></div>');
            var clipButton = '<div class="btn-copy-wrap"><button type="button" class="clipbtn" aria-label="Copy URL to clipboard"><i class="far fa-copy" aria-hidden="true"></i></button></div>';
            $(this).after(clipButton);
        }
    });

    var clipboard = new Clipboard('.clipbtn', {
        text: function (trigger) {
            var container = trigger.closest('pre') || trigger.closest('.highlight') || trigger.closest('.codeblock-pre') || trigger.closest('.rest-api-example') || trigger.closest('.rest-api-try-item__box') || trigger.closest('.rest-api-try-item');
            if (!container) return '';
            var content = container.querySelector('.clipboardjs') || container.querySelector('code') || container.querySelector('.codeblock') || container;
            var clone = content.cloneNode(true);
            var btn = clone.querySelector('.btn-copy-wrap');
            if (btn) btn.remove();
            var text = clone.textContent || '';
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
            button.setAttribute('aria-label', 'Copy code to clipboard');
            button.removeAttribute("data-is-copying");
            if (button.isConnected === false) {
                return;
            }
            if (text.parentNode === button) {
                button.replaceChild(originalIcon, text);
            }
        }, 2000);
    });
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

