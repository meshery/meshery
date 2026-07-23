import { renderHook, act } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useTableUrlState } from './useTableUrlState';

const routerReplace = vi.fn();
const router = {
  query: {} as Record<string, unknown>,
  pathname: '/configuration/filters',
  replace: (...args: unknown[]) => routerReplace(...args),
};

vi.mock('next/router', () => ({
  useRouter: () => router,
}));

vi.mock('./useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

describe('useTableUrlState', () => {
  beforeEach(() => {
    router.query = {};
    routerReplace.mockReset();
    window.history.replaceState({}, '', '/configuration/filters');
  });

  it('preserves filter params when a stale router.query races a later update', () => {
    const { result } = renderHook(() =>
      useTableUrlState({
        tableKey: 'fil',
        defaults: {
          page: 0,
          pageSize: 10,
          sortOrder: '',
          search: '',
          filters: { vis: '' },
        },
      }),
    );

    act(() => {
      result.current.updateTableState({
        filters: { vis: 'public' },
        page: 0,
      });
    });

    window.history.replaceState({}, '', '/configuration/filters?fil_vis=public');

    act(() => {
      result.current.updateTableState({ search: '', page: 0 });
    });

    expect(routerReplace).toHaveBeenLastCalledWith(
      expect.objectContaining({
        query: expect.objectContaining({ fil_vis: 'public' }),
      }),
      undefined,
      { shallow: true },
    );
  });
});
