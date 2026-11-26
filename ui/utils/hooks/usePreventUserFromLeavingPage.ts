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
    // Prevents tab quit / tab refresh
    if (typeof window !== 'undefined') {
      if (shouldPreventLeaving) {
        // Adding window alert if the user quits without saving
        window.onbeforeunload = function () {
          return confirmationMsg;
        };
      } else {
        window.onbeforeunload = () => {};
      }
    }

    if (shouldPreventLeaving && SingletonRouter.router) {
      // Prevents next routing
      SingletonRouter.router.change = (...args) => {
        if (typeof window !== 'undefined' && typeof window.confirm === 'function') {
          if (window.confirm(confirmationMsg)) {
            return Router.prototype.change.apply(SingletonRouter.router, args);
          }
          return Promise.resolve(false);
        }
        return Router.prototype.change.apply(SingletonRouter.router, args);
      };
    }

    // unmount
    return () => {
      if (typeof window !== 'undefined') {
        window.onbeforeunload = null;
      }
      if (SingletonRouter?.router) {
        delete SingletonRouter.router.change;
      }
    };
  }, [shouldPreventLeaving]);

  return setShouldPreventLeaving;
}
