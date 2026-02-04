import { useCallback, useRef, useEffect } from 'react';

type TimeoutId = ReturnType<typeof setTimeout> | null;

const useDebouncedCallback = <T extends (...args: any[]) => void>(
  callback: T,
  delay = 300,
): ((...args: Parameters<T>) => void) => {
  const timeoutRef = useRef<TimeoutId>(null);

  const debouncedFunction = useCallback(
    (...args: Parameters<T>) => {
      // Clear any existing timeout
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      // Set a new timeout to call the callback
      timeoutRef.current = setTimeout(() => {
        callback(...args);
      }, delay);
    },
    [callback, delay],
  );

  // Clean up timeout on unmount
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return debouncedFunction;
};

export default useDebouncedCallback;
