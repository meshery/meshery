import React from 'react';
import { render, screen, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const passthrough =
    (Tag: any) =>
    ({ children, sx: _sx, size: _size, ...props }: any) => <Tag {...props}>{children}</Tag>;
  return {
    Typography: passthrough('span'),
    Grid2: passthrough('div'),
    Box: passthrough('div'),
    List: passthrough('ul'),
    ListItem: passthrough('li'),
    Chip: ({ label, title }: any) => (
      <span data-testid="chip" title={title}>
        {label}
      </span>
    ),
    RenderMarkdown: ({ content }: any) => <span data-testid="markdown">{content}</span>,
    // styled(Component, options)(styles) returns a component that drops the
    // styling props.
    styled:
      (Component: any) =>
      () =>
      ({ matched: _matched, ...props }: any) => <Component {...props} />,
  };
});

vi.mock('../../../data-formatter', () => ({
  FormatStructuredData: ({ data, propertyFormatters }: any) => {
    if (typeof data === 'string') return <div data-testid="fsd-string">{data}</div>;
    return (
      <div data-testid="format-structured">
        {Object.entries(data || {}).map(([key, value]) => {
          const Formatter = propertyFormatters?.[key];
          return (
            <div key={key} data-key={key}>
              {Formatter ? Formatter(value) : <span>{String(value)}</span>}
            </div>
          );
        })}
      </div>
    );
  },
}));

vi.mock('./common', () => ({
  TitleLink: ({ href, children }: any) => (
    <a href={href} target="_blank" rel="noopener noreferrer">
      {children}
    </a>
  ),
}));

vi.mock('../notification', () => ({
  canTruncateDescription: (description: string) => description.length > 62,
}));

import { ErrorMetadataFormatter } from './error';

const reconnect = {
  id: 'reconnect',
  description: 'Reconnect the connection once the cluster is reachable.',
  action: 'connection.reconnect',
  actionPolicy: 'requiresApproval',
};
const eventError = {
  code: 'meshery-server-1304',
  longDescription: 'connection refused',
  probableCauses: [
    {
      id: 'cluster-unreachable',
      description: 'The API server is not reachable.',
      remediations: [{ id: 'check-network', description: 'Check the network.' }, reconnect],
    },
    {
      id: 'credentials-expired',
      description: 'The credentials have expired.',
      remediations: [
        {
          id: 'refresh-credentials',
          description: 'Refresh the credentials.',
          action: 'connection.refresh-credentials',
          actionPolicy: 'automatic',
        },
        reconnect,
      ],
    },
  ],
  matchedCauses: [{ causeId: 'credentials-expired', confidence: 0.75 }],
  referenceUrl: 'https://docs.meshery.io/reference/error-codes/component-meshery-server/x/',
};

describe('ErrorMetadataFormatter', () => {
  it('renders each cause with its own remediations, the matched cause first', () => {
    render(<ErrorMetadataFormatter metadata={eventError} event={{ description: 'short' }} />);

    const causes = screen.getAllByTestId('error-cause');
    expect(causes.map((cause) => cause.getAttribute('data-cause-id'))).toEqual([
      'credentials-expired',
      'cluster-unreachable',
    ]);
    expect(within(causes[0]).getByText('Likely cause · 75%')).toBeInTheDocument();
    expect(within(causes[1]).queryByText(/Likely cause/)).not.toBeInTheDocument();

    const firstRemediations = within(causes[0])
      .getAllByTestId('error-remediation')
      .map((el) => el.textContent);
    expect(firstRemediations).toEqual([
      'Refresh the credentials.Runs automatically',
      'Reconnect the connection once the cluster is reachable.Requires approval',
    ]);
    // The shared remediation is shown under both causes; a manual one has no
    // policy label.
    const secondRemediations = within(causes[1])
      .getAllByTestId('error-remediation')
      .map((el) => el.textContent);
    expect(secondRemediations).toEqual([
      'Check the network.',
      'Reconnect the connection once the cluster is reachable.Requires approval',
    ]);

    // No flat lists when the causes are paired.
    expect(screen.queryByText('Suggested Remediation')).not.toBeInTheDocument();
  });

  it('links the error code reference page', () => {
    render(<ErrorMetadataFormatter metadata={eventError} />);
    const link = screen.getByRole('link', { name: 'meshery-server-1304 reference' });
    expect(link).toHaveAttribute('href', eventError.referenceUrl);
    expect(link).toHaveAttribute('target', '_blank');
  });

  it('falls back to the flat lists for a legacy MeshKit error', () => {
    render(
      <ErrorMetadataFormatter
        metadata={{
          Code: 'meshery-server-1033',
          LongDescription: ['First line', 'Second line'],
          ProbableCause: ['Cause A'],
          SuggestedRemediation: ['Fix B', 'Fix C'],
        }}
        event={{ description: 'short' }}
      />,
    );

    expect(screen.queryByTestId('error-causes')).not.toBeInTheDocument();
    expect(screen.getByText('Details')).toBeInTheDocument();
    expect(screen.getByText('Probable Cause')).toBeInTheDocument();
    expect(screen.getByText('Suggested Remediation')).toBeInTheDocument();
    expect(screen.getAllByTestId('markdown').map((el) => el.textContent)).toEqual([
      'First line',
      'Second line',
      'Cause A',
      'Fix B',
      'Fix C',
    ]);
    expect(screen.queryByRole('link')).not.toBeInTheDocument();
  });

  it('falls back to the flat lists for an EventError that carries only them', () => {
    render(
      <ErrorMetadataFormatter
        metadata={{ code: 'meshery-server-1', probableCause: ['A'], suggestedRemediation: ['B'] }}
      />,
    );
    expect(screen.queryByTestId('error-causes')).not.toBeInTheDocument();
    expect(screen.getByText('Probable Cause')).toBeInTheDocument();
  });

  it('shows the event description only when the title had to truncate it', () => {
    const long = 'x'.repeat(80);
    const { rerender } = render(
      <ErrorMetadataFormatter metadata={eventError} event={{ description: long }} />,
    );
    expect(screen.getByTestId('fsd-string')).toHaveTextContent(long);

    rerender(<ErrorMetadataFormatter metadata={eventError} event={{ description: 'short' }} />);
    expect(screen.queryByTestId('fsd-string')).not.toBeInTheDocument();
  });

  it('renders nothing for a value that is not an error object', () => {
    const { container } = render(<ErrorMetadataFormatter metadata="boom" />);
    expect(container).toBeEmptyDOMElement();
  });
});
