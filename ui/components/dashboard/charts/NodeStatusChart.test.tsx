import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const bbChartSpy = vi.fn();

vi.mock('billboard.js', () => ({
  donut: () => 'donut',
}));

vi.mock('../../general/BBChart', () => ({
  default: (props: { options: unknown }) => {
    bbChartSpy(props.options);
    return <div data-testid="bb-chart" />;
  },
}));

vi.mock('@sistent/sistent', () => ({
  CircularProgress: () => <div data-testid="spinner" />,
  KEPPEL: '#1ECEC0',
  Typography: ({ children, variant }: { children?: React.ReactNode; variant?: string }) => (
    <p data-testid="typography" data-variant={variant}>
      {children}
    </p>
  ),
  Stack: ({ children }: { children?: React.ReactNode }) => (
    <div data-testid="stack">{children}</div>
  ),
}));

vi.mock('@/theme', () => ({
  useTheme: () => ({
    palette: { error: { main: '#d32f2f' } },
  }),
}));

vi.mock('./utils', () => ({
  getLegendTemplate: () => '<span></span>',
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
  LegendSection: ({ id }: { id?: string }) => <div data-testid="legend" id={id} />,
}));

import { NodeStatusChart } from './NodeStatusChart';

describe('NodeStatusChart', () => {
  it('renders the loading spinner when isClusterLoading is true', () => {
    render(<NodeStatusChart isClusterLoading={true} />);
    expect(screen.getByTestId('loading-container')).toBeInTheDocument();
    expect(screen.getByTestId('spinner')).toBeInTheDocument();
  });

  it('renders the ConnectCluster fallback when nodeData is undefined', () => {
    render(<NodeStatusChart />);
    expect(screen.getByTestId('connect-cluster')).toHaveTextContent(
      'No workloads found in your cluster(s).',
    );
  });

  it('renders the empty-cluster message when total nodes are 0', () => {
    render(<NodeStatusChart nodeData={[]} />);
    expect(screen.getByText('No nodes are currently in the cluster')).toBeInTheDocument();
    expect(screen.queryByTestId('bb-chart')).not.toBeInTheDocument();
  });

  it('renders the BBChart and legend when nodes are present', () => {
    render(
      <NodeStatusChart
        nodeData={[
          { status: 'Ready', count: 3 },
          { status: 'Not Ready', count: 1 },
        ]}
      />,
    );
    expect(screen.getByTestId('bb-chart')).toBeInTheDocument();
    expect(screen.getByTestId('legend')).toHaveAttribute('id', 'nodeLegend');
  });

  it('passes columns derived from nodeData into the BBChart options', () => {
    bbChartSpy.mockReset();
    render(
      <NodeStatusChart
        nodeData={[
          { status: 'Ready', count: 5 },
          { status: 'Not Ready', count: 2 },
        ]}
      />,
    );
    expect(bbChartSpy).toHaveBeenCalled();
    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: { columns: Array<[string, number]> };
      donut: { title: string };
    };
    expect(options.data.columns).toEqual([
      ['Ready', 5],
      ['Not Ready', 2],
    ]);
    expect(options.donut.title).toBe('7\nNodes');
  });
});
