import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Modal: ({ children, open, title, closeModal }: any) =>
    open ? (
      <div data-testid="modal" data-title={title}>
        {children}
        <button data-testid="modal-close" onClick={closeModal}>
          close
        </button>
      </div>
    ) : null,
  importFilterSchema: { type: 'object' },
  importFilterUiSchema: {},
}));

vi.mock('../shared/Modal/Modal', () => ({
  RJSFModalWrapper: ({ submitBtnText, handleSubmit, handleClose }: any) => (
    <div data-testid="rjsf-modal">
      <button data-testid="submit" onClick={() => handleSubmit({ foo: 'bar' })}>
        {submitBtnText}
      </button>
      <button data-testid="rjsf-close" onClick={handleClose}>
        cancel
      </button>
    </div>
  ),
}));

vi.mock('../../public/static/img/drawer-icons/filter_svg', () => ({
  default: () => <svg data-testid="filter-icon" />,
}));

import ImportModal from './ImportModal';

describe('ImportModal', () => {
  it('renders the modal with title "Import Design"', () => {
    render(<ImportModal handleClose={() => {}} handleImportFilter={() => {}} />);
    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'Import Design');
    expect(screen.getByTestId('rjsf-modal')).toBeInTheDocument();
  });

  it('forwards handleImportFilter to the RJSF handleSubmit and handleClose for cancel', async () => {
    const user = userEvent.setup();
    const handleClose = vi.fn();
    const handleImportFilter = vi.fn();
    render(<ImportModal handleClose={handleClose} handleImportFilter={handleImportFilter} />);

    await user.click(screen.getByTestId('submit'));
    expect(handleImportFilter).toHaveBeenCalledWith({ foo: 'bar' });

    await user.click(screen.getByTestId('rjsf-close'));
    expect(handleClose).toHaveBeenCalledTimes(1);

    await user.click(screen.getByTestId('modal-close'));
    expect(handleClose).toHaveBeenCalledTimes(2);
  });
});
