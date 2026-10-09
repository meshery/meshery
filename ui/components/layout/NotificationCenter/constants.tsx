import AlertIcon from '../../../assets/icons/AlertIcon';
import ErrorIcon from '../../../assets/icons/ErrorIcon';
import ReadIcon from '../../../assets/icons/ReadIcon';
import { InfoIcon, notificationColors } from '@sistent/sistent';
import type { Theme } from '@sistent/sistent';

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

// Theme-aware status styles. Resolves to live palette tokens so dark/light
// mode picks the correct value at render time instead of using hard-coded hex.
export const getStatusStyle = (theme: Theme) => ({
  [STATUS.READ]: {
    icon: ReadIcon,
    color: theme.palette.text.primary,
    darkColor: theme.palette.text.primary,
  },
});

export const SEVERITY_STYLE = {
  [SEVERITY.INFO]: {
    icon: InfoIcon,
    color: notificationColors.info.main,
    darkColor: notificationColors.info.main,
  },
  [SEVERITY.ERROR]: {
    icon: ErrorIcon,
    color: notificationColors.error.main,
    darkColor: notificationColors.error.dark,
  },
  [SEVERITY.WARNING]: {
    icon: AlertIcon,
    color: notificationColors.warning.main,
    darkColor: notificationColors.warning.main,
  },
  [SEVERITY.SUCCESS]: {
    icon: InfoIcon,
    color: notificationColors.success.main,
    darkColor: notificationColors.success.main,
  },
};

export const eventDetailFormatterKey = ({
  action,
  category,
}: {
  action: string;
  category: string;
}): string => `${action}-${category}`;

export const EVENT_TYPE = {
  CatalogManagementDeployDesign: {
    category: 'pattern',
    action: 'deploy',
  },
  CatalogManagementUndeployDesign: {
    category: 'pattern',
    action: 'undeploy',
  },
  CatalogManagementValidateDesign: {
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
