import { describe, expect, it, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';

const notify = vi.fn();
const handleSuccess = vi.fn();
const handleInfo = vi.fn();
const notifyApiError = vi.fn();
const deletePatternByIdMock = vi.fn();
const publishCatalogContentMock = vi.fn();
const unassignDesignFromWorkspaceMock = vi.fn();
const routerPushMock = vi.fn();

vi.mock('@/rtk-query/design', () => ({
  useDeletePatternMutation: () => [deletePatternByIdMock],
}));

vi.mock('@/rtk-query/filter', () => ({
  usePublishFilterMutation: () => [publishCatalogContentMock],
}));

vi.mock('@/rtk-query/workspace', () => ({
  useUnassignDesignFromWorkspaceMutation: () => [unassignDesignFromWorkspaceMock],
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
  useNotificationHandlers: () => ({
    handleSuccess,
    handleInfo,
    notifyApiError,
  }),
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: {
    SUCCESS: 'SUCCESS',
    ERROR: 'ERROR',
    INFO: 'INFO',
  },
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ push: routerPushMock }),
}));

import { useDeletePattern, usePublishPattern } from './hooks';

const wrap = (resolved: any = {}) => ({
  unwrap: () => Promise.resolve(resolved),
});

const reject = (error: any) => ({
  unwrap: () => Promise.reject(error),
});

describe('useDeletePattern', () => {
  beforeEach(() => {
    notify.mockReset();
    deletePatternByIdMock.mockReset();
    unassignDesignFromWorkspaceMock.mockReset();
    routerPushMock.mockReset();
  });

  it('deletes a pattern, notifies, and routes to /catalog/my-designs', async () => {
    deletePatternByIdMock.mockReturnValue(wrap({}));
    const { result } = renderHook(() => useDeletePattern());

    await act(async () => {
      await result.current.deletePattern({ id: 'p-1', name: 'pat' });
    });

    expect(deletePatternByIdMock).toHaveBeenCalledWith({ id: 'p-1' });
    expect(notify).toHaveBeenCalledWith({
      message: '"pat" deleted',
      event_type: 'SUCCESS',
    });
    expect(routerPushMock).toHaveBeenCalledWith({ pathname: '/catalog/my-designs' });
  });

  it('surfaces deletion errors via the notify hook', async () => {
    deletePatternByIdMock.mockReturnValue(reject({ data: { error: 'nope' } }));

    const { result } = renderHook(() => useDeletePattern());

    await act(async () => {
      await result.current.deletePattern({ id: 'p-1', name: 'pat' });
    });

    expect(notify).toHaveBeenCalledWith({
      message: 'nope',
      event_type: 'ERROR',
    });
  });

  it('surfaces generic deletion errors with default message', async () => {
    deletePatternByIdMock.mockReturnValue(reject({}));

    const { result } = renderHook(() => useDeletePattern());

    await act(async () => {
      await result.current.deletePattern({ id: 'p-1', name: 'pat' });
    });

    expect(notify).toHaveBeenCalledWith({
      message: 'Error deleting design',
      event_type: 'ERROR',
    });
  });

  it('handleDeleteModal runs deletion only when modal returns "Delete"', async () => {
    deletePatternByIdMock.mockReturnValue(wrap({}));
    const showMock = vi.fn().mockResolvedValue('Delete');
    const modalRef = { current: { show: showMock } };

    const { result } = renderHook(() => useDeletePattern());
    const handler = result.current.handleDeleteModal({ id: 'p-2', name: 'item' }, modalRef);

    await act(async () => {
      await handler();
    });

    expect(showMock).toHaveBeenCalledWith(
      expect.objectContaining({
        title: 'Delete Catalog item?',
        subtitle: 'Are you sure that you want to delete "item"?',
        primaryOption: 'Delete',
      }),
    );
    expect(deletePatternByIdMock).toHaveBeenCalledWith({ id: 'p-2' });
  });

  it('handleDeleteModal does NOT delete when modal returns something else', async () => {
    deletePatternByIdMock.mockReturnValue(wrap({}));
    const showMock = vi.fn().mockResolvedValue('Cancel');
    const modalRef = { current: { show: showMock } };

    const { result } = renderHook(() => useDeletePattern());
    const handler = result.current.handleDeleteModal({ id: 'p-2', name: 'item' }, modalRef);
    await act(async () => {
      await handler();
    });
    expect(deletePatternByIdMock).not.toHaveBeenCalled();
  });

  it('handleWorkspaceDesignDeleteModal calls the unassign mutator and notifies on success', async () => {
    unassignDesignFromWorkspaceMock.mockReturnValue(wrap({}));
    const { result } = renderHook(() => useDeletePattern());

    await act(async () => {
      await result.current.handleWorkspaceDesignDeleteModal('d-1', 'ws-1');
    });

    expect(unassignDesignFromWorkspaceMock).toHaveBeenCalledWith({
      workspaceId: 'ws-1',
      designId: 'd-1',
    });
    expect(notify).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Design removed from workspace',
        event_type: 'SUCCESS',
      }),
    );
  });

  it('handleBulkDeleteModal triggers a delete per pattern when modal returns "Delete"', async () => {
    deletePatternByIdMock.mockReturnValue(wrap({}));
    const showMock = vi.fn().mockResolvedValue('Delete');
    const modalRef = { current: { show: showMock } };

    const { result } = renderHook(() => useDeletePattern());
    const patterns = [
      { id: 'p-1', name: 'a' },
      { id: 'p-2', name: 'b' },
    ];

    await act(async () => {
      await result.current.handleBulkDeleteModal(patterns, modalRef);
    });

    expect(deletePatternByIdMock).toHaveBeenCalledTimes(2);
  });

  it('handleBulkDeleteModal is a no-op when modal returns "Cancel"', async () => {
    const showMock = vi.fn().mockResolvedValue('Cancel');
    const modalRef = { current: { show: showMock } };
    const { result } = renderHook(() => useDeletePattern());

    await act(async () => {
      await result.current.handleBulkDeleteModal([{ id: 'a', name: 'a' }], modalRef);
    });

    expect(deletePatternByIdMock).not.toHaveBeenCalled();
  });

  it('handleBulkWorkspaceDesignDeleteModal unassigns each design when "Delete" is returned', async () => {
    unassignDesignFromWorkspaceMock.mockReturnValue(wrap({}));
    const showMock = vi.fn().mockResolvedValue('Delete');
    const modalRef = { current: { show: showMock } };

    const { result } = renderHook(() => useDeletePattern());
    await act(async () => {
      await result.current.handleBulkWorkspaceDesignDeleteModal(
        [{ id: 'd-1' }, { id: 'd-2' }],
        modalRef,
        'wsName',
        'ws-1',
      );
    });

    expect(unassignDesignFromWorkspaceMock).toHaveBeenCalledTimes(2);
    expect(unassignDesignFromWorkspaceMock).toHaveBeenNthCalledWith(1, {
      workspaceId: 'ws-1',
      designId: 'd-1',
    });
    expect(unassignDesignFromWorkspaceMock).toHaveBeenNthCalledWith(2, {
      workspaceId: 'ws-1',
      designId: 'd-2',
    });
  });
});

describe('usePublishPattern', () => {
  beforeEach(() => {
    publishCatalogContentMock.mockReset();
    handleSuccess.mockReset();
    handleInfo.mockReset();
    notifyApiError.mockReset();
  });

  it('handlePublishModal returns the pattern as-is', () => {
    const { result } = renderHook(() => usePublishPattern(undefined, undefined));
    expect(result.current.handlePublishModal({ id: 'p' })).toEqual({ id: 'p' });
  });

  it('handlePublish maps compatibility names from displayNames and publishes', async () => {
    publishCatalogContentMock.mockReturnValue(wrap({ status: 'approved' }));
    const refetch = vi.fn();
    const models = {
      models: [
        { displayName: 'Kubernetes', name: 'kubernetes' },
        { displayName: 'kubernetes', name: 'kubernetes-duplicate' },
        { displayName: 'Istio', name: 'istio' },
      ],
    };

    const { result } = renderHook(() => usePublishPattern(models, refetch));

    await act(async () => {
      await result.current.handlePublish(
        { pattern: { id: 'p-1', name: 'design' } },
        {
          compatibility: ['Kubernetes', 'IsTIO'],
          type: 'WASM',
        },
      );
    });

    expect(publishCatalogContentMock).toHaveBeenCalledWith({
      type: 'pattern',
      mesheryCatalogRequestBody: {
        id: 'p-1',
        catalogType: 'pattern',
        catalogData: {
          compatibility: ['kubernetes', 'istio'],
          type: 'wasm',
        },
      },
    });
    expect(handleSuccess).toHaveBeenCalledWith('design published successfully');
    expect(refetch).toHaveBeenCalled();
  });

  it('handlePublish surfaces info when status is "pending"', async () => {
    publishCatalogContentMock.mockReturnValue(wrap({ status: 'PENDING' }));
    const { result } = renderHook(() => usePublishPattern({ models: [] }, undefined));

    await act(async () => {
      await result.current.handlePublish(
        { pattern: { id: 'p-1', name: 'design' } },
        { compatibility: [] },
      );
    });

    expect(handleInfo).toHaveBeenCalledWith(
      expect.stringContaining('design queued for publishing'),
    );
  });

  it('handlePublish reports errors via notifyApiError', async () => {
    publishCatalogContentMock.mockReturnValue(reject({ data: { message: 'fail' } }));
    const { result } = renderHook(() => usePublishPattern({ models: [] }, undefined));

    await act(async () => {
      await result.current.handlePublish(
        { pattern: { id: 'p-1', name: 'design' } },
        { compatibility: [] },
      );
    });

    expect(notifyApiError).toHaveBeenCalledWith(
      expect.objectContaining({ data: { message: 'fail' } }),
      'Failed to publish "design" to the catalog',
    );
  });
});
