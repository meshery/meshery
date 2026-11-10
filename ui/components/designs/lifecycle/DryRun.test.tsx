import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

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
    Collapse: ({ children }: any) => <div>{children}</div>,
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
    tapOnError: vi.fn(),
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
  DryRunErrorContainer: ({ children, onClick }: any) => <div onClick={onClick}>{children}</div>,
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
