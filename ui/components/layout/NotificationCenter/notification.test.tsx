import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const dispatchMock = vi.fn();
const updateStatusMutation = vi.fn();
const deleteEventMutation = vi.fn();

const selectors = {
  event: {
    id: 'evt-1',
    description: 'Short description',
    severity: 'informational',
    status: 'unread',
    metadata: {},
    checked: false,
    userID: 'u1',
    systemID: 's1',
    createdAt: '2024-01-01T00:00:00.000Z',
  } as any,
  visible: true,
  ui: { history_mode: false },
};

vi.mock('react-redux', () => ({
  useSelector: (sel: any) =>
    sel({
      events: { ui: selectors.ui },
      // The selector helpers below ignore the slice shape.
    }),
  useDispatch: () => dispatchMock,
}));

vi.mock('../../../store/slices/events', () => ({
  selectEventById: () => selectors.event,
  selectIsEventVisible: () => selectors.visible,
  updateIsEventChecked: (payload: any) => ({ type: 'events/updateIsEventChecked', payload }),
}));

vi.mock('../../../rtk-query/notificationCenter', () => ({
  useUpdateStatusMutation: () => [updateStatusMutation],
  useDeleteEventMutation: () => [deleteEventMutation],
}));

vi.mock('../../../rtk-query/user', () => ({
  useGetUserByIdQuery: () => ({
    data: { firstName: 'Alice', lastName: 'Doe', avatarUrl: '/a.png' },
  }),
}));

vi.mock('@sistent/sistent', () => {
  return {
    alpha: (color: string) => color,
    Avatar: ({ alt, src }: any) => <img data-testid="avatar" alt={alt} src={src} />,
    Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    Checkbox: ({ checked, onClick, onChange }: any) => (
      <input
        type="checkbox"
        data-testid="event-checkbox"
        checked={Boolean(checked)}
        readOnly
        onClick={onClick}
        onChange={(e) => onChange?.(e, e.target.checked)}
      />
    ),
    Collapse: ({ in: isIn, children }: any) =>
      isIn ? <div data-testid="collapse-open">{children}</div> : null,
    CustomTooltip: ({ children, title }: any) => (
      <div data-testid="tooltip" data-title={title}>
        {children}
      </div>
    ),
    FormattedTime: ({ date }: any) => <span data-testid="formatted-time">{String(date)}</span>,
    IconButton: ({ children, onClick, ...props }: any) => (
      <button type="button" onClick={onClick} {...props}>
        {children}
      </button>
    ),
    MoreVertIcon: () => <svg data-testid="more-vert" />,
    Popover: ({ open, children }: any) =>
      open ? <div data-testid="popover">{children}</div> : null,
    Slide: ({ children, in: isIn }: any) =>
      isIn !== false ? <div data-testid="slide">{children}</div> : null,
    Typography: ({ children, ...props }: any) => <span {...props}>{children}</span>,
    useTheme: () => ({
      palette: {
        icon: { secondary: '#888' },
        text: { default: '#000' },
      },
    }),
    DeleteIcon: () => <svg data-testid="delete-icon" />,
  };
});

vi.mock('./notificationCenter.style', () => ({
  OptionList: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  OptionListItem: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  MenuPaper: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  SocialListItem: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  ListButton: ({ children, onClick, ...props }: any) => (
    <button type="button" onClick={onClick} {...props}>
      {children}
    </button>
  ),
  ActorAvatar: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Expanded: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  GridItem: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Message: ({ children, ...props }: any) => <span {...props}>{children}</span>,
  StyledAvatarStack: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Root: ({ children, ...props }: any) => (
    <div data-testid="notification-root" {...props}>
      {children}
    </div>
  ),
  Summary: ({ children, onClick, ...props }: any) => (
    <div data-testid="notification-summary" onClick={onClick} {...props}>
      {children}
    </div>
  ),
}));

vi.mock('./constants', () => ({
  SEVERITY: {
    INFO: 'informational',
    ERROR: 'error',
    WARNING: 'warning',
    SUCCESS: 'success',
  },
  STATUS: { READ: 'read', UNREAD: 'unread' },
  SEVERITY_STYLE: {
    informational: {
      icon: ({ fill }: any) => <svg data-testid="severity-icon" data-fill={fill} />,
      color: '#info',
    },
  },
}));

vi.mock('../../../css/icons.styles', () => ({
  iconLarge: {},
  iconMedium: {},
}));

vi.mock('../../../assets/icons/FacebookIcon', () => ({
  default: () => <svg data-testid="facebook-icon" />,
}));
vi.mock('../../../assets/icons/LinkedInIcon', () => ({
  default: () => <svg data-testid="linkedin-icon" />,
}));
vi.mock('../../../assets/icons/TwitterIcon', () => ({
  default: () => <svg data-testid="twitter-icon" />,
}));
vi.mock('../../../assets/icons/ShareIcon', () => ({
  default: () => <svg data-testid="share-icon" />,
}));

vi.mock('@/assets/icons/ErrorIcon', () => ({ default: () => <svg data-testid="error-icon" /> }));
vi.mock('../../../assets/icons/ReadIcon', () => ({
  default: () => <svg data-testid="read-icon" />,
}));
vi.mock('../../../assets/icons/UnreadIcon', () => ({
  default: () => <svg data-testid="unread-icon" />,
}));

vi.mock('./metadata', () => ({
  FormattedLinkMetadata: () => <div data-testid="link-meta" />,
  FormattedMetadata: () => <div data-testid="formatted-meta" />,
  PropertyLinkFormatters: {
    doc: (value: string) => ({ label: 'Doc', href: value }),
  },
}));

vi.mock('react-share', () => ({
  FacebookShareButton: ({ children }: any) => <div data-testid="fb-share">{children}</div>,
  LinkedinShareButton: ({ children }: any) => <div data-testid="li-share">{children}</div>,
  TwitterShareButton: ({ children }: any) => <div data-testid="tw-share">{children}</div>,
}));

vi.mock('@/constants/endpoints', () => ({
  MESHERY_DOCS_URL: 'https://docs.meshery.io',
}));

import {
  canTruncateDescription,
  MAX_NOTIFICATION_DESCRIPTION_LENGTH,
  eventPreventDefault,
  eventstopPropagation,
  Notification,
  DeleteEvent,
  ChangeStatus,
  getErrorCodesFromEvent,
} from './notification';

describe('helpers', () => {
  it('canTruncateDescription returns true when over max length', () => {
    expect(canTruncateDescription('x')).toBe(false);
    expect(canTruncateDescription('x'.repeat(MAX_NOTIFICATION_DESCRIPTION_LENGTH + 1))).toBe(true);
  });

  it('eventPreventDefault calls preventDefault on the passed event', () => {
    const preventDefault = vi.fn();
    eventPreventDefault({ preventDefault } as any);
    expect(preventDefault).toHaveBeenCalledTimes(1);
  });

  it('eventstopPropagation calls stopPropagation on the passed event', () => {
    const stopPropagation = vi.fn();
    eventstopPropagation({ stopPropagation } as any);
    expect(stopPropagation).toHaveBeenCalledTimes(1);
  });
});

describe('getErrorCodesFromEvent', () => {
  it('returns null when there is no event or metadata', () => {
    expect(getErrorCodesFromEvent(null)).toBeNull();
    expect(getErrorCodesFromEvent({} as any)).toBeNull();
  });

  it('collects codes from a single metadata.error object', () => {
    expect(getErrorCodesFromEvent({ metadata: { error: { Code: 'E1' } } } as any)).toEqual(['E1']);
  });

  it('collects codes from an array metadata.error and ModelDetails errors', () => {
    const codes = getErrorCodesFromEvent({
      metadata: {
        error: [{ Code: 'A' }, { Code: 'B' }, { Code: undefined }],
        ModelDetails: {
          a: { Errors: [{ error: { Code: 'C' } }, { error: { Code: 'A' } }] },
        },
      },
    } as any);
    expect(codes).toEqual(expect.arrayContaining(['A', 'B', 'C']));
    // Deduped via Set
    expect(codes!.filter((c) => c === 'A')).toHaveLength(1);
  });
});

describe('DeleteEvent', () => {
  beforeEach(() => deleteEventMutation.mockReset());

  it('fires the delete mutation when clicked', () => {
    render(<DeleteEvent event={{ id: 'evt-1' } as any} />);
    fireEvent.click(screen.getByText('Delete'));
    expect(deleteEventMutation).toHaveBeenCalledWith({ id: 'evt-1' });
  });
});

describe('ChangeStatus', () => {
  beforeEach(() => updateStatusMutation.mockReset());

  it('marks an unread event as read', () => {
    render(<ChangeStatus event={{ id: 'evt-1', status: 'unread' } as any} />);
    fireEvent.click(screen.getByText(/mark as read/i));
    expect(updateStatusMutation).toHaveBeenCalledWith({ id: 'evt-1', status: 'read' });
  });

  it('marks a read event as unread', () => {
    render(<ChangeStatus event={{ id: 'evt-2', status: 'read' } as any} />);
    fireEvent.click(screen.getByText(/mark as unread/i));
    expect(updateStatusMutation).toHaveBeenCalledWith({ id: 'evt-2', status: 'unread' });
  });
});

describe('Notification component', () => {
  beforeEach(() => {
    dispatchMock.mockReset();
    updateStatusMutation.mockReset();
    deleteEventMutation.mockReset();
    selectors.event = {
      id: 'evt-1',
      description: 'Short description',
      severity: 'informational',
      status: 'unread',
      metadata: {},
      checked: false,
      userID: 'u1',
      systemID: 's1',
      createdAt: '2024-01-01T00:00:00.000Z',
    };
    selectors.visible = true;
    selectors.ui = { history_mode: false };
  });

  it('renders summary content for the event', () => {
    render(<Notification event_id="evt-1" />);
    expect(screen.getByTestId('notification-summary')).toBeInTheDocument();
    expect(screen.getByText('Short description')).toBeInTheDocument();
    expect(screen.getByTestId('formatted-time')).toHaveTextContent('2024-01-01T00:00:00.000Z');
  });

  it('toggles expanded detail panel on summary click', () => {
    render(<Notification event_id="evt-1" />);
    // collapse closed initially
    expect(screen.queryByTestId('collapse-open')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('notification-summary'));
    expect(screen.getByTestId('collapse-open')).toBeInTheDocument();
    expect(screen.getByTestId('formatted-meta')).toBeInTheDocument();
  });

  it('dispatches updateIsEventChecked when the checkbox is changed', () => {
    render(<Notification event_id="evt-1" />);
    fireEvent.click(screen.getByTestId('event-checkbox'));
    fireEvent.change(screen.getByTestId('event-checkbox'), { target: { checked: true } });
    expect(dispatchMock).toHaveBeenCalled();
    const lastCall = dispatchMock.mock.calls[dispatchMock.mock.calls.length - 1][0];
    expect(lastCall.type).toBe('events/updateIsEventChecked');
  });

  it('prefers history_title over description when in history mode', () => {
    selectors.ui = { history_mode: true };
    selectors.event = {
      ...selectors.event,
      metadata: { history_title: 'History Title' },
    };
    render(<Notification event_id="evt-1" />);
    expect(screen.getByText('History Title')).toBeInTheDocument();
  });
});
