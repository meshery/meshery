import React, { useEffect } from 'react';

type Loader = {
  show?: () => void;
  hide?: () => void;
  setMessage?: (message: string) => void;
  resetMessage?: () => void;
};

declare global {
  interface Window {
    Loader?: Loader;
  }
}

type DynamicFullScreenLoaderProps = {
  children: React.ReactNode;
  isLoading: boolean;
  message?: string;
};

export const DynamicFullScreenLoader = ({
  children,
  isLoading,
  message,
}: DynamicFullScreenLoaderProps) => {
  useEffect(() => {
    const loader = window.Loader;
    if (!loader) {
      return;
    }

    if (isLoading) {
      loader.show();
      if (message) {
        loader.setMessage?.(message);
      } else {
        loader.resetMessage?.();
      }
      return;
    }

    loader.resetMessage?.();
    loader.hide();
  }, [isLoading, message]);

  if (!isLoading) return children;

  return null;
};
