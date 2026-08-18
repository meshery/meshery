package handlers

import (
	"testing"

	"github.com/meshery/meshery/server/models/pattern/patterns"
	"github.com/stretchr/testify/assert"
)

func TestDeploymentFailedComponents(t *testing.T) {
	ctxWith := func(comps ...patterns.DeploymentMessagePerComp) []patterns.DeploymentMessagePerContext {
		return []patterns.DeploymentMessagePerContext{{Summary: comps}}
	}

	t.Run("empty response", func(t *testing.T) {
		assert.Empty(t, deploymentFailedComponents(map[string]interface{}{}))
	})

	t.Run("all components succeed", func(t *testing.T) {
		resp := map[string]interface{}{
			"nginx": ctxWith(
				patterns.DeploymentMessagePerComp{CompName: "nginx", Success: true},
				patterns.DeploymentMessagePerComp{CompName: "gateway", Success: true},
			),
		}
		assert.Empty(t, deploymentFailedComponents(resp))
	})

	t.Run("single failed component", func(t *testing.T) {
		resp := map[string]interface{}{
			"nginx": ctxWith(
				patterns.DeploymentMessagePerComp{CompName: "nginx", Success: false},
				patterns.DeploymentMessagePerComp{CompName: "gateway", Success: true},
			),
		}
		assert.Equal(t, []string{"nginx"}, deploymentFailedComponents(resp))
	})

	t.Run("fall back to kind when comp name is empty", func(t *testing.T) {
		resp := map[string]interface{}{
			"nginx": ctxWith(
				patterns.DeploymentMessagePerComp{Kind: "Deployment", Success: false},
			),
		}
		assert.Equal(t, []string{"Deployment"}, deploymentFailedComponents(resp))
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
		assert.Equal(t, []string{"gateway", "nginx"}, deploymentFailedComponents(resp))
	})

	t.Run("non-summary values are skipped", func(t *testing.T) {
		resp := map[string]interface{}{
			"dryRun": map[string]interface{}{"success": true},
			"nginx":  "unexpected",
		}
		assert.Empty(t, deploymentFailedComponents(resp))
	})

	t.Run("empty names are skipped", func(t *testing.T) {
		resp := map[string]interface{}{
			"nginx": ctxWith(
				patterns.DeploymentMessagePerComp{Success: false},
			),
		}
		assert.Empty(t, deploymentFailedComponents(resp))
	})
}
