import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

let stepperState: any;

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
    TextField: ({ label, value, onChange, id }: any) => (
      <label>
        {label}
        <input
          data-testid={`textfield-${id || label}`}
          value={value || ''}
          onChange={(e) => onChange?.(e)}
        />
      </label>
    ),
    ModalButtonSecondary: ({ children, onClick, disabled }: any) => (
      <button data-testid="back" onClick={onClick} disabled={disabled}>
        {children}
      </button>
    ),
    ModalButtonPrimary: ({ children, onClick, disabled, 'data-testid': testId }: any) => (
      <button data-testid={testId || 'next'} onClick={onClick} disabled={disabled}>
        {children}
      </button>
    ),
    Select: ({ children, value, onChange }: any) => (
      <select data-testid="select" value={value || ''} onChange={onChange}>
        {children}
      </select>
    ),
    InputLabel: ({ children }: any) => <label>{children}</label>,
    FormControlLabel: ({ label, control }: any) => (
      <label>
        {control}
        {label}
      </label>
    ),
    Checkbox: ({ checked, onChange }: any) => (
      <input type="checkbox" checked={!!checked} onChange={onChange} />
    ),
    Typography: ({ children }: any) => <span>{children}</span>,
    FormControl: ({ children }: any) => <div>{children}</div>,
    RadioGroup: ({ children }: any) => <div>{children}</div>,
    MenuItem: ({ children, value }: any) => <option value={value}>{children}</option>,
    Radio: () => <input type="radio" />,
    Grid2: ({ children }: any) => <div>{children}</div>,
    AppRegistrationIcon: () => <svg />,
    BrushIcon: () => <svg />,
    CategoryIcon: () => <svg />,
    DescriptionIcon: () => <svg />,
  };
});

vi.mock('./style', () => ({
  StyledSummaryBox: ({ children }: any) => <div>{children}</div>,
  StyledSummaryItem: ({ children }: any) => <div>{children}</div>,
  SectionHeading: ({ children }: any) => <h2>{children}</h2>,
  StyledColorBox: () => <div />,
  StyledDocsRedirectLink: ({ children, href }: any) => <a href={href}>{children}</a>,
}));

vi.mock('@/assets/icons/SourceIcon', () => ({ default: () => <svg /> }));
vi.mock('@/assets/icons/FinishFlagIcon', () => ({ default: () => <svg /> }));
vi.mock('@/assets/icons/DeploymentSelectorIcon', () => ({
  DeploymentSelectorIcon: () => <svg />,
}));

vi.mock('@meshery/schemas', () => {
  const stringField = {
    pattern: '^[a-zA-Z0-9-_ ]+$',
    helperText: 'a hint',
    description: 'desc',
    examples: ['example'],
    default: '',
  };
  return {
    CategoryDefinitionV1Beta1OpenApiSchema: {
      components: {
        schemas: {
          CategoryDefinition: {
            properties: {
              name: {
                default: 'Uncategorized',
                enum: ['Uncategorized', 'A', 'B'],
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
              metadata: {
                properties: {
                  shape: {
                    default: 'circle',
                    description: 'shape',
                    enum: ['circle', 'square'],
                  },
                  primaryColor: {
                    default: '#000',
                    description: 'primary color',
                  },
                  secondaryColor: {
                    default: '#fff',
                    description: 'secondary color',
                  },
                },
              },
              name: stringField,
              displayName: stringField,
              version: stringField,
              category: { ...stringField, description: 'category desc' },
              subCategory: { ...stringField, description: 'subcategory desc' },
            },
          },
        },
      },
    },
    SubCategoryDefinitionV1Beta1OpenApiSchema: {
      components: {
        schemas: {
          SubCategoryDefinition: {
            default: 'Uncategorized',
            enum: ['Uncategorized', 'X', 'Y'],
          },
        },
      },
    },
  };
});

vi.mock('lodash', () => ({
  capitalize: (s: string) => s,
}));

vi.mock('./FinishModelGenerateStep', () => ({
  default: () => <div data-testid="finish-step" />,
}));

import UrlStepper from './UrlStepper';

describe('UrlStepper', () => {
  beforeEach(() => {
    stepperState = {
      activeStep: 0,
      steps: [],
      handleNext: vi.fn(),
      goBack: vi.fn(),
      canGoBack: false,
      activeStepComponent: null,
    };
  });

  it('renders modal body and footer', () => {
    render(<UrlStepper handleClose={vi.fn()} />);
    expect(screen.getByTestId('modal-body')).toBeInTheDocument();
    expect(screen.getByTestId('modal-footer')).toBeInTheDocument();
  });

  it('disables Next on the first step when model name/displayName are empty', () => {
    render(<UrlStepper handleClose={vi.fn()} />);
    expect(screen.getByTestId('UrlStepper-Button-Next')).toBeDisabled();
  });

  it('disables Back on the first step', () => {
    render(<UrlStepper handleClose={vi.fn()} />);
    expect(screen.getByTestId('back')).toBeDisabled();
  });

  it('Generate button on step 5 calls handleNext on click', () => {
    stepperState.activeStep = 5;
    render(<UrlStepper handleClose={vi.fn()} />);
    const btn = screen.getByTestId('UrlStepper-Button-Generate');
    expect(btn).toBeInTheDocument();
    fireEvent.click(btn);
    expect(stepperState.handleNext).toHaveBeenCalled();
  });

  it('Finish button on step 6 calls handleClose', () => {
    stepperState.activeStep = 6;
    const handleClose = vi.fn();
    render(<UrlStepper handleClose={handleClose} />);
    fireEvent.click(screen.getByTestId('UrlStepper-Button-Finish'));
    expect(handleClose).toHaveBeenCalled();
  });
});
