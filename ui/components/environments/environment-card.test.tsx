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

vi.mock('@sistent/sistent', () => ({
  DeleteIcon: () => <svg data-testid="delete-icon" />,
  EditIcon: () => <svg data-testid="edit-icon" />,
  Grid2: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  SyncAltIcon: () => <svg data-testid="sync-icon" />,
  useTheme: () => ({
    palette: {
      background: { neutral: { default: 'neutral' } },
    },
  }),
}));

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
  PopupButton: ({ children, onClick, disabled }: any) => (
    <button data-testid="popup-button" onClick={onClick} disabled={disabled}>
      {children}
    </button>
  ),
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
    render(<TransferButton title="Connections" count={4} onAssign={onAssign} disabled={false} />);

    expect(screen.getByTestId('tab-count')).toHaveTextContent('4');
    expect(screen.getByTestId('tab-title')).toHaveTextContent('Connections');
    await user.click(screen.getByTestId('popup-button'));
    expect(onAssign).toHaveBeenCalledTimes(1);
  });

  it('respects the disabled prop', async () => {
    const user = userEvent.setup();
    const onAssign = vi.fn();
    render(<TransferButton title="Connections" count={1} onAssign={onAssign} disabled={true} />);
    const btn = screen.getByTestId('popup-button');
    expect(btn).toBeDisabled();
    await user.click(btn);
    expect(onAssign).not.toHaveBeenCalled();
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
    // The edit and delete buttons should all be disabled when this env is selected
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
    // Edit button is the first non-popup that contains an EditIcon.
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
