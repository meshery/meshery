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

if (toggleBtnSidebarNav) {
    toggleBtnSidebarNav.addEventListener("click",()=>{
        let sidebarNav=document.querySelector(".left-container")
        if(sidebarNav){
            sidebarNav.classList.toggle("left-container--active")
        }
    })
}

const toggleBtnMainNav=document.querySelector(".nav-toggle-btn--main");

if (toggleBtnMainNav) {
    toggleBtnMainNav.addEventListener("click",()=>{
        let sidebarNav=document.getElementById("main_navbar")
        if(sidebarNav){
            sidebarNav.classList.toggle("main-navbar--active")
        }
    })
}

document.addEventListener("click", (event) => {
    let sidebarNav = document.getElementById("main_navbar")
    if (sidebarNav && toggleBtnMainNav) {
        let isClickInsideSidebar = sidebarNav.contains(event.target)
        let isClickOnToggleButton = toggleBtnMainNav.contains(event.target)

        if (!isClickInsideSidebar && !isClickOnToggleButton) {
            sidebarNav.classList.remove("main-navbar--active")
        }
    }
})

window.addEventListener("load", () => {
    const docsSidebar = document.getElementById("td-section-nav");
    if (!docsSidebar) {
        return;
    }
    const sidebarMenu = document.getElementById("td-sidebar-menu");
    if (sidebarMenu) {
        sidebarMenu.classList.remove("d-none");
    }
    const stateKey = "meshery-docs-sidebar-fold-state-v1";
    const scrollKey = "meshery-docs-sidebar-scroll-v1";
    const sidebarScroller = document.querySelector(".sidebar-container");
    const foldInputs = Array.from(
        docsSidebar.querySelectorAll("li.with-child > input[type='checkbox'][id]")
    );

    const readStoredState = () => {
        try {
            return JSON.parse(sessionStorage.getItem(stateKey) || "{}");
        } catch (error) {
            return {};
        }
    };

    const persistState = () => {
        const nextState = {};
        foldInputs.forEach((input) => {
            nextState[input.id] = input.checked;
        });
        sessionStorage.setItem(stateKey, JSON.stringify(nextState));
    };

    const normalizePath = (value) => {
        try {
            const parsedUrl = new URL(value, window.location.origin);
            let path = parsedUrl.pathname || "/";
            path = path.replace(/\/{2,}/g, "/");
            return path.endsWith("/") ? path : `${path}/`;
        } catch (error) {
            return "/";
        }
    };

    const setBranchVisibility = (input) => {
        const branch = input ? input.closest("li.with-child") : null;
        const childList = branch ? branch.querySelector(":scope > ul") : null;
        if (!childList) {
            return;
        }
        childList.style.removeProperty("display");
    };

    const clearActiveState = () => {
        docsSidebar.querySelectorAll("a.td-sidebar-link.active").forEach((node) => node.classList.remove("active"));
        docsSidebar.querySelectorAll(".td-sidebar-nav-active-item").forEach((node) =>
            node.classList.remove("td-sidebar-nav-active-item")
        );
        docsSidebar.querySelectorAll("li.active-path").forEach((node) => node.classList.remove("active-path"));
    };

    const applyCurrentPathState = () => {
        const currentPath = normalizePath(window.location.pathname);
        const links = Array.from(docsSidebar.querySelectorAll("a.td-sidebar-link[href]")).filter((link) => {
            const hrefPath = normalizePath(link.getAttribute("href") || "");
            return hrefPath !== "/search/";
        });
        let bestLink = null;
        let bestScore = -1;
        links.forEach((link) => {
            const hrefPath = normalizePath(link.getAttribute("href") || "");
            let score = -1;
            if (hrefPath === currentPath) {
                score = 10000 + hrefPath.length;
            } else if (hrefPath !== "/" && currentPath.startsWith(hrefPath)) {
                score = hrefPath.length;
            }
            if (score > bestScore) {
                bestScore = score;
                bestLink = link;
            }
        });
        clearActiveState();
        if (!bestLink) {
            return;
        }
        bestLink.classList.add("active");
        const activeSpan = bestLink.querySelector("span");
        if (activeSpan) {
            activeSpan.classList.add("td-sidebar-nav-active-item");
        }
        let parentNode = bestLink.closest("li");
        while (parentNode) {
            parentNode.classList.add("active-path");
            parentNode = parentNode.parentElement ? parentNode.parentElement.closest("li") : null;
        }
    };

    const setToggleButtonState = (input) => {
        const branch = input ? input.closest("li.with-child") : null;
        const label = branch ? branch.querySelector(":scope > label") : null;
        const toggleButton = label ? label.querySelector(".sidebar-branch-toggle") : null;
        if (toggleButton) {
            toggleButton.setAttribute("aria-expanded", input.checked ? "true" : "false");
        }
    };

    const ensureToggleButtons = () => {
        foldInputs.forEach((input) => {
            const branch = input.closest("li.with-child");
            const label = branch ? branch.querySelector(":scope > label") : null;
            if (!label || label.querySelector(".sidebar-branch-toggle")) {
                setToggleButtonState(input);
                return;
            }
            const toggleButton = document.createElement("button");
            toggleButton.type = "button";
            toggleButton.className = "sidebar-branch-toggle";
            toggleButton.setAttribute("aria-label", "Toggle section");
            label.appendChild(toggleButton);
            setToggleButtonState(input);
        });
    };

    const storedState = readStoredState();
    const restoreSidebarState = () => {
        foldInputs.forEach((input) => {
            if (typeof storedState[input.id] === "boolean") {
                input.checked = storedState[input.id];
            }
            setBranchVisibility(input);
            setToggleButtonState(input);
        });

        docsSidebar
            .querySelectorAll("li.active-path > input[type='checkbox']")
            .forEach((input) => {
                input.checked = true;
                setBranchVisibility(input);
                setToggleButtonState(input);
            });
    };

    applyCurrentPathState();
    ensureToggleButtons();
    restoreSidebarState();

    const toggleInputState = (input) => {
        if (!input) {
            return;
        }
        input.checked = !input.checked;
        setBranchVisibility(input);
        setToggleButtonState(input);
        persistState();
    };

    docsSidebar.addEventListener("click", (event) => {
        const toggleButton = event.target.closest("button.sidebar-branch-toggle");
        if (toggleButton) {
            event.preventDefault();
            event.stopPropagation();
            const branch = toggleButton.closest("li.with-child");
            const branchInput = branch ? branch.querySelector(":scope > input[type='checkbox']") : null;
            toggleInputState(branchInput);
            return;
        }

        const link = event.target.closest("label > a.td-sidebar-link");
        if (!link) {
            return;
        }
        if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) {
            return;
        }

        const href = link.getAttribute("href");
        if (!href || href.startsWith("#")) {
            return;
        }

        event.preventDefault();
        event.stopPropagation();
        if (sidebarScroller) {
            sessionStorage.setItem(scrollKey, String(sidebarScroller.scrollTop));
        }
        window.location.assign(href);
    });

    foldInputs.forEach((input) => {
        input.addEventListener("change", () => {
            setBranchVisibility(input);
            setToggleButtonState(input);
            persistState();
        });
    });

    if (sidebarScroller) {
        const savedScroll = Number(sessionStorage.getItem(scrollKey) || "0");
        if (!Number.isNaN(savedScroll) && savedScroll > 0) {
            requestAnimationFrame(() => {
                sidebarScroller.scrollTop = savedScroll;
            });
        }

        let scrollRafId = null;
        const saveScroll = () => {
            if (scrollRafId) {
                return;
            }
            scrollRafId = window.requestAnimationFrame(() => {
                sessionStorage.setItem(scrollKey, String(sidebarScroller.scrollTop));
                scrollRafId = null;
            });
        };

        sidebarScroller.addEventListener("scroll", saveScroll, { passive: true });
        window.addEventListener("beforeunload", () => {
            sessionStorage.setItem(scrollKey, String(sidebarScroller.scrollTop));
        });
    }

    const sidebarSearchInput = document.querySelector(
        ".td-sidebar__search input[type='search'], .td-sidebar__search .td-search-input, #sidebar-search-input"
    );
    if (sidebarSearchInput && !sidebarSearchInput.id) {
        sidebarSearchInput.id = "sidebar-search-input";
    }

    const sidebarSearchForm = document.querySelector(".td-sidebar__search");
    if (sidebarSearchForm && sidebarSearchInput) {
        const baseHref = sidebarSearchInput.getAttribute("data-offline-search-base-href") || "/";
        const normalizedBase = baseHref.endsWith("/") ? baseHref : `${baseHref}/`;
        sidebarSearchForm.setAttribute("action", `${normalizedBase}search/`);
        sidebarSearchForm.setAttribute("method", "get");
        sidebarSearchInput.setAttribute("name", "q");
    }
});
