import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const { notify, resultsState, fetchPage } = vi.hoisted(() => ({
  notify: vi.fn(),
  fetchPage: vi.fn(),
  resultsState: { current: {} as Record<string, unknown> },
}));

vi.mock('react-big-calendar', () => ({ momentLocalizer: () => ({}) }));
vi.mock('react-big-calendar/lib/css/react-big-calendar.css', () => ({}));
vi.mock('react-redux', () => ({ useDispatch: () => vi.fn() }));
vi.mock('@/store/slices/mesheryUi', () => ({ updateProgressAction: vi.fn() }));
vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));
vi.mock('@meshery/schemas/mesheryApi', () => ({
  useLazyGetPerformanceResultsQuery: () => [fetchPage, resultsState.current],
}));
vi.mock('./style', () => ({
  CalendarComponent: ({ events }) =>
    events.map((event) => <span key={event.resource}>{event.title}</span>),
}));
vi.mock('../shared/Modal/GenericModal', () => ({ default: () => null }));
vi.mock('../MesheryChart', () => ({ default: () => null }));
vi.mock('@sistent/sistent', () => ({
  Typography: ({ children }: any) => <span>{children}</span>,
  Paper: ({ children }: any) => <div>{children}</div>,
}));

import PerformanceCalendar from './PerformanceCalendar';

describe('PerformanceCalendar', () => {
  beforeEach(() => {
    notify.mockClear();
    fetchPage.mockReset().mockResolvedValue({ error: {} });
  });

  it('shows the server error instead of [object Object]', () => {
    resultsState.current = {
      isFetching: false,
      isError: true,
      error: { status: 500, data: 'failed to fetch results: database is locked' },
    };

    render(<PerformanceCalendar />);

    expect(notify).toHaveBeenCalledTimes(1);
    const { message } = notify.mock.calls[0][0];
    expect(message).toBe('failed to fetch results: database is locked');
  });

  it('falls back to a readable message when the error has no details', () => {
    resultsState.current = { isFetching: false, isError: true, error: { status: 'FETCH_ERROR' } };

    render(<PerformanceCalendar />);

    expect(notify).toHaveBeenCalledTimes(1);
    expect(notify.mock.calls[0][0].message).toBe('Failed to Fetch Results');
  });

  it('renders runs beyond the first page', async () => {
    const runs = Array.from({ length: 11 }, (_, i) => ({
      name: `Run ${i + 1}`,
      testStartTime: '2026-11-01T12:00:00Z',
      runnerResults: { ActualDuration: 1e9 },
    }));
    resultsState.current = {};
    fetchPage.mockImplementation(async ({ page }) => ({
      data: { results: page === '0' ? runs.slice(0, 10) : runs.slice(10), totalCount: 11 },
    }));

    render(<PerformanceCalendar />);

    expect(await screen.findByText('Run 11')).toBeInTheDocument();
    expect(screen.getAllByText(/^Run /)).toHaveLength(11);
    expect(fetchPage).toHaveBeenCalledTimes(2);
  });

  it('keeps fetched runs when a later page fails', async () => {
    const runs = Array.from({ length: 10 }, (_, i) => ({
      name: `Run ${i + 1}`,
      testStartTime: '2026-11-01T12:00:00Z',
      runnerResults: { ActualDuration: 1e9 },
    }));
    resultsState.current = {};
    fetchPage.mockImplementation(async ({ page }) =>
      page === '0' ? { data: { results: runs, totalCount: 11 } } : { isError: true },
    );

    render(<PerformanceCalendar />);

    expect(await screen.findByText('Run 10')).toBeInTheDocument();
    expect(screen.getAllByText(/^Run /)).toHaveLength(10);
  });
});
