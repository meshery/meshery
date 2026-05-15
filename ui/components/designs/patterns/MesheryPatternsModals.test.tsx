import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  importDesignSchema: { type: 'object', __mockId: 'import' },
  importDesignUiSchema: { __mockUi: 'import' },
  publishCatalogItemSchema: { type: 'object', __mockId: 'publish' },
  publishCatalogItemUiSchema: { __mockUi: 'publish' },
  Modal: ({ title, children, closeModal, 'data-testid': testId }: any) => (
    <div data-testid={testId || 'modal'} data-title={title}>
      <button type="button" onClick={closeModal}>
        close
      </button>
      {children}
    </div>
  ),
}));

vi.mock('../../shared/Modal/Modal', () => ({
  RJSFModalWrapper: ({ submitBtnText, schema, handleSubmit, handleClose, helpText }: any) => (
    <div data-testid="rjsf-modal-wrapper">
      <button onClick={() => handleSubmit?.({ ok: true })} type="button">
        {submitBtnText}
      </button>
      <button onClick={handleClose} type="button">
        close-wrapper
      </button>
      <span data-testid="schema-id">{schema.__mockId}</span>
      {helpText ? <span data-testid="help-text">{helpText}</span> : null}
    </div>
  ),
}));

vi.mock('../../../public/static/img/drawer-icons/pattern_svg', () => ({
  default: () => <svg data-testid="pattern-svg" />,
}));

import { ImportDesignModal, PublishModal } from './MesheryPatternsModals';

describe('ImportDesignModal', () => {
  it('renders the Import Design modal with the import schema', () => {
    render(<ImportDesignModal handleClose={vi.fn()} handleImportDesign={vi.fn()} />);
    expect(screen.getByTestId('import-design-modal')).toHaveAttribute(
      'data-title',
      'Import Design',
    );
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
