import SingletonRouter, { Router } from 'next/router';
import React from 'react';

/**
 * Restrict the user to navigate away to another page
 * @param {bool} preventLeave
 * @returns
 */
export default function usePreventUserFromLeavingPage(preventLeave) {
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

    if (shouldPreventLeaving) {
      // Prevents next routing
      SingletonRouter.router.change = (...args) => {
        if (confirm(confirmationMsg)) {
          return Router.prototype.change.apply(SingletonRouter.router, args);
        } else {
          /* eslint-disable */
          return new Promise((resolve, reject) => resolve(false));
        }
      };
    }

    // unmount
    return () => {
      delete SingletonRouter.router.change;
    };
  }, [shouldPreventLeaving]);

  return setShouldPreventLeaving;
}
