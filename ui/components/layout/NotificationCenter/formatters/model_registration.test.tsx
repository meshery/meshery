import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Typography: ({ children, ...props }: any) => <span {...props}>{children}</span>,
  Grid2: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  LaunchIcon: () => <svg data-testid="launch-icon" />,
}));

vi.mock('./error', () => ({
  ErrorMetadataFormatter: ({ metadata }: any) => (
    <div data-testid="error-metadata">{JSON.stringify(metadata)}</div>
  ),
}));

vi.mock('./common', () => ({
  TitleLink: ({ href, children, target }: any) => (
    <a data-testid="title-link" href={href} target={target}>
      {children}
    </a>
  ),
}));

vi.mock('@/constants/common', () => ({
  FALLBACK_MESHERY_IMAGE_PATH: '/fallback.png',
}));

vi.mock('@/utils/fallback', () => ({
  normalizeStaticImagePath: (src: string) => src,
}));

vi.mock('css/icons.styles', () => ({
  iconMedium: {},
  iconSmall: {},
}));

// jsdom doesn't implement fetch by default; provide a stub that always rejects
// so the component falls back to the default image path.
beforeEach(() => {
  (global as any).fetch = vi.fn(() => Promise.reject(new Error('no fetch')));
});

import { ModelImportedSection, ModelImportMessages } from './model_registration';
import { beforeEach } from 'vitest';

describe('ModelImportMessages', () => {
  it('renders the summary heading and message body', () => {
    render(<ModelImportMessages message="All components imported" />);
    const wrapper = screen.getByTestId('ModelImportMessages-Wrapper');
    expect(wrapper).toHaveTextContent(/summary/i);
    expect(wrapper).toHaveTextContent('All components imported');
  });
});

describe('ModelImportedSection', () => {
  it('renders null for non-object input', () => {
    const { container } = render(<ModelImportedSection modelDetails={'not-an-object'} />);
    expect(container).toBeEmptyDOMElement();
  });

  it('renders null for an array (the type guard rejects arrays)', () => {
    const { container } = render(<ModelImportedSection modelDetails={[]} />);
    expect(container).toBeEmptyDOMElement();
  });

  it('renders a MODEL header with a Registry TitleLink for normal models', () => {
    render(
      <ModelImportedSection
        modelDetails={{
          'core-mesh': {
            Components: [
              { DisplayName: 'PodComp', Metadata: 'Pod', Model: 'core-mesh', Version: '1.0.0' },
            ],
            Relationships: [],
            Errors: [],
          },
        }}
      />,
    );

    expect(screen.getByTestId('ModelImportedSection-ModelHeader-core-mesh')).toBeInTheDocument();
    expect(screen.getByText(/MODEL:/)).toBeInTheDocument();
    expect(screen.getByText('core-mesh')).toBeInTheDocument();
    expect(screen.getByTestId('title-link')).toHaveAttribute(
      'href',
      'settings?settingsCategory=Registry&tab=Models&searchText=core-mesh',
    );
    expect(screen.getByText('PodComp')).toBeInTheDocument();
  });

  it('renders FILE NAME header with no TitleLink for yaml/yml/json entity files', () => {
    render(
      <ModelImportedSection
        modelDetails={{
          'design.yaml': {
            Components: [],
            Relationships: [],
            Errors: [],
          },
        }}
      />,
    );

    expect(screen.getByText(/FILE NAME:/)).toBeInTheDocument();
    expect(screen.queryByTestId('title-link')).not.toBeInTheDocument();
  });

  it('renders relationship sections with FROM / TO information', () => {
    render(
      <ModelImportedSection
        modelDetails={{
          'foo-mesh': {
            Components: [],
            Relationships: [
              {
                Kind: 'edge',
                Subtype: 'network',
                RelationshipType: 'tcp',
                Selectors: [
                  {
                    allow: {
                      from: [{ kind: 'pod' }],
                      to: [{ kind: 'service' }],
                    },
                  },
                ],
              },
            ],
            Errors: [],
          },
        }}
      />,
    );

    expect(screen.getByText('FROM')).toBeInTheDocument();
    expect(screen.getByText('TO')).toBeInTheDocument();
    expect(screen.getByText('pod')).toBeInTheDocument();
    expect(screen.getByText('service')).toBeInTheDocument();
    expect(screen.getByText(/Kind of edge, sub type network and type tcp/i)).toBeInTheDocument();
  });

  it('renders an error summary for failed imports of normal models', () => {
    render(
      <ModelImportedSection
        modelDetails={{
          'broken-model': {
            Components: [],
            Relationships: [],
            Errors: [
              {
                name: ['broken-model', 'broken-model'],
                entityType: ['component', 'component'],
                error: { Code: 'E1' },
              } as any,
            ],
          },
        }}
      />,
    );

    expect(
      screen.getByText(/Import did not occur for 2 entities of type component\./),
    ).toBeInTheDocument();
    expect(screen.getByTestId('error-metadata')).toHaveTextContent('"Code":"E1"');
  });

  it('renders a file-import error block for yaml/json file errors', () => {
    render(
      <ModelImportedSection
        modelDetails={{
          'broken.yaml': {
            Components: [],
            Relationships: [],
            Errors: [
              {
                name: ['broken.yaml'],
                entityType: ['component'],
                error: { Code: 'E2' },
              } as any,
            ],
          },
        }}
      />,
    );

    expect(
      screen.getByText(/Import process for file broken\.yaml encountered error/),
    ).toBeInTheDocument();
    expect(screen.getByTestId('error-metadata')).toHaveTextContent('"Code":"E2"');
  });

  it('uses plural COMPONENTS heading when there is more than one component', () => {
    render(
      <ModelImportedSection
        modelDetails={{
          'multi-model': {
            Components: [
              { DisplayName: 'C1', Metadata: 'Pod', Model: 'multi-model' },
              { DisplayName: 'C2', Metadata: 'Service', Model: 'multi-model' },
            ],
            Relationships: [],
            Errors: [],
          },
        }}
      />,
    );

    expect(screen.getByText('COMPONENTS:')).toBeInTheDocument();
    expect(screen.getByText('C1')).toBeInTheDocument();
    expect(screen.getByText('C2')).toBeInTheDocument();
  });
});
