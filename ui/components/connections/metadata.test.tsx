import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import FormatConnectionMetadata from './metadata';

const pingKubernetes = vi.fn();
const pingMesheryOperator = vi.fn();
const pingMeshSync = vi.fn();
const pingNats = vi.fn();
const getControllerStatesByConnectionID = vi.fn();
const triggerPrometheusPing = vi.fn();
const triggerGrafanaPing = vi.fn();

vi.mock('@sistent/sistent', () => {
  const styled = (Component) => () => {
    const StyledComponent = ({ children, ...props }) => (
      <Component {...props}>{children}</Component>
    );
    StyledComponent.displayName = 'StyledSistentMock';
    return StyledComponent;
  };

  return {
    Grid2: ({ children }) => <div>{children}</div>,
    List: ({ children }) => <div>{children}</div>,
    ListItem: ({ children }) => <div>{children}</div>,
    ListItemText: ({ primary, secondary }) => (
      <div>
        <span>{primary}</span>
        <span>{secondary}</span>
      </div>
    ),
    InfoIcon: () => <svg data-testid="info-icon" />,
    EventBus: class {
      publish() {}
      subscribe() {}
      on() {
        return { subscribe() {} };
      }
    },
    Box: ({ children }) => <div>{children}</div>,
    Typography: ({ children }) => <span>{children}</span>,
    styled,
    createTheme: () => ({
      breakpoints: {
        up: () => 'up',
        down: () => 'down',
      },
    }),
    useTheme: () => ({
      palette: {
        background: { card: 'black', default: '#111' },
        text: { tertiary: 'gray', secondary: 'silver' },
        error: { main: 'red' },
        warning: { main: 'orange' },
        info: { main: 'blue' },
        divider: '#333',
      },
    }),
  };
});

vi.mock('css/icons.styles', () => ({
  iconMedium: {},
  iconSmall: {},
}));

vi.mock('../../utils/utils', () => ({
  formatToTitleCase: (value) => value,
}));

vi.mock('../data-formatter', () => ({
  FormatId: ({ id }) => <span data-testid="format-id">{id}</span>,
  FormatStructuredData: ({ data }) => (
    <div data-testid="structured-data">{JSON.stringify(data || {})}</div>
  ),
  FormattedDate: ({ date }) => <span>{String(date)}</span>,
  KeyValue: ({ Key, Value }) => (
    <div>
      <span>{Key}</span>
      <span>{Value}</span>
    </div>
  ),
  Link: ({ title, href }) => <a href={href}>{title}</a>,
  createColumnUiSchema: ({ metadata }) => ({ fields: Object.keys(metadata || {}) }),
}));

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => pingKubernetes,
  useControllerStatus: () => ({ getControllerStatesByConnectionID }),
  useMesheryOperator: () => ({ ping: pingMesheryOperator }),
  useMeshsSyncController: () => ({ ping: pingMeshSync }),
  useNatsController: () => ({ ping: pingNats }),
}));

// ControllerDiagnosticsSection calls this RTK Query hook, which needs a Redux
// <Provider>. Stub it so the metadata component renders without a store.
vi.mock('@/rtk-query/connection', () => ({
  useGetControllerDiagnosticsQuery: () => ({
    data: { diagnostics: [] },
    isFetching: false,
    refetch: vi.fn(),
  }),
}));

const emptyPingState = {
  data: undefined,
  isError: false,
  isFetching: false,
  isUninitialized: true,
  isSuccess: false,
};

vi.mock('@/rtk-query/telemetryPrometheus', () => ({
  // Lazy: trigger only - cache reads go through useQueryState.
  useLazyPingPrometheusConnectionQuery: () => [triggerPrometheusPing, emptyPingState],
  default: {
    endpoints: {
      pingPrometheusConnection: {
        useQueryState: () => emptyPingState,
      },
    },
  },
}));

vi.mock('@/rtk-query/telemetryGrafana', () => ({
  useLazyPingGrafanaConnectionQuery: () => [triggerGrafanaPing, emptyPingState],
  default: {
    endpoints: {
      pingGrafanaConnection: {
        useQueryState: () => emptyPingState,
      },
    },
  },
}));

vi.mock('./ConnectionChip', () => ({
  TooltipWrappedConnectionChip: ({ title, handlePing, disabled }) => (
    <button disabled={disabled} onClick={handlePing} type="button">
      {String(title)}
    </button>
  ),
}));

vi.mock('./styles', () => ({
  ColumnWrapper: ({ children }) => <div>{children}</div>,
  ContentContainer: ({ children }) => <div>{children}</div>,
  OperationButton: ({ children }) => <div>{children}</div>,
  FormatterWrapper: ({ children }) => <div>{children}</div>,
}));

describe('FormatConnectionMetadata', () => {
  beforeEach(() => {
    pingKubernetes.mockReset();
    pingMesheryOperator.mockReset();
    pingMeshSync.mockReset();
    pingNats.mockReset();
    triggerPrometheusPing.mockReset();
    triggerGrafanaPing.mockReset();
    getControllerStatesByConnectionID.mockReset();
    getControllerStatesByConnectionID.mockReturnValue({
      operatorState: 'DEPLOYED',
      meshSyncState: 'CONNECTED',
      natsState: 'RUNNING',
      operatorVersion: 'v1.0.0',
      meshSyncVersion: 'v2.0.0',
      natsVersion: 'v3.0.0',
    });
  });

  it('wires kubernetes metadata chips to the correct ping handlers', () => {
    render(
      <FormatConnectionMetadata
        meshsyncControllerState={{}}
        connection={{
          id: 'connection-1',
          kind: 'kubernetes',
          status: 'connected',
          metadata: {
            name: 'cluster-a',
            server: 'https://cluster-a.local',
            meshsync_deployment_mode: 'operator',
          },
        }}
      />,
    );

    fireEvent.click(screen.getByRole('button', { name: 'cluster-a' }));
    fireEvent.click(screen.getByRole('button', { name: 'Operator' }));
    fireEvent.click(screen.getByRole('button', { name: 'MeshSync' }));
    fireEvent.click(screen.getByRole('button', { name: 'BROKER' }));

    expect(pingKubernetes).toHaveBeenCalledWith(
      'cluster-a',
      'https://cluster-a.local',
      'connection-1',
    );
    expect(pingMesheryOperator).toHaveBeenCalledWith({ connectionID: 'connection-1' });
    expect(pingMeshSync).toHaveBeenCalledWith({ connectionID: 'connection-1' });
    expect(pingNats).toHaveBeenCalledWith({ connectionID: 'connection-1' });
  });

  it('renders meshery platform metadata from canonical server-* keys', () => {
    render(
      <FormatConnectionMetadata
        connection={{
          kind: 'meshery',
          name: 'My Meshery',
          status: 'connected',
          createdAt: '2024-01-01',
          updatedAt: '2024-01-02',
          metadata: {
            serverId: 'instance-123',
            serverVersion: 'v0.8.0',
            serverBuildSha: 'deadbeefcafebabe',
            serverLocation: 'https://meshery.local',
          },
        }}
      />,
    );

    expect(screen.getByText('Server Name')).toBeInTheDocument();
    expect(screen.getByText('Server Version')).toBeInTheDocument();
    expect(screen.getByText('v0.8.0')).toBeInTheDocument();
    expect(screen.getByText('Server Location')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'https://meshery.local' })).toHaveAttribute(
      'href',
      'https://meshery.local',
    );
    expect(screen.getByText('Server ID')).toBeInTheDocument();
    expect(screen.getByTestId('format-id')).toHaveTextContent('instance-123');
    expect(screen.getByRole('link', { name: 'deadbee' })).toHaveAttribute(
      'href',
      'https://github.com/meshery/meshery/commit/deadbeefcafebabe',
    );
  });

  it('renders prometheus details and only pings when the connection chip is clicked', () => {
    render(
      <FormatConnectionMetadata
        connection={{
          id: 'prom-1',
          kind: 'prometheus',
          name: 'Prom Prod',
          status: 'connected',
          metadata: {
            url: 'https://prometheus.example:9090',
            telemetryPrometheusPanels: [{ id: 'p1' }, { id: 'p2' }],
          },
        }}
      />,
    );

    // Expanding / mounting the detail panel must not ping.
    expect(triggerPrometheusPing).not.toHaveBeenCalled();
    expect(screen.getByText('Server')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'https://prometheus.example:9090' })).toHaveAttribute(
      'href',
      'https://prometheus.example:9090',
    );
    expect(screen.getByText('Version')).toBeInTheDocument();
    // No ping yet - version stays placeholder.
    expect(screen.getByText('-')).toBeInTheDocument();
    expect(screen.getByText('Saved Panels')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Prom Prod' }));
    expect(triggerPrometheusPing).toHaveBeenCalledTimes(1);
    expect(triggerPrometheusPing).toHaveBeenCalledWith({ connectionID: 'prom-1' });
  });

  it('renders grafana details and only pings when the connection chip is clicked', () => {
    render(
      <FormatConnectionMetadata
        connection={{
          id: 'graf-1',
          kind: 'grafana',
          name: 'Grafana Ops',
          status: 'connected',
          metadata: {
            url: 'https://grafana.example',
            telemetryPinnedBoards: [{ uid: 'b1' }],
          },
        }}
      />,
    );

    expect(triggerGrafanaPing).not.toHaveBeenCalled();
    expect(screen.getByRole('link', { name: 'https://grafana.example' })).toHaveAttribute(
      'href',
      'https://grafana.example',
    );
    expect(screen.getByText('Pinned Boards')).toBeInTheDocument();
    expect(screen.getByText('1')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Grafana Ops' }));
    expect(triggerGrafanaPing).toHaveBeenCalledTimes(1);
    expect(triggerGrafanaPing).toHaveBeenCalledWith({ connectionID: 'graf-1' });
  });

  it('renders github app metadata from installationId and snapshotPaths', () => {
    render(
      <FormatConnectionMetadata
        connection={{
          kind: 'github',
          name: 'GitHub App',
          type: 'source',
          subType: 'git',
          status: 'connected',
          metadata: {
            installationId: '12345678',
            snapshotPaths: [{ 'layer5io/meshery': '/.github' }],
          },
        }}
      />,
    );

    expect(screen.getByText('Installation ID')).toBeInTheDocument();
    expect(screen.getByTestId('format-id')).toHaveTextContent('12345678');
    expect(screen.getByText('Type')).toBeInTheDocument();
    expect(screen.getByText('source / git')).toBeInTheDocument();
    expect(screen.getByText('Snapshot Paths')).toBeInTheDocument();
    expect(screen.getByText('1 configured')).toBeInTheDocument();
    // Must not surface design-file location fields as if they were connection metadata.
    expect(screen.queryByText('Account Owner')).not.toBeInTheDocument();
    expect(screen.queryByText('Target Path')).not.toBeInTheDocument();
  });

  it('falls back to the generic structured formatter for unknown connection kinds', () => {
    render(
      <FormatConnectionMetadata
        connection={{
          kind: 'slack',
          metadata: { workspace: 'layer5' },
        }}
      />,
    );

    expect(screen.getByTestId('structured-data')).toHaveTextContent('layer5');
  });
});
