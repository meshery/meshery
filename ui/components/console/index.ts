/**
 * The public surface of Meshery's interactive consoles.
 *
 * This barrel is what remote components consume: it is registered in
 * `ui/remote-component.config.js` under `@meshery/consoles`, so an extension
 * imports the very components and hooks Meshery itself renders rather than
 * carrying a second implementation of terminals and log streaming.
 *
 * The two entry points an extension normally wants:
 *
 *   - `consoles` — the same handle for callers that are not React components,
 *     such as a Kanvas context-menu action dispatched from a plain function.
 *
 *   - `useConsole()` — an imperative handle on the shared consoles panel
 *     (`openConsole` / `closeConsole` / `hasConsole`). Use it to launch a shell
 *     or a log tail from a graph node, a context menu, or a button, and to render
 *     an indicator on a resource that already has a console attached. The panel is
 *     mounted app-wide in `pages/_app.tsx`, so this works from any page: it docks
 *     to the foot of the content area, detaches into a floating window, and
 *     minimizes to a control at the foot of the Navigator.
 *
 *   - `ConsolePanel` — one embedded console, for a host that wants to own the
 *     chrome itself rather than use the shared panel.
 *
 * Wire types come from `@meshery/schemas` (the `v1beta1/console` construct);
 * they are re-exported here so an extension need not depend on schemas directly.
 */

export { default as ConsolePanel } from './ConsolePanel';
export type { ConsolePanelProps } from './ConsolePanel';

export { default as ConsoleProvider, useConsole, consoles } from './ConsoleProvider';
export type { OpenConsoleRequest, ConsoleContextValue } from './ConsoleProvider';

export { default as MinimizedConsoles } from './MinimizedConsoles';
export { restoreConsoles, useConsoleDock } from './dock-store';
export type { ConsoleDockState } from './dock-store';
export type { ConsolePanelMode } from './ConsoleShell';

export { default as ConsoleActionsCell } from './ConsoleActionsCell';

export { default as TerminalConsole } from './TerminalConsole';
export { default as LogConsole } from './LogConsole';

export type {
  ControlMessage,
  ControlType,
  ConsoleCapabilities,
  ConsoleKind,
  ConsoleTarget,
} from 'lib/console/protocol';
