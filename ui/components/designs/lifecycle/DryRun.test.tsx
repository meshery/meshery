import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (_tag: any) => () => {
    const Styled = ({ children, ...rest }: any) => <div {...rest}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
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
      },
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

import { FormatDryRunResponse, getTotalCountOfDeploymentErrors } from './DryRun';

describe('getTotalCountOfDeploymentErrors', () => {
  it('returns 0 when there are no errors', () => {
    expect(getTotalCountOfDeploymentErrors(undefined as any)).toBe(0);
    expect(getTotalCountOfDeploymentErrors([])).toBe(0);
  });

  it('sums up errors across entries', () => {
    const result = getTotalCountOfDeploymentErrors([
      { errors: [1, 2] },
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
});
