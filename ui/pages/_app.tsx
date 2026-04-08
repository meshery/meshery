import { CheckCircle, Error, Info, Warning } from '@mui/icons-material';
import { Footer, KubernetesSubscription, NavigationBar } from '../components/AppComponents';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterMoment } from '@mui/x-date-pickers/AdapterMoment';
import { CacheProvider } from '@emotion/react';
import createCache from '@emotion/cache';
import 'billboard.js/dist/theme/dark.min.css';
import _ from 'lodash';
import Head from 'next/head';
import { SnackbarProvider } from 'notistack';
import React, { useEffect, useMemo, useCallback, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import Header from '../components/Header';
import MesheryProgressBar from '../components/MesheryProgressBar';
import getPageContext from '../components/PageContext';
import { MESHERY_CONTROLLER_SUBSCRIPTION } from '../components/subscription/helpers';
import { GQLSubscription } from '../components/subscription/subscriptionhandler';
import { useLazyGetSystemSyncQuery, useLazyGetKubernetesContextsQuery } from '../rtk-query/system';
import { useGetUserPrefQuery } from '../rtk-query/user';
import { api } from '../rtk-query';
import { useLazyGetConnectionsQuery } from '../rtk-query/connection';
import 'react-grid-layout/css/styles.css';
import 'react-resizable/css/styles.css';
import { getConnectionIDsFromContextIds, getK8sConfigIdsFromK8sConfig } from '../utils/multi-ctx';
import './../public/static/style/index.css';
import './styles/AnimatedFilter.css';
import './styles/AnimatedMeshery.css';
import './styles/AnimatedMeshPattern.css';
import './styles/AnimatedMeshSync.css';
import PlaygroundMeshDeploy from './extension/AccessMesheryModal';
import Router from 'next/router';
import { RelayEnvironmentProvider } from 'react-relay';
import { createRelayEnvironment } from '../lib/relayEnvironment';
import './styles/charts.css';
import uiConfig from '../ui.config';
import { NotificationCenterProvider } from '../components/NotificationCenter';
import { getMeshModelComponentByName } from '../api/meshmodel';
import { CONNECTION_KINDS, CONNECTION_KINDS_DEF, CONNECTION_STATES } from '../utils/Enum';
import { ability } from '../utils/can';
import { useLazyGetCredentialByIdQuery } from '@/rtk-query/credentials';
import { DynamicComponentProvider } from '@/utils/context/dynamicContext';
import { formatToTitleCase } from '@/utils/utils';
import { useThemePreference } from '@/themes/hooks';
import { ErrorBoundary, SistentThemeProvider, CssBaseline, NoSsr } from '@sistent/sistent';
import { LoadSessionGuard } from '@/rtk-query/ability';
import CustomErrorFallback from '@/components/General/ErrorBoundary';
import {
  StyledAppContent,
  StyledMainContent,
  StyledContentWrapper,
  StyledRoot,
  ThemeResponsiveSnackbar,
} from '../themes/App.styles';
import {
  setConnectionMetadata,
  setControllerState,
  setK8sContexts,
  setKeys,
  setOrganization,
  toggleCatalogContent,
  updateExtensionType,
} from '@/store/slices/mesheryUi';
import { updateLoadTestPref } from '@/store/slices/prefTest';
import { updateGrafanaConfig, updatePrometheusConfig } from '@/store/slices/telemetry';
import { updateAdaptersInfo } from '@/store/slices/adapter';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';
import WorkspaceModalContextProvider from '@/utils/context/WorkspaceModalContextProvider';
import RegistryModalContextProvider from '@/utils/context/RegistryModalContextProvider';
import { DynamicFullScreenLoader } from '@/components/LoadingComponents/DynamicFullscreenLoader';

export const mesheryExtensionRoute = '/extension/meshmap';
function isMesheryUiRestrictedAndThePageIsNotPlayground(capabilitiesRegistry) {
  return (
    !window.location.pathname.startsWith(mesheryExtensionRoute) &&
    capabilitiesRegistry?.restrictedAccess?.isMesheryUiRestricted
  );
}

export function isExtensionOpen() {
  return window.location.pathname.startsWith(mesheryExtensionRoute);
}

const MesheryApp = ({ Component, pageProps, relayEnvironment, emotionCache }) => {
  const pageContext = useMemo(() => getPageContext(), []);
  const { k8sConfig } = useSelector((state) => state.ui);
  const { capabilitiesRegistry } = useSelector((state) => state.ui);
  const { isDrawerCollapsed } = useSelector((state) => state.ui);
  const [fetchCredentialById] = useLazyGetCredentialByIdQuery();
  const [fetchSystemSync] = useLazyGetSystemSyncQuery();
  const [fetchKubernetesContexts] = useLazyGetKubernetesContextsQuery();
  const [fetchOrganizations] = api.endpoints.getOrgs.useLazyQuery();
  const [fetchUserKeys] = api.endpoints.getUserKeys.useLazyQuery();
  const [fetchConnections] = useLazyGetConnectionsQuery();
  const { data: userPrefData } = useGetUserPrefQuery();
  const dispatch = useDispatch();
  const [state, setState] = useState({
    mobileOpen: false,
    isDrawerCollapsed: false,
    isFullScreenMode: false,
    isLoading: true,
    k8sContexts: [],
    activeK8sContexts: [],
    mesheryControllerSubscription: null,
    disposeK8sContextSubscription: null,
    theme: 'light',
    isOpen: false,
    relayEnvironment: createRelayEnvironment(),
    connectionMetadata: {},
    keys: [],
    abilities: [],
    abilityUpdated: false,
  });

  const setAppState = useCallback((partialState, callback) => {
    setState((prevState) => {
      const newState = { ...prevState, ...partialState };
      if (callback) {
        setTimeout(callback, 0);
      }
      return newState;
    });
  }, []);

  const loadPromGrafanaConnection = useCallback(async () => {
    try {
      const res = await fetchConnections({
        page: 0,
        pagesize: 2,
        status: JSON.stringify([CONNECTION_STATES.CONNECTED, CONNECTION_STATES.REGISTERED]),
        kind: JSON.stringify([CONNECTION_KINDS.PROMETHEUS, CONNECTION_KINDS.GRAFANA]),
      }).unwrap();

      res?.connections?.forEach((connection) => {
        if (connection.kind == CONNECTION_KINDS.PROMETHEUS) {
          const promCfg = {
            prometheusURL: connection?.metadata?.url || '',
            selectedPrometheusBoardsConfigs: connection?.metadata['prometheus_boards'] || [],
            connectionID: connection?.id,
            connectionName: connection?.name,
          };
          dispatch(updatePrometheusConfig(promCfg));
        } else {
          const credentialID = connection?.credential_id;
          fetchCredentialById(credentialID)
            .unwrap()
            .then((credRes) => {
              const grafanaCfg = {
                grafanaURL: connection?.metadata?.url || '',
                grafanaAPIKey: credRes?.secret?.secret || '',
                grafanaBoardSearch: '',
                grafanaBoards: connection?.metadata['grafana_boards'] || [],
                selectedBoardsConfigs: [],
                connectionID: connection?.id,
                connectionName: connection?.name,
              };
              dispatch(updateGrafanaConfig(grafanaCfg));
            });
        }
      });
    } catch (err) {
      console.error('Failed to load telemetry connections:', err);
    }
  }, [dispatch, fetchConnections, fetchCredentialById]);

  const fullScreenChanged = useCallback(() => {
    setState((prevState) => {
      return { ...prevState, isFullScreenMode: !prevState.isFullScreenMode };
    });
  }, []);

  const loadMeshModelComponent = useCallback(async () => {
    const connectionDef = {};

    const promises = CONNECTION_KINDS_DEF.map(async (kind) => {
      try {
        const res = await getMeshModelComponentByName(formatToTitleCase(kind).concat('Connection'));
        if (res?.components) {
          connectionDef[CONNECTION_KINDS[kind]] = {
            transitions: res?.components[0].metadata.transitions,
            icon: res?.components[0].styles.svgColor,
          };
        }
      } catch (error) {
        console.error(`Error fetching component for ${kind}:`, error);
      }
    });

    await Promise.all(promises);

    setState((prevState) => ({ ...prevState, connectionMetadata: connectionDef }));

    dispatch(
      setConnectionMetadata({
        connectionMetadataState: connectionDef,
      }),
    );
  }, [dispatch]);

  const initSubscriptions = useCallback(
    (contexts) => {
      if (!k8sConfig?.length) {
        return;
      }
      const connectionIDs = getConnectionIDsFromContextIds(contexts, k8sConfig);
      // No need to create a controller subscription if there are no connections
      if (connectionIDs.length < 1) {
        setState((prevState) => ({ ...prevState, mesheryControllerSubscription: () => {} }));
        return;
      }

      const mesheryControllerSubscription = new GQLSubscription({
        type: MESHERY_CONTROLLER_SUBSCRIPTION,
        connectionIDs: connectionIDs,
        callbackFunction: (data) => {
          dispatch(setControllerState({ controllerState: data }));
        },
      });
      mesheryControllerSubscription.initSubscription();

      setState((prevState) => ({ ...prevState, mesheryControllerSubscription }));
    },
    [k8sConfig],
  );

  const handleDrawerToggle = useCallback(() => {
    setState((prevState) => ({ ...prevState, mobileOpen: !prevState.mobileOpen }));
  }, []);

  const handleMesheryCommunityClick = useCallback(() => {
    setState((prevState) => ({ ...prevState, isOpen: !prevState.isOpen }));
  }, []);

  /**
   * Sets the selected k8s context on global level.
   * @param {Array.<string>} activeK8sContexts
   */
  const activeContextChangeCallback = useCallback((activeK8sContexts) => {
    if (activeK8sContexts.includes('all')) {
      activeK8sContexts = ['all'];
    }
    dispatch(setK8sContexts({ selectedK8sContexts: activeK8sContexts }));
  }, []);

  const setActiveContexts = useCallback(
    (id) => {
      if (state.k8sContexts?.contexts) {
        if (id === 'all') {
          let activeContexts = [];
          state.k8sContexts.contexts.forEach((ctx) => activeContexts.push(ctx.id));
          activeContexts.push('all');
          setState((prevState) => ({ ...prevState, activeK8sContexts: activeContexts }));
          activeContextChangeCallback(activeContexts);
          return;
        }

        // if id is an empty array, clear all active contexts
        if (Array.isArray(id) && id.length === 0) {
          setState((prevState) => ({ ...prevState, activeK8sContexts: [] }));
          activeContextChangeCallback([]);

          return;
        }

        setState((prevState) => {
          let ids = [...(prevState.activeK8sContexts || [])];
          //pop event
          if (ids.includes(id)) {
            ids = ids.filter((cid) => cid !== 'all');
            const filteredIds = ids.filter((cid) => cid !== id);
            activeContextChangeCallback(filteredIds);
            return { ...prevState, activeK8sContexts: filteredIds };
          }

          //push event
          if (ids.length === prevState.k8sContexts.contexts.length - 1) {
            ids.push('all');
          }
          const newIds = [...ids, id];
          activeContextChangeCallback(newIds);
          return { ...prevState, activeK8sContexts: newIds };
        });
      }
    },
    [state.k8sContexts, state.activeK8sContexts, activeContextChangeCallback],
  );

  const searchContexts = useCallback(
    async (search = '') => {
      try {
        const ctx = await fetchKubernetesContexts({ pagesize: 10, search }).unwrap();
        setState((prevState) => ({ ...prevState, k8sContexts: ctx }));
        const active = ctx?.contexts?.find((c) => c.is_current_context === true);
        if (active) {
          setState((prevState) => ({ ...prevState, activeK8sContexts: [active?.id] }));
          activeContextChangeCallback([active?.id]);
        }
      } catch (err) {
        console.error(err);
      }
    },
    [activeContextChangeCallback, fetchKubernetesContexts],
  );

  const updateCurrentExtensionType = useCallback(
    (type) => {
      dispatch(updateExtensionType({ extensionType: type }));
    },
    [dispatch],
  );
  const setCurrentOrganization = useCallback(
    (org) => {
      dispatch(setOrganization({ organization: org }));
    },
    [dispatch],
  );

  const updateAbility = useCallback(() => {
    ability.update(
      state.keys?.map((key) => ({ action: key.id, subject: _.lowerCase(key.function) })),
    );
    setState((prevState) => ({ ...prevState, abilityUpdated: true }));
  }, [state.keys]);

  const loadAbility = useCallback(
    async (orgID, reFetchKeys) => {
      const storedKeys = sessionStorage.getItem('keys');
      if (storedKeys !== null && !reFetchKeys && storedKeys !== 'undefined') {
        setState((prevState) => ({ ...prevState, keys: JSON.parse(storedKeys) }));
        updateAbility();
      } else {
        try {
          const result = await fetchUserKeys({ orgId: orgID }).unwrap();
          if (result) {
            setState((prevState) => ({ ...prevState, keys: result.keys }));
            dispatch(setKeys({ keys: result.keys }));
            updateAbility();
          }
        } catch (err) {
          console.log('There was an error fetching user keys:', err);
        }
      }
    },
    [dispatch, updateAbility, fetchUserKeys],
  );

  const loadOrg = useCallback(async () => {
    const currentOrg = sessionStorage.getItem('currentOrg');
    let reFetchKeys = false;

    if (currentOrg && currentOrg !== 'undefined') {
      let org = JSON.parse(currentOrg);
      await loadAbility(org.id, reFetchKeys);
      setCurrentOrganization(org);
    }

    try {
      const result = await fetchOrganizations({}).unwrap();
      let organizationToSet;
      const sessionOrg = currentOrg ? JSON.parse(currentOrg) : null;

      if (currentOrg) {
        const indx = result.organizations.findIndex((org) => org.id === sessionOrg.id);
        if (indx === -1) {
          organizationToSet = result.organizations[0];
          reFetchKeys = true;
          await loadAbility(organizationToSet.id, reFetchKeys);
          setCurrentOrganization(organizationToSet);
        }
      } else {
        organizationToSet = result.organizations[0];
        reFetchKeys = true;
        await loadAbility(organizationToSet.id, reFetchKeys);
        setCurrentOrganization(organizationToSet);
      }
    } catch (err) {
      console.log('There was an error fetching available orgs:', err);
    }
  }, [loadAbility, setCurrentOrganization, fetchOrganizations]);

  const loadConfigFromServer = useCallback(async () => {
    try {
      const result = await fetchSystemSync().unwrap();
      if (result) {
        if (result.meshAdapters && result.meshAdapters !== null && result.meshAdapters.length > 0) {
          dispatch(updateAdaptersInfo({ meshAdapters: result.meshAdapters }));
        }
        if (result.loadTestPrefs) {
          const loadTestPref = Object.assign(
            {
              c: 0,
              qps: 0,
              t: 0,
              gen: 0,
            },
            result.loadTestPrefs,
          );
          dispatch(updateLoadTestPref({ loadTestPref }));
        }
      }
    } catch (error) {
      console.log(`there was an error fetching user config data: ${error}`);
    }
  }, [dispatch, fetchSystemSync]);

  useEffect(() => {
    const { startSessionTimer } = require('../lib/sessionTimer');
    startSessionTimer();

    const loadAll = async () => {
      try {
        loadConfigFromServer();
        loadPromGrafanaConnection();
        await loadOrg();

        initSubscriptions([]);

        // Catalog content preference is loaded via useGetUserPrefQuery (reactive)
        if (typeof userPrefData?.usersExtensionPreferences?.catalogContent !== 'undefined') {
          dispatch(
            toggleCatalogContent({
              catalogVisibility: userPrefData?.usersExtensionPreferences?.catalogContent,
            }),
          );
        }

        document.addEventListener('fullscreenchange', fullScreenChanged);
        await loadMeshModelComponent();
      } catch (error) {
        console.error('[Meshery bootstrap] Failed to initialize the application shell', error);
      } finally {
        setState((prevState) => ({ ...prevState, isLoading: false }));
      }
    };
    loadAll();

    return () => {
      document.removeEventListener('fullscreenchange', fullScreenChanged);
      if (state.disposeK8sContextSubscription) {
        state.disposeK8sContextSubscription();
      }
    };
  }, []);

  // Update effect for k8sConfig
  useEffect(() => {
    // in case the meshery-ui is restricted, the user will be redirected to signup/extension page
    if (
      typeof window !== 'undefined' &&
      isMesheryUiRestrictedAndThePageIsNotPlayground(capabilitiesRegistry)
    ) {
      Router.push(mesheryExtensionRoute);
    }

    if (k8sConfig?.length > 0) {
      const { mesheryControllerSubscription } = state;
      const ids = getK8sConfigIdsFromK8sConfig(k8sConfig);
      if (mesheryControllerSubscription) {
        mesheryControllerSubscription.updateSubscription(
          getConnectionIDsFromContextIds(ids, k8sConfig),
        );
      } else {
        initSubscriptions(ids);
      }
    }
  }, [k8sConfig, capabilitiesRegistry]);

  const canShowNav = !state.isFullScreenMode && uiConfig?.components?.navigator !== false;
  const { extensionType } = useSelector((state) => state.ui);

  return (
    <DynamicFullScreenLoader isLoading={state.isLoading}>
      <DynamicComponentProvider>
        <RelayEnvironmentProvider environment={relayEnvironment}>
          <MesheryThemeProvider emotionCache={emotionCache}>
            <NoSsr>
              <ErrorBoundary customFallback={CustomErrorFallback}>
                <LoadSessionGuard>
                  <WorkspaceModalContextProvider>
                    <RegistryModalContextProvider>
                      <StyledRoot>
                        <CssBaseline />
                        <NavigationBar
                          isDrawerCollapsed={isDrawerCollapsed}
                          mobileOpen={state.mobileOpen}
                          handleDrawerToggle={handleDrawerToggle}
                          updateExtensionType={updateCurrentExtensionType}
                          canShowNav={canShowNav}
                        />
                        <StyledAppContent
                          canShowNav={canShowNav}
                          isDrawerCollapsed={isDrawerCollapsed}
                        >
                          <SnackbarProvider
                            anchorOrigin={{
                              vertical: 'bottom',
                              horizontal: 'right',
                            }}
                            iconVariant={{
                              success: <CheckCircle style={{ marginRight: '0.5rem' }} />,
                              error: <Error style={{ marginRight: '0.5rem' }} />,
                              warning: <Warning style={{ marginRight: '0.5rem' }} />,
                              info: <Info style={{ marginRight: '0.5rem' }} />,
                            }}
                            Components={{
                              info: ThemeResponsiveSnackbar,
                              success: ThemeResponsiveSnackbar,
                              error: ThemeResponsiveSnackbar,
                              warning: ThemeResponsiveSnackbar,
                              loading: ThemeResponsiveSnackbar,
                            }}
                            maxSnack={10}
                          >
                            <NotificationCenterProvider>
                              <MesheryProgressBar />
                              <KubernetesSubscription setAppState={setAppState} />
                              {!state.isFullScreenMode && (
                                <Header
                                  onDrawerToggle={handleDrawerToggle}
                                  onDrawerCollapse={isDrawerCollapsed}
                                  contexts={state.k8sContexts}
                                  activeContexts={state.activeK8sContexts}
                                  setActiveContexts={setActiveContexts}
                                  searchContexts={searchContexts}
                                  updateExtensionType={updateCurrentExtensionType}
                                  abilityUpdated={state.abilityUpdated}
                                />
                              )}
                              <StyledContentWrapper>
                                <StyledMainContent
                                  id="meshery-main"
                                  style={{
                                    padding: extensionType === 'navigator' && '0px',
                                  }}
                                >
                                  <LocalizationProvider dateAdapter={AdapterMoment}>
                                    <ErrorBoundary customFallback={CustomErrorFallback}>
                                      <Component
                                        pageContext={pageContext}
                                        contexts={state.k8sContexts}
                                        activeContexts={state.activeK8sContexts}
                                        setActiveContexts={setActiveContexts}
                                        searchContexts={searchContexts}
                                        {...pageProps}
                                      />
                                    </ErrorBoundary>
                                  </LocalizationProvider>
                                </StyledMainContent>
                                <Footer
                                  handleMesheryCommunityClick={handleMesheryCommunityClick}
                                  capabilitiesRegistry={capabilitiesRegistry}
                                />
                              </StyledContentWrapper>
                            </NotificationCenterProvider>
                          </SnackbarProvider>
                        </StyledAppContent>
                      </StyledRoot>
                      <PlaygroundMeshDeploy
                        closeForm={() => setState((prevState) => ({ ...prevState, isOpen: false }))}
                        isOpen={state.isOpen}
                      />
                    </RegistryModalContextProvider>
                  </WorkspaceModalContextProvider>
                </LoadSessionGuard>
              </ErrorBoundary>
            </NoSsr>
          </MesheryThemeProvider>
        </RelayEnvironmentProvider>
      </DynamicComponentProvider>
    </DynamicFullScreenLoader>
  );
};

// Keep the static getInitialProps method
MesheryApp.getInitialProps = async ({ Component, ctx }) => {
  if (!Component) {
    return { pageProps: {} };
  }
  const pageProps = Component.getInitialProps ? await Component.getInitialProps(ctx) : {};
  return { pageProps };
};

// Client-side Emotion cache with prepend: true ensures CSS-in-JS styles
// are inserted before other stylesheets for correct specificity
const clientSideEmotionCache = createCache({ key: 'css', prepend: true });

const MesheryThemeProvider = ({ children, emotionCache }) => {
  const themePref = useThemePreference();
  const mode = themePref?.data?.mode || 'dark';
  return (
    <SistentThemeProvider initialMode={mode} emotionCache={emotionCache}>
      {children}
    </SistentThemeProvider>
  );
};

const MesheryAppWrapper = ({ emotionCache = clientSideEmotionCache, ...props }) => {
  return (
    <CacheProvider value={emotionCache}>
      <ProviderStoreWrapper>
        <Head>
          <link rel="shortcut icon" href="/static/img/meshery-logo/meshery-logo.svg" />
          <title>Meshery</title>
        </Head>
        <LocalizationProvider dateAdapter={AdapterMoment}>
          <MesheryApp {...props} emotionCache={emotionCache} />
        </LocalizationProvider>
      </ProviderStoreWrapper>
    </CacheProvider>
  );
};

export default MesheryAppWrapper;
