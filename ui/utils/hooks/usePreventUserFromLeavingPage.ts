import SingletonRouter, { Router } from 'next/router';
import React from 'react';

const CONFIRMATION_MSG = 'You might have some unsaved changes. Are you sure you want to leave?';
const activeGuardTokens = new Set<symbol>();

const beforeUnloadHandler = () => CONFIRMATION_MSG;

const routerChangeHandler: typeof Router.prototype.change = (...args) => {
  if (typeof window !== 'undefined' && typeof window.confirm === 'function') {
    if (window.confirm(CONFIRMATION_MSG)) {
      return Router.prototype.change.apply(SingletonRouter.router, args);
    }
  }
  return Promise.resolve(false);
};

/**
 * Restrict the user to navigate away to another page
 * @param {boolean} [preventLeave]
 * @returns {React.Dispatch<React.SetStateAction<boolean>>}
 */
export default function usePreventUserFromLeavingPage(
  preventLeave?: boolean,
): React.Dispatch<React.SetStateAction<boolean>> {
  const [shouldPreventLeaving, setShouldPreventLeaving] = React.useState(!!preventLeave);

  React.useEffect(() => {
    if (!shouldPreventLeaving) {
      return undefined;
    }

    const token = Symbol('preventLeaveGuard');
    activeGuardTokens.add(token);

    // Prevents tab quit / tab refresh
    if (typeof window !== 'undefined') {
      window.onbeforeunload = beforeUnloadHandler;
    }

    if (SingletonRouter.router) {
      // Prevents next routing
      SingletonRouter.router.change = routerChangeHandler;
    }

    // unmount
    return () => {
      activeGuardTokens.delete(token);
      if (activeGuardTokens.size > 0) {
        return;
      }
      if (typeof window !== 'undefined' && window.onbeforeunload === beforeUnloadHandler) {
        window.onbeforeunload = null;
      }
      if (SingletonRouter?.router?.change === routerChangeHandler) {
        delete SingletonRouter.router.change;
      }
    };
  }, [shouldPreventLeaving]);

  return setShouldPreventLeaving;
}
