import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import { NoSsr, Tooltip, MenuItem, IconButton } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import LoadTestTimerDialog from '../components/load-test-timer-dialog';
import MesheryChart from '../components/MesheryChart';
import { withSnackbar } from 'notistack';
import dataFetch from '../lib/data-fetch';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import { updateLoadTestData } from '../lib/store';
import GrafanaCharts from './GrafanaCharts';
import CloseIcon from '@material-ui/icons/Close';


const meshes = [
  'AspenMesh',
  'Consul Connect',
  'Grey Matter',
  'Istio',
  'Kong',
  'Linkerd 1.x',
  'Linkerd 2.x',
  'Mesher',
  'Rotor',
  'SOFAMesh',
  'Zuul',
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
  chartContent: {
    minHeight: window.innerHeight * 0.7,
  },
});

class MesheryPerformanceComponent extends React.Component {
  constructor(props){
    super(props);
    const {testName, meshName, url, qps, c, t, result} = props;

    this.state = {
      testName, 
      meshName, 
      url,
      qps,
      c,
      t,
      result,

      timerDialogOpen: false,
      urlError: false,
      tError: false,
      testNameError: false,
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
    this.setState({timerDialogOpen: true});
  }

  submitLoadTest = () => {
    const {testName, meshName, url, qps, c, t} = this.state;

    let computedTestName = testName;
    if (testName.trim() === '') {
      const mesh = meshName === ''?'No mesh': meshName;
      computedTestName = `${mesh}_${(new Date()).getTime()}`;
    }

    const t1 = t.substring(0, t.length - 1);
    const dur = t.substring(t.length - 1, t.length).toLowerCase();

    const data = {
      name: computedTestName, 
      mesh: meshName, 
      url,
      qps,
      c,
      t: t1, 
      dur,
    };
    const params = Object.keys(data).map((key) => {
      return encodeURIComponent(key) + '=' + encodeURIComponent(data[key]);
    }).join('&');
    // console.log(`data to be submitted for load test: ${params}`);
    let self = this;
    dataFetch('/api/load-test', { 
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
      },
      body: params
    }, result => {
      if (typeof result !== 'undefined' && typeof result.runner_results !== 'undefined'){
        this.setState({result, timerDialogOpen: false});
        this.props.enqueueSnackbar('Successfully fetched the data.', {
          variant: 'success',
          autoHideDuration: 4000,
        });
        this.props.updateLoadTestData({loadTest: {
          testName,
          meshName,
          url,
          qps,
          c,
          t, 
          result,
        }});
      }
    }, self.handleError);
  }

  handleError = error => {
    this.setState({timerDialogOpen: false });
    const self = this;
    this.props.enqueueSnackbar(`Load test did not run successfully with msg: ${error}`, {
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

  handleTimerDialogClose = () => {
    this.setState({timerDialogOpen: false});
  }

  render() {
    const { classes, grafana } = this.props;
    const { timerDialogOpen, qps, url, testName, testNameError, meshName, t, c, result, urlError, tError } = this.state;

    let chartStyle = {}
    if (timerDialogOpen) {
      chartStyle = {opacity: .3};
    }

    let displayGCharts = '';
    if (grafana.selectedBoardsConfigs.length > 0) {
      displayGCharts = (
        <React.Fragment>
          <Typography variant="h6" gutterBottom className={classes.chartTitle}>
            Grafana charts
          </Typography>
        <GrafanaCharts 
          boardPanelConfigs={grafana.selectedBoardsConfigs} 
          grafanaURL={grafana.grafanaURL} />
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
              value={meshName}
              margin="normal"
              variant="outlined"
              onChange={this.handleChange('meshName')}
          >
                <MenuItem key={'mh_-_none'} value={''}>None</MenuItem>
              {meshes && meshes.map((mesh) => (
                  <MenuItem key={'mh_-_'+mesh} value={mesh}>{mesh}</MenuItem>
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
            disabled={timerDialogOpen}
          >
           Run Test
          </Button>
        </div>
      </React.Fragment>

      <LoadTestTimerDialog open={timerDialogOpen} 
      t={t}
      onClose={this.handleTimerDialogClose} 
      countDownComplete={this.handleTimerDialogClose}
       />

      <Typography variant="h6" gutterBottom className={classes.chartTitle} id="timerAnchor">
        Test Results
      </Typography>
        <div className={classes.chartContent} style={chartStyle}>
          <MesheryChart data={[result && result.runner_results?result.runner_results:{}]} />    
        </div>
        
      
      </div>
    </React.Fragment>
    
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
    updateLoadTestData: bindActionCreators(updateLoadTestData, dispatch)
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
  return {...loadTest, grafana};
}


export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps
)(withSnackbar(MesheryPerformanceComponent)));
