package handlers

import (
	"context"
	"errors"
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
	openConns    atomic.Int32
	onProvision  func(context.Context)
	provisionErr error
}

func (a *provisionAdapter) Provision(ctx context.Context, _ *meshes.ProvisionRequest) (*meshes.ProvisionResponse, error) {
	if a.onProvision != nil {
		a.onProvision(ctx)
	}
	if a.provisionErr != nil {
		return nil, a.provisionErr
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
	second.onProvision = func(ctx context.Context) {
		firstStillOpen.Store(first.waitForNoOpenConns())
	}

	comp := mutateTestComponent("workload", "Deployment")
	comp.Model = &modelv1beta1.ModelDefinition{Name: "kubernetes"}

	sap := &serviceActionProvider{
		ctx:         context.Background(),
		log:         log,
		userID:      "user",
		patternName: "design",
	}
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
func TestProvisionPropagatesRequestContextCancellation(t *testing.T) {
	log, err := logger.New("provision-cancel-test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build a logger: %v", err)
	}

	adapter, host := startProvisionAdapter(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	adapterStarted := make(chan struct{})
	adapterReceivedCancellation := make(chan struct{})

	adapter.onProvision = func(rpcCtx context.Context) {
		close(adapterStarted)
		<-rpcCtx.Done()
		close(adapterReceivedCancellation)
	}

	comp := mutateTestComponent("workload", "Deployment")
	comp.Model = &modelv1beta1.ModelDefinition{Name: "kubernetes"}

	sap := &serviceActionProvider{
		log:         log,
		userID:      "user",
		patternName: "design",
		ctx:         ctx,
	}

	done := make(chan error, 1)
	go func() {
		_, err := sap.Provision(stages.CompConfigPair{
			Component: *comp,
			Hosts:     []connection.Connection{host},
		})
		done <- err
	}()

	select {
	case <-adapterStarted:
		cancel()
	case <-time.After(3 * time.Second):
		t.Fatal("adapter RPC did not start")
	}

	select {
	case <-adapterReceivedCancellation:
		// The adapter received cancellation.
	case <-time.After(3 * time.Second):
		t.Fatal("adapter RPC did not receive context cancellation")
	}

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Provision returned nil error after context cancellation")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Provision did not return after cancellation")
	}
}

func TestProvisionDoesNotStartNextAdapterAfterCancellation(t *testing.T) {
	log, err := logger.New("provision-next-host-cancel-test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build a logger: %v", err)
	}

	first, firstHost := startProvisionAdapter(t)
	second, secondHost := startProvisionAdapter(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firstStarted := make(chan struct{})
	var secondCalls atomic.Int32

	first.onProvision = func(rpcCtx context.Context) {
		close(firstStarted)
		<-rpcCtx.Done()
	}
	second.onProvision = func(context.Context) {
		secondCalls.Add(1)
	}

	comp := mutateTestComponent("workload", "Deployment")
	comp.Model = &modelv1beta1.ModelDefinition{Name: "kubernetes"}

	sap := &serviceActionProvider{
		log:         log,
		userID:      "user",
		patternName: "design",
		ctx:         ctx,
	}

	done := make(chan error, 1)
	go func() {
		_, err := sap.Provision(stages.CompConfigPair{
			Component: *comp,
			Hosts:     []connection.Connection{firstHost, secondHost},
		})
		done <- err
	}()

	select {
	case <-firstStarted:
		cancel()
	case <-time.After(3 * time.Second):
		t.Fatal("first adapter RPC did not start")
	}

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Provision returned nil error after cancellation")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Provision did not return after cancellation")
	}

	if calls := secondCalls.Load(); calls != 0 {
		t.Errorf("second adapter received %d Provision call(s) after cancellation, want 0", calls)
	}

	for name, adapter := range map[string]*provisionAdapter{
		"first":  first,
		"second": second,
	} {
		if open := adapter.waitForNoOpenConns(); open != 0 {
			t.Errorf("the %s adapter still holds %d client connection(s) after Provision returned, want 0", name, open)
		}
	}
}
func TestProvisionContinuesAfterAdapterRPCError(t *testing.T) {
	log, err := logger.New("provision-rpc-error-test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build a logger: %v", err)
	}

	first, firstHost := startProvisionAdapter(t)
	second, secondHost := startProvisionAdapter(t)

	first.provisionErr = errors.New("first adapter failed")

	var secondCalls atomic.Int32
	second.onProvision = func(context.Context) {
		secondCalls.Add(1)
	}

	comp := mutateTestComponent("workload", "Deployment")
	comp.Model = &modelv1beta1.ModelDefinition{Name: "kubernetes"}

	sap := &serviceActionProvider{
		ctx:         context.Background(),
		log:         log,
		userID:      "user",
		patternName: "design",
	}

	msgs, err := sap.Provision(stages.CompConfigPair{
		Component: *comp,
		Hosts:     []connection.Connection{firstHost, secondHost},
	})
	if err != nil {
		t.Fatalf("Provision returned unexpected error: %v", err)
	}

	if calls := secondCalls.Load(); calls != 1 {
		t.Errorf("second adapter received %d Provision call(s), want 1", calls)
	}

	if len(msgs) != 2 {
		t.Fatalf("Provision returned %d messages, want 2", len(msgs))
	}

	if msgs[0].Summary[0].Success {
		t.Error("first adapter summary marked successful after RPC failure")
	}

	if !msgs[1].Summary[0].Success {
		t.Error("second adapter summary marked failed, want success")
	}
}
