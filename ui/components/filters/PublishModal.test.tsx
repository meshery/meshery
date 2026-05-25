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
  publishCatalogItemSchema: {},
  publishCatalogItemUiSchema: {},
}));

vi.mock('../shared/Modal/Modal', () => ({
  RJSFModalWrapper: ({ submitBtnText, handleSubmit }: any) => (
    <button data-testid="submit" onClick={() => handleSubmit({ ok: true })}>
      {submitBtnText}
    </button>
  ),
}));

vi.mock('../../public/static/img/drawer-icons/filter_svg', () => ({
  default: () => <svg data-testid="filter-icon" />,
}));

import PublishModal from './PublishModal';

describe('PublishModal', () => {
  it('renders the publish modal with provided title and submit button label', () => {
    render(<PublishModal title="Publish foo" handleClose={() => {}} handleSubmit={() => {}} />);
    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'Publish foo');
    expect(screen.getByTestId('submit')).toHaveTextContent('Submit for Approval');
  });

  it('calls handleSubmit when the inner submit button is clicked', async () => {
    const user = userEvent.setup();
    const handleSubmit = vi.fn();
    render(<PublishModal title="Publish bar" handleClose={() => {}} handleSubmit={handleSubmit} />);
    await user.click(screen.getByTestId('submit'));
    expect(handleSubmit).toHaveBeenCalledWith({ ok: true });
  });

  it('forwards closeModal events to handleClose', async () => {
    const user = userEvent.setup();
    const handleClose = vi.fn();
    render(<PublishModal title="Publish bar" handleClose={handleClose} handleSubmit={() => {}} />);
    await user.click(screen.getByTestId('modal-close'));
    expect(handleClose).toHaveBeenCalledTimes(1);
  });
});
