import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';

/**
 *  Show Snackbar when error occurs. Can be used in catch blocks
 *  of functions.
 * @deprecated use the useNotification hook instead
 *
 */
function HandleError() {
  const { notify } = useNotification();
  /**
   *
   * @param {Object} err
   * @param {string} prefixMessage
   * @param {("error"|"warning")} variant
   */
  const errorH = (err, prefixMessage, variant) => {
    console.error('an error occured with severity: ', variant, { err });
    return notify({
      message: `${prefixMessage}: ${err?.message}`,
      event_type: EVENT_TYPES.ERROR,
      details: err.toString(),
    });
  };

  return errorH;
}

export default HandleError;
