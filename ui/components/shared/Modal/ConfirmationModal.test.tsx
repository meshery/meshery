import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import ConfirmModal from './ConfirmModal';

vi.mock('./Modal', () => ({
  Modal: ({ isOpen, onClose, title, children, actions }: any) =>
    isOpen ? (
      <div data-testid="modal" data-title={title}>
        <button type="button" data-testid="modal-close" onClick={onClose}>
          close
        </button>
        {children}
        {actions}
      </div>
    ) : null,
}));

vi.mock('@sistent/sistent', () => ({
  ModalButtonPrimary: ({ children, onClick, disabled }: any) => (
    <button type="button" onClick={onClick} disabled={disabled}>
      {children}
    </button>
  ),
  ModalButtonSecondary: ({ children, onClick }: any) => (
    <button type="button" onClick={onClick}>
      {children}
    </button>
  ),
  ModalButtonDanger: ({ children, onClick, disabled }: any) => (
    <button type="button" data-testid="danger-button" onClick={onClick} disabled={disabled}>
      {children}
    </button>
  ),
  Typography: ({ children }: any) => <p>{children}</p>,
}));

describe('ConfirmModal', () => {
  it('renders the title and message when open', () => {
    render(
      <ConfirmModal
        isOpen
        onClose={vi.fn()}
        onConfirm={vi.fn()}
        title="Delete environment"
        message="This cannot be undone."
      />,
    );

    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'Delete environment');
    expect(screen.getByText('This cannot be undone.')).toBeInTheDocument();
  });

  it('invokes onClose from both the close affordances', async () => {
    const user = userEvent.setup();
    const onClose = vi.fn();

    render(
      <ConfirmModal
        isOpen
        onClose={onClose}
        onConfirm={vi.fn()}
        title="Delete environment"
        message="This cannot be undone."
      />,
    );

    await user.click(screen.getByTestId('modal-close'));
    await user.click(screen.getByRole('button', { name: 'Cancel' }));
    expect(onClose).toHaveBeenCalledTimes(2);
  });

  it('invokes onConfirm and uses the danger variant when requested', async () => {
    const user = userEvent.setup();
    const onConfirm = vi.fn();

    render(
      <ConfirmModal
        isOpen
        onClose={vi.fn()}
        onConfirm={onConfirm}
        title="Delete environment"
        message="This cannot be undone."
        variant="danger"
        confirmText="Delete"
      />,
    );

    await user.click(screen.getByRole('button', { name: 'Delete' }));
    expect(onConfirm).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('danger-button')).toBeInTheDocument();
  });
});
