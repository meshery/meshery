// @ts-nocheck
import React from "react";
import PropTypes from "prop-types";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";
import { Autocomplete } from "@material-ui/lab";
import { withStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import { URLValidator } from "../../utils/URLValidator";
import {
  NoSsr,
  Tooltip,
  MenuItem,
  IconButton,
  CircularProgress,
  FormControl,
  FormLabel,
  RadioGroup,
  FormControlLabel,
  Radio,
  Divider,
  Link,
  ExpansionPanel,
  ExpansionPanelSummary,
  ExpansionPanelDetails,
} from "@material-ui/core";
import TextField from "@material-ui/core/TextField";
import { withSnackbar } from "notistack";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import CloseIcon from "@material-ui/icons/Close";
import GetAppIcon from "@material-ui/icons/GetApp";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import HelpOutlineOutlinedIcon from "@material-ui/icons/HelpOutlineOutlined";
import SaveOutlinedIcon from "@material-ui/icons/SaveOutlined";
import {
  updateLoadTestData,
  updateStaticPrometheusBoardConfig,
  updateLoadTestPref,
  updateProgress,
} from "../../lib/store";
import dataFetch from "../../lib/data-fetch";
import MesheryChart from "../MesheryChart";
import LoadTestTimerDialog from "../load-test-timer-dialog";
import GrafanaCustomCharts from "../telemetry/grafana/GrafanaCustomCharts";
import { durationOptions } from "../../lib/prePopulatedOptions";
import fetchControlPlanes from "../graphql/queries/ControlPlanesQuery";
import { ctxUrl, getK8sClusterIdsFromCtxId } from "../../utils/multi-ctx";
import { iconMedium } from "../../css/icons.styles";

// =============================== HELPER FUNCTIONS ===========================

/**
 * generatePerformanceProfile takes in data and generate a performance
 * profile object from it
 * @param {*} data
 */
function generatePerformanceProfile(data) {
  const {
    id,
    name,
    loadGenerator,
    additional_options,
    endpoint,
    serviceMesh,
    concurrentRequest,
    qps,
    duration,
    requestHeaders,
    requestCookies,
    requestBody,
    contentType,
    caCertificate
  } = data;

  const performanceProfileName = MesheryPerformanceComponent.generateTestName(name, serviceMesh);

  return {
    ...(id && { id }),
    name : performanceProfileName,
    load_generators : [loadGenerator],
    endpoints : [endpoint],
    service_mesh : serviceMesh,
    concurrent_request : concurrentRequest,
    qps,
    duration,
    request_headers : requestHeaders,
    request_body : requestBody,
    request_cookies : requestCookies,
    content_type : contentType,
    metadata : {
      additional_options : [additional_options],
      ca_certificate : {
        file : caCertificate.file,
        name : caCertificate.name
      }
    }
  };
}

// =============================== PERFORMANCE COMPONENT =======================

const loadGenerators = ["fortio", "wrk2", "nighthawk"];

const infoFlags = (
  <>
  Only .json files are supported.
  </>
)

const infoCRTCertificates = (
  <>
  Only .crt files are supported.
  </>
);

const infoloadGenerators = (
  <>
    Which load generators does Meshery support?
    <ul>
      <li>
        fortio - Fortio load testing library, command line tool, advanced echo server and web UI in go (golang). Allows
        to specify a set query-per-second load and record latency histograms and other useful stats.{" "}
      </li>
      <li> wrk2 - A constant throughput, correct latency recording variant of wrk.</li>
      <li>
        {" "}
        nighthawk - Enables users to run distributed performance tests to better mimic real-world, distributed systems
        scenarios.
      </li>
    </ul>
    <Link
      style={{ textDecoration : "underline" }}
      color="inherit"
      href="https://docs.meshery.io/functionality/performance-management"
    >
      {" "}
      Performance Management
    </Link>
  </>
);
const styles = (theme) => ({
  title : {
    textAlign : 'center',
    minWidth : 400,
    padding : '10px',
    color : '#fff',
    backgroundColor : theme.palette.type === 'dark' ? theme.palette.secondary.headerColor : theme.palette.secondary.mainBackground,
  },
  wrapperClss : { padding : theme.spacing(10), position : "relative", paddingTop : theme.spacing(5) },
  buttons : { display : "flex", justifyContent : "flex-end" },
  spacing : {
    marginTop : theme.spacing(3), marginLeft : theme.spacing(1)
  },
  button : {
    backgroundColor : theme.palette.type === "dark" ? "#00B39F" : "#607d8b",
    "&:hover" : {
      backgroundColor : theme.palette.type === "dark" ? "#00B39F" : "#607d8b"
    },
    color : "#fff" },
  upload : {
    paddingLeft : "0.7rem",
    paddingTop : "8px"
  },
  expansionPanel : { boxShadow : "none", border : "1px solid rgb(196,196,196)" },
  margin : { margin : theme.spacing(1) },
  chartTitle : { textAlign : "center" },
  chartTitleGraf : {
    textAlign : "center",
    // marginTop: theme.spacing(5),
  },
  chartContent : {
    // minHeight: window.innerHeight * 0.7,
  },
  centerTimer : {
    width : "100%",
    height : "100%",
    position : "absolute",
    top : "0",
    left : "0",
    backgroundColor : "rgba(0, 0, 0, 0.6)",
    zIndex : 1201,
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
  },
  paper : {
    backgroundColor : theme.palette.background.paper,
    border : "2px solid #000",
    boxShadow : theme.shadows[5],
    padding : theme.spacing(2, 4, 3),
  },
  smallIcons : {
    width : "15px",
    height : "18px",
    marginBottom : theme.spacing(1),
    marginLeft : theme.spacing(0.3),
  },
  radio : {
    '&.Mui-checked' : {
      color : theme.palette.type === 'dark' ? theme.palette.secondary.focused : theme.palette.primary
    },
  },
});

class MesheryPerformanceComponent extends React.Component {
  constructor(props) {
    super(props);
    const {
      testName = "",
      meshName = "",
      url = "",
      qps = "0",
      c = "0",
      t = "30s",
      result,
      staticPrometheusBoardConfig,
      performanceProfileID,
      profileName,
      loadGenerator,
      additional_options,
      headers,
      cookies,
      reqBody,
      contentType,
      metadata,
    } = props;

    this.state = {
      testName,
      meshName,
      url,
      qps,
      c,
      t,
      tValue : t,
      loadGenerator : loadGenerator || "fortio",
      additional_options : additional_options || "",
      result,
      headers : headers || "",
      cookies : cookies || "",
      reqBody : reqBody || "",
      contentType : contentType || "",
      certificate : {},
      certificateKey : {},
      caCertificate : {},
      profileName : profileName || "",
      performanceProfileID : performanceProfileID || "",
      timerDialogOpen : false,
      blockRunTest : false,
      urlError : false,
      tError : "",
      disableTest : !(URLValidator(url) || this.isJsonString(additional_options)),
      testUUID : this.generateUUID(),
      staticPrometheusBoardConfig,
      selectedMesh : "",
      availableAdapters : [],
      availableSMPMeshes : [],
      disableAvailableOptionsUploadButton : false,
      disableAvailableOptionsInputField : false,
      metadata
    };
  }

  isJsonString(str) {
    try {
      JSON.parse(str);
    } catch (e) {
      return false;
    }
    return true;
  }

  handleChange = (name) => (event) => {

    if (name === "caCertificate") {
      if (!event.target.files?.length) return;

      const file = event.target.files[0];
      const reader = new FileReader();
      reader.addEventListener("load", (evt) => {
        this.setState({
          caCertificate : {
            name : file.name,
            file : evt.target.result
          }
        })
        console.log("test: ", name);
      })
      reader.readAsText(file);
    }

    if (name === "url" && event.target.value !== "") {
      let urlPattern = event.target.value;

      let val = URLValidator(urlPattern);
      if (!val) {
        this.setState({ disableTest : true });
        this.setState({ urlError : true });
      } else {
        this.setState({ disableTest : false });
        this.setState({ urlError : false });
      }
    } else this.setState({ urlError : false });

    if (name === "additional_options" ) {
      const { value } = event.target;

      // Check if the target event is an input element (typing) or a file input (upload)
      const isFileUpload = event.target.getAttribute("type") === "file";

      if (isFileUpload) {
        // Handle file upload
        const file = event.target.files[0];
        if (file) {
          const reader = new FileReader();
          reader.onload = (event) => {
            try {
              const fileContent = event.target.result;
              // Validate JSON
              JSON.parse(fileContent);
              this.setState({
                additional_options : fileContent,
                jsonError : false,
              });
            } catch (error) {
              this.setState({
                additional_options : event.target.result,
                jsonError : true,
              });
            }
          };
          reader.readAsText(file);
        }
      } else {
        // Handle text input
        try {
          // empty text input exception
          if (value !== "") JSON.parse(value);
          this.setState({
            additional_options : value,
            jsonError : false,
          });
        } catch (error) {
          this.setState({
            additional_options : value,
            jsonError : true,
          });
        }
      }
    }
    this.setState({ [name] : event.target.value });
  };

  handleDurationChange = (event, newValue) => {
    this.setState({ tValue : newValue });
    if (newValue !== null) {
      this.setState({ tError : "" });
    }
  };

  handleInputDurationChange = (event, newValue) => {
    this.setState({ t : newValue });
  };

  handleSubmit = () => {
    const { url, t } = this.state;

    if (url === "") {
      this.setState({ urlError : true });
      return;
    }

    let err = false;
    let tNum = 0;
    try {
      tNum = parseInt(t.substring(0, t.length - 1));
    } catch (ex) {
      err = true;
    }

    if (
      t === "" ||
      t === null ||
      !(t.toLowerCase().endsWith("h") || t.toLowerCase().endsWith("m") || t.toLowerCase().endsWith("s")) ||
      err ||
      tNum <= 0
    ) {
      this.setState({ tError : "error-autocomplete-value" });
      return;
    }

    if (!this.state.performanceProfileID) {
      this.submitProfile(({ id }) => this.submitLoadTest(id));
      return;
    }

    this.submitLoadTest(this.state.performanceProfileID);
  };

  submitProfile = (cb) => {
    const self = this.state;
    const profile = generatePerformanceProfile({
      name : self.profileName,
      loadGenerator : self.loadGenerator,
      additional_options : self.additional_options,
      endpoint : self.url,
      serviceMesh : self.meshName,
      concurrentRequest : +self.c || 0,
      qps : +self.qps || 0,
      duration : self.t,
      requestHeaders : self.headers,
      requestCookies : self.cookies,
      requestBody : self.reqBody,
      contentType : self.contentType,
      caCertificate : self.caCertificate,
      testName : self.testName,
      id : self.performanceProfileID,
    });

    this.handleProfileUpload(profile, true, cb);
  };

  handleAbort = () => {
    this.setState({
      profileName : "",
      loadGenerator : "",
      additional_options : "",
      url : "",
      meshName : "",
      c : 0,
      qps : 0,
      t : "30s",
      headers : "",
      cookies : "",
      reqBody : "",
      contentType : "",
      testName : "",
      performanceProfileID : "",
    });
    this.setState({ disableTest : true });
    return;
  };

  handleProfileUpload = (body, generateNotif, cb) => {
    if (generateNotif) this.props.updateProgress({ showProgress : true });

    dataFetch(
      "/api/user/performance/profiles",
      { method : "POST", credentials : "include", body : JSON.stringify(body) },
      (result) => {
        if (typeof result !== "undefined") {
          this.props.updateProgress({ showProgress : false });

          this.setState({ performanceProfileID : result.id }, () => {
            if (cb) cb(result);
          });

          if (generateNotif) {
            this.props.enqueueSnackbar("Performance Profile Created!", {
              variant : "success",
              autoHideDuration : 2000,
              action : (key) => (
                <IconButton
                  key="close"
                  aria-label="Close"
                  color="inherit"
                  onClick={() => this.props.closeSnackbar(key)}
                >
                  <CloseIcon />
                </IconButton>
              ),
            });
          }
        }
      },
      (err) => {
        console.error(err);
        this.props.updateProgress({ showProgress : false });
        this.props.enqueueSnackbar("Failed to create performance profile", {
          variant : "error",
          autoHideDuration : 2000,
          action : (key) => (
            <IconButton style={iconMedium} key="close" aria-label="Close" color="inherit" onClick={() => this.props.closeSnackbar(key)}>
              <CloseIcon />
            </IconButton>
          ),
        });
      }
    );
  };

  /**
   * generateTestName takes in test name and service mesh name
   * and generates a random name (if test name is an empty string or is falsy) or
   * will return the given name
   *
   * @param {string} name
   * @param {string} meshName
   * @returns {string}
   */
  static generateTestName = (name, meshName) => {
    if (!name || name.trim() === "") {
      const mesh = meshName === "" || meshName === "None" ? "No mesh" : meshName;
      return `${mesh}_${new Date().getTime()}`;
    }

    return name;
  };

  submitLoadTest = (id) => {
    const { testName, meshName, url, qps, c, t, loadGenerator, additional_options, testUUID, headers, cookies, reqBody, contentType } =
      this.state;

    const computedTestName = MesheryPerformanceComponent.generateTestName(testName, meshName);
    this.setState({ testName : computedTestName });

    const t1 = t.substring(0, t.length - 1);
    const dur = t.substring(t.length - 1, t.length).toLowerCase();

    const data = {
      name : computedTestName,
      mesh : meshName === "" || meshName === "None" ? "" : meshName, // to prevent None from getting to the DB
      url,
      qps,
      c,
      t : t1,
      dur,
      uuid : testUUID,
      loadGenerator,
      additional_options : additional_options,
      headers : headers,
      cookies : cookies,
      reqBody : reqBody,
      contentType : contentType,
    };
    const params = Object.keys(data)
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
      .join("&");
    console.log(params);

    const runURL = ctxUrl(`/api/user/performance/profiles/${id}/run`, this.props?.selectedK8sContexts) + "&cert=true";
    this.startEventStream(`${runURL}${this.props?.selectedK8sContexts?.length > 0 ? "&" : "?"}${params}`);
    this.setState({ blockRunTest : true }); // to block the button
  };

  handleSuccess() {
    const self = this;
    return (result) => {
      const { testName, meshName, url, qps, c, t, loadGenerator } = this.state;
      if (typeof result !== "undefined" && typeof result.runner_results !== "undefined") {
        self.props.enqueueSnackbar("fetched the data.", {
          variant : "success",
          autoHideDuration : 2000,
          action : (key) => (
            <IconButton style={iconMedium} key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
              <CloseIcon style={iconMedium} />
            </IconButton>
          ),
        });
        self.props.updateLoadTestData({
          loadTest : {
            testName,
            meshName,
            url,
            qps,
            c,
            t,
            loadGenerator,
            result,
          },
        });
        self.setState({ result, testUUID : self.generateUUID() });
      }
      self.closeEventStream();
      self.setState({ blockRunTest : false, timerDialogOpen : false });
    };
  }

  async startEventStream(url) {
    this.closeEventStream();
    this.eventStream = new EventSource(url);
    this.eventStream.onmessage = this.handleEvents();
    this.eventStream.onerror = this.handleError(
      "Connection to the server got disconnected. Load test might be running in the background. Please check the results page in a few."
    );
    this.props.enqueueSnackbar("Load test has been submitted", {
      variant : "info",
      autoHideDuration : 1000,
      action : (key) => (
        <IconButton style={iconMedium} key="close" aria-label="Close" color="inherit" onClick={() => this.props.closeSnackbar(key)}>
          <CloseIcon style={iconMedium} />
        </IconButton>
      ),
    });
  }

  handleEvents() {
    const self = this;
    let track = 0;
    return (e) => {
      const data = JSON.parse(e.data);
      switch (data.status) {
        case "info":
          self.props.enqueueSnackbar(data.message, {
            variant : "info",
            autoHideDuration : 1000,
            action : (key) => (
              <IconButton style={iconMedium} key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon style={iconMedium} />
              </IconButton>
            ),
          });
          if (track === 0) {
            self.setState({ timerDialogOpen : true, result : {} });
            track++;
          }
          break;
        case "error":
          self.handleError("Load test did not run with msg")(data.message);
          break;
        case "success":
          self.handleSuccess()(data.result);
          break;
      }
    };
  }

  closeEventStream() {
    if (this.eventStream && this.eventStream.close) {
      this.eventStream.close();
      this.eventStream = null;
    }
  }

  componentDidMount() {
    this.getStaticPrometheusBoardConfig();
    this.scanForMeshes();
    this.getLoadTestPrefs();
    this.getSMPMeshes();

    if (this.props.runTestOnMount) this.handleSubmit();
  }

  getLoadTestPrefs = () => {
    dataFetch(
      ctxUrl("/api/user/prefs", this.props?.selectedK8sContexts),
      { credentials : "same-origin", method : "GET" },
      (result) => {
        if (typeof result !== "undefined") {
          this.setState({
            qps : result.loadTestPrefs.qps,
            c : result.loadTestPrefs.c,
            t : result.loadTestPrefs.t,
            loadGenerator : result.loadTestPrefs.gen,
          });
        }
      },
      () => { }
    ); //error is already captured from the handler, also we have a redux-store for same & hence it's not needed here.
  };

  getStaticPrometheusBoardConfig = () => {
    const self = this;
    if (
      (self.props.staticPrometheusBoardConfig &&
        self.props.staticPrometheusBoardConfig !== null &&
        Object.keys(self.props.staticPrometheusBoardConfig).length > 0) ||
      (self.state.staticPrometheusBoardConfig &&
        self.state.staticPrometheusBoardConfig !== null &&
        Object.keys(self.state.staticPrometheusBoardConfig).length > 0)
    ) {
      return;
    }
    dataFetch(
      "/api/telemetry/metrics/static-board",
      { credentials : "include" },
      (result) => {
        if (
          typeof result !== "undefined" &&
          typeof result.cluster !== "undefined" &&
          typeof result.node !== "undefined" &&
          typeof result.cluster.panels !== "undefined" &&
          result.cluster.panels.length > 0 &&
          typeof result.node.panels !== "undefined" &&
          result.node.panels.length > 0
        ) {
          self.props.updateStaticPrometheusBoardConfig({
            staticPrometheusBoardConfig : result, // will contain both the cluster and node keys for the respective boards
          });
          self.setState({ staticPrometheusBoardConfig : result });
        }
      },
      self.handleError("unable to fetch pre-configured boards")
    );
  };


  getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(this.props.selectedK8sContexts, this.props.k8sconfig)
  }

  scanForMeshes = () => {
    const self = this;

    if (typeof self.props.k8sConfig === "undefined" || !self.props.k8sConfig.clusterConfigured) {
      return;
    }
    /**
     * ALL_MESH indicates that we are interested in control plane
     * component of all of the service meshes supported by meshsync v2
     */

    const ALL_MESH = { type : "ALL_MESH", k8sClusterIDs : this.getK8sClusterIds() };

    fetchControlPlanes(ALL_MESH).subscribe({
      next : (res) => {
        let result = res?.controlPlanesState;
        if (typeof result !== "undefined" && Object.keys(result).length > 0) {
          const adaptersList = [];
          result.forEach((mesh) => {
            if (mesh?.members.length > 0) {
              let name = mesh?.name;
              adaptersList.push(
                // Capatilize First Letter and replace undersocres
                name
                  .split(/ |_/i)
                  .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
                  .join(" ")
              );
            }
          });
          self.setState({ availableAdapters : adaptersList });
          result.forEach((mesh) => {
            self.setState({ selectedMesh : mesh?.name });
          });
        }
      },
      error : (err) => console.error(err),
    });
  };

  getSMPMeshes = () => {
    const self = this;
    dataFetch(
      "/api/mesh",
      { credentials : "include" },
      (result) => {
        if (result && Array.isArray(result.available_meshes)) {
          self.setState({ availableSMPMeshes : result.available_meshes.sort((m1, m2) => m1.localeCompare(m2)) });
        }
      },
      self.handleError("unable to fetch SMP meshes")
    );
  };

  generateUUID() {
    const { v4 : uuid } = require("uuid");
    return uuid();
  }

  handleError(msg) {
    const self = this;
    return (error) => {
      self.setState({ blockRunTest : false, timerDialogOpen : false });
      self.closeEventStream();
      let finalMsg = msg;
      if (typeof error === "string") {
        finalMsg = `${msg}: ${error}`;
      }
      self.props.enqueueSnackbar(finalMsg, {
        variant : "error",
        action : (key) => (
          <IconButton style={iconMedium} key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
            <CloseIcon style={iconMedium}/>
          </IconButton>
        ),
        autoHideDuration : 4000,
      });
    };
  }

  handleCertificateUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
      const newMetadata = {
        ...this.state.metadata,
        ca_certificate : {
          ...this.state.metadata.ca_certificate,
          name : file.name,
        },
      };
      this.setState({ metadata : newMetadata });
    }
  };


  handleTimerDialogClose = () => {
    this.setState({ timerDialogOpen : false });
  };

  render() {
    const { classes, grafana, prometheus } = this.props;
    const {
      timerDialogOpen,
      blockRunTest,
      url,
      qps,
      c,
      t,
      loadGenerator,
      additional_options,
      meshName,
      result,
      urlError,
      jsonError,
      tError,
      testUUID,
      selectedMesh,
      availableAdapters,
      headers,
      cookies,
      reqBody,
      contentType,
      tValue,
      disableTest,
      profileName,
      metadata
    } = this.state;

    let staticPrometheusBoardConfig;
    if (
      this.props.staticPrometheusBoardConfig &&
      this.props.staticPrometheusBoardConfig != null &&
      Object.keys(this.props.staticPrometheusBoardConfig).length > 0
    ) {
      staticPrometheusBoardConfig = this.props.staticPrometheusBoardConfig;
    } else {
      staticPrometheusBoardConfig = this.state.staticPrometheusBoardConfig;
    }
    let chartStyle = {};
    if (timerDialogOpen) {
      chartStyle = { opacity : 0.3 };
    }
    let displayStaticCharts = null;
    let displayGCharts = null;
    let displayPromCharts = null;

    availableAdapters.forEach((item) => {
      const index = this.state.availableSMPMeshes.indexOf(item);
      if (index !== -1) this.state.availableSMPMeshes.splice(index, 1);
    });

    if (
      staticPrometheusBoardConfig &&
      staticPrometheusBoardConfig !== null &&
      Object.keys(staticPrometheusBoardConfig).length > 0 &&
      prometheus.prometheusURL !== ""
    ) {
      // only add testUUID to the board that should be persisted
      if (staticPrometheusBoardConfig.cluster) {
        staticPrometheusBoardConfig.cluster.testUUID = testUUID;
      }
      displayStaticCharts = (
        <React.Fragment>
          <Typography variant="h6" gutterBottom className={classes.chartTitle}>
            Node Metrics
          </Typography>
          <GrafanaCustomCharts
            boardPanelConfigs={[staticPrometheusBoardConfig.cluster, staticPrometheusBoardConfig.node]}
            prometheusURL={prometheus.prometheusURL}
          />
        </React.Fragment>
      );
    }
    if (prometheus.selectedPrometheusBoardsConfigs.length > 0) {
      displayPromCharts = (
        <React.Fragment>
          <Typography variant="h6" gutterBottom className={classes.chartTitleGraf}>
            Prometheus charts
          </Typography>
          <GrafanaCustomCharts
            boardPanelConfigs={prometheus.selectedPrometheusBoardsConfigs}
            prometheusURL={prometheus.prometheusURL}
          />
        </React.Fragment>
      );
    }
    if (grafana.selectedBoardsConfigs.length > 0) {
      displayGCharts = (
        <React.Fragment>
          <Typography variant="h6" gutterBottom className={classes.chartTitleGraf}>
            Grafana charts
          </Typography>
          <GrafanaCustomCharts
            boardPanelConfigs={grafana.selectedBoardsConfigs}
            grafanaURL={grafana.grafanaURL}
            grafanaAPIKey={grafana.grafanaAPIKey}
          />
        </React.Fragment>
      );
    }
    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.wrapperClss} style={this.props.style || {}}>
            <Grid container spacing={1}>
              <Grid item xs={12} md={6}>

                <TextField
                  id="profileName"
                  name="profileName"
                  label="Profile Name"
                  fullWidth
                  value={profileName}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleChange("profileName")}
                  inputProps={{
                    maxLength : 300,
                  }}
                  InputProps={{
                    endAdornment : (
                      <Tooltip title="Create a profile providing a name, if a profile name is not provided, a random one will be generated for you.">
                        <HelpOutlineOutlinedIcon style={{ color : "#929292" }} />
                      </Tooltip>
                    )
                  }}
                />
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  select
                  id="meshName"
                  name="meshName"
                  label="Service Mesh"
                  fullWidth
                  value={meshName === "" && selectedMesh !== "" ? selectedMesh : meshName}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleChange("meshName")}
                >
                  {availableAdapters &&
                    availableAdapters.map((mesh) => (
                      <MenuItem key={`mh_-_${mesh}`} value={mesh.toLowerCase()}>
                        {mesh}
                      </MenuItem>
                    ))}
                  {availableAdapters && availableAdapters.length > 0 && <Divider />}
                  <MenuItem key="mh_-_none" value="None">
                    None
                  </MenuItem>
                  {this.state.availableSMPMeshes &&
                    this.state.availableSMPMeshes.map((mesh) => (
                      <MenuItem key={`mh_-_${mesh}`} value={mesh.toLowerCase()}>
                        {mesh}
                      </MenuItem>
                    ))}
                </TextField>
              </Grid>
              <Grid item xs={12}>
                <TextField
                  required
                  id="url"
                  name="url"
                  label="URL to test"
                  type="url"
                  fullWidth
                  value={url}
                  error={urlError}
                  helperText={urlError ? "Please enter a valid URL along with protocol" : ""}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleChange("url")}
                  InputProps={{
                    endAdornment : (
                      <Tooltip title="The Endpoint where the load will be generated and the perfromance test will run against.">
                        <HelpOutlineOutlinedIcon style={{ color : "#929292" }} />
                      </Tooltip>
                    )
                  }}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <TextField
                  required
                  id="c"
                  name="c"
                  label="Concurrent requests"
                  type="number"
                  fullWidth
                  value={c}
                  inputProps={{ min : "0", step : "1" }}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleChange("c")}
                  InputLabelProps={{ shrink : true }}
                  InputProps={{
                    endAdornment : (
                      <Tooltip title="Load Testing tool will create this many concurrent request against the endpoint.">
                        <HelpOutlineOutlinedIcon style={{ color : "#929292" }} />
                      </Tooltip>
                    )
                  }}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <TextField
                  required
                  id="qps"
                  name="qps"
                  label="Queries per second"
                  type="number"
                  fullWidth
                  value={qps}
                  inputProps={{ min : "0", step : "1" }}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleChange("qps")}
                  InputLabelProps={{ shrink : true }}
                  InputProps={{
                    endAdornment : (
                      <Tooltip title="The Number of queries/second. If not provided then the MAX number of queries/second will be requested">
                        <HelpOutlineOutlinedIcon style={{ color : "#929292" }} />
                      </Tooltip>
                    )
                  }}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <Tooltip title={"Please use 'h', 'm' or 's' suffix for hour, minute or second respectively."}>
                  <Autocomplete
                    required
                    id="t"
                    name="t"
                    freeSolo
                    label="Duration*"
                    fullWidth
                    variant="outlined"
                    className={classes.errorValue}
                    classes={{ root : tError }}
                    value={tValue}
                    inputValue={t}
                    onChange={this.handleDurationChange}
                    onInputChange={this.handleInputDurationChange}
                    options={durationOptions}
                    style={{ marginTop : "16px", marginBottom : "8px" }}
                    renderInput={(params) => <TextField {...params} label="Duration*" variant="outlined" />}
                    InputProps={{
                      endAdornment : (
                        <Tooltip title="Default duration is 30 seconds">
                          <HelpOutlineOutlinedIcon style={{ color : "#929292" }} />
                        </Tooltip>
                      )
                    }}
                  />
                </Tooltip>
              </Grid>
              <Grid item xs={12} md={12}>
                <ExpansionPanel className={classes.expansionPanel}>
                  <ExpansionPanelSummary expanded={true} expandIcon={<ExpandMoreIcon />}>
                    <Typography align="center" color="textSecondary" variant="h6">
                      Advanced Options
                    </Typography>
                  </ExpansionPanelSummary>
                  <ExpansionPanelDetails>
                    <Grid container spacing={1}>
                      <Grid item xs={12}>
                        <TextField
                          id="headers"
                          name="headers"
                          label='Request Headers e.g. {"host":"bookinfo.meshery.io"}'
                          fullWidth
                          value={headers}
                          multiline
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange("headers")}
                        ></TextField>
                      </Grid>
                      <Grid item xs={12}>
                        <TextField
                          id="cookies"
                          name="cookies"
                          label='Request Cookies e.g. {"yummy_cookie":"choco_chip"}'
                          fullWidth
                          value={cookies}
                          multiline
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange("cookies")}
                        ></TextField>
                      </Grid>
                      <Grid item xs={12}>
                        <TextField
                          id="contentType"
                          name="contentType"
                          label="Content Type e.g. application/json"
                          fullWidth
                          value={contentType}
                          multiline
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange("contentType")}
                        ></TextField>
                      </Grid>
                      <Grid item xs={12} md={12}>
                        <TextField
                          id="cookies"
                          name="cookies"
                          label='Request Body e.g. {"method":"post","url":"http://bookinfo.meshery.io/test"}'
                          fullWidth
                          value={reqBody}
                          multiline
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange("reqBody")}
                        ></TextField>
                      </Grid>
                      <Grid item xs={12} md={12}>
                        <Grid item xs={6}>
                          <TextField
                            id="additional_options"
                            name="additional_options"
                            label="Additional Options e.g. { `requestPerSecond`: 20 }"
                            fullWidth
                            error={jsonError}
                            helperText={jsonError ? "Please enter a valid JSON string" : ""}
                            value={additional_options.length > 150 ? `${additional_options.slice(0,150)} .....` : additional_options}
                            multiline
                            margin="normal"
                            variant="outlined"
                            onChange={this.handleChange("additional_options")}
                          />
                          <label htmlFor="upload-additional-options"
                            style={{ paddingLeft : '0' }}
                            className={classes.upload}
                          >
                            <Button
                              variant="outlined"
                              onChange={this.handleChange("additional_options")}
                              aria-label="Upload Button"
                              component="span"
                              className={classes.button}
                            >
                              <input id="upload-additional-options"  type="file" accept={".json"} name="upload-button" hidden  data-cy="additional-options-upload-button" />
                              Browse
                            </Button>
                            <Tooltip title={infoFlags} interactive>
                              <HelpOutlineOutlinedIcon className={classes.smallIcons} />
                            </Tooltip>
                          </label>
                        </Grid>
                      </Grid>
                      <Grid container xs={12} md={12}>
                        <Grid item xs={6}>
                          <TextField
                            size="small"
                            variant="outlined"
                            label={this.state.caCertificate?.name || "Upload SSL Certificate e.g. .crt file"}
                            style={{ width : "100%", margin : '0.5rem 0' }}
                            value={metadata?.ca_certificate.name}
                          />
                          <label htmlFor="upload-cacertificate"
                            className={classes.upload}
                            style={{ paddingLeft : '0' }}
                          >
                            <Button
                              variant="outlined"
                              aria-label="Upload Button"
                              onChange={this.handleChange("caCertificate")} component="span"
                              className={classes.button}
                            >
                              <input id="upload-cacertificate" type="file" accept={".crt"} name="upload-button"  hidden data-cy="cacertificate-upload-button" onChange={this.handleCertificateUpload}/>
                              Browse
                            </Button>
                            <Tooltip title={infoCRTCertificates} interactive>
                              <HelpOutlineOutlinedIcon className={classes.smallIcons} />
                            </Tooltip>
                          </label>
                        </Grid>
                      </Grid>
                    </Grid>
                  </ExpansionPanelDetails>
                </ExpansionPanel>
              </Grid>
              <Grid item xs={12} md={4}>
                <FormControl component="loadGenerator" className={classes.margin}>
                  <FormLabel
                    component="loadGenerator"
                    style={{ display : "flex", alignItems : "center", flexWrap : "wrap" }}
                  >
                    Load generator
                    <Tooltip title={infoloadGenerators} interactive>
                      <HelpOutlineOutlinedIcon className={classes.smallIcons} />
                    </Tooltip>
                  </FormLabel>
                  <RadioGroup
                    aria-label="loadGenerator"
                    name="loadGenerator"
                    value={loadGenerator}
                    onChange={this.handleChange("loadGenerator")}
                    row
                  >
                    {loadGenerators.map((lg, index) => (
                      <FormControlLabel key={index} value={lg} disabled={lg==="wrk2"} control={<Radio color="primary" className={classes.radio} />} label={lg} />
                    ))}
                  </RadioGroup>
                </FormControl>
              </Grid>
            </Grid>
            <React.Fragment>
              <div className={classes.buttons}>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  className={classes.spacing}
                  disabled={disableTest}
                  onClick={() => this.handleAbort()}
                >
                  Clear
                </Button>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  onClick={() => this.submitProfile()}
                  className={classes.spacing}
                  disabled={disableTest}
                  startIcon={<SaveOutlinedIcon />}
                >
                  Save Profile
                </Button>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  onClick={this.handleSubmit}
                  className={classes.spacing}
                  disabled={blockRunTest || disableTest}
                >
                  {blockRunTest ? <CircularProgress size={30} /> : "Run Test"}
                </Button>
              </div>
            </React.Fragment>

            {timerDialogOpen ? (
              <div className={classes.centerTimer}>
                <LoadTestTimerDialog
                  open={timerDialogOpen}
                  t={t}
                  onClose={this.handleTimerDialogClose}
                  countDownComplete={this.handleTimerDialogClose}
                />
              </div>
            ) : null}

            {result && result.runner_results && (
              <div>
                <Typography variant="h6" gutterBottom className={classes.chartTitle} id="timerAnchor">
                  Test Results
                  <IconButton
                    key="download"
                    aria-label="download"
                    color="inherit"
                    // onClick={() => self.props.closeSnackbar(key) }
                    href={`/api/perf/profile/result/${encodeURIComponent(result.meshery_id)}`}
                  >
                    <GetAppIcon style={iconMedium}/>
                  </IconButton>
                </Typography>
                <div className={classes.chartContent} style={chartStyle}>
                  <MesheryChart
                    rawdata={[result && result.runner_results ? result : {}]}
                    data={[result && result.runner_results ? result.runner_results : {}]} />
                </div>
              </div>
            )}
          </div>
        </React.Fragment>

        {displayStaticCharts}

        {displayPromCharts}

        {displayGCharts}
      </NoSsr>
    );
  }
}

MesheryPerformanceComponent.propTypes = { classes : PropTypes.object.isRequired };

const mapDispatchToProps = (dispatch) => ({
  updateLoadTestData : bindActionCreators(updateLoadTestData, dispatch),
  updateStaticPrometheusBoardConfig : bindActionCreators(updateStaticPrometheusBoardConfig, dispatch),
  updateLoadTestPref : bindActionCreators(updateLoadTestPref, dispatch),
  updateProgress : bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const grafana = state.get("grafana").toJS();
  const prometheus = state.get("prometheus").toJS();
  const k8sConfig = state.get("k8sConfig");
  const staticPrometheusBoardConfig = state.get("staticPrometheusBoardConfig").toJS();
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    grafana,
    prometheus,
    staticPrometheusBoardConfig,
    k8sConfig,
    selectedK8sContexts,
  };
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryPerformanceComponent))
);
