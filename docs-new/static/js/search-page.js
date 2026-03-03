(function ($) {
    'use strict';

    // 1. HELPER: Escape Regex characters (Fixes Crash on symbols)
    function escapeRegExp(string) {
        return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    }

    // 2. HELPER: Sanitize HTML (Fixes XSS Security Risk)
    function escapeHTML(str) {
        return $('<div>').text(str).html();
    }

    $(document).ready(function () {
        const urlParams = new URLSearchParams(window.location.search);
        const searchQuery = urlParams.get('q');
        const $resultsContainer = $('#search-results');
        const $loading = $('#search-loading');
        const $title = $('#search-title');

        // 3. FIX: Get Dynamic Path (Fixes Hardcoded URL issue)
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

        $.ajax(searchIndexUrl).then((data) => {
            // 4. FIX: Create Map for Fast Lookup (Fixes Performance O(n) issue)
            const dataMap = new Map(data.map(item => [item.ref, item]));

            const idx = lunr(function () {
                this.ref('ref');
                this.field('title', { boost: 100 });
                this.field('description', { boost: 50 });
                this.field('body');
                data.forEach((doc) => { this.add(doc); });
            });

            // Search with escaped query
            const results = idx.search(`${searchQuery} ${searchQuery}*`);
            $loading.hide();

            if (results.length === 0) {
                $resultsContainer.html('<p>No results found.</p>');
                return;
            }

            const $list = $('<div class="search-results-list" style="padding-left: 10px;">');

            results.forEach((r) => {
                // Use Map for instant lookup
                const doc = dataMap.get(r.ref);
                
                // Guard: Skip if doc is missing or is a release note
                if (!doc || doc.ref.includes("/project/releases/")) return;

                let snippet = doc.description || doc.body.substring(0, 300);
                if (snippet.length > 300) snippet = snippet.substring(0, 300) + "...";

                // Prepare Logic
                const escapedQuery = escapeRegExp(searchQuery);
                const regex = new RegExp(`(${escapedQuery})`, "gi");

                // 5. FIX: Sanitize text BEFORE highlighting (Fixes XSS)
                const safeTitle = escapeHTML(doc.title);
                const safeSnippet = escapeHTML(snippet);

                // Apply Highlight to the Safe Title
                const highlightedTitle = safeTitle.replace(
                    regex,
                    "<span style='color:#00b39f; font-weight:bold;'>$1</span>"
                );

                const item = `
                    <div class="search-result-item" style="margin-bottom: 2rem;">
                      <h3 style="margin-bottom: 0.5rem; font-size: 1.3rem; line-height: 1.4;">
                          <span style="color: #00b39f; margin-right: 8px;">&bull;</span>
                          <a href="${doc.ref}" style="color: #00b39f; text-decoration: none; font-weight: 500;">${highlightedTitle}</a>
                      </h3>
                      <p style="color: #ccc; font-size: 1rem; line-height: 1.6; padding-left: 20px;">${safeSnippet}</p>
                    </div>`;
                
                $list.append(item);
            });

            $resultsContainer.append($list);
        }).catch((err) => {
            console.error("Error loading search index:", err);
            $loading.hide();
            $resultsContainer.html('<p>Error loading search index.</p>');
        });
    });
})(jQuery);