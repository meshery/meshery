import MomentUtils from '@date-io/moment';
import { NoSsr, Typography } from '@material-ui/core';
import CssBaseline from '@material-ui/core/CssBaseline';
import Hidden from '@material-ui/core/Hidden';
import { ThemeProvider, withStyles } from '@material-ui/core/styles';
import { CheckCircle, Error, Info, Warning } from '@material-ui/icons';
import FavoriteIcon from '@material-ui/icons/Favorite';
import { MuiPickersUtilsProvider } from '@material-ui/pickers';
// import 'billboard.js/dist/theme/insight.min.css';
import 'billboard.js/dist/theme/dark.min.css';
// import 'billboard.js/dist/billboard. min.css';
import 'codemirror/addon/lint/lint.css';
// codemirror + js-yaml imports when added to a page was preventing to navigating to that page using nextjs
// link clicks, hence attempting to add them here
import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/material.css';
import { fromJS } from 'immutable';
import _ from 'lodash';
import withRedux from 'next-redux-wrapper';
import App from 'next/app';
import Head from 'next/head';
import { SnackbarProvider } from 'notistack';
import PropTypes from 'prop-types';
import React from 'react';
import { connect, Provider } from 'react-redux';
import Header from '../components/Header';
import MesheryProgressBar from '../components/MesheryProgressBar';
import Navigator from '../components/Navigator';
import getPageContext from '../components/PageContext';
import { MESHERY_CONTROLLER_SUBSCRIPTION } from '../components/subscription/helpers';
import { GQLSubscription } from '../components/subscription/subscriptionhandler';
import dataFetch, { promisifiedDataFetch } from '../lib/data-fetch';
import {
  actionTypes,
  makeStore,
  toggleCatalogContent,
  updateTelemetryUrls,
  setConnectionMetadata,
} from '../lib/store';
import theme, { styles } from '../themes';
import { getK8sConfigIdsFromK8sConfig } from '../utils/multi-ctx';
import './../public/static/style/index.css';
import subscribeK8sContext from '../components/graphql/subscriptions/K8sContextSubscription';
import { bindActionCreators } from 'redux';
import { darkTheme } from '../themes/app';
import './styles/AnimatedFilter.css';
import './styles/AnimatedMeshery.css';
import './styles/AnimatedMeshPattern.css';
import './styles/AnimatedMeshSync.css';
import PlaygroundMeshDeploy from './extension/AccessMesheryModal';
import Router from 'next/router';
import subscribeMeshSyncEvents from '../components/graphql/subscriptions/MeshSyncEventsSubscription';
import { isTelemetryComponent, TelemetryComps } from '../utils/nameMapper';
import { extractURLFromScanData } from '../components/ConnectionWizard/helpers/metrics';
import { updateURLs } from '../utils/utils';
import { RelayEnvironmentProvider } from 'react-relay';
import { createRelayEnvironment } from '../lib/relayEnvironment';
import './styles/charts.css';

import { ErrorBoundary } from '../components/General/ErrorBoundary';
import { NotificationCenterProvider } from '../components/NotificationCenter';
import { getMeshModelComponent } from '../api/meshmodel';
import { CONNECTION_KINDS } from '../utils/Enum';

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

class MesheryApp extends App {
  constructor() {
    super();
    this.pageContext = getPageContext();
    this.meshsyncEventsSubscriptionRef = React.createRef();
    this.eventsSubscriptionRef = React.createRef();
    this.fullScreenChanged = this.fullScreenChanged.bind(this);
    this.state = {
      mobileOpen: false,
      isDrawerCollapsed: false,
      isFullScreenMode: false,
      k8sContexts: [],
      activeK8sContexts: [],
      operatorSubscription: null,
      mesheryControllerSubscription: null,
      meshSyncSubscription: null,
      disposeK8sContextSubscription: null,
      theme: 'light',
      isOpen: false,
      relayEnvironment: createRelayEnvironment(),
      connectionMetadata: {},
    };
  }

  initMeshSyncEventsSubscription(contexts = []) {
    if (this.meshsyncEventsSubscriptionRef.current) {
      this.meshsyncEventsSubscriptionRef.current.dispose();
    }

    const meshSyncEventsSubscription = subscribeMeshSyncEvents(
      (result) => {
        if (result.meshsyncevents.object.kind === 'Service') {
          const telemetryCompName = isTelemetryComponent(
            result.meshsyncevents.object?.metadata?.name,
          );
          let prometheusURLs = [];
          let grafanaURLs = [];

          const grafanaUrlsSet = new Set(this.props.telemetryURLs.grafana);
          const promUrlsSet = new Set(this.props.telemetryURLs.prometheus);

          const eventType = result.meshsyncevents.type;
          const spec = result?.meshsyncevents?.object?.spec?.attribute;
          const status = result?.meshsyncevents?.object?.status?.attribute;
          const data = { spec: JSON.parse(spec), status: JSON.parse(status) };

          if (telemetryCompName === TelemetryComps.GRAFANA) {
            grafanaURLs = grafanaURLs.concat(extractURLFromScanData(data));
            updateURLs(grafanaUrlsSet, grafanaURLs, eventType);
          } else if (telemetryCompName === TelemetryComps.PROMETHEUS) {
            prometheusURLs = new Set(prometheusURLs.concat(extractURLFromScanData(data)));
            updateURLs(promUrlsSet, prometheusURLs, eventType);
          }

          this.props.updateTelemetryUrls({
            telemetryURLs: {
              grafana: Array.from(grafanaUrlsSet),
              prometheus: Array.from(promUrlsSet),
            },
          });
        }
      },
      {
        k8scontextIDs: contexts,
        eventTypes: ['ADDED', 'DELETED'],
      },
    );

    this.meshsyncEventsSubscriptionRef.current = meshSyncEventsSubscription;
  }

  fullScreenChanged = () => {
    this.setState((state) => {
      return { isFullScreenMode: !state.isFullScreenMode };
    });
  };

  componentDidMount() {
    this.loadConfigFromServer(); // this works, but sometimes other components which need data load faster than this data is obtained.
    this.initSubscriptions([]);
    dataFetch(
      '/api/user/prefs',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (typeof result?.usersExtensionPreferences?.catalogContent !== 'undefined') {
          this.props.toggleCatalogContent({
            catalogVisibility: result?.usersExtensionPreferences?.catalogContent,
          });
        }
      },
      (err) => console.error(err),
    );

    this.initMeshSyncEventsSubscription(this.state.activeK8sContexts);
    // this.initEventsSubscription()
    const k8sContextSubscription = (page = '', search = '', pageSize = '10', order = '') => {
      return subscribeK8sContext(
        (result) => {
          this.setState({ k8sContexts: result.k8sContext }, () => this.setActiveContexts('all'));
          this.props.store.dispatch({
            type: actionTypes.UPDATE_CLUSTER_CONFIG,
            k8sConfig: result.k8sContext.contexts,
          });
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
    const disposeK8sContextSubscription = k8sContextSubscription();
    this.setState({ disposeK8sContextSubscription });

    document.addEventListener('fullscreenchange', this.fullScreenChanged);
    this.loadMeshModelComponent();
  }

  loadMeshModelComponent = () => {
    const connectionDef = {};
    Object.keys(CONNECTION_KINDS).map(async (kind) => {
      const connectionKind =
        CONNECTION_KINDS[kind] === 'meshery' ? 'meshery-core' : CONNECTION_KINDS[kind];
      const res = await getMeshModelComponent(connectionKind, 'Connection');
      if (res?.components) {
        connectionDef[CONNECTION_KINDS[kind]] = {
          transitions: res?.components[0].model.metadata.transitions,
          icon: res?.components[0].metadata.svgColor,
        };
      }
      this.setState({ connectionMetadata: connectionDef });
      this.props.setConnectionMetadata({
        connectionMetadataState: connectionDef,
      });
    });
  };

  componentWillUnmount() {
    document.removeEventListener('fullscreenchange', this.fullScreenChanged);
  }

  componentDidUpdate(prevProps) {
    const { k8sConfig, capabilitiesRegistry } = this.props;

    // in case the meshery-ui is restricted, the user will be redirected to signup/extension page
    if (isMesheryUiRestrictedAndThePageIsNotPlayground(capabilitiesRegistry)) {
      Router.push(mesheryExtensionRoute);
    }

    if (!_.isEqual(prevProps.k8sConfig, k8sConfig)) {
      const { meshSyncSubscription, mesheryControllerSubscription } = this.state;
      console.log(
        'k8sconfig changed, re-initialising subscriptions',
        k8sConfig,
        this.state.activeK8sContexts,
      );
      const ids = getK8sConfigIdsFromK8sConfig(k8sConfig);
      // if (operatorSubscription) {
      //   operatorSubscription.updateSubscription(ids);
      // }

      if (mesheryControllerSubscription) {
        mesheryControllerSubscription.updateSubscription(ids);
      }
      if (meshSyncSubscription) {
        meshSyncSubscription.updateSubscription(ids);
      }

      if (this.meshsyncEventsSubscriptionRef.current) {
        this.initMeshSyncEventsSubscription(ids);
      }
    }
  }

  initSubscriptions = (contexts) => {
    const mesheryControllerSubscription = new GQLSubscription({
      type: MESHERY_CONTROLLER_SUBSCRIPTION,
      contextIds: contexts,
      callbackFunction: (data) => {
        this.props.store.dispatch({
          type: actionTypes.SET_CONTROLLER_STATE,
          controllerState: data,
        });
        console.log('CONTROLLER TEST CONTROLLER ', data);
      },
    });
    // const meshSyncSubscription = new GQLSubscription({ type : MESHSYNC_EVENT_SUBSCRIPTION, contextIds : contexts, callbackFunction : meshSyncCallback }) above uses old listenToMeshSyncEvents subscription, instead new subscribeMeshSyncEvents is used
    this.setState({ mesheryControllerSubscription });
    // this.setState({ operatorSubscription });
  };

  handleDrawerToggle = () => {
    this.setState((state) => ({ mobileOpen: !state.mobileOpen }));
  };

  handleL5CommunityClick = () => {
    this.setState((state) => ({ isOpen: !state.isOpen }));
  };

  /**
   * Sets the selected k8s context on global level.
   * @param {Array.<string>} activeK8sContexts
   */
  activeContextChangeCallback = (activeK8sContexts) => {
    if (activeK8sContexts.includes('all')) {
      activeK8sContexts = ['all'];
    }
    this.props.store.dispatch({
      type: actionTypes.SET_K8S_CONTEXT,
      selectedK8sContexts: activeK8sContexts,
    });
  };

  setActiveContexts = (id) => {
    if (this.state.k8sContexts?.contexts) {
      if (id === 'all') {
        let activeContexts = [];
        this.state.k8sContexts.contexts.forEach((ctx) => activeContexts.push(ctx.id));
        activeContexts.push('all');
        this.setState({ activeK8sContexts: activeContexts }, () =>
          this.activeContextChangeCallback(this.state.activeK8sContexts),
        );
        return;
      }

      this.setState(
        (state) => {
          let ids = [...(state.activeK8sContexts || [])];
          //pop event
          if (ids.includes(id)) {
            ids = ids.filter((id) => id !== 'all');
            return { activeK8sContexts: ids.filter((cid) => cid !== id) };
          }

          //push event
          if (ids.length === this.state.k8sContexts.contexts.length - 1) {
            ids.push('all');
          }
          return { activeK8sContexts: [...ids, id] };
        },
        () => this.activeContextChangeCallback(this.state.activeK8sContexts),
      );
    }
  };

  searchContexts = (search = '') => {
    fetchContexts(10, search)
      .then((ctx) => {
        this.setState({ k8sContexts: ctx });
        const active = ctx?.contexts?.find((c) => c.is_current_context === true);
        if (active) this.setState({ activeK8sContexts: [active?.id] });
      })
      .catch((err) => console.error(err));
  };

  updateExtensionType = (type) => {
    this.props.store.dispatch({ type: actionTypes.UPDATE_EXTENSION_TYPE, extensionType: type });
  };

  async loadConfigFromServer() {
    const { store } = this.props;
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
          if (result.grafana) {
            const grafanaCfg = Object.assign(
              {
                grafanaURL: '',
                grafanaAPIKey: '',
                grafanaBoardSearch: '',
                grafanaBoards: [],
                selectedBoardsConfigs: [],
              },
              result.grafana,
            );
            store.dispatch({ type: actionTypes.UPDATE_GRAFANA_CONFIG, grafana: grafanaCfg });
          }
          if (result.prometheus) {
            if (typeof result.prometheus.prometheusURL === 'undefined') {
              result.prometheus.prometheusURL = '';
            }
            if (typeof result.prometheus.selectedPrometheusBoardsConfigs === 'undefined') {
              result.prometheus.selectedPrometheusBoardsConfigs = [];
            }
            const promCfg = Object.assign(
              {
                prometheusURL: '',
                selectedPrometheusBoardsConfigs: [],
              },
              result.prometheus,
            );
            store.dispatch({ type: actionTypes.UPDATE_PROMETHEUS_CONFIG, prometheus: promCfg });
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
  }

  static async getInitialProps({ Component, ctx }) {
    const pageProps = Component.getInitialProps ? await Component.getInitialProps(ctx) : {};
    return { pageProps };
  }
  themeSetter = (thememode) => {
    this.setState({ theme: thememode });
  };
  render() {
    const { Component, pageProps, classes, isDrawerCollapsed, relayEnvironment } = this.props;

    return (
      <RelayEnvironmentProvider environment={relayEnvironment}>
        <ThemeProvider theme={this.state.theme === 'dark' ? darkTheme : theme}>
          <NoSsr>
            <ErrorBoundary>
              <div className={classes.root}>
                <CssBaseline />
                {!this.state.isFullScreenMode && (
                  <nav
                    className={isDrawerCollapsed ? classes.drawerCollapsed : classes.drawer}
                    data-test="navigation"
                  >
                    <Hidden smUp implementation="js">
                      <Navigator
                        variant="temporary"
                        open={this.state.mobileOpen}
                        onClose={this.handleDrawerToggle}
                        onCollapseDrawer={(open = null) => this.handleCollapseDrawer(open)}
                        isDrawerCollapsed={isDrawerCollapsed}
                        updateExtensionType={this.updateExtensionType}
                      />
                    </Hidden>
                    <Hidden xsDown implementation="css">
                      <Navigator
                        onCollapseDrawer={(open = null) => this.handleCollapseDrawer(open)}
                        isDrawerCollapsed={isDrawerCollapsed}
                        updateExtensionType={this.updateExtensionType}
                      />
                    </Hidden>
                  </nav>
                )}
                <div className={classes.appContent}>
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
                    classes={{
                      variantSuccess:
                        this.state.theme === 'dark'
                          ? classes.darknotifSuccess
                          : classes.notifSuccess,
                      variantError:
                        this.state.theme === 'dark' ? classes.darknotifError : classes.notifError,
                      variantWarning:
                        this.state.theme === 'dark' ? classes.darknotifWarn : classes.notifWarn,
                      variantInfo:
                        this.state.theme === 'dark' ? classes.darknotifInfo : classes.notifInfo,
                    }}
                    maxSnack={10}
                  >
                    <NotificationCenterProvider>
                      <MesheryProgressBar />
                      {!this.state.isFullScreenMode && (
                        <Header
                          onDrawerToggle={this.handleDrawerToggle}
                          onDrawerCollapse={isDrawerCollapsed}
                          contexts={this.state.k8sContexts}
                          activeContexts={this.state.activeK8sContexts}
                          setActiveContexts={this.setActiveContexts}
                          searchContexts={this.searchContexts}
                          updateExtensionType={this.updateExtensionType}
                          theme={this.state.theme}
                          themeSetter={this.themeSetter}
                        />
                      )}
                    </NotificationCenterProvider>
                    <main className={classes.mainContent}>
                      <MuiPickersUtilsProvider utils={MomentUtils}>
                        <ErrorBoundary>
                          <Component
                            pageContext={this.pageContext}
                            contexts={this.state.k8sContexts}
                            activeContexts={this.state.activeK8sContexts}
                            setActiveContexts={this.setActiveContexts}
                            searchContexts={this.searchContexts}
                            theme={this.state.theme}
                            themeSetter={this.themeSetter}
                            {...pageProps}
                          />
                        </ErrorBoundary>
                      </MuiPickersUtilsProvider>
                    </main>
                  </SnackbarProvider>
                  <footer
                    className={
                      this.props.capabilitiesRegistry?.restrictedAccess?.isMesheryUiRestricted
                        ? classes.playgroundFooter
                        : this.state.theme === 'dark'
                        ? classes.footerDark
                        : classes.footer
                    }
                  >
                    <Typography
                      variant="body2"
                      align="center"
                      color="textSecondary"
                      component="p"
                      style={
                        this.props.capabilitiesRegistry?.restrictedAccess?.isMesheryUiRestricted
                          ? { color: '#000' }
                          : {}
                      }
                    >
                      <span onClick={this.handleL5CommunityClick} className={classes.footerText}>
                        {this.props.capabilitiesRegistry?.restrictedAccess
                          ?.isMesheryUiRestricted ? (
                          'ACCESS LIMITED IN MESHERY PLAYGROUND. DEPLOY MESHERY TO ACCESS ALL FEATURES.'
                        ) : (
                          <>
                            {' '}
                            Built with{' '}
                            <FavoriteIcon
                              style={{
                                color:
                                  this.state.theme === 'dark'
                                    ? theme.palette.secondary.focused
                                    : '#00b39f',
                              }}
                              className={classes.footerIcon}
                            />{' '}
                            by the Layer5 Community
                          </>
                        )}
                      </span>
                    </Typography>
                  </footer>
                </div>
              </div>
              <PlaygroundMeshDeploy
                closeForm={() => this.setState({ isOpen: false })}
                isOpen={this.state.isOpen}
              />
            </ErrorBoundary>
          </NoSsr>
        </ThemeProvider>
      </RelayEnvironmentProvider>
    );
  }
}

MesheryApp.propTypes = { classes: PropTypes.object.isRequired };

const mapStateToProps = (state) => ({
  isDrawerCollapsed: state.get('isDrawerCollapsed'),
  k8sConfig: state.get('k8sConfig'),
  operatorSubscription: state.get('operatorSubscription'),
  meshSyncSubscription: state.get('meshSyncSubscription'),
  capabilitiesRegistry: state.get('capabilitiesRegistry'),
  telemetryURLs: state.get('telemetryURLs'),
  connectionMetadata: state.get('connectionMetadata'),
});

const mapDispatchToProps = (dispatch) => ({
  toggleCatalogContent: bindActionCreators(toggleCatalogContent, dispatch),
  updateTelemetryUrls: bindActionCreators(updateTelemetryUrls, dispatch),
  setConnectionMetadata: bindActionCreators(setConnectionMetadata, dispatch),
});

const MesheryWithRedux = withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(MesheryApp),
);

const MesheryAppWrapper = (props) => {
  return (
    <Provider store={props.store}>
      <Head>
        <link rel="shortcut icon" href="/static/img/meshery-logo/meshery-logo.svg" />
        <title>Meshery</title>
      </Head>
      <MuiPickersUtilsProvider utils={MomentUtils}>
        <MesheryWithRedux {...props} />
      </MuiPickersUtilsProvider>
    </Provider>
  );
};

// export default withStyles(styles)(withRedux(makeStore, {
//   serializeState : state => state.toJS(),
//   deserializeState : state => fromJS(state)
// })(MesheryAppWrapper));
export default withStyles(styles)(
  withRedux(makeStore, {
    serializeState: (state) => state.toJS(),
    deserializeState: (state) => fromJS(state),
  })(MesheryAppWrapper),
);
