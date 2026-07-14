import { createContext, useContext } from 'react';

/**
 * The element in the panel's header bar that the focused console projects its
 * controls into.
 *
 * The container select and the search box belong to a *console*, but they read as
 * chrome of the *panel*, so they are rendered by the console and portaled into the
 * header. A portal rather than lifted state because the controls are backed by
 * per-console state — which container this shell is attached to, what this log
 * tail is searching for — and hoisting that into the shell would mean the shell
 * holding a map of state keyed by console id, and handing it back down.
 *
 * Null when a console is embedded without the shell around it (an extension
 * rendering `ConsolePanel` on its own), in which case the controls stay inline in
 * the console's own toolbar.
 */
export const ConsoleHeaderSlotContext = createContext<HTMLElement | null>(null);

export const useConsoleHeaderSlot = (): HTMLElement | null => useContext(ConsoleHeaderSlotContext);
