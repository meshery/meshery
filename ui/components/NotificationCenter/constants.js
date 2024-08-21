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
