/**
 * Registry modal.
 *
 * Full-screen-capable shell hosting the nested model-registry view
 * (`MeshModelComponent`) with a left-rail navigator over Models, Components,
 * Relationships, and Registrants. Switching the active rail item updates
 * both the header copy and the embedded view.
 *
 * Built atop the shared `Modal` primitive with `disableBodyWrap` so this
 * file can render its own `ModalBody` matching the original 90%/80% layout.
 * Migrated to shared primitives as part of Phase 5.b.6 (#18754).
 */
import { useContext, useState, useEffect, ReactNode, FC } from 'react';
import {
  ModalBody,
  List,
  LeftArrowIcon,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Box,
  FileIcon,
  CustomTooltip,
  Drawer,
  DatabaseIcon,
  DARK_BLUE_GRAY,
  useMediaQuery,
} from '@sistent/sistent';
import { styled, useTheme } from '@/theme';
import { Modal } from '@/components/shared/Modal';
import { ChevronButtonWrapper } from '../general/style';
import ConnectionIcon from '@/assets/icons/Connection';
import ComponentIcon from '@/assets/icons/Component';
import MeshModelComponent from './MeshModelComponent';
import { iconMedium, iconSmall } from 'css/icons.styles';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import { RegistryModalContext } from '@/utils/context/RegistryModalContextProvider';
import {
  useGetMeshModelsQuery,
  useGetComponentsQuery,
  useGetRelationshipsQuery,
  useGetRegistrantsQuery,
} from '@/rtk-query/meshModel';
import { removeDuplicateVersions } from './helper';

const DRAWER_WIDTH = 250;

const StyledDrawer = styled(Drawer, { shouldForwardProp: (prop) => prop !== 'open' })<{
  open?: boolean;
}>(({ theme, open }) => ({
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
}));

const StyledMainContent = styled(Box)(() => ({
  flexGrow: 1,
  height: '100%',
  overflow: 'auto',
  display: 'flex',
  flexDirection: 'column',
}));

// Override Sistent/MUI Dialog defaults so the modal occupies ~90%/80% of the
// viewport in normal mode and collapses to a full-width sheet on smaller
// screens. Mirrors the legacy `StyledModal` behaviour.
const StyledRegistryModal = styled(Modal)(({ theme }) => ({
  zIndex: 1500,
  '& .MuiDialog-paperFullScreen': {
    margin: 0,
  },
  '& .MuiDialog-paperFullWidth': {
    width: '90%',
    height: '80%',
  },
  '& .MuiDialog-paper': {
    maxWidth: '100%',
    [theme.breakpoints.down('md')]: {
      margin: 0,
      width: '100%',
      maxWidth: '100%',
      height: '100%',
      maxHeight: '100%',
    },
  },
}));

const StyledModalBody = styled(ModalBody)(() => ({
  height: '100%',
  padding: 0,
}));

interface NavItemDef {
  id: string;
  label: string;
  icon: ReactNode;
}

interface CountSummary {
  models: number;
  components: number;
  relationships: number;
  registrants: number;
}

const getNavItems = (theme: ReturnType<typeof useTheme>, counts: CountSummary): NavItemDef[] => [
  {
    id: MODELS,
    label: `Models (${counts.models?.toLocaleString() || 0})`,
    icon: <FileIcon {...iconSmall} fill={theme.palette.icon.default} />,
  },
  {
    id: COMPONENTS,
    label: `Components (${counts.components?.toLocaleString() || 0})`,
    icon: <ComponentIcon {...iconSmall} fill={theme.palette.icon.default} />,
  },
  {
    id: RELATIONSHIPS,
    label: `Relationships (${counts.relationships?.toLocaleString() || 0})`,
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
    label: `Registrants (${counts.registrants?.toLocaleString() || 0})`,
    icon: <DatabaseIcon {...iconSmall} fill={theme.palette.icon.default} />,
  },
];

interface NavItemProps {
  item: NavItemDef;
  open: boolean;
  selectedId: string;
  onSelect: (_id: string) => void;
}

const NavItem: FC<NavItemProps> = ({ item, open, selectedId, onSelect }) => (
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

interface RegistryContentProps {
  selectedView: string;
  searchText: string | null;
  selectedItemUUID: string | null;
}

const RegistryContentWrapper: FC<RegistryContentProps> = ({
  selectedView,
  searchText,
  selectedItemUUID,
}) => (
  <MeshModelComponent
    externalView={selectedView}
    externalSearchText={searchText}
    externalSelectedItemUUID={selectedItemUUID}
  />
);

interface HeaderInfo {
  title: string;
  icon: ReactNode;
}

interface NavigationProps {
  setHeaderInfo: (_info: HeaderInfo) => void;
}

export const Navigation: FC<NavigationProps> = ({ setHeaderInfo }) => {
  const theme = useTheme();
  const closeList = useMediaQuery(theme.breakpoints.down('xl'));
  const [open, setOpen] = useState(!closeList);
  const registryContext = useContext(RegistryModalContext);
  const { selectedView, searchText, selectedItemUUID } = registryContext;
  const [selectedId, setSelectedId] = useState<string>(selectedView || MODELS);

  useEffect(() => {
    if (selectedView && selectedView !== selectedId) {
      setSelectedId(selectedView);
    }
  }, [selectedView]);

  const { data: modelsData } = useGetMeshModelsQuery({
    params: { pagesize: 'all' },
  });
  const { data: componentsData } = useGetComponentsQuery({
    params: { pagesize: 'all' },
  });
  const { data: relationshipsData } = useGetRelationshipsQuery({
    params: { pagesize: 'all' },
  });
  const { data: registrantsData } = useGetRegistrantsQuery({
    params: { pagesize: 'all' },
  });
  const counts: CountSummary = {
    models: modelsData ? removeDuplicateVersions(modelsData.models || []).length : 0,
    components: componentsData?.totalCount ?? componentsData?.total_count ?? 0,
    relationships: relationshipsData?.totalCount ?? relationshipsData?.total_count ?? 0,
    registrants: registrantsData?.totalCount ?? registrantsData?.total_count ?? 0,
  };
  const navConfig = getNavItems(theme, counts);

  const handleDrawerToggle = () => setOpen((prev) => !prev);

  const updateHeaderInfo = (id: string) => {
    const item = navConfig.find((nav) => nav.id === id);
    if (item) {
      setHeaderInfo({
        title: `Registry - ${id}`,
        icon: <FileIcon {...iconMedium} fill={theme.palette.common.white} />,
      });
    }
  };

  const handleItemSelect = (id: string) => {
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
      </StyledDrawer>
      <StyledMainContent>
        <RegistryContentWrapper
          selectedView={selectedId}
          searchText={searchText}
          selectedItemUUID={selectedItemUUID}
        />
      </StyledMainContent>
      <ChevronButtonWrapper
        isCollapsed={!open}
        onClick={handleDrawerToggle}
        sx={{
          position: 'absolute',
          bottom: '12%',
          left: open
            ? `${DRAWER_WIDTH}px`
            : {
                xs: `calc(${theme.spacing(7)} + 1px - 1.2rem)`,
                sm: `calc(${theme.spacing(8)} + 1px - 1.2rem)`,
              },
          right: 'auto',
          top: 'auto',
          zIndex: 1400,
        }}
      >
        <LeftArrowIcon
          aria-label="Sidebar collapse toggle"
          style={{
            cursor: 'pointer',
            verticalAlign: 'middle',
          }}
          fill={theme.palette.icon.default}
          stroke={theme.palette.icon.default}
          width="1.2rem"
          height="2.8rem"
        />
      </ChevronButtonWrapper>
    </Box>
  );
};

const RegistryModal: FC = () => {
  const theme = useTheme();
  const isSmallScreen = useMediaQuery(theme.breakpoints.down('md'));
  const registryContext = useContext(RegistryModalContext);
  const [headerInfo, setHeaderInfo] = useState<HeaderInfo>({
    title: 'Registry',
    icon: <FileIcon {...iconMedium} fill={theme.palette.icon.default} />,
  });

  return (
    <StyledRegistryModal
      isOpen={registryContext.open}
      onClose={registryContext.closeModal}
      headerIcon={headerInfo.icon}
      title={headerInfo.title}
      isFullScreenModeAllowed={!isSmallScreen}
      disableBodyWrap
    >
      <StyledModalBody>
        {registryContext.open && <Navigation setHeaderInfo={setHeaderInfo} />}
      </StyledModalBody>
    </StyledRegistryModal>
  );
};

export default RegistryModal;
