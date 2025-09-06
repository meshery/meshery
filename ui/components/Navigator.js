import React, { useEffect, useState, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useRouter } from 'next/router';
import Link from 'next/link';
import { useMediaQuery } from '@mui/material';

// Material-UI and Sistent components
import {
  Collapse,
  Grow,
  List,
  ListItem,
  NoSsr,
  Zoom,
  useTheme,
  ListItemIcon,
  CustomTooltip,
  CatalogIcon,
  LeftArrowIcon,
  TachographDigitalIcon,
} from '@sistent/sistent';

// Icons
import RemoveIcon from '@mui/icons-material/Remove';
import HelpIcon from '@mui/icons-material/Help';
import LifecycleIcon from '../public/static/img/drawer-icons/lifecycle_mgmt_svg';
import PerformanceIcon from '../public/static/img/drawer-icons/performance_svg';
import ExtensionIcon from '../public/static/img/drawer-icons/extensions_svg';
import PatternIcon from '../public/static/img/drawer-icons/pattern_svg';
import LifecycleHover from '../public/static/img/drawer-icons/lifecycle_hover_svg';
import PerformanceHover from '../public/static/img/drawer-icons/performance_hover_svg';
import ConfigurationHover from '../public/static/img/drawer-icons/configuration_hover_svg';
import ConfigurationIcon from '../assets/icons/ConfigurationIcon';
import DocumentIcon from '../assets/icons/DocumentIcon';
import SlackIcon from '../assets/icons/SlackIcon';
import GithubIcon from '../assets/icons/GithubIcon';
import ChatIcon from '../assets/icons/ChatIcon';
import ServiceMeshIcon from '../assets/icons/ServiceMeshIcon';
import DashboardIcon from '@/assets/icons/DashboardIcon';

// Redux Actions and Selectors
import {
  toggleDrawer,
  updateBetaBadge,
  updateCapabilities,
  updateTitle,
} from '@/store/slices/mesheryUi';
import { setAdapter } from '@/store/slices/adapter';
import { getProviderCapabilities, getSystemVersion } from '@/rtk-query/user';

// Utils and Constants
import { CapabilitiesRegistry } from '../utils/disabledComponents';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import {
  DESIGN,
  CONFIGURATION,
  DASHBOARD,
  CATALOG,
  LIFECYCLE,
  SERVICE_MESH,
  PERFORMANCE,
  PROFILES,
  TOGGLER,
  CONNECTION,
  ENVIRONMENT,
  WORKSPACE,
  EXTENSIONS,
} from '../constants/navigator';

// Styles and Styled Components
import { iconSmall } from '../css/icons.styles';
import { cursorNotAllowed, disabledStyle } from '../css/disableComponent.styles';
import { CustomTextTooltip } from './MesheryMeshInterface/PatternService/CustomTextTooltip';
import {
  HideScrollbar,
  LinkContainer,
  ListIconSide,
  MainListIcon,
  MainLogo,
  MainLogoCollapsed,
  MainLogoText,
  MainLogoTextCollapsed,
  NavigatorList,
  NavigatorListItem,
  NavigatorListItemII,
  NavigatorListItemIII,
  RootDiv,
  SecondaryDivider,
  SideBarListItem,
  SideBarText,
  StyledListItem,
  NavigatorLink,
  NavigatorHelpIcons,
  HelpListItem,
  HelpButton,
  ChevronButtonWrapper,
  FixedSidebarFooter,
  SidebarDrawer,
  ExpandMore,
} from './General/style';

// --- Constants moved to the top for clarity ---

const drawerIconsStyle = { height: '19.36px', width: '19.36px', fontSize: '1.45rem', ...iconSmall };

const getBaseNavigatorComponents = (capabilityRegistryObj) => [
  {
    id: DASHBOARD,
    icon: <DashboardIcon style={drawerIconsStyle} />,
    hovericon: <DashboardIcon style={drawerIconsStyle} />,
    href: '/',
    title: 'Dashboard',
    show: capabilityRegistryObj.isNavigatorComponentEnabled([DASHBOARD]),
    link: true,
    submenu: true,
  },
  {
    id: LIFECYCLE,
    icon: <LifecycleIcon style={drawerIconsStyle} />,
    hovericon: <LifecycleHover style={drawerIconsStyle} />,
    title: 'Lifecycle',
    link: true,
    href: '/management/connections',
    show: capabilityRegistryObj.isNavigatorComponentEnabled([LIFECYCLE]),
    submenu: true,
    children: [
      {
        id: CONNECTION,
        href: '/management/connections',
        title: 'Connections',
        show: capabilityRegistryObj.isNavigatorComponentEnabled([LIFECYCLE, CONNECTION]),
        link: true,
        permission: {
          action: keys.VIEW_CONNECTIONS.action,
          subject: keys.VIEW_CONNECTIONS.subject,
        },
      },
      {
        id: ENVIRONMENT,
        href: '/management/environments',
        title: 'Environments',
        show: capabilityRegistryObj.isNavigatorComponentEnabled([LIFECYCLE, ENVIRONMENT]),
        link: true,
        permission: {
          action: keys.VIEW_ENVIRONMENTS.action,
          subject: keys.VIEW_ENVIRONMENTS.subject,
        },
      },
      {
        id: WORKSPACE,
        href: '/management/workspaces',
        title: 'Workspaces',
        show: capabilityRegistryObj.isNavigatorComponentEnabled([LIFECYCLE, WORKSPACE]),
        link: true,
        permission: {
          action: keys.VIEW_WORKSPACE.action,
          subject: keys.VIEW_WORKSPACE.subject,
        },
      },
      {
        id: SERVICE_MESH,
        href: '/management/adapter',
        title: 'Adapters',
        link: true,
        icon: <ServiceMeshIcon style={{ ...drawerIconsStyle }} />,
        show: true,
        permission: {
          action: keys.VIEW_CLOUD_NATIVE_INFRASTRUCTURE.action,
          subject: keys.VIEW_CLOUD_NATIVE_INFRASTRUCTURE.subject,
        },
      },
    ],
  },
  {
    id: CONFIGURATION,
    icon: <ConfigurationIcon {...drawerIconsStyle} />,
    hovericon: <ConfigurationHover style={drawerIconsStyle} />,
    href: '/configuration/designs',
    title: 'Configuration',
    show: capabilityRegistryObj.isNavigatorComponentEnabled([CONFIGURATION]),
    link: true,
    submenu: true,
    children: [
      {
        id: CATALOG,
        icon: <CatalogIcon style={{ ...drawerIconsStyle }} />,
        href: '/configuration/catalog',
        title: 'Catalog',
        show: capabilityRegistryObj.isNavigatorComponentEnabled([CONFIGURATION, CATALOG]),
        link: true,
        isBeta: true,
        permission: {
          action: keys.VIEW_CATALOG.action,
          subject: keys.VIEW_CATALOG.subject,
        },
      },
      {
        id: DESIGN,
        icon: <PatternIcon style={{ ...drawerIconsStyle }} />,
        href: '/configuration/designs',
        title: 'Designs',
        show: capabilityRegistryObj.isNavigatorComponentEnabled([CONFIGURATION, DESIGN]),
        link: true,
        isBeta: true,
        permission: {
          action: keys.VIEW_DESIGNS.action,
          subject: keys.VIEW_DESIGNS.subject,
        },
      },
    ],
  },
  {
    id: PERFORMANCE,
    icon: <PerformanceIcon style={{ transform: 'scale(1.3)', ...drawerIconsStyle }} />,
    hovericon: <PerformanceHover style={drawerIconsStyle} />,
    href: '/performance',
    title: 'Performance',
    show: capabilityRegistryObj.isNavigatorComponentEnabled([PERFORMANCE]),
    link: true,
    submenu: true,
    children: [
      {
        id: PROFILES,
        icon: <TachographDigitalIcon />,
        href: '/performance/profiles',
        title: 'Profiles',
        show: capabilityRegistryObj.isNavigatorComponentEnabled([PERFORMANCE, PROFILES]),
        link: true,
        permission: {
          action: keys.VIEW_PERFORMANCE_PROFILES.action,
          subject: keys.VIEW_PERFORMANCE_PROFILES.subject,
        },
      },
    ],
  },
  {
    id: EXTENSIONS,
    icon: <ExtensionIcon style={drawerIconsStyle} />,
    hovericon: <ExtensionIcon style={drawerIconsStyle} />,
    title: 'Extensions',
    show: capabilityRegistryObj.isNavigatorComponentEnabled([EXTENSIONS]),
    width: 12,
    link: true,
    href: '/extensions',
    submenu: false,
    permission: {
      action: keys.VIEW_EXTENSIONS.action,
      subject: keys.VIEW_EXTENSIONS.subject,
    },
  },
];

// --- Sub-components for better decomposition ---

const SidebarTitle = ({ isDrawerCollapsed, isDashboardEnabled, handleTitleClick }) => (
  <div style={!isDashboardEnabled ? cursorNotAllowed : {}}>
    <StyledListItem component="a" onClick={handleTitleClick} disableLogo={!isDashboardEnabled}>
      {isDrawerCollapsed ? (
        <>
          <MainLogoCollapsed src="/static/img/meshery-logo.png" onClick={handleTitleClick} />
          <MainLogoTextCollapsed
            src="/static/img/meshery-logo-text.png"
            onClick={handleTitleClick}
          />
        </>
      ) : (
        <>
          <MainLogo src="/static/img/meshery-logo.png" onClick={handleTitleClick} />
          <MainLogoText src="/static/img/meshery-logo-text.png" onClick={handleTitleClick} />
        </>
      )}
    </StyledListItem>
  </div>
);

const HelpSection = ({ isDrawerCollapsed }) => {
  const [showHelperButton, setShowHelperButton] = useState(false);

  const externalLinks = [
    {
      id: 'doc',
      href: 'https://docs.meshery.io',
      title: 'Documentation',
      icon: <DocumentIcon style={drawerIconsStyle} />,
    },
    {
      id: 'community',
      href: 'https://slack.meshery.io',
      title: 'Community',
      icon: <SlackIcon style={{ ...drawerIconsStyle, height: '24px', width: '24px' }} />,
    },
    {
      id: 'forum',
      href: 'https://meshery.io/community#community-forums',
      title: 'Discussion Forum',
      icon: <ChatIcon style={drawerIconsStyle} />,
    },
    {
      id: 'issues',
      href: 'https://github.com/meshery/meshery/issues/new/choose',
      title: 'Issues',
      icon: <GithubIcon style={drawerIconsStyle} />,
    },
  ];

  return (
    <NavigatorHelpIcons
      isCollapsed={isDrawerCollapsed}
      size="large"
      orientation={isDrawerCollapsed ? 'vertical' : 'horizontal'}
    >
      {externalLinks.map(({ id, icon, title, href }, index) => (
        <HelpListItem
          key={id}
          style={isDrawerCollapsed && !showHelperButton ? { display: 'none' } : {}}
        >
          <Grow
            in={showHelperButton || !isDrawerCollapsed}
            timeout={{ enter: 600 - index * 200, exit: 100 * index }}
          >
            <a
              href={href}
              target="_blank"
              rel="noreferrer"
              style={
                isDrawerCollapsed
                  ? { display: 'flex', alignItems: 'center', justifyContent: 'center' }
                  : {}
              }
            >
              <CustomTextTooltip title={title} placement={isDrawerCollapsed ? 'right' : 'top'}>
                <ListIconSide>{icon}</ListIconSide>
              </CustomTextTooltip>
            </a>
          </Grow>
        </HelpListItem>
      ))}
      <ListItem style={{ display: isDrawerCollapsed ? 'inherit' : 'none' }}>
        <CustomTextTooltip title="Help" placement={isDrawerCollapsed ? 'right' : 'top'}>
          <HelpButton
            isCollapsed={isDrawerCollapsed}
            onClick={() => setShowHelperButton((p) => !p)}
          >
            <HelpIcon
              style={{
                fontSize: '1.45rem',
                ...iconSmall,
                color: '#fff',
                opacity: '0.7',
                transition: 'opacity 200ms linear',
                '&:hover': { opacity: 1 },
                '&:focus': { opacity: 1 },
              }}
            />
          </HelpButton>
        </CustomTextTooltip>
      </ListItem>
    </NavigatorHelpIcons>
  );
};

// --- Main Navigator Component ---

const Navigator_ = () => {
  const router = useRouter();
  const dispatch = useDispatch();
  const theme = useTheme();

  // State from Redux
  const { isDrawerCollapsed, capabilitiesRegistry, catalogVisibility } = useSelector(
    (state) => state.ui,
  );
  const { meshAdapters } = useSelector((state) => state.adapter);

  // Local Component State using individual hooks
  const [openItems, setOpenItems] = useState([]);
  const [hoveredId, setHoveredId] = useState(null);
  const [navigatorExtensions, setNavigatorExtensions] = useState([]);
  const [versionDetail, setVersionDetail] = useState({
    build: '',
    latest: '',
    outdated: false,
    commitsha: '',
    release_channel: 'NA',
  });

  const capabilitiesRegistryObj = useMemo(
    () => (capabilitiesRegistry ? new CapabilitiesRegistry(capabilitiesRegistry) : null),
    [capabilitiesRegistry],
  );

  // Fetch initial data on mount
  useEffect(() => {
    const fetchInitialData = async () => {
      const { data: capsData } = await getProviderCapabilities();
      if (capsData) {
        dispatch(updateCapabilities({ capabilitiesRegistry: capsData }));
        const extensions = ExtensionPointSchemaValidator('navigator')(
          capsData?.extensions?.navigator,
        );
        setNavigatorExtensions(extensions);
      }

      const { data: versionData } = await getSystemVersion();
      if (versionData) {
        setVersionDetail(versionData);
      }
    };
    fetchInitialData();
  }, [dispatch]);

  // Effect to update page title based on route changes
  useEffect(() => {
    const findItemByPath = (items, path) => {
      for (const item of items) {
        if (item.href === path) return item;
        if (item.children) {
          const found = findItemByPath(item.children, path);
          if (found) return found;
        }
      }
      return null;
    };

    if (capabilitiesRegistryObj) {
      const baseNavComponents = getBaseNavigatorComponents(capabilitiesRegistryObj);
      const activeItem = findItemByPath(baseNavComponents, router.pathname);
      if (activeItem) {
        dispatch(updateTitle({ title: activeItem.title }));
        dispatch(updateBetaBadge({ isBeta: activeItem.isBeta || false }));
      }
    }
  }, [router.pathname, capabilitiesRegistryObj, dispatch]);

  const handleAdapterClick = (id, link) => {
    dispatch(setAdapter({ selectedAdapter: id }));
    if (id !== -1 && !link) {
      router.push('/management');
    }
  };

  const toggleItemCollapse = (itemId) => {
    setOpenItems((prevOpen) => (prevOpen.includes(itemId) ? [] : [itemId]));
  };

  // IMMUTABLE derived state for navigator components using useMemo
  const navigatorComponents = useMemo(() => {
    if (!capabilitiesRegistryObj) return [];

    const newNavComponents = JSON.parse(
      JSON.stringify(getBaseNavigatorComponents(capabilitiesRegistryObj)),
    );

    const lifecycleCategory = newNavComponents.find((cat) => cat.id === LIFECYCLE);
    if (lifecycleCategory) {
      lifecycleCategory.children.forEach((child) => {
        if (child.id === SERVICE_MESH) return;
        child.children = meshAdapters
          .filter((adapter) => adapter.name.toLowerCase() === child.id.toLowerCase())
          .map((adapter) => ({
            id: adapter.adapter_location,
            icon: <RemoveIcon />,
            href: `/management?adapter=${adapter.adapter_location}`,
            title: `Management - ${adapter.adapter_location}`,
            link: true,
            show: true,
          }));
      });
    }

    const configCategory = newNavComponents.find((cat) => cat.id === CONFIGURATION);
    if (configCategory) {
      const catalogItem = configCategory.children.find((ch) => ch.id === CATALOG);
      if (catalogItem) {
        catalogItem.show = catalogVisibility;
      }
    }

    return newNavComponents;
  }, [capabilitiesRegistryObj, meshAdapters, catalogVisibility]);

  const renderChildren = (idname, children, depth) => {
    if (!children || children.length === 0) return null;
    const ListComponent = idname === LIFECYCLE ? NavigatorListItemIII : NavigatorListItemII;

    return (
      <List disablePadding>
        {children.map(
          ({ id, title, icon, href, show, link, children: grandChildren, permission }) => {
            if (show === false) return null;
            const isActive = router.pathname === href;
            return (
              <div key={id}>
                <ListComponent
                  button
                  data-testid={id}
                  depth={depth}
                  isDrawerCollapsed={isDrawerCollapsed}
                  isActive={isActive}
                  onClick={() => {
                    if (idname === LIFECYCLE) handleAdapterClick(id, link);
                    if (link && href) router.push(href);
                  }}
                  disabled={permission ? !CAN(permission.action, permission.subject) : false}
                >
                  <LinkContainer>
                    <MainListIcon>{icon}</MainListIcon>
                    <SideBarText drawerCollapsed={isDrawerCollapsed}>{title}</SideBarText>
                  </LinkContainer>
                </ListComponent>
                {renderChildren(idname, grandChildren, depth + 1)}
              </div>
            );
          },
        )}
      </List>
    );
  };

  const renderNavigatorExtensions = (children, depth) => {
    if (!children || children.length === 0) return null;
    return (
      <NavigatorList disablePadding>
        {children.map(({ id, icon, href, title, show }) => {
          if (show === false) return null;
          const isActive = router.pathname === href;
          return (
            <NavigatorListItem
              button
              depth={depth}
              key={id}
              isDrawerCollapsed={isDrawerCollapsed}
              isActive={isActive}
            >
              <Link href={href} style={{ textDecoration: 'none', color: 'inherit', width: '100%' }}>
                <LinkContainer>
                  <MainListIcon>
                    <img src={icon} style={{ width: '20px' }} />
                  </MainListIcon>
                  <SideBarText drawerCollapsed={isDrawerCollapsed}>{title}</SideBarText>
                </LinkContainer>
              </Link>
            </NavigatorListItem>
          );
        })}
      </NavigatorList>
    );
  };

  if (!capabilitiesRegistryObj) {
    return null; // or a loading skeleton
  }

  return (
    <NoSsr>
      <SidebarDrawer isCollapsed={isDrawerCollapsed} variant="permanent">
        <SidebarTitle
          isDrawerCollapsed={isDrawerCollapsed}
          isDashboardEnabled={capabilitiesRegistryObj.isNavigatorComponentEnabled([DASHBOARD])}
          handleTitleClick={() => router.push('/')}
        />

        <HideScrollbar disablePadding>
          {navigatorComponents.map(
            ({ id, title, icon, href, show, link, children, hovericon, submenu, permission }) => {
              if (!show) return null;
              return (
                <RootDiv key={id}>
                  <SideBarListItem
                    button={!!link}
                    dense
                    link={!!link}
                    isActive={router.pathname === href}
                    isShow={!show}
                    onClick={() => toggleItemCollapse(id)}
                    onMouseOver={() => (isDrawerCollapsed ? setHoveredId(id) : null)}
                    onMouseLeave={() =>
                      !submenu || !openItems.includes(id) ? setHoveredId(null) : null
                    }
                    disabled={permission ? !CAN(permission.action, permission.subject) : false}
                  >
                    <Link href={link ? href : ''}>
                      <NavigatorLink data-testid={id}>
                        <CustomTooltip
                          title={title}
                          placement="right"
                          disableFocusListener={!isDrawerCollapsed}
                          disableHoverListener={!isDrawerCollapsed}
                          disableTouchListener={!isDrawerCollapsed}
                          TransitionComponent={Zoom}
                        >
                          {isDrawerCollapsed && (hoveredId === id || openItems.includes(id)) ? (
                            <ListItemIcon style={{ marginLeft: '20%', marginBottom: '0.4rem' }}>
                              {hovericon}
                            </ListItemIcon>
                          ) : (
                            <MainListIcon>{icon}</MainListIcon>
                          )}
                        </CustomTooltip>
                        <SideBarText drawerCollapsed={isDrawerCollapsed}>{title}</SideBarText>
                      </NavigatorLink>
                    </Link>
                    <ExpandMore
                      onClick={() => toggleItemCollapse(id)}
                      isCollapsed={openItems.includes(id)}
                      isDrawerCollapsed={isDrawerCollapsed}
                      theme={theme}
                      hasChildren={!!children}
                    />
                  </SideBarListItem>
                  <Collapse
                    in={openItems.includes(id)}
                    style={{ backgroundColor: '#396679', opacity: '100%' }}
                  >
                    {renderChildren(id, children, 1)}
                  </Collapse>
                </RootDiv>
              );
            },
          )}
          {navigatorExtensions && navigatorExtensions.length > 0 && (
            <React.Fragment>
              <SecondaryDivider />
              {renderNavigatorExtensions(navigatorExtensions, 1)}
            </React.Fragment>
          )}
        </HideScrollbar>

        <FixedSidebarFooter>
          <ChevronButtonWrapper
            isCollapsed={isDrawerCollapsed}
            style={
              capabilitiesRegistryObj.isNavigatorComponentEnabled([TOGGLER]) ? {} : cursorNotAllowed
            }
          >
            <div
              style={
                capabilitiesRegistryObj.isNavigatorComponentEnabled([TOGGLER]) ? {} : disabledStyle
              }
              onClick={() => dispatch(toggleDrawer({ isDrawerCollapsed: !isDrawerCollapsed }))}
            >
              <LeftArrowIcon
                style={{ cursor: 'pointer', verticalAlign: 'middle' }}
                fill={theme.palette.icon.default}
                width="1.2rem"
                height="2.8rem"
              />
            </div>
          </ChevronButtonWrapper>
          <HelpSection isDrawerCollapsed={isDrawerCollapsed} />
          <ListItem style={{ color: '#eeeeee', fontSize: '0.75rem' }}>
            {isDrawerCollapsed ? (
              <div style={{ textAlign: 'center', width: '100%' }}>{versionDetail.build}</div>
            ) : (
              <Grow in={!isDrawerCollapsed} timeout={{ enter: 800, exit: 100 }}>
                <span style={{ textAlign: 'center', width: '100%' }}>v{versionDetail.build}</span>
              </Grow>
            )}
          </ListItem>
        </FixedSidebarFooter>
      </SidebarDrawer>
    </NoSsr>
  );
};

// Wrapper to handle mobile view logic
const NavigatorWrapper = () => {
  const isMobile = useMediaQuery('(max-width:599px)');
  const dispatch = useDispatch();
  const { isDrawerCollapsed } = useSelector((state) => state.ui);

  useEffect(() => {
    if (isMobile && !isDrawerCollapsed) {
      dispatch(toggleDrawer({ isDrawerCollapsed: true }));
    }
  }, [isMobile, isDrawerCollapsed, dispatch]);

  return <Navigator_ />;
};

export const Navigator = NavigatorWrapper;
export default Navigator;
