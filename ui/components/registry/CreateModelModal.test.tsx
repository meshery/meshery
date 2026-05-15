import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Modal: ({ children, open, closeModal, title }: any) =>
    open ? (
      <div data-testid="modal" data-title={title}>
        <button onClick={closeModal} data-testid="modal-close" type="button">
          x
        </button>
        {children}
      </div>
    ) : null,
}));

vi.mock('./Stepper/UrlStepper', () => ({
  default: ({ handleClose }: any) => (
    <button data-testid="url-stepper-close" onClick={handleClose} type="button">
      stepper-close
    </button>
  ),
}));

import CreateModelModal from './CreateModelModal';

describe('CreateModelModal', () => {
  it('renders the modal when open', () => {
    render(<CreateModelModal isCreateModalOpen={true} setIsCreateModalOpen={vi.fn()} />);

    expect(screen.getByTestId('modal')).toBeInTheDocument();
    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'Create Model');
    expect(screen.getByTestId('url-stepper-close')).toBeInTheDocument();
  });

  it('does not render the modal when closed', () => {
    render(<CreateModelModal isCreateModalOpen={false} setIsCreateModalOpen={vi.fn()} />);

    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });

  it('calls setIsCreateModalOpen(false) when the modal close fires', () => {
    const setIsCreateModalOpen = vi.fn();
    render(
      <CreateModelModal isCreateModalOpen={true} setIsCreateModalOpen={setIsCreateModalOpen} />,
    );

    fireEvent.click(screen.getByTestId('modal-close'));
    expect(setIsCreateModalOpen).toHaveBeenCalledWith(false);
  });

  it('calls setIsCreateModalOpen(false) when the inner stepper handleClose fires', () => {
    const setIsCreateModalOpen = vi.fn();
    render(
      <CreateModelModal isCreateModalOpen={true} setIsCreateModalOpen={setIsCreateModalOpen} />,
    );

    fireEvent.click(screen.getByTestId('url-stepper-close'));
    expect(setIsCreateModalOpen).toHaveBeenCalledWith(false);
  });
});
