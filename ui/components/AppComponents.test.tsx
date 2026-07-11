import React from 'react';
import { render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { KubernetesSubscription } from './AppComponents';

const dispatchMock = vi.fn();
const setAppStateMock = vi.fn();

// Return value of useGetConnectionsQuery, swapped per test.
let connectionsResult: { data?: { connections?: unknown[] } } = { data: undefined };

vi.mock('react-redux', () => ({
  useDispatch: () => dispatchMock,
  useSelector: (selector: (state: Record<string, unknown>) => unknown) =>
    selector({ ui: { extensionType: '' } }),
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

describe('KubernetesSubscription', () => {
  beforeEach(() => {
    dispatchMock.mockReset();
    setAppStateMock.mockReset();
    connectionsResult = { data: undefined };
  });

  it('maps kubernetes connections into contexts and connection config', () => {
    connectionsResult = {
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
});
