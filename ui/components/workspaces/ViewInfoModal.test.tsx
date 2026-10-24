import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ViewInfoModal_, ViewInfoModal, UserChip } from './ViewInfoModal';

const notify = vi.fn();
const mockUpdateView = vi.fn().mockReturnValue({ unwrap: () => Promise.resolve() });

let mockViewQuery: any = {
  data: {
    id: 'v1',
    userId: 'user-1',
    visibility: 'public',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-02-01T00:00:00Z',
  },
  isLoading: false,
  isFetching: false,
};

let mockUserQuery: any = {
  data: { id: 'user-1' },
  isLoading: false,
  isFetching: false,
};

let mockUserProfile: any = {
  data: { firstName: 'Alice', lastName: 'Smith', avatarUrl: '/avatar.png' },
  isError: false,
  isLoading: false,
};

vi.mock('@/rtk-query/view', () => ({
  useGetViewQuery: () => mockViewQuery,
  useUpdateViewVisibilityMutation: () => [mockUpdateView],
}));

vi.mock('@/rtk-query/user', () => ({
  useGetLoggedInUserQuery: () => mockUserQuery,
  useGetUserProfileSummaryByIdQuery: () => mockUserProfile,
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: { INFO: 'info', ERROR: 'error' },
}));

vi.mock('./SpacesSwitcher/hooks', () => ({
  handleUpdateViewVisibility: vi.fn(),
  viewPath: (view: any) => `/views/${view?.id}`,
}));

vi.mock('@/store/ProviderStoreWrapper', () => ({
  default: ({ children }: any) => <div data-testid="provider-wrapper">{children}</div>,
}));

vi.mock('css/icons.styles', () => ({ iconLarge: {} }));

vi.mock('@/utils/Enum', () => ({
  VIEW_VISIBILITY: { PUBLIC: 'public', PRIVATE: 'private' },
}));

vi.mock('@/assets/icons', () => ({
  Lock: () => <svg data-testid="lock" />,
  Public: () => <svg data-testid="public" />,
}));

vi.mock('rehype-sanitize', () => ({ default: () => null }));

vi.mock('../general/Markdown', () => ({
  MDEditor: ({ value, onChange }: any) => (
    <textarea data-testid="md-editor" value={value} onChange={(e) => onChange?.(e.target.value)} />
  ),
}));

vi.mock('../meshery-mesh-interface/PatternService/RJSF_wrapper', () => ({
  default: ({ formData, onChange, widgets }: any) => {
    const Widget = widgets?.markdown;
    return (
      <div data-testid="rjsf-wrapper">
        <button onClick={() => onChange?.({ ...(formData || {}), notes: 'updated' })}>
          change
        </button>
        {Widget ? <Widget label="Notes" value={formData?.notes} onChange={vi.fn()} /> : null}
      </div>
    );
  },
}));

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = (props: any) => <Component {...props}>{props.children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };

  return {
    Avatar: ({ src }: any) => <img data-testid="avatar" alt="avatar" src={src} />,
    Chip: ({ avatar, label }: any) => (
      <div data-testid="chip">
        {avatar}
        <span>{label}</span>
      </div>
    ),
    CircularProgress: () => <div data-testid="progress" />,
    FormLabel: ({ children }: any) => <label>{children}</label>,
    Typography: ({ children, variant }: any) => <span data-variant={variant}>{children}</span>,
    ViewIcon: () => <svg data-testid="view-icon" />,
    getFullFormattedTime: (t: string) => `formatted:${t}`,
    Box: ({ children, style, ...rest }: any) => (
      <div data-testid="box" style={style} {...rest}>
        {children}
      </div>
    ),
    Modal: ({ open, closeModal, title, children }: any) =>
      open ? (
        <div data-testid="modal" data-title={title}>
          <button onClick={closeModal} aria-label="modal-close">
            close
          </button>
          {children}
        </div>
      ) : null,
    ModalBody: ({ children }: any) => <div data-testid="modal-body">{children}</div>,
    ModalFooter: ({ children }: any) => <div data-testid="modal-footer">{children}</div>,
    VisibilityChipMenu: ({ value, onChange, enabled }: any) => (
      <button
        data-testid="visibility-menu"
        data-value={value}
        data-enabled={String(enabled)}
        onClick={() => onChange('private')}
      >
        {value}
      </button>
    ),
    ModalButtonSecondary: ({ children, onClick, disabled }: any) => (
      <button onClick={onClick} disabled={disabled} data-testid="copy-button">
        {children}
      </button>
    ),
    ModalButtonPrimary: ({ children, onClick, disabled }: any) => (
      <button onClick={onClick} disabled={disabled} data-testid="save-button">
        {children}
      </button>
    ),
    createTheme: () => ({
      breakpoints: {
        down: () => '',
        up: () => '',
      },
    }),
    useTheme: () => ({
      palette: {
        text: { primary: '#000' },
        mode: 'light',
        common: { white: '#fff' },
        background: { paper: '#fafafa' },
      },
    }),
    styled,
  };
});

describe('ViewInfoModal_', () => {
  beforeEach(() => {
    notify.mockReset();
    mockUpdateView.mockClear();
    mockViewQuery = {
      data: {
        id: 'v1',
        userId: 'user-1',
        visibility: 'public',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-02-01T00:00:00Z',
      },
      isLoading: false,
      isFetching: false,
    };
    mockUserQuery = {
      data: { id: 'user-1' },
      isLoading: false,
      isFetching: false,
    };
    mockUserProfile = {
      data: { firstName: 'Alice', lastName: 'Smith', avatarUrl: '/a.png' },
      isError: false,
      isLoading: false,
    };
  });

  const baseProps = {
    open: true,
    closeModal: vi.fn(),
    viewId: 'v1',
    viewName: 'My View',
    metadata: { notes: 'init' },
  };

  it('renders modal with view name', () => {
    render(<ViewInfoModal_ {...baseProps} />);
    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'My View');
  });

  it('renders the loading state when view query is loading', () => {
    mockViewQuery = { ...mockViewQuery, isLoading: true };
    render(<ViewInfoModal_ {...baseProps} />);
    expect(screen.getByTestId('progress')).toBeInTheDocument();
  });

  it('renders the formatted created and updated times', () => {
    render(<ViewInfoModal_ {...baseProps} />);
    expect(screen.getByText('formatted:2024-01-01T00:00:00Z')).toBeInTheDocument();
    expect(screen.getByText('formatted:2024-02-01T00:00:00Z')).toBeInTheDocument();
  });

  it('shows the Save button when the user is the owner', () => {
    render(<ViewInfoModal_ {...baseProps} />);
    expect(screen.getByTestId('save-button')).toBeInTheDocument();
  });

  it('hides the Save button when not the owner', () => {
    mockUserQuery = { ...mockUserQuery, data: { id: 'user-2' } };
    render(<ViewInfoModal_ {...baseProps} />);
    expect(screen.queryByTestId('save-button')).not.toBeInTheDocument();
  });

  it('calls updateView when save is clicked', async () => {
    const user = userEvent.setup();
    const closeModal = vi.fn();

    render(<ViewInfoModal_ {...baseProps} closeModal={closeModal} />);

    await user.click(screen.getByTestId('save-button'));

    expect(mockUpdateView).toHaveBeenCalledWith({
      id: 'v1',
      body: { metadata: { notes: 'init' } },
    });
  });

  it('disables copy link when view does not exist', () => {
    mockViewQuery = { ...mockViewQuery, data: undefined };
    render(<ViewInfoModal_ {...baseProps} />);
    expect(screen.getByTestId('copy-button')).toBeDisabled();
  });

  it('copies the link to the clipboard when copy is clicked', async () => {
    const user = userEvent.setup();
    const writeText = vi.fn();
    Object.defineProperty(navigator, 'clipboard', {
      value: { writeText },
      writable: true,
      configurable: true,
    });

    render(<ViewInfoModal_ {...baseProps} />);

    await user.click(screen.getByTestId('copy-button'));

    expect(writeText).toHaveBeenCalledWith('/views/v1');
    expect(notify).toHaveBeenCalledWith({
      message: 'Link copied to clipboard',
      event_type: 'info',
    });
  });
});

describe('UserChip', () => {
  beforeEach(() => {
    mockUserProfile = {
      data: { firstName: 'Alice', lastName: 'Smith', avatarUrl: '/a.png' },
      isError: false,
      isLoading: false,
    };
  });

  it('renders the username and avatar when data is loaded', () => {
    render(<UserChip userId="user-1" />);
    expect(screen.getByTestId('avatar')).toHaveAttribute('src', '/a.png');
    expect(screen.getByText('Alice Smith')).toBeInTheDocument();
  });

  it('renders nothing while loading', () => {
    mockUserProfile = { isLoading: true, isError: false, data: null };
    const { container } = render(<UserChip userId="user-1" />);
    expect(container.textContent).toBe('');
  });

  it('renders nothing on error', () => {
    mockUserProfile = { isLoading: false, isError: true, data: null };
    const { container } = render(<UserChip userId="user-1" />);
    expect(container.textContent).toBe('');
  });
});

describe('ViewInfoModal', () => {
  it('wraps the view modal in a provider store wrapper', () => {
    render(
      <ViewInfoModal open={true} closeModal={vi.fn()} viewId="v1" viewName="Test" metadata={{}} />,
    );

    expect(screen.getByTestId('provider-wrapper')).toBeInTheDocument();
  });
});
