import React, { useEffect, useRef, useState } from 'react';

const layoutClassName = 'react-grid-layout';
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
    const { measureBeforeMount = false, className, style, debounceTimeout = 0, ...rest } = props;

    const [width, setWidth] = useState(DEFAULT_WIDTH);
    const [mounted, setMounted] = useState(false);

    const elementRef = useRef<HTMLDivElement>(null);
    const rafIdRef = useRef<number | null>(null);
    const debounceTimerRef = useRef<number | null>(null);

    useEffect(() => setMounted(true), []);

    useEffect(() => {
      const node = elementRef.current;
      if (!(node instanceof HTMLElement)) return;

      const scheduleWidthUpdate = (newWidth: number) => {
        if (debounceTimeout > 0) {
          if (debounceTimerRef.current) clearTimeout(debounceTimerRef.current);
          debounceTimerRef.current = window.setTimeout(() => {
            setWidth((w) => (w === newWidth ? w : newWidth));
            debounceTimerRef.current = null;
          }, debounceTimeout);
        } else {
          setWidth((w) => (w === newWidth ? w : newWidth));
        }
      };

      const observer = new ResizeObserver((entries) => {
        if (!entries[0]) return;
        if (rafIdRef.current) cancelAnimationFrame(rafIdRef.current);
        rafIdRef.current = requestAnimationFrame(() => {
          scheduleWidthUpdate(Math.floor(entries[0].contentRect.width));
          rafIdRef.current = null;
        });
      });

      observer.observe(node);

      return () => {
        if (rafIdRef.current) cancelAnimationFrame(rafIdRef.current);
        if (debounceTimerRef.current) clearTimeout(debounceTimerRef.current);
        observer.disconnect();
      };
    }, [debounceTimeout, mounted]);

    return (
      <div
        className={[className, layoutClassName].filter(Boolean).join(' ')}
        style={style}
        ref={elementRef}
      >
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
