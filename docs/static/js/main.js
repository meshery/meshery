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
function copyTextToClipboard(text) {
    if (navigator.clipboard && window.isSecureContext) {
        return navigator.clipboard.writeText(text);
    }
    return new Promise((resolve, reject) => {
        var textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.setAttribute('readonly', '');
        textArea.style.position = 'absolute';
        textArea.style.left = '-9999px';
        document.body.appendChild(textArea);
        textArea.select();
        try {
            var successful = document.execCommand('copy');
            document.body.removeChild(textArea);
            if (successful) {
                resolve();
                return;
            }
            reject(new Error('Copy command was not successful'));
        } catch (error) {
            document.body.removeChild(textArea);
            reject(error);
        }
    });
}

function updateCopyButton(button, state) {
    button.setAttribute('data-copy-state', state);
    if (state === 'copied') {
        button.innerHTML = '<span>Copied!</span>';
        return;
    }
    if (state === 'error') {
        button.innerHTML = '<span>Press Ctrl+C</span>';
        return;
    }
    button.innerHTML = '<i class="far fa-copy" aria-hidden="true"></i>';
}

function getCodeBlockText(preElement) {
    var clipboardElement = preElement.querySelector('.clipboardjs');
    if (clipboardElement) {
        return clipboardElement.textContent;
    }
    var codeElement = preElement.querySelector('code');
    if (codeElement) {
        return codeElement.textContent;
    }
    return preElement.textContent;
}

function addCopyButtonsToCodeBlocks() {
    var codeBlocks = document.querySelectorAll('.td-content pre');
    codeBlocks.forEach(function(preElement) {
        if (preElement.dataset.copyButtonAttached === 'true') {
            return;
        }
        var codeText = getCodeBlockText(preElement);
        if (!codeText || !codeText.trim()) {
            return;
        }
        var buttonWrapper = document.createElement('div');
        buttonWrapper.className = 'btn-copy-wrap';
        var button = document.createElement('button');
        button.className = 'clipbtn';
        button.type = 'button';
        button.setAttribute('aria-label', 'Copy code to clipboard');
        updateCopyButton(button, 'default');
        button.addEventListener('click', function() {
            var textToCopy = getCodeBlockText(preElement);

            copyTextToClipboard(textToCopy).then(function() {
                updateCopyButton(button, 'copied');
                setTimeout(function() {
                    updateCopyButton(button, 'default');
                }, 2000);
            }).catch(function() {
                updateCopyButton(button, 'error');
                setTimeout(function() {
                    updateCopyButton(button, 'default');
                }, 2000);
            });
        });
        buttonWrapper.appendChild(button);
        preElement.appendChild(buttonWrapper);
        preElement.dataset.copyButtonAttached = 'true';
    });
}
document.addEventListener('DOMContentLoaded', addCopyButtonsToCodeBlocks);

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
