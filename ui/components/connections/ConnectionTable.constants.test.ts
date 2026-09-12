import { describe, expect, it } from 'vitest';
import {
  getTransitionDescription,
  toServerSortOrder,
  toUiSortOrder,
  type ConnectionTransitionMap,
} from './ConnectionTable.constants';

describe('getTransitionDescription', () => {
  const transitionMap: ConnectionTransitionMap = {
    connected: [
      { nextState: 'disconnected', description: 'Stops managing the cluster.' },
      { nextState: 'deleted' },
    ],
  };

  it('returns the definition-authored description for a known transition', () => {
    expect(getTransitionDescription(transitionMap, 'connected', 'disconnected')).toBe(
      'Stops managing the cluster.',
    );
  });

  it('normalizes status casing from callers that do not lowercase', () => {
    expect(getTransitionDescription(transitionMap, 'CONNECTED', 'DISCONNECTED')).toBe(
      'Stops managing the cluster.',
    );
  });

  it('returns undefined when the transition has no authored description', () => {
    expect(getTransitionDescription(transitionMap, 'connected', 'deleted')).toBeUndefined();
  });

  it('returns undefined for unknown states, missing maps, and unknown current status', () => {
    expect(getTransitionDescription(transitionMap, 'discovered', 'registered')).toBeUndefined();
    expect(getTransitionDescription(undefined, 'connected', 'disconnected')).toBeUndefined();
    expect(getTransitionDescription(transitionMap, undefined, 'disconnected')).toBeUndefined();
  });
});

describe('toServerSortOrder', () => {
  it('maps camelCase wire fields to the server DB sort columns', () => {
    expect(toServerSortOrder('createdAt desc')).toBe('created_at desc');
    expect(toServerSortOrder('createdAt asc')).toBe('created_at asc');
    expect(toServerSortOrder('updatedAt desc')).toBe('updated_at desc');
  });

  it('passes through fields the server already understands', () => {
    expect(toServerSortOrder('name asc')).toBe('name asc');
    // Older bookmarked URLs may still carry the snake_case form.
    expect(toServerSortOrder('created_at desc')).toBe('created_at desc');
  });

  it('defaults the direction to desc when missing (a bare column is dropped server-side)', () => {
    expect(toServerSortOrder('createdAt')).toBe('created_at desc');
  });

  it('normalizes extra whitespace between column and direction', () => {
    expect(toServerSortOrder('createdAt  desc')).toBe('created_at desc');
  });

  it('returns a valid default when the input is empty or whitespace-only', () => {
    expect(toServerSortOrder('')).toBe('created_at desc');
    expect(toServerSortOrder('   ')).toBe('created_at desc');
  });
});

describe('toUiSortOrder', () => {
  it('maps a bookmarked snake_case param back to the table column name', () => {
    expect(toUiSortOrder('created_at desc')).toBe('createdAt desc');
    expect(toUiSortOrder('updated_at asc')).toBe('updatedAt asc');
  });

  it('leaves an order that already uses column names untouched', () => {
    expect(toUiSortOrder('createdAt desc')).toBe('createdAt desc');
    expect(toUiSortOrder('name asc')).toBe('name asc');
  });

  it('defaults the direction and guards empty input', () => {
    expect(toUiSortOrder('created_at')).toBe('createdAt desc');
    expect(toUiSortOrder('')).toBe('createdAt desc');
    expect(toUiSortOrder('   ')).toBe('createdAt desc');
  });

  it('round-trips with toServerSortOrder so the sort survives a bookmarked URL', () => {
    // The legacy URL still drives the correct server column...
    expect(toServerSortOrder('created_at desc')).toBe('created_at desc');
    // ...and resolves to a real column for the active-sort indicator.
    expect(toUiSortOrder('created_at desc')).toBe('createdAt desc');
  });
});
