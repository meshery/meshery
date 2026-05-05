import { isRejectedWithValue, Middleware } from '@reduxjs/toolkit';
import { pushEvent } from '../slices/events';

/**
 * RTK Query middleware that dispatches error events into the notification
 * system for failed API calls. Skips 401 errors (handled by authMiddleware).
 */
export const rtkErrorMiddleware: Middleware = (storeApi) => (next) => (action) => {
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

    storeApi.dispatch(
      pushEvent({
        id: `rtk-error-${Date.now()}`,
        severity: 'error',
        description: `${endpointName}: ${errorMessage}`,
        action: 'api_error',
        category: 'api',
        createdAt: new Date().toISOString(),
      }),
    );
  }

  return next(action);
};
