import React, { useEffect, useRef, useState } from 'react';
const DEFAULT_WIDTH = 1280;

interface DebounceWidthProviderProps {
  measureBeforeMount?: boolean;
  className?: string;
  style?: React.CSSProperties;
  debounceTimeout?: number;
}

const debounceWidthProvider = <P extends object>(
  ComposedComponent: React.ComponentType<P & { width: number }>,
) => {
  const Wrapper = (props: Omit<P, 'width'> & DebounceWidthProviderProps) => {
    const { measureBeforeMount = false, className, style, debounceTimeout = 100, ...rest } = props;

    const [width, setWidth] = useState(DEFAULT_WIDTH);
    const [mounted, setMounted] = useState(!measureBeforeMount);

    const elementRef = useRef<HTMLDivElement>(null);
    const rafIdRef = useRef<number | null>(null);
    const debounceTimerRef = useRef<number | null>(null);
    const hasMeasuredRef = useRef(false);

    useEffect(() => {
      const node = elementRef.current;
      if (!(node instanceof HTMLElement)) return;

      const scheduleWidthUpdate = (newWidth: number) => {
        const performUpdate = () => {
          setWidth(newWidth);
          setMounted(true);
          hasMeasuredRef.current = true;
        };

        if (hasMeasuredRef.current && debounceTimeout > 0) {
          if (debounceTimerRef.current !== null) clearTimeout(debounceTimerRef.current);
          debounceTimerRef.current = window.setTimeout(() => {
            performUpdate();
            debounceTimerRef.current = null;
          }, debounceTimeout);
        } else {
          performUpdate();
        }
      };

      const observer = new ResizeObserver((entries) => {
        if (!entries[0]) return;
        if (rafIdRef.current !== null) cancelAnimationFrame(rafIdRef.current);
        rafIdRef.current = requestAnimationFrame(() => {
          scheduleWidthUpdate(Math.floor(entries[0].contentRect.width));
          rafIdRef.current = null;
        });
      });

      observer.observe(node);

      return () => {
        if (rafIdRef.current !== null) cancelAnimationFrame(rafIdRef.current);
        if (debounceTimerRef.current !== null) clearTimeout(debounceTimerRef.current);
        observer.disconnect();
      };
    }, [debounceTimeout]);

    return (
      <div className={className} style={style} ref={elementRef}>
        {(!measureBeforeMount || mounted) && <ComposedComponent {...(rest as P)} width={width} />}
      </div>
    );
  };

  Wrapper.displayName = `DebouncedWidthProvider(${
    ComposedComponent.displayName || ComposedComponent.name || 'Component'
  })`;
  return Wrapper;
};

export default debounceWidthProvider;
