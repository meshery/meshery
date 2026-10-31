import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('./general/MesheryChartDialog', () => ({
  default: ({ open, content, title }: any) =>
    open ? (
      <div data-testid="chart-dialog" data-title={title}>
        {content}
      </div>
    ) : null,
}));

import MesheryResultDialog from './MesheryResultDialog';

const runnerResults = {
  URL: 'http://productpage:9080',
  RequestedQPS: '10',
  ActualQPS: 9.5,
  NumThreads: 2,
  SocketCount: 4,
  RequestedDuration: '30s',
  ActualDuration: 30000000000,
  DurationHistogram: { Avg: 0.01, Max: 0.2, Percentiles: [{ Percentile: 50, Value: 0.009 }] },
};

describe('MesheryResultDialog', () => {
  it('renders the load profile of a camelCase result', () => {
    render(<MesheryResultDialog rowData={{ name: 'run-1', runnerResults }} close={vi.fn()} />);

    expect(screen.getByText('http://productpage:9080')).toBeInTheDocument();
    expect(screen.getByText('30.0')).toBeInTheDocument();
    expect(screen.getByText('p50 Response Time:')).toBeInTheDocument();
  });

  it('still renders a snake_case result', () => {
    render(
      <MesheryResultDialog
        rowData={{ name: 'run-1', runner_results: runnerResults }}
        close={vi.fn()}
      />,
    );

    expect(screen.getByText('http://productpage:9080')).toBeInTheDocument();
  });

  it('renders the load profile when the result has no duration histogram', () => {
    const { DurationHistogram, ...withoutHistogram } = runnerResults;
    render(
      <MesheryResultDialog
        rowData={{ name: 'run-1', runnerResults: withoutHistogram }}
        close={vi.fn()}
      />,
    );

    expect(screen.getByText('http://productpage:9080')).toBeInTheDocument();
    expect(screen.queryByText('Average Response Time:')).not.toBeInTheDocument();
  });
});
