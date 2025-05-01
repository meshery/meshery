//@ts-check
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
import ViewsContent from './ViewsContent';
import DesignsContent from './DesignsContent';
import RecentContent from './RecentContent';
import { useGetWorkspacesQuery } from '../../rtk-query/workspace';
import { useLegacySelector } from 'lib/store';
import { DrawerHeader, StyledDrawer } from './styles';
import { WorkspaceSwitcherContext } from './WorkspaceSwitcher';
import WorkspaceDesignContent from './WorkspaceDesignContent';
import WorkspaceViewContent from './WorkspaceViewContent';

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
      icon: <DesignIcon fill="white" secondaryFill="white" width="20" height="20" />,
      content: <DesignsContent filterByAuthor={true} />,
    },
    {
      id: 'My-Views',
      label: 'My Views',
      icon: <ViewIcon height="24" width="24" fill="white" />,
      content: <ViewsContent filterByAuthor={true} />,
    },
    // {
    //   id: 'SharedWithMe',
    //   label: 'Shared With Me',
    //   icon: <PeopleAltIcon />,
    //   content: (
    //     <div>
    //       <h2>Shared With Me</h2>
    //       <p>All the resources that are shared with me</p>
    //     </div>
    //   ),
    // },
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

const NestedNavItem = ({ item, open, selectedId, onSelect, level = 0 }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const hasChildren = item.items && item.items.length > 0;

  const handleClick = () => {
    // Only select the item, don't toggle expansion
    onSelect(item.id);
  };

  const handleToggleExpand = (e) => {
    // Stop event propagation to prevent triggering the parent button's onClick
    e.stopPropagation();
    setIsExpanded(!isExpanded);
  };

  const handleDesignsClick = () => {
    // setFetchDesigns(true);
    onSelect(`${item.id}_designs`);
  };

  const handleViewsClick = () => {
    // setFetchViews(true);
    onSelect(`${item.id}_views`);
  };

  return (
    <>
      <ListItem disablePadding sx={{ display: 'block' }}>
        <ListItemButton
          selected={selectedId === item.id}
          onClick={handleClick}
          sx={{
            minHeight: 48,
            px: 2.5,
            pl: level > 0 ? `${(level + 2) * 16}px` : 2.5, // Increased padding for better indentation
            justifyContent: open ? 'initial' : 'center',
            position: 'relative', // For absolute positioning of the chevron
          }}
        >
          {/* Expansion chevron - only visible when the drawer is open and has children */}
          {open && hasChildren && (
            <Box
              onClick={handleToggleExpand}
              size="small"
              sx={{
                position: 'absolute',
                left: level > 0 ? `${level * 16}px` : 0,
                paddingLeft: '4px',
                paddingTop: '4px',
              }}
            >
              {isExpanded ? <ExpandLess fontSize="small" /> : <ExpandMore fontSize="small" />}
            </Box>
          )}

          <ListItemIcon
            sx={{
              minWidth: 0,
              justifyContent: 'center',
              mr: open ? 3 : 'auto',
            }}
          >
            {item.icon}
          </ListItemIcon>
          <ListItemText primary={item.name} sx={{ opacity: open ? 1 : 0 }} />
        </ListItemButton>
      </ListItem>

      {open && hasChildren && (
        <Collapse in={isExpanded} timeout="auto" unmountOnExit>
          <ListItem key={`${item.id}_designs`} disablePadding sx={{ display: 'block' }}>
            <ListItemButton
              selected={selectedId === `${item.id}_designs`}
              onClick={handleDesignsClick}
              sx={{
                pl: `${(level + 3) * 16}px`, // Consistent indentation
                minHeight: 40,
                justifyContent: open ? 'initial' : 'center',
              }}
            >
              <ListItemIcon
                sx={{
                  minWidth: 0,
                  mr: open ? 2 : 'auto',
                  justifyContent: 'center',
                }}
              >
                <DesignIcon fill="white" secondaryFill="white" width="20" height="20" />
              </ListItemIcon>
              <ListItemText
                primary="Designs"
                sx={{ opacity: open ? 1 : 0 }}
                primaryTypographyProps={{ fontSize: '0.9rem' }}
              />
            </ListItemButton>
          </ListItem>

          <ListItem key={`${item.id}_views`} disablePadding sx={{ display: 'block' }}>
            <ListItemButton
              selected={selectedId === `${item.id}_views`}
              onClick={handleViewsClick}
              sx={{
                pl: `${(level + 3) * 16}px`, // Consistent indentation
                minHeight: 40,
                justifyContent: open ? 'initial' : 'center',
              }}
            >
              <ListItemIcon
                sx={{
                  minWidth: 0,
                  mr: open ? 2 : 'auto',
                  justifyContent: 'center',
                }}
              >
                <ViewIcon height="20" width="20" fill="white" />
              </ListItemIcon>
              <ListItemText
                primary="Views"
                sx={{ opacity: open ? 1 : 0 }}
                primaryTypographyProps={{ fontSize: '0.9rem' }}
              />
            </ListItemButton>
          </ListItem>
        </Collapse>
      )}
    </>
  );
};

// Update WorkspacesSection to receive workspacesData as prop
const WorkspacesSection = ({ open, selectedId, onSelect, workspacesData, isLoading }) => {
  // Always expanded by default - no toggle needed
  const [isExpanded, setIsExpanded] = useState(true);

  const handleWorkspacesClick = () => {
    // Only select workspaces without expanding
    onSelect('All Workspaces');
  };

  const workspaces = workspacesData?.workspaces?.map((workspace) => ({
    id: workspace.id,
    name: workspace.name,
    icon: <WorkspaceIcon fill="white" secondaryFill="white" width="20" height="20" />,
    items: [
      {
        id: `${workspace.id}_designs`,
        name: 'Designs',
        icon: <DesignIcon fill="white" secondaryFill="white" width="20" height="20" />,
      },
      {
        id: `${workspace.id}_views`,
        name: 'Views',
        icon: <ViewIcon height="20" width="20" fill="white" />,
      },
    ],
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
        </ListItemButton>
      </ListItem>

      {open && (
        <Collapse in={isExpanded} timeout="auto" unmountOnExit>
          {isLoading ? (
            <ListItem sx={{ pl: 4 }}>
              <ListItemText primary="Loading..." />
            </ListItem>
          ) : (
            workspaces.map((workspace) => (
              <NestedNavItem
                key={workspace.id}
                item={workspace}
                open={open}
                selectedId={selectedId}
                onSelect={onSelect}
                level={1}
              />
            ))
          )}
        </Collapse>
      )}
    </>
  );
};

const WorkspaceContent = ({ id, workspacesData }) => {
  const workspaceSwitcherContext = useContext(WorkspaceSwitcherContext);

  useEffect(() => {
    if (id !== 'All Workspaces' && workspacesData) {
      const workspace = workspacesData.workspaces?.find((workspace) => workspace.id === id);
      if (workspace) {
        workspaceSwitcherContext.setSelectedWorkspace({
          id: id,
          name: workspace.name,
        });
      }
    } else {
      workspaceSwitcherContext.setSelectedWorkspace({
        id: null,
        name: null,
      });
    }
  }, [id, workspaceSwitcherContext, workspacesData]);

  return <WorkspacesComponent />;
};

// Update getContentById to receive workspacesData as parameter
const getContentById = (id, workspacesData) => {
  const mainItem = navConfig.mainItems.find((item) => item.id === id);
  if (mainItem && mainItem.content) {
    return mainItem.content;
  }

  if (id && (id.endsWith('_designs') || id.endsWith('_views'))) {
    const parts = id.split('_');
    const workspaceId = parts[0];
    const contentType = parts[parts.length - 1];

    if (contentType === 'designs') {
      return <WorkspaceDesignContent workspaceId={workspaceId} />;
      // return <div>Designs</div>;
    } else if (contentType === 'views') {
      return <WorkspaceViewContent workspaceId={workspaceId} />;
      // return <WorkspaceViewContent workspacesData={workspacesData} />;
    }
  }

  return <WorkspaceContent id={id} workspacesData={workspacesData} />;
};

const Navigation = () => {
  const [open, setOpen] = useState(true);
  const workspaceSwitcherContext = useContext(WorkspaceSwitcherContext);
  const { selectedWorkspace } = workspaceSwitcherContext;
  // const [selectedId, setSelectedId] = useState(selectedWorkspace.id || 'Recent');
  const [selectedId, setSelectedId] = useState('Recent');
  const currentOrganization = useLegacySelector((state) => state.get('organization'));

  const { data: workspacesData, isLoading } = useGetWorkspacesQuery(
    {
      page: 0,
      pagesize: 'all',
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
      title="Workspaces"
    >
      <ModalBody style={{ maxHeight: '80vh', overflowY: 'hidden', padding: '0' }}>
        {workspaceModal && <Navigation />}
      </ModalBody>
    </Modal>
  );
};

export default WorkspaceModal;
