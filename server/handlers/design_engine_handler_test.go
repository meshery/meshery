package handlers

import (
	"reflect"
	"testing"

	"github.com/meshery/meshery/server/models/pattern/patterns"
)

func TestDeploymentFailedComponents(t *testing.T) {
	ctxWith := func(comps ...patterns.DeploymentMessagePerComp) []patterns.DeploymentMessagePerContext {
		return []patterns.DeploymentMessagePerContext{{Summary: comps}}
	}

	t.Run("empty response", func(t *testing.T) {
		if got := deploymentFailedComponents(map[string]interface{}{}); len(got) != 0 {
			t.Fatalf("got %v, want empty", got)
		}
	})

	t.Run("all components succeed", func(t *testing.T) {
		resp := map[string]interface{}{
			"nginx": ctxWith(
				patterns.DeploymentMessagePerComp{CompName: "nginx", Success: true},
				patterns.DeploymentMessagePerComp{CompName: "gateway", Success: true},
			),
		}
		if got := deploymentFailedComponents(resp); len(got) != 0 {
			t.Fatalf("got %v, want empty", got)
		}
	})

	t.Run("single failed component", func(t *testing.T) {
		resp := map[string]interface{}{
			"nginx": ctxWith(
				patterns.DeploymentMessagePerComp{CompName: "nginx", Success: false},
				patterns.DeploymentMessagePerComp{CompName: "gateway", Success: true},
			),
		}
		want := []string{"nginx"}
		if got := deploymentFailedComponents(resp); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("fall back to kind when comp name is empty", func(t *testing.T) {
		resp := map[string]interface{}{
			"nginx": ctxWith(
				patterns.DeploymentMessagePerComp{Kind: "Deployment", Success: false},
			),
		}
		want := []string{"Deployment"}
		if got := deploymentFailedComponents(resp); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("dedupe and sort across contexts", func(t *testing.T) {
		resp := map[string]interface{}{
			"nginx": ctxWith(
				patterns.DeploymentMessagePerComp{CompName: "gateway", Success: false},
				patterns.DeploymentMessagePerComp{CompName: "nginx", Success: false},
			),
			"gateway": ctxWith(
				patterns.DeploymentMessagePerComp{CompName: "gateway", Success: false},
				patterns.DeploymentMessagePerComp{CompName: "bookinfo", Success: true},
			),
		}
		want := []string{"gateway", "nginx"}
		if got := deploymentFailedComponents(resp); !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("non-summary values are skipped", func(t *testing.T) {
		resp := map[string]interface{}{
			"dryRun": map[string]interface{}{"success": true},
			"nginx":  "unexpected",
		}
		if got := deploymentFailedComponents(resp); len(got) != 0 {
			t.Fatalf("got %v, want empty", got)
		}
	})

	t.Run("empty names are skipped", func(t *testing.T) {
		resp := map[string]interface{}{
			"nginx": ctxWith(
				patterns.DeploymentMessagePerComp{Success: false},
			),
		}
		if got := deploymentFailedComponents(resp); len(got) != 0 {
			t.Fatalf("got %v, want empty", got)
		}
	})
}
