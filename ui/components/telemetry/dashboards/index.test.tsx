import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const useGetConnectionsQuery = vi.fn();

vi.mock('@/rtk-query/connection', () => ({
  useGetConnectionsQuery: (...args: unknown[]) => useGetConnectionsQuery(...args),
}));

vi.mock('@/rtk-query/telemetryGrafana', () => ({
  useGetPinnedBoardsQuery: () => ({ data: [] }),
  usePingGrafanaConnectionQuery: () => ({ data: { reachable: true } }),
  useUpdatePinnedBoardsMutation: () => [vi.fn()],
}));

vi.mock('@/utils/context/ConnectionWizardContextProvider', () => ({
  useConnectionWizardModal: () => ({ openCreateConnection: vi.fn() }),
}));

vi.mock('../../general/error-404/index', () => ({
  default: () => <div data-testid="default-error" />,
}));

vi.mock('../common/ConnectionPicker', () => ({ default: () => <div /> }));
vi.mock('../common/PingStatus', () => ({ default: () => <div /> }));
vi.mock('../common/TimeRangePicker', () => ({ default: () => <div /> }));
vi.mock('../common/RefreshControl', () => ({ default: () => <div /> }));
vi.mock('./BoardLibrary', () => ({ default: () => <div /> }));
vi.mock('./BoardView', () => ({ default: () => <div /> }));

vi.mock('@sistent/sistent', () => ({
  AddIcon: () => <span />,
  Box: ({ children, ...props }: { children?: React.ReactNode }) => <div {...props}>{children}</div>,
  Button: ({ children }: { children?: React.ReactNode }) => <button>{children}</button>,
  CircularProgress: () => <div />,
  Drawer: ({ children }: { children?: React.ReactNode }) => <div>{children}</div>,
  InsertChartIcon: () => <span />,
  Typography: ({ children }: { children?: React.ReactNode }) => <p>{children}</p>,
  styled: (Component: unknown) => () => Component,
  useHasPermission: () => true,
  useTheme: () => ({
    palette: {
      background: { default: '#000' },
      border: { default: '#333' },
      icon: { secondary: '#888' },
    },
  }),
}));

import TelemetryDashboards from './index';

describe('TelemetryDashboards', () => {
  beforeEach(() => {
    useGetConnectionsQuery.mockReset();
    useGetConnectionsQuery.mockReturnValue({ data: { connections: [] }, isLoading: false });
  });

  it('shows the empty state and lists grafana connections with a plain kind filter', () => {
    render(<TelemetryDashboards />);

    expect(screen.getByTestId('telemetry-grafana-empty')).toBeInTheDocument();
    expect(screen.getByText('No Grafana connections yet')).toBeInTheDocument();
    // #20617: plain ?kind=grafana, not JSON.stringify(['grafana'])
    expect(useGetConnectionsQuery).toHaveBeenCalledWith(
      { kind: 'grafana', pageSize: 100 },
      { skip: false },
    );
  });
});
