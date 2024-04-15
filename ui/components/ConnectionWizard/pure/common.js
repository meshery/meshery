/* eslint-disable react/display-name */
import { IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';
import { EVENT_TYPES } from '../../../lib/event-types';

/**
 * A function that generates a close button component for a snackbar.
 * @param {function} closeSnackbar - A function to close the snackbar.
 * @returns {function} - A function that generates a close button component.
 */
export const closeButtonForSnackbarAction = (closeSnackbar) => (key) => (
  <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
    <CloseIcon />
  </IconButton>
);

/**
 * A function that generates a success handler for notifying.
 * @param {function} notify - A function to notify.
 * @param {string} msg - The message to be displayed.
 * @param {function} cb - An optional callback function.
 * @returns {function} - A success handler function.
 */
export const successHandlerGenerator = (notify, msg, cb) => (res) => {
  if (res !== undefined) {
    if (cb !== undefined) cb(res);
    if (typeof res == 'object') {
      res = JSON.stringify(res);
    }
    notify({ message: `${msg}`, details: `${res}`, event_type: EVENT_TYPES.SUCCESS });
  }
};

/**
 * A function that generates an error handler for notifying.
 * @param {function} notify - A function to notify.
 * @param {string} msg - The message to be displayed.
 * @param {function} cb - An optional callback function.
 * @returns {function} - An error handler function
 */
export const errorHandlerGenerator = (notify, msg, cb) => (err) => {
  if (cb !== undefined) cb(err);
  err = typeof err !== 'string' ? err.toString() : err;
  notify({ message: `${msg}`, details: err, event_type: EVENT_TYPES.ERROR });
};
