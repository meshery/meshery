import AlertIcon from '../../../assets/icons/AlertIcon';
import BellIcon from '../../../assets/icons/BellIcon';
import ErrorIcon from '../../../assets/icons/ErrorIcon';
import InfoIcon from '../../../assets/icons/InfoIcon';

/**
 * Icons the Notification Center header can be opened with, keyed by a
 * serializable name.
 *
 * `openNotificationCenter` is dispatched into the Meshery store by extensions
 * over the event bus (see `store/index.ts`), so whatever it carries is held in
 * Redux state. A React component is not serializable, so the store carries the
 * name and the header resolves it to a component at render time.
 *
 * Extensions that want an icon here pass one of these names. Adding a new one
 * is a one-line entry in this map.
 */
export const NOTIFICATION_ICONS = {
  bell: BellIcon,
  alert: AlertIcon,
  error: ErrorIcon,
  info: InfoIcon,
};

export const DEFAULT_NOTIFICATION_ICON = 'bell';

/**
 * Resolves a stored icon name to its component, falling back to the bell for
 * an unknown name, a missing one, or anything that is not a string - including
 * a component passed by an older extension, which no longer reaches the store.
 */
export const resolveNotificationIcon = (name) => {
  if (typeof name !== 'string') return NOTIFICATION_ICONS[DEFAULT_NOTIFICATION_ICON];
  return NOTIFICATION_ICONS[name] || NOTIFICATION_ICONS[DEFAULT_NOTIFICATION_ICON];
};
