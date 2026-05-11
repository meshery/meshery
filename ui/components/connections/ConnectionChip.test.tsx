import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  ConnectionChip,
  ConnectionStateChip,
  TooltipWrappedConnectionChip,
} from './ConnectionChip';

const normalizeStaticImagePath = vi.fn((src?: string) => src);

vi.mock('@sistent/sistent', () => {
  const styled = (Component) => () => {
    const StyledComponent = ({ children, ...props }) => (
      <Component {...props}>{children}</Component>
    );
    StyledComponent.displayName = 'StyledSistentMock';
    return StyledComponent;
  };

  return {
    Avatar: ({ children, src, sx }) => (
      <div data-testid="avatar" data-src={src} data-size={JSON.stringify(sx || {})}>
        {children}
      </div>
    ),
    AssignmentTurnedInIcon: () => <svg data-testid="assignment-icon" />,
    CustomTooltip: ({ title, children }) => (
      <div data-testid="tooltip" data-title={String(title)}>
        {children}
      </div>
    ),
    Typography: ({ children }) => <span>{children}</span>,
    styled,
    createTheme: () => ({ breakpoints: {} }),
    useTheme: () => ({
      palette: {
        background: {
          brand: { default: 'brand' },
          warning: { default: 'warning' },
        },
        text: { disabled: 'disabled' },
      },
    }),
  };
});

vi.mock('@/utils/fallback', () => ({
  normalizeStaticImagePath: (...args) => normalizeStaticImagePath(...args),
}));

vi.mock('../../themes', () => ({
  notificationColors: {
    lightwarning: 'warning',
    info: 'info',
  },
}));

vi.mock('../CustomAvatar', () => ({
  default: ({ children, color }) => (
    <div data-testid="badge-avatar" data-color={color}>
      {children}
    </div>
  ),
}));

vi.mock('./styles', () => {
  const makeStateChip = (testId: string) => {
    const StateChip = ({ avatar, label }) => (
      <div data-testid={testId}>
        {avatar}
        <span>{label}</span>
      </div>
    );

    StateChip.displayName = `${testId}Mock`;
    return StateChip;
  };

  return {
    ChipWrapper: ({ label, onClick, onDelete, disabled, avatar, style }) => (
      <div data-testid="chip-wrapper" data-width={style?.width}>
        {avatar}
        <button onClick={onClick} disabled={disabled}>
          {label}
        </button>
        {onDelete ? <button onClick={onDelete}>delete</button> : null}
      </div>
    ),
    ConnectedChip: makeStateChip('connected-state-chip'),
    DeletedChip: makeStateChip('deleted-state-chip'),
    DisconnectedChip: makeStateChip('disconnected-state-chip'),
    DiscoveredChip: makeStateChip('discovered-state-chip'),
    IgnoredChip: makeStateChip('ignored-state-chip'),
    MaintainanceChip: makeStateChip('maintenance-state-chip'),
    NotFoundChip: makeStateChip('not-found-state-chip'),
    RegisteredChip: makeStateChip('registered-state-chip'),
  };
});

vi.mock('css/icons.styles', () => ({
  iconMedium: {},
  iconSmall: {},
}));

vi.mock('@/assets/icons/Connection', () => ({
  default: () => <svg data-testid="connection-icon" />,
}));

vi.mock('../../assets/icons/disconnect', () => ({
  default: () => <svg data-testid="disconnect-icon" />,
}));

describe('ConnectionChip', () => {
  beforeEach(() => {
    normalizeStaticImagePath.mockClear();
  });

  it('calls handlePing when clicked and renders a normalized avatar source', async () => {
    const user = userEvent.setup();
    const handlePing = vi.fn();

    render(
      <ConnectionChip
        title="cluster-a"
        handlePing={handlePing}
        iconSrc="/static/img/kubernetes.svg"
        status="connected"
        width="12rem"
      />,
    );

    await user.click(screen.getByRole('button', { name: 'cluster-a' }));

    expect(handlePing).toHaveBeenCalledTimes(1);
    expect(normalizeStaticImagePath).toHaveBeenCalledWith('/static/img/kubernetes.svg');
    expect(screen.getByTestId('avatar')).toHaveAttribute('data-src', '/static/img/kubernetes.svg');
    expect(screen.getByTestId('chip-wrapper')).toHaveAttribute('data-width', '12rem');
    expect(screen.getByTestId('badge-avatar')).toHaveAttribute('data-color', 'brand');
  });

  it('respects disabled chips by blocking ping and delete actions', async () => {
    const user = userEvent.setup();
    const handlePing = vi.fn();
    const onDelete = vi.fn();

    render(
      <ConnectionChip
        title="operator"
        handlePing={handlePing}
        onDelete={onDelete}
        disabled={true}
        status="running"
      />,
    );

    const button = screen.getByRole('button', { name: 'operator' });
    expect(button).toBeDisabled();

    await user.click(button);

    expect(handlePing).not.toHaveBeenCalled();
    expect(screen.queryByRole('button', { name: 'delete' })).not.toBeInTheDocument();
  });
});

describe('TooltipWrappedConnectionChip', () => {
  it('prefers explicit tooltip content over the title', () => {
    render(<TooltipWrappedConnectionChip title="cluster-a" tooltip="Server: demo" />);

    expect(screen.getByTestId('tooltip')).toHaveAttribute('data-title', 'Server: demo');
  });
});

describe('ConnectionStateChip', () => {
  it('maps actionable states to transition labels', () => {
    render(<ConnectionStateChip status="connected" actionable={true} />);

    expect(screen.getByTestId('connected-state-chip')).toHaveTextContent('Connect');
  });

  it('renders known states and falls back to the discovered chip for unknown statuses', () => {
    const { rerender } = render(<ConnectionStateChip status="not found" />);
    expect(screen.getByTestId('not-found-state-chip')).toHaveTextContent('not found');

    rerender(<ConnectionStateChip status="mystery-state" />);
    expect(screen.getByTestId('discovered-state-chip')).toHaveTextContent('mystery-state');
  });
});
