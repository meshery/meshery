import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const {
  getUserPrefQuerySpy,
  getMeshQuerySpy,
  canSpy,
  saveProfileSpy,
  formPropsSpy,
  actionsPropsSpy,
  USER_PREF_RESULT,
  MESH_RESULT,
} = vi.hoisted(() => ({
  getUserPrefQuerySpy: vi.fn(),
  getMeshQuerySpy: vi.fn(),
  canSpy: vi.fn(() => true),
  saveProfileSpy: vi.fn(),
  formPropsSpy: vi.fn(),
  actionsPropsSpy: vi.fn(),
  USER_PREF_RESULT: { data: {}, isSuccess: true },
  MESH_RESULT: { data: [], isSuccess: true, isError: false, error: undefined },
}));

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
  useSavePerformanceProfileMutation: () => [saveProfileSpy],
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
  Link: ({ children, ...rest }: any) => <a {...rest}>{children}</a>,
  ModalBody: ({ children }: any) => <div>{children}</div>,
  ModalFooter: ({ children }: any) => <div>{children}</div>,
  NoSsr: ({ children }: any) => <>{children}</>,
}));

vi.mock('./PerformanceForm', () => ({
  default: (props: any) => {
    formPropsSpy(props);
    return null;
  },
}));
vi.mock('./PerformanceFormActions', () => ({
  default: (props: any) => {
    actionsPropsSpy(props);
    return null;
  },
}));
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

describe('MesheryPerformanceComponent load parameter validation', () => {
  const lastFormProps = () => formPropsSpy.mock.calls.at(-1)![0];
  const lastActionsProps = () => actionsPropsSpy.mock.calls.at(-1)![0];

  beforeEach(() => {
    canSpy.mockReturnValue(true);
    saveProfileSpy.mockReset();
    formPropsSpy.mockClear();
    actionsPropsSpy.mockClear();
  });

  it.each([
    ['negative concurrent requests', { c: '-10', qps: '10' }, { cError: true, qpsError: false }],
    ['zero concurrent requests', { c: '0', qps: '10' }, { cError: true, qpsError: false }],
    ['negative QPS', { c: '5', qps: '-5' }, { cError: false, qpsError: true }],
  ])('flags %s and blocks Save Profile and Run Test', (_name, values, expected) => {
    render(
      <MesheryPerformanceComponent
        url="http://example.com"
        performanceProfileID="p1"
        {...values}
      />,
    );

    expect(!!lastFormProps().cError).toBe(expected.cError);
    expect(!!lastFormProps().qpsError).toBe(expected.qpsError);
    expect(lastActionsProps().disableTest).toBe(true);

    lastActionsProps().onSaveProfile();
    lastActionsProps().onRunTest();
    expect(saveProfileSpy).not.toHaveBeenCalled();
  });

  it('keeps the form valid when QPS is 0 (unbounded)', () => {
    render(
      <MesheryPerformanceComponent
        url="http://example.com"
        performanceProfileID="p1"
        c="1"
        qps="0"
      />,
    );

    expect(lastFormProps().cError).toBe('');
    expect(lastFormProps().qpsError).toBe('');
    expect(lastActionsProps().disableTest).toBe(false);
  });

  it('opens a new profile form with a valid default for concurrent requests', () => {
    render(<MesheryPerformanceComponent url="http://example.com" />);

    expect(lastFormProps().cError).toBe('');
    expect(lastActionsProps().disableTest).toBe(false);
  });
});
