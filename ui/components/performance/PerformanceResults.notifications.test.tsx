import React from 'react';
import { render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const { queryState, notify } = vi.hoisted(() => ({
  queryState: { current: {} as Record<string, unknown> },
  notify: vi.fn(),
}));

vi.mock('@meshery/schemas/mesheryApi', () => ({
  useGetPerformanceProfileResultsQuery: () => queryState.current,
}));
vi.mock('@sistent/mui-datatables', () => ({ default: () => null }));
vi.mock('react-redux', () => ({
  useDispatch: () => vi.fn(),
  useSelector: (selector: any) =>
    selector({ ui: { providerCapabilities: {} }, prefTest: { results_selection: {} } }),
}));
vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));
vi.mock('@/store/slices/mesheryUi', () => ({ updateProgressAction: vi.fn() }));
vi.mock('@/store/slices/prefTest', () => ({ updateResultsSelection: vi.fn() }));
vi.mock('../shared/Modal/GenericModal', () => ({ default: () => null }));
vi.mock('../MesheryChart', () => ({ default: () => null }));
vi.mock('../CustomToolbarSelect', () => ({ default: () => null }));
vi.mock('./NodeDetails', () => ({ default: () => null }));

import MesheryResults from './PerformanceResults';

const endpoint = '/api/user/performance/profiles/profile-1/results';

describe('PerformanceResults error notification', () => {
  beforeEach(() => {
    notify.mockClear();
  });

  it('shows the server error instead of [object Object]', () => {
    queryState.current = {
      isFetching: false,
      isError: true,
      error: { status: 500, data: 'failed to fetch results: database is locked' },
    };

    render(<MesheryResults endpoint={endpoint} />);

    expect(notify).toHaveBeenCalledTimes(1);
    const { message } = notify.mock.calls[0][0];
    expect(message).toBe('failed to fetch results: database is locked');
    expect(message).not.toContain('[object Object]');
  });

  it('falls back to a readable message when the error has no details', () => {
    queryState.current = {
      isFetching: false,
      isError: true,
      error: { status: 'FETCH_ERROR' },
    };

    render(<MesheryResults endpoint={endpoint} />);

    expect(notify.mock.calls[0][0].message).toBe('There was an error fetching results');
  });
});
