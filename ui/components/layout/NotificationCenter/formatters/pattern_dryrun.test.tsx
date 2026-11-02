import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const formatDryRunResponseMock = vi.fn();

vi.mock('../../../designs/lifecycle/DryRun', () => ({
  FormatDryRunResponse: ({ dryRunErrors }: any) => (
    <div data-testid="dry-run-format">{JSON.stringify(dryRunErrors)}</div>
  ),
}));

vi.mock('machines/validator/designValidator', () => ({
  formatDryRunResponse: (...args: unknown[]) => formatDryRunResponseMock(...args),
}));

vi.mock('@/components/designs/lifecycle/ValidateDesign', () => ({
  ValidationResults: ({
    validationResults,
    errorCount,
    compCount,
    annotationCount,
    design,
  }: any) => (
    <div data-testid="validation-results">
      <span data-testid="vr-design">{design}</span>
      <span data-testid="vr-errors">{errorCount}</span>
      <span data-testid="vr-comp-count">{compCount}</span>
      <span data-testid="vr-annotation-count">{annotationCount}</span>
      <span data-testid="vr-raw">{JSON.stringify(validationResults)}</span>
    </div>
  ),
}));

import { DryRunResponse, SchemaValidationFormatter } from './pattern_dryrun';

describe('DryRunResponse', () => {
  beforeEach(() => {
    formatDryRunResponseMock.mockReset();
  });

  it('formats the response via formatDryRunResponse before rendering', () => {
    formatDryRunResponseMock.mockReturnValue([{ kind: 'Pod' }]);
    render(<DryRunResponse response={{ raw: 'response' }} />);
    expect(formatDryRunResponseMock).toHaveBeenCalledWith({ raw: 'response' });
    expect(screen.getByTestId('dry-run-format')).toHaveTextContent('Pod');
  });
});

const errorLines = () =>
  screen.getAllByTestId('design-event-error').map((line) => line.textContent);

describe('SchemaValidationFormatter', () => {
  it('renders the design placeholder for an empty event', () => {
    render(<SchemaValidationFormatter event={{} as any} />);
    expect(screen.getByText(/Design:/)).toHaveTextContent('Design: Unknown Design');
    expect(screen.queryByTestId('validation-results')).not.toBeInTheDocument();
  });

  it('renders status, errorCount and the errors of a new-format validate event', () => {
    render(
      <SchemaValidationFormatter
        event={{
          metadata: {
            designName: 'my-design',
            status: 'failed',
            errorCount: 3,
            errors: [
              { code: 'required', componentId: 'c1' },
              { code: 'type', componentId: 'c2' },
              { code: 'meshery-server-1480', componentId: 'c3' },
            ],
            truncated: false,
          },
        }}
      />,
    );

    expect(screen.getByText(/Design:/)).toHaveTextContent('Design: my-design');
    expect(screen.getByText(/Status:/)).toHaveTextContent('Status: failed - 3 errors');
    expect(errorLines()).toEqual([
      'required - component c1',
      'type - component c2',
      'meshery-server-1480 - component c3',
    ]);
    // A JSON-schema keyword is shown as-is; only a MeshKit code is linked.
    expect(screen.getByText(/required/).closest('a')).toBeNull();
    expect(screen.getByText('meshery-server-1480').closest('a')).toHaveAttribute(
      'href',
      'https://docs.meshery.io/reference/error-codes#meshery-server-1480',
    );
    expect(screen.queryByTestId('validation-results')).not.toBeInTheDocument();
  });

  it('falls back to the legacy design_name of a rewritten validate event', () => {
    render(
      <SchemaValidationFormatter
        event={{
          metadata: {
            design_name: 'legacy-design',
            status: 'failed',
            errorCount: 1,
            errors: [{ code: 'required', componentId: 'c1' }],
            truncated: false,
          },
        }}
      />,
    );
    expect(screen.getByText(/Design:/)).toHaveTextContent('Design: legacy-design');
    expect(screen.queryByText(/undefined/)).not.toBeInTheDocument();
  });

  it('prefers designName over the legacy design_name', () => {
    render(
      <SchemaValidationFormatter
        event={{
          metadata: {
            designName: 'new-name',
            design_name: 'old-name',
            status: 'success',
            errorCount: 0,
            errors: [],
            truncated: false,
          },
        }}
      />,
    );
    expect(screen.getByText(/Design:/)).toHaveTextContent('Design: new-name');
    expect(screen.getByText(/Status:/)).toHaveTextContent('Status: success - 0 errors');
  });

  it('renders the new-format summary even when a legacy validationResult remains', () => {
    render(
      <SchemaValidationFormatter
        event={{
          metadata: {
            designName: 'my-design',
            status: 'failed',
            errorCount: 1,
            errors: [{ code: 'required', componentId: 'c1' }],
            truncated: false,
            validationResult: { c1: { errors: [1] } },
          },
        }}
      />,
    );
    expect(errorLines()).toEqual(['required - component c1']);
    expect(screen.queryByTestId('validation-results')).not.toBeInTheDocument();
  });

  it('passes legacy validationResult counts and totals through', () => {
    render(
      <SchemaValidationFormatter
        event={{
          metadata: {
            validationResult: {
              serviceA: { errors: [1, 2] },
              serviceB: { errors: [1] },
              serviceC: {},
            },
            design_name: 'my-design',
            total_components: 5,
            configurable_components: 3,
          },
        }}
      />,
    );

    expect(screen.getByTestId('vr-design')).toHaveTextContent('my-design');
    expect(screen.getByTestId('vr-errors')).toHaveTextContent('3');
    expect(screen.getByTestId('vr-comp-count')).toHaveTextContent('3');
    expect(screen.getByTestId('vr-annotation-count')).toHaveTextContent('2');
  });
});
