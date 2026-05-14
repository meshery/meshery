import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi, beforeEach } from 'vitest';

const pushSpy = vi.fn();

vi.mock('next/router', () => ({
  useRouter: () => ({ push: pushSpy }),
}));

vi.mock('./ResponsiveHoneycomb', () => ({
  default: ({
    items,
    renderItem,
  }: {
    items: unknown[];
    renderItem: (item: unknown, index: number) => React.ReactNode;
  }) => (
    <div data-testid="responsive-honeycomb">
      {items.map((it, i) => (
        <div key={i} data-testid="rendered-item">
          {renderItem(it, i)}
        </div>
      ))}
    </div>
  ),
}));

vi.mock('./Hexagon', () => ({
  default: ({
    children,
    onClick,
  }: {
    children?: React.ReactNode;
    onClick?: React.MouseEventHandler;
  }) => (
    <button data-testid="hex" onClick={onClick} type="button">
      {children}
    </button>
  ),
}));

vi.mock('../../charts/ConnectCluster', () => ({
  default: ({ message }: { message: React.ReactNode }) => (
    <div data-testid="connect-cluster">{message}</div>
  ),
}));

vi.mock('../../resources/config', () => ({
  generateDynamicURL: (kind: string) => `?kind=${kind}`,
}));

vi.mock('./useResourceOptions', () => ({
  DEFAULT_GROUP_BY: 'all',
  SORT_DIRECTIONS: { ASC: 'asc', DESC: 'desc' },
  useResourceOptions: () => [
    { value: 'all', label: 'All Resources' },
    { value: 'workload', label: 'Workload' },
  ],
  useResourceFiltering: (kinds: Array<{ Kind: string; Count?: number }> | undefined) => kinds ?? [],
}));

vi.mock('../../utils', () => ({
  default: ({ kind }: { kind: string }) => <svg data-testid={`k8s-icon-${kind}`} />,
}));

vi.mock('../../style', () => ({
  HoneycombRoot: ({ children }: { children?: React.ReactNode }) => (
    <div data-testid="honeycomb-root">{children}</div>
  ),
  IconWrapper: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
  ResourceCount: ({ children }: { children?: React.ReactNode }) => (
    <span data-testid="count">{children}</span>
  ),
  SelectedHexagon: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
  SkeletonHexagon: ({ children }: { children?: React.ReactNode }) => (
    <div data-testid="skeleton-hex">{children}</div>
  ),
  HeaderContainer: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
  ControlsContainer: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
  NoResourcesText: ({ children }: { children?: React.ReactNode }) => (
    <div data-testid="no-resources">{children}</div>
  ),
}));

vi.mock('@sistent/sistent', () => ({
  ArrowDownwardIcon: () => <svg data-testid="arrow-down" />,
  ArrowUpwardIcon: () => <svg data-testid="arrow-up" />,
  CustomTooltip: ({ children }: { children?: React.ReactNode }) => <>{children}</>,
  ErrorBoundary: ({ children }: { children?: React.ReactNode }) => <>{children}</>,
  IconButton: ({
    children,
    onClick,
  }: {
    children?: React.ReactNode;
    onClick?: React.MouseEventHandler;
  }) => (
    <button data-testid="icon-button" type="button" onClick={onClick}>
      {children}
    </button>
  ),
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
      data-testid="group-select"
      value={value}
      onChange={(e) => onChange?.({ target: { value: e.target.value } } as never)}
    >
      {children}
    </select>
  ),
  Skeleton: () => <div data-testid="skeleton" />,
  Typography: ({ children }: { children?: React.ReactNode }) => <p>{children}</p>,
}));

import HoneycombComponent from './HoneyCombComponent';

describe('HoneycombComponent', () => {
  beforeEach(() => {
    pushSpy.mockReset();
  });

  it('renders skeleton hexagons while loading', () => {
    render(<HoneycombComponent isClusterLoading={true} />);
    expect(screen.getAllByTestId('skeleton-hex').length).toBeGreaterThan(0);
  });

  it('renders the ConnectCluster fallback when kinds is undefined', () => {
    render(<HoneycombComponent />);
    expect(screen.getByTestId('connect-cluster')).toHaveTextContent(
      'No workloads found in your cluster(s).',
    );
  });

  it('shows "No resources" message when filtered list is empty (but kinds is defined)', () => {
    render(<HoneycombComponent kinds={[]} />);
    expect(screen.getByTestId('no-resources')).toHaveTextContent(
      'No resources found for the selected group',
    );
  });

  it('renders one hex per filtered kind and routes when a hex is clicked', async () => {
    render(
      <HoneycombComponent
        kinds={[
          { Kind: 'Pod', Count: 5 },
          { Kind: 'Deployment', Count: 2 },
        ]}
      />,
    );
    const hexes = screen.getAllByTestId('hex');
    expect(hexes).toHaveLength(2);

    await userEvent.click(hexes[0]);
    expect(pushSpy).toHaveBeenCalledWith('?kind=Pod');
  });

  it('renders the arrow-down icon by default and toggles when the sort button is pressed', async () => {
    render(<HoneycombComponent kinds={[{ Kind: 'Pod', Count: 5 }]} />);
    expect(screen.getByTestId('arrow-down')).toBeInTheDocument();

    await userEvent.click(screen.getByTestId('icon-button'));
    expect(screen.getByTestId('arrow-up')).toBeInTheDocument();
  });

  it('renders the group-by select with the documented options', () => {
    render(<HoneycombComponent kinds={[{ Kind: 'Pod', Count: 1 }]} />);
    expect(screen.getByTestId('group-select')).toBeInTheDocument();
    expect(screen.getByRole('option', { name: 'All Resources' })).toBeInTheDocument();
    expect(screen.getByRole('option', { name: 'Workload' })).toBeInTheDocument();
  });
});
