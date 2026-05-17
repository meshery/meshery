import { useEffect, useRef } from 'react';

type TimeoutCallback = () => void;

const useTimeout = (
  callback: TimeoutCallback,
  delay: number | null,
  deps: ReadonlyArray<unknown> = [],
): void => {
  const savedCallback = useRef<TimeoutCallback>(callback);

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  useEffect(() => {
    if (delay === null) {
      return undefined;
    }
    const id = setTimeout(() => savedCallback.current(), delay);
    return () => clearTimeout(id);
  }, [delay, ...deps]);
};

export default useTimeout;
