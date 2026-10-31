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

vi.mock('./common', () => ({
  TitleLink: ({ href, children }: any) => (
    <a data-testid="title-link" href={href}>
      {children}
    </a>
  ),
  formatErrorCode: (code: any, component = 'meshery-server') => {
    if (!code) return '';
    const trimmed = String(code).trim();
    if (trimmed.includes('-')) return trimmed;
    if (/^\d+$/.test(trimmed)) return `${component}-${trimmed}`;
    return trimmed;
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

  it('renders the universal error code link when Code is provided in metadata', () => {
    render(
      <ErrorMetadataFormatter
        metadata={{
          Code: '1562',
          LongDescription: ['Something went wrong'],
        }}
        event={{ description: 'An error occurred' }}
      />,
    );
    const link = screen.getByTestId('title-link');
    expect(link).toHaveAttribute(
      'href',
      'https://docs.meshery.io/reference/error-codes#meshery-server-1562',
    );
    expect(link).toHaveTextContent('meshery-server-1562');
  });

  it('renders existing full mnemonic error code without double-prefixing', () => {
    render(
      <ErrorMetadataFormatter
        metadata={{
          Code: 'meshery-server-1133',
        }}
        event={{ description: 'Conversion error' }}
      />,
    );
    const link = screen.getByTestId('title-link');
    expect(link).toHaveAttribute(
      'href',
      'https://docs.meshery.io/reference/error-codes#meshery-server-1133',
    );
    expect(link).toHaveTextContent('meshery-server-1133');
  });
});
