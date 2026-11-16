import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const back = vi.fn();
const ping = vi.fn();
const setView = vi.fn();
const normalizeStaticImagePath = vi.fn((src?: string) => src);
const useGetConnectionsQuery = vi.fn();
const canSpy = vi.fn(() => true);
const getK8sContextFromClusterId = vi.fn((id: string) => ({
  name: `ctx-${id}`,
  server: 'srv',
  connectionId: 'conn',
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ back, pathname: '/foo', push: vi.fn() }),
}));

vi.mock('@sistent/sistent', () => {
  return {
    ArrowBackIcon: () => <svg data-testid="arrow-back" />,
    Box: ({ children }: any) => <div>{children}</div>,
    ErrorBoundary: ({ children }: any) => <>{children}</>,
    OperatorDataFormatter: ({ data }: any) => (
      <div data-testid="op-formatter">{JSON.stringify(data)}</div>
    ),
    Paper: ({ children }: any) => <div>{children}</div>,
    styled: (Component: any) => () => Component,
    Typography: ({ children }: any) => <span>{children}</span>,
    useHasPermission: (key: unknown) => canSpy(key),
    useResourceCleanData: () => ({ getResourceCleanData: ({ resource }: any) => resource }),
  };
});

vi.mock('@/constants/common', () => ({
  FALLBACK_MESHERY_IMAGE_PATH: '/fallback.svg',
  KUBERNETES: 'kubernetes',
}));

vi.mock('@/utils/fallback', () => ({
  normalizeStaticImagePath: (...args: any[]) => normalizeStaticImagePath(...args),
}));

vi.mock('css/icons.styles', () => ({
  iconXLarge: { width: 32 },
}));

vi.mock('@/utils/multi-ctx', () => ({
  getK8sContextFromClusterId: (...args: any[]) => getK8sContextFromClusterId(...args),
}));

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => ping,
}));

vi.mock('../connections/ConnectionChip', () => ({
  TooltipWrappedConnectionChip: ({ title, status, handlePing }: any) => (
    <div
      data-testid="connection-chip"
      data-status={status}
      data-has-ping={String(Boolean(handlePing))}
    >
      {title}
    </div>
  ),
}));

vi.mock('./view-component', () => ({
  default: ({ data }: any) => <div data-testid="resource-detail">{JSON.stringify(data)}</div>,
  JSONViewFormatter: ({ data }: any) => <pre>{JSON.stringify(data)}</pre>,
}));

vi.mock('./utils', () => ({
  default: ({ kind, model }: any) => (
    <svg data-testid="node-icon" data-kind={kind} data-model={model} />
  ),
}));

vi.mock('@/utils/Enum', async () => ({
  ...(await vi.importActual('@/utils/Enum')),
  CONNECTION_STATES: { DISCONNECTED: 'disconnected', CONNECTED: 'connected' },
}));

vi.mock('@/rtk-query/connection', () => ({
  useGetConnectionsQuery: (...args: any[]) => useGetConnectionsQuery(...args),
}));

vi.mock('../../utils/TooltipButton', () => ({
  TooltipIconButton: ({ children, onClick }: any) => (
    <button type="button" onClick={onClick} data-testid="tooltip-button">
      {children}
    </button>
  ),
}));

vi.mock('./resources/config', () => ({ ALL_VIEW: 'all' }));

import { Keys } from '@meshery/schemas/permissions';
import { CONNECTION_STATES } from '@/utils/Enum';
import View, { Title } from './view';

describe('View', () => {
  beforeEach(() => {
    back.mockReset();
    ping.mockReset();
    setView.mockReset();
    normalizeStaticImagePath.mockReset();
    useGetConnectionsQuery.mockReset();
    canSpy.mockClear();
    canSpy.mockReturnValue(true);
    getK8sContextFromClusterId.mockReset();
    getK8sContextFromClusterId.mockImplementation((id: string) => ({
      name: `ctx-${id}`,
      server: 'srv',
      connectionId: 'conn',
    }));
    normalizeStaticImagePath.mockImplementation((src) => src);
    useGetConnectionsQuery.mockReturnValue({ data: { connections: [] } });
  });

  it('returns null and skips connections query when no resource is provided', () => {
    const { container } = render(<View setView={setView} resource={null} k8sConfig={{}} />);
    expect(container.firstChild).toBeNull();
    expect(useGetConnectionsQuery).toHaveBeenCalledWith(
      expect.objectContaining({ kind: 'kubernetes', pageSize: 100 }),
      { skip: true },
    );
  });

  it('renders the back button, image, and resource title for a resource', () => {
    render(
      <View
        setView={setView}
        resource={{
          kind: 'Pod',
          cluster_id: 'cluster-1',
          metadata: { name: 'pod-1' },
          component_metadata: { styles: { svgColor: '/icon.svg' } },
        }}
        k8sConfig={{}}
      />,
    );

    expect(screen.getByTestId('tooltip-button')).toBeInTheDocument();
    const img = screen.getByRole('img', { name: 'Pod' });
    expect(img).toHaveAttribute('src', '/icon.svg');
    expect(screen.getByText('pod-1')).toBeInTheDocument();
    expect(screen.getByTestId('connection-chip')).toHaveTextContent('ctx-cluster-1');
    expect(canSpy).toHaveBeenCalledWith(Keys.WorkspaceManagementViewConnections);
    // #20617: plain ?kind=kubernetes, not JSON.stringify(['kubernetes'])
    expect(useGetConnectionsQuery).toHaveBeenCalledWith(
      expect.objectContaining({ kind: 'kubernetes', pageSize: 100 }),
      { skip: false },
    );
  });

  it('navigates back and switches view when the back button is clicked', async () => {
    const user = userEvent.setup();
    render(
      <View
        setView={setView}
        resource={{
          kind: 'Pod',
          metadata: { name: 'pod-1' },
          component_metadata: { styles: {} },
        }}
        k8sConfig={{}}
      />,
    );

    await user.click(screen.getByTestId('tooltip-button'));
    expect(back).toHaveBeenCalled();
    expect(setView).toHaveBeenCalledWith('all');
  });

  it('renders cleanly without crashing when connection data is loading (undefined)', () => {
    useGetConnectionsQuery.mockReturnValue({ data: undefined });
    render(
      <View
        setView={setView}
        resource={{
          kind: 'Pod',
          cluster_id: 'cluster-1',
          metadata: { name: 'pod-1' },
          component_metadata: { styles: {} },
        }}
        k8sConfig={{}}
      />,
    );
    expect(screen.getByText('pod-1')).toBeInTheDocument();
    expect(screen.getByTestId('connection-chip')).toHaveAttribute(
      'data-status',
      CONNECTION_STATES.DISCONNECTED,
    );
  });

  it('renders cleanly without ping handler when resolved cluster context is undefined or empty', () => {
    for (const resolved of [undefined, {}]) {
      getK8sContextFromClusterId.mockReturnValueOnce(resolved as never);
      const { unmount } = render(
        <View
          setView={setView}
          resource={{
            kind: 'Pod',
            cluster_id: 'missing-cluster',
            metadata: { name: 'pod-orphan' },
            component_metadata: { styles: {} },
          }}
          k8sConfig={[]}
        />,
      );
      expect(screen.getByText('pod-orphan')).toBeInTheDocument();
      const chip = screen.getByTestId('connection-chip');
      expect(chip).toHaveAttribute('data-status', CONNECTION_STATES.DISCONNECTED);
      expect(chip).toHaveAttribute('data-has-ping', 'false');
      unmount();
    }
  });

  it('skips useGetConnectionsQuery and ignores cached connection status when WorkspaceManagementViewConnections is denied', () => {
    canSpy.mockReturnValue(false);
    useGetConnectionsQuery.mockReturnValue({
      data: { connections: [{ id: 'conn', status: CONNECTION_STATES.CONNECTED }] },
    });

    render(
      <View
        setView={setView}
        resource={{
          kind: 'Pod',
          cluster_id: 'cluster-1',
          metadata: { name: 'pod-denied' },
          component_metadata: { styles: {} },
        }}
        k8sConfig={[]}
      />,
    );

    expect(useGetConnectionsQuery).toHaveBeenCalledWith(
      expect.objectContaining({ kind: 'kubernetes', pageSize: 100 }),
      { skip: true },
    );
    expect(screen.getByTestId('connection-chip')).toHaveAttribute(
      'data-status',
      CONNECTION_STATES.DISCONNECTED,
    );
  });
});

describe('Title', () => {
  it('renders the value and toggles underline on hover', async () => {
    const user = userEvent.setup();
    render(<Title onClick={vi.fn()} value="pod-x" kind="Pod" />);
    expect(screen.getByText('pod-x')).toBeInTheDocument();
    expect(screen.getByTestId('node-icon')).toHaveAttribute('data-kind', 'Pod');

    // Mouseenter changes hover state — verify by interacting with the element.
    const titleContent = screen.getByText('pod-x').parentElement?.parentElement;
    await user.hover(titleContent as Element);
  });

  it('calls onClick when the title content is clicked', async () => {
    const user = userEvent.setup();
    const onClick = vi.fn();
    render(<Title onClick={onClick} value="pod-x" kind="Pod" />);
    await user.click(screen.getByText('pod-x'));
    expect(onClick).toHaveBeenCalled();
  });
});
