/* eslint-disable */
import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import {
  NoSsr, Card, CardHeader, CardContent, Grid, IconButton, Tooltip,
} from '@material-ui/core';

import Dialog from '@material-ui/core/Dialog';
import MuiDialogTitle from '@material-ui/core/DialogTitle';
import MuiDialogContent from '@material-ui/core/DialogContent';
import MuiDialogActions from '@material-ui/core/DialogActions';
import TextField from '@material-ui/core/TextField';
import { connect } from 'react-redux';
import dataFetch from '../lib/data-fetch';
import PlayCircleFilledIcon from '@material-ui/icons/PlayCircleFilled';
import TimerIcon from '@material-ui/icons/Timer';
import CloseIcon from '@material-ui/icons/Close';
import DeleteIcon from '@material-ui/icons/Delete';
import VisibilityIcon from '@material-ui/icons/Visibility';
import { bindActionCreators } from 'redux';

const styles = (theme) => ({
  root: {
    padding: theme.spacing(10),
  },
  expansionPanel: {
    boxShadow:'none',
    border: '1px solid rgb(196,196,196)',
  },
  margin: {
    margin: theme.spacing(1),
  },
  button: {
    float: 'right',
    marginRight: '10px'
  },
  profileButtons: {
    padding: '15px'
  },
  closeButton: {
    position: 'absolute',
    right: theme.spacing(1),
    top: theme.spacing(1),
    color: theme.palette.grey[500],
  },
});

const DialogTitle = withStyles(styles)(props => {
  const { children, classes, onClose, ...other } = props;
  return (
    <MuiDialogTitle disableTypography className={classes.modalHeading} {...other}>
      <Typography variant="h6">{children}</Typography>
      {onClose ? (
        <IconButton aria-label="close" className={classes.closeButton} onClick={onClose}>
          <CloseIcon />
        </IconButton>
      ) : null}
    </MuiDialogTitle>
  );
});

const DialogContent = withStyles(theme => ({
  root: {
    padding: theme.spacing(2),
  },
}))(MuiDialogContent);

const DialogActions = withStyles(theme => ({
  root: {
    padding: theme.spacing(1),
  },
}))(MuiDialogActions);

class MesheryPerfProfileComponent extends React.Component {
  constructor(props) {
    super(props);
    this.state={
      modalOpen : false,
      updatedProfile : true,
    }
    /*this.state = {
      testName,
      meshName,
      url,
      qps,
      c,
      t,
      loadGenerator: 'fortio',
      result,
      headers: "",
      cookies: "",
      reqBody: "",
      contentType: "",

      timerDialogOpen: false,
      blockRunTest: false,
      urlError: false,
      tError: false,

      testUUID: this.generateUUID(),
      staticPrometheusBoardConfig,
      selectedMesh: '',
      availableAdapters: [],
    };*/
  }

  /*handleChange = (name) => (event) => {
    if (name === 'url' && event.target.value !== '') {
      this.setState({ urlError: false });
    }
    if (name === 't' && (event.target.value.toLowerCase().endsWith('h')
      || event.target.value.toLowerCase().endsWith('m') || event.target.value.toLowerCase().endsWith('s'))) {
      this.setState({ tError: false });
    }

    this.setState({ [name]: event.target.value });
  };


  handleSubmit = () => {

    const {
      url, t
    } = this.state;

    if (url === '') {
      this.setState({ urlError: true });
      return;
    }

    let err = false;
    let tNum = 0;
    try {
      tNum = parseInt(t.substring(0, t.length - 1));
    } catch (ex) {
      err = true;
    }

    if (t === '' || !(t.toLowerCase().endsWith('h')
      || t.toLowerCase().endsWith('m') || t.toLowerCase().endsWith('s')) || err || tNum <= 0) {
      this.setState({ tError: true });
      return;
    }

    this.submitLoadTest();
  }
  */
  handleSubmit = () => {
    const data = {
      rps:'10',
      t:'10s',
      c:'10',
      gen:'fortio',
      protocol:'TCP',
      headers:{"h1":"v1"},
      cookies: {"h1":"v1"},
      reqBody: 'reqBody',
      contentType: 'type',
      endpoint:"https://google.com",
      labels:{"h1":"v1"}
    };

    const params = Object.keys(data).map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`).join('&');
    dataFetch('/api/user/test-prefs', {
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': 'meshery-provider=None',
      },
      body: params,
    }, (result) => {
      console.log(result);
      alert(result);
    }, console.error("Fetch Fail"));
    //*/
    /*const data1 = {
      qps:'10',
      c:'10',
      t:'10s',
      gen:'fortio',
    };
    const params1 = Object.keys(data1).map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data1[key])}`).join('&');
    console.log(params1);
    dataFetch('/api/perf/load-test-prefs', {
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
      },
      body: params1,
    }, (result) => {
      console.log(result);
    }, console.error("Fetch Fail"));*/
  }
 /*
  handleSuccess() {
    const self = this;
    return (result) => {
      const {
        testName, meshName, url, qps, c, t, loadGenerator,
      } = this.state;
      if (typeof result !== 'undefined' && typeof result.runner_results !== 'undefined') {
        self.props.enqueueSnackbar('Successfully fetched the data.', {
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
        self.props.updateLoadTestData({
          loadTest: {
            testName,
            meshName,
            url,
            qps,
            c,
            t,
            loadGenerator,
            result,
          },
        });
        self.setState({ result, testUUID: self.generateUUID() });
      }
      self.closeEventStream();
      self.setState({ blockRunTest: false, timerDialogOpen: false });
    };
  }
  */
  async startEventStream(url) {
    this.closeEventStream();
    this.eventStream = new EventSource(url);
    this.eventStream.onmessage = this.handleEvents();
    this.eventStream.onerror = this.handleError('Connection to the server got disconnected. Load test might be running in the background. Please check the results page in a few.');
    this.props.enqueueSnackbar('Load test has been successfully submitted', {
      variant: 'info',
      autoHideDuration: 1000,
      action: (key) => (
        <IconButton
          key="close"
          aria-label="Close"
          color="inherit"
          onClick={() => this.props.closeSnackbar(key)}
        >
          <CloseIcon />
        </IconButton>
      ),
    });
  }
  /*
  handleEvents() {
    const self = this;
    let track = 0;
    return (e) => {
      const data = JSON.parse(e.data);
      switch (data.status) {
        case 'info':
          self.props.enqueueSnackbar(data.message, {
            variant: 'info',
            autoHideDuration: 1000,
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
          if (track === 0) {
            self.setState({ timerDialogOpen: true, result: {} });
            track++;
          }
          break;
        case 'error':
          self.handleError('Load test did not run successfully with msg')(data.message);
          break;
        case 'success':
          self.handleSuccess()(data.result);
          break;
      }
    };
  }

  closeEventStream() {
    if (this.eventStream && this.eventStream.close) {
      this.eventStream.close();
      this.eventStream = null;
    }
  }

  componentDidMount() {
    this.getStaticPrometheusBoardConfig();
    this.scanForMeshes();
    this.getLoadTestPrefs();
  }

  getLoadTestPrefs = () => {
    const self = this;
    dataFetch('/api/perf/load-test-prefs', {
      credentials: 'same-origin',
      method: 'GET',
      credentials: 'include',
    }, (result) => {
      if (typeof result !== 'undefined') {
        this.setState({
          qps: result.loadTestPrefs.qps,
          c: result.loadTestPrefs.c,
          t: result.loadTestPrefs.t,
          loadGenerator: result.loadTestPrefs.gen,
        });
      }
    }, () => {}); //error is already captured from the handler, also we have a redux-store for same & hence it's not needed here.

  }

  getStaticPrometheusBoardConfig = () => {
    const self = this;
    if ((self.props.staticPrometheusBoardConfig && self.props.staticPrometheusBoardConfig !== null && Object.keys(self.props.staticPrometheusBoardConfig).length > 0)
      || (self.state.staticPrometheusBoardConfig && self.state.staticPrometheusBoardConfig !== null && Object.keys(self.state.staticPrometheusBoardConfig).length > 0)) {
      return;
    }
    dataFetch('/api/prometheus/static_board', {
      credentials: 'same-origin',
      credentials: 'include',
    }, (result) => {
      if (typeof result !== 'undefined' && typeof result.cluster !== 'undefined' && typeof result.node !== 'undefined'
        && typeof result.cluster.panels !== 'undefined' && result.cluster.panels.length > 0
        && typeof result.node.panels !== 'undefined' && result.node.panels.length > 0) {
        self.props.updateStaticPrometheusBoardConfig({
          staticPrometheusBoardConfig: result, // will contain both the cluster and node keys for the respective boards
        });
        self.setState({ staticPrometheusBoardConfig: result });
      }
    }, self.handleError('unable to fetch pre-configured boards'));
  }

  scanForMeshes = () => {
    const self = this;

    if (typeof self.props.k8sConfig === 'undefined' || !self.props.k8sConfig.clusterConfigured) {
      return;
    }
    dataFetch('/api/mesh/scan', {
      credentials: 'same-origin',
      credentials: 'include',
    }, (result) => {
      if (typeof result !== 'undefined' && Object.keys(result).length > 0) {
        const adaptersList = [];
        Object.keys(result).forEach((mesh) => {
          adaptersList.push(mesh);
        });
        self.setState({ availableAdapters: adaptersList });
        Object.keys(result).forEach((mesh) => {
          self.setState({ selectedMesh: mesh });
        });
      }
    // }, self.handleError("unable to scan the kubernetes cluster"));
    }, () => {});
  }

  generateUUID() {
    const { v4: uuid } = require('uuid');
    return uuid();
  }

  handleError(msg) {
    const self = this;
    return (error) => {
      self.setState({ blockRunTest: false, timerDialogOpen: false });
      self.closeEventStream();
      let finalMsg = msg;
      if (typeof error === 'string') {
        finalMsg = `${msg}: ${error}`;
      }
      self.props.enqueueSnackbar(finalMsg, {
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
        autoHideDuration: 4000,
      });
    };
  }

  handleTimerDialogClose = () => {
    this.setState({ timerDialogOpen: false });
  }*/

  handleChange = (val) => {
    const self = this;
    return () => {
      self.setState({
        updatedProfile : false,
      });
    };
  }

  handleModalOpen() {
    const self = this;
    return () => {
      self.setState({ modalOpen : true });
    };
  }

  handleModalClose() {
    const self=this;
    return () => {
      self.setState({ modalOpen : false });
    };
  }

  render() {
    const { classes } = this.props;
    /*const {
      timerDialogOpen, blockRunTest, url, qps, c, t, loadGenerator, testName, meshName, result, urlError,
      tError, testUUID, selectedMesh, availableAdapters, headers, cookies, reqBody, contentType
    } = this.state;
    let staticPrometheusBoardConfig;
    if (this.props.staticPrometheusBoardConfig && this.props.staticPrometheusBoardConfig != null && Object.keys(this.props.staticPrometheusBoardConfig).length > 0) {
      staticPrometheusBoardConfig = this.props.staticPrometheusBoardConfig;
    } else {
      staticPrometheusBoardConfig = this.state.staticPrometheusBoardConfig;
    }
    let chartStyle = {};
    if (timerDialogOpen) {
      chartStyle = { opacity: 0.3 };
    }
    let displayStaticCharts = '';
    let displayGCharts = '';
    let displayPromCharts = '';

    availableAdapters.forEach((item) => {
      const index = meshes.indexOf(item);
      if (index !== -1) meshes.splice(index, 1);
    });

    if (staticPrometheusBoardConfig && staticPrometheusBoardConfig !== null && Object.keys(staticPrometheusBoardConfig).length > 0 && prometheus.prometheusURL !== '') {
      // only add testUUID to the board that should be persisted
      if (staticPrometheusBoardConfig.cluster) {
        staticPrometheusBoardConfig.cluster.testUUID = testUUID;
      }
      displayStaticCharts = (
        <React.Fragment>
          <Typography variant="h6" gutterBottom className={classes.chartTitle}>
            Node Metrics
          </Typography>
          <GrafanaCustomCharts
            boardPanelConfigs={[staticPrometheusBoardConfig.cluster, staticPrometheusBoardConfig.node]}
            prometheusURL={prometheus.prometheusURL}
          />
        </React.Fragment>
      );
    }
    if (prometheus.selectedPrometheusBoardsConfigs.length > 0) {
      displayPromCharts = (
        <React.Fragment>
          <Typography variant="h6" gutterBottom cclassName={classes.chartTitleGraf}>
            Prometheus charts
          </Typography>
          <GrafanaCustomCharts
            boardPanelConfigs={prometheus.selectedPrometheusBoardsConfigs}
            prometheusURL={prometheus.prometheusURL}
          />
        </React.Fragment>
      );
    }
    if (grafana.selectedBoardsConfigs.length > 0) {
      displayGCharts = (
        <React.Fragment>
          <Typography variant="h6" gutterBottom className={classes.chartTitleGraf}>
            Grafana charts
          </Typography>
          <GrafanaCustomCharts
            boardPanelConfigs={grafana.selectedBoardsConfigs}
            grafanaURL={grafana.grafanaURL}
            grafanaAPIKey={grafana.grafanaAPIKey}
          />
        </React.Fragment>
      );
    }*/
    const self = this;
    const { modalOpen, updatedProfile } = this.state;
    const profiles = [
      {
        name: 'Test1',
        id: '1001'
      },
      {
        name: 'Test2',
        id: '1002'
      },
      {
        name: 'Test3',
        id: '1002'
      },
    ]
    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.root}>
          <Grid container spacing={3}>
            {profiles.map( ({ name,id }, index) => {
              return (
                <>
                  <Grid item xs={12} md={6} key={index} gutterBottom>
                    <Card>
                      <CardHeader 
                        title={name}
                        subheader={id}
                        action={
                          <Tooltip title="View or Edit Profile">
                          <IconButton onClick={self.handleModalOpen()}>
                            <VisibilityIcon />
                          </IconButton>
                          </Tooltip>
                        }
                      />
                      <CardContent className={classes.profileButtons}>
                        <Tooltip title="Delete Profile">
                          <IconButton>
                            <DeleteIcon />
                          </IconButton>
                        </Tooltip>
                        <Button
                            type="submit"
                            variant="contained"
                            color="primary"
                            size="medium"
                            className={classes.button}
                            onClick={this.handleSubmit}
                            startIcon={<PlayCircleFilledIcon />}
                        >
                          Run Test
                        </Button>
                        <Button
                            type="submit"
                            variant="contained"
                            color="primary"
                            size="medium"
                            className={classes.button}
                            startIcon={<TimerIcon />}
                        >
                          Schedule Test
                        </Button>
                      </CardContent>
                    </Card>
                  </Grid>
                  <Dialog onClose={self.handleModalClose()} aria-labelledby="customized-dialog-title" open={modalOpen} disableScrollLock={true}>
                  <DialogTitle id="customized-dialogs-title" onClose={self.handleModalClose()}>
                    <b>Profile Details</b>
                  </DialogTitle>
                  <DialogContent dividers>
                    <Grid container spacing={1}>
                    <Grid item xs={12} md={6}>
                      <Tooltip title="If a test name is not provided, a random one will be generated for you.">
                        <TextField
                          id="testName"
                          name="testName"
                          label="Test Name"
                          fullWidth
                          value={name}
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange('testName')}
                          inputProps={{ maxLength: 300 }}
                        />
                      </Tooltip>
                    </Grid>
                    <Grid item xs={12} md={6}>
                      <TextField
                        select
                        id="meshName"
                        name="meshName"
                        label="Service Mesh"
                        fullWidth
                        value="Service Mesh"
                        margin="normal"
                        variant="outlined"
                        onChange={this.handleChange('meshName')}
                      >
                      </TextField>
                    </Grid>
                    <Grid item xs={12}>
                      <TextField
                        required
                        id="url"
                        name="url"
                        label="URL to test"
                        type="url"
                        fullWidth
                        value="URL"
                        margin="normal"
                        variant="outlined"
                        onChange={this.handleChange('url')}
                      />
                    </Grid>
                    <Grid item xs={12} md={4}>
                      <TextField
                        required
                        id="c"
                        name="c"
                        label="Concurrent requests"
                        type="number"
                        fullWidth
                        value="Concurrent Req"
                        inputProps={{ min: '0', step: '1' }}
                        margin="normal"
                        variant="outlined"
                        onChange={this.handleChange('c')}
                      />
                    </Grid>
                    <Grid item xs={12} md={4}>
                      <TextField
                        required
                        id="qps"
                        name="qps"
                        label="Queries per second"
                        type="number"
                        fullWidth
                        value="QPS"
                        inputProps={{ min: '0', step: '1' }}
                        margin="normal"
                        variant="outlined"
                        onChange={this.handleChange('qps')}
                      />
                    </Grid>
                    <Grid item xs={12} md={4}>
                      <Tooltip title={"Please use 'h', 'm' or 's' suffix for hour, minute or second respectively."}>
                        <TextField
                          required
                          id="t"
                          name="t"
                          label="Duration"
                          fullWidth
                          value="Duration"
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange('t')}
                        />
                      </Tooltip>
                    </Grid>
                    <Grid item xs={12} md={12} gutterBottom>
                          <Grid container spacing={1}>
                            <Grid item xs={12}>
                              <TextField
                                id="headers"
                                name="headers"
                                label="Request Headers"
                                fullWidth
                                value="Req Headers"
                                multiline
                                margin="normal"
                                variant="outlined"
                                onChange={this.handleChange('headers')}
                              >
                              </TextField>
                            </Grid>
                            <Grid item xs={12}>
                              <TextField
                                id="cookies"
                                name="cookies"
                                label="Request Cookies"
                                fullWidth
                                value=""
                                multiline
                                margin="normal"
                                variant="outlined"
                                onChange={this.handleChange('cookies')}
                              >
                              </TextField>
                            </Grid>
                            <Grid item xs={12}>
                              <TextField
                                id="contentType"
                                name="contentType"
                                label="Content Type"
                                fullWidth
                                value="Content Type"
                                multiline
                                margin="normal"
                                variant="outlined"
                                onChange={this.handleChange('contentType')}
                              >
                              </TextField>
                            </Grid>
                            <Grid item xs={12} md={12}>
                              <TextField
                                id="cookies"
                                name="cookies"
                                label="Request Body"
                                fullWidth
                                value="Req Body"
                                multiline
                                margin="normal"
                                variant="outlined"
                                onChange={this.handleChange('reqBody')}
                              >
                              </TextField>
                            </Grid>
                          </Grid>
                    </Grid>
                  </Grid>
                </DialogContent>
                <DialogActions>
                  <Button color="primary" disabled={updatedProfile}>
                    Update Profile
                  </Button>
                  <Button autoFocus onClick={self.handleModalClose()} color="primary">
                    OK
                  </Button>
                </DialogActions>
              </Dialog>
              </>
            )
            })}
            </Grid>
          </div>
        </React.Fragment>
      </NoSsr>
    );
  }
}

MesheryPerfProfileComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

/*const mapDispatchToProps = (dispatch) => ({
  updateLoadTestData: bindActionCreators(updateLoadTestData, dispatch),
  updateStaticPrometheusBoardConfig: bindActionCreators(updateStaticPrometheusBoardConfig, dispatch),
  updateLoadTestPref: bindActionCreators(updateLoadTestPref, dispatch),
});
const mapStateToProps = (state) => {
  const loadTest = state.get('loadTest').toJS();
  const grafana = state.get('grafana').toJS();
  const prometheus = state.get('prometheus').toJS();
  const k8sConfig = state.get('k8sConfig').toJS();
  const staticPrometheusBoardConfig = state.get('staticPrometheusBoardConfig').toJS();
  const loadTestPref = state.get('loadTestPref').toJS();
  return {
    ...loadTest, grafana, prometheus, staticPrometheusBoardConfig, k8sConfig,
  };
};*/


export default withStyles(styles)(connect(
  //mapStateToProps,
  //mapDispatchToProps,
  null,
  null
)(MesheryPerfProfileComponent));
