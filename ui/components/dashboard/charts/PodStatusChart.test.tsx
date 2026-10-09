import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const bbChartSpy = vi.fn();

vi.mock('billboard.js', () => ({
  donut: () => 'donut',
}));

vi.mock('@/components/general/BBChart', () => ({
  default: (props: { options: unknown }) => {
    bbChartSpy(props.options);
    return <div data-testid="bb-chart" />;
  },
}));

vi.mock('@sistent/sistent', () => ({
  Typography: ({ children, variant }: { children?: React.ReactNode; variant?: string }) => (
    <p data-testid="typography" data-variant={variant}>
      {children}
    </p>
  ),
  SAFFRON: '#F0BC00',
  CircularProgress: () => <div data-testid="spinner" />,
  KEPPEL: '#1ECEC0',
  DARK_SLATE_GRAY: '#222',
  TEAL_BLUE: '#477E96',
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

import { PodStatusChart } from './PodStatusChart';

describe('PodStatusChart', () => {
  it('shows spinner when loading', () => {
    render(<PodStatusChart isClusterLoading={true} />);
    expect(screen.getByTestId('spinner')).toBeInTheDocument();
  });

  it('renders ConnectCluster fallback when podData is undefined', () => {
    render(<PodStatusChart />);
    expect(screen.getByTestId('connect-cluster')).toHaveTextContent(
      'No workloads found in your cluster(s).',
    );
  });

  it('renders the empty-cluster message when total pods are 0', () => {
    render(<PodStatusChart podData={[]} />);
    expect(screen.getByText('No pods are currently in the cluster')).toBeInTheDocument();
  });

  it('renders BBChart and legend with pod data', () => {
    bbChartSpy.mockReset();
    render(
      <PodStatusChart
        podData={[
          { status: 'Running', count: 5 },
          { status: 'Pending', count: 1 },
        ]}
      />,
    );
    expect(screen.getByTestId('bb-chart')).toBeInTheDocument();
    expect(screen.getByTestId('legend')).toHaveAttribute('id', 'podLegend');

    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: { columns: Array<[string, number]>; colors: Record<string, string> };
      donut: { title: string };
    };
    expect(options.data.columns).toEqual([
      ['Running', 5],
      ['Pending', 1],
    ]);
    expect(options.donut.title).toBe('6\nPods');
    expect(options.data.colors.Running).toBe('#1ECEC0');
    expect(options.data.colors.Pending).toBe('#F0BC00');
  });
});
