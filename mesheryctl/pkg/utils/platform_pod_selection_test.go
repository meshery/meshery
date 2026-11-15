package utils

import (
	"testing"

	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetRequiredPodsSelection(t *testing.T) {
	names := []string{"meshery-broker-0", "meshery-d6f7c8b9-x2k4m", "meshery-operator-5f6b-qwert", "meshery-operator-5f6b-asdfg"}
	var pods []v1core.Pod
	for _, name := range names {
		pods = append(pods, v1core.Pod{ObjectMeta: metav1.ObjectMeta{Name: name}})
	}
	for _, tc := range []struct {
		input, want string
		fail        bool
	}{
		{"meshery", names[1], false}, {names[0], names[0], false},
		{"broker", names[0], false}, {"meshery-operator", "", true},
		{"mesh", "", true}, {"missing", "", true},
	} {
		t.Run(tc.input, func(t *testing.T) {
			got, err := GetRequiredPods([]string{tc.input}, pods)
			if (err != nil) != tc.fail {
				t.Fatalf("got error %v, want failure %v", err, tc.fail)
			}
			if !tc.fail && (len(got) != 1 || got[tc.want] != tc.want) {
				t.Fatalf("got %v, want %s", got, tc.want)
			}
		})
	}

	// Exact matches win over clean-name matches: a bare pod named
	// "meshery-operator" resolves directly even though the operator
	// replicas also clean-match that name.
	t.Run("exact beats clean", func(t *testing.T) {
		pods := []v1core.Pod{
			{ObjectMeta: metav1.ObjectMeta{Name: "meshery-operator"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "meshery-operator-5f6b-qwert"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "meshery-operator-5f6b-asdfg"}},
		}
		got, err := GetRequiredPods([]string{"meshery-operator"}, pods)
		if err != nil {
			t.Fatalf("got error %v, want the exact pod", err)
		}
		if len(got) != 1 || got["meshery-operator"] != "meshery-operator" {
			t.Fatalf("got %v, want meshery-operator", got)
		}
	})
}
