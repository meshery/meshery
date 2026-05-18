import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Typography: ({ children, ...props }: any) => <span {...props}>{children}</span>,
  Grid2: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  List: ({ children, ...props }: any) => <ul {...props}>{children}</ul>,
  ListItem: ({ children, ...props }: any) => <li {...props}>{children}</li>,
  RenderMarkdown: ({ content }: any) => <span data-testid="markdown">{content}</span>,
}));

vi.mock('../../../data-formatter', () => ({
  FormatStructuredData: ({ data, propertyFormatters }: any) => {
    // Render keys + the formatter output so tests can assert which formatter
    // was applied and the source data flowed in correctly.
    return (
      <div data-testid="format-structured">
        {Object.entries(data || {}).map(([key, value]) => {
          const Formatter = propertyFormatters?.[key];
          return (
            <div key={key} data-key={key}>
              <span data-testid="fsd-key">{key}</span>
              {Formatter ? Formatter(value) : <span data-testid="fsd-raw">{String(value)}</span>}
            </div>
          );
        })}
      </div>
    );
  },
}));

import { ErrorMetadataFormatter } from './error';

describe('ErrorMetadataFormatter', () => {
  it('renders the event description, details, cause, and remediation sections', () => {
    render(
      <ErrorMetadataFormatter
        metadata={{
          LongDescription: ['First line', 'Second line'],
          ProbableCause: ['Cause A'],
          SuggestedRemediation: ['Fix B', 'Fix C'],
        }}
        event={{ description: 'Top-level description' }}
      />,
    );

    expect(screen.getAllByTestId('format-structured').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Details').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Probable Cause').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Suggested Remediation').length).toBeGreaterThan(0);
    expect(screen.getAllByTestId('markdown').map((el) => el.textContent)).toEqual(
      expect.arrayContaining(['First line', 'Second line', 'Cause A', 'Fix B', 'Fix C']),
    );
  });

  it('defaults empty arrays when metadata fields are missing', () => {
    render(<ErrorMetadataFormatter metadata={{} as any} event={{}} />);
    // The headings should still render even with empty value arrays.
    expect(screen.getAllByText('Details').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Probable Cause').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Suggested Remediation').length).toBeGreaterThan(0);
  });
});
