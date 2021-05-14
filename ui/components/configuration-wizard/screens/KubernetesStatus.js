import React from "react";
import { connect } from "react-redux";
import { withStyles, Typography, IconButton } from "@material-ui/core/";
import CloseIcon from "@material-ui/icons/Close";
import { bindActionCreators } from "redux";
import { withRouter } from "next/router";
import { withSnackbar } from "notistack";
import { updateK8SConfig, updateProgress } from "../../../lib/store";
import dataFetch from "../../../lib/data-fetch";

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
    const { inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer } = props;
    this.state = {
      inClusterConfig, // read from store
      inClusterConfigForm: inClusterConfig,
      k8sfile, // read from store
      k8sfileElementVal: "",
      contextName, // read from store
      contextNameForForm: "",
      contextsFromFile: [],
      clusterConfigured, // read from store
      configuredServer,
      k8sfileError: false,
      ts: new Date(),
    };
    console.log("cc", clusterConfigured);
  }

  static getDerivedStateFromProps(props, state) {
    const { inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer } = props;
    if (props.ts > state.ts) {
      return {
        inClusterConfig,
        k8sfile,
        k8sfileElementVal: "",
        contextName,
        clusterConfigured,
        configuredServer,
        ts: props.ts,
      };
    }
    return {};
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

  render() {
    const { classes } = this.props;
    const { inClusterConfig, contextName, clusterConfigured, configuredServer } = this.state;
    return (
      <>
        <div className={classes.infoContainer}>
          <Typography className={classes.infoStatus}>Status</Typography>
          <Typography className={classes.infoContext}>
            Current-Context:{inClusterConfig ? "Using In Cluster Config" : contextName}
          </Typography>
          <Typography className={classes.infoContext}>
            Cluster: {inClusterConfig ? "Using In Cluster Config" : "Using Out Of Cluster Config"}
          </Typography>
        </div>
      </>
    );
  }
}

const mapDispatchToProps = (dispatch) => ({
  updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});
const mapStateToProps = (state) => {
  const k8sconfig = state.get("k8sConfig").toJS();
  return k8sconfig;
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(KubernetesStatus)))
);
