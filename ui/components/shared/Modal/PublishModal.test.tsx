import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import PublishModal from '../../designs/PublishDesignModal';

let modalProps: any = null;

vi.mock('@/components/shared/Modal', () => ({
  FormModal: (props: any) => {
    modalProps = props;
    return (
      <div data-testid="rjsf-modal" data-title={props.title} data-submit-btn={props.submitText}>
        <button type="button" onClick={() => props.onSubmit?.({ ok: true })}>
          {props.submitText}
        </button>
        <button type="button" onClick={props.onClose}>
          close
        </button>
      </div>
    );
  },
}));

vi.mock('@sistent/sistent', () => ({
  publishCatalogItemSchema: { type: 'object', __mockId: 'publish' },
  publishCatalogItemUiSchema: { __mockUi: 'publish' },
}));

vi.mock('../../designs/design-modal-header', () => ({
  DesignModalHeaderIcon: () => <svg data-testid="design-modal-header" />,
}));

describe('PublishModal', () => {
  beforeEach(() => {
    modalProps = null;
  });

  it('renders the publish modal with the fallback schema', () => {
    render(<PublishModal title="Publish Design" handleClose={vi.fn()} handleSubmit={vi.fn()} />);

    expect(screen.getByTestId('rjsf-modal')).toHaveAttribute('data-title', 'Publish Design');
    expect(screen.getByTestId('rjsf-modal')).toHaveAttribute(
      'data-submit-btn',
      'Submit for Approval',
    );
    expect(modalProps.schema).toEqual({ type: 'object', __mockId: 'publish' });
  });

  it('uses the provided publish form schema when supplied', () => {
    const publishFormSchema = {
      rjsfSchema: { type: 'object', __mockId: 'custom' },
      uiSchema: { __mockUi: 'custom' },
    };

    render(
      <PublishModal
        title="Publish Design"
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
        publishFormSchema={publishFormSchema}
      />,
    );

    expect(modalProps.schema).toEqual({ type: 'object', __mockId: 'custom' });
    expect(modalProps.uiSchema).toEqual({ __mockUi: 'custom' });
  });

  it('calls handleSubmit and handleClose', async () => {
    const user = userEvent.setup();
    const handleSubmit = vi.fn();
    const handleClose = vi.fn();

    render(
      <PublishModal title="Publish Design" handleClose={handleClose} handleSubmit={handleSubmit} />,
    );

    await user.click(screen.getByRole('button', { name: 'Submit for Approval' }));
    await user.click(screen.getByRole('button', { name: 'close' }));

    expect(handleSubmit).toHaveBeenCalledWith({ ok: true });
    expect(handleClose).toHaveBeenCalledTimes(1);
  });
});
