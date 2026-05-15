// sseClient.ts — thin wrapper around the browser's EventSource API that
// normalises error handling and reconnection for Meshery's SSE streams.
//
// The event envelope shape (name + data) mirrors server/internal/sse/stream.go.
//
// TODO(schemas-canonical): Once meshery/schemas publishes an AsyncAPI document
// for the SSE channels (see docs/asyncapi/controller-status-sse.yaml) and
// generates TypeScript client stubs, replace the inline EventEnvelope type
// here with the generated import.
// Tracked in: meshery/meshery#19424

/** Wire envelope emitted by server/internal/sse/stream.go */
export interface SseEvent<T = unknown> {
  name: string;
  data: T;
}

export type SseEventHandler<T = unknown> = (_event: SseEvent<T>) => void;
export type SseErrorHandler = (_err: Event) => void;

export interface SseClientOptions {
  /** Absolute or relative URL of the SSE endpoint. */
  url: string;
  /**
   * Map from SSE event names to their handlers.
   * Use "*" as key to receive all events regardless of name.
   */
  handlers: Record<string, SseEventHandler>;
  /** Called on connection error.  The client does NOT auto-reconnect. */
  onError?: SseErrorHandler;
}

/**
 * SseClient opens an EventSource connection and dispatches events to the
 * registered handlers.  Call close() to tear down the connection.
 */
export class SseClient {
  private es: EventSource | null = null;
  private readonly opts: SseClientOptions;

  constructor(opts: SseClientOptions) {
    this.opts = opts;
  }

  open(): void {
    if (this.es) {
      return;
    }

    this.es = new EventSource(this.opts.url, { withCredentials: true });

    // Register a listener for each named event.
    for (const [eventName, handler] of Object.entries(this.opts.handlers)) {
      if (eventName === '*') continue;
      this.es.addEventListener(eventName, (raw: MessageEvent) => {
        try {
          const data = JSON.parse(raw.data);
          handler({ name: eventName, data });
        } catch {
          handler({ name: eventName, data: raw.data });
        }
      });
    }

    // Wildcard handler receives every named event.
    if (this.opts.handlers['*']) {
      const wildcardHandler = this.opts.handlers['*'];
      this.es.onmessage = (raw: MessageEvent) => {
        try {
          const data = JSON.parse(raw.data);
          wildcardHandler({ name: raw.type, data });
        } catch {
          wildcardHandler({ name: raw.type, data: raw.data });
        }
      };
    }

    if (this.opts.onError) {
      this.es.onerror = this.opts.onError;
    }
  }

  close(): void {
    this.es?.close();
    this.es = null;
  }

  get readyState(): number {
    return this.es?.readyState ?? EventSource.CLOSED;
  }
}
