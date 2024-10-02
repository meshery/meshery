/* eslint-disable no-unused-vars */
//NOTE: This file is being refactored to use the new notification center
import { IconButton } from '@material-ui/core';
import { ToggleButtonGroup } from '@mui/material';
import { useSnackbar } from 'notistack';
import { iconMedium } from '../../css/icons.styles';
import CloseIcon from '@material-ui/icons/Close';
import { useDispatch } from 'react-redux';
import moment from 'moment';
import { v4 } from 'uuid';
import { NOTIFICATION_STATUS } from '../../lib/event-types';
import { store as rtkStore } from '../../store/index';
import { toggleNotificationCenter } from '../../store/slices/events';
import { NOTIFICATION_CENTER_TOGGLE_CLASS } from '../../components/NotificationCenter/constants';
import React from 'react';
import BellIcon from '../../assets/icons/BellIcon';
import { AddClassRecursively } from '../Elements';

/**
 * A React hook to facilitate emitting events from the client.
 * The hook takes care of storing the events on the client through Redux
 * and also notifying the user through snackbars and the notification center.
 *
 * @returns {Object} An object with the `notify` property.
 */
export const useNotification = () => {
  const x = useSnackbar();
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  /**
   * Opens an event in the notification center.
   *
   * @param {string} eventId - The ID of the event to be opened.
   */
  const openEvent = (eventId) => {
    rtkStore.dispatch(toggleNotificationCenter());
  };

  /**
   * Notifies and stores the event.
   *
   * @param {Object} options - Options for the event notification.
   * @param {string} options.id - A unique ID for the event. If not provided, a random ID will be generated.
   * @param {string} options.message - Summary of the event.
   * @param {string} options.details - Description of the event.
   * @param {Object} options.event_type - The type of the event.
   * @param {number} options.timestamp - UTC timestamp for the event. If not provided, it is generated on the client.
   * @param {Object} options.customEvent - Additional properties related to the event.
   * @param {boolean} options.showInNotificationCenter - Whether to show the event in the notification center. Defaults to `true`.
   * @param {boolean} options.pushToServer - Whether to push the event to the server. Defaults to `false`.
   */
  const notify = ({
    id = null,
    message,
    dataTestID = 'notify',
    details = null,
    event_type,
    timestamp = null,
    customEvent = null,
    showInNotificationCenter = false,
    pushToServer = false,
  }) => {
    timestamp = timestamp ?? moment.utc().valueOf();
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
                  onClick={() => openEvent(id)}
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
  };

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
