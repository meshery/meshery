import { describe, expect, it, vi } from 'vitest';
import globalEnvironmentContextReducer, {
  selectEnv,
  unselectEnv,
  selectConnection,
  unselectConnection,
  toggleConnection,
  toggleEnvSelection,
  selectIsConnectionSelected,
  selectIsEnvSelected,
  selectSelectedEnvs,
  selectSelectedConnections,
  selectAllSelectedConnections,
  selectAllSelectedK8sConnections,
  selectSelectedK8sConnections,
} from '../globalEnvironmentContext';

const env1 = { id: 'env-1', name: 'Production' };
const env2 = { id: 'env-2', name: 'Staging' };
const k8sConn = { id: 'c-1', kind: 'kubernetes', name: 'cluster-a' };
const k8sConn2 = { id: 'c-2', kind: 'kubernetes', name: 'cluster-b' };
const dockerConn = { id: 'd-1', kind: 'docker', name: 'docker-host' };

describe('globalEnvironmentContext slice', () => {
  it('returns the initial state for unknown actions', () => {
    const state = globalEnvironmentContextReducer(undefined, { type: 'init' } as any);
    expect(state).toEqual({ selectedEnvs: {} });
  });

  describe('selectEnv', () => {
    it('adds an environment with no selected connections by default', () => {
      const state = globalEnvironmentContextReducer(undefined, selectEnv({ environment: env1 }));
      expect(state.selectedEnvs[env1.id]).toEqual({
        ...env1,
        selectedConnections: [],
      });
    });

    it('adds an environment with provided connections', () => {
      const state = globalEnvironmentContextReducer(
        undefined,
        selectEnv({ environment: env1, selectedConnections: [k8sConn] }),
      );
      expect(state.selectedEnvs[env1.id].selectedConnections).toEqual([k8sConn]);
    });

    it('overwrites a previously selected environment with the same id', () => {
      let state = globalEnvironmentContextReducer(
        undefined,
        selectEnv({ environment: env1, selectedConnections: [k8sConn] }),
      );
      state = globalEnvironmentContextReducer(
        state,
        selectEnv({ environment: env1, selectedConnections: [k8sConn2] }),
      );
      expect(state.selectedEnvs[env1.id].selectedConnections).toEqual([k8sConn2]);
    });

    it('supports multiple environments at once', () => {
      let state = globalEnvironmentContextReducer(
        undefined,
        selectEnv({ environment: env1, selectedConnections: [k8sConn] }),
      );
      state = globalEnvironmentContextReducer(
        state,
        selectEnv({ environment: env2, selectedConnections: [k8sConn2] }),
      );
      expect(Object.keys(state.selectedEnvs)).toHaveLength(2);
    });
  });

  describe('unselectEnv', () => {
    it('removes an environment from selectedEnvs', () => {
      let state = globalEnvironmentContextReducer(undefined, selectEnv({ environment: env1 }));
      state = globalEnvironmentContextReducer(state, unselectEnv({ envId: env1.id }));
      expect(state.selectedEnvs[env1.id]).toBeUndefined();
    });

    it('is a no-op when env was not previously selected', () => {
      const initial = globalEnvironmentContextReducer(undefined, { type: 'init' } as any);
      const state = globalEnvironmentContextReducer(initial, unselectEnv({ envId: 'missing' }));
      expect(state.selectedEnvs).toEqual({});
    });
  });

  describe('selectConnection', () => {
    it('creates the env entry if it does not exist yet', () => {
      const state = globalEnvironmentContextReducer(
        undefined,
        selectConnection({ env: env1, connection: k8sConn }),
      );
      expect(state.selectedEnvs[env1.id].selectedConnections).toEqual([k8sConn]);
    });

    it('appends to existing selectedConnections', () => {
      let state = globalEnvironmentContextReducer(
        undefined,
        selectEnv({ environment: env1, selectedConnections: [k8sConn] }),
      );
      state = globalEnvironmentContextReducer(
        state,
        selectConnection({ env: env1, connection: k8sConn2 }),
      );
      expect(state.selectedEnvs[env1.id].selectedConnections).toEqual([k8sConn, k8sConn2]);
    });
  });

  describe('unselectConnection', () => {
    it('removes a connection by id', () => {
      let state = globalEnvironmentContextReducer(
        undefined,
        selectEnv({ environment: env1, selectedConnections: [k8sConn, k8sConn2] }),
      );
      state = globalEnvironmentContextReducer(
        state,
        unselectConnection({ envId: env1.id, connectionId: k8sConn.id }),
      );
      expect(state.selectedEnvs[env1.id].selectedConnections).toEqual([k8sConn2]);
    });

    it('leaves the env entry behind even when the last connection is removed', () => {
      let state = globalEnvironmentContextReducer(
        undefined,
        selectEnv({ environment: env1, selectedConnections: [k8sConn] }),
      );
      state = globalEnvironmentContextReducer(
        state,
        unselectConnection({ envId: env1.id, connectionId: k8sConn.id }),
      );
      expect(state.selectedEnvs[env1.id].selectedConnections).toEqual([]);
    });
  });

  describe('toggleConnection thunk', () => {
    it('dispatches selectConnection when not previously selected', () => {
      const dispatch = vi.fn();
      const getState = () => ({ globalEnvironmentContext: { selectedEnvs: {} } }) as any;
      toggleConnection(env1, k8sConn)(dispatch, getState);
      expect(dispatch).toHaveBeenCalledWith(selectConnection({ env: env1, connection: k8sConn }));
    });

    it('dispatches unselectConnection when already selected', () => {
      const dispatch = vi.fn();
      const getState = () =>
        ({
          globalEnvironmentContext: {
            selectedEnvs: {
              [env1.id]: { ...env1, selectedConnections: [k8sConn] },
            },
          },
        }) as any;
      toggleConnection(env1, k8sConn)(dispatch, getState);
      expect(dispatch).toHaveBeenCalledWith(
        unselectConnection({ envId: env1.id, connectionId: k8sConn.id }),
      );
    });
  });

  describe('toggleEnvSelection thunk', () => {
    it('dispatches selectEnv when env was not selected', () => {
      const dispatch = vi.fn();
      const getState = () => ({ globalEnvironmentContext: { selectedEnvs: {} } }) as any;
      toggleEnvSelection(env1, [k8sConn])(dispatch, getState);
      expect(dispatch).toHaveBeenCalledWith(
        selectEnv({ environment: env1, selectedConnections: [k8sConn] }),
      );
    });

    it('dispatches unselectEnv when env was selected', () => {
      const dispatch = vi.fn();
      const getState = () =>
        ({
          globalEnvironmentContext: {
            selectedEnvs: {
              [env1.id]: { ...env1, selectedConnections: [] },
            },
          },
        }) as any;
      toggleEnvSelection(env1, [])(dispatch, getState);
      expect(dispatch).toHaveBeenCalledWith(unselectEnv({ envId: env1.id }));
    });
  });

  describe('selectors', () => {
    const stateWith = (envs: Record<string, any>) =>
      ({ globalEnvironmentContext: { selectedEnvs: envs } }) as any;

    it('selectIsEnvSelected returns true when env is in state', () => {
      expect(
        selectIsEnvSelected(
          stateWith({ [env1.id]: { ...env1, selectedConnections: [] } }),
          env1.id,
        ),
      ).toBe(true);
    });

    it('selectIsEnvSelected returns false when env is missing', () => {
      expect(selectIsEnvSelected(stateWith({}), env1.id)).toBe(false);
    });

    it('selectIsConnectionSelected returns false when env is not selected', () => {
      expect(selectIsConnectionSelected(stateWith({}), env1.id, k8sConn.id)).toBe(false);
    });

    it('selectIsConnectionSelected returns true for selected connection in selected env', () => {
      const state = stateWith({
        [env1.id]: { ...env1, selectedConnections: [k8sConn] },
      });
      expect(selectIsConnectionSelected(state, env1.id, k8sConn.id)).toBe(true);
    });

    it('selectIsConnectionSelected returns false when connection is not selected', () => {
      const state = stateWith({
        [env1.id]: { ...env1, selectedConnections: [k8sConn] },
      });
      expect(selectIsConnectionSelected(state, env1.id, 'other-conn-id')).toBe(false);
    });

    it('selectSelectedEnvs returns the selectedEnvs map', () => {
      const envs = { [env1.id]: { ...env1, selectedConnections: [] } };
      expect(selectSelectedEnvs(stateWith(envs))).toBe(envs);
    });

    it('selectSelectedConnections returns connections for a selected env', () => {
      const state = stateWith({
        [env1.id]: { ...env1, selectedConnections: [k8sConn, dockerConn] },
      });
      expect(selectSelectedConnections(state, env1.id)).toEqual([k8sConn, dockerConn]);
    });

    it('selectSelectedConnections returns [] for unknown env', () => {
      expect(selectSelectedConnections(stateWith({}), 'missing')).toEqual([]);
    });

    it('selectAllSelectedConnections aggregates from every env', () => {
      const state = stateWith({
        [env1.id]: { ...env1, selectedConnections: [k8sConn] },
        [env2.id]: { ...env2, selectedConnections: [dockerConn, k8sConn2] },
      });
      const all = selectAllSelectedConnections(state);
      expect(all).toHaveLength(3);
      expect(all).toEqual(expect.arrayContaining([k8sConn, k8sConn2, dockerConn]));
    });

    it('selectAllSelectedConnections returns [] for empty state', () => {
      expect(selectAllSelectedConnections(stateWith({}))).toEqual([]);
    });

    it('selectSelectedK8sConnections returns only kubernetes connections for an env', () => {
      const state = stateWith({
        [env1.id]: { ...env1, selectedConnections: [k8sConn, dockerConn] },
      });
      expect(selectSelectedK8sConnections(state, env1.id)).toEqual([k8sConn]);
    });

    it('selectAllSelectedK8sConnections filters out non-kubernetes connections', () => {
      const state = stateWith({
        [env1.id]: { ...env1, selectedConnections: [k8sConn, dockerConn] },
        [env2.id]: { ...env2, selectedConnections: [k8sConn2] },
      });
      const all = selectAllSelectedK8sConnections(state);
      expect(all).toHaveLength(2);
      expect(all).toEqual(expect.arrayContaining([k8sConn, k8sConn2]));
      expect(all).not.toContain(dockerConn);
    });
  });

  it('action types have the expected slice prefix', () => {
    expect(selectEnv({ environment: env1 }).type).toBe('globalEnvironmentContext/selectEnv');
    expect(unselectEnv({ envId: 'x' }).type).toBe('globalEnvironmentContext/unselectEnv');
    expect(selectConnection({ env: env1, connection: k8sConn }).type).toBe(
      'globalEnvironmentContext/selectConnection',
    );
    expect(unselectConnection({ envId: 'x', connectionId: 'y' }).type).toBe(
      'globalEnvironmentContext/unselectConnection',
    );
  });
});
