package appoptics

import (
	"testing"
)

func TestNewAOClient(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *AppOptics
		wantErr bool
	}{
		{
			name: "valid test",
			args: args{
				token: "ad6c84be90e7e16ef9150e0c0d809644956d5df6897b73d2340b3238fda40d9d",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAOClient(tt.args.token, "istio")
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAOClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
