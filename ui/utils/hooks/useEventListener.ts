import { useEffect, useLayoutEffect, useRef } from 'react';

type EventHandler<E extends Event = Event> = (event: E) => void;

const useIsomorphicLayoutEffect = typeof window !== 'undefined' ? useLayoutEffect : useEffect;

/**
 * useEventListener - Subscribes to DOM or custom EventTarget events declaratively.
 *
 * Guards against:
 * 1. Non-string, empty, or whitespace-only event names (logging a dev-only warning).
 * 2. Missing, nullish, or non-callable handler callbacks.
 * 3. Safely restarts listener when handler transitions between non-callable and callable.
 * 4. Gracefully no-ops when target is null, undefined, or lacks addEventListener.
 * 5. Safely cleans up listener even if target lacks removeEventListener.
 */
const useEventListener = <E extends Event = Event>(
  eventName: string,
  handler: EventHandler<E> | null | undefined,
  target: EventTarget | null = typeof window !== 'undefined' ? window : null,
): void => {
  const savedHandler = useRef<EventHandler<E> | null | undefined>(handler);

  useIsomorphicLayoutEffect(() => {
    savedHandler.current = handler;
  }, [handler]);

  const normalizedEventName = typeof eventName === 'string' ? eventName.trim() : '';
  const hasHandler = typeof handler === 'function';

  useEffect(() => {
    if (normalizedEventName === '') {
      if (typeof process !== 'undefined' && process.env?.NODE_ENV !== 'production') {
        console.warn('useEventListener: expected a non-empty string eventName.');
      }
      return undefined;
    }

    if (!hasHandler || !target || typeof target.addEventListener !== 'function') {
      return undefined;
    }

    const listener = (event: Event) => {
      if (typeof savedHandler.current === 'function') {
        savedHandler.current(event as E);
      }
    };

    target.addEventListener(normalizedEventName, listener);
    return () => {
      if (typeof target.removeEventListener === 'function') {
        target.removeEventListener(normalizedEventName, listener);
      }
    };
  }, [normalizedEventName, hasHandler, target]);
};

export default useEventListener;
