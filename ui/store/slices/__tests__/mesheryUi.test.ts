import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

// Mock dependencies — these are pulled in transitively when the slice module
// loads. They must be mocked before importing the slice.
vi.mock('@/utils/eventBus', () => ({
  mesheryEventBus: {
    publish: vi.fn(),
    on: vi.fn(() => ({ subscribe: vi.fn() })),
  },
}));

const persistSelectedK8sContextsMock = vi.fn();

vi.mock('@/utils/multi-ctx', () => ({
  persistSelectedK8sContexts: (...args: unknown[]) => persistSelectedK8sContextsMock(...args),
  getK8sClusterIdsFromCtxId: vi.fn((selectedContexts, k8sconfig) => {
    if (!selectedContexts || !k8sconfig || selectedContexts.length === 0) return [];
    if (selectedContexts.includes('all')) {
      return k8sconfig.map((c: any) => c.kubernetesServerId);
    }
    return selectedContexts.map((ctx: any) => {
      const cfg = k8sconfig.find((c: any) => c.id === ctx);
      return cfg?.kubernetesServerId;
    });
  }),
}));

vi.mock('../..', () => ({
  store: {
    dispatch: vi.fn(),
  },
}));

import mesheryUiReducer, {
  updatePage,
  updatePagePath,
  updateTitle,
  updateBetaBadge,
  updateUser,
  updateK8SConfig,
  setK8sContextsAction,
  updateProgressAction,
  toggleDrawer,
  toggleCatalogContent,
  setControllerState,
  updateExtensionType,
  updateProviderCapabilities,
  setConnectionMetadata,
  setOrganization,
  setWorkspace,
  setKeys,
  setK8sContexts,
  updateProgress,
  selectSelectedK8sClusters,
  selectK8sConfig,
  selectedOrg,
} from '../mesheryUi';
import { mesheryEventBus } from '@/utils/eventBus';
import { store } from '../..';

describe('mesheryUi slice', () => {
  beforeEach(() => {
    // jsdom provides sessionStorage; clear before each test
    sessionStorage.clear();
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('returns the initial state', () => {
    const state = mesheryUiReducer(undefined, { type: 'init' } as any);
    expect(state.page).toEqual({ path: '', title: '', isBeta: false });
    expect(state.user).toEqual({});
    expect(state.k8sConfig).toEqual([]);
    expect(state.selectedK8sContexts).toEqual(['all']);
    expect(state.showProgress).toBe(false);
    expect(state.isDrawerCollapsed).toBe(false);
    expect(state.catalogVisibility).toBe(true);
    expect(state.extensionType).toBe('');
    expect(state.providerCapabilities).toBeNull();
    expect(state.controllerState).toBeNull();
    expect(state.connectionMetadataState).toBeNull();
    expect(state.organization).toBeNull();
    expect(state.workspace).toBeNull();
    expect(state.keys).toBeNull();
  });

  it('updatePage merges into page', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    state = mesheryUiReducer(state, updatePage({ path: '/foo', title: 'Foo' }));
    expect(state.page).toEqual({ path: '/foo', title: 'Foo', isBeta: false });
  });

  it('updatePagePath updates only the path', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    state = mesheryUiReducer(state, updatePagePath({ path: '/x' }));
    expect(state.page.path).toBe('/x');
    expect(state.page.title).toBe('');
  });

  it('updateTitle updates the title', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    state = mesheryUiReducer(state, updateTitle({ title: 'Settings' }));
    expect(state.page.title).toBe('Settings');
  });

  it('updateBetaBadge updates isBeta', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    state = mesheryUiReducer(state, updateBetaBadge({ isBeta: true }));
    expect(state.page.isBeta).toBe(true);
    state = mesheryUiReducer(state, updateBetaBadge({ isBeta: false }));
    expect(state.page.isBeta).toBe(false);
  });

  it('updateUser replaces the user object', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    state = mesheryUiReducer(state, updateUser({ user: { id: 'u1', name: 'alice' } }));
    expect(state.user).toEqual({ id: 'u1', name: 'alice' });
  });

  it('updateK8SConfig replaces k8sConfig', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    const cfg = [{ id: 'c1', name: 'cluster-a' }];
    state = mesheryUiReducer(state, updateK8SConfig({ k8sConfig: cfg }));
    expect(state.k8sConfig).toBe(cfg);
  });

  it('setK8sContextsAction reducer updates selectedK8sContexts', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    state = mesheryUiReducer(
      state,
      setK8sContextsAction({ selectedK8sContexts: ['ctx-1', 'ctx-2'] }),
    );
    expect(state.selectedK8sContexts).toEqual(['ctx-1', 'ctx-2']);
  });

  it('updateProgressAction updates showProgress', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    state = mesheryUiReducer(state, updateProgressAction({ showProgress: true }));
    expect(state.showProgress).toBe(true);
  });

  it('toggleDrawer updates isDrawerCollapsed', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    state = mesheryUiReducer(state, toggleDrawer({ isDrawerCollapsed: true }));
    expect(state.isDrawerCollapsed).toBe(true);
  });

  it('toggleCatalogContent updates catalogVisibility', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    state = mesheryUiReducer(state, toggleCatalogContent({ catalogVisibility: false }));
    expect(state.catalogVisibility).toBe(false);
  });

  it('setControllerState updates controllerState', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    state = mesheryUiReducer(state, setControllerState({ controllerState: { ok: true } }));
    expect(state.controllerState).toEqual({ ok: true });
  });

  it('updateExtensionType updates extensionType', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    state = mesheryUiReducer(state, updateExtensionType({ extensionType: 'plugin' }));
    expect(state.extensionType).toBe('plugin');
  });

  it('updateProviderCapabilities updates providerCapabilities', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    const caps = { foo: 'bar' };
    state = mesheryUiReducer(state, updateProviderCapabilities({ providerCapabilities: caps }));
    expect(state.providerCapabilities).toBe(caps);
  });

  it('setConnectionMetadata updates connectionMetadataState', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    const meta = { kind: 'k8s' };
    state = mesheryUiReducer(state, setConnectionMetadata({ connectionMetadataState: meta }));
    expect(state.connectionMetadataState).toBe(meta);
  });

  it('setOrganization updates state and writes to sessionStorage', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    const org = { id: 'org-1', name: 'acme' };
    state = mesheryUiReducer(state, setOrganization({ organization: org }));
    expect(state.organization).toBe(org);
    expect(sessionStorage.getItem('currentOrg')).toBe(JSON.stringify(org));
  });

  it('setWorkspace updates state and writes to sessionStorage', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    const ws = { id: 'ws-1', name: 'team' };
    state = mesheryUiReducer(state, setWorkspace({ workspace: ws }));
    expect(state.workspace).toBe(ws);
    expect(sessionStorage.getItem('currentWorkspace')).toBe(JSON.stringify(ws));
  });

  it('setKeys updates state and writes to sessionStorage', () => {
    let state = mesheryUiReducer(undefined, { type: 'init' } as any);
    const keys = { canEdit: true };
    state = mesheryUiReducer(state, setKeys({ keys }));
    expect(state.keys).toBe(keys);
    expect(sessionStorage.getItem('keys')).toBe(JSON.stringify(keys));
  });

  describe('setK8sContexts thunk', () => {
    it('dispatches the reducer action and publishes to the event bus', () => {
      const dispatch = vi.fn();
      const payload = { selectedK8sContexts: ['ctx-1'] };
      setK8sContexts(payload)(dispatch);

      expect(dispatch).toHaveBeenCalledTimes(1);
      const dispatched = dispatch.mock.calls[0][0];
      expect(dispatched.type).toBe('core/setK8sContexts');
      expect(dispatched.payload).toEqual(payload);

      expect(mesheryEventBus.publish).toHaveBeenCalledWith({
        type: 'K8S_CONTEXTS_UPDATED',
        data: { selectedK8sContexts: ['ctx-1'] },
      });
    });

    it('session-persists the selection (side effect lives in the thunk, not the reducer)', () => {
      const dispatch = vi.fn();
      setK8sContexts({ selectedK8sContexts: ['ctx-1', 'ctx-2'] })(dispatch);

      expect(persistSelectedK8sContextsMock).toHaveBeenCalledWith(['ctx-1', 'ctx-2']);
    });
  });

  describe('updateProgress non-thunk', () => {
    it('dispatches updateProgressAction to the global store', () => {
      const payload = { showProgress: true };
      updateProgress(payload);
      expect(store.dispatch).toHaveBeenCalledTimes(1);
      const dispatched = (store.dispatch as ReturnType<typeof vi.fn>).mock.calls[0][0];
      expect(dispatched.type).toBe('core/updateProgress');
      expect(dispatched.payload).toEqual(payload);
    });
  });

  describe('selectors', () => {
    it('selectK8sConfig returns ui.k8sConfig', () => {
      const cfg = [{ id: 'a' }];
      expect(selectK8sConfig({ ui: { k8sConfig: cfg } } as any)).toBe(cfg);
    });

    it('selectedOrg returns ui.organization', () => {
      const org = { id: 'o' };
      expect(selectedOrg({ ui: { organization: org } } as any)).toBe(org);
    });

    it('selectSelectedK8sClusters delegates to getK8sClusterIdsFromCtxId', () => {
      const state = {
        ui: {
          selectedK8sContexts: ['all'],
          k8sConfig: [{ id: 'c1', kubernetesServerId: 'srv-1' }],
        },
      };
      const result = selectSelectedK8sClusters(state as any);
      expect(result).toEqual(['srv-1']);
    });
  });

  it('immutability check: reducers return new state, not the same reference', () => {
    const initial = mesheryUiReducer(undefined, { type: 'init' } as any);
    const next = mesheryUiReducer(initial, updateTitle({ title: 'foo' }));
    expect(next).not.toBe(initial);
    expect(initial.page.title).toBe('');
  });

  it('action types have the expected slice prefix', () => {
    expect(updatePage({}).type).toBe('core/updatePage');
    expect(updatePagePath({ path: '' }).type).toBe('core/updatePagePath');
    expect(updateUser({ user: {} }).type).toBe('core/updateUser');
    expect(setK8sContextsAction({ selectedK8sContexts: [] }).type).toBe('core/setK8sContexts');
    expect(updateProgressAction({ showProgress: false }).type).toBe('core/updateProgress');
  });
});
