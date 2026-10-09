import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  InfoIcon: () => null,
  notificationColors: {
    info: { main: '#info' },
    error: { main: '#error', dark: '#error_dark' },
    warning: { main: '#warning', light: '#warning_light' },
    success: { main: '#success' },
  },
}));

vi.mock('../../../assets/icons/AlertIcon', () => ({
  default: () => null,
}));
vi.mock('../../../assets/icons/ErrorIcon', () => ({
  default: () => null,
}));
vi.mock('../../../assets/icons/ReadIcon', () => ({
  default: () => null,
}));

import {
  SEVERITY,
  STATUS,
  NOTIFICATION_CENTER_TOGGLE_CLASS,
  SEVERITY_TO_NOTIFICATION_TYPE_MAPPING,
  SEVERITY_STYLE,
  EVENT_TYPE,
  eventDetailFormatterKey,
  getStatusStyle,
} from './constants';
import { notificationColors } from '@sistent/sistent';

describe('NotificationCenter constants', () => {
  it('defines severity levels', () => {
    expect(SEVERITY.INFO).toBe('informational');
    expect(SEVERITY.ERROR).toBe('error');
    expect(SEVERITY.WARNING).toBe('warning');
    expect(SEVERITY.SUCCESS).toBe('success');
  });

  it('defines read/unread statuses', () => {
    expect(STATUS.READ).toBe('read');
    expect(STATUS.UNREAD).toBe('unread');
  });

  it('exposes the notification center toggle class name', () => {
    expect(NOTIFICATION_CENTER_TOGGLE_CLASS).toBe('toggle-notification-center');
  });

  it('maps severities to notification type strings', () => {
    expect(SEVERITY_TO_NOTIFICATION_TYPE_MAPPING[SEVERITY.INFO]).toBe('info');
    expect(SEVERITY_TO_NOTIFICATION_TYPE_MAPPING[SEVERITY.ERROR]).toBe('error');
    expect(SEVERITY_TO_NOTIFICATION_TYPE_MAPPING[SEVERITY.WARNING]).toBe('warning');
    expect(SEVERITY_TO_NOTIFICATION_TYPE_MAPPING[SEVERITY.SUCCESS]).toBe('success');
  });

  it('provides severity styles that resolve to palette colors', () => {
    expect(SEVERITY_STYLE[SEVERITY.INFO].color).toBe(notificationColors.info.main);
    expect(SEVERITY_STYLE[SEVERITY.ERROR].color).toBe(notificationColors.error.main);
    expect(SEVERITY_STYLE[SEVERITY.ERROR].darkColor).toBe(notificationColors.error.dark);
    expect(SEVERITY_STYLE[SEVERITY.WARNING].color).toBe(notificationColors.warning.main);
    expect(SEVERITY_STYLE[SEVERITY.SUCCESS].color).toBe(notificationColors.success.main);
    Object.values(SEVERITY_STYLE).forEach((style) => {
      expect(style.icon).toBeDefined();
    });
  });

  it('describes well-known event types', () => {
    expect(EVENT_TYPE.CatalogManagementDeployDesign).toEqual({
      category: 'pattern',
      action: 'deploy',
    });
    expect(EVENT_TYPE.CatalogManagementUndeployDesign).toEqual({
      category: 'pattern',
      action: 'undeploy',
    });
    expect(EVENT_TYPE.CatalogManagementValidateDesign).toEqual({
      category: 'pattern',
      action: 'validate',
    });
    expect(EVENT_TYPE.EVALUATE_DESIGN).toEqual({ category: 'relationship', action: 'evaluation' });
    expect(EVENT_TYPE.REGISTRANT_SUMMARY).toEqual({ category: 'entity', action: 'get_summary' });
    expect(EVENT_TYPE.ACADEMY_QUIZ_EVALUATION).toEqual({
      category: 'academy',
      action: 'academy_quiz_evaluation',
    });
  });

  it('joins action and category into a stable formatter key', () => {
    expect(eventDetailFormatterKey({ action: 'foo', category: 'bar' })).toBe('foo-bar');
    expect(eventDetailFormatterKey(EVENT_TYPE.CatalogManagementDeployDesign)).toBe(
      'deploy-pattern',
    );
  });

  it('produces a theme-aware status style map', () => {
    const fakeTheme = {
      palette: { text: { primary: '#fff' } },
    } as any;
    const styles = getStatusStyle(fakeTheme);
    expect(styles[STATUS.READ].color).toBe('#fff');
    expect(styles[STATUS.READ].darkColor).toBe('#fff');
    expect(styles[STATUS.READ].icon).toBeDefined();
  });
});
