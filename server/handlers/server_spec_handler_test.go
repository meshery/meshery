package handlers

import "testing"

func TestIsCoreVersionOutdated(t *testing.T) {
	tests := []struct {
		name          string
		serverVersion string
		latestVersion string
		outdated      bool
	}{
		{
			name:          "kanvas revision suffix is not outdated against the same core release",
			serverVersion: "v1.0.82-1",
			latestVersion: "v1.0.82",
			outdated:      false,
		},
		{
			name:          "later kanvas revision on the same core is not outdated",
			serverVersion: "v1.0.82-2",
			latestVersion: "v1.0.82",
			outdated:      false,
		},
		{
			name:          "exact stable match is not outdated",
			serverVersion: "v1.0.82",
			latestVersion: "v1.0.82",
			outdated:      false,
		},
		{
			name:          "newer core than the published release is not outdated",
			serverVersion: "v1.0.83-1",
			latestVersion: "v1.0.82",
			outdated:      false,
		},
		{
			name:          "older kanvas revision is outdated against a newer core release",
			serverVersion: "v1.0.82-1",
			latestVersion: "v1.0.83",
			outdated:      true,
		},
		{
			name:          "older stable release is outdated",
			serverVersion: "v1.0.81",
			latestVersion: "v1.0.82",
			outdated:      true,
		},
		{
			name:          "older kanvas revision is outdated against a later core",
			serverVersion: "v1.0.81-1",
			latestVersion: "v1.0.82",
			outdated:      true,
		},
		{
			name:          "v-prefix on only one side still compares the same core",
			serverVersion: "1.0.82-1",
			latestVersion: "v1.0.82",
			outdated:      false,
		},
		{
			name:          "build metadata on the same core is not outdated",
			serverVersion: "v1.0.82+ca86624a",
			latestVersion: "v1.0.82",
			outdated:      false,
		},
		{
			name:          "unparsable versions fall back to exact string inequality",
			serverVersion: "edge-latest",
			latestVersion: "v1.0.82",
			outdated:      true,
		},
		{
			name:          "identical unparsable versions are not outdated",
			serverVersion: "dev",
			latestVersion: "dev",
			outdated:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCoreVersionOutdated(tt.serverVersion, tt.latestVersion)
			if got != tt.outdated {
				t.Errorf("isCoreVersionOutdated(%q, %q) = %v, want %v", tt.serverVersion, tt.latestVersion, got, tt.outdated)
			}
		})
	}
}
