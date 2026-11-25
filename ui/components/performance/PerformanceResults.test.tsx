import React from 'react';
import { act, render } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

const { useGetPerformanceProfileResultsQuery, tableProps } = vi.hoisted(() => ({
  useGetPerformanceProfileResultsQuery: vi.fn(
    (_args: Record<string, unknown>, _opts?: unknown) => ({
      data: undefined,
      isFetching: false,
      isError: false,
      error: undefined,
    }),
  ),
  tableProps: { current: null as any },
}));

vi.mock('@meshery/schemas/mesheryApi', () => ({ useGetPerformanceProfileResultsQuery }));
vi.mock('@sistent/mui-datatables', () => ({
  default: (props: any) => {
    tableProps.current = props;
    return <div data-testid="results-table" />;
  },
}));
vi.mock('react-redux', () => ({
  useDispatch: () => vi.fn(),
  useSelector: (selector: any) =>
    selector({ ui: { providerCapabilities: {} }, prefTest: { results_selection: {} } }),
}));
vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));
vi.mock('@/store/slices/mesheryUi', () => ({ updateProgressAction: vi.fn() }));
vi.mock('@/store/slices/prefTest', () => ({ updateResultsSelection: vi.fn() }));
vi.mock('../shared/Modal/GenericModal', () => ({ default: () => null }));
vi.mock('../MesheryChart', () => ({ default: () => null }));
vi.mock('../CustomToolbarSelect', () => ({ default: () => null }));
vi.mock('./NodeDetails', () => ({ default: () => null }));

import MesheryResults from './PerformanceResults';

const lastQueryArgs = () => useGetPerformanceProfileResultsQuery.mock.calls.at(-1)![0];

describe('PerformanceResults', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    useGetPerformanceProfileResultsQuery.mockClear();
    tableProps.current = null;
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('goes back to the first page when searching from a later page', () => {
    render(<MesheryResults endpoint="/api/user/performance/profiles/profile-1/results" />);
    expect(lastQueryArgs()).toMatchObject({ performanceProfileId: 'profile-1', page: '0' });

    act(() => tableProps.current.options.onTableChange('changePage', { page: 2 }));
    expect(lastQueryArgs()).toMatchObject({ page: '2', search: '' });

    act(() => {
      tableProps.current.options.onTableChange('search', { searchText: 'istio' });
      vi.advanceTimersByTime(500);
    });
    expect(lastQueryArgs()).toMatchObject({ page: '0', search: 'istio' });
  });
});
