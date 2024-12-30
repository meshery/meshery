import React from 'react';
import classNames from 'classnames';
import {
  Divider,
  Drawer,
  List,
  ListItem,
  Grow,
  ListItemIcon,
  ListItemText,
  NoSsr,
  Zoom,
  ButtonGroup,
  IconButton,
  Collapse,
  styled,
  ButtonGroup,
  IconButton
} from '@layer5/sistent';
import RemoveIcon from '@mui/icons-material/Remove';
import HelpIcon from '@mui/icons-material/Help';
import DashboardIcon from '@mui/icons-material/Dashboard';
import OpenInNewIcon from '@mui/icons-material/OpenInNew';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faAngleLeft,
  faCaretDown,
  faExternalLinkAlt,
  faDigitalTachograph,
} from '@fortawesome/free-solid-svg-icons';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import HelpIcon from '@material-ui/icons/Help';
import DashboardIcon from '@material-ui/icons/Dashboard';
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
import OpenInNewIcon from '@material-ui/icons/OpenInNew';
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
import { CatalogIcon, CustomTooltip } from '@layer5/sistent';
import { UsesSistent } from './SistentWrapper';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import dataFetch from '../lib/data-fetch';
import {
  cursorNotAllowed,
  disabledStyle,
  disabledStyleWithOutOpacity,
} from '../css/disableComponent.styles';
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

// Styled Components
const StyledDrawer = styled(Drawer)(({ theme, isCollapsed }) => ({
  width: isCollapsed ? theme.spacing(8) + 4 : 256,
  flexShrink: 0,
  whiteSpace: 'nowrap',
  boxSizing: 'border-box',
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: isCollapsed
      ? theme.transitions.duration.leavingScreen
      : theme.transitions.duration.enteringScreen,
  }),
  '& .MuiDrawer-paper': {
    width: isCollapsed ? theme.spacing(8) + 4 : 256,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: isCollapsed
        ? theme.transitions.duration.leavingScreen
        : theme.transitions.duration.enteringScreen,
    }),
    overflowX: 'hidden',
  },
}));

const MainLogo = styled('img')(({ isCollapsed }) => ({
  marginRight: isCollapsed ? 8 : theme.spacing(1),
  marginTop: theme.spacing(1),
  marginLeft: isCollapsed ? -0.5 * theme.spacing(1) : -theme.spacing(1),
  width: isCollapsed ? 40 : 40,
  height: 40,
  borderRadius: 'unset',
}));

const MainLogoText = styled('img')(({ isCollapsed }) => ({
  marginLeft: isCollapsed ? theme.spacing(1) : theme.spacing(0.5),
  marginTop: theme.spacing(1),
  width: isCollapsed ? 170 : 170,
  borderRadius: 'unset',
  display: isCollapsed ? 'none' : 'block',
}));

const CollapseButtonWrapper = styled('div')(({ theme, isCollapsed }) => ({
  boxShadow:
    '0.5px 0px 0px 0px rgb(0 0 0 / 20%), 1.5px 0px 0px 0px rgb(0 0 0 / 14%), 2.5px 1px 3px 0px rgb(0 0 0 / 12%)',
  borderRadius: '0 5px 5px 0',
  position: 'fixed',
  cursor: 'pointer',
  display: 'flex',
  justifyContent: 'center',
  bottom: '12%',
  left: isCollapsed ? '49px' : '257px',
  zIndex: 1400,
  width: 'auto',
  transition: 'left 225ms',
  backgroundColor: isCollapsed ? '#515b60' : 'transparent',
  color: isCollapsed ? '#ffffff' : 'inherit',
  transform: isCollapsed ? 'rotate(180deg)' : 'none',
  '&:hover, &:focus': {
    opacity: 1,
    background: 'transparent',
  },
}));

class Navigator_ extends React.Component {
  constructor(props) {
    super(props);
    const { meshAdapters } = props;
    this.state = {
      path: '',
      meshAdapters,
      mts: new Date(),
      navigator: ExtensionPointSchemaValidator('navigator')(),
      showHelperButton: false,
      openItems: [],
      hoveredId: null,
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
              release_channel: 'NA',
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

  renderNavigatorExtensions(children, depth) {
    const { isDrawerCollapsed } = this.props;
    const { path } = this.state;
    if (children && children.length > 0) {
      return (
        <List disablePadding>
          {children.map(({ id, icon, href, title, children, show: showc }) => {
            if (typeof showc !== 'undefined' && !showc) {
              return null;
            }
            return (
              <React.Fragment key={id}>
                <ListItem
                  button
                  className={classNames(
                    depth === 1 ? '' : 'nested1',
                    'item',
                    'itemActionable',
                    path === href && 'itemActiveItem',
                    isDrawerCollapsed && 'noPadding',
                  )}
                >
                  {this.extensionPointContent(icon, href, title, isDrawerCollapsed)}
                </ListItem>
                {this.renderNavigatorExtensions(children, depth + 1)}
              </React.Fragment>
            );
          })}
        </List>
      );
    }
  }

  extensionPointContent(icon, href, name, drawerCollapsed) {
    let content = (
      <div className="link" data-cy={name}>
        <CustomTooltip
          title={name}
          placement="right"
          disableFocusListener={!drawerCollapsed}
          disableTouchListener={!drawerCollapsed}
        >
          <ListItemIcon className="listIcon">
            <img
              src={icon}
              className="icon"
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
          </ListItemIcon>
        </CustomTooltip>
        <ListItemText className={drawerCollapsed ? 'isHidden' : 'isDisplayed'} primary={name} />
      </div>
    );

    if (href) {
      content = (
        <Link href={href} passHref>
          <span className="link" onClick={() => this.props.updateExtensionType('navigator')}>
            {content}
          </span>
        </Link>
      );
    }

    return content;
  }

  updatenavigatorComponentsMenus() {
    const { navigatorComponents } = this.state;
    navigatorComponents.forEach((cat, ind) => {
      if (cat.id === LIFECYCLE) {
        cat.children.forEach((catc, ind1) => {
          if (catc.id === SERVICE_MESH) {
            return;
          }
          const icon = this.pickIcon(catc.id, catc.href);
          navigatorComponents[ind].children[ind1].icon = icon;

          const cr = this.fetchChildren(catc.id);
          navigatorComponents[ind].children[ind1].children = cr;
        });
      }

      if (cat.id === 'Configuration') {
        let show = false;
        cat.children?.forEach((ch) => {
          if (ch.id === 'Designs') {
            const idx = this.props.capabilitiesRegistry?.capabilities?.findIndex(
              (cap) => cap.feature === 'persist-meshery-patterns',
            );
            if (idx !== -1) {
              ch.show = true;
              show = true;
            }
          }
        });

        cat.show = show;
      }

      // To Toggle Catalog Extension
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
            navigatorComponents[ind].children[ind1].children?.[0]?.href
          ) {
            const newhref = `${navigatorComponents[ind].children[ind1].children[0].href}`;
            navigatorComponents[ind].children[ind1].link = true;
            navigatorComponents[ind].children[ind1].href = newhref;
          }
        });
      }
    });
  }

  static getDerivedStateFromProps(props, state) {
    const { meshAdapters, meshAdaptersts } = props;
    const path = typeof window !== 'undefined' ? window.location.pathname : '';
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

  pickIcon(aName, href) {
    aName = aName.toLowerCase();
    let image = '/static/img/meshery-logo.png';
    let filter =
      typeof window !== 'undefined' && window.location.pathname === href
        ? 'invert(50%) sepia(78%) saturate(2392%) hue-rotate(160deg) brightness(93%) contrast(101%)'
        : '';
    let logoIcon = <img src={image} className="icon" />;
    if (aName) {
      image = `/static/img/${aName}-light.svg`;
      logoIcon = <img src={image} className="icon" style={{ filter }} />;
    }
    return logoIcon;
  }

  handleTitleClick = () => {
    this.props.router.push('/');
  };

  handleAdapterClick(id, link) {
    const { setAdapter } = this.props;
    setAdapter({ selectedAdapter: id });
    if (id !== -1 && !link) {
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

  getMesheryVersionText() {
    const { build, outdated, release_channel } = this.state.versionDetail;

    // If the version is outdated => specify the build
    if (release_channel === 'edge' && outdated) return `${build}`;
    // If it is not outdated => "edge-latest"
    if (release_channel === 'edge' && !outdated) return `${release_channel}-latest`;
    // If stable => "stable-build"
    if (release_channel === 'stable') return `${release_channel}-${build}`;
    return `${build}`;
  }

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
    const { showHelperButton, navigatorComponents } = this.state;
    this.updatenavigatorComponentsMenus();

    // Define styles using sx or styled components
    return (
      <NoSsr>
        <StyledDrawer variant="permanent" isCollapsed={isDrawerCollapsed}>
          {/* Title */}
          <div
            style={
              !this.state.capabilitiesRegistryObj?.isNavigatorComponentEnabled([DASHBOARD])
                ? cursorNotAllowed
                : {}
            }
          >
            <ListItem
              onClick={this.handleTitleClick}
              sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                height: '30px',
                cursor: 'pointer',
                backgroundColor: '#263238',
                color: '#ffffff',
                position: 'sticky',
                top: 0,
                zIndex: 5,
              }}
            >
              <MainLogo
                src="/static/img/meshery-logo.png"
                isCollapsed={isDrawerCollapsed}
                onClick={this.handleTitleClick}
              />
              <MainLogoText
                src="/static/img/meshery-logo-text.png"
                isCollapsed={isDrawerCollapsed}
                onClick={this.handleTitleClick}
              />
            </ListItem>
          </div>

          {/* Menu */}
          <List disablePadding sx={{ overflow: 'hidden auto', '::-webkit-scrollbar': { display: 'none' } }}>
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
                  <div key={childId} style={!show ? cursorNotAllowed : {}}>
                    <ListItem
                      button={!!link}
                      dense
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
                      sx={{
                        paddingTop: 0.5,
                        paddingBottom: 0.5,
                        color: 'rgba(255, 255, 255, 0.7)',
                        '&.Mui-selected': {
                          color: '#4fc3f7',
                          fill: '#4fc3f7',
                        },
                        '&:hover': {
                          backgroundColor: 'rgb(0, 187, 166, 0.5)',
                        },
                      }}
                    >
                      <Link href={link ? href : ''} passHref>
                        <div data-cy={childId} className="link" style={{ display: 'flex', alignItems: 'center', width: '100%' }}>
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
                              <CustomTooltip
                                title={title}
                                placement="right"
                                TransitionComponent={Zoom}
                              >
                                <ListItemIcon
                                  onClick={() => this.toggleItemCollapse(childId)}
                                  sx={{ marginLeft: '20%', marginBottom: '0.4rem' }}
                                >
                                  {hovericon}
                                </ListItemIcon>
                              </CustomTooltip>
                            ) : (
                              <ListItemIcon>{icon}</ListItemIcon>
                            )}
                          </CustomTooltip>
                          {!isDrawerCollapsed && (
                            <ListItemText primary={title} />
                          )}
                        </div>
                      </Link>
                      {!isDrawerCollapsed && children && (
                        <FontAwesomeIcon
                          icon={faCaretDown}
                          onClick={() => this.toggleItemCollapse(childId)}
                          className={classNames({
                            collapsed: this.state.openItems.includes(childId),
                          })}
                          style={isDrawerCollapsed || !children ? { opacity: 0 } : {}}
                        />
                      )}
                    </ListItem>
                    <Collapse in={this.state.openItems.includes(childId)} sx={{ backgroundColor: '#396679' }}>
                      {this.renderChildren(childId, children, 1)}
                    </Collapse>
                  </div>
                );
              },
            )}
            {this.state.navigator && this.state.navigator.length > 0 && (
              <>
                <Divider sx={{ my: 1 }} />
                {this.renderNavigatorExtensions(this.state.navigator, 1)}
              </>
            )}
            <Divider sx={{ my: 1 }} />
          </List>

          {/* Fixed Sidebar Footer */}
          <div style={{ flexGrow: 1, display: 'flex', flexDirection: 'column', marginTop: 'auto', marginBottom: '0.5rem' }}>
            {/* Collapse Button */}
            <CollapseButtonWrapper isCollapsed={isDrawerCollapsed}>
              <div
                style={
                  this.state?.capabilitiesRegistryObj?.isNavigatorComponentEnabled?.([TOGGLER])
                    ? {}
                    : disabledStyle
                }
                onClick={this.toggleMiniDrawer}
              >
                <FontAwesomeIcon
                  icon={faAngleLeft}
                  fixedWidth
                  size="2x"
                  style={{ margin: '0.75rem 0.2rem ', width: '0.8rem', verticalAlign: 'middle' }}
                  alt="Sidebar collapse toggle icon"
                />
              </div>
            </CollapseButtonWrapper>

            {/* Help Icons */}
            <ButtonGroup
              size="large"
              sx={{
                marginLeft: !isDrawerCollapsed ? '5px' : '4px',
                flexDirection: isDrawerCollapsed ? 'column' : 'row',
                alignItems: 'center',
              }}
            >
              {externlinks.map(({ id, icon, title, href }, index) => (
                <ListItem
                  key={id}
                  sx={{
                    display: isDrawerCollapsed && !showHelperButton ? 'none' : 'block',
                    padding: 0,
                  }}
                >
                  <Grow
                    in={showHelperButton || !isDrawerCollapsed}
                    timeout={{ enter: 600 - index * 200, exit: 100 * index }}
                  >
                    <a href={href} target="_blank" rel="noreferrer" style={{ display: 'flex', alignItems: 'center', color: 'inherit', textDecoration: 'none' }}>
                      <CustomTextTooltip title={title} placement={isDrawerCollapsed ? 'right' : 'top'}>
                        <ListItemIcon sx={{ color: '#fff', opacity: 0.7, transition: 'opacity 200ms linear' }} className="helpIcon">
                          {icon}
                        </ListItemIcon>
                      </CustomTextTooltip>
                    </a>
                  </Grow>
                </ListItem>
              ))}
              <ListItem
                sx={{
                  display: isDrawerCollapsed ? 'block' : 'none',
                  marginLeft: '4px',
                  padding: 0,
                }}
              >
                <CustomTextTooltip title="Help" placement={isDrawerCollapsed ? 'right' : 'top'}>
                  <IconButton
                    onClick={this.toggleSpacing}
                    sx={{
                      height: '1.45rem',
                      marginTop: '-4px',
                      transform: 'translateX(0px)',
                    }}
                  >
                    <HelpIcon sx={{ fontSize: '1.45rem', ...iconSmall, color: '#fff', opacity: 0.7, transition: 'opacity 200ms linear', '&:hover': { opacity: 1 }, '&:focus': { opacity: 1 } }} />
                  </IconButton>
                </CustomTextTooltip>
              </ListItem>
            </ButtonGroup>

            {/* Version Information */}
            <ListItem
              sx={{
                position: 'sticky',
                paddingLeft: 0,
                paddingRight: 0,
                color: '#eeeeee',
                fontSize: '0.75rem',
                textAlign: 'center',
                width: '100%',
              }}
            >
              {isDrawerCollapsed ? (
                <div>{this.state.versionDetail.build}</div>
              ) : (
                <Grow in={!isDrawerCollapsed} timeout={{ enter: 800, exit: 100 }}>
                  <span>
                    {this.getMesheryVersionText()} {'  '}
                    <span style={{ cursor: 'pointer' }}>{this.openReleaseNotesInNew()}</span>
                    {this.versionUpdateMsg()}
                  </span>
                </Grow>
              )}
            </ListItem>
          </div>
        </StyledDrawer>
      </NoSsr>
    );
  }
}

// Utility function to get navigator components
const getNavigatorComponents = (capabilityRegistryObj) => [
  {
    id: DASHBOARD,
    icon: <DashboardIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
    hovericon: <DashboardIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
    href: '/',
    title: 'Dashboard',
    show: capabilityRegistryObj.isNavigatorComponentEnabled([DASHBOARD]),
    link: true,
    submenu: true,
  },
  {
    id: LIFECYCLE,
    icon: <LifecycleIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
    hovericon: <LifecycleHover sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
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
        icon: <ServiceMeshIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
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
    icon: <ConfigurationIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
    hovericon: <ConfigurationHover sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
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
              sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }}
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
        icon: <FilterIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
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
        icon: <PatternIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
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
    icon: <PerformanceIcon sx={{ transform: 'scale(1.3)', height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
    hovericon: <PerformanceHover sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
    href: '/performance',
    title: 'Performance',
    show: capabilityRegistryObj.isNavigatorComponentEnabled([PERFORMANCE]),
    link: true,
    submenu: true,
    children: [
      {
        id: PROFILES,
        icon: <FontAwesomeIcon icon={faDigitalTachograph} style={{ fontSize: 24 }} />,
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
    icon: <ExtensionIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
    hovericon: <ExtensionIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
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
  <FontAwesomeIcon style={{ width: '1.11rem', fontSize: '1.11rem' }} icon={faExternalLinkAlt} transform="shrink-7" />
);

const externlinks = [
  {
    id: 'doc',
    href: 'https://docs.meshery.io',
    title: 'Documentation',
    icon: <DocumentIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
    external_icon: ExternalLinkIcon,
  },
  {
    id: 'community',
    href: 'https://slack.meshery.io',
    title: 'Community',
    icon: (
      <SlackIcon sx={{ height: '1.5rem', width: '1.5rem', fontSize: '1.45rem' }} />
    ),
    external_icon: ExternalLinkIcon,
  },
  {
    id: 'forum',
    href: 'http://discuss.meshery.io',
    title: 'Discussion Forum',
    icon: <ChatIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
    external_icon: ExternalLinkIcon,
  },
  {
    id: 'issues',
    href: 'https://github.com/meshery/meshery/issues/new/choose',
    title: 'Issues',
    icon: <GithubIcon sx={{ height: '1.21rem', width: '1.21rem', fontSize: '1.45rem' }} />,
    external_icon: ExternalLinkIcon,
  },
];

export class Navigator_ extends React.Component {
  // ... [constructor, lifecycle methods, and other methods remain unchanged]

  render() {
    const { isDrawerCollapsed } = this.props;
    const { showHelperButton, navigatorComponents } = this.state;
    this.updatenavigatorComponentsMenus();

    return (
      <NoSsr>
        <StyledDrawer variant="permanent" isCollapsed={isDrawerCollapsed}>
          {/* Title */}
          <div
            style={
              !this.state.capabilitiesRegistryObj?.isNavigatorComponentEnabled([DASHBOARD])
                ? cursorNotAllowed
                : {}
            }
          >
            <ListItem
              onClick={this.handleTitleClick}
              sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                height: '30px',
                cursor: 'pointer',
                backgroundColor: '#263238',
                color: '#ffffff',
                position: 'sticky',
                top: 0,
                zIndex: 5,
              }}
            >
              <MainLogo
                src="/static/img/meshery-logo.png"
                isCollapsed={isDrawerCollapsed}
                onClick={this.handleTitleClick}
              />
              <MainLogoText
                src="/static/img/meshery-logo-text.png"
                isCollapsed={isDrawerCollapsed}
                onClick={this.handleTitleClick}
              />
            </ListItem>
          </div>

          {/* Menu */}
          <List disablePadding sx={{ overflow: 'hidden auto', '::-webkit-scrollbar': { display: 'none' } }}>
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
                  <div key={childId} style={!show ? cursorNotAllowed : {}}>
                    <ListItem
                      button={!!link}
                      dense
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
                      sx={{
                        paddingTop: 0.5,
                        paddingBottom: 0.5,
                        color: 'rgba(255, 255, 255, 0.7)',
                        '&.Mui-selected': {
                          color: '#4fc3f7',
                          fill: '#4fc3f7',
                        },
                        '&:hover': {
                          backgroundColor: 'rgb(0, 187, 166, 0.5)',
                        },
                      }}
                    >
                      <Link href={link ? href : ''} passHref>
                        <div data-cy={childId} className="link" style={{ display: 'flex', alignItems: 'center', width: '100%' }}>
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
                              <CustomTooltip
                                title={title}
                                placement="right"
                                TransitionComponent={Zoom}
                              >
                                <ListItemIcon
                                  onClick={() => this.toggleItemCollapse(childId)}
                                  sx={{ marginLeft: '20%', marginBottom: '0.4rem' }}
                                >
                                  {hovericon}
                                </ListItemIcon>
                              </CustomTooltip>
                            ) : (
                              <ListItemIcon>{icon}</ListItemIcon>
                            )}
                          </CustomTooltip>
                          {!isDrawerCollapsed && (
                            <ListItemText primary={title} />
                          )}
                        </div>
                      </Link>
                      {!isDrawerCollapsed && children && (
                        <FontAwesomeIcon
                          icon={faCaretDown}
                          onClick={() => this.toggleItemCollapse(childId)}
                          className={classNames({
                            collapsed: this.state.openItems.includes(childId),
                          })}
                          style={isDrawerCollapsed || !children ? { opacity: 0 } : {}}
                        />
                      )}
                    </ListItem>
                    <Collapse in={this.state.openItems.includes(childId)} sx={{ backgroundColor: '#396679' }}>
                      {this.renderChildren(childId, children, 1)}
                    </Collapse>
                  </div>
                );
              },
            )}
            {this.state.navigator && this.state.navigator.length > 0 && (
              <>
                <Divider sx={{ my: 1 }} />
                {this.renderNavigatorExtensions(this.state.navigator, 1)}
              </>
            )}
            <Divider sx={{ my: 1 }} />
          </List>

          {/* Fixed Sidebar Footer */}
          <div style={{ flexGrow: 1, display: 'flex', flexDirection: 'column', marginTop: 'auto', marginBottom: '0.5rem' }}>
            {/* Collapse Button */}
            <CollapseButtonWrapper isCollapsed={isDrawerCollapsed}>
              <div
                style={
                  this.state?.capabilitiesRegistryObj?.isNavigatorComponentEnabled?.([TOGGLER])
                    ? {}
                    : disabledStyle
                }
                onClick={this.toggleMiniDrawer}
              >
                <FontAwesomeIcon
                  icon={faAngleLeft}
                  fixedWidth
                  size="2x"
                  style={{ margin: '0.75rem 0.2rem ', width: '0.8rem', verticalAlign: 'middle' }}
                  alt="Sidebar collapse toggle icon"
                />
              </div>
            </CollapseButtonWrapper>

            {/* Help Icons */}
            <ButtonGroup
              size="large"
              sx={{
                marginLeft: !isDrawerCollapsed ? '5px' : '4px',
                flexDirection: isDrawerCollapsed ? 'column' : 'row',
                alignItems: 'center',
              }}
            >
              {externlinks.map(({ id, icon, title, href }, index) => (
                <ListItem
                  key={id}
                  sx={{
                    display: isDrawerCollapsed && !showHelperButton ? 'none' : 'block',
                    padding: 0,
                  }}
                >
                  <Grow
                    in={showHelperButton || !isDrawerCollapsed}
                    timeout={{ enter: 600 - index * 200, exit: 100 * index }}
                  >
                    <a href={href} target="_blank" rel="noreferrer" style={{ display: 'flex', alignItems: 'center', color: 'inherit', textDecoration: 'none' }}>
                      <CustomTextTooltip
                        title={title}
                        placement={isDrawerCollapsed ? 'right' : 'top'}
                      >
                        <ListItemIcon sx={{ color: '#fff', opacity: 0.7, transition: 'opacity 200ms linear' }} className="helpIcon">
                          {icon}
                        </ListItemIcon>
                      </CustomTextTooltip>
                    </a>
                  </Grow>
                </ListItem>
              ))}
              <ListItem
                sx={{
                  display: isDrawerCollapsed ? 'block' : 'none',
                  marginLeft: '4px',
                  padding: 0,
                }}
              >
                <CustomTextTooltip title="Help" placement={isDrawerCollapsed ? 'right' : 'top'}>
                  <IconButton
                    onClick={this.toggleSpacing}
                    sx={{
                      height: '1.45rem',
                      marginTop: '-4px',
                      transform: 'translateX(0px)',
                    }}
                  >
                    <HelpIcon sx={{ fontSize: '1.45rem', ...iconSmall, color: '#fff', opacity: 0.7, transition: 'opacity 200ms linear', '&:hover': { opacity: 1 }, '&:focus': { opacity: 1 } }} />
                  </IconButton>
                </CustomTextTooltip>
              </ListItem>
            </ButtonGroup>

            {/* Version Information */}
            <ListItem
              sx={{
                position: 'sticky',
                paddingLeft: 0,
                paddingRight: 0,
                color: '#eeeeee',
                fontSize: '0.75rem',
                textAlign: 'center',
                width: '100%',
              }}
            >
              {isDrawerCollapsed ? (
                <div>{this.state.versionDetail.build}</div>
              ) : (
                <Grow in={!isDrawerCollapsed} timeout={{ enter: 800, exit: 100 }}>
                  <span>
                    {this.getMesheryVersionText()} {'  '}
                    <span style={{ cursor: 'pointer' }}>{this.openReleaseNotesInNew()}</span>
                    {this.versionUpdateMsg()}
                  </span>
                </Grow>
              )}
            </ListItem>
          </div>
        </StyledDrawer>
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

export const NavigatorWithRedux = connect(mapStateToProps, mapDispatchToProps)(withRouter(Navigator_));

export const Navigator = NavigatorWithRedux;

export default Navigator;
