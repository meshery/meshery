//NOTE: This file is being refactored to use the new notification center

import { useCallback } from 'react';
import type { ComponentType } from 'react';
import { CloseIcon, IconButton, ToggleButtonGroup } from '@sistent/sistent';
import { useSnackbar } from 'notistack';
import type { VariantType } from 'notistack';
import { iconMedium } from '../../css/iconSizes';
import { v4 } from 'uuid';
import { store as rtkStore } from '../../store/index';
import { toggleNotificationCenter } from '../../store/slices/events';
import { NOTIFICATION_CENTER_TOGGLE_CLASS } from '../../components/layout/NotificationCenter/constants';
import BellIcon from '../../assets/icons/BellIcon';
import { AddClassRecursively } from '../Elements';
import { formatApiError } from '../helpers/meshkitError';

type NotifyEventType = string | { type?: string | null } | null | undefined;

type NotifyPayload = {
  id?: string | null;
  message: string;
  dataTestID?: string;
  event_type?: NotifyEventType;
  showInNotificationCenter?: boolean;
};

const openEvent = () => {
  rtkStore.dispatch(toggleNotificationCenter());
};

/**
 * A React hook to facilitate emitting events from the client.
 * The hook takes care of storing the events on the client through Redux
 * and also notifying the user through snackbars and the notification center.
 *
 * @returns {Object} An object with the `notify` property.
 */
export const useNotification = () => {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  // Memoized so consumers can list `notify` in hook dep arrays without
  // invalidating their memos every render.
  const notify = useCallback(
    ({
      id = null,
      message,
      dataTestID = 'notify',
      event_type,
      showInNotificationCenter = false,
    }: NotifyPayload) => {
      id = id || v4();

      enqueueSnackbar(message, {
        //NOTE: Need to Consolidate the variant and event_type
        variant: typeof event_type === 'string' ? event_type : event_type?.type,
        action: function Action(key) {
          return (
            <ToggleButtonGroup data-testid={dataTestID}>
              {showInNotificationCenter && (
                <AddClassRecursively className={NOTIFICATION_CENTER_TOGGLE_CLASS}>
                  <IconButton
                    key={`openevent-${id}`}
                    aria-label="Open"
                    color="inherit"
                    onClick={() => openEvent()}
                  >
                    <BellIcon {...iconMedium} />
                  </IconButton>
                </AddClassRecursively>
              )}
              <IconButton
                key={`closeevent-${id}`}
                aria-label="Close"
                color="inherit"
                onClick={() => closeSnackbar(key)}
              >
                <CloseIcon style={iconMedium} />
              </IconButton>
            </ToggleButtonGroup>
          );
        },
      });
    },
    [enqueueSnackbar, closeSnackbar],
  );

  return {
    notify,
  };
};

type NotifyFn = ReturnType<typeof useNotification>['notify'];

/**
 * A higher-order component that provides the `notify` function as a prop to a class-based component.
 *
 * @param Component - The class-based component to be wrapped.
 * @returns The wrapped component with the `notify` prop injected.
 */
export function withNotify<P extends { notify: NotifyFn }>(
  Component: ComponentType<P>,
): ComponentType<Omit<P, 'notify'>> {
  return function WrappedWithNotify(props: Omit<P, 'notify'>) {
    const { notify } = useNotification();
    return <Component {...(props as P)} notify={notify} />;
  };
}

export const useNotificationHandlers = () => {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  const handleNotification = useCallback(
    (type: VariantType, msg: string | { response?: { data?: string } }) => {
      let message = typeof msg === 'string' ? msg : msg?.response?.data;
      enqueueSnackbar(message, {
        variant: type,
        action: (key) => (
          <IconButton
            key="close"
            aria-label="Close"
            color="inherit"
            onClick={() => closeSnackbar(key)}
          >
            <CloseIcon />
          </IconButton>
        ),
        autoHideDuration: 8000,
        style: {
          display: 'flex',
          flexWrap: 'nowrap',
        },
      });
    },
    [enqueueSnackbar, closeSnackbar],
  );

  const handleSuccess = useCallback(
    (message: string) => {
      handleNotification('success', message);
    },
    [handleNotification],
  );

  const handleError = useCallback(
    (message: string) => {
      handleNotification('error', message);
    },
    [handleNotification],
  );

  const handleInfo = useCallback(
    (message: string) => {
      handleNotification('info', message);
    },
    [handleNotification],
  );

  const handleWarn = useCallback(
    (message: string) => {
      handleNotification('warning', message);
    },
    [handleNotification],
  );

  /**
   * Surface an RTK Query error in a toast, automatically consuming the
   * structured `meshkit` envelope set by `@meshery/schemas` v1.2.2's
   * `transformErrorResponse` wrapper. When MeshKit metadata is present the
   * toast renders:
   *   - `meshkit.message` as a bold title,
   *   - `meshkit.suggestedRemediation` as a bullet list (one entry per line),
   *   - `meshkit.code` as a muted reference for support tickets.
   * When the error is not a MeshKit envelope (network failure, legacy
   * endpoint, etc.) the helper falls back to the prior single-line behavior
   * — `error.data` / `error.message` / the supplied fallback title — so it
   * is safe to call unconditionally.
   */
  const notifyApiError = useCallback(
    (error: unknown, fallbackTitle?: string) => {
      const { message } = formatApiError(error, fallbackTitle);
      handleNotification('error', message);
    },
    [handleNotification],
  );

  return { handleSuccess, handleError, handleInfo, handleWarn, notifyApiError };
};
