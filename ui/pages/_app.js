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
import React, { useEffect } from 'react';
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
  updateTelemetryUrls,
  setConnectionMetadata,
  LegacyStoreContext,
  updateK8sContexts,
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
import { RTKContext } from '@/store/hooks';
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

const KubernetesSubscription = ({ store, setActiveContexts }) => {
  const k8sContextSubscription = (page = '', search = '', pageSize = '10', order = '') => {
    if (!CAN(keys.VIEW_ALL_KUBERNETES_CLUSTERS.action, keys.VIEW_ALL_KUBERNETES_CLUSTERS.subject)) {
      return () => {};
    }

    return subscribeK8sContext(
      (result) => {
        setActiveContexts('all');

        store.dispatch({
          type: actionTypes.UPDATE_CLUSTER_CONFIG,
          k8sConfig: result.k8sContext.contexts,
        });
      },
      { selector: { page, pageSize, order, search } },
    );
  };

  useEffect(() => {
    const disposeK8sContextSubscription = k8sContextSubscription();

    return () => {
      if (disposeK8sContextSubscription) {
        disposeK8sContextSubscription();
      }
    };
  }, []);

  return null;
};

const MesheryApp = ({ Component, pageProps, ...props }) => {
  const pageContext = getPageContext();
  const [mobileOpen, setMobileOpen] = React.useState(false);
  const [isFullScreenMode, setIsFullScreenMode] = React.useState(false);
  const [isLoading, setIsLoading] = React.useState(true);
  const [k8sContexts, setK8sContexts] = React.useState([]);
  const [activeK8sContexts, setActiveK8sContexts] = React.useState([]);
  const [isOpen, setIsOpen] = React.useState(false);
  const [keys, setKeys] = React.useState([]);
  const [abilityUpdated, setAbilityUpdated] = React.useState(false);
  const relayEnvironment = React.useState(createRelayEnvironment());
  const [mesheryControllerSubscription, setMesheryControllerSubscription] = React.useState(null);
  const prevK8sConfig = React.useRef(props.k8sConfig);

  const loadPromGrafanaConnection = () => {
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
  };

  const fullScreenChanged = () => {
    setIsFullScreenMode((state) => !state);
  };

  React.useEffect(() => {
    const loadData = async () => {
      await loadConfigFromServer(); // Ensure this finishes before proceeding
      loadPromGrafanaConnection();
      loadOrg();
      initSubscriptions([]);
      dataFetch(
        '/api/user/prefs',
        { method: 'GET', credentials: 'include' },
        (result) => {
          if (typeof result?.usersExtensionPreferences?.catalogContent !== 'undefined') {
            props.toggleCatalogContent({
              catalogVisibility: result?.usersExtensionPreferences?.catalogContent,
            });
          }
        },
        (err) => console.error(err),
      );
      loadMeshModelComponent();
      setIsLoading(false);
    };

    loadData();
    document.addEventListener('fullscreenchange', fullScreenChanged);

    return () => {
      document.removeEventListener('fullscreenchange', fullScreenChanged);
    };
  }, []);

  React.useEffect(() => {
    if (isMesheryUiRestrictedAndThePageIsNotPlayground(props.capabilitiesRegistry)) {
      Router.push(mesheryExtensionRoute);
    }
  }, [props.capabilitiesRegistry]); // Runs when `capabilitiesRegistry` changes

  React.useEffect(() => {
    if (!_.isEqual(prevK8sConfig.current, props.k8sConfig)) {
      console.log(
        'k8sconfig changed, re-initializing subscriptions',
        props.k8sConfig,
        activeK8sContexts,
      );
      const ids = getK8sConfigIdsFromK8sConfig(props.k8sConfig);
      if (mesheryControllerSubscription) {
        mesheryControllerSubscription.updateSubscription(
          getConnectionIDsFromContextIds(ids, props.k8sConfig),
        );
      }
    }
    prevK8sConfig.current = props.k8sConfig;
  }, [props.k8sConfig]);

  const loadMeshModelComponent = () => {
    const connectionDef = {};
    CONNECTION_KINDS_DEF.map(async (kind) => {
      const res = await getMeshModelComponentByName(formatToTitleCase(kind).concat('Connection'));
      if (res?.components) {
        connectionDef[CONNECTION_KINDS[kind]] = {
          transitions: res?.components[0].metadata.transitions,
          icon: res?.components[0].styles.svgColor,
        };
      }
      setConnectionMetadata(connectionDef);
    });
    props.setConnectionMetadata({
      connectionMetadataState: connectionDef,
    });
  };

  const initSubscriptions = (contexts) => {
    const connectionIDs = getConnectionIDsFromContextIds(contexts, props.k8sConfig);

    // No need to create a controller subscription if there are no connections
    if (connectionIDs.length < 1) {
      setMesheryControllerSubscription(() => {});
      return;
    }

    const mesheryControllerSubscription = new GQLSubscription({
      type: MESHERY_CONTROLLER_SUBSCRIPTION,
      connectionIDs: getConnectionIDsFromContextIds(contexts, props.k8sConfig),
      callbackFunction: (data) => {
        props.store.dispatch({
          type: actionTypes.SET_CONTROLLER_STATE,
          controllerState: data,
        });
      },
    });

    setMesheryControllerSubscription(mesheryControllerSubscription);
  };

  const handleDrawerToggle = () => {
    setMobileOpen((state) => !state);
  };

  const handleL5CommunityClick = () => {
    setIsOpen((state) => !state);
  };

  const activeContextChangeCallback = (activeK8sContexts) => {
    if (activeK8sContexts.includes('all')) {
      activeK8sContexts = ['all'];
    }
    props.store.dispatch(updateK8sContexts({ selectedK8sContexts: activeK8sContexts }));
  };

  const setActiveContexts = (id) => {
    if (k8sContexts?.contexts) {
      if (id === 'all') {
        let activeContexts = [];
        k8sContexts.contexts.forEach((ctx) => activeContexts.push(ctx.id));
        activeContexts.push('all');
        setActiveK8sContexts(activeContexts);
        activeContextChangeCallback(activeK8sContexts);
        return;
      }

      // if id is an empty array, clear all active contexts
      if (Array.isArray(id) && id.length === 0) {
        setActiveK8sContexts([]);
        activeContextChangeCallback(activeK8sContexts);
        return;
      }

      let ids = [...(activeK8sContexts || [])];
      //pop event
      if (ids.includes(id)) {
        ids = ids.filter((id) => id !== 'all');
        setActiveK8sContexts(ids.filter((cid) => cid !== id));
      }

      //push event
      if (ids.length === k8sContexts.contexts.length - 1) {
        ids.push('all');
      }
      setActiveK8sContexts([...ids, id]);
      activeContextChangeCallback(activeK8sContexts);
    }
  };

  const searchContexts = (search = '') => {
    fetchContexts(10, search)
      .then((ctx) => {
        setK8sContexts(ctx);
        const active = ctx?.contexts?.find((c) => c.is_current_context === true);
        if (active) setActiveK8sContexts([active?.id]);
      })
      .catch((err) => console.error(err));
  };

  const updateExtensionType = (type) => {
    props.store.dispatch({ type: actionTypes.UPDATE_EXTENSION_TYPE, extensionType: type });
  };

  const loadOrg = async () => {
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
        const sessionOrg = JSON.parse(currentOrg);

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
  };

  const setOrganization = (org) => {
    const { store } = props;
    store.dispatch({
      type: actionTypes.SET_ORGANIZATION,
      organization: org,
    });
  };

  const setWorkspace = (workspace) => {
    const { store } = props;
    store.dispatch({
      type: actionTypes.SET_WORKSPACE,
      workspace: workspace,
    });
  };

  const loadWorkspace = async (orgId) => {
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
  };

  const loadAbility = async (orgID, reFetchKeys) => {
    const storedKeys = sessionStorage.getItem('keys');
    const { store } = props;
    if (storedKeys !== null && !reFetchKeys && storedKeys !== 'undefined') {
      setKeys(JSON.parse(storedKeys));
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
            setKeys(result.keys);
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
  };

  const updateAbility = () => {
    ability.update(keys?.map((key) => ({ action: key.id, subject: _.lowerCase(key.function) })));
    setAbilityUpdated(true);
  };

  const loadConfigFromServer = () => {
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
            props.store.dispatch({
              type: actionTypes.UPDATE_ADAPTERS_INFO,
              meshAdapters: result.meshAdapters,
            });
          }
          // if (result.grafana) {
          //   const grafanaCfg = Object.assign(
          //     {
          //       grafanaURL: '',
          //       grafanaAPIKey: '',
          //       grafanaBoardSearch: '',
          //       grafanaBoards: [],
          //       selectedBoardsConfigs: [],
          //     },
          //     result.grafana,
          //   );
          //   props.store.dispatch({ type: actionTypes.UPDATE_GRAFANA_CONFIG, grafana: grafanaCfg });
          // }
          // if (result.prometheus) {
          //   if (typeof result.prometheus.prometheusURL === 'undefined') {
          //     result.prometheus.prometheusURL = '';
          //   }
          //   if (typeof result.prometheus.selectedPrometheusBoardsConfigs === 'undefined') {
          //     result.prometheus.selectedPrometheusBoardsConfigs = [];
          //   }
          //   const promCfg = Object.assign(
          //     {
          //       prometheusURL: '',
          //       selectedPrometheusBoardsConfigs: [],
          //     },
          //     result.prometheus,
          //   );
          //   props.store.dispatch({ type: actionTypes.UPDATE_PROMETHEUS_CONFIG, prometheus: promCfg });
          // }
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
            props.store.dispatch({ type: actionTypes.UPDATE_LOAD_GEN_CONFIG, loadTestPref });
          }
          if (typeof result.anonymousUsageStats !== 'undefined') {
            props.store.dispatch({
              type: actionTypes.UPDATE_ANONYMOUS_USAGE_STATS,
              anonymousUsageStats: result.anonymousUsageStats,
            });
          }
          if (typeof result.anonymousPerfResults !== 'undefined') {
            props.store.dispatch({
              type: actionTypes.UPDATE_ANONYMOUS_PERFORMANCE_RESULTS,
              anonymousPerfResults: result.anonymousPerfResults,
            });
          }
        }
      },
    );
  };

  return (
    <LoadingScreen message={randomLoadingMessage} isLoading={isLoading}>
      <DynamicComponentProvider>
        <RelayEnvironmentProvider environment={relayEnvironment}>
          <MesheryThemeProvider>
            <NoSsr>
              <ErrorBoundary customFallback={CustomErrorFallback}>
                <LoadSessionGuard>
                  <StyledRoot>
                    <CssBaseline />
                    <NavigationBar
                      isDrawerCollapsed={props.isDrawerCollapsed}
                      mobileOpen={mobileOpen}
                      handleDrawerToggle={handleDrawerToggle}
                      isFullScreenMode={isFullScreenMode}
                      updateExtensionType={updateExtensionType}
                    />
                    <StyledAppContent>
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
                          <KubernetesSubscription
                            store={props.store}
                            setActiveContexts={setActiveContexts}
                          />
                          {!isFullScreenMode && (
                            <Header
                              onDrawerToggle={handleDrawerToggle}
                              onDrawerCollapse={props.isDrawerCollapsed}
                              contexts={k8sContexts}
                              activeContexts={activeK8sContexts}
                              setActiveContexts={setActiveContexts}
                              searchContexts={searchContexts}
                              updateExtensionType={updateExtensionType}
                              abilityUpdated={abilityUpdated}
                            />
                          )}
                          <StyledContentWrapper>
                            <StyledMainContent
                              style={{
                                padding: props.extensionType === 'navigator' && '0px',
                              }}
                            >
                              <LocalizationProvider dateAdapter={AdapterMoment}>
                                <ErrorBoundary customFallback={CustomErrorFallback}>
                                  <Component
                                    pageContext={pageContext}
                                    contexts={k8sContexts}
                                    activeContexts={activeK8sContexts}
                                    setActiveContexts={setActiveContexts}
                                    searchContexts={searchContexts}
                                    {...pageProps}
                                  />
                                </ErrorBoundary>
                              </LocalizationProvider>
                            </StyledMainContent>
                            <Footer
                              handleL5CommunityClick={handleL5CommunityClick}
                              capabilitiesRegistry={props.capabilitiesRegistry}
                            />
                          </StyledContentWrapper>
                        </NotificationCenterProvider>
                      </SnackbarProvider>
                    </StyledAppContent>
                  </StyledRoot>
                  <PlaygroundMeshDeploy closeForm={() => setIsOpen(false)} isOpen={isOpen} />
                </LoadSessionGuard>
              </ErrorBoundary>
            </NoSsr>
          </MesheryThemeProvider>
        </RelayEnvironmentProvider>
      </DynamicComponentProvider>
    </LoadingScreen>
  );
};

MesheryApp.getInitialProps = async ({ Component, ctx }) => {
  const pageProps = Component.getInitialProps ? await Component.getInitialProps(ctx) : {};
  return { pageProps };
};

const mapStateToProps = (state) => ({
  isDrawerCollapsed: state.get('isDrawerCollapsed'),
  k8sConfig: state.get('k8sConfig'),
  operatorSubscription: state.get('operatorSubscription'),
  capabilitiesRegistry: state.get('capabilitiesRegistry'),
  telemetryURLs: state.get('telemetryURLs'),
  connectionMetadata: state.get('connectionMetadata'),
  extensionType: state.get('extensionType'),
  organization: state.get('organization'),
});

const mapDispatchToProps = (dispatch) => ({
  toggleCatalogContent: bindActionCreators(toggleCatalogContent, dispatch),
  updateTelemetryUrls: bindActionCreators(updateTelemetryUrls, dispatch),
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
  isFullScreenMode,
  updateExtensionType,
}) => {
  const canShowNav = !isFullScreenMode && uiConfig?.components?.navigator !== false;

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
