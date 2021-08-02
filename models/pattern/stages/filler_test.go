package stages

import (
	"testing"

	"github.com/layer5io/meshery/models/pattern/core"
)

func TestFiller(t *testing.T) {
	type args struct {
		data *Data
		err  error
		next ChainStageNextFunction
	}

	var samplePattern = `
name: InvalidPattern
services:
  entity1:
    type: Entity
    namespace: istio-system
    settings:
      some: "$(#ref.services.entity2.namespace)"
  entity2:
    type: XYZ
    namespace: test
  entity3:
    type: ABCD
    namespace: "$(#ref.services.entity1.namespace)"
    settings:
`
	p, err := core.NewPatternFile([]byte(samplePattern))
	if err != nil {
		t.Fatal("failed to generate pattern file: ", err)
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				data: &Data{
					Pattern: &p,
				},
				err: nil,
				next: func(data *Data, err error) {
					if data.Pattern.Services["entity1"].Settings["some"] != "test" {
						t.Errorf("expected: %s\nGot: %s", "test", data.Pattern.Services["entity1"].Settings["some"])
					}

					if data.Pattern.Services["entity3"].Namespace != "istio-system" {
						t.Errorf(
							"expected: %s\nGot: %s", "istio-system",
							data.Pattern.Services["entity3"].Namespace,
						)
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Filler(tt.args.data, tt.args.err, tt.args.next)
		})
	}
}

func Test_matchPattern(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "When pattern matches",
			args: args{
				str: "$(#ref.services.test)",
			},
			want:  "services.test",
			want1: true,
		},
		{
			name: "When pattern doesn't matches",
			args: args{
				str: "$(#ref.services.test",
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchPattern(tt.args.str)
			if got != tt.want {
				t.Errorf("matchPattern() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("matchPattern() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
