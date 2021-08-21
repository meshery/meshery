import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import {
  NoSsr,
  Button,
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
  TextField,
  List,
  ListItem,
  ListItemText,
  CircularProgress,
} from "@material-ui/core";
import blue from "@material-ui/core/colors/blue";
import CloudUploadIcon from "@material-ui/icons/CloudUpload";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { withRouter } from "next/router";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import { updateK8SConfig, updateProgress } from "../lib/store";
import dataFetch from "../lib/data-fetch";
import subscribeOperatorStatusEvents from "./graphql/subscriptions/OperatorStatusSubscription";
import subscribeMeshSyncStatusEvents from "./graphql/subscriptions/MeshSyncStatusSubscription";
import changeOperatorState from "./graphql/mutations/OperatorStatusMutation";
import fetchMesheryOperatorStatus from "./graphql/queries/OperatorStatusQuery";
import deployMeshSync from "./graphql/queries/DeployMeshSyncQuery";
import connectToNats from "./graphql/queries/DeployNatsQuery";
import PromptComponent from "./PromptComponent";

const styles = (theme) => ({
  root : { padding : theme.spacing(5), },
  buttons : { display : "flex",
    justifyContent : "flex-end", },
  button : { marginTop : theme.spacing(3),
    marginLeft : theme.spacing(1), },
  buttonsCluster : { display : "flex",
    justifyContent : "center", },
  margin : { margin : theme.spacing(1), },
  alreadyConfigured : { textAlign : "center",
    padding : theme.spacing(20), },
  colorSwitchBase : { color : blue[300],
    "&$colorChecked" : { color : blue[500],
      "& + $colorBar" : { backgroundColor : blue[500], }, }, },
  colorBar : {},
  colorChecked : {},
  fileLabel : { width : "100%", },
  fileLabelText : { cursor : "pointer",
    "& *" : { cursor : "pointer", }, },
  inClusterLabel : { paddingRight : theme.spacing(2), },
  alignCenter : { textAlign : "center", },
  alignLeft : { textAlign : "left",
    marginBottom : theme.spacing(2), },
  fileInputStyle : { display : "none", },
  icon : { width : theme.spacing(2.5), },
  configure : {
    display : "inline-block",
    width : "48%",
    wordWrap : "break-word",
    [theme.breakpoints.down(945)] : { width : "100%", },
  },
  vertical : {
    display : "inline-block",
    height : 150,
    marginBottom : -60,
    [theme.breakpoints.down(945)] : { display : "none", },
  },
  horizontal : { display : "none",
    [theme.breakpoints.down(945)] : { display : "block", }, },
  formconfig : {
    display : "inline-block",
    marginLeft : 30,
    [theme.breakpoints.up(946)] : { width : "45%", },
    [theme.breakpoints.down(945)] : { width : "100%",
      marginLeft : 0, },
  },
  currentConfigHeading : {
    display : "inline-block",
    width : "48%",
    textAlign : "center",
    [theme.breakpoints.down(945)] : { width : "100%", },
  },
  changeConfigHeading : {
    display : "inline-block",
    width : "48%",
    textAlign : "center",
    [theme.breakpoints.down(945)] : { display : "none", },
  },
  operationButton : {
    [theme.breakpoints.down(1180)] : {
      marginRight : "25px",
    },
  },
  contentContainer : {
    [theme.breakpoints.down(1050)] : {
      flexDirection : "column",
    },
  },
  configBoxContainer : {
    [theme.breakpoints.down(1050)] : {
      flexGrow : 0,
      maxWidth : '100%',
      flexBasis : '100%',
    },
  },
  changeConfigHeadingOne : { display : "none",
    [theme.breakpoints.down(945)] : { display : "inline-block",
      width : "100%",
      textAlign : "center", }, },
  buttonconfig : { display : "inline-block",
    width : "48%",
    [theme.breakpoints.down(945)] : { width : "100%", }, },
  paper : { padding : theme.spacing(2), },
  heading : { textAlign : "center", },
  grey : { background : "WhiteSmoke",
    padding : theme.spacing(2),
    borderRadius : "inherit", },
});

class MeshConfigComponent extends React.Component {
  constructor(props) {
    super(props);
    const {
      inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer
    } = props;
    this.state = {
      inClusterConfig, // read from store
      inClusterConfigForm : inClusterConfig,
      k8sfile, // read from store
      k8sfileElementVal : "",
      contextName, // read from store
      contextNameForForm : "",
      contextsFromFile : [],
      clusterConfigured, // read from store
      configuredServer,
      k8sfileError : false,
      ts : new Date(),

      operatorInstalled: false,
      operatorVersion: "N/A",
      meshSyncInstalled: false,
      meshSyncVersion: "N/A",
      NATSState: "UNKNOWN",
      NATSVersion: "N/A",

      operatorSwitch: false,
      operatorProcessing: false,


      meshSyncStatusEventsSubscription: null,
      operatorStatusEventsSubscription: null,

    };
    this.ref = React.createRef();
  }

  static getDerivedStateFromProps(props, state) {
    const {
      inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer
    } = props;
    if (props.ts > state.ts) {
      let newState = {
        inClusterConfig,
        k8sfile,
        k8sfileElementVal : "",
        contextName,
        clusterConfigured,
        configuredServer,
        ts : props.ts,
      };

      // If contextsFromFile is empty then add the default value to it
      if (!state.contextsFromFile?.length)
        newState = { ...newState, contextsFromFile : [{ contextName, currentContext : true }] };
      return newState;
    }
    return {};
  }


  componentDidMount() {
    const self = this;
    // Subscribe to the operator events
    let meshSyncStatusEventsSubscription = subscribeMeshSyncStatusEvents((res) => {
      if (res.meshsync?.error) {
        self.handleError(res.meshsync?.error?.description || "MeshSync could not be reached");
        return;
      }
    });


    let operatorStatusEventsSubscription = subscribeOperatorStatusEvents(self.setOperatorState);
    fetchMesheryOperatorStatus().subscribe({
      next: (res) => {
        self.setOperatorState(res);
      },
      error: (err) => console.log("error at operator scan: " + err),
    });

    self.setState({meshSyncStatusEventsSubscription, operatorStatusEventsSubscription})
  }


  componentWillUnmount () {
    this.state.meshSyncStatusEventsSubscription.dispose()
    this.state.operatorStatusEventsSubscription.dispose()
  }

  setOperatorState = (res) => {
    console.log("incoming change")
    const self = this;
    if (res.operator?.error) {
      self.handleError("Operator could not be reached")(res.operator?.error?.description);
      self.setState({operatorProcessing: false})
      return false;
    }

    if (res.operator?.status === "ENABLED") {
      self.setState({operatorProcessing: false})
      res.operator?.controllers?.forEach((controller) => {
        if (controller.name === "broker" && controller.status == "CONNECTED") {
          self.setState({
            NATSState: controller.status,
            NATSVersion: controller.version,
          });
        } else if (controller.name === "meshsync" && controller.status == "ENABLED") {
          self.setState({ meshSyncInstalled : true,
            meshSyncVersion : controller.version, });
        }
      });
      self.setState({ operatorInstalled : true,
        operatorSwitch : true,
        operatorVersion : res.operator?.version, });
      return true;
    }

    if (res.operator?.status === "DISABLED") self.setState({operatorProcessing: false})

    if(res.operator?.status === "PROCESSING") {
      console.log("setting to processing")
      self.setState({operatorProcessing: true})
    }


    self.setState({
      operatorInstalled: false,
      NATSState: "UNKNOWN",
      meshSyncInstalled: false,
      operatorSwitch: false,
      operatorVersion: "N/A",
      meshSyncVersion: "N/A",
      NATSVersion: "N/A",
    });

    return false;
  };

  handleOperatorSwitch = () => {
    const self = this;
    const variables = { status : `${!self.state.operatorSwitch
      ? "ENABLED"
      : "DISABLED"}`, };
    self.props.updateProgress({ showProgress : true });

    changeOperatorState((response, errors) => {
      self.props.updateProgress({ showProgress : false });
      if (errors !== undefined) {
        self.handleError("Unable to install operator");
      }
      self.props.enqueueSnackbar("Operator " + response.operatorStatus.toLowerCase(), { variant : "success",
        autoHideDuration : 2000,
        action : (key) => (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        ), });
      self.setState((state) => ({ operatorSwitch : !state.operatorSwitch }));
    }, variables);
  };

  handleChange = (name) => {
    const self = this;
    return (event) => {
      if (name === "inClusterConfigForm") {
        self.setState({ [name] : event.target.checked, ts : new Date() });
        return;
      }
      if (name === "k8sfile") {
        if (event.target.value !== "") {
          self.setState({ k8sfileError : false });
        }
        self.setState({ k8sfileElementVal : event.target.value });
        self.fetchContexts();
      }
      self.setState({ [name] : event.target.value, ts : new Date() });
      this.handleSubmit();
    };
  };

  handleSubmit = () => {
    const { inClusterConfigForm, k8sfile } = this.state;
    if (!inClusterConfigForm && k8sfile === "") {
      this.setState({ k8sfileError : true });
      return;
    }
    this.submitConfig();
  };

  fetchContexts = () => {
    const { inClusterConfigForm } = this.state;
    const fileInput = document.querySelector("#k8sfile");
    const formData = new FormData();
    if (inClusterConfigForm) {
      return;
    }
    if (fileInput.files.length == 0) {
      this.setState({ contextsFromFile : [], contextNameForForm : "" });
      return;
    }
    // formData.append('contextName', contextName);
    formData.append("k8sfile", fileInput.files[0]);
    this.props.updateProgress({ showProgress : true });
    const self = this;
    dataFetch(
      "/api/system/kubernetes/contexts",
      {
        credentials : "same-origin",
        method : "POST",
        credentials : "include",
        body : formData,
      },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          let ctName = "";
          result.forEach(({ contextName, currentContext }) => {
            if (currentContext) {
              ctName = contextName;
            }
          });
          self.setState({ contextsFromFile : result, contextNameForForm : ctName });
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
    formData.append("inClusterConfig", inClusterConfigForm
      ? "on"
      : ""); // to simulate form behaviour of a checkbox
    if (!inClusterConfigForm) {
      formData.append("contextName", contextNameForForm);
      formData.append("k8sfile", fileInput.files[0]);
    }
    this.props.updateProgress({ showProgress : true });
    const self = this;
    dataFetch(
      "/api/system/kubernetes",
      {
        credentials : "same-origin",
        method : "POST",
        credentials : "include",
        body : formData,
      },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          //prompt
          const modal = this.ref.current;
          const self = this;
          if (self.state.operatorSwitch) {
            setTimeout(async () => {
              let response = await modal.show({ title : "Remove Meshery Operator from this cluster?",
                subtitle :
                  "Meshery is now disconnected from your Kubernetes cluster. Do you want to remove the Meshery Operator from your cluster as well?",
                options : ["yes", "no"], });
              if (response == "yes") {
                const variables = { status : "DISABLED", };
                self.props.updateProgress({ showProgress : true });

                changeOperatorState((response, errors) => {
                  self.props.updateProgress({ showProgress : false });
                  if (errors !== undefined) {
                    self.handleError("Operator action failed");
                  }
                  self.props.enqueueSnackbar("Operator " + response.operatorStatus.toLowerCase(), { variant : "success",
                    autoHideDuration : 2000,
                    action : (key) => (
                      <IconButton
                        key="close"
                        aria-label="Close"
                        color="inherit"
                        onClick={() => self.props.closeSnackbar(key)}
                      >
                        <CloseIcon />
                      </IconButton>
                    ), });
                  self.setState((state) => ({ operatorSwitch : !state.operatorSwitch }));
                }, variables);
              }
            }, 100);
          }
          this.setState({ clusterConfigured : true,
            configuredServer : result.configuredServer,
            contextName : result.contextName, });
          this.props.enqueueSnackbar("Kubernetes config was successfully validated!", { variant : "success",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ), });
          this.props.updateK8SConfig({ k8sConfig : {
            inClusterConfig : inClusterConfigForm,
            k8sfile,
            contextName : result.contextName,
            clusterConfigured : true,
            configuredServer : result.configuredServer,
          }, });
        }
      },
      self.handleError("Kubernetes config could not be validated")
    );
  };

  handleKubernetesClick = () => {
    this.props.updateProgress({ showProgress : true });
    const self = this;
    dataFetch(
      "/api/system/kubernetes/ping",
      { credentials : "same-origin",
        credentials : "include", },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Kubernetes was successfully pinged!", {
            variant : "success",
            "data-cy" : "k8sSuccessSnackbar",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ),
          });
        }
      },
      self.handleError("Kubernetes config could not be validated")
    );
  };

  handleOperatorClick = () => {
    this.props.updateProgress({ showProgress : true });
    const self = this;
    fetchMesheryOperatorStatus().subscribe({ next : (res) => {
      console.log(res);
      let state = self.setOperatorState(res);
      self.props.updateProgress({ showProgress : false });
      if (state == true) {
        this.props.enqueueSnackbar("Operator was successfully pinged!", { variant : "success",
          autoHideDuration : 2000,
          action : (key) => (
            <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
              <CloseIcon />
            </IconButton>
          ), });
      } else {
        self.handleError("Operator could not be reached")("Operator is disabled");
      }
    },
    error : self.handleError("Operator could not be pinged"), });
  };

handleNATSClick = () => {
  this.props.updateProgress({ showProgress: true });
  const self = this;

  connectToNats().subscribe({
    next: (res) => {
      if(res.connectToNats === "PROCESSING") {
        this.props.updateProgress({ showProgress: false });
        this.props.enqueueSnackbar(`Reconnecting to NATS...`, {
          variant: "info",
          action: (key) => (
            <IconButton key="close" aria-label="close" color="inherit" onClick={() => self.props.closesnackbar(key)}>
              <CloseIcon />
            </IconButton>
          ),
          autohideduration: 7000,
        })
      }
      if(res.connectToNats === "CONNECTED") {
        this.props.updateProgress({ showProgress: false });
        this.props.enqueueSnackbar(`Successfully connected to NATS`, {
          variant: "success",
          action: (key) => (
            <IconButton key="close" aria-label="close" color="inherit" onClick={() => self.props.closesnackbar(key)}>
              <CloseIcon />
            </IconButton>
          ),
          autohideduration: 7000,
        })
      }
          
    },
    error: self.handleError("Failed to request reconnection with NATS"),
  });

};

  handleMeshSyncClick = () => {
    this.props.updateProgress({ showProgress : true });
    const self = this;

    deployMeshSync().subscribe({
      next: (res) => {
        if(res.deployMeshsync === "PROCESSING") {
          this.props.updateProgress({ showProgress: false });
          this.props.enqueueSnackbar(`Meshsync deployment in progress`, {
            variant: "info",
            action: (key) => (
              <IconButton key="close" aria-label="close" color="inherit" onClick={() => self.props.closesnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ),
            autohideduration: 7000,
          })
        }

      },
      error: self.handleError("Failed to request Meshsync redeployment"),
    });
  };

  handleError = (msg) => (error) => {
    this.props.updateProgress({ showProgress: false });
    self.setOperatorState(res);
    const self = this;
    this.props.enqueueSnackbar(`${msg}: ${error}`, { variant : "error",
      action : (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration : 7000, });
  };

  handleReconfigure = () => {
    const self = this;
    dataFetch(
      "/api/system/kubernetes",
      { credentials : "same-origin",
        method : "DELETE",
        credentials : "include", },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          //prompt
          const modal = this.ref.current;
          const self = this;
          if (self.state.operatorSwitch) {
            setTimeout(async () => {
              let response = await modal.show({ title : "Remove Meshery Operator from this cluster?",
                subtitle :
                  "Meshery is now disconnected from your Kubernetes cluster. Do you want to remove the Meshery Operator from your cluster as well?",
                options : ["yes", "no"], });
              if (response == "yes") {
                const variables = { status : "DISABLED", };
                self.props.updateProgress({ showProgress : true });

                changeOperatorState((response, errors) => {
                  self.props.updateProgress({ showProgress : false });
                  if (errors !== undefined) {
                    self.handleError("Operator action failed");
                  }
                  self.props.enqueueSnackbar("Operator " + response.operatorStatus.toLowerCase(), { variant : "success",
                    autoHideDuration : 2000,
                    action : (key) => (
                      <IconButton
                        key="close"
                        aria-label="Close"
                        color="inherit"
                        onClick={() => self.props.closeSnackbar(key)}
                      >
                        <CloseIcon />
                      </IconButton>
                    ), });
                  self.setState((state) => ({ operatorSwitch : !state.operatorSwitch }));
                }, variables);
              }
            }, 100);
          }
          this.setState({
            inClusterConfigForm : false,
            inClusterConfig : false,
            k8sfile : "",
            k8sfileElementVal : "",
            k8sfileError : false,
            contextName : "",
            contextNameForForm : "",
            clusterConfigured : false,
          });
          this.props.updateK8SConfig({ k8sConfig : {
            inClusterConfig : false,
            k8sfile : "",
            contextName : "",
            clusterConfigured : false,
          }, });
          this.props.enqueueSnackbar("Kubernetes config was successfully removed!", { variant : "success",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ), });
        }
      },
      self.handleError("Kubernetes config could not be validated")
    );
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
      NATSState,
      NATSVersion,
      operatorSwitch,
    } = this.state;
    let showConfigured = "";
    const self = this;
    if (clusterConfigured) {
      let chp = (
        <Chip
          // label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
          label={inClusterConfig
            ? "Using In Cluster Config"
            : contextName}
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
            <ListItemText
              primary="Context Name"
              secondary={inClusterConfig
                ? "Using In Cluster Config"
                : contextName}
              data-cy="itemListContextName"
            />
          </ListItem>
          <ListItem>
            <ListItemText
              primary="Server Name"
              secondary={inClusterConfig
                ? "In Cluster Server"
                : configuredServer || ""}
              data-cy="itemListServerName"
            />
          </ListItem>
        </List>
      );
      if (configuredServer) {
        chp = <Tooltip title={`Server: ${configuredServer}`}>{chp}</Tooltip>;
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
      showConfigured = <div>{lst}</div>;
    }

    const operator = (
      <React.Fragment>
        <div>
          <Grid container spacing={1} >
            <Grid item xs={12} md={4} className={classes.operationButton}>
              <List>
                <ListItem>
                  <Tooltip
                    title={operatorInstalled
                      ? `Version: ${operatorVersion}`
                      : "Not Available"}
                    aria-label="meshSync"
                  >
                    <Chip
                      // label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
                      label={"Operator"}
                      // onDelete={self.handleReconfigure}
                      onClick={self.handleOperatorClick}
                      icon={<img src="/static/img/meshery-operator.svg" className={classes.icon} />}
                      variant="outlined"
                      data-cy="chipOperator"
                    />
                  </Tooltip>
                </ListItem>
              </List>
            </Grid>
            {operatorInstalled && 
            <>
              <Grid item xs={12} md={4}>
                <List>
                  <ListItem>
                    <Tooltip
                      title={meshSyncInstalled ? `Redeploy MeshSync` : "Not Available"}
                      aria-label="meshSync"
                    >
                      <Chip
                        label={"MeshSync"}
                        onClick={self.handleMeshSyncClick}
                        icon={<img src="/static/img/meshsync.svg" className={classes.icon} />}
                        variant="outlined"
                        data-cy="chipMeshSync"
                      />
                    </Tooltip>
                  </ListItem>
                </List>
              </Grid>
              <Grid item xs={12} md={4}>
                <List>
                  <ListItem>
                    <Tooltip
                      title={NATSState === "CONNECTED" ? `Reconnect NATS` : "Not Available"}
                      aria-label="nats"
                    >
                      <Chip
                        label={"NATS"}
                        onClick={self.handleNATSClick}
                        icon={<img src="/static/img/meshsync.svg" className={classes.icon} />}
                        variant="outlined"
                        data-cy="chipNATS"
                      />
                    </Tooltip>
                  </ListItem>
                </List>
              </Grid>
            </>
            }
          </Grid>
          <Grid container spacing={1}>
            <Grid item xs={12} md={4}>
              <List>
                <ListItem>
                  <ListItemText primary="Operator State" secondary={operatorInstalled
                    ? "Active"
                    : "Disabled"} />
                </ListItem>
                <ListItem>
                  <ListItemText primary="Operator Version" secondary={operatorVersion} />
                </ListItem>
              </List>
            </Grid>
            <Grid item xs={12} md={4}>
              <List>
                <ListItem>
                  <ListItemText primary="MeshSync State" secondary={meshSyncInstalled
                    ? "Active"
                    : "Disabled"} />
                </ListItem>
                <ListItem>
                  <ListItemText primary="MeshSync Version" secondary={meshSyncVersion} />
                </ListItem>
              </List>
            </Grid>
            <Grid item xs={12} md={4}>
              <List>
                <ListItem>
                  <ListItemText primary="NATS State" secondary={NATSState} />
                </ListItem>
                <ListItem>
                  <ListItemText primary="NATS Version" secondary={NATSVersion} />
                </ListItem>
              </List>
            </Grid>
          </Grid>
        </div>
        <div className={classes.grey}>
          <FormGroup>
            <FormControlLabel
              control={
                <Switch
                  checked={operatorSwitch}
                  onClick={self.handleOperatorSwitch}
                  disabled={self.state.operatorProcessing}
                  name="OperatorSwitch"
                  color="primary"
                />
              }
              label="Meshery Operator"
            />
            {self.state.operatorProcessing && <CircularProgress />}
          </FormGroup>
        </div>
      </React.Fragment>
    );

    if (this.props.tabs == 0) {
      return this.meshOut(showConfigured, operator);
    }
    return this.meshIn(showConfigured, operator);
  };

  meshOut = (showConfigured, operator) => {
    const { classes } = this.props;
    const {
      k8sfile, k8sfileElementVal, contextNameForForm, contextsFromFile, contextName
    } = this.state;

    return (
      <NoSsr>
        <PromptComponent ref={this.ref} />
        <div className={classes.root}>
          <Grid container spacing={5} className={classes.contentContainer}>
            <Grid item spacing={1} xs={12} md={6} className={classes.configBoxContainer}>
              <div className={classes.heading}>
                <h4>Cluster Configuration</h4>
              </div>
              <Paper className={classes.paper}>
                <div>{showConfigured}</div>
                <div className={classes.grey}>
                  <FormGroup>
                    <input
                      id="k8sfile"
                      type="file"
                      value={k8sfileElementVal}
                      onChange={this.handleChange("k8sfile")}
                      className={classes.fileInputStyle}
                    />
                    <TextField
                      id="k8sfileLabelText"
                      name="k8sfileLabelText"
                      className={classes.fileLabelText}
                      label="Upload kubeconfig"
                      variant="outlined"
                      fullWidth
                      value={k8sfile.replace("C:\\fakepath\\", "")}
                      onClick={() => document.querySelector("#k8sfile").click()}
                      margin="normal"
                      InputProps={{ readOnly : true,
                        endAdornment : (
                          <InputAdornment position="end">
                            <CloudUploadIcon />
                          </InputAdornment>
                        ), }}
                    />
                  </FormGroup>
                  <TextField
                    select
                    id="contextName"
                    name="contextName"
                    label="Context Name"
                    fullWidth
                    value={contextNameForForm || contextName}
                    margin="normal"
                    variant="outlined"
                    // disabled={inClusterConfigForm === true}
                    onChange={this.handleChange("contextNameForForm")}
                  >
                    {contextsFromFile &&
                      contextsFromFile.map((ct) => (
                        <MenuItem key={`ct_---_${ct.contextName}`} value={ct.contextName}>
                          {ct.contextName}
                          {ct.currentContext
                            ? " (default)"
                            : ""}
                        </MenuItem>
                      ))}
                  </TextField>
                </div>
              </Paper>
            </Grid>

            <Grid item spacing={1} xs={12} md={6} className={classes.configBoxContainer}>
              <div className={classes.heading}>
                <h4>Operator Configuration</h4>
              </div>
              <Paper className={classes.paper}>{operator}</Paper>
            </Grid>
          </Grid>
        </div>
      </NoSsr>
    );
  };

  meshIn = (showConfigured, operator) => {
    const { classes } = this.props;

    return (
      <NoSsr>
        <PromptComponent ref={this.ref} />
        <div className={classes.root}>
          <Grid container spacing={5}>
            <Grid item spacing={1} xs={12} md={6}>
              <div className={classes.heading}>
                <h4>Cluster Configuration</h4>
              </div>
              <Paper className={classes.paper}>
                <div>{showConfigured}</div>
                <div className={classes.grey}>
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
              </Paper>
            </Grid>

            <Grid item xs={12} md={6} spacing={1}>
              <div className={classes.heading}>
                <h4>Operator Configuration</h4>
              </div>
              <Paper className={classes.paper}>{operator}</Paper>
            </Grid>
          </Grid>
        </div>
      </NoSsr>
    );
  };

  render() {
    return this.configureTemplate();
  }
}

MeshConfigComponent.propTypes = { classes : PropTypes.object.isRequired, };

const mapDispatchToProps = (dispatch) => ({ updateK8SConfig : bindActionCreators(updateK8SConfig, dispatch),
  updateProgress : bindActionCreators(updateProgress, dispatch), });
const mapStateToProps = (state) => {
  const k8sconfig = state.get("k8sConfig").toJS();
  return k8sconfig;
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(MeshConfigComponent)))
);
