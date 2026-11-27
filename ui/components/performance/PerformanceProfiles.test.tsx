import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// PerformanceProfiles is the table+grid view manager for performance
// profiles. Its children (grid, table, modals, the performance wizard) are
// stubbed here; the grid itself is unit-tested in
// PerformanceProfileGrid.test.tsx and the full page is covered by the
// playwright performance suite. These tests cover the page's own state:
// the page and search values it sends to the profiles query.

// Query results must be referentially stable across renders: the component
// copies them into state inside an effect keyed on the data object.
const { getProfilesQuery, canSpy, PROFILES_QUERY_RESULT } = vi.hoisted(() => ({
  getProfilesQuery: vi.fn(),
  canSpy: vi.fn(() => true),
  PROFILES_QUERY_RESULT: {
    data: {
      profiles: [{ id: 'p1', name: 'profile-1' }],
      totalCount: 40,
      pageSize: 10,
    },
    isFetching: false,
    isError: false,
    error: undefined,
    refetch: () => {},
  },
}));

vi.mock('@/rtk-query/performance-profile', () => ({
  useDeletePerformanceProfileMutation: () => [vi.fn()],
  useGetPerformanceProfilesQuery: (args: unknown, options: unknown) => {
    getProfilesQuery(args, options);
    return PROFILES_QUERY_RESULT;
  },
}));

vi.mock('react-redux', () => ({
  useDispatch: () => vi.fn(),
  useSelector: (selector: (state: unknown) => unknown) =>
    selector({ ui: { providerCapabilities: {} } }),
}));

vi.mock('@/store/slices/mesheryUi', () => ({ updateProgressAction: vi.fn() }));
vi.mock('@/utils/hooks/useNotification', () => ({ useNotification: () => ({ notify: vi.fn() }) }));
vi.mock('@/utils/dimension', () => ({ useWindowDimensions: () => ({ width: 1024 }) }));
vi.mock('@/utils/responsive-column', () => ({ updateVisibleColumns: () => ({}) }));
vi.mock('@/utils/provider', () => ({ isLocalProvider: () => false }));
vi.mock('@/utils/utils', () => ({ ConditionalTooltip: ({ children }: any) => <>{children}</> }));
vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({ id: 'id', function: 'fn' }) }),
}));

vi.mock('./index', () => ({ default: () => null }));
vi.mock('./PerformanceResults', () => ({ default: () => null }));
vi.mock('../general/PromptComponent', () => ({ default: React.forwardRef(() => null) }));
vi.mock('../general/ViewSwitch', () => ({ default: () => null }));
vi.mock('../general/error-404', () => ({ default: () => <div>no access</div> }));
vi.mock('../connections/common', () => ({
  DefaultTableCell: () => null,
  SortableTableCell: () => null,
}));
vi.mock('./style', () => ({
  ButtonTextWrapper: ({ children }: any) => <span>{children}</span>,
  ProfileContainer: ({ children }: any) => <div>{children}</div>,
  ViewSwitchBUtton: ({ children }: any) => <div>{children}</div>,
}));
vi.mock('@/assets/styles/general/tool.styles', () => ({
  ToolWrapper: ({ children }: any) => <div>{children}</div>,
}));
vi.mock('./PerformanceProfileGrid', () => ({
  default: ({ page, setPage }: any) => (
    <div>
      <span data-testid="grid-page">{page}</span>
      <button data-testid="go-to-page-3" onClick={() => setPage(2)}>
        page 3
      </button>
    </div>
  ),
}));

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@sistent/sistent')>();
  return {
    ...actual,
    ResponsiveDataTable: () => null,
    CustomColumnVisibilityControl: () => null,
    SearchBar: ({ onSearch }: any) => (
      <button data-testid="search-load" onClick={() => onSearch('load')}>
        search
      </button>
    ),
    useHasPermission: (key: any) => canSpy(key),
  };
});

import PerformanceProfiles from './PerformanceProfiles';

const lastQueryArgs = () => getProfilesQuery.mock.calls.at(-1)?.[0];

describe('PerformanceProfiles', () => {
  beforeEach(() => {
    canSpy.mockReturnValue(true);
    getProfilesQuery.mockClear();
  });

  it('goes back to the first page when the search changes', async () => {
    const user = userEvent.setup();
    render(<PerformanceProfiles />);

    await user.click(screen.getByTestId('go-to-page-3'));
    expect(lastQueryArgs()).toMatchObject({ page: 2, search: '' });
    expect(screen.getByTestId('grid-page')).toHaveTextContent('2');

    await user.click(screen.getByTestId('search-load'));
    expect(lastQueryArgs()).toMatchObject({ page: 0, search: 'load' });
    expect(screen.getByTestId('grid-page')).toHaveTextContent('0');
  });

  it('renders 404 fallback and skips query when view permission is denied', () => {
    canSpy.mockReturnValue(false);
    render(<PerformanceProfiles />);

    expect(screen.getByText('no access')).toBeInTheDocument();
    expect(getProfilesQuery).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ skip: true }),
    );
  });
});
