import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { createActor, createMachine } from 'xstate';

vi.mock('@sistent/sistent', () => {
  const styled = (_tag: any, _options?: any) => () => {
    const Styled = ({ children, ...rest }: any) => <div {...rest}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
    Box: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
    Button: ({ children, onClick, ...rest }: any) => (
      <button type="button" onClick={onClick} {...rest}>
        {children}
      </button>
    ),
    ButtonGroup: ({ children, ...rest }: any) => <div>{children}</div>,
    Card: ({ children, ...rest }: any) => <div>{children}</div>,
    Chip: ({ label, ...rest }: any) => <span {...rest}>{label}</span>,
    List: ({ children, ...rest }: any) => <ul {...rest}>{children}</ul>,
    ListItemText: ({ primary }: any) => <span>{primary}</span>,
    ListItemIcon: ({ children }: any) => <span>{children}</span>,
    Typography: ({ children, 'data-testid': testId, ...rest }: any) => (
      <span data-testid={testId} {...rest}>
        {children}
      </span>
    ),
    Collapse: ({ children, in: isOpen }: any) => (isOpen ? <div>{children}</div> : null),
    useTheme: () => ({
      palette: {
        text: { default: '#000' },
        background: { error: { default: '#f00' } },
        success: { main: '#0f0' },
        error: { main: '#f00' },
        warning: { main: '#f80' },
        primary: { main: '#00B39F' },
      },
      spacing: (n: number) => `${n * 8}px`,
    }),
    ErrorIcon: () => <svg data-testid="error-icon" />,
  };
});

vi.mock('@/assets/icons', () => ({
  ExpandLess: () => <svg data-testid="expand-less" />,
  ExpandMore: () => <svg data-testid="expand-more" />,
}));

vi.mock('./common', () => ({
  ComponentIcon: ({ iconSrc }: any) => <span data-testid="component-icon" data-src={iconSrc} />,
  DEPLOYMENT_TYPE: { DEPLOY: 'deploy', UNDEPLOY: 'undeploy' },
  Loading: ({ message }: any) => <div data-testid="loading">{message}</div>,
  processDesign: vi.fn(() => ({ configurableComponents: [], annotationComponents: [] })),
}));

vi.mock('../../../machines/validator/designValidator', () => ({
  designValidatorCommands: {
    dryRunDesignDeployment: vi.fn(),
    dryRunDesignUnDeployment: vi.fn(),
  },
  designValidatorEvents: {
    tapOnError: vi.fn((data: any) => ({ type: 'TAP_ON_ERROR', data })),
  },
  useDryRunValidationResults: vi.fn(),
  useIsValidatingDryRun: vi.fn(),
  useDryRunRawResponse: vi.fn(),
}));

vi.mock('react-redux', () => ({
  useSelector: vi.fn(() => []),
}));

vi.mock('./CrossClusterCompatibilityAnalysis', () => ({
  CrossClusterCompatibilityAnalysis: () => <div data-testid="compatibility-matrix" />,
}));

vi.mock('../../data-formatter', () => ({
  FormatStructuredData: ({ data }: any) => (
    <pre data-testid="structured">{JSON.stringify(data)}</pre>
  ),
}));

vi.mock('./styles', () => ({
  DryRunComponentLabel: ({ children, onClick }: any) => (
    <button type="button" onClick={onClick}>
      {children}
    </button>
  ),
  DryRunComponentStyled: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
  DryRunErrorContainer: ({ children, clickable, disablePadding: _dp, ...rest }: any) => (
    <div data-clickable={String(clickable)} {...rest}>
      {children}
    </div>
  ),
  DryRunSignleError: ({ primary }: any) => <div>{primary}</div>,
  DryRunRootListStyled: ({ children, subheader }: any) => (
    <div>
      {subheader}
      {children}
    </div>
  ),
  ValidationSubHeader: ({ children }: any) => <div>{children}</div>,
}));

import { FormatDryRunResponse, DryRunDesign, getTotalCountOfDeploymentErrors } from './DryRun';
import {
  designValidatorCommands,
  useDryRunValidationResults,
  useIsValidatingDryRun,
  useDryRunRawResponse,
} from '../../../machines/validator/designValidator';

describe('getTotalCountOfDeploymentErrors', () => {
  it('returns 0 when there are no errors', () => {
    expect(getTotalCountOfDeploymentErrors(undefined as any)).toBe(0);
    expect(getTotalCountOfDeploymentErrors([])).toBe(0);
  });

  it('sums up errors across entries and skips entries without errors array', () => {
    const result = getTotalCountOfDeploymentErrors([
      { errors: [1, 2] },
      { type: 'CompatibilityMatrix', rawDryRunResponse: {} },
      { errors: [3] },
      { errors: [] },
    ] as any);
    expect(result).toBe(3);
  });
});

describe('FormatDryRunResponse', () => {
  it('renders the success state when there are no errors', () => {
    render(
      <FormatDryRunResponse
        dryRunErrors={[]}
        configurableComponentsCount={2}
        annotationComponentsCount={3}
        validationMachine={{}}
        currentComponentName={null}
      />,
    );
    expect(screen.getByTestId('dry-run-summary-success')).toBeInTheDocument();
    expect(screen.getByTestId('dry-run-summary-errors')).toHaveTextContent('0 error(s)');
  });

  it('renders FormatStructuredData entries for RequestError type errors', () => {
    render(
      <FormatDryRunResponse
        dryRunErrors={[
          {
            type: 'RequestError',
            errors: [{ data: { foo: 'bar' } }],
          },
        ]}
        configurableComponentsCount={2}
        annotationComponentsCount={3}
        validationMachine={{}}
        currentComponentName={null}
      />,
    );
    expect(screen.getByTestId('structured')).toHaveTextContent('"foo":"bar"');
  });

  it('shows the component count + annotation count when supplied', () => {
    render(
      <FormatDryRunResponse
        dryRunErrors={[]}
        configurableComponentsCount={5}
        annotationComponentsCount={4}
        validationMachine={{}}
        currentComponentName={null}
      />,
    );
    expect(screen.getByText('5 components and 4 annotations')).toBeInTheDocument();
  });

  it('ignores CompatibilityMatrix entries and renders success state if no real errors', () => {
    render(
      <FormatDryRunResponse
        dryRunErrors={[
          {
            type: 'CompatibilityMatrix',
            rawDryRunResponse: { frontend: { 'c-1': { success: true } } },
          },
        ]}
        configurableComponentsCount={2}
        annotationComponentsCount={1}
        validationMachine={{}}
        currentComponentName={null}
      />,
    );
    expect(screen.getByTestId('dry-run-summary-success')).toBeInTheDocument();
    expect(screen.getByTestId('dry-run-summary-errors')).toHaveTextContent('0 error(s)');
  });
});

describe('DryRunDesign view fallback', () => {
  const baseProps = {
    handleClose: vi.fn(),
    currentComponentName: null,
    design: { name: 'test-design' },
    validationMachine: { send: vi.fn() },
    selectedK8sContexts: [],
    deployment_type: 'deploy',
    includeDependencies: false,
  };

  it('falls back to the list view and hides the view toggle when there is no raw dry-run response', () => {
    vi.mocked(useDryRunValidationResults).mockReturnValue([]);
    vi.mocked(useIsValidatingDryRun).mockReturnValue(false);
    vi.mocked(useDryRunRawResponse).mockReturnValue(null);

    render(<DryRunDesign {...baseProps} />);

    expect(screen.queryByTestId('toggle-view-matrix')).not.toBeInTheDocument();
    expect(screen.queryByTestId('compatibility-matrix')).not.toBeInTheDocument();
    expect(screen.getByTestId('dry-run-summary-errors')).toHaveTextContent('0 error(s)');
  });

  it('renders the matrix view and the view toggle when a raw dry-run response exists', () => {
    vi.mocked(useDryRunValidationResults).mockReturnValue([]);
    vi.mocked(useIsValidatingDryRun).mockReturnValue(false);
    vi.mocked(useDryRunRawResponse).mockReturnValue({ frontend: {} });

    render(<DryRunDesign {...baseProps} />);

    expect(screen.getByTestId('toggle-view-matrix')).toBeInTheDocument();
    expect(screen.getByTestId('compatibility-matrix')).toBeInTheDocument();
  });
});

describe('dry-run error deep links', () => {
  const deployment = { id: 'comp-1', displayName: 'web', styles: {} };
  const replicasError = {
    Type: 'FieldValueInvalid',
    FieldPath: 'Deployment.spec.replicas',
    Message: 'bad',
  };
  const componentErrors = [
    {
      type: 'ComponentError',
      compName: 'web',
      component: deployment,
      contextId: 'ctx',
      errors: [replicasError],
    },
  ];
  const tapsSent = (send: ReturnType<typeof vi.fn>) =>
    send.mock.calls.map(([event]) => event).filter((event) => event?.type === 'TAP_ON_ERROR');

  // jsdom implements no layout, so it has no scrollIntoView.
  const scrollIntoView = vi.fn();

  beforeEach(() => {
    scrollIntoView.mockClear();
    Element.prototype.scrollIntoView = scrollIntoView;
    vi.mocked(useDryRunValidationResults).mockReturnValue(componentErrors);
    vi.mocked(useIsValidatingDryRun).mockReturnValue(false);
    vi.mocked(useDryRunRawResponse).mockReturnValue(null);
  });

  it('expands the component Kanvas opened the dry run for and links each error to it', () => {
    const send = vi.fn();
    render(
      <DryRunDesign
        handleClose={vi.fn()}
        currentComponentId="comp-1"
        design={{ name: 'd' }}
        validationMachine={{ send }}
        selectedK8sContexts={[]}
        deployment_type="deploy"
        includeDependencies={false}
      />,
    );

    // Expanded and scrolled to: Kanvas identifies it by id, not by name.
    expect(scrollIntoView).toHaveBeenCalledTimes(1);
    const error = screen.getByTestId('dry-run-error');
    expect(error).toHaveAttribute('role', 'button');
    expect(error).toHaveAttribute('data-clickable', 'true');

    fireEvent.click(error);
    expect(tapsSent(send)).toEqual([
      {
        type: 'TAP_ON_ERROR',
        data: { error: replicasError, type: 'dryRun', component: deployment },
      },
    ]);
  });

  it('lists errors as plain text where there is no designer to navigate into', () => {
    const send = vi.fn();
    render(
      <DryRunDesign
        handleClose={vi.fn()}
        currentComponentId="comp-1"
        design={{ name: 'd' }}
        validationMachine={{ send }}
        canNavigateToErrors={false}
        selectedK8sContexts={[]}
        deployment_type="deploy"
        includeDependencies={false}
      />,
    );

    const error = screen.getByTestId('dry-run-error');
    expect(error).not.toHaveAttribute('role');
    expect(error).toHaveAttribute('data-clickable', 'false');
    fireEvent.click(error);
    expect(tapsSent(send)).toEqual([]);
  });

  it('defaults FormatDryRunResponse to plain text, as the notification center renders stored events', () => {
    render(<FormatDryRunResponse dryRunErrors={componentErrors} currentComponentName="web" />);

    const error = screen.getByTestId('dry-run-error');
    expect(error).not.toHaveAttribute('role');
    expect(error).toHaveAttribute('data-clickable', 'false');
  });
});

describe('DryRunDesign mount-time request', () => {
  const props = {
    handleClose: vi.fn(),
    design: { name: 'test-design' },
    selectedK8sContexts: ['ctx-1'],
    includeDependencies: true,
  };

  beforeEach(() => {
    vi.mocked(useDryRunValidationResults).mockReturnValue([]);
    vi.mocked(useIsValidatingDryRun).mockReturnValue(false);
    vi.mocked(useDryRunRawResponse).mockReturnValue(null);
  });

  it('puts the returnAddress prop on its dry-run request', () => {
    const returnAddress = createActor(createMachine({ id: 'designValidationManager' }));
    const send = vi.fn();
    render(
      <DryRunDesign
        {...props}
        deployment_type="deploy"
        validationMachine={{ send }}
        returnAddress={returnAddress}
      />,
    );
    expect(designValidatorCommands.dryRunDesignDeployment).toHaveBeenCalledWith({
      design: props.design,
      k8sContexts: ['ctx-1'],
      includeDependencies: true,
      returnAddress,
    });
    expect(send).toHaveBeenCalledTimes(1);
  });

  it('still requests an undeploy dry run when no returnAddress is given', () => {
    const send = vi.fn();
    render(<DryRunDesign {...props} deployment_type="undeploy" validationMachine={{ send }} />);
    expect(designValidatorCommands.dryRunDesignUnDeployment).toHaveBeenCalledWith({
      design: props.design,
      k8sContexts: ['ctx-1'],
      includeDependencies: true,
      returnAddress: undefined,
    });
    expect(send).toHaveBeenCalledTimes(1);
  });
});
