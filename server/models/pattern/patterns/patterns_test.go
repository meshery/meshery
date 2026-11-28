package patterns

import (
	"context"
	"errors"
	"testing"

	"github.com/meshery/meshery/server/models/pattern/patterns/k8s"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta2/component"
)

func TestProcess_PreCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	comp := component.ComponentDefinition{
		DisplayName: "test-comp",
	}

	msgs, err := Process(
		ctx,
		[]string{},
		[]component.ComponentDefinition{comp},
		false,
		"test-pattern",
		nil,
		"test-user",
		nil,
		connection.Connection{Kind: "kubernetes"},
		false,
		false,
	)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Process with pre-cancelled context returned error %v, want %v", err, context.Canceled)
	}

	if msgs != nil {
		t.Fatalf("Process with pre-cancelled context returned msgs %v, want nil", msgs)
	}
}

func TestDeploy_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	comp := component.ComponentDefinition{
		DisplayName: "test-comp",
	}

	err := k8s.Deploy(ctx, nil, comp, false)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Deploy with cancelled context returned error %v, want %v", err, context.Canceled)
	}
}
