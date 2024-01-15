import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Typography } from '@material-ui/core';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import dataFetch from '../../../lib/data-fetch';
import PrometheusSelectionComponent from './PrometheusSelectionComponent';
import GrafanaDisplaySelection from '../grafana/GrafanaDisplaySelection';
import { updateGrafanaConfig, updateProgress, updatePrometheusConfig } from '../../../lib/store';
import GrafanaCustomCharts from '../grafana/GrafanaCustomCharts';
import PrometheusConfigComponent from './PrometheusConfigComponent';
import { getK8sClusterIdsFromCtxId } from '../../../utils/multi-ctx';
import fetchAvailableAddons from '../../graphql/queries/AddonsStatusQuery';
import { withNotify } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import { CONNECTION_KINDS } from '@/utils/Enum';
import { withTelemetryHook } from '@/components/hooks/useTelemetryHook';

const promStyles = (theme) => ({
  buttons: {
    display: 'flex',
    //   justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    //   marginLeft: theme.spacing(1),
  },
  margin: { margin: theme.spacing(1) },
  icon: { width: theme.spacing(2.5) },
  alignRight: { textAlign: 'right' },
  formControl: { margin: theme.spacing(1), minWidth: 180 },
  panelChips: { display: 'flex', flexWrap: 'wrap' },
  panelChip: { margin: theme.spacing(0.25) },
  chartTitle: { marginLeft: theme.spacing(3), marginTop: theme.spacing(2), textAlign: 'center' },
});

// todo
export const submitPrometheusConfigure = (self, cb = () => {}) => {
  const { prometheusURL, selectedPrometheusBoardsConfigs } = self.state;
  if (prometheusURL === '') {
    return;
  }
  const data = { prometheusURL };
  const params = Object.keys(data)
    .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
    .join('&');
  self.props.updateProgress({ showProgress: true });
  dataFetch(
    `/api/telemetry/metrics/config`,
    {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
      body: params,
    },
    (result) => {
      self.props.updateProgress({ showProgress: false });
      if (typeof result !== 'undefined') {
        const notify = self.props.notify;
        notify({ message: 'Prometheus was configured!', event_type: EVENT_TYPES.SUCCESS });
        self.setState({ prometheusConfigSuccess: true });
        self.props.updatePrometheusConfig({
          prometheus: { prometheusURL, selectedPrometheusBoardsConfigs },
        });
        cb();
      }
    },
    self.handleError,
  );
};

class PrometheusComponent extends Component {
  constructor(props) {
    super(props);
    // const {selectedPrometheusBoardsConfigs} = props.grafana;
    const { prometheusURL, selectedPrometheusBoardsConfigs, connectionID, connectionName } =
      props.prometheus;
    let prometheusConfigSuccess = false;
    if (prometheusURL !== '') {
      prometheusConfigSuccess = true;
    }

    this.state = {
      urlError: false,

      prometheusConfigSuccess,
      selectedPrometheusBoardsConfigs: selectedPrometheusBoardsConfigs
        ? selectedPrometheusBoardsConfigs
        : [],
      prometheusURL,
      ts: new Date(),
      connectionID,
      connectionName,
    };
  }

  componentDidMount() {
    const self = this;
    if (self.props.isMeshConfigured)
      // todo
      dataFetch(
        `/api/telemetry/metrics/config`,
        {
          method: 'GET',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
          },
        },
        (result) => {
          self.props.updateProgress({ showProgress: false });
          if (
            typeof result !== 'undefined' &&
            result?.prometheusURL &&
            result?.prometheusURL != ''
          ) {
            let selector = {
              type: 'ALL_MESH',
              k8sClusterIDs: this.getK8sClusterIds(),
            };
            fetchAvailableAddons(selector).subscribe({
              next: (res) => {
                res?.addonsState?.forEach((addon) => {
                  if (addon.name === 'prometheus' && self.state.prometheusURL === '') {
                    self.setState({ prometheusURL: 'http://' + addon.endpoint });
                    submitPrometheusConfigure(self);
                  }
                });
              },
              error: (err) => console.log('error registering Prometheus: ' + err),
            });
          }
        },
        self.handleError,
      );
  }

  static getDerivedStateFromProps(props, state) {
    const { prometheusURL, selectedPrometheusBoardsConfigs } = props.prometheus;
    let configs = selectedPrometheusBoardsConfigs ? [] : selectedPrometheusBoardsConfigs;
    if (prometheusURL !== state.prometheusURL && prometheusURL !== '') {
      return {
        prometheusURL,
        configs,
        prometheusConfigSuccess: prometheusURL !== '',
        ts: props.ts,
      };
    }
    return {};
  }
  getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(this.props.selectedK8sContexts, this.props.k8sconfig);
  };

  handleChange = (name) => (data) => {
    if (name === 'prometheusURL' && !!data) {
      this.setState({ urlError: false });
    }
    const promCfg = {
      prometheusURL: data?.value || '',
      selectedPrometheusBoardsConfigs: data?.metadata['prometheus_boards'] || [],
      connectionID: data?.id,
      connectionName: data?.name,
    };
    this.props.updatePrometheusConfig({ prometheus: promCfg });
    // this.setState({ [name]: value });
  };

  handlePrometheusConfigure = () => {
    const { prometheusURL } = this.state;
    if (
      prometheusURL === '' ||
      !(
        prometheusURL.toLowerCase().startsWith('http://') ||
        prometheusURL.toLowerCase().startsWith('https://')
      )
    ) {
      this.setState({ urlError: true });
      return;
    }
    submitPrometheusConfigure(this);
  };

  handleError = () => {
    this.props.updateProgress({ showProgress: false });
    const notify = this.props.notify;
    notify({
      message: 'There was an error communicating with Prometheus',
      event_type: EVENT_TYPES.ERROR,
    });
  };

  // done
  handlePrometheusChipDelete = (e) => {
    e.preventDefault();
    const self = this;
    self.props.updateProgress({ showProgress: true });
    dataFetch(
      `/api/integrations/connections/${CONNECTION_KINDS.PROMETHEUS}/status`,
      {
        method: 'PUT',
        credentials: 'include',
        body: JSON.stringify({ id: self?.state?.connectionID }),
      },
      (result) => {
        self.props.updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          self.setState({
            prometheusConfigSuccess: false,
            prometheusURL: '',
            selectedPrometheusBoardsConfigs: [],
          });
          self.props.updatePrometheusConfig({
            prometheus: { prometheusURL: '', selectedPrometheusBoardsConfigs: [] },
          });
        }
      },
      self.handleError,
    );
  };

  // done
  handlePrometheusClick = () => {
    const self = this;
    this.props.ping(self.state.connectionName, self.state.prometheusURL, self.state.connectionID);
  };

  addSelectedBoardPanelConfig = (boardsSelection) => {
    const { prometheusURL, selectedPrometheusBoardsConfigs } = this.state;
    if (boardsSelection && boardsSelection.panels && boardsSelection.panels.length) {
      selectedPrometheusBoardsConfigs.push(boardsSelection);
      this.setState({ selectedPrometheusBoardsConfigs });
      this.props.updatePrometheusConfig({
        prometheus: { prometheusURL, selectedPrometheusBoardsConfigs },
      });
    }
  };

  deleteSelectedBoardPanelConfig = (indexes) => {
    const { prometheusURL, selectedPrometheusBoardsConfigs } = this.state;
    indexes.sort();
    for (let i = indexes.length - 1; i >= 0; i--) {
      selectedPrometheusBoardsConfigs.splice(indexes[i], 1);
    }
    this.setState({ selectedPrometheusBoardsConfigs });
    this.props.updatePrometheusConfig({
      prometheus: { prometheusURL, selectedPrometheusBoardsConfigs },
    });
  };

  render() {
    const { classes } = this.props;
    const {
      urlError,
      prometheusURL,
      prometheusConfigSuccess,
      selectedPrometheusBoardsConfigs,
      connectionID,
    } = this.state;
    if (prometheusConfigSuccess) {
      let displaySelec = '';
      if (selectedPrometheusBoardsConfigs.length > 0) {
        displaySelec = (
          <React.Fragment>
            <GrafanaDisplaySelection
              boardPanelConfigs={selectedPrometheusBoardsConfigs}
              deleteSelectedBoardPanelConfig={this.deleteSelectedBoardPanelConfig}
            />

            <Typography variant="h6" gutterBottom className={classes.chartTitle}>
              Prometheus charts
            </Typography>
            {/* <GrafanaCharts
                  boardPanelConfigs={selectedPrometheusBoardsConfigs}
                  prometheusURL={prometheusURL} /> */}
            <GrafanaCustomCharts
              boardPanelConfigs={selectedPrometheusBoardsConfigs}
              prometheusURL={prometheusURL}
              connectionID={connectionID}
            />
          </React.Fragment>
        );
      }

      return (
        <NoSsr>
          <React.Fragment>
            <PrometheusSelectionComponent
              prometheusURL={prometheusURL}
              handlePrometheusChipDelete={this.handlePrometheusChipDelete}
              addSelectedBoardPanelConfig={this.addSelectedBoardPanelConfig}
              handlePrometheusClick={this.handlePrometheusClick}
              handleError={this.handleError}
              connectionID={connectionID}
            />
            {displaySelec}
          </React.Fragment>
        </NoSsr>
      );
    }
    return (
      <NoSsr>
        <PrometheusConfigComponent
          prometheusURL={prometheusURL && { label: prometheusURL, value: prometheusURL }}
          urlError={urlError}
          handleChange={this.handleChange}
          handlePrometheusConfigure={this.handlePrometheusConfigure}
        />
      </NoSsr>
    );
  }
}

PrometheusComponent.propTypes = {
  classes: PropTypes.object.isRequired,
  scannedPrometheus: PropTypes.array.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  updateGrafanaConfig: bindActionCreators(updateGrafanaConfig, dispatch),
  updatePrometheusConfig: bindActionCreators(updatePrometheusConfig, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (st) => {
  const grafana = st.get('grafana').toJS();
  const prometheus = st.get('prometheus').toJS();
  const selectedK8sContexts = st.get('selectedK8sContexts');
  const k8sconfig = st.get('k8sConfig');

  return { grafana, prometheus, selectedK8sContexts, k8sconfig };
};

export default withStyles(promStyles)(
  connect(
    mapStateToProps,
    mapDispatchToProps,
  )(withTelemetryHook(withNotify(PrometheusComponent), CONNECTION_KINDS.PROMETHEUS)),
);
