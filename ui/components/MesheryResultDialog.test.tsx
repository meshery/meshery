import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const passThrough =
    (tag: string) =>
    ({ children }: any) =>
      React.createElement(tag, null, children);
  return {
    styled: (Comp: any) => () => Comp,
    Grid2: passThrough('div'),
    Table: passThrough('table'),
    Typography: passThrough('span'),
    NoSsr: passThrough('div'),
    TableRow: passThrough('tr'),
    TableCell: passThrough('td'),
    TableBody: passThrough('tbody'),
  };
});

vi.mock('./general/MesheryChartDialog', () => ({
  default: ({ content }: any) => <div data-testid="chart-dialog">{content}</div>,
}));

import MesheryResultDialog from './MesheryResultDialog';

describe('MesheryResultDialog', () => {
  it('renders without crashing when kubernetes.nodes or runner_results.DurationHistogram is missing', () => {
    render(
      <MesheryResultDialog
        close={vi.fn()}
        rowData={{
          name: 'test-run',
          runner_results: {
            URL: 'http://example.com',
            RequestedQPS: '10',
            ActualQPS: 10,
            NumThreads: 2,
            SocketCount: 2,
            RequestedDuration: '5s',
            ActualDuration: 5000000000,
            kubernetes: {
              server_version: 'v1.30.0',
            },
          },
        }}
      />,
    );

    expect(screen.getByText('Environment')).toBeInTheDocument();
    expect(screen.getByText('v1.30.0')).toBeInTheDocument();
    expect(screen.getByText('Load Profile')).toBeInTheDocument();
    expect(screen.getByText('Average Response Time:')).toBeInTheDocument();
    expect(screen.getByText('Maximum Response Time:')).toBeInTheDocument();
  });

  it('renders without crashing when kubernetes.nodes is a non-array and renders DurationHistogram values when present', () => {
    render(
      <MesheryResultDialog
        close={vi.fn()}
        rowData={{
          name: 'non-array-nodes-run',
          runner_results: {
            URL: 'http://example.com/2',
            RequestedQPS: '20',
            ActualQPS: 20,
            NumThreads: 4,
            SocketCount: 4,
            RequestedDuration: '10s',
            ActualDuration: 10000000000,
            DurationHistogram: {
              Avg: 0.015,
              Max: 0.09,
            },
            kubernetes: {
              server_version: 'v1.31.0',
              nodes: { unexpected: true },
            },
          },
        }}
      />,
    );

    expect(screen.getByText('v1.31.0')).toBeInTheDocument();
    expect(screen.getByText('0.015')).toBeInTheDocument();
    expect(screen.getByText('0.09')).toBeInTheDocument();
  });
});
