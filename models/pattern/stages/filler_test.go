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
      some2:
        - some2Some1:
            data:
              - name: 123
                value: "1234"
  entity2:
    type: XYZ
    namespace: test
    settings:
      some: "$(#ref.services.entity1.settings.some2.0.some2Some1.data.0.value)"
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
					if err != nil {
						t.Fatalf("got error: %s", err)
					}

					if data.Pattern.Services["entity1"].Settings["some"] != "test" {
						t.Errorf("expected: %s\nGot: %s", "test", data.Pattern.Services["entity1"].Settings["some"])
					}

					if data.Pattern.Services["entity3"].Namespace != "istio-system" {
						t.Errorf(
							"expected: %s\nGot: %s", "istio-system",
							data.Pattern.Services["entity3"].Namespace,
						)
					}

					if data.Pattern.Services["entity2"].Settings["some"] != "1234" {
						t.Errorf(
							"expected: %s\nGot: %s", "1234",
							data.Pattern.Services["entity2"].Settings["some"],
						)
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Filler(true)(tt.args.data, tt.args.err, tt.args.next)
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
