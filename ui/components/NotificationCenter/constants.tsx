import AlertIcon from '../../assets/icons/AlertIcon';
import ErrorIcon from '../../assets/icons/ErrorIcon';
import ReadIcon from '../../assets/icons/ReadIcon';
import { InfoIcon } from '@sistent/sistent';

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
    colorKey: 'text.secondary',
  },
};

export const SEVERITY_STYLE = {
  [SEVERITY.INFO]: {
    icon: InfoIcon,
    colorKey: 'info.main',
  },
  [SEVERITY.ERROR]: {
    icon: ErrorIcon,
    colorKey: 'error.main',
  },
  [SEVERITY.WARNING]: {
    icon: AlertIcon,
    colorKey: 'warning.main',
  },
  [SEVERITY.SUCCESS]: {
    icon: InfoIcon,
    colorKey: 'success.main',
  },
};

export const eventDetailFormatterKey = ({ action, category }) => `${action}-${category}`;

export const EVENT_TYPE = {
  DEPLOY_DESIGN: {
    category: 'pattern',
    action: 'deploy',
  },
  UNDEPLOY_DESIGN: {
    category: 'pattern',
    action: 'undeploy',
  },
  VALIDATE_DESIGN: {
    category: 'pattern',
    action: 'validate',
  },
  EVALUATE_DESIGN: {
    category: 'relationship',
    action: 'evaluation',
  },
  REGISTRANT_SUMMARY: {
    category: 'entity',
    action: 'get_summary',
  },
  ACADEMY_QUIZ_EVALUATION: {
    category: 'academy',
    action: 'academy_quiz_evaluation',
  },
};
