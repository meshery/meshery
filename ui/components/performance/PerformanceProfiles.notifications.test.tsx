import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Covers the notifications PerformanceProfiles raises for profile actions.
// Its children (grid, table, modals, the performance wizard) are stubbed; the
// grid stub exposes a delete button wired to the real deleteHandler.

// Query results must be referentially stable across renders: the component
// copies them into state inside an effect keyed on the data object.
const { notify, deleteProfile, refetch, PROFILES_QUERY_RESULT } = vi.hoisted(() => {
  const refetch = vi.fn();
  return {
    notify: vi.fn(),
    deleteProfile: vi.fn(),
    refetch,
    PROFILES_QUERY_RESULT: {
      data: { profiles: [{ id: 'p1', name: 'profile-1' }], totalCount: 1, pageSize: 10 },
      isFetching: false,
      isError: false,
      error: undefined,
      refetch,
    },
  };
});

vi.mock('@/rtk-query/performance-profile', () => ({
  useDeletePerformanceProfileMutation: () => [deleteProfile],
  useGetPerformanceProfilesQuery: () => PROFILES_QUERY_RESULT,
}));

vi.mock('react-redux', () => ({
  useDispatch: () => vi.fn(),
  useSelector: (selector: (state: unknown) => unknown) =>
    selector({ ui: { providerCapabilities: {} } }),
}));

vi.mock('@/store/slices/mesheryUi', () => ({ updateProgressAction: vi.fn() }));
vi.mock('@/utils/hooks/useNotification', () => ({ useNotification: () => ({ notify }) }));
vi.mock('@/utils/dimension', () => ({ useWindowDimensions: () => ({ width: 1024 }) }));
vi.mock('@/utils/responsive-column', () => ({ updateVisibleColumns: () => ({}) }));
vi.mock('@/utils/provider', () => ({ isLocalProvider: () => false }));
vi.mock('@/utils/utils', () => ({ ConditionalTooltip: ({ children }: any) => <>{children}</> }));
vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({ id: 'id', function: 'fn' }) }),
}));

vi.mock('./index', () => ({ default: () => null }));
vi.mock('./PerformanceResults', () => ({ default: () => null }));
// The delete confirmation prompt always answers DELETE.
vi.mock('../general/PromptComponent', () => ({
  default: React.forwardRef(function Prompt(_props, ref) {
    React.useImperativeHandle(ref, () => ({ show: () => Promise.resolve('DELETE') }));
    return null;
  }),
}));
vi.mock('../general/ViewSwitch', () => ({
  default: ({ changeView }: any) => (
    <button data-testid="table-view" onClick={() => changeView('table')}>
      table
    </button>
  ),
}));
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
  default: ({ deleteHandler }: any) => (
    <button data-testid="delete-p1" onClick={() => deleteHandler('p1')}>
      delete
    </button>
  ),
}));

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@sistent/sistent')>();
  return {
    ...actual,
    // Render the toolbar the table shows when the first row is selected.
    ResponsiveDataTable: ({ options }: any) => (
      <div data-testid="selected-toolbar">
        {options.customToolbarSelect({ data: [{ index: 0 }] }, [], vi.fn())}
      </div>
    ),
    CustomTooltip: ({ children }: any) => <>{children}</>,
    IconButton: ({ children, onClick }: any) => (
      <button data-testid="bulk-delete" onClick={onClick}>
        {children}
      </button>
    ),
    CustomColumnVisibilityControl: () => null,
    SearchBar: () => null,
    useHasPermission: () => true,
  };
});

import PerformanceProfiles from './PerformanceProfiles';

const notifiedEventTypes = () =>
  notify.mock.calls.map(([arg]) => arg?.event_type?.type ?? arg?.event_type);

describe('PerformanceProfiles delete notifications', () => {
  beforeEach(() => {
    notify.mockReset();
    deleteProfile.mockReset();
    refetch.mockReset();
  });

  it('shows an error notification when the delete fails', async () => {
    deleteProfile.mockReturnValue({
      unwrap: () => Promise.reject({ status: 500, data: 'profile is in use' }),
    });
    const user = userEvent.setup();
    render(<PerformanceProfiles />);

    await user.click(screen.getByTestId('delete-p1'));

    await waitFor(() => expect(notify).toHaveBeenCalledTimes(1));
    expect(deleteProfile).toHaveBeenCalledWith({ id: 'p1' });
    expect(notifiedEventTypes()).toEqual(['error']);
    expect(notify.mock.calls[0][0].message).toContain('profile is in use');
    expect(refetch).not.toHaveBeenCalled();
  });

  it('shows a success notification and refetches when the delete succeeds', async () => {
    deleteProfile.mockReturnValue({ unwrap: () => Promise.resolve({}) });
    const user = userEvent.setup();
    render(<PerformanceProfiles />);

    await user.click(screen.getByTestId('delete-p1'));

    await waitFor(() => expect(notify).toHaveBeenCalledTimes(1));
    expect(notifiedEventTypes()).toEqual(['success']);
    expect(refetch).toHaveBeenCalledTimes(1);
  });

  it('includes the server reason when a bulk delete from the table fails', async () => {
    deleteProfile.mockReturnValue({
      unwrap: () => Promise.reject({ status: 500, data: 'profile is in use' }),
    });
    const user = userEvent.setup();
    render(<PerformanceProfiles />);

    await user.click(screen.getByTestId('table-view'));
    await user.click(screen.getByTestId('bulk-delete'));

    await waitFor(() => expect(notify).toHaveBeenCalledTimes(1));
    expect(deleteProfile).toHaveBeenCalledWith({ id: 'p1' });
    expect(notifiedEventTypes()).toEqual(['error']);
    expect(notify.mock.calls[0][0].message).toBe('Failed to delete 1 Profile: profile is in use');
  });
});
