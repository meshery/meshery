import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const dispatch = vi.fn();
const notify = vi.fn();

let mockGetUserQuery: any = {
  data: undefined,
  isSuccess: false,
  isError: false,
  error: undefined,
};
let mockProviderCapabilities: any = { providerUrl: 'https://provider.test' };

vi.mock('react-redux', () => ({
  useDispatch: () => dispatch,
  useSelector: (selector: any) =>
    selector({ ui: { providerCapabilities: mockProviderCapabilities } }),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetLoggedInUserQuery: () => mockGetUserQuery,
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateUser: (payload: any) => ({ type: 'mesheryUi/updateUser', payload }),
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: { ERROR: { type: 'error' } },
}));

const ExtensionPointSchemaValidator = vi.fn((type: string) => (input: any) => input ?? []);

vi.mock('../utils/ExtensionPointSchemaValidator', () => ({
  default: (type: string) => ExtensionPointSchemaValidator(type),
}));

vi.mock('next/link', () => ({
  default: ({ href, children }: any) => (
    <a data-testid="next-link" href={href}>
      {children}
    </a>
  ),
}));

vi.mock('@sistent/sistent', () => ({
  Avatar: ({ src, sx, slotProps }: any) => (
    <img
      data-testid="avatar"
      src={src || ''}
      data-size={JSON.stringify(sx || {})}
      data-referrer={slotProps?.img?.referrerPolicy}
    />
  ),
  Button: ({ children, variant, color, ...props }: any) => (
    <button data-testid="button" data-variant={variant} data-color={color} {...props}>
      {children}
    </button>
  ),
  NoSsr: ({ children }: any) => <>{children}</>,
}));

vi.mock('./layout/Header/Header.styles', () => ({
  IconButtonAvatar: ({ children, onClick, color, 'aria-haspopup': hasPopup }: any) => (
    <button
      data-testid="icon-button-avatar"
      data-color={color}
      aria-haspopup={hasPopup}
      onClick={onClick}
    >
      {children}
    </button>
  ),
}));

import UserProvider from './User';

const setWindowLocation = (href: string) => {
  // jsdom: replace window.location with a stub whose setter captures string
  // assignments (e.g. `window.location = profileUrl`) into the .href getter.
  // This lets tests assert via window.location.href regardless of whether
  // the source writes to .location directly or to .location.href.
  let current = href;
  const stub: any = {
    get href() {
      return current;
    },
    set href(v: string) {
      current = String(v);
    },
    get origin() {
      try {
        return new URL(current).origin;
      } catch {
        return '';
      }
    },
    get pathname() {
      try {
        return new URL(current).pathname;
      } catch {
        return '';
      }
    },
    toString() {
      return current;
    },
  };
  Object.defineProperty(window, 'location', {
    configurable: true,
    set(v: any) {
      if (typeof v === 'string') {
        current = v;
      } else if (v && typeof v.href === 'string') {
        current = v.href;
      }
    },
    get() {
      return stub;
    },
  });
};

describe('User component', () => {
  beforeEach(() => {
    dispatch.mockClear();
    notify.mockClear();
    ExtensionPointSchemaValidator.mockClear();
    mockGetUserQuery = {
      data: undefined,
      isSuccess: false,
      isError: false,
      error: undefined,
    };
    mockProviderCapabilities = { providerUrl: 'https://provider.test' };
    setWindowLocation('http://localhost:9081/app');
  });

  it('renders an avatar button for a logged-in user', () => {
    mockGetUserQuery = {
      data: { status: 'authenticated', avatarUrl: 'https://cdn.test/me.png' },
      isSuccess: true,
      isError: false,
      error: undefined,
    };

    render(<UserProvider color="primary" />);

    expect(screen.getByTestId('profile-button')).toBeInTheDocument();
    expect(screen.getByTestId('avatar')).toHaveAttribute('src', 'https://cdn.test/me.png');
  });

  it('renders a Sign In button when the user is anonymous', () => {
    mockGetUserQuery = {
      data: { status: 'anonymous', id: 'anon-1' },
      isSuccess: true,
      isError: false,
      error: undefined,
    };

    render(<UserProvider />);

    expect(screen.getByTestId('sign-in-button')).toBeInTheDocument();
    const link = screen.getByTestId('next-link') as HTMLAnchorElement;
    expect(link.getAttribute('href')).toContain('https://provider.test');
    expect(link.getAttribute('href')).toContain('anonymousUserID=anon-1');
  });

  it('dispatches updateUser on successful login', async () => {
    mockGetUserQuery = {
      data: { status: 'authenticated', id: 'u-1' },
      isSuccess: true,
      isError: false,
      error: undefined,
    };

    render(<UserProvider />);

    await waitFor(() =>
      expect(dispatch).toHaveBeenCalledWith({
        type: 'mesheryUi/updateUser',
        payload: { user: { status: 'authenticated', id: 'u-1' } },
      }),
    );
  });

  it('calls notify on getUser error', async () => {
    mockGetUserQuery = {
      data: undefined,
      isSuccess: false,
      isError: true,
      error: { data: 'oops' },
    };

    render(<UserProvider />);

    await waitFor(() => expect(notify).toHaveBeenCalled());
    const [payload] = notify.mock.calls[0];
    expect(payload.message).toBe('Error fetching user');
    expect(payload.details).toBe('oops');
  });

  it('navigates to the profile URL when the avatar is clicked', async () => {
    const user = userEvent.setup();
    mockGetUserQuery = {
      data: { status: 'authenticated' },
      isSuccess: true,
      isError: false,
      error: undefined,
    };
    mockProviderCapabilities = {
      providerUrl: 'https://provider.test',
      extensions: {
        account: [{ title: 'Cloud Account', href: 'https://cloud.test/profile' }],
      },
    };

    render(<UserProvider />);

    // Trigger a re-render of the providerCapabilities effect.
    await waitFor(() => expect(ExtensionPointSchemaValidator).toHaveBeenCalledWith('account'));

    await user.click(screen.getByTestId('icon-button-avatar'));
    expect(window.location.href).toContain('https://cloud.test/profile');
  });

  it('does not redirect when no profile URL is present', async () => {
    const user = userEvent.setup();
    mockGetUserQuery = {
      data: { status: 'authenticated' },
      isSuccess: true,
      isError: false,
      error: undefined,
    };

    const startingHref = window.location.href;
    render(<UserProvider />);

    await user.click(screen.getByTestId('icon-button-avatar'));
    // window.location.href should still be the same since profileUrl is undefined
    expect(window.location.href).toBe(startingHref);
  });
});
