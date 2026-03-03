(function ($) {
    'use strict';

    // Helper: Escape Regex characters (For Highlighting)
    function escapeRegExp(string) {
        return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    }

    // Helper: Sanitize HTML (For XSS Security)
    function escapeHTML(str) {
        return $('<div>').text(str).html();
    }

    // Helper: Escape Lunr Special Characters (Fixes Crash on ':', '~', '^')
    function escapeLunrSearch(string) {
        return string.replace(/[+\-&|!(){}\[\]^"~*?:\\]/g, '\\$&');
    }

    // Helper: Safe Highlight Strategy
    function getSafeHighlightedText(text, query) {
        if (!query) return escapeHTML(text);
        
        const escapedQuery = escapeRegExp(query);
        const regex = new RegExp(`(${escapedQuery})`, "gi");
        
        // Split text by the query, map over parts to escape them, then wrap matches
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
        
        // Dynamic Index Path
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
            dataType: "json" // Explicitly expect JSON
        }).then((data) => {
            const dataMap = new Map(data.map(item => [item.ref, item]));

            const idx = lunr(function () {
                this.ref('ref');
                this.field('title', { boost: 100 });
                this.field('description', { boost: 50 });
                this.field('body');
                data.forEach((doc) => { this.add(doc); });
            });

            // Use the Lunr-Specific Escape function here to prevent crashes
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
                
                // Guard: Skip missing docs, release notes, or unsafe URLs
                if (!doc || doc.ref.includes("/project/releases/")) return;
                
                // Security: Ensure URL starts with / to prevent external redirects
                if (!doc.ref.startsWith("/")) return;

                // Guard: Handle missing body/description gracefully
                let snippetRaw = doc.description || (doc.body ? doc.body.substring(0, 300) : "");
                if (snippetRaw.length > 300) snippetRaw = snippetRaw.substring(0, 300) + "...";

                // Safe Highlighting
                const highlightedTitle = getSafeHighlightedText(doc.title, searchQuery);
                
                // DOM Construction (Safe against XSS)
                const $item = $('<div class="search-result-item" style="margin-bottom: 2rem;">');
                
                const $titleElem = $('<h3 style="margin-bottom: 0.5rem; font-size: 1.3rem; line-height: 1.4;">');
                const $bullet = $('<span style="color: #00b39f; margin-right: 8px;">&bull;</span>');
                
                // IMPORTANT: We use .attr() and .html() here safely
                const $link = $('<a>')
                    .attr('href', doc.ref)
                    .attr('style', 'color: #00b39f; text-decoration: none; font-weight: 500;')
                    .html(highlightedTitle); // Safe because we escaped parts in getSafeHighlightedText
                
                const $snippetElem = $('<p style="color: #ccc; font-size: 1rem; line-height: 1.6; padding-left: 20px;">')
                    .text(snippetRaw); // .text() automatically escapes HTML for us

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