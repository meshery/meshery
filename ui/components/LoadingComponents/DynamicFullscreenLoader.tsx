// @ts-nocheck
import { useEffect } from 'react';

export const DynamicFullScreenLoader = ({ children, isLoading, message }) => {
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
