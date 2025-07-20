package design

import (
	"testing"

	coreV1 "github.com/meshery/schemas/models/v1alpha1/core"
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
			name:    "Import invalid source type",
			args:    args{"invalid source type", "file.yaml", "", false},
			want:    ErrInValidSource("invalid source type", []string{string(coreV1.HelmChart), string(coreV1.K8sManifest), string(coreV1.DockerCompose), string(coreV1.MesheryDesign)}),
			wantErr: true,
		},
		{
			name:    "Import missing file flag",
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
