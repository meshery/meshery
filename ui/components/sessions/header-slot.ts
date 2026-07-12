import { createContext, useContext } from 'react';

/**
 * The element in the panel's header bar that the focused session projects its
 * controls into.
 *
 * The container select and the search box belong to a *session*, but they read as
 * chrome of the *panel*, so they are rendered by the session and portaled into the
 * header. A portal rather than lifted state because the controls are backed by
 * per-session state — which container this shell is attached to, what this log
 * tail is searching for — and hoisting that into the shell would mean the shell
 * holding a map of state keyed by session id, and handing it back down.
 *
 * Null when a session is embedded without the shell around it (an extension
 * rendering `SessionPanel` on its own), in which case the controls stay inline in
 * the session's own toolbar.
 */
export const SessionHeaderSlotContext = createContext<HTMLElement | null>(null);

export const useSessionHeaderSlot = (): HTMLElement | null => useContext(SessionHeaderSlotContext);
