import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';

// `usePageTitle` dispatches the redux `updatePage` action wrapping the current
// path + the title we pass in. Wire mocks before importing the hook so module
// state binds to our spies, not the real store / window.location bits.

const dispatchMock = vi.fn();
const updatePageMock = vi.fn((payload) => ({ type: 'core/updatePage', payload }));
const getPathMock = vi.fn(() => '/dashboard');

vi.mock('react-redux', () => ({
  useDispatch: () => dispatchMock,
}));
vi.mock('@/store/slices/mesheryUi', () => ({
  updatePage: (payload: unknown) => updatePageMock(payload),
}));
vi.mock('../../../lib/path', () => ({
  getPath: () => getPathMock(),
}));

import { usePageTitle } from '../usePageTitle';

describe('usePageTitle', () => {
  beforeEach(() => {
    dispatchMock.mockClear();
    updatePageMock.mockClear();
    getPathMock.mockClear();
    getPathMock.mockReturnValue('/dashboard');
  });

  it('dispatches the page update on mount with the current path and supplied title', () => {
    renderHook(() => usePageTitle('Dashboard'));

    expect(getPathMock).toHaveBeenCalled();
    expect(updatePageMock).toHaveBeenCalledWith({ path: '/dashboard', title: 'Dashboard' });
    expect(dispatchMock).toHaveBeenCalledWith({
      type: 'core/updatePage',
      payload: { path: '/dashboard', title: 'Dashboard' },
    });
  });

  it('dispatches again when the title changes', () => {
    const { rerender } = renderHook(({ t }) => usePageTitle(t), {
      initialProps: { t: 'First' },
    });
    expect(dispatchMock).toHaveBeenCalledTimes(1);

    rerender({ t: 'Second' });

    expect(dispatchMock).toHaveBeenCalledTimes(2);
    expect(updatePageMock).toHaveBeenLastCalledWith({ path: '/dashboard', title: 'Second' });
  });

  it('does not dispatch on title-stable re-renders', () => {
    const { rerender } = renderHook(({ t }) => usePageTitle(t), {
      initialProps: { t: 'Stable' },
    });
    rerender({ t: 'Stable' });
    expect(dispatchMock).toHaveBeenCalledTimes(1);
  });

  it('reads the path freshly on each dispatching effect', () => {
    getPathMock.mockReturnValueOnce('/first');
    const { rerender } = renderHook(({ t }) => usePageTitle(t), {
      initialProps: { t: 'A' },
    });
    expect(updatePageMock).toHaveBeenLastCalledWith({ path: '/first', title: 'A' });

    getPathMock.mockReturnValueOnce('/second');
    act(() => {
      rerender({ t: 'B' });
    });
    expect(updatePageMock).toHaveBeenLastCalledWith({ path: '/second', title: 'B' });
  });
});
