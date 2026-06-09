import { beforeEach, describe, expect, it, vi } from 'vitest';

const encodeDesignFile = vi.fn((d: any) => `encoded:${JSON.stringify(d)}`);
const getUnit8ArrayDecodedFile = vi.fn((f: any) => `decoded:${f}`);
const downloadContent = vi.fn();
const buildImportDesignRequestBody = vi.fn();

vi.mock('@sistent/sistent', () => ({
  PROMPT_VARIANTS: { DANGER: 'danger' },
}));

vi.mock('../../../utils/utils', () => ({
  encodeDesignFile: (d: any) => encodeDesignFile(d),
  getUnit8ArrayDecodedFile: (f: any) => getUnit8ArrayDecodedFile(f),
}));

vi.mock('../../../utils/Enum', () => ({
  FILE_OPS: {
    DELETE: 'DELETE',
    UPDATE: 'UPDATE',
    FILE_UPLOAD: 'FILE_UPLOAD',
    URL_UPLOAD: 'URL_UPLOAD',
  },
}));

vi.mock('../../../lib/event-types', () => ({
  EVENT_TYPES: {
    ERROR: 'error',
    SUCCESS: 'success',
    WARNING: 'warning',
    INFO: 'info',
  },
}));

vi.mock('../../../utils/fileDownloader', () => ({
  default: (...args: any[]) => downloadContent(...args),
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: vi.fn(),
}));

vi.mock('../import-design-request', () => ({
  buildImportDesignRequestBody: (data: any) => buildImportDesignRequestBody(data),
}));

import { createPatternsActions } from './patterns-actions';

const baseDeps = () => {
  const setSelectedRowData = vi.fn();
  const setImportModal = vi.fn();
  const setInfoModal = vi.fn();
  const setPublishModal = vi.fn();
  const notify = vi.fn();
  const getPatterns = vi.fn();
  const sistentInfoModal = {
    openModal: vi.fn(),
    closeModal: vi.fn(),
  };
  return {
    clonePattern: vi.fn(),
    publishCatalog: vi.fn(),
    unpublishCatalog: vi.fn(),
    deletePattern: vi.fn(),
    deletePatternFile: vi.fn(),
    importPattern: vi.fn(),
    updatePattern: vi.fn(),
    uploadPatternFile: vi.fn(),
    deployPatternMutation: vi.fn().mockResolvedValue({}),
    undeployPatternMutation: vi.fn().mockResolvedValue({}),
    modalRef: { current: { show: vi.fn() } },
    meshModels: [],
    infoModal: { selectedResource: { id: 'res-1' } },
    publishModal: { name: 'Catalog Item' },
    user: { roleNames: ['admin'] },
    setImportModal,
    setPublishModal,
    setSelectedRowData,
    setInfoModal,
    notify,
    sistentInfoModal,
    getPatterns,
  };
};

describe('createPatternsActions', () => {
  beforeEach(() => {
    buildImportDesignRequestBody.mockReset();
  });

  it('returns the canonical set of action methods', () => {
    const actions = createPatternsActions(baseDeps());
    const keys = [
      'handleError',
      'resetSelectedRowData',
      'handleDeploy',
      'handleUndeploy',
      'handleUploadImport',
      'handleUploadImportClose',
      'handleInfoModalClose',
      'handleInfoModal',
      'handleUnpublishModal',
      'handlePublishModalClose',
      'handlePublish',
      'handleClone',
      'handleSubmit',
      'handleImportDesign',
      'deletePatterns',
      'handleDownload',
      'showModal',
    ];
    keys.forEach((key) => expect(typeof (actions as any)[key]).toBe('function'));
  });

  it('handleDeploy encodes the design file and invokes the deploy mutation', async () => {
    const deps = baseDeps();
    const actions = createPatternsActions(deps);

    await actions.handleDeploy({ design: { id: 'd1', body: 'x' }, selectedK8sContexts: ['ctx-1'] });

    expect(encodeDesignFile).toHaveBeenCalledWith({ id: 'd1', body: 'x' });
    expect(deps.deployPatternMutation).toHaveBeenCalledWith({
      patternFile: 'encoded:{"id":"d1","body":"x"}',
      patternId: 'd1',
      selectedK8sContexts: ['ctx-1'],
    });
  });

  it('handleUndeploy invokes the undeploy mutation', async () => {
    const deps = baseDeps();
    const actions = createPatternsActions(deps);

    await actions.handleUndeploy({ design: { id: 'd2' }, selectedK8sContexts: ['ctx-2'] });

    expect(deps.undeployPatternMutation).toHaveBeenCalledWith({
      patternFile: 'encoded:{"id":"d2"}',
      patternId: 'd2',
      selectedK8sContexts: ['ctx-2'],
    });
  });

  it('handleUploadImport opens the import modal', () => {
    const deps = baseDeps();
    const actions = createPatternsActions(deps);
    actions.handleUploadImport();
    expect(deps.setImportModal).toHaveBeenCalledWith({ open: true });
  });

  it('handleUploadImportClose closes the import modal', () => {
    const deps = baseDeps();
    const actions = createPatternsActions(deps);
    actions.handleUploadImportClose();
    expect(deps.setImportModal).toHaveBeenCalledWith({ open: false });
  });

  it('handleInfoModal opens both the sistent modal and the in-app info modal', () => {
    const deps = baseDeps();
    const actions = createPatternsActions(deps);
    actions.handleInfoModal({ name: 'My Pattern', userId: 'u-1', id: 'p-1' });

    expect(deps.sistentInfoModal.openModal).toHaveBeenCalledWith({ title: 'My Pattern' });
    expect(deps.setInfoModal).toHaveBeenCalledWith({
      open: true,
      ownerID: 'u-1',
      selectedResource: { name: 'My Pattern', userId: 'u-1', id: 'p-1' },
    });
  });

  it('handleInfoModalClose closes both sistent and the in-app info modal', () => {
    const deps = baseDeps();
    const actions = createPatternsActions(deps);
    actions.handleInfoModalClose();
    expect(deps.sistentInfoModal.closeModal).toHaveBeenCalled();
    expect(deps.setInfoModal).toHaveBeenCalledWith({ open: false });
  });

  it('handlePublishModalClose resets publish modal state', () => {
    const deps = baseDeps();
    const actions = createPatternsActions(deps);
    actions.handlePublishModalClose();
    expect(deps.setPublishModal).toHaveBeenCalledWith({ open: false, pattern: {}, name: '' });
  });

  it('resetSelectedRowData returns a callable that clears the selection', () => {
    const deps = baseDeps();
    const actions = createPatternsActions(deps);
    actions.resetSelectedRowData()();
    expect(deps.setSelectedRowData).toHaveBeenCalledWith(null);
  });

  it('handleDownload triggers download and notifies', () => {
    const deps = baseDeps();
    const actions = createPatternsActions(deps);
    const stop = vi.fn();
    actions.handleDownload(
      { stopPropagation: stop } as any,
      { id: 'd1', name: 'Design 1' },
      'manifest',
      { param: 'v' },
    );

    expect(stop).toHaveBeenCalled();
    expect(downloadContent).toHaveBeenCalledWith({
      id: 'd1',
      name: 'Design 1',
      type: 'pattern',
      source_type: 'manifest',
      params: { param: 'v' },
    });
    expect(deps.notify).toHaveBeenCalled();
  });

  it('handleError surfaces the failure as a notification', () => {
    const deps = baseDeps();
    const actions = createPatternsActions(deps);
    actions.handleError({ name: 'FAIL', error_msg: 'Something went wrong' })('boom');
    expect(deps.notify).toHaveBeenCalledWith({
      message: 'Something went wrong: boom',
      event_type: 'error',
    });
  });

  it('handleImportDesign uploads file imports using the resolved file metadata', async () => {
    const deps = baseDeps();
    deps.importPattern.mockReturnValue({ unwrap: () => Promise.resolve({}) });
    buildImportDesignRequestBody.mockResolvedValue({
      requestBody: JSON.stringify({
        name: 'Imported design',
        file_name: 'imported-design.yaml',
        file: [1, 2, 3],
      }),
    });

    const actions = createPatternsActions(deps);
    await actions.handleImportDesign({
      uploadType: 'File Upload',
      name: 'Imported design',
      file: 'data:text/plain;base64,QQ==',
    });

    expect(buildImportDesignRequestBody).toHaveBeenCalledWith({
      uploadType: 'File Upload',
      name: 'Imported design',
      file: 'data:text/plain;base64,QQ==',
    });
    expect(deps.importPattern).toHaveBeenCalledWith({
      importBody: JSON.stringify({
        name: 'Imported design',
        file_name: 'imported-design.yaml',
        file: [1, 2, 3],
      }),
    });
  });

  it('handleImportDesign surfaces a missing file instead of dereferencing a null input', async () => {
    const deps = baseDeps();
    buildImportDesignRequestBody.mockResolvedValue({
      errorMessage: 'Please choose a design file before continuing.',
    });

    const actions = createPatternsActions(deps);
    await actions.handleImportDesign({
      uploadType: 'File Upload',
      name: 'Imported design',
      file: undefined,
    });

    expect(deps.importPattern).not.toHaveBeenCalled();
    expect(deps.notify).toHaveBeenCalledWith({
      message: 'Please choose a design file before continuing.',
      event_type: 'error',
    });
  });
});
