import { useEffect, useRef } from 'react';

const useEventListener = <E extends Event = Event>(
  eventName: string,
  handler: (event: E) => void,
  target: EventTarget | null = typeof window !== 'undefined' ? window : null,
): void => {
  const savedHandler = useRef(handler);

  useEffect(() => {
    savedHandler.current = handler;
  }, [handler]);

  useEffect(() => {
    if (!target || typeof target.addEventListener !== 'function') {
      return undefined;
    }
    const listener = (event: Event) => savedHandler.current(event as E);
    target.addEventListener(eventName, listener);
    return () => target.removeEventListener(eventName, listener);
  }, [eventName, target]);
};

export default useEventListener;
