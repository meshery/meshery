import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const useMeshModelComponentsMock = vi.fn();
const useStepperMock = vi.fn();
const downloadFileMock = vi.fn();

vi.mock('@meshery/schemas', () => ({
  RelationshipDefinitionV1Beta2OpenApiSchema: {
    components: {
      schemas: {
        RelationshipDefinition: {
          type: 'object',
          properties: {
            name: { type: 'string' },
            kind: { type: 'string' },
            capabilities: { type: 'array' },
            model: {
              description: 'Model Reference to the specific registered model...',
            },
            selectors: {
              items: {
                properties: {
                  allow: {
                    properties: { from: { items: { properties: { kind: {}, model: {} } } } },
                  },
                },
              },
            },
          },
        },
      },
    },
  },
  ModelDefinitionV1Beta1OpenApiSchema: {
    components: {
      schemas: {
        ModelDefinition: {
          properties: {
            category: {
              properties: {
                name: {
                  description: 'The category of the model that determines the main grouping.',
                },
              },
            },
          },
        },
      },
    },
  },
}));

vi.mock('@/theme', () => ({
  GlobalStyles: () => null,
}));

vi.mock('@sistent/sistent', () => {
  return {
    ModalFooter: ({ children }: any) => <div data-testid="modal-footer">{children}</div>,
    ModalBody: ({ children }: any) => <div data-testid="modal-body">{children}</div>,
    ModalButtonPrimary: ({ children, onClick, disabled, ...props }: any) => (
      <button
        data-testid={props['data-testid'] || 'primary-btn'}
        onClick={onClick}
        disabled={!!disabled}
      >
        {children}
      </button>
    ),
    ModalButtonSecondary: ({ children, onClick, disabled }: any) => (
      <button data-testid="secondary-btn" onClick={onClick} disabled={!!disabled}>
        {children}
      </button>
    ),
    CustomizedStepper: ({ children }: any) => <div data-testid="stepper">{children}</div>,
    Box: ({ children }: any) => <div>{children}</div>,
    TextField: ({ value, onChange, label, children, select, id }: any) => (
      <div>
        <label htmlFor={id}>{label}</label>
        {select ? (
          <select data-testid={`select-${id}`} value={value || ''} onChange={onChange}>
            {children}
          </select>
        ) : (
          <input data-testid={`input-${id}`} id={id} value={value || ''} onChange={onChange} />
        )}
      </div>
    ),
    MenuItem: ({ children, value, disabled }: any) => (
      <option value={value} disabled={!!disabled}>
        {children}
      </option>
    ),
    FormControl: ({ children }: any) => <div>{children}</div>,
    Grid2: ({ children }: any) => <div>{children}</div>,
    Typography: ({ children }: any) => <span>{children}</span>,
    useStepper: (...args: any[]) => useStepperMock(...args),
    styled:
      () =>
      () =>
      ({ children }: any) => <div>{children}</div>,
    DescriptionIcon: () => null,
    CodeIcon: () => null,
  };
});

vi.mock('@/assets/icons', () => ({
  Link: () => null,
  FilterAlt: () => null,
  SaveAlt: () => null,
}));

vi.mock('@/utils/fileDownloader', () => ({
  downloadFileFromContent: (...args: any[]) => downloadFileMock(...args),
}));

vi.mock('../meshery-mesh-interface/PatternService/RJSF_wrapper', () => ({
  default: () => <div data-testid="rjsf-wrapper" />,
}));

vi.mock('@/utils/hooks/useMeshModelComponents', () => ({
  useMeshModelComponents: () => useMeshModelComponentsMock(),
}));

vi.mock('./SelectorsForm', () => ({
  default: () => <div data-testid="selectors-form" />,
}));

import RelationshipFormStepper from './RelationshipFormStepper';

const makeStepperMock = (overrides: any = {}) => {
  const defaultMock = {
    activeStep: 0,
    activeStepComponent: <div data-testid="active-step" />,
    handleNext: vi.fn(),
    goBack: vi.fn(),
    canGoBack: false,
    steps: [
      { helpText: 'help 1' },
      { helpText: 'help 2' },
      { helpText: 'help 3' },
      { helpText: 'help 4' },
    ],
  };
  return { ...defaultMock, ...overrides };
};

describe('RelationshipFormStepper', () => {
  beforeEach(() => {
    downloadFileMock.mockReset();
    useMeshModelComponentsMock.mockReturnValue({
      models: { Kubernetes: [{ name: 'pod', displayName: 'Pod', id: '1' }] },
      meshmodelComponents: {},
      getModelFromCategory: vi.fn(),
      getComponentsFromModel: vi.fn(),
      categories: [{ name: 'Kubernetes' }],
    });
  });

  it('renders the active step content and a Next button', () => {
    useStepperMock.mockReturnValue(makeStepperMock());
    render(<RelationshipFormStepper handleClose={vi.fn()} />);

    expect(screen.getByTestId('active-step')).toBeInTheDocument();
    expect(screen.getByTestId('RelationshipStepper-Button-Next')).toBeInTheDocument();
    expect(screen.getByTestId('secondary-btn')).toBeInTheDocument();
  });

  it('disables Next on the first step when name/category/model are not set', () => {
    useStepperMock.mockReturnValue(makeStepperMock());
    render(<RelationshipFormStepper handleClose={vi.fn()} />);

    expect(screen.getByTestId('RelationshipStepper-Button-Next')).toBeDisabled();
  });

  it('enables Next on the relationship properties step', () => {
    useStepperMock.mockReturnValue(makeStepperMock({ activeStep: 1 }));
    render(<RelationshipFormStepper handleClose={vi.fn()} />);

    expect(screen.getByTestId('RelationshipStepper-Button-Next')).not.toBeDisabled();
  });

  it('shows Finish on the last step and invokes handleClose', () => {
    const handleClose = vi.fn();
    const stepper = makeStepperMock({ activeStep: 3 });
    useStepperMock.mockReturnValue(stepper);

    render(<RelationshipFormStepper handleClose={handleClose} />);

    const finish = screen.getByTestId('RelationshipStepper-Button-Finish');
    fireEvent.click(finish);
    expect(handleClose).toHaveBeenCalled();
  });

  it('invokes handleNext when Next is clicked on a valid step', () => {
    const stepper = makeStepperMock({ activeStep: 1 });
    useStepperMock.mockReturnValue(stepper);

    render(<RelationshipFormStepper handleClose={vi.fn()} />);

    fireEvent.click(screen.getByTestId('RelationshipStepper-Button-Next'));
    expect(stepper.handleNext).toHaveBeenCalled();
  });

  it('invokes goBack when Back is clicked and canGoBack is true', () => {
    const stepper = makeStepperMock({ canGoBack: true });
    useStepperMock.mockReturnValue(stepper);

    render(<RelationshipFormStepper handleClose={vi.fn()} />);

    fireEvent.click(screen.getByTestId('secondary-btn'));
    expect(stepper.goBack).toHaveBeenCalled();
  });
});
