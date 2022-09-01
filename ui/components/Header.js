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
import subscribeBrokerStatusEvents from "./graphql/subscriptions/BrokerStatusSubscription"
import mesheryControllersStatusSubcription from "./graphql/subscriptions/MesheryControllersStatusSubscription"
import meshSyncEventsSub from "./graphql/subscriptions/MeshSyncEventsSubscription"
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
  successHandlerGenerator, errorHandlerGenerator, closeButtonForSnackbarAction } from './ConnectionWizard/helpers/common';
import { getFirstCtxIdFromSelectedCtxIds } from '../utils/multi-ctx';
import { promisifiedDataFetch } from '../lib/data-fetch';
import { updateK8SConfig, updateProgress } from '../lib/store';
import { bindActionCreators } from 'redux';
import BadgeAvatars from './CustomAvatar';
const lightColor = 'rgba(255, 255, 255, 0.7)';
const styles = (theme) => ({
  secondaryBar : { zIndex : 0, },
  menuButton : { marginLeft : -theme.spacing(1), },
  iconButtonAvatar : { padding : 4, },
  link : {
    textDecoration : 'none',
    color : lightColor,
    '&:hover' : { color : theme.palette.common.white, },
  },
  button : { borderColor : lightColor, },
  notifications : {
    paddingLeft : theme.spacing(4),
    paddingRight : theme.spacing(0),
    marginLeft : theme.spacing(4),
  },
  userContainer : {
    paddingLeft : 1,
    display : 'flex',
    backgroundColor : "#396679",
    alignItems : 'center'
  },
  userSpan : { marginLeft : theme.spacing(1), },
  pageTitleWrapper : {
    flexGrow : 1,
    marginRight : 'auto',
  },
  betaBadge : { color : '#EEEEEE', fontWeight : '300', fontSize : '13px' },
  pageTitle : {
    paddingLeft : theme.spacing(2),
    fontSize : '1.25rem',
    [theme.breakpoints.up('sm')] : { fontSize : '1.65rem', },
  },
  appBarOnDrawerOpen : {
    backgroundColor : "#396679",
    shadowColor : " #808080",
    zIndex : theme.zIndex.drawer + 1,
    [theme.breakpoints.between(635, 732)] : { padding : theme.spacing(0.75, 1.4), },
    [theme.breakpoints.between(600, 635)] : { padding : theme.spacing(0.4, 1.4), },
  },
  appBarOnDrawerClosed : {
    backgroundColor : "#396679",
    zIndex : theme.zIndex.drawer + 1,
  },
  toolbarOnDrawerClosed : {
    minHeight : 59,
    padding : theme.spacing(2.4),
    paddingLeft : 34,
    paddingRight : 34,
    backgroundColor : "#396679"
  },
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
    width : 24,
    height : 24
  },
  Chip : {
    backgroundColor : "white",
    flexGrow : 1,
    cursor : "pointer",
    "& .MuiChip-label" : {
      flexGrow : 1
    }
  },
  cMenuContainer : {
    backgroundColor : "revert",
    marginTop : "-0.7rem",
    borderRadius : "3px",
    padding : "1rem",
    zIndex : 1201,
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

async function loadActiveK8sContexts() {
  try {
    const res = await promisifiedDataFetch("/api/system/sync")
    if (res?.k8sConfig) {
      return res.k8sConfig
    } else {
      throw new Error("No kubernetes configurations found")
    }
  } catch (e) {
    console.error("An error occurred while loading k8sconfig", e)
  }
}

function K8sContextMenu({
  classes = {},
  contexts = {},
  activeContexts = [],
  runningStatus,
  updateK8SConfig,
  updateProgress,

  setActiveContexts = () => { },
  searchContexts = () => { }
}) {
  const [anchorEl, setAnchorEl] = React.useState(false);
  const [showFullContextMenu, setShowFullContextMenu] = React.useState(false);
  const [transformProperty, setTransformProperty] = React.useState(100)
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  const styleSlider = {
    backgroundColor : "#EEEEEE",
    position : "absolute",
    left : "-5rem",
    zIndex : "-1",
    bottom : "-55%",
    transform : showFullContextMenu ? `translateY(${transformProperty}%)` : "translateY(0)"
  }

  const getOperatorStatus = (contextId) => {
    const state = runningStatus.operatorStatus;
    if (!state) {
      return false;
    }

    const context = state.find(st => st.contextID === contextId)
    if (!context) {
      return false;
    }

    return context.operatorStatus.status === "ENABLED";
  }

  const getMeshSyncStatus = (contextId) => {
    const state = runningStatus.meshSyncStatus;
    if (!state) {
      return false;
    }

    const context = state.find(st => st.contextID === contextId)
    if (!context) {
      return false;
    }

    return context.OperatorControllerStatus.status?.includes("ENABLED");
  }

  const handleKubernetesClick = (id) => {
    updateProgress({ showProgress : true })

    pingKubernetes(
      successHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Kubernetes succesfully pinged", () => updateProgress({ showProgress : false })),
      errorHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Kubernetes not pinged successfully", () => updateProgress({ showProgress : false })),
      id
    )
  }

  const handleKubernetesDelete = (name, ctxId) => () => {
    if (contexts?.total_count === 1) {
      alert("You have only one kubernetes context connected to Meshery, Aborting delete operation");
      return;
    }

    if (confirm(`Are you sure you want to delete "${name}" cluster from Meshery?`)) {
      const successCallback = async () => {
        const updatedConfig = await loadActiveK8sContexts()
        if (Array.isArray(updatedConfig)) {
          updateK8SConfig({ k8sConfig : updatedConfig })
        }
      }
      deleteKubernetesConfig(
        successHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Kubernetes config successfully removed", successCallback),
        errorHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Not able to remove config"),
        ctxId
      )
    }
  }


  let open = Boolean(anchorEl);
  if (showFullContextMenu) {
    open = showFullContextMenu;
  }

  useEffect(() => {
    setTransformProperty(prev => (prev + (contexts.total_count ? contexts.total_count * 3.125 : 0)))
  }, [])

  return (
    <>
      <IconButton
        aria-label="contexts"
        className="k8s-icon-button"
        onClick={(e) => {
          e.preventDefault();
          setShowFullContextMenu(prev => !prev);
        }}
        onMouseOver={(e) => {
          e.preventDefault();
          setAnchorEl(true);
        }}

        onMouseLeave={(e) => {
          e.preventDefault();
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

      <Slide direction="down" style={styleSlider} timeout={400} in={open} mountOnEnter unmountOnExit>
        <div>
          <ClickAwayListener onClickAway={(e) => {

            if (!e.target.className?.includes("cbadge") && e.target?.className != "k8s-image" && !e.target.className.includes("k8s-icon-button")) {
              setAnchorEl(false)
              setShowFullContextMenu(false)
            }
          }}>

            <Paper className={classes.cMenuContainer}>
              <div>
                <TextField
                  id="search-ctx"
                  label="Search"
                  size="small"
                  variant="outlined"
                  onChange={ev => searchContexts(ev.target.value)}
                  style={{ width : "100%", backgroundColor : "rgba(102, 102, 102, 0.12)", margin : "1px 0px" }}
                  InputProps={{
                    endAdornment :
                      (
                        <Search className={classes.searchIcon} />
                      )
                  }}
                />
              </div>
              <div>
                {
                  contexts?.total_count
                    ?
                    <>
                      <Checkbox
                        checked={activeContexts.includes("all")}
                        onChange={() => setActiveContexts("all")}
                        color="primary"
                      />
                      <span style={{ fontWeight : "bolder" }}>select all</span>
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
                        <AddIcon className={classes.AddIcon} />
                        Connect Clusters
                      </Button>
                    </Link>
                }
                {contexts?.contexts?.map(ctx => {
                  const meshStatus = getMeshSyncStatus(ctx.id);
                  const operStatus = getOperatorStatus(ctx.id);

                  function getStatus(status) {
                    if (status) {
                      return "Active"
                    } else {
                      return "InActive"
                    }
                  }

                  return <div id={ctx.id} className={classes.chip}>
                    <Tooltip title={`Server: ${ctx.server}, Meshsync: ${getStatus(meshStatus)}, Operator: ${getStatus(operStatus)}`}>
                      <div style={{ display : "flex", justifyContent : "flex-start", alignItems : "center" }}>
                        <Checkbox
                          checked={activeContexts.includes(ctx.id)}
                          onChange={() => setActiveContexts(ctx.id)}
                          color="primary"
                        />
                        <Chip
                          label={ctx?.name}
                          onDelete={handleKubernetesDelete(ctx.name, ctx.id)}
                          onClick={() => handleKubernetesClick(ctx.id)}
                          avatar={
                            meshStatus ?
                              <BadgeAvatars>
                                <Avatar src="/static/img/kubernetes.svg" className={classes.icon}
                                  style={operStatus ? {} : { opacity : 0.2 }}
                                />
                              </BadgeAvatars> :
                              <Avatar src="/static/img/kubernetes.svg" className={classes.icon}
                                style={operStatus ? {} : { opacity : 0.2 }}
                              />
                          }
                          variant="filled"
                          className={classes.Chip}
                          data-cy="chipContextName"
                        />
                      </div>
                    </Tooltip>
                  </div>
                })}

              </div>
            </Paper>

          </ClickAwayListener>
        </div>
      </Slide>
    </>
  )
}

class Header extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      brokerStatusSubscription : null,
      brokerStatus : false,
    }
  }

  componentDidMount() {
    this._isMounted = true;
    const brokerStatusSub = subscribeBrokerStatusEvents(data => {
      console.log({ brokerData : data })
      this.setState({ brokerStatus : data?.subscribeBrokerConnection })
    });
    this.setState({ brokerStatusSubscription : brokerStatusSub })

    mesheryControllersStatusSubcription(data => {
      console.log({ status : data })
      // this.setState({ brokerStatus: data?.subscribeBrokerConnection })
    });

    meshSyncEventsSub(data => {
      console.log({ event : data })
      // this.setState({ brokerStatus: data?.subscribeBrokerConnection })
    });
  }



  getSelectedContextId = () => {
    return getFirstCtxIdFromSelectedCtxIds(["all"], this.props.k8sconfig)
  }

  componentWillUnmount = () => {
    this._isMounted = false;
  }


  render() {
    const { classes, title, onDrawerToggle, onDrawerCollapse, isBeta } = this.props;
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
                <Grid item className={classes.userContainer}>
                  <div className={classes.userSpan} style={{ position : "relative" }}>
                    <K8sContextMenu
                      classes={classes}
                      contexts={this.props.contexts}
                      activeContexts={this.props.activeContexts}
                      setActiveContexts={this.props.setActiveContexts}
                      searchContexts={this.props.searchContexts}
                      runningStatus={{ operatorStatus : this.props.operatorState, meshSyncStatus : this.props.meshSyncState }}
                      updateK8SConfig={this.props.updateK8SConfig}
                      updateProgress={this.props.updateProgress}
                    />
                  </div>

                  <div data-test="settings-button">
                    <IconButton color="inherit">
                      <Link href="/settings">
                        <SettingsIcon className={classes.headerIcons + " " + (title === 'Settings'
                          ? classes.itemActiveItem
                          : '')} />
                      </Link>
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
        </React.Fragment>
      </NoSsr>
    );
  }
}

Header.propTypes = {
  classes : PropTypes.object.isRequired,
  onDrawerToggle : PropTypes.func.isRequired,
};

const mapStateToProps = (state) => {
  return ({
    title : state.get('page').get('title'),
    isBeta : state.get('page').get('isBeta'),
    selectedK8sContexts : state.get('selectedK8sContexts'),
    k8sconfig : state.get('k8sConfig'),
    operatorState : state.get('operatorState'),
    meshSyncState : state.get('meshSyncState')
  })
};

const mapDispatchToProps = (dispatch) => ({
  updateK8SConfig : bindActionCreators(updateK8SConfig, dispatch),
  updateProgress : bindActionCreators(updateProgress, dispatch),
});


export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(Header));
