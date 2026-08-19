package ai

// Middleware wraps a provider with shared cross-cutting behavior.
type Middleware func(LLMProvider) LLMProvider

// ChainMiddleware applies middlewares in order so the first middleware passed
// becomes the outermost wrapper.
func ChainMiddleware(provider LLMProvider, middlewares ...Middleware) LLMProvider {
	if provider == nil {
		return nil
	}

	wrapped := provider
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		if middleware == nil {
			continue
		}

		wrapped = middleware(wrapped)
	}

	return wrapped
}
