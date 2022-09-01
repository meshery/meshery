import MomentUtils from '@date-io/moment';
import { NoSsr, Typography } from '@material-ui/core';
import CssBaseline from '@material-ui/core/CssBaseline';
import Hidden from '@material-ui/core/Hidden';
import { MuiThemeProvider, withStyles } from '@material-ui/core/styles';
import {
  CheckCircle, Error, Info, Warning
} from '@material-ui/icons';
import FavoriteIcon from '@material-ui/icons/Favorite';
import { MuiPickersUtilsProvider } from '@material-ui/pickers';
// import 'billboard.js/dist/theme/insight.min.css';
// import 'billboard.js/dist/theme/graph.min.css';
import 'billboard.js/dist/billboard.min.css';
import 'codemirror/addon/lint/lint.css';
// codemirror + js-yaml imports when added to a page was preventing to navigating to that page using nextjs
// link clicks, hence attempting to add them here
import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/material.css';
import { fromJS } from 'immutable';
import _ from 'lodash';
import withRedux from "next-redux-wrapper";
import App from 'next/app';
import Head from 'next/head';
import { SnackbarProvider } from 'notistack';
import PropTypes from 'prop-types';
import React from 'react';
import { connect, Provider } from "react-redux";
import Header from '../components/Header';
import MesheryProgressBar from '../components/MesheryProgressBar';
import Navigator from '../components/Navigator';
import getPageContext from '../components/PageContext';
import { MESHSYNC_EVENT_SUBSCRIPTION, OPERATOR_EVENT_SUBSCRIPTION } from '../components/subscription/helpers';
import { GQLSubscription } from '../components/subscription/subscriptionhandler';
import dataFetch, { promisifiedDataFetch } from '../lib/data-fetch';
import { actionTypes, makeStore, toggleCatalogContent } from '../lib/store';
import theme, { styles } from "../themes";
import { getK8sConfigIdsFromK8sConfig } from '../utils/multi-ctx';
import './../public/static/style/index.css';
import subscribeK8sContext from "../components/graphql/subscriptions/K8sContextSubscription";
import { bindActionCreators } from 'redux';

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

async function fetchContexts(number = 10, search = "") {
  return await promisifiedDataFetch(`/api/system/kubernetes/contexts?pageSize=${number}&search=${encodeURIComponent(search)}`)
}

class MesheryApp extends App {
  constructor() {
    super();
    this.pageContext = getPageContext();

    this.state = {
      mobileOpen : false,
      isDrawerCollapsed : false,
      k8sContexts : [],
      activeK8sContexts : [],
      operatorSubscription : null,
      meshSyncSubscription : null,
      disposeK8sContextSubscription : null
    };
  }

  componentDidMount() {
    this.loadConfigFromServer(); // this works, but sometimes other components which need data load faster than this data is obtained.
    this.initSubscriptions([]);
    dataFetch(
      "/api/user/prefs",
      { credentials : "same-origin",
        method : "GET",
        credentials : "include", },
      (result) => {
        if (result) {
          this.props.toggleCatalogContent({
            catalogVisibility : result?.usersExtensionPreferences?.catalogContent || true
          })
        }
      },
      err => console.error(err)
    )
    const k8sContextSubscription = (page="", search="", pageSize="10", order="") => {
      return subscribeK8sContext((result) => {
        this.setState({ k8sContexts : result.k8sContext }, () => this.setActiveContexts("all"))
        this.props.store.dispatch({ type : actionTypes.UPDATE_CLUSTER_CONFIG, k8sConfig : result.k8sContext.contexts });
      },
      {
        selector : {
          page : page,
          pageSize : pageSize,
          order : order,
          search : search
        }
      })
    }
    const disposeK8sContextSubscription = k8sContextSubscription();
    this.setState({ disposeK8sContextSubscription })
  }

  componentDidUpdate(prevProps) {
    const { k8sConfig } = this.props;
    if (!_.isEqual(prevProps.k8sConfig, k8sConfig)) {
      const { operatorSubscription, meshSyncSubscription } = this.state;
      console.log("k8sconfig changed, re-initialising subscriptions");
      const ids = getK8sConfigIdsFromK8sConfig(k8sConfig)
      if (operatorSubscription) {
        operatorSubscription.updateSubscription(ids);
      }

      if (meshSyncSubscription) {
        meshSyncSubscription.updateSubscription(ids);
      }
    }
  }

  initSubscriptions = (contexts) => {
    const operatorCallback = (data) => {
      this.props.store.dispatch({ type : actionTypes.SET_OPERATOR_SUBSCRIPTION, operatorState : data });
    }

    const meshSyncCallback = (data) => {
      this.props.store.dispatch({ type : actionTypes.SET_MESHSYNC_SUBSCRIPTION, meshSyncState : data });
    }

    const operatorSubscription = new GQLSubscription({ type : OPERATOR_EVENT_SUBSCRIPTION, contextIds : contexts, callbackFunction : operatorCallback })
    const meshSyncSubscription = new GQLSubscription({ type : MESHSYNC_EVENT_SUBSCRIPTION, contextIds : contexts, callbackFunction : meshSyncCallback })

    this.setState({ operatorSubscription, meshSyncSubscription });
  }

  handleDrawerToggle = () => {
    this.setState(state => ({ mobileOpen : !state.mobileOpen }));
  }

  handleL5CommunityClick = () => {
    if (typeof window !== 'undefined') {
      const w = window.open('https://layer5.io', '_blank');
      w.focus();
    }
  };

  /**
   * Sets the selected k8s context on global level.
   * @param {Array.<string>} activeK8sContexts
   */
  activeContextChangeCallback = (activeK8sContexts)  => {
    if (activeK8sContexts.includes("all")) {
      activeK8sContexts = ["all"];
    }
    this.props.store.dispatch({ type : actionTypes.SET_K8S_CONTEXT, selectedK8sContexts : activeK8sContexts });
  }

  setActiveContexts = (id) => {
    if (this.state.k8sContexts?.contexts) {
      if (id === "all") {
        let activeContexts = [];
        this.state.k8sContexts.contexts.forEach(ctx =>
          activeContexts.push(ctx.id)
        );
        activeContexts.push("all");
        this.setState(state => {
          if (state.activeK8sContexts?.includes("all")) return { activeK8sContexts : [] };
          return { activeK8sContexts : activeContexts };
        },
        () => this.activeContextChangeCallback(this.state.activeK8sContexts));
        return;
      }

      this.setState(state => {
        let ids = [...(state.activeK8sContexts || [])];
        //pop event
        if (ids.includes(id)) {
          ids  = ids.filter(id => id != "all")
          return { activeK8sContexts : ids.filter(cid => cid !== id) }
        }

        //push event
        if (ids.length === this.state.k8sContexts.contexts.length - 1) {
          ids.push("all");
        }
        return { activeK8sContexts : [...ids, id] }
      }, () => this.activeContextChangeCallback(this.state.activeK8sContexts))
    }
  }

  searchContexts = (search = "") => {
    fetchContexts(10, search)
      .then(ctx => {
        this.setState({ k8sContexts : ctx })
        const active = ctx?.contexts?.find(c => c.is_current_context === true);
        if (active) this.setState({ activeK8sContexts : [active?.id] })
      })
      .catch(err => console.error(err))
  }

  async loadConfigFromServer() {
    const { store } = this.props;
    dataFetch('/api/system/sync', {
      credentials : 'same-origin',
      method : 'GET',
      credentials : 'include',
    }, result => {
      if (result) {
        if (result.k8sConfig && result.k8sConfig.length != 0) {
          const kubeConfigs = result.k8sConfig.map(config => Object.assign({
            inClusterConfig : false,
            k8sfile : "",
            name : "",
            clusterConfigured : "",
            server : "",
            created_at : "",
            updated_at : "",
            ts : new Date()
          }, config));
          store.dispatch({ type : actionTypes.UPDATE_CLUSTER_CONFIG, k8sConfig : kubeConfigs });
        }
        if (result.meshAdapters && result.meshAdapters !== null && result.meshAdapters.length > 0) {
          store.dispatch({ type : actionTypes.UPDATE_ADAPTERS_INFO, meshAdapters : result.meshAdapters });
        }
        if (result.grafana) {
          const grafanaCfg = Object.assign({
            grafanaURL : "",
            grafanaAPIKey : "",
            grafanaBoardSearch : "",
            grafanaBoards : [],
            selectedBoardsConfigs : []
          }, result.grafana)
          store.dispatch({ type : actionTypes.UPDATE_GRAFANA_CONFIG, grafana : grafanaCfg });
        }
        if (result.prometheus) {
          if (typeof result.prometheus.prometheusURL === 'undefined') {
            result.prometheus.prometheusURL = '';
          }
          if (typeof result.prometheus.selectedPrometheusBoardsConfigs === 'undefined') {
            result.prometheus.selectedPrometheusBoardsConfigs = [];
          }
          const promCfg = Object.assign({
            prometheusURL : "",
            selectedPrometheusBoardsConfigs : []
          }, result.prometheus)
          store.dispatch({ type : actionTypes.UPDATE_PROMETHEUS_CONFIG, prometheus : promCfg });
        }
        if (result.loadTestPrefs) {
          const loadTestPref = Object.assign({
            c : 0,
            qps : 0,
            t : 0,
            gen : 0
          }, result.loadTestPrefs)
          store.dispatch({ type : actionTypes.UPDATE_LOAD_GEN_CONFIG, loadTestPref });
        }
        if (typeof result.anonymousUsageStats !== 'undefined') {
          store.dispatch({ type : actionTypes.UPDATE_ANONYMOUS_USAGE_STATS, anonymousUsageStats : result.anonymousUsageStats });
        }
        if (typeof result.anonymousPerfResults !== 'undefined') {
          store.dispatch({ type : actionTypes.UPDATE_ANONYMOUS_PERFORMANCE_RESULTS, anonymousPerfResults : result.anonymousPerfResults });
        }
      }
    }, error => {
      console.log(`there was an error fetching user config data: ${error}`);
    });
  }

  static async getInitialProps({ Component, ctx }) {
    const pageProps = Component.getInitialProps
      ? await Component.getInitialProps(ctx)
      : {}
    return { pageProps };
  }

  render() {
    const {
      Component, pageProps, classes, isDrawerCollapsed
    } = this.props;
    return (
      <NoSsr>
        <div className={classes.root}>
          <CssBaseline />
          <nav className={isDrawerCollapsed
            ? classes.drawerCollapsed
            : classes.drawer} data-test="navigation">
            <Hidden smUp implementation="js">
              <Navigator
                variant="temporary"
                open={this.state.mobileOpen}
                onClose={this.handleDrawerToggle}
                onCollapseDrawer={(open = null) => this.handleCollapseDrawer(open)}
                isDrawerCollapsed={isDrawerCollapsed}
              />
            </Hidden>
            <Hidden xsDown implementation="css">
              <Navigator
                onCollapseDrawer={(open = null) => this.handleCollapseDrawer(open)}
                isDrawerCollapsed={isDrawerCollapsed} />
            </Hidden>
          </nav>
          <div className={classes.appContent}>
            <SnackbarProvider
              anchorOrigin={{
                vertical : 'bottom',
                horizontal : 'right',
              }}
              iconVariant={{
                success : <CheckCircle style={{ marginRight : "0.5rem" }} />,
                error : <Error style={{ marginRight : "0.5rem" }} />,
                warning : <Warning style={{ marginRight : "0.5rem" }} />,
                info : <Info style={{ marginRight : "0.5rem" }} />
              }}
              classes={{
                variantSuccess : classes.notifSuccess,
                variantError : classes.notifError,
                variantWarning : classes.notifWarn,
                variantInfo : classes.notifInfo,
              }}
              maxSnack={10}
            >
              <MesheryProgressBar />
              <Header
                onDrawerToggle={this.handleDrawerToggle}
                onDrawerCollapse={isDrawerCollapsed}
                contexts={this.state.k8sContexts}
                activeContexts={this.state.activeK8sContexts}
                setActiveContexts={this.setActiveContexts}
                searchContexts={this.searchContexts}
              />
              <main className={classes.mainContent}>
                <MuiPickersUtilsProvider utils={MomentUtils}>
                  <Component
                    pageContext={this.pageContext}
                    contexts={this.state.k8sContexts}
                    activeContexts={this.state.activeK8sContexts}
                    setActiveContexts={this.setActiveContexts}
                    searchContexts={this.searchContexts}
                    {...pageProps}
                  />
                </MuiPickersUtilsProvider>
              </main>
            </SnackbarProvider>
            <footer className={classes.footer}>
              <Typography variant="body2" align="center" color="textSecondary" component="p">
                <span onClick={this.handleL5CommunityClick} className={classes.footerText}>
                  Built with <FavoriteIcon className={classes.footerIcon} /> by the Layer5 Community</span>
              </Typography>
            </footer>
          </div>
        </div>
      </NoSsr>
    );
  }
}

MesheryApp.propTypes = { classes : PropTypes.object.isRequired, };

const mapStateToProps = state => ({
  isDrawerCollapsed : state.get("isDrawerCollapsed"),
  k8sConfig : state.get("k8sConfig"),
  operatorSubscription : state.get("operatorSubscription"),
  meshSyncSubscription : state.get("meshSyncSubscription")
})

const mapDispatchToProps = dispatch => ({
  toggleCatalogContent : bindActionCreators(toggleCatalogContent, dispatch)
})

const MesheryWithRedux = connect(mapStateToProps, mapDispatchToProps)(MesheryApp);

const MesheryAppWrapper = (props) => {
  return (
    <Provider store={props.store}>
      <Head>
        <link rel="shortcut icon" href="/static/img/meshery-logo/meshery-logo.svg" />
        <title>Meshery</title>
      </Head>
      <MuiThemeProvider theme={theme}>
        <MuiPickersUtilsProvider utils={MomentUtils}>
          <MesheryWithRedux {...props} />
        </MuiPickersUtilsProvider>
      </MuiThemeProvider>
    </Provider>
  );
}

export default withStyles(styles)(withRedux(makeStore, {
  serializeState : state => state.toJS(),
  deserializeState : state => fromJS(state)
})(MesheryAppWrapper));


