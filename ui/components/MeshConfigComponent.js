import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import { NoSsr,  FormGroup, InputAdornment } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import Snackbar from '@material-ui/core/Snackbar';
import MesherySnackbarWrapper from '../components/MesherySnackbarWrapper';
import dataFetch from '../lib/data-fetch';
import Switch from '@material-ui/core/Switch';
import blue from '@material-ui/core/colors/blue';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import { updateK8SConfig, updateMeshInfo } from '../lib/store';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';


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
  alreadyConfigured: {
    textAlign: 'center',
    padding: theme.spacing(20),
  },
  colorSwitchBase: {
    color: blue[300],
    '&$colorChecked': {
      color: blue[500],
      '& + $colorBar': {
        backgroundColor: blue[500],
      },
    },
  },
  colorBar: {},
  colorChecked: {},
  // uploadButton: {
  //   // margin: theme.spacing(1),
  //   // marginTop: theme.spacing(3),
  //   float: 'right',
  // },
  // rightIcon: {
  //   // marginLeft: theme.spacing(1),
  // },
  fileLabel: {
    width: '100%',
    // paddingBottom: theme.spacing(1),
  },
  fileLabelText: {
    // width: '79%',
  },
  inClusterLabel: {
    paddingRight: theme.spacing(2),
  },
  alignCenter: {
    textAlign: 'center',
  },
  alignRight: {
    textAlign: 'right',
  }
});

class MeshConfigComponent extends React.Component {

  constructor(props) {
    super(props);
    const {inClusterConfig, contextName, meshLocationURL, reconfigureCluster} = props;
    this.state = {
        showSnackbar: false,
        snackbarVariant: '',
        snackbarMessage: '',
    
        
        inClusterConfig, // read from store
        k8sfile: '', // leaving this one out just to play it safe for now
        contextName, // read from store
        meshLocationURL, // read from store
    
        reconfigureCluster, // read from store
    
        k8sfileError: false,
        meshLocationURLError: false,
      };
  }

  handleSnackbarClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }

    this.setState({ showSnackbar: false });
  };

  handleChange = name => event => {
    if (name === 'inClusterConfig'){
        this.setState({ [name]: event.target.checked });
        return;
    }
    if (name === 'k8sfile' && event.target.value !== ''){
        this.setState({ k8sfileError: false });    
    }
    if (name === 'meshLocationURL' && event.target.value !== '') {
        this.setState({meshLocationURLError: false})
    }
    this.setState({ [name]: event.target.value });
  };

  handleSubmit = () => {
    const { inClusterConfig, k8sfile, meshLocationURL } = this.state;
    if (!inClusterConfig && k8sfile === '') {
        this.setState({k8sfileError: true});
        return;
    }
    if (meshLocationURL === ''){
        this.setState({meshLocationURLError: true})
        return;
      }

    this.submitConfig()
  }

  submitConfig = () => {
    const { inClusterConfig, k8sfile, meshLocationURL, contextName, reconfigureCluster } = this.state;
    // const data = {
    //     inClusterConfig, k8sfile, meshLocationURL, contextName
    // }
    const fileInput = document.querySelector('#k8sfile') ;
    const formData = new FormData();
    formData.append('inClusterConfig', inClusterConfig?"on":''); // to simulate form behaviour of a checkbox
    if (!inClusterConfig) {
        formData.append('contextName', contextName);
        formData.append('k8sfile', fileInput.files[0]);
    }
    formData.append('meshLocationURL', meshLocationURL);

    // console.log(`data to be submitted for load test: ${params}`);
    let self = this;
    dataFetch('/api/k8sconfig', { 
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      body: formData
    }, result => {
      if (typeof result !== 'undefined'){
        this.setState({reconfigureCluster: false, showSnackbar: true, snackbarVariant: 'success', snackbarMessage: 'Kubernetes config was successfully validated!'});
        this.props.updateK8SConfig({k8sConfig: {inClusterConfig, k8sfile, meshLocationURL, contextName, reconfigureCluster: false}});
        this.props.updateMeshInfo({mesh: result});
      }
    }, self.handleError);
  }

  handleError = error => {
    this.setState({showSnackbar: true, snackbarVariant: 'error', snackbarMessage: `Kubernetes config could not be validated: ${error}`});
  }

//   handleTimerDialogClose = () => {
//     this.setState({timerDialogOpen: false});
//   }

  handleReconfigure = () => {
    const { inClusterConfig, k8sfile, meshLocationURL, contextName, reconfigureCluster } = this.state;
      this.setState({
        // inClusterConfig: false,
        // k8sfile: '', 
        // k8sfileError: false,
        // contextName: '', 
        // meshLocationURL: '', 
        // meshLocationURLError: false,
        reconfigureCluster: true,
      })
      this.props.updateK8SConfig({k8sConfig: {inClusterConfig, k8sfile, meshLocationURL, contextName, reconfigureCluster: true}});
  }

  alreadyConfiguredTemplate = () =>{
    const { classes } = this.props;
    const {showSnackbar, 
        snackbarVariant, snackbarMessage } = this.state;
      return (
    <NoSsr>
    <React.Fragment>
    <div className={classes.root}>
    <Grid container spacing={5} alignItems="flex-end">
      <Grid item xs={12} className={classes.alignCenter}>
          <Typography variant="h5" gutterBottom>
            Mesh is configured
            </Typography>
      </Grid>
      <Grid item xs={12} sm={6} className={classes.alignRight}>
          <Button variant="contained" color="secondary" size="large" onClick={this.handleReconfigure}>
            Reconfigure
            </Button>
      </Grid>
      <Grid item xs={12} sm={6}>
          <Button variant="contained" color="primary" size="large" onClick={() => this.props.router.push('/play')}>
            Continue to Play
            </Button>
      </Grid>
      </Grid>
      </div>
{/*         
        <div className={classes.alreadyConfigured}>
            <Typography variant="h4" gutterBottom>
            Mesh is already configured
            </Typography>
            <Button variant="contained" color="secondary" onClick={this.handleReconfigure}>
            Reconfigure
            </Button>
            <Button variant="contained" color="primary" onClick={() => this.props.router.push('/play')}>
            Continue to Play
            </Button>
        </div> */}
    </React.Fragment>

    <Snackbar
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        open={showSnackbar}
        autoHideDuration={6000}
        onClose={this.handleSnackbarClose}
      >
      <MesherySnackbarWrapper 
        variant={snackbarVariant}
        message={snackbarMessage}
        onClose={this.handleSnackbarClose}
        />
    </Snackbar>
    </NoSsr>
  );
      }

  configureTemplate = () => {
    const { classes } = this.props;
    const { inClusterConfig, k8sfile, k8sfileError, contextName, meshLocationURL, meshLocationURLError, showSnackbar, 
        snackbarVariant, snackbarMessage } = this.state;
    
      return (
    <NoSsr>
    <React.Fragment>
    <div className={classes.root}>
    <Grid container spacing={5} alignItems="flex-end">
      <Grid item xs={12} className={classes.alignCenter}>
      <FormControlLabel
            key="inCluster"
            control={<FormControlLabel
                      key="notInCluster"
                      control={
                          <Switch
                              checked={inClusterConfig}
                              onChange={this.handleChange('inClusterConfig')}
                              color="default"
                              //   value="checkedA"
                              // classes={{
                              //     switchBase: classes.colorSwitchBase,
                              //     checked: classes.colorChecked,
                              //     bar: classes.colorBar,
                              // }}
                          />
                      }
                      label="Use attached config"          
                      labelPlacement="start"
                      className={classes.inClusterLabel}
                  />}
            labelPlacement="end"
            label="Use in-cluster Kubernetes config"
      />
      </Grid>
      <Grid item xs={12} sm={6}>
        <TextField
          id="contextName"
          name="contextName"
          label="Context Name"
          fullWidth
          value={contextName}
          margin="normal"
          variant="outlined"
          disabled={inClusterConfig == true}
          onChange={this.handleChange('contextName')}
        />
      </Grid>
      <Grid item xs={12} sm={6}>
      <FormGroup row>
        <input
            hidden
            className={classes.input}
            id="k8sfile"
            multiple
            type="file"
            value={k8sfile}
            onChange={this.handleChange('k8sfile')}
            disabled={inClusterConfig == true}
        />
        <label htmlFor="k8sfile" id="k8sfileLabel" className={classes.fileLabel}>
            {/* <Typography variant="body1" inline>
            {k8sfile.replace('C:\\fakepath\\', '')}
            </Typography> */}
            <TextField
                id="k8sfileLabelText"
                name="k8sfileLabelText"
                className={classes.fileLabelText}
                label="Upload config"
                variant="outlined"
                fullWidth
                value={k8sfile.replace('C:\\fakepath\\', '')}
                onClick={e => document.querySelector('#k8sfileLabel').click()}
                margin="normal"
                // InputProps={{
                //     readOnly: true,
                //   }}
                InputProps={{
                    readOnly: true,
                    startAdornment: (
                      <InputAdornment position="start">
                        <CloudUploadIcon />
                      </InputAdornment>
                    ),
                  }}
                // variant="outlined"
                disabled
                />
                {/* <Button component="span" variant="outlined" size="large" color={k8sfileError?"secondary":"primary"} disabled={inClusterConfig == true} className={classes.button, classes.uploadButton}>
                  <CloudUploadIcon className={classes.rightIcon} />
              </Button> */}
        </label>
        </FormGroup>
      </Grid>
    {/* <Grid item xs={12} sm={6}>
      <FormControl className={classes.formControl} disabled={inClusterConfig == true} fullWidth variant="outlined">
          <InputLabel htmlFor="k8sfile">Name</InputLabel>
          <input
            hidden
            className={classes.input}
            id="k8sfile"
            multiple
            type="file"
            value={k8sfile}
            onChange={this.handleChange('k8sfile')}
            disabled={inClusterConfig == true}
        />
          <Input id="k8sfile" value={k8sfile} readOnly />
        </FormControl>
    </Grid> */}
      <Grid item xs={12}>
        <TextField
          required
          id="meshLocationURL"
          name="meshLocationURL"
          label="Mesh Adapter Location"
          type="url"
          fullWidth
          value={meshLocationURL}
          error={meshLocationURLError}
          margin="normal"
          variant="outlined"
          onChange={this.handleChange('meshLocationURL')}
        />
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
        >
         Submit
        </Button>
      </div>
    </React.Fragment>
    </div>
  </React.Fragment>
  
  {/* <LoadTestTimerDialog open={timerDialogOpen} 
    t={t}
    onClose={this.handleTimerDialogClose} 
    countDownComplete={this.handleTimerDialogClose} />

  <Typography variant="h6" gutterBottom className={classes.chartTitle}>
      Results
    </Typography>
  <MesheryChart data={result} />     */}
  
  <Snackbar
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        open={showSnackbar}
        autoHideDuration={6000}
        onClose={this.handleSnackbarClose}
      >
      <MesherySnackbarWrapper 
        variant={snackbarVariant}
        message={snackbarMessage}
        onClose={this.handleSnackbarClose}
        />
    </Snackbar>
    </NoSsr>
  );
    }

  render() {
    const { reconfigureCluster } = this.state;
    if (reconfigureCluster) {
        return this.configureTemplate();
    }
    return this.alreadyConfiguredTemplate();
  }
}

MeshConfigComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
    return {
        updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch),
        updateMeshInfo: bindActionCreators(updateMeshInfo, dispatch),
    }
}
const mapStateToProps = state => {
    // console.log("header - mapping state to props. . . new title: "+ state.get("page").get("title"));
    // console.log("state: " + JSON.stringify(state));
    const k8sconfig = state.get("k8sConfig");
    let newprops = {};
    if (typeof k8sconfig !== 'undefined'){
      newprops = { 
        inClusterConfig: k8sconfig.get('inClusterConfig'),
        // k8sfile: '', 
        contextName: k8sconfig.get('contextName'), 
        meshLocationURL: k8sconfig.get('meshLocationURL'), 

        reconfigureCluster: k8sconfig.get('reconfigureCluster'),
      }
    }
    return newprops;
}

export default withStyles(styles)(connect(
    mapStateToProps,
    mapDispatchToProps
  )(withRouter(MeshConfigComponent)));
