package helpers

import (
	"context"
	"testing"

	"github.com/meshery/meshery/server/models"
)

func TestNewAdaptersTracker(t *testing.T) {
	tests := []struct {
		name        string
		adapterURLs []string
		want        []models.Adapter
	}{
		{
			name:        "empty adapter URL list",
			adapterURLs: nil,
			want:        []models.Adapter{},
		},
		{
			name:        "initializes adapters by location",
			adapterURLs: []string{"meshery-istio:10000", "meshery-linkerd:10001"},
			want: []models.Adapter{
				{Location: "meshery-istio:10000"},
				{Location: "meshery-linkerd:10001"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewAdaptersTracker(tt.adapterURLs)

			assertAdaptersEqual(t, tracker.GetAdapters(context.Background()), tt.want)
		})
	}
}

func TestAdaptersTrackerLifecycle(t *testing.T) {
	ctx := context.Background()
	tracker := NewAdaptersTracker([]string{"meshery-istio:10000"})

	assertAdaptersEqual(t, tracker.GetAdapters(ctx), []models.Adapter{
		{Location: "meshery-istio:10000"},
	})

	tracker.AddAdapter(ctx, models.Adapter{
		Location:     "meshery-linkerd:10001",
		Name:         "meshery-linkerd",
		Version:      "v1.0.0",
		GitCommitSHA: "abc123",
	})
	assertAdaptersEqual(t, tracker.GetAdapters(ctx), []models.Adapter{
		{Location: "meshery-istio:10000"},
		{
			Location:     "meshery-linkerd:10001",
			Name:         "meshery-linkerd",
			Version:      "v1.0.0",
			GitCommitSHA: "abc123",
		},
	})

	tracker.AddAdapter(ctx, models.Adapter{
		Location:     "meshery-linkerd:10001",
		Name:         "meshery-linkerd",
		Version:      "v1.1.0",
		GitCommitSHA: "def456",
	})
	assertAdaptersEqual(t, tracker.GetAdapters(ctx), []models.Adapter{
		{Location: "meshery-istio:10000"},
		{
			Location:     "meshery-linkerd:10001",
			Name:         "meshery-linkerd",
			Version:      "v1.1.0",
			GitCommitSHA: "def456",
		},
	})

	tracker.RemoveAdapter(ctx, models.Adapter{Location: "meshery-istio:10000"})
	assertAdaptersEqual(t, tracker.GetAdapters(ctx), []models.Adapter{
		{
			Location:     "meshery-linkerd:10001",
			Name:         "meshery-linkerd",
			Version:      "v1.1.0",
			GitCommitSHA: "def456",
		},
	})

	tracker.RemoveAdapter(ctx, models.Adapter{Location: "meshery-linkerd:10001"})
	assertAdaptersEqual(t, tracker.GetAdapters(ctx), []models.Adapter{})
}

func TestAdaptersTrackerGetAdaptersReturnsSnapshot(t *testing.T) {
	ctx := context.Background()
	tracker := NewAdaptersTracker([]string{"meshery-istio:10000"})

	adapters := tracker.GetAdapters(ctx)
	if len(adapters) != 1 {
		t.Fatalf("expected one adapter, got %d", len(adapters))
	}

	adapters[0].Location = "mutated-location"

	assertAdaptersEqual(t, tracker.GetAdapters(ctx), []models.Adapter{
		{Location: "meshery-istio:10000"},
	})
}

func assertAdaptersEqual(t *testing.T, got, want []models.Adapter) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expected %d adapters, got %d: %#v", len(want), len(got), got)
	}

	gotByLocation := make(map[string]models.Adapter, len(got))
	for _, adapter := range got {
		gotByLocation[adapter.Location] = adapter
	}

	for _, wantAdapter := range want {
		gotAdapter, ok := gotByLocation[wantAdapter.Location]
		if !ok {
			t.Fatalf("expected adapter at location %q, got %#v", wantAdapter.Location, got)
		}
		if gotAdapter.Name != wantAdapter.Name {
			t.Fatalf("expected adapter %q name %q, got %q", wantAdapter.Location, wantAdapter.Name, gotAdapter.Name)
		}
		if gotAdapter.Version != wantAdapter.Version {
			t.Fatalf("expected adapter %q version %q, got %q", wantAdapter.Location, wantAdapter.Version, gotAdapter.Version)
		}
		if gotAdapter.GitCommitSHA != wantAdapter.GitCommitSHA {
			t.Fatalf("expected adapter %q git SHA %q, got %q", wantAdapter.Location, wantAdapter.GitCommitSHA, gotAdapter.GitCommitSHA)
		}
	}
}
