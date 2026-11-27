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
  });
});
