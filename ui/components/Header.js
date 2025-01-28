import React, { useEffect } from 'react';
import PropTypes from 'prop-types';
import AppBar from '@material-ui/core/AppBar';
import Grid from '@material-ui/core/Grid';
import Hidden from '@material-ui/core/Hidden';
import IconButton from '@material-ui/core/IconButton';
import Toolbar from '@material-ui/core/Toolbar';
import { withStyles } from '@material-ui/core/styles';
import { connect, useSelector } from 'react-redux';
import NoSsr from '@material-ui/core/NoSsr';
import { NotificationDrawerButton } from './NotificationCenter';
import User from './User';
import Slide from '@material-ui/core/Slide';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import { Edit, Search } from '@mui/icons-material';
import { TextField } from '@material-ui/core';
import { Paper } from '@material-ui/core';
import { deleteKubernetesConfig } from '../utils/helpers/kubernetesHelpers';
import { successHandlerGenerator, errorHandlerGenerator } from '../utils/helpers/common';
import { _ConnectionChip } from './connections/ConnectionChip';
import { promisifiedDataFetch } from '../lib/data-fetch';
import { updateK8SConfig, updateProgress, updateCapabilities } from '../lib/store';
import { bindActionCreators } from 'redux';
import _PromptComponent from './PromptComponent';
import { iconMedium, iconSmall } from '../css/icons.styles';
import ExtensionSandbox from './ExtensionSandbox';
import RemoteComponent from './RemoteComponent';
import { CapabilitiesRegistry } from '../utils/disabledComponents';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import dataFetch from '../lib/data-fetch';
import { useNotification, withNotify } from '../utils/hooks/useNotification';
import useKubernetesHook, { useControllerStatus } from './hooks/useKubernetesHook';
import { formatToTitleCase } from '../utils/utils';
import { CONNECTION_KINDS } from '../utils/Enum';
import {
  Checkbox,
  MenuIcon,
  OutlinedSettingsIcon,
  Box,
  CustomTooltip,
  Typography,
  styled,
  PROMPT_VARIANTS,
} from '@layer5/sistent';
import { CustomTextTooltip } from './MesheryMeshInterface/PatternService/CustomTextTooltip';
import { Colors } from '@/themes/app';

import { CanShow } from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import SpaceSwitcher from './SpacesSwitcher/SpaceSwitcher';
import { UsesSistent } from './SistentWrapper';
import Router from 'next/router';
import HeaderMenu from './HeaderMenu';
import ConnectionModal from './Modals/ConnectionModal';
import MesherySettingsEnvButtons from './MesherySettingsEnvButtons';

const lightColor = 'rgba(255, 255, 255, 0.7)';
const styles = (theme) => ({
  secondaryBar: { zIndex: 0 },
  menuButton: { marginLeft: -theme.spacing(1) },
  iconButtonAvatar: { padding: 4 },
  link: {
    textDecoration: 'none',
    color: theme.palette.secondary.link,
  },
  button: { borderColor: lightColor },
  notifications: {
    paddingLeft: theme.spacing(4),
    paddingRight: theme.spacing(0),
    marginLeft: theme.spacing(4),
  },
  userContainer: {
    paddingLeft: 1,
    display: 'flex',
    alignItems: 'center',
    [theme.breakpoints.down('xs')]: {
      width: '100%',
      justifyContent: 'flex-end',
      marginRight: '1rem',
      marginBlock: '0.5rem',
    },
  },
  userSpan: { marginLeft: theme.spacing(1) },
  pageTitleWrapper: {
    flexGrow: 1,
    marginRight: 'auto',
    flexWrap: 'nowrap',
    marginBlock: '0.5rem',
  },
  appBarOnDrawerOpen: {
    backgroundColor: theme.palette.secondary.mainBackground,
    shadowColor: ' #808080',
    [theme.breakpoints.between(635, 732)]: { padding: theme.spacing(0.75, 1.4) },
    [theme.breakpoints.between(600, 635)]: { padding: theme.spacing(0.4, 1.4) },
  },
  appBarOnDrawerClosed: {
    backgroundColor: theme.palette.secondary.mainBackground,
  },
  addClusterButtonClass: {
    borderRadius: 5,
    marginRight: '2rem',
    width: '100%',
    marginTop: '1rem',
  },
  toolbarOnDrawerClosed: {
    minHeight: 59,
    padding: theme.spacing(2),
    paddingLeft: 0,
    paddingRight: 34,
    backgroundColor: theme.palette.secondary.mainBackground,
    boxShadow: `3px 0px 4px ${theme.palette.secondary.focused}`,
    [theme.breakpoints.down('xs')]: {
      padding: 0,
    },
  },
  toolbarOnDrawerOpen: {
    minHeight: 58,
    padding: theme.spacing(2),
    paddingLeft: 34,
    paddingRight: 34,
    backgroundColor: theme.palette.secondary.mainBackground,
    [theme.breakpoints.between(620, 732)]: { minHeight: 68, paddingLeft: 20, paddingRight: 20 },
    boxShadow: `3px 0px 4px ${theme.palette.secondary.focused}`,
  },
  itemActiveItem: { fill: '#00B39F' },
  headerIcons: {
    fontSize: '1.5rem',
    height: '1.5rem',
    width: '1.5rem',
    fill: theme.palette.secondary.whiteIcon,
    '&:hover': {
      fill: Colors.keppelGreen,
    },
  },
  cbadge: {
    fontSize: '0.65rem',
    backgroundColor: 'white',
    borderRadius: '50%',
    color: 'black',
    height: '1.30rem',
    width: '1.30rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    position: 'absolute',
    zIndex: 1,
    right: '-0.75rem',
    top: '-0.29rem',
  },
  cbadgeContainer: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    position: 'relative',
  },
  icon: {
    width: 24,
    height: 24,
  },
  Chip: {
    width: '12.8rem',
    textAlign: 'center',
    cursor: 'pointer',
    '& .MuiChip-label': {
      flexGrow: 1,
    },
    overflow: 'hidden',
    whiteSpace: 'nowrap',
    textOverflow: 'ellipsis',
  },
  cMenuContainer: {
    backgroundColor: theme.palette.secondary.headerColor,
    marginTop: '-0.7rem',
    borderRadius: '3px',
    padding: '1rem',
    // zIndex: 1201,
    boxShadow: '20px #979797',
    transition: 'linear .2s',
    transitionProperty: 'height',
  },
  alertEnter: {
    opacity: '0',
    transform: 'scale(0.9)',
  },
  alertEnterActive: {
    opacity: '1',
    transform: 'translateX(0)',
    transition: 'opacity 300ms, transform 300ms',
  },
  chip: {
    margin: '0.25rem 0',
  },
  AddIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  searchIcon: {
    width: theme.spacing(3.5),
  },
  darkThemeToggle: {
    marginLeft: '1.5em',
  },

  toggle: {
    appearance: 'none',
    outline: 'none',
    cursor: 'pointer',
    width: '1.5rem',
    height: '1.5rem',
    boxShadow: 'inset calc(1.5rem * 0.33) calc(1.5rem * -0.25) 0',
    borderRadius: '999px',
    color: '#00B39F',
    transition: 'all 500ms',
    zIndex: '1',
    '&:checked': {
      width: '1.5rem',
      height: '1.5rem',
      borderRadius: '50%',
      background: 'orange',
      boxShadow: '0 0 10px orange, 0 0 60px orange,0 0 200px yellow, inset 0 0 80px yellow',
    },
  },
  addButton: {
    width: '100%',
    whiteSpace: 'nowrap',
    color: theme.palette.secondary.text,
    margin: '0.5rem',
  },
  mesherySettingsEnvButtons: {
    marginTop: '1rem',
  },
});

async function loadActiveK8sContexts() {
  try {
    const res = await promisifiedDataFetch('/api/system/sync');
    if (res?.k8sConfig) {
      return res.k8sConfig;
    } else {
      throw new Error('No kubernetes configurations found');
    }
  } catch (e) {
    console.error('An error occurred while loading k8sconfig', e);
  }
}

const K8sContextConnectionChip_ = ({
  ctx,
  classes,
  selectable = false,
  onSelectChange,
  connectionMetadataState,
  meshsyncControllerState,
  selected,
  onDelete,
}) => {
  const ping = useKubernetesHook();
  const { getControllerStatesByConnectionID } = useControllerStatus(meshsyncControllerState);

  const { operatorState, meshSyncState, natsState } = getControllerStatesByConnectionID(
    ctx.connection_id,
  );

  return (
    <div id={ctx.id} className={classes.chip}>
      <CustomTextTooltip
        placement="left-end"
        leaveDelay={200}
        interactive={true}
        title={`Server: ${ctx.server},  Operator: ${formatToTitleCase(
          operatorState,
        )}, MeshSync: ${formatToTitleCase(meshSyncState)}, Broker: ${formatToTitleCase(natsState)}`}
      >
        <div
          style={{
            display: 'flex',
            justifyContent: 'flex-start',
            alignItems: 'center',
          }}
        >
          {selectable && (
            <UsesSistent>
              <Checkbox checked={selected} onChange={() => onSelectChange(ctx.id)} />
            </UsesSistent>
          )}
          <_ConnectionChip
            title={ctx?.name}
            onDelete={onDelete ? () => onDelete(ctx.name, ctx.connection_id) : null}
            handlePing={() => ping(ctx.name, ctx.server, ctx.connection_id)}
            iconSrc={
              connectionMetadataState && connectionMetadataState[CONNECTION_KINDS.KUBERNETES]?.icon
                ? `/${connectionMetadataState[CONNECTION_KINDS.KUBERNETES]?.icon}`
                : '/static/img/kubernetes.svg'
            } // TODO: Make this a dyanmic referernce to the respective Connection's color SVG
            status={operatorState}
          />
        </div>
      </CustomTextTooltip>
    </div>
  );
};

export const K8sContextConnectionChip = withStyles(styles)(K8sContextConnectionChip_);

function K8sContextMenu({
  classes = {},
  contexts = {},
  activeContexts = [],
  updateK8SConfig,
  setActiveContexts = () => {},
  searchContexts = () => {},
}) {
  const [anchorEl, setAnchorEl] = React.useState(false);
  const [showFullContextMenu, setShowFullContextMenu] = React.useState(false);
  const [transformProperty, setTransformProperty] = React.useState(100);
  const deleteCtxtRef = React.createRef();
  const { notify } = useNotification();
  const connectionMetadataState = useSelector((state) => state.get('connectionMetadataState'));
  const meshsyncControllerState = useSelector((state) => state.get('controllerState'));

  const styleSlider = {
    position: 'absolute',
    left: '-7rem',
    zIndex: '-1',
    bottom: showFullContextMenu ? '40%' : '-110%',
    transform: showFullContextMenu ? `translateY(${transformProperty}%)` : 'translateY(0)',
  };

  const StateTransitionDetails = styled(Box)(({ theme }) => ({
    backgroundColor: theme.palette.background.secondary,
    padding: '1rem',
    borderRadius: '0.5rem',
    textAlign: 'left',
  }));
  const handleKubernetesDelete = async (name, connectionID) => {
    let responseOfDeleteK8sCtx = await deleteCtxtRef.current.show({
      title: `Delete Kubernetes connection?`,
      subtitle: (
        <>
          <Typography variant="body">
            {' '}
            Are you sure you want to delete Kubernetes connection &quot;{name}&quot; and associated
            credential?
          </Typography>
          <details>
            <summary style={{ textAlign: 'left', marginTop: '1rem', cursor: 'pointer' }}>
              <strong>What does this mean?</strong>
            </summary>

            <StateTransitionDetails>
              <Typography variant="body2">
                Deleting a connection administratively removes the cluster from Meshery&apos;s
                purview of management, which includes the removal of Meshery Operator from the
                cluster. Record of this Kubernetes connection and all associated data collected
                through MeshSync for this connection will be purged from Meshery&apos;s database.
                Note: By deleting this connection, you are not deleting the Kubernetes cluster
                itself.
              </Typography>
              <Typography variant="body2" sx={{ marginTop: '1rem' }}>
                <strong>Reconnecting:</strong> You can always reconnect Meshery to the cluster
                again. By default, Meshery will automatically reconnect to the cluster when next
                presented with the same kubeconfig file / context. If you wish to prevent
                reconnection, *disconnect* this connection instead of *deleting* this connection.
              </Typography>
            </StateTransitionDetails>
          </details>
        </>
      ),
      primaryOption: 'CONFIRM',
      variant: PROMPT_VARIANTS.DANGER,
      showInfoIcon: `Learn more about the [lifecycle of connections](https://docs.meshery.io/concepts/logical/connections) and what it means to delete a connection.`,
    });
    if (responseOfDeleteK8sCtx === 'CONFIRM') {
      const successCallback = async () => {
        const updatedConfig = await loadActiveK8sContexts();
        if (Array.isArray(updatedConfig)) {
          updateK8SConfig({ k8sConfig: updatedConfig });
        }
      };
      deleteKubernetesConfig(
        successHandlerGenerator(notify, `Kubernetes connection "${name}" removed`, successCallback),
        errorHandlerGenerator(
          notify,
          `Failed to remove Kubernetes connection "
          ${name}"`,
        ),
        connectionID,
      );
    }
  };

  let open = Boolean(anchorEl);
  if (showFullContextMenu) {
    open = showFullContextMenu;
  }

  useEffect(() => {
    setTransformProperty(
      (prev) => prev + (contexts.total_count ? contexts.total_count * 3.125 : 0),
    );
  }, []);
  const [isConnectionOpenModal, setIsConnectionOpenModal] = React.useState(false);

  return (
    <>
      <div>
        <CanShow Key={keys.VIEW_ALL_KUBERNETES_CLUSTERS}>
          <IconButton
            aria-label="contexts"
            className="k8s-icon-button"
            onClick={(e) => {
              e.preventDefault();
              setShowFullContextMenu((prev) => !prev);
            }}
            onMouseOver={(e) => {
              e.preventDefault();
              setAnchorEl(true);
            }}
            onMouseLeave={(e) => {
              e.preventDefault();
              setAnchorEl(false);
            }}
            aria-owns={open ? 'menu-list-grow' : undefined}
            aria-haspopup="true"
            style={{
              marginRight: '0.5rem',
            }}
          >
            <div className={classes.cbadgeContainer}>
              <img
                className="k8s-image"
                src={
                  connectionMetadataState &&
                  connectionMetadataState[CONNECTION_KINDS.KUBERNETES]?.icon
                    ? `/${connectionMetadataState[CONNECTION_KINDS.KUBERNETES]?.icon}`
                    : '/static/img/kubernetes.svg'
                }
                onError={(e) => {
                  e.target.src = '/static/img/kubernetes.svg';
                }}
                width="24px"
                height="24px"
                style={{ objectFit: 'contain' }}
              />
              <div className={classes.cbadge}>{contexts?.total_count || 0}</div>
            </div>
          </IconButton>
        </CanShow>

        <Slide
          direction="down"
          style={styleSlider}
          timeout={400}
          in={open}
          mountOnEnter
          unmountOnExit
        >
          <div>
            <CanShow Key={keys.VIEW_ALL_KUBERNETES_CLUSTERS} invert_action={['hide']}>
              <ClickAwayListener
                onClickAway={(e) => {
                  if (
                    typeof e.target.className == 'string' &&
                    !e.target.className?.includes('cbadge') &&
                    e.target?.className != 'k8s-image' &&
                    !e.target.className.includes('k8s-icon-button')
                  ) {
                    setAnchorEl(false);
                    setShowFullContextMenu(false);
                  }
                }}
              >
                <Paper className={classes.cMenuContainer}>
                  <div>
                    <TextField
                      id="search-ctx"
                      label="Search"
                      size="small"
                      variant="outlined"
                      onChange={(ev) => searchContexts(ev.target.value)}
                      style={{
                        width: '100%',
                        backgroundColor: 'rgba(102, 102, 102, 0.12)',
                        margin: '1px 0px',
                      }}
                      InputProps={{
                        endAdornment: <Search className={classes.searchIcon} style={iconMedium} />,
                      }}
                    />
                  </div>
                  <div>
                    {contexts?.total_count > 0 && (
                      <div
                        style={{
                          display: 'flex',
                          justifyContent: 'space-between',
                          alignItems: 'center',
                        }}
                      >
                        <div>
                          <UsesSistent>
                            <Checkbox
                              checked={activeContexts.includes('all')}
                              onChange={() =>
                                activeContexts.includes('all')
                                  ? setActiveContexts([])
                                  : setActiveContexts('all')
                              }
                              color="primary"
                            />
                          </UsesSistent>
                          <span style={{ fontWeight: 'bolder' }}>select all</span>
                        </div>
                        <CustomTooltip title="Configure Connections">
                          <div>
                            <IconButton size="small" onClick={() => setIsConnectionOpenModal(true)}>
                              <Edit style={{ ...iconSmall }} />
                            </IconButton>
                          </div>
                        </CustomTooltip>
                      </div>
                    )}
                    {contexts?.contexts?.map((ctx) => {
                      return (
                        <K8sContextConnectionChip
                          key={ctx.id}
                          classes={classes}
                          ctx={ctx}
                          selectable
                          onDelete={handleKubernetesDelete}
                          selected={activeContexts.includes(ctx.id)}
                          onSelectChange={() => setActiveContexts(ctx.id)}
                          meshsyncControllerState={meshsyncControllerState}
                          connectionMetadataState={connectionMetadataState}
                        />
                      );
                    })}
                    <Box className={classes.mesherySettingsEnvButtons}>
                      <MesherySettingsEnvButtons />
                    </Box>
                  </div>
                </Paper>
              </ClickAwayListener>
            </CanShow>
          </div>
        </Slide>
      </div>
      <_PromptComponent ref={deleteCtxtRef} />
      <ConnectionModal
        isOpenModal={isConnectionOpenModal}
        setIsOpenModal={setIsConnectionOpenModal}
        meshsyncControllerState={meshsyncControllerState}
        connectionMetadataState={connectionMetadataState}
      />
    </>
  );
}

class Header extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      /** @type {CapabilityRegistryClass} */
      capabilityregistryObj: null,
      collaboratorExt: null,
    };
  }
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

          this.setState({
            collaboratorExt: ExtensionPointSchemaValidator('collaborator')(
              result?.extensions?.collaborator,
            ),
            capabilityregistryObj: capabilitiesRegistryObj,
          });
          this.props.updateCapabilities({ capabilitiesRegistry: result });
        }
      },
      (err) => console.error(err),
    );
    this._isMounted = true;
  }
  componentWillUnmount = () => {
    this._isMounted = false;
  };

  render() {
    const { classes, title, onDrawerToggle, isBeta, onDrawerCollapse, abilityUpdated } = this.props;
    const loaderType = 'circular';

    return (
      <NoSsr>
        <React.Fragment>
          <AppBar
            id="top-navigation-bar"
            color="primary"
            position="sticky"
            // elevation={1}
            className={onDrawerCollapse ? classes.appBarOnDrawerClosed : classes.appBarOnDrawerOpen}
          >
            <Toolbar
              className={
                onDrawerCollapse ? classes.toolbarOnDrawerClosed : classes.toolbarOnDrawerOpen
              }
            >
              <Grid container alignItems="center">
                <Hidden smUp>
                  <Grid item style={{ display: 'none' }}>
                    <IconButton
                      color="inherit"
                      aria-label="Open drawer"
                      onClick={onDrawerToggle}
                      className={classes.menuButton}
                    >
                      <MenuIcon className={classes.headerIcons} style={iconMedium} />
                    </IconButton>
                  </Grid>
                </Hidden>
                <Grid item xs container alignItems="center" className={classes.pageTitleWrapper}>
                  {/* Extension Point for   Logo */}
                  <div
                    id="nav-header-logo"
                    style={{
                      height: '100%',
                      display: 'flex',
                      alignItems: 'center',
                      minWidth: '34px',
                      justifyContent: 'center',
                    }}
                  ></div>
                  <SpaceSwitcher title={title} isBeta={isBeta} />
                </Grid>
                <Grid item className={classes.userContainer} style={{ position: 'relative' }}>
                  {/* According to the capabilities load the component */}
                  {this.state.collaboratorExt && (
                    <ExtensionSandbox
                      type="collaborator"
                      Extension={(url) => RemoteComponent({ url, loaderType })}
                    />
                  )}
                  <div className={classes.userSpan} style={{ position: 'relative' }}>
                    <K8sContextMenu
                      classes={classes}
                      contexts={this.props.contexts}
                      activeContexts={this.props.activeContexts}
                      setActiveContexts={this.props.setActiveContexts}
                      searchContexts={this.props.searchContexts}
                      runningStatus={{
                        operatorStatus: this.props.operatorState,
                        meshSyncStatus: this.props.meshSyncState,
                      }}
                      updateK8SConfig={this.props.updateK8SConfig}
                      updateProgress={this.props.updateProgress}
                    />
                  </div>

                  <div data-testid="settings-button" aria-describedby={abilityUpdated}>
                    <CanShow Key={keys.VIEW_SETTINGS}>
                      <IconButton onClick={() => Router.push('/settings')} color="inherit">
                        <OutlinedSettingsIcon
                          className={
                            classes.headerIcons +
                            ' ' +
                            (title === 'Settings' ? classes.itemActiveItem : '')
                          }
                          style={iconMedium}
                        />
                      </IconButton>
                    </CanShow>
                  </div>

                  <div data-testid="notification-button">
                    <NotificationDrawerButton />
                  </div>

                  <span className={classes.userSpan}>
                    <User
                      classes={classes}
                      color="inherit"
                      iconButtonClassName={classes.iconButtonAvatar}
                      avatarClassName={classes.avatar}
                      updateExtensionType={this.props.updateExtensionType}
                    />
                  </span>
                  <span className={classes.userSpan}>
                    <HeaderMenu
                      classes={classes}
                      color="inherit"
                      iconButtonClassName={classes.iconButtonAvatar}
                      avatarClassName={classes.avatar}
                      updateExtensionType={this.props.updateExtensionType}
                    />
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
  classes: PropTypes.object.isRequired,
  onDrawerToggle: PropTypes.func.isRequired,
};

const mapStateToProps = (state) => {
  return {
    title: state.get('page').get('title'),
    isBeta: state.get('page').get('isBeta'),
    selectedK8sContexts: state.get('selectedK8sContexts'),
    k8sconfig: state.get('k8sConfig'),
    operatorState: state.get('operatorState'),
    meshSyncState: state.get('meshSyncState'),
    capabilitiesRegistry: state.get('capabilitiesRegistry'),
  };
};

const mapDispatchToProps = (dispatch) => ({
  updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
  updateCapabilities: bindActionCreators(updateCapabilities, dispatch),
});

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withNotify(Header)));
