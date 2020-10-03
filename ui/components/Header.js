import React from 'react';
import PropTypes from 'prop-types';
import AppBar from '@material-ui/core/AppBar';
import Grid from '@material-ui/core/Grid';
import Hidden from '@material-ui/core/Hidden';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import { connect } from 'react-redux';
import NoSsr from '@material-ui/core/NoSsr';
import Link from 'next/link';
import SettingsIcon from '@material-ui/icons/Settings';
import DashboardIcon from '@material-ui/icons/Dashboard';
import MesheryNotification from './MesheryNotification';
import User from './User';
// import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
// import { faHome } from '@fortawesome/free-solid-svg-icons';

const lightColor = 'rgba(255, 255, 255, 0.7)';

const styles = (theme) => ({
  secondaryBar: {
    zIndex: 0,
  },
  menuButton: {
    marginLeft: -theme.spacing(1),
  },
  iconButtonAvatar: {
    padding: 4,
  },
  link: {
    textDecoration: 'none',
    color: lightColor,
    '&:hover': {
      color: theme.palette.common.white,
    },
  },
  button: {
    borderColor: lightColor,
  },
  notifications: {
    paddingLeft: theme.spacing(4),
    paddingRight: theme.spacing(0),
    marginLeft: theme.spacing(4),
  },
  userContainer: {
    paddingLeft: 1,
    display: 'flex',
  },
  userSpan: {
    marginLeft: theme.spacing(1),
  },
  pageTitleWrapper: {
    flexGrow: 1,
    marginRight: 'auto',
  },
  pageTitle: {
    paddingLeft: theme.spacing(2),
    fontSize: '22px',
    [theme.breakpoints.up('sm')]: {
      fontSize: '26px',
    },
  },
  appBar: {
    padding: theme.spacing(1.4),
  },
  itemActiveItem: {
    color: '#4fc3f7',
  },
});

class Header extends React.Component {
  render() {
    const { classes, title, onDrawerToggle } = this.props;
    return (
      <NoSsr>
        <React.Fragment>
          <AppBar color="primary" position="sticky" elevation={0} className={classes.appBar}>
            <Toolbar>
              <Grid container alignItems="center">
                <Hidden smUp>
                  <Grid item>
                    <IconButton
                      color="inherit"
                      aria-label="Open drawer"
                      onClick={onDrawerToggle}
                      className={classes.menuButton}
                    >
                      <MenuIcon />
                    </IconButton>
                  </Grid>
                </Hidden>
                <Grid item xs container alignItems="center" className={classes.pageTitleWrapper}>
                  <Typography color="inherit" variant="h5" className={classes.pageTitle}>
                    {title === 'Dashboard' ? 'System Dashboard' : title}
                  </Typography>
                </Grid>
                {/* <Grid item className={classes.notifications}>
                <MesheryNotification />
              </Grid> */}
                <Grid item className={classes.userContainer}>
                  {/* <IconButton color="inherit" className={classes.iconButtonAvatar}>
                  <Avatar className={classes.avatar} src="/static/images/avatar/1.jpg" />
                </IconButton> */}

                  <IconButton color="inherit">
                    <Link href="/">
                      <DashboardIcon className={title === 'Dashboard' ? classes.itemActiveItem : ''} />
                      {/* <FontAwesomeIcon icon={faHome} transform="shrink-2" fixedWidth className={title === 'Dashboard' && classes.itemActiveItem} /> */}
                    </Link>
                  </IconButton>

                  <IconButton color="inherit">
                    <Link href="/settings">
                      <SettingsIcon className={title === 'Settings' ? classes.itemActiveItem : ''} />
                    </Link>
                  </IconButton>

                  <MesheryNotification />
                  <span className={classes.userSpan}>
                    <User color="inherit" iconButtonClassName={classes.iconButtonAvatar} avatarClassName={classes.avatar} />
                  </span>
                </Grid>
              </Grid>
            </Toolbar>
          </AppBar>
          {/* <AppBar
          component="div"
          className={classes.secondaryBar}
          color="primary"
          position="static"
          elevation={0}
        >
          <Toolbar>
            <Grid container alignItems="center" spacing={8}>
              <Grid item xs>

              </Grid>
              {/* <Grid item>
                <Button className={classes.button} variant="outlined" color="inherit" size="small">
                  Web setup
                </Button>
              </Grid> * /}
              {/* <Grid item>
                <Tooltip title="Help">
                  <IconButton color="inherit">
                    <HelpIcon />
                  </IconButton>
                </Tooltip>
              </Grid> * /}
            </Grid>
          </Toolbar>
        </AppBar> */}
          {/* <AppBar
          component="div"
          className={classes.secondaryBar}
          color="primary"
          position="static"
          elevation={0}
        >
          <Tabs value={0} textColor="inherit">
            <Tab textColor="inherit" label="Users" />
            <Tab textColor="inherit" label="Sign-in method" />
            <Tab textColor="inherit" label="Templates" />
            <Tab textColor="inherit" label="Usage" />
          </Tabs>
        </AppBar> */}
        </React.Fragment>
      </NoSsr>
    );
  }
}

Header.propTypes = {
  classes: PropTypes.object.isRequired,
  onDrawerToggle: PropTypes.func.isRequired,
};

const mapStateToProps = (state) =>
  // console.log("header - mapping state to props. . . new title: "+ state.get("page").get("title"));
  // console.log("state: " + JSON.stringify(state));
  ({ title: state.get('page').get('title') })
;

// const mapDispatchToProps = dispatch => {
//   return {
//     updatePageAndTitle: bindActionCreators(updatePageAndTitle, dispatch)
//   }
// }

export default withStyles(styles)(connect(mapStateToProps)(Header));
