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
import {
  useGetMeshModelsQuery,
  useGetComponentsQuery,
  useGetRelationshipsQuery,
  useGetRegistrantsQuery,
} from '@/rtk-query/meshModel';
import { removeDuplicateVersions } from '../Settings/Registry/helper';

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
  zIndex: '1500',
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

const RegistryContentWrapper = ({ selectedView, searchText, selectedItemUUID }) => {
  return (
    <MeshModelComponent
      externalView={selectedView}
      externalSearchText={searchText}
      externalSelectedItemUUID={selectedItemUUID}
    />
  );
};

export const Navigation = ({ setHeaderInfo }) => {
  const theme = useTheme();
  const closeList = useMediaQuery(theme.breakpoints.down('xl'));
  const [open, setOpen] = useState(!closeList);
  const registryContext = useContext(RegistryModalContext);
  const { selectedView, searchText, selectedItemUUID } = registryContext;
  const [selectedId, setSelectedId] = useState(selectedView || MODELS);

  useEffect(() => {
    if (selectedView && selectedView !== selectedId) {
      setSelectedId(selectedView);
    }
  }, [selectedView]);

  const { data: modelsData, isLoading: modelsLoading } = useGetMeshModelsQuery({
    params: { pagesize: 'all' },
  });
  const { data: componentsData, isLoading: componentsLoading } = useGetComponentsQuery({
    params: { pagesize: 'all' },
  });
  const { data: relationshipsData, isLoading: relationshipsLoading } = useGetRelationshipsQuery({
    params: { pagesize: 'all' },
  });
  const { data: registrantsData, isLoading: registrantsLoading } = useGetRegistrantsQuery({
    params: { pagesize: 'all' },
  });
  const counts = {
    models: modelsData ? removeDuplicateVersions(modelsData.models || []).length : 0,
    components: componentsData?.total_count || 0,
    relationships: relationshipsData?.total_count || 0,
    registrants: registrantsData?.total_count || 0,
  };
  const navConfig = getNavItems(theme, counts);

  const handleDrawerToggle = () => {
    setOpen(!open);
  };

  const updateHeaderInfo = (id) => {
    const item = navConfig.find((nav) => nav.id === id);
    if (item) {
      setHeaderInfo({
        title: `Registry - ${id}`,
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
          searchText={searchText}
          selectedItemUUID={selectedItemUUID}
          counts={counts}
          isLoading={
            modelsLoading || componentsLoading || relationshipsLoading || registrantsLoading
          }
        />
      </StyledMainContent>
    </Box>
  );
};

const RegistryModal = () => {
  const theme = useTheme();
  const isSmallScreen = useMediaQuery(theme.breakpoints.down('md'));
  const registryContext = useContext(RegistryModalContext);
  const [headerInfo, setHeaderInfo] = useState({
    title: 'Registry',
    icon: <FileIcon {...iconMedium} fill={theme.palette.icon.default} />,
  });

  return (
    <StyledModal
      closeModal={registryContext.closeModal}
      open={registryContext.open}
      headerIcon={headerInfo.icon}
      title={headerInfo.title}
      isFullScreenModeAllowed={!isSmallScreen}
    >
      <ModalBody style={{ height: '100%', padding: '0' }}>
        {registryContext.open && <Navigation setHeaderInfo={setHeaderInfo} />}
      </ModalBody>
    </StyledModal>
  );
};

export default RegistryModal;
