import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const { getUserPrefQuerySpy, getMeshQuerySpy, canSpy, USER_PREF_RESULT, MESH_RESULT } = vi.hoisted(
  () => ({
    getUserPrefQuerySpy: vi.fn(),
    getMeshQuerySpy: vi.fn(),
    canSpy: vi.fn(() => true),
    USER_PREF_RESULT: { data: {}, isSuccess: true },
    MESH_RESULT: { data: [], isSuccess: true, isError: false, error: undefined },
  }),
);

vi.mock('react-redux', () => ({
  useDispatch: () => vi.fn(),
  useSelector: (fn: any) => fn({ ui: { selectedK8sContexts: [], k8sConfig: {} } }),
}));

vi.mock('@/store/slices/mesheryUi', () => ({ updateProgress: vi.fn() }));
vi.mock('@/store/slices/prefTest', () => ({ updateLoadTest: vi.fn() }));

vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetUserPrefWithContextQuery: (args: unknown, options: unknown) => {
    getUserPrefQuerySpy(args, options);
    return USER_PREF_RESULT;
  },
}));

vi.mock('@/rtk-query/performance-profile', () => ({
  useSavePerformanceProfileMutation: () => [vi.fn()],
}));

vi.mock('@/rtk-query/mesh', () => ({
  useGetMeshQuery: (args: unknown, options: unknown) => {
    getMeshQuerySpy(args, options);
    return MESH_RESULT;
  },
}));

vi.mock('@/api/meshmodel', () => ({
  getMeshModels: vi.fn().mockResolvedValue({ models: [] }),
}));

vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({ id: 'id', function: 'fn' }) }),
}));

vi.mock('@/components/general/error-404/index', () => ({
  default: () => <div data-testid="default-error">no access</div>,
}));

vi.mock('@sistent/sistent', () => ({
  useHasPermission: (key: any) => canSpy(key),
  ModalBody: ({ children }: any) => <div>{children}</div>,
  ModalFooter: ({ children }: any) => <div>{children}</div>,
  NoSsr: ({ children }: any) => <>{children}</>,
}));

vi.mock('./PerformanceForm', () => ({ default: () => null }));
vi.mock('./PerformanceFormActions', () => ({ default: () => null }));
vi.mock('./PerformanceTestResults', () => ({ default: () => null }));
vi.mock('../load-test-timer-dialog', () => ({ default: () => null }));
vi.mock('./style', () => ({ CenterTimer: ({ children }: any) => <div>{children}</div> }));

import { getMeshModels } from '@/api/meshmodel';
import MesheryPerformanceComponent from './index';

describe('MesheryPerformanceComponent permission and query skipping', () => {
  beforeEach(() => {
    canSpy.mockReturnValue(true);
    getUserPrefQuerySpy.mockClear();
    getMeshQuerySpy.mockClear();
    vi.mocked(getMeshModels).mockClear();
  });

  it('renders 404 fallback and skips queries and background mesh model fetch when permission is denied', () => {
    canSpy.mockReturnValue(false);
    render(<MesheryPerformanceComponent />);

    expect(screen.getByTestId('default-error')).toBeInTheDocument();
    expect(getUserPrefQuerySpy).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ skip: true }),
    );
    expect(getMeshQuerySpy).toHaveBeenCalledWith(
      undefined,
      expect.objectContaining({ skip: true }),
    );
    expect(getMeshModels).not.toHaveBeenCalled();
  });

  it('executes queries without skip when permission is granted', () => {
    canSpy.mockReturnValue(true);
    render(<MesheryPerformanceComponent />);

    expect(getUserPrefQuerySpy).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ skip: false }),
    );
    expect(getMeshQuerySpy).toHaveBeenCalledWith(
      undefined,
      expect.objectContaining({ skip: false }),
    );
    expect(getMeshModels).toHaveBeenCalled();
  });
});
