document.addEventListener("DOMContentLoaded", function () {
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
        if ((e.key === "Escape" || e.keyCode === 27) && lightbox.classList.contains("active")) {
            e.preventDefault();
            closeLightbox();
        }
    });

    function attachLightbox() {
        var contentSelectors = [
            "article img",
            ".content img",
            ".main-content img",
            "img.content-image",
            "main img",
            ".td-content img",
            ".page-content img"
        ];
        
        contentSelectors.forEach(function(selector) {
            var images = document.querySelectorAll(selector);
            
            images.forEach(function (contentImg) {
                if (!contentImg.dataset.lightboxReady) {
                    if (contentImg.height < 50 && contentImg.width < 50) {
                        return;
                    }
                    
                    contentImg.dataset.lightboxReady = "true";
                    contentImg.style.cursor = "zoom-in";
                    contentImg.addEventListener("click", function (e) {
                        e.preventDefault();
                        e.stopPropagation();
                        openLightbox(this.getAttribute("src"), this.getAttribute("alt") || "");
                    });
                }
            });
        });
    }
    
    attachLightbox();
    
    setTimeout(attachLightbox, 500);
    setTimeout(attachLightbox, 1000);
});
