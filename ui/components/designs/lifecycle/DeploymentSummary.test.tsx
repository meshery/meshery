import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (tag: any) => () => {
    const Styled = ({ children, ...rest }: any) =>
      typeof tag === 'string' ? (
        React.createElement(tag, rest, children)
      ) : (
        <div {...rest}>{children}</div>
      );
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
    Box: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
    Stack: ({ children }: any) => <div>{children}</div>,
    Typography: ({ children, 'data-testid': testId }: any) => (
      <span data-testid={testId}>{children}</span>
    ),
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

vi.mock('../../layout/NotificationCenter/formatters/pattern_dryrun', () => ({
  DryRunResponse: ({ response }: any) => (
    <div data-testid="dry-run-response">{JSON.stringify(response)}</div>
  ),
}));

const { openViewScopedToDesignInOperator, useIsOperatorEnabled, router } = vi.hoisted(() => ({
  openViewScopedToDesignInOperator: vi.fn(),
  useIsOperatorEnabled: vi.fn(() => false),
  router: { push: () => {} },
}));

vi.mock('@/utils/utils', () => ({
  openViewScopedToDesignInOperator,
  useIsOperatorEnabled,
}));

vi.mock('next/router', () => ({
  useRouter: () => router,
}));

import { DeploymentSummaryFormatter } from './DeploymentSummary';

const errorLines = () =>
  screen.getAllByTestId('design-event-error').map((line) => line.textContent);

describe('DeploymentSummaryFormatter', () => {
  afterEach(() => {
    useIsOperatorEnabled.mockReturnValue(false);
    openViewScopedToDesignInOperator.mockClear();
  });

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
    expect(errorLines()).toEqual([
      'meshery-server-1480 - component c1',
      'FieldValueInvalid - component c2',
    ]);
    // Only the MeshKit code links to the error code reference.
    expect(screen.getByText('meshery-server-1480').closest('a')).toHaveAttribute(
      'href',
      'https://docs.meshery.io/reference/error-codes#meshery-server-1480',
    );
    expect(screen.getByText(/FieldValueInvalid/).closest('a')).toBeNull();
    expect(screen.queryByText(/Showing the first/)).not.toBeInTheDocument();
  });

  it('renders a rewritten deploy error that carries a code but no componentId', () => {
    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Deployment Failed',
          severity: 'error',
          action: 'deploy',
          metadata: {
            design_name: 'Legacy Design',
            status: 'failed',
            errorCount: 1,
            errors: [{ code: 'meshery-server-1480' }],
            truncated: false,
          },
        }}
      />,
    );
    expect(errorLines()).toEqual(['meshery-server-1480']);
    expect(screen.queryByText(/undefined|unknown/i)).not.toBeInTheDocument();
  });

  it('renders the error summary of a dry-run event and its legacy dryRunResponse otherwise', () => {
    const { unmount } = render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Dry run complete',
          severity: 'informational',
          action: 'dry-run',
          metadata: {
            status: 'failed',
            errorCount: 1,
            errors: [{ code: 'FieldValueRequired', componentId: 'c1' }],
            truncated: false,
          },
        }}
      />,
    );
    expect(errorLines()).toEqual(['FieldValueRequired - component c1']);
    expect(screen.queryByTestId('dry-run-response')).not.toBeInTheDocument();
    unmount();

    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Dry run complete',
          severity: 'informational',
          action: 'dry-run',
          metadata: { summary: { dryRunResponse: { pod: { ctx: { success: true } } } } },
        }}
      />,
    );
    expect(screen.getByTestId('dry-run-response')).toHaveTextContent('pod');
  });

  it('opens the design in operator from the legacy design_name/design_id keys', () => {
    useIsOperatorEnabled.mockReturnValue(true);
    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Deployment Complete',
          severity: 'success',
          action: 'deploy',
          metadata: { design_name: 'Legacy Design', design_id: 'legacy-id' },
        }}
      />,
    );
    fireEvent.click(screen.getByText(/Open In Operator/));
    expect(openViewScopedToDesignInOperator).toHaveBeenLastCalledWith(
      'Legacy Design',
      'legacy-id',
      router,
    );
  });

  it('prefers the camelCase design keys over the legacy ones', () => {
    useIsOperatorEnabled.mockReturnValue(true);
    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Deployment Complete',
          severity: 'success',
          action: 'deploy',
          metadata: {
            designName: 'New Name',
            design_name: 'Old Name',
            designId: 'new-id',
            design_id: 'old-id',
          },
        }}
      />,
    );
    fireEvent.click(screen.getByText(/Open In Operator/));
    expect(openViewScopedToDesignInOperator).toHaveBeenLastCalledWith('New Name', 'new-id', router);
  });

  it('offers Open In Operator for a deploy event but not for a dry-run event', () => {
    useIsOperatorEnabled.mockReturnValue(true);
    const metadata = { designId: 'd1', designName: 'My Design' };
    const { unmount } = render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Dry run complete',
          severity: 'informational',
          action: 'dry-run',
          metadata,
        }}
      />,
    );
    expect(screen.queryByText(/Open In Operator/)).not.toBeInTheDocument();
    unmount();

    render(
      <DeploymentSummaryFormatter
        event={{
          description: 'Deployment Complete',
          severity: 'success',
          action: 'deploy',
          metadata,
        }}
      />,
    );
    expect(screen.getByText(/Open In Operator/)).toBeInTheDocument();
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
