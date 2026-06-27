import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const generateMock = vi.fn(() => ({ destroy: vi.fn() }));

vi.mock('billboard.js', () => ({
  default: {
    generate: (...args: any[]) => generateMock(...args),
    areaStep: () => 'area-step',
    line: () => 'line',
  },
  bb: undefined,
  // also export named so import { areaStep, line } works
  areaStep: () => 'area-step',
  line: () => 'line',
}));

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => (props: any) => <Component {...props} />;
  return {
    Grid2: ({ children }: any) => <div>{children}</div>,
    Typography: ({ children }: any) => <span>{children}</span>,
    IconButton: ({ children, onClick, 'aria-label': aria }: any) => (
      <button onClick={onClick} aria-label={aria}>
        {children}
      </button>
    ),
    Paper: ({ children }: any) => <div>{children}</div>,
    ClickAwayListener: ({ children }: any) => (
      <div>{typeof children === 'function' ? null : children}</div>
    ),
    Fade: ({ children }: any) => <div>{children}</div>,
    Popper: ({ children, open }: any) =>
      open ? (
        <div data-testid="popper">{typeof children === 'function' ? children({}) : children}</div>
      ) : null,
    NoSsr: ({ children }: any) => <>{children}</>,
    ReplyIcon: () => <svg data-testid="reply-icon" />,
    styled,
  };
});

vi.mock('../lib/chartjs-formatter', () => ({
  fortioResultToJsChartData: (rawdata: any, data: any) => ({ rawdata, data }),
  makeChart: (data: any) => ({
    data: { datasets: [{ label: 'd1', data: [{ x: 1, y: 2 }], cubicInterpolationMode: false }] },
    options: {
      scales: { xAxes: [], yAxes: [] },
      title: { text: ['title-line'] },
      metadata: { qps: { display: { value: 'rps 100' } } },
    },
    percentiles: [
      { Percentile: 50, Value: 0.1 },
      { Percentile: 75, Value: 0.2 },
      { Percentile: 90, Value: 0.3 },
      { Percentile: 99, Value: 0.4 },
      { Percentile: 99.9, Value: 0.5 },
    ],
    raw: data,
  }),
  makeOverlayChart: () => ({
    data: { datasets: [] },
    options: { scales: {}, title: { text: [] } },
  }),
  makeMultiChart: () => ({
    data: { datasets: [], labels: [] },
    options: { scales: { yAxes: [] }, title: { text: 'multi' } },
  }),
}));

vi.mock('react-share', () => ({
  TwitterShareButton: ({ children }: any) => <div>{children}</div>,
  LinkedinShareButton: ({ children }: any) => <div>{children}</div>,
  FacebookShareButton: ({ children }: any) => <div>{children}</div>,
  TwitterIcon: () => <svg data-testid="twitter-icon" />,
  LinkedinIcon: () => <svg data-testid="linkedin-icon" />,
  FacebookIcon: () => <svg data-testid="facebook-icon" />,
}));

import MesheryChart from './MesheryChart';

describe('MesheryChart', () => {
  beforeEach(() => {
    generateMock.mockClear();
  });

  it('renders a share icon and percentile summary for a single result', () => {
    render(
      <MesheryChart
        rawdata={{}}
        hideTitle
        data={[
          {
            StartTime: 1,
            DurationHistogram: {},
          },
        ]}
      />,
    );
    expect(screen.getByLabelText('Share')).toBeInTheDocument();
    expect(screen.getByText('Percentile Summary')).toBeInTheDocument();
  });

  it('does not render the percentile summary for overlay (2-dataset) charts', () => {
    render(<MesheryChart rawdata={{}} hideTitle data={[{ StartTime: 1 }, { StartTime: 2 }]} />);
    expect(screen.queryByText('Percentile Summary')).not.toBeInTheDocument();
  });

  it('opens the share popper when the share button is clicked', () => {
    render(<MesheryChart rawdata={{}} hideTitle data={[{ StartTime: 1 }]} />);
    fireEvent.click(screen.getByLabelText('Share'));
    expect(screen.getByTestId('popper')).toBeInTheDocument();
  });

  it('renders without crashing when data is undefined', () => {
    expect(() => render(<MesheryChart rawdata={{}} data={[]} />)).not.toThrow();
  });
});
