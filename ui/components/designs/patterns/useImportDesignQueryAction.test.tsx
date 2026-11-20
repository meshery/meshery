import React, { StrictMode } from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useImportDesignQueryAction } from './useImportDesignQueryAction';

interface MockRouter {
  isReady: boolean;
  pathname: string;
  query: Record<string, string>;
  replace: ReturnType<typeof vi.fn>;
  push: ReturnType<typeof vi.fn>;
}

let mockRouter: MockRouter;

vi.mock('next/router', () => ({
  useRouter: () => mockRouter,
}));

describe('useImportDesignQueryAction', () => {
  beforeEach(() => {
    mockRouter = {
      isReady: true,
      pathname: '/configuration/designs',
      query: {},
      replace: vi.fn(),
      push: vi.fn(),
    };
  });

  it('opens modal and strips action when permitted', () => {
    mockRouter.query = { action: 'import' };
    const openImportModal = vi.fn();

    renderHook(() =>
      useImportDesignQueryAction({
        canImport: true,
        openImportModal,
      }),
    );

    expect(openImportModal).toHaveBeenCalledTimes(1);
    expect(mockRouter.replace).toHaveBeenCalledTimes(1);
    expect(mockRouter.replace).toHaveBeenCalledWith(
      {
        pathname: '/configuration/designs',
        query: {},
      },
      undefined,
      { shallow: true },
    );
  });

  it('does NOT open when permission denied but still strips action', () => {
    mockRouter.query = { action: 'import' };
    const openImportModal = vi.fn();

    renderHook(() =>
      useImportDesignQueryAction({
        canImport: false,
        openImportModal,
      }),
    );

    expect(openImportModal).not.toHaveBeenCalled();
    expect(mockRouter.replace).toHaveBeenCalledTimes(1);
    expect(mockRouter.replace).toHaveBeenCalledWith(
      {
        pathname: '/configuration/designs',
        query: {},
      },
      undefined,
      { shallow: true },
    );
  });

  it('preserves other query params (e.g. view=grid)', () => {
    mockRouter.query = { action: 'import', view: 'grid', search: 'k8s' };
    const openImportModal = vi.fn();

    renderHook(() =>
      useImportDesignQueryAction({
        canImport: true,
        openImportModal,
      }),
    );

    expect(openImportModal).toHaveBeenCalledTimes(1);
    expect(mockRouter.replace).toHaveBeenCalledTimes(1);
    expect(mockRouter.replace).toHaveBeenCalledWith(
      {
        pathname: '/configuration/designs',
        query: { view: 'grid', search: 'k8s' },
      },
      undefined,
      { shallow: true },
    );
  });

  it('does nothing without action', () => {
    mockRouter.query = { view: 'table' };
    const openImportModal = vi.fn();

    renderHook(() =>
      useImportDesignQueryAction({
        canImport: true,
        openImportModal,
      }),
    );

    expect(openImportModal).not.toHaveBeenCalled();
    expect(mockRouter.replace).not.toHaveBeenCalled();
  });

  it('does nothing before router.isReady', () => {
    mockRouter.isReady = false;
    mockRouter.query = { action: 'import' };
    const openImportModal = vi.fn();

    renderHook(() =>
      useImportDesignQueryAction({
        canImport: true,
        openImportModal,
      }),
    );

    expect(openImportModal).not.toHaveBeenCalled();
    expect(mockRouter.replace).not.toHaveBeenCalled();
  });

  it('fires only once under <React.StrictMode>', () => {
    mockRouter.query = { action: 'import' };
    const openImportModal = vi.fn();

    renderHook(
      () =>
        useImportDesignQueryAction({
          canImport: true,
          openImportModal,
        }),
      {
        wrapper: ({ children }: { children: React.ReactNode }) => (
          <StrictMode>{children}</StrictMode>
        ),
      },
    );

    expect(openImportModal).toHaveBeenCalledTimes(1);
    expect(mockRouter.replace).toHaveBeenCalledTimes(1);
    expect(mockRouter.replace).toHaveBeenCalledWith(
      {
        pathname: '/configuration/designs',
        query: {},
      },
      undefined,
      { shallow: true },
    );
  });

  it('does not fire for action=other', () => {
    mockRouter.query = { action: 'export' };
    const openImportModal = vi.fn();

    renderHook(() =>
      useImportDesignQueryAction({
        canImport: true,
        openImportModal,
      }),
    );

    expect(openImportModal).not.toHaveBeenCalled();
    expect(mockRouter.replace).not.toHaveBeenCalled();
  });

  it('strips action when a competing effect re-merges stale query with default view', () => {
    mockRouter.query = { action: 'import' };
    const openImportModal = vi.fn();

    const { rerender } = renderHook(() =>
      useImportDesignQueryAction({
        canImport: true,
        openImportModal,
      }),
    );

    expect(openImportModal).toHaveBeenCalledTimes(1);

    // Competing effect ran with a stale snapshot and merged action=import back in with view=grid
    mockRouter.query = { action: 'import', view: 'grid' };
    rerender();

    // Modal must NOT be opened a second time
    expect(openImportModal).toHaveBeenCalledTimes(1);

    // But action must be stripped from the updated query, preserving view=grid
    expect(mockRouter.replace).toHaveBeenLastCalledWith(
      {
        pathname: '/configuration/designs',
        query: { view: 'grid' },
      },
      undefined,
      { shallow: true },
    );
  });
});
