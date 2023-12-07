import React from 'react';
import { NoSsr, Grid } from '@material-ui/core';

import Popup from '../Popup';
import { withStyles } from '@material-ui/core/styles';
import { connect } from 'react-redux';
import { withRouter } from 'next/router';
import { withNotify } from '../../utils/hooks/useNotification';
import { bindActionCreators } from 'redux';
import {
  updateGrafanaConfig,
  updateProgress,
  updatePrometheusConfig,
  updateTelemetryUrls,
} from '../../lib/store.js';
import blue from '@material-ui/core/colors/blue';

import DashboardMeshModelGraph from './charts/DashboardMeshModelGraph.js';
import ConnectionStatsChart from './charts/ConnectionCharts.js';
import MesheryConfigurationChart from './charts/MesheryConfigurationCharts.js';

const styles = (theme) => ({
  rootClass: { backgroundColor: theme.palette.secondary.elevatedComponents2, marginTop: '1rem' },
  datatable: {
    boxShadow: 'none',
  },
  chip: {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  link: {
    cursor: 'pointer',
    textDecoration: 'none',
  },
  metricsButton: { width: '240px' },
  alreadyConfigured: { textAlign: 'center' },
  margin: { margin: theme.spacing(1) },
  colorSwitchBase: {
    color: blue[300],
    '&$colorChecked': {
      color: blue[500],
      '& + $colorBar': { backgroundColor: blue[500] },
    },
  },
  colorBar: {},
  colorChecked: {},
  fileLabel: { width: '100%' },
  fileLabelText: {},
  inClusterLabel: { paddingRight: theme.spacing(2) },
  alignCenter: { textAlign: 'center' },
  icon: { width: theme.spacing(2.5) },
  istioIcon: { width: theme.spacing(1.5) },
  settingsIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  addIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  cardHeader: { fontSize: theme.spacing(2) },
  card: {
    height: '100%',
    marginTop: theme.spacing(2),
  },
  cardContent: { height: '100%' },
  redirectButton: {
    marginLeft: '-.5em',
    color: '#000',
  },
  dashboardSection: {
    backgroundColor: theme.palette.secondary.elevatedComponents,
    padding: theme.spacing(2),
    borderRadius: 4,
    height: '100%',
    marginBottom: theme.spacing(2),
  },
});

const Overview = (props) => {
  const { classes } = props;
  return (
    <NoSsr>
      <Popup />
      <div className={classes.rootClass}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={12}>
            <DashboardMeshModelGraph classes={classes} />
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={12} md={6}>
                <ConnectionStatsChart classes={classes} />
              </Grid>
              <Grid item xs={12} md={6}>
                <MesheryConfigurationChart classes={classes} />
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </div>
    </NoSsr>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  updateGrafanaConfig: bindActionCreators(updateGrafanaConfig, dispatch),
  updatePrometheusConfig: bindActionCreators(updatePrometheusConfig, dispatch),
  updateTelemetryUrls: bindActionCreators(updateTelemetryUrls, dispatch),
});

const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const meshAdapters = state.get('meshAdapters');
  const meshAdaptersts = state.get('meshAdaptersts');
  const grafana = state.get('grafana').toJS();
  const prometheus = state.get('prometheus').toJS();
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    meshAdapters,
    meshAdaptersts,
    k8sconfig,
    grafana,
    prometheus,
    selectedK8sContexts,
  };
};

export default withStyles(styles, { withTheme: true })(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withNotify(Overview))),
);
