import { useEffect, useRef } from 'react';

type EventTargetLike = EventTarget | null | undefined;

const useEventListener = <K extends keyof WindowEventMap>(
  eventName: K,
  handler: (event: WindowEventMap[K]) => void,
  target: EventTargetLike = typeof window !== 'undefined' ? window : null,
  options?: boolean | AddEventListenerOptions,
): void => {
  const savedHandler = useRef(handler);

  useEffect(() => {
    savedHandler.current = handler;
  }, [handler]);

  useEffect(() => {
    if (!target || typeof target.addEventListener !== 'function') {
      return undefined;
    }
    const listener = (event: Event) => savedHandler.current(event as WindowEventMap[K]);
    target.addEventListener(eventName, listener, options);
    return () => target.removeEventListener(eventName, listener, options);
  }, [eventName, target, options]);
};

export default useEventListener;
