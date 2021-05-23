import React from "react";
import { connect } from "react-redux";
import { withStyles, Typography, IconButton } from "@material-ui/core/";
import CloseIcon from "@material-ui/icons/Close";
import { bindActionCreators } from "redux";
import { withRouter } from "next/router";
import { withSnackbar } from "notistack";
import { updateK8SConfig, updateProgress } from "../../../lib/store";

const styles = (theme) => (
  console.log("theme", theme),
  {
    infoContainer: {
      position: "relative",
      display: "flex",
      flexDirection: "column",
      width: "20rem",
      padding: "5rem 1rem",
      boxShadow: "0px 1px 6px 1px rgba(0,0,0,0.75)",
    },
    infoTitle: {
      position: "absolute",
      bottom: "12.50rem",
      right: "10rem",
      color: "#647881",
      background: "#F1F3F4",
      padding: ".5rem 5rem .75rem 1.5rem",
      borderRadius: "0.25rem",
      fontSize: ".8rem",
    },
    infoItemContainer: {
      display: "flex",
      flexDirection: "row",
    },
    infoLabel: {
      fontSize: ".9rem",
      color: theme.palette.text.primary,
      marginRight: "1rem",
    },
    infoData: {
      fontSize: ".9rem",
      color: theme.palette.text.secondary,
    },
  }
);

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
    const { inClusterConfig, contextName } = this.state;
    return (
      <>
        <div className={classes.infoContainer}>
          <Typography className={classes.infoTitle}>Status</Typography>
          <div className={classes.infoItemContainer}>
            <Typography className={classes.infoLabel}>Current-Context:</Typography>
            <Typography className={classes.infoData}>
              {inClusterConfig ? "Using In Cluster Config" : contextName}
            </Typography>
          </div>
          <div className={classes.infoItemContainer}>
            <Typography className={classes.infoLabel}>Cluster:</Typography>
            <Typography className={classes.infoData}>
              {inClusterConfig ? "Using In Cluster Config" : "Using Out Of Cluster Config"}
            </Typography>
          </div>
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
