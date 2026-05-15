import React from 'react';
import { act, render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { KubernetesSubscription } from './AppComponents';

const dispatchMock = vi.fn();
const setAppStateMock = vi.fn();
const disposeMock = vi.fn();
const sseSubscribeMock = vi.fn();
let subscriptionCallback: ((result: Record<string, unknown>) => void) | undefined;

vi.mock('react-redux', () => ({
  useDispatch: () => dispatchMock,
  useSelector: (selector: (state: Record<string, unknown>) => unknown) =>
    selector({ ui: { extensionType: '' } }),
}));

vi.mock('@/lib/sseClient', () => ({
  sseSubscribe: (opts: {
    path: string;
    onMessage: (result: Record<string, unknown>) => void;
    subscriptionName?: string;
  }) => {
    subscriptionCallback = opts.onMessage;
    return sseSubscribeMock(opts);
  },
}));

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

vi.mock('@/utils/permission_constants', () => ({
  keys: {
    VIEW_ALL_KUBERNETES_CLUSTERS: { action: 'view', subject: 'k8s' },
  },
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
    disposeMock.mockReset();
    sseSubscribeMock.mockReset();
    subscriptionCallback = undefined;

    sseSubscribeMock.mockReturnValue({ dispose: disposeMock, rebind: vi.fn() });
  });

  it('normalizes subscription payloads before storing contexts and connection config', () => {
    const { unmount } = render(<KubernetesSubscription setAppState={setAppStateMock} />);

    expect(sseSubscribeMock).toHaveBeenCalledWith(
      expect.objectContaining({
        path: '/api/system/kubernetes/contexts/stream',
        subscriptionName: 'K8sContext',
        onMessage: expect.any(Function),
      }),
    );
    expect(setAppStateMock).toHaveBeenCalledWith({
      disposeK8sContextSubscription: expect.any(Function),
    });

    act(() => {
      subscriptionCallback?.({
        k8sContext: {
          total_count: 2,
          contexts: [
            {
              id: 'ctx-1',
              name: 'prod-us',
              connection_id: 'conn-1',
              created_by: 'meshery',
            },
            {
              id: 'ctx-2',
              name: 'prod-eu',
              connection_id: 'conn-2',
              created_by: 'meshery',
            },
          ],
        },
      });
    });

    expect(setAppStateMock).toHaveBeenCalledWith({
      k8sContexts: expect.objectContaining({
        totalCount: 2,
        contexts: [
          expect.objectContaining({
            id: 'ctx-1',
            connectionId: 'conn-1',
            createdBy: 'meshery',
          }),
          expect.objectContaining({
            id: 'ctx-2',
            connectionId: 'conn-2',
            createdBy: 'meshery',
          }),
        ],
      }),
      activeK8sContexts: ['ctx-1', 'ctx-2', 'all'],
    });
    expect(dispatchMock).toHaveBeenCalledWith(
      expect.objectContaining({
        type: 'core/updateK8SConfig',
        payload: {
          k8sConfig: [
            expect.objectContaining({
              id: 'ctx-1',
              connectionId: 'conn-1',
            }),
            expect.objectContaining({
              id: 'ctx-2',
              connectionId: 'conn-2',
            }),
          ],
        },
      }),
    );

    unmount();
    expect(disposeMock).toHaveBeenCalledTimes(1);
  });

  it('passes through an already-camelCased subscription payload (GraphQL alias path)', () => {
    render(<KubernetesSubscription setAppState={setAppStateMock} />);

    act(() => {
      subscriptionCallback?.({
        k8sContext: {
          totalCount: 4,
          contexts: [
            { id: 'ctx-1', name: 'a', connectionId: 'conn-1', createdBy: 'meshery' },
            { id: 'ctx-2', name: 'b', connectionId: 'conn-2', createdBy: 'meshery' },
          ],
        },
      });
    });

    expect(setAppStateMock).toHaveBeenCalledWith(
      expect.objectContaining({
        k8sContexts: expect.objectContaining({
          totalCount: 4,
          contexts: expect.arrayContaining([
            expect.objectContaining({ id: 'ctx-1', connectionId: 'conn-1' }),
            expect.objectContaining({ id: 'ctx-2', connectionId: 'conn-2' }),
          ]),
        }),
      }),
    );
  });
});
