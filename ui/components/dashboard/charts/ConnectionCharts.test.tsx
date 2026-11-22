import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi, beforeEach } from 'vitest';

let connectionsQueryReturn: { data?: { connections?: Array<{ status?: string }> } } = {
  data: { connections: [] },
};
const canSpy = vi.fn((_key?: unknown) => true);
const bbChartSpy = vi.fn();
const routerPushSpy = vi.fn();
const getConnectionsQuerySpy = vi.fn();

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
  isValidColumnName: (s: unknown) => Boolean(s) && s !== 'invalid' && s !== 'unknown',
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
  CustomTextTooltip: ({ children }: { children?: React.ReactNode }) => (
    <div data-testid="tooltip">{children}</div>
  ),
}));

vi.mock('@/rtk-query/connection', () => ({
  useGetConnectionsQuery: (...args: unknown[]) => {
    getConnectionsQuerySpy(...args);
    return connectionsQueryReturn;
  },
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ push: routerPushSpy }),
}));

vi.mock('../style', () => ({
  DashboardSection: ({ children }: { children?: React.ReactNode }) => (
    <section data-testid="dashboard-section">{children}</section>
  ),
}));

vi.mock('./ConnectCluster', () => ({
  default: ({ message }: { message: React.ReactNode }) => (
    <div data-testid="connect-cluster">{message}</div>
  ),
}));

vi.mock('@sistent/sistent', () => ({
  Box: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
  InfoOutlinedIcon: () => <svg data-testid="info-icon" />,
  Typography: ({ children }: { children?: React.ReactNode }) => <p>{children}</p>,
  useTheme: () => ({ palette: { mode: 'light', icon: { default: '#000' } } }),
  // Forward the requested key so tests can pin *which* permission the chart
  // authorizes against, not just that it asked for one.
  useHasPermission: (key: unknown) => canSpy(key),
}));

import { Keys } from '@meshery/schemas/permissions';
import ConnectionStatsChart from './ConnectionCharts';

describe('ConnectionStatsChart (Connections donut chart)', () => {
  beforeEach(() => {
    bbChartSpy.mockReset();
    routerPushSpy.mockReset();
    getConnectionsQuerySpy.mockClear();
    canSpy.mockClear();
    canSpy.mockReturnValue(true);
    connectionsQueryReturn = { data: { connections: [] } };
  });

  it('renders ConnectCluster fallback when there are no connections', () => {
    render(<ConnectionStatsChart />);
    expect(screen.getByTestId('connect-cluster')).toHaveTextContent(
      'No connections found in your clusters',
    );
  });

  it('omits the "unknown" bucket when isValidColumnName rejects undefined statuses', () => {
    connectionsQueryReturn = {
      data: {
        connections: [
          { status: 'connected' },
          { status: 'connected' },
          { status: 'discovered' },
          {}, // becomes 'unknown'
          { status: 'invalid' },
        ],
      },
    };
    render(<ConnectionStatsChart />);
    expect(screen.getByTestId('bb-chart')).toBeInTheDocument();
    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: { columns: Array<[string, number]> };
    };
    const map = Object.fromEntries(options.data.columns);
    expect(map.connected).toBe(2);
    expect(map.discovered).toBe(1);
    expect(map.unknown).toBeUndefined();
    expect(map.invalid).toBeUndefined();
  });

  it('omits the Link and blocks keyboard focus when the user lacks permission', async () => {
    const user = userEvent.setup();
    canSpy.mockReturnValue(false);
    connectionsQueryReturn = { data: { connections: [{ status: 'connected' }] } };
    render(<ConnectionStatsChart />);
    expect(screen.queryByRole('link')).not.toBeInTheDocument();
    await user.tab();
    expect(document.activeElement).toBe(document.body);
  });

  it('authorizes the connections link against WorkspaceManagementViewConnections', () => {
    render(<ConnectionStatsChart />);
    expect(canSpy).toHaveBeenCalledWith(Keys.WorkspaceManagementViewConnections);
    expect(getConnectionsQuerySpy).toHaveBeenCalledWith(
      { page: 0, pagesize: 'all' },
      { skip: false },
    );
    expect(screen.getByRole('link')).toHaveAttribute('href', '/management/connections');
  });

  it('skips useGetConnectionsQuery and hides chart data when WorkspaceManagementViewConnections is denied', () => {
    canSpy.mockReturnValue(false);
    connectionsQueryReturn = { data: { connections: [{ status: 'connected' }] } };
    render(<ConnectionStatsChart />);
    expect(getConnectionsQuerySpy).toHaveBeenCalledWith(
      { page: 0, pagesize: 'all' },
      { skip: true },
    );
    expect(screen.queryByRole('link')).not.toBeInTheDocument();
    expect(screen.getByTestId('connect-cluster')).toBeInTheDocument();
    expect(screen.queryByTestId('bb-chart')).not.toBeInTheDocument();
  });

  it('navigates on chart onclick only when permitted', () => {
    connectionsQueryReturn = { data: { connections: [{ status: 'connected' }] } };
    render(<ConnectionStatsChart />);
    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: { onclick: () => void };
    };
    options.data.onclick();
    expect(routerPushSpy).toHaveBeenCalledWith('/management/connections');
  });

  it('renders the Connections heading regardless of state', () => {
    render(<ConnectionStatsChart />);
    expect(screen.getByText('Connections')).toBeInTheDocument();
  });
});
