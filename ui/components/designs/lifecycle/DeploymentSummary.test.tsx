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
    Box: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
    Stack: ({ children }: any) => <div>{children}</div>,
    Typography: ({ children }: any) => <span>{children}</span>,
    useTheme: () => ({
      palette: {
        text: { default: '#000' },
        background: { constant: { white: '#fff' } },
        success: { main: '#0f0' },
        error: { main: '#f00' },
        warning: { main: '#f80' },
      },
    }),
    Button: ({ children, onClick }: any) => (
      <button type="button" onClick={onClick}>
        {children}
      </button>
    ),
    ExternalLinkIcon: () => <svg data-testid="external-link" />,
    componentIcon: ({ kind }: any) => `/icons/${kind}.svg`,
  };
});

vi.mock('../../data-formatter', () => ({
  FormatStructuredData: ({ data }: any) => (
    <pre data-testid="structured">{JSON.stringify(data)}</pre>
  ),
  TextWithLinks: ({ text }: any) => <span data-testid="text-with-links">{text}</span>,
}));

vi.mock('../../layout/NotificationCenter/constants', () => ({
  SEVERITY_STYLE: {
    success: { color: '#0f0' },
    error: { color: '#f00' },
  },
}));

vi.mock('./common', () => ({
  ComponentIcon: ({ iconSrc, alt }: any) => (
    <img data-testid="component-icon" src={iconSrc} alt={alt} />
  ),
}));

vi.mock('../../layout/NotificationCenter/formatters/error', () => ({
  ErrorMetadataFormatter: ({ metadata }: any) => (
    <div data-testid="error-metadata">{JSON.stringify(metadata)}</div>
  ),
}));

vi.mock('@/utils/utils', () => ({
  openViewScopedToDesignInOperator: vi.fn(),
  useIsOperatorEnabled: () => false,
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

import { DeploymentSummaryFormatter } from './DeploymentSummary';

describe('DeploymentSummaryFormatter', () => {
  it('renders the event description and skips Open In Operator when operator is disabled', () => {
    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Deployment Complete',
          severity: 'success',
          action: 'deploy',
          metadata: {},
        }}
      />,
    );

    expect(screen.getByTestId('text-with-links')).toHaveTextContent('Deployment Complete');
    expect(screen.queryByText('Open In Operator')).not.toBeInTheDocument();
  });

  it('renders the success component details from legacy summary metadata', () => {
    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Deployment Complete',
          severity: 'success',
          action: 'deploy',
          metadata: {
            summary: {
              ctx1: [
                {
                  Location: 'us-east',
                  Summary: [{ CompName: 'svc', Kind: 'Service', Model: 'core', Success: true }],
                },
              ],
            },
          },
        }}
      />,
    );
    expect(screen.getByText('Deployed Service "svc"')).toBeInTheDocument();
  });

  it('renders the failure variant for failed components', () => {
    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Deployment Failed',
          severity: 'error',
          action: 'deploy',
          metadata: {
            summary: {
              ctx1: [
                {
                  Location: 'eu-west',
                  Summary: [
                    {
                      CompName: 'pod',
                      Kind: 'Pod',
                      Model: 'core',
                      Success: false,
                    },
                  ],
                },
              ],
            },
          },
        }}
      />,
    );
    expect(screen.getByText('Failed to deploy Pod "pod"')).toBeInTheDocument();
  });

  it('renders status, errorCount and the per-component errors from the slim metadata', () => {
    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Deployment Failed',
          severity: 'error',
          action: 'deploy',
          metadata: {
            designId: 'd1',
            designName: 'My Design',
            status: 'partial',
            errorCount: 2,
            errors: [
              { code: 'meshery-server-1480', componentId: 'c1' },
              { code: 'FieldValueInvalid', componentId: 'c2' },
            ],
            truncated: false,
          },
        }}
      />,
    );
    expect(screen.getByText(/Status:/)).toHaveTextContent('Status: partial - 2 errors');
    expect(screen.getByText('meshery-server-1480 - component c1')).toBeInTheDocument();
    expect(screen.getByText('FieldValueInvalid - component c2')).toBeInTheDocument();
    expect(screen.queryByText(/Showing the first/)).not.toBeInTheDocument();
  });

  it('notes when the errors list is truncated', () => {
    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Deployment Failed',
          severity: 'error',
          action: 'deploy',
          metadata: {
            status: 'failed',
            errorCount: 25,
            errors: Array.from({ length: 20 }, (_, i) => ({ code: 'E', componentId: `c${i}` })),
            truncated: true,
          },
        }}
      />,
    );
    expect(screen.getByText('Showing the first 20 of 25 errors.')).toBeInTheDocument();
  });

  it('renders ErrorMetadataFormatter when errors are present at the top level', () => {
    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'oh no',
          severity: 'error',
          action: 'deploy',
          metadata: {
            error: { code: 'E_FAIL' },
          },
        }}
      />,
    );
    expect(screen.getByTestId('error-metadata')).toHaveTextContent('E_FAIL');
  });
});
