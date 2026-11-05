package socket

import (
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

// ListenUnix wraps `winio.ListenUnix`.
// It provides API compatibility for named pipes with the Unix domain socket API.
func ListenUnix(path string) (net.Listener, error) {
	return winio.ListenPipe(path, &winio.PipeConfig{
		MessageMode:      true,  // Use message mode so that CloseWrite() is supported
		InputBufferSize:  65536, // Use 64KB buffers to improve performance
		OutputBufferSize: 65536,
	})
}

func DialSocket(socket string) (net.Conn, error) {
	timeout := 1 * time.Second
	return winio.DialPipe(socket, &timeout)
}
