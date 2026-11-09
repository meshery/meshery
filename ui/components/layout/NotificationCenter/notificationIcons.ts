import AlertIcon from '../../../assets/icons/AlertIcon';
import BellIcon from '../../../assets/icons/BellIcon';
import ErrorIcon from '../../../assets/icons/ErrorIcon';
import InfoIcon from '../../../assets/icons/InfoIcon';

export const NOTIFICATION_ICONS = {
  bell: BellIcon,
  alert: AlertIcon,
  error: ErrorIcon,
  info: InfoIcon,
};

export const DEFAULT_NOTIFICATION_ICON = 'bell';

export const resolveNotificationIcon = (name) => {
  if (typeof name !== 'string') {
    return NOTIFICATION_ICONS[DEFAULT_NOTIFICATION_ICON];
  }
  return Object.hasOwn(NOTIFICATION_ICONS, name)
    ? NOTIFICATION_ICONS[name]
    : NOTIFICATION_ICONS[DEFAULT_NOTIFICATION_ICON];
};
