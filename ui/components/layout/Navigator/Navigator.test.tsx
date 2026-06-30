import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const dispatchMock = vi.fn();
const useMediaQueryMock = vi.fn(() => false);

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
  useRouter: () => ({ asPath: '/', push: vi.fn() }),
}));

vi.mock('next/link', () => ({
  default: ({ href, children }: any) => <a href={href}>{children}</a>,
}));

vi.mock('@sistent/sistent', () => ({
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
  Zoom: ({ children, in: open }: any) => (open ? <div>{children}</div> : null),
  HelpOutlinedIcon: () => <svg data-testid="help-outlined" />,
  LeftArrowIcon: () => <svg data-testid="left-arrow" />,
  ExternalLinkIcon: () => <svg data-testid="external-link" />,
  OpenInNewIcon: () => <svg data-testid="open-in-new" />,
  RemoveIcon: () => <svg data-testid="remove" />,
  useTheme: () => ({
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
    constructor(_caps: any) {}
    isNavigatorComponentEnabled() {
      return true;
    }
    providerCapabilities = [];
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
    MainLogo: ({ src, onClick }: any) => (
      <img data-testid="main-logo" src={src} onClick={onClick} />
    ),
    MainLogoCollapsed: ({ src, onClick }: any) => (
      <img data-testid="main-logo-collapsed" src={src} onClick={onClick} />
    ),
    MainLogoText: ({ src, onClick }: any) => (
      <img data-testid="main-logo-text" src={src} onClick={onClick} />
    ),
    MainLogoTextCollapsed: ({ src, onClick }: any) => (
      <img data-testid="main-logo-text-collapsed" src={src} onClick={onClick} />
    ),
    NavigatorList: make('nav-list'),
    NavigatorListItem: make(),
    NavigatorListItemII: make(),
    NavigatorListItemIII: make(),
    RootDiv: make(),
    SecondaryDivider: make(),
    SideBarListItem: make(),
    SideBarText: make(),
    StyledListItem: ({ children, component, onClick, ...props }: any) => {
      const Comp = component || 'div';
      return (
        <Comp {...props} onClick={onClick} data-testid="title-link">
          {children}
        </Comp>
      );
    },
    NavigatorLink: make(),
    NavigatorHelpIcons: make(),
    HelpListItem: make(),
    HelpButton: ({ children, onClick }: any) => (
      <button data-testid="help-button" type="button" onClick={onClick}>
        {children}
      </button>
    ),
    ChevronButtonWrapper: ({ children, ...props }: any) => (
      <div data-testid="chevron-wrapper" {...props}>
        {children}
      </div>
    ),
    FixedSidebarFooter: make('sidebar-footer'),
    SidebarDrawer: ({ children, isCollapsed }: any) => (
      <aside data-testid="sidebar-drawer" data-collapsed={String(Boolean(isCollapsed))}>
        {children}
      </aside>
    ),
    ExpandMore: () => <svg data-testid="expand-more" />,
  };
});

// fake user RTK-style helpers
vi.mock('@/rtk-query/user', () => ({
  getProviderCapabilities: vi.fn(() =>
    Promise.resolve({ data: { providerUrl: '', extensions: {} }, isSuccess: true, isError: false }),
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
      hovericon: <span data-testid="dash-hover-icon">D</span>,
      href: '/',
      show: true,
      link: true,
      submenu: true,
    },
    {
      id: 'CONFIGURATION',
      title: 'Configuration',
      icon: <span data-testid="config-icon">C</span>,
      hovericon: <span data-testid="config-hover-icon">C</span>,
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

describe('Navigator', () => {
  beforeEach(() => {
    dispatchMock.mockReset();
    useMediaQueryMock.mockReset().mockReturnValue(false);
    stateContainer.isDrawerCollapsed = false;
    stateContainer.catalogVisibility = true;
    stateContainer.meshAdapters = [];
  });

  it('renders the sidebar drawer wrapper and logo', () => {
    render(<Navigator />);
    expect(screen.getByTestId('sidebar-drawer')).toBeInTheDocument();
    expect(screen.getByTestId('main-logo')).toHaveAttribute(
      'src',
      '/static/img/meshery-logo/meshery-logo.png',
    );
  });

  it('renders the collapsed logos when the drawer is collapsed', () => {
    stateContainer.isDrawerCollapsed = true;
    render(<Navigator />);
    expect(screen.getByTestId('main-logo-collapsed')).toBeInTheDocument();
    expect(screen.getByTestId('main-logo-text-collapsed')).toBeInTheDocument();
  });

  it('marks the sidebar as collapsed via data-collapsed when collapsed', () => {
    stateContainer.isDrawerCollapsed = true;
    render(<Navigator />);
    expect(screen.getByTestId('sidebar-drawer')).toHaveAttribute('data-collapsed', 'true');
  });

  it('dispatches toggleDrawer when on a small screen and the drawer is not yet collapsed', () => {
    useMediaQueryMock.mockReturnValue(true);
    stateContainer.isDrawerCollapsed = false;
    render(<Navigator />);
    expect(dispatchMock.mock.calls.some(([action]) => action?.type === 'ui/toggleDrawer')).toBe(
      true,
    );
  });

  it('renders the chevron and version footer area', () => {
    render(<Navigator />);
    expect(screen.getByTestId('sidebar-footer')).toBeInTheDocument();
    expect(screen.getByTestId('chevron-wrapper')).toBeInTheDocument();
  });

  it('renders an expand/collapse caret only for nav items that have children', async () => {
    render(<Navigator />);
    // The list populates once provider capabilities resolve; wait for the parent row.
    await screen.findByTestId('config-icon');
    // The mocked navigator list has one leaf item (DASHBOARD, no children) and one
    // parent item (CONFIGURATION, with children). The caret is gated on children, so
    // leaf items must not render one - this keeps the caret button out of leaf anchors.
    expect(screen.getByTestId('dash-icon')).toBeInTheDocument();
    expect(screen.queryAllByTestId('expand-more')).toHaveLength(1);
  });
});
