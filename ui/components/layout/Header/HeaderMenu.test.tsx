import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const canSpy = vi.fn((_key?: unknown) => true);

const routerPush = vi.fn();
const dispatchMock = vi.fn();
const notifyMock = vi.fn();
const openRegistryModalMock = vi.fn();

const getLoggedInUserState: { data: any; isSuccess: boolean; isError: boolean; error: any } = {
  data: { firstName: 'Alice', status: 'active' },
  isSuccess: true,
  isError: false,
  error: null,
};

const tokenState: { isError: boolean; error: any } = {
  isError: false,
  error: null,
};

const triggerGetTokenMock = vi.fn(() => ({
  unwrap: () => Promise.resolve({ token: 'jwt-token' }),
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ push: routerPush }),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetLoggedInUserQuery: () => getLoggedInUserState,
  useLazyGetTokenQuery: () => [triggerGetTokenMock, tokenState],
}));

// Note: HeaderMenu only adds "Get Token" when there are NO account extension
// items. Default behavior here: return an empty array so Get Token is included.
const extensionPointResult: { items: any[] } = { items: [] };
vi.mock('../../../utils/ExtensionPointSchemaValidator', () => ({
  default: () => () => extensionPointResult.items,
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: notifyMock }),
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: { ERROR: 'error' },
}));

vi.mock('@sistent/sistent', () => ({
  MenuIcon: () => <svg data-testid="menu-icon" />,
  // Mirrors the real NavigationNavbar contract: an item whose `permission` is
  // false is rendered but disabled (`permission ?? true`), not hidden.
  NavigationNavbar: ({ navigationItems }: any) => (
    <ul data-testid="navigation-navbar">
      {(navigationItems || []).map((item: any) => (
        <li key={item.id}>
          <button
            type="button"
            data-testid={`nav-item-${item.id}`}
            disabled={!(item.permission ?? true)}
            onClick={item.onClick}
          >
            {item.title}
          </button>
        </li>
      ))}
    </ul>
  ),
  Popover: ({ open, children }: any) =>
    open ? <div data-testid="header-popover">{children}</div> : null,
  useHasPermission: (key: { id?: string }) => canSpy(key),
}));

vi.mock('./Header.styles', () => ({
  IconButtonMenu: ({ children, onClick }: any) => (
    <button data-testid="menu-button" type="button" onClick={onClick}>
      {children}
    </button>
  ),
}));

vi.mock('react-redux', () => ({
  useDispatch: () => dispatchMock,
  useSelector: (sel: any) => sel({ ui: { providerCapabilities: { extensions: { account: [] } } } }),
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateExtensionType: (payload: any) => ({ type: 'ui/updateExtensionType', payload }),
  updateUser: (payload: any) => ({ type: 'ui/updateUser', payload }),
}));

vi.mock('@/utils/hooks/useRegistryModal', () => ({
  useRegistryModal: () => ({ openModal: openRegistryModalMock }),
}));

import { Keys } from '@meshery/schemas/permissions';
import HeaderMenu from './HeaderMenu';

describe('HeaderMenu', () => {
  beforeEach(() => {
    canSpy.mockClear();
    canSpy.mockReturnValue(true);
    routerPush.mockReset();
    dispatchMock.mockReset();
    notifyMock.mockReset();
    openRegistryModalMock.mockReset();
    triggerGetTokenMock.mockClear();
    getLoggedInUserState.data = { firstName: 'Alice', status: 'active' };
    getLoggedInUserState.isSuccess = true;
    getLoggedInUserState.isError = false;
    getLoggedInUserState.error = null;
    tokenState.isError = false;
    tokenState.error = null;
    extensionPointResult.items = [];

    // Stub window.location.assign so logout doesn't throw in jsdom.
    Object.defineProperty(window, 'location', {
      value: { href: '', assign: vi.fn() },
      writable: true,
      configurable: true,
    });
  });

  it('renders the menu button and opens the popover on click', () => {
    render(<HeaderMenu />);
    expect(screen.queryByTestId('header-popover')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('menu-button'));
    expect(screen.getByTestId('header-popover')).toBeInTheDocument();
  });

  it('returns null for anonymous users', () => {
    getLoggedInUserState.data = { status: 'anonymous' };
    const { container } = render(<HeaderMenu />);
    expect(container).toBeEmptyDOMElement();
  });

  it('dispatches updateUser when the user query resolves', () => {
    render(<HeaderMenu />);
    expect(dispatchMock).toHaveBeenCalledWith({
      type: 'ui/updateUser',
      payload: { user: getLoggedInUserState.data },
    });
  });

  it('notifies on user query error', () => {
    dispatchMock.mockReset();
    notifyMock.mockReset();
    getLoggedInUserState.isSuccess = false;
    getLoggedInUserState.isError = true;
    getLoggedInUserState.error = { data: 'bad' };
    render(<HeaderMenu />);
    expect(notifyMock).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Error fetching user',
        event_type: 'error',
        details: 'bad',
      }),
    );
  });

  it('notifies on token query error via a useEffect', () => {
    notifyMock.mockReset();
    tokenState.isError = true;
    tokenState.error = { data: 'token-bad' };
    render(<HeaderMenu />);
    expect(notifyMock).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Error fetching token',
        event_type: 'error',
        details: 'token-bad',
      }),
    );
  });

  it('triggers preferences navigation when the Preferences item is clicked', () => {
    render(<HeaderMenu />);
    fireEvent.click(screen.getByTestId('menu-button'));
    fireEvent.click(screen.getByTestId('nav-item-preferences'));
    expect(routerPush).toHaveBeenCalledWith('/user/preferences');
  });

  it('triggers settings navigation when the Settings item is clicked', () => {
    render(<HeaderMenu />);
    fireEvent.click(screen.getByTestId('menu-button'));
    fireEvent.click(screen.getByTestId('nav-item-settings'));
    expect(routerPush).toHaveBeenCalledWith('/settings');
  });

  it('triggers connections navigation when the Connections item is clicked', () => {
    render(<HeaderMenu />);
    fireEvent.click(screen.getByTestId('menu-button'));
    fireEvent.click(screen.getByTestId('nav-item-Connections'));
    expect(routerPush).toHaveBeenCalledWith('/management/connections');
  });

  it('opens the registry modal when the Registry item is clicked', () => {
    render(<HeaderMenu />);
    fireEvent.click(screen.getByTestId('menu-button'));
    fireEvent.click(screen.getByTestId('nav-item-registry'));
    expect(openRegistryModalMock).toHaveBeenCalled();
  });

  it('navigates to logout on Logout click', () => {
    render(<HeaderMenu />);
    fireEvent.click(screen.getByTestId('menu-button'));
    fireEvent.click(screen.getByTestId('nav-item-logout'));
    expect(window.location.href).toBe('/user/logout');
  });

  it('disables the Get Token item when the download-token permission is denied', () => {
    canSpy.mockReturnValue(false);
    render(<HeaderMenu />);
    fireEvent.click(screen.getByTestId('menu-button'));

    // NavigationNavbar renders a permission-denied item disabled rather than
    // hiding it, so the item stays present but is not actionable.
    expect(screen.getByTestId('nav-item-get-token')).toBeDisabled();
    expect(canSpy).toHaveBeenCalledWith(Keys.SecurityManagementDownloadToken);
  });

  it('triggers a token download when the Get Token item is clicked', async () => {
    // Create a link factory we can spy on.
    const click = vi.fn();
    const remove = vi.fn();
    const setAttribute = vi.fn();
    const createElementSpy = vi
      .spyOn(document, 'createElement')
      .mockImplementation((tag: string) => {
        if (tag === 'a') {
          return { setAttribute, click, remove } as any;
        }
        return document.createElementNS('http://www.w3.org/1999/xhtml', tag);
      });

    render(<HeaderMenu />);
    fireEvent.click(screen.getByTestId('menu-button'));
    fireEvent.click(screen.getByTestId('nav-item-get-token'));

    await waitFor(() => expect(triggerGetTokenMock).toHaveBeenCalled());
    await waitFor(() => expect(click).toHaveBeenCalled());
    expect(setAttribute).toHaveBeenCalledWith('download', 'auth.json');
    createElementSpy.mockRestore();
  });
});
