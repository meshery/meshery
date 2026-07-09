import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';

let connectionsQueryReturn: { data?: { connections?: Array<{ status?: string }> } } = {
  data: { connections: [] },
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
  useGetConnectionsQuery: () => connectionsQueryReturn,
}));

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => canSpy(...args),
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ push: vi.fn() }),
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
}));

import ConnectionStatsChart from './ConnectionCharts';

describe('ConnectionStatsChart (Connections donut chart)', () => {
  beforeEach(() => {
    bbChartSpy.mockReset();
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

  it('applies pointer-events:none to the link when the user lacks permission', () => {
    canSpy.mockReturnValue(false);
    connectionsQueryReturn = { data: { connections: [{ status: 'connected' }] } };
    render(<ConnectionStatsChart />);
    expect(screen.getByRole('link')).toHaveStyle({ pointerEvents: 'none' });
  });

  it('renders the Connections heading regardless of state', () => {
    render(<ConnectionStatsChart />);
    expect(screen.getByText('Connections')).toBeInTheDocument();
  });
});
