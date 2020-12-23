import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import { withStyles } from '@material-ui/core/styles';
import {
  NoSsr, FormGroup, InputAdornment, Chip, IconButton, MenuItem, Tooltip,
} from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Divider from '@material-ui/core/Divider';
import blue from '@material-ui/core/colors/blue';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { withSnackbar } from 'notistack';
import CloseIcon from '@material-ui/icons/Close';
import { updateK8SConfig, updateProgress } from '../lib/store';
import dataFetch from '../lib/data-fetch';

const styles = (theme) => ({
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
  buttonsCluster: {
    display: 'flex',
    justifyContent: 'center',
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
  alignLeft: {
    textAlign: 'left',
    marginBottom: theme.spacing(2),
  },
  fileInputStyle: {
    opacity: '0.01',
  },
  icon: {
    width: theme.spacing(2.5),
  },
  configure: {
    display: 'inline-block',
    width: '48%',
    wordWrap: 'break-word',
    [theme.breakpoints.down(945)]: {
      width: '100%',
    },
  },
  vertical: {
    display: 'inline-block',
    height: 150,
    marginBottom: -60,
    [theme.breakpoints.down(945)]: {
      display: 'none',
    },
  },
  horizontal: {
    display: 'none',
    [theme.breakpoints.down(945)]: {
      display: 'block',
    },
  },
  formconfig: {
    display: 'inline-block',
    marginLeft: 30,
    [theme.breakpoints.up(946)]: {
      width: '45%',
    },
    [theme.breakpoints.down(945)]: {
      width: '100%',
      marginLeft: 0,
    },
  },
  currentConfigHeading: {
    display: 'inline-block',
    width: '48%',
    textAlign: 'center',
    [theme.breakpoints.down(945)]: {
      width: '100%',
    },
  },
  changeConfigHeading: {
    display: 'inline-block',
    width: '48%',
    textAlign: 'center',
    [theme.breakpoints.down(945)]: {
      display: 'none',
    },
  },
  changeConfigHeadingOne: {
    display: 'none',
    [theme.breakpoints.down(945)]: {
      display: 'inline-block',
      width: '100%',
      textAlign: 'center',
    },
  },
  buttonconfig: {
    display: 'inline-block',
    width: '48%',
    [theme.breakpoints.down(945)]: {
      width: '100%',
    },
  },
});

class MeshConfigComponent extends React.Component {
  constructor(props) {
    super(props);
    const {
      inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer,
    } = props;
    this.state = {
      inClusterConfig, // read from store
      inClusterConfigForm: inClusterConfig,
      k8sfile, // read from store
      k8sfileElementVal: '',
      contextName, // read from store
      contextNameForForm: '',
      contextsFromFile: [],
      clusterConfigured, // read from store
      configuredServer,
      k8sfileError: false,
      ts: new Date(),
    };
  }

  static getDerivedStateFromProps(props, state) {
    const {
      inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer,
    } = props;
    if (props.ts > state.ts) {
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

  handleChange = (name) => {
    const self = this;
    return (event) => {
      if (name === 'inClusterConfigForm') {
        self.setState({ [name]: event.target.checked, ts: new Date() });
        return;
      }
      if (name === 'k8sfile') {
        if (event.target.value !== '') {
          self.setState({ k8sfileError: false });
        }
        self.setState({ k8sfileElementVal: event.target.value });
        self.fetchContexts();
      }
      self.setState({ [name]: event.target.value, ts: new Date() });
      this.handleSubmit();
    };
  }

  handleSubmit = () => {
    const { inClusterConfigForm, k8sfile } = this.state;
    if (!inClusterConfigForm && k8sfile === '') {
      this.setState({ k8sfileError: true });
      return;
    }
    this.submitConfig();
  }

  fetchContexts = () => {
    const { inClusterConfigForm } = this.state;
    const fileInput = document.querySelector('#k8sfile');
    const formData = new FormData();
    if (inClusterConfigForm) {
      return;
    }
    if (fileInput.files.length == 0) {
      this.setState({ contextsFromFile: [], contextNameForForm: '' });
      return;
    }
    // formData.append('contextName', contextName);
    formData.append('k8sfile', fileInput.files[0]);
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch('/api/k8sconfig/contexts', {
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      body: formData,
    }, (result) => {
      this.props.updateProgress({ showProgress: false });
      if (typeof result !== 'undefined') {
        let ctName = '';
        result.forEach(({ contextName, currentContext }) => {
          if (currentContext) {
            ctName = contextName;
          }
        });
        self.setState({ contextsFromFile: result, contextNameForForm: ctName });
        self.submitConfig();
      }
    }, self.handleError);
  }

  submitConfig = () => {
    const { inClusterConfigForm, k8sfile, contextNameForForm } = this.state;
    const fileInput = document.querySelector('#k8sfile');
    const formData = new FormData();
    formData.append('inClusterConfig', inClusterConfigForm ? 'on' : ''); // to simulate form behaviour of a checkbox
    if (!inClusterConfigForm) {
      formData.append('contextName', contextNameForForm);
      formData.append('k8sfile', fileInput.files[0]);
    }
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch('/api/k8sconfig', {
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      body: formData,
    }, (result) => {
      this.props.updateProgress({ showProgress: false });
      if (typeof result !== 'undefined') {
        this.setState({ clusterConfigured: true, configuredServer: result.configuredServer, contextName: result.contextName });
        this.props.enqueueSnackbar('Kubernetes config was successfully validated!', {
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
        this.props.updateK8SConfig({
          k8sConfig: {
            inClusterConfig: inClusterConfigForm, k8sfile, contextName: result.contextName, clusterConfigured: true, configuredServer: result.configuredServer,
          },
        });
      }
    }, self.handleError);
  }

  handleKubernetesClick = () => {
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch('/api/k8sconfig/ping', {
      credentials: 'same-origin',
      credentials: 'include',
    }, (result) => {
      this.props.updateProgress({ showProgress: false });
      if (typeof result !== 'undefined') {
        this.props.enqueueSnackbar('Kubernetes was successfully pinged!', {
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

  handleError = (error) => {
    this.props.updateProgress({ showProgress: false });
    const self = this;
    this.props.enqueueSnackbar(`Kubernetes config could not be validated: ${error}`, {
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


  handleReconfigure = () => {
    const self = this;
    dataFetch('/api/k8sconfig', {
      credentials: 'same-origin',
      method: 'DELETE',
      credentials: 'include',
    }, (result) => {
      this.props.updateProgress({ showProgress: false });
      if (typeof result !== 'undefined') {
        this.setState({
          inClusterConfigForm: false,
          inClusterConfig: false,
          k8sfile: '',
          k8sfileElementVal: '',
          k8sfileError: false,
          contextName: '',
          contextNameForForm: '',
          clusterConfigured: false,
        });
        this.props.updateK8SConfig({
          k8sConfig: {
            inClusterConfig: false, k8sfile: '', contextName: '', clusterConfigured: false,
          },
        });
        this.props.enqueueSnackbar('Kubernetes config was successfully removed!', {
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

  configureTemplate = () => {
    const { classes } = this.props;
    const {
      inClusterConfig, contextName, clusterConfigured, configuredServer,
    } = this.state;
    let showConfigured = '';
    const self = this;
    if (clusterConfigured) {
      let chp = (
        <Chip
          // label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
          label={inClusterConfig ? 'Using In Cluster Config' : contextName}
          onDelete={self.handleReconfigure}
          onClick={self.handleKubernetesClick}
          icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />}
          variant="outlined"
          data-cy="chipContextName"
        />
      );
      const lst = (
        <List>
          <ListItem>
            <ListItemText primary="Context Name" secondary={inClusterConfig ? 'Using In Cluster Config' : contextName} data-cy="itemListContextName" />
          </ListItem>
          <ListItem>
            <ListItemText primary="Server Name" secondary={inClusterConfig ? 'In Cluster Server' : (configuredServer || '')} data-cy="itemListServerName" />
          </ListItem>
        </List>
      );
      if (configuredServer) {
        chp = (
          <Tooltip title={`Server: ${configuredServer}`}>
            {chp}
          </Tooltip>
        );
      }
      showConfigured = (
        <div>
          {chp}
          {lst}
        </div>
      );
    }
    if (!clusterConfigured) {
      const lst = (
        <List>
          <ListItem>
            <ListItemText primary="Context Name" secondary="Not Configured" />
          </ListItem>
          <ListItem>
            <ListItemText primary="Server Name" secondary="Not Configured" />
          </ListItem>
        </List>
      );
      showConfigured = (
        <div>
          {lst}
        </div>
      );
    }
    if (this.props.tabs == 0) {
      return this.meshOut(showConfigured);
    }
    return this.meshIn(showConfigured);
  }

  meshOut = (showConfigured) => {
    const { classes } = this.props;
    const {
      k8sfile, k8sfileElementVal, contextNameForForm, contextsFromFile,
    } = this.state;

    return (
      <NoSsr>
        <div className={classes.root}>
          <div className={classes.currentConfigHeading}>
            <h4>
              Current Configuration Details
            </h4>
          </div>
          <div className={classes.changeConfigHeading}>
            <h4>
              Change Configuration...
            </h4>
          </div>
          <div className={classes.configure}>
            {showConfigured}
          </div>
          <Divider className={classes.vertical} orientation="vertical" />
          <Divider className={classes.horizontal} orientation="horizontal" />
          <div className={classes.changeConfigHeadingOne}>
            <h4>
              Change Configuration...
            </h4>
          </div>
          <div className={classes.formconfig}>
            <FormGroup>
              <input
                id="k8sfile"
                type="file"
                value={k8sfileElementVal}
                onChange={this.handleChange('k8sfile')}
                className={classes.fileInputStyle}
              />
              <TextField
                id="k8sfileLabelText"
                name="k8sfileLabelText"
                className={classes.fileLabelText}
                label="Upload kubeconfig"
                variant="outlined"
                fullWidth
                value={k8sfile.replace('C:\\fakepath\\', '')}
                onClick={() => document.querySelector('#k8sfile').click()}
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
            <TextField
              select
              id="contextName"
              name="contextName"
              label="Context Name"
              fullWidth
              value={contextNameForForm}
              margin="normal"
              variant="outlined"
              // disabled={inClusterConfigForm === true}
              onChange={this.handleChange('contextNameForForm')}
            >
              {contextsFromFile && contextsFromFile.map((ct) => (
                <MenuItem key={`ct_---_${ct.contextName}`} value={ct.contextName}>
                  {ct.contextName}
                  {ct.currentContext ? ' (default)' : ''}
                </MenuItem>
              ))}
            </TextField>
          </div>
        </div>
      </NoSsr>
    );
  }

  meshIn = (showConfigured) => {
    const { classes } = this.props;

    return (
      <NoSsr>
        <div className={classes.root}>
          <div className={classes.currentConfigHeading}>
            <h4>
              Current Configuration Details
            </h4>
          </div>
          <div className={classes.changeConfigHeading}>
            <h4>
              Change Configuration...
            </h4>
          </div>

          <div className={classes.configure}>
            {showConfigured}
          </div>
          <Divider className={classes.vertical} orientation="vertical" />
          <Divider className={classes.horizontal} orientation="horizontal" />
          <div className={classes.changeConfigHeadingOne}>
            <h4>
              Change Configuration...
            </h4>
          </div>
          <div className={classes.buttonconfig}>
            <div className={classes.buttonsCluster}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                size="large"
                onClick={() => window.location.reload(false)}
                className={classes.button}
                data-cy="btnDiscoverCluster"
              >
                Discover Cluster
              </Button>
            </div>
          </div>
        </div>
      </NoSsr>
    );
  }

  render() {
    return this.configureTemplate();
  }
}

MeshConfigComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});
const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig').toJS();
  return k8sconfig;
};

export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps,
)(withRouter(withSnackbar(MeshConfigComponent))));
