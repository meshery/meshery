import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import {
  Container,
  Fade,
  IconButton,
  Typography,
  Switch,
  FormControlLabel,
  Card,
  CardContent,
} from "@material-ui/core/";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { withRouter } from "next/router";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import { updateK8SConfig, updateProgress } from "../../../lib/store";
import subscribeOperatorStatusEvents from "../../graphql/subscriptions/OperatorStatusSubscription";
import subscribeMeshSyncStatusEvents from "../../graphql/subscriptions/MeshSyncStatusSubscription";
import changeOperatorState from "../../graphql/mutations/OperatorStatusMutation";
import fetchMesheryOperatorStatus from "../../graphql/queries/OperatorStatusQuery";
import MesheryOperatorIcon from "../icons/MesheryOperatorIcon";

const styles = () => ({
  // Container
  cardContainer: {
    position: "relative",
    display: "flex",
    flexDirection: "row",
    justifyContent: "space-between",
    padding: "2rem 6rem",
  },
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
  // Card
  card: {
    position: "relative",
    width: "12rem",
    minWidth: "10rem",
    border: "1px solid gray",
    borderRadius: "0.75rem",
    //top: "2rem",
    margin: "0rem 0rem 6rem 0rem",
    ["@media (max-width:1024px)"]: {
      //eslint-disable-line no-useless-computed-key
      margin: "0rem 0rem 6rem 0",
    },
  },
  cardChecked: {
    height: "15rem",
    marginBottom: "1rem",
  },
  cardUnchecked: {
    height: "10rem",
  },
  cardContent: {
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    alignItems: "center",
    height: "100%",
    width: "100%",
    padding: "0",
  },
  contentTop: {
    background: "#434343",
    height: "12rem",
    width: "100%",
    display: "flex",
    alignItems: "center",
  },
  contentTopUnchecked: {
    background: "#434343",
    height: "100%",
    width: "100%",
    display: "flex",
    alignItems: "center",
  },
  contentTopSwitcher: {
    paddingLeft: "2rem",
  },
  iconContainer: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    marginLeft: "1rem",
  },
  cardIcon: {
    width: "3rem",
  },
  cardIconText: {
    color: "white",
    fontSize: "0.85rem",
    textAlign: "center",
    "&:first-letter": {
      textTransform: "capitalize",
    },
  },
  contentBottomChecked: {
    background: "white",
    height: "6rem",
    width: "100%",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
  },
  contentBottomUnchecked: {
    display: "none",
  },
  contentBottomInput: {
    border: "1px solid lightgray",
    borderRadius: "5px",
    width: "9rem",
    height: "2rem",
    marginBottom: "0.15rem",
    fontSize: "0.75rem",
    padding: "0.50rem",
  },
  topInputIcon: {
    position: "absolute",
    fontSize: "1.25rem",
    color: "lightgray",
    bottom: "4.25rem",
    left: "9rem",
    cursor: "pointer",
    zIndex: "99999",
    "&:hover": {
      color: "grey",
    },
  },
  file: {
    display: "none",
  },
  // Meshery Operator Status
  statusContainer: {
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
  statusTitle: {
    fontSize: ".85rem",
  },
  statusData: {
    fontSize: ".75rem",
    marginTop: ".5rem",
    color: "#CACACA",
  },
});

const MeshySwitch = withStyles({
  switchBase: {
    color: "grey",
    "&$checked": {
      color: "#00B39F",
    },
    "&$checked + $track": {
      backgroundColor: "#00B39F",
    },
  },
  checked: {},
  track: {},
})(Switch);

class MesheryOperator extends React.Component {
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
      isChecked: false,
    };
    this.ref = React.createRef();
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

  handleChange = (e) => {
    this.handleOperatorSwitch();
    this.setState({ isChecked: e.target.checked });
  };
  render() {
    const { classes } = this.props;
    const { operatorInstalled, meshSyncInstalled, NATSInstalled } = this.state;
    return (
      <Fade timeout={{ enter: "500ms" }} in="true">
        <Container className={classes.cardContainer}>
          {" "}
          <Card className={`${classes.card} ${classes.cardChecked}`} variant="outlined">
            <CardContent className={classes.cardContent}>
              <div className={classes.contentTop}>
                <div className={classes.iconContainer}>
                  <MesheryOperatorIcon className={classes.cardIcon} alt={`Meshery Operator icon`} />
                  <Typography className={classes.cardIconText} color="primary">
                    Meshery Operator
                  </Typography>
                </div>
                <FormControlLabel
                  className={classes.contentTopSwitcher}
                  control={<MeshySwitch checked={this.state.isChecked} name="Meshery Operator" />}
                  onChange={this.handleChange}
                />
              </div>
              <div className={classes.contentBottomChecked}></div>
            </CardContent>
          </Card>
          {!this.state.isChecked ? null : (
            <div className={classes.statusContainer}>
              <div>
                <Typography className={classes.statusTitle}>Operator State</Typography>
                <Typography className={classes.statusData}>{operatorInstalled ? "Active" : "Disabled"}</Typography>
              </div>
              <div>
                <Typography className={classes.statusTitle}>MeshSync State</Typography>
                <Typography className={classes.statusData}>{meshSyncInstalled ? "Active" : "Disabled"}</Typography>
              </div>
              <div>
                <Typography className={classes.statusTitle}>NATS State</Typography>
                <Typography className={classes.statusData}>{NATSInstalled ? "Active" : "Disabled"}</Typography>
              </div>
            </div>
          )}
        </Container>
      </Fade>
    );
  }
}

MesheryOperator.propTypes = {
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
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(MesheryOperator)))
);
