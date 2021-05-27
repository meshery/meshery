import React from "react";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { withRouter } from "next/router";
import { withSnackbar } from "notistack";
import { updateK8SConfig, updateProgress } from "../../../lib/store";
import dataFetch from "../../../lib/data-fetch";
import {
  withStyles,
  FormGroup,
  TextField,
  InputAdornment,
  MenuItem,
  Container,
  Switch,
  FormControlLabel,
  Card,
  CardContent,
  Typography,
  IconButton,
} from "@material-ui/core/";
import BackupIcon from "@material-ui/icons/Backup";
import CloseIcon from "@material-ui/icons/Close";

import KubernetesIcon from "../icons/KubernetesIcon";

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

const styles = (theme) => ({
  // Container
  cardContainer: {
    position: "relative",
    display: "flex",
    flexDirection: "row",
    justifyContent: "space-between",
    padding: "2rem 6rem",
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
  contentBottomInputChecked: {
    background: "white",
    height: "6rem",
    width: "100%",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
  },
  contentBottomInputUnchecked: {
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
  // Inputs
  fileInputStyle: {
    opacity: "0.01",
  },
  contentBottomChecked: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "flex-end",
  },
  contentBottomUpperInput: {
    width: "11rem",
    fontSize: "0.75rem",
    marginLeft: "2.4rem",
    marginTop: "-1rem",
    marginBottom: "0rem",
  },
  contentBottomLowerInput: {
    width: "11rem",
    marginBottom: "-1rem",
    fontSize: "0.75rem",
    marginTop: "0",
  },
  // Status
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
});

class KubernetesScreen extends React.Component {
  constructor(props) {
    super(props);
    const { inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer } = this.props;
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
      isChecked: false,
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
  componentDidMount() {
    const { inClusterConfig, contextName } = this.state;
    const { handleConnectToKubernetes } = this.props;
    if (inClusterConfig || contextName) {
      this.setState({ isChecked: true });
      handleConnectToKubernetes(true);
    }
  }
  handleSwitch = (name, checked) => {
    const { inClusterConfig, contextName } = this.state;
    if (inClusterConfig || contextName) {
      this.setState({ isChecked: checked });
      if (this.props.handleConnectToKubernetes) {
        this.props.handleConnectToKubernetes(checked);
      }
    }
  };
  handleChange = (name) => {
    const self = this;
    return (event) => {
      if (name === "inClusterConfigForm") {
        self.setState({ [name]: event.target.checked, ts: new Date() });
        return;
      }
      if (name === "k8sfile") {
        if (event.target.value !== "") {
          self.setState({ k8sfileError: false });
        }
        self.setState({ k8sfileElementVal: event.target.value });
        self.fetchContexts();
      }
      self.setState({ [name]: event.target.value, ts: new Date() });
      this.handleSubmit();
    };
  };
  fetchContexts = () => {
    const { inClusterConfigForm } = this.state;
    const fileInput = document.querySelector("#k8sfile");
    const formData = new FormData();
    if (inClusterConfigForm) {
      return;
    }
    if (fileInput.files.length == 0) {
      this.setState({ contextsFromFile: [], contextNameForForm: "" });
      return;
    }
    // formData.append('contextName', contextName);
    formData.append("k8sfile", fileInput.files[0]);
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch(
      "/api/k8sconfig/contexts",
      {
        credentials: "same-origin",
        method: "POST",
        credentials: "include",
        body: formData,
      },
      (result) => {
        this.props.updateProgress({ showProgress: false });
        if (typeof result !== "undefined") {
          let ctName = "";
          result.forEach(({ contextName, currentContext }) => {
            if (currentContext) {
              ctName = contextName;
            }
          });
          self.setState({ contextsFromFile: result, contextNameForForm: ctName });
          self.submitConfig();
        }
      },
      self.handleError("Kubernetes config could not be validated")
    );
  };

  submitConfig = () => {
    const { inClusterConfigForm, k8sfile, contextNameForForm } = this.state;
    const fileInput = document.querySelector("#k8sfile");
    const formData = new FormData();
    formData.append("inClusterConfig", inClusterConfigForm ? "on" : ""); // to simulate form behaviour of a checkbox
    if (!inClusterConfigForm) {
      formData.append("contextName", contextNameForForm);
      formData.append("k8sfile", fileInput.files[0]);
    }
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch(
      "/api/k8sconfig",
      {
        credentials: "same-origin",
        method: "POST",
        credentials: "include",
        body: formData,
      },
      (result) => {
        this.props.updateProgress({ showProgress: false });
        if (typeof result !== "undefined") {
          //prompt
          const modal = this.ref.current;
          const self = this;
          if (self.state.operatorSwitch) {
            setTimeout(async () => {
              let response = await modal.show({
                title: "Do you wanna remove Operator from this cluster?",
                subtitle: "The Meshery Operator will be uninstalled from the cluster if responded with 'yes'",
                options: ["yes", "no"],
              });
              if (response == "yes") {
                const variables = {
                  status: "DISABLED",
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
                  self.setState((state) => ({ operatorSwitch: !state.operatorSwitch }));
                }, variables);
              }
            }, 100);
          }
          this.setState({
            clusterConfigured: true,
            configuredServer: result.configuredServer,
            contextName: result.contextName,
          });
          this.props.enqueueSnackbar("Kubernetes config was successfully validated!", {
            variant: "success",
            autoHideDuration: 2000,
            action: (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ),
          });
          this.props.updateK8SConfig({
            k8sConfig: {
              inClusterConfig: inClusterConfigForm,
              k8sfile,
              contextName: result.contextName,
              clusterConfigured: true,
              configuredServer: result.configuredServer,
            },
          });
        }
      },
      self.handleError("Kubernetes config could not be validated")
    );
  };
  handleSubmit = () => {
    const { inClusterConfigForm, k8sfile } = this.state;
    if (!inClusterConfigForm && k8sfile === "") {
      this.setState({ k8sfileError: true });
      return;
    }
    this.submitConfig();
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

  render() {
    const { classes } = this.props;
    const { inClusterConfig, contextName } = this.state;
    return (
      <Container className={classes.cardContainer}>
        <Card className={`${classes.card} ${classes.cardChecked}`} variant="outlined">
          <CardContent className={classes.cardContent}>
            <div className={classes.contentTop}>
              <div className={classes.iconContainer}>
                <KubernetesIcon className={classes.cardIcon} alt="Kubernetes icon" />
                <Typography className={classes.cardIconText} color="primary">
                  Kubernetes
                </Typography>
              </div>
              <FormControlLabel
                className={classes.contentTopSwitcher}
                control={<MeshySwitch checked={this.state.isChecked} name="Kubernetes" />}
                onChange={this.handleSwitch}
              />
            </div>
            <div className={classes.contentBottomInputChecked}>
              <FormGroup>
                <input
                  id="k8sfile"
                  type="file"
                  value={this.state.k8sfileElementVal}
                  onChange={this.handleChange("k8sfile")}
                  className={classes.fileInputStyle}
                />
                <TextField
                  id="k8sfileLabelText"
                  name="k8sfileLabelText"
                  className={classes.contentBottomUpperInput}
                  label="Upload kubeconfig"
                  variant="outlined"
                  fullWidth
                  value={this.state.k8sfile.replace("C:\\fakepath\\", "")}
                  onClick={() => document.querySelector("#k8sfile").click()}
                  margin="normal"
                  InputProps={{
                    readOnly: true,
                    endAdornment: (
                      <InputAdornment position="end">
                        <BackupIcon style={{ marginRight: "-.8rem" }} />
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
                value={this.state.contextNameForForm}
                margin="normal"
                variant="outlined"
                // disabled={inClusterConfigForm === true}
                onChange={this.handleChange("contextNameForForm")}
                className={classes.contentBottomLowerInput}
              >
                {this.state.contextsFromFile &&
                  this.state.contextsFromFile.map((ct) => (
                    <MenuItem key={`ct_---_${ct.contextName}`} value={ct.contextName}>
                      {ct.contextName}
                      {ct.currentContext ? " (default)" : ""}
                    </MenuItem>
                  ))}
              </TextField>
            </div>
          </CardContent>
        </Card>
        {this.props.k8sconfig?.k8sfile || this.state?.contextName ? (
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
        ) : null}
      </Container>
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
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(KubernetesScreen)))
);
