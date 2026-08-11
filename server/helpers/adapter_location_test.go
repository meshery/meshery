package helpers

import "testing"

func TestAdapterPort(t *testing.T) {
	tests := []struct {
		name     string
		location string
		want     string
		wantErr  bool
	}{
		{
			name:     "host port",
			location: "meshery-istio:10000",
			want:     "10000",
		},
		{
			name:     "http url",
			location: "http://localhost:11434",
			want:     "11434",
		},
		{
			name:     "https url",
			location: "https://adapter.example.com:10001",
			want:     "10001",
		},
		{
			name:     "ipv6 host port",
			location: "[::1]:10000",
			want:     "10000",
		},
		{
			name:     "http ipv6 url",
			location: "http://[::1]:10000",
			want:     "10000",
		},
		{
			name:     "missing port panics previously",
			location: "localhost",
			wantErr:  true,
		},
		{
			name:     "empty",
			location: "",
			wantErr:  true,
		},
		{
			name:     "scheme only",
			location: "http://",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AdapterPort(tt.location)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("AdapterPort(%q) error = nil, want error", tt.location)
				}
				return
			}
			if err != nil {
				t.Fatalf("AdapterPort(%q) unexpected error: %v", tt.location, err)
			}
			if got != tt.want {
				t.Fatalf("AdapterPort(%q) = %q, want %q", tt.location, got, tt.want)
			}
		})
	}
}
