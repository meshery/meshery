import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';

let patternsQueryReturn: {
  data?: { patterns?: unknown[]; totalCount?: number };
  error?: unknown;
} = { data: { patterns: [], totalCount: 0 } };
let filtersQueryReturn: {
  data?: { filters?: unknown[]; totalCount?: number };
  error?: unknown;
} = { data: { filters: [], totalCount: 0 } };

const bbChartSpy = vi.fn();

vi.mock('billboard.js', () => ({ donut: () => 'donut' }));

vi.mock('../../general/BBChart', () => ({
  default: (props: { options: unknown }) => {
    bbChartSpy(props.options);
    return <div data-testid="bb-chart" />;
  },
}));

vi.mock('../../../utils/charts', () => ({ dataToColors: () => ({}) }));

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

vi.mock('../../../css/icons.styles', () => ({ iconSmall: {} }));

vi.mock('@/components/meshery-mesh-interface/PatternService/CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
}));

vi.mock('@/rtk-query/design', () => ({
  useGetPatternsQuery: () => patternsQueryReturn,
}));

vi.mock('@/rtk-query/filter', () => ({
  useGetFiltersQuery: () => filtersQueryReturn,
}));

vi.mock('@/utils/can', () => ({ default: () => true }));

vi.mock('next/router', () => ({ useRouter: () => ({ push: vi.fn() }) }));

vi.mock('../style', () => ({
  DashboardSection: ({ children }: { children?: React.ReactNode }) => <section>{children}</section>,
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
  useTheme: () => ({ palette: { icon: { default: '#000' } } }),
}));

import MesheryConfigurationChart from './MesheryConfigurationCharts';

describe('MesheryConfigurationChart', () => {
  beforeEach(() => {
    bbChartSpy.mockReset();
    patternsQueryReturn = { data: { patterns: [], totalCount: 0 } };
    filtersQueryReturn = { data: { filters: [], totalCount: 0 } };
  });

  it('renders BBChart with Designs and Filters when data is present', () => {
    patternsQueryReturn = { data: { patterns: [{}], totalCount: 3 } };
    filtersQueryReturn = { data: { filters: [{}], totalCount: 5 } };

    render(<MesheryConfigurationChart />);
    expect(screen.getByTestId('bb-chart')).toBeInTheDocument();
    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: { columns: Array<[string, number]> };
    };
    expect(options.data.columns).toEqual([
      ['Designs', 3],
      ['Filters', 5],
    ]);
  });

  it('renders the ConnectCluster fallback when both queries return errored data (no buckets)', () => {
    patternsQueryReturn = { error: { status: 500 } };
    filtersQueryReturn = { error: { status: 500 } };
    render(<MesheryConfigurationChart />);
    expect(screen.getByTestId('connect-cluster')).toBeInTheDocument();
    expect(bbChartSpy).not.toHaveBeenCalled();
  });

  it('skips a query bucket when patterns/filters query errors', () => {
    patternsQueryReturn = { error: { status: 500 } };
    filtersQueryReturn = { data: { filters: [{}], totalCount: 2 } };

    render(<MesheryConfigurationChart />);
    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: { columns: Array<[string, number]> };
    };
    expect(options.data.columns).toEqual([['Filters', 2]]);
  });

  it('falls back to 0 for totalCount when totalCount is missing on either query response', () => {
    patternsQueryReturn = { data: { patterns: [{}] } };
    filtersQueryReturn = { data: { filters: [{}] } };
    render(<MesheryConfigurationChart />);
    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: { columns: Array<[string, number]> };
    };
    expect(options.data.columns).toEqual([
      ['Designs', 0],
      ['Filters', 0],
    ]);
  });
});
