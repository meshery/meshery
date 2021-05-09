import React from "react";
import { connect } from "react-redux";
import { withStyles, Typography, IconButton } from "@material-ui/core/";
import CloseIcon from '@material-ui/icons/Close';
import { bindActionCreators } from 'redux';
import { withRouter } from "next/router";
import { withSnackbar } from "notistack";
import { updateK8SConfig, updateProgress } from "../lib/store";
import dataFetch from '../lib/data-fetch';

const styles = () => ({
  infoContainer: {
    position: "relative",
    display: "flex",
    flexDirection: "column",
    width: "20rem",
    padding: "5rem 1rem",
    boxShadow: "0px 1px 6px 1px rgba(0,0,0,0.75)",
  },
  infoStatus: {
    position: "absolute",
    bottom: "12.50rem",
    right: "10rem",
    color: "#647881",
    background: "#F1F3F4",
    padding: ".5rem 5rem .75rem 1.5rem",
    borderRadius: "0.25rem",
    fontSize: ".8rem",
  },
  infoContext: {
    fontSize: ".9rem",
  },
  infoKind: {
    fontSize: ".75rem",
    color: "#CACACA",
  },
});

class KubernetesStatus extends React.Component {
  constructor(props) {
    super(props);
    const { inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer, } = props;
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
    }, self.handleError("Kubernetes config could not be validated"));
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
        //prompt
        const modal = this.ref.current;
        const self = this;
        if (self.state.operatorSwitch) {
          setTimeout(async () => {
            let response = await modal.show({ title: "Do you wanna remove Operator from this cluster?", subtitle: "The Meshery Operator will be uninstalled from the cluster if responded with 'yes'", options: ["yes", "no"] });
            if (response == "yes") {
              const variables = {
                status: "DISABLED",
              }
              self.props.updateProgress({ showProgress: true })
                    
              changeOperatorState((response, errors) => {
                self.props.updateProgress({ showProgress: false });
                if (errors !== undefined) {
                  self.handleError("Operator action failed")
                }
                self.props.enqueueSnackbar('Operator '+response.operatorStatus.toLowerCase(), {
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
                self.setState((state) => ({ operatorSwitch: !state.operatorSwitch }))
              }, variables);
            }
          }, 100);
        }
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
    }, self.handleError("Kubernetes config could not be validated"));
  }

  handleError = (msg) => (error) => {
    this.props.updateProgress({ showProgress: false });
    const self = this;
    this.props.enqueueSnackbar(`${msg}: ${error}`, {
      variant: "error",
      action: (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration: 7000,
    });
  };

  configureTemplate = () => {
    const { classes } = this.props;
    const { inClusterConfig, contextName, clusterConfigured, configuredServer } = this.state;
    let showConfigured = "";
    const self = this;
    if (clusterConfigured) {
      const lst = (
        <div className={classes.infoContainer}>
          <Typography className={classes.infoStatus}>Status</Typography>
          <Typography className={classes.infoContext}>
            Current-Context:{inClusterConfig ? "Using In Cluster Config" : contextName}
          </Typography>
          <Typography className={classes.infoContext}>
            Cluster: {inClusterConfig ? "Using In Cluster Config" : "Using Out Of Cluster Config"}
          </Typography>
        </div>
      );
      showConfigured = <div>{lst}</div>;
    }
    if (!clusterConfigured) {
      const lst = <div>Cluster not configured</div>;
      showConfigured = <div>{lst}</div>;
    }
  };

  render() {
    this.configureTemplete();
  }
}


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
)(withRouter(withSnackbar(KubernetesStatus))));
