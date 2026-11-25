import { useCallback, useRef, useEffect } from 'react';

const DEFAULT_DEBOUNCE_DELAY = 300;

const useDebouncedCallback = <TArgs extends unknown[]>(
  callback: ((...args: TArgs) => void) | null | undefined,
  delay?: number | null | undefined,
): ((...args: TArgs) => void) => {
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const savedCallback = useRef(callback);

  const clearPendingTimeout = useCallback(() => {
    if (timeoutRef.current !== null) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = null;
    }
  }, []);

  useEffect(() => {
    savedCallback.current = callback;
    if (typeof callback !== 'function') {
      clearPendingTimeout();
    }
  }, [callback, clearPendingTimeout]);

  const normalizedDelay =
    typeof delay === 'number' && Number.isFinite(delay) && delay >= 0
      ? delay
      : DEFAULT_DEBOUNCE_DELAY;

  const debouncedFunction = useCallback(
    (...args: TArgs) => {
      clearPendingTimeout();

      if (typeof savedCallback.current !== 'function') {
        return;
      }

      timeoutRef.current = setTimeout(() => {
        timeoutRef.current = null;
        if (typeof savedCallback.current === 'function') {
          savedCallback.current(...args);
        }
      }, normalizedDelay);
    },
    [clearPendingTimeout, normalizedDelay],
  );

  // Runs on unmount and whenever delay or normalizedDelay changes, so a
  // pending invocation is also cancelled when the delay changes. Keep them in deps.
  useEffect(() => {
    return () => {
      clearPendingTimeout();
    };
  }, [clearPendingTimeout, delay, normalizedDelay]);

  return debouncedFunction;
};

export default useDebouncedCallback;
