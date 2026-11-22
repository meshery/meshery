import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi, beforeEach } from 'vitest';

let connectionsQueryReturn: {
  data?: { connections?: Array<{ status?: string }> };
  isFetching?: boolean;
  isLoading?: boolean;
  isError?: boolean;
} = {
  data: { connections: [] },
  isFetching: false,
  isLoading: false,
  isError: false,
};
const canSpy = vi.fn((_key?: unknown) => true);
const bbChartSpy = vi.fn();
const routerPushSpy = vi.fn();
const useGetConnectionsQuery = vi.fn();

vi.mock('billboard.js', () => ({
  donut: () => 'donut',
}));

vi.mock('../../general/BBChart', () => ({
  default: (props: { options: unknown }) => {
    bbChartSpy(props.options);
    return <div data-testid="bb-chart" />;
  },
}));

vi.mock('../../../utils/charts', () => ({
  dataToColors: () => ({}),
  isValidColumnName: (s: unknown) => Boolean(s) && s !== 'invalid',
}));

vi.mock('next/link', () => ({
  default: ({
    children,
    href,
    style,
  }: {
    children: React.ReactNode;
    href: string;
    style?: React.CSSProperties;
  }) => (
    <a href={href} style={style}>
      {children}
    </a>
  ),
}));

vi.mock('@/rtk-query/connection', () => ({
  useGetConnectionsQuery: (...args: unknown[]) => useGetConnectionsQuery(...args),
}));

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => canSpy(...args),
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ push: routerPushSpy }),
}));

vi.mock('./ConnectCluster', () => ({
  default: ({ message }: { message: React.ReactNode }) => (
    <div data-testid="connect-cluster">{message}</div>
  ),
}));

vi.mock('../widgets/WidgetErrorFallback', () => ({
  default: ({ widgetTitle, message }: { widgetTitle: string; message?: string }) => (
    <div data-testid="widget-error-fallback" data-title={widgetTitle}>
      {message}
    </div>
  ),
}));

vi.mock('../style', () => ({
  DashboardSection: ({ children }: { children?: React.ReactNode }) => (
    <section data-testid="dashboard-section">{children}</section>
  ),
  LoadingContainer: ({ children }: { children?: React.ReactNode }) => (
    <div data-testid="loading-container">{children}</div>
  ),
}));

vi.mock('@sistent/sistent', () => ({
  Box: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
  CircularProgress: () => <div data-testid="circular-progress" />,
  InfoTooltip: () => <svg data-testid="info-icon" />,
  KubernetesIcon: () => <svg data-testid="k8s-icon" />,
  Typography: ({ children }: { children?: React.ReactNode }) => <p>{children}</p>,
  useTheme: () => ({
    palette: { mode: 'light', icon: { default: '#000', disabled: '#777' } },
  }),
  useHasPermission: (key: unknown) => canSpy(key),
}));

import { Keys } from '@meshery/schemas/permissions';
import KubernetesConnectionStatsChart from './KubernetesConnectionChart';

describe('KubernetesConnectionStatsChart', () => {
  beforeEach(() => {
    bbChartSpy.mockReset();
    routerPushSpy.mockReset();
    canSpy.mockClear();
    canSpy.mockReturnValue(true);
    connectionsQueryReturn = {
      data: { connections: [] },
      isFetching: false,
      isLoading: false,
      isError: false,
    };
    useGetConnectionsQuery.mockReset();
    useGetConnectionsQuery.mockImplementation(() => connectionsQueryReturn);
  });

  it('shows ConnectCluster fallback when there are no kubernetes connections', () => {
    render(<KubernetesConnectionStatsChart />);
    expect(screen.getByTestId('connect-cluster')).toHaveTextContent(
      'No connections found in your clusters',
    );
    // #20617: plain ?kind=kubernetes, not JSON.stringify(['kubernetes'])
    expect(useGetConnectionsQuery).toHaveBeenCalledWith(
      expect.objectContaining({ kind: 'kubernetes', pageSize: 'all' }),
      { skip: false },
    );
  });

  it('authorizes against WorkspaceManagementViewConnections', () => {
    render(<KubernetesConnectionStatsChart />);
    expect(canSpy).toHaveBeenCalledWith(Keys.WorkspaceManagementViewConnections);
  });

  it('skips useGetConnectionsQuery and renders ConnectCluster fallback when view permission is denied', () => {
    canSpy.mockReturnValue(false);
    connectionsQueryReturn = {
      data: undefined,
      isFetching: false,
      isLoading: false,
      isError: false,
    };
    render(<KubernetesConnectionStatsChart />);
    expect(useGetConnectionsQuery).toHaveBeenCalledWith(
      expect.objectContaining({ kind: 'kubernetes', pageSize: 'all' }),
      { skip: true },
    );
    expect(screen.getByTestId('connect-cluster')).toHaveTextContent(
      'No connections found in your clusters',
    );
    expect(screen.queryByTestId('widget-error-fallback')).not.toBeInTheDocument();
    expect(screen.queryByTestId('bb-chart')).not.toBeInTheDocument();
  });

  it('hides retained connections, errors, and Links when view permission is denied', async () => {
    const user = userEvent.setup();
    connectionsQueryReturn = {
      data: {
        connections: [{ status: 'connected' }],
      },
      isLoading: false,
      isError: false,
    };
    canSpy.mockReturnValue(false);
    render(<KubernetesConnectionStatsChart />);
    expect(useGetConnectionsQuery).toHaveBeenCalledWith(
      expect.objectContaining({ kind: 'kubernetes', pageSize: 'all' }),
      { skip: true },
    );
    expect(screen.getByTestId('connect-cluster')).toBeInTheDocument();
    expect(screen.queryByTestId('bb-chart')).not.toBeInTheDocument();
    expect(screen.queryByRole('link')).not.toBeInTheDocument();
    await user.tab();
    expect(document.activeElement).toBe(document.body);

    connectionsQueryReturn = {
      data: undefined,
      isLoading: false,
      isError: true,
    };
    render(<KubernetesConnectionStatsChart />);
    expect(screen.queryByTestId('widget-error-fallback')).not.toBeInTheDocument();
  });

  it('renders Link and navigates on chart onclick when view permission is granted', () => {
    connectionsQueryReturn = {
      data: {
        connections: [{ status: 'connected' }],
      },
    };
    render(<KubernetesConnectionStatsChart />);
    expect(screen.getByRole('link')).toHaveAttribute('href', '/management/connections');
    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: { onclick: () => void };
    };
    options.data.onclick();
    expect(routerPushSpy).toHaveBeenCalledWith('/management/connections');
  });

  it('renders BBChart with normalized status counts when connections exist', () => {
    connectionsQueryReturn = {
      data: {
        connections: [
          { status: 'connected' },
          { status: 'connected' },
          { status: 'disconnected' },
          { status: 'invalid' },
        ],
      },
    };
    render(<KubernetesConnectionStatsChart />);
    expect(screen.getByTestId('bb-chart')).toBeInTheDocument();
    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: { columns: Array<[string, number]> };
      donut: { title: string };
    };
    const map = Object.fromEntries(options.data.columns);
    expect(map.connected).toBe(2);
    expect(map.disconnected).toBe(1);
    expect(map.invalid).toBeUndefined();
    expect(options.donut.title).toBe('Clusters\n  Status');
  });

  it('omits BBChart and Link when the user lacks permission even with retained data', () => {
    connectionsQueryReturn = {
      data: {
        connections: [{ status: 'connected' }],
      },
    };
    canSpy.mockReturnValue(false);
    render(<KubernetesConnectionStatsChart />);
    expect(screen.queryByRole('link')).not.toBeInTheDocument();
    expect(screen.queryByTestId('bb-chart')).not.toBeInTheDocument();
    expect(screen.getByTestId('connect-cluster')).toBeInTheDocument();
  });

  it('renders the Kubernetes header regardless of data state', () => {
    render(<KubernetesConnectionStatsChart />);
    expect(screen.getByTestId('k8s-icon')).toBeInTheDocument();
    expect(screen.getByText('KUBERNETES CLUSTER STATUS')).toBeInTheDocument();
  });

  it('shows a loading indicator on the initial load', () => {
    connectionsQueryReturn = { isFetching: true, isLoading: true };
    render(<KubernetesConnectionStatsChart />);
    expect(screen.getByTestId('loading-container')).toBeInTheDocument();
    expect(screen.getByTestId('circular-progress')).toBeInTheDocument();
    expect(screen.queryByTestId('connect-cluster')).not.toBeInTheDocument();
  });

  it('keeps showing the chart during a background refetch instead of the loading indicator', () => {
    connectionsQueryReturn = {
      data: { connections: [{ status: 'connected' }] },
      isFetching: true,
      isLoading: false,
    };
    render(<KubernetesConnectionStatsChart />);
    expect(screen.queryByTestId('loading-container')).not.toBeInTheDocument();
    expect(screen.getByTestId('bb-chart')).toBeInTheDocument();
  });

  it('shows the error fallback when the connections query fails', () => {
    connectionsQueryReturn = { isError: true };
    render(<KubernetesConnectionStatsChart />);
    expect(screen.getByTestId('widget-error-fallback')).toHaveAttribute(
      'data-title',
      'Kubernetes Cluster Status',
    );
    expect(screen.queryByTestId('connect-cluster')).not.toBeInTheDocument();
  });
});
