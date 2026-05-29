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
                return `<span class="search-highlight">${escapeHTML(part)}</span>`;
            }
            return escapeHTML(part);
        }).join('');
    }

    let searchIndex = null;
    let searchDataMap = null;

    function initSearchIndex(url) {
        return $.ajax({ url: url, dataType: "json" }).then(function (data) {
            searchDataMap = new Map(data.map(item => [item.ref, item]));
            searchIndex = lunr(function () {
                this.ref('ref');
                this.field('title', { boost: 100 });
                this.field('description', { boost: 50 });
                this.field('body');
                data.forEach((doc) => { this.add(doc); });
            });
        });
    }

    function performSearch(query) {
        const $resultsContainer = $('#search-results');
        const $loading = $('#search-loading');
        const safeLunrQuery = escapeLunrSearch(query);
        const results = searchIndex.search(`${safeLunrQuery} ${safeLunrQuery}*`);
        $loading.hide();

        if (results.length === 0) {
            $resultsContainer.html('<p>No results found.</p>');
            return;
        }

        const $list = $('<div class="search-results-list">');

        results.forEach((r) => {
            const doc = searchDataMap.get(r.ref);

            if (!doc || doc.ref.includes("/project/releases/")) return;
            if (!doc.ref.startsWith("/")) return;

            let snippetRaw = doc.description || (doc.body ? doc.body.substring(0, 300) : "");
            if (snippetRaw.length > 300) snippetRaw = snippetRaw.substring(0, 300) + "...";

            const highlightedTitle = getSafeHighlightedText(doc.title, query);

            const $item      = $('<div class="search-result-item">');
            const $titleElem = $('<h3 class="search-result-title">');
            const $bullet    = $('<span class="search-result-bullet">&bull;</span>');
            const $link      = $('<a class="search-result-link">').attr('href', doc.ref).html(highlightedTitle);
            const $snippet   = $('<p class="search-result-snippet">').text(snippetRaw);

            $titleElem.append($bullet).append($link);
            $item.append($titleElem).append($snippet);
            $list.append($item);
        });

        if ($list.children().length === 0) {
            $resultsContainer.html('<p>No results found.</p>');
            return;
        }

        $resultsContainer.append($list);
    }

    $(document).ready(function () {
        const urlParams = new URLSearchParams(window.location.search);
        const searchQuery = urlParams.get('q');
        const $resultsContainer = $('#search-results');
        const $loading = $('#search-loading');
        const $title = $('#search-title');

        const el = document.querySelector('[data-search-index]');
        const searchIndexUrl = (el && el.getAttribute('data-search-index'))
            ? el.getAttribute('data-search-index')
            : "/offline-search-index.json";

        if (!searchQuery) {
            $loading.hide();
            $resultsContainer.html('<p>Please enter a search term.</p>');
            return;
        }

        $title.text(`Showing results for "${searchQuery}"`);

        initSearchIndex(searchIndexUrl).then(function () {
            performSearch(searchQuery);
        }).catch(function (err) {
            console.error("Error loading search index:", err);
            $loading.hide();
            $resultsContainer.html('<p>Error loading search index. Please check the console.</p>');
        });
    });
})(jQuery);
