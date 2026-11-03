import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useTableUrlState } from '../useTableUrlState';

const router = {
  query: {} as Record<string, string | string[] | undefined>,
  pathname: '/management/connections',
  replace: vi.fn(),
};

vi.mock('next/router', () => ({
  useRouter: () => router,
}));

vi.mock('../useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

// Next.js exposes `router.query` values already URL-decoded, so the hook must
// consume them verbatim rather than decoding them a second time (#22148).
describe('useTableUrlState', () => {
  beforeEach(() => {
    router.query = {};
  });

  it('does not throw on a search containing a bare "%"', () => {
    router.query = { con_q: '100%' };

    const { result } = renderHook(() => useTableUrlState({ tableKey: 'con', defaults: {} }));

    expect(result.current.tableState.search).toBe('100%');
  });

  it('does not decode a valid-looking escape sequence a second time', () => {
    router.query = { con_q: '50%25' };

    const { result } = renderHook(() => useTableUrlState({ tableKey: 'con', defaults: {} }));

    expect(result.current.tableState.search).toBe('50%25');
  });

  it('preserves "%" in sort and filter params verbatim', () => {
    router.query = { con_sort: 'name%20asc', con_status: '100%' };

    const { result } = renderHook(() =>
      useTableUrlState({ tableKey: 'con', defaults: { filters: { status: '' } } }),
    );

    expect(result.current.tableState.sortOrder).toBe('name%20asc');
    expect(result.current.tableState.filters.status).toBe('100%');
  });
});
