import { isRejectedWithValue, Middleware } from '@reduxjs/toolkit';

/**
 * RTK Query middleware that intercepts 401 responses.
 * Instead of doing a hard page reload (legacy behavior), it dispatches
 * a SESSION_EXPIRED action so the UI can show a re-auth modal.
 */
export const authMiddleware: Middleware = (storeApi) => (next) => (action) => {
  if (isRejectedWithValue(action)) {
    const status = action.payload?.status || action.payload?.originalStatus;

    if (status === 401) {
      storeApi.dispatch({ type: 'SESSION_EXPIRED' });
    }
  }

  return next(action);
};
