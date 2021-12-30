import React from 'react';
import App from 'next/app';
import Head from 'next/head';
import { MuiThemeProvider,  withStyles } from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';
import getPageContext from '../components/PageContext';
import Navigator from '../components/Navigator';
import Header from '../components/Header';
import PropTypes from 'prop-types';
import Hidden from '@material-ui/core/Hidden';
import withRedux from "next-redux-wrapper";
import { makeStore, actionTypes } from '../lib/store';
import { Provider } from "react-redux";
import { fromJS } from 'immutable';
import { NoSsr, Typography } from '@material-ui/core';
import FavoriteIcon from '@material-ui/icons/Favorite';
import { SnackbarProvider } from 'notistack';
import { MuiPickersUtilsProvider } from '@material-ui/pickers';
import MomentUtils from '@date-io/moment';
import {
  CheckCircle,
  Info,
  Error,
  Warning
} from '@material-ui/icons';

// codemirror + js-yaml imports when added to a page was preventing to navigating to that page using nextjs
// link clicks, hence attempting to add them here
import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/material.css';
import 'codemirror/addon/lint/lint.css';
import './../public/static/style/index.css';

// import 'billboard.js/dist/theme/insight.min.css';
// import 'billboard.js/dist/theme/graph.min.css';
import 'billboard.js/dist/billboard.min.css';

import MesheryProgressBar from '../components/MesheryProgressBar';
import dataFetch from '../lib/data-fetch';
import theme, { styles } from "../themes"

if (typeof window !== 'undefined') {
  require('codemirror/mode/yaml/yaml');
  require('codemirror/mode/javascript/javascript');
  require('codemirror/addon/lint/lint');
  require('codemirror/addon/lint/yaml-lint');
  require('codemirror/addon/lint/json-lint');
  if (typeof window.jsyaml === 'undefined'){
    window.jsyaml = require('js-yaml');
  }
  if (typeof window.jsonlint === 'undefined'){
    // jsonlint did not work well with codemirror json-lint. Hence, found an alternative (jsonlint-mod) based on https://github.com/scniro/react-codemirror2/issues/21
    window.jsonlint = require('jsonlint-mod');
  }
}

class MesheryApp extends App {
  constructor() {
    super();
    this.pageContext = getPageContext();

    this.state = { mobileOpen : false,
      isDrawerCollapsed : false };
  }

  handleDrawerToggle = () => {
    this.setState(state => ({ mobileOpen : !state.mobileOpen }));
  }

  handleCollapseDrawer = (open = null) => {
    if (typeof(open) === 'boolean')
      return this.setState({ isDrawerCollapsed : open });
    this.setState(state => ({ isDrawerCollapsed : !state.isDrawerCollapsed }));
  };

  handleL5CommunityClick = () => {
    if (typeof window !== 'undefined'){
      const w = window.open('https://layer5.io', '_blank');
      w.focus();
    }
  };

  async loadConfigFromServer() {
    const { store } = this.props;
    dataFetch('/api/system/sync', { credentials : 'same-origin',
      method : 'GET',
      credentials : 'include', }, result => {
      if (typeof result !== 'undefined'){
        if (result.k8sConfig){
          if (typeof result.k8sConfig.inClusterConfig === 'undefined'){
            result.k8sConfig.inClusterConfig = false;
          }
          if (typeof result.k8sConfig.k8sfile === 'undefined'){
            result.k8sConfig.k8sfile = '';
          }
          if (typeof result.k8sConfig.contextName === 'undefined'){
            result.k8sConfig.contextName = '';
          }
          if (typeof result.k8sConfig.clusterConfigured === 'undefined'){
            result.k8sConfig.clusterConfigured = false;
          }
          if (typeof result.k8sConfig.configuredServer === 'undefined'){
            result.k8sConfig.configuredServer = '';
          }
          store.dispatch({ type : actionTypes.UPDATE_CLUSTER_CONFIG, k8sConfig : result.k8sConfig });
        }
        if (result.meshAdapters && result.meshAdapters !== null && result.meshAdapters.length > 0) {
          store.dispatch({ type : actionTypes.UPDATE_ADAPTERS_INFO, meshAdapters : result.meshAdapters });
        }
        if (result.grafana){
          if (typeof result.grafana.grafanaURL === 'undefined'){
            result.grafana.grafanaURL = '';
          }
          if (typeof result.grafana.grafanaAPIKey === 'undefined'){
            result.grafana.grafanaAPIKey = '';
          }
          if (typeof result.grafana.grafanaBoardSearch === 'undefined'){
            result.grafana.grafanaBoardSearch = '';
          }
          if (typeof result.grafana.grafanaBoards === 'undefined'){
            result.grafana.grafanaBoards = [];
          }
          if (typeof result.grafana.selectedBoardsConfigs === 'undefined'){
            result.grafana.selectedBoardsConfigs = [];
          }
          store.dispatch({ type : actionTypes.UPDATE_GRAFANA_CONFIG, grafana : result.grafana });
        }
        if (result.prometheus){
          if (typeof result.prometheus.prometheusURL === 'undefined'){
            result.prometheus.prometheusURL = '';
          }
          if (typeof result.prometheus.selectedPrometheusBoardsConfigs === 'undefined'){
            result.prometheus.selectedPrometheusBoardsConfigs = [];
          }
          store.dispatch({ type : actionTypes.UPDATE_PROMETHEUS_CONFIG, prometheus : result.prometheus });
        }
        if (result.loadTestPrefs){
          if (typeof result.loadTestPrefs.c === 'undefined'){
            result.loadTestPrefs.c = 0;
          }
          if (typeof result.loadTestPrefs.qps === 'undefined'){
            result.loadTestPrefs.qps = 0;
          }
          if (typeof result.loadTestPrefs.t === 'undefined'){
            result.loadTestPrefs.t = '30s';
          }
          if (typeof result.loadTestPrefs.gen === 'undefined'){
            result.loadTestPrefs.gen = '';
          }
          store.dispatch({ type : actionTypes.UPDATE_LOAD_GEN_CONFIG, loadTestPref : result.loadTestPrefs });
        }
        if (typeof result.anonymousUsageStats !== 'undefined'){
          store.dispatch({ type : actionTypes.UPDATE_ANONYMOUS_USAGE_STATS, anonymousUsageStats : result.anonymousUsageStats });
        }
        if (typeof result.anonymousPerfResults !== 'undefined'){
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

  componentDidMount(){
    this.loadConfigFromServer(); // this works, but sometimes other components which need data load faster than this data is obtained.
  }

  render() {
    const {
      Component, store, pageProps, classes
    } = this.props;
    const { isDrawerCollapsed } = this.state;
    return (
      <NoSsr>
        <Provider store={store}>
          <Head>
            <link rel="shortcut icon" href="/ui/public/static/img/meshery-logo/meshery-logo.svg" />
            <title>Meshery</title>
          </Head>
          <MuiThemeProvider theme={theme}>
            <MuiPickersUtilsProvider utils={MomentUtils}>
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
                    anchorOrigin={{ vertical : 'bottom',
                      horizontal : 'right', }}
                    iconVariant = {{
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
                    <Header onDrawerToggle={this.handleDrawerToggle} onDrawerCollapse={isDrawerCollapsed}/>
                    <main className={classes.mainContent}>
                      <MuiPickersUtilsProvider utils={MomentUtils}>
                        <Component pageContext={this.pageContext} {...pageProps} />
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
            </MuiPickersUtilsProvider>
          </MuiThemeProvider>
        </Provider>
      </NoSsr>
    );
  }
}

MesheryApp.propTypes = { classes : PropTypes.object.isRequired, };

export default withStyles(styles)(withRedux(makeStore, { serializeState : state => state.toJS(),
  deserializeState : state => fromJS(state) })(MesheryApp));
