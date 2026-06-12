import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

let stepperState: any = {
  activeStep: 0,
  steps: [],
  handleNext: vi.fn(),
  goBack: vi.fn(),
  canGoBack: false,
  activeStepComponent: <div data-testid="step-content" />,
};

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) =>
      typeof Component === 'string'
        ? React.createElement(Component, props, children)
        : React.createElement(Component, props, children);
    return StyledComponent;
  };

  return {
    styled,
    ModalFooter: ({ children, helpText }: any) => (
      <div data-testid="modal-footer">
        <div data-testid="help-text">{helpText}</div>
        {children}
      </div>
    ),
    useStepper: (config: any) => {
      stepperState.steps = config.steps;
      stepperState.activeStepComponent = config.steps[stepperState.activeStep]?.component;
      return stepperState;
    },
    CustomizedStepper: ({ children }: any) => (
      <div data-testid="customized-stepper">{children}</div>
    ),
    ModalBody: ({ children }: any) => <div data-testid="modal-body">{children}</div>,
    Box: ({ children }: any) => <div>{children}</div>,
    ModalButtonSecondary: ({ children, onClick, disabled }: any) => (
      <button data-testid="back" onClick={onClick} disabled={disabled}>
        {children}
      </button>
    ),
    ModalButtonPrimary: ({ children, onClick, disabled }: any) => (
      <button data-testid="next" onClick={onClick} disabled={disabled}>
        {children}
      </button>
    ),
    ComponentIcon: () => <svg data-testid="component-icon" />,
    Typography: ({ children }: any) => <span>{children}</span>,
    FormControl: ({ children }: any) => <div>{children}</div>,
    OutlinedInput: (props: any) => <input data-testid={props.id} {...props} />,
    DownloadIcon: () => <svg data-testid="download-icon" />,
    useTheme: () => ({
      palette: { icon: { default: '#333' }, border: { brand: '#0bf' }, text: { brand: '#0bf' } },
      shape: { borderRadius: 4 },
    }),
    Chip: ({ label, onDelete, icon }: any) => (
      <div data-testid="file-chip">
        {icon}
        {label}
        {onDelete && (
          <button data-testid="delete-chip" onClick={onDelete}>
            x
          </button>
        )}
      </div>
    ),
  };
});

vi.mock('@/assets/icons/ModelIcon', () => ({
  default: () => <svg data-testid="model-icon" />,
}));

vi.mock('@/assets/icons', () => ({
  LanOutlined: () => <svg data-testid="lan-icon" />,
  InsertDriveFile: () => <svg data-testid="file-drive-icon" />,
}));

vi.mock('@/utils/TooltipButton', () => ({
  TooltipIconButton: ({ children, onClick, title }: any) => (
    <button data-testid="tooltip-button" data-title={title} onClick={onClick}>
      {children}
    </button>
  ),
}));

vi.mock('@/constants/endpoints', () => ({
  MESHERY_DOCS_URL: 'https://docs.meshery.io',
}));

vi.mock('./style', () => ({
  StyledDocsRedirectLink: ({ children, href }: any) => <a href={href}>{children}</a>,
}));

vi.mock('@meshery/schemas', () => ({
  ComponentDefinitionV1Beta1OpenApiSchema: {
    components: {
      schemas: { ComponentDefinition: { description: 'Component def' } },
    },
  },
  ModelDefinitionV1Beta1OpenApiSchema: {
    components: {
      schemas: { ModelDefinition: { description: 'Model def' } },
    },
  },
}));

vi.mock('@/assets/icons/FinishFlagIcon', () => ({
  default: () => <svg data-testid="finish-flag" />,
}));

vi.mock('./FinishModelGenerateStep', () => ({
  default: () => <div data-testid="finish-step" />,
}));

import CsvStepper from './CSVStepper';

describe('CsvStepper', () => {
  beforeEach(() => {
    stepperState = {
      activeStep: 0,
      steps: [],
      handleNext: vi.fn(),
      goBack: vi.fn(),
      canGoBack: false,
      activeStepComponent: <div data-testid="step-content" />,
    };
  });

  it('renders the modal body with the stepper and footer', () => {
    render(<CsvStepper handleClose={vi.fn()} />);

    expect(screen.getByTestId('modal-body')).toBeInTheDocument();
    expect(screen.getByTestId('customized-stepper')).toBeInTheDocument();
    expect(screen.getByTestId('modal-footer')).toBeInTheDocument();
  });

  it('Back button is disabled on the first step', () => {
    render(<CsvStepper handleClose={vi.fn()} />);
    expect(screen.getByTestId('back')).toBeDisabled();
  });

  it('Next button is disabled when no model CSV is selected on step 0', () => {
    render(<CsvStepper handleClose={vi.fn()} />);
    expect(screen.getByTestId('next')).toBeDisabled();
    expect(screen.getByTestId('next')).toHaveTextContent('Next');
  });

  it('Finish button calls handleClose on the last step', () => {
    stepperState.activeStep = 3;
    const handleClose = vi.fn();
    render(<CsvStepper handleClose={handleClose} />);

    expect(screen.getByTestId('next')).toHaveTextContent('Finish');
    fireEvent.click(screen.getByTestId('next'));
    expect(handleClose).toHaveBeenCalled();
  });

  it('renders four steps total (Model, Component, Relationship CSVs, Finish)', () => {
    render(<CsvStepper handleClose={vi.fn()} />);
    expect(stepperState.steps).toHaveLength(4);
    expect(stepperState.steps[0].label).toBe('Model CSV');
    expect(stepperState.steps[1].label).toBe('Component CSV');
    expect(stepperState.steps[2].label).toBe('Relationship CSV');
    expect(stepperState.steps[3].label).toBe('Finish');
  });

  it('step 2 (Generate) has nextButtonText "Generate"', () => {
    stepperState.activeStep = 2;
    render(<CsvStepper handleClose={vi.fn()} />);
    expect(screen.getByTestId('next')).toHaveTextContent('Generate');
  });

  it('Next is enabled when on step 2 but disabled without relationship file', () => {
    stepperState.activeStep = 2;
    render(<CsvStepper handleClose={vi.fn()} />);
    expect(screen.getByTestId('next')).toBeDisabled();
  });
});
