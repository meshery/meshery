import { useEffect, useRef } from 'react';

type IntervalCallback = () => void;

const useInterval = (callback: IntervalCallback, delay: number | null): void => {
  const savedCallback = useRef<IntervalCallback>(callback);

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  useEffect(() => {
    if (delay === null) {
      return undefined;
    }
    const id = setInterval(() => savedCallback.current(), delay);
    return () => clearInterval(id);
  }, [delay]);
};

export default useInterval;
