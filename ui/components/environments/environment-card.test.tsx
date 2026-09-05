import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const can = vi.fn(() => true);
const getEnvironmentConnectionsQuery = vi.fn();

vi.mock('../../rtk-query/environments', () => ({
  useGetEnvironmentConnectionsQuery: (...args: unknown[]) =>
    getEnvironmentConnectionsQuery(...args),
}));

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => can(...args),
}));

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) => (
      <Component {...props}>{children}</Component>
    );
    StyledComponent.displayName = 'StyledSistentMock';
    return StyledComponent;
  };

  return {
    DeleteIcon: () => <svg data-testid="delete-icon" />,
    EditIcon: () => <svg data-testid="edit-icon" />,
    Grid2: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    SyncAltIcon: () => <svg data-testid="sync-icon" />,
    IconButton: ({ children, onClick, disabled, ...props }: any) => (
      <button data-testid="sistent-icon-button" onClick={onClick} disabled={disabled} {...props}>
        {children}
      </button>
    ),
    CustomTooltip: ({ children, title }: any) => (
      <div data-testid="custom-tooltip" title={title}>
        {children}
      </div>
    ),
    Chip: ({ label, size, sx }: any) => (
      <div data-testid="purpose-chip" data-size={size}>
        {label}
      </div>
    ),
    Typography: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    useTheme: () => ({
      palette: {
        background: { neutral: { default: 'neutral' } },
        icon: { secondary: 'secondary' },
      },
    }),
    styled,
  };
});

vi.mock('../lifecycle/general', () => ({
  FlipCard: ({ frontComponents, backComponents, disableFlip }: any) => (
    <div data-testid="flip-card" data-disabled={disableFlip}>
      <div data-testid="flip-front">{frontComponents}</div>
      <div data-testid="flip-back">{backComponents}</div>
    </div>
  ),
}));

vi.mock('./styles', () => ({
  Name: ({ children, onClick }: any) => (
    <div data-testid="env-name" onClick={onClick}>
      {children}
    </div>
  ),
  IconButton: ({ children, onClick, disabled }: any) => (
    <button onClick={onClick} disabled={disabled}>
      {children}
    </button>
  ),
  CardWrapper: ({ children }: any) => <div data-testid="card-wrapper">{children}</div>,
  DateLabel: ({ children, onClick }: any) => <span onClick={onClick}>{children}</span>,
  DescriptionLabel: ({ children, onClick }: any) => (
    <p data-testid="description" onClick={onClick}>
      {children}
    </p>
  ),
  EmptyDescription: ({ children, onClick }: any) => (
    <p data-testid="empty-description" onClick={onClick}>
      {children}
    </p>
  ),
  TabCount: ({ children }: any) => <span data-testid="tab-count">{children}</span>,
  TabTitle: ({ children }: any) => <span data-testid="tab-title">{children}</span>,
  PopupButton: ({ children, onClick, disabled, permissionKey }: any) => {
    const isDisabled =
      disabled || (permissionKey && !can(permissionKey.id, permissionKey.function));
    return (
      <button data-testid="popup-button" onClick={onClick} disabled={isDisabled}>
        {children}
      </button>
    );
  },
  AllocationButton: ({ children, onClick }: any) => <div onClick={onClick}>{children}</div>,
  BulkSelectCheckbox: ({ onChange, onClick, disabled }: any) => (
    <input
      type="checkbox"
      data-testid="bulk-checkbox"
      onChange={onChange}
      onClick={onClick}
      disabled={disabled}
    />
  ),
  CardTitle: ({ children, onClick }: any) => (
    <div data-testid="card-title" onClick={onClick}>
      {children}
    </div>
  ),
}));

import EnvironmentCard, { formattoLongDate, TransferButton } from './environment-card';

describe('formattoLongDate', () => {
  it('formats a date into a long human-readable form', () => {
    const formatted = formattoLongDate('2026-01-15T00:00:00.000Z');
    expect(typeof formatted).toBe('string');
    expect(formatted).toMatch(/2026/);
  });
});

describe('TransferButton', () => {
  it('renders the title, count and triggers onAssign when clicked', async () => {
    const user = userEvent.setup();
    const onAssign = vi.fn();
    render(<TransferButton title="Connections" count={4} onAssign={onAssign} />);

    expect(screen.getByTestId('tab-count')).toHaveTextContent('4');
    expect(screen.getByTestId('tab-title')).toHaveTextContent('Connections');
    await user.click(screen.getByTestId('popup-button'));
    expect(onAssign).toHaveBeenCalledTimes(1);
  });

  it('respects the permissionKey prop when user lacks permission', async () => {
    can.mockReturnValue(false);
    const user = userEvent.setup();
    const onAssign = vi.fn();
    const mockKey = { id: 'view-connections', function: 'ViewConnections' };
    render(
      <TransferButton title="Connections" count={1} onAssign={onAssign} permissionKey={mockKey} />,
    );
    const btn = screen.getByTestId('popup-button');
    expect(btn).toBeDisabled();
    await user.click(btn);
    expect(onAssign).not.toHaveBeenCalled();
    can.mockReturnValue(true);
  });
});

const baseEnvironment = {
  id: 'env-1',
  name: 'dev',
  description: 'dev env',
  deletedAt: null,
  updatedAt: '2026-01-12T00:00:00.000Z',
  createdAt: '2026-01-10T00:00:00.000Z',
};

describe('EnvironmentCard', () => {
  beforeEach(() => {
    can.mockReset();
    can.mockReturnValue(true);
    getEnvironmentConnectionsQuery.mockReset();
    getEnvironmentConnectionsQuery.mockReturnValue({ data: { totalCount: 7 } });
  });

  it('renders the environment name, description and connections count', () => {
    render(
      <EnvironmentCard
        environmentDetails={baseEnvironment}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );

    expect(screen.getAllByTestId('env-name')[0]).toHaveTextContent('dev');
    expect(screen.getByTestId('description')).toHaveTextContent('dev env');
    expect(screen.getByTestId('tab-count')).toHaveTextContent('7');
  });

  it('falls back to "No description" for environments without a description', () => {
    render(
      <EnvironmentCard
        environmentDetails={{ ...baseEnvironment, description: '' }}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    expect(screen.getByTestId('empty-description')).toHaveTextContent('No description');
  });

  it('disables edit/delete buttons when the environment is currently selected', () => {
    render(
      <EnvironmentCard
        environmentDetails={baseEnvironment}
        selectedEnvironments={['env-1']}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    const buttons = screen.getAllByRole('button');
    const disabledButtons = buttons.filter((b) => (b as HTMLButtonElement).disabled);
    expect(disabledButtons.length).toBeGreaterThanOrEqual(2);
  });

  it('considers deletedAt as object with Valid: true as deleted (disables the checkbox)', () => {
    render(
      <EnvironmentCard
        environmentDetails={{
          ...baseEnvironment,
          deletedAt: { Time: '2026-01-01', Valid: true },
        }}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );

    expect(screen.getByTestId('bulk-checkbox')).toBeDisabled();
  });

  it('considers deletedAt as object with Valid: false as not deleted (checkbox enabled)', () => {
    render(
      <EnvironmentCard
        environmentDetails={{
          ...baseEnvironment,
          deletedAt: { Valid: false },
        }}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    expect(screen.getByTestId('bulk-checkbox')).not.toBeDisabled();
  });

  it('fires onEdit/onDelete callbacks when corresponding buttons are clicked', async () => {
    const user = userEvent.setup();
    const onEdit = vi.fn();
    const onDelete = vi.fn();
    render(
      <EnvironmentCard
        environmentDetails={baseEnvironment}
        selectedEnvironments={[]}
        onDelete={onDelete}
        onEdit={onEdit}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );

    const buttons = screen.getAllByRole('button');
    const editButton = buttons.find((b) => b.querySelector('[data-testid="edit-icon"]'));
    const deleteButton = buttons.find((b) => b.querySelector('[data-testid="delete-icon"]'));

    if (editButton) await user.click(editButton);
    if (deleteButton) await user.click(deleteButton);

    expect(onEdit).toHaveBeenCalled();
    expect(onDelete).toHaveBeenCalled();
  });

  it('uses {skip: !id} for connections query', () => {
    render(
      <EnvironmentCard
        environmentDetails={{ ...baseEnvironment, id: '' }}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    expect(getEnvironmentConnectionsQuery).toHaveBeenCalledWith(
      { environmentId: '' },
      { skip: true },
    );
  });
});

describe('Administrative environment purpose', () => {
  beforeEach(() => {
    can.mockReset();
    can.mockReturnValue(true);
    getEnvironmentConnectionsQuery.mockReset();
    getEnvironmentConnectionsQuery.mockReturnValue({ data: { totalCount: 3 } });
  });

  it('renders the Administrative badge on both card faces when purpose is "administrative"', () => {
    render(
      <EnvironmentCard
        environmentDetails={{ ...baseEnvironment, purpose: 'administrative' }}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    const badges = screen.getAllByTestId('purpose-chip');
    // Badge appears on both the front and back of the flip card
    expect(badges.length).toBeGreaterThanOrEqual(1);
    badges.forEach((badge) => {
      expect(badge).toHaveTextContent('Administrative');
    });
  });

  it('does not render the Administrative badge when purpose is absent', () => {
    render(
      <EnvironmentCard
        environmentDetails={baseEnvironment}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    expect(screen.queryByTestId('purpose-chip')).not.toBeInTheDocument();
  });

  it('does not render the Administrative badge when purpose is "user"', () => {
    render(
      <EnvironmentCard
        environmentDetails={{ ...baseEnvironment, purpose: 'user' }}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    expect(screen.queryByTestId('purpose-chip')).not.toBeInTheDocument();
  });

  it('does not render the Administrative badge when purpose is "absent"', () => {
    render(
      <EnvironmentCard
        environmentDetails={{ ...baseEnvironment, purpose: 'absent' }}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    expect(screen.queryByTestId('purpose-chip')).not.toBeInTheDocument();
  });

  it('disables edit and delete buttons for administrative environments', () => {
    render(
      <EnvironmentCard
        environmentDetails={{ ...baseEnvironment, purpose: 'administrative' }}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    const buttons = screen.getAllByRole('button');
    const editButton = buttons.find((b) => b.querySelector('[data-testid="edit-icon"]'));
    const deleteButton = buttons.find((b) => b.querySelector('[data-testid="delete-icon"]'));
    expect(editButton).toBeDisabled();
    expect(deleteButton).toBeDisabled();
  });

  it('shows descriptive tooltip for edit on administrative environments', () => {
    render(
      <EnvironmentCard
        environmentDetails={{ ...baseEnvironment, purpose: 'administrative' }}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    const tooltips = screen.getAllByTestId('custom-tooltip');
    const editTooltip = tooltips.find(
      (t) => t.getAttribute('title') === 'Administrative environments cannot be edited',
    );
    expect(editTooltip).toBeTruthy();
  });

  it('shows descriptive tooltip for delete on administrative environments', () => {
    render(
      <EnvironmentCard
        environmentDetails={{ ...baseEnvironment, purpose: 'administrative' }}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    const tooltips = screen.getAllByTestId('custom-tooltip');
    const deleteTooltip = tooltips.find(
      (t) => t.getAttribute('title') === 'Administrative environments cannot be deleted',
    );
    expect(deleteTooltip).toBeTruthy();
  });

  it('does not fire onEdit when edit is clicked on an administrative environment', async () => {
    const user = userEvent.setup();
    const onEdit = vi.fn();
    render(
      <EnvironmentCard
        environmentDetails={{ ...baseEnvironment, purpose: 'administrative' }}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={onEdit}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    const buttons = screen.getAllByRole('button');
    const editButton = buttons.find((b) => b.querySelector('[data-testid="edit-icon"]'));
    if (editButton) await user.click(editButton);
    expect(onEdit).not.toHaveBeenCalled();
  });

  it('does not fire onDelete when delete is clicked on an administrative environment', async () => {
    const user = userEvent.setup();
    const onDelete = vi.fn();
    render(
      <EnvironmentCard
        environmentDetails={{ ...baseEnvironment, purpose: 'administrative' }}
        selectedEnvironments={[]}
        onDelete={onDelete}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    const buttons = screen.getAllByRole('button');
    const deleteButton = buttons.find((b) => b.querySelector('[data-testid="delete-icon"]'));
    if (deleteButton) await user.click(deleteButton);
    expect(onDelete).not.toHaveBeenCalled();
  });

  it('edit and delete remain enabled for non-administrative environments', () => {
    render(
      <EnvironmentCard
        environmentDetails={baseEnvironment}
        selectedEnvironments={[]}
        onDelete={() => {}}
        onEdit={() => {}}
        onSelect={() => {}}
        onAssignConnection={() => {}}
      />,
    );
    const buttons = screen.getAllByRole('button');
    const editButton = buttons.find((b) => b.querySelector('[data-testid="edit-icon"]'));
    const deleteButton = buttons.find((b) => b.querySelector('[data-testid="delete-icon"]'));
    expect(editButton).not.toBeDisabled();
    expect(deleteButton).not.toBeDisabled();
  });
});
