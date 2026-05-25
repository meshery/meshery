import React from 'react';
import { renderHook, act, render } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';

const isInOperatorMode = vi.fn();
const JsonParse = vi.fn((s: string) => JSON.parse(s));
const notify = vi.fn();
const deleteView = vi.fn();
const deletePatternFile = vi.fn();
const downloadContent = vi.fn();
const downloadFileFromContent = vi.fn();
const updateProgress = vi.fn();

vi.mock('@/utils/Enum', () => ({
  APP_MODE: { DESIGN: 'design', OPERATOR: 'operator' },
  RESOURCE_TYPE: { DESIGN: 'design', VIEW: 'view', CATALOG: 'catalog' },
}));

vi.mock('@/utils/utils', () => ({
  isInOperatorMode: () => isInOperatorMode(),
  JsonParse: (s: string) => JsonParse(s),
}));

vi.mock('@/rtk-query/design', () => ({
  useDeletePatternFileMutation: () => [deletePatternFile],
}));

vi.mock('@/rtk-query/view', () => ({
  useDeleteViewMutation: () => [deleteView],
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('@sistent/sistent', () => ({
  DesignIcon: ({ ...props }: any) => <svg data-testid="design-icon" {...props} />,
  ViewIcon: ({ ...props }: any) => <svg data-testid="view-icon" {...props} />,
  PROMPT_VARIANTS: { DANGER: 'danger' },
  useTheme: () => ({ palette: { icon: { brand: 'brand' } } }),
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: { SUCCESS: 'SUCCESS', ERROR: 'ERROR', INFO: 'INFO' },
}));

vi.mock('css/icons.styles', () => ({ iconMedium: {} }));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: (...args: unknown[]) => updateProgress(...args),
}));

vi.mock('@/utils/fileDownloader', () => ({
  default: (...args: unknown[]) => downloadContent(...args),
  downloadFileFromContent: (...args: unknown[]) => downloadFileFromContent(...args),
}));

import useInfiniteScroll, {
  getDesignPath,
  viewPath,
  catalogPath,
  getShareableResourceRoute,
  getModelNamesBasedOnDisplayNames,
  handleUpdatePatternVisibility,
  handleUpdateViewVisibility,
  getDefaultFilterType,
  useGetIconBasedOnMode,
  useContentDelete,
  useContentDownload,
} from './hooks';

beforeEach(() => {
  isInOperatorMode.mockReset();
  JsonParse.mockClear();
  notify.mockReset();
  deletePatternFile.mockReset();
  deleteView.mockReset();
  downloadContent.mockReset();
  downloadFileFromContent.mockReset();
  updateProgress.mockReset();
  Object.defineProperty(window, 'location', {
    value: new URL('https://example.com/foo'),
    writable: true,
  });
});

describe('URL path helpers', () => {
  it('getDesignPath returns a URL with mode=design and design id', () => {
    const url = getDesignPath('d-1');
    expect(url).toContain('mode=design');
    expect(url).toContain('design=d-1');
  });

  it('viewPath returns a URL with mode=operator and view type/id/name', () => {
    const url = viewPath({ id: 'v-1', name: 'My View' });
    expect(url).toContain('mode=operator');
    expect(url).toContain('type=view');
    expect(url).toContain('id=v-1');
    expect(url).toContain('name=My+View');
  });

  it('catalogPath returns a URL with mode=design and catalog type/id/name', () => {
    const url = catalogPath({ id: 'c-1', name: 'A' });
    expect(url).toContain('type=catalog');
    expect(url).toContain('id=c-1');
  });

  it('getShareableResourceRoute routes by type and throws for unknown types', () => {
    expect(getShareableResourceRoute('design', 'd-1', 'n')).toContain('mode=design');
    expect(getShareableResourceRoute('view', 'v-1', 'n')).toContain('mode=operator');
    expect(getShareableResourceRoute('catalog', 'c-1', 'n')).toContain('type=catalog');
    expect(() => getShareableResourceRoute('mystery', 'x', 'n')).toThrow(/Unknown resource type/);
  });
});

describe('getModelNamesBasedOnDisplayNames', () => {
  it('returns model names matching the provided display names (case-insensitive, de-duplicated)', () => {
    const models = [
      { displayName: 'Kubernetes', name: 'kubernetes' },
      { displayName: 'kubernetes', name: 'kubernetes-dup' },
      { displayName: 'Istio', name: 'istio' },
    ];
    const out = getModelNamesBasedOnDisplayNames(models, ['Kubernetes', 'ISTIO']);
    expect(out).toEqual(['kubernetes', 'istio']);
  });
});

describe('visibility update helpers', () => {
  it('handleUpdatePatternVisibility forwards id, name, file, visibility to updatePatterns', async () => {
    const updatePatterns = vi.fn().mockResolvedValue({});
    const selectedResource = {
      id: 'p-1',
      name: 'design',
      catalogData: { foo: 'bar' },
      patternFile: JSON.stringify({ k: 'v' }),
    };
    await handleUpdatePatternVisibility({
      value: 'public',
      updatePatterns,
      selectedResource,
    });
    expect(updatePatterns).toHaveBeenCalledWith({
      updateBody: {
        id: 'p-1',
        name: 'design',
        catalogData: { foo: 'bar' },
        designFile: { k: 'v' },
        visibility: 'public',
      },
    });
  });

  it('handleUpdateViewVisibility calls updateView with id and visibility body', async () => {
    const updateView = vi.fn().mockResolvedValue({});
    await handleUpdateViewVisibility({
      value: 'private',
      updateView,
      selectedResource: { id: 'v-1' },
    });
    expect(updateView).toHaveBeenCalledWith({
      id: 'v-1',
      body: { visibility: 'private' },
    });
  });

  it('propagates errors via the returned error field', async () => {
    const updateView = vi.fn().mockResolvedValue({ error: { error: 'oops' } });
    const out = await handleUpdateViewVisibility({
      value: 'private',
      updateView,
      selectedResource: { id: 'v-1' },
    });
    expect(out.error).toBe('oops');
  });
});

describe('getDefaultFilterType', () => {
  it('returns VIEW when in operator mode', () => {
    isInOperatorMode.mockReturnValue(true);
    expect(getDefaultFilterType()).toBe('view');
  });

  it('returns DESIGN otherwise', () => {
    isInOperatorMode.mockReturnValue(false);
    expect(getDefaultFilterType()).toBe('design');
  });
});

describe('useGetIconBasedOnMode', () => {
  it('renders DesignIcon when mode is DESIGN', () => {
    const Component = () => useGetIconBasedOnMode({ mode: 'design' });
    const { container } = render(<Component />);
    expect(container.querySelector('[data-testid="design-icon"]')).toBeInTheDocument();
  });

  it('renders ViewIcon when mode is VIEW', () => {
    const Component = () => useGetIconBasedOnMode({ mode: 'view' });
    const { container } = render(<Component />);
    expect(container.querySelector('[data-testid="view-icon"]')).toBeInTheDocument();
  });

  it('returns undefined for an unknown mode', () => {
    const Component = () => useGetIconBasedOnMode({ mode: 'banana' });
    const { container } = render(<Component />);
    expect(container.innerHTML).toBe('');
  });
});

describe('useContentDelete', () => {
  it('does NOT delete when modal returns something other than "DELETE"', async () => {
    const modalRef = { current: { show: vi.fn().mockResolvedValue('Cancel') } };
    const { result } = renderHook(() => useContentDelete(modalRef));

    await act(async () => {
      await result.current.handleDelete([{ id: 'a', name: 'a' }], 'design', vi.fn());
    });

    expect(deletePatternFile).not.toHaveBeenCalled();
  });

  it('deletes designs (calls deletePatternFile per item) and notifies on success', async () => {
    deletePatternFile.mockReturnValue({ unwrap: () => Promise.resolve({}) });
    const modalRef = { current: { show: vi.fn().mockResolvedValue('DELETE') } };
    const refetch = vi.fn();

    const { result } = renderHook(() => useContentDelete(modalRef));
    await act(async () => {
      await result.current.handleDelete(
        [
          { id: 'a', name: 'one' },
          { id: 'b', name: 'two' },
        ],
        'design',
        refetch,
      );
    });

    expect(deletePatternFile).toHaveBeenCalledTimes(2);
    expect(notify).toHaveBeenCalledWith({
      message: '"one" Design deleted',
      event_type: 'SUCCESS',
    });
    expect(refetch).toHaveBeenCalled();
  });

  it('notifies an ERROR when a deletion mutation throws', async () => {
    deletePatternFile.mockReturnValue({ unwrap: () => Promise.reject(new Error('boom')) });
    const modalRef = { current: { show: vi.fn().mockResolvedValue('DELETE') } };

    const { result } = renderHook(() => useContentDelete(modalRef));
    await act(async () => {
      await result.current.handleDelete([{ id: 'a', name: 'one' }], 'design', vi.fn());
    });
    expect(notify).toHaveBeenCalledWith({
      message: 'Unable to delete "one" Design',
      event_type: 'ERROR',
    });
  });

  it('deletes views when type is VIEW', async () => {
    deleteView.mockReturnValue({ unwrap: () => Promise.resolve({}) });
    const modalRef = { current: { show: vi.fn().mockResolvedValue('DELETE') } };

    const { result } = renderHook(() => useContentDelete(modalRef));
    await act(async () => {
      await result.current.handleDelete([{ id: 'v-1', name: 'View1' }], 'view', vi.fn());
    });
    expect(deleteView).toHaveBeenCalledWith({ id: 'v-1' });
    expect(notify).toHaveBeenCalledWith({
      message: '"View1" View deleted',
      event_type: 'SUCCESS',
    });
  });
});

describe('useContentDownload', () => {
  it('handleDesignDownload calls downloadContent and notifies', () => {
    const { result } = renderHook(() => useContentDownload());
    const ev = { stopPropagation: vi.fn() } as any;
    result.current.handleDesignDownload(ev, { id: 'd-1', name: 'a' }, 'src', {});
    expect(downloadContent).toHaveBeenCalledWith(
      expect.objectContaining({ id: 'd-1', name: 'a', type: 'pattern' }),
    );
    expect(notify).toHaveBeenCalledWith({
      message: '"a" design downloaded',
      event_type: 'INFO',
    });
  });

  it('handleDesignDownload supports an array of designs', () => {
    const { result } = renderHook(() => useContentDownload());
    result.current.handleDesignDownload(
      { stopPropagation: vi.fn() } as any,
      [
        { id: 'd-1', name: 'a' },
        { id: 'd-2', name: 'b' },
      ],
      'src',
      {},
    );
    expect(downloadContent).toHaveBeenCalledTimes(2);
  });

  it('handleViewDownload writes each view to a JSON file and notifies', () => {
    const { result } = renderHook(() => useContentDownload());
    result.current.handleViewDownload([{ id: 'v-1', name: 'one' }]);
    expect(downloadFileFromContent).toHaveBeenCalled();
    expect(downloadFileFromContent.mock.calls[0][1]).toBe('one.json');
    expect(notify).toHaveBeenCalledWith({
      message: '"one" view downloaded',
      event_type: 'INFO',
    });
  });
});

describe('useInfiniteScroll', () => {
  beforeEach(() => {
    // jsdom doesn't ship IntersectionObserver; stub it for these tests.
    class MockIntersectionObserver {
      observe() {}
      unobserve() {}
      disconnect() {}
    }
    (globalThis as any).IntersectionObserver = MockIntersectionObserver;
  });

  it('returns a loadingRef and reacts to hasMore/isLoading without throwing', () => {
    const onLoadMore = vi.fn();
    const { result, rerender } = renderHook(
      ({ isLoading, hasMore }) => useInfiniteScroll({ isLoading, hasMore, onLoadMore }),
      { initialProps: { isLoading: true, hasMore: true } },
    );
    expect(result.current.loadingRef).toBeDefined();

    // Toggle states to exercise the hook's effect branches.
    rerender({ isLoading: false, hasMore: true });
    rerender({ isLoading: false, hasMore: false });
  });
});
