import React from 'react';
import { render, screen } from '@testing-library/react';
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
const canSpy = vi.fn(() => true);
const bbChartSpy = vi.fn();

vi.mock('billboard.js', () => ({
  donut: () => 'donut',
}));

vi.mock('../../BBChart', () => ({
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

vi.mock('../../../css/icons.styles', () => ({
  iconSmall: {},
}));

vi.mock('@/components/meshery-mesh-interface/PatternService/CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
}));

vi.mock('@/rtk-query/connection', () => ({
  useGetConnectionsQuery: () => connectionsQueryReturn,
}));

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => canSpy(...args),
}));

vi.mock('@/utils/permission_constants', () => ({
  keys: { VIEW_CONNECTIONS: { action: 'view', subject: 'connections' } },
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ push: vi.fn() }),
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
  InfoOutlinedIcon: () => <svg data-testid="info-icon" />,
  KubernetesIcon: () => <svg data-testid="k8s-icon" />,
  Typography: ({ children }: { children?: React.ReactNode }) => <p>{children}</p>,
  useTheme: () => ({
    palette: { mode: 'light', icon: { default: '#000', disabled: '#777' } },
  }),
}));

import KubernetesConnectionStatsChart from './KubernetesConnectionChart';

describe('KubernetesConnectionStatsChart', () => {
  beforeEach(() => {
    bbChartSpy.mockReset();
    canSpy.mockClear();
    canSpy.mockReturnValue(true);
    connectionsQueryReturn = {
      data: { connections: [] },
      isFetching: false,
      isLoading: false,
      isError: false,
    };
  });

  it('shows ConnectCluster fallback when there are no kubernetes connections', () => {
    render(<KubernetesConnectionStatsChart />);
    expect(screen.getByTestId('connect-cluster')).toHaveTextContent(
      'No connections found in your clusters',
    );
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

  it('applies pointer-events:none to the Link when the user lacks permission', () => {
    connectionsQueryReturn = {
      data: {
        connections: [{ status: 'connected' }],
      },
    };
    canSpy.mockReturnValue(false);
    render(<KubernetesConnectionStatsChart />);
    expect(screen.getByRole('link')).toHaveStyle({ pointerEvents: 'none' });
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
