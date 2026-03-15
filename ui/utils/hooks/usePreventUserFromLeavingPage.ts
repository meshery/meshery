import SingletonRouter from 'next/router';
import React from 'react';

/**
 * Restrict the user to navigate away to another page
 * @param {boolean} preventLeave
 * @returns
 */
export default function usePreventUserFromLeavingPage(preventLeave: boolean) {
  const confirmationMsg = 'You might have some unsaved changes. Are you sure you want to leave?';
  const [shouldPreventLeaving, setShouldPreventLeaving] = React.useState(!!preventLeave);

  React.useEffect(() => {
    // Prevents tab quit / tab refresh
    if (shouldPreventLeaving) {
      // Adding window alert if the shop quits without saving
      window.onbeforeunload = function () {
        return confirmationMsg;
      };
    } else {
      window.onbeforeunload = () => {};
    }

    let originalChange: ((..._args: any[]) => any) | undefined;

    if (shouldPreventLeaving && SingletonRouter.router) {
      // Prevents next routing by monkey-patching the internal change method
      const routerAny = SingletonRouter.router as any;
      originalChange = routerAny.change;
      routerAny.change = (...args: any[]) => {
        if (confirm(confirmationMsg)) {
          return originalChange?.apply(routerAny, args);
        }
        return Promise.resolve(false);
      };
    }

    // unmount / dependency cleanup
    return () => {
      if (SingletonRouter.router && originalChange) {
        (SingletonRouter.router as any).change = originalChange;
      }
      window.onbeforeunload = null;
    };
  }, [shouldPreventLeaving]);

  return setShouldPreventLeaving;
}
