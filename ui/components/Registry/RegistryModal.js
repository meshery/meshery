import React, { useContext, useState, useEffect } from 'react';
import {
  ModalBody,
  useTheme,
  List,
  IconButton,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Box,
  FileIcon,
  useMediaQuery,
  CustomTooltip,
  Modal,
  styled,
  Drawer,
  DatabaseIcon,
  DARK_BLUE_GRAY,
} from '@sistent/sistent';
import ConnectionIcon from '@/assets/icons/Connection';
import ComponentIcon from '@/assets/icons/Component';
import MeshModelComponent from '../Settings/Registry/MeshModelComponent';
import { iconMedium, iconSmall } from 'css/icons.styles';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import { RegistryModalContext } from '@/utils/context/RegistryModalContextProvider';

const DRAWER_WIDTH = 250;

const DrawerHeader = styled('div')(({ theme, open }) => ({
  display: 'flex',
  alignItems: 'flex-end',
  justifyContent: open ? 'flex-end' : 'center',
  marginBottom: '1rem',
  height: '100%',
  ...theme.mixins.toolbar,
}));

const StyledDrawer = styled(Drawer, { shouldForwardProp: (prop) => prop !== 'open' })(
  ({ theme, open }) => ({
    width: DRAWER_WIDTH,
    flexShrink: 0,
    whiteSpace: 'nowrap',
    boxSizing: 'border-box',
    '& .MuiDrawer-paper': {
      position: 'relative',
      height: '100%',
    },
    ...(open && {
      width: DRAWER_WIDTH,
      transition: theme.transitions.create('width', {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.enteringScreen,
      }),
      overflowX: 'hidden',
      '& .MuiDrawer-paper': {
        position: 'relative',
        height: '100%',
        width: DRAWER_WIDTH,
        backgroundColor:
          theme.palette.mode == 'light' ? theme.palette.background.paper : DARK_BLUE_GRAY,
        transition: theme.transitions.create('width', {
          easing: theme.transitions.easing.sharp,
          duration: theme.transitions.duration.enteringScreen,
        }),
        overflowX: 'hidden',
      },
    }),
    ...(!open && {
      transition: theme.transitions.create('width', {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
      }),
      overflowX: 'hidden',
      width: `calc(${theme.spacing(7)} + 1px)`,
      [theme.breakpoints.up('sm')]: {
        width: `calc(${theme.spacing(8)} + 1px)`,
      },
      '& .MuiDrawer-paper': {
        position: 'relative',
        height: '100%',
        backgroundColor:
          theme.palette.mode == 'light' ? theme.palette.background.paper : DARK_BLUE_GRAY,
        transition: theme.transitions.create('width', {
          easing: theme.transitions.easing.sharp,
          duration: theme.transitions.duration.leavingScreen,
        }),
        overflowX: 'hidden',
        width: `calc(${theme.spacing(7)} + 1px)`,
        [theme.breakpoints.up('sm')]: {
          width: `calc(${theme.spacing(8)} + 1px)`,
        },
      },
    }),
  }),
);

const StyledMainContent = styled(Box)(() => ({
  flexGrow: 1,
  height: '100%',
  overflow: 'auto',
  display: 'flex',
  flexDirection: 'column',
}));

const StyledModal = styled(Modal)(({ theme }) => ({
  '& .MuiDialog-paperFullScreen': {
    margin: '0',
  },

  '& .MuiDialog-paperFullWidth': {
    width: '90%',
    height: '80%',
  },

  '& .MuiDialog-paper': {
    maxWidth: '100%',

    [theme.breakpoints.down('md')]: {
      margin: '0',
      width: '100%',
      maxWidth: '100%',
      height: '100%',
      maxHeight: '100%',
    },
  },
}));

const getNavItems = (theme, counts) => {
  return [
    {
      id: MODELS,
      label: `Models (${counts.models})`,
      icon: <FileIcon {...iconSmall} fill={theme.palette.icon.default} />,
    },
    {
      id: COMPONENTS,
      label: `Components (${counts.components})`,
      icon: <ComponentIcon {...iconSmall} fill={theme.palette.icon.default} />,
    },
    {
      id: RELATIONSHIPS,
      label: `Relationships (${counts.relationships})`,
      icon: (
        <ConnectionIcon
          {...iconSmall}
          fill={theme.palette.icon.default}
          primaryFill={theme.palette.icon.default}
          secondaryFill={theme.palette.icon.default}
        />
      ),
    },
    {
      id: REGISTRANTS,
      label: `Registrants (${counts.registrants})`,
      icon: <DatabaseIcon {...iconSmall} fill={theme.palette.icon.default} />,
    },
  ];
};

const NavItem = ({ item, open, selectedId, onSelect }) => {
  return (
    <CustomTooltip title={item.label} disableHoverListener={open} placement="right">
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
    </CustomTooltip>
  );
};

const RegistryContentWrapper = ({
  selectedView,
  modelsCount,
  componentsCount,
  relationshipsCount,
  registrantCount,
}) => {
  // Create a mock settingsRouter function that returns the expected structure
  const mockSettingsRouter = () => ({
    selectedSettingsCategory: 'registry',
    selectedTab: selectedView,
    handleChangeSelectedTab: () => {
      // This will be handled by the modal's own state in the parent Navigation component
    },
    handleChangeSettingsCategory: () => {},
    handleChangeSelectedTabCustomCategory: () => {},
  });

  return (
    <MeshModelComponent
      modelsCount={modelsCount}
      componentsCount={componentsCount}
      relationshipsCount={relationshipsCount}
      registrantCount={registrantCount}
      settingsRouter={mockSettingsRouter}
      externalView={selectedView} // Pass the external view
      hideInternalTabs={true} // Hide the internal tab cards
    />
  );
};

export const Navigation = ({
  setHeaderInfo,
  modelsCount,
  componentsCount,
  relationshipsCount,
  registrantCount,
}) => {
  const theme = useTheme();
  const closeList = useMediaQuery(theme.breakpoints.down('xl'));
  const [open, setOpen] = useState(!closeList);
  const registryContext = useContext(RegistryModalContext);
  const { selectedView } = registryContext;
  const [selectedId, setSelectedId] = useState(selectedView || MODELS);

  const counts = {
    models: modelsCount,
    components: componentsCount,
    relationships: relationshipsCount,
    registrants: registrantCount,
  };

  const navConfig = getNavItems(theme, counts);

  const handleDrawerToggle = () => {
    setOpen(!open);
  };

  const updateHeaderInfo = (id) => {
    const item = navConfig.find((nav) => nav.id === id);
    if (item) {
      setHeaderInfo({
        title: `Registry - ${selectedView}`,
        icon: <FileIcon {...iconMedium} fill={theme.palette.common.white} />,
      });
    }
  };

  const handleItemSelect = (id) => {
    setSelectedId(id);
    registryContext.setSelectedView(id);
    updateHeaderInfo(id);
  };

  // Set initial header info on component mount
  useEffect(() => {
    updateHeaderInfo(selectedId);
  }, [selectedId, theme]);

  return (
    <Box sx={{ display: 'flex', position: 'relative', height: '100%', gap: '1rem' }}>
      <StyledDrawer variant="permanent" open={open}>
        <List>
          {navConfig.map((item) => (
            <NavItem
              key={item.id}
              item={item}
              open={open}
              selectedId={selectedId}
              onSelect={handleItemSelect}
            />
          ))}
        </List>

        <DrawerHeader open={open}>
          <IconButton onClick={handleDrawerToggle}>
            {open ? <ChevronLeftIcon /> : <ChevronRightIcon />}
          </IconButton>
        </DrawerHeader>
      </StyledDrawer>
      <StyledMainContent>
        <RegistryContentWrapper
          selectedView={selectedId}
          modelsCount={modelsCount}
          componentsCount={componentsCount}
          relationshipsCount={relationshipsCount}
          registrantCount={registrantCount}
        />
      </StyledMainContent>
    </Box>
  );
};

const RegistryModal = ({
  registryModal,
  closeRegistryModal,
  modelsCount,
  componentsCount,
  relationshipsCount,
  registrantCount,
}) => {
  const theme = useTheme();
  const isSmallScreen = useMediaQuery(theme.breakpoints.down('md'));
  const [headerInfo, setHeaderInfo] = useState({
    title: 'Registry',
    icon: <FileIcon {...iconMedium} fill={theme.palette.icon.default} />,
  });

  return (
    <StyledModal
      closeModal={closeRegistryModal}
      open={registryModal}
      headerIcon={headerInfo.icon}
      title={headerInfo.title}
      isFullScreenModeAllowed={!isSmallScreen}
    >
      <ModalBody style={{ height: '100%', padding: '0' }}>
        {registryModal && (
          <Navigation
            setHeaderInfo={setHeaderInfo}
            modelsCount={modelsCount}
            componentsCount={componentsCount}
            relationshipsCount={relationshipsCount}
            registrantCount={registrantCount}
          />
        )}
      </ModalBody>
    </StyledModal>
  );
};

export default RegistryModal;
