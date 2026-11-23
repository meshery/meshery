/* eslint-disable react/display-name */
import { CloseIcon, IconButton } from '@sistent/sistent';
import { EVENT_TYPES } from '../../lib/event-types';

/**
 * A function that generates a close button component for a snackbar.
 * @param {function} closeSnackbar - A function to close the snackbar.
 * @returns {function} - A function that generates a close button component.
 */
export const closeButtonForSnackbarAction = (closeSnackbar) => (key) => (
  <IconButton
    key="close"
    aria-label="Close"
    color="inherit"
    onClick={() => typeof closeSnackbar === 'function' && closeSnackbar(key)}
  >
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
    if (typeof cb === 'function') cb(res);
    let details = res;
    if (typeof res === 'object' && res !== null) {
      try {
        details = JSON.stringify(res);
      } catch {
        details = String(res);
      }
    }
    if (typeof notify === 'function') {
      notify({
        message: `${msg ?? ''}`,
        details: `${details ?? ''}`,
        event_type: EVENT_TYPES.SUCCESS,
      });
    }
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
  if (typeof cb === 'function') cb(err);
  let errorDetails = '';
  if (typeof err === 'string') {
    errorDetails = err;
  } else if (err !== null && err !== undefined) {
    try {
      errorDetails = typeof err.toString === 'function' ? err.toString() : String(err);
    } catch {
      errorDetails = String(err);
    }
  }
  if (typeof notify === 'function') {
    notify({ message: `${msg ?? ''}`, details: errorDetails, event_type: EVENT_TYPES.ERROR });
  }
};
