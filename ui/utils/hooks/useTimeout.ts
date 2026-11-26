import { useEffect, useRef } from 'react';

type TimeoutCallback = () => void;

const useTimeout = (
  callback: TimeoutCallback | null | undefined,
  delay: number | null | undefined,
  deps: ReadonlyArray<unknown> = [],
): void => {
  const savedCallback = useRef<TimeoutCallback | null | undefined>(callback);

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  const hasCallback = typeof callback === 'function';
  const normalizedDeps = Array.isArray(deps) ? deps : [];

  useEffect(() => {
    if (typeof delay !== 'number' || !Number.isFinite(delay) || delay < 0 || !hasCallback) {
      return undefined;
    }

    const id = setTimeout(() => {
      if (typeof savedCallback.current === 'function') {
        savedCallback.current();
      }
    }, delay);

    return () => clearTimeout(id);
  }, [delay, hasCallback, ...normalizedDeps]);
};

export default useTimeout;
