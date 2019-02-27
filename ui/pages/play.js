import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import { NoSsr, IconButton, FormGroup, FormControl, InputLabel, Input } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import LoadTestTimerDialog from '../components/load-test-timer-dialog';
import MesheryChart from '../components/MesheryChart';
import Snackbar from '@material-ui/core/Snackbar';
import MesherySnackbarWrapper from '../components/MesherySnackbarWrapper';
import dataFetch from '../lib/data-fetch';
import Switch from '@material-ui/core/Switch';
import blue from '@material-ui/core/colors/blue';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import { updateK8SConfig } from '../lib/store';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import MesheryPlayComponent from '../components/MesheryPlayComponent';


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
  uploadButton: {
    margin: theme.spacing(1),
    marginTop: theme.spacing(3),
  },
  rightIcon: {
    marginLeft: theme.spacing(1),
  },
});

class Play extends React.Component {

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
      this.setState({
        // inClusterConfig: false,
        // k8sfile: '', 
        // k8sfileError: false,
        // contextName: '', 
        // meshLocationURL: '', 
        // meshLocationURLError: false,
        reconfigureCluster: true,
      })
  }


  render() {
    const { classes } = this.props;
      return (
    <NoSsr>
    <React.Fragment>
        <MesheryPlayComponent />
    </React.Fragment>
    </NoSsr>
  );
      }
}

// Play.propTypes = {
//   classes: PropTypes.object.isRequired,
// };

const mapDispatchToProps = dispatch => {
    return {
        // updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch)
    }
}
const mapStateToProps = state => {
    // console.log("header - mapping state to props. . . new title: "+ state.get("page").get("title"));
    // console.log("state: " + JSON.stringify(state));
    // const k8sconfig = state.get("k8sConfig");
    let newprops = {};
    // if (typeof k8sconfig !== 'undefined'){
    //   newprops = { 
    //     inClusterConfig: k8sconfig.get('inClusterConfig'),
    //     // k8sfile: '', 
    //     contextName: k8sconfig.get('contextName'), 
    //     meshLocationURL: k8sconfig.get('meshLocationURL'), 

    //     reconfigureCluster: k8sconfig.get('reconfigureCluster'),
    //   }
    // }
    return newprops;
}

export default withStyles(styles)(connect(
    mapStateToProps,
    mapDispatchToProps
  )(Play));
