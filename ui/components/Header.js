import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { NotificationDrawerButton } from './NotificationCenter';
import User from './User';
import { Search } from '@mui/icons-material';
import { deleteKubernetesConfig } from '../utils/helpers/kubernetesHelpers';
import { successHandlerGenerator, errorHandlerGenerator } from '../utils/helpers/common';
import { ConnectionChip } from './connections/ConnectionChip';
import { promisifiedDataFetch } from '../lib/data-fetch';
import _PromptComponent from './PromptComponent';
import { iconMedium, iconSmall } from '../css/icons.styles';
import ExtensionSandbox from './ExtensionSandbox';
import RemoteComponent from './RemoteComponent';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import { useNotification } from '../utils/hooks/useNotification';
import useKubernetesHook, { useControllerStatus } from './hooks/useKubernetesHook';
import { formatToTitleCase } from '../utils/utils';
import { CONNECTION_KINDS } from '../utils/Enum';
import SettingsIcon from '@mui/icons-material/Settings';
import {
  Checkbox,
  Box,
  CustomTooltip,
  Typography,
  styled,
  PROMPT_VARIANTS,
  TextField,
  ClickAwayListener,
  IconButton,
  Grid,
  Slide,
  Hidden,
  NoSsr,
  useTheme,
  useMediaQuery,
} from '@layer5/sistent';
import { CustomTextTooltip } from './MesheryMeshInterface/PatternService/CustomTextTooltip';
import { CanShow } from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import SpaceSwitcher from './SpacesSwitcher/SpaceSwitcher';
import Router from 'next/router';
import HeaderMenu from './HeaderMenu';
import ConnectionModal from './Modals/ConnectionModal';
import MesherySettingsEnvButtons from './MesherySettingsEnvButtons';
import {
  HeaderAppBar,
  UserContainer,
  PageTitleWrapper,
  CBadgeContainer,
  CMenuContainer,
  HeaderIcons,
  MenuIconButton,
  UserSpan,
  CBadge,
  StyledToolbar,
  UserInfoContainer,
} from './Header.styles';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';
import { EVENT_TYPES } from 'lib/event-types';
import { useDispatch, useSelector } from 'react-redux';
import { updateCapabilities, updateK8SConfig } from '@/store/slices/mesheryUi';

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
    <Box id={ctx.id} sx={{ margin: '0.25rem 0' }}>
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
            <>
              <Checkbox checked={selected} onChange={() => onSelectChange(ctx.id)} />
            </>
          )}
          <ConnectionChip
            title={ctx?.name}
            onDelete={onDelete ? () => onDelete(ctx.name, ctx.connection_id) : null}
            handlePing={() => ping(ctx.name, ctx.server, ctx.connection_id)}
            iconSrc={
              connectionMetadataState && connectionMetadataState[CONNECTION_KINDS.KUBERNETES]?.icon
                ? `/${connectionMetadataState[CONNECTION_KINDS.KUBERNETES]?.icon}`
                : '/static/img/kubernetes.svg'
            }
            status={operatorState}
          />
        </div>
      </CustomTextTooltip>
    </Box>
  );
};

export const K8sContextConnectionChip = K8sContextConnectionChip_;

function K8sContextMenu({
  contexts = {},
  activeContexts = [],
  setActiveContexts = () => {},
  searchContexts = () => {},
}) {
  const [anchorEl, setAnchorEl] = useState(false);
  const [showFullContextMenu, setShowFullContextMenu] = useState(false);
  const [transformProperty, setTransformProperty] = useState(100);
  const deleteCtxtRef = React.createRef();
  const { notify } = useNotification();
  const { controllerState: meshsyncControllerState } = useSelector((state) => state.ui);
  const dispatch = useDispatch();
  const { connectionMetadataState } = useSelector((state) => state.ui);

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
          dispatch(updateK8SConfig({ k8sConfig: updatedConfig }));
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
  const [isConnectionOpenModal, setIsConnectionOpenModal] = useState(false);

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
            <CBadgeContainer>
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
              <CBadge>{contexts?.total_count || 0}</CBadge>
            </CBadgeContainer>
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
                <CMenuContainer>
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
                        endAdornment: <Search style={iconMedium} width={24} />,
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
                          <>
                            <Checkbox
                              checked={activeContexts.includes('all')}
                              onChange={() =>
                                activeContexts.includes('all')
                                  ? setActiveContexts([])
                                  : setActiveContexts('all')
                              }
                            />
                          </>
                          <span style={{ fontWeight: 'bolder' }}>select all</span>
                        </div>
                        <CustomTooltip title="Configure Connections">
                          <div>
                            <IconButton size="small" onClick={() => setIsConnectionOpenModal(true)}>
                              <SettingsIcon style={{ ...iconSmall }} />
                            </IconButton>
                          </div>
                        </CustomTooltip>
                      </div>
                    )}
                    {contexts?.contexts?.map((ctx) => {
                      return (
                        <K8sContextConnectionChip
                          key={ctx.id}
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
                    <Box sx={{ marginTop: '1rem' }}>
                      <MesherySettingsEnvButtons />
                    </Box>
                  </div>
                </CMenuContainer>
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

const Header = ({
  onDrawerToggle,
  onDrawerCollapse,
  abilityUpdated,
  contexts,
  activeContexts,
  setActiveContexts,
  searchContexts,
}) => {
  const { notify } = useNotification;

  const dispatch = useDispatch();
  const {
    data: providerCapabilities,
    isSuccess: isProviderCapabilitiesSuccess,
    isError: isProviderCapabilitiesError,
    error: providerCapabilitiesError,
  } = useGetProviderCapabilitiesQuery();

  const collaboratorExt = () => {
    if (isProviderCapabilitiesSuccess) {
      dispatch(updateCapabilities({ capabilitiesRegistry: providerCapabilities }));
      return ExtensionPointSchemaValidator('collaborator')(
        providerCapabilities?.extensions?.collaborator,
      );
    }
    return null;
  };

  if (isProviderCapabilitiesError) {
    notify({
      message: 'Error fetching provider capabilities',
      event_type: EVENT_TYPES.ERROR,
      details: providerCapabilitiesError?.data,
    });
  }

  const loaderType = 'circular';
  const theme = useTheme();
  const isDesktop = useMediaQuery(theme.breakpoints.up('sm'));
  return (
    <NoSsr>
      <>
        <HeaderAppBar id="top-navigation-bar" color="primary" position="sticky">
          <StyledToolbar disableGutters isDrawerCollapsed={onDrawerCollapse}>
            <Grid container alignItems="center">
              <Hidden smUp>
                <Grid item style={{ display: 'none' }}>
                  <MenuIconButton aria-label="Open drawer" onClick={onDrawerToggle}>
                    <HeaderIcons style={iconMedium} />
                  </MenuIconButton>
                </Grid>
              </Hidden>
              <Grid item xs container alignItems="center" component={PageTitleWrapper}>
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
                <SpaceSwitcher />
              </Grid>
              <Box
                component={UserContainer}
                style={{
                  position: 'relative',
                  display: 'flex',
                  flexWrap: 'wrap',
                  gap: '1rem 0.5rem',
                }}
              >
                {/* According to the capabilities load the component */}
                {collaboratorExt && (
                  <ExtensionSandbox
                    type="collaborator"
                    Extension={(url) => RemoteComponent({ url, loaderType })}
                  />
                )}
                <UserInfoContainer>
                  <UserSpan style={{ position: 'relative' }}>
                    <K8sContextMenu
                      contexts={contexts}
                      activeContexts={activeContexts}
                      setActiveContexts={setActiveContexts}
                      searchContexts={searchContexts}
                    />
                  </UserSpan>

                  {isDesktop && (
                    <div data-testid="settings-button" aria-describedby={abilityUpdated}>
                      <CanShow Key={keys.VIEW_SETTINGS}>
                        <IconButton onClick={() => Router.push('/settings')}>
                          <SettingsIcon
                            style={{ ...iconMedium, fill: theme.palette.common.white }}
                          />
                        </IconButton>
                      </CanShow>
                    </div>
                  )}

                  <div data-testid="notification-button">
                    <NotificationDrawerButton />
                  </div>

                  <UserSpan>
                    <User />
                  </UserSpan>
                  <UserSpan data-testid="header-menu">
                    <HeaderMenu />
                  </UserSpan>
                </UserInfoContainer>
              </Box>
            </Grid>
          </StyledToolbar>
        </HeaderAppBar>
      </>
    </NoSsr>
  );
};

Header.propTypes = {
  onDrawerToggle: PropTypes.func.isRequired,
};

export default Header;
