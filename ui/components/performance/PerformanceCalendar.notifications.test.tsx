import React from 'react';
import { render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const { notify, resultsState } = vi.hoisted(() => ({
  notify: vi.fn(),
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
  useGetPerformanceResultsQuery: () => resultsState.current,
}));
vi.mock('./style', () => ({ CalendarComponent: () => null }));
vi.mock('../shared/Modal/GenericModal', () => ({ default: () => null }));
vi.mock('../MesheryChart', () => ({ default: () => null }));
vi.mock('@sistent/sistent', () => ({
  Typography: ({ children }: any) => <span>{children}</span>,
  Paper: ({ children }: any) => <div>{children}</div>,
}));

import PerformanceCalendar from './PerformanceCalendar';

describe('PerformanceCalendar error notification', () => {
  beforeEach(() => {
    notify.mockClear();
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
});
