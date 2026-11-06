package utils

import (
	"errors"
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEphemeralPort(t *testing.T) {
	port, err := GetEphemeralPort()
	assert.NoError(t, err)
	assert.Greater(t, port, 0)
	assert.LessOrEqual(t, port, 65535)
}

func TestPortForwardLocalPort(t *testing.T) {
	pf := &PortForward{
		host:      "localhost",
		localPort: 9081,
	}
	assert.Equal(t, 9081, pf.LocalPort())
	assert.Equal(t, "http://localhost:9081", pf.URLFor(""))
}

func TestIsBindError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil error", err: nil, want: false},
		{name: "generic error", err: assert.AnError, want: false},
		{name: "direct syscall.EADDRINUSE", err: syscall.EADDRINUSE, want: true},
		{
			name: "net.OpError wrapping syscall.EADDRINUSE",
			err: &net.OpError{
				Op:  "listen",
				Net: "tcp",
				Err: &os.SyscallError{
					Syscall: "bind",
					Err:     syscall.EADDRINUSE,
				},
			},
			want: true,
		},
		{
			name: "unix address in use string",
			err:  errors.New("listen tcp 127.0.0.1:9081: bind: address already in use"),
			want: true,
		},
		{
			name: "windows address in use string",
			err:  errors.New("bind: Only one usage of each socket address (protocol/network address/port) is normally permitted."),
			want: true,
		},
		{
			name: "client-go listener error detail with address already in use",
			err:  errors.New("unable to listen on any of the requested ports: [{9081 8080}]: Unable to listen on port 9081: Listeners failed to create with the following errors: [unable to create listener: Error listen tcp4 127.0.0.1:9081: bind: address already in use] for meshery/meshery-abc"),
			want: true,
		},
		{
			name: "generic client-go listener failure without detail",
			err:  errors.New("unable to listen on any of the requested ports: [{9081 8080}]"),
			want: false,
		},
		{
			name: "permission denied string",
			err:  errors.New("bind: permission denied"),
			want: false,
		},
		{
			name: "listen tcp permission denied",
			err:  errors.New("listen tcp 0.0.0.0:80: bind: permission denied"),
			want: false,
		},
		{
			name: "client-go permission denied listener error",
			err:  errors.New("unable to create listener: Error listen tcp4 127.0.0.1:80: bind: permission denied"),
			want: false,
		},
		{
			name: "operation not permitted",
			err:  errors.New("operation not permitted"),
			want: false,
		},
		{
			name: "net.OpError wrapping syscall.EACCES",
			err: &net.OpError{
				Op:  "listen",
				Net: "tcp",
				Err: &os.SyscallError{
					Syscall: "bind",
					Err:     syscall.EACCES,
				},
			},
			want: false,
		},
		{
			name: "spdy transport error",
			err:  errors.New("error upgrading connection: SPDY error"),
			want: false,
		},
		{
			name: "connection refused",
			err:  errors.New("dial tcp 10.96.0.1:443: connect: connection refused"),
			want: false,
		},
		{
			name: "pod not found",
			err:  errors.New("pods \"meshery-xyz\" not found"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsBindError(tt.err), "IsBindError failed for case: %s", tt.name)
			assert.Equal(t, tt.want, IsPortBusyError(tt.err), "IsPortBusyError failed for case: %s", tt.name)
		})
	}
}


