/**
 * Drawer + content layout that switches between Recents, Designs, Views,
 * Shared, and individual workspaces inside `WorkspaceFormModal`.
 */
import React, { useContext, useState, useEffect, FC } from 'react';
import {
  AccessTimeFilledIcon,
  Box,
  ChevronLeftIcon,
  ChevronRightIcon,
  DesignIcon,
  Divider,
  ErrorBoundary,
  IconButton,
  List,
  PeopleIcon,
  ViewIcon,
  WorkspaceIcon,
  useHasPermission,
  useMediaQuery,
} from '@sistent/sistent';
import { styled, useTheme } from '@/theme';
import { WorkspacesComponent } from '@/components/lifecycle';
import { iconMedium, iconSmall } from 'css/icons.styles';
import MyViewsContent from './SpacesSwitcher/MyViewsContent';
import MyDesignsContent from './SpacesSwitcher/MyDesignsContent';
import RecentContent from './SpacesSwitcher/RecentContent';
import { useGetWorkspacesQuery } from '@/rtk-query/workspace';
import { DrawerHeader, StyledDrawer, StyledMainContent } from './SpacesSwitcher/styles';
import WorkspaceContent from './SpacesSwitcher/WorkspaceContent';
import { useGetProviderCapabilitiesQuery, useGetSelectedOrganization } from '@/rtk-query/user';
import { isLocalProvider } from '@/utils/provider';
import SharedContent from './SpacesSwitcher/SharedContent';
import { Keys } from '@meshery/schemas/permissions';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import type { Theme } from '@/theme';
import { NavItem, WorkspacesSection, NavConfigItem } from './WorkspaceFormModalSections';

/** Nav item selected on open, and the fallback when a selection loses its gate. */
const DEFAULT_NAV_ID = 'Recents (Global)';

const Layout = styled(Box)({
  display: 'flex',
  position: 'relative',
  height: '100%',
});

export type HeaderInfo = {
  title: string;
  icon: React.ReactNode;
};

/**
 * Single source of truth for the drawer's nav items and their permission gates.
 * Both the drawer (which filters disabled items out) and the content wrapper
 * (which resolves the selected item to its body) read from this hook, so the
 * navigation and the content it renders can never diverge on the permission key
 * they authorize against.
 */
const useNavConfig = (theme: Theme): NavConfigItem[] => {
  const canViewViews = useHasPermission(Keys.KanvasViewViews);
  return [
    {
      id: 'Recents (Global)',
      label: 'Recents (Global)',
      icon: <AccessTimeFilledIcon />,
      content: <RecentContent />,
    },
    {
      id: 'My-Designs',
      label: 'My Designs',
      icon: (
        <DesignIcon
          fill={theme.palette.icon.default}
          secondaryFill={theme.palette.icon.default}
          {...iconSmall}
          primaryFill={theme.palette.icon.default}
        />
      ),
      content: <MyDesignsContent />,
    },
    {
      id: 'My-Views',
      label: 'My Views',
      icon: <ViewIcon {...iconSmall} fill={theme.palette.icon.default} />,
      enabled: canViewViews,
      content: <MyViewsContent />,
    },
    {
      id: 'Shared-With-Me',
      label: 'Shared With Me',
      icon: <PeopleIcon {...iconSmall} />,
      content: <SharedContent />,
    },
  ];
};

type WorkspaceContentWrapperProps = {
  id: string;
  workspacesData?: { workspaces?: Array<{ id: string; name: string }> };
  onSelectWorkspace: (args: { id: string; name: string | null }) => void;
};

const WorkspaceContentWrapper: FC<WorkspaceContentWrapperProps> = ({
  id,
  workspacesData,
  onSelectWorkspace,
}) => {
  const workspaceSwitcherContext = useContext(WorkspaceModalContext);
  const theme = useTheme();
  const navConfig = useNavConfig(theme);

  useEffect(() => {
    if (id === 'All Workspaces') {
      workspaceSwitcherContext.setSelectedWorkspace({
        id: null,
        name: null,
      });
    }
  }, [id, workspacesData]);

  // Match the drawer's own filter: an item the user cannot see must not be
  // renderable by id either, so a stale/deep-linked selection cannot bypass the gate.
  const mainItem = navConfig.find((item) => item.id === id && item.enabled !== false);

  if (mainItem && mainItem.content) {
    return <>{mainItem.content}</>;
  }
  if (id === 'All Workspaces') {
    return <WorkspacesComponent onSelectWorkspace={onSelectWorkspace} />;
  }

  const foundWorkspace = workspacesData?.workspaces?.find((workspace) => workspace.id === id);
  if (foundWorkspace) {
    return <WorkspaceContent workspace={foundWorkspace} />;
  }

  return <RecentContent />;
};

type NavigationProps = {
  setHeaderInfo: (info: HeaderInfo) => void;
};

export const Navigation: FC<NavigationProps> = ({ setHeaderInfo }) => {
  const theme = useTheme();
  const closeList = useMediaQuery(theme.breakpoints.down('xl'));
  const [open, setOpen] = useState(!closeList);
  const { data: capabilitiesData } = useGetProviderCapabilitiesQuery();
  const isLocal = isLocalProvider(capabilitiesData);
  const workspaceSwitcherContext = useContext(WorkspaceModalContext);
  const { selectedWorkspace } = workspaceSwitcherContext;
  const [selectedId, setSelectedId] = useState<string>(selectedWorkspace?.id || DEFAULT_NAV_ID);
  const allNavItems = useNavConfig(theme);
  const navConfig = allNavItems.filter((item) => item.enabled !== false);

  // A selection can outlive its permission (revoked mid-session). The drawer
  // drops the item and the content wrapper falls back, but `selectedId` would
  // otherwise keep the header naming a view that no longer renders. Reconcile
  // to the default instead. Guarded on a boolean so this cannot loop, and only
  // nav ids are considered - a selected workspace id is never in `allNavItems`.
  const selectedNavItemRevoked =
    allNavItems.find((item) => item.id === selectedId)?.enabled === false;
  useEffect(() => {
    if (selectedNavItemRevoked) {
      setSelectedId(DEFAULT_NAV_ID);
    }
  }, [selectedNavItemRevoked]);
  const { selectedOrganization } = useGetSelectedOrganization();
  const { data: workspacesData, isLoading } = useGetWorkspacesQuery(
    {
      page: 0,
      pagesize: 'all',
      order: 'updated_at desc',
      orgId: selectedOrganization?.id,
    },
    {
      skip: !selectedOrganization?.id,
    },
  );
  const onSelectWorkspace = ({ id, name }: { id: string; name: string | null }) => {
    setSelectedId(id);
    workspaceSwitcherContext.setSelectedWorkspace({
      id: id,
      name: name,
    });
  };
  useEffect(() => {
    setOpen(!closeList);
  }, [closeList]);

  const handleDrawerToggle = () => {
    setOpen(!open);
  };

  const updateHeaderInfo = (id: string) => {
    const mainItem = navConfig.find((item) => item.id === id);

    if (mainItem) {
      setHeaderInfo({
        title: mainItem.label,
        icon: mainItem.icon,
      });
    } else if (id === 'All Workspaces') {
      setHeaderInfo({
        title: 'All Workspaces',
        icon: <WorkspaceIcon {...iconMedium} secondaryFill={theme.palette.icon.neutral.default} />,
      });
    } else {
      const foundWorkspace = workspacesData?.workspaces?.find((workspace) => workspace.id === id);
      if (foundWorkspace) {
        setHeaderInfo({
          title: `Workspace "${foundWorkspace.name}"`,
          icon: (
            <WorkspaceIcon {...iconMedium} secondaryFill={theme.palette.icon.neutral.default} />
          ),
        });
      }
    }
  };

  const handleItemSelect = (id: string) => {
    setSelectedId(id);
    updateHeaderInfo(id);
  };

  // Set initial header info on component mount or when data changes
  useEffect(() => {
    updateHeaderInfo(selectedId);
  }, [selectedId, workspacesData, theme]);

  return (
    <Layout>
      <ErrorBoundary>
        <StyledDrawer
          variant="permanent"
          open={open}
          sx={{
            '& .MuiDrawer-paper': {
              position: 'relative',
              height: '100%',
            },
          }}
        >
          <List>
            {!isLocal &&
              navConfig.map((item) => (
                <NavItem
                  key={item.id}
                  item={item}
                  open={open}
                  selectedId={selectedId}
                  onSelect={handleItemSelect}
                />
              ))}
            <Divider sx={{ marginBlock: '0.5rem' }} />
            <WorkspacesSection
              open={open}
              selectedId={selectedId}
              onSelect={handleItemSelect}
              workspacesData={workspacesData}
              isLoading={isLoading}
            />
          </List>

          <DrawerHeader open={open}>
            <IconButton onClick={handleDrawerToggle}>
              {open ? <ChevronLeftIcon /> : <ChevronRightIcon />}
            </IconButton>
          </DrawerHeader>
        </StyledDrawer>
      </ErrorBoundary>
      <ErrorBoundary>
        <StyledMainContent>
          <WorkspaceContentWrapper
            id={selectedId}
            workspacesData={workspacesData}
            onSelectWorkspace={onSelectWorkspace}
          />
        </StyledMainContent>
      </ErrorBoundary>
    </Layout>
  );
};
