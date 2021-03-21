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
  Link,
  Box,
} from "@material-ui/core";
import blue from "@material-ui/core/colors/blue";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import SettingsIcon from "@material-ui/icons/Settings";
import AddIcon from "@material-ui/icons/AddCircleOutline";
import { withRouter } from "next/router";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import { updateProgress } from "../lib/store";
import dataFetch from "../lib/data-fetch";
import subscribeControlPlaneEvents from "./graphql/subscriptions/ControlPlaneSubscription";
import subscribeOperatorStatusEvents from "./graphql/subscriptions/OperatorStatusSubscription";
import subscribeMeshSyncStatusEvents from "./graphql/subscriptions/MeshSyncStatusSubscription";
import fetchControlPlanes from "./graphql/queries/ControlPlanesQuery";

const styles = (theme) => ({
  root: {
    backgroundColor: "#eaeff1",
  },
  chip: {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
  buttons: {
    display: "flex",
    justifyContent: "flex-end",
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  metricsButton: {
    width: "240px",
  },
  alreadyConfigured: {
    textAlign: "center",
  },
  margin: {
    margin: theme.spacing(1),
  },
  colorSwitchBase: {
    color: blue[300],
    "&$colorChecked": {
      color: blue[500],
      "& + $colorBar": {
        backgroundColor: blue[500],
      },
    },
  },
  colorBar: {},
  colorChecked: {},
  fileLabel: {
    width: "100%",
  },
  fileLabelText: {},
  inClusterLabel: {
    paddingRight: theme.spacing(2),
  },
  alignCenter: {
    textAlign: "center",
  },
  icon: {
    width: theme.spacing(2.5),
  },
  istioIcon: {
    width: theme.spacing(1.5),
  },
  settingsIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  addIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  cardHeader: {
    fontSize: theme.spacing(2),
  },
  card: {
    height: "100%",
    marginTop: theme.spacing(2),
  },
  cardContent: {
    height: "100%",
  },
  redirectButton: {
    marginLeft: "-.5em",
    color: "#000",
  },
  dashboardSection: {
    backgroundColor: "#fff",
    padding: theme.spacing(2),
    borderRadius: 4,
    height: "100%",
  },
});

/**
 * capitalize takes in a string and returns
 * capitalized string
 * @param {string} str - string to be capitalized
 */
function capitalize(str) {
  return `${str?.charAt(0).toUpperCase()}${str?.substring(1)}`;
}

class DashboardComponent extends React.Component {
  constructor(props) {
    super(props);
    const { meshAdapters, k8sconfig, grafana, prometheus } = props;
    this.state = {
      meshAdapters,
      availableAdapters: [],
      mts: new Date(),
      meshLocationURLError: false,

      inClusterConfig: k8sconfig.inClusterConfig, // read from store
      k8sfile: k8sconfig.k8sfile, // read from store
      contextName: k8sconfig.contextName, // read from store

      clusterConfigured: k8sconfig.clusterConfigured, // read from store
      configuredServer: k8sconfig.configuredServer,
      grafanaUrl: grafana.grafanaURL,
      prometheusUrl: prometheus.prometheusURL,
      k8sfileError: false,
      kts: new Date(),

      grafana,
      prometheus,

      versionDetail: {
        build: "",
        latest: "",
        outdated: false,
        commitsha: "",
        release_channel: "NA",
      },

      meshScan: {},
      activeMeshScanNamespace: {},
      meshScanNamespaces: {},
    };
  }

  static getDerivedStateFromProps(props, state) {
    const { meshAdapters, meshAdaptersts, k8sconfig, grafana, prometheus } = props;
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

  initMeshSyncControlPlaneSubscription = () => {
    /**
     * ALL_MESH indicates that we are interested in control plane
     * component of all of the service meshes supported by meshsync v2
     */
    const ALL_MESH = {};

    const self = this;
    subscribeMeshSyncStatusEvents((res) => {
      if (res.meshsync?.error) {
        self.handleError(res.meshsync?.error?.description || "MeshSync could not be reached");
        return;
      }
    });
    subscribeOperatorStatusEvents(self.setOperatorState);
    subscribeControlPlaneEvents(self.setMeshScanData, ALL_MESH);

    fetchControlPlanes(ALL_MESH).subscribe({
      next: (res) => {
        self.setMeshScanData(res);
      },
      error: (err) => console.error(err),
    });
  };

  componentDidMount = () => {
    this.fetchAvailableAdapters();
    this.fetchVersionDetails();
    this.initMeshSyncControlPlaneSubscription();
  };

  fetchAvailableAdapters = () => {
    const self = this;
    this.props.updateProgress({ showProgress: true });
    dataFetch(
      "/api/mesh/adapters",
      {
        credentials: "same-origin",
        method: "GET",
        credentials: "include",
      },
      (result) => {
        this.props.updateProgress({ showProgress: false });
        if (typeof result !== "undefined") {
          const options = result.map((res) => ({
            value: res,
            label: res,
          }));
          this.setState({ availableAdapters: options });
        }
      },
      self.handleError("Unable to fetch list of adapters.")
    );
  };

  fetchVersionDetails = () => {
    const self = this;
    this.props.updateProgress({ showProgress: true });
    dataFetch(
      "/api/server/version",
      {
        credentials: "same-origin",
        method: "GET",
        credentials: "include",
      },
      (result) => {
        this.props.updateProgress({ showProgress: false });
        if (typeof result !== "undefined") {
          this.setState({ versionDetail: result });
        } else {
          this.setState({
            versionDetail: {
              build: "Unknown",
              latest: "Unknown",
              outdated: false,
              commitsha: "Unknown",
            },
          });
        }
      },
      self.handleError("Unable to fetch Meshery version.")
    );
  };

  setMeshScanData = (data) => {
    const self = this;
    const namespaces = {};
    const activeNamespaces = {};

    data?.controlPlanesState?.map((mesh) => {
      if (!mesh?.members?.length) {
        return;
      }
      mesh?.members?.map((member) => {
        if (namespaces[mesh.name]) {
          namespaces[mesh.name].add(member.namespace);
        } else {
          namespaces[mesh.name] = new Set([member.namespace]);
        }
      });

      namespaces[mesh.name] = [...namespaces[mesh.name]];
      activeNamespaces[mesh.name] = namespaces[mesh.name][0] || "";

      const meshData = mesh?.members?.map((member) => ({
        name: member.name,
        component: member.component,
        version: member.version,
        namespace: member.namespace,
      }));
      self.setState((state) => ({ meshScan: { ...state.meshScan, [mesh.name]: meshData } }));
    });

    self.setState({ meshScanNamespaces: namespaces, activeMeshScanNamespace: activeNamespaces });
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
    return {
      full: podname,
      trimmed: str.substring(0, (hash ? str.indexOf(hash) : str.length) - 1),
    };
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

  /**
   * redirectErrorToConsole returns a function which redirects
   * ther error to the console under the group labelled by the "msg"
   * param
   * @param {string} msg
   */
  redirectErrorToConsole = (msg) => (error) => {
    this.props.updateProgress({ showProgress: false });
    console.group(msg);
    console.error(error);
    console.groupEnd();
  };

  handleAdapterPingError = (msg) => () => {
    const { classes } = this.props;
    this.props.updateProgress({ showProgress: false });
    const self = this;
    this.props.enqueueSnackbar(`${msg}. To configure an adapter, visit`, {
      variant: "error",
      autoHideDuration: 2000,
      action: (key) => (
        <>
          <Button
            key="configure-close"
            aria-label="Configure"
            className={classes.redirectButton}
            onClick={() => {
              self.props.router.push("/settings#service-mesh");
              self.props.closeSnackbar(key);
            }}
          >
            Settings
          </Button>

          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        </>
      ),
    });
  };

  handleDelete() {
    return false;
  }

  handleAdapterClick = (adapterLoc) => () => {
    // const { meshAdapters } = this.state;
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch(
      `/api/mesh/adapter/ping?adapter=${encodeURIComponent(adapterLoc)}`,
      {
        credentials: "same-origin",
        credentials: "include",
      },
      (result) => {
        this.props.updateProgress({ showProgress: false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Adapter successfully pinged!", {
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
      self.handleAdapterPingError("Could not ping adapter.")
    );
  };

  handleConfigure = (val) => {
    this.props.router.push(`/settings#metrics/${val}`);
  };

  handleKubernetesClick = () => {
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch(
      "/api/k8sconfig/ping",
      {
        credentials: "same-origin",
        credentials: "include",
      },
      (result) => {
        this.props.updateProgress({ showProgress: false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Kubernetes successfully pinged!", {
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
      self.handleError("Could not ping Kubernetes.")
    );
  };

  handleGrafanaClick = () => {
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch(
      "/api/grafana/ping",
      {
        credentials: "same-origin",
        credentials: "include",
      },
      (result) => {
        this.props.updateProgress({ showProgress: false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Grafana successfully pinged!", {
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
      self.handleError("Could not ping Grafana.")
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
        <Paper elevation={1} style={{ padding: "2rem", marginTop: "1rem" }}>
          <Grid container justify="space-between" spacing={1}>
            <Grid item>
              <div style={{ display: "flex", alignItems: "center", marginBottom: "1rem" }}>
                <img src={mesh.icon} className={this.props.classes.icon} style={{ marginRight: "0.75rem" }} />
                <Typography variant="h6">{mesh.tag}</Typography>
              </div>
            </Grid>
            <Grid item>
              {self.state.activeMeshScanNamespace[mesh.name] && (
                <Select
                  value={self.state.activeMeshScanNamespace[mesh.name]}
                  onChange={(e) =>
                    self.setState((state) => ({
                      activeMeshScanNamespace: { ...state.activeMeshScanNamespace, [mesh.name]: e.target.value },
                    }))
                  }
                >
                  {self.state.meshScanNamespaces[mesh.name] &&
                    self.state.meshScanNamespaces[mesh.name].map((ns) => <MenuItem value={ns}>{ns}</MenuItem>)}
                </Select>
              )}
            </Grid>
          </Grid>
          <TableContainer>
            <Table aria-label="mesh details table">
              <TableHead>
                <TableRow>
                  <TableCell align="center">Control Plane Pods</TableCell>
                  <TableCell align="center">Component</TableCell>
                  <TableCell align="center">Version</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {components
                  .filter((comp) => comp.namespace === self.state.activeMeshScanNamespace[mesh.name])
                  .map((component) => (
                    <TableRow key={component.name.full}>
                      {/* <TableCell scope="row" align="center">
                        <Tooltip title={component.name.full}>
                          <div style={{ textAlign: "center" }}>
                            {component.name.trimmed}
                          </div>
                        </Tooltip>
                      </TableCell> */}
                      <TableCell align="center">{component.name}</TableCell>
                      <TableCell align="center">{component.component}</TableCell>
                      <TableCell align="center">{component.version}</TableCell>
                    </TableRow>
                  ))}
              </TableBody>
            </Table>
          </TableContainer>
        </Paper>
      );

    return null;
  };

  handlePrometheusClick = () => {
    this.props.updateProgress({ showProgress: true });
    const self = this;
    dataFetch(
      "/api/prometheus/ping",
      {
        credentials: "same-origin",
        credentials: "include",
      },
      (result) => {
        this.props.updateProgress({ showProgress: false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Prometheus successfully pinged!", {
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
      self.handleError("Could not ping Prometheus.")
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
          label={inClusterConfig ? "Using In Cluster Config" : contextName}
          onClick={self.handleKubernetesClick}
          icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />}
          className={classes.chip}
          key="k8s-key"
          variant="outlined"
        />
      );

      if (configuredServer) {
        chp = <Tooltip title={`Server: ${configuredServer}`}>{chp}</Tooltip>;
      }

      showConfigured = <div showConfigured>{chp}</div>;
    }

    let showAdapters = "No adapters configured.";
    if (availableAdapters.length > 0) {
      availableAdapters.sort((a1, a2) => (a1.value < a2.value ? -1 : a1.value > a2.value ? 1 : 0));

      showAdapters = (
        <div>
          {availableAdapters.map((aa, ia) => {
            let isDisabled = true;
            let image = "/static/img/meshery-logo.png";
            let logoIcon = <img src={image} className={classes.icon} />;
            let adapterType = "";
            meshAdapters.forEach((adapter) => {
              if (aa.value === adapter.adapter_location) {
                isDisabled = false;
                adapterType = adapter.name;
                switch (adapter.name.toLowerCase()) {
                  case "istio":
                    image = "/static/img/istio.svg";
                    logoIcon = <img src={image} className={classes.istioIcon} />;
                    break;
                  case "linkerd":
                    image = "/static/img/linkerd.svg";
                    logoIcon = <img src={image} className={classes.icon} />;
                    break;
                  case "consul":
                    image = "/static/img/consul.svg";
                    logoIcon = <img src={image} className={classes.icon} />;
                    break;
                  case "network service mesh":
                    image = "/static/img/nsm.svg";
                    logoIcon = <img src={image} className={classes.icon} />;
                    break;
                  case "octarine":
                    image = "/static/img/octarine.svg";
                    logoIcon = <img src={image} className={classes.icon} />;
                    break;
                  case "citrix service mesh":
                    image = "/static/img/citrix.svg";
                    logoIcon = <img src={image} className={classes.icon} />;
                    break;
                  case "osm":
                    image = "/static/img/osm.svg";
                    logoIcon = <img src={image} className={classes.icon} />;
                    break;
                  case "kuma":
                    image = "/static/img/kuma.svg";
                    logoIcon = <img src={image} className={classes.icon} />;
                    break;
                  // Disable support for NGINX SM
                  // case "nginx service mesh":
                  //   image = "/static/img/nginx-sm.svg";
                  //   logoIcon = <img src={image} className={classes.icon} />;
                  //   break;
                  case "traefik mesh":
                    image = "/static/img/traefikmesh.svg";
                    logoIcon = <img src={image} className={classes.icon} />;
                    break;
                }
              }
            });

            return (
              <Tooltip
                key={`adapters-${ia}`}
                title={
                  isDisabled
                    ? "This adapter is inactive"
                    : `${adapterType
                      .toLowerCase()
                      .split(" ")
                      .map((s) => s.charAt(0).toUpperCase() + s.substring(1))
                      .join(" ")} adapter on port ${aa.label.split(":")[1]}`
                }
              >
                <Chip
                  label={aa.label.split(":")[0]}
                  onClick={self.handleAdapterClick(aa.value)}
                  icon={logoIcon}
                  className={classes.chip}
                  key={`adapters-${ia}`}
                  variant={isDisabled ? "default" : "outlined"}
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
      </Grid>
    );

    const showServiceMesh = (
      <>
        {Object.keys(self.state.meshScan).length ? (
          <>
            {self.Meshcard({ name: "consul", tag: "Consul", icon: "/static/img/consul.svg" }, self.state.meshScan.consul)}
            {self.Meshcard({ name: "istio", tag: "Istio", icon: "/static/img/istio.svg" }, self.state.meshScan.istio)}
            {self.Meshcard({ name: "linkerd", tag: "Linkerd", icon: "/static/img/linkerd.svg" }, self.state.meshScan.linkerd)}
            {self.Meshcard({ name: "osm", tag: "Open Service Mesh", icon: "/static/img/osm.svg" }, self.state.meshScan.osm)}
            {self.Meshcard({ name: "osm", tag: "Network Service Mesh", icon: "/static/img/nsm.svg" }, self.state.meshScan.nsm)}
            {self.Meshcard({ name: "octarine", tag: "Octarine", icon: "/static/img/octarine.svg" }, self.state.meshScan.octarine)}
            {self.Meshcard({ name: "traefikmesh", tag: "Traefik Mesh", icon: "/static/img/traefikmesh.svg" }, self.state.meshScan.traefikmesh)}
            {self.Meshcard({ name: "kuma", tag: "Kuma", icon: "/static/img/kuma.svg" }, self.state.meshScan.kuma)}
            {/**self.Meshcard({ name: "nginx-sm", tag: "Nginx Service Mesh", icon: "/static/img/nginx-sm.svg" }, self.state.meshScan.nginx-sm) */}
            {self.Meshcard({ name: "citrix", tag: "Citrix", icon: "/static/img/citrix.svg" }, self.state.meshScan.citrix)}
          </>
        ) : (
          <div
            style={{
              padding: "2rem",
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              flexDirection: "column",
            }}
          >
            <Typography style={{ fontSize: "1.5rem", marginBottom: "2rem" }} align="center" color="textSecondary">
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

    /**
     * getMesheryVersionText returs a well formatted version text
     *
     * If the meshery is running latest version then and is using "edge" channel
     * then it will just show "edge-latest". However, if the meshery is on edge and
     * is running an outdated version then it will return "edge-$version".
     *
     * If on stable channel, then it will always show "stable-$version"
     */
    const getMesheryVersionText = () => {
      const { build, outdated, release_channel } = this.state.versionDetail;

      // If the version is outdated then no matter what the
      // release channel is, specify the build
      if (outdated) return `${release_channel}-${build}`;

      if (release_channel === "edge") return `${release_channel}-latest`;
      if (release_channel === "stable") return `${release_channel}-${build}`;

      return ``;
    };

    /**
     * versionUpdateMsg returns the appropriate message
     * based on the meshery's current running version and latest available
     * version.
     *
     * @returns {React.ReactNode} react component to display
     */
    const versionUpdateMsg = () => {
      const { outdated, latest } = this.state.versionDetail;

      if (outdated)
        return (
          <>
            Newer version of Meshery available:{" "}
            <Link href={`https://docs.meshery.io/project/releases/${latest}`}>{`${latest}`}</Link>
          </>
        );

      return <>Running latest Meshery version.</>;
    };

    const showRelease = (
      <>
        <Grid container justify="space-between" spacing={1}>
          <Grid item xs={12} md={6}>
            <Typography style={{ fontWeight: "bold", paddingBottom: "4px" }}>Channel Subscription</Typography>
            <Typography style={{ paddingTop: "2px", paddingBottom: "8px" }}>
              {capitalize(this.state.versionDetail.release_channel)}
            </Typography>
          </Grid>
          <Grid item xs={12} md={6} style={{ padding: "0" }}>
            <Typography style={{ fontWeight: "bold", paddingBottom: "4px" }}>Version</Typography>
            <Typography style={{ paddingTop: "2px", paddingBottom: "8px" }}>{getMesheryVersionText()}</Typography>
          </Grid>
        </Grid>
        <Typography component="div" style={{ marginTop: "1.5rem" }}>
          <Box fontStyle="italic" fontSize={14}>
            {versionUpdateMsg()}
          </Box>
        </Typography>
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
                <div>{self.showCard("Release", showRelease)}</div>
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

DashboardComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});
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
