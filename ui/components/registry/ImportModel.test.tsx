import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const importModelReq = vi.fn();

vi.mock('@sistent/sistent', () => ({
  Modal: ({ children, open, title }: any) =>
    open ? <div data-testid={`modal-${title}`}>{children}</div> : null,
  FormControlLabel: ({ label, value, control }: any) => (
    <label data-radio-value={value}>
      {control}
      {label}
    </label>
  ),
  FormControl: ({ children }: any) => <div>{children}</div>,
  RadioGroup: ({ children, value, onChange }: any) => (
    <div data-testid="radio-group" data-value={value} onChange={onChange}>
      {children}
    </div>
  ),
  Radio: () => <input type="radio" />,
  Typography: ({ children }: any) => <span>{children}</span>,
  importModelUiSchema: { uploadType: {} },
  importModelSchema: {
    title: 'Import',
    type: 'object',
    properties: { uploadType: { enum: ['File Import', 'URL Import', 'CSV Import'] } },
  },
}));

vi.mock('../shared/Modal/Modal', () => ({
  RJSFModalWrapper: ({ handleSubmit, schema, widgets, helpText }: any) => (
    <div data-testid="rjsf-form">
      {helpText}
      <button
        data-testid="submit-file"
        onClick={() =>
          handleSubmit({ uploadType: 'File Import', file: 'data:application/x-tar;base64,xxx' })
        }
      >
        submit file
      </button>
      <button
        data-testid="submit-url"
        onClick={() => handleSubmit({ uploadType: 'URL Import', url: 'https://x.com' })}
      >
        submit url
      </button>
      <button
        data-testid="submit-empty-url"
        onClick={() => handleSubmit({ uploadType: 'URL Import', url: '' })}
      >
        submit empty url
      </button>
      <button data-testid="submit-csv" onClick={() => handleSubmit({ uploadType: 'CSV Import' })}>
        submit csv
      </button>
      <button data-testid="submit-invalid" onClick={() => handleSubmit({ uploadType: 'Invalid' })}>
        submit invalid
      </button>
    </div>
  ),
}));

vi.mock('./Stepper/CSVStepper', () => ({
  default: () => <div data-testid="csv-stepper" />,
  StyledDocsRedirectLink: ({ children, href }: any) => <a href={href}>{children}</a>,
}));

vi.mock('@/constants/endpoints', () => ({
  MESHERY_DOCS_URL: 'https://docs.meshery.io',
}));

vi.mock('@/utils/utils', () => ({
  getUnit8ArrayDecodedFile: (data: any) => (data ? [1, 2, 3] : null),
}));

vi.mock('@/rtk-query/meshModel', () => ({
  useImportMeshModelMutation: () => [importModelReq],
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: vi.fn(),
}));

import ImportModelModal from './ImportModel';

describe('ImportModel (ImportModelModal)', () => {
  beforeEach(() => {
    importModelReq.mockReset();
  });

  it('renders the import modal when open', () => {
    render(<ImportModelModal isImportModalOpen={true} setIsImportModalOpen={vi.fn()} />);
    expect(screen.getByTestId('modal-Import Model')).toBeInTheDocument();
  });

  it('does not render the modal when closed', () => {
    render(<ImportModelModal isImportModalOpen={false} setIsImportModalOpen={vi.fn()} />);
    expect(screen.queryByTestId('modal-Import Model')).not.toBeInTheDocument();
  });

  it('opens the CSV modal and closes the import modal on CSV Import submit', () => {
    const setIsImportModalOpen = vi.fn();
    // Need a file element for File Import path; not relevant here but
    // also assert CSV path works.
    render(
      <ImportModelModal isImportModalOpen={true} setIsImportModalOpen={setIsImportModalOpen} />,
    );

    fireEvent.click(screen.getByTestId('submit-csv'));

    expect(setIsImportModalOpen).toHaveBeenCalledWith(false);
  });

  it('calls importModelReq for URL Import submit', () => {
    render(<ImportModelModal isImportModalOpen={true} setIsImportModalOpen={vi.fn()} />);

    fireEvent.click(screen.getByTestId('submit-url'));

    expect(importModelReq).toHaveBeenCalledWith({
      importBody: expect.objectContaining({
        importBody: { url: 'https://x.com' },
        uploadType: 'urlImport',
      }),
    });
  });

  it('does not call importModelReq for empty URL Import submit', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    render(<ImportModelModal isImportModalOpen={true} setIsImportModalOpen={vi.fn()} />);

    fireEvent.click(screen.getByTestId('submit-empty-url'));

    expect(importModelReq).not.toHaveBeenCalled();
    consoleSpy.mockRestore();
  });

  it('does not call importModelReq for invalid upload type', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    render(<ImportModelModal isImportModalOpen={true} setIsImportModalOpen={vi.fn()} />);

    fireEvent.click(screen.getByTestId('submit-invalid'));

    expect(importModelReq).not.toHaveBeenCalled();
    consoleSpy.mockRestore();
  });
});
