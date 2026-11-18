import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi, beforeEach } from 'vitest';

let patternsQueryReturn: {
  data?: { patterns?: unknown[]; totalCount?: number };
  error?: unknown;
} = { data: { patterns: [], totalCount: 0 } };
let filtersQueryReturn: {
  data?: { filters?: unknown[]; totalCount?: number };
  error?: unknown;
} = { data: { filters: [], totalCount: 0 } };

const permissionOverrides = new Map<string, boolean>();
const canSpy = vi.fn((key?: { id?: string }) => permissionOverrides.get(key?.id ?? '') ?? true);
const bbChartSpy = vi.fn();
const routerPushSpy = vi.fn();
const getPatternsQuerySpy = vi.fn();
const getFiltersQuerySpy = vi.fn();

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
  useGetPatternsQuery: (...args: unknown[]) => {
    getPatternsQuerySpy(...args);
    return patternsQueryReturn;
  },
}));

vi.mock('@/rtk-query/filter', () => ({
  useGetFiltersQuery: (...args: unknown[]) => {
    getFiltersQuerySpy(...args);
    return filtersQueryReturn;
  },
}));

vi.mock('@/utils/can', () => ({ default: () => true }));

vi.mock('next/router', () => ({ useRouter: () => ({ push: routerPushSpy }) }));

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
  useHasPermission: (key?: { id?: string }) => canSpy(key),
}));

import { Keys } from '@meshery/schemas/permissions';
import MesheryConfigurationChart from './MesheryConfigurationCharts';

describe('MesheryConfigurationChart', () => {
  beforeEach(() => {
    bbChartSpy.mockReset();
    permissionOverrides.clear();
    canSpy.mockClear();
    routerPushSpy.mockReset();
    getPatternsQuerySpy.mockClear();
    getFiltersQuerySpy.mockClear();
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

  it('passes { skip: false } to both queries when Designs and Filters permissions are granted', () => {
    render(<MesheryConfigurationChart />);
    expect(canSpy).toHaveBeenCalledWith(Keys.CatalogManagementViewDesigns);
    expect(canSpy).toHaveBeenCalledWith(Keys.CatalogManagementViewFilters);
    expect(getPatternsQuerySpy).toHaveBeenCalledWith({ page: 0, pagesize: 1 }, { skip: false });
    expect(getFiltersQuerySpy).toHaveBeenCalledWith({ page: 0, pagesize: 1 }, { skip: false });
    expect(screen.getByRole('link')).toHaveAttribute('href', '/configuration/designs');
  });

  it('skips useGetPatternsQuery, omits Designs, and links to /configuration/filters when only CatalogManagementViewDesigns is denied', () => {
    permissionOverrides.set(Keys.CatalogManagementViewDesigns.id, false);
    patternsQueryReturn = { data: { patterns: [{}], totalCount: 3 } };
    filtersQueryReturn = { data: { filters: [{}], totalCount: 5 } };

    render(<MesheryConfigurationChart />);
    expect(getPatternsQuerySpy).toHaveBeenCalledWith({ page: 0, pagesize: 1 }, { skip: true });
    expect(getFiltersQuerySpy).toHaveBeenCalledWith({ page: 0, pagesize: 1 }, { skip: false });
    expect(screen.getByRole('link')).toHaveAttribute('href', '/configuration/filters');

    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: {
        columns: Array<[string, number]>;
        onclick: (dataPoint: { name: string }) => void;
      };
    };
    expect(options.data.columns).toEqual([['Filters', 5]]);
    options.data.onclick({ name: 'Designs' });
    expect(routerPushSpy).not.toHaveBeenCalled();
    options.data.onclick({ name: 'Filters' });
    expect(routerPushSpy).toHaveBeenCalledWith('/configuration/filters');
  });

  it('skips useGetFiltersQuery and omits Filters when CatalogManagementViewFilters is denied', () => {
    permissionOverrides.set(Keys.CatalogManagementViewFilters.id, false);
    patternsQueryReturn = { data: { patterns: [{}], totalCount: 3 } };
    filtersQueryReturn = { data: { filters: [{}], totalCount: 5 } };

    render(<MesheryConfigurationChart />);
    expect(getPatternsQuerySpy).toHaveBeenCalledWith({ page: 0, pagesize: 1 }, { skip: false });
    expect(getFiltersQuerySpy).toHaveBeenCalledWith({ page: 0, pagesize: 1 }, { skip: true });
    expect(screen.getByRole('link')).toHaveAttribute('href', '/configuration/designs');

    const options = bbChartSpy.mock.calls.at(-1)?.[0] as {
      data: {
        columns: Array<[string, number]>;
        onclick: (dataPoint: { name: string }) => void;
      };
    };
    expect(options.data.columns).toEqual([['Designs', 3]]);
    options.data.onclick({ name: 'Filters' });
    expect(routerPushSpy).not.toHaveBeenCalled();
    options.data.onclick({ name: 'Designs' });
    expect(routerPushSpy).toHaveBeenCalledWith('/configuration/designs');
  });

  it('skips both queries and renders permission message without ConnectCluster, Link, or keyboard focus when both view permissions are denied', async () => {
    const user = userEvent.setup();
    permissionOverrides.set(Keys.CatalogManagementViewDesigns.id, false);
    permissionOverrides.set(Keys.CatalogManagementViewFilters.id, false);
    patternsQueryReturn = { data: { patterns: [{}], totalCount: 3 } };
    filtersQueryReturn = { data: { filters: [{}], totalCount: 5 } };

    render(<MesheryConfigurationChart />);
    expect(getPatternsQuerySpy).toHaveBeenCalledWith({ page: 0, pagesize: 1 }, { skip: true });
    expect(getFiltersQuerySpy).toHaveBeenCalledWith({ page: 0, pagesize: 1 }, { skip: true });
    expect(
      screen.getByText('You do not have permission to view Designs or Filters.'),
    ).toBeInTheDocument();
    expect(screen.queryByTestId('connect-cluster')).not.toBeInTheDocument();
    expect(screen.queryByRole('link')).not.toBeInTheDocument();
    await user.tab();
    expect(document.activeElement).toBe(document.body);
    expect(bbChartSpy).not.toHaveBeenCalled();
  });
});
