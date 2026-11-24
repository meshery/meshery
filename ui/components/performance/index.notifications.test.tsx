import React from 'react';
import { render, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Covers the notification MesheryPerformanceComponent raises when the SMP mesh
// list can't be fetched. The form, its actions, the timer dialog and the
// network-facing modules are stubbed; only the component's own effects run.

// Query results must be referentially stable across renders: the component
// keys effects on them.
const { notify, MESH_QUERY_RESULT, USER_PREF_RESULT, REDUX_STATE } = vi.hoisted(() => ({
  notify: vi.fn(),
  MESH_QUERY_RESULT: {
    data: undefined as unknown,
    isSuccess: false,
    isError: false,
    error: undefined as unknown,
  },
  USER_PREF_RESULT: { data: undefined, isSuccess: false },
  REDUX_STATE: { ui: { selectedK8sContexts: [], k8sConfig: [] } },
}));

vi.mock('@/rtk-query/mesh', () => ({ useGetMeshQuery: () => MESH_QUERY_RESULT }));
vi.mock('@/rtk-query/user', () => ({ useGetUserPrefWithContextQuery: () => USER_PREF_RESULT }));
vi.mock('@/rtk-query/performance-profile', () => ({
  useSavePerformanceProfileMutation: () => [vi.fn()],
}));
vi.mock('../../rtk-query', () => ({ api: { util: { invalidateTags: vi.fn() } } }));
vi.mock('@/api/meshmodel', () => ({ getMeshModels: () => Promise.resolve({ models: [] }) }));
vi.mock('@/graphql/queries/ControlPlanesQuery', () => ({ default: vi.fn() }));
vi.mock('../../utils/hooks/useNotification', () => ({ useNotification: () => ({ notify }) }));
vi.mock('react-redux', () => ({
  useDispatch: () => vi.fn(),
  useSelector: (selector: (state: unknown) => unknown) => selector(REDUX_STATE),
}));
vi.mock('@/store/slices/mesheryUi', () => ({ updateProgress: vi.fn() }));
vi.mock('@/store/slices/prefTest', () => ({ updateLoadTest: vi.fn() }));
vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({ id: 'id', function: 'fn' }) }),
}));
vi.mock('@/components/general/error-404/index', () => ({ default: () => <div>no access</div> }));
vi.mock('../load-test-timer-dialog', () => ({ default: () => null }));
vi.mock('./PerformanceForm', () => ({ default: () => null }));
vi.mock('./PerformanceFormActions', () => ({ default: () => null }));
vi.mock('./PerformanceTestResults', () => ({ default: () => null }));
vi.mock('./style', () => ({ CenterTimer: ({ children }: any) => <div>{children}</div> }));

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@sistent/sistent')>();
  return {
    ...actual,
    ModalBody: ({ children }: any) => <div>{children}</div>,
    ModalFooter: ({ children }: any) => <div>{children}</div>,
    useHasPermission: () => true,
  };
});

import MesheryPerformanceComponent from './index';

const notifiedEventTypes = () =>
  notify.mock.calls.map(([arg]) => arg?.event_type?.type ?? arg?.event_type);

describe('MesheryPerformanceComponent SMP mesh notifications', () => {
  beforeEach(() => {
    notify.mockReset();
    Object.assign(MESH_QUERY_RESULT, {
      data: undefined,
      isSuccess: false,
      isError: false,
      error: undefined,
    });
  });

  it('shows an error notification when the SMP mesh list fails to load', async () => {
    Object.assign(MESH_QUERY_RESULT, {
      isError: true,
      error: { status: 500, data: 'mesh list unavailable' },
    });

    render(<MesheryPerformanceComponent />);

    await waitFor(() => expect(notify).toHaveBeenCalledTimes(1));
    expect(notifiedEventTypes()).toEqual(['error']);
    expect(notify.mock.calls[0][0].message).toContain('mesh list unavailable');
  });

  it('does not notify when the SMP mesh list loads', async () => {
    Object.assign(MESH_QUERY_RESULT, {
      data: { availableMeshes: ['istio'] },
      isSuccess: true,
    });

    render(<MesheryPerformanceComponent />);

    // Let mount effects (including the async mesh model fetch) settle.
    await new Promise((resolve) => setTimeout(resolve, 0));
    expect(notify).not.toHaveBeenCalled();
  });
});
