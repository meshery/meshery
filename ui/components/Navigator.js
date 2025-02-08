import React from 'react';
import NoSsr from '@mui/material/NoSsr';
import RemoveIcon from '@mui/icons-material/Remove';
import Zoom from '@mui/material/Zoom';
import Link from 'next/link';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import HelpIcon from '@mui/icons-material/Help';
import LifecycleIcon from '../public/static/img/drawer-icons/lifecycle_mgmt_svg';
import PerformanceIcon from '../public/static/img/drawer-icons/performance_svg';
import ExtensionIcon from '../public/static/img/drawer-icons/extensions_svg';
import FilterIcon from '../public/static/img/drawer-icons/filter_svg';
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
import OpenInNewIcon from '@mui/icons-material/OpenInNew';
import {
  faAngleLeft,
  faCaretDown,
  faExternalLinkAlt,
  faDigitalTachograph,
} from '@fortawesome/free-solid-svg-icons';
import {
  updatepagetitle,
  updatebetabadge,
  toggleDrawer,
  setAdapter,
  updateCapabilities,
} from '../lib/store';
import {
  CatalogIcon,
  CustomTooltip,
  ListItemIcon,
  Grow,
  ListItem,
  List,
  Collapse,
  Box,
} from '@layer5/sistent';
import { UsesSistent } from './SistentWrapper';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import dataFetch from '../lib/data-fetch';
import { cursorNotAllowed, disabledStyle } from '../css/disableComponent.styles';
import { CapabilitiesRegistry } from '../utils/disabledComponents';
import {
  DESIGN,
  CONFIGURATION,
  DASHBOARD,
  CATALOG,
  FILTER,
  LIFECYCLE,
  SERVICE_MESH,
  PERFORMANCE,
  PROFILES,
  TOGGLER,
  CONNECTION,
  ENVIRONMENT,
  WORKSPACE,
} from '../constants/navigator';
import { iconSmall } from '../css/icons.styles';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { CustomTextTooltip } from './MesheryMeshInterface/PatternService/CustomTextTooltip';
import {
  ChevronIcon,
  ExpandMoreIcon,
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
} from './General/style';
import DashboardIcon from '@/assets/icons/DashboardIcon';

const drawerIconsStyle = { height: '1.21rem', width: '1.21rem', fontSize: '1.45rem', ...iconSmall };
const externalLinkIconStyle = { width: '1.11rem', fontSize: '1.11rem' };

const getNavigatorComponents = (/** @type {CapabilitiesRegistry} */ capabilityRegistryObj) => [
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
        icon: (
          <UsesSistent>
            <CatalogIcon
              primaryFill="#FFFFFF"
              secondaryFill="#FFFFFFb3"
              tertiaryFill="transparent"
              style={{ ...drawerIconsStyle }}
            />
          </UsesSistent>
        ),
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
        id: FILTER,
        icon: <FilterIcon style={{ ...drawerIconsStyle }} />,
        href: '/configuration/filters',
        title: 'Filters',
        show: capabilityRegistryObj.isNavigatorComponentEnabled([CONFIGURATION, FILTER]),
        link: true,
        isBeta: true,
        permission: {
          action: keys.VIEW_FILTERS.action,
          subject: keys.VIEW_FILTERS.subject,
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
        icon: (
          <FontAwesomeIcon icon={faDigitalTachograph} style={{ fontSize: 24, color: 'white' }} />
        ),
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
    id: 'Extensions',
    icon: <ExtensionIcon style={drawerIconsStyle} />,
    hovericon: <ExtensionIcon style={drawerIconsStyle} />,
    title: 'Extensions',
    show: capabilityRegistryObj.isNavigatorComponentEnabled(['Extensions']),
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

const ExternalLinkIcon = (
  <FontAwesomeIcon style={externalLinkIconStyle} icon={faExternalLinkAlt} transform="shrink-7" />
);

const externlinks = [
  {
    id: 'doc',
    href: 'https://docs.meshery.io',
    title: 'Documentation',
    icon: <DocumentIcon style={drawerIconsStyle} />,
    external_icon: ExternalLinkIcon,
  },
  {
    id: 'community',
    href: 'https://slack.meshery.io',
    title: 'Community',
    icon: (
      <SlackIcon
        style={{ ...drawerIconsStyle, height: '1.5rem', width: '1.5rem', marginTop: '' }}
      />
    ),
    external_icon: ExternalLinkIcon,
  },
  {
    id: 'forum',
    href: 'https://meshery.io/community#community-forums',
    title: 'Discussion Forum',
    icon: <ChatIcon style={drawerIconsStyle} />,
    external_icon: ExternalLinkIcon,
  },
  {
    id: 'issues',
    href: 'https://github.com/meshery/meshery/issues/new/choose',
    title: 'Issues',
    icon: <GithubIcon style={drawerIconsStyle} />,
    external_icon: ExternalLinkIcon,
  },
];

class Navigator_ extends React.Component {
  constructor(props) {
    super(props);
    const { meshAdapters } = props;
    this.state = {
      path: '',
      meshAdapters,
      mts: new Date(),

      // ExtensionPointSchemaValidator will return a navigator schema
      // decoder which in turn will return an empty array when there is no content
      // passed into it
      navigator: ExtensionPointSchemaValidator('navigator')(),
      showHelperButton: false,
      openItems: [],
      hoveredId: null,
      /** @type {CapabilitiesRegistry} */
      capabilitiesRegistryObj: null,
      versionDetail: {
        build: '',
        latest: '',
        outdated: false,
        commitsha: '',
        release_channel: 'NA',
      },
      navigatorComponents: [],
    };
  }

  isServiceMeshActive() {
    return this.state.meshAdapters.length > 0;
  }

  componentId = 'navigator';

  componentDidMount() {
    this.fetchCapabilities();
    this.fetchVersionDetails();
  }

  fetchCapabilities() {
    dataFetch(
      '/api/provider/capabilities',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (result) {
          const capabilitiesRegistryObj = new CapabilitiesRegistry(result);
          const navigatorComponents = this.createNavigatorComponents(capabilitiesRegistryObj);

          this.setState({
            navigator: ExtensionPointSchemaValidator('navigator')(result?.extensions?.navigator),
            capabilitiesRegistryObj,
            navigatorComponents,
          });
          this.props.updateCapabilities({ capabilitiesRegistry: result });
        }
      },
      (err) => console.error(err),
    );
  }

  fetchVersionDetails() {
    dataFetch(
      '/api/system/version',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (typeof result !== 'undefined') {
          this.setState({ versionDetail: result });
        } else {
          this.setState({
            versionDetail: {
              build: 'Unknown',
              latest: 'Unknown',
              outdated: false,
              commitsha: 'Unknown',
            },
          });
        }
      },
      (err) => console.error(err),
    );
  }

  createNavigatorComponents(capabilityRegistryObj) {
    return getNavigatorComponents(capabilityRegistryObj);
  }

  /**
   * @param {import("../utils/ExtensionPointSchemaValidator").NavigatorSchema[]} children
   * @param {number} depth
   */
  renderNavigatorExtensions(children, depth) {
    const { isDrawerCollapsed } = this.props;
    const { path } = this.state;
    if (children && children.length > 0) {
      return (
        <NavigatorList disablePadding>
          {children.map(({ id, icon, href, title, children, show: showc }) => {
            if (typeof showc !== 'undefined' && !showc) {
              return '';
            }
            const isActive = path === href;
            return (
              <React.Fragment key={id}>
                <NavigatorListItem
                  button
                  depth={depth}
                  key={id}
                  isDrawerCollapsed={isDrawerCollapsed}
                  isActive={isActive}
                >
                  {this.extensionPointContent(icon, href, title, isDrawerCollapsed)}
                </NavigatorListItem>
                {this.renderNavigatorExtensions(children, depth + 1)}
              </React.Fragment>
            );
          })}
        </NavigatorList>
      );
    }
  }

  extensionPointContent(icon, href, name, drawerCollapsed) {
    let content = (
      <UsesSistent>
        <LinkContainer data-cy={name}>
          <CustomTooltip
            title={name}
            placement="right"
            disableFocusListener={!drawerCollapsed}
            disableTouchListener={!drawerCollapsed}
          >
            <MainListIcon>
              <img
                src={icon}
                style={{
                  width: '20px',
                }}
                onMouseOver={(e) => {
                  e.target.style.transform = 'translate(-20%, -25%)';
                  e.target.style.top = '0';
                  e.target.style.right = '0';
                }}
                onMouseOut={(e) => {
                  e.target.style.transform = 'translate(0, 0)';
                  e.target.style.top = 'auto';
                  e.target.style.right = 'auto';
                }}
              />
            </MainListIcon>
          </CustomTooltip>
          <SideBarText drawerCollapsed={drawerCollapsed}>{name}</SideBarText>
        </LinkContainer>
      </UsesSistent>
    );

    if (href) {
      content = (
        <Link
          href={href}
          style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            width: '100%',
            height: '30px',
          }}
        >
          <Box width="100%" onClick={() => this.props.updateExtensionType('navigator')}>
            {content}
          </Box>
        </Link>
      );
    }

    return content;
  }
  updatenavigatorComponentsMenus() {
    const self = this;
    const { navigatorComponents } = this.state;
    navigatorComponents.forEach((cat, ind) => {
      if (cat.id === LIFECYCLE) {
        cat.children.forEach((catc, ind1) => {
          if (catc.id == SERVICE_MESH) {
            return;
          }
          const icon = self.pickIcon(catc.id, catc.href);
          navigatorComponents[ind].children[ind1].icon = icon;

          const cr = self.fetchChildren(catc.id);
          navigatorComponents[ind].children[ind1].children = cr;
        });
      }

      if (cat.id === 'Configuration') {
        let show = false;
        cat.children?.forEach((ch) => {
          if (ch.id === 'Designs') {
            const idx = self.props.capabilitiesRegistry?.capabilities?.findIndex(
              (cap) => cap.feature === 'persist-meshery-patterns',
            );
            if (idx != -1) {
              ch.show = true;
              show = true;
            }
          }
        });

        cat.show = show;
      }

      //To Toggle Catalog Extension
      if (cat.id === CONFIGURATION) {
        cat.children?.forEach((ch) => {
          if (ch.id === CATALOG) {
            ch.show = this.props.catalogVisibility;
          }
        });
      }
    });
  }
  updateAdaptersLink() {
    const { navigatorComponents } = this.state;
    navigatorComponents.forEach((cat, ind) => {
      if (cat.id === LIFECYCLE) {
        cat.children.forEach((catc, ind1) => {
          if (
            typeof navigatorComponents[ind].children[ind1].children[0] !== 'undefined' &&
            typeof navigatorComponents[ind].children[ind1].children[0].href !== 'undefined'
          ) {
            const val = true;
            const newhref = `${navigatorComponents[ind].children[ind1].children[0].href}`;
            navigatorComponents[ind].children[ind1].link = val;
            navigatorComponents[ind].children[ind1].href = newhref;
          }
        });
      }
    });
  }

  static getDerivedStateFromProps(props, state) {
    const { meshAdapters, meshAdaptersts } = props;
    const path = window.location.pathname;
    const st = {};
    if (meshAdaptersts > state.mts) {
      st.meshAdapters = meshAdapters;
      st.mts = meshAdaptersts;
    }

    const fetchNestedPathAndTitle = (path, title, href, children, isBeta) => {
      if (href === path) {
        props.updatepagetitle({ title });
        props.updatebetabadge({ isBeta });
        return;
      }
      if (children && children.length > 0) {
        children.forEach(({ title, href, children, isBeta }) => {
          fetchNestedPathAndTitle(path, title, href, children, isBeta);
        });
      }
    };

    state.navigatorComponents.forEach(({ title, href, children, isBeta }) => {
      fetchNestedPathAndTitle(path, title, href, children, isBeta);
    });
    st.path = path;
    return st;
  }

  /**
   * @param {String} category
   *
   * Format and return the meshadapters
   *
   * @returns {Array<{id : Number, icon : JSX.Element, href : String, title : String, link : Boolean, show : Boolean}>} children
   */
  fetchChildren(category) {
    const { meshAdapters } = this.state;
    const children = [];
    category = category.toLowerCase();
    meshAdapters.forEach((adapter) => {
      let aName = adapter.name.toLowerCase();
      if (category !== aName) {
        return;
      }
      children.push({
        id: adapter.adapter_location,
        icon: <RemoveIcon />,
        href: `/management?adapter=${adapter.adapter_location}`,
        title: `Management - ${adapter.adapter_location}`,
        link: true,
        show: true,
      });
    });
    return children;
  }

  /**
   * @param {String} aName
   *
   * @returns {JSX.Element} image to display
   */
  pickIcon(aName, href) {
    aName = aName.toLowerCase();
    let image = '/static/img/meshery-logo.png';
    let filter =
      window.location.pathname === href
        ? 'invert(50%) sepia(78%) saturate(2392%) hue-rotate(160deg) brightness(93%) contrast(101%)'
        : '';
    let logoIcon = <img src={image} style={{ width: '20px' }} />;
    if (aName) {
      image = '/static/img/' + aName + '-light.svg';
      logoIcon = <img src={image} style={{ filter: filter, width: '20px' }} />;
    }
    return logoIcon;
  }

  /**
   * Changes the route to "/"
   */
  handleTitleClick = () => {
    this.props.router.push('/');
  };

  /**
   * @param {number} id
   * @param {Boolean link
   *
   * Changes the route to "/management"
   */
  handleAdapterClick(id, link) {
    const { setAdapter } = this.props;
    setAdapter({ selectedAdapter: id });
    if (id != -1 && !link) {
      this.props.router.push('/management');
    }
  }

  toggleMiniDrawer = () => {
    const { toggleDrawer, isDrawerCollapsed } = this.props;
    toggleDrawer({ isDrawerCollapsed: !isDrawerCollapsed });
  };

  toggleSpacing = () => {
    const { showHelperButton } = this.state;
    this.setState({ showHelperButton: !showHelperButton });
  };

  /**
   * @param {number} id
   *
   * Removes id from openitems if present
   * Adds id in openitems if not present already
   */
  toggleItemCollapse = (itemId) => {
    const isItemOpen = this.state.openItems.includes(itemId);
    const activeItems = [...this.state.openItems];
    if (isItemOpen) {
      this.setState({ openItems: activeItems.filter((item) => item !== itemId) });
    } else {
      activeItems.push(itemId);
      this.setState({ openItems: [itemId] });
    }
  };

  /**
   * @param {String} idname
   * @param {Array<{id : Number, icon : JSX.Element, href : String, title : String, link : Boolean, show : Boolean}>} children
   * @param {Number} depth
   *
   * Renders children of the menu
   *
   * @returns {JSX.Element}
   */
  renderChildren(idname, children, depth) {
    const { isDrawerCollapsed } = this.props;
    const { path } = this.state;

    if (idname != LIFECYCLE && children && children.length > 0) {
      return (
        <UsesSistent>
          <List disablePadding>
            {children.map(
              ({
                id: idc,
                title: titlec,
                icon: iconc,
                href: hrefc,
                show: showc,
                link: linkc,
                children: childrenc,
                permission: permissionc,
              }) => {
                if (typeof showc !== 'undefined' && !showc) {
                  return '';
                }
                const isActive = path === hrefc;
                return (
                  <div key={idc}>
                    <NavigatorListItemII
                      button
                      key={idc}
                      depth={depth}
                      isDrawerCollapsed={isDrawerCollapsed}
                      isActive={isActive}
                      onClick={() => {
                        if (linkc && hrefc) {
                          this.props.router.push(hrefc);
                        }
                      }}
                      disabled={permissionc ? !CAN(permissionc.action, permissionc.subject) : false}
                    >
                      {this.linkContent(iconc, titlec, hrefc, false, isDrawerCollapsed)}
                    </NavigatorListItemII>
                    {this.renderChildren(idname, childrenc, depth + 1)}
                  </div>
                );
              },
            )}
          </List>
        </UsesSistent>
      );
    }

    if (idname == LIFECYCLE) {
      if (children && children.length > 0) {
        return (
          <UsesSistent>
            <List disablePadding>
              {children.map(
                ({
                  id: idc,
                  title: titlec,
                  icon: iconc,
                  href: hrefc,
                  show: showc,
                  link: linkc,
                  children: childrenc,
                  permission: permissionc,
                }) => {
                  if (typeof showc !== 'undefined' && !showc) {
                    return '';
                  }
                  const isActive = path === hrefc;
                  return (
                    <div key={idc} style={!showc ? cursorNotAllowed : null}>
                      <NavigatorListItemIII
                        component="a"
                        data-cy={idc}
                        button
                        key={idc}
                        depth={depth}
                        isDrawerCollapsed={isDrawerCollapsed}
                        isActive={isActive}
                        isShow={!showc}
                        onClick={() => {
                          this.handleAdapterClick(idc, linkc);
                          if (linkc && hrefc) {
                            this.props.router.push(hrefc);
                          }
                        }}
                        disabled={
                          permissionc ? !CAN(permissionc.action, permissionc.subject) : false
                        }
                      >
                        {this.linkContent(iconc, titlec, hrefc, false, isDrawerCollapsed)}{' '}
                      </NavigatorListItemIII>
                      {this.renderChildren(idname, childrenc, depth + 1)}
                    </div>
                  );
                },
              )}
            </List>
          </UsesSistent>
        );
      }
      if (children && children.length === 1) {
        this.updateAdaptersLink();
      }
    }
    return '';
  }

  /**
   * @param {JSX.Element} iconc
   * @param {String} titlec
   * @param {String} hrefc
   * @param {Boolean} linkc
   * @param {Boolean} drawerCollapsed
   *
   * @return {JSX.Element} content
   */
  linkContent(iconc, titlec, hrefc, linkc, drawerCollapsed) {
    let linkContent = (
      <UsesSistent>
        <LinkContainer>
          <CustomTooltip
            title={titlec}
            placement="right"
            disableFocusListener={!drawerCollapsed}
            disableHoverListener={!drawerCollapsed}
            disableTouchListener={!drawerCollapsed}
          >
            <MainListIcon>{iconc}</MainListIcon>
          </CustomTooltip>
          <SideBarText drawerCollapsed={drawerCollapsed}>{titlec}</SideBarText>
        </LinkContainer>
      </UsesSistent>
    );

    if (linkc && hrefc) {
      linkContent = <Link href={hrefc}>{linkContent}</Link>;
    }
    return linkContent;
  }

  /**
   * getMesheryVersionText returs a well formatted version text
   *
   * If the meshery is running latest version then and is using "edge" channel
   * then it will just show "edge-latest". However, if the meshery is on edge and
   * is running an outdated version then it will return "edge-$version".
   *
   * If on stable channel, then it will always show "stable-$version"
   */
  getMesheryVersionText() {
    const { build, outdated, release_channel } = this.state.versionDetail;

    // If the version is outdated then no matter what the
    // release channel is, specify the build which gets covered in the default case

    if (release_channel === 'edge' && outdated) return `${build}`;
    //if it is not outdated which means running on latest, return edge-latest

    if (release_channel === 'edge' && !outdated) return `${release_channel}-latest`;

    if (release_channel === 'stable') return `${release_channel}-${build}`;

    return `${build}`;
  }

  /**
   * versionUpdateMsg returns the appropriate message
   * based on the meshery's current running version and latest available
   * version.
   *
   * @returns {React.ReactNode} react component to display
   */
  versionUpdateMsg() {
    const { outdated, latest } = this.state.versionDetail;

    if (outdated)
      return (
        <span style={{ marginLeft: '15px' }}>
          {'Update available '}
          <a
            href={`https://docs.meshery.io/project/releases/${latest}`}
            target="_blank"
            rel="noreferrer"
            style={{ color: 'white' }}
          >
            <CustomTextTooltip
              title={`Newer version of Meshery available: ${latest}`}
              placement="right"
            >
              <OpenInNewIcon style={{ width: '0.85rem', verticalAlign: 'middle' }} />
            </CustomTextTooltip>
          </a>
        </span>
      );

    return <span style={{ marginLeft: '15px' }}>Running latest</span>;
  }

  /**
   * openReleaseNotesInNew returns the appropriate link to the release note
   * based on the meshery's current running channel and version.
   *
   * @returns {React.ReactNode} react component to display
   */
  openReleaseNotesInNew() {
    const { release_channel, build } = this.state.versionDetail;

    if (release_channel === 'edge')
      return (
        <a
          href="https://docs.meshery.io/project/releases"
          target="_blank"
          rel="noreferrer"
          style={{ color: 'white' }}
        >
          <OpenInNewIcon style={{ width: '0.85rem', verticalAlign: 'middle' }} />
        </a>
      );

    return (
      <a
        href={`https://docs.meshery.io/project/releases/${build}`}
        target="_blank"
        rel="noreferrer"
        style={{ color: 'white' }}
      >
        <OpenInNewIcon style={{ width: '0.85rem', verticalAlign: 'middle' }} />
      </a>
    );
  }

  render() {
    const { isDrawerCollapsed } = this.props;
    const { path, showHelperButton, navigatorComponents } = this.state;
    this.updatenavigatorComponentsMenus();

    const Title = (
      <div
        style={
          !this.state.capabilitiesRegistryObj?.isNavigatorComponentEnabled([DASHBOARD])
            ? cursorNotAllowed
            : {}
        }
      >
        <UsesSistent>
          <StyledListItem
            component="a"
            onClick={this.handleTitleClick}
            disableLogo={
              !this.state.capabilitiesRegistryObj?.isNavigatorComponentEnabled([DASHBOARD])
            }
          >
            {isDrawerCollapsed ? (
              <>
                <MainLogoCollapsed
                  src="/static/img/meshery-logo.png"
                  onClick={this.handleTitleClick}
                />
                <MainLogoTextCollapsed
                  src="/static/img/meshery-logo-text.png"
                  onClick={this.handleTitleClick}
                />
              </>
            ) : (
              <>
                <MainLogo src="/static/img/meshery-logo.png" onClick={this.handleTitleClick} />
                <MainLogoText
                  src="/static/img/meshery-logo-text.png"
                  onClick={this.handleTitleClick}
                />
              </>
            )}
          </StyledListItem>
        </UsesSistent>
      </div>
    );
    const Menu = (
      <UsesSistent>
        <HideScrollbar disablePadding>
          {navigatorComponents.map(
            ({
              id: childId,
              title,
              icon,
              href,
              show,
              link,
              children,
              hovericon,
              submenu,
              permission,
            }) => {
              return (
                <RootDiv key={childId}>
                  <SideBarListItem
                    button={!!link}
                    dense
                    key={childId}
                    link={!!link}
                    isActive={path === href}
                    isShow={!show}
                    onClick={() => this.toggleItemCollapse(childId)}
                    onMouseOver={() =>
                      isDrawerCollapsed ? this.setState({ hoveredId: childId }) : null
                    }
                    onMouseLeave={() =>
                      !submenu || !this.state.openItems.includes(childId)
                        ? this.setState({ hoveredId: false })
                        : null
                    }
                    disabled={permission ? !CAN(permission.action, permission.subject) : false}
                  >
                    <Link href={link ? href : ''}>
                      <NavigatorLink data-cy={childId}>
                        <CustomTooltip
                          title={childId}
                          placement="right"
                          disableFocusListener={!isDrawerCollapsed}
                          disableHoverListener={true}
                          disableTouchListener={!isDrawerCollapsed}
                          TransitionComponent={Zoom}
                        >
                          {isDrawerCollapsed &&
                          (this.state.hoveredId === childId ||
                            (this.state.openItems.includes(childId) && submenu)) ? (
                            <div>
                              <CustomTooltip
                                title={title}
                                placement="right"
                                TransitionComponent={Zoom}
                              >
                                <ListItemIcon
                                  onClick={() => this.toggleItemCollapse(childId)}
                                  style={{ marginLeft: '20%', marginBottom: '0.4rem' }}
                                >
                                  {hovericon}
                                </ListItemIcon>
                              </CustomTooltip>
                            </div>
                          ) : (
                            <MainListIcon>{icon}</MainListIcon>
                          )}
                        </CustomTooltip>
                        <SideBarText drawerCollapsed={isDrawerCollapsed}>{title}</SideBarText>
                      </NavigatorLink>
                    </Link>
                    <ExpandMoreIcon
                      icon={faCaretDown}
                      onClick={() => this.toggleItemCollapse(childId)}
                      isCollapsed={this.state.openItems.includes(childId)} // Pass collapsed state
                      isDrawerCollapsed={isDrawerCollapsed} // Pass drawer state
                      hasChildren={!!children}
                    />
                  </SideBarListItem>
                  <Collapse
                    in={this.state.openItems.includes(childId)}
                    style={{ backgroundColor: '#396679', opacity: '100%' }}
                  >
                    {this.renderChildren(childId, children, 1)}
                  </Collapse>
                </RootDiv>
              );
            },
          )}
          {this.state.navigator && this.state.navigator.length ? (
            <React.Fragment>
              <SecondaryDivider />
              {this.renderNavigatorExtensions(this.state.navigator, 1)}
            </React.Fragment>
          ) : null}
          <SecondaryDivider />
        </HideScrollbar>
      </UsesSistent>
    );
    const HelpIcons = (
      <UsesSistent>
        <NavigatorHelpIcons
          isCollapsed={isDrawerCollapsed}
          size="large"
          orientation={isDrawerCollapsed ? 'vertical' : 'horizontal'}
        >
          {externlinks.map(({ id, icon, title, href }, index) => {
            return (
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
                    <CustomTextTooltip
                      title={title}
                      placement={isDrawerCollapsed ? 'right' : 'top'}
                    >
                      <ListIconSide>{icon}</ListIconSide>
                    </CustomTextTooltip>
                  </a>
                </Grow>
              </HelpListItem>
            );
          })}
          <ListItem>
            <CustomTextTooltip title="Help" placement={isDrawerCollapsed ? 'right' : 'top'}>
              <HelpButton isCollapsed={isDrawerCollapsed} onClick={this.toggleSpacing}>
                <HelpIcon
                  style={{
                    fontSize: '1.45rem',
                    ...iconSmall,
                    color: '#fff',
                    opacity: '0.7',
                    transition: 'opacity 200ms linear',
                    '&:hover': {
                      opacity: 1,
                      background: 'transparent',
                    },
                    '&:focus': {
                      opacity: 1,
                      background: 'transparent',
                    },
                  }}
                />
              </HelpButton>
            </CustomTextTooltip>
          </ListItem>
        </NavigatorHelpIcons>
      </UsesSistent>
    );
    const Version = (
      <ListItem
        style={{
          position: 'sticky',
          paddingLeft: 0,
          paddingRight: 0,
          color: '#eeeeee',
          fontSize: '0.75rem',
        }}
      >
        {isDrawerCollapsed ? (
          <div style={{ textAlign: 'center', width: '100%' }}>{this.state.versionDetail.build}</div>
        ) : (
          <Grow
            in={!isDrawerCollapsed}
            timeout={{ enter: 800, exit: 100 }}
            style={{ textAlign: 'center', width: '100%' }}
          >
            <span>
              {this.getMesheryVersionText()} {'  '}
              <span style={{ cursor: 'pointer' }}>{this.openReleaseNotesInNew()}</span>
              {this.versionUpdateMsg()}
            </span>
          </Grow>
        )}
      </ListItem>
    );

    const Chevron = (
      <ChevronButtonWrapper
        isCollapsed={isDrawerCollapsed}
        style={
          this.state?.capabilitiesRegistryObj?.isNavigatorComponentEnabled?.([TOGGLER])
            ? {}
            : cursorNotAllowed
        }
      >
        <div
          style={
            this.state?.capabilitiesRegistryObj?.isNavigatorComponentEnabled?.([TOGGLER])
              ? {}
              : disabledStyle
          }
          onClick={this.toggleMiniDrawer}
        >
          <UsesSistent>
            <ChevronIcon
              icon={faAngleLeft}
              fixedWidth
              size="2x"
              style={{ margin: '0.75rem 0.2rem ', width: '0.8rem', verticalAlign: 'middle' }}
              alt="Sidebar collapse toggle icon"
            />
          </UsesSistent>
        </div>
      </ChevronButtonWrapper>
    );

    return (
      <NoSsr>
        <SidebarDrawer isCollapsed={isDrawerCollapsed} variant="permanent">
          {Title}
          {Menu}
          <FixedSidebarFooter>
            {Chevron}
            {HelpIcons}
            {Version}
          </FixedSidebarFooter>
        </SidebarDrawer>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = (dispatch) => ({
  updatepagetitle: bindActionCreators(updatepagetitle, dispatch),
  updatebetabadge: bindActionCreators(updatebetabadge, dispatch),
  toggleDrawer: bindActionCreators(toggleDrawer, dispatch),
  setAdapter: bindActionCreators(setAdapter, dispatch),
  updateCapabilities: bindActionCreators(updateCapabilities, dispatch),
});

const mapStateToProps = (state) => {
  const meshAdapters = state.get('meshAdapters').toJS();
  const meshAdaptersts = state.get('meshAdaptersts');
  const path = state.get('page').get('path');
  const isDrawerCollapsed = state.get('isDrawerCollapsed');
  const capabilitiesRegistry = state.get('capabilitiesRegistry');
  const organization = state.get('organization');
  const keys = state.get('keys');
  const catalogVisibility = state.get('catalogVisibility');
  return {
    meshAdapters,
    meshAdaptersts,
    path,
    isDrawerCollapsed,
    capabilitiesRegistry,
    organization,
    keys,
    catalogVisibility,
  };
};

export const NavigatorWithRedux = connect(
  mapStateToProps,
  mapDispatchToProps,
)(withRouter(Navigator_));

export const Navigator = NavigatorWithRedux;

export default Navigator;
