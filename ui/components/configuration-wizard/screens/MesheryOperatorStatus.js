import React from "react";
import PropTypes from "prop-types";
import Button from "@material-ui/core/Button";
import { withStyles } from "@material-ui/core/styles";
import {
  NoSsr,
  FormGroup,
  InputAdornment,
  Chip,
  IconButton,
  MenuItem,
  Tooltip,
  Paper,
  Grid,
  FormControlLabel,
  Switch,
} from "@material-ui/core";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { withRouter } from "next/router";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import { updateK8SConfig, updateProgress } from "../../../lib/store";
import dataFetch from "../../../lib/data-fetch";
import subscribeOperatorStatusEvents from "./graphql/subscriptions/OperatorStatusSubscription";
import subscribeMeshSyncStatusEvents from "./graphql/subscriptions/MeshSyncStatusSubscription";
import changeOperatorState from "./graphql/mutations/OperatorStatusMutation";
import fetchMesheryOperatorStatus from "./graphql/queries/OperatorStatusQuery";


const styles = (theme) => ({
  infoContainer: {
    position: "relative",
    display: "flex",
    flexDirection: "row",
    justifyContent: "space-between",
    width: "20rem",
    height: "12.5rem",
    padding: "2rem 2rem",
    marginRight: "3rem",
    boxShadow: "0px 1px 6px 1px rgba(0,0,0,0.75)",
  },
  infoContext: {
    fontSize: ".85rem",
  },
  infoKind: {
    fontSize: ".75rem",
    marginTop: ".5rem",
    color: "#CACACA",
  },
});

class MesheryOperatorStatus extends React.Component {
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

      operatorInstalled: false,
      operatorVersion: "N/A",
      meshSyncInstalled: false,
      meshSyncVersion: "N/A",
      NATSInstalled: false,
      NATSVersion: "N/A",

      operatorSwitch: false,
    };
    this.ref = React.createRef();
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

  componentDidMount() {
    const self = this;
    // Subscribe to the operator events
    subscribeMeshSyncStatusEvents((res) => {
      if (res.meshsync?.error) {
        self.handleError(res.meshsync?.error?.description || "MeshSync could not be reached");
        return;
      }
    });

    subscribeOperatorStatusEvents(self.setOperatorState);
    fetchMesheryOperatorStatus().subscribe({
      next: (res) => {
        self.setOperatorState(res);
      },
      error: (err) => console.log("error at operator scan: " + err),
    });
  }

  setOperatorState = (res) => {
    const self = this;
    if (res.operator?.error) {
      self.handleError(res.operator?.error?.description || "Operator could not be reached");
      return;
    }

    if (res.operator?.status === "ENABLED") {
      res.operator?.controllers?.forEach((controller) => {
        if (controller.name === "broker" && controller.status == "ENABLED") {
          self.setState({
            NATSInstalled: true,
            NATSVersion: controller.version,
          });
        } else if (controller.name === "meshsync" && controller.status == "ENABLED") {
          self.setState({
            meshSyncInstalled: true,
            meshSyncVersion: controller.version,
          });
        }
      });
      self.setState({
        operatorInstalled: true,
        operatorSwitch: true,
        operatorVersion: res.operator?.version,
      });
      return;
    }

    self.setState({
      operatorInstalled: false,
      NATSInstalled: false,
      meshSyncInstalled: false,
      operatorSwitch: false,
      operatorVersion: "N/A",
      meshSyncVersion: "N/A",
      NATSVersion: "N/A",
    });
  };

  handleOperatorSwitch = () => {
    const self = this;
    const variables = {
      status: `${!self.state.operatorSwitch ? "ENABLED" : "DISABLED"}`,
    };
    self.props.updateProgress({ showProgress: true });

    changeOperatorState((response, errors) => {
      self.props.updateProgress({ showProgress: false });
      if (errors !== undefined) {
        self.handleError("Operator action failed");
      }
      self.props.enqueueSnackbar("Operator " + response.operatorStatus.toLowerCase(), {
        variant: "success",
        autoHideDuration: 2000,
        action: (key) => (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        ),
      });
      self.setState((state) => ({ operatorSwitch: !state.operatorSwitch }));
    }, variables);
  };

  handleOperatorClick = () => {
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch(
      "/api/system/operator/status",
      {
        credentials: "same-origin",
        credentials: "include",
      },
      (result) => {
        this.setState({
          operatorInstalled: result["operator-installed"] == "true" ? true : false,
          NATSInstalled: result["broker-installed"] == "true" ? true : false,
          meshSyncInstalled: result["meshsync-installed"] == "true" ? true : false,
        });
        this.props.updateProgress({ showProgress: false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Operator was successfully pinged!", {
            variant: "success",
            autoHideDuration: 2000,
            action: (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ),
          });
        }
      },
      self.handleError("Operator could not be pinged")
    );
  };

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
    const {
      inClusterConfig,
      contextName,
      clusterConfigured,
      configuredServer,
      operatorInstalled,
      operatorVersion,
      meshSyncInstalled,
      meshSyncVersion,
      NATSInstalled,
      NATSVersion,
      operatorSwitch,
    } = this.state;
    

    return (
      <div className={classes.infoContainer}>
        <div>
          <Typography className={classes.infoContext}>Operator State</Typography>
          <Typography className={classes.infoKind}>{operatorInstalled ? "Active" : "Disabled"}</Typography>
        </div>
        <div>
          <Typography className={classes.infoContext}>MeshSync State</Typography>
          <Typography className={classes.infoKind}>{meshSyncInstalled ? "Active" : "Disabled"}</Typography>
        </div>
        <div>
          <Typography className={classes.infoContext}>NATS State</Typography>
          <Typography className={classes.infoKind}>{NATSInstalled ? "Active" : "Disabled"}</Typography>
        </div>
      </div>
    );
  };

  render() {
    this.operatorSwitch()
    return this.configureTemplate();
  }
}

MesheryOperatorStatus.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});
const mapStateToProps = (state) => {
  const k8sconfig = state.get("k8sConfig").toJS();
  return k8sconfig;
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(MesheryOperatorStatus)))
);
