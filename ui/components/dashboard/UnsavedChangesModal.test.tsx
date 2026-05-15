import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import UnsavedChangesModal from '../shared/Modal/UnsavedChangesModal';

vi.mock('@sistent/sistent', () => ({
  CheckCircleIcon: () => <svg data-testid="check-icon" />,
  DeleteIcon: ({ fill }: any) => <svg data-testid="delete-icon" data-fill={fill} />,
  ModalButtonPrimary: ({ children, onClick, startIcon, ...props }: any) => (
    <button type="button" onClick={onClick} {...props}>
      {startIcon}
      <span>{children}</span>
    </button>
  ),
  Typography: ({ children }: any) => <span>{children}</span>,
}));

vi.mock('@/theme', () => ({
  styled:
    (Component: any) =>
    () =>
    ({ children, ...props }: any) =>
      typeof Component === 'string' ? (
        React.createElement(Component, props, children)
      ) : (
        <Component {...props}>{children}</Component>
      ),
  useTheme: () => ({
    palette: {
      common: { white: '#fff' },
      background: { error: { default: '#f00' } },
    },
  }),
}));

vi.mock('../shared/Modal/Modal', () => ({
  Modal: ({ isOpen, onClose, title, children, actions }: any) =>
    isOpen ? (
      <div data-testid="modal" data-title={title}>
        <button type="button" onClick={onClose} data-testid="modal-close">
          close
        </button>
        {children}
        {actions}
      </div>
    ) : null,
}));

describe('UnsavedChangesModal', () => {
  const onClose = vi.fn();
  const onDiscard = vi.fn();
  const onSave = vi.fn();

  beforeEach(() => {
    onClose.mockReset();
    onDiscard.mockReset();
    onSave.mockReset();
  });

  it('does not render when open is false', () => {
    render(
      <UnsavedChangesModal open={false} onClose={onClose} onDiscard={onDiscard} onSave={onSave} />,
    );
    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });

  it('renders with the title and warning copy when open', () => {
    render(
      <UnsavedChangesModal open={true} onClose={onClose} onDiscard={onDiscard} onSave={onSave} />,
    );
    expect(screen.getByTestId('modal')).toHaveAttribute(
      'data-title',
      'Unsaved dashboard layout changes',
    );
    expect(
      screen.getByText(/You have unsaved changes to your dashboard layout/i),
    ).toBeInTheDocument();
  });

  it('invokes onDiscard when the discard button is clicked', async () => {
    const user = userEvent.setup();
    render(
      <UnsavedChangesModal open={true} onClose={onClose} onDiscard={onDiscard} onSave={onSave} />,
    );
    await user.click(screen.getByRole('button', { name: /discard changes/i }));
    expect(onDiscard).toHaveBeenCalledTimes(1);
  });

  it('invokes onSave when the save button is clicked', async () => {
    const user = userEvent.setup();
    render(
      <UnsavedChangesModal open={true} onClose={onClose} onDiscard={onDiscard} onSave={onSave} />,
    );
    await user.click(screen.getByRole('button', { name: /save changes/i }));
    expect(onSave).toHaveBeenCalledTimes(1);
  });

  it('invokes onClose when the modal close action fires', async () => {
    const user = userEvent.setup();
    render(
      <UnsavedChangesModal open={true} onClose={onClose} onDiscard={onDiscard} onSave={onSave} />,
    );
    await user.click(screen.getByTestId('modal-close'));
    expect(onClose).toHaveBeenCalledTimes(1);
  });
});
