package telemetry

import (
	"errors"
	"io"
)

const maxResponseSize = 16 << 20 // 16 MiB

// ErrResponseTooLarge identifies an upstream response that cannot be returned intact.
var ErrResponseTooLarge = errors.New("response exceeds 16 MiB limit")

// ReadBody reads one extra byte so an oversized response cannot look successful.
func ReadBody(body io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(body, maxResponseSize+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxResponseSize {
		return nil, ErrResponseTooLarge
	}
	return data, nil
}
