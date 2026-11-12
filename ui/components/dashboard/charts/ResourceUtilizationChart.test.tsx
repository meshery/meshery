import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const bbChartSpy = vi.fn();

vi.mock('billboard.js', () => ({
  gauge: () => 'gauge',
}));

vi.mock('@/components/general/BBChart', () => ({
  default: (props: { options: unknown }) => {
    bbChartSpy(props.options);
    return <div data-testid="bb-chart" />;
  },
}));

vi.mock('@sistent/sistent', () => ({
  Box: ({ children }: { children?: React.ReactNode }) => <div data-testid="box">{children}</div>,
  Typography: ({ children, variant }: { children?: React.ReactNode; variant?: string }) => (
    <p data-variant={variant}>{children}</p>
  ),
  Stack: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
  KEPPEL: '#1ECEC0',
  SAFFRON: '#F0BC00',
  CircularProgress: () => <div data-testid="spinner" />,
}));

vi.mock('./ConnectCluster', () => ({
  default: ({ message }: { message: React.ReactNode }) => (
    <div data-testid="connect-cluster">{message}</div>
  ),
}));

vi.mock('../style', () => ({
  LoadingContainer: ({ children }: { children?: React.ReactNode }) => (
    <div data-testid="loading-container">{children}</div>
  ),
  ChartSectionWithColumn: ({ children }: { children?: React.ReactNode }) => (
    <div data-testid="chart-section">{children}</div>
  ),
}));

import { ResourceUtilizationChart } from './ResourceUtilizationChart';

describe('ResourceUtilizationChart', () => {
  it('shows the spinner while loading', () => {
    render(<ResourceUtilizationChart isClusterLoading={true} />);
    expect(screen.getByTestId('spinner')).toBeInTheDocument();
  });

  it('renders the ConnectCluster fallback when usageData is missing', () => {
    render(<ResourceUtilizationChart />);
    expect(screen.getByTestId('connect-cluster')).toHaveTextContent(
      'No workloads found in your cluster(s).',
    );
  });

  it('renders three gauge charts for CPU, Memory, and Disk utilization', () => {
    bbChartSpy.mockReset();
    render(
      <ResourceUtilizationChart
        usageData={[
          { resource: 'CPU', percentage: 12 },
          { resource: 'Memory', percentage: 45 },
          { resource: 'Disk', percentage: 90 },
        ]}
      />,
    );

    expect(screen.getByText('CPU Utilization')).toBeInTheDocument();
    expect(screen.getByText('Memory Utilization')).toBeInTheDocument();
    expect(screen.getByText('Disk Utilization')).toBeInTheDocument();
    expect(screen.getAllByTestId('bb-chart')).toHaveLength(3);
    // One call per chart
    expect(bbChartSpy).toHaveBeenCalledTimes(3);
    // Each call should specify a different gauge column.
    const columns = bbChartSpy.mock.calls.map(
      (args) => (args[0] as { data: { columns: unknown } }).data.columns,
    );
    expect(columns).toEqual([[['CPU', 12]], [['Memory', 45]], [['Disk', 90]]]);
  });
});
