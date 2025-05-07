import React, { useContext, useState, useEffect } from 'react';
import {
  Modal,
  ModalBody,
  useTheme,
  WorkspaceIcon,
  List,
  IconButton,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Box,
  DesignIcon,
  ViewIcon,
  Collapse,
  useMediaQuery,
  Divider,
  ErrorBoundary,
} from '@layer5/sistent';
import { WorkspacesComponent } from '../Lifecycle';
import { iconMedium, iconSmall } from 'css/icons.styles';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import AccessTimeFilledIcon from '@mui/icons-material/AccessTimeFilled';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import MyViewsContent from './MyViewsContent';
import MyDesignsContent from './MyDesignsContent';
import RecentContent from './RecentContent';
import { useGetWorkspacesQuery } from '../../rtk-query/workspace';
import { DrawerHeader, StyledDrawer, StyledMainContent } from './styles';
import { WorkspaceSwitcherContext } from './WorkspaceSwitcher';
import WorkspaceContent from './WorkspaceContent';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';
import PeopleIcon from '@mui/icons-material/People';
import SharedContent from './SharedContent';
import { useCurrentOrganization } from '@/utils/hooks/useCurrentOrganization';

const getNavItem = (theme) => {
  return [
    {
      id: 'Recent',
      label: 'Recent',
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

const NavItem = ({ item, open, selectedId, onSelect }) => {
  return (
    <ListItem disablePadding sx={{ display: 'block' }}>
      <ListItemButton
        selected={selectedId === item.id}
        onClick={() => onSelect(item.id)}
        sx={{
          minHeight: 48,
          px: 2.5,
          justifyContent: open ? 'initial' : 'center',
        }}
      >
        <ListItemIcon
          sx={{
            minWidth: 0,
            justifyContent: 'center',
            mr: open ? 3 : 'auto',
          }}
        >
          {item.icon}
        </ListItemIcon>
        <ListItemText primary={item.label} sx={{ opacity: open ? 1 : 0 }} />
      </ListItemButton>
    </ListItem>
  );
};

const WorkspacesSection = ({ open, selectedId, onSelect, workspacesData, isLoading }) => {
  const theme = useTheme();
  const [isExpanded, setIsExpanded] = useState(true);

  const handleWorkspacesClick = () => {
    onSelect('All Workspaces');
    setIsExpanded(!isExpanded);
  };

  const workspaces = workspacesData?.workspaces?.map((workspace) => ({
    id: workspace.id,
    name: workspace.name,
    icon: (
      <WorkspaceIcon
        fill={theme.palette.icon.default}
        secondaryFill={theme.palette.icon.default}
        {...iconSmall}
      />
    ),
  }));

  return (
    <>
      <ListItem disablePadding sx={{ display: 'block' }}>
        <ListItemButton
          selected={selectedId === 'All Workspaces'}
          onClick={handleWorkspacesClick}
          sx={{
            minHeight: 48,
            px: 2.5,
            justifyContent: open ? 'initial' : 'center',
          }}
        >
          <ListItemIcon
            sx={{
              minWidth: 0,
              justifyContent: 'center',
              mr: open ? 3 : 'auto',
            }}
          >
            <WorkspaceIcon
              fill={theme.palette.icon.default}
              secondaryFill={theme.palette.icon.default}
              {...iconSmall}
            />
          </ListItemIcon>
          <ListItemText primary="All Workspaces" sx={{ opacity: open ? 1 : 0 }} />
          {open && workspaces && workspaces.length > 0 && (
            <Box component="span">{isExpanded ? <ExpandLess /> : <ExpandMore />}</Box>
          )}
        </ListItemButton>
      </ListItem>

      {open && (
        <Collapse in={isExpanded} timeout="auto" unmountOnExit>
          {isLoading ? (
            <ListItem sx={{ pl: 4 }}>
              <ListItemText primary="Loading..." />
            </ListItem>
          ) : (
            workspaces &&
            workspaces.map((workspace) => (
              <ListItem
                key={workspace.id}
                disablePadding
                sx={{ display: 'block', backgroundColor: theme.palette.background.secondary }}
              >
                <ListItemButton
                  selected={selectedId === workspace.id}
                  onClick={() => onSelect(workspace.id)}
                  sx={{
                    minHeight: 48,
                    px: 2.5,
                    pl: '2.5rem',
                    justifyContent: open ? 'initial' : 'center',
                  }}
                >
                  <ListItemIcon
                    sx={{
                      minWidth: 0,
                      justifyContent: 'center',
                      mr: open ? 3 : 'auto',
                    }}
                  >
                    {workspace.icon}
                  </ListItemIcon>
                  <ListItemText primary={workspace.name} sx={{ opacity: open ? 1 : 0 }} />
                </ListItemButton>
              </ListItem>
            ))
          )}
        </Collapse>
      )}
    </>
  );
};

const WorkspaceContentWrapper = ({ id, workspacesData }) => {
  const workspaceSwitcherContext = useContext(WorkspaceSwitcherContext);
  const theme = useTheme();

  useEffect(() => {
    if (id === 'All Workspaces') {
      workspaceSwitcherContext.setSelectedWorkspace({
        id: null,
        name: null,
      });
    }
  }, [id, workspacesData]);

  const navConfig = getNavItem(theme);
  const mainItem = navConfig.find((item) => item.id === id);

  if (mainItem && mainItem.content) {
    return mainItem.content;
  }

  if (id === 'All Workspaces') {
    return <WorkspacesComponent />;
  }

  const foundWorkspace = workspacesData?.workspaces?.find((workspace) => workspace.id === id);
  if (foundWorkspace) {
    return <WorkspaceContent workspace={foundWorkspace} />;
  }

  return <RecentContent />;
};

const Navigation = ({ setHeaderInfo }) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('lg'));
  const [open, setOpen] = useState(!isMobile);
  const { data: capabilitiesData } = useGetProviderCapabilitiesQuery();
  const isLocalProvider = capabilitiesData?.provider_type === 'local';
  const workspaceSwitcherContext = useContext(WorkspaceSwitcherContext);
  const { selectedWorkspace } = workspaceSwitcherContext;
  const [selectedId, setSelectedId] = useState(selectedWorkspace?.id || 'Recent');
  const currentOrganization = useCurrentOrganization();
  const navConfig = getNavItem(theme);

  const { data: workspacesData, isLoading } = useGetWorkspacesQuery(
    {
      page: 0,
      pagesize: 'all',
      order: 'updated_at desc',
      orgId: currentOrganization?.id,
    },
    {
      skip: !currentOrganization?.id,
    },
  );

  useEffect(() => {
    setOpen(!isMobile);
  }, [isMobile]);

  const handleDrawerToggle = () => {
    setOpen(!open);
  };

  const updateHeaderInfo = (id) => {
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

  const handleItemSelect = (id) => {
    setSelectedId(id);
    updateHeaderInfo(id);
  };

  // Set initial header info on component mount or when data changes
  useEffect(() => {
    updateHeaderInfo(selectedId);
  }, [selectedId, workspacesData, theme]);

  return (
    <Box sx={{ display: 'flex', position: 'relative', height: '100%' }}>
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
            {!isLocalProvider &&
              navConfig.map((item) => (
                <NavItem
                  key={item.id}
                  item={item}
                  open={open}
                  selectedId={selectedId}
                  onSelect={handleItemSelect}
                />
              ))}
            <Divider
              sx={{
                marginBlock: '0.5rem',
              }}
            />
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
          <WorkspaceContentWrapper id={selectedId} workspacesData={workspacesData} />
        </StyledMainContent>
      </ErrorBoundary>
    </Box>
  );
};

const WorkspaceModal = ({ setWorkspaceModal, workspaceModal }) => {
  const theme = useTheme();
  const [headerInfo, setHeaderInfo] = useState({
    title: 'All Workspaces',
    icon: <WorkspaceIcon {...iconMedium} secondaryFill={theme.palette.icon.neutral.default} />,
  });

  return (
    <Modal
      closeModal={() => setWorkspaceModal(false)}
      fullScreen
      fullWidth
      sx={{ margin: '5rem 8rem' }}
      open={workspaceModal}
      headerIcon={headerInfo.icon}
      title={headerInfo.title}
    >
      <ModalBody style={{ height: '100%', padding: '0' }}>
        {workspaceModal && <Navigation setHeaderInfo={setHeaderInfo} />}
      </ModalBody>
    </Modal>
  );
};

export default WorkspaceModal;
