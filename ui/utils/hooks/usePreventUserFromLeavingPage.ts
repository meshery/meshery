import SingletonRouter, { Router } from 'next/router';
import React from 'react';

/**
 * Restrict the user to navigate away to another page
 * @param {boolean} [preventLeave]
 * @returns {React.Dispatch<React.SetStateAction<boolean>>}
 */
export default function usePreventUserFromLeavingPage(
  preventLeave?: boolean,
): React.Dispatch<React.SetStateAction<boolean>> {
  const confirmationMsg = 'You might have some unsaved changes. Are you sure you want to leave?';
  const [shouldPreventLeaving, setShouldPreventLeaving] = React.useState(!!preventLeave);

  React.useEffect(() => {
    let beforeUnloadHandler: (() => string) | null = null;
    let routerChangeHandler: typeof Router.prototype.change | null = null;

    if (shouldPreventLeaving) {
      // Prevents tab quit / tab refresh
      if (typeof window !== 'undefined') {
        beforeUnloadHandler = function () {
          return confirmationMsg;
        };
        window.onbeforeunload = beforeUnloadHandler;
      }

      if (SingletonRouter.router) {
        // Prevents next routing
        routerChangeHandler = (...args) => {
          if (typeof window !== 'undefined' && typeof window.confirm === 'function') {
            if (window.confirm(confirmationMsg)) {
              return Router.prototype.change.apply(SingletonRouter.router, args);
            }
          }
          return Promise.resolve(false);
        };
        SingletonRouter.router.change = routerChangeHandler;
      }
    } else if (typeof window !== 'undefined') {
      window.onbeforeunload = () => {};
    }

    // unmount
    return () => {
      if (
        beforeUnloadHandler &&
        typeof window !== 'undefined' &&
        window.onbeforeunload === beforeUnloadHandler
      ) {
        window.onbeforeunload = null;
      }
      if (routerChangeHandler && SingletonRouter?.router?.change === routerChangeHandler) {
        delete SingletonRouter.router.change;
      }
    };
  }, [shouldPreventLeaving]);

  return setShouldPreventLeaving;
}
