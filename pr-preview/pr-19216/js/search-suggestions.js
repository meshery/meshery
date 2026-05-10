(function ($) {
    'use strict';

    function getSiteRootUrl() {
        const rootLink = document.querySelector("a.navbar-brand, #td-section-nav a.tree-root");
        const href = rootLink ? rootLink.getAttribute("href") : "/";
        return new URL(href || "/", window.location.href);
    }

    function resolveSitePath(pathname) {
        if (!pathname || typeof pathname !== "string") return pathname;
        if (/^(?:[a-z]+:)?\/\//i.test(pathname)) return pathname;
        if (/^(?:data:|mailto:|tel:|#)/i.test(pathname)) return pathname;
        if (!pathname.startsWith("/")) return pathname;
        return new URL(pathname.slice(1), getSiteRootUrl()).toString();
    }

    function escapeHTML(str) {
        return $('<div>').text(str).html();
    }

    let fuse = null;
    let indexFetchPromise = null;

    function fetchSearchIndex() {
        if (indexFetchPromise) return indexFetchPromise;

        const url = resolveSitePath("/offline-search-index.json");
        indexFetchPromise = $.ajax({ url: url, dataType: "json" }).then(function (data) {
            const options = {
                includeMatches: true,
                threshold: 0.3,
                minMatchCharLength: 2,
                keys: [
                    { name: 'title', weight: 0.7 },
                    { name: 'description', weight: 0.2 },
                    { name: 'body', weight: 0.1 }
                ]
            };
            fuse = new Fuse(data, options);
            return true;
        }).catch(function (err) {
            console.error("Error loading search index:", err);
            return false;
        });

        return indexFetchPromise;
    }

    function highlightMatches(text, matches, key) {
        const match = matches.find(m => m.key === key);
        if (!match || !match.indices || match.indices.length === 0) {
            return escapeHTML(text);
        }

        let result = '';
        let lastIndex = 0;
        
        let mergedIndices = [];
        let current = [...match.indices[0]];
        for (let i = 1; i < match.indices.length; i++) {
            const next = match.indices[i];
            if (next[0] <= current[1] + 1) {
                current[1] = Math.max(current[1], next[1]);
            } else {
                mergedIndices.push(current);
                current = [...next];
            }
        }
        mergedIndices.push(current);

        for (const [start, end] of mergedIndices) {
            result += escapeHTML(text.substring(lastIndex, start));
            result += `<mark class="search-match">${escapeHTML(text.substring(start, end + 1))}</mark>`;
            lastIndex = end + 1;
        }
        result += escapeHTML(text.substring(lastIndex));
        return result;
    }

    $(document).ready(function () {
        const $searchInputs = $('.td-search-input');
        if ($searchInputs.length === 0) return;

        let activeIndex = -1;

        $searchInputs.on('focus', function() {
            activeIndex = -1;
            if (!fuse) fetchSearchIndex();
            const $this = $(this);
            const $suggestionsContainer = $this.siblings('.search-suggestions-dropdown');
            if ($this.val().trim().length >= 2) {
                $suggestionsContainer.show();
            }
        });

        $searchInputs.on('input', function() {
            const $this = $(this);
            const query = $this.val().trim();
            const $suggestionsContainer = $this.siblings('.search-suggestions-dropdown');

            if (query.length < 2) {
                $suggestionsContainer.hide();
                return;
            }

            if (!fuse) {
                fetchSearchIndex().then(ready => {
                    if (ready) renderSuggestions(query, $suggestionsContainer);
                });
            } else {
                renderSuggestions(query, $suggestionsContainer);
            }
        });

        function renderSuggestions(query, $suggestionsContainer) {
            if (!fuse) return;
            const results = fuse.search(query);
            $suggestionsContainer.empty();
            activeIndex = -1;

            if (results.length === 0) {
                $suggestionsContainer.append('<div class="suggestion-item no-results">No results found</div>');
                $suggestionsContainer.show();
                return;
            }

            const filteredResults = results
                .filter(result => {
                    const doc = result.item;
                    return doc.ref && doc.ref.startsWith("/") && !doc.ref.includes("/project/releases/");
                })
                .slice(0, 5);

            filteredResults.forEach((result, idx) => {
                const doc = result.item;

                const titleHtml = highlightMatches(doc.title, result.matches, 'title');
                
                const $item = $('<a class="suggestion-item">')
                    .attr('href', resolveSitePath(doc.ref))
                    .attr('data-index', idx);
                
                const $title = $('<div class="suggestion-title">').html(titleHtml);
                $item.append($title);
                
                $suggestionsContainer.append($item);
            });

            if ($suggestionsContainer.children().length > 0) {
                $suggestionsContainer.show();
            } else {
                $suggestionsContainer.hide();
            }
        }

        $searchInputs.on('keydown', function(e) {
            const $this = $(this);
            const $suggestionsContainer = $this.siblings('.search-suggestions-dropdown');
            const $items = $suggestionsContainer.find('a.suggestion-item');
            
            if ($items.length === 0 || !$suggestionsContainer.is(':visible')) return;

            if (e.key === 'ArrowDown') {
                e.preventDefault();
                activeIndex = (activeIndex + 1) % $items.length;
                updateActiveItem($items);
            } else if (e.key === 'ArrowUp') {
                e.preventDefault();
                activeIndex = (activeIndex - 1 + $items.length) % $items.length;
                updateActiveItem($items);
            } else if (e.key === 'Enter') {
                if (activeIndex >= 0 && activeIndex < $items.length) {
                    e.preventDefault();
                    window.location.href = $items.eq(activeIndex).attr('href');
                }
            } else if (e.key === 'Escape') {
                $suggestionsContainer.hide();
            }
        });

        function updateActiveItem($items) {
            $items.removeClass('active');
            if (activeIndex >= 0) {
                $items.eq(activeIndex).addClass('active');
            }
        }

        $(document).on('click', function(e) {
            if (!$(e.target).closest('.td-search').length) {
                $('.search-suggestions-dropdown').hide();
            }
        });

        // Prevent form submission to avoid standard docsy search redirect
        $('.td-sidebar__search, .td-search').closest('form').on('submit', function(e) {
            e.preventDefault();
            const $suggestionsContainer = $(this).find('.search-suggestions-dropdown');
            const $items = $suggestionsContainer.find('a.suggestion-item');
            if (activeIndex >= 0 && activeIndex < $items.length) {
                window.location.href = $items.eq(activeIndex).attr('href');
            } else if ($items.length > 0) {
                window.location.href = $items.first().attr('href');
            }
        });
    });

})(jQuery);
