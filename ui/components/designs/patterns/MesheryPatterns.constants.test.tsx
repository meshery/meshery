import { describe, expect, it, vi } from 'vitest';
import {
  ACTION_TYPES,
  genericClickHandler,
  resetSelectedPattern,
} from './MesheryPatterns.constants';

describe('ACTION_TYPES', () => {
  it('exposes the design CRUD action descriptors with name and error_msg', () => {
    const keys = [
      'FETCH_PATTERNS',
      'UPDATE_PATTERN',
      'DELETE_PATTERN',
      'DEPLOY_PATTERN',
      'UNDEPLOY_PATTERN',
      'UPLOAD_PATTERN',
      'CLONE_PATTERN',
      'PUBLISH_CATALOG',
      'UNPUBLISH_CATALOG',
      'SCHEMA_FETCH',
    ];

    keys.forEach((key) => {
      const entry = (ACTION_TYPES as Record<string, { name: string; error_msg: string }>)[key];
      expect(entry).toBeDefined();
      expect(entry.name).toBe(key);
      expect(typeof entry.error_msg).toBe('string');
      expect(entry.error_msg.length).toBeGreaterThan(0);
    });
  });
});

describe('genericClickHandler', () => {
  it('stops propagation and invokes the callback with the event', () => {
    const stopPropagation = vi.fn();
    const cb = vi.fn();
    const ev = { stopPropagation } as any;

    genericClickHandler(ev, cb);

    expect(stopPropagation).toHaveBeenCalledTimes(1);
    expect(cb).toHaveBeenCalledWith(ev);
  });
});

describe('resetSelectedPattern', () => {
  it('returns the canonical initial state for selectedPattern', () => {
    expect(resetSelectedPattern()).toEqual({ show: false, pattern: null });
  });
});
