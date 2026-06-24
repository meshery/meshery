(function($) {
    'use strict';
    $(function() {
        $('[data-toggle="tooltip"]').tooltip();
        $('[data-toggle="popover"]').popover();
        $('.popover-dismiss').popover({
            trigger: 'focus'
        })
    });

    function bottomPos(element) {
        return element.offset().top + element.outerHeight();
    }
    $(function() {
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
        $(window).on('scroll', function() {
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
(function($) {
    'use strict';
    var Search = {
        init: function() {
            $(document).ready(function() {
                $(document).on('keypress', '.sidebar__search-input', function(e) {
                    if (e.keyCode !== 13) {
                        return
                    }
                    var query = $(this).val();
                    var searchPage = "/search/?q=" + query;
                    document.location = searchPage;
                    return false;
                });
            });
        },
    };
    Search.init();
}(jQuery));


/*popup-hidden-div*/

function HideToggleFunction() {
    var hide = document.getElementById("hiddendiv");
    if (hide.style.display === "block") {
      hide.style.display = "none";
    } else {
      hide.style.display = "block";
    }
  }

/*clipboard*/

var getcodeelement = $('.clipboardjs'); /*create custom id*/

getcodeelement.each(function(i) {  
    /*target*/
    var currentId = 'codeblock' + (i + 1);
    $(this).attr('id', currentId);

    /*trigger*/
    var text = $(this).text();
    text = text.replace(/\$ /gi, '')
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/'/g, "&#39;")
    .replace(/"/g, "&quot;");
    var clipButton = '<div class="btn-copy-wrap"><button class="clipbtn" data-clipboard-text="' + text + '" data-clipboard-target="#' + currentId + '"><i class="far fa-copy"></i></button></div>';
       $(this).after(clipButton);
});

var clipboard = new Clipboard('.clipbtn');

/* Change copy icon to text when successfully copied*/
clipboard.on("success", (e)=>{
    console.info(e.trigger);
    console.info(e.trigger.childNodes[0]);
    let originalIcon = e.trigger.childNodes[0];

    var icon = e.trigger.childNodes[0];
    var text = document.createElement('span');
    text.textContent = "Copied!";
    text.style.color = "white";

    e.trigger.replaceChild(text, icon);

    setTimeout(()=>{
        e.trigger.replaceChild(originalIcon, text);
    },2000)
})

const toggleBtnSidebarNav=document.querySelector(".nav-toggle-btn--document");

toggleBtnSidebarNav.addEventListener("click",()=>{
    let sidebarNav=document.querySelector(".left-container")
    if(sidebarNav){
        sidebarNav.classList.toggle("left-container--active")
    }
})

const toggleBtnMainNav=document.querySelector(".nav-toggle-btn--main");

toggleBtnMainNav.addEventListener("click",()=>{
    let sidebarNav=document.getElementById("main_navbar")
    if(sidebarNav){
        sidebarNav.classList.toggle("main-navbar--active")
    }
})

document.addEventListener("click", (event) => {
    let sidebarNav = document.getElementById("main_navbar")
    if (sidebarNav) {
        let isClickInsideSidebar = sidebarNav.contains(event.target)
        let isClickOnToggleButton = toggleBtnMainNav.contains(event.target)

        if (!isClickInsideSidebar && !isClickOnToggleButton) {
            sidebarNav.classList.remove("main-navbar--active")
        }
    }
})

// Bring focus to search bar when "/" is pressed
const searchInput = document.getElementById("sidebar-search-input");

if (searchInput) {
  document.addEventListener("keydown", function (event) {
    const activeEl = document.activeElement;
    const isEditable =
      activeEl.isContentEditable ||
      ["INPUT", "TEXTAREA"].includes(activeEl.tagName);

    if (
      event.key === "/" &&
      !event.ctrlKey &&
      !event.metaKey &&
      !event.altKey &&
      !isEditable
    ) {
      event.preventDefault();
      searchInput.focus();
    }
  });
}

/*
 * Version banner – injects a dismissible warning when the user is
 * viewing a pre-built static snapshot (e.g. /v0.8/, /v0.9/).
 */
;(function () {
    var LATEST_BASE = 'https://docs.meshery.io';
    var STORAGE_KEY = 'version-banner-dismissed';
    var VERSION_PREFIXES = ['v0.8', 'v0.9'];

    if (sessionStorage.getItem(STORAGE_KEY) === '1') return;

    var path = window.location.pathname;
    var matchedPrefix = null;

    for (var i = 0; i < VERSION_PREFIXES.length; i++) {
        var prefix = VERSION_PREFIXES[i];
        if (
            path === '/' + prefix ||
            path === '/' + prefix + '/' ||
            path.indexOf('/' + prefix + '/') === 0
        ) {
            matchedPrefix = prefix;
            break;
        }
    }

    if (!matchedPrefix) return;

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
        'background-color:#EBC017',
        'color:#1a1a1a',
        'font-size:0.875rem',
        'font-weight:500',
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
        '<a href="' + latestUrl + '" style="color:#1a1a1a;font-weight:700;text-decoration:underline;text-underline-offset:2px;margin-left:0.25rem;">View the latest version &rarr;</a>';

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
        sessionStorage.setItem(STORAGE_KEY, '1');
    });

    banner.appendChild(icon);
    banner.appendChild(text);
    banner.appendChild(closeBtn);

    document.body.insertBefore(banner, document.body.firstChild);

    var bh = banner.offsetHeight;
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
