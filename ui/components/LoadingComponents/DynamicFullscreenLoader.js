import { useEffect } from 'react';

export const DynamicFullScreenLoader = ({ children, isLoading }) => {
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
