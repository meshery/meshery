import React from "react";
import PropTypes from "prop-types";
import classNames from "classnames";
import { withStyles } from "@material-ui/core/styles";
import Divider from "@material-ui/core/Divider";
import Drawer from "@material-ui/core/Drawer";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import RemoveIcon from "@material-ui/icons/Remove";
import Link from "next/link";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import NoSsr from "@material-ui/core/NoSsr";
import Avatar from "@material-ui/core/Avatar";
import { withRouter } from "next/router";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faTerminal,
  faTachometerAlt,
  faExternalLinkAlt,
  faChevronCircleLeft,
  faPollH,
} from "@fortawesome/free-solid-svg-icons";
import { updatepagetitle } from "../lib/store";
import { Tooltip } from "@material-ui/core";

const styles = (theme) => ({
  categoryHeader: {
    paddingTop: 16,
    paddingBottom: 16,
  },
  categoryHeaderPrimary: {
    color: theme.palette.common.white,
  },
  item: {
    paddingTop: 4,
    paddingBottom: 4,
    color: "rgba(255, 255, 255, 0.7)",
  },
  itemCategory: {
    backgroundColor: "#263238",
    boxShadow: "0 -1px 0 #404854 inset",
    paddingTop: 16,
    paddingBottom: 16,
  },
  firebase: {
    fontSize: 24,
    fontFamily: theme.typography.fontFamily,
    color: theme.palette.common.white,
  },
  link: {
    display: "inline-flex",
    width: "100%",
    height: "30px",
  },
  itemActionable: {
    "&:hover": {
      backgroundColor: "rgba(255, 255, 255, 0.08)",
    },
  },
  itemActiveItem: {
    color: "#4fc3f7",
  },
  itemPrimary: {
    color: "inherit",
    fontSize: theme.typography.fontSize,
    "&$textDense": {
      fontSize: theme.typography.fontSize,
    },
  },
  textDense: {},
  divider: {
    marginTop: theme.spacing(2),
  },
  mainLogo: {
    marginRight: theme.spacing(1),
    marginTop: theme.spacing(1),
    marginLeft: theme.spacing(-1),
    width: 40,
    height: 40,
    borderRadius: "unset",
  },
  mainLogoText: {
    marginLeft: theme.spacing(0.5),
    marginTop: theme.spacing(1),
    width: 170,
    height: "100%",
    borderRadius: "unset",
  },
  community: {
    marginTop: theme.spacing(2),
  },
  settingsIcon: {
    marginLeft: theme.spacing(2),
  },
  cursorPointer: {
    cursor: "pointer",
  },
  listIcon: {
    minWidth: theme.spacing(3.5),
    paddingTop: theme.spacing(0.5),
    textAlign: "center",
    display: "inline-table",
    paddingRight: theme.spacing(0.5),
  },
  nested1: {
    paddingLeft: theme.spacing(3),
  },
  nested2: {
    paddingLeft: theme.spacing(5),
  },
  icon: {
    width: theme.spacing(2.5),
  },
  istioIcon: {
    width: theme.spacing(1.8),
  },
  isHidden: {
    opacity: 0,
    transition: "opacity 200ms ease-in-out",
  },
  isDisplayed: {
    opacity: 1,
    transition: "opacity 200ms ease-in-out",
  },
  sidebarCollapsed: {
    transition: theme.transitions.create("width", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    overflowX: "hidden",
    width: theme.spacing(7) + 1,
  },
  sidebarExpanded: {
    width: "256px",
    overflowX: "hidden",
    transition: theme.transitions.create("width", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  fixedSidebarFooter: {
    display: "none",
    [theme.breakpoints.up("sm")]: {
      display: "flex",
      "margin-top": "auto",
      "margin-bottom": "0.5rem",
    },
  },
  collapseButtonWrapper: {
    width: "auto",
    "margin-left": "auto",
    opacity: "0.7",
    transition: "opacity 200ms linear",
    "&:hover": {
      opacity: 1,
      background: "transparent",
    },
    "&:focus": {
      opacity: 1,
      background: "transparent",
    },
  },
  collapseButtonWrapperRotated: {
    width: "auto",
    "margin-left": "auto",
    opacity: "0.7",
    transition: "opacity 200ms linear",
    transform: "rotate(180deg)",
    "&:hover": {
      opacity: 1,
      background: "transparent",
    },
    "&:focus": {
      opacity: 1,
      background: "transparent",
    },
  },
  noPadding: {
    paddingLeft: "16px",
    paddingRight: "16px",
  },
});

const categories = [
  {
    id: "Dashboard",
    href: "/",
    title: "System Dashboard",
    show: false,
    link: true,
  },
  {
    id: "Performance",
    icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
    href: "/performance",
    title: "Performance Test",
    show: true,
    link: true,
    children: [
      {
        id: "Results",
        icon: <FontAwesomeIcon icon={faPollH} fixedWidth />,
        href: "/results",
        title: "View & Compare Results",
        show: true,
        link: true,
      },
    ],
  },
  {
    id: "Settings",
    href: "/settings",
    title: "Settings",
    show: false,
    link: true,
  }, // title is used for comparison in the Header.js file as well
  {
    id: "Conformance",
    icon: <FontAwesomeIcon icon={faTasks} transform="shrink-2" fixedWidth />,
    href: "/smi_results", //Temp
    title: "Conformance",
    show: true,
    link: true,
    children: [
      {
        id: "SMI Results",
        icon: <FontAwesomeIcon icon={faPollH} fixedWidth />,
        href: "/smi_results",
        title: "Service Mesh Interface Results",
        show: true,
        link: true,
      },
    ],
  },
  {
    id: "Management",
    icon: <FontAwesomeIcon icon={faTerminal} transform="shrink-4" fixedWidth />,
    href: "/management",
    title: "Management",
    show: true,
    link: true,
    children: [
      {
        id: "Citrix Service Mesh",
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
        href: "/management/citrix",
        title: "Citrix Service Mesh",
        link: false,
        show: true,
      },
      {
        id: "Consul",
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
        href: "/management/consul",
        title: "Consul",
        link: false,
        show: true,
      },
      {
        id: "Istio",
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
        href: "/management/istio",
        title: "Istio",
        link: false,
        show: true,
      },
      {
        id: "Kuma",
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
        href: "/management/kuma",
        title: "Kuma",
        link: false,
        show: true,
      },
      {
        id: "Linkerd",
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
        href: "/management/linkerd",
        title: "Linkerd",
        link: false,
        show: true,
      },
      {
        id: "Network Service Mesh",
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
        href: "/management/nsm",
        title: "Network Service Mesh",
        link: false,
        show: true,
      },
      {
        id: "NGINX Service Mesh",
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
        href: "/management/nginx",
        title: "NGINX Service Mesh",
        link: false,
        show: true,
      },
      {
        id: "Octarine",
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
        href: "/management/octarine",
        title: "Octarine",
        link: false,
        show: true,
      },
      {
        id: "Open Service Mesh",
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
        href: "/management/osm",
        title: "Open Service Mesh",
        link: false,
        show: true,
      },
    ],
  },
];

class Navigator extends React.Component {
  constructor(props) {
    super(props);
    const { meshAdapters } = props;
    this.state = {
      path: "",
      meshAdapters,
      mts: new Date(),
    };
  }

  updateCategoriesMenus() {
    const self = this;
    categories.forEach((cat, ind) => {
      if (cat.id === "Management") {
        cat.children.forEach((catc, ind1) => {
          const cr = self.fetchChildren(catc.id);
          const icon = self.pickIcon(catc.id);
          categories[ind].children[ind1].icon = icon;
          categories[ind].children[ind1].children = cr;
        });
      }
    });
  }

  updateAdaptersLink() {
    categories.forEach((cat, ind) => {
      if (cat.id === "Management") {
        cat.children.forEach((catc, ind1) => {
          if (
            typeof categories[ind].children[ind1].children[0] !== "undefined" &&
            typeof categories[ind].children[ind1].children[0].href !== "undefined"
          ) {
            const val = true;
            const newhref = `${categories[ind].children[ind1].children[0].href}`;
            categories[ind].children[ind1].link = val;
            categories[ind].children[ind1].href = newhref;
          }
        });
      }
    });
  }

  static getDerivedStateFromProps(props, state) {
    const { meshAdapters, meshAdaptersts, path } = props;
    const st = {};
    if (meshAdaptersts > state.mts) {
      st.meshAdapters = meshAdapters;
      st.mts = meshAdaptersts;
    }

    const fetchNestedPathAndTitle = (path, title, href, children) => {
      if (href === path) {
        props.updatepagetitle({ title });
        return;
      }
      if (children && children.length > 0) {
        children.forEach(({ title, href, children }) => {
          fetchNestedPathAndTitle(path, title, href, children);
        });
      }
    };

    categories.forEach(({ title, href, children }) => {
      fetchNestedPathAndTitle(path, title, href, children);
    });
    st.path = path;
    return st;
  }

  fetchChildren(category) {
    const { meshAdapters } = this.state;
    const children = [];
    category = category.toLowerCase();
    meshAdapters.forEach((adapter) => {
      const aName = adapter.name.toLowerCase();
      if (category !== aName) {
        return;
      }
      children.push({
        id: adapter.adapter_location,
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
        icon: <RemoveIcon />,
        href: `/management?adapter=${adapter.adapter_location}`,
        title: `Management - ${adapter.adapter_location}`,
        link: true,
        show: true,
      });
    });
    return children;
  }

  pickIcon(aName) {
    aName = aName.toLowerCase();
    const { classes } = this.props;
    let image = "/static/img/meshery-logo.png";
    let logoIcon = <img src={image} className={classes.icon} />;
    switch (aName) {
      case "istio":
        image = "/static/img/istio-light.svg";
        logoIcon = <img src={image} className={classes.istioIcon} />;
        break;
      case "linkerd":
        image = "/static/img/linkerd-light.svg";
        logoIcon = <img src={image} className={classes.icon} />;
        break;
      case "consul":
        image = "/static/img/consul-light.svg";
        logoIcon = <img src={image} className={classes.icon} />;
        break;
      case "network service mesh":
        image = "/static/img/nsm-light.svg";
        logoIcon = <img src={image} className={classes.icon} />;
        break;
      case "octarine":
        image = "/static/img/octarine-white.svg";
        logoIcon = <img src={image} className={classes.icon} />;
        break;
      case "citrix service mesh":
        image = "/static/img/citrix-light.svg";
        logoIcon = <img src={image} className={classes.icon} />;
        break;
      case "open service mesh":
        image = "/static/img/osm-white.svg";
        logoIcon = <img src={image} className={classes.icon} />;
        break;
      case "kuma":
        image = "/static/img/kuma-light.svg";
        logoIcon = <img src={image} className={classes.icon} />;
        break;
      case "nginx service mesh":
        image = "/static/img/nginx-sm-light.svg";
        logoIcon = <img src={image} className={classes.icon} />;
        break;
    }
    return logoIcon;
  }

  handleTitleClick = () => {
    this.props.router.push("/");
  };

  handleAdapterClick = (id, link) => {
    let allowedId = [
      "Consul",
      "Istio",
      "Linkerd",
      "Network Service Mesh",
      "Octarine",
      "Citrix Service Mesh",
      "Open Service Mesh",
      "Kuma",
      "NGINX Service Mesh"
    ];
    let index = allowedId.indexOf(id);
    if (index != -1 && !link) {
      this.props.router.push("/management");
    }
  };

  toggleMiniDrawer = () => {
    const { onCollapseDrawer } = this.props;
    onCollapseDrawer();
  };

  renderChildren(idname, children, depth) {
    const { classes, isDrawerCollapsed } = this.props;
    const { path } = this.state;

    if (idname != "Management" && children && children.length > 0) {
      return (
        <List disablePadding>
          {children.map(({ id: idc, icon: iconc, href: hrefc, show: showc, link: linkc, children: childrenc }) => {
            if (typeof showc !== "undefined" && !showc) {
              return "";
            }
            return (
              <React.Fragment key={idc}>
                <ListItem
                  button
                  key={idc}
                  className={classNames(
                    depth === 1 ? classes.nested1 : classes.nested2,
                    classes.item,
                    classes.itemActionable,
                    path === hrefc && classes.itemActiveItem,
                    isDrawerCollapsed && classes.noPadding
                  )}
                >
                  {this.linkContent(iconc, idc, hrefc, linkc, isDrawerCollapsed)}
                </ListItem>
                {this.renderChildren(idname, childrenc, depth + 1)}
              </React.Fragment>
            );
          })}
        </List>
      );
    }
    if (idname == "Management") {
      if (children && children.length > 1) {
        return (
          <List disablePadding>
            {children.map(({ id: idc, icon: iconc, href: hrefc, show: showc, link: linkc, children: childrenc }) => {
              if (typeof showc !== "undefined" && !showc) {
                return "";
              }
              return (
                <React.Fragment key={idc}>
                  <ListItem
                    button
                    key={idc}
                    className={classNames(
                      depth === 1 ? classes.nested1 : classes.nested2,
                      classes.item,
                      classes.itemActionable,
                      path === hrefc && classes.itemActiveItem,
                      isDrawerCollapsed && classes.noPadding
                    )}
                    onClick={() => this.handleAdapterClick(idc, linkc)}
                  >
                    {this.linkContent(iconc, idc, hrefc, linkc, isDrawerCollapsed)}
                  </ListItem>
                  {this.renderChildren(idname, childrenc, depth + 1)}
                </React.Fragment>
              );
            })}
          </List>
        );
      }
      if (children && children.length == 1) {
        this.updateAdaptersLink();
      }
    }
    return "";
  }

  linkContent(iconc, idc, hrefc, linkc, drawerCollapsed) {
    const { classes } = this.props;

    let linkContent = (
      <div className={classNames(classes.link)}>
        <Tooltip
          title={idc}
          placement="right"
          disableFocusListener={!drawerCollapsed}
          disableHoverListener={!drawerCollapsed}
          disableTouchListener={!drawerCollapsed}
        >
          <ListItemIcon className={classes.listIcon}>{iconc}</ListItemIcon>
        </Tooltip>
        <ListItemText
          className={drawerCollapsed ? classes.isHidden : classes.isDisplayed}
          classes={{
            primary: classes.itemPrimary,
          }}
        >
          {idc}
        </ListItemText>
      </div>
    );
    if (linkc) {
      linkContent = <Link href={hrefc}>{linkContent}</Link>;
    }
    return linkContent;
  }

  render() {
    const { classes, isDrawerCollapsed, ...other } = this.props;
    const { path } = this.state;
    this.updateCategoriesMenus();
    let classname;
    if (isDrawerCollapsed) {
      classname = classes.collapseButtonWrapperRotated;
    } else {
      classname = classes.collapseButtonWrapper;
    }
    return (
      <NoSsr>
        <Drawer
          variant="permanent"
          {...other}
          className={isDrawerCollapsed ? classes.sidebarCollapsed : classes.sidebarExpanded}
          classes={{
            paper: isDrawerCollapsed ? classes.sidebarCollapsed : classes.sidebarExpanded,
          }}
          style={{ width: "inherit" }}
        >
          <List disablePadding>
            <ListItem
              component="a"
              onClick={this.handleTitleClick}
              className={classNames(classes.firebase, classes.item, classes.itemCategory, classes.cursorPointer)}
            >
              <Avatar className={classes.mainLogo} src="/static/img/meshery-logo.png" onClick={this.handleTitleClick} />
              <Avatar
                className={classes.mainLogoText}
                src="/static/img/meshery-logo-text.png"
                onClick={this.handleTitleClick}
              />

              {/* <span className={isDrawerCollapsed ? classes.isHidden : classes.isDisplayed}>Meshery</span> */}
            </ListItem>
            {categories.map(({ id: childId, icon, href, show, link, children }) => {
              if (typeof show !== "undefined" && !show) {
                return "";
              }
              return (
                <React.Fragment key={childId}>
                  <ListItem
                    button
                    dense
                    key={childId}
                    className={classNames(
                      classes.item,
                      classes.itemActionable,
                      path === href && classes.itemActiveItem
                    )}
                  >
                    <Link href={link ? href : ""}>
                      <div className={classNames(classes.link)}>
                        <Tooltip
                          title={childId}
                          placement="right"
                          disableFocusListener={!isDrawerCollapsed}
                          disableHoverListener={!isDrawerCollapsed}
                          disableTouchListener={!isDrawerCollapsed}
                        >
                          <ListItemIcon className={classes.listIcon}>{icon}</ListItemIcon>
                        </Tooltip>
                        <ListItemText
                          className={isDrawerCollapsed ? classes.isHidden : classes.isDisplayed}
                          classes={{
                            primary: classes.itemPrimary,
                          }}
                        >
                          {childId}
                        </ListItemText>
                      </div>
                    </Link>
                  </ListItem>
                  {this.renderChildren(childId, children, 1)}
                </React.Fragment>
              );
            })}
            <Divider className={classes.divider} />
            <ListItem
              component="a"
              href="https://meshery.io/"
              target="_blank"
              key="about"
              className={classNames(classes.item, classes.itemActionable, classes.community)}
            >
              <div className={classNames(classes.link)}>
                <Tooltip
                  title="Community"
                  placement="right"
                  disableFocusListener={!isDrawerCollapsed}
                  disableHoverListener={!isDrawerCollapsed}
                  disableTouchListener={!isDrawerCollapsed}
                >
                  <ListItemIcon className={classes.listIcon}>
                    <FontAwesomeIcon icon={faExternalLinkAlt} transform="shrink-2" fixedWidth />
                  </ListItemIcon>
                </Tooltip>
                <ListItemText
                  className={isDrawerCollapsed ? classes.isHidden : classes.isDisplayed}
                  classes={{
                    primary: classes.itemPrimary,
                  }}
                >
                  Community
                </ListItemText>
              </div>
            </ListItem>
          </List>
          <div className={classes.fixedSidebarFooter}>
            <ListItem button onClick={() => this.toggleMiniDrawer()} className={classname}>
              <FontAwesomeIcon
                icon={faChevronCircleLeft}
                fixedWidth
                color="#FFFFFF"
                size="lg"
                alt="Sidebar collapse toggle icon"
              />
            </ListItem>
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
});

const mapStateToProps = (state) => {
  const meshAdapters = state.get("meshAdapters").toJS();
  const meshAdaptersts = state.get("meshAdaptersts");
  const path = state.get("page").get("path");
  return { meshAdapters, meshAdaptersts, path };
};

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withRouter(Navigator)));
