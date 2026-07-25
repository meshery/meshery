//NOTE: This file is being refactored to use the new notification center

import { CloseIcon, IconButton, styled, ToggleButtonGroup } from '@sistent/sistent';
import { useSnackbar } from 'notistack';
import { iconMedium } from '../../css/icons.styles';
import moment from 'moment';
import { v4 } from 'uuid';
import Router from 'next/router';
import { store as rtkStore } from '../../store/index';
import { toggleNotificationCenter } from '../../store/slices/events';
import { NOTIFICATION_CENTER_TOGGLE_CLASS } from '../../components/layout/NotificationCenter/constants';
import React from 'react';
import BellIcon from '../../assets/icons/BellIcon';
import { AddClassRecursively } from '../Elements';
import { useCallback } from 'react';
import { formatApiError } from '../helpers/meshkitError';

const openEvent = () => {
  rtkStore.dispatch(toggleNotificationCenter());
};

/**
 * Only allow same-app relative paths for snackbar navigation actions.
 * Blocks open redirects: `https://…`, `//evil`, `javascript:`, etc.
 */
export const isSafeInternalNavPath = (href) => {
  if (typeof href !== 'string') {
    return false;
  }
  const path = href.trim();
  if (!path.startsWith('/') || path.startsWith('//') || path.includes('\\')) {
    return false;
  }
  // Reject scheme-like prefixes that could slip past a leading slash after decode tricks.
  const lower = path.toLowerCase();
  if (lower.includes('javascript:') || lower.includes('data:')) {
    return false;
  }
  return true;
};

/** Same-tab action link in snackbars (underlined so it reads as a link). */
const SnackbarActionLink = styled('button')(() => ({
  color: 'inherit',
  background: 'none',
  border: 'none',
  cursor: 'pointer',
  textDecoration: 'underline',
  font: 'inherit',
  fontWeight: 600,
  padding: '0 0.5rem',
  whiteSpace: 'nowrap',
}));

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
      details = null,
      event_type,
      timestamp = null,
      customEvent = null,
      showInNotificationCenter = false,
      pushToServer = false,
      /**
       * Optional same-tab navigation action. Prefer this over markdown links in
       * `message`: ThemeResponsiveSnackbar renders messages via BasicMarkdown,
       * which always opens anchors in a new tab.
       */
      link = null,
    }) => {
      timestamp = timestamp ?? moment.utc().valueOf();
      id = id || v4();
      const safeLink =
        link?.href && isSafeInternalNavPath(link.href)
          ? { href: link.href.trim(), label: link.label || 'Open' }
          : null;

      enqueueSnackbar(message, {
        //NOTE: Need to Consolidate the variant and event_type
        variant: typeof event_type === 'string' ? event_type : event_type?.type,
        action: function Action(key) {
          return (
            <ToggleButtonGroup data-testid={dataTestID}>
              {safeLink && (
                <SnackbarActionLink
                  type="button"
                  key={`link-${id}`}
                  data-testid={`${dataTestID}-link`}
                  onClick={() => {
                    closeSnackbar(key);
                    // Same-tab SPA navigation; href already allowlisted as internal.
                    Router.push(safeLink.href);
                  }}
                >
                  {safeLink.label}
                </SnackbarActionLink>
              )}
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

/**
 * A higher-order component that provides the `notify` function as a prop to a class-based component.
 *
 * @param {React.Component} Component - The class-based component to be wrapped.
 * @returns {React.Component} The wrapped component with the `notify` prop.
 */
export function withNotify(Component) {
  return function WrappedWithNotify(props) {
    const { notify } = useNotification();
    return <Component {...props} notify={notify} />;
  };
}

export const useNotificationHandlers = () => {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  const handleNotification = useCallback(
    (type, msg) => {
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
    (message) => {
      handleNotification('success', message);
    },
    [handleNotification],
  );

  const handleError = useCallback(
    (message) => {
      handleNotification('error', message);
    },
    [handleNotification],
  );

  const handleInfo = useCallback(
    (message) => {
      handleNotification('info', message);
    },
    [handleNotification],
  );

  const handleWarn = useCallback(
    (message) => {
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
