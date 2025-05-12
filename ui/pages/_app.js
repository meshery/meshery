import { CheckCircle, Error, Info, Warning } from '@mui/icons-material';
import FavoriteIcon from '@mui/icons-material/Favorite';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterMoment } from '@mui/x-date-pickers/AdapterMoment';
import 'billboard.js/dist/theme/dark.min.css';
import 'codemirror/addon/lint/lint.css';
// codemirror + js-yaml imports when added to a page was preventing to navigating to that page using nextjs
// link clicks, hence attempting to add them here
import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/material.css';
import { fromJS } from 'immutable';
import _ from 'lodash';
import withRedux from 'next-redux-wrapper';
import Head from 'next/head';
import { SnackbarProvider } from 'notistack';
import React, { useEffect, useMemo, useCallback, useState } from 'react';
import { connect, Provider, useSelector } from 'react-redux';
import Header from '../components/Header';
import MesheryProgressBar from '../components/MesheryProgressBar';
import Navigator from '../components/Navigator';
import getPageContext from '../components/PageContext';
import { MESHERY_CONTROLLER_SUBSCRIPTION } from '../components/subscription/helpers';
import { GQLSubscription } from '../components/subscription/subscriptionhandler';
import dataFetch, { promisifiedDataFetch } from '../lib/data-fetch';
import 'react-grid-layout/css/styles.css';
import 'react-resizable/css/styles.css';
import {
  actionTypes,
  makeStore,
  toggleCatalogContent,
  setConnectionMetadata,
  LegacyStoreContext,
} from '../lib/store';
import { getConnectionIDsFromContextIds, getK8sConfigIdsFromK8sConfig } from '../utils/multi-ctx';
import './../public/static/style/index.css';
import subscribeK8sContext from '../components/graphql/subscriptions/K8sContextSubscription';
import { bindActionCreators } from 'redux';
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
import CAN, { ability } from '../utils/can';
import { getCredentialByID } from '@/api/credentials';
import { DynamicComponentProvider } from '@/utils/context/dynamicContext';
import { store } from '../store';
import { RTKContext, useDispatchRtk, useSelectorRtk } from '@/store/hooks';
import { formatToTitleCase } from '@/utils/utils';
import { useThemePreference } from '@/themes/hooks';
import {
  ErrorBoundary,
  useTheme,
  SistentThemeProvider,
  CssBaseline,
  Typography,
  Hidden,
  NoSsr,
} from '@layer5/sistent';
import LoadingScreen from '@/components/LoadingComponents/LoadingComponentServer';
import { LoadSessionGuard } from '@/rtk-query/ability';
import { randomLoadingMessage } from '@/components/LoadingComponents/loadingMessages';
import { keys } from '@/utils/permission_constants';
import CustomErrorFallback from '@/components/General/ErrorBoundary';
import {
  StyledAppContent,
  StyledDrawer,
  StyledFooterBody,
  StyledFooterText,
  StyledMainContent,
  StyledContentWrapper,
  StyledRoot,
  ThemeResponsiveSnackbar,
} from '../themes/App.styles';
import { setK8sContexts, updateK8SConfig } from '@/store/slices/mesheryUi';

if (typeof window !== 'undefined') {
  require('codemirror/mode/yaml/yaml');
  require('codemirror/mode/javascript/javascript');
  require('codemirror/addon/lint/lint');
  require('codemirror/addon/lint/yaml-lint');
  require('codemirror/addon/lint/json-lint');
  if (typeof window.jsyaml === 'undefined') {
    window.jsyaml = require('js-yaml');
  }
  if (typeof window.jsonlint === 'undefined') {
    // jsonlint did not work well with codemirror json-lint. Hence, found an alternative (jsonlint-mod) based on https://github.com/scniro/react-codemirror2/issues/21
    window.jsonlint = require('jsonlint-mod');
  }
}

async function fetchContexts(number = 10, search = '') {
  return await promisifiedDataFetch(
    `/api/system/kubernetes/contexts?pagesize=${number}&search=${encodeURIComponent(search)}`,
  );
}

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

const Footer = ({ capabilitiesRegistry, handleL5CommunityClick }) => {
  const theme = useTheme();

  const extension = useSelector((state) => {
    return state.get('extensionType');
  });

  if (extension == 'navigator') {
    return null;
  }

  return (
    <StyledFooterBody>
      <Typography
        variant="body2"
        align="center"
        component="p"
        style={{
          color:
            theme.palette.mode === 'light'
              ? theme.palette.text.default
              : theme.palette.text.disabled,
        }}
      >
        <StyledFooterText onClick={handleL5CommunityClick}>
          {capabilitiesRegistry?.restrictedAccess?.isMesheryUiRestricted ? (
            'ACCESS LIMITED IN MESHERY PLAYGROUND. DEPLOY MESHERY TO ACCESS ALL FEATURES.'
          ) : (
            <>
              {' '}
              Built with{' '}
              <FavoriteIcon
                style={{
                  color: theme.palette.background.brand.default,
                  display: 'inline',
                  verticalAlign: 'bottom',
                }}
              />{' '}
              by the Layer5 Community
            </>
          )}
        </StyledFooterText>
      </Typography>
    </StyledFooterBody>
  );
};

const KubernetesSubscription = ({ setAppState }) => {
  const dispatch = useDispatchRtk();
  const k8sContextSubscription = (page = '', search = '', pageSize = '10', order = '') => {
    // Don't fetch k8s contexts if user doesn't have permission
    if (!CAN(keys.VIEW_ALL_KUBERNETES_CLUSTERS.action, keys.VIEW_ALL_KUBERNETES_CLUSTERS.subject)) {
      return () => {};
    }

    return subscribeK8sContext(
      (result) => {
        // Initialize activeContexts with all context IDs plus "all"
        const allContexts = [];
        if (result.k8sContext?.contexts?.length > 0) {
          result.k8sContext.contexts.forEach((ctx) => allContexts.push(ctx.id));
          allContexts.push('all');
        }

        // TODO: Remove local state and only use redux store
        setAppState({
          k8sContexts: result.k8sContext,
          activeK8sContexts: allContexts,
        });
        dispatch(updateK8SConfig({ k8sConfig: result.k8sContext.contexts }));
      },
      {
        selector: {
          page: page,
          pageSize: pageSize,
          order: order,
          search: search,
        },
      },
    );
  };

  useEffect(() => {
    const disposeK8sContextSubscription = k8sContextSubscription();
    setAppState({ disposeK8sContextSubscription });
    return () => {
      disposeK8sContextSubscription();
    };
  }, []);

  return null;
};

const MesheryApp = ({
  Component,
  pageProps,
  isDrawerCollapsed,
  relayEnvironment,
  store,
  toggleCatalogContent,
  setConnectionMetadata,
  extensionType,
}) => {
  const pageContext = useMemo(() => getPageContext(), []);
  const { k8sConfig } = useSelectorRtk((state) => state.ui);
  const { capabilitiesRegistry } = useSelectorRtk((state) => state.ui);
  const [state, setState] = useState({
    mobileOpen: false,
    isDrawerCollapsed: false,
    isFullScreenMode: false,
    isLoading: true,
    k8sContexts: [],
    activeK8sContexts: [],
    operatorSubscription: null,
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

  const loadPromGrafanaConnection = useCallback(() => {
    dataFetch(
      `/api/integrations/connections?page=0&pagesize=2&status=${encodeURIComponent(
        JSON.stringify([CONNECTION_STATES.CONNECTED, CONNECTION_STATES.REGISTERED]),
      )}&kind=${encodeURIComponent(
        JSON.stringify([CONNECTION_KINDS.PROMETHEUS, CONNECTION_KINDS.GRAFANA]),
      )}`,
      {
        credentials: 'include',
        method: 'GET',
      },
      (res) => {
        res?.connections?.forEach((connection) => {
          if (connection.kind == CONNECTION_KINDS.PROMETHEUS) {
            const promCfg = {
              prometheusURL: connection?.metadata?.url || '',
              selectedPrometheusBoardsConfigs: connection?.metadata['prometheus_boards'] || [],
              connectionID: connection?.id,
              connectionName: connection?.name,
            };

            store.dispatch({ type: actionTypes.UPDATE_PROMETHEUS_CONFIG, prometheus: promCfg });
          } else {
            const credentialID = connection?.credential_id;

            getCredentialByID(credentialID).then((res) => {
              const grafanaCfg = {
                grafanaURL: connection?.metadata?.url || '',
                grafanaAPIKey: res?.secret?.secret || '',
                grafanaBoardSearch: '',
                grafanaBoards: connection?.metadata['grafana_boards'] || [],
                selectedBoardsConfigs: [],
                connectionID: connection?.id,
                connectionName: connection?.name,
              };
              store.dispatch({ type: actionTypes.UPDATE_GRAFANA_CONFIG, grafana: grafanaCfg });
            });
          }
        });
      },
    );
  }, [store]);

  const fullScreenChanged = useCallback(() => {
    setState((prevState) => {
      return { ...prevState, isFullScreenMode: !prevState.isFullScreenMode };
    });
  }, []);

  const loadMeshModelComponent = useCallback(() => {
    const connectionDef = {};
    CONNECTION_KINDS_DEF.map(async (kind) => {
      const res = await getMeshModelComponentByName(formatToTitleCase(kind).concat('Connection'));
      if (res?.components) {
        connectionDef[CONNECTION_KINDS[kind]] = {
          transitions: res?.components[0].metadata.transitions,
          icon: res?.components[0].styles.svgColor,
        };
      }
      setState((prevState) => ({ ...prevState, connectionMetadata: connectionDef }));
    });
    setConnectionMetadata({
      connectionMetadataState: connectionDef,
    });
  }, [setConnectionMetadata]);

  const initSubscriptions = useCallback(
    (contexts) => {
      console.log('amit this called in initSubscriptions', contexts);
      const connectionIDs = getConnectionIDsFromContextIds(contexts, k8sConfig);
      console.log('amit connectionIDs', connectionIDs);
      // No need to create a controller subscription if there are no connections
      if (connectionIDs.length < 1) {
        setState((prevState) => ({ ...prevState, mesheryControllerSubscription: () => {} }));
        return;
      }

      const mesheryControllerSubscription = new GQLSubscription({
        type: MESHERY_CONTROLLER_SUBSCRIPTION,
        connectionIDs: connectionIDs,
        callbackFunction: (data) => {
          store.dispatch({
            type: actionTypes.SET_CONTROLLER_STATE,
            controllerState: data,
          });
        },
      });

      setState((prevState) => ({ ...prevState, mesheryControllerSubscription }));
    },
    [k8sConfig, store],
  );

  const handleDrawerToggle = useCallback(() => {
    setState((prevState) => ({ ...prevState, mobileOpen: !prevState.mobileOpen }));
  }, []);

  const handleL5CommunityClick = useCallback(() => {
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
    setK8sContexts({ selectedK8sContexts: activeK8sContexts });
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
    (search = '') => {
      fetchContexts(10, search)
        .then((ctx) => {
          setState((prevState) => ({ ...prevState, k8sContexts: ctx }));
          const active = ctx?.contexts?.find((c) => c.is_current_context === true);
          if (active) {
            setState((prevState) => ({ ...prevState, activeK8sContexts: [active?.id] }));
            activeContextChangeCallback([active?.id]);
          }
        })
        .catch((err) => console.error(err));
    },
    [activeContextChangeCallback],
  );

  const updateExtensionType = useCallback(
    (type) => {
      store.dispatch({ type: actionTypes.UPDATE_EXTENSION_TYPE, extensionType: type });
    },
    [store],
  );
  const setOrganization = useCallback(
    (org) => {
      store.dispatch({
        type: actionTypes.SET_ORGANIZATION,
        organization: org,
      });
    },
    [store],
  );

  const setWorkspace = useCallback(
    (workspace) => {
      store.dispatch({
        type: actionTypes.SET_WORKSPACE,
        workspace: workspace,
      });
    },
    [store],
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
        dataFetch(
          `/api/identity/orgs/${orgID}/users/keys`,
          {
            method: 'GET',
            credentials: 'include',
          },
          (result) => {
            if (result) {
              setState((prevState) => ({ ...prevState, keys: result.keys }));
              store.dispatch({
                type: actionTypes.SET_KEYS,
                keys: result.keys,
              });
              updateAbility();
            }
          },
          (err) => console.log('There was an error fetching available orgs:', err),
        );
      }
    },
    [store, updateAbility],
  );

  const loadWorkspace = useCallback(
    async (orgId) => {
      const currentWorkspace = sessionStorage.getItem('currentWorkspace');
      if (currentWorkspace && currentWorkspace !== 'undefined') {
        let workspace = JSON.parse(currentWorkspace);
        setWorkspace(workspace);
      } else {
        dataFetch(
          `/api/workspaces?search=&order=&page=0&pagesize=10&orgID=${orgId}`,
          {
            method: 'GET',
            credentials: 'include',
          },
          async (result) => {
            setWorkspace(result.workspaces[0]);
          },
          (err) => console.log('There was an error fetching workspaces:', err),
        );
      }
    },
    [setWorkspace],
  );

  const loadOrg = useCallback(async () => {
    const currentOrg = sessionStorage.getItem('currentOrg');
    let reFetchKeys = false;

    if (currentOrg && currentOrg !== 'undefined') {
      let org = JSON.parse(currentOrg);
      await loadAbility(org.id, reFetchKeys);
      setOrganization(org);
      await loadWorkspace(org.id);
    }

    dataFetch(
      '/api/identity/orgs',
      {
        method: 'GET',
        credentials: 'include',
      },
      async (result) => {
        let organizationToSet;
        const sessionOrg = currentOrg ? JSON.parse(currentOrg) : null;

        if (currentOrg) {
          const indx = result.organizations.findIndex((org) => org.id === sessionOrg.id);
          if (indx === -1) {
            organizationToSet = result.organizations[0];
            reFetchKeys = true;
            await loadAbility(organizationToSet.id, reFetchKeys);
            await loadWorkspace(organizationToSet.id);
            setOrganization(organizationToSet);
          }
        } else {
          organizationToSet = result.organizations[0];
          reFetchKeys = true;
          await loadWorkspace(organizationToSet.id);
          await loadAbility(organizationToSet.id, reFetchKeys);
          setOrganization(organizationToSet);
        }
      },
      (err) => console.log('There was an error fetching available orgs:', err),
    );
  }, [loadAbility, loadWorkspace, setOrganization]);

  const loadConfigFromServer = useCallback(async () => {
    dataFetch(
      '/api/system/sync',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (result) {
          if (
            result.meshAdapters &&
            result.meshAdapters !== null &&
            result.meshAdapters.length > 0
          ) {
            store.dispatch({
              type: actionTypes.UPDATE_ADAPTERS_INFO,
              meshAdapters: result.meshAdapters,
            });
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
            store.dispatch({ type: actionTypes.UPDATE_LOAD_GEN_CONFIG, loadTestPref });
          }
          if (typeof result.anonymousUsageStats !== 'undefined') {
            store.dispatch({
              type: actionTypes.UPDATE_ANONYMOUS_USAGE_STATS,
              anonymousUsageStats: result.anonymousUsageStats,
            });
          }
          if (typeof result.anonymousPerfResults !== 'undefined') {
            store.dispatch({
              type: actionTypes.UPDATE_ANONYMOUS_PERFORMANCE_RESULTS,
              anonymousPerfResults: result.anonymousPerfResults,
            });
          }
        }
      },
      (error) => {
        console.log(`there was an error fetching user config data: ${error}`);
      },
    );
  }, [store]);

  useEffect(() => {
    loadConfigFromServer();
    loadPromGrafanaConnection();
    loadOrg();
    initSubscriptions([]);

    dataFetch(
      '/api/user/prefs',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (typeof result?.usersExtensionPreferences?.catalogContent !== 'undefined') {
          toggleCatalogContent({
            catalogVisibility: result?.usersExtensionPreferences?.catalogContent,
          });
        }
      },
      (err) => console.error(err),
    );

    document.addEventListener('fullscreenchange', fullScreenChanged);
    loadMeshModelComponent();
    setState((prevState) => ({ ...prevState, isLoading: false }));

    // Cleanup
    return () => {
      document.removeEventListener('fullscreenchange', fullScreenChanged);
      if (state.disposeK8sContextSubscription) {
        state.disposeK8sContextSubscription();
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
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

    const { mesheryControllerSubscription } = state;
    console.log('amit mesheryControllerSubscription', mesheryControllerSubscription);
    if (mesheryControllerSubscription) {
      const ids = getK8sConfigIdsFromK8sConfig(k8sConfig);
      mesheryControllerSubscription.updateSubscription(
        getConnectionIDsFromContextIds(ids, k8sConfig),
      );
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [k8sConfig, capabilitiesRegistry]);

  const canShowNav = !state.isFullScreenMode && uiConfig?.components?.navigator !== false;

  return (
    <LoadingScreen message={randomLoadingMessage} isLoading={state.isLoading}>
      <DynamicComponentProvider>
        <RelayEnvironmentProvider environment={relayEnvironment}>
          <MesheryThemeProvider>
            <NoSsr>
              <ErrorBoundary customFallback={CustomErrorFallback}>
                <LoadSessionGuard>
                  <StyledRoot>
                    <CssBaseline />
                    <NavigationBar
                      isDrawerCollapsed={isDrawerCollapsed}
                      mobileOpen={state.mobileOpen}
                      handleDrawerToggle={handleDrawerToggle}
                      updateExtensionType={updateExtensionType}
                      canShowNav={canShowNav}
                    />
                    <StyledAppContent canShowNav={canShowNav}>
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
                              updateExtensionType={updateExtensionType}
                              abilityUpdated={state.abilityUpdated}
                            />
                          )}
                          <StyledContentWrapper>
                            <StyledMainContent
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
                              handleL5CommunityClick={handleL5CommunityClick}
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
                </LoadSessionGuard>
              </ErrorBoundary>
            </NoSsr>
          </MesheryThemeProvider>
        </RelayEnvironmentProvider>
      </DynamicComponentProvider>
    </LoadingScreen>
  );
};

// Keep the static getInitialProps method
MesheryApp.getInitialProps = async ({ Component, ctx }) => {
  const pageProps = Component.getInitialProps ? await Component.getInitialProps(ctx) : {};
  return { pageProps };
};

const mapStateToProps = (state) => ({
  isDrawerCollapsed: state.get('isDrawerCollapsed'),
  operatorSubscription: state.get('operatorSubscription'),
  telemetryURLs: state.get('telemetryURLs'),
  connectionMetadata: state.get('connectionMetadata'),
  extensionType: state.get('extensionType'),
  organization: state.get('organization'),
});

const mapDispatchToProps = (dispatch) => ({
  toggleCatalogContent: bindActionCreators(toggleCatalogContent, dispatch),
  setConnectionMetadata: bindActionCreators(setConnectionMetadata, dispatch),
});

const MesheryWithRedux = connect(mapStateToProps, mapDispatchToProps)(MesheryApp);

const MesheryThemeProvider = ({ children }) => {
  const themePref = useThemePreference();
  const mode = themePref?.data?.mode || 'dark';
  return <SistentThemeProvider initialMode={mode}>{children}</SistentThemeProvider>;
};

const MesheryAppWrapper = (props) => {
  return (
    <Provider store={store} context={RTKContext}>
      <Provider store={props.store} context={LegacyStoreContext}>
        <Provider store={props.store}>
          <Head>
            <link rel="shortcut icon" href="/static/img/meshery-logo/meshery-logo.svg" />
            <title>Meshery</title>
          </Head>
          <LocalizationProvider dateAdapter={AdapterMoment}>
            <MesheryWithRedux {...props} />
          </LocalizationProvider>
        </Provider>
      </Provider>
    </Provider>
  );
};

export default withRedux(makeStore, {
  serializeState: (state) => state.toJS(),
  deserializeState: (state) => fromJS(state),
})(MesheryAppWrapper);

const NavigationBar = ({
  isDrawerCollapsed,
  mobileOpen,
  handleDrawerToggle,

  updateExtensionType,
  canShowNav,
}) => {
  if (!canShowNav) {
    return null;
  }

  return (
    <StyledDrawer
      isDrawerCollapsed={isDrawerCollapsed}
      data-testid="navigation"
      id="left-navigation-bar"
    >
      <Hidden smUp implementation="js">
        <Navigator
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          isDrawerCollapsed={isDrawerCollapsed}
          updateExtensionType={updateExtensionType}
        />
      </Hidden>
      <Hidden xsDown implementation="css">
        <Navigator
          isDrawerCollapsed={isDrawerCollapsed}
          updateExtensionType={updateExtensionType}
        />
      </Hidden>
    </StyledDrawer>
  );
};
