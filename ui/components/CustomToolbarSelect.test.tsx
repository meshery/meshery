import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const dispatch = vi.fn();
let mockResultsSelection: Record<string, Record<string, any>> = {};

vi.mock('react-redux', () => ({
  useDispatch: () => dispatch,
  useSelector: (selector: any) =>
    selector({ prefTest: { results_selection: mockResultsSelection } }),
}));

vi.mock('@/store/slices/prefTest', () => ({
  clearResultsSelection: () => ({ type: 'prefTest/clearResultsSelection' }),
}));

vi.mock('@sistent/sistent', () => ({
  IndeterminateCheckBoxIcon: () => <span data-testid="indeterminate-checkbox" />,
  CompareArrowsIcon: () => <span data-testid="compare-arrows" />,
  GetAppIcon: () => <span data-testid="get-app" />,
  IconButton: ({ children, onClick, href, download, 'aria-label': ariaLabel }: any) => (
    <a
      data-testid={ariaLabel ? `icon-button-${ariaLabel}` : 'icon-button'}
      onClick={onClick}
      href={href}
      download={download}
    >
      {children}
    </a>
  ),
  Tooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
  styled: () => () => {
    const Styled = ({ children, onClick, href, download, 'aria-label': ariaLabel }: any) => (
      <a
        data-testid={ariaLabel ? `styled-icon-button-${ariaLabel}` : 'styled-icon-button'}
        onClick={onClick}
        href={href}
        download={download}
      >
        {children}
      </a>
    );
    Styled.displayName = 'StyledMock';
    return Styled;
  },
  NoSsr: ({ children }: any) => <>{children}</>,
}));

vi.mock('./MesheryChartDialog', () => ({
  default: ({ open, content, handleClose }: any) =>
    open ? (
      <div data-testid="chart-dialog">
        <button type="button" onClick={handleClose}>
          Close Dialog
        </button>
        {content}
      </div>
    ) : null,
}));

vi.mock('./MesheryChart', () => ({
  default: ({ data }: any) => <div data-testid="meshery-chart">{JSON.stringify(data)}</div>,
}));

import CustomToolbarSelect from './CustomToolbarSelect';

describe('CustomToolbarSelect', () => {
  beforeEach(() => {
    dispatch.mockClear();
    mockResultsSelection = {};
  });

  it('renders deselect and compare buttons by default', () => {
    render(<CustomToolbarSelect setSelectedRows={vi.fn()} />);

    expect(screen.getByTestId('indeterminate-checkbox')).toBeInTheDocument();
    expect(screen.getByTestId('compare-arrows')).toBeInTheDocument();
    expect(screen.queryByTestId('get-app')).not.toBeInTheDocument();
  });

  it('calls setSelectedRows([]) and dispatches clearResultsSelection on deselect', async () => {
    const user = userEvent.setup();
    const setSelectedRows = vi.fn();
    render(<CustomToolbarSelect setSelectedRows={setSelectedRows} />);

    const tooltips = screen.getAllByTestId('tooltip');
    const deselectTooltip = tooltips.find((t) => t.getAttribute('data-title') === 'Deselect ALL');
    // The styled icon button is inside the tooltip.
    const button = deselectTooltip!.querySelector('a');
    await user.click(button!);

    expect(setSelectedRows).toHaveBeenCalledWith([]);
    expect(dispatch).toHaveBeenCalledWith({ type: 'prefTest/clearResultsSelection' });
  });

  it('shows the download button when exactly one result is selected', () => {
    mockResultsSelection = {
      k1: {
        k2: {
          meshery_id: 'meshery-abc',
          name: 'profile-1',
          runner_results: { StartTime: '2023-01-01', ActualDuration: 1, Labels: 'L' },
          server_metrics: {},
          server_board_config: {},
        },
      },
    };

    render(<CustomToolbarSelect setSelectedRows={vi.fn()} />);

    const downloadAnchor = screen.getByText(
      (_text, el) =>
        el?.tagName === 'A' &&
        (el as HTMLAnchorElement).getAttribute('download') === 'profile-1_test_result.json',
    ) as HTMLAnchorElement;
    expect(downloadAnchor).toBeInTheDocument();
    expect(downloadAnchor.getAttribute('href')).toContain('meshery-abc');
  });

  it('opens the chart dialog when compare button is clicked', async () => {
    mockResultsSelection = {
      k1: {
        k2: {
          meshery_id: 'meshery-abc',
          name: 'profile-1',
          runner_results: {
            StartTime: '2023-01-01T00:00:00Z',
            ActualDuration: 1000000,
            Labels: 'env=prod',
          },
          server_metrics: { cpu: 0 },
          server_board_config: { panels: [] },
        },
      },
    };

    const user = userEvent.setup();
    render(<CustomToolbarSelect setSelectedRows={vi.fn()} />);

    const tooltips = screen.getAllByTestId('tooltip');
    const compareTooltip = tooltips.find(
      (t) => t.getAttribute('data-title') === 'Compare selected',
    );
    const compareButton = compareTooltip!.querySelector('a');
    await user.click(compareButton!);

    expect(screen.getByTestId('chart-dialog')).toBeInTheDocument();
    expect(screen.getByTestId('meshery-chart')).toBeInTheDocument();
  });

  it('skips undefined entries when iterating results_selection', () => {
    mockResultsSelection = {
      k1: { k2: undefined as any },
    };

    expect(() => render(<CustomToolbarSelect setSelectedRows={vi.fn()} />)).not.toThrow();
    // Download tooltip should not render with zero "filled" rows.
    expect(screen.queryByTestId('get-app')).not.toBeInTheDocument();
  });
});
