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
  Link: ({ title }) => <span>{title}</span>,
  createColumnUiSchema: ({ metadata }) => ({ fields: Object.keys(metadata || {}) }),
}));

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => pingKubernetes,
  useControllerStatus: () => ({ getControllerStatesByConnectionID }),
  useMesheryOperator: () => ({ ping: pingMesheryOperator }),
  useMeshsSyncController: () => ({ ping: pingMeshSync }),
  useNatsController: () => ({ ping: pingNats }),
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

  it('renders structured metadata for meshery connections', () => {
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
