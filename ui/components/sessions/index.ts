/**
 * The public surface of Meshery's interactive sessions.
 *
 * This barrel is what remote components consume: it is registered in
 * `ui/remote-component.config.js` under `@meshery/sessions`, so an extension
 * imports the very components and hooks Meshery itself renders rather than
 * carrying a second implementation of terminals and log streaming.
 *
 * The two entry points an extension normally wants:
 *
 *   - `sessions` — the same handle for callers that are not React components,
 *     such as a Kanvas context-menu action dispatched from a plain function.
 *
 *   - `useSessions()` — an imperative handle on the shared bottom drawer
 *     (`openSession` / `closeSession` / `hasSession`). Use it to launch a shell
 *     or a log tail from a graph node, a context menu, or a button, and to render
 *     an indicator on a resource that already has a session attached. The drawer
 *     is mounted app-wide in `pages/_app.tsx`, so this works from any page.
 *
 *   - `SessionPanel` — one embedded session, for a host that wants to own the
 *     chrome itself rather than use the shared drawer.
 *
 * Wire types come from `@meshery/schemas` (the `v1beta1/session` construct);
 * they are re-exported here so an extension need not depend on schemas directly.
 */

export { default as SessionPanel } from './SessionPanel';
export type { SessionPanelProps } from './SessionPanel';

export { default as SessionsProvider, useSessions, sessions } from './SessionsProvider';
export type { OpenSessionRequest, SessionsContextValue } from './SessionsProvider';

export { default as SessionActionsCell } from './SessionActionsCell';

export { default as TerminalSession } from './TerminalSession';
export { default as LogSession } from './LogSession';

export type {
  ControlMessage,
  ControlType,
  SessionCapabilities,
  SessionKind,
  SessionTarget,
} from 'lib/sessions/protocol';
