package design

import (
	"reflect"
	"testing"
)

func TestParseDeploymentSummary(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantFailed []string
		wantOK     bool
	}{
		{
			name: "failed components from map iteration are returned sorted",
			body: `{
  "zeta": [{"Summary": [{"Kind": "Deployment", "CompName": "zeta", "Success": false}]}],
  "alpha": [{"Summary": [{"Kind": "Deployment", "CompName": "alpha", "Success": false}]}],
  "web": [{"Summary": [{"Kind": "Deployment", "CompName": "web", "Success": true}]}]
}`,
			wantFailed: []string{"alpha", "zeta"},
			wantOK:     true,
		},
		{
			name:       "all components succeed",
			body:       `{"web": [{"Summary": [{"Kind": "Deployment", "CompName": "web", "Success": true}]}]}`,
			wantFailed: nil,
			wantOK:     true,
		},
		{
			name: "duplicate failed components are deduped",
			body: `{
  "web": [
    {"Summary": [{"Kind": "Deployment", "CompName": "web", "Success": false}]},
    {"Summary": [{"Kind": "Deployment", "CompName": "web", "Success": false}]}
  ]
}`,
			wantFailed: []string{"web"},
			wantOK:     true,
		},
		{
			name:       "non-summary body is reported as not a summary",
			body:       `deployed application myapp`,
			wantFailed: nil,
			wantOK:     false,
		},
		{
			name:       "empty summary has no failures",
			body:       `{}`,
			wantFailed: nil,
			wantOK:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFailed, gotOK := parseDeploymentSummary([]byte(tt.body))
			if gotOK != tt.wantOK {
				t.Errorf("parseDeploymentSummary() ok = %v, want %v", gotOK, tt.wantOK)
			}
			if !reflect.DeepEqual(gotFailed, tt.wantFailed) {
				t.Errorf("parseDeploymentSummary() failed = %v, want %v", gotFailed, tt.wantFailed)
			}
		})
	}
}
