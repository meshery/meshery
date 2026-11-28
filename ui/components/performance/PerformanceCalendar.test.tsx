import React from 'react';
import { render, act } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import PerformanceCalendar, { generateDateRange } from './PerformanceCalendar';

const { resultsQuerySpy, calendarComponentSpy, mockResults } = vi.hoisted(() => ({
  resultsQuerySpy: vi.fn(),
  calendarComponentSpy: vi.fn(),
  mockResults: { data: { results: [] }, isFetching: false, isError: false },
}));

vi.mock('react-big-calendar', () => ({ momentLocalizer: () => ({}) }));
vi.mock('react-big-calendar/lib/css/react-big-calendar.css', () => ({}));
vi.mock('react-redux', () => ({ useDispatch: () => vi.fn() }));
vi.mock('@/store/slices/mesheryUi', () => ({ updateProgressAction: vi.fn() }));
vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));
vi.mock('@meshery/schemas/mesheryApi', () => ({
  useGetPerformanceResultsQuery: (args: any) => {
    resultsQuerySpy(args);
    return mockResults;
  },
}));
vi.mock('./style', () => ({
  CalendarComponent: (props: any) => {
    calendarComponentSpy(props);
    return <div data-testid="calendar-mock" />;
  },
}));
vi.mock('../shared/Modal/GenericModal', () => ({ default: () => null }));
vi.mock('../MesheryChart', () => ({ default: () => null }));
vi.mock('@sistent/sistent', () => ({
  Typography: ({ children }: any) => <span>{children}</span>,
  Paper: ({ children }: any) => <div>{children}</div>,
}));

describe('generateDateRange', () => {
  it('handles Day view date array with a single date', () => {
    const day = new Date('2024-03-15T12:00:00Z');
    const range = generateDateRange([day]);
    expect(range).toEqual({
      start: '2024-03-15',
      end: '2024-03-15',
    });
  });

  it('handles Week view date array spanning a month boundary', () => {
    const days = [
      new Date('2024-01-28T12:00:00Z'),
      new Date('2024-01-29T12:00:00Z'),
      new Date('2024-01-30T12:00:00Z'),
      new Date('2024-01-31T12:00:00Z'),
      new Date('2024-02-01T12:00:00Z'),
      new Date('2024-02-02T12:00:00Z'),
      new Date('2024-02-03T12:00:00Z'),
    ];
    const range = generateDateRange(days);
    expect(range).toEqual({
      start: '2024-01-28',
      end: '2024-02-03',
    });
  });

  it('handles Month view range object with start and end', () => {
    const range = generateDateRange({
      start: new Date('2024-02-01T12:00:00Z'),
      end: new Date('2024-02-29T12:00:00Z'),
    });
    expect(range).toEqual({
      start: '2024-02-01',
      end: '2024-02-29',
    });
  });

  it('falls back to current month boundaries when no range is provided', () => {
    const range = generateDateRange();
    expect(range.start).toMatch(/^\d{4}-\d{2}-01$/);
    expect(range.end).toMatch(/^\d{4}-\d{2}-01$/);
  });
});

describe('PerformanceCalendar onRangeChange', () => {
  let unmountComponent: (() => void) | null = null;

  beforeEach(() => {
    resultsQuerySpy.mockClear();
    calendarComponentSpy.mockClear();
  });

  afterEach(() => {
    if (unmountComponent) {
      unmountComponent();
      unmountComponent = null;
    }
  });

  it('updates query date range when onRangeChange is triggered with week view dates', async () => {
    const { unmount } = render(<PerformanceCalendar />);
    unmountComponent = unmount;

    const latestProps = calendarComponentSpy.mock.calls[0][0];
    expect(typeof latestProps.onRangeChange).toBe('function');

    const weekDates = [
      new Date('2024-01-28T12:00:00Z'),
      new Date('2024-01-29T12:00:00Z'),
      new Date('2024-01-30T12:00:00Z'),
      new Date('2024-01-31T12:00:00Z'),
      new Date('2024-02-01T12:00:00Z'),
      new Date('2024-02-02T12:00:00Z'),
      new Date('2024-02-03T12:00:00Z'),
    ];

    await act(async () => {
      latestProps.onRangeChange(weekDates);
    });

    expect(resultsQuerySpy).toHaveBeenLastCalledWith(
      expect.objectContaining({
        from: '2024-01-28',
        to: '2024-02-03',
      }),
    );
  });

  it('updates query date range when onRangeChange is triggered with a single day view date', async () => {
    const { unmount } = render(<PerformanceCalendar />);
    unmountComponent = unmount;

    const latestProps = calendarComponentSpy.mock.calls[0][0];

    await act(async () => {
      latestProps.onRangeChange([new Date('2024-05-10T10:00:00Z')]);
    });

    expect(resultsQuerySpy).toHaveBeenLastCalledWith(
      expect.objectContaining({
        from: '2024-05-10',
        to: '2024-05-10',
      }),
    );
  });
});
