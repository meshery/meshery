import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Typography, IconButton } from '@material-ui/core';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import CloseIcon from '@material-ui/icons/Close';
import { withSnackbar } from 'notistack';
import dataFetch from '../lib/data-fetch';
import PrometheusSelectionComponent from './PrometheusSelectionComponent';
import GrafanaDisplaySelection from './GrafanaDisplaySelection';
import { updateGrafanaConfig, updateProgress, updatePrometheusConfig } from '../lib/store';
import GrafanaCustomCharts from './GrafanaCustomCharts';
import PrometheusConfigComponent from './PrometheusConfigComponent';

const promStyles = (theme) => ({
  root: {
    padding: theme.spacing(5),
  },
  buttons: {
    display: 'flex',
    //   justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    //   marginLeft: theme.spacing(1),
  },
  margin: {
    margin: theme.spacing(1),
  },
  chartTitle: {
    textAlign: 'center',
  },
  icon: {
    width: theme.spacing(2.5),
  },
  alignRight: {
    textAlign: 'right',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 180,
  },
  panelChips: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  panelChip: {
    margin: theme.spacing(0.25),
  },
  chartTitle: {
    marginLeft: theme.spacing(3),
    marginTop: theme.spacing(2),
  },
});

class PrometheusComponent extends Component {
  constructor(props) {
    super(props);

    // const {selectedPrometheusBoardsConfigs} = props.grafana;
    const { prometheusURL, selectedPrometheusBoardsConfigs } = props.prometheus;
    let prometheusConfigSuccess = false;
    if (prometheusURL !== '') {
      prometheusConfigSuccess = true;
    }

    this.state = {
      urlError: false,

      prometheusConfigSuccess,
      selectedPrometheusBoardsConfigs,
      prometheusURL,
      ts: new Date(),
    };
  }

  static getDerivedStateFromProps(props, state) {
    // const {prometheusURL, selectedPrometheusBoardsConfigs} = props.grafana;
    // if(prometheusURL !== state.prometheusURL || JSON.stringify(selectedPrometheusBoardsConfigs) !== JSON.stringify(state.selectedPrometheusBoardsConfigs)) { // JSON.stringify is not the best. Will leave it for now until a better solution is found
    if (props.ts > state.ts) {
      return {
        prometheusURL, selectedPrometheusBoardsConfigs, prometheusConfigSuccess: (prometheusURL !== ''), ts: props.ts,
      };
    }
    return {};
  }

      handleChange = (name) => (value) => {
        if (name === 'prometheusURL' && value !== '') {
          this.setState({ urlError: false });
        }
        this.setState({ [name]: value });
      };

      handlePrometheusConfigure = () => {
        const { prometheusURL } = this.state;
        if (prometheusURL === '' || !(prometheusURL.toLowerCase().startsWith('http://') || prometheusURL.toLowerCase().startsWith('https://'))) {
          this.setState({ urlError: true });
          return;
        }
        this.submitPrometheusConfigure();
      }

      submitPrometheusConfigure = () => {
        const { prometheusURL, selectedPrometheusBoardsConfigs } = this.state;
        const data = {
          prometheusURL,
        };
        const params = Object.keys(data).map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`).join('&');
        this.props.updateProgress({ showProgress: true });
        const self = this;
        dataFetch('/api/prometheus/config', {
          credentials: 'same-origin',
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
          },
          body: params,
        }, (result) => {
          this.props.updateProgress({ showProgress: false });
          if (typeof result !== 'undefined') {
            this.props.enqueueSnackbar('Prometheus was successfully configured!', {
              variant: 'success',
              autoHideDuration: 2000,
              action: (key) => (
                <IconButton
                  key="close"
                  aria-label="Close"
                  color="inherit"
                  onClick={() => self.props.closeSnackbar(key)}
                >
                  <CloseIcon />
                </IconButton>
              ),
            });
            this.setState({ prometheusConfigSuccess: true });
            // this.props.updateGrafanaConfig({
            //   grafana: {
            //     prometheusURL,
            //     grafanaAPIKey,
            //     grafanaBoardSearch,
            //     grafanaBoards,
            //     selectedPrometheusBoardsConfigs,
            //   },
            // })
            this.props.updatePrometheusConfig({
              prometheus: {
                prometheusURL,
                selectedPrometheusBoardsConfigs,
              },
            });
          }
        }, self.handleError);
      }

      handleError = () => {
        const self = this;
        this.props.updateProgress({ showProgress: false });
        this.props.enqueueSnackbar('There was an error communicating with Prometheus', {
          variant: 'error',
          action: (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key)}
            >
              <CloseIcon />
            </IconButton>
          ),
          autoHideDuration: 8000,
        });
      }

      handlePrometheusChipDelete = () => {
        const self = this;
        self.props.updateProgress({ showProgress: true });
        dataFetch('/api/prometheus/config', {
          credentials: 'same-origin',
          method: 'DELETE',
          credentials: 'include',
        }, (result) => {
          self.props.updateProgress({ showProgress: false });
          if (typeof result !== 'undefined') {
            self.setState({
              prometheusConfigSuccess: false,
              prometheusURL: '',
              selectedPrometheusBoardsConfigs: [],
            });
            self.props.updatePrometheusConfig({
              prometheus: {
                prometheusURL: '',
                selectedPrometheusBoardsConfigs: [],
              },
            });
          }
        }, self.handleError);
      }

    handlePrometheusClick = () => {
      this.props.updateProgress({ showProgress: true });
      const self = this;
      dataFetch('/api/prometheus/ping', {
        credentials: 'same-origin',
        credentials: 'include',
      }, (result) => {
        self.props.updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          self.props.enqueueSnackbar('Prometheus successfully pinged!', {
            variant: 'success',
            autoHideDuration: 2000,
            action: (key) => (
              <IconButton
                key="close"
                aria-label="Close"
                color="inherit"
                onClick={() => self.props.closeSnackbar(key)}
              >
                <CloseIcon />
              </IconButton>
            ),
          });
        }
      }, self.handleError);
    }


    addSelectedBoardPanelConfig = (boardsSelection) => {
      const { prometheusURL, selectedPrometheusBoardsConfigs } = this.state;
      if (boardsSelection && boardsSelection.panels && boardsSelection.panels.length) {
        selectedPrometheusBoardsConfigs.push(boardsSelection);
        this.setState({ selectedPrometheusBoardsConfigs });
        this.props.updatePrometheusConfig({
          prometheus: {
            prometheusURL,
            selectedPrometheusBoardsConfigs,
          },
        });
      }
    }

    deleteSelectedBoardPanelConfig = (indexes) => {
      const { prometheusURL, selectedPrometheusBoardsConfigs } = this.state;
      indexes.sort();
      for (let i = indexes.length - 1; i >= 0; i--) {
        selectedPrometheusBoardsConfigs.splice(indexes[i], 1);
      }
      this.setState({ selectedPrometheusBoardsConfigs });
      this.props.updatePrometheusConfig({
        prometheus: {
          prometheusURL,
          selectedPrometheusBoardsConfigs,
        },
      });
    }


    render() {
      const { classes } = this.props;
      const {
        urlError, prometheusURL, prometheusConfigSuccess, selectedPrometheusBoardsConfigs,
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
            options={this.props.scannedPrometheus.map(url => ({label: url, value: url}))}
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
  scannedPrometheus: PropTypes.array.isRequired
};

const mapDispatchToProps = (dispatch) => ({
  updateGrafanaConfig: bindActionCreators(updateGrafanaConfig, dispatch),
  updatePrometheusConfig: bindActionCreators(updatePrometheusConfig, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});
const mapStateToProps = (st) => {
  const grafana = st.get('grafana').toJS();
  const prometheus = st.get('prometheus').toJS();
  return { grafana, prometheus };
};

export default withStyles(promStyles)(connect(
  mapStateToProps,
  mapDispatchToProps,
)(withSnackbar(PrometheusComponent)));
