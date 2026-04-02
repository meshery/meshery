import { isRejectedWithValue, Middleware } from '@reduxjs/toolkit';

/**
 * RTK Query middleware that automatically dispatches error events
 * for failed API calls, so they appear in the notification system.
 *
 * Skips 401 errors (handled by authMiddleware).
 */
export const rtkErrorMiddleware: Middleware = () => (next) => (action) => {
  if (isRejectedWithValue(action)) {
    const status = action.payload?.status || action.payload?.originalStatus;

    // 401 is handled by authMiddleware
    if (status === 401) {
      return next(action);
    }

    const endpointName = action.meta?.arg?.endpointName || 'API call';
    const errorMessage =
      action.payload?.data?.message ||
      action.payload?.data?.error ||
      action.payload?.error ||
      `Request failed (${status || 'unknown status'})`;

    console.error(`[RTK Query] ${endpointName} failed:`, errorMessage);
  }

  return next(action);
};
