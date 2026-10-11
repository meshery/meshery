package utils

import (
	"errors"
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
		localPort: 9081,
	}
	assert.Equal(t, 9081, pf.LocalPort())
}

func TestIsBindError(t *testing.T) {
	assert.False(t, IsBindError(nil))
	assert.False(t, IsBindError(assert.AnError))
	assert.True(t, IsBindError(errors.New("listen tcp 127.0.0.1:9081: bind: address already in use")))
	assert.True(t, IsBindError(errors.New("unable to listen on any of the requested ports: [{9081 8080}]")))
	assert.True(t, IsBindError(errors.New("bind: Only one usage of each socket address (protocol/network address/port) is normally permitted.")))
	assert.True(t, IsBindError(errors.New("bind: permission denied")))
	assert.False(t, IsBindError(errors.New("error upgrading connection: SPDY error")))
	assert.False(t, IsBindError(errors.New("dial tcp 10.96.0.1:443: connect: connection refused")))
	assert.False(t, IsBindError(errors.New("pods \"meshery-xyz\" not found")))
}
