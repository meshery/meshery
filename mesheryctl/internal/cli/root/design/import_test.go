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
			name:    "given invalid source type when design import then error is thrown",
			args:    args{"invalid source type", "file.yaml", "", false},
			want:    ErrInValidSource("invalid source type", validDesignSourceTypes),
			wantErr: true,
		},
		{
			name:    "given missing file flag when design import then error is thrown",
			args:    args{"helm", "", "", false},
			want:    ErrDesignFileNotProvided(),
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
