import React from 'react';
import { act, cleanup, fireEvent, render, screen } from '@testing-library/react';
import { SistentThemeProvider } from '@sistent/sistent';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import MetricExplorer from './MetricExplorer';

const { previewHook, metricsHook } = vi.hoisted(() => ({
  previewHook: vi.fn(),
  metricsHook: vi.fn(),
}));
vi.mock('@/rtk-query/telemetryPrometheus', () => ({
  useLazyQueryPrometheusRangeQuery: previewHook,
  useListPrometheusMetricsQuery: metricsHook,
  usePrometheusMetricMetadataQuery: () => ({}),
}));
vi.mock('../common/TimeSeriesChart', () => ({ default: () => null }));

const sizeError = "Telemetry response exceeds Meshery's 16 MiB limit";
const renderExplorer = () =>
  render(
    <SistentThemeProvider initialMode="light">
      <MetricExplorer
        connectionID="prometheus-connection"
        timeWindow={{ start: 1, end: 2, step: 15 }}
        onAdd={vi.fn()}
      />
    </SistentThemeProvider>,
  );

afterEach(cleanup);
beforeEach(() => metricsHook.mockReturnValue({ data: [] }));

describe('MetricExplorer response errors', () => {
  it('shows the size limit and remedies without developer detail', () => {
    previewHook.mockReturnValue([
      vi.fn(),
      {
        isFetching: false,
        error: {
          status: 502,
          data: {
            error: sizeError,
            suggestedRemediation: ['Shorten the time range.', 'Use selective label matchers.'],
            longDescription: ['internal query path'],
          },
        },
      },
    ]);
    renderExplorer();
    expect(screen.getByRole('alert')).toHaveTextContent(sizeError);
    expect(screen.getAllByRole('listitem')).toHaveLength(2);
    expect(screen.queryByText('internal query path')).not.toBeInTheDocument();
  });

  it('does not mistake a rejected preview for an empty result', async () => {
    const trigger = vi.fn(() => ({ unwrap: () => Promise.reject(new Error('too large')) }));
    previewHook.mockReturnValue([
      trigger,
      { isFetching: false, error: { status: 502, data: { error: sizeError } } },
    ]);
    renderExplorer();
    fireEvent.change(screen.getByPlaceholderText('PromQL expression'), { target: { value: 'up' } });
    await act(async () => fireEvent.click(screen.getByRole('button', { name: 'Preview' })));
    expect(trigger).toHaveBeenCalled();
    expect(screen.getByRole('alert')).toHaveTextContent(sizeError);
    expect(screen.queryByText('No data for this query/range')).not.toBeInTheDocument();
  });

  it('shows a metric discovery error instead of an empty list', () => {
    metricsHook.mockReturnValue({ error: { status: 502, data: { error: sizeError } } });
    previewHook.mockReturnValue([vi.fn(), { isFetching: false }]);
    renderExplorer();
    expect(screen.getByRole('alert')).toHaveTextContent(sizeError);
    expect(screen.queryByText('Type to search metric names')).not.toBeInTheDocument();
  });
});
