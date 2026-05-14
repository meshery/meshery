import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const importModelReq = vi.fn();

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<any>();
  return {
    ...actual,
    FormControlLabel: ({ label, control }: any) => (
      <label>
        {control}
        {label}
      </label>
    ),
    Button: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
    FormControl: ({ children }: any) => <div>{children}</div>,
    RadioGroup: ({ children }: any) => <div>{children}</div>,
    Radio: () => <input type="radio" />,
    Typography: ({ children }: any) => <span>{children}</span>,
    ModalButtonPrimary: ({ children, onClick }: any) => (
      <button onClick={onClick}>{children}</button>
    ),
    Box: ({ children }: any) => <div>{children}</div>,
    EventBus: class {
      publish() {}
      on() {
        return { subscribe: () => ({ unsubscribe() {} }) };
      }
      onAny() {
        return { subscribe: () => ({ unsubscribe() {} }) };
      }
    },
    InfoIcon: () => <svg data-testid="info-icon" />,
    styled: (Component: any) => () => {
      const Styled = ({ children, ...props }: any) =>
        typeof Component === 'string' ? (
          React.createElement(Component, props, children)
        ) : (
          <Component {...props}>{children}</Component>
        );
      return Styled;
    },
    createTheme: (theme: any = {}) => ({
      ...theme,
      breakpoints: theme.breakpoints ?? {
        values: { xs: 0, sm: 600, md: 900, lg: 1200, xl: 1536 },
        up: () => '',
        down: () => '',
        between: () => '',
      },
    }),
    useTheme: () => ({
      palette: { background: { surfaces: '#fff' } },
    }),
  };
});

vi.mock('@meshery/schemas', () => ({
  ModelImportRjsfSchemaV1Beta2: {
    title: 'Import',
    properties: { uploadType: { enumNames: ['File', 'URL', 'CSV'] } },
    allOf: [],
  },
  ModelImportRjsfUiSchemaV1Beta2: {},
}));

vi.mock('@/components/shared/Modal', () => ({
  FormModal: ({ isOpen, title, onSubmit, helpText }: any) =>
    isOpen ? (
      <div data-testid={`modal-${title}`}>
        <div data-testid="rjsf-wrapper">
          {helpText}
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
      </div>
    ) : null,
  Modal: ({ children, isOpen, title, actions }: any) =>
    isOpen ? (
      <div data-testid={`modal-${title}`}>
        {children}
        {actions}
      </div>
    ) : null,
}));

vi.mock('./Stepper/CSVStepper', () => ({
  default: () => <div data-testid="csv-stepper" />,
}));

vi.mock('@/constants/endpoints', () => ({
  MESHERY_DOCS_URL: 'https://docs.meshery.io',
}));

vi.mock('@/utils/utils', () => ({
  getUnit8ArrayDecodedFile: (d: any) => (d ? [1, 2, 3] : null),
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

vi.mock('@/store/slices/mesheryUi', () => ({
  default: () => ({}),
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

  it('calls importModelReq on URL submission with a non-empty url', () => {
    render(<ImportModelModal isImportModalOpen={true} setIsImportModalOpen={vi.fn()} />);

    fireEvent.click(screen.getByTestId('submit-url'));
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
    // The CSV modal should now be open since handleClose closes the import modal
    // and sets the CSV one open.
    expect(screen.getByTestId('modal-Import CSV')).toBeInTheDocument();
  });
});
