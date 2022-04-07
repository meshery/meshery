import React, { useEffect } from 'react';
import PropTypes from 'prop-types';
import AppBar from '@material-ui/core/AppBar';
import Grid from '@material-ui/core/Grid';
import Hidden from '@material-ui/core/Hidden';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import Toolbar from '@material-ui/core/Toolbar';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import { connect } from 'react-redux';
import NoSsr from '@material-ui/core/NoSsr';
import Link from 'next/link';
import SettingsIcon from '@material-ui/icons/Settings';
import Chip from '@material-ui/core/Chip';
import MesheryNotification from './MesheryNotification';
import User from './User';
import subscribeMeshSyncStatusEvents from './graphql/subscriptions/MeshSyncStatusSubscription';
import subscribeBrokerStatusEvents from "./graphql/subscriptions/BrokerStatusSubscription"
import Slide from '@material-ui/core/Slide';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import { Checkbox, Button } from '@material-ui/core';
import AddIcon from '@material-ui/icons/Add';
import { Search } from '@material-ui/icons';
import { TextField } from '@material-ui/core';
import Avatar from '@material-ui/core/Avatar';
import { Paper } from '@material-ui/core';
import { useSnackbar } from "notistack";
import { deleteKubernetesConfig, pingKubernetes } from './ConnectionWizard/helpers/kubernetesHelpers';
import {
  successHandlerGenerator, errorHandlerGenerator, closeButtonForSnackbarAction, showProgress, hideProgress
} from './ConnectionWizard/helpers/common';

const lightColor = 'rgba(255, 255, 255, 0.7)';
const styles = (theme) => ({
  secondaryBar : { zIndex : 0, },
  menuButton : { marginLeft : -theme.spacing(1), },
  iconButtonAvatar : { padding : 4, },
  link : { textDecoration : 'none',
    color : lightColor,
    '&:hover' : { color : theme.palette.common.white, }, },
  button : { borderColor : lightColor, },
  notifications : { paddingLeft : theme.spacing(4),
    paddingRight : theme.spacing(0),
    marginLeft : theme.spacing(4), },
  userContainer : { paddingLeft : 1,
    display : 'flex', backgroundColor : "#396679" },
  userSpan : { marginLeft : theme.spacing(1), },
  pageTitleWrapper : { flexGrow : 1,
    marginRight : 'auto', },
  betaBadge : { color : '#EEEEEE', fontWeight : '300', fontSize : '13px' },
  pageTitle : { paddingLeft : theme.spacing(2),
    fontSize : '1.25rem',
    [theme.breakpoints.up('sm')] : { fontSize : '1.65rem', }, },
  appBarOnDrawerOpen : {
    backgroundColor : "#396679",
    shadowColor : " #808080",
    zIndex : theme.zIndex.drawer+1,
    [theme.breakpoints.between(635,732)] : { padding : theme.spacing(0.75,1.4), },
    [theme.breakpoints.between(600,635)] : { padding : theme.spacing(0.4,1.4), },
  },
  appBarOnDrawerClosed : {
    backgroundColor : "#396679",
    zIndex : theme.zIndex.drawer+1, },
  toolbarOnDrawerClosed : { minHeight : 59,
    padding : theme.spacing(2.4),
    paddingLeft : 34,
    paddingRight : 34,
    backgroundColor : "#396679" },
  toolbarOnDrawerOpen : {
    minHeight : 58,
    padding : theme.spacing(2.4),
    paddingLeft : 34,
    paddingRight : 34,
    backgroundColor : "#396679",
    [theme.breakpoints.between(620, 732)] : { minHeight : 68, paddingLeft : 20, paddingRight : 20 },
  },
  itemActiveItem : { color : "#00B39F" },
  headerIcons : { fontSize : "1.5rem", height : "1.5rem", width : "1.5rem" },
  cbadge : {
    fontSize : "0.65rem",
    backgroundColor : "white",
    borderRadius : "50%",
    color : "black",
    height : "1.30rem",
    width : "1.30rem",
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    position : "absolute",
    zIndex : 1,
    right : "-0.75rem",
    top : "-0.29rem"
  },
  cbadgeContainer : {
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    position : "relative"
  },
  icon : {
    width : theme.spacing(2.5)
  },
  Chip : {
    backgroundColor : "white",
    cursor : "pointer"
  },
  cMenuContainer : {
    backgroundColor : "#EEEEEE",
    borderRadius : "3px",
    padding : "1rem",
    zIndex : 1201,
    marginTop : "1.3rem",
    boxShadow : "20px #979797",
    transition : "linear .2s",
    transitionProperty : "height"
  },
  alertEnter : {
    opacity : "0",
    transform : "scale(0.9)",
  },
  alertEnterActive : {
    opacity : "1",
    transform : "translateX(0)",
    transition : "opacity 300ms, transform 300ms"
  },
  chip : {
    margin : "0.25rem 0",
  },
  AddIcon : {
    width : theme.spacing(2.5),
    paddingRight : theme.spacing(0.5),
  },
  searchIcon : {
    width : theme.spacing(3.5),
  }
});

function K8sContextMenu({
  classes = {},
  contexts = {},
  activeContexts = [],
  setActiveContexts = () => {},
  searchContexts = () => {}
}) {
  const [anchorEl, setAnchorEl] = React.useState(false);
  const [showFullContextMenu, setShowFullContextMenu] = React.useState(false);
  const [transformProperty, setTransformProperty] = React.useState(100)
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  const handleKubernetesClick = () => {
    showProgress()
    pingKubernetes(
      successHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Kubernetes succesfully pinged", () => hideProgress()),
      errorHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Kubernetes not pinged successfully", () => hideProgress())
    )

  }
  const handleKubernetesDelete = () => {
    showProgress()
    deleteKubernetesConfig(
      successHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Kubernetes config successfully removed"),
      errorHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Not able to remove config")
    )
  }


  let open = Boolean(anchorEl);
  if (showFullContextMenu) {
    open = showFullContextMenu;
  }


  useEffect(() => {
    setTransformProperty(prev => (prev + ( contexts.total_count ? contexts.total_count * 3.125 : 0 )))
  }, [])

  return (
    <>
      <IconButton
        aria-label="contexts"
        className="k8s-icon-button"
        onClick={(e) => {
          e.preventDefault();
          console.log(contexts);
          setShowFullContextMenu(prev => !prev);
        }}
        onMouseOver={(e) => {
          e.preventDefault();
          console.log(contexts);
          setAnchorEl(true);
        }}

        onMouseLeave={(e) => {
          e.preventDefault();
          console.log(contexts);
          setAnchorEl(false)
        }}

        aria-owns={open
          ? 'menu-list-grow'
          : undefined}
        aria-haspopup="true"
        style={{ marginRight : "0.5rem" }}
      >
        <div className={classes.cbadgeContainer}>
          <img className="k8s-image" src="/static/img/kubernetes.svg" width="24px" height="24px" style={{ zIndex : "2" }} />
          <div className={classes.cbadge}>{contexts?.total_count || 0}</div>
        </div>
      </IconButton>

      <Slide direction="down"  timeout={400} in={open} style={{ position : "absolute", left : "-5rem", zIndex : "-1", bottom : "-100%", transform : showFullContextMenu ? `translateY(${transformProperty}%)`: "translateY(0)" }} mountOnEnter unmountOnExit>
        <div>
          <ClickAwayListener onClickAway={(e) => {

            if (!e.target.className.includes("cbadge") && e.target.className !="k8s-image" && !e.target.className.includes("k8s-icon-button")) {
              setAnchorEl(false)
              setShowFullContextMenu(false)
            }
          }}>
            <Paper style={{ marginTop : "-1.5rem", }} className={classes.cMenuContainer}>
              <div>
                <TextField
                  id="search-ctx"
                  label="Search"
                  size="small"
                  variant="outlined"
                  onChange={ev => searchContexts(ev.target.value)}
                  style={{ width : "100%", backgroundColor : "rgba(102, 102, 102, 0.12)", margin : "1px 0px" }}
                  InputProps={{ endAdornment :
                (
                  <Search className={classes.searchIcon} />
                ) }}
                />
              </div>
              <div>
                {
                  contexts?.total_count
                    ?
                    <>
                      <Checkbox
                        checked={activeContexts.includes(".all")}
                        onChange={() => setActiveContexts(".all")}
                        color="primary"
                      />
                      <span>Select All</span>
                    </>
                    :
                    <Link href="/settings">
                      <Button
                        type="submit"
                        variant="contained"
                        color="primary"
                        size="large"
                        style={{ margin : "0.5rem 0.5rem", whiteSpace : "nowrap" }}
                      >
                        <AddIcon className={classes.AddIcon}/>
                      Connect Clusters
                      </Button>
                    </Link>
                }
                {contexts?.contexts?.map(ctx => (
                  <div id={ctx.id} className={classes.chip}>
                    <Tooltip title={`Server: ${ctx.server}`}>
                      <div style={{ display : "flex", justifyContent : "center", alignItems : "center" }}>
                        <Checkbox
                          checked={activeContexts.includes(ctx.id)}
                          onChange={() => setActiveContexts(ctx.id)}
                          color="primary"
                        />
                        <Chip
                          label={ctx?.name}
                          onDelete={handleKubernetesDelete}
                          onClick={handleKubernetesClick}
                          avatar={<Avatar src="/static/img/kubernetes.svg" className={classes.icon} />}
                          variant="filled"
                          className={classes.Chip}
                          data-cy="chipContextName"
                        />
                      </div>
                    </Tooltip>
                  </div>
                ))}

              </div>
            </Paper>

          </ClickAwayListener>
        </div>
      </Slide>
    </>
  )
}

class Header extends React.Component {
  constructor(props){
    super(props);
    this.state = {
      meshSyncStatusSubscription : null,
      brokerStatusSubscription : null,
      meshSyncStatus : {
        name : "meshsync",
        status : "DISABLED",
        version : "v0.0.0",
        error : null,
      },
      brokerStatus : false,
    }
  }
  componentDidMount(){
    this._isMounted = true;
    const meshSyncStatusSub = subscribeMeshSyncStatusEvents(data => this.setState({ meshSyncStatus : data?.listenToMeshSyncEvents }));
    const brokerStatusSub = subscribeBrokerStatusEvents(data => {
      console.log({ brokerData : data })
      this.setState({ brokerStatus : data?.subscribeBrokerConnection })
    });
    this.setState({ meshSyncStatusSubscription : meshSyncStatusSub, brokerStatusSubscription : brokerStatusSub })
  }

  componentWillUnmount = () => {
    this._isMounted = false;
    this.disposeSubscriptions()
  }

  disposeSubscriptions = () => {
    if (this.state.meshSyncStatusSubscription) {
      this.state.meshSyncStatusSubscription.dispose();
    }
  }

  render() {
    const {
      classes, title, onDrawerToggle ,onDrawerCollapse ,isBeta
    } = this.props;
    return (
      <NoSsr>
        <React.Fragment>
          <AppBar color="primary" position="sticky" elevation={2} className={onDrawerCollapse
            ? classes.appBarOnDrawerClosed
            : classes.appBarOnDrawerOpen}>
            <Toolbar className={onDrawerCollapse
              ? classes.toolbarOnDrawerClosed
              : classes.toolbarOnDrawerOpen}>
              <Grid container alignItems="center" style={{ backgroundColor : "#396679" }}>
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
                  <Typography color="inherit" variant="h5" className={classes.pageTitle} data-cy="headerPageTitle">
                    {title}{isBeta ? <sup className={classes.betaBadge}>BETA</sup> : ""}
                  </Typography>
                </Grid>

                {/* <Grid item className={classes.notifications}>
                <MesheryNotification />
              </Grid> */}
                <Grid item className={classes.userContainer}>
                  {/* <IconButton color="inherit" className={classes.iconButtonAvatar}>
                  <Avatar className={classes.avatar} src="/static/images/avatar/1.jpg" />
                </IconButton>
                  <div data-test="index-button">
                    <IconButton color="inherit">
                      <Link href="/">
                        <DashboardIcon className={ classes.headerIcons +" "+(title === 'System Dashboard' ? classes.itemActiveItem : '')} />
                        <FontAwesomeIcon icon={faHome} transform="shrink-2" fixedWidth className={title === 'Dashboard' && classes.itemActiveItem} />
                      </Link>
                    </IconButton>
                  </div>*/}
                  {/* <div data-test="connection-wizard-button">
                    <IconButton color="inherit">
                      <Link href="/system/connections">
                        <img src={title === 'Connection Wizard'
                          ? "/static/img/connection_wizard/connection-wizard-green.svg"
                          : "/static/img/connection_wizard/connection-wizard-white.svg"} className={ classes.headerIcons +" "+(title === 'Connection Wizard'
                          ? classes.itemActiveItem
                          : '')} />
                        {/* <FontAwesomeIcon icon={faHome} transform="shrink-2" fixedWidth className={title === 'Dashboard' && classes.itemActiveItem} /> */}
                  {/* </Link>
                    </IconButton>
                  </div> */}
                  <div className={classes.userSpan} style={{ position : "relative" }}>
                    <K8sContextMenu
                      classes={classes}
                      contexts={this.props.contexts}
                      activeContexts={this.props.activeContexts}
                      setActiveContexts={this.props.setActiveContexts}
                      searchContexts={this.props.searchContexts}
                    />
                  </div>

                  <div data-test="settings-button">
                    <IconButton color="inherit">
                      <Link href="/settings">
                        <SettingsIcon className={classes.headerIcons +" "+(title === 'Settings'
                          ? classes.itemActiveItem
                          : '')} />
                      </Link>
                    </IconButton>
                  </div>

                  <div data-test="notification-button">
                    <MesheryNotification />
                  </div>

                  <Tooltip title={this.state?.meshSyncStatus?.status === "ENABLED" ? "Active" : "Inactive" }>
                    <IconButton>
                      <Link href="/settings#environment">
                        <img className={classes.headerIcons} src={this.state?.meshSyncStatus?.status === "ENABLED" ? "/static/img/meshsync.svg" : "/static/img/meshsync-white.svg"} />
                      </Link>
                    </IconButton>
                  </Tooltip>

                  {/* <Tooltip title="Broker Status">
                    <div style={{ padding : "1rem", height : "2rem", width : "2rem", borderRadius : "50%", backgroundColor : this.state.brokerStatus ? "green" : "red" }} />
                  </Tooltip>  */}

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

Header.propTypes = { classes : PropTypes.object.isRequired,
  onDrawerToggle : PropTypes.func.isRequired, };

const mapStateToProps = (state) =>
  // console.log("header - mapping state to props. . . new title: "+ state.get("page").get("title"));
  // console.log("state: " + JSON.stringify(state));
  ({ title : state.get('page').get('title'), isBeta : state.get('page').get('isBeta') })
;

// const mapDispatchToProps = dispatch => {
//   return {
//     updatePageAndTitle: bindActionCreators(updatePageAndTitle, dispatch)
//   }
// }

export default withStyles(styles)(connect(mapStateToProps)(Header));
