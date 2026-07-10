/**
 * Client half of the session wire protocol.
 *
 * Every type below is re-exported from `@meshery/schemas` — the `v1beta1/session`
 * construct is the single source of truth for this contract, and the Go server
 * aliases the same generated models. Nothing here re-declares a wire shape; this
 * module contributes only the URL builders, which are client concerns.
 *
 * A session socket carries two interleaved channels, distinguished by WebSocket
 * frame type rather than by an envelope field. Binary frames carry raw payload
 * bytes: stdin keystrokes to the server (terminal sessions only), and stdout or
 * log bytes to the client. Text frames carry a JSON `SessionControlMessage`.
 *
 * Terminal output is not valid UTF-8 in general — escape sequences, and
 * arbitrary octets from a program writing raw bytes — so it travels as binary
 * and is handed to xterm as a Uint8Array rather than being decoded here.
 *
 * Every session ends with exactly one terminal control frame (`error`, `exit`,
 * or `eof`) before the socket closes, so a client never has to infer why a
 * stream stopped from the close code alone.
 */

import type { components } from '@meshery/schemas/constructs/v1beta1/session/Session';
import type { OpenLogSessionApiArg, OpenTerminalSessionApiArg } from '@meshery/schemas/mesheryApi';

type SessionSchemas = components['schemas'];

/** Session kinds a resource may support. */
export type SessionKind = SessionSchemas['SessionKind'];

/** Control frame types. `resize` is client-to-server; the rest are server-to-client. */
export type ControlType = SessionSchemas['SessionControlType'];

/** What a driver can do with one specific target, as resolved against live state. */
export type SessionCapabilities = SessionSchemas['SessionCapabilities'];

/** Addresses the resource a session attaches to, within one connection. */
export type SessionTarget = SessionSchemas['SessionTarget'];

/** The JSON payload of a session socket's text frames. */
export type ControlMessage = SessionSchemas['SessionControlMessage'];

/**
 * Parameters of a log session: the generated query-argument type, minus the
 * fields that address the target. They are fixed when the remote log stream is
 * created, so changing one means opening a new session.
 */
export type LogSessionParams = Omit<
  OpenLogSessionApiArg,
  'connectionId' | 'resource' | 'name' | 'namespace' | 'container'
>;

/** Parameters of a terminal session, on the same basis. */
export type TerminalSessionParams = Omit<
  OpenTerminalSessionApiArg,
  'connectionId' | 'resource' | 'name' | 'namespace' | 'container'
>;

/**
 * Close codes the server sends, from the private-use range RFC 6455 §7.4.2
 * reserves for applications.
 */
export const CLOSE_SESSION_ENDED = 4000;
export const CLOSE_SESSION_FAILED = 4001;

const SESSIONS_BASE = '/api/integrations/connections';

/**
 * Builds the absolute WebSocket URL for a session.
 *
 * The socket is same-origin, so the browser sends the Meshery session cookie on
 * the handshake and the server's usual auth chain applies. The server also
 * enforces a same-origin check on that handshake, which is why this must be
 * built from `window.location` rather than from a configured base URL.
 */
export const sessionSocketUrl = (
  connectionId: string,
  kind: SessionKind,
  target: SessionTarget,
  params: LogSessionParams & TerminalSessionParams = {},
): string => {
  const scheme = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const url = new URL(
    `${scheme}//${window.location.host}${SESSIONS_BASE}/${connectionId}/sessions/${kind}`,
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
