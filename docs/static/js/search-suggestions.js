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

    // Compute the canonical index URL once from the single source of truth
    function getIndexUrl() {
        var el = document.querySelector('[data-offline-search-index-json-src]');
        if (el) {
            var src = el.getAttribute('data-offline-search-index-json-src');
            return resolveSitePath(src);
        }
        return resolveSitePath('/offline-search-index.json');
    }

    var fuse = null;
    var indexFetchPromise = null;

    function fetchSearchIndex() {
        if (indexFetchPromise) return indexFetchPromise;

        var url = getIndexUrl();
        indexFetchPromise = $.ajax({ url: url, dataType: "json" }).then(
            function (data) {
                if (typeof window.Fuse === 'undefined') {
                    console.error("Fuse.js library is not loaded.");
                    indexFetchPromise = null;
                    return false;
                }
                var options = {
                    includeMatches: true,
                    threshold: 0.3,
                    minMatchCharLength: 2,
                    keys: ['title']
                };
                fuse = new window.Fuse(data, options);
                return true;
            },
            function (err) {
                console.error("Error loading search index:", err);
                indexFetchPromise = null;
                fuse = null;
                return false;
            }
        );

        return indexFetchPromise;
    }

    function highlightMatches(text, matches, key) {
        if (!Array.isArray(matches)) {
            return escapeHTML(text);
        }
        var match = matches.find(function (m) { return m.key === key; });
        if (!match || !match.indices || match.indices.length === 0) {
            return escapeHTML(text);
        }

        var result = '';
        var lastIndex = 0;

        var mergedIndices = [];
        var current = [match.indices[0][0], match.indices[0][1]];
        for (var i = 1; i < match.indices.length; i++) {
            var next = match.indices[i];
            if (next[0] <= current[1] + 1) {
                current[1] = Math.max(current[1], next[1]);
            } else {
                mergedIndices.push(current);
                current = [next[0], next[1]];
            }
        }
        mergedIndices.push(current);

        for (var j = 0; j < mergedIndices.length; j++) {
            var start = mergedIndices[j][0];
            var end = mergedIndices[j][1];
            result += escapeHTML(text.substring(lastIndex, start));
            result += '<mark class="search-match">' + escapeHTML(text.substring(start, end + 1)) + '</mark>';
            lastIndex = end + 1;
        }
        result += escapeHTML(text.substring(lastIndex));
        return result;
    }

    $(document).ready(function () {
        var $searchInputs = $('.td-search--offline .td-search-input');
        if ($searchInputs.length === 0) return;

        // Per-input state stored via WeakMap to avoid cross-input interference
        var inputState = new WeakMap();
        var instanceCounter = 0;

        // Assign unique IDs and initialize per-input state
        $searchInputs.each(function () {
            var $input = $(this);
            var $container = $input.closest('.td-search--offline').find('.search-suggestions-dropdown');
            var uid = 'search-suggestions-' + (instanceCounter++);

            // Give the input a unique ID if it doesn't have one (for aria references)
            if (!$input.attr('id')) {
                $input.attr('id', uid + '-input');
            }

            $container.attr('id', uid)
                      .attr('role', 'listbox')
                      .attr('aria-label', 'Search suggestions');

            $input.attr('role', 'combobox')
                  .attr('aria-autocomplete', 'list')
                  .attr('aria-expanded', 'false')
                  .attr('aria-controls', uid);

            // Store per-input state
            inputState.set(this, {
                activeIndex: -1,
                debounceTimer: null
            });
        });

        function getState(inputEl) {
            return inputState.get(inputEl);
        }

        function showDropdown($input, $container) {
            $container.show();
            $input.attr('aria-expanded', 'true');
        }

        function hideDropdown($input, $container) {
            var state = getState($input[0]);
            $container.hide();
            $input.attr('aria-expanded', 'false').removeAttr('aria-activedescendant');
            // Reset selection state so stale activeIndex doesn't cause unexpected submit redirects
            if (state) {
                state.activeIndex = -1;
            }
            $container.find('.suggestion-item').removeClass('active').attr('aria-selected', 'false');
        }

        $searchInputs.on('focus', function () {
            var state = getState(this);
            if (state) state.activeIndex = -1;
            if (!fuse) fetchSearchIndex();
            var $this = $(this);
            var $suggestionsContainer = $this.closest('.td-search--offline').find('.search-suggestions-dropdown');
            if ($this.val().trim().length >= 2) {
                showDropdown($this, $suggestionsContainer);
            }
        });

        $searchInputs.on('input', function () {
            var $this = $(this);
            var state = getState(this);
            var query = $this.val().trim();
            var $suggestionsContainer = $this.closest('.td-search--offline').find('.search-suggestions-dropdown');

            if (query.length < 2) {
                hideDropdown($this, $suggestionsContainer);
                if (state) {
                    clearTimeout(state.debounceTimer);
                }
                return;
            }

            if (state) {
                clearTimeout(state.debounceTimer);
                state.debounceTimer = setTimeout(function () {
                    if (!fuse) {
                        fetchSearchIndex().then(function (ready) {
                            // Re-read the current input value to avoid rendering stale suggestions
                            var currentQuery = $this.val().trim();
                            if (ready && currentQuery.length >= 2) {
                                renderSuggestions(currentQuery, $this, $suggestionsContainer);
                            }
                        });
                    } else {
                        renderSuggestions(query, $this, $suggestionsContainer);
                    }
                }, 200);
            }
        });

        function renderSuggestions(query, $input, $suggestionsContainer) {
            if (!fuse) return;
            var state = getState($input[0]);
            var results = fuse.search(query);
            $suggestionsContainer.empty();
            if (state) state.activeIndex = -1;
            $input.removeAttr('aria-activedescendant');

            var listboxId = $suggestionsContainer.attr('id');

            if (results.length === 0) {
                $suggestionsContainer.append(
                    '<div class="suggestion-item no-results" role="option" aria-disabled="true" aria-selected="false">No results found</div>'
                );
                showDropdown($input, $suggestionsContainer);
                return;
            }

            var filteredResults = results
                .filter(function (result) {
                    var doc = result.item;
                    return doc.ref && doc.ref.startsWith("/") && !doc.ref.includes("/project/releases/");
                })
                .slice(0, 5);

            filteredResults.forEach(function (result, idx) {
                var doc = result.item;
                var optionId = listboxId + '-option-' + idx;

                var titleHtml = highlightMatches(doc.title, result.matches, 'title');

                var $item = $('<a class="suggestion-item">')
                    .attr('href', resolveSitePath(doc.ref))
                    .attr('tabindex', '-1')
                    .attr('data-index', idx)
                    .attr('id', optionId)
                    .attr('role', 'option')
                    .attr('aria-selected', 'false');

                var $title = $('<div class="suggestion-title">').html(titleHtml);
                $item.append($title);

                $suggestionsContainer.append($item);
            });

            if ($suggestionsContainer.children().length > 0) {
                showDropdown($input, $suggestionsContainer);
            } else {
                hideDropdown($input, $suggestionsContainer);
            }
        }

        $searchInputs.on('keydown', function (e) {
            var $this = $(this);
            var state = getState(this);
            var $suggestionsContainer = $this.closest('.td-search--offline').find('.search-suggestions-dropdown');

            if (e.key === 'Escape') {
                hideDropdown($this, $suggestionsContainer);
                return;
            }

            var $items = $suggestionsContainer.find('a.suggestion-item');

            if ($items.length === 0 || !$suggestionsContainer.is(':visible')) return;
            if (!state) return;

            if (e.key === 'ArrowDown') {
                e.preventDefault();
                state.activeIndex = (state.activeIndex + 1) % $items.length;
                updateActiveItem($this, $items, state);
            } else if (e.key === 'ArrowUp') {
                e.preventDefault();
                state.activeIndex = (state.activeIndex - 1 + $items.length) % $items.length;
                updateActiveItem($this, $items, state);
            } else if (e.key === 'Enter') {
                if (state.activeIndex >= 0 && state.activeIndex < $items.length) {
                    e.preventDefault();
                    window.location.href = $items.eq(state.activeIndex).attr('href');
                }
            }
        });

        function updateActiveItem($input, $items, state) {
            $items.removeClass('active').attr('aria-selected', 'false');
            if (state.activeIndex >= 0) {
                var $active = $items.eq(state.activeIndex);
                $active.addClass('active').attr('aria-selected', 'true');
                $input.attr('aria-activedescendant', $active.attr('id'));
            } else {
                $input.removeAttr('aria-activedescendant');
            }
        }

        $(document).on('click', function (e) {
            if (!$(e.target).closest('.td-search--offline').length) {
                $searchInputs.each(function () {
                    var $input = $(this);
                    var $container = $input.closest('.td-search--offline').find('.search-suggestions-dropdown');
                    hideDropdown($input, $container);
                });
            }
        });

        $('.td-sidebar__search').on('submit', function (e) {
            var $suggestionsContainer = $(this).find('.td-search--offline .search-suggestions-dropdown');
            var $items = $suggestionsContainer.find('a.suggestion-item');
            var $input = $(this).find('.td-search-input');

            // Only redirect to a suggestion if the dropdown is actually visible
            if (!$suggestionsContainer.is(':visible')) return;

            var state = $input.length ? getState($input[0]) : null;
            if (state && state.activeIndex >= 0 && state.activeIndex < $items.length) {
                e.preventDefault();
                e.stopImmediatePropagation();
                window.location.href = $items.eq(state.activeIndex).attr('href');
            }
        });
    });

})(jQuery);
