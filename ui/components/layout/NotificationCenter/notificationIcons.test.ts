import { describe, expect, it } from 'vitest';
import AlertIcon from '../../../assets/icons/AlertIcon';
import BellIcon from '../../../assets/icons/BellIcon';
import {
  DEFAULT_NOTIFICATION_ICON,
  NOTIFICATION_ICONS,
  resolveNotificationIcon,
} from './notificationIcons';

describe('resolveNotificationIcon', () => {
  it('resolves a registered name to its component', () => {
    expect(resolveNotificationIcon('bell')).toBe(BellIcon);
    expect(resolveNotificationIcon('alert')).toBe(AlertIcon);
  });

  it('falls back to the bell for an unknown name', () => {
    expect(resolveNotificationIcon('no-such-icon')).toBe(BellIcon);
  });

  it('falls back to the bell for prototype chain names', () => {
    expect(resolveNotificationIcon('constructor')).toBe(BellIcon);
    expect(resolveNotificationIcon('__proto__')).toBe(BellIcon);
    expect(resolveNotificationIcon('toString')).toBe(BellIcon);
  });

  it('falls back to the bell when no icon is set', () => {
    expect(resolveNotificationIcon(undefined)).toBe(BellIcon);
    expect(resolveNotificationIcon(null)).toBe(BellIcon);
    expect(resolveNotificationIcon('')).toBe(BellIcon);
  });

  it('falls back to the bell for a component, which the store no longer keeps', () => {
    const LegacyIcon = () => null;
    expect(resolveNotificationIcon(LegacyIcon)).toBe(BellIcon);
  });

  it('never returns undefined', () => {
    const inputs = [
      'bell',
      'alert',
      'error',
      'info',
      'nope',
      '',
      null,
      undefined,
      42,
      {},
      () => null,
    ];
    inputs.forEach((input) => {
      expect(resolveNotificationIcon(input)).toBeTypeOf('function');
    });
  });

  it('exposes a default that is itself a registered name', () => {
    expect(NOTIFICATION_ICONS[DEFAULT_NOTIFICATION_ICON]).toBe(BellIcon);
  });
});
