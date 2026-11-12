import { describe, expect, it } from 'vitest';
import { NOTIFICATION_STATUS, EVENT_TYPES, SERVER_EVENT_TYPES } from '../event-types';

describe('NOTIFICATION_STATUS', () => {
  it('exposes VIEWED with value "viewed"', () => {
    expect(NOTIFICATION_STATUS.VIEWED).toBe('viewed');
  });

  it('exposes NEW with value "new"', () => {
    expect(NOTIFICATION_STATUS.NEW).toBe('new');
  });

  it('has exactly the expected keys', () => {
    expect(Object.keys(NOTIFICATION_STATUS).sort()).toEqual(['NEW', 'VIEWED']);
  });
});

describe('EVENT_TYPES', () => {
  it('exposes SUCCESS with type "success"', () => {
    expect(EVENT_TYPES.SUCCESS).toEqual({ type: 'success' });
  });

  it('exposes DEFAULT with type "default"', () => {
    expect(EVENT_TYPES.DEFAULT).toEqual({ type: 'default' });
  });

  it('exposes INFO with type "info"', () => {
    expect(EVENT_TYPES.INFO).toEqual({ type: 'info' });
  });

  it('exposes WARNING with type "warning"', () => {
    expect(EVENT_TYPES.WARNING).toEqual({ type: 'warning' });
  });

  it('exposes ERROR with type "error"', () => {
    expect(EVENT_TYPES.ERROR).toEqual({ type: 'error' });
  });

  it('has exactly the expected keys', () => {
    expect(Object.keys(EVENT_TYPES).sort()).toEqual([
      'DEFAULT',
      'ERROR',
      'INFO',
      'SUCCESS',
      'WARNING',
    ]);
  });
});

describe('SERVER_EVENT_TYPES', () => {
  it('maps 0 to SUCCESS', () => {
    expect(SERVER_EVENT_TYPES[0]).toBe(EVENT_TYPES.SUCCESS);
    expect(SERVER_EVENT_TYPES[0]).toEqual({ type: 'success' });
  });

  it('maps 1 to WARNING', () => {
    expect(SERVER_EVENT_TYPES[1]).toBe(EVENT_TYPES.WARNING);
    expect(SERVER_EVENT_TYPES[1]).toEqual({ type: 'warning' });
  });

  it('maps 2 to ERROR', () => {
    expect(SERVER_EVENT_TYPES[2]).toBe(EVENT_TYPES.ERROR);
    expect(SERVER_EVENT_TYPES[2]).toEqual({ type: 'error' });
  });

  it('returns undefined for unknown numeric codes', () => {
    // @ts-expect-error – intentionally testing absent key
    expect(SERVER_EVENT_TYPES[3]).toBeUndefined();
    // @ts-expect-error – intentionally testing absent key
    expect(SERVER_EVENT_TYPES[99]).toBeUndefined();
  });

  it('has exactly three numeric entries', () => {
    expect(Object.keys(SERVER_EVENT_TYPES).sort()).toEqual(['0', '1', '2']);
  });
});
