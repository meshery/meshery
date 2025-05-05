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
} from '@layer5/sistent';
import { WorkspacesComponent } from '../Lifecycle';
import { iconMedium } from 'css/icons.styles';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import AccessTimeFilledIcon from '@mui/icons-material/AccessTimeFilled';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import MyViewsContent from './MyViewsContent';
import MyDesignsContent from './MyDesignsContent';
import RecentContent from './RecentContent';
import { useGetWorkspacesQuery } from '../../rtk-query/workspace';
import { useLegacySelector } from 'lib/store';
import { DrawerHeader, StyledDrawer } from './styles';
import { WorkspaceSwitcherContext } from './WorkspaceSwitcher';
import WorkspaceContent from './WorkspaceContent';

const navConfig = {
  mainItems: [
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
        <DesignIcon fill="white" secondaryFill="white" width="20" height="20" primaryFill="white" />
      ),
      content: <MyDesignsContent />,
    },
    {
      id: 'My-Views',
      label: 'My Views',
      icon: <ViewIcon height="24" width="24" fill="white" />,
      content: <MyViewsContent />,
    },
  ],
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
  const [isExpanded, setIsExpanded] = useState(true);

  const handleWorkspacesClick = () => {
    onSelect('All Workspaces');
    setIsExpanded(!isExpanded);
  };

  const workspaces = workspacesData?.workspaces?.map((workspace) => ({
    id: workspace.id,
    name: workspace.name,
    icon: <WorkspaceIcon fill="white" secondaryFill="white" width="20" height="20" />,
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
            <WorkspaceIcon fill="white" secondaryFill="white" />
          </ListItemIcon>
          <ListItemText primary="All Workspaces" sx={{ opacity: open ? 1 : 0 }} />
          {open && workspaces && workspaces.length > 0 && (
            <Box component="span" sx={{ color: 'white' }}>
              {isExpanded ? <ExpandLess /> : <ExpandMore />}
            </Box>
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
              <ListItem key={workspace.id} disablePadding sx={{ display: 'block' }}>
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

const getContentById = (id, workspacesData) => {
  const workspaceSwitcherContext = useContext(WorkspaceSwitcherContext);

  useEffect(() => {
    if (id === 'All Workspaces') {
      workspaceSwitcherContext.setSelectedWorkspace({
        id: null,
        name: null,
      });
    }
  }, [id, workspaceSwitcherContext, workspacesData]);

  const mainItem = navConfig.mainItems.find((item) => item.id === id);
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

const Navigation = () => {
  const [open, setOpen] = useState(true);
  const workspaceSwitcherContext = useContext(WorkspaceSwitcherContext);
  const { selectedWorkspace } = workspaceSwitcherContext;

  // const [selectedId, setSelectedId] = useState(selectedWorkspace?.id || 'Recent');
  const [selectedId, setSelectedId] = useState('Recent');

  const currentOrganization = useLegacySelector((state) => state.get('organization'));

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

  const handleDrawerToggle = () => {
    setOpen(!open);
  };

  const handleItemSelect = (id) => {
    setSelectedId(id);
  };

  return (
    <Box sx={{ display: 'flex', position: 'relative', height: '70vh' }}>
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
          {navConfig.mainItems.map((item) => (
            <NavItem
              key={item.id}
              item={item}
              open={open}
              selectedId={selectedId}
              onSelect={handleItemSelect}
            />
          ))}

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
      <Box
        component="div"
        sx={{
          display: 'flex',
          width: '100%',
          height: '100%',
          flexDirection: 'column',
          padding: '1rem 2rem',
          overflowY: 'auto',
        }}
      >
        {getContentById(selectedId, workspacesData)}
      </Box>
    </Box>
  );
};

const WorkspaceModal = ({ setWorkspaceModal, workspaceModal }) => {
  const theme = useTheme();

  return (
    <Modal
      closeModal={() => setWorkspaceModal(false)}
      open={workspaceModal}
      maxWidth="xl"
      headerIcon={
        <WorkspaceIcon {...iconMedium} secondaryFill={theme.palette.icon.neutral.default} />
      }
      title="All Workspaces"
    >
      <ModalBody style={{ maxHeight: '80vh', overflowY: 'hidden', padding: '0' }}>
        {workspaceModal && <Navigation />}
      </ModalBody>
    </Modal>
  );
};

export default WorkspaceModal;
