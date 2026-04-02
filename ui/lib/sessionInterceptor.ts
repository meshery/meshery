/**
 * Global fetch interceptor that catches 401/redirect responses from ALL
 * HTTP requests — dataFetch, RTK Query, Relay, raw fetch(), everything.
 *
 * Call installSessionInterceptor(store) once during app bootstrap
 * (in store/index.ts after configureStore).
 */

const INTERCEPTED = Symbol.for('meshery_session_intercepted');

let currentStore: { dispatch: (action: unknown) => void } | null = null;

export function installSessionInterceptor(store: { dispatch: (action: unknown) => void }) {
  if (typeof window === 'undefined') return;

  currentStore = store;

  // Don't wrap twice
  if ((window.fetch as unknown as Record<symbol, boolean>)[INTERCEPTED]) return;

  const originalFetch = window.fetch.bind(window);

  const interceptedFetch = async (...args: Parameters<typeof fetch>): Promise<Response> => {
    const response = await originalFetch(...args);

    if (response.status === 401 || response.redirected) {
      currentStore?.dispatch({ type: 'SESSION_EXPIRED' });
    }

    return response;
  };

  (interceptedFetch as unknown as Record<symbol, boolean>)[INTERCEPTED] = true;
  window.fetch = interceptedFetch;
}
