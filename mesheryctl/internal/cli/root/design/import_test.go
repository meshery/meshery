package design

import (
	"testing"
)

func Test_importPattern_DisplayErrorsMissingFlags(t *testing.T) {
	type args struct {
		sourceType string
		file       string
		patternURL string
		save       bool
	}

	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name:    "Import missing source type flag",
			args:    args{"", "file.yaml", "", false},
			want:    ErrDesignSourceType(),
			wantErr: true,
		},
		{
			name:    "Import missing file flag",
			args:    args{"helm", "", "", false},
			want:    ErrDesignManifest(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := importPattern(tt.args.sourceType, tt.args.file, tt.args.patternURL, tt.args.save)
			if (err != nil) != tt.wantErr {
				t.Errorf("importPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
