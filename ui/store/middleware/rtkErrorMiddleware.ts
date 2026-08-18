import { isRejectedWithValue, Middleware } from '@reduxjs/toolkit';
import { pushEvent } from '../slices/events';

const ENDPOINT_RESOURCE_LABELS: Record<string, string> = {
  getEvents: 'notifications',
  getEventsSummary: 'notifications',
  deleteEvent: 'notification',
  deleteEvents: 'notifications',
  updateEvents: 'notifications',
  updateStatus: 'notification status',
  getEventFilters: 'notification filters',
  getEventConfig: 'notification settings',
  updateEventConfig: 'notification settings',
};

const getEndpointOperation = (endpointName: string) => {
  if (/^(get|list|fetch)/i.test(endpointName)) {
    return 'load';
  }

  if (/^(create|add)/i.test(endpointName)) {
    return 'create';
  }

  if (/^(update|edit|set)/i.test(endpointName)) {
    return 'update';
  }

  if (/^(delete|remove)/i.test(endpointName)) {
    return 'delete';
  }

  return 'complete';
};

const getEndpointResource = (endpointName: string) => {
  if (ENDPOINT_RESOURCE_LABELS[endpointName]) {
    return ENDPOINT_RESOURCE_LABELS[endpointName];
  }

  const withoutVerb = endpointName.replace(
    /^(get|list|fetch|create|add|update|edit|set|delete|remove)/i,
    '',
  );
  const resource = withoutVerb
    .replace(/([a-z0-9])([A-Z])/g, '$1 $2')
    .replace(/[_-]+/g, ' ')
    .trim()
    .toLowerCase();

  return resource || 'request';
};

const buildErrorDescription = (endpointName: string, payload: any) => {
  const status = payload?.status || payload?.originalStatus;
  const resource = getEndpointResource(endpointName);
  const operation = getEndpointOperation(endpointName);
  const rawErrorMessage = payload?.data?.message || payload?.data?.error || payload?.error;

  if (status === 'FETCH_ERROR') {
    return `Unable to ${operation} ${resource}. Check your connection and try again.`;
  }

  if (status === 'TIMEOUT_ERROR') {
    return `Timed out while trying to ${operation} ${resource}. Please try again.`;
  }

  if (rawErrorMessage) {
    return `Unable to ${operation} ${resource}: ${rawErrorMessage}`;
  }

  return `Unable to ${operation} ${resource}. Request failed (${status || 'unknown status'}).`;
};

/**
 * RTK Query middleware that dispatches error events into the notification
 * system for failed API calls. Every rejected status is reported, including
 * 401: this used to skip 401 and defer to an `authMiddleware` that has never
 * existed in this codebase, so an unauthorized response was swallowed and the
 * user saw nothing at all.
 */
export const rtkErrorMiddleware: Middleware = (storeApi) => (next) => (action) => {
  if (isRejectedWithValue(action)) {
    const status = action.payload?.status || action.payload?.originalStatus;
    const endpointName = action.meta?.arg?.endpointName || 'API call';
    const errorMessage =
      action.payload?.data?.message ||
      action.payload?.data?.error ||
      action.payload?.error ||
      `Request failed (${status || 'unknown status'})`;
    const description = buildErrorDescription(endpointName, action.payload);

    console.error(`[RTK Query] ${endpointName} failed:`, errorMessage);

    storeApi.dispatch(
      pushEvent({
        id: `rtk-error-${Date.now()}`,
        severity: 'error',
        description,
        action: 'api_error',
        category: 'api',
        createdAt: new Date().toISOString(),
      }),
    );
  }

  return next(action);
};
