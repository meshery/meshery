document.addEventListener('DOMContentLoaded', function() {
    var searchInput = document.querySelector('.sidebar__search-input');
    var searchSuggestions = document.getElementById('search-suggestions');
    var fuse = null;
    var searchData = null;
    var isFetching = false;

    if (!searchInput || !searchSuggestions) return;

    // Utility to get base URL
    function getBaseUrl() {
        var baseUrl = window.location.origin;
        if (window.location.pathname.startsWith('/meshery/')) {
            baseUrl += '/meshery';
        }
        return baseUrl;
    }

    // Utility to HTML escape against XSS
    function escapeHTML(str) {
        if (!str) return "";
        return str.replace(/[&<>'"]/g, function(tag) {
            var charsToReplace = {
                '&': '&amp;',
                '<': '&lt;',
                '>': '&gt;',
                "'": '&#39;',
                '"': '&quot;'
            };
            return charsToReplace[tag] || tag;
        });
    }

    // Utility to prevent javascript: URI injection
    function isValidUrl(url) {
        if (!url) return false;
        // Block javascript:, vbscript:, data: URIs
        if (/^(javascript|vbscript|data):/i.test(url.trim())) {
            return false;
        }
        return true;
    }

    // Helper to highlight matched characters in the title
    function highlight(resultItem) {
        var item = resultItem.item;
        var title = escapeHTML(item.title || "");
        
        if (resultItem.matches) {
           var titleMatchObj = resultItem.matches.find(function(m) { return m.key === 'title'; });
           if (titleMatchObj && titleMatchObj.indices.length > 0) {
               var indices = titleMatchObj.indices;
               var res = "";
               var last = 0;
               indices.forEach(function(idx) {
                   res += title.substring(last, idx[0]);
                   res += "<mark class='search-highlight'>" + title.substring(idx[0], idx[1] + 1) + "</mark>";
                   last = idx[1] + 1;
               });
               res += title.substring(last);
               title = res;
           }
        }
        return title;
    }

    function renderResults(results) {
        if (!results || results.length === 0) {
            searchSuggestions.style.display = 'none';
            return;
        }

        var html = '<ul style="list-style: none; padding: 0; margin: 0;">';
        var maxResults = Math.min(results.length, 5);
        
        for (var i = 0; i < maxResults; i++) {
            var item = results[i].item;
            var url = item.url;
            if (!isValidUrl(url)) continue; // Skip items with malicious URLs
            
            var titleHTML = highlight(results[i]) || escapeHTML(item.title);
            
            var cats = item.categories ? item.categories.split(', ') : [];
            var catsHTML = '';
            if (cats.length > 0 && cats[0]) {
               catsHTML = '<span class="search-category-badge">' + escapeHTML(cats[0]) + '</span>';
            }

            html += '<li class="search-result-item">';
            html += '<a href="' + escapeHTML(url) + '" class="search-result-link">';
            html += '<div class="search-result-content">';
            html += '<span class="search-result-title">' + titleHTML + '</span>';
            html += catsHTML;
            html += '</div>';
            html += '</a></li>';
        }
        html += '</ul>';
        
        searchSuggestions.innerHTML = html;
        searchSuggestions.style.display = 'block';
    }

    function initSearch() {
        if (fuse || isFetching) return;
        isFetching = true;
        
        var baseUrl = getBaseUrl();

        fetch(baseUrl + '/search.json')
            .then(function(res) { return res.json(); })
            .then(function(data) {
                searchData = data;
                var options = {
                    includeScore: true,
                    includeMatches: true,
                    threshold: 0.3, // Fuzzy sensitivity (0.0 is perfect match, 1.0 is everything)
                    ignoreLocation: true,
                    useExtendedSearch: true,
                    keys: [
                        { name: 'title', weight: 2.0 },
                        { name: 'categories', weight: 1.0 },
                        { name: 'content', weight: 0.3 }
                    ]
                };
                fuse = new window.Fuse(data, options);
                isFetching = false;
                
                // If user typed while loading, trigger search immediately
                if (searchInput.value) {
                    var res = fuse.search(searchInput.value);
                    renderResults(res);
                }
            })
            .catch(function(err) {
                console.error("Error fetching search data:", err);
                isFetching = false;
            });
    }

    // Trigger loading of JSON on first click/focus to save bandwidth
    searchInput.addEventListener('focus', function() {
        initSearch();
        if (searchInput.value && fuse) {
            renderResults(fuse.search(searchInput.value));
        }
    });

    searchInput.addEventListener('input', function(e) {
        var query = e.target.value;
        if (!query) {
            searchSuggestions.style.display = 'none';
            return;
        }
        // Initialize if not done
        initSearch();
        
        if (fuse) {
            var res = fuse.search(query);
            renderResults(res);
        }
    });

    // Hijack 'Enter' to go to first result if dropdown is open
    searchInput.addEventListener('keydown', function(e) {
        if (e.key === 'Enter') {
            e.preventDefault();
            var firstResult = searchSuggestions.querySelector('a');
            if (firstResult && searchSuggestions.style.display === 'block') {
                window.location.href = firstResult.href;
            } else {
                // If nothing in suggestions, use standard search logic
                var searchPageUrl = getBaseUrl() + '/search/?q=' + encodeURIComponent(searchInput.value);
                window.location.href = searchPageUrl;
            }
        }
    });

    // Close dropdown on click outside
    document.addEventListener('click', function(e) {
        if (!searchInput.contains(e.target) && !searchSuggestions.contains(e.target)) {
            searchSuggestions.style.display = 'none';
        }
    });
});
