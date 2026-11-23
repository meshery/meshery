/**
 * Client half of the console wire protocol.
 *
 * Every type below is re-exported from `@meshery/schemas` — the `v1beta1/console`
 * construct is the single source of truth for this contract, and the Go server
 * aliases the same generated models. Nothing here re-declares a wire shape; this
 * module contributes only the URL builders, which are client concerns.
 *
 * A console socket carries two interleaved channels, distinguished by WebSocket
 * frame type rather than by an envelope field. Binary frames carry raw payload
 * bytes: stdin keystrokes to the server (terminal consoles only), and stdout or
 * log bytes to the client. Text frames carry a JSON `ConsoleControlMessage`.
 *
 * Terminal output is not valid UTF-8 in general — escape sequences, and
 * arbitrary octets from a program writing raw bytes — so it travels as binary
 * and is handed to xterm as a Uint8Array rather than being decoded here.
 *
 * Every console ends with exactly one terminal control frame (`error`, `exit`,
 * or `eof`) before the socket closes, so a client never has to infer why a
 * stream stopped from the close code alone.
 */

import type { components } from '@meshery/schemas/constructs/v1beta1/console/Console';
import type { OpenLogConsoleApiArg, OpenTerminalConsoleApiArg } from '@meshery/schemas/mesheryApi';

type ConsoleSchemas = components['schemas'];

/** Console kinds a resource may support. */
export type ConsoleKind = ConsoleSchemas['ConsoleKind'];

/** Control frame types. `resize` is client-to-server; the rest are server-to-client. */
export type ControlType = ConsoleSchemas['ConsoleControlType'];

/** What a driver can do with one specific target, as resolved against live state. */
export type ConsoleCapabilities = ConsoleSchemas['ConsoleCapabilities'];

/** Addresses the resource a console attaches to, within one connection. */
export type ConsoleTarget = ConsoleSchemas['ConsoleTarget'];

/** The JSON payload of a console socket's text frames. */
export type ControlMessage = ConsoleSchemas['ConsoleControlMessage'];

/**
 * Parameters of a log console: the generated query-argument type, minus the
 * fields that address the target. They are fixed when the remote log stream is
 * created, so changing one means opening a new console.
 */
export type LogConsoleParams = Omit<
  OpenLogConsoleApiArg,
  'connectionId' | 'resource' | 'name' | 'namespace' | 'container'
>;

/** Parameters of a terminal console, on the same basis. */
export type TerminalConsoleParams = Omit<
  OpenTerminalConsoleApiArg,
  'connectionId' | 'resource' | 'name' | 'namespace' | 'container'
>;

/**
 * Close codes the server sends, from the private-use range RFC 6455 §7.4.2
 * reserves for applications.
 */
export const CLOSE_CONSOLE_ENDED = 4000;
export const CLOSE_CONSOLE_FAILED = 4001;

const CONNECTIONS_BASE = '/api/integrations/connections';

/**
 * Builds the absolute WebSocket URL for a console.
 *
 * The socket is same-origin, so the browser sends the Meshery session cookie on
 * the handshake and the server's usual auth chain applies. The server also
 * enforces a same-origin check on that handshake, which is why this must be
 * built from `window.location` rather than from a configured base URL.
 */
export const consoleSocketUrl = (
  connectionId: string,
  kind: ConsoleKind,
  target: ConsoleTarget,
  params: LogConsoleParams & TerminalConsoleParams = {},
): string => {
  const scheme = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const url = new URL(
    `${scheme}//${window.location.host}${CONNECTIONS_BASE}/${connectionId}/console/${kind}`,
  );

  url.searchParams.set('resource', target.resource);
  url.searchParams.set('name', target.name);
  if (target.namespace) url.searchParams.set('namespace', target.namespace);
  if (target.container) url.searchParams.set('container', target.container);

  const { command, ...rest } = params;
  // `command` repeats, one query parameter per argv element, because that is
  // how Go's url.Values decodes a []string.
  command?.forEach((arg) => url.searchParams.append('command', arg));

  Object.entries(rest).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      url.searchParams.set(key, String(value));
    }
  });

  return url.toString();
};
