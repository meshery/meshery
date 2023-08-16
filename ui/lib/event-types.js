import WarningIcon from '@material-ui/icons/Warning';
import ErrorIcon from '@material-ui/icons/Error';
import InfoIcon from '@material-ui/icons/Info';

// TODO: why are we using numbers as key ?
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
    type: 'success',
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