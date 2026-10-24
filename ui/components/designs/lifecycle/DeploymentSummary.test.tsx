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
  SEVERITY: {
    INFO: 'informational',
    ERROR: 'error',
    WARNING: 'warning',
    SUCCESS: 'success',
  },
  SEVERITY_STYLE: {
    informational: { color: '#00f' },
    success: { color: '#0f0' },
    warning: { color: '#f80' },
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
import { SEVERITY } from '../../layout/NotificationCenter/constants';

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

  it('renders the success component details from the summary metadata', () => {
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

  it.each([
    [SEVERITY.WARNING, '1 of 3 components failed', 3, 1],
    [SEVERITY.ERROR, '2 of 2 components failed', 2, 2],
  ])('renders counts for %s', (severity, expectedText, total, failed) => {
    const summary = {
      c: [{ Location: 'loc', Summary: [{ CompName: 'db', Kind: 'Database', Success: false }] }],
    };
    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Outcome',
          severity,
          action: 'deploy',
          metadata: { total, failed, summary },
        }}
      />,
    );
    expect(screen.getByText(expectedText)).toBeInTheDocument();
    expect(screen.getByText('Failed to deploy Database "db"')).toBeInTheDocument();
  });

  it('renders old-server event without the new keys as today without crashing', () => {
    const summary = {
      c: [
        { Location: 'cluster-1', Summary: [{ CompName: 'web', Kind: 'Service', Success: true }] },
      ],
    };
    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Deployment completed',
          severity: SEVERITY.INFO,
          action: 'deploy',
          metadata: { summary },
        }}
      />,
    );
    expect(screen.getByText('Deployed Service "web"')).toBeInTheDocument();
    expect(screen.queryByText(/components failed/)).not.toBeInTheDocument();
  });
});
