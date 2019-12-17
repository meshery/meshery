package helpers

import (
	"io/ioutil"
	"testing"
)

func TestScanIstio(t *testing.T) {
	type args struct {
		kubeconfig  []byte
		contextName string
	}
	config, _ := ioutil.ReadFile("../10.199.75.50.config")
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				kubeconfig:  config,
				contextName: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ScanIstio(tt.args.kubeconfig, tt.args.contextName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanIstio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("ScanIstio() = %v, want %v", got, tt.want)
			// }
		})
	}
}
