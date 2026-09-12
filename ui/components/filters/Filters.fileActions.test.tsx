import { describe, expect, it, vi, beforeEach } from 'vitest';

const updateProgress = vi.fn();
const downloadContent = vi.fn();
const trueRandomMock = vi.fn(() => 0.5);

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: (...args: unknown[]) => updateProgress(...args),
}));

vi.mock('../../utils/Enum', () => ({
  FILE_OPS: {
    DELETE: 'DELETE',
    FILE_UPLOAD: 'FILE_UPLOAD',
    URL_UPLOAD: 'URL_UPLOAD',
    UPDATE: 'UPDATE',
  },
}));

vi.mock('../../lib/event-types', () => ({
  EVENT_TYPES: {
    SUCCESS: 'SUCCESS',
    ERROR: 'ERROR',
    INFO: 'INFO',
  },
}));

vi.mock('../../utils/utils', () => ({
  getUnit8ArrayDecodedFile: (s: any) => `decoded:${s}`,
}));

vi.mock('../../lib/trueRandom', () => ({
  trueRandom: () => trueRandomMock(),
}));

vi.mock('../../utils/fileDownloader', () => ({
  default: (...args: unknown[]) => downloadContent(...args),
}));

import {
  createHandleSubmit,
  createHandleDownload,
  createUploadHandler,
  createHandleImportFilter,
  createHandleUnpublishModal,
  createHandlePublish,
  createHandleClone,
  createDeleteFilter,
} from './Filters.fileActions';

const flush = () => new Promise((resolve) => setTimeout(resolve, 0));

const wrap = (resolved: any = {}) => ({ unwrap: () => Promise.resolve(resolved) });
const reject = () => ({ unwrap: () => Promise.reject(new Error('boom')) });

beforeEach(() => {
  updateProgress.mockReset();
  downloadContent.mockReset();
  trueRandomMock.mockReset();
  trueRandomMock.mockReturnValue(0.5);
});

describe('createHandleSubmit', () => {
  it('deletes via deleteFilterFile only when modal returns "Delete"', async () => {
    const deleteFilterFile = vi.fn(() => wrap());
    const showmodal = vi.fn().mockResolvedValueOnce('Delete');
    const notify = vi.fn();
    const handleError = vi.fn(() => () => {});
    const resetSelectedRowData = vi.fn(() => vi.fn());

    const handleSubmit = createHandleSubmit({
      notify,
      handleError,
      showmodal,
      resetSelectedRowData,
      deleteFilterFile,
      uploadFilterFile: vi.fn(),
      updateFilterFile: vi.fn(),
    });

    await handleSubmit({ data: '', name: 'foo', id: 'i-1', type: 'DELETE' });
    await flush();

    expect(showmodal).toHaveBeenCalledWith(1, 'foo');
    expect(deleteFilterFile).toHaveBeenCalledWith({ id: 'i-1' });
    expect(notify).toHaveBeenCalledWith({
      message: '"foo" filter deleted',
      event_type: 'SUCCESS',
    });
  });

  it('aborts when modal cancels the delete', async () => {
    const deleteFilterFile = vi.fn(() => wrap());
    const showmodal = vi.fn().mockResolvedValueOnce('Cancel');
    const resetSelectedRowData = vi.fn(() => vi.fn());

    const handleSubmit = createHandleSubmit({
      notify: vi.fn(),
      handleError: vi.fn(() => () => {}),
      showmodal,
      resetSelectedRowData,
      deleteFilterFile,
      uploadFilterFile: vi.fn(),
      updateFilterFile: vi.fn(),
    });

    await handleSubmit({ data: '', name: 'foo', id: 'i-1', type: 'DELETE' });
    expect(deleteFilterFile).not.toHaveBeenCalled();
  });

  it('handles a file upload by sending a JSON body to uploadFilterFile', async () => {
    const uploadFilterFile = vi.fn(() => wrap());

    const handleSubmit = createHandleSubmit({
      notify: vi.fn(),
      handleError: vi.fn(() => () => {}),
      showmodal: vi.fn(),
      resetSelectedRowData: vi.fn(() => vi.fn()),
      deleteFilterFile: vi.fn(),
      uploadFilterFile,
      updateFilterFile: vi.fn(),
    });

    await handleSubmit({
      data: 'binary',
      name: 'foo',
      type: 'FILE_UPLOAD',
      metadata: { name: 'foo', config: { x: 1 } },
    });

    expect(uploadFilterFile).toHaveBeenCalled();
    const body = JSON.parse(uploadFilterFile.mock.calls[0][0].uploadBody);
    expect(body).toEqual({
      save: true,
      filterData: { filterFile: 'binary', name: 'foo' },
      config: { x: 1 },
    });
  });

  it('handles a URL upload', async () => {
    const uploadFilterFile = vi.fn(() => wrap());
    const handleSubmit = createHandleSubmit({
      notify: vi.fn(),
      handleError: vi.fn(() => () => {}),
      showmodal: vi.fn(),
      resetSelectedRowData: vi.fn(() => vi.fn()),
      deleteFilterFile: vi.fn(),
      uploadFilterFile,
      updateFilterFile: vi.fn(),
    });

    await handleSubmit({
      data: 'http://example.com',
      name: 'foo',
      type: 'URL_UPLOAD',
      metadata: { name: 'foo', config: { y: 2 } },
    });

    const body = JSON.parse(uploadFilterFile.mock.calls[0][0].uploadBody);
    expect(body).toEqual({
      save: true,
      url: 'http://example.com',
      name: 'foo',
      config: { y: 2 },
    });
  });

  it('handles an UPDATE with the filter id, name and catalogData', async () => {
    const updateFilterFile = vi.fn(() => wrap());
    const handleSubmit = createHandleSubmit({
      notify: vi.fn(),
      handleError: vi.fn(() => () => {}),
      showmodal: vi.fn(),
      resetSelectedRowData: vi.fn(() => vi.fn()),
      deleteFilterFile: vi.fn(),
      uploadFilterFile: vi.fn(),
      updateFilterFile,
    });

    await handleSubmit({
      data: 'apiVersion: v1',
      id: 'i-1',
      name: 'foo',
      type: 'UPDATE',
      catalogData: { kind: 'wasm' },
    });

    const body = JSON.parse(updateFilterFile.mock.calls[0][0].updateBody);
    expect(body).toEqual({
      filterData: { id: 'i-1', name: 'foo', catalogData: { kind: 'wasm' } },
      config: 'apiVersion: v1',
      save: true,
    });
  });
});

describe('createHandleDownload', () => {
  it('downloads via downloadContent and emits an info notification', () => {
    const notify = vi.fn();
    const handler = createHandleDownload({ notify });
    const ev = { stopPropagation: vi.fn() } as any;
    handler(ev, 'i-1', 'foo');
    expect(ev.stopPropagation).toHaveBeenCalled();
    expect(downloadContent).toHaveBeenCalledWith({ id: 'i-1', name: 'foo', type: 'filter' });
    expect(notify).toHaveBeenCalledWith({
      message: '"foo" filter downloaded',
      event_type: 'INFO',
    });
  });
});

describe('createUploadHandler', () => {
  it('reads the file as ArrayBuffer and calls handleSubmit with the byte array', async () => {
    const handleSubmit = vi.fn();
    const upload = createUploadHandler({ handleSubmit });

    const file = new File([new Uint8Array([1, 2, 3])], 'a.wasm');
    const ev = { target: { files: [file] } } as any;

    upload(ev, undefined, { name: 'a.wasm', config: { c: 1 } });

    // FileReader is async; wait briefly
    await new Promise((resolve) => setTimeout(resolve, 50));

    expect(handleSubmit).toHaveBeenCalledWith(
      expect.objectContaining({
        type: 'FILE_UPLOAD',
        name: 'a.wasm',
        metadata: { name: 'a.wasm', config: { c: 1 } },
      }),
    );
  });

  it('no-ops if there are no files in the input event', () => {
    const handleSubmit = vi.fn();
    const upload = createUploadHandler({ handleSubmit });
    upload({ target: { files: null } } as any, undefined, {});
    expect(handleSubmit).not.toHaveBeenCalled();
  });
});

describe('createHandleImportFilter', () => {
  it('imports a File Upload payload and refreshes filters on success', async () => {
    const updateFilterFile = vi.fn(() => wrap());
    const notify = vi.fn();
    const getFilters = vi.fn();
    const handler = createHandleImportFilter({
      notify,
      handleError: vi.fn(() => () => {}),
      updateFilterFile,
      getFilters,
    });

    handler({
      uploadType: 'File Upload',
      config: 'cfg',
      name: 'n',
      url: '',
      file: 'rawdata' as any,
    });
    await flush();

    const body = JSON.parse(updateFilterFile.mock.calls[0][0].updateBody);
    expect(body).toEqual({
      config: 'cfg',
      save: true,
      filter_data: { name: 'n', filter_file: 'decoded:rawdata' },
    });
    expect(notify).toHaveBeenCalledWith({
      message: '"n" filter uploaded',
      event_type: 'SUCCESS',
    });
    expect(getFilters).toHaveBeenCalled();
  });

  it('imports a URL Upload payload', async () => {
    const updateFilterFile = vi.fn(() => wrap());
    const handler = createHandleImportFilter({
      notify: vi.fn(),
      handleError: vi.fn(() => () => {}),
      updateFilterFile,
      getFilters: vi.fn(),
    });

    handler({
      uploadType: 'URL Upload',
      config: 'cfg',
      name: 'n',
      url: 'http://e.com',
      file: '' as any,
    });
    await flush();

    const body = JSON.parse(updateFilterFile.mock.calls[0][0].updateBody);
    expect(body).toEqual({
      config: 'cfg',
      save: true,
      url: 'http://e.com',
      filter_data: { name: 'n' },
    });
  });

  it('handles errors by calling handleError(UPLOAD_FILTERS)', async () => {
    const updateFilterFile = vi.fn(() => reject());
    const errorHandler = vi.fn();
    const handleError = vi.fn(() => errorHandler);
    const handler = createHandleImportFilter({
      notify: vi.fn(),
      handleError,
      updateFilterFile,
      getFilters: vi.fn(),
    });

    handler({
      uploadType: 'File Upload',
      config: 'cfg',
      name: 'n',
      url: '',
      file: 'rd' as any,
    });
    await flush();

    expect(handleError).toHaveBeenCalledWith(expect.objectContaining({ name: 'UPLOAD_FILTERS' }));
  });
});

describe('createHandleUnpublishModal', () => {
  it('returns a function that unpublishes only when modal yields "Yes" and canPublishFilter is true', async () => {
    const unpublishFilter = vi.fn(() => wrap());
    const showMock = vi.fn().mockResolvedValueOnce('Yes');
    const modalRef = { current: { show: showMock } } as any;
    const notify = vi.fn();

    const handleUnpublishModal = createHandleUnpublishModal({
      canPublishFilter: true,
      modalRef,
      unpublishFilter,
      notify,
      handleError: vi.fn(() => () => {}),
    });

    const fn = handleUnpublishModal({} as any, { id: 'f-1', name: 'flt' });
    expect(typeof fn).toBe('function');
    await fn!();
    await flush();

    expect(unpublishFilter).toHaveBeenCalledWith({
      unpublishBody: JSON.stringify({ id: 'f-1' }),
    });
    expect(notify).toHaveBeenCalledWith({
      message: '"flt" filter unpublished',
      event_type: 'SUCCESS',
    });
  });

  it('returns undefined when canPublishFilter is false', () => {
    const handler = createHandleUnpublishModal({
      canPublishFilter: false,
      modalRef: { current: null } as any,
      unpublishFilter: vi.fn(),
      notify: vi.fn(),
      handleError: vi.fn(() => () => {}),
    });
    expect(handler({} as any, {})).toBeUndefined();
  });
});

describe('createHandlePublish', () => {
  it('builds the catalog payload using normalised compatibility names', async () => {
    const publishFilter = vi.fn(() => wrap());
    const notify = vi.fn();
    const meshModels = [
      { displayName: 'Kubernetes', name: 'kubernetes' },
      { displayName: 'kubernetes', name: 'kubernetes-dup' },
      { displayName: 'Istio', name: 'istio' },
    ];
    const handler = createHandlePublish({
      meshModels,
      publishModal: { filter: { id: 'f', name: 'n' } } as any,
      user: { roleNames: ['admin'] } as any,
      publishFilter,
      notify,
      handleError: vi.fn(() => () => {}),
    });

    handler({ compatibility: ['Kubernetes', 'ISTIO'], type: 'WASM' });
    await flush();

    expect(publishFilter).toHaveBeenCalled();
    const body = JSON.parse(publishFilter.mock.calls[0][0].publishBody);
    expect(body).toEqual({
      id: 'f',
      catalogData: {
        compatibility: ['kubernetes', 'istio'],
        type: 'wasm',
      },
    });
    expect(notify).toHaveBeenCalledWith({
      message: 'n filter published to Meshery Catalog',
      event_type: 'SUCCESS',
    });
  });

  it('uses the moderation-flow notification for non-admin users', async () => {
    const publishFilter = vi.fn(() => wrap());
    const notify = vi.fn();
    const handler = createHandlePublish({
      meshModels: [],
      publishModal: { filter: { id: 'f', name: 'n' } } as any,
      user: { roleNames: ['user'] } as any,
      publishFilter,
      notify,
      handleError: vi.fn(() => () => {}),
    });

    handler({});
    await flush();
    expect(notify).toHaveBeenCalledWith(
      expect.objectContaining({ message: expect.stringContaining('Maintainers notified') }),
    );
  });
});

describe('createHandleClone', () => {
  it('clones a filter and notifies on success', async () => {
    const cloneFilter = vi.fn(() => wrap());
    const notify = vi.fn();
    const clone = createHandleClone({
      cloneFilter,
      notify,
      handleError: vi.fn(() => () => {}),
    });
    clone('f-1', 'flt');
    await flush();
    expect(cloneFilter).toHaveBeenCalledWith({
      body: JSON.stringify({ name: 'flt (Copy)' }),
      filterID: 'f-1',
    });
    expect(notify).toHaveBeenCalledWith({
      message: '"flt" filter cloned',
      event_type: 'SUCCESS',
    });
  });
});

describe('createDeleteFilter', () => {
  it('deletes a filter and notifies on success', async () => {
    const deleteFilterFile = vi.fn(() => wrap());
    const notify = vi.fn();
    const del = createDeleteFilter({
      deleteFilterFile,
      notify,
      handleError: vi.fn(() => () => {}),
    });
    del('f-1');
    await flush();
    expect(deleteFilterFile).toHaveBeenCalledWith({ id: 'f-1' });
    expect(notify).toHaveBeenCalledWith({
      message: 'Filter deleted',
      event_type: 'SUCCESS',
    });
  });
});
