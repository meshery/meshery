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
import _ from 'lodash';
import Head from 'next/head';
import { SnackbarProvider } from 'notistack';
import React, { useEffect, useMemo, useCallback, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import Header from '../components/Header';
import MesheryProgressBar from '../components/MesheryProgressBar';
import Navigator from '../components/Navigator';
import getPageContext from '../components/PageContext';
import { MESHERY_CONTROLLER_SUBSCRIPTION } from '../components/subscription/helpers';
import { GQLSubscription } from '../components/subscription/subscriptionhandler';
import dataFetch, { promisifiedDataFetch } from '../lib/data-fetch';
import 'react-grid-layout/css/styles.css';
import 'react-resizable/css/styles.css';
import { getConnectionIDsFromContextIds, getK8sConfigIdsFromK8sConfig } from '../utils/multi-ctx';
import './../public/static/style/index.css';
import subscribeK8sContext from '../components/graphql/subscriptions/K8sContextSubscription';
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
} from '@sistent/sistent';
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
import {
  setConnectionMetadata,
  setControllerState,
  setK8sContexts,
  setKeys,
  setOrganization,
  toggleCatalogContent,
  updateExtensionType,
  updateK8SConfig,
} from '@/store/slices/mesheryUi';
import { updateLoadTestPref } from '@/store/slices/prefTest';
import { updateGrafanaConfig, updatePrometheusConfig } from '@/store/slices/telemetry';
import { updateAdaptersInfo } from '@/store/slices/adapter';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';
import WorkspaceModalContextProvider from '@/utils/context/WorkspaceModalContextProvider';

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

const Footer = ({ capabilitiesRegistry, handleMesheryCommunityClick }) => {
  const theme = useTheme();

  const { extensionType: extension } = useSelector((state) => state.ui);

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
        <StyledFooterText onClick={handleMesheryCommunityClick}>
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
              by the Meshery Community
            </>
          )}
        </StyledFooterText>
      </Typography>
    </StyledFooterBody>
  );
};

const KubernetesSubscription = ({ setAppState }) => {
  const dispatch = useDispatch();
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

const MesheryApp = ({ Component, pageProps, relayEnvironment }) => {
  const pageContext = useMemo(() => getPageContext(), []);
  const { k8sConfig } = useSelector((state) => state.ui);
  const { capabilitiesRegistry } = useSelector((state) => state.ui);
  const { isDrawerCollapsed } = useSelector((state) => state.ui);
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

            dispatch(updatePrometheusConfig(promCfg));
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
              dispatch(updateGrafanaConfig(grafanaCfg));
            });
          }
        });
      },
    );
  }, [dispatch]);

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
      if (!k8sConfig?.length) return;
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
        dataFetch(
          `/api/identity/orgs/${orgID}/users/keys`,
          {
            method: 'GET',
            credentials: 'include',
          },
          (result) => {
            if (result) {
              setState((prevState) => ({ ...prevState, keys: result.keys }));
              dispatch(setKeys({ keys: result.keys }));
              updateAbility();
            }
          },
          (err) => console.log('There was an error fetching available orgs:', err),
        );
      }
    },
    [dispatch, updateAbility],
  );

  const loadOrg = useCallback(async () => {
    const currentOrg = sessionStorage.getItem('currentOrg');
    let reFetchKeys = false;

    if (currentOrg && currentOrg !== 'undefined') {
      let org = JSON.parse(currentOrg);
      await loadAbility(org.id, reFetchKeys);
      setCurrentOrganization(org);
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
            setCurrentOrganization(organizationToSet);
          }
        } else {
          organizationToSet = result.organizations[0];
          reFetchKeys = true;
          await loadAbility(organizationToSet.id, reFetchKeys);
          setCurrentOrganization(organizationToSet);
        }
      },
      (err) => console.log('There was an error fetching available orgs:', err),
    );
  }, [loadAbility, setCurrentOrganization]);

  const loadConfigFromServer = useCallback(() => {
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
      },
      (error) => {
        console.log(`there was an error fetching user config data: ${error}`);
      },
    );
  }, [dispatch]);

  useEffect(() => {
    // todo further refactoring required for data fetch
    const loadAll = async () => {
      loadConfigFromServer();
      loadPromGrafanaConnection();
      await loadOrg();

      initSubscriptions([]);

      dataFetch(
        '/api/user/prefs',
        {
          method: 'GET',
          credentials: 'include',
        },
        (result) => {
          if (typeof result?.usersExtensionPreferences?.catalogContent !== 'undefined') {
            dispatch(
              toggleCatalogContent({
                catalogVisibility: result?.usersExtensionPreferences?.catalogContent,
              }),
            );
          }
        },
        (err) => console.error(err),
      );

      document.addEventListener('fullscreenchange', fullScreenChanged);
      await loadMeshModelComponent();
      setState((prevState) => ({ ...prevState, isLoading: false }));
    };
    loadAll();

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
    if (mesheryControllerSubscription) {
      const ids = getK8sConfigIdsFromK8sConfig(k8sConfig);
      mesheryControllerSubscription.updateSubscription(
        getConnectionIDsFromContextIds(ids, k8sConfig),
      );
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [k8sConfig, capabilitiesRegistry]);

  const canShowNav = !state.isFullScreenMode && uiConfig?.components?.navigator !== false;
  const { extensionType } = useSelector((state) => state.ui);

  return (
    <LoadingScreen message={randomLoadingMessage} isLoading={state.isLoading}>
      <DynamicComponentProvider>
        <RelayEnvironmentProvider environment={relayEnvironment}>
          <MesheryThemeProvider>
            <NoSsr>
              <ErrorBoundary customFallback={CustomErrorFallback}>
                <LoadSessionGuard>
                  <WorkspaceModalContextProvider>
                    <StyledRoot>
                      <CssBaseline />
                      <NavigationBar
                        isDrawerCollapsed={isDrawerCollapsed}
                        mobileOpen={state.mobileOpen}
                        handleDrawerToggle={handleDrawerToggle}
                        updateExtensionType={updateCurrentExtensionType}
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
                                updateExtensionType={updateCurrentExtensionType}
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
                  </WorkspaceModalContextProvider>
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

const MesheryThemeProvider = ({ children }) => {
  const themePref = useThemePreference();
  const mode = themePref?.data?.mode || 'dark';
  return <SistentThemeProvider initialMode={mode}>{children}</SistentThemeProvider>;
};

const MesheryAppWrapper = (props) => {
  return (
    <ProviderStoreWrapper>
      <Head>
        <link rel="shortcut icon" href="/static/img/meshery-logo/meshery-logo.svg" />
        <title>Meshery</title>
      </Head>
      <LocalizationProvider dateAdapter={AdapterMoment}>
        <MesheryApp {...props} />
      </LocalizationProvider>
    </ProviderStoreWrapper>
  );
};

export default MesheryAppWrapper;

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
