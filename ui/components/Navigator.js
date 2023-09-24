import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { withStyles } from '@material-ui/core/styles';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Grow from '@material-ui/core/Grow';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import NoSsr from '@material-ui/core/NoSsr';
import RemoveIcon from '@material-ui/icons/Remove';
import Zoom from '@material-ui/core/Zoom';
import Link from 'next/link';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import HelpIcon from '@material-ui/icons/Help';
import DashboardIcon from '@material-ui/icons/Dashboard';
import LifecycleIcon from '../public/static/img/drawer-icons/lifecycle_mgmt_svg';
import PerformanceIcon from '../public/static/img/drawer-icons/performance_svg';
import ApplicationIcon from '../public/static/img/drawer-icons/application_svg';
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
import { ButtonGroup, IconButton, Tooltip } from '@layer5/sistent-components';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import dataFetch from '../lib/data-fetch';
import { Collapse } from '@material-ui/core';
import {
  cursorNotAllowed,
  disabledStyle,
  disabledStyleWithOutOpacity,
} from '../css/disableComponent.styles';
import { CapabilitiesRegistry } from '../utils/disabledComponents';
import {
  APPLICATION,
  DESIGN,
  CONFIGURATION,
  DASHBOARD,
  FILTER,
  LIFECYCLE,
  SERVICE_MESH,
  PERFORMANCE,
  PROFILES,
  TOGGLER,
} from '../constants/navigator';
import { iconSmall } from '../css/icons.styles';

const styles = (theme) => ({
  root: {
    '& svg': {
      width: '1.21rem',
      height: '1.21rem',
    },
  },
  categoryHeader: {
    paddingTop: 16,
    paddingBottom: 16,
  },
  categoryHeaderPrimary: { color: theme.palette.common.white },
  item: {
    paddingTop: 4,
    paddingBottom: 4,
    color: 'rgba(255, 255, 255, 0.7)',
    fill: '#fff',
    '&:hover': {
      '& $expandMoreIcon': {
        opacity: 1,
        transition: 'opacity 200ms ease-in',
      },
    },
  },
  itemCategory: {
    backgroundColor: '#263238',
    boxShadow: '0 -1px 0 #404854 inset',
    paddingTop: 16,
    paddingBottom: 16,
  },
  firebase: {
    top: 0,
    position: 'sticky',
    zIndex: 5,
  },
  link: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    width: '100%',
    height: '30px',
  },

  itemActionable: { '&:hover': { backgroundColor: 'rgb(0, 187, 166, 0.5)' } },
  itemActiveItem: {
    color: '#4fc3f7',
    fill: '#4fc3f7',
  },
  itemPrimary: {
    color: 'inherit',
    fontSize: theme.typography.fontSize,
    '&$textDense': { fontSize: theme.typography.fontSize },
  },
  textDense: {},
  divider: {
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
  mainLogo: {
    marginRight: theme.spacing(1),
    marginTop: theme.spacing(1),
    marginLeft: theme.spacing(-1),
    width: 40,
    height: 40,
    borderRadius: 'unset',
  },
  mainLogoText: {
    marginLeft: theme.spacing(0.5),
    marginTop: theme.spacing(1),
    width: 170,
    borderRadius: 'unset',
  },
  mainLogoCollapsed: {
    marginRight: theme.spacing(1),
    marginTop: theme.spacing(1),
    marginLeft: theme.spacing(-0.5),
    width: 40,
    height: 40,
    borderRadius: 'unset',
  },
  mainLogoTextCollapsed: {
    marginLeft: theme.spacing(1),
    marginTop: theme.spacing(1),
    width: 170,
    borderRadius: 'unset',
  },
  settingsIcon: { marginLeft: theme.spacing(2) },
  cursorPointer: { cursor: 'pointer' },
  listIcon: {
    minWidth: theme.spacing(3.5),
    paddingTop: theme.spacing(0.5),
    textAlign: 'center',
    display: 'inline-table',
    paddingRight: theme.spacing(0.5),
    marginLeft: theme.spacing(0.8),
  },
  listIcon1: {
    minWidth: theme.spacing(3.5),
    paddingTop: theme.spacing(0.5),
    textAlign: 'center',
    display: 'inline-table',
    paddingRight: theme.spacing(0.5),
    opacity: 0.5,
  },
  listIconSlack: {
    minWidth: theme.spacing(3.5),
    paddingTop: theme.spacing(0.5),
    textAlign: 'center',
    display: 'inline-table',
    marginLeft: theme.spacing(-0.1),
    paddingRight: theme.spacing(0.5),
    opacity: 0.5,
  },
  nested1: { paddingLeft: theme.spacing(3) },
  nested2: { paddingLeft: theme.spacing(5) },
  icon: { width: theme.spacing(2.5) },
  istioIcon: { width: theme.spacing(1.8) },
  isHidden: {
    opacity: 0,
    transition: 'opacity 200ms ease-in-out',
  },
  isDisplayed: {
    opacity: 1,
    transition: 'opacity 200ms ease-in-out',
  },
  sidebarCollapsed: {
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    overflowX: 'hidden',
    width: theme.spacing(8) + 4,
  },
  sidebarExpanded: {
    width: '256px',
    overflowX: 'hidden',
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  fixedSidebarFooter: {
    display: 'flex',
    flexDirection: 'column',
    marginTop: 'auto',
    marginBottom: '0.5rem',
  },
  collapseButtonWrapper: {
    boxShadow:
      '0.5px 0px 0px 0px rgb(0 0 0 / 20%), 1.5px 0px 0px 0px rgb(0 0 0 / 14%), 2.5px 1px 3px 0px rgb(0 0 0 / 12%)',
    borderRadius: '0 5px 5px 0',
    position: 'fixed',
    cursor: 'pointer',
    display: 'flex',
    justifyContent: 'center',
    bottom: '12%',
    left: '257px',
    zIndex: '1400',
    width: 'auto',
    transition: 'left 195ms',
    '&:hover': {
      opacity: 1,
      background: 'transparent',
    },
    '&:focus': {
      opacity: 1,
      background: 'transparent',
    },
  },
  collapseButtonWrapperRotated: {
    backgroundColor: '#515b60',
    color: '#ffffff',
    position: 'fixed',
    borderRadius: '0 5px 5px 0',
    cursor: 'pointer',
    bottom: '12%',
    left: '49px',
    zIndex: '1400',
    width: 'auto',
    transition: 'left 225ms',
    transform: 'rotate(180deg)',
    display: 'flex',
    justifyContent: 'center',

    '&:hover': { opacity: 1 },
    '&:focus': { opacity: 1 },
  },
  noPadding: {
    paddingLeft: '16px',
    paddingRight: '16px',
  },
  drawerIcons: {
    height: '1.21rem',
    width: '1.21rem',
    fontSize: '1.21rem',
  },
  avatarGroup: { '& .MuiAvatarGroup-avatar': { border: 'none' } },
  marginLeft: {
    padding: '5px',
    '& > li': {
      padding: '0',
    },
  },
  rightMargin: { marginRight: 8 },
  btnGrpMarginRight: {
    marginRight: 4,
    alignItems: 'center',
  },
  helpIcon: {
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
  },
  extraPadding: {
    // paddingTop : 4,
    // paddingBottom : 4,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  restrictPointer: { pointerEvents: 'none' },
  expandMoreIcon: {
    opacity: 0,
    cursor: 'pointer',
    transform: 'translateX(3px)',
    '&:hover': { color: '#4fc3f7' },
  },
  collapsed: { transform: 'rotate(180deg) translateX(-0.8px)' },
  collapsedHelpButton: {
    height: '1.45rem',
    marginTop: '-4px',
    transform: 'translateX(0px)',
  },
  rightTranslate: { transform: 'translateX(0.5px)' },
  hideScrollbar: {
    overflow: 'hidden auto',
    'scrollbar-width': 'none',
    '-ms-overflow-style': 'none',
    '&::-webkit-scrollbar': {
      display: 'none',
    },
  },
  disabled: disabledStyle,
  disableLogo: disabledStyleWithOutOpacity,
  cursorNotAllowed: cursorNotAllowed,
});

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
    submenu: false,
  },
  {
    id: LIFECYCLE,
    icon: <LifecycleIcon style={drawerIconsStyle} />,
    hovericon: <LifecycleHover style={drawerIconsStyle} />,
    title: 'Lifecycle',
    link: true,
    href: '/management',
    show: capabilityRegistryObj.isNavigatorComponentEnabled([LIFECYCLE]),
    submenu: true,
    children: [
      // {
      //   id : CONNECTION,
      //   href : "/management/connections",
      //   title : "Connections",
      //   show : capabilityRegistryObj.isNavigatorComponentEnabled([LIFECYCLE, CONNECTION]),
      //   link : true,
      // },
      {
        id: SERVICE_MESH,
        href: '/management',
        title: 'Service Mesh',
        link: true,
        icon: <ServiceMeshIcon style={{ ...drawerIconsStyle }} />,
        show: true,
      },
    ],
  },
  {
    id: CONFIGURATION,
    icon: <ConfigurationIcon {...drawerIconsStyle} />,
    hovericon: <ConfigurationHover style={drawerIconsStyle} />,
    href: '/configuration/applications',
    title: 'Configuration',
    show: capabilityRegistryObj.isNavigatorComponentEnabled([CONFIGURATION]),
    link: true,
    submenu: true,
    children: [
      {
        id: APPLICATION,
        icon: <ApplicationIcon style={{ ...drawerIconsStyle }} />,
        href: '/configuration/applications',
        title: 'Applications',
        show: capabilityRegistryObj.isNavigatorComponentEnabled([CONFIGURATION, APPLICATION]),
        link: true,
        isBeta: true,
      },
      {
        id: FILTER,
        icon: <FilterIcon style={{ ...drawerIconsStyle }} />,
        href: '/configuration/filters',
        title: 'Filters',
        show: capabilityRegistryObj.isNavigatorComponentEnabled([CONFIGURATION, FILTER]),
        link: true,
        isBeta: true,
      },
      {
        id: DESIGN,
        icon: <PatternIcon style={{ ...drawerIconsStyle }} />,
        href: '/configuration/designs',
        title: 'Designs',
        show: capabilityRegistryObj.isNavigatorComponentEnabled([CONFIGURATION, DESIGN]),
        link: true,
        isBeta: true,
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
        icon: <FontAwesomeIcon icon={faDigitalTachograph} style={{ fontSize: 24 }} />,
        href: '/performance/profiles',
        title: 'Profiles',
        show: capabilityRegistryObj.isNavigatorComponentEnabled([PERFORMANCE, PROFILES]),
        link: true,
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
    href: 'http://discuss.meshery.io',
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

class Navigator extends React.Component {
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
    dataFetch(
      '/api/provider/capabilities',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (result) {
          const capabilitiesRegistryObj = new CapabilitiesRegistry(result);
          const navigatorComponents = getNavigatorComponents(capabilitiesRegistryObj);

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

  /**
   * @param {import("../utils/ExtensionPointSchemaValidator").NavigatorSchema[]} children
   * @param {number} depth
   */
  renderNavigatorExtensions(children, depth) {
    const { classes, isDrawerCollapsed } = this.props;
    const { path } = this.state;
    if (children && children.length > 0) {
      return (
        <List disablePadding>
          {children.map(({ id, icon, href, title, children, show: showc }) => {
            if (typeof showc !== 'undefined' && !showc) {
              return '';
            }
            return (
              <React.Fragment key={id}>
                <ListItem
                  button
                  key={id}
                  className={classNames(
                    depth === 1 ? '' : classes.nested1,
                    classes.item,
                    classes.itemActionable,
                    path === href && classes.itemActiveItem,
                    isDrawerCollapsed && classes.noPadding,
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
    const { classes } = this.props;

    let content = (
      <div className={classNames(classes.link)} data-cy={name}>
        <Tooltip
          title={name}
          placement="right"
          disableFocusListener={!drawerCollapsed}
          disableTouchListener={!drawerCollapsed}
        >
          <ListItemIcon className={classes.listIcon}>
            <img src={icon} className={classes.icon} />
          </ListItemIcon>
        </Tooltip>
        <ListItemText
          className={drawerCollapsed ? classes.isHidden : classes.isDisplayed}
          classes={{ primary: classes.itemPrimary }}
        >
          {name}
        </ListItemText>
      </div>
    );

    if (href) {
      content = (
        <Link href={href}>
          <span
            className={classNames(classes.link)}
            onClick={() => this.props.updateExtensionType('navigator')}
          >
            {content}
          </span>
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
    const { classes } = this.props;
    let image = '/static/img/meshery-logo.png';
    let filter =
      window.location.pathname === href
        ? 'invert(50%) sepia(78%) saturate(2392%) hue-rotate(160deg) brightness(93%) contrast(101%)'
        : '';
    let logoIcon = <img src={image} className={classes.icon} />;
    if (aName) {
      image = '/static/img/' + aName + '-light.svg';
      logoIcon = <img src={image} className={classes.icon} style={{ filter: filter }} />;
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
  toggleItemCollapse(id) {
    const activeItems = [...this.state.openItems];
    if (this.state.openItems.includes(id)) {
      this.setState({ openItems: activeItems.filter((item) => item !== id) });
    } else {
      activeItems.push(id);
      this.setState({ openItems: activeItems });
    }
  }

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
    const { classes, isDrawerCollapsed } = this.props;
    const { path } = this.state;

    if (idname != LIFECYCLE && children && children.length > 0) {
      return (
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
            }) => {
              if (typeof showc !== 'undefined' && !showc) {
                return '';
              }
              return (
                <div key={idc}>
                  <ListItem
                    button
                    key={idc}
                    className={classNames(
                      depth === 1 ? classes.nested1 : classes.nested2,
                      classes.item,
                      classes.itemActionable,
                      path === hrefc && classes.itemActiveItem,
                      isDrawerCollapsed && classes.noPadding,
                    )}
                  >
                    {this.linkContent(iconc, titlec, hrefc, linkc, isDrawerCollapsed)}
                  </ListItem>
                  {this.renderChildren(idname, childrenc, depth + 1)}
                </div>
              );
            },
          )}
        </List>
      );
    }
    if (idname == LIFECYCLE) {
      if (children && children.length > 0) {
        return (
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
              }) => {
                if (typeof showc !== 'undefined' && !showc) {
                  return '';
                }
                return (
                  <div key={idc} className={!showc ? classes.cursorNotAllowed : null}>
                    <ListItem
                      data-cy={idc}
                      button
                      key={idc}
                      className={classNames(
                        depth === 1 ? classes.nested1 : classes.nested2,
                        classes.item,
                        classes.itemActionable,
                        path === hrefc && classes.itemActiveItem,
                        isDrawerCollapsed && classes.noPadding,
                        !showc && classes.disabled,
                      )}
                      onClick={() => this.handleAdapterClick(idc, linkc)}
                    >
                      {this.linkContent(iconc, titlec, hrefc, linkc, isDrawerCollapsed)}
                    </ListItem>
                    {this.renderChildren(idname, childrenc, depth + 1)}
                  </div>
                );
              },
            )}
          </List>
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
    const { classes } = this.props;
    let linkContent = (
      <div className={classNames(classes.link)}>
        <Tooltip
          title={titlec}
          placement="right"
          disableFocusListener={!drawerCollapsed}
          disableHoverListener={!drawerCollapsed}
          disableTouchListener={!drawerCollapsed}
        >
          <ListItemIcon className={classes.listIcon}>{iconc} </ListItemIcon>
        </Tooltip>
        <ListItemText
          className={drawerCollapsed ? classes.isHidden : classes.isDisplayed}
          classes={{ primary: classes.itemPrimary }}
        >
          {titlec}
        </ListItemText>
      </div>
    );
    if (linkc) {
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
            <Tooltip title={`Newer version of Meshery available: ${latest}`} placement="right">
              <OpenInNewIcon style={{ width: '0.85rem', verticalAlign: 'middle' }} />
            </Tooltip>
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
    const { classes, isDrawerCollapsed, ...other } = this.props;
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
        <ListItem
          component="a"
          onClick={this.handleTitleClick}
          className={classNames(
            classes.firebase,
            classes.item,
            classes.itemCategory,
            classes.cursorPointer,
            !this.state.capabilitiesRegistryObj?.isNavigatorComponentEnabled([DASHBOARD]) &&
              classes.disableLogo,
          )}
        >
          <img
            className={isDrawerCollapsed ? classes.mainLogoCollapsed : classes.mainLogo}
            src="/static/img/meshery-logo.png"
            onClick={this.handleTitleClick}
          />
          <img
            className={isDrawerCollapsed ? classes.mainLogoTextCollapsed : classes.mainLogoText}
            src="/static/img/meshery-logo-text.png"
            onClick={this.handleTitleClick}
          />

          {/* <span className={isDrawerCollapsed ? classes.isHidden : classes.isDisplayed}>Meshery</span> */}
        </ListItem>
      </div>
    );
    const Menu = (
      <List disablePadding className={classes.hideScrollbar}>
        {navigatorComponents.map(
          ({ id: childId, title, icon, href, show, link, children, hovericon, submenu }) => {
            // if (typeof show !== "undefined" && !show) {
            //   return "";
            // }
            return (
              <div key={childId} style={!show ? cursorNotAllowed : {}} className={classes.root}>
                <ListItem
                  button={!!link}
                  dense
                  key={childId}
                  className={classNames(
                    classes.item,
                    link ? classes.itemActionable : '',
                    path === href && classes.itemActiveItem,
                    !show && classes.disabled,
                  )}
                  onClick={() => this.toggleItemCollapse(childId)}
                  onMouseOver={() =>
                    isDrawerCollapsed ? this.setState({ hoveredId: childId }) : null
                  }
                  onMouseLeave={() =>
                    !submenu || !this.state.openItems.includes(childId)
                      ? this.setState({ hoveredId: false })
                      : null
                  }
                >
                  <Link href={link ? href : ''}>
                    <div data-cy={childId} className={classNames(classes.link)}>
                      <Tooltip
                        title={childId}
                        placement="right"
                        disableFocusListener={!isDrawerCollapsed}
                        disableHoverListener={true}
                        disableTouchListener={!isDrawerCollapsed}
                        TransitionComponent={Zoom}
                        arrow
                      >
                        {isDrawerCollapsed &&
                        (this.state.hoveredId === childId ||
                          (this.state.openItems.includes(childId) && submenu)) ? (
                          <div>
                            <Tooltip
                              title={title}
                              placement="right"
                              TransitionComponent={Zoom}
                              arrow
                            >
                              <ListItemIcon
                                onClick={() => this.toggleItemCollapse(childId)}
                                style={{ marginLeft: '20%', marginBottom: '0.4rem' }}
                              >
                                {hovericon}
                              </ListItemIcon>
                            </Tooltip>
                          </div>
                        ) : (
                          <ListItemIcon className={classes.listIcon}>{icon}</ListItemIcon>
                        )}
                      </Tooltip>
                      <ListItemText
                        className={isDrawerCollapsed ? classes.isHidden : classes.isDisplayed}
                        classes={{ primary: classes.itemPrimary }}
                      >
                        {title}
                      </ListItemText>
                    </div>
                  </Link>
                  <FontAwesomeIcon
                    icon={faCaretDown}
                    onClick={() => this.toggleItemCollapse(childId)}
                    className={classNames(classes.expandMoreIcon, {
                      [classes.collapsed]: this.state.openItems.includes(childId),
                    })}
                    style={isDrawerCollapsed || !children ? { opacity: 0 } : {}}
                  />
                </ListItem>
                <Collapse
                  in={this.state.openItems.includes(childId)}
                  style={{ backgroundColor: '#396679', opacity: '100%' }}
                >
                  {this.renderChildren(childId, children, 1)}
                </Collapse>
              </div>
            );
          },
        )}
        {this.state.navigator && this.state.navigator.length ? (
          <React.Fragment>
            <Divider className={classes.divider} />
            {this.renderNavigatorExtensions(this.state.navigator, 1)}
          </React.Fragment>
        ) : null}
        <Divider className={classes.divider} />
      </List>
    );
    const HelpIcons = (
      <ButtonGroup
        size="large"
        className={!isDrawerCollapsed ? classes.marginLeft : classes.btnGrpMarginRight}
        orientation={isDrawerCollapsed ? 'vertical' : 'horizontal'}
      >
        {externlinks.map(({ id, icon, title, href }, index) => {
          return (
            <ListItem
              key={id}
              className={classes.item}
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
                  className={classNames(
                    classes.link,
                    isDrawerCollapsed ? classes.extraPadding : '',
                  )}
                >
                  <Tooltip title={title} placement={isDrawerCollapsed ? 'right' : 'top'}>
                    <ListItemIcon className={classNames(classes.listIcon, classes.helpIcon)}>
                      {icon}
                    </ListItemIcon>
                  </Tooltip>
                </a>
              </Grow>
            </ListItem>
          );
        })}
        <ListItem
          className={classes.rightMargin}
          style={!isDrawerCollapsed ? { display: 'none' } : { marginLeft: '4px' }}
        >
          <Tooltip title="Help" placement={isDrawerCollapsed ? 'right' : 'top'}>
            <IconButton
              className={isDrawerCollapsed ? classes.collapsedHelpButton : classes.rightTranslate}
              onClick={this.toggleSpacing}
            >
              <HelpIcon
                className={classes.helpIcon}
                style={{ fontSize: '1.45rem', ...iconSmall }}
              />
            </IconButton>
          </Tooltip>
        </ListItem>
      </ButtonGroup>
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
      <div
        className={classNames(
          isDrawerCollapsed ? classes.collapseButtonWrapperRotated : classes.collapseButtonWrapper,
        )}
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
          <FontAwesomeIcon
            icon={faAngleLeft}
            fixedWidth
            size="1.5x"
            style={{ margin: '0.75rem 0.2rem ', width: '0.8rem', verticalAlign: 'middle' }}
            alt="Sidebar collapse toggle icon"
          />
        </div>
      </div>
    );

    return (
      <NoSsr>
        <Drawer
          variant="permanent"
          {...other}
          className={isDrawerCollapsed ? classes.sidebarCollapsed : classes.sidebarExpanded}
          classes={{
            paper: isDrawerCollapsed ? classes.sidebarCollapsed : classes.sidebarExpanded,
          }}
          style={{ width: 'inherit' }}
        >
          {Title}
          {Menu}
          <div className={classes.fixedSidebarFooter}>
            {Chevron}
            {HelpIcons}
            {Version}
          </div>
        </Drawer>
      </NoSsr>
    );
  }
}

Navigator.propTypes = {
  classes: PropTypes.object.isRequired,
  onCollapseDrawer: PropTypes.func.isRequired,
};

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
  return { meshAdapters, meshAdaptersts, path, isDrawerCollapsed, capabilitiesRegistry };
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(Navigator)),
);
