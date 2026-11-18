import { beforeEach, describe, expect, it, vi } from 'vitest';
import { renderHook, act } from '@testing-library/react';

const { listConnectionDefinitions, lazyTrigger, lazyResult, lazyLastInfo } = vi.hoisted(() => ({
  listConnectionDefinitions: vi.fn(),
  lazyTrigger: vi.fn(),
  lazyResult: { isFetching: false },
  lazyLastInfo: { lastArg: null },
}));

vi.mock('@meshery/schemas/mesheryApi', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@meshery/schemas/mesheryApi')>()),
  useListConnectionDefinitionsQuery: listConnectionDefinitions,
  useLazyListConnectionDefinitionsQuery: () => [lazyTrigger, lazyResult, lazyLastInfo],
}));

import {
  useGetConnectionDefinitionsQuery,
  useLazyGetConnectionDefinitionsQuery,
} from '../meshModel';

describe('useGetConnectionDefinitionsQuery wrapper', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('defaults pageSize to "all" when called with no arguments', () => {
    renderHook(() => useGetConnectionDefinitionsQuery());

    expect(listConnectionDefinitions).toHaveBeenCalledWith(
      expect.objectContaining({
        pageSize: 'all',
        page: 0,
        search: '',
        trim: false,
      }),
      undefined,
    );
  });

  it('forwards custom params while maintaining pageSize: "all" default', () => {
    const options = { skip: true };
    renderHook(() =>
      useGetConnectionDefinitionsQuery({ params: { search: 'kubernetes', page: 2 } }, options),
    );

    expect(listConnectionDefinitions).toHaveBeenCalledWith(
      expect.objectContaining({
        search: 'kubernetes',
        page: 2,
        pageSize: 'all',
        trim: false,
      }),
      options,
    );
  });

  it('maps legacy pagesize parameter to canonical pageSize wire key', () => {
    renderHook(() => useGetConnectionDefinitionsQuery({ params: { pagesize: 25 } }));

    const callArg = listConnectionDefinitions.mock.calls[0]?.[0];
    expect(callArg).toEqual(
      expect.objectContaining({
        pageSize: 25,
        page: 0,
        search: '',
        trim: false,
      }),
    );
    expect(callArg).not.toHaveProperty('pagesize');
  });
});

describe('useLazyGetConnectionDefinitionsQuery wrapper', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('maps arguments through listArg with pageSize: "all" on lazy trigger invocation', () => {
    const { result } = renderHook(() => useLazyGetConnectionDefinitionsQuery());
    const [trigger, queryResult, lastInfo] = result.current;

    expect(queryResult).toBe(lazyResult);
    expect(lastInfo).toBe(lazyLastInfo);

    act(() => {
      trigger({ params: { search: 'prometheus' } }, true);
    });

    expect(lazyTrigger).toHaveBeenCalledWith(
      expect.objectContaining({
        search: 'prometheus',
        page: 0,
        pageSize: 'all',
        trim: false,
      }),
      true,
    );
  });
});
