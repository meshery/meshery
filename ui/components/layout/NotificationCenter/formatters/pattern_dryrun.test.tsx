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

describe('SchemaValidationFormatter', () => {
  it('renders defaults for an empty event', () => {
    render(<SchemaValidationFormatter event={{} as any} />);
    expect(screen.getByTestId('vr-design')).toHaveTextContent('Unknown Design');
    expect(screen.getByTestId('vr-errors')).toHaveTextContent('0');
    expect(screen.getByTestId('vr-comp-count')).toHaveTextContent('0');
    expect(screen.getByTestId('vr-annotation-count')).toHaveTextContent('0');
  });

  it('passes through metadata counts and totals', () => {
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
