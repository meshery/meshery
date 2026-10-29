import { DependencyList, RefObject, useLayoutEffect, useState } from 'react';

export function useIsTextTruncated<T extends HTMLElement | null>(
  ref: RefObject<T>,
  deps: DependencyList = [],
): boolean {
  const [isTruncated, setIsTruncated] = useState<boolean>(false);

  useLayoutEffect(() => {
    const element = ref?.current;
    if (!element) {
      setIsTruncated(false);
      return;
    }

    const checkTruncation = () => {
      setIsTruncated(element.scrollWidth > element.clientWidth);
    };

    checkTruncation();

    const resizeObserver = new ResizeObserver(() => {
      checkTruncation();
    });

    resizeObserver.observe(element);

    return () => {
      resizeObserver.disconnect();
    };
  }, deps);

  return isTruncated;
}
