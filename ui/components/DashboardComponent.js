import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import {
  NoSsr,
  Chip,
  IconButton,
  Button,
  Card,
  CardContent,
  Typography,
  CardHeader,
  Tooltip,
  TableContainer,
  Table,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
  Paper,
  Select,
  MenuItem,
} from "@material-ui/core";
import blue from "@material-ui/core/colors/blue";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import SettingsIcon from "@material-ui/icons/Settings";
import AddIcon from "@material-ui/icons/AddCircleOutline";
import { withRouter } from "next/router";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import { updateGrafanaConfig, updatePrometheusConfig, updateProgress } from "../lib/store";
import dataFetch from "../lib/data-fetch";
// import subscribeServiceMeshEvents from "./graphql/subscriptions/ServiceMeshSubscription";
import subscribeDataPlaneEvents from "./graphql/subscriptions/DataPlanesSubscription";
import subscribeControlPlaneEvents from "./graphql/subscriptions/ControlPlaneSubscription";
import subscribeOperatorStatusEvents from "./graphql/subscriptions/OperatorStatusSubscription";
import fetchControlPlanes from "./graphql/queries/ControlPlanesQuery";
import fetchDataPlanes from "./graphql/queries/DataPlanesQuery";
import fetchAvailableAddons from "./graphql/queries/AddonsStatusQuery";
import { submitPrometheusConfigure } from "./PrometheusComponent";
import { submitGrafanaConfigure } from "./GrafanaComponent";
import { versionMapper } from "../utils/nameMapper";
// podNameMapper,
//import MesheryMetrics from "./MesheryMetrics";

const styles = (theme) => ({
  root : { backgroundColor : "#eaeff1", },
  chip : { marginRight : theme.spacing(1),
    marginBottom : theme.spacing(1), },
  buttons : { display : "flex",
    justifyContent : "flex-end", },
  button : { marginTop : theme.spacing(3),
    marginLeft : theme.spacing(1), },
  metricsButton : { width : "240px", },
  alreadyConfigured : { textAlign : "center", },
  margin : { margin : theme.spacing(1), },
  colorSwitchBase : { color : blue[300],
    "&$colorChecked" : { color : blue[500],
      "& + $colorBar" : { backgroundColor : blue[500], }, }, },
  colorBar : {},
  colorChecked : {},
  fileLabel : { width : "100%", },
  fileLabelText : {},
  inClusterLabel : { paddingRight : theme.spacing(2), },
  alignCenter : { textAlign : "center", },
  icon : { width : theme.spacing(2.5), },
  istioIcon : { width : theme.spacing(1.5), },
  settingsIcon : { width : theme.spacing(2.5),
    paddingRight : theme.spacing(0.5), },
  addIcon : { width : theme.spacing(2.5),
    paddingRight : theme.spacing(0.5), },
  cardHeader : { fontSize : theme.spacing(2), },
  card : { height : "100%",
    marginTop : theme.spacing(2), },
  cardContent : { height : "100%", },
  redirectButton : { marginLeft : "-.5em",
    color : "#000", },
  dashboardSection : {
    backgroundColor : "#fff",
    padding : theme.spacing(2),
    borderRadius : 4,
    height : "100%",
  },
});
class DashboardComponent extends React.Component {
  constructor(props) {
    super(props);
    const {
      meshAdapters, k8sconfig, grafana, prometheus
    } = props;
    this._isMounted = false;
    this.state = {
      meshAdapters,
      availableAdapters : [],
      mts : new Date(),
      meshLocationURLError : false,

      inClusterConfig : k8sconfig.inClusterConfig, // read from store
      k8sfile : k8sconfig.k8sfile, // read from store
      contextName : k8sconfig.contextName, // read from store

      clusterConfigured : k8sconfig.clusterConfigured, // read from store
      configuredServer : k8sconfig.configuredServer,
      grafanaUrl : grafana.grafanaURL,
      prometheusUrl : prometheus.prometheusURL,
      k8sfileError : false,
      kts : new Date(),

      grafana,
      prometheus,

      urlError : false,
      grafanaConfigSuccess : props.grafana.grafanaURL !== "",
      grafanaBoardSearch : "",
      grafanaURL : props.grafana.grafanaURL,
      grafanaAPIKey : props.grafana.grafanaAPIKey,
      grafanaBoards : props.grafana.grafanaBoards,
      selectedBoardsConfigs : props.grafana.selectedBoardsConfigs,
      ts : props.grafana.ts,

      meshScan : [],
      activeMeshScanNamespace : {},
      meshScanNamespaces : {},

      isMetricsConfigured : grafana.grafanaURL !== '' && prometheus.prometheusURL !== '' && k8sconfig.clusterConfigured,
      controlPlaneState : "",
      dataPlaneState : "",

      // subscriptions disposable
      dataPlaneSubscription : null,
      controlPlaneSubscription : null,
      operatorStatusSubscription : null,
    };
  }

  static getDerivedStateFromProps(props, state) {
    const {
      meshAdapters, meshAdaptersts, k8sconfig, grafana, prometheus
    } = props;
    const st = {};
    if (meshAdaptersts > state.mts) {
      st.meshAdapters = meshAdapters;
      st.mts = meshAdaptersts;
    }
    if (k8sconfig.ts > state.kts) {
      st.inClusterConfig = k8sconfig.inClusterConfig;
      st.k8sfile = k8sconfig.k8sfile;
      st.contextName = k8sconfig.contextName;
      st.clusterConfigured = k8sconfig.clusterConfigured;
      st.configuredServer = k8sconfig.configuredServer;
      st.kts = props.ts;
    }

    st.grafana = grafana;
    st.prometheus = prometheus;

    return st;
  }

  disposeSubscriptions = () => {
    if (this.state.operatorStatusSubscription) {
      this.state.operatorStatusSubscription.dispose()
    }
    if (this.state.dataPlaneSubscription) {
      this.state.dataPlaneSubscription.dispose()
    }
    if (this.state.controlPlaneSubscription) {
      this.state.controlPlaneSubscription.dispose()
    }
  }

  initMeshSyncControlPlaneSubscription = () => {
    /**
     * ALL_MESH indicates that we are interested in control plane
     * component of all of the service meshes supported by meshsync v2
     */
    const ALL_MESH = {};

    const self = this;

    if (self._isMounted){
      const opSub = subscribeOperatorStatusEvents(self.setOperatorState);
      // subscribeServiceMeshEvents(self.setMeshScanData, ALL_MESH, this.state, e => this.setState({ ...e }));
      const cpSub = subscribeControlPlaneEvents((res) => {
        if (res?.controlPlanesState !== undefined){
          this.setState({ controlPlaneState : res })
        }
      }, ALL_MESH)
      const dpSub = subscribeDataPlaneEvents((res) => {
        if (res?.dataPlanesState !== undefined){
          this.setState({ dataPlaneState : res })
        }
      }, ALL_MESH)
      this.disposeSubscriptions()
      this.setState({ operatorStatusSubscription : opSub, dataPlaneSubscription : dpSub, controlPlaneSubscription : cpSub })
      fetchControlPlanes(ALL_MESH).subscribe({
        next : (controlPlaneRes) => {
          this.setState({ controlPlaneState : controlPlaneRes })
          // self.setMeshScanData(controlPlaneRes, null);
          fetchDataPlanes(ALL_MESH).subscribe({
            next : (dataPlaneRes) => {
              this.setState({ dataPlaneState : dataPlaneRes })
              // if (controlPlaneRes) self.setMeshScanData(controlPlaneRes, dataPlaneRes);
            },
            error : (err) => console.error(err),
          });
        },
        error : (err) => console.error(err),
      });
    }

  }

  componentWillUnmount = () => {
    this._isMounted = false
    this.disposeSubscriptions()
  }

  componentDidMount = () => {
    this._isMounted = true
    this.fetchAvailableAdapters();

    if (this.state.isMetricsConfigured){
      this.fetchMetricComponents();
    }

    if (this._isMounted){
      this.initMeshSyncControlPlaneSubscription();
    }
  };

  componentDidUpdate(prevProps, prevState){
    let updateControlPlane = false;
    let updateDataPlane = false;

    // deep compare very limited, order of object fields is important
    if (JSON.stringify(prevState.controlPlaneState) !== JSON.stringify(this.state.controlPlaneState)){
      updateControlPlane = true
    }
    if (JSON.stringify(prevState.dataPlaneState) !== JSON.stringify(this.state.dataPlaneState)){
      updateDataPlane = true
    }

    if (updateDataPlane || updateControlPlane){
      this.setMeshScanData(
        updateControlPlane ? this.state.controlPlaneState : prevState.controlPlaneState,
        updateDataPlane ? this.state.dataPlaneState : prevState.dataPlaneState
      )
    }

  }

  fetchMetricComponents = () => {
    const self = this;

    dataFetch(
      "/api/telemetry/metrics/config",
      { method : "GET",
        credentials : "include",
        headers : { "Content-Type" : "application/x-www-form-urlencoded;charset=UTF-8", }, },
      (result) => {
        self.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined" && result?.prometheusURL && result?.prometheusURL != "") {
          let selector = { serviceMesh : "ALL_MESH", };
          fetchAvailableAddons(selector).subscribe({ next : (res) => {
            res?.addonsState?.forEach((addon) => {
              if (addon.name === "prometheus" && ( self.state.prometheusURL === "" || self.state.prometheusURL == undefined )) {
                self.setState({ prometheusURL : "http://" + addon.endpoint })
                submitPrometheusConfigure(self, () => console.log("Prometheus added"));
              }
            });
          },
          error : (err) => console.log("error registering prometheus: " + err), });
        }
      },
      self.handleError("Error getting prometheus config")
    );

    dataFetch(
      "/api/telemetry/metrics/grafana/config",
      { method : "GET",
        credentials : "include",
        headers : { "Content-Type" : "application/x-www-form-urlencoded;charset=UTF-8", }, },
      (result) => {
        self.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined" && result?.grafanaURL && result?.grafanaURL !="") {
          let selector = { serviceMesh : "ALL_MESH", };
          fetchAvailableAddons(selector).subscribe({ next : (res) => {
            res?.addonsState?.forEach((addon) => {
              if (addon.name === "grafana" && ( self.state.grafanaURL === "" || self.state.grafanaURL == undefined )) {
                self.setState({ grafanaURL : "http://" + addon.endpoint })
                submitGrafanaConfigure(self, () => {
                  self.state.selectedBoardsConfigs.push(self.state.boardConfigs);
                  console.log("Grafana added");
                });
              }
            });
          },
          error : (err) => console.log("error registering grafana: " + err), });
        }
      },
      self.handleError("There was an error communicating with grafana config")
    );

    let selector = { serviceMesh : "ALL_MESH", };

    fetchAvailableAddons(selector).subscribe({ next : (res) => {
      res?.addonsState?.forEach((addon) => {
        if (addon.name === "prometheus" && ( self.state.prometheusURL === "" || self.state.prometheusURL == undefined )) {
          self.setState({ prometheusURL : "http://" + addon.endpoint })
          submitPrometheusConfigure(self, () => console.log("Prometheus connected"));
        } else if (addon.name === "grafana" && ( self.state.grafanaURL === "" || self.state.grafanaURL == undefined )) {
          self.setState({ grafanaURL : "http://" + addon.endpoint })
          submitGrafanaConfigure(self, () => {
            self.state.selectedBoardsConfigs.push(self.state.boardConfigs);
            console.log("Grafana added");
          });
        }
      });
    },
    error : (err) => console.log("error registering addons: " + err), });
  };

  fetchAvailableAdapters = () => {
    const self = this;
    this.props.updateProgress({ showProgress : true });
    dataFetch(
      "/api/system/adapters",
      { credentials : "same-origin",
        method : "GET",
        credentials : "include", },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          const options = result.map((res) => ({ value : res.adapter_location,
            label : res.adapter_location, }));
          this.setState({ availableAdapters : options });
        }
      },
      self.handleError("Unable to fetch list of adapters.")
    );
  };

  setOperatorState = (res) => {
    const self = this;
    if (res.operator?.error) {
      self.handleError("Operator could not be reached")(res.operator?.error?.description);
      return false;
    }

    if (res.operator?.status === "ENABLED") {
      res.operator?.controllers?.forEach((controller) => {
        if (controller.name === "broker" && controller.status == "ENABLED") {
          self.setState({ NATSInstalled : true,
            NATSVersion : controller.version, });
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

    self.setState({
      operatorInstalled : false,
      NATSInstalled : false,
      meshSyncInstalled : false,
      operatorSwitch : false,
      operatorVersion : "N/A",
      meshSyncVersion : "N/A",
      NATSVersion : "N/A",
    });

    return false;
  };

  setMeshScanData = (controlPlanesData, dataPlanesData) => {
    const self = this;
    const namespaces = {};
    const activeNamespaces = {};
    const processedControlPlanesData = controlPlanesData?.controlPlanesState?.map((mesh) => {
      if (!mesh?.members?.length) {
        return;
      }
      let proxies = []

      if (Array.isArray(dataPlanesData?.dataPlanesState)){
        const dataplane = dataPlanesData.dataPlanesState.find(mesh_ => mesh_.name === mesh.name)

        if (Array.isArray(dataplane?.proxies)) proxies = dataplane.proxies
      }
      const processedMember = mesh?.members?.map((member) => {
        if (namespaces[mesh.name]) {
          namespaces[mesh.name].add(member.namespace);
        } else {
          namespaces[mesh.name] = new Set([member.namespace]);
        }

        // retrieve data planes according to mesh name
        if (proxies.length > 0){
          const controlPlaneMemberProxies = proxies.filter(proxy => proxy.controlPlaneMemberName === member.name)

          if (controlPlaneMemberProxies.length > 0){
            member = {
              ...member,
              data_planes : controlPlaneMemberProxies
            }
          }
        }

        return member
      });
      namespaces[mesh.name] = [...namespaces[mesh.name]];
      activeNamespaces[mesh.name] = namespaces[mesh.name][0] || "";

      return {
        ...mesh,
        members : processedMember
      }
    });

    self.setState({ meshScan : processedControlPlanesData?.filter(data => !!data).filter((data) => data.members?.length > 0) });
    self.setState({ meshScanNamespaces : namespaces, activeMeshScanNamespace : activeNamespaces });
  };

  /**
   * generateMeshScanPodName takes in the podname and the hash
   * and returns the trimmed pod name
   * @param {string} podname
   * @param {string} hash
   * @param {string | undefined} custom
   * @returns {{full, trimmed}}
   */
  generateMeshScanPodName = (podname, hash, custom) => {
    const str = custom || podname;
    return { full : podname,
      trimmed : str.substring(0, (hash ? str.indexOf(hash)
        : str.length) - 1), };
  };

  /**
   * generateMeshScanVersion takes in the string from which version
   * is to be extracted and returns the version. If the version string
   * is undefined then it returns "NA"
   * @param {string | undefined} versionStr is the string from which version is to be extracted
   * @returns {string}
   */
  generateMeshScanVersion = (versionStr) => {
    if (typeof versionStr !== "string") return "NA";

    const matchResult = versionStr.match(/\d+(\.\d+){2,}/g);
    if (!matchResult) return "NA";

    // Add "v" iff we have a valid match result
    return `v${matchResult[0]}`;
  };

  handleError = (msg) => (error) => {
    this.props.updateProgress({ showProgress : false });
    const self = this;
    this.props.enqueueSnackbar(`${msg}: ${error}`, { variant : "error", preventDuplicate : true,
      action : (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration : 7000, });
  };

  /**
   * redirectErrorToConsole returns a function which redirects
   * ther error to the console under the group labelled by the "msg"
   * param
   * @param {string} msg
   */
  redirectErrorToConsole = (msg) => (error) => {
    this.props.updateProgress({ showProgress : false });
    console.group(msg);
    console.error(error);
    console.groupEnd();
  };

  handleAdapterPingError = (msg) => () => {
    const { classes } = this.props;
    this.props.updateProgress({ showProgress : false });
    const self = this;
    this.props.enqueueSnackbar(`${msg}. To configure an adapter, visit`, { variant : "error",
      autoHideDuration : 3000,
      action : (key) => (
        <>
          <Button
            variant="contained"
            key="configure-close"
            aria-label="Configure"
            className={classes.redirectButton}
            onClick={() => {
              self.props.router.push("/settings#service-mesh");
              self.props.closeSnackbar(key);
            }}
          >
            <SettingsIcon className={classes.settingsIcon} />
            Settings
          </Button>

          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        </>
      ), });
  };

  handleDelete() {
    return false;
  }

  handleAdapterClick = (adapterLoc) => () => {
    // const { meshAdapters } = this.state;
    this.props.updateProgress({ showProgress : true });
    const self = this;
    dataFetch(
      `/api/system/adapters?adapter=${encodeURIComponent(adapterLoc)}`,
      { credentials : "same-origin",
        credentials : "include", },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Meshery Adapter connected at " + adapterLoc, { variant : "success",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ), });
        }
      },
      self.handleAdapterPingError("Could not connect to " + adapterLoc)
    );
  };

  handleConfigure = (val) => {
    this.props.router.push(`/settings#metrics/${val}`);
  };

  handleKubernetesClick = () => {
    this.props.updateProgress({ showProgress : true });
    const self = this;
    const { configuredServer } = this.state;
    dataFetch(
      "/api/system/kubernetes/ping",
      { credentials : "same-origin",
        credentials : "include", },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Kubernetes connected at "  + `${configuredServer}` , { variant : "success",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ), });
        }
      },
      self.handleError("Could not connect to Kubernetes")
    );
  };

  handleGrafanaClick = () => {
    this.props.updateProgress({ showProgress : true });
    const self = this;
    const { grafanaUrl } = this.state;
    dataFetch(
      "/api/telemetry/metrics/grafana/ping",
      { credentials : "same-origin",
        credentials : "include", },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Grafana connected at "  + `${grafanaUrl}`, { variant : "success",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ), });
        }
      },
      self.handleError("Could not connect to Grafana")
    );
  };

  /**
   * Meshcard takes in the mesh related data
   * and renders a table along with other information of
   * the mesh
   * @param {{name, icon, tag}} mesh
   * @param {{name, component, version, namespace}[]} components Array of components data
   */
  Meshcard = (mesh, components = []) => {
    const self = this;
    if (Array.isArray(components) && components.length)
      return (
        <Paper elevation={1} style={{ padding : "2rem", marginTop : "1rem" }}>
          <Grid container justify="space-between" spacing={1}>
            <Grid item>
              <div style={{ display : "flex", alignItems : "center", marginBottom : "1rem" }}>
                <img src={mesh.icon} className={this.props.classes.icon} style={{ marginRight : "0.75rem" }} />
                <Typography variant="h6">{mesh.tag}</Typography>
              </div>
            </Grid>
            <Grid item>
              {self.state.activeMeshScanNamespace[mesh.name] && (
                <Select
                  value={self.state.activeMeshScanNamespace[mesh.name]}
                  onChange={(e) =>
                    self.setState((state) => ({ activeMeshScanNamespace : { ...state.activeMeshScanNamespace, [mesh.name] : e.target.value }, }))
                  }
                >
                  {self.state.meshScanNamespaces[mesh.name] &&
                    self.state.meshScanNamespaces[mesh.name].map((ns) => <MenuItem value={ns}>{ns}</MenuItem>)}
                </Select>
              )}
            </Grid>
          </Grid>
          <TableContainer>
            <Table aria-label="Deployed service mesh details">
              <TableHead>
                <TableRow>
                  {/* <TableCell align="center">Control Plane</TableCell> */}
                  <TableCell align="center">Component</TableCell>
                  <TableCell align="center">Version</TableCell>
                  <TableCell align="center">Proxy</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {components
                  .filter((comp) => comp.namespace === self.state.activeMeshScanNamespace[mesh.name])
                  .map((component) => {
                    return (
                      <TableRow key={component.name.full}>
                        {/* <TableCell scope="row" align="center">
                          <Tooltip title={component.name.full}>
                            <div style={{ textAlign: "center" }}>
                              {component.name.trimmed}
                            </div>
                          </Tooltip>
                        </TableCell> */}
                        {/* <TableCell align="center">{podNameMapper(component.component, component.name)}</TableCell> */}
                        <TableCell align="center">{component.component}</TableCell>
                        <TableCell align="center">{versionMapper(component.version)}</TableCell>
                        <Tooltip
                          key={`component-${component.name}`}
                          title={
                            Array.isArray(component?.data_planes) && component.data_planes.length > 0 ? (
                              component.data_planes.map((cont) => {
                                return (
                                  <div key={cont.name} style={{ fontSize : "15px", color : '#fff', paddingBottom : '10px', padding : '1vh' } }>
                                    <p>Name: {cont?.containerName ? cont.containerName : 'Unspecified'}</p>
                                    <p>Status: {cont?.status?.ready ? 'ready' : 'not ready'}</p>
                                    {!cont?.status?.ready && (
                                      typeof cont?.status?.lastState === 'object' && cont?.status?.lastState !== null && Object.keys(cont.status.lastState).length > 0 && (
                                        <div>
                                          <p>Last state: {Object.keys(cont?.status?.lastState)[0]} <br/> Error: {Object.values(cont?.status?.lastState)[0]?.exitCode} <br/> Finished at: {Object.values(cont?.status?.lastState)[0]?.finishedAt}</p>
                                        </div>
                                      )
                                    )}
                                    {typeof cont?.status?.state === 'object' && cont?.status?.state !== null && Object.keys(cont.status.state).length > 0 && (
                                      <p>State: {Object.keys(cont.status.state)[0]}</p>
                                    )}
                                    {cont?.status?.restartCount && (
                                      <p>Restart count: {cont?.status.restartCount}</p>
                                    )}
                                    <p>Image: {cont.image}</p>
                                    <p>Ports: <br/> {cont?.ports && cont.ports.map(port => `[ ${port?.name ? port.name : 'Unknown'}, ${port?.containerPort ? port.containerPort : 'Unknown'}, ${port?.protocol ? port.protocol : 'Unknown'} ]`).join(', ')}</p>
                                    {cont?.resources && (
                                      <div>
                                        Resources used: <br/>

                                        <div style={{ paddingLeft : '2vh' }}>
                                          {cont?.resources?.limits && (
                                            <div>
                                              <p>Limits: <br/>
                                              CPU: {cont?.resources?.limits?.cpu} - Memory: {cont?.resources?.limits?.memory}</p>
                                            </div>
                                          )}
                                          {cont?.resources?.requests && (
                                            <div>
                                              <p>Requests: <br/>
                                              CPU: {cont?.resources?.requests?.cpu} - Memory: {cont?.resources?.requests?.memory}</p>
                                            </div>
                                          )}
                                        </div>
                                      </div>
                                    )}
                                  </div>
                                )
                              })
                            ) : "No proxy attached"}
                        >
                          <TableCell align="center">{component?.data_planes?.length || 0}</TableCell>
                        </Tooltip>
                      </TableRow>
                    )
                  })}
              </TableBody>
            </Table>
          </TableContainer>
        </Paper>
      );

    return null;
  };

  handlePrometheusClick = () => {
    this.props.updateProgress({ showProgress : true });
    const self = this;
    const { prometheusUrl } = this.state;
    dataFetch(
      "/api/telemetry/metrics/ping",
      { credentials : "same-origin",
        credentials : "include", },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Prometheus connected at" + ` ${prometheusUrl}`, { variant : "success",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ), });
        }
      },
      self.handleError("Could not connect to Prometheus")
    );
  };

  showCard(title, content) {
    const { classes } = this.props;
    return (
      <Card className={classes.card}>
        <CardHeader
          disableTypography
          title={title}
          // action={iconComponent}
          className={classes.cardHeader}
        />
        <CardContent className={classes.cardContent}>{content}</CardContent>
      </Card>
    );
  }

  configureTemplate = () => {
    const { classes } = this.props;
    const {
      inClusterConfig,
      contextName,
      clusterConfigured,
      configuredServer,
      meshAdapters,
      grafanaUrl,
      prometheusUrl,
      availableAdapters,
      grafana,
      prometheus,
    } = this.state;
    const self = this;
    let showConfigured = "Not connected to Kubernetes.";
    if (clusterConfigured) {
      let chp = (
        <Chip
          label={inClusterConfig
            ? "Using In Cluster Config"
            : contextName}
          onClick={self.handleKubernetesClick}
          icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />}
          className={classes.chip}
          key="k8s-key"
          variant="outlined"
        />
      );

      if (configuredServer) {
        chp = <Tooltip title={`Server: ${configuredServer}`}>{chp}</Tooltip>;
      } else {
        chp=showConfigured;
      }

      showConfigured = <div showConfigured>{chp}</div>;
    }


    let showAdapters = "No adapters configured.";
    if (availableAdapters.length > 0) {
      availableAdapters.sort((a1, a2) => (a1.value < a2.value
        ? -1
        : a1.value > a2.value
          ? 1
          : 0));

      showAdapters = (
        <div>
          {availableAdapters.map((aa, ia) => {
            let isDisabled = true;
            let image = "/static/img/meshery-logo.png";
            let logoIcon = <img src={image} className={classes.icon} />;
            let adapterType = "";
            let adapterVersion = "";
            meshAdapters.forEach((adapter) => {
              if (aa.value === adapter.adapter_location) {
                isDisabled = false;
                adapterType = adapter.name;
                adapterVersion = adapter.version;
                image = "/static/img/" + adapter.name.toLowerCase() + ".svg";
                logoIcon = <img src={image} className={classes.icon} />;
              }
            });

            return (
              <Tooltip
                key={`adapters-${ia}`}
                title={
                  isDisabled
                    ? "Inactive Meshery Adapter"
                    : `Meshery Adapter for 
                      ${adapterType
                      .toLowerCase()
                      .split(" ")
                      .map((s) => s.charAt(0).toUpperCase() + s.substring(1))
                      .join(" ")} on ${aa.label.split(":")[1] }/tcp (${adapterVersion})`}
              >
                <Chip
                  label={
                    isDisabled
                      ? aa.label.split(":")[0] + ":" + aa.label.split(":")[1]
                      : adapterType.toLowerCase()
                        .split("_")
                        .map((s) => s.charAt(0).toUpperCase() + s.substring(1) + " ")}
                  onClick={self.handleAdapterClick(aa.value)}
                  icon={logoIcon}
                  className={classes.chip}
                  key={`adapters-${ia}`}
                  variant={isDisabled
                    ? "default"
                    : "outlined"}
                />
              </Tooltip>
            );
          })}
        </div>
      );
    }
    let showGrafana;
    if (grafanaUrl === "") {
      showGrafana = (
        <div className={classes.alreadyConfigured}>
          <Button
            variant="contained"
            color="primary"
            size="large"
            className={classes.metricsButton}
            onClick={() => this.handleConfigure("grafana")}
          >
            <SettingsIcon className={classes.settingsIcon} />
            Configure Grafana
          </Button>
        </div>
      );
    }
    if (grafana && grafana.grafanaURL && grafana.grafanaURL !== "") {
      showGrafana = (
        <Chip
          label={grafana.grafanaURL}
          onClick={self.handleGrafanaClick}
          icon={<img src="/static/img/grafana_icon.svg" className={classes.icon} />}
          className={classes.chip}
          key="graf-key"
          variant="outlined"
        />
      );
    }

    let showPrometheus;
    if (prometheusUrl === "") {
      showPrometheus = (
        <div className={classes.alreadyConfigured}>
          <Button
            variant="contained"
            color="primary"
            size="large"
            className={classes.metricsButton}
            onClick={() => this.handleConfigure("prometheus")}
          >
            <SettingsIcon className={classes.settingsIcon} />
            Configure Prometheus
          </Button>
        </div>
      );
    }
    if (prometheus && prometheus.prometheusURL && prometheus.prometheusURL !== "") {
      showPrometheus = (
        <Chip
          label={prometheus.prometheusURL}
          onClick={self.handlePrometheusClick}
          icon={<img src="/static/img/prometheus_logo_orange_circle.svg" className={classes.icon} />}
          className={classes.chip}
          key="prom-key"
          variant="outlined"
        />
      );
    }

    const showMetrics = (
      <Grid container justify="center" spacing={2}>
        <Grid item>{showPrometheus}</Grid>
        <Grid item>{showGrafana}</Grid>
        {/*<Grid item>
          <Paper className={classes.paper}>
            <MesheryMetrics
              boardConfigs={grafana.selectedBoardsConfigs}
              grafanaURL={grafana.grafanaURL}
              grafanaAPIKey={grafana.grafanaAPIKey}
              handleGrafanaChartAddition={() => router.push("/settings/#metrics")}
            />
          </Paper>
        </Grid>*/}
      </Grid>
    );

    const showServiceMesh = (
      <>
        {self?.state?.meshScan && Object.keys(self.state.meshScan).length
          ? (
            <>
              {self.state.meshScan.map((mesh) => {
                let tag = "";
                mesh.name
                  .split("_")
                  .forEach((element) => {
                    tag = tag + " " + element[0].toUpperCase() + element.slice(1, element.length);
                  });
                return self.Meshcard(
                  { name : mesh.name, tag : tag, icon : "/static/img/" + mesh.name + ".svg" },
                  mesh.members
                );
              })}
            </>
          )
          : (
            <div
              style={{
                padding : "2rem",
                display : "flex",
                justifyContent : "center",
                alignItems : "center",
                flexDirection : "column",
              }}
            >
              <Typography style={{ fontSize : "1.5rem", marginBottom : "2rem" }} align="center" color="textSecondary">
              No service meshes detected in the {self.state.contextName} cluster.
              </Typography>
              <Button
                aria-label="Add Meshes"
                variant="contained"
                color="primary"
                size="large"
                onClick={() => self.props.router.push("/management")}
              >
                <AddIcon className={classes.addIcon} />
              Install Service Mesh
              </Button>
            </div>
          )}
      </>
    );
    return (
      <NoSsr>
        <div className={classes.root}>
          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <div className={classes.dashboardSection} data-test="service-mesh">
                <Typography variant="h6" gutterBottom className={classes.chartTitle}>
                  Service Mesh
                </Typography>
                {showServiceMesh}
              </div>
            </Grid>
            <Grid item xs={12} md={6}>
              <div className={classes.dashboardSection} data-test="connection-status">
                <Typography variant="h6" gutterBottom className={classes.chartTitle}>
                  Connection Status
                </Typography>
                <div>{self.showCard("Kubernetes", showConfigured)}</div>
                <div>{self.showCard("Adapters", showAdapters)}</div>
                <div>{self.showCard("Metrics", showMetrics)}</div>
              </div>
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

DashboardComponent.propTypes = { classes : PropTypes.object.isRequired, };

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch),
  updateGrafanaConfig : bindActionCreators(updateGrafanaConfig, dispatch),
  updatePrometheusConfig : bindActionCreators(updatePrometheusConfig, dispatch), });
const mapStateToProps = (state) => {
  const k8sconfig = state.get("k8sConfig").toJS();
  const meshAdapters = state.get("meshAdapters").toJS();
  const meshAdaptersts = state.get("meshAdaptersts");
  const grafana = state.get("grafana").toJS();
  const prometheus = state.get("prometheus").toJS();
  return {
    meshAdapters,
    meshAdaptersts,
    k8sconfig,
    grafana,
    prometheus,
  };
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(DashboardComponent)))
);