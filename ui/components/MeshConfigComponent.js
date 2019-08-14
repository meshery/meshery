import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import { NoSsr,  FormGroup, InputAdornment, Chip, IconButton, MenuItem } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import dataFetch from '../lib/data-fetch';
import blue from '@material-ui/core/colors/blue';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import { updateK8SConfig, updateProgress } from '../lib/store';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { withSnackbar } from 'notistack';
import CloseIcon from '@material-ui/icons/Close';

const styles = theme => ({
  root: {
    padding: theme.spacing(5),
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
  fileLabel: {
    width: '100%',
  },
  fileLabelText: {
  },
  inClusterLabel: {
    paddingRight: theme.spacing(2),
  },
  alignCenter: {
    textAlign: 'center',
  },
  alignRight: {
    textAlign: 'right',
    marginBottom: theme.spacing(2),
  },
  fileInputStyle: {
    opacity: '0.01',
  },
  icon: {
    width: theme.spacing(2.5),
  },
});

class MeshConfigComponent extends React.Component {

  constructor(props) {
    super(props);
    const {inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer } = props;
    this.state = {
        inClusterConfig, // read from store
        k8sfile, // read from store
        k8sfileElementVal: '',
        contextName, // read from store
        contextNameForForm: null,
        contextsFromFile: [],
    
        clusterConfigured, // read from store
        configuredServer,
        k8sfileError: false,
        ts: new Date(),
      };
  }

  static getDerivedStateFromProps(props, state){
    const {inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer } = props;
    // if(inClusterConfig !== state.inClusterConfig || clusterConfigured !== state.clusterConfigured || k8sfile !== state.k8sfile 
        // || configuredServer !== state.configuredServer){
    if(props.ts > state.ts){
      return {
        inClusterConfig,
          k8sfile,
          k8sfileElementVal: '',
          contextName, 
          clusterConfigured,
          configuredServer,
          ts: props.ts,
      };
    }
    return {};
  }

  handleChange = name => {
    const self = this;
    return event => {
      if (name === 'inClusterConfig'){
        self.setState({ [name]: event.target.checked, ts: new Date() });
        return;
      }
      if (name === 'k8sfile'){
        if (event.target.value !== ''){
          self.setState({ k8sfileError: false });
        }
        self.setState({k8sfileElementVal: event.target.value});
        self.fetchContexts();
      }
      self.setState({ [name]: event.target.value, ts: new Date() });
    };
  }

  handleSubmit = () => {
    const { inClusterConfig, k8sfile } = this.state;
    if (!inClusterConfig && k8sfile === '') {
        this.setState({k8sfileError: true});
        return;
    }
    this.submitConfig()
  }

  fetchContexts = () => {
    const { inClusterConfig, k8sfile } = this.state;
    const fileInput = document.querySelector('#k8sfile') ;
    const formData = new FormData();
    if (inClusterConfig) {
      return;
    }
    if(fileInput.files.length == 0){
      this.setState({contextsFromFile: [], contextNameForForm: null});
      return;
    }
    // formData.append('contextName', contextName);
    formData.append('k8sfile', fileInput.files[0]);
    this.props.updateProgress({showProgress: true});
    let self = this;
    dataFetch('/api/k8sconfig/contexts', { 
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      body: formData
    }, result => {
      this.props.updateProgress({showProgress: false});
      if (typeof result !== 'undefined'){
          let ctName = '';
          result.forEach(({contextName, currentContext}) => {
            if(currentContext){
              ctName = contextName
            }
          });
          self.setState({contextsFromFile: result, contextNameForForm: ctName});
      }
    }, self.handleError);
  }

  submitConfig = () => {
    const { inClusterConfig, k8sfile, contextNameForForm } = this.state;
    const fileInput = document.querySelector('#k8sfile') ;
    const formData = new FormData();
    formData.append('inClusterConfig', inClusterConfig?"on":''); // to simulate form behaviour of a checkbox
    if (!inClusterConfig) {
        formData.append('contextName', contextNameForForm);
        formData.append('k8sfile', fileInput.files[0]);
    }
    this.props.updateProgress({showProgress: true});
    let self = this;
    dataFetch('/api/k8sconfig', { 
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      body: formData
    }, result => {
      this.props.updateProgress({showProgress: false});
      if (typeof result !== 'undefined'){
        this.setState({clusterConfigured: true, configuredServer: result.configuredServer, contextName: result.contextName});
        this.props.enqueueSnackbar('Kubernetes config was successfully validated!', {
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
        this.props.updateK8SConfig({k8sConfig: {inClusterConfig, k8sfile, contextName: result.contextName, clusterConfigured: true, configuredServer: result.configuredServer}});
      }
    }, self.handleError);
  }

  handleError = error => {
    this.props.updateProgress({showProgress: false});
    const self = this;
    this.props.enqueueSnackbar(`Kubernetes config could not be validated: ${error}`, {
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

//   handleTimerDialogClose = () => {
//     this.setState({timerDialogOpen: false});
//   }

  handleReconfigure = () => {
	let self = this;
    dataFetch('/api/k8sconfig', { 
      credentials: 'same-origin',
      method: 'DELETE',
      credentials: 'include',
    }, result => {
      this.props.updateProgress({showProgress: false});
      if (typeof result !== 'undefined'){
        this.setState({
        inClusterConfig: false,
        k8sfile: '', 
        k8sfileElementVal: '',
        k8sfileError: false,
        contextName: '', 
        contextNameForForm: null,
        clusterConfigured: false,
      })
      this.props.updateK8SConfig({k8sConfig: {inClusterConfig: false, k8sfile:'', contextName:'', contextNameForForm: null, clusterConfigured: false}});
        
      this.props.enqueueSnackbar('Kubernetes config was successfully removed!', {
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
     }
    }, self.handleError);
  }

  configureTemplate = () => {
    const { classes } = this.props;
    const { inClusterConfig, k8sfile, k8sfileElementVal, contextName, contextNameForForm, contextsFromFile, clusterConfigured, configuredServer } = this.state;
    
    let showConfigured = '';
    const self = this;
    if (clusterConfigured) {
      showConfigured = (
        <div className={classes.alignRight}>
          <Chip 
              label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
              onDelete={self.handleReconfigure} 
              icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />} 
              variant="outlined" />
        </div>
      )
    }


      return (
    <NoSsr>
    <div className={classes.root}>
    
    {showConfigured}
    
    <Grid container spacing={1} alignItems="flex-end">
      {/* <Grid item xs={12} className={classes.alignCenter}>
      <FormControlLabel
            hidden={true} // hiding this component for now
            key="inCluster"
            control={
              <Switch
                    hidden={true} // hiding this component for now
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
            labelPlacement="end"
            label="Use in-cluster Kubernetes config"
      />
      </Grid> */}
      <Grid item xs={12} sm={6}>
      <FormGroup row>
        <input
            className={classes.input}
            id="k8sfile"
            type="file"
            // value={k8sfile}
            value={k8sfileElementVal}
            onChange={this.handleChange('k8sfile')}
            disabled={inClusterConfig == true}
            className={classes.fileInputStyle}
        />
            <TextField
                id="k8sfileLabelText"
                name="k8sfileLabelText"
                className={classes.fileLabelText}
                label="Upload config"
                variant="outlined"
                fullWidth
                value={k8sfile.replace('C:\\fakepath\\', '')}
                onClick={e => document.querySelector('#k8sfile').click()}
                margin="normal"
                InputProps={{
                    readOnly: true,
                    endAdornment: (
                      <InputAdornment position="end">
                        <CloudUploadIcon />
                      </InputAdornment>
                    ),
                  }}
                disabled
                />
        </FormGroup>
      </Grid>
      <Grid item xs={12} sm={6}>
        <TextField
          select
          id="contextName"
          name="contextName"
          label="Context Name"
          fullWidth
          value={contextNameForForm}
          margin="normal"
          variant="outlined"
          disabled={inClusterConfig == true}
          onChange={this.handleChange('contextName')}
        >
          {contextsFromFile && contextsFromFile.map((ct) => (
              <MenuItem key={'ct_---_'+ct.contextName} value={ct.contextName}>{ct.contextName}{ct.currentContext?' (default)':''}</MenuItem>
          ))}
        </TextField>
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
  
  {/* <LoadTestTimerDialog open={timerDialogOpen} 
    t={t}
    onClose={this.handleTimerDialogClose} 
    countDownComplete={this.handleTimerDialogClose} />

  <Typography variant="h6" gutterBottom className={classes.chartTitle}>
      Results
    </Typography>
  <MesheryChart data={result} />     */}
    </NoSsr>
  );
    }

  render() {
    const { reconfigureCluster } = this.state;
    // if (reconfigureCluster) {
    return this.configureTemplate();
    // }
    // return this.alreadyConfiguredTemplate();
  }
}   

MeshConfigComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
    return {
        updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch),
        updateProgress: bindActionCreators(updateProgress, dispatch),
    }
}
const mapStateToProps = state => {
    const k8sconfig = state.get("k8sConfig").toJS();
    return k8sconfig;
}

export default withStyles(styles)(connect(
    mapStateToProps,
    mapDispatchToProps
  )(withRouter(withSnackbar(MeshConfigComponent))));
