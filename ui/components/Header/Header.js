/* eslint-disable react/prop-types */
import React from "react";
import SettingsIcon from "@mui/icons-material/Settings";
import { AppBar, Grid, Hidden, IconButton, Toolbar, Typography } from "@mui/material";
import Link from "next/link";
import MenuIcon from "@mui/icons-material/Menu";
import { useStyles } from "./Header.styles";
import K8sContextMenu from "../K8sContextSwitcher/K8sContextSwitcher";

const HeaderComponent = ({ drawerOpen, onDrawerToggle, pageTitle }) => {
  const classes = useStyles();

  return (
    <AppBar
      color="primary"
      position="sticky"
      elevation={0}
      className={drawerOpen ? classes.appBarOnDrawerClosed : classes.appBarOnDrawerOpen}
    >
      <Toolbar className={drawerOpen ? classes.toolbarOnDrawerClosed : classes.toolbarOnDrawerOpen}>
        <Grid container alignItems="center">
          <Hidden smUp>
            <Grid item>
              <IconButton
                color="inherit"
                aria-label="Open drawer"
                onClick={onDrawerToggle}
                className={classes.menuButton}
                size="large"
              >
                <MenuIcon className={classes.headerIcons} />
              </IconButton>
            </Grid>
          </Hidden>
          <Grid item xs container alignItems="center" className={classes.pageTitleWrapper}>
            <Typography color="inherit" variant="h5" className={classes.pageTitle}>
              {pageTitle}
            </Typography>
          </Grid>

          {/* <Grid item className={classes.notifications}>
                <MesheryNotification />
              </Grid> */}
          <Grid item className={classes.userContainer}>
            <K8sContextMenu contexts={contexts}/>
            <IconButton color="inherit" size="large">
              <Link href="/system/connections">
                <img
                  src={
                    pageTitle === "Connection Wizard"
                      ? "/static/img/connection_wizard/connection-wizard-green.svg"
                      : "/static/img/connection_wizard/connection-wizard-white.svg"
                  }
                  className={
                    classes.headerIcons + " " + (pageTitle === "Connection Wizard" ? classes.itemActiveItem : "")
                  }
                />
              </Link>
            </IconButton>

            <IconButton color="inherit" size="large">
              <Link href="/settings">
                <SettingsIcon
                  className={classes.headerIcons + " " + (pageTitle === "Settings" ? classes.itemActiveItem : "")}
                />
              </Link>
            </IconButton>

            {/* <div data-test="notification-button">
              <MesheryNotification />
            </div> */}
            {/* <span className={classes.userSpan}>
              <User color="inherit" iconButtonClassName={classes.iconButtonAvatar} avatarClassName={classes.avatar} />
            </span> */}
          </Grid>
        </Grid>
      </Toolbar>
    </AppBar>
  );
};

export default HeaderComponent;
