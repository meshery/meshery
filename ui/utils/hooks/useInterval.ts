import { useEffect, useRef } from 'react';

type IntervalCallback = () => void;

const useInterval = (callback: IntervalCallback, delay: number | null | undefined): void => {
  const savedCallback = useRef<IntervalCallback>(callback);

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  useEffect(() => {
    if (
      typeof delay !== 'number' ||
      Number.isNaN(delay) ||
      delay < 0 ||
      typeof savedCallback.current !== 'function'
    ) {
      return undefined;
    }
    const id = setInterval(() => {
      if (typeof savedCallback.current === 'function') {
        savedCallback.current();
      }
    }, delay);
    return () => clearInterval(id);
  }, [delay]);
};

export default useInterval;
