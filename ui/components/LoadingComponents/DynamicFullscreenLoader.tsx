import React, { useEffect } from 'react';

type DynamicFullScreenLoaderProps = {
  children: React.ReactNode;
  isLoading: boolean;
};

/* eslint-disable no-unused-vars */
declare global {
  interface Window {
    Loader?: {
      show: () => void;
      hide: () => void;
    };
  }
}
/* eslint-enable no-unused-vars */

export const DynamicFullScreenLoader = ({ children, isLoading }: DynamicFullScreenLoaderProps) => {
  useEffect(() => {
    const loader = window.Loader;
    if (!loader) {
      return;
    }
    if (isLoading) {
      loader.show();
    }

    if (loader && !isLoading) {
      loader.hide();
    }
  }, [isLoading]);

  if (!isLoading) return children;

  return null;
};
