import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/components/shared/Modal', () => ({
  Modal: ({ isOpen, onClose, title, children }: any) =>
    isOpen ? (
      <div data-testid="modal" data-title={title}>
        {children}
        <button data-testid="close-btn" onClick={onClose}>
          close
        </button>
      </div>
    ) : null,
}));

vi.mock('./RelationshipFormStepper', () => ({
  default: ({ handleClose }: any) => (
    <div data-testid="stepper">
      <button data-testid="trigger-close" onClick={handleClose}>
        finish
      </button>
    </div>
  ),
}));

import CreateRelationshipModal from '../registry/CreateRelationshipModal';

describe('CreateRelationshipModal', () => {
  it('does not render the modal when closed', () => {
    render(
      <CreateRelationshipModal
        isRelationshipModalOpen={false}
        setIsRelationshipModalOpen={vi.fn()}
      />,
    );
    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });

  it('renders the stepper inside the modal when open', () => {
    render(
      <CreateRelationshipModal
        isRelationshipModalOpen={true}
        setIsRelationshipModalOpen={vi.fn()}
      />,
    );
    expect(screen.getByTestId('modal')).toBeInTheDocument();
    expect(screen.getByTestId('stepper')).toBeInTheDocument();
    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'Create Relationship');
  });

  it('calls setIsRelationshipModalOpen(false) when the modal is closed', () => {
    const setOpen = vi.fn();
    render(
      <CreateRelationshipModal
        isRelationshipModalOpen={true}
        setIsRelationshipModalOpen={setOpen}
      />,
    );
    fireEvent.click(screen.getByTestId('close-btn'));
    expect(setOpen).toHaveBeenCalledWith(false);
  });

  it('passes a handleClose to the stepper that closes the modal', () => {
    const setOpen = vi.fn();
    render(
      <CreateRelationshipModal
        isRelationshipModalOpen={true}
        setIsRelationshipModalOpen={setOpen}
      />,
    );
    fireEvent.click(screen.getByTestId('trigger-close'));
    expect(setOpen).toHaveBeenCalledWith(false);
  });
});
