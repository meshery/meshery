import { NOTIFICATIONCOLORS } from '../../themes';
import AlertIcon from '../../assets/icons/AlertIcon';
import ErrorIcon from '../../assets/icons/ErrorIcon.js';
import { Colors } from '../../themes/app';
import ReadIcon from '../../assets/icons/ReadIcon';
import { InfoIcon } from '@layer5/sistent';

export const SEVERITY = {
  INFO: 'informational',
  ERROR: 'error',
  WARNING: 'warning',
  SUCCESS: 'success',
};

// This class is required to add to any svg or button that opens notification center
// To prevent the clickaway listner from blocking it
export const NOTIFICATION_CENTER_TOGGLE_CLASS = 'toggle-notification-center';

export const SEVERITY_TO_NOTIFICATION_TYPE_MAPPING = {
  [SEVERITY.INFO]: 'info',
  [SEVERITY.ERROR]: 'error',
  [SEVERITY.WARNING]: 'warning',
  [SEVERITY.SUCCESS]: 'success',
};

export const STATUS = {
  READ: 'read',
  UNREAD: 'unread',
};

export const STATUS_STYLE = {
  [STATUS.READ]: {
    icon: ReadIcon,
    color: Colors.charcoal,
    darkColor: '#BCC7CC',
  },
};

export const SEVERITY_STYLE = {
  [SEVERITY.INFO]: {
    icon: InfoIcon,
    color: NOTIFICATIONCOLORS.INFO,
    darkColor: NOTIFICATIONCOLORS.INFO,
  },
  [SEVERITY.ERROR]: {
    icon: ErrorIcon,
    color: NOTIFICATIONCOLORS.ERROR,
    darkColor: NOTIFICATIONCOLORS.ERROR_DARK,
  },
  [SEVERITY.WARNING]: {
    icon: AlertIcon,
    color: NOTIFICATIONCOLORS.WARNING,
    darkColor: NOTIFICATIONCOLORS.WARNING,
  },
  [SEVERITY.SUCCESS]: {
    icon: InfoIcon,
    color: NOTIFICATIONCOLORS.SUCCESS,
    darkColor: NOTIFICATIONCOLORS.SUCCESS,
  },
};

// Validate event against EVENT_SCHEMA and return [isValid,validatedEvent]
export const validateEvent = (event) => {
  // const eventCopy = { ...event };
  // event.status = eventCopy.status.trim() || STATUS.UNREAD;
  // event.severity = eventCopy.severity.trim() || SEVERITY.INFO;
  const valid = event !== null && event !== undefined;
  return [valid, event];
};

// return validated events (adds default values if not present)
export const validateEvents = (events) => {
  return events
    .map((event) => {
      const [isValid, validatedEvent] = validateEvent(event);
      return isValid ? validatedEvent : null;
    })
    .filter((event) => event);
};

export const validateEventMetadata = (metadata) => {
  const valid = metadata !== null && metadata !== undefined;
  return [valid, metadata];
};
