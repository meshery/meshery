(function ($) {
    'use strict';

    function escapeRegExp(string) {
        return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    }

    function escapeHTML(str) {
        return $('<div>').text(str).html();
    }

    function escapeLunrSearch(string) {
        return string.replace(/[+\-&|!(){}\[\]^"~*?:\\]/g, '\\$&');
    }

    function getSafeHighlightedText(text, query) {
        if (!query) return escapeHTML(text);
        
        const escapedQuery = escapeRegExp(query);
        const regex = new RegExp(`(${escapedQuery})`, "gi");
        
        return text.split(regex).map(part => {
            if (part.toLowerCase() === query.toLowerCase()) {
                return `<span style='color:#00b39f; font-weight:bold;'>${escapeHTML(part)}</span>`;
            }
            return escapeHTML(part);
        }).join('');
    }

    $(document).ready(function () {
        const urlParams = new URLSearchParams(window.location.search);
        const searchQuery = urlParams.get('q');
        const $resultsContainer = $('#search-results');
        const $loading = $('#search-loading');
        const $title = $('#search-title');
        
        const searchIndexElement = document.querySelector('[data-search-index]');
        const searchIndexUrl = (searchIndexElement && searchIndexElement.getAttribute('data-search-index'))
                             ? searchIndexElement.getAttribute('data-search-index')
                             : "/offline-search-index.json";

        if (!searchQuery) {
            $loading.hide();
            $resultsContainer.html('<p>Please enter a search term.</p>');
            return;
        }

        $title.text(`Showing results for "${searchQuery}"`);

        $.ajax({
            url: searchIndexUrl,
            dataType: "json" 
        }).then((data) => {
            const dataMap = new Map(data.map(item => [item.ref, item]));

            const idx = lunr(function () {
                this.ref('ref');
                this.field('title', { boost: 100 });
                this.field('description', { boost: 50 });
                this.field('body');
                data.forEach((doc) => { this.add(doc); });
            });

            const safeLunrQuery = escapeLunrSearch(searchQuery);
            const results = idx.search(`${safeLunrQuery} ${safeLunrQuery}*`);
            $loading.hide();

            if (results.length === 0) {
                $resultsContainer.html('<p>No results found.</p>');
                return;
            }

            const $list = $('<div class="search-results-list" style="padding-left: 10px;">');

            results.forEach((r) => {
                const doc = dataMap.get(r.ref);
                
                if (!doc || doc.ref.includes("/project/releases/")) return;
                
                if (!doc.ref.startsWith("/")) return;

                let snippetRaw = doc.description || (doc.body ? doc.body.substring(0, 300) : "");
                if (snippetRaw.length > 300) snippetRaw = snippetRaw.substring(0, 300) + "...";

                const highlightedTitle = getSafeHighlightedText(doc.title, searchQuery);
                
                const $item = $('<div class="search-result-item" style="margin-bottom: 2rem;">');
                
                const $titleElem = $('<h3 style="margin-bottom: 0.5rem; font-size: 1.3rem; line-height: 1.4;">');
                const $bullet = $('<span style="color: #00b39f; margin-right: 8px;">&bull;</span>');
                
                const $link = $('<a>')
                    .attr('href', doc.ref)
                    .attr('style', 'color: #00b39f; text-decoration: none; font-weight: 500;')
                    .html(highlightedTitle); 
                
                const $snippetElem = $('<p style="color: #ccc; font-size: 1rem; line-height: 1.6; padding-left: 20px;">')
                    .text(snippetRaw);

                $titleElem.append($bullet).append($link);
                $item.append($titleElem).append($snippetElem);

                $list.append($item);
            });

            $resultsContainer.append($list);
        }).catch((err) => {
            console.error("Error loading search index:", err);
            $loading.hide();
            $resultsContainer.html('<p>Error loading search index. Please check the console.</p>');
        });
    });
})(jQuery);
