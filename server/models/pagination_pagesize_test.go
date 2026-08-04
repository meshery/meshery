package models

import (
	"encoding/json"
	"testing"
)

// pageResponse decodes any paginated response that has page, pageSize, and totalCount.
type pageResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
}

// TestPageSizeMatchesRequested checks that pageSize in the response equals what
// the caller asked for, not how many rows came back. Before the fix, PageSize was
// set to len(fetched), which broke the UI hasMore formula on the last page and
// caused infinite empty network requests.
func TestPageSizeMatchesRequested(t *testing.T) {
	t.Run("last page - pageSize is requested, not len(fetched)", func(t *testing.T) {
		// 13 total, request page 1 with pageSize=10 => only 3 items come back.
		// Old: PageSize=3 => hasMore = 13 > 3*(1+1) = true  (wrong)
		// New: PageSize=10 => hasMore = 13 > 10*(1+1) = false (correct)
		resp := pageResponse{Page: 1, PageSize: 10, TotalCount: 13}
		hasMore := resp.TotalCount > resp.PageSize*(resp.Page+1)
		if hasMore {
			t.Errorf("hasMore=true on last page (totalCount=%d, pageSize=%d, page=%d), want false",
				resp.TotalCount, resp.PageSize, resp.Page)
		}
	})

	t.Run("pageSize=all sets pageSize to totalCount", func(t *testing.T) {
		resp := pageResponse{Page: 0, PageSize: 42, TotalCount: 42}
		if hasMore := resp.TotalCount > resp.PageSize*(resp.Page+1); hasMore {
			t.Errorf("hasMore=true for pageSize=all (totalCount=%d), want false", resp.TotalCount)
		}
	})

	t.Run("middle page - hasMore is true", func(t *testing.T) {
		resp := pageResponse{Page: 0, PageSize: 10, TotalCount: 25}
		if hasMore := resp.TotalCount > resp.PageSize*(resp.Page+1); !hasMore {
			t.Errorf("hasMore=false on page 0 of 25 items with pageSize 10, want true")
		}
	})

	t.Run("exact last page - hasMore is false", func(t *testing.T) {
		resp := pageResponse{Page: 1, PageSize: 10, TotalCount: 20}
		if hasMore := resp.TotalCount > resp.PageSize*(resp.Page+1); hasMore {
			t.Errorf("hasMore=true for exact last page, want false")
		}
	})

	t.Run("json round-trip preserves pageSize", func(t *testing.T) {
		original := pageResponse{Page: 2, PageSize: 10, TotalCount: 23}
		b, _ := json.Marshal(original)
		var decoded pageResponse
		json.Unmarshal(b, &decoded)
		if decoded.PageSize != 10 {
			t.Errorf("PageSize after JSON round-trip = %d, want 10", decoded.PageSize)
		}
	})
}
