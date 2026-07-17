import React from 'react';
import { render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { KubernetesSubscription } from './AppComponents';

const dispatchMock = vi.fn();
const setAppStateMock = vi.fn();

// Return value of useGetConnectionsQuery, swapped per test.
let connectionsResult: { data?: { connections?: unknown[] } } = { data: undefined };

// Redux boots with ['all'] selected; swapped per test to exercise the sync path.
let reduxSelectedK8sContexts: string[] = ['all'];

vi.mock('react-redux', () => ({
  useDispatch: () => dispatchMock,
  useSelector: (selector: (state: Record<string, unknown>) => unknown) =>
    selector({ ui: { extensionType: '', selectedK8sContexts: reduxSelectedK8sContexts } }),
}));

// The k8s context list is now driven by the connections REST API
// (kind=kubernetes) instead of the subscribeK8sContext GraphQL subscription.
vi.mock('@/rtk-query/connection', () => ({
  useGetConnectionsQuery: () => connectionsResult,
}));

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

vi.mock('@/utils/Enum', () => ({
  CONNECTION_KINDS: { KUBERNETES: 'kubernetes' },
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateK8SConfig: (payload: Record<string, unknown>) => ({
    type: 'core/updateK8SConfig',
    payload,
  }),
  setK8sContexts: (payload: Record<string, unknown>) => ({
    type: 'core/setK8sContexts',
    payload,
  }),
}));

vi.mock('@sistent/sistent', () => ({
  FavoriteIcon: () => null,
  Hidden: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  Typography: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  useTheme: () => ({
    palette: {
      mode: 'light',
      text: { default: '#000', disabled: '#666' },
      background: { brand: { default: '#00b39f' } },
    },
  }),
}));

vi.mock('./layout/Navigator/Navigator', () => ({
  default: () => null,
}));

vi.mock('../themes/App.styles', () => ({
  StyledDrawer: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  StyledFooterBody: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  StyledFooterText: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

const twoConnections = {
  data: {
    connections: [
      {
        id: 'conn-1',
        metadata: { id: 'ctx-1', name: 'prod-us', kubernetes_server_id: 's1' },
      },
      {
        id: 'conn-2',
        metadata: { id: 'ctx-2', name: 'prod-eu' },
      },
    ],
  },
};

describe('KubernetesSubscription', () => {
  beforeEach(() => {
    dispatchMock.mockReset();
    setAppStateMock.mockReset();
    connectionsResult = { data: undefined };
    reduxSelectedK8sContexts = ['all'];
    window.sessionStorage.clear();
  });

  it('maps kubernetes connections into contexts and connection config', () => {
    connectionsResult = twoConnections;

    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    expect(setAppStateMock).toHaveBeenCalledWith({
      k8sContexts: expect.objectContaining({
        totalCount: 2,
        contexts: [
          expect.objectContaining({ id: 'ctx-1', connectionId: 'conn-1', name: 'prod-us' }),
          expect.objectContaining({ id: 'ctx-2', connectionId: 'conn-2', name: 'prod-eu' }),
        ],
      }),
      activeK8sContexts: ['ctx-1', 'ctx-2', 'all'],
    });
    expect(dispatchMock).toHaveBeenCalledWith(
      expect.objectContaining({
        type: 'core/updateK8SConfig',
        payload: {
          k8sConfig: [
            expect.objectContaining({ id: 'ctx-1', connectionId: 'conn-1' }),
            expect.objectContaining({ id: 'ctx-2', connectionId: 'conn-2' }),
          ],
        },
      }),
    );
  });

  it('stores an empty list when there are no kubernetes connections', () => {
    connectionsResult = { data: { connections: [] } };

    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    expect(dispatchMock).toHaveBeenCalledWith(
      expect.objectContaining({
        type: 'core/updateK8SConfig',
        payload: { k8sConfig: [] },
      }),
    );
  });

  it('honors a persisted partial selection instead of force-selecting every context', () => {
    connectionsResult = twoConnections;
    window.sessionStorage.setItem('selectedK8sContexts', JSON.stringify(['ctx-2']));

    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    expect(setAppStateMock).toHaveBeenCalledWith(
      expect.objectContaining({ activeK8sContexts: ['ctx-2'] }),
    );
    // Redux boots with ['all'], so the restored selection must be synced.
    expect(dispatchMock.mock.calls.map(([action]) => action)).toContainEqual(
      expect.objectContaining({
        type: 'core/setK8sContexts',
        payload: { selectedK8sContexts: ['ctx-2'] },
      }),
    );
  });

  it('honors a persisted explicit empty selection', () => {
    connectionsResult = twoConnections;
    window.sessionStorage.setItem('selectedK8sContexts', JSON.stringify([]));

    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    expect(setAppStateMock).toHaveBeenCalledWith(
      expect.objectContaining({ activeK8sContexts: [] }),
    );
  });

  it('selects every context when the persisted selection is ["all"]', () => {
    connectionsResult = twoConnections;
    window.sessionStorage.setItem('selectedK8sContexts', JSON.stringify(['all']));

    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    expect(setAppStateMock).toHaveBeenCalledWith(
      expect.objectContaining({ activeK8sContexts: ['ctx-1', 'ctx-2', 'all'] }),
    );
    // Redux already holds ['all'] — no redundant sync dispatch.
    expect(dispatchMock.mock.calls.map(([action]) => action)).not.toContainEqual(
      expect.objectContaining({ type: 'core/setK8sContexts' }),
    );
  });

  it('falls back to selecting all contexts when every persisted id is stale', () => {
    connectionsResult = twoConnections;
    window.sessionStorage.setItem('selectedK8sContexts', JSON.stringify(['ctx-gone']));

    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    expect(setAppStateMock).toHaveBeenCalledWith(
      expect.objectContaining({ activeK8sContexts: ['ctx-1', 'ctx-2', 'all'] }),
    );
  });

  it('restores the implicit "all" when the persisted ids cover every context', () => {
    connectionsResult = twoConnections;
    window.sessionStorage.setItem('selectedK8sContexts', JSON.stringify(['ctx-1', 'ctx-2']));

    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    expect(setAppStateMock).toHaveBeenCalledWith(
      expect.objectContaining({ activeK8sContexts: ['ctx-1', 'ctx-2', 'all'] }),
    );
    // Storage is resynced to the resolved selection so it does not keep a
    // stale variant of the same selection.
    expect(JSON.parse(window.sessionStorage.getItem('selectedK8sContexts') || 'null')).toEqual([
      'all',
    ]);
  });

  it('does not touch the selection or storage while connections are still loading', () => {
    connectionsResult = { data: undefined };
    window.sessionStorage.setItem('selectedK8sContexts', JSON.stringify(['ctx-2']));

    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    expect(setAppStateMock).not.toHaveBeenCalled();
    expect(dispatchMock).not.toHaveBeenCalled();
    // The persisted selection survives untouched until data arrives.
    expect(JSON.parse(window.sessionStorage.getItem('selectedK8sContexts') || 'null')).toEqual([
      'ctx-2',
    ]);
  });

  it('leaves redux and storage untouched when no contexts exist', () => {
    connectionsResult = { data: { connections: [] } };
    window.sessionStorage.setItem('selectedK8sContexts', JSON.stringify(['ctx-1']));

    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    // k8sConfig still syncs (empty), but the selection is not rewritten.
    expect(dispatchMock.mock.calls.map(([action]) => action)).not.toContainEqual(
      expect.objectContaining({ type: 'core/setK8sContexts' }),
    );
    expect(JSON.parse(window.sessionStorage.getItem('selectedK8sContexts') || 'null')).toEqual([
      'ctx-1',
    ]);
  });

  it('drops stale ids but keeps the surviving persisted selection', () => {
    connectionsResult = twoConnections;
    window.sessionStorage.setItem('selectedK8sContexts', JSON.stringify(['ctx-1', 'ctx-gone']));

    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    expect(setAppStateMock).toHaveBeenCalledWith(
      expect.objectContaining({ activeK8sContexts: ['ctx-1'] }),
    );
  });

  it('ignores corrupt persisted values and defaults to all contexts', () => {
    connectionsResult = twoConnections;
    window.sessionStorage.setItem('selectedK8sContexts', 'not-json');

    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    expect(setAppStateMock).toHaveBeenCalledWith(
      expect.objectContaining({ activeK8sContexts: ['ctx-1', 'ctx-2', 'all'] }),
    );
  });
});
