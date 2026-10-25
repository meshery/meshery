package utils

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/meshery/meshkit/logger"
)

// SafeGoOption defines optional configuration for SafeGo.
type SafeGoOption func(*safeGoConfig)

type safeGoConfig struct {
	ctx     context.Context
	onPanic func(r any, stack []byte)
}

// WithContext attaches a context to the SafeGo routine config.
// If the context is already cancelled when the goroutine starts, fn is skipped.
func WithContext(ctx context.Context) SafeGoOption {
	return func(c *safeGoConfig) {
		c.ctx = ctx
	}
}

// WithOnPanic registers a callback invoked when a panic is recovered.
func WithOnPanic(cb func(r any, stack []byte)) SafeGoOption {
	return func(c *safeGoConfig) {
		c.onPanic = cb
	}
}

// SafeGo runs fn in a background goroutine behind a panic-recovery boundary.
// If fn panics, the panic is recovered, formatted with its full stack trace,
// logged via log, and optionally passed to an onPanic callback.
// This prevents unhandled panics in background goroutines from crashing the meshery-server process (#20445).
func SafeGo(log logger.Handler, name string, fn func(), opts ...SafeGoOption) {
	var cfg safeGoConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				if log != nil {
					log.Error(fmt.Errorf("panic recovered in background goroutine [%s]: %v\n%s", name, r, stack))
				}
				if cfg.onPanic != nil {
					func() {
						defer func() {
							if cbErr := recover(); cbErr != nil {
								if log != nil {
									log.Error(fmt.Errorf("panic recovered in onPanic callback for [%s]: %v\n%s", name, cbErr, debug.Stack()))
								}
							}
						}()
						cfg.onPanic(r, stack)
					}()
				}
			}
		}()

		if cfg.ctx != nil && cfg.ctx.Err() != nil {
			return
		}

		fn()
	}()
}

// RunSafe executes fn synchronously behind a panic-recovery boundary.
// It returns an error if fn panics, capturing the panic and stack trace in the log.
func RunSafe(log logger.Handler, name string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			panicErr := fmt.Errorf("panic recovered in [%s]: %v", name, r)
			if log != nil {
				log.Error(fmt.Errorf("%s\n%s", panicErr.Error(), stack))
			}
			err = panicErr
		}
	}()

	return fn()
}
