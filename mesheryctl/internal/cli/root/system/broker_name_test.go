package system

import (
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// TestIsBrokerPodName pins that `mesheryctl system check --operator` recognises
// the broker under the name the operator actually deploys it as.
//
// Meshery Operator >= 1.0.0 renders the broker from the official NATS chart, so
// the pod is meshery-nats-0; matching only the pre-1.0.0 meshery-broker name
// reported a healthy, Running broker as "!! Meshery Broker is not running" on
// every current cluster. Both names are accepted because the operator version
// belongs to the cluster, not to the CLI.
func TestIsBrokerPodName(t *testing.T) {
	cases := []struct {
		podName string
		want    bool
	}{
		// Operator >= 1.0.0 (NATS chart).
		{podName: "meshery-nats-0", want: true},
		{podName: "meshery-nats-1", want: true},
		// Operator < 1.0.0.
		{podName: "meshery-broker-0", want: true},
		// Everything else in the namespace.
		{podName: "meshery-operator-84945fd54c-45cck", want: false},
		{podName: "meshery-meshsync-c4b96d647-kl5l6", want: false},
		{podName: "meshery-6c54448cc8-fhvwj", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.podName, func(t *testing.T) {
			// Go through GetCleanPodName, because that is what the health check
			// feeds the matcher - asserting on the raw pod name would pass while
			// the real call site failed.
			if got := isBrokerPodName(utils.GetCleanPodName(tc.podName)); got != tc.want {
				t.Errorf("isBrokerPodName(GetCleanPodName(%q)) = %v, want %v", tc.podName, got, tc.want)
			}
		})
	}
}
