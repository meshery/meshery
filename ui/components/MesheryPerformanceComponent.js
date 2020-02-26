import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import { NoSsr, Tooltip, MenuItem, IconButton, CircularProgress, FormControl, FormLabel, RadioGroup, FormControlLabel, Radio, Divider } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import LoadTestTimerDialog from '../components/load-test-timer-dialog';
import MesheryChart from '../components/MesheryChart';
import { withSnackbar } from 'notistack';
import dataFetch from '../lib/data-fetch';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import { updateLoadTestData, updateStaticPrometheusBoardConfig } from '../lib/store';
// import GrafanaCharts from './GrafanaCharts';
import CloseIcon from '@material-ui/icons/Close';
import GetAppIcon from '@material-ui/icons/GetApp';
import GrafanaCustomCharts from './GrafanaCustomCharts';

let uuid;
if (typeof window !== 'undefined') { 
  uuid = require('uuid/v4');
}


const meshes = [
  'Istio',
  'Linkerd',
  'App Mesh',
  'Aspen Mesh',
  'Citrix Service Mesh',
  'Consul Connect',
  'Grey Matter',
  'Kong',
  'Mesher',
  'Network Service Mesh',
  'Octarine',
  'Rotor',
  'SOFAMesh',
  'Zuul',
]

const loadGenerators = [
  'fortio',
  'wrk2'
]

const styles = theme => ({
  root: {
    padding: theme.spacing(10),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  margin: {
    margin: theme.spacing(1),
  },
  chartTitle: {
    textAlign: 'center',
  },
  chartTitleGraf: {
    textAlign: 'center',
    // marginTop: theme.spacing(5),
  },
  chartContent: {
    // minHeight: window.innerHeight * 0.7,
  },
  centerTimer: {
    width: '100%',
  }
});

class MesheryPerformanceComponent extends React.Component {
  constructor(props){
    super(props);
    const {testName, meshName, url, qps, c, t, result, staticPrometheusBoardConfig, k8sConfig} = props;

    this.state = {
      testName, 
      meshName, 
      url,
      qps,
      c,
      t,
      loadGenerator : 'fortio',
      result,

      timerDialogOpen: false,
      blockRunTest: false,
      urlError: false,
      tError: false,
      testNameError: false,

      testUUID: this.generateUUID(),
      staticPrometheusBoardConfig,
      selectedMesh: '',
      availableAdapters: [],
    };
  }

  handleChange = name => event => {
    if (name === 'url' && event.target.value !== ''){
      this.setState({urlError: false});
    }
    if (name === 'testName`' && event.target.value !== ''){
      this.setState({testNameError: false});
    }
    if (name === 't' && (event.target.value.toLowerCase().endsWith('h') || 
      event.target.value.toLowerCase().endsWith('m') || event.target.value.toLowerCase().endsWith('s'))){
      this.setState({tError: false});
    }
    this.setState({ [name]: event.target.value });
  };

  handleSubmit = () => {

    const { url, t, testName, meshName } = this.state;
    if (url === ''){
      this.setState({urlError: true})
      return;
    }

    // if (testName === ''){
    //   this.setState({testNameError: true})
    //   return;
    // }

    let err = false, tNum = 0;
    try {
      tNum = parseInt(t.substring(0, t.length - 1))
    }catch(ex){
      err = true;
    }

    if (t === '' || !(t.toLowerCase().endsWith('h') || 
      t.toLowerCase().endsWith('m') || t.toLowerCase().endsWith('s')) || err || tNum <= 0){
      this.setState({tError: true})
      return;
    }

    this.submitLoadTest();
  }

  submitLoadTest = () => {
    const {testName, meshName, url, qps, c, t, loadGenerator, testUUID} = this.state;

    let computedTestName = testName;
    if (testName.trim() === '') {
      const mesh = meshName === '' || meshName === 'None'?'No mesh': meshName;
      computedTestName = `${mesh}_${(new Date()).getTime()}`;
    }

    const t1 = t.substring(0, t.length - 1);
    const dur = t.substring(t.length - 1, t.length).toLowerCase();

    const data = {
      name: computedTestName, 
      mesh: meshName === '' || meshName === 'None'?'': meshName, // to prevent None from getting to the DB
      url,
      qps,
      c,
      t: t1, 
      dur,
      uuid: testUUID,
      loadGenerator,
    };
    const params = Object.keys(data).map((key) => {
      return encodeURIComponent(key) + '=' + encodeURIComponent(data[key]);
    }).join('&');
    this.startEventStream(`/api/load-test?${params}`);
    this.setState({blockRunTest: true}); // to block the button
  }

  handleSuccess() {
    const self = this;
    return (result) => {
      const {testName, meshName, url, qps, c, t, loadGenerator, testUUID} = this.state;
      if (typeof result !== 'undefined' && typeof result.runner_results !== 'undefined'){
        self.props.enqueueSnackbar('Successfully fetched the data.', {
          variant: 'success',
          autoHideDuration: 2000,
          action: (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key) }
            >
              <CloseIcon />
            </IconButton>
          ),
        });
        self.props.updateLoadTestData({loadTest: {
          testName,
          meshName,
          url,
          qps,
          c,
          t, 
          loadGenerator,
          result,
        }});
        self.setState({result, testUUID: self.generateUUID()});
      }
      self.closeEventStream();
      self.setState({blockRunTest: false, timerDialogOpen: false});
    }
  }

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
          onClick={() => this.props.closeSnackbar(key) }
        >
          <CloseIcon />
        </IconButton>
      ),
    });
  }

  handleEvents(){
    const self = this;
    let track = 0;
    return e => {
      const data = JSON.parse(e.data);
      switch(data.status){
      case 'info':
        self.props.enqueueSnackbar(data.message, {
          variant: 'info',
          autoHideDuration: 1000,
          action: (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key) }
            >
              <CloseIcon />
            </IconButton>
          ),
        });
        if (track === 0){
          self.setState({timerDialogOpen: true, result: {}});
          track++;
        }
        break;
      case 'error':
        self.handleError("Load test did not run successfully with msg")(data.message);
        break;
      case 'success':
        self.handleSuccess()(data.result);
        break;
      }
    }
  }

  closeEventStream() {
    if(this.eventStream && this.eventStream.close){
      this.eventStream.close();
      this.eventStream = null;
    }
  }

  componentDidMount() {
    this.getStaticPrometheusBoardConfig();
    this.scanForMeshes();
  }

  getStaticPrometheusBoardConfig = () => {
    let self = this;
    if ((self.props.staticPrometheusBoardConfig && self.props.staticPrometheusBoardConfig !== null && Object.keys(self.props.staticPrometheusBoardConfig).length > 0) || 
      (self.state.staticPrometheusBoardConfig && self.state.staticPrometheusBoardConfig !==null && Object.keys(self.state.staticPrometheusBoardConfig).length > 0)) {
      return;
    }
    dataFetch('/api/prometheus/static_board', { 
      credentials: 'same-origin',
      credentials: 'include',
    }, result => {
      if (typeof result !== 'undefined' && typeof result.cluster !== 'undefined' && typeof result.node !== 'undefined' && 
        typeof result.cluster.panels !== 'undefined' && result.cluster.panels.length > 0 && 
        typeof result.node.panels !== 'undefined' && result.node.panels.length > 0){
        self.props.updateStaticPrometheusBoardConfig({
          staticPrometheusBoardConfig: result, // will contain both the cluster and node keys for the respective boards
        });
        self.setState({staticPrometheusBoardConfig: result});
      }
    }, self.handleError("unable to fetch pre-configured boards"));
  }

  scanForMeshes = () => {
    const self = this;
    const {selectedMesh} = this.state;
    const {availableAdapters} = this.state;

    if (typeof self.props.k8sConfig === 'undefined' || !self.props.k8sConfig.clusterConfigured){
      return;
    }
    dataFetch('/api/mesh/scan', { 
      credentials: 'same-origin',
      credentials: 'include',
    }, result => {
      if (typeof result !== 'undefined' && Object.keys(result).length > 0){
        let adaptersList = [];
        Object.keys(result).forEach(mesh => {
          adaptersList.push(mesh);
        });
        self.setState({availableAdapters: adaptersList});
        Object.keys(result).forEach(mesh => {
          self.setState({selectedMesh: mesh});
          return;
        })
      }
    // }, self.handleError("unable to scan the kubernetes cluster"));
    }, () => {});
  }

  generateUUID(){
    return uuid();
  }

  handleError (msg){
    const self = this;
    return error => {
      self.setState({blockRunTest: false, timerDialogOpen: false});
      self.closeEventStream();
      let finalMsg = msg;
      if (typeof error === 'string'){
        finalMsg = `${msg}: ${error}`;
      }
      self.props.enqueueSnackbar(finalMsg, {
        variant: 'error',
        action: (key) => (
          <IconButton
            key="close"
            aria-label="Close"
            color="inherit"
            onClick={() => self.props.closeSnackbar(key) }
          >
            <CloseIcon />
          </IconButton>
        ),
        autoHideDuration: 8000,
      });
    }
  }

  handleTimerDialogClose = () => {
    this.setState({timerDialogOpen: false});
  }

  render() {
    const { classes, grafana, prometheus } = this.props;
    const { timerDialogOpen, blockRunTest, qps, url, testName, testNameError, meshName, t, c, result, loadGenerator, 
        urlError, tError, testUUID, selectedMesh, availableAdapters } = this.state;
    let staticPrometheusBoardConfig;
    if(this.props.staticPrometheusBoardConfig && this.props.staticPrometheusBoardConfig != null && Object.keys(this.props.staticPrometheusBoardConfig).length > 0){
      staticPrometheusBoardConfig = this.props.staticPrometheusBoardConfig;
    } else {
      staticPrometheusBoardConfig = this.state.staticPrometheusBoardConfig;
    }
    let chartStyle = {}
    if (timerDialogOpen) {
      chartStyle = {opacity: .3};
    }
    let displayStaticCharts = '';
    let displayGCharts = '';
    let displayPromCharts = '';

    availableAdapters.forEach((item) => {
      let index = meshes.indexOf(item);
      if (index !== -1) meshes.splice(index, 1);
    });

    if (staticPrometheusBoardConfig && staticPrometheusBoardConfig !== null && Object.keys(staticPrometheusBoardConfig).length > 0 && prometheus.prometheusURL !== '') {
      // only add testUUID to the board that should be persisted
      if (staticPrometheusBoardConfig.cluster) {
        staticPrometheusBoardConfig.cluster.testUUID = testUUID
      }
      displayStaticCharts = (
        <React.Fragment>
          <Typography variant="h6" gutterBottom className={classes.chartTitle}>
            Node Metrics
          </Typography>
          <GrafanaCustomCharts
            boardPanelConfigs={[staticPrometheusBoardConfig.cluster, staticPrometheusBoardConfig.node]} 
            prometheusURL={prometheus.prometheusURL} />
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
            prometheusURL={prometheus.prometheusURL} />
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
            grafanaAPIKey={grafana.grafanaAPIKey} />
        </React.Fragment>
      );
    }
    return (
      <NoSsr>
      <React.Fragment>
      <div className={classes.root}>
      <Grid container spacing={1}>
        <Grid item xs={12} sm={6}>
          <Tooltip title={"If a test name is not provided, a random one will be generated for you."}>
            <TextField
              id="testName"
              name="testName"
              label="Test Name"
              autoFocus
              fullWidth
              value={testName}
              error={testNameError}
              margin="normal"
              variant="outlined"
              onChange={this.handleChange('testName')}
              inputProps={{ maxLength: 300 }}
            />
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={6}>
          <TextField
              select
              id="meshName"
              name="meshName"
              label="Service Mesh"
              fullWidth
              value={meshName === '' && selectedMesh !== ''?selectedMesh:meshName}
              margin="normal"
              variant="outlined"
              onChange={this.handleChange('meshName')}
          >
            
              {availableAdapters && availableAdapters.map((mesh) => (
                <MenuItem key={'mh_-_'+mesh} value={mesh.toLowerCase()}>{mesh}</MenuItem>
              ))}
              {availableAdapters && (availableAdapters.length > 0) && <Divider />}
              <MenuItem key={'mh_-_none'} value={'None'}>None</MenuItem>
              {meshes && meshes.map((mesh) => (
                  <MenuItem key={'mh_-_'+mesh} value={mesh.toLowerCase()}>{mesh}</MenuItem>
              ))}
          </TextField>
        </Grid>
        <Grid item xs={12}>
          <TextField
            required
            id="url"
            name="url"
            label="URL to test"
            type="url"
            autoFocus
            fullWidth
            value={url}
            error={urlError}
            margin="normal"
            variant="outlined"
            onChange={this.handleChange('url')}
          />
        </Grid>
        <Grid item xs={12} sm={4}>
          <TextField
            required
            id="c"
            name="c"
            label="Concurrent requests"
            type="number"
            fullWidth
            value={c}
            inputProps={{ min: "0", step: "1" }}
            margin="normal"
            variant="outlined"
            onChange={this.handleChange('c')}
          />
        </Grid>
        <Grid item xs={12} sm={4}>
          <TextField
            required
            id="qps"
            name="qps"
            label="Queries per second"
            type="number"
            fullWidth
            value={qps}
            inputProps={{ min: "0", step: "1" }}
            margin="normal"
            variant="outlined"
            onChange={this.handleChange('qps')}
          />
        </Grid>
        <Grid item xs={12} sm={4}>
          <Tooltip title={"Please use 'h', 'm' or 's' suffix for hour, minute or second respectively."}>
            <TextField
              required
              id="t"
              name="t"
              label="Duration"
              fullWidth
              value={t}
              error={tError}
              margin="normal"
              variant="outlined"
              onChange={this.handleChange('t')}
            />
          </Tooltip>
        </Grid>
        <Grid item xs={12} sm={4}>
          <FormControl component="loadGenerator" className={classes.formControl}>
            <FormLabel component="loadGenerator">Load generator</FormLabel>
            <RadioGroup aria-label="loadGenerator" name="loadGenerator" value={loadGenerator} onChange={this.handleChange('loadGenerator')} row>
              {loadGenerators.map(lg => (
                <FormControlLabel value={lg} control={<Radio color="primary" />} label={lg} />
              ))}
            </RadioGroup>
          </FormControl>
        </Grid>
      </Grid>
      <React.Fragment>
        <div className={classes.buttons}>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            size="large"
            onClick={this.handleSubmit}
            className={classes.button}
            disabled={blockRunTest}
          >
           {blockRunTest?<CircularProgress size={30} />:'Run Test'}
          </Button>
        </div>
      </React.Fragment>

      <div className={classes.centerTimer}>
        <LoadTestTimerDialog open={timerDialogOpen} 
        t={t}
        onClose={this.handleTimerDialogClose} 
        countDownComplete={this.handleTimerDialogClose}
        />
       </div>

      {result && result.runner_results && 
        (<div>
          <Typography variant="h6" gutterBottom className={classes.chartTitle} id="timerAnchor">
            Test Results 
            <IconButton
              key="download"
              aria-label="download"
              color="inherit"
              // onClick={() => self.props.closeSnackbar(key) }
              href={`/api/result?id=${encodeURIComponent(result.meshery_id)}`}
            >
              <GetAppIcon />
            </IconButton>
          </Typography>
          <div className={classes.chartContent} style={chartStyle}>
            <MesheryChart data={[result && result.runner_results?result.runner_results:{}]} />    
          </div>
        </div>)
            }
        
      
          </div>
        </React.Fragment>

        {displayStaticCharts}

        {displayPromCharts}

        {displayGCharts}

      </NoSsr>
    );
  }
}

MesheryPerformanceComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    updateLoadTestData: bindActionCreators(updateLoadTestData, dispatch),
    updateStaticPrometheusBoardConfig: bindActionCreators(updateStaticPrometheusBoardConfig, dispatch),
  }
}
const mapStateToProps = state => {
  
  const loadTest = state.get("loadTest").toJS();
  // let newprops = {};
  // if (typeof loadTest !== 'undefined'){
  //   newprops = { 
  //     url: loadTest.get('url'),
  //     qps: loadTest.get('qps'), 
  //     c: loadTest.get('c'), 
  //     t: loadTest.get('t'),
  //     result: loadTest.get('result'),
  //   }
  // }
  const grafana = state.get("grafana").toJS();
  const prometheus = state.get("prometheus").toJS();
  const k8sConfig = state.get("k8sConfig").toJS();
  const staticPrometheusBoardConfig = state.get("staticPrometheusBoardConfig").toJS();
  return {...loadTest, grafana, prometheus, staticPrometheusBoardConfig, k8sConfig};
}


export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps
)(withSnackbar(MesheryPerformanceComponent)));
