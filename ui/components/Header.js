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
import BuildIcon from "@material-ui/icons/Build";
import MesheryNotification from './MesheryNotification';
import User from './User';
import Modal from './ConfigurationWizard/utilities/Modal'
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
  appBarOnDrawerOpen: {
    padding: theme.spacing(1.4),
    zIndex: theme.zIndex.drawer+1,
    [theme.breakpoints.between(635,732)]: {
      padding: theme.spacing(0.75,1.4),
    },
    [theme.breakpoints.between(600,635)]: {
      padding: theme.spacing(0.4,1.4),
    },
  },
  appBarOnDrawerClosed: {
    padding: theme.spacing(1.4),
    zIndex: theme.zIndex.drawer+1,
  },
  toolbarOnDrawerClosed: {
    minHeight: 59,
    paddingLeft:24,
    paddingRight:24,
  },
  toolbarOnDrawerOpen: {
    minHeight: 58,
    paddingLeft:20,
    paddingRight:20,
    [theme.breakpoints.between(620,732)]: {
      minHeight: 68,
      paddingLeft:20,
      paddingRight:20,
    },
  },
  itemActiveItem: {
    color: '#00B39F',
  },
  headerIcons:{
    fontSize: '1.5rem',
    height: '1.5rem',
    width: '1.5rem',
  }
});

class Header extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      open: false,
    };
  }
  handleOpen = () => {
    this.setState({ open: true })
  }
  handleClose = () => {
    this.setState({ open: false })
  }
  
  render() {
    const { classes, title, onDrawerToggle ,onDrawerCollapse} = this.props;
    console.log('title', title)
   
    return (
      <NoSsr>
        <React.Fragment>
          <AppBar color="primary" position="sticky" elevation={0} className={onDrawerCollapse ? classes.appBarOnDrawerClosed : classes.appBarOnDrawerOpen}>
            <Toolbar className={onDrawerCollapse ? classes.toolbarOnDrawerClosed : classes.toolbarOnDrawerOpen}>
              <Grid container alignItems="center">
                <Hidden smUp>
                  <Grid item>
                    <IconButton
                      color="inherit"
                      aria-label="Open drawer"
                      onClick={onDrawerToggle}
                      className={classes.menuButton}
                    >
                      <MenuIcon className={classes.headerIcons} />
                    </IconButton>
                  </Grid>
                </Hidden>
                <Grid item xs container alignItems="center" className={classes.pageTitleWrapper}>
                  <Typography color="inherit" variant="h5" className={classes.pageTitle}>
                    {title}
                  </Typography>
                </Grid>

                {/* <Grid item className={classes.notifications}>
                <MesheryNotification />
              </Grid> */}
                <Grid item className={classes.userContainer}>
                  {/* <IconButton color="inherit" className={classes.iconButtonAvatar}>
                  <Avatar className={classes.avatar} src="/static/images/avatar/1.jpg" />
                </IconButton> */}
                  <div data-test="index-button">
                    <IconButton color="inherit">
                      <Link href="/">
                        <DashboardIcon className={ classes.headerIcons +" "+(title === 'System Dashboard' ? classes.itemActiveItem : '')} />
                        {/* <FontAwesomeIcon icon={faHome} transform="shrink-2" fixedWidth className={title === 'Dashboard' && classes.itemActiveItem} /> */}
                      </Link>
                    </IconButton>
                  </div>

                  <div data-test="settings-button">
                    <IconButton color="inherit">
                      <Link href="/settings">
                        <SettingsIcon className={classes.headerIcons +" "+(title === 'Settings' ? classes.itemActiveItem : '')} />
                      </Link>
                    </IconButton>
                  </div>

                  <div data-test="wizard-button">
                    <IconButton color="inherit">
                      <BuildIcon onClick={this.handleOpen} className={title === 'Wizard' ? classes.itemActiveItem : ''} />
                      <Modal open={this.state.open} handleClose={this.handleClose} />
                    </IconButton>
                  </div>

                  <div data-test="notification-button">
                    <MesheryNotification />
                  </div>
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
  // console.log("state: " + JSON.stringify(state))
  ({ title: state.get('page').get('title') })
;

// const mapDispatchToProps = dispatch => {
//   return {
//     updatePageAndTitle: bindActionCreators(updatePageAndTitle, dispatch)
//   }
// }

export default withStyles(styles)(connect(mapStateToProps)(Header));
