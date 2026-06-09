import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const importModelReq = vi.fn();

vi.mock('@sistent/sistent', () => ({
  FormControlLabel: ({ label, control }: any) => (
    <label>
      {control}
      {label}
    </label>
  ),
  FormControl: ({ children }: any) => <div>{children}</div>,
  RadioGroup: ({ children }: any) => <div>{children}</div>,
  Radio: () => <input type="radio" />,
  Typography: ({ children }: any) => <span>{children}</span>,
  Box: ({ children }: any) => <div>{children}</div>,
  ModalButtonPrimary: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
}));

vi.mock('@/components/shared/Modal', () => ({
  Modal: ({ isOpen, title, children }: any) =>
    isOpen ? <div data-testid={`modal-${title}`}>{children}</div> : null,
  FormModal: ({ isOpen, title, helpText, onSubmit }: any) =>
    isOpen ? (
      <div data-testid={`modal-${title}`}>
        <div data-testid="rjsf-wrapper">{helpText}</div>
        <button
          data-testid="submit-url"
          onClick={() => onSubmit({ uploadType: 'urlImport', url: 'https://x.com' })}
        >
          submit-url
        </button>
        <button
          data-testid="submit-empty-url"
          onClick={() => onSubmit({ uploadType: 'urlImport', url: '' })}
        >
          submit-empty-url
        </button>
        <button data-testid="submit-csv" onClick={() => onSubmit({ uploadType: 'csv' })}>
          submit-csv
        </button>
        <button data-testid="submit-invalid" onClick={() => onSubmit({ uploadType: 'unknown' })}>
          submit-invalid
        </button>
      </div>
    ) : null,
}));

vi.mock('./Stepper/CSVStepper', () => ({
  default: () => <div data-testid="csv-stepper" />,
}));

vi.mock('@/constants/endpoints', () => ({
  MESHERY_DOCS_URL: 'https://docs.meshery.io',
}));

vi.mock('@/theme', () => ({
  useTheme: () => ({
    palette: { background: { surfaces: '#fff' } },
  }),
}));

vi.mock('@/rtk-query/meshModel', () => ({
  useImportMeshModelMutation: () => [importModelReq],
}));

vi.mock('lodash', () => ({
  capitalize: (s: string) => s,
}));

vi.mock('@/components/designs/lifecycle/common', () => ({
  Loading: ({ message }: any) => <div data-testid="loading">{message}</div>,
}));

vi.mock('@/components/layout/NotificationCenter', () => {
  const Ctx = (require('react') as typeof import('react')).createContext({
    operationsCenterActorRef: {
      on: () => ({ unsubscribe: () => {} }),
    },
  });
  return { NotificationCenterContext: Ctx };
});

vi.mock('machines/operationsCenter', () => ({
  OPERATION_CENTER_EVENTS: { EVENT_RECEIVED_FROM_SERVER: 'event' },
}));

vi.mock('@/components/layout/NotificationCenter/formatters/model_registration', () => ({
  ModelImportedSection: () => <div data-testid="imported-section" />,
  ModelImportMessages: () => <div data-testid="import-messages" />,
}));

vi.mock('./Stepper/style', () => ({
  StyledDocsRedirectLink: ({ children, href }: any) => <a href={href}>{children}</a>,
}));

vi.mock('./importModelSchema', () => ({
  UPLOAD_TYPE_CSV: 'csv',
  UPLOAD_TYPE_FILE: 'file',
  UPLOAD_TYPE_URL: 'urlImport',
  decodeDataUrlToBytes: (d: any) => (d ? [1, 2, 3] : null),
  filenameFromDataUrl: () => 'model.yml',
  findSelectedModelFile: () => null,
  importModelSchema: { properties: { uploadType: { description: 'Import help text' } } },
  importModelUiSchema: {},
}));

vi.mock('@/utils/fileUpload', () => ({
  readFileAsBytes: vi.fn(),
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: vi.fn(),
}));

import ImportModelModal from './ImportModelModal';

describe('ImportModelModal', () => {
  beforeEach(() => {
    importModelReq.mockReset();
    importModelReq.mockReturnValue({ unwrap: () => Promise.resolve({}) });
  });

  it('renders the main Import Model modal when open', () => {
    render(<ImportModelModal isImportModalOpen={true} setIsImportModalOpen={vi.fn()} />);

    expect(screen.getByTestId('modal-Import Model')).toBeInTheDocument();
    expect(screen.getByTestId('rjsf-wrapper')).toBeInTheDocument();
  });

  it('does not render the main modal when closed', () => {
    render(<ImportModelModal isImportModalOpen={false} setIsImportModalOpen={vi.fn()} />);
    expect(screen.queryByTestId('modal-Import Model')).not.toBeInTheDocument();
  });

  it('calls importModelReq on URL submission with a non-empty url', async () => {
    render(<ImportModelModal isImportModalOpen={true} setIsImportModalOpen={vi.fn()} />);

    fireEvent.click(screen.getByTestId('submit-url'));
    await Promise.resolve();
    expect(importModelReq).toHaveBeenCalled();
  });

  it('does not call importModelReq for an empty URL submit', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    render(<ImportModelModal isImportModalOpen={true} setIsImportModalOpen={vi.fn()} />);
    fireEvent.click(screen.getByTestId('submit-empty-url'));
    expect(importModelReq).not.toHaveBeenCalled();
    consoleSpy.mockRestore();
  });

  it('does not call importModelReq for an invalid upload type', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    render(<ImportModelModal isImportModalOpen={true} setIsImportModalOpen={vi.fn()} />);
    fireEvent.click(screen.getByTestId('submit-invalid'));
    expect(importModelReq).not.toHaveBeenCalled();
    consoleSpy.mockRestore();
  });

  it('opens the CSV modal on CSV upload type submit', () => {
    const setIsImportModalOpen = vi.fn();
    render(
      <ImportModelModal isImportModalOpen={true} setIsImportModalOpen={setIsImportModalOpen} />,
    );

    fireEvent.click(screen.getByTestId('submit-csv'));
    expect(setIsImportModalOpen).toHaveBeenCalledWith(false);
    expect(screen.getByTestId('modal-Import CSV')).toBeInTheDocument();
  });
});
