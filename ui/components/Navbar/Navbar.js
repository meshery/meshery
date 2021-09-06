/* eslint-disable react/prop-types */
import { Avatar, Drawer, Hidden, Link, List, ListItem, ListItemIcon, ListItemText, Tooltip } from "@material-ui/core";
import React from "react";
import { useStyles } from "./Navbar.styles";
import classNames from "classnames";
import { ExpandMoreIcon } from "@material-ui/icons";

const Navigator = ({ isDrawerCollapsed, classes, handleTitleClick, categories }) => {
  return (
    <Drawer
      variant="permanent"
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
          <Avatar
            className={isDrawerCollapsed ? classes.mainLogoCollapsed : classes.mainLogo}
            src="/static/img/meshery-logo.png"
            onClick={handleTitleClick}
          />
          <Avatar
            className={isDrawerCollapsed ? classes.mainLogoTextCollapsed : classes.mainLogoText}
            src="/static/img/meshery-logo-text.png"
            onClick={handleTitleClick}
          />

          {/* <span className={isDrawerCollapsed ? classes.isHidden : classes.isDisplayed}>Meshery</span> */}
        </ListItem>
        {categories.map(({ id: childId, title, icon, href, show, link, children }) => {
          if (typeof show !== "undefined" && !show) {
            return "";
          }
          return (
            <React.Fragment key={childId}>
              <ListItem
                button={!!link}
                dense
                key={childId}
                className={classNames(
                  classes.item,
                  link ? classes.itemActionable : "",
                  path === href && classes.itemActiveItem
                )}
                onClick={() => this.toggleItemCollapse(childId)}
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
                      {title}
                    </ListItemText>
                  </div>
                </Link>
                <ExpandMoreIcon
                  onClick={() => this.toggleItemCollapse(childId)}
                  className={classNames(classes.expandMoreIcon, {
                    [classes.collapsed]: this.state.openItems.includes(childId),
                  })}
                  style={isDrawerCollapsed || !children ? { opacity: 0 } : {}}
                />
              </ListItem>
              <Collapse in={isDrawerCollapsed || this.state.openItems.includes(childId)}>
                {this.renderChildren(childId, children, 1)}
              </Collapse>
            </React.Fragment>
          );
        })}
        {this.state.navigator && this.state.navigator.length ? (
          <React.Fragment>
            <Divider className={classes.divider} />
            {this.renderNavigatorExtensions(this.state.navigator, 1)}
          </React.Fragment>
        ) : null}
        <Divider className={classes.divider} />
      </List>
      <div className={classes.fixedSidebarFooter}>
        <ButtonGroup
          size="large"
          className={!isDrawerCollapsed ? classes.marginLeft : classes.btnGrpMarginRight}
          orientation={isDrawerCollapsed ? "vertical" : "horizontal"}
        >
          {externlinks.map(({ id, icon, title, href }, index) => {
            return (
              <ListItem
                key={id}
                className={classes.item}
                style={isDrawerCollapsed && !showHelperButton ? { display: "none" } : {}}
              >
                <Grow in={showHelperButton} timeout={{ enter: 600 - index * 200, exit: 100 * index }}>
                  <a
                    href={href}
                    target="_blank"
                    rel="noreferrer"
                    className={classNames(classes.link, isDrawerCollapsed ? classes.extraPadding : "")}
                  >
                    <Tooltip title={title} placement={isDrawerCollapsed ? "right" : "top"}>
                      <ListItemIcon className={classNames(classes.listIcon, classes.helpIcon)}>{icon}</ListItemIcon>
                    </Tooltip>
                  </a>
                </Grow>
              </ListItem>
            );
          })}
          <ListItem className={classes.rightMargin}>
            <Tooltip title="Help" placement={isDrawerCollapsed ? "right" : "top"}>
              <IconButton
                className={isDrawerCollapsed ? classes.collapsedHelpButton : classes.rightTranslate}
                onClick={() => this.toggleSpacing()}
              >
                <HelpIcon className={classes.helpIcon} style={{ fontSize: "1.45rem" }} />
              </IconButton>
            </Tooltip>
          </ListItem>
        </ButtonGroup>

        <ListItem
          button
          className={classname}
          onClick={() => this.toggleMiniDrawer()}
          style={{ position: "sticky", zIndex: "1", bottom: "0", right: "0" }}
        >
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
  );
};

export const Navbar = () => {
  const classes = useStyles();
  return (
    <nav className={isDrawerCollapsed ? classes.drawerCollapsed : classes.drawer} data-test="navigation">
      <Hidden smUp implementation="js">
        <Navigator
          variant="temporary"
          open={this.state.mobileOpen}
          onClose={this.handleDrawerToggle}
          onCollapseDrawer={this.handleCollapseDrawer}
          isDrawerCollapsed={isDrawerCollapsed}
        />
      </Hidden>
      <Hidden xsDown implementation="css">
        <Navigator onCollapseDrawer={this.handleCollapseDrawer} isDrawerCollapsed={isDrawerCollapsed} />
      </Hidden>
    </nav>
  );
};
