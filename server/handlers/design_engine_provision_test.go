package handlers

import (
	"context"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/meshery/meshery/server/meshes"
	"github.com/meshery/meshery/server/models/pattern/stages"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta1/connection"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

// provisionAdapter is an in-process adapter that counts its open client connections.
type provisionAdapter struct {
	meshes.UnimplementedMeshServiceServer
	openConns   atomic.Int32
	onProvision func()
}

func (a *provisionAdapter) Provision(context.Context, *meshes.ProvisionRequest) (*meshes.ProvisionResponse, error) {
	if a.onProvision != nil {
		a.onProvision()
	}
	return &meshes.ProvisionResponse{Message: "provisioned"}, nil
}

func (a *provisionAdapter) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
	return ctx
}
func (a *provisionAdapter) HandleRPC(context.Context, stats.RPCStats) {}
func (a *provisionAdapter) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

func (a *provisionAdapter) HandleConn(_ context.Context, s stats.ConnStats) {
	switch s.(type) {
	case *stats.ConnBegin:
		a.openConns.Add(1)
	case *stats.ConnEnd:
		a.openConns.Add(-1)
	}
}

// waitForNoOpenConns waits briefly, since a client close reaches the server asynchronously.
func (a *provisionAdapter) waitForNoOpenConns() int32 {
	deadline := time.Now().Add(2 * time.Second)
	for a.openConns.Load() != 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	return a.openConns.Load()
}

// startProvisionAdapter serves a provisionAdapter on a free local port and
// returns it with the host Provision reaches it through.
func startProvisionAdapter(t *testing.T) (*provisionAdapter, connection.Connection) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen for the adapter: %v", err)
	}

	adapter := &provisionAdapter{}
	server := grpc.NewServer(grpc.StatsHandler(adapter))
	meshes.RegisterMeshServiceServer(server, adapter)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)

	host := connection.Connection{
		Kind: "adapter",
		Metadata: map[string]interface{}{
			"host_name": "127.0.0.1",
			"port":      listener.Addr().(*net.TCPAddr).Port,
		},
	}
	return adapter, host
}

// Each host's adapter client must close before the next host is provisioned (#20874).
func TestProvisionClosesEachAdapterClientBeforeTheNextHost(t *testing.T) {
	log, err := logger.New("provision-test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build a logger: %v", err)
	}

	first, firstHost := startProvisionAdapter(t)
	second, secondHost := startProvisionAdapter(t)

	var firstStillOpen atomic.Int32
	second.onProvision = func() {
		firstStillOpen.Store(first.waitForNoOpenConns())
	}

	comp := mutateTestComponent("workload", "Deployment")
	comp.Model = &modelv1beta1.ModelDefinition{Name: "kubernetes"}

	sap := &serviceActionProvider{log: log, userID: "user", patternName: "design"}
	msgs, err := sap.Provision(stages.CompConfigPair{
		Component: *comp,
		Hosts:     []connection.Connection{firstHost, secondHost},
	})
	if err != nil {
		t.Fatalf("Provision: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("Provision returned %d message(s), want 2", len(msgs))
	}

	if open := firstStillOpen.Load(); open != 0 {
		t.Errorf("the first adapter still held %d client connection(s) while the second host was provisioned, want 0", open)
	}
	for name, adapter := range map[string]*provisionAdapter{"first": first, "second": second} {
		if open := adapter.waitForNoOpenConns(); open != 0 {
			t.Errorf("the %s adapter still holds %d client connection(s) after Provision returned, want 0", name, open)
		}
	}
}
