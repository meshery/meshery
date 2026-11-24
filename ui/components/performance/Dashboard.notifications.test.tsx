import React from 'react';
import { render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const { notify, profilesState, resultsState } = vi.hoisted(() => ({
  notify: vi.fn(),
  profilesState: { current: {} as Record<string, unknown> },
  resultsState: { current: {} as Record<string, unknown> },
}));

vi.mock('next/router', () => ({ useRouter: () => ({ push: vi.fn() }) }));
vi.mock('react-redux', () => ({ useDispatch: () => vi.fn() }));
vi.mock('@/store/slices/mesheryUi', () => ({ updateProgressAction: vi.fn() }));
vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));
vi.mock('@/rtk-query/performance-profile', () => ({
  useGetPerformanceProfilesQuery: () => profilesState.current,
}));
vi.mock('@meshery/schemas/mesheryApi', () => ({
  useGetPerformanceResultsQuery: () => resultsState.current,
}));
vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({ id: 'id', function: 'fn' }) }),
}));
vi.mock('@/components/general/error-404/index', () => ({ default: () => null }));
vi.mock('./PerformanceCalendar', () => ({ default: () => null }));
vi.mock('./index', () => ({ default: () => null }));
vi.mock('@sistent/sistent', () => ({
  useHasPermission: () => true,
  Modal: ({ children, open }: any) => (open ? <div>{children}</div> : null),
  Button: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  Paper: ({ children }: any) => <div>{children}</div>,
  Typography: ({ children }: any) => <span>{children}</span>,
  useTheme: () => ({
    palette: { border: { normal: '#ccc' }, text: { default: '#000' }, icon: { default: '#000' } },
    breakpoints: { down: () => '@media (max-width: 600px)' },
  }),
  styled: (Component: any) => () => {
    return ({ children, ...rest }: any) => React.createElement(Component, rest, children);
  },
}));

import Dashboard from './Dashboard';

const ok = { data: undefined, isFetching: false, isError: false, error: undefined };

describe('Performance Dashboard error notifications', () => {
  beforeEach(() => {
    notify.mockClear();
    profilesState.current = ok;
    resultsState.current = ok;
  });

  it('shows the server error when profiles fail to load, not [object Object]', () => {
    profilesState.current = {
      ...ok,
      isError: true,
      error: { status: 500, data: 'failed to fetch profiles: database is locked' },
    };

    render(<Dashboard />);

    expect(notify).toHaveBeenCalledTimes(1);
    const { message } = notify.mock.calls[0][0];
    expect(message).toBe('failed to fetch profiles: database is locked');
    expect(message).not.toContain('[object Object]');
  });

  it('falls back to a readable message when the results error has no details', () => {
    resultsState.current = { ...ok, isError: true, error: { status: 'FETCH_ERROR' } };

    render(<Dashboard />);

    expect(notify).toHaveBeenCalledTimes(1);
    const { message } = notify.mock.calls[0][0];
    expect(message).toBe('Failed to Fetch Results');
  });
});
