import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { PINNED_NAV_ITEMS_STORAGE_KEY } from '@/utils/hooks/usePinnedNavItems';

const dispatchMock = vi.fn();
const useMediaQueryMock = vi.fn(() => false);
const routerPushMock = vi.fn();

const stateContainer: {
  isDrawerCollapsed: boolean;
  catalogVisibility: boolean;
  meshAdapters: any[];
} = {
  isDrawerCollapsed: false,
  catalogVisibility: true,
  meshAdapters: [],
};

vi.mock('react-redux', () => ({
  useDispatch: () => dispatchMock,
  useSelector: (sel: any) =>
    sel({
      ui: {
        isDrawerCollapsed: stateContainer.isDrawerCollapsed,
        catalogVisibility: stateContainer.catalogVisibility,
      },
      adapter: { meshAdapters: stateContainer.meshAdapters },
    }),
  Provider: ({ children }: any) => <>{children}</>,
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ asPath: '/', push: routerPushMock }),
}));

vi.mock('next/link', () => ({
  default: ({ href, children, onClick, ...props }: any) => (
    <a href={href} onClick={onClick} {...props}>
      {children}
    </a>
  ),
}));

vi.mock('@sistent/sistent', () => ({
  styled: (Base: any) => () => Base,
  CustomTooltip: ({ children }: any) => <>{children}</>,
  ListItemIcon: ({ children }: any) => <div>{children}</div>,
  Grow: ({ children, in: open }: any) => (open ? <div>{children}</div> : null),
  ListItem: ({ children, button, ...props }: any) => (
    <div {...(button ? { role: 'button' } : {})} {...props}>
      {children}
    </div>
  ),
  List: ({ children }: any) => <ul>{children}</ul>,
  Collapse: ({ children, in: open }: any) => (open ? <div>{children}</div> : null),
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  NoSsr: ({ children }: any) => <>{children}</>,
  ClickAwayListener: ({ children }: any) => <div>{children}</div>,
  Zoom: ({ children, in: open }: any) => (open ? <div>{children}</div> : null),
  HelpOutlinedIcon: () => <svg data-testid="help-outlined" />,
  LeftArrowIcon: () => <svg data-testid="left-arrow" />,
  ExternalLinkIcon: () => <svg data-testid="external-link" />,
  OpenInNewIcon: () => <svg data-testid="open-in-new" />,
  RemoveIcon: () => <svg data-testid="remove" />,
  useTheme: () => ({
    breakpoints: { down: vi.fn((breakpoint: string) => breakpoint) },
    palette: {
      icon: { default: '#000', brand: '#brand' },
      background: {
        constant: { white: '#fff' },
        brand: { default: '#brand-default' },
      },
      common: { white: '#fff' },
      navigation: { secondary: '#nav-secondary', active: '#nav-active', hover: '#nav-hover' },
    },
  }),
  SlackIcon: () => <svg data-testid="slack-icon" />,
  FileIcon: () => <svg data-testid="file-icon" />,
  GithubIcon: () => <svg data-testid="github-icon" />,
  DiscussForumIcon: () => <svg data-testid="discuss-forum-icon" />,
  useMediaQuery: () => useMediaQueryMock(),
}));

vi.mock('../../meshery-mesh-interface/PatternService/CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children }: any) => <>{children}</>,
}));

vi.mock('../../../utils/ExtensionPointSchemaValidator', () => ({
  default: () => () => [],
}));

vi.mock('../../../css/disableComponent.styles', () => ({
  cursorNotAllowed: {},
  disabledStyle: {},
}));

vi.mock('../../../utils/disabledComponents', () => ({
  ProviderUiAccessControl: class {
    isNavigatorComponentEnabled() {
      return true;
    }
    providerPayload = {
      providerCapabilities: [{ feature: 'persist-meshery-patterns' }],
    };
    providerCapabilities = [{ feature: 'persist-meshery-patterns' }];
  },
}));

vi.mock('../../../constants/navigator', () => ({
  CONFIGURATION: 'CONFIGURATION',
  DASHBOARD: 'DASHBOARD',
  CATALOG: 'CATALOG',
  LIFECYCLE: 'LIFECYCLE',
  SERVICE_MESH: 'SERVICE_MESH',
  TOGGLER: 'TOGGLER',
}));

vi.mock('../../../css/icons.styles', () => ({
  iconSmall: {},
}));

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

vi.mock('../../general/style', () => {
  const make = (testId?: string) => {
    const Mock = ({ children, ...props }: any) => (
      <div {...(testId ? { 'data-testid': testId } : {})} {...props}>
        {children}
      </div>
    );
    return Mock;
  };
  return {
    HideScrollbar: make('hide-scrollbar'),
    LinkContainer: make(),
    ListIconSide: make(),
    MainListIcon: make(),
    MainLogo: () => <img data-testid="main-logo" />,
    MainLogoCollapsed: () => <img data-testid="main-logo-collapsed" />,
    MainLogoText: () => <img data-testid="main-logo-text" />,
    MainLogoTextCollapsed: () => <img data-testid="main-logo-text-collapsed" />,
    NavigatorList: make('nav-list'),
    NavigatorListItem: make(),
    NavigatorListItemII: make(),
    NavigatorListItemIII: make(),
    RootDiv: make(),
    SecondaryDivider: make('secondary-divider'),
    SideBarListItem: ({ children, component: Component, href, onClick, ...props }: any) => {
      if (Component) {
        return (
          <Component href={href} onClick={onClick} {...props}>
            {children}
          </Component>
        );
      }
      return (
        <div onClick={onClick} {...props}>
          {children}
        </div>
      );
    },
    SideBarText: make(),
    StyledListItem: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    NavigatorLink: make(),
    NavigatorHelpIcons: make(),
    HelpListItem: make(),
    HelpButton: make(),
    ChevronButtonWrapper: make(),
    FixedSidebarFooter: make('sidebar-footer'),
    SidebarDrawer: ({ children }: any) => <aside>{children}</aside>,
    ExpandMore: () => <svg data-testid="expand-more" />,
  };
});

vi.mock('@/rtk-query/user', () => ({
  getProviderCapabilities: vi.fn(() =>
    Promise.resolve({
      data: { providerUrl: '', extensions: { navigator: [] } },
      isSuccess: true,
      isError: false,
    }),
  ),
  getSystemVersion: vi.fn(() =>
    Promise.resolve({
      data: {
        build: '1.0.0',
        latest: '1.0.0',
        outdated: false,
        commitsha: 'abc',
        release_channel: 'stable',
      },
      isSuccess: true,
      isError: false,
    }),
  ),
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  toggleDrawer: (payload: any) => ({ type: 'ui/toggleDrawer', payload }),
  updateBetaBadge: (payload: any) => ({ type: 'ui/updateBetaBadge', payload }),
  updateExtensionType: (payload: any) => ({ type: 'ui/updateExtensionType', payload }),
  updateProviderCapabilities: (payload: any) => ({
    type: 'ui/updateProviderCapabilities',
    payload,
  }),
  updateTitle: (payload: any) => ({ type: 'ui/updateTitle', payload }),
}));

vi.mock('@/store/slices/adapter', () => ({
  setAdapter: (payload: any) => ({ type: 'adapter/setAdapter', payload }),
}));

vi.mock('./navigatorComponents', () => ({
  getNavigatorComponents: () => [
    {
      id: 'DASHBOARD',
      title: 'Dashboard',
      icon: <span data-testid="dash-icon">D</span>,
      href: '/',
      show: true,
      link: true,
      submenu: true,
    },
    {
      id: 'CONFIGURATION',
      title: 'Configuration',
      icon: <span data-testid="config-icon">C</span>,
      href: '/configuration',
      show: true,
      link: true,
      submenu: true,
      children: [
        {
          id: 'Designs',
          title: 'Designs',
          icon: <span data-testid="design-icon">d</span>,
          href: '/configuration/designs',
          link: true,
        },
        {
          id: 'CATALOG',
          title: 'Catalog',
          icon: <span data-testid="catalog-icon">c</span>,
          href: '/configuration/catalog',
          link: true,
        },
      ],
    },
  ],
}));

import Navigator from './Navigator';

describe('Navigator Favorites and Pinned Navigation', () => {
  beforeEach(() => {
    window.localStorage.clear();
    dispatchMock.mockReset();
    routerPushMock.mockReset();
    useMediaQueryMock.mockReturnValue(false);
    stateContainer.isDrawerCollapsed = false;
    stateContainer.catalogVisibility = true;
    stateContainer.meshAdapters = [];
  });

  it('renders the Favorites section and empty state when nothing is pinned', async () => {
    render(<Navigator />);
    expect(await screen.findByTestId('favorites-section')).toBeInTheDocument();
    expect(screen.getByTestId('favorites-header')).toHaveTextContent('Favorites');
    expect(screen.getByTestId('favorites-empty-state')).toHaveTextContent('No pinned favorites');
  });

  it('renders collapsed empty state indicator when drawer is collapsed', async () => {
    stateContainer.isDrawerCollapsed = true;
    render(<Navigator />);
    expect(await screen.findByTestId('favorites-empty-state-collapsed')).toBeInTheDocument();
  });

  it('pins a leaf item when its pin button is clicked', async () => {
    render(<Navigator />);
    await screen.findByTestId('config-icon');

    const pinButton = screen.getByTestId('pin-button-DASHBOARD');
    expect(pinButton).toHaveAttribute('aria-pressed', 'false');
    expect(pinButton).toHaveAttribute('aria-label', 'Pin Dashboard to Favorites');

    fireEvent.click(pinButton);

    expect(await screen.findByTestId('favorite-DASHBOARD')).toBeInTheDocument();
    expect(pinButton).toHaveAttribute('aria-pressed', 'true');
    expect(pinButton).toHaveAttribute('aria-label', 'Unpin Dashboard from Favorites');
    expect(JSON.parse(window.localStorage.getItem(PINNED_NAV_ITEMS_STORAGE_KEY)!)).toEqual([
      'DASHBOARD',
    ]);
  });

  it('unpins an item from Favorites using the unpin button in Favorites', async () => {
    window.localStorage.setItem(PINNED_NAV_ITEMS_STORAGE_KEY, JSON.stringify(['DASHBOARD']));

    render(<Navigator />);
    expect(await screen.findByTestId('favorite-DASHBOARD')).toBeInTheDocument();

    const unpinButton = screen.getByTestId('unpin-button-DASHBOARD');
    expect(unpinButton).toHaveAttribute('aria-pressed', 'true');
    expect(unpinButton).toHaveAttribute('aria-label', 'Unpin Dashboard from Favorites');

    fireEvent.click(unpinButton);

    expect(screen.queryByTestId('favorite-DASHBOARD')).not.toBeInTheDocument();
    expect(screen.getByTestId('favorites-empty-state')).toBeInTheDocument();
    expect(JSON.parse(window.localStorage.getItem(PINNED_NAV_ITEMS_STORAGE_KEY)!)).toEqual([]);
  });

  it('persists pinned items across refresh/remount', async () => {
    window.localStorage.setItem(PINNED_NAV_ITEMS_STORAGE_KEY, JSON.stringify(['CATALOG']));

    const { unmount } = render(<Navigator />);
    expect(await screen.findByTestId('favorite-CATALOG')).toBeInTheDocument();
    unmount();

    render(<Navigator />);
    expect(await screen.findByTestId('favorite-CATALOG')).toBeInTheDocument();
  });

  it('prevents duplicate IDs from creating multiple favorite entries', async () => {
    window.localStorage.setItem(
      PINNED_NAV_ITEMS_STORAGE_KEY,
      JSON.stringify(['DASHBOARD', 'DASHBOARD']),
    );

    render(<Navigator />);
    await screen.findByTestId('favorite-DASHBOARD');
    const favorites = screen.getAllByTestId('favorite-DASHBOARD');
    expect(favorites).toHaveLength(1);
  });

  it('silently skips pinned IDs that no longer resolve without crashing', async () => {
    window.localStorage.setItem(
      PINNED_NAV_ITEMS_STORAGE_KEY,
      JSON.stringify(['NON_EXISTENT_ID', 'DASHBOARD']),
    );

    render(<Navigator />);
    expect(await screen.findByTestId('favorite-DASHBOARD')).toBeInTheDocument();
    expect(screen.queryByTestId('favorite-NON_EXISTENT_ID')).not.toBeInTheDocument();
    expect(JSON.parse(window.localStorage.getItem(PINNED_NAV_ITEMS_STORAGE_KEY)!)).toContain(
      'NON_EXISTENT_ID',
    );
  });

  it('renders favorite rows as real links with matching href', async () => {
    window.localStorage.setItem(PINNED_NAV_ITEMS_STORAGE_KEY, JSON.stringify(['DASHBOARD']));

    render(<Navigator />);
    const favoriteItem = await screen.findByTestId('favorite-DASHBOARD');
    const linkElement = favoriteItem.closest('a');
    expect(linkElement).toBeInTheDocument();
    expect(linkElement).toHaveAttribute('href', '/');
  });

  it('operates pin button with keyboard Enter and Space keys via userEvent', async () => {
    const user = userEvent.setup();
    render(<Navigator />);
    await screen.findByTestId('config-icon');

    const pinButton = screen.getByTestId('pin-button-DASHBOARD');
    expect(pinButton.tagName).toBe('BUTTON');

    for (let i = 0; i < 10; i++) {
      if (document.activeElement === pinButton) break;
      await user.tab();
    }
    expect(pinButton).toHaveFocus();

    // Activate with Enter key to pin
    await user.keyboard('{Enter}');
    expect(await screen.findByTestId('favorite-DASHBOARD')).toBeInTheDocument();
    expect(pinButton).toHaveAttribute('aria-pressed', 'true');

    // Activate with Space key to unpin
    await user.keyboard('{ }');
    expect(screen.queryByTestId('favorite-DASHBOARD')).not.toBeInTheDocument();
    expect(pinButton).toHaveAttribute('aria-pressed', 'false');
  });

  it('pins a nested child after expanding its parent category', async () => {
    render(<Navigator />);
    await screen.findByTestId('config-icon');

    expect(screen.queryByTestId('pin-button-Designs')).not.toBeInTheDocument();

    const configCategory = screen.getByTestId('CONFIGURATION');
    fireEvent.click(configCategory);

    const designsPinButton = await screen.findByTestId('pin-button-Designs');
    expect(designsPinButton).toHaveAttribute('aria-pressed', 'false');

    fireEvent.click(designsPinButton);

    expect(await screen.findByTestId('favorite-Designs')).toBeInTheDocument();
    expect(designsPinButton).toHaveAttribute('aria-pressed', 'true');
    expect(JSON.parse(window.localStorage.getItem(PINNED_NAV_ITEMS_STORAGE_KEY)!)).toContain(
      'Designs',
    );
  });

  it('hides the pin button when the drawer is collapsed', async () => {
    stateContainer.isDrawerCollapsed = true;
    window.localStorage.setItem(PINNED_NAV_ITEMS_STORAGE_KEY, JSON.stringify(['DASHBOARD']));

    render(<Navigator />);
    await screen.findByTestId('favorite-DASHBOARD');

    expect(screen.queryByTestId('pin-button-DASHBOARD')).not.toBeInTheDocument();
    expect(screen.queryByTestId('unpin-button-DASHBOARD')).not.toBeInTheDocument();
  });

  it('preserves the order of favorites based on stored pinned IDs', async () => {
    window.localStorage.setItem(
      PINNED_NAV_ITEMS_STORAGE_KEY,
      JSON.stringify(['CATALOG', 'DASHBOARD']),
    );

    render(<Navigator />);
    await screen.findByTestId('favorite-CATALOG');
    await screen.findByTestId('favorite-DASHBOARD');

    const favoriteElements = screen.getAllByTestId(/^favorite-/);
    expect(favoriteElements[0]).toHaveAttribute('data-testid', 'favorite-CATALOG');
    expect(favoriteElements[1]).toHaveAttribute('data-testid', 'favorite-DASHBOARD');
  });
});
