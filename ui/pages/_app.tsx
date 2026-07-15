import {
  CheckCircleIcon as CheckCircle,
  ErrorIcon as Error,
  InfoIcon as Info,
  PermissionProvider,
  WarningIcon as Warning,
} from '@sistent/sistent';
import { Footer, KubernetesSubscription, NavigationBar } from '../components/AppComponents';
import { AdapterMoment, LocalizationProvider } from '@/components/shared/DatePicker';
import { CacheProvider } from '@emotion/react';
import createCache from '@emotion/cache';
import 'billboard.js/dist/theme/dark.min.css';
import _ from 'lodash';
import Head from 'next/head';
import { SnackbarProvider } from 'notistack';
import React, { useEffect, useMemo, useCallback, useRef, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { startSessionTimer } from '../lib/sessionTimer';
import Header from '../components/layout/Header/Header';
import MesheryProgressBar from '../components/MesheryProgressBar';
import getPageContext from '../components/PageContext';
import { subscribeToControllersStatus } from 'lib/controllersStatusSubscription';
import { useLazyGetSystemSyncQuery, useLazyGetKubernetesContextsQuery } from '../rtk-query/system';
import { useGetUserPrefQuery } from '../rtk-query/user';
import { api } from '../rtk-query';
import 'react-grid-layout/css/styles.css';
import 'react-resizable/css/styles.css';
// Host-side CSS for packages shared with extensions via remote-component.
// Next.js (pages router) only permits global CSS imports from _app; remote
// plugins cannot inject their own global stylesheets through the bundler
// pipeline. Import the full set of tippy.js themes/animations and xterm CSS
// here so any extension that references them via remote-component-stubbed
// subpaths (see remote-component.config.js) has the styles already on-page.
import 'tippy.js/dist/tippy.css';
import 'tippy.js/dist/svg-arrow.css';
import 'tippy.js/dist/border.css';
import 'tippy.js/dist/backdrop.css';
import 'tippy.js/themes/light.css';
import 'tippy.js/themes/light-border.css';
import 'tippy.js/themes/material.css';
import 'tippy.js/themes/translucent.css';
import 'tippy.js/animations/shift-away.css';
import 'tippy.js/animations/shift-away-subtle.css';
import 'tippy.js/animations/shift-away-extreme.css';
import 'tippy.js/animations/shift-toward.css';
import 'tippy.js/animations/shift-toward-subtle.css';
import 'tippy.js/animations/shift-toward-extreme.css';
import 'tippy.js/animations/scale.css';
import 'tippy.js/animations/scale-subtle.css';
import 'tippy.js/animations/scale-extreme.css';
import 'tippy.js/animations/perspective.css';
import 'tippy.js/animations/perspective-subtle.css';
import 'tippy.js/animations/perspective-extreme.css';
import '@xterm/xterm/css/xterm.css';
import {
  getControllerPollConnectionIDsFromContextIds,
  getK8sConfigIdsFromK8sConfig,
} from '../utils/multi-ctx';
import './../public/static/style/index.css';
import './styles/AnimatedFilter.css';
import './styles/AnimatedMeshery.css';
import './styles/AnimatedMeshPattern.css';
import './styles/AnimatedMeshSync.css';
import PlaygroundMeshDeploy from '../components/layout/AccessMesheryModal';
import Router from 'next/router';
import { RelayEnvironmentProvider } from 'react-relay';
import { createRelayEnvironment } from '../lib/relayEnvironment';
import './styles/charts.css';
import uiConfig from '../ui.config';
import { NotificationCenterProvider } from '../components/layout/NotificationCenter';
import { getConnectionDefinitions, getMeshModelComponentByName } from '../api/meshmodel';
import { CONNECTION_KINDS, CONNECTION_KINDS_DEF } from '../utils/Enum';
import { ability } from '../utils/can';
import { DynamicComponentProvider } from '@/utils/context/dynamicContext';
import { formatToTitleCase } from '@/utils/utils';
import { useThemePreference } from '@/themes/hooks';
import { CssBaseline, NoSsr, SistentThemeProvider } from '@/theme';
import { ErrorBoundary } from '@sistent/sistent';
import { LoadSessionGuard } from '@/rtk-query/ability';
import { useGetLoggedInUserQuery, useGetSelectedOrganization } from '@/rtk-query/user';
import CustomErrorFallback from '@/components/shared/ErrorBoundary/ErrorBoundary';
import { normalizeLoadTestPrefs } from '../lib/load-test-prefs';
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
import { updateAdaptersInfo } from '@/store/slices/adapter';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';
import WorkspaceModalContextProvider from '@/utils/context/WorkspaceModalContextProvider';
import RegistryModalContextProvider from '@/utils/context/RegistryModalContextProvider';
import { DynamicFullScreenLoader } from '@/components/shared/LoadingState/DynamicFullscreenLoader';

export const mesheryExtensionRoute = '/extension/meshmap';
function isMesheryUIRestrictedAndThePageIsNotPlayground(providerCapabilities) {
  return (
    !window.location.pathname.startsWith(mesheryExtensionRoute) &&
    providerCapabilities?.restrictedAccess?.isMesheryUIRestricted
  );
}

export function isExtensionOpen() {
  return window.location.pathname.startsWith(mesheryExtensionRoute);
}

const MesheryApp = ({ Component, pageProps, relayEnvironment, emotionCache }) => {
  const pageContext = useMemo(() => getPageContext(), []);
  const { k8sConfig } = useSelector((state) => state.ui);
  const { providerCapabilities } = useSelector((state) => state.ui);
  const { isDrawerCollapsed } = useSelector((state) => state.ui);
  const [fetchSystemSync] = useLazyGetSystemSyncQuery();
  const [fetchKubernetesContexts] = useLazyGetKubernetesContextsQuery();
  const [fetchOrganizations] = api.endpoints.getOrgs.useLazyQuery();
  const [fetchUserKeys] = api.endpoints.getUserKeys.useLazyQuery();
  const { data: userPrefData } = useGetUserPrefQuery();
  const dispatch = useDispatch();
  const [state, setState] = useState({
    mobileOpen: false,
    isDrawerCollapsed: false,
    isFullScreenMode: false,
    isLoading: true,
    k8sContexts: { totalCount: 0, contexts: [] },
    activeK8sContexts: [],
    mesheryControllerSubscription: null,
    theme: 'light',
    isOpen: false,
    relayEnvironment: createRelayEnvironment(),
    connectionMetadata: {},
    keys: [],
    abilities: [],
    abilityUpdated: false,
  });

  // ── PermissionProvider: CASL adapter ──────────────────────
  // Sistent permission checks are delegated to the existing CASL `ability` instance here.
  // If CASL is replaced later, only this adapter should need to change.
  const userHasPermission = useCallback(
    (key) => ability.can(key.id, _.lowerCase(key.function)),
    // `ability` is a module-level singleton; the reference never changes.
    // Re-creating this callback is intentionally avoided.
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [],
  );

  const { data: loggedInUser } = useGetLoggedInUserQuery({});
  const { selectedOrganization } = useGetSelectedOrganization();

  const permissionUserContext = useMemo(() => {
    const firstName = loggedInUser?.firstName || loggedInUser?.first_name || '';
    const lastName = loggedInUser?.lastName || loggedInUser?.last_name || '';
    const userName = `${firstName} ${lastName}`.trim() || loggedInUser?.name || loggedInUser?.email;
    return {
      userName,
      orgName: selectedOrganization?.name,
      roleNames: loggedInUser?.roleNames || [],
    };
  }, [loggedInUser, selectedOrganization]);

  // Holds the live controller-status SSE subscription ({ dispose }) so
  // initSubscriptions can tear down the previous stream and the bootstrap
  // cleanup can dispose it on unmount, without racing a stale state closure.
  const mesheryControllerSubscriptionRef = useRef<null | { dispose: () => void }>(null);

  const setAppState = useCallback((partialState, callback) => {
    setState((prevState) => {
      const newState = { ...prevState, ...partialState };
      if (callback) {
        setTimeout(callback, 0);
      }
      return newState;
    });
  }, []);

  const fullScreenChanged = useCallback(() => {
    setState((prevState) => {
      return { ...prevState, isFullScreenMode: !prevState.isFullScreenMode };
    });
  }, []);

  const loadMeshModelComponent = useCallback(async () => {
    const connectionDef = {};

    // Connection definitions are the source of truth for a kind's icon
    // (styles.svgColor) and its state machine (transitionMap). Seed the metadata
    // from them — these entries must survive even when the legacy
    // `<Kind>Connection` *component* no longer exists (connection definitions
    // replaced those components), otherwise the kind has no transitionMap and
    // the status dropdown shows "No transitions Available".
    try {
      const res = await getConnectionDefinitions();
      (res?.connectionDefinitions || []).forEach((definition) => {
        if (definition?.kind) {
          connectionDef[definition.kind] = {
            transitionMap: definition.transitionMap,
            icon: definition.styles?.svgColor,
          };
        }
      });
    } catch (error) {
      console.error('Error fetching connection definitions:', error);
    }

    // Fall back to the legacy `<Kind>Connection` component for kinds without a
    // first-class connection definition yet (e.g. meshery, github), and to
    // backfill the flat `transitions` list / icon the definition did not provide.
    const promises = CONNECTION_KINDS_DEF.map(async (kind) => {
      try {
        const res = await getMeshModelComponentByName(formatToTitleCase(kind).concat('Connection'));
        if (res?.components?.length) {
          const kindKey = CONNECTION_KINDS[kind];
          const existing = connectionDef[kindKey] || {};
          connectionDef[kindKey] = {
            ...existing,
            transitions: existing.transitions ?? res.components[0].metadata?.transitions,
            icon: existing.icon || res.components[0].styles?.svgColor,
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
      // Only watch controller status for connections that are BOTH in operator
      // mode AND connected: the operator/broker/meshsync controllers exist
      // in-cluster only in operator mode, and only a connected connection has
      // live controllers to poll. Embedded or not-connected connections have
      // nothing to poll. This re-scopes automatically — a mode/status change
      // invalidates the connections cache, refetching k8sConfig and re-running
      // initSubscriptions with the updated eligible set.
      const connectionIDs = getControllerPollConnectionIDsFromContextIds(contexts, k8sConfig);

      // Tear down any prior controller-status stream before opening a new one,
      // so re-subscribing on a context change never leaks an EventSource.
      mesheryControllerSubscriptionRef.current?.dispose?.();
      mesheryControllerSubscriptionRef.current = null;

      // No operator-mode connections → no controller-status stream to open.
      if (connectionIDs.length < 1) {
        setState((prevState) => ({ ...prevState, mesheryControllerSubscription: null }));
        return;
      }

      // SSE stream (replaces the subscribeMesheryControllersStatus GraphQL
      // subscription). The server sends the full controller-status array on
      // every change, so we just replace the redux state — no client merge.
      const mesheryControllerSubscription = subscribeToControllersStatus(connectionIDs, (data) => {
        dispatch(setControllerState({ controllerState: data }));
      });
      mesheryControllerSubscriptionRef.current = mesheryControllerSubscription;

      setState((prevState) => ({ ...prevState, mesheryControllerSubscription }));
    },
    [k8sConfig, dispatch],
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
        const active = ctx?.contexts?.find((c) => c.isCurrentContext === true);
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
          const loadTestPref = normalizeLoadTestPrefs(result.loadTestPrefs);
          dispatch(updateLoadTestPref({ loadTestPref }));
        }
      }
    } catch (error) {
      console.log(`there was an error fetching user config data: ${error}`);
    }
  }, [dispatch, fetchSystemSync]);

  useEffect(() => {
    startSessionTimer();

    const loadAll = async () => {
      try {
        loadConfigFromServer();
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
      mesheryControllerSubscriptionRef.current?.dispose?.();
    };
  }, []);

  // Update effect for k8sConfig
  useEffect(() => {
    // in case the meshery-ui is restricted, the user will be redirected to signup/extension page
    if (
      typeof window !== 'undefined' &&
      isMesheryUIRestrictedAndThePageIsNotPlayground(providerCapabilities)
    ) {
      Router.push(mesheryExtensionRoute);
    }

    if (k8sConfig?.length > 0) {
      // initSubscriptions disposes any existing stream and re-subscribes with
      // the current connection set, so it is safe to call on every k8sConfig
      // change.
      const ids = getK8sConfigIdsFromK8sConfig(k8sConfig);
      initSubscriptions(ids);
    }
  }, [k8sConfig, providerCapabilities]);

  const canShowNav = !state.isFullScreenMode && uiConfig?.components?.navigator !== false;
  const { extensionType } = useSelector((state) => state.ui);

  return (
    <DynamicFullScreenLoader isLoading={state.isLoading}>
      <PermissionProvider userHasPermission={userHasPermission} userContext={permissionUserContext}>
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
                                    providerCapabilities={providerCapabilities}
                                  />
                                </StyledContentWrapper>
                              </NotificationCenterProvider>
                            </SnackbarProvider>
                          </StyledAppContent>
                        </StyledRoot>
                        <PlaygroundMeshDeploy
                          closeForm={() =>
                            setState((prevState) => ({ ...prevState, isOpen: false }))
                          }
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
      </PermissionProvider>
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
