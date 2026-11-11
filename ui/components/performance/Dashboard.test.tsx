import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const {
  getProfilesQuerySpy,
  getResultsQuerySpy,
  canSpy,
  PROFILES_QUERY_RESULT,
  RESULTS_QUERY_RESULT,
} = vi.hoisted(() => ({
  getProfilesQuerySpy: vi.fn(),
  getResultsQuerySpy: vi.fn(),
  canSpy: vi.fn(() => true),
  PROFILES_QUERY_RESULT: {
    data: { profiles: [], totalCount: 0 },
    isFetching: false,
    isError: false,
    error: undefined,
  },
  RESULTS_QUERY_RESULT: {
    data: { results: [], total_count: 0 },
    isFetching: false,
    isError: false,
    error: undefined,
  },
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock('react-redux', () => ({
  useDispatch: () => vi.fn(),
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgressAction: vi.fn(),
}));

vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

vi.mock('@/rtk-query/performance-profile', () => ({
  useGetPerformanceProfilesQuery: (args: unknown, options: unknown) => {
    getProfilesQuerySpy(args, options);
    return PROFILES_QUERY_RESULT;
  },
}));

vi.mock('@meshery/schemas/mesheryApi', () => ({
  useGetPerformanceResultsQuery: (args: unknown, options: unknown) => {
    getResultsQuerySpy(args, options);
    return RESULTS_QUERY_RESULT;
  },
}));

vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({ id: 'id', function: 'fn' }) }),
}));

vi.mock('@/components/general/error-404/index', () => ({
  default: () => <div data-testid="default-error">no access</div>,
}));

vi.mock('./PerformanceCalendar', () => ({
  default: () => <div data-testid="performance-calendar" />,
}));

vi.mock('./index', () => ({
  default: () => <div data-testid="performance-wizard" />,
}));

vi.mock('@sistent/sistent', () => ({
  useHasPermission: (key: any) => canSpy(key),
  Modal: ({ children, open }: any) => (open ? <div>{children}</div> : null),
  Button: ({ children, onClick, ...rest }: any) => (
    <button onClick={onClick} {...rest}>
      {children}
    </button>
  ),
  Grid2: ({ children }: any) => <div>{children}</div>,
  Paper: ({ children }: any) => <div>{children}</div>,
  Typography: ({ children }: any) => <span>{children}</span>,
  useTheme: () => ({
    palette: { border: { normal: '#ccc' } },
    breakpoints: { down: () => '@media (max-width: 600px)' },
  }),
  styled: (Component: any) => () => {
    return ({ children, ...rest }: any) =>
      typeof Component === 'string'
        ? React.createElement(Component, rest, children)
        : React.createElement(Component, rest, children);
  },
}));

import Dashboard from './Dashboard';

describe('Performance Dashboard', () => {
  beforeEach(() => {
    canSpy.mockReturnValue(true);
    getProfilesQuerySpy.mockClear();
    getResultsQuerySpy.mockClear();
  });

  it('renders the 404 fallback and skips both profile and result queries when view permission is denied', () => {
    canSpy.mockReturnValue(false);
    render(<Dashboard />);

    expect(screen.getByTestId('default-error')).toBeInTheDocument();
    expect(screen.queryByTestId('performance-calendar')).not.toBeInTheDocument();

    expect(getProfilesQuerySpy).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ skip: true }),
    );
    expect(getResultsQuerySpy).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ skip: true }),
    );
  });

  it('renders dashboard content and executes queries without skip when view permission is granted', () => {
    canSpy.mockReturnValue(true);
    render(<Dashboard />);

    expect(screen.getByTestId('performance-calendar')).toBeInTheDocument();
    expect(screen.queryByTestId('default-error')).not.toBeInTheDocument();

    expect(getProfilesQuerySpy).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ skip: false }),
    );
    expect(getResultsQuerySpy).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ skip: false }),
    );
  });
});
