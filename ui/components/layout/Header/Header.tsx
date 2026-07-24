import React, { useState, useContext, useMemo } from 'react';
import PropTypes from 'prop-types';
import { NotificationDrawerButton } from '../NotificationCenter/index';
import User from '../../User';
import { successHandlerGenerator, errorHandlerGenerator } from '../../../utils/helpers/common';
import { ConnectionChip } from '../../connections/ConnectionChip';
import { normalizeStaticImagePath } from '../../../utils/fallback';
import { useLazyGetSystemSyncQuery } from '../../../rtk-query/system';
import { useUpdateConnectionStatusMutation } from '../../../rtk-query/connection';
import { CONNECTION_KINDS, CONNECTION_STATES } from '../../../utils/Enum';
import ConnectionStateTransitionModal from '../../connections/ConnectionStateTransitionModal';
import type { ConnectionStateTransitionModalRef } from '../../connections/ConnectionStateTransitionModal';
import { iconMedium, iconSmall } from '../../../css/icons.styles';
import { createPathForRemoteComponent } from '../../ExtensionSandbox';
import RemoteComponent from '../../general/RemoteComponent';
import { useNotification } from '../../../utils/hooks/useNotification';
import useKubernetesHook, { useControllerStatus } from '@/utils/hooks/useKubernetesHook';
import { formatToTitleCase } from '../../../utils/utils';
import RegistryModal from '../../registry/RegistryModal';

import {
  Checkbox,
  Box,
  CustomTooltip,
  TextField,
  ClickAwayListener,
  IconButton,
  Slide,
  Grid2,
  Hidden,
  NoSsr,
  useTheme,
  useMediaQuery,
  SearchIcon,
  SettingsIcon,
  FilterAllIcon,
  useHasPermission,
} from '@sistent/sistent';
import { Keys } from '@meshery/schemas/permissions';
import OrganizationAndWorkSpaceSwitcher from '../../workspaces/SpacesSwitcher/SpaceSwitcher';
import HeaderMenu from './HeaderMenu';
import ConnectionModal from '../../connections/ConnectionFormModal';
import MesherySettingsEnvButtons from '../../MesherySettingsEnvButtons';
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
import {
  getUserAccessToken,
  getUserProfile,
  useGetProviderCapabilitiesQuery,
} from '@/rtk-query/user';
import { useGetConnectionsQuery } from '@/rtk-query/connection';
import { EVENT_TYPES } from 'lib/event-types';
import { useDispatch, useSelector } from 'react-redux';
import { updateK8SConfig } from '@/store/slices/mesheryUi';
import { ErrorBoundary } from '@sistent/sistent';
import { WorkspaceModalContext } from '../../../utils/context/WorkspaceModalContextProvider';

const K8sContextConnectionChip_ = ({
  ctx,
  selectable = false,
  onSelectChange,
  connectionMetadataState,
  meshsyncControllerState,
  selected,
  onDelete,
  connections = [],
}) => {
  const ping = useKubernetesHook();
  const { getControllerStatesByConnectionID } = useControllerStatus(meshsyncControllerState);

  const { operatorState, meshSyncState, natsState } = getControllerStatesByConnectionID(
    ctx.connectionId,
  );

  // Prefer status already mapped onto the context (connectionsToK8sContexts
  // sets `connectionStatus`). Fall back to a connections-list lookup for
  // callers that only pass connectionId + a connections array.
  const connectionStatus = useMemo(() => {
    if (ctx.connectionStatus) {
      return ctx.connectionStatus;
    }
    if (!connections?.length || !ctx.connectionId) {
      return null;
    }
    const connection = connections.find((conn) => conn.id === ctx.connectionId);
    return connection?.status || null;
  }, [connections, ctx.connectionId, ctx.connectionStatus]);

  return (
    <Box id={ctx.id} sx={{ margin: '0.25rem 0' }}>
      <CustomTooltip
        placement="left-end"
        leaveDelay={200}
        interactive={true}
        title={`Server: ${ctx.server}, Connection: ${formatToTitleCase(
          connectionStatus || 'Unknown',
        )}, Operator: ${formatToTitleCase(
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
            onDelete={onDelete ? () => onDelete(ctx.name, ctx.connectionId) : null}
            handlePing={() => ping(ctx.name, ctx.server, ctx.connectionId)}
            // Pass the raw icon (inline SVG markup from the connection
            // definition's styles.svgColor, or a path) — ConnectionChip runs it
            // through normalizeStaticImagePath, which turns SVG markup into a
            // data URI. Prefixing with "/" here would corrupt the SVG markup.
            iconSrc={
              connectionMetadataState?.[CONNECTION_KINDS.KUBERNETES]?.icon ||
              '/static/img/integrations/kubernetes.svg'
            }
            status={connectionStatus}
          />
        </div>
      </CustomTooltip>
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
  const theme = useTheme();
  const hasK8sPermission = useHasPermission(Keys.IdentityAccessManagementViewAllKubernetesClusters);
  const [showFullContextMenu, setShowFullContextMenu] = useState(false);
  const anchorRef = React.useRef(null);
  // The dropdown slides up from below; its translate distance scales with the
  // number of context rows it will render so it ends up flush against the badge.
  // useRef (not createRef) so the same ref instance survives re-renders.
  const deleteCtxtRef = React.useRef<ConnectionStateTransitionModalRef | null>(null);
  const { notify } = useNotification();
  const [fetchSystemSync] = useLazyGetSystemSyncQuery();
  const [updateConnectionStatus] = useUpdateConnectionStatusMutation();
  const { controllerState: meshsyncControllerState, connectionMetadataState } = useSelector(
    (state) => state.ui,
  );
  const dispatch = useDispatch();

  // Same filter shape as KubernetesSubscription / the connections table:
  // plain kind=kubernetes (not JSON-encoded) and pageSize=all so status dots
  // resolve for every cluster in the switcher.
  const { data: connectionData } = useGetConnectionsQuery({
    kind: CONNECTION_KINDS.KUBERNETES,
    pageSize: 'all',
  });

  const connections = connectionData?.connections || [];

  const styleSlider = {
    position: 'absolute',
    left: '-7rem',
    zIndex: '-1',
    top: '60px',
  };

  const handleKubernetesDelete = async (name, connectionID) => {
    // The shared transition modal explains what deleting this connection means
    // by resolving the Kubernetes connection definition's transition copy for
    // (current state → deleted); pass the connection's current state so that
    // lookup can resolve (connectionsToK8sContexts maps it onto each ctx).
    const currentStatus = contexts?.contexts?.find(
      (ctx) => ctx.connectionId === connectionID,
    )?.connectionStatus;
    const confirmed = await deleteCtxtRef.current?.show({
      targetStatus: CONNECTION_STATES.DELETED,
      kind: CONNECTION_KINDS.KUBERNETES,
      currentStatus,
      connections: [{ id: connectionID, name, status: currentStatus }],
    });
    if (confirmed) {
      const successCallback = async () => {
        try {
          const res = await fetchSystemSync().unwrap();
          if (Array.isArray(res?.k8sConfig)) {
            dispatch(updateK8SConfig({ k8sConfig: res.k8sConfig }));
          }
        } catch (e) {
          console.error('An error occurred while loading k8sconfig', e);
        }
      };
      try {
        await updateConnectionStatus({
          kind: CONNECTION_KINDS.KUBERNETES,
          body: { [connectionID]: CONNECTION_STATES.DELETED },
        }).unwrap();
        successHandlerGenerator(
          notify,
          `Kubernetes connection "${name}" removed`,
          successCallback,
        )();
      } catch (err) {
        errorHandlerGenerator(notify, `Failed to remove Kubernetes connection "${name}"`)(err);
      }
    }
  };

  const [isConnectionOpenModal, setIsConnectionOpenModal] = useState(false);

  return (
    <>
      <div>
        <IconButton
          ref={anchorRef}
          aria-label="contexts"
          className="k8s-icon-button"
          onClick={(e) => {
            e.preventDefault();
            setShowFullContextMenu((prev) => !prev);
          }}
          aria-controls={showFullContextMenu ? 'menu-list-grow' : undefined}
          aria-haspopup="true"
          style={{
            marginRight: '0.5rem',
          }}
          permissionKey={Keys.IdentityAccessManagementViewAllKubernetesClusters}
          permissionAction="hide"
        >
          <CBadgeContainer>
            <img
              className="k8s-image"
              src={
                normalizeStaticImagePath(
                  connectionMetadataState?.[CONNECTION_KINDS.KUBERNETES]?.icon,
                ) || '/static/img/integrations/kubernetes.svg'
              }
              onError={(e) => {
                e.target.src = '/static/img/integrations/kubernetes.svg';
              }}
              width="24px"
              height="24px"
              style={{ objectFit: 'contain' }}
            />
            <CBadge
              onClick={(e) => {
                e.stopPropagation();
                setShowFullContextMenu((prev) => !prev);
              }}
            >
              {contexts?.totalCount || 0}
            </CBadge>
          </CBadgeContainer>
        </IconButton>

        <Slide
          direction="down"
          style={styleSlider}
          timeout={400}
          in={showFullContextMenu}
          mountOnEnter
          unmountOnExit
        >
          <div>
            {hasK8sPermission && (
              <ClickAwayListener
                onClickAway={(e) => {
                  if (anchorRef.current && anchorRef.current.contains(e.target as Node)) {
                    return;
                  }
                  setShowFullContextMenu(false);
                }}
              >
                <CMenuContainer id="menu-list-grow">
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
                        endAdornment: <SearchIcon style={iconMedium} width={24} />,
                      }}
                    />
                  </div>
                  <div>
                    {contexts?.totalCount > 0 && (
                      <div
                        style={{
                          display: 'flex',
                          justifyContent: 'space-between',
                          alignItems: 'center',
                          marginTop: '1rem',
                        }}
                      >
                        <div style={{ display: 'flex', alignItems: 'center' }}>
                          <>
                            <Checkbox
                              checked={activeContexts.includes('all')}
                              onChange={() =>
                                activeContexts.includes('all')
                                  ? setActiveContexts([])
                                  : setActiveContexts('all')
                              }
                              icon={
                                <FilterAllIcon
                                  fill={theme.palette.background.brand.default}
                                  style={{ opacity: 0.4 }}
                                />
                              }
                              inputProps={{ 'aria-label': 'select all contexts' }}
                            />
                          </>
                          <span style={{ fontWeight: 'bolder', whiteSpace: 'nowrap' }}>
                            select all
                          </span>
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
                          connections={connections}
                        />
                      );
                    })}
                    <Box sx={{ marginTop: '1rem' }}>
                      <MesherySettingsEnvButtons onOpened={() => setShowFullContextMenu(false)} />
                    </Box>
                  </div>
                </CMenuContainer>
              </ClickAwayListener>
            )}
          </div>
        </Slide>
      </div>
      <ConnectionStateTransitionModal ref={deleteCtxtRef} />
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
  contexts,
  activeContexts,
  setActiveContexts,
  searchContexts,
}) => {
  const { notify } = useNotification;
  const { openModal } = useContext(WorkspaceModalContext) || {};
  const theme = useTheme();
  const isSmallScreen = useMediaQuery(theme.breakpoints.up('md'));

  const {
    data: providerCapabilities,
    isError: isProviderCapabilitiesError,
    error: providerCapabilitiesError,
  } = useGetProviderCapabilitiesQuery();

  if (isProviderCapabilitiesError) {
    notify({
      message: 'Error fetching provider capabilities',
      event_type: EVENT_TYPES.ERROR,
      details: providerCapabilitiesError?.data,
    });
  }

  const remoteProviderUrl = providerCapabilities?.providerUrl;
  const collaboratorExtensionUri = providerCapabilities?.extensions?.collaborator?.[0]?.component;

  const loaderType = 'circular';
  return (
    <NoSsr>
      <>
        <HeaderAppBar id="top-navigation-bar" color="primary" position="sticky">
          <StyledToolbar disableGutters isDrawerCollapsed={onDrawerCollapse}>
            <Grid2 container size="grow" sx={{ alignItems: 'center' }}>
              <Hidden smUp>
                <Grid2 style={{ display: 'none' }}>
                  <MenuIconButton aria-label="Open drawer" onClick={onDrawerToggle}>
                    <HeaderIcons style={iconMedium} />
                  </MenuIconButton>
                </Grid2>
              </Hidden>
              <Grid2
                container
                component={PageTitleWrapper}
                size="grow"
                sx={{ alignItems: 'center' }}
              >
                {/* Extension Point for   Logo */}
                <div
                  id="nav-header-logo"
                  style={{
                    height: '100%',
                    display: 'flex',
                    alignItems: 'center',
                    width: 'fit-content',
                    justifyContent: 'center',
                    position: 'relative',
                  }}
                ></div>
                <OrganizationAndWorkSpaceSwitcher />
              </Grid2>
              <Box
                component={UserContainer}
                style={{
                  position: 'relative',
                  display: 'flex',
                  gap: '1rem 0.5rem',
                  width: 'fit-content',
                }}
              >
                {/* According to the capabilities load the component */}
                <ErrorBoundary customFallback={() => null}>
                  {collaboratorExtensionUri && isSmallScreen && (
                    <RemoteComponent
                      url={{ url: createPathForRemoteComponent(collaboratorExtensionUri) }}
                      loaderType={loaderType}
                      providerUrl={remoteProviderUrl}
                      getUserAccessToken={getUserAccessToken}
                      getUserProfile={getUserProfile}
                      onOpenWorkspace={openModal}
                    />
                  )}
                </ErrorBoundary>
                <UserInfoContainer>
                  <UserSpan
                    sx={{
                      display: {
                        xs: 'none',
                        sm: 'inline-flex',
                      },
                    }}
                    style={{ position: 'relative' }}
                  >
                    <K8sContextMenu
                      contexts={contexts}
                      activeContexts={activeContexts}
                      setActiveContexts={setActiveContexts}
                      searchContexts={searchContexts}
                    />
                  </UserSpan>
                  <CustomTooltip title="Notifications">
                    <div data-testid="notification-button">
                      <NotificationDrawerButton />
                    </div>
                  </CustomTooltip>
                  <CustomTooltip title={'User Profile'}>
                    <UserSpan>
                      <User />
                    </UserSpan>
                  </CustomTooltip>
                  <UserSpan data-testid="header-menu">
                    <HeaderMenu />
                  </UserSpan>
                </UserInfoContainer>
              </Box>
            </Grid2>
          </StyledToolbar>
        </HeaderAppBar>
        <RegistryModal />
      </>
    </NoSsr>
  );
};

Header.propTypes = {
  onDrawerToggle: PropTypes.func.isRequired,
};

export default Header;
