import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const bbChartSpy = vi.fn();
const canSpy = vi.fn(() => true);

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
  isValidColumnName: (name: unknown) => Boolean(name) && name !== 'invalid',
}));

vi.mock('../../general/ConnectClustersBtn', () => ({
  default: () => <button type="button">Connect Cluster</button>,
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

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => canSpy(...args),
}));

vi.mock('@sistent/sistent', () => ({
  Box: ({ children }: { children?: React.ReactNode }) => <div data-testid="box">{children}</div>,
  MenuItem: ({ children, value }: { children?: React.ReactNode; value?: string }) => (
    <option value={value}>{children}</option>
  ),
  Select: ({
    children,
    value,
    onChange,
  }: {
    children?: React.ReactNode;
    value?: string;
    onChange?: (event: { target: { value: string } }) => void;
  }) => (
    <select
      data-testid="namespace-select"
      value={value}
      onChange={(e) => onChange?.({ target: { value: e.target.value } } as never)}
    >
      {children}
    </select>
  ),
  Typography: ({ children, className }: { children?: React.ReactNode; className?: string }) => (
    <p className={className}>{children}</p>
  ),
}));

vi.mock('@/theme', () => ({
  useTheme: () => ({ palette: {} }),
}));

import WorkloadChart from './WorkloadChart';

const classes = {
  dashboardSection: 'dashboard-section',
  link: 'link',
};

describe('WorkloadChart', () => {
  it('renders the empty state when there are no resources', () => {
    bbChartSpy.mockReset();
    render(
      <WorkloadChart
        classes={classes}
        handleSetNamespace={() => {}}
        resourses={[]}
        namespaces={[]}
      />,
    );
    expect(screen.getByText('No workloads found in your cluster(s).')).toBeInTheDocument();
    expect(bbChartSpy).not.toHaveBeenCalled();
  });

  it('renders the BBChart when there is workload data and filters out invalid kinds', () => {
    bbChartSpy.mockReset();
    render(
      <WorkloadChart
        classes={classes}
        handleSetNamespace={() => {}}
        resourses={[
          { kind: 'Pod', count: 5 },
          { kind: 'invalid', count: 1 },
          { kind: 'Deployment', count: 2 },
        ]}
        namespaces={[]}
      />,
    );
    expect(screen.getByTestId('bb-chart')).toBeInTheDocument();
    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: { columns: Array<[string, number]> };
    };
    expect(options.data.columns).toEqual([
      ['Pod', 5],
      ['Deployment', 2],
    ]);
  });

  it('renders namespace selector when namespaces are provided', () => {
    const handleSetNamespace = vi.fn();
    render(
      <WorkloadChart
        classes={classes}
        handleSetNamespace={handleSetNamespace}
        resourses={[{ kind: 'Pod', count: 1 }]}
        namespaces={['default', 'kube-system']}
        selectedNamespace="default"
      />,
    );
    const select = screen.getByTestId('namespace-select');
    expect(select).toBeInTheDocument();
    fireEvent.change(select, { target: { value: 'kube-system' } });
    expect(handleSetNamespace).toHaveBeenCalledWith('kube-system');
  });

  it('blocks pointer events on the Workloads link when the user lacks permission', () => {
    canSpy.mockReturnValueOnce(false);
    render(
      <WorkloadChart
        classes={classes}
        handleSetNamespace={() => {}}
        resourses={[{ kind: 'Pod', count: 1 }]}
      />,
    );
    const link = screen.getByRole('link');
    expect(link).toHaveStyle({ pointerEvents: 'none' });
  });
});
