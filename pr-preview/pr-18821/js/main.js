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

toggleBtnSidebarNav.addEventListener("click",()=>{
    let sidebar=document.querySelector(".left-container")
    if(sidebar){
        isActiveState = localStorage.getItem('sidebar-state')
        if (isActiveState === 'active') {
            localStorage.setItem('sidebar-state', 'inactive')
        } else {
            localStorage.setItem('sidebar-state', 'active')
        }
        sidebar.classList.toggle('left-container--active')
    }
})

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
    if (sidebarNav) {
        let isClickInsideSidebar = sidebarNav.contains(event.target)
        let isClickOnToggleButton = toggleBtnMainNav.contains(event.target)

        if (!isClickInsideSidebar && !isClickOnToggleButton) {
            sidebarNav.classList.remove("main-navbar--active")
        }
    }
});


;(function() {
  'use strict';
  
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initResize);
  } else {
    initResize();
  }
  
  function initResize() {
    const MIN_SIDEBAR_WIDTH = 250;
    const MAX_SIDEBAR_WIDTH = 400;
    const RESIZE_HANDLE_WIDTH = 10;
    const MOBILE_BREAKPOINT = 1200;

    const leftContainer = document.querySelector('.left-container');
    const sidebar = document.querySelector('.sidebar-container');

    if (!leftContainer || !sidebar) return;
        
    let isResizing = false;

    leftContainer.addEventListener('mousedown', function(e) {
      if(window.innerWidth <= MOBILE_BREAKPOINT) return;

      const rect = leftContainer.getBoundingClientRect();
      if (e.clientX >= rect.right - RESIZE_HANDLE_WIDTH) {
        isResizing = true;
        document.body.style.cursor = 'col-resize';
        document.body.style.userSelect = 'none';
        e.preventDefault();
      }
    });

    document.addEventListener('mousemove', function(e) {
  if (!isResizing) return;

  requestAnimationFrame(function() {
    const containerRect = leftContainer.getBoundingClientRect();
    const newWidth = Math.min(Math.max(e.clientX - containerRect.left, 250), 400);

    leftContainer.style.width = newWidth + 'px';
    leftContainer.style.flexBasis = newWidth + 'px';
    leftContainer.style.flexShrink = '0';
    
    const wrapper = leftContainer.querySelector('.sidebar-container-wrapper');
    if (wrapper) {
      wrapper.style.width = newWidth + 'px';
    }
    
    sidebar.style.width = newWidth + 'px';
  });
});

    document.addEventListener('mouseup', function() {
      if (!isResizing) return;
      isResizing = false;
      document.body.style.cursor = '';
      document.body.style.userSelect = '';
    });
  }
})();