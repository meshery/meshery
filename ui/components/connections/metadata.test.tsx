import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import FormatConnectionMetadata from './metadata';

const pingKubernetes = vi.fn();
const pingMesheryOperator = vi.fn();
const pingMeshSync = vi.fn();
const pingNats = vi.fn();
const getControllerStatesByConnectionID = vi.fn();

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
      palette: { background: { card: 'black' }, text: { tertiary: 'gray' } },
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
  FormatId: ({ id }) => <span>{id}</span>,
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
  Link: ({ href, title }) => <a href={href}>{title}</a>,
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

  it('renders the well-known meshery server fields as labeled key-values', () => {
    render(
      <FormatConnectionMetadata
        connection={{
          kind: 'meshery',
          createdAt: '2026-05-08T10:00:00Z',
          updatedAt: '2026-05-09T10:00:00Z',
          metadata: {
            serverId: 'server-uuid-1',
            serverVersion: 'v0.9.0',
            serverBuildSha: 'abc1234',
            serverLocation: 'https://meshery.local:9081',
          },
        }}
      />,
    );

    expect(screen.getByText('Server ID')).toBeInTheDocument();
    expect(screen.getByText('server-uuid-1')).toBeInTheDocument();
    expect(screen.getByText('Server Version')).toBeInTheDocument();
    expect(screen.getByText('v0.9.0')).toBeInTheDocument();
    expect(screen.getByText('Server Build SHA')).toBeInTheDocument();
    expect(screen.getByText('abc1234')).toBeInTheDocument();
    expect(screen.getByText('Server Location')).toBeInTheDocument();
    // The location renders as a real, navigable link.
    expect(screen.getByRole('link', { name: 'https://meshery.local:9081' })).toHaveAttribute(
      'href',
      'https://meshery.local:9081',
    );
    expect(screen.getByText('Discovered At')).toBeInTheDocument();
    expect(screen.getByText('Updated At')).toBeInTheDocument();
    // Only well-known fields are present, so the generic fallback is omitted.
    expect(screen.queryByTestId('structured-data')).not.toBeInTheDocument();
  });

  it('tolerates snake_case metadata from older meshery connection records', () => {
    render(
      <FormatConnectionMetadata
        connection={{
          kind: 'meshery',
          metadata: {
            server_id: 'server-uuid-2',
            server_version: 'v0.8.0',
          },
        }}
      />,
    );

    expect(screen.getByText('server-uuid-2')).toBeInTheDocument();
    expect(screen.getByText('v0.8.0')).toBeInTheDocument();
  });

  it('still renders unrecognized meshery metadata through the structured formatter', () => {
    render(
      <FormatConnectionMetadata
        connection={{
          kind: 'meshery',
          metadata: { endpoint: 'https://meshery.local' },
        }}
      />,
    );

    expect(screen.getByTestId('structured-data')).toHaveTextContent('meshery.local');
  });

  it('shows a dash when kubernetes metadata.server is missing', () => {
    render(
      <FormatConnectionMetadata
        meshsyncControllerState={{}}
        connection={{
          id: 'connection-2',
          kind: 'kubernetes',
          status: 'connected',
          metadata: {
            name: 'cluster-b',
            // no server field
          },
        }}
      />,
    );

    // The Server row should render a plain dash, not an invalid link.
    expect(screen.queryByRole('link')).toBeNull();
    expect(screen.getByText('-')).toBeInTheDocument();
  });

  it('falls back to the generic structured formatter for other connection kinds', () => {
    render(
      <FormatConnectionMetadata
        connection={{
          kind: 'github',
          metadata: { owner: 'meshery' },
        }}
      />,
    );

    expect(screen.getByTestId('structured-data')).toHaveTextContent('meshery');
  });
});
