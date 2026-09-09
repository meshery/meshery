import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Only the My-Views gate is keyed; every other capability stays granted so the
// navigation itself is what the assertions exercise.
let canViewViews = true;

vi.mock('@sistent/sistent', () => ({
  useHasPermission: () => canViewViews,
  useMediaQuery: () => false,
  Box: ({ children }: any) => <div>{children}</div>,
  Divider: () => <hr />,
  ErrorBoundary: ({ children }: any) => <>{children}</>,
  IconButton: ({ children, onClick }: any) => (
    <button type="button" onClick={onClick}>
      {children}
    </button>
  ),
  List: ({ children }: any) => <ul>{children}</ul>,
  AccessTimeFilledIcon: () => <svg />,
  ChevronLeftIcon: () => <svg />,
  ChevronRightIcon: () => <svg />,
  DesignIcon: () => <svg />,
  PeopleIcon: () => <svg />,
  ViewIcon: () => <svg />,
  WorkspaceIcon: () => <svg />,
}));

vi.mock('@/theme', () => ({
  styled: () => () =>
    function Styled({ children }: any) {
      return <div>{children}</div>;
    },
  useTheme: () => ({
    palette: { icon: { default: '#000', neutral: { default: '#000' } } },
    breakpoints: { down: () => '' },
  }),
}));

vi.mock('./SpacesSwitcher/styles', () => ({
  DrawerHeader: ({ children }: any) => <div>{children}</div>,
  StyledDrawer: ({ children }: any) => <aside>{children}</aside>,
  StyledMainContent: ({ children }: any) => <main>{children}</main>,
}));

vi.mock('./WorkspaceFormModalSections', () => ({
  NavItem: ({ item, onSelect }: any) => (
    <li>
      <button type="button" data-testid={`nav-${item.id}`} onClick={() => onSelect(item.id)}>
        {item.label}
      </button>
    </li>
  ),
  WorkspacesSection: () => <li data-testid="workspaces-section" />,
}));

vi.mock('./SpacesSwitcher/MyViewsContent', () => ({
  default: () => <div data-testid="my-views-content" />,
}));
vi.mock('./SpacesSwitcher/MyDesignsContent', () => ({
  default: () => <div data-testid="my-designs-content" />,
}));
vi.mock('./SpacesSwitcher/RecentContent', () => ({
  default: () => <div data-testid="recent-content" />,
}));
vi.mock('./SpacesSwitcher/SharedContent', () => ({
  default: () => <div data-testid="shared-content" />,
}));
vi.mock('./SpacesSwitcher/WorkspaceContent', () => ({
  default: () => <div data-testid="workspace-content" />,
}));
vi.mock('@/components/lifecycle', () => ({
  WorkspacesComponent: () => <div data-testid="all-workspaces" />,
}));

vi.mock('@/rtk-query/workspace', () => ({
  useGetWorkspacesQuery: () => ({ data: { workspaces: [] }, isLoading: false }),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetProviderCapabilitiesQuery: () => ({ data: { provider_type: 'remote' } }),
  useGetSelectedOrganization: () => ({ selectedOrganization: { id: 'org-1' } }),
}));

vi.mock('@/utils/provider', () => ({
  isLocalProvider: () => false,
}));

vi.mock('css/icons.styles', () => ({ iconMedium: {}, iconSmall: {} }));

vi.mock('@/utils/context/WorkspaceModalContextProvider', () => ({
  WorkspaceModalContext: React.createContext({
    selectedWorkspace: { id: null, name: null },
    setSelectedWorkspace: vi.fn(),
  }),
}));

import { Navigation } from './WorkspaceFormModalNav';

const renderNav = () => {
  const setHeaderInfo = vi.fn();
  const view = render(<Navigation setHeaderInfo={setHeaderInfo} />);
  return { ...view, setHeaderInfo };
};

const lastHeaderTitle = (setHeaderInfo: ReturnType<typeof vi.fn>) =>
  setHeaderInfo.mock.calls.at(-1)?.[0]?.title;

describe('WorkspaceFormModalNav', () => {
  beforeEach(() => {
    canViewViews = true;
  });

  it('lists My Views when the view permission is granted', () => {
    renderNav();
    expect(screen.getByTestId('nav-My-Views')).toBeInTheDocument();
  });

  it('omits My Views from the drawer when the view permission is denied', () => {
    canViewViews = false;
    renderNav();

    expect(screen.queryByTestId('nav-My-Views')).not.toBeInTheDocument();
    // The remaining items are unaffected - only the gated one disappears.
    expect(screen.getByTestId('nav-My-Designs')).toBeInTheDocument();
    expect(screen.getByTestId('nav-Shared-With-Me')).toBeInTheDocument();
  });

  it('renders the selected item content', async () => {
    const { getByTestId } = renderNav();
    getByTestId('nav-My-Views').click();
    expect(await screen.findByTestId('my-views-content')).toBeInTheDocument();
  });

  // A selection can outlive its permission. Leaving `selectedId` alone left the
  // header naming "My Views" while the content had already fallen back.
  it('falls back to the default when the selected item loses its permission', async () => {
    const { getByTestId, setHeaderInfo, rerender } = renderNav();

    getByTestId('nav-My-Views').click();
    expect(await screen.findByTestId('my-views-content')).toBeInTheDocument();
    expect(lastHeaderTitle(setHeaderInfo)).toBe('My Views');

    canViewViews = false;
    rerender(<Navigation setHeaderInfo={setHeaderInfo} />);

    expect(screen.queryByTestId('my-views-content')).not.toBeInTheDocument();
    expect(screen.getByTestId('recent-content')).toBeInTheDocument();
    expect(lastHeaderTitle(setHeaderInfo)).toBe('Recents (Global)');
  });
});
