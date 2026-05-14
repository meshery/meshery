import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  importDesignSchema: { type: 'object', __mockId: 'import' },
  importDesignUiSchema: { __mockUi: 'import' },
  publishCatalogItemSchema: { type: 'object', __mockId: 'publish' },
  publishCatalogItemUiSchema: { __mockUi: 'publish' },
}));

vi.mock('@/components/shared/Modal', () => ({
  FormModal: ({ title, schema, submitText, onSubmit, onClose, helpText }: any) => (
    <div data-testid="form-modal" data-title={title}>
      <button onClick={() => onSubmit?.({ ok: true })} type="button">
        {submitText}
      </button>
      <button onClick={onClose} type="button">
        close-wrapper
      </button>
      <span data-testid="schema-id">{schema.__mockId}</span>
      {helpText ? <span data-testid="help-text">{helpText}</span> : null}
    </div>
  ),
}));

vi.mock('../design-modal-header', () => ({
  DesignModalHeaderIcon: () => <svg data-testid="pattern-svg" />,
}));

import ImportDesignModal from '../ImportDesignModal';
import PublishModal from '../PublishDesignModal';

describe('ImportDesignModal', () => {
  it('renders the Import Design modal with the import schema', () => {
    render(<ImportDesignModal handleClose={vi.fn()} handleImportDesign={vi.fn()} />);
    expect(screen.getByTestId('form-modal')).toHaveAttribute('data-title', 'Import Design');
    expect(screen.getByTestId('schema-id')).toHaveTextContent('import');
  });

  it('invokes handleImportDesign when submitting the form', () => {
    const handleImportDesign = vi.fn();
    render(<ImportDesignModal handleClose={vi.fn()} handleImportDesign={handleImportDesign} />);
    screen.getByRole('button', { name: 'Import' }).click();
    expect(handleImportDesign).toHaveBeenCalledWith({ ok: true });
  });
});

describe('PublishModal', () => {
  it('renders the Publish modal with the supplied title and publish schema', () => {
    render(<PublishModal handleClose={vi.fn()} handleSubmit={vi.fn()} title="Publish X" />);
    expect(screen.getByTestId('schema-id')).toHaveTextContent('publish');
    expect(screen.getByTestId('help-text')).toBeInTheDocument();
  });

  it('invokes handleSubmit when submitting the form', () => {
    const handleSubmit = vi.fn();
    render(<PublishModal handleClose={vi.fn()} handleSubmit={handleSubmit} title="t" />);
    screen.getByRole('button', { name: 'Submit for Approval' }).click();
    expect(handleSubmit).toHaveBeenCalledWith({ ok: true });
  });
});
