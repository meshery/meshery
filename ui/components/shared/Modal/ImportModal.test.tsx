import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ImportModal from '../../designs/ImportDesignModal';

let lastFormModalProps: any = null;

vi.mock('@/components/shared/Modal', () => ({
  FormModal: (props: any) => {
    lastFormModalProps = props;
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
  importDesignSchema: { title: 'Import', __mockId: 'import' },
  importDesignUiSchema: { __mockUi: 'import' },
}));

vi.mock('../../designs/design-modal-header', () => ({
  DesignModalHeaderIcon: () => <svg data-testid="design-modal-header" />,
}));

describe('ImportModal', () => {
  beforeEach(() => {
    lastFormModalProps = null;
  });

  it('renders the design import modal with the current schema', () => {
    render(<ImportModal handleClose={vi.fn()} handleImportDesign={vi.fn()} />);

    expect(screen.getByTestId('rjsf-modal')).toHaveAttribute('data-title', 'Import Design');
    expect(screen.getByTestId('rjsf-modal')).toHaveAttribute('data-submit-btn', 'Import');
    expect(lastFormModalProps.schema).toEqual({ title: 'Import', __mockId: 'import' });
  });

  it('calls handleImportDesign when submitting', async () => {
    const user = userEvent.setup();
    const handleImportDesign = vi.fn();

    render(<ImportModal handleClose={vi.fn()} handleImportDesign={handleImportDesign} />);

    await user.click(screen.getByRole('button', { name: 'Import' }));
    expect(handleImportDesign).toHaveBeenCalledWith({ ok: true });
  });

  it('calls handleClose when closed', async () => {
    const user = userEvent.setup();
    const handleClose = vi.fn();

    render(<ImportModal handleClose={handleClose} handleImportDesign={vi.fn()} />);

    await user.click(screen.getByRole('button', { name: 'close' }));
    expect(handleClose).toHaveBeenCalledTimes(1);
  });
});
