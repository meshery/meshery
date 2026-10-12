// The Notification Center header used to read a React component straight out of
// Redux (`uiConfig.icon || BellIcon`), which is what put a non-serializable
// value at `events.ui.icon` (issue #16873). The store now carries a name and the
// header resolves it here, so this is the half that has to stay total: every
// input, including the component an older extension might still send, must
// produce a renderable component rather than undefined.

import { describe, expect, it } from 'vitest';
import BellIcon from '../../../../assets/icons/BellIcon';
import AlertIcon from '../../../../assets/icons/AlertIcon';
import {
  DEFAULT_NOTIFICATION_ICON,
  NOTIFICATION_ICONS,
  resolveNotificationIcon,
} from '../notificationIcons';

describe('resolveNotificationIcon', () => {
  it('resolves a registered name to its component', () => {
    expect(resolveNotificationIcon('bell')).toBe(BellIcon);
    expect(resolveNotificationIcon('alert')).toBe(AlertIcon);
  });

  it('falls back to the bell for an unknown name', () => {
    expect(resolveNotificationIcon('no-such-icon')).toBe(BellIcon);
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

  it('never returns undefined, so the header always has something to render', () => {
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

  it('registers only serializable names, so any of them can live in Redux', () => {
    Object.keys(NOTIFICATION_ICONS).forEach((name) => {
      expect(typeof name).toBe('string');
    });
  });
});
