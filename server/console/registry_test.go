package console

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/errors"
)

// stubDriver is a Driver that serves no consoles; the registry never inspects
// what a driver can do, only that it exists.
type stubDriver struct{}

func (stubDriver) Capabilities(context.Context, TargetRef) (Capabilities, error) {
	return Capabilities{}, nil
}

type stubFactory struct {
	kind string
	err  error
}

func (f stubFactory) Kind() string { return f.kind }

func (f stubFactory) NewDriver(context.Context, ConnectionContext) (Driver, error) {
	if f.err != nil {
		return nil, f.err
	}
	return stubDriver{}, nil
}

func connectionOfKind(t *testing.T, kind string) *connections.Connection {
	t.Helper()
	id, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("uuid.NewV4: %v", err)
	}
	return &connections.Connection{ID: id, Kind: kind}
}

func TestRegistryResolvesDriverByConnectionKind(t *testing.T) {
	registry := NewRegistry()
	registry.Register(stubFactory{kind: "kubernetes"})

	driver, err := registry.Driver(context.Background(), ConnectionContext{
		Connection: connectionOfKind(t, "kubernetes"),
	})
	if err != nil {
		t.Fatalf("Driver: %v", err)
	}
	if driver == nil {
		t.Fatal("Driver returned a nil driver with a nil error")
	}
}

// TestRegistryIsOpenToNewKinds is the scalability claim in test form: a new
// resource universe is one Register call, with nothing else changing.
func TestRegistryIsOpenToNewKinds(t *testing.T) {
	registry := NewRegistry()
	registry.Register(stubFactory{kind: "kubernetes"})
	registry.Register(stubFactory{kind: "docker"})

	for _, kind := range []string{"kubernetes", "docker"} {
		if _, err := registry.Driver(context.Background(), ConnectionContext{
			Connection: connectionOfKind(t, kind),
		}); err != nil {
			t.Errorf("Driver for kind %q: %v", kind, err)
		}
	}

	kinds := registry.Kinds()
	if len(kinds) != 2 || kinds[0] != "docker" || kinds[1] != "kubernetes" {
		t.Errorf("Kinds() = %v, want a sorted [docker kubernetes]", kinds)
	}
}

func TestRegistryUnknownKind(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Driver(context.Background(), ConnectionContext{
		Connection: connectionOfKind(t, "nomad"),
	})
	if err == nil {
		t.Fatal("Driver for an unregistered kind succeeded")
	}
	if code := errors.GetCode(err); code != ErrNoDriverCode {
		t.Errorf("error code = %q, want %q", code, ErrNoDriverCode)
	}
}

func TestRegistryNilConnection(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Driver(context.Background(), ConnectionContext{})
	if err == nil {
		t.Fatal("Driver with a nil connection succeeded")
	}
	if code := errors.GetCode(err); code != ErrInvalidTargetCode {
		t.Errorf("error code = %q, want %q", code, ErrInvalidTargetCode)
	}
}

// TestRegistryWrapsFactoryFailure keeps a credential-resolution failure from
// masquerading as "this kind is unsupported".
func TestRegistryWrapsFactoryFailure(t *testing.T) {
	registry := NewRegistry()
	registry.Register(stubFactory{kind: "kubernetes", err: ErrInvalidTarget("no credential")})

	_, err := registry.Driver(context.Background(), ConnectionContext{
		Connection: connectionOfKind(t, "kubernetes"),
	})
	if err == nil {
		t.Fatal("Driver succeeded despite a failing factory")
	}
	if code := errors.GetCode(err); code != ErrDriverInitCode {
		t.Errorf("error code = %q, want %q", code, ErrDriverInitCode)
	}
}

func TestRegistryRegisterReplaces(t *testing.T) {
	registry := NewRegistry()
	registry.Register(stubFactory{kind: "kubernetes", err: ErrInvalidTarget("first")})
	registry.Register(stubFactory{kind: "kubernetes"})

	if _, err := registry.Driver(context.Background(), ConnectionContext{
		Connection: connectionOfKind(t, "kubernetes"),
	}); err != nil {
		t.Errorf("re-registering a kind did not replace the previous factory: %v", err)
	}
	if kinds := registry.Kinds(); len(kinds) != 1 {
		t.Errorf("Kinds() = %v, want a single entry", kinds)
	}
}

func TestKindValid(t *testing.T) {
	for _, kind := range []Kind{KindTerminal, KindLogs} {
		if !ValidKind(kind) {
			t.Errorf("ValidKind(%q) = false, want true", kind)
		}
	}
	for _, kind := range []Kind{"", "attach", "portforward"} {
		if ValidKind(kind) {
			t.Errorf("ValidKind(%q) = true, want false", kind)
		}
	}
}

func TestCapabilitiesSupports(t *testing.T) {
	tests := []struct {
		name string
		caps Capabilities
		kind Kind
		want bool
	}{
		{"terminal supported", Capabilities{Terminal: true}, KindTerminal, true},
		{"terminal unsupported", Capabilities{Logs: true}, KindTerminal, false},
		{"logs supported", Capabilities{Logs: true}, KindLogs, true},
		{"logs unsupported", Capabilities{Terminal: true}, KindLogs, false},
		{"unknown kind", Capabilities{Terminal: true, Logs: true}, Kind("attach"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Supports(tt.caps, tt.kind); got != tt.want {
				t.Errorf("Supports(%q) = %v, want %v", tt.kind, got, tt.want)
			}
		})
	}
}
