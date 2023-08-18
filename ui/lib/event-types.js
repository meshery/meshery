import WarningIcon from '@material-ui/icons/Warning';
import ErrorIcon from '@material-ui/icons/Error';
import InfoIcon from '@material-ui/icons/Info';

// TODO: rename to serverEventType or better integrate wiht EVENT_TYPES
export const eventTypes = {
  0: {
    icon: InfoIcon,
    type: 'success',
  },
  1: {
    icon: WarningIcon,
    type: 'warning',
  },
  2: {
    icon: ErrorIcon,
    type: 'error',
  },
};

export const ALL_EVENT_TYPES = new Set(Object.values(eventTypes).map(e=>e.type))

export const EVENT_TYPES = {
  SUCCESS: {
    icon: InfoIcon,
    type: 'success',
  },
  DEFAULT: {
    icon: InfoIcon,
    type: 'default',
  },
  INFO: {
    icon: InfoIcon,
    type: 'info',
  },
  WARNING: {
    icon: WarningIcon,
    type: 'warning',
  },
  ERROR: {
    icon: ErrorIcon,
    type: 'error',
  },
}