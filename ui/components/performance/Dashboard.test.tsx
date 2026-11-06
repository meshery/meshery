import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Dashboard is a container page; only its access gate is unit tested here.
// Children, routing, redux and the two RTK queries are mocked so the test can
// assert what the page renders and which requests it would issue.
const hoisted = vi.hoisted(() => ({
  canView: true,
  profilesQuery: vi.fn(),
  resultsQuery: vi.fn(),
}));

vi.mock('@sistent/sistent', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@sistent/sistent')>()),
  useHasPermission: () => hoisted.canView,
}));
vi.mock('@/rtk-query/performance-profile', () => ({
  useGetPerformanceProfilesQuery: (...args: unknown[]) => hoisted.profilesQuery(...args),
}));
vi.mock('@meshery/schemas/mesheryApi', () => ({
  useGetPerformanceResultsQuery: (...args: unknown[]) => hoisted.resultsQuery(...args),
}));
vi.mock('@/components/general/error-404/index', () => ({
  default: () => <div data-testid="default-error" />,
}));
vi.mock('./PerformanceCalendar', () => ({ default: () => null }));
vi.mock('./index', () => ({ default: () => null }));
vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));
vi.mock('next/router', () => ({ useRouter: () => ({ push: vi.fn() }) }));
vi.mock('react-redux', () => ({ useDispatch: () => vi.fn() }));
vi.mock('@/store/slices/mesheryUi', () => ({ updateProgressAction: vi.fn() }));

import Dashboard from './Dashboard';

const queryResult = { data: undefined, isFetching: false, isError: false, error: undefined };

describe('Performance Dashboard access gate', () => {
  beforeEach(() => {
    hoisted.profilesQuery.mockReset().mockReturnValue(queryResult);
    hoisted.resultsQuery.mockReset().mockReturnValue(queryResult);
  });

  it('shows DefaultError and issues no requests without the view permission', () => {
    hoisted.canView = false;

    render(<Dashboard />);

    expect(screen.getByTestId('default-error')).toBeInTheDocument();
    expect(hoisted.profilesQuery).toHaveBeenCalledWith(expect.anything(), { skip: true });
    expect(hoisted.resultsQuery).toHaveBeenCalledWith(expect.anything(), { skip: true });
  });

  it('fetches profiles and results with the view permission', () => {
    hoisted.canView = true;

    render(<Dashboard />);

    expect(screen.queryByTestId('default-error')).not.toBeInTheDocument();
    expect(hoisted.profilesQuery).toHaveBeenCalledWith(expect.anything(), { skip: false });
    expect(hoisted.resultsQuery).toHaveBeenCalledWith(expect.anything(), { skip: false });
  });
});
