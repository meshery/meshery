import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import { NoSsr, Tooltip, MenuItem, IconButton, CircularProgress, FormControl, FormLabel, RadioGroup, FormControlLabel, Radio } from '@material-ui/core';
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
import GrafanaCustomCharts from './GrafanaCustomCharts';

let uuid;
if (typeof window !== 'undefined') { 
  uuid = require('uuid/v4');
}


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

class MesherySettingsPerformanceComponent extends React.Component {
  constructor(props){
    super(props);
    const {testName, meshName, url, qps, c, t, result, staticPrometheusBoardConfig, k8sConfig} = props;

    this.state = {
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
    };
  }

  handleChange = name => event => {
   
    if ((event.target.value.toLowerCase().endsWith('h') || 
      event.target.value.toLowerCase().endsWith('m') || event.target.value.toLowerCase().endsWith('s'))){
      this.setState({tError: false});
    }
    this.setState({ [name]: event.target.value });
  };

  handleSubmit = () => {

    const {testName, meshName, url, qps, c, t, loadGenerator, testUUID} = this.state;
  

    if (t === '' || !(t.toLowerCase().endsWith('h') || 
      t.toLowerCase().endsWith('m') || t.toLowerCase().endsWith('s'))){
      this.setState({tError: true})
      return;
    }

    this.submitLoadTest();
  }

  submitLoadTest = () => {
    const {testName, meshName, url, qps, c, t, loadGenerator, testUUID} = this.state;

   
    const t1 = t.substring(0, t.length);
    const dur = t.substring(t.length - 1, t.length).toLowerCase();

    const data = {
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
   
    
    fetch(`/api/load-test-prefs?${params}`, {
    method: 'POST',
    body: data
})
  
  }

  componentDidMount() {
    this.scanForMeshes();
  }

  scanForMeshes = () => {
    const self = this;
    const {selectedMesh} = this.state;
    if (typeof self.props.k8sConfig === 'undefined' || !self.props.k8sConfig.clusterConfigured){
      return;
    }
    dataFetch('/api/mesh/scan', { 
      credentials: 'same-origin',
      credentials: 'include',
    }, result => {
      if (typeof result !== 'undefined' && Object.keys(result).length > 0){
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

  handleTimerDialogClose = () => {
    this.setState({timerDialogOpen: false});
  }

  render() {
    const { classes, grafana, prometheus } = this.props;
    const { timerDialogOpen, blockRunTest, qps, url, testName, testNameError, meshName, t, c, result, loadGenerator, 
        urlError, tError, testUUID, selectedMesh } = this.state;
    
    return (
      <NoSsr>
      <React.Fragment>
      <div className={classes.root}>
      <Grid container spacing={1}>
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
        
      
      </div>
    </React.Fragment>

      </NoSsr>
    );
  }
}

MesherySettingsPerformanceComponent.propTypes = {
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
 
  const grafana = state.get("grafana").toJS();
  const prometheus = state.get("prometheus").toJS();
  const k8sConfig = state.get("k8sConfig").toJS();
  const staticPrometheusBoardConfig = state.get("staticPrometheusBoardConfig").toJS();
  return {...loadTest, grafana, prometheus, staticPrometheusBoardConfig, k8sConfig};
}


export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps
)(withSnackbar(MesherySettingsPerformanceComponent)));
