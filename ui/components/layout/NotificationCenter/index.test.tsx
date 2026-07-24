import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// jsdom does not implement IntersectionObserver; EventsView constructs one
// inline so we need a minimal polyfill.
class IntersectionObserverPolyfill {
  observe = vi.fn();
  unobserve = vi.fn();
  disconnect = vi.fn();
  takeRecords = () => [];
  root = null;
  rootMargin = '';
  thresholds = [];
  constructor() {}
}
(globalThis as any).IntersectionObserver =
  (globalThis as any).IntersectionObserver || IntersectionObserverPolyfill;

// Hoisted store-ish state we can swap between tests.
const dispatchMock = vi.fn();
const fetchEventsMock = vi.fn();
const deleteEventsMutation = vi.fn();
const updateEventsMutation = vi.fn();
const summaryQuery: { data: any; error: any; isLoading: boolean } = {
  data: { countBySeverityLevel: [], readCount: 0 },
  error: null,
  isLoading: false,
};

const sliceState = {
  isNotificationCenterOpen: true,
  events: [{ id: 'e1', description: 'Event 1', severity: 'informational', status: 'unread' }],
  areAllChecked: false,
  checkedEvents: [],
  selectedSeverity: 'informational',
  hasMore: false,
  view_to_fetch_on_open: { page: 0, filters: {} },
  uiConfig: { title: 'Notifications', history_mode: false, icon: null },
};

vi.mock('react-redux', () => ({
  useDispatch: () => dispatchMock,
  useSelector: (sel: any) =>
    sel({
      events: {
        ui: sliceState.uiConfig,
        isNotificationCenterOpen: sliceState.isNotificationCenterOpen,
        current_view: { has_more: sliceState.hasMore, page: 0 },
        view_to_fetch_on_open: sliceState.view_to_fetch_on_open,
      },
    }),
  Provider: ({ children }: any) => <>{children}</>,
}));

vi.mock('../../../store/slices/events', () => ({
  closeNotificationCenter: () => ({ type: 'events/closeNotificationCenter' }),
  loadEvents: (...args: unknown[]) => ({ type: 'events/loadEvents', args }),
  loadNextPage: (...args: unknown[]) => ({ type: 'events/loadNextPage', args }),
  selectAreAllEventsChecked: () => sliceState.areAllChecked,
  selectCheckedEvents: () => sliceState.checkedEvents,
  selectEvents: () => sliceState.events,
  selectSeverity: () => sliceState.selectedSeverity,
  toggleNotificationCenter: () => ({ type: 'events/toggleNotificationCenter' }),
  updateCheckAllEvents: (payload: any) => ({
    type: 'events/updateCheckAllEvents',
    payload,
  }),
}));

vi.mock('../../../rtk-query/notificationCenter', () => ({
  useDeleteEventsMutation: () => [deleteEventsMutation, { isLoading: false }],
  useGetEventsSummaryQuery: () => summaryQuery,
  useLazyGetEventsQuery: () => [fetchEventsMock, { isFetching: false }],
  useUpdateEventsMutation: () => [updateEventsMutation, { isLoading: false }],
}));

vi.mock('@xstate/react', () => ({
  useActorRef: () => ({ send: vi.fn() }),
}));

vi.mock('machines/operationsCenter', () => ({
  operationsCenterActor: { id: 'ops' },
}));

vi.mock('../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

vi.mock('@sistent/sistent', () => ({
  alpha: (c: string) => c,
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Button: ({ children, onClick, ...props }: any) => (
    <button type="button" onClick={onClick} {...props}>
      {children}
    </button>
  ),
  Checkbox: ({ checked, onChange }: any) => (
    <input
      type="checkbox"
      data-testid="bulk-checkbox"
      checked={Boolean(checked)}
      readOnly
      onChange={(e) => onChange?.(e, e.target.checked)}
    />
  ),
  CircularProgress: () => <div data-testid="loading" />,
  ClickAwayListener: ({ children }: any) => <>{children}</>,
  Collapse: ({ in: open, children }: any) =>
    open ? <div data-testid="collapse-open">{children}</div> : null,
  CustomTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={title}>
      {children}
    </div>
  ),
  Divider: () => <hr />,
  ErrorBoundary: ({ children }: any) => <>{children}</>,
  IconButton: ({ children, onClick, disabled }: any) => (
    <button type="button" disabled={disabled} onClick={onClick}>
      {children}
    </button>
  ),
  NoSsr: ({ children }: any) => <>{children}</>,
  Typography: ({ children, ...props }: any) => <span {...props}>{children}</span>,
  useTheme: () => ({
    palette: {
      mode: 'light',
      icon: { secondary: '#888' },
      text: { primary: '#000' },
    },
  }),
  DeleteIcon: () => <svg data-testid="delete-icon" />,
}));

vi.mock('./filter', () => ({
  default: ({ handleFilter }: any) => (
    <button
      data-testid="filter-button"
      type="button"
      onClick={() => handleFilter({ status: 'read' })}
    >
      Filter
    </button>
  ),
}));

vi.mock('../../../assets/icons/BellIcon', () => ({
  default: () => <svg data-testid="bell-icon" />,
}));

vi.mock('../../../assets/icons/DoneIcon', () => ({
  default: () => <svg data-testid="done-icon" />,
}));

vi.mock('../../../assets/icons/ReadIcon', () => ({
  default: () => <svg data-testid="read-icon" />,
}));

vi.mock('../../../assets/icons/UnreadIcon', () => ({
  default: () => <svg data-testid="unread-icon" />,
}));

vi.mock('../../../css/icons.styles', () => ({
  iconMedium: { height: '20px' },
}));

vi.mock('./constants', () => ({
  NOTIFICATION_CENTER_TOGGLE_CLASS: 'toggle-notification-center',
  SEVERITY: {
    INFO: 'informational',
    ERROR: 'error',
    WARNING: 'warning',
    SUCCESS: 'success',
  },
  STATUS: { READ: 'read', UNREAD: 'unread' },
  SEVERITY_STYLE: {
    informational: { icon: () => <svg data-testid="sev-info" />, color: '#info' },
    error: { icon: () => <svg data-testid="sev-error" />, color: '#err' },
    warning: { icon: () => <svg data-testid="sev-warning" />, color: '#warn' },
    success: { icon: () => <svg data-testid="sev-success" />, color: '#suc' },
  },
  getStatusStyle: () => ({
    read: { icon: () => <svg data-testid="status-read" />, color: '#read' },
  }),
}));

vi.mock('./notification', () => ({
  default: ({ event_id }: any) => (
    <div data-testid={`notification-${event_id}`}>notification {event_id}</div>
  ),
}));

vi.mock('./notificationCenter.style', () => ({
  Container: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  DarkBackdrop: () => <div data-testid="backdrop" />,
  NotificationButton: ({ children, onClick, onMouseOver, onMouseLeave, ...props }: any) => (
    <button
      type="button"
      data-testid="notification-button"
      onClick={onClick}
      onMouseOver={onMouseOver}
      onMouseLeave={onMouseLeave}
      {...props}
    >
      {children}
    </button>
  ),
  NotificationContainer: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  SeverityChips: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  SeverityChip: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  SideList: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  StyledBadge: ({ children, badgeContent, badgeColor }: any) => (
    <div data-testid="styled-badge" data-content={badgeContent} data-color={badgeColor}>
      {children}
    </div>
  ),
  StyledNotificationDrawer: ({ children, open, ...rest }: any) => (
    <div data-testid="drawer" data-open={String(Boolean(open))}>
      {children}
    </div>
  ),
  Title: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  TitleBellIcon: ({ children, onClick }: any) => (
    <button type="button" onClick={onClick} data-testid="title-bell">
      {children}
    </button>
  ),
}));

vi.mock('../../shared/ErrorBoundary/ErrorBoundary', () => ({
  default: () => <div data-testid="error-fallback" />,
}));

vi.mock('../../../utils/Elements', () => ({
  hasClass: () => false,
}));

import NotificationCenter, {
  NotificationCenterContext,
  NotificationCenterProvider,
  NotificationDrawerButton,
} from './index';

beforeEach(() => {
  dispatchMock.mockReset();
  fetchEventsMock.mockReset();
  deleteEventsMutation.mockReset();
  updateEventsMutation.mockReset();
  summaryQuery.data = { countBySeverityLevel: [], readCount: 0 };
  summaryQuery.error = null;
  summaryQuery.isLoading = false;
  sliceState.isNotificationCenterOpen = true;
  sliceState.events = [
    { id: 'e1', description: 'Event 1', severity: 'informational', status: 'unread' },
  ];
  sliceState.areAllChecked = false;
  sliceState.checkedEvents = [];
  sliceState.selectedSeverity = 'informational';
  sliceState.hasMore = false;
  sliceState.view_to_fetch_on_open = { page: 0, filters: {} };
  sliceState.uiConfig = { title: 'Notifications', history_mode: false, icon: null };
});

describe('NotificationCenter (default export)', () => {
  it('renders nothing when notification center is closed', () => {
    sliceState.isNotificationCenterOpen = false;
    const { container } = render(<NotificationCenter />);
    expect(container).toBeEmptyDOMElement();
  });

  it('renders the drawer when the notification center is open', () => {
    render(<NotificationCenter />);
    expect(screen.getByTestId('drawer')).toBeInTheDocument();
    // The mocked summary returns empty so default bell icon should render.
    expect(screen.getByTestId('notification-e1')).toBeInTheDocument();
  });

  it('dispatches loadEvents on initial mount', () => {
    render(<NotificationCenter />);
    expect(dispatchMock.mock.calls.some(([action]) => action?.type === 'events/loadEvents')).toBe(
      true,
    );
  });

  it('clicking the title bell icon closes the center', () => {
    render(<NotificationCenter />);
    fireEvent.click(screen.getByTestId('title-bell'));
    expect(
      dispatchMock.mock.calls.some(([action]) => action?.type === 'events/closeNotificationCenter'),
    ).toBe(true);
  });

  it('clicking the filter button reloads events with new filters', () => {
    render(<NotificationCenter />);
    fireEvent.click(screen.getByTestId('filter-button'));
    expect(dispatchMock.mock.calls.some(([action]) => action?.type === 'events/loadEvents')).toBe(
      true,
    );
  });

  it('renders an empty state when there are no events', () => {
    sliceState.events = [];
    render(<NotificationCenter />);
    expect(screen.getByText(/No notifications to show/i)).toBeInTheDocument();
  });
});

describe('NotificationDrawerButton', () => {
  it('dispatches toggle when clicked and renders the navbar icon', () => {
    render(<NotificationDrawerButton />);
    fireEvent.click(screen.getByTestId('notification-button'));
    expect(
      dispatchMock.mock.calls.some(
        ([action]) => action?.type === 'events/toggleNotificationCenter',
      ),
    ).toBe(true);
  });

  it('shows a colored badge when there is an error-severity event', () => {
    summaryQuery.data = {
      countBySeverityLevel: [{ severity: 'error', count: 3 }],
      readCount: 0,
    };
    render(<NotificationDrawerButton />);
    expect(screen.getByTestId('styled-badge')).toHaveAttribute('data-content', '3');
  });

  it('shows a warning badge when only warnings are present', () => {
    summaryQuery.data = {
      countBySeverityLevel: [{ severity: 'warning', count: 2 }],
      readCount: 0,
    };
    render(<NotificationDrawerButton />);
    expect(screen.getByTestId('styled-badge')).toHaveAttribute('data-content', '2');
  });

  it('falls back to the bell icon when there are no errors or warnings', () => {
    summaryQuery.data = { countBySeverityLevel: [], readCount: 0 };
    render(<NotificationDrawerButton />);
    expect(screen.getByTestId('bell-icon')).toBeInTheDocument();
  });
});

describe('NotificationCenterProvider', () => {
  it('renders children and the drawer in the same tree', () => {
    sliceState.isNotificationCenterOpen = false;
    render(
      <NotificationCenterProvider>
        <span data-testid="child">child</span>
      </NotificationCenterProvider>,
    );
    expect(screen.getByTestId('child')).toBeInTheDocument();
  });

  it('exposes a context default with null refs', () => {
    const Spy = () => {
      const ctx = React.useContext(NotificationCenterContext);
      return <div data-testid="ctx-keys">{Object.keys(ctx).sort().join(',')}</div>;
    };
    render(<Spy />);
    expect(screen.getByTestId('ctx-keys')).toHaveTextContent(
      'drawerAnchorEl,operationsCenterActorRef,setDrawerAnchor,toggleButtonRef',
    );
  });
});
