import React from 'react';
import {
  CatalogIcon,
  TachographDigitalIcon,
  DashboardIcon,
  LifeCycleIcon,
  ConfigurationIcon,
  PerformanceIcon,
  ExtensionIcon,
} from '@sistent/sistent';
import ConnectionIcon from '../../../assets/icons/Connection';
import EnvironmentIcon from '../../../assets/icons/Environment';
import ServiceMeshIcon from '../../../assets/icons/ServiceMeshIcon';
import WorkspaceOutlinedIcon from '../../../assets/icons/WorkspaceOutlined';
import PatternIcon from '../../../public/static/img/drawer-icons/pattern_svg';
import {
  CATALOG,
  CONFIGURATION,
  CONNECTION,
  DASHBOARD,
  DESIGN,
  ENVIRONMENT,
  EXTENSIONS,
  LIFECYCLE,
  PERFORMANCE,
  PROFILES,
  SERVICE_MESH,
  WORKSPACE,
} from '../../../constants/navigator';
import { iconSmall } from '../../../css/icons.styles';
import { keys } from '@/utils/permission_constants';

export const drawerIconsStyle = {
  height: '19.36px',
  width: '19.36px',
  fontSize: '1.45rem',
  ...iconSmall,
};

export const getNavigatorComponents = (
  /** @type {ProviderUiAccessControl} */ providerUiAccessControl,
  theme,
) => [
  {
    id: DASHBOARD,
    icon: <DashboardIcon style={drawerIconsStyle} />,
    hovericon: <DashboardIcon isHoverEffect={true} style={drawerIconsStyle} />,
    href: '/',
    title: 'Dashboard',
    show: providerUiAccessControl.isNavigatorComponentEnabled([DASHBOARD]),
    link: true,
    submenu: true,
  },
  {
    id: LIFECYCLE,
    icon: <LifeCycleIcon style={drawerIconsStyle} />,
    hovericon: <LifeCycleIcon isHoverEffect={true} style={drawerIconsStyle} />,
    title: 'Lifecycle',
    link: true,
    href: '/management/connections',
    show: providerUiAccessControl.isNavigatorComponentEnabled([LIFECYCLE]),
    submenu: true,
    children: [
      {
        id: CONNECTION,
        icon: <ConnectionIcon style={{ ...drawerIconsStyle }} />,
        href: '/management/connections',
        title: 'Connections',
        show: providerUiAccessControl.isNavigatorComponentEnabled([LIFECYCLE, CONNECTION]),
        link: true,
        permission: {
          action: keys.VIEW_CONNECTIONS.action,
          subject: keys.VIEW_CONNECTIONS.subject,
        },
      },
      {
        id: ENVIRONMENT,
        icon: <EnvironmentIcon style={{ ...drawerIconsStyle }} />,
        href: '/management/environments',
        title: 'Environments',
        show: providerUiAccessControl.isNavigatorComponentEnabled([LIFECYCLE, ENVIRONMENT]),
        link: true,
        permission: {
          action: keys.VIEW_ENVIRONMENTS.action,
          subject: keys.VIEW_ENVIRONMENTS.subject,
        },
      },
      {
        id: WORKSPACE,
        icon: <WorkspaceOutlinedIcon style={{ ...drawerIconsStyle }} />,
        href: '/management/workspaces',
        title: 'Workspaces',
        show: providerUiAccessControl.isNavigatorComponentEnabled([LIFECYCLE, WORKSPACE]),
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
    hovericon: <ConfigurationIcon isHoverEffect={true} style={drawerIconsStyle} />,
    href: '/configuration/designs',
    title: 'Configuration',
    show: providerUiAccessControl.isNavigatorComponentEnabled([CONFIGURATION]),
    link: true,
    submenu: true,
    children: [
      {
        id: CATALOG,
        icon: (
          <CatalogIcon
            primaryFill={
              window.location.pathname === '/configuration/catalog'
                ? theme.palette.background.constant.white
                : ''
            }
            secondaryFill={
              window.location.pathname === '/configuration/catalog'
                ? theme.palette.background.constant.white
                : ''
            }
            tertiaryFill="transparent"
            style={{ ...drawerIconsStyle }}
          />
        ),
        href: '/configuration/catalog',
        title: 'Catalog',
        show: providerUiAccessControl.isNavigatorComponentEnabled([CONFIGURATION, CATALOG]),
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
        show: providerUiAccessControl.isNavigatorComponentEnabled([CONFIGURATION, DESIGN]),
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
    hovericon: (
      <PerformanceIcon
        isHoverEffect={true}
        style={{ transform: 'scale(1.3)', ...drawerIconsStyle }}
      />
    ),
    href: '/performance',
    title: 'Performance',
    show: providerUiAccessControl.isNavigatorComponentEnabled([PERFORMANCE]),
    link: true,
    submenu: true,
    children: [
      {
        id: PROFILES,
        icon: <TachographDigitalIcon fill={theme.palette.icon.default} />,
        href: '/performance/profiles',
        title: 'Profiles',
        show: providerUiAccessControl.isNavigatorComponentEnabled([PERFORMANCE, PROFILES]),
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
    hovericon: <ExtensionIcon isHoverEffect={true} style={drawerIconsStyle} />,
    title: 'Extensions',
    show: providerUiAccessControl.isNavigatorComponentEnabled([EXTENSIONS]),
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
