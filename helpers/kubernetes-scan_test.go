package helpers

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_detectServiceForDeploymentImage(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	type args struct {
		kubeconfig  []byte
		contextName string
		imageNames  []string
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
				imageNames:  []string{"prometheus", "grafana"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := detectServiceForDeploymentImage(tt.args.kubeconfig, tt.args.contextName, tt.args.imageNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("detectServiceForDeploymentImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("detectServiceForDeploymentImage() = %v, want %v", got, tt.want)
			// }
		})
	}
}
