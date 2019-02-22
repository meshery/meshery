import React from 'react';
import PropTypes from 'prop-types';
import AppBar from '@material-ui/core/AppBar';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import HelpIcon from '@material-ui/icons/Help';
import Hidden from '@material-ui/core/Hidden';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import NotificationsIcon from '@material-ui/icons/Notifications';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Toolbar from '@material-ui/core/Toolbar';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import {connect} from "react-redux";

const lightColor = 'rgba(255, 255, 255, 0.7)';

const styles = theme => ({
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
});

class Header extends React.Component {
  // state = {
  //   title: '',
  // }

  // shouldComponentUpdate(nextProps, nextState) {
  //   console.log("shouldComponentUpdate");
  //   return (this.props.title !== nextProps.title);
  // }

  // componentDidUpdate(prevProps) {
  //   console.log("componentDidUpdate");
  //   // // Typical usage (don't forget to compare props):
  //   // if (this.props.title !== prevProps.title) {
  //   //   this.setState({title});
  //   // }
  // }

  render() {
    // console.log("header - props: " + JSON.stringify(this.props));
    const { classes, title, onDrawerToggle } = this.props;

    // console.log("header - retrieved title: "+ title);
    return (
      <React.Fragment>
        <AppBar color="primary" position="sticky" elevation={0}>
          <Toolbar>
            <Grid container spacing={8} alignItems="center">
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
              <Grid item xs />
              {/* <Grid item>
                <Typography className={classes.link} component="a" href="#">
                  Go to docs
                </Typography>
              </Grid> */}
              {/* <Grid item>
                <Tooltip title="Alerts • No alters">
                  <IconButton color="inherit">
                    <NotificationsIcon />
                  </IconButton>
                </Tooltip>
              </Grid> */}
              <Grid item>
                <IconButton color="inherit" className={classes.iconButtonAvatar}>
                  <Avatar className={classes.avatar} src="/static/images/avatar/1.jpg" />
                </IconButton>
              </Grid>
            </Grid>
          </Toolbar>
        </AppBar>
        <AppBar
          component="div"
          className={classes.secondaryBar}
          color="primary"
          position="static"
          elevation={0}
        >
          <Toolbar>
            <Grid container alignItems="center" spacing={8}>
              <Grid item xs>
                <Typography color="inherit" variant="h5">
                  {title}
                </Typography>
              </Grid>
              {/* <Grid item>
                <Button className={classes.button} variant="outlined" color="inherit" size="small">
                  Web setup
                </Button>
              </Grid> */}
              {/* <Grid item>
                <Tooltip title="Help">
                  <IconButton color="inherit">
                    <HelpIcon />
                  </IconButton>
                </Tooltip>
              </Grid> */}
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
          <Tabs value={0} textColor="inherit">
            <Tab textColor="inherit" label="Users" />
            <Tab textColor="inherit" label="Sign-in method" />
            <Tab textColor="inherit" label="Templates" />
            <Tab textColor="inherit" label="Usage" />
          </Tabs>
        </AppBar> */}
      </React.Fragment>
    );
  }
}

Header.propTypes = {
  classes: PropTypes.object.isRequired,
  onDrawerToggle: PropTypes.func.isRequired,
};

const mapStateToProps = state => {
  console.log("header - mapping state to props. . . new title: "+ state.get("page").get("title"));
  console.log("state: " + JSON.stringify(state));
  return { title: state.get("page").get("title") }
}

// const mapDispatchToProps = dispatch => {
//   return {
//     updatePageAndTitle: bindActionCreators(updatePageAndTitle, dispatch)
//   }
// }

export default withStyles(styles)(connect(
  mapStateToProps
)(Header));