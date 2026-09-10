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
    TextField: ({ label, value, onChange, id, helperText, disabled }: any) => (
      <label>
        {label}
        <input
          data-testid={`textfield-${id || label}`}
          value={value || ''}
          disabled={disabled}
          onChange={(e) => onChange?.(e)}
        />
        {helperText ? <span data-testid={`helper-${id || label}`}>{helperText}</span> : null}
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
    // A RadioGroup option carries its value here; render it as a real radio so
    // selecting a source drives the component the way the browser does. The
    // checkbox form (control, no value) keeps its original shape.
    FormControlLabel: ({ label, control, value, onChange, ...rest }: any) =>
      value !== undefined ? (
        <label>
          <input
            type="radio"
            data-testid={rest['data-testid']}
            value={value}
            onChange={(e) => onChange?.(e)}
          />
          {label}
        </label>
      ) : (
        <label>
          {control}
          <span>{label}</span>
        </label>
      ),
    Checkbox: ({ checked, onChange }: any) => (
      <input type="checkbox" checked={!!checked} onChange={onChange} />
    ),
    Typography: ({ children }: any) => <span>{children}</span>,
    FormControl: ({ children }: any) => <div>{children}</div>,
    RadioGroup: ({ children, onChange }: any) => (
      <div data-testid="radio-group">
        {React.Children.map(children, (child: any) =>
          React.isValidElement(child) ? React.cloneElement(child, { onChange } as any) : child,
        )}
      </div>
    ),
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

  describe('Source step URL validation', () => {
    const INVALID = 'this-is-not-a-url-at-all';
    const VALID_AH = 'https://artifacthub.io/packages/search?ts_query_web=meshery-operator';
    const VALID_GH = 'git://github.com/cert-manager/cert-manager/master/deploy/crds';
    const SOURCE_STEP = 3;

    const renderSourceStep = () => {
      stepperState.activeStep = SOURCE_STEP;
      render(<UrlStepper handleClose={vi.fn()} />);
    };
    // Click, not fireEvent.change: the mocked radio already holds this exact
    // value, so React's controlled-input de-duplication means setting it to the
    // same string dispatches nothing and the source is never selected. A click
    // is what React turns into a radio's change event anyway.
    const selectSource = (label: string) =>
      fireEvent.click(screen.getByTestId(`UrlStepper-Select-Source-${label}`));
    const typeUrl = (value: string) =>
      fireEvent.change(screen.getByTestId('textfield-model-url'), { target: { value } });
    const helperText = () => screen.queryByTestId('helper-model-url')?.textContent ?? '';

    // The regression. The radio stores "artifact hub" (lower-cased label, space
    // included) while validateUrl used to compare against "artifacthub", so no
    // pattern was chosen and `new RegExp(undefined)` matched everything.
    it('rejects an invalid Artifact Hub URL', () => {
      renderSourceStep();
      selectSource('Artifact Hub');
      typeUrl(INVALID);

      expect(helperText()).toContain('Invalid ArtifactHub URL');
    });

    it('accepts a well-formed Artifact Hub URL', () => {
      renderSourceStep();
      selectSource('Artifact Hub');
      typeUrl(VALID_AH);

      expect(helperText()).toBe('');
    });

    it('still rejects an invalid GitHub URL', () => {
      renderSourceStep();
      selectSource('GitHub');
      typeUrl(INVALID);

      expect(helperText()).toContain('Invalid GitHub URL');
    });

    // Guards the gap that would reopen if a source were added without a matching
    // pattern and message: Next is gated on `!urlError`, so an empty error for an
    // invalid URL would let it through. Driven off the rendered options so a
    // third source is covered automatically.
    it('reports an error for every selectable source when the URL is invalid', () => {
      renderSourceStep();
      const radios = screen.getAllByRole('radio');
      expect(radios.length).toBe(2);

      radios.forEach((radio) => {
        fireEvent.click(radio);
        typeUrl(INVALID);
        expect(helperText()).not.toBe('');
      });
    });

    // A URL is only valid relative to a source. Without revalidation on switch,
    // the GitHub URL below keeps its cleared error and Next stays enabled, so
    // the wizard would submit a GitHub URL with registrant "artifact hub".
    it('revalidates the entered URL when the source changes', () => {
      renderSourceStep();
      selectSource('GitHub');
      typeUrl(VALID_GH);
      expect(helperText()).toBe('');

      selectSource('Artifact Hub');

      expect(helperText()).toContain('Invalid ArtifactHub URL');
    });

    it('clears the error when switching to a source the URL is valid for', () => {
      renderSourceStep();
      selectSource('Artifact Hub');
      typeUrl(VALID_GH);
      expect(helperText()).toContain('Invalid ArtifactHub URL');

      selectSource('GitHub');

      expect(helperText()).toBe('');
    });

    it('accepts a well-formed GitHub URL', () => {
      renderSourceStep();
      selectSource('GitHub');
      typeUrl(VALID_GH);

      expect(helperText()).toBe('');
    });
  });
});
