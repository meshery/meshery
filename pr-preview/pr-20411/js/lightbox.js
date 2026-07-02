(function initLightbox() {
    if (document.getElementById("lightbox")) {
        return;
    }

    var lightbox = document.createElement("div");
    lightbox.setAttribute("id", "lightbox");

    var closeBtn = document.createElement("button");
    closeBtn.setAttribute("id", "lightbox-close");
    closeBtn.setAttribute("aria-label", "Close lightbox");

    var img = document.createElement("img");
    img.setAttribute("id", "lightbox-img");

    var caption = document.createElement("div");
    caption.setAttribute("id", "lightbox-caption");

    lightbox.appendChild(closeBtn);
    lightbox.appendChild(img);
    lightbox.appendChild(caption);
    document.body.appendChild(lightbox);

    function closeLightbox() {
        lightbox.classList.remove("active");
        document.body.style.overflow = "";
    }

    function openLightbox(src, alt) {
        img.setAttribute("src", src);
        img.setAttribute("alt", alt);
        caption.textContent = alt;
        caption.style.display = alt ? "block" : "none";
        lightbox.classList.add("active");
        document.body.style.overflow = "hidden";
    }

    function isLightboxImage(element) {
        if (element.tagName !== "IMG") return false;
        if (element.dataset.noLightbox !== undefined) return false;
        if (element.height < 50 && element.width < 50) return false;
        
        var contentSelectors = [
            "article",
            ".content",
            ".main-content",
            "main",
            ".td-content",
            ".page-content"
        ];
        
        for (var i = 0; i < contentSelectors.length; i++) {
            if (element.closest(contentSelectors[i])) {
                return true;
            }
        }
        
        return false;
    }

    document.addEventListener("click", function (e) {
        if (isLightboxImage(e.target)) {
            e.preventDefault();
            e.stopPropagation();
            openLightbox(e.target.getAttribute("src"), e.target.getAttribute("alt") || "");
        }
    });

    lightbox.addEventListener("click", function (e) {
        if (e.target === lightbox || e.target.id === "lightbox-caption") {
            closeLightbox();
        }
    });

    closeBtn.addEventListener("click", function (e) {
        e.preventDefault();
        e.stopPropagation();
        closeLightbox();
    });

    document.addEventListener("keydown", function (e) {
        if (e.key === "Escape" && lightbox.classList.contains("active")) {
            e.preventDefault();
            closeLightbox();
        }
    });
})();

if (window.addEventListener) {
    window.addEventListener("popstate", function() {
        setTimeout(initLightbox, 100);
    });
}


