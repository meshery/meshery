import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  useTheme: () => ({
    palette: {
      text: { default: '#000' },
    },
  }),
  DownloadIcon: () => <svg data-testid="download-icon" />,
  InfoIcon: () => <svg data-testid="info-icon" />,
  createTheme: () => ({ breakpoints: { up: () => '', down: () => '' } }),
  notificationColors: {
    info: { main: '#info' },
    error: { main: '#error', dark: '#error_dark' },
    warning: { main: '#warning', light: '#warning_light' },
    success: { main: '#success' },
  },
}));

vi.mock('../../../assets/icons/AlertIcon', () => ({ default: () => null }));
vi.mock('../../../assets/icons/ErrorIcon', () => ({ default: () => null }));
vi.mock('../../../assets/icons/ReadIcon', () => ({ default: () => null }));

vi.mock('../../data-formatter', () => ({
  FormatStructuredData: ({ data, order, propertyFormatters }: any) => (
    <div data-testid="fsd">
      <span data-testid="fsd-order">{JSON.stringify(order || [])}</span>
      <span data-testid="fsd-keys">{Object.keys(data || {}).join(',')}</span>
      <span data-testid="fsd-data">{JSON.stringify(data)}</span>
      {Object.entries(data || {}).map(([key, value]) => {
        const Formatter = propertyFormatters?.[key];
        return (
          <div key={key} data-key={key}>
            {Formatter ? Formatter(value) : <span>{String(value)}</span>}
          </div>
        );
      })}
    </div>
  ),
  SectionBody: ({ body }: any) => <div data-testid="section-body">{body}</div>,
  reorderObjectProperties: (obj: any) => obj,
}));

vi.mock('../../../utils/objects', () => ({
  isEmptyAtAllDepths: (obj: any) => !obj || Object.keys(obj).length === 0,
}));

vi.mock('./notification', () => ({
  canTruncateDescription: (s: string) => (s || '').length > 62,
}));

vi.mock('../../designs/lifecycle/DeploymentSummary', () => ({
  DeploymentSummaryFormatter: ({ event }: any) => (
    <div data-testid="deployment-summary">{event?.description ?? ''}</div>
  ),
}));

vi.mock('./formatters/common', () => ({
  TitleLink: ({ href, children }: any) => (
    <a data-testid="title-link" href={href}>
      {children}
    </a>
  ),
  DataToFileLink: ({ data }: any) => (
    <div data-testid="data-to-file-link">
      {typeof data === 'string' ? data : JSON.stringify(data)}
    </div>
  ),
  EmptyState: ({ event }: any) => <div data-testid="empty-state">{event?.description ?? ''}</div>,
}));

vi.mock('./formatters/error', () => ({
  ErrorMetadataFormatter: ({ metadata }: any) => (
    <div data-testid="error-metadata">{JSON.stringify(metadata)}</div>
  ),
}));

vi.mock('./formatters/pattern_dryrun', () => ({
  DryRunResponse: ({ response }: any) => (
    <div data-testid="dry-run">{JSON.stringify(response)}</div>
  ),
  SchemaValidationFormatter: ({ event }: any) => (
    <div data-testid="schema-validation">{event?.metadata?.design_name || ''}</div>
  ),
}));

vi.mock('./formatters/model_registration', () => ({
  ModelImportMessages: ({ message }: any) => (
    <div data-testid="model-import-messages">{message}</div>
  ),
  ModelImportedSection: ({ modelDetails }: any) => (
    <div data-testid="model-imported-section">{Object.keys(modelDetails || {}).join(',')}</div>
  ),
}));

vi.mock('./formatters/relationship_evaluation', () => ({
  RelationshipEvaluationEventFormatter: ({ event }: any) => (
    <div data-testid="rel-eval">{event?.description ?? ''}</div>
  ),
}));

vi.mock('./formatters/meshsync_events', () => ({
  MeshSyncPropertyFormatters: {
    connectionID: (value: any) => <span data-testid="meshsync-conn-id">{value}</span>,
  },
}));

vi.mock('./formatters/academy_events', () => ({
  AcademyEventsFormatter: ({ event }: any) => (
    <div data-testid="academy">{event?.metadata?.result?.quiz?.parent?.title ?? ''}</div>
  ),
}));

vi.mock('../../connections/styles', () => ({
  ChipWrapper: ({ label, href, target }: any) => (
    <a data-testid="chip-wrapper" href={href} target={target}>
      {label}
    </a>
  ),
}));

import {
  PropertyFormatters,
  LinkFormatters,
  PropertyLinkFormatters,
  FormattedMetadata,
  FormattedLinkMetadata,
} from './metadata';

describe('PropertyFormatters', () => {
  it('formats trace via DataToFileLink', () => {
    render(<>{PropertyFormatters.trace('trace-data')}</>);
    expect(screen.getByTestId('data-to-file-link')).toHaveTextContent('trace-data');
  });

  it('renders ShortDescription via SectionBody', () => {
    render(<>{PropertyFormatters.ShortDescription('short')}</>);
    expect(screen.getByTestId('section-body')).toHaveTextContent('short');
  });

  it('renders saved design link with encoded id and the design name', () => {
    render(<>{PropertyFormatters.design({ name: 'My Design', id: 'abc/123' })}</>);
    const link = screen.getByTestId('title-link');
    expect(link).toHaveAttribute(
      'href',
      '/extension/meshmap?mode=design&design=' + encodeURIComponent('abc/123'),
    );
    expect(link).toHaveTextContent('Saved design My Design');
  });

  it('renders connectionName as a ChipWrapper link', () => {
    render(<>{PropertyFormatters.connectionName('my-conn name')}</>);
    const chip = screen.getByTestId('chip-wrapper');
    expect(chip).toHaveAttribute(
      'href',
      `/management/connections?tab=connections&searchText=${encodeURIComponent('my-conn name')}`,
    );
    expect(chip).toHaveTextContent('my-conn name');
  });

  it('forwards ModelImportMessage and ModelDetails when non-empty', () => {
    render(<>{PropertyFormatters.ModelImportMessage('msg')}</>);
    expect(screen.getByTestId('model-import-messages')).toHaveTextContent('msg');

    render(<>{PropertyFormatters.ModelDetails({ 'core-mesh': {} })}</>);
    expect(screen.getByTestId('model-imported-section')).toHaveTextContent('core-mesh');
  });

  it('returns falsy for empty model-related fields', () => {
    expect(PropertyFormatters.ModelImportMessage(null)).toBe(null);
    expect(PropertyFormatters.ModelDetails(undefined)).toBe(undefined);
  });

  it('returns null for history_title field', () => {
    expect(PropertyFormatters.history_title('anything')).toBeNull();
  });

  it('exposes dryRunResponse formatter', () => {
    render(<>{PropertyFormatters.dryRunResponse({ raw: 1 })}</>);
    expect(screen.getByTestId('dry-run')).toHaveTextContent('"raw":1');
  });
});

describe('LinkFormatters and PropertyLinkFormatters', () => {
  it('renders doclink as a title link with end-aligned styling', () => {
    render(<>{LinkFormatters.doclink('https://docs.example')}</>);
    const link = screen.getByTestId('title-link');
    expect(link).toHaveAttribute('href', 'https://docs.example');
    expect(link).toHaveTextContent('Doc');
  });

  it('builds property link descriptors with the right URLs and labels', () => {
    expect(PropertyLinkFormatters.doc('http://x')).toEqual({
      label: 'Doc',
      href: 'http://x',
    });
    expect(PropertyLinkFormatters.DownloadLink('/path file')).toMatchObject({
      label: 'Download File',
      href: '/api/system/fileDownload?file=' + encodeURIComponent('/path file'),
    });
    expect(PropertyLinkFormatters.ViewLink('/log.txt')).toMatchObject({
      label: 'Get Logs',
      href: '/api/system/fileView?file=' + encodeURIComponent('/log.txt'),
    });
  });
});

describe('FormattedMetadata', () => {
  it('delegates to a registered event type formatter when matched', () => {
    render(
      <FormattedMetadata
        event={{
          action: 'deploy',
          category: 'pattern',
          description: 'deployed!',
        }}
      />,
    );
    expect(screen.getByTestId('deployment-summary')).toHaveTextContent('deployed!');
  });

  it('renders EmptyState when metadata is empty or missing', () => {
    render(<FormattedMetadata event={{ description: 'nothing here' }} />);
    expect(screen.getByTestId('empty-state')).toHaveTextContent('nothing here');
  });

  it('renders structured metadata via FormatStructuredData for generic events', () => {
    render(
      <FormattedMetadata
        event={{
          action: 'other',
          category: 'other',
          description: 'short',
          metadata: { connectionName: 'cluster-1' },
        }}
      />,
    );
    expect(screen.getByTestId('fsd')).toBeInTheDocument();
    expect(screen.getByTestId('fsd-keys')).toHaveTextContent('connectionName');
  });

  it('excludes ShortDescription when ImportedModelName is present', () => {
    const longDescription = 'x'.repeat(80);
    render(
      <FormattedMetadata
        event={{
          action: 'whatever',
          category: 'whatever',
          description: longDescription,
          metadata: { ImportedModelName: 'core-mesh' },
        }}
      />,
    );
    // ShortDescription must not be in the data passed to FSD
    const dataJson = screen.getByTestId('fsd-data').textContent || '{}';
    const parsed = JSON.parse(dataJson);
    expect(parsed.ShortDescription).toBeNull();
  });
});

describe('FormattedLinkMetadata', () => {
  it('only forwards doclink to FormatStructuredData', () => {
    render(
      <FormattedLinkMetadata
        event={{
          metadata: { doclink: 'https://docs.example', other: 'ignored' },
        }}
      />,
    );
    expect(screen.getByTestId('fsd-keys')).toHaveTextContent('doclink');
    expect(screen.getByTestId('fsd-keys')).not.toHaveTextContent('other');
  });
});
