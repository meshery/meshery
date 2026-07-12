import { useSyncExternalStore } from 'react';

/**
 * The consoles panel's dock state, mirrored outside React.
 *
 * The restore affordance for a minimized panel lives at the foot of the
 * Navigator, and the Navigator is mounted *above* `ConsoleProvider` in the tree
 * (`pages/_app.tsx`), so it cannot read the provider's context. Hoisting the
 * provider over the whole shell just to feed one button would drag the panel out
 * of the content area with it — and the docked panel has to be a child of the
 * content area to share space with the page rather than cover it. So the provider
 * publishes here instead, and the button subscribes.
 */
export interface ConsoleDockState {
  /** How many consoles are open. Zero means there is nothing to restore. */
  count: number;
  /** Whether the panel is currently minimized out of view. */
  minimized: boolean;
}

const NO_CONSOLES: ConsoleDockState = { count: 0, minimized: false };

let state: ConsoleDockState = NO_CONSOLES;
let restore: (() => void) | null = null;
const listeners = new Set<() => void>();

const getState = () => state;

const subscribe = (listener: () => void) => {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
};

/**
 * Publishes dock state, skipping the notify when nothing actually changed.
 *
 * `useSyncExternalStore` re-renders on any snapshot whose *identity* differs, so
 * a fresh object per publish would re-render the Navigator on every provider
 * render.
 */
export const publishDockState = (next: ConsoleDockState) => {
  if (next.count === state.count && next.minimized === state.minimized) return;
  state = next;
  listeners.forEach((listener) => listener());
};

/**
 * Binds the mounted provider's restore action. The returned unbind clears the
 * published state as well, so a Navigator button can never outlive the panel it
 * would restore. The identity guard matters under StrictMode's double-mount,
 * where the second provider binds before the first one's cleanup runs.
 */
export const bindDock = (onRestore: () => void) => {
  restore = onRestore;
  return () => {
    if (restore !== onRestore) return;
    restore = null;
    publishDockState(NO_CONSOLES);
  };
};

export const restoreConsoles = () => restore?.();

export const useConsoleDock = (): ConsoleDockState =>
  useSyncExternalStore(subscribe, getState, getState);
