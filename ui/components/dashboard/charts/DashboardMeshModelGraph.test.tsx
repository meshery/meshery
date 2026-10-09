import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';

let modelsCount = 1;
let componentsCount = 2;
let relationshipsCount = 3;
let registrantsCount = 4;
let categorySummary: Record<string, number> = {};

const bbChartSpy = vi.fn();
const canSpy = vi.fn((_key?: unknown) => true);

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

vi.mock('@/rtk-query/meshModel', () => ({
  useGetMeshModelsQuery: () => ({ data: { totalCount: modelsCount } }),
  useGetComponentsQuery: () => ({ data: { totalCount: componentsCount } }),
  useGetRelationshipsQuery: () => ({ data: { totalCount: relationshipsCount } }),
  useGetRegistrantsQuery: () => ({ data: { totalCount: registrantsCount } }),
  useGetCategoriesSummary: () => categorySummary,
}));

vi.mock('next/router', () => ({ useRouter: () => ({ push: vi.fn() }) }));

vi.mock('../style', () => ({
  DashboardSection: ({ children }: { children?: React.ReactNode }) => <section>{children}</section>,
}));

vi.mock('@sistent/sistent', () => ({
  Grid2: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
  InfoOutlinedIcon: () => <svg data-testid="info-icon" />,
  Typography: ({ children }: { children?: React.ReactNode }) => <p>{children}</p>,
  useTheme: () => ({ palette: { icon: { default: '#000' } } }),
  // Forward the requested key so tests can pin *which* permission the Registry
  // link authorizes against, not just that it asked for one.
  useHasPermission: (key: unknown) => canSpy(key),
}));

import { Keys } from '@meshery/schemas/permissions';
import MeshModelGraph from './DashboardMeshModelGraph';

describe('DashboardMeshModelGraph', () => {
  beforeEach(() => {
    bbChartSpy.mockReset();
    canSpy.mockClear();
    canSpy.mockReturnValue(true);
    modelsCount = 1;
    componentsCount = 2;
    relationshipsCount = 3;
    registrantsCount = 4;
    categorySummary = { Networking: 5, Security: 3 };
  });

  it('authorizes the Registry link against MesherySystemViewRegistry', () => {
    render(<MeshModelGraph />);
    expect(canSpy).toHaveBeenCalledWith(Keys.MesherySystemViewRegistry);
  });

  it('disables the Registry link when the registry permission is denied', () => {
    canSpy.mockReturnValue(false);
    render(<MeshModelGraph />);
    const registryLink = screen
      .getAllByRole('link')
      .find((link) => link.getAttribute('href') === '/settings?settingsCategory=Registry');
    expect(registryLink).toHaveStyle({ pointerEvents: 'none' });
  });

  it('renders both subcharts (Models by Category and Registry)', () => {
    render(<MeshModelGraph />);
    expect(screen.getAllByTestId('bb-chart')).toHaveLength(2);
    expect(screen.getByText('Models by Category')).toBeInTheDocument();
    expect(screen.getByText('Registry')).toBeInTheDocument();
  });

  it('passes the registry counts into the Registry donut chart', () => {
    render(<MeshModelGraph />);
    // The Registry chart is the second BBChart rendered (after Models by Category).
    const callArgs = bbChartSpy.mock.calls.map((c) => c[0]) as Array<{
      data: { columns: Array<[string, number]> };
    }>;
    // One of the calls includes Models/Components/Relationships/Registrants
    const registryCall = callArgs.find((call) => call.data.columns.some(([k]) => k === 'Models'));
    expect(registryCall).toBeDefined();
    expect(registryCall!.data.columns).toEqual([
      ['Models', 1],
      ['Components', 2],
      ['Relationships', 3],
      ['Registrants', 4],
    ]);
  });

  it('passes the category summary into the Models by Category chart', () => {
    render(<MeshModelGraph />);
    const callArgs = bbChartSpy.mock.calls.map((c) => c[0]) as Array<{
      data: { columns: Array<[string, number]> };
    }>;
    const categoryCall = callArgs.find((call) =>
      call.data.columns.some(([k]) => k === 'Networking'),
    );
    expect(categoryCall).toBeDefined();
    const map = Object.fromEntries(categoryCall!.data.columns);
    expect(map).toEqual({ Networking: 5, Security: 3 });
  });

  it('falls back to 0 counts when the queries return undefined data', () => {
    modelsCount = 0;
    componentsCount = 0;
    relationshipsCount = 0;
    registrantsCount = 0;
    render(<MeshModelGraph />);
    const callArgs = bbChartSpy.mock.calls.map((c) => c[0]) as Array<{
      data: { columns: Array<[string, number]> };
    }>;
    const registryCall = callArgs.find((call) => call.data.columns.some(([k]) => k === 'Models'));
    expect(registryCall!.data.columns).toEqual([
      ['Models', 0],
      ['Components', 0],
      ['Relationships', 0],
      ['Registrants', 0],
    ]);
  });
});
