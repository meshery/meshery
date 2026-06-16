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

        const siteRootPath = getSiteRootUrl().pathname.replace(/\/+$/, "/");
        const siteRootWithoutTrailingSlash = siteRootPath.replace(/\/$/, "");

        if (
            siteRootPath === "/" ||
            pathname === siteRootWithoutTrailingSlash ||
            pathname.startsWith(siteRootPath)
        ) {
            return new URL(pathname, window.location.origin).toString();
        }

        return new URL(pathname.slice(1), getSiteRootUrl()).toString();
    }

    function escapeHTML(str) {
        return $('<div>').text(str).html();
    }

    let fuse = null;
    let indexFetchPromise = null;

    function fetchSearchIndex(indexUrl) {
        if (indexFetchPromise) return indexFetchPromise;

        const url = indexUrl || resolveSitePath("/offline-search-index.json");
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
            indexFetchPromise = null;
            fuse = null;
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
        const $searchInputs = $('.td-search--offline .td-search-input');
        if ($searchInputs.length === 0) return;

        let activeIndex = -1;
        let debounceTimer = null;
        let instanceCounter = 0;

        // Assign unique IDs to each search input / dropdown pair for ARIA
        $searchInputs.each(function () {
            const $input = $(this);
            const $container = $input.closest('.td-search--offline').find('.search-suggestions-dropdown');
            const uid = 'search-suggestions-' + (instanceCounter++);

            $container.attr('id', uid).attr('role', 'listbox');
            $input.attr('role', 'combobox')
                  .attr('aria-autocomplete', 'list')
                  .attr('aria-expanded', 'false')
                  .attr('aria-controls', uid);
        });

        // Read the index URL from the data attribute set by Hugo
        function getIndexUrl() {
            const el = document.querySelector('[data-offline-search-index-json-src]');
            if (el) {
                const src = el.getAttribute('data-offline-search-index-json-src');
                return resolveSitePath(src);
            }
            return resolveSitePath('/offline-search-index.json');
        }

        function showDropdown($input, $container) {
            $container.show();
            $input.attr('aria-expanded', 'true');
        }

        function hideDropdown($input, $container) {
            $container.hide();
            $input.attr('aria-expanded', 'false').removeAttr('aria-activedescendant');
        }

        $searchInputs.on('focus', function() {
            activeIndex = -1;
            if (!fuse) fetchSearchIndex(getIndexUrl());
            const $this = $(this);
            const $suggestionsContainer = $this.closest('.td-search--offline').find('.search-suggestions-dropdown');
            if ($this.val().trim().length >= 2) {
                showDropdown($this, $suggestionsContainer);
            }
        });

        $searchInputs.on('input', function() {
            const $this = $(this);
            const query = $this.val().trim();
            const $suggestionsContainer = $this.closest('.td-search--offline').find('.search-suggestions-dropdown');

            if (query.length < 2) {
                hideDropdown($this, $suggestionsContainer);
                clearTimeout(debounceTimer);
                return;
            }

            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(function () {
                if (!fuse) {
                    fetchSearchIndex(getIndexUrl()).then(function (ready) {
                        if (ready) renderSuggestions(query, $this, $suggestionsContainer);
                    });
                } else {
                    renderSuggestions(query, $this, $suggestionsContainer);
                }
            }, 200);
        });

        function renderSuggestions(query, $input, $suggestionsContainer) {
            if (!fuse) return;
            const results = fuse.search(query);
            $suggestionsContainer.empty();
            activeIndex = -1;
            $input.removeAttr('aria-activedescendant');

            const listboxId = $suggestionsContainer.attr('id');

            if (results.length === 0) {
                $suggestionsContainer.append(
                    '<div class="suggestion-item no-results" role="option" aria-disabled="true">No results found</div>'
                );
                showDropdown($input, $suggestionsContainer);
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
                const optionId = listboxId + '-option-' + idx;

                const titleHtml = highlightMatches(doc.title, result.matches, 'title');
                
                const $item = $('<a class="suggestion-item">')
                    .attr('href', resolveSitePath(doc.ref))
                    .attr('data-index', idx)
                    .attr('id', optionId)
                    .attr('role', 'option');
                
                const $title = $('<div class="suggestion-title">').html(titleHtml);
                $item.append($title);
                
                $suggestionsContainer.append($item);
            });

            if ($suggestionsContainer.children().length > 0) {
                showDropdown($input, $suggestionsContainer);
            } else {
                hideDropdown($input, $suggestionsContainer);
            }
        }

        $searchInputs.on('keydown', function(e) {
            const $this = $(this);
            const $suggestionsContainer = $this.closest('.td-search--offline').find('.search-suggestions-dropdown');
            const $items = $suggestionsContainer.find('a.suggestion-item');
            
            if ($items.length === 0 || !$suggestionsContainer.is(':visible')) return;

            if (e.key === 'ArrowDown') {
                e.preventDefault();
                activeIndex = (activeIndex + 1) % $items.length;
                updateActiveItem($this, $items);
            } else if (e.key === 'ArrowUp') {
                e.preventDefault();
                activeIndex = (activeIndex - 1 + $items.length) % $items.length;
                updateActiveItem($this, $items);
            } else if (e.key === 'Enter') {
                if (activeIndex >= 0 && activeIndex < $items.length) {
                    e.preventDefault();
                    window.location.href = $items.eq(activeIndex).attr('href');
                }
            } else if (e.key === 'Escape') {
                hideDropdown($this, $suggestionsContainer);
            }
        });

        function updateActiveItem($input, $items) {
            $items.removeClass('active').attr('aria-selected', 'false');
            if (activeIndex >= 0) {
                const $active = $items.eq(activeIndex);
                $active.addClass('active').attr('aria-selected', 'true');
                $input.attr('aria-activedescendant', $active.attr('id'));
            } else {
                $input.removeAttr('aria-activedescendant');
            }
        }

        $(document).on('click', function(e) {
            if (!$(e.target).closest('.td-search').length) {
                $('.search-suggestions-dropdown').hide();
                $searchInputs.attr('aria-expanded', 'false').removeAttr('aria-activedescendant');
            }
        });

        $('.td-sidebar__search').on('submit', function(e) {
            const $suggestionsContainer = $(this).find('.td-search--offline .search-suggestions-dropdown');
            const $items = $suggestionsContainer.find('a.suggestion-item');

            if (activeIndex >= 0 && activeIndex < $items.length) {
                e.preventDefault();
                e.stopImmediatePropagation();
                window.location.href = $items.eq(activeIndex).attr('href');
            }
        });
    });

})(jQuery);
