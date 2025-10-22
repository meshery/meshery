package models

import (
	"net/http"
	"sync"
	"time"
)

// HTTPClientPool manages a pool of HTTP clients for better performance
type HTTPClientPool struct {
	clients chan *http.Client
	maxSize int
}

// NewHTTPClientPool creates a new HTTP client pool
func NewHTTPClientPool(maxSize int) *HTTPClientPool {
	pool := &HTTPClientPool{
		clients: make(chan *http.Client, maxSize),
		maxSize: maxSize,
	}

	// Pre-populate the pool with optimized HTTP clients
	for i := 0; i < maxSize; i++ {
		client := &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		}
		pool.clients <- client
	}

	return pool
}

// Get retrieves an HTTP client from the pool
func (p *HTTPClientPool) Get() *http.Client {
	select {
	case client := <-p.clients:
		return client
	default:
		// If pool is empty, create a new client
		return &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		}
	}
}

// Put returns an HTTP client to the pool
func (p *HTTPClientPool) Put(client *http.Client) {
	select {
	case p.clients <- client:
		// Successfully returned to pool
	default:
		// Pool is full, let the client be garbage collected
	}
}

// Global HTTP client pool instance
var (
	globalHTTPClientPool *HTTPClientPool
	poolOnce             sync.Once
)

// GetHTTPClientPool returns the global HTTP client pool instance
func GetHTTPClientPool() *HTTPClientPool {
	poolOnce.Do(func() {
		globalHTTPClientPool = NewHTTPClientPool(10)
	})
	return globalHTTPClientPool
}
