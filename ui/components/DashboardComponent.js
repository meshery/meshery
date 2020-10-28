import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import RadioGroup from "@material-ui/core/RadioGroup";
import Radio from "@material-ui/core/Radio";
import FormControlLabel from "@material-ui/core/FormControlLabel";
import FormControl from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";
import { NoSsr, Chip, IconButton, Button, Card, CardContent, Typography, CardHeader, Tooltip } from "@material-ui/core";
import blue from "@material-ui/core/colors/blue";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import SettingsIcon from "@material-ui/icons/Settings";
import { withRouter } from "next/router";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import { updateProgress } from "../lib/store";
import dataFetch from "../lib/data-fetch";

const styles = (theme) => ({
  root: {
    padding: theme.spacing(5),
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
    padding: theme.spacing(2),
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
  cardHeader: {
    fontSize: theme.spacing(2),
  },
  card: {
    height: "100%",
  },
  cardContent: {
    height: "100%",
  },
  redirectButton: {
    marginLeft: "-.5em",
    color: "#000",
  },
});

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

      versionDetail: { build: "", latest: "", outdated: false },
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

  componentDidMount = () => {
    this.fetchAvailableAdapters();
    this.fetchVersionDetails();
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
    const self = this
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
          this.setState({versionDetail: {
            build: "Unknown",
            latest: "Unknown",
            outdated: false
          }});
        }
      },
      self.handleError("Unable to fetch meshery version.")
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

      showConfigured = <div className={classes.alignRight}>{chp}</div>;
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
                  case "open service mesh":
                    image = "/static/img/osm.svg";
                    logoIcon = <img src={image} className={classes.icon} />;
                    break;
                  case "kuma":
                    image = "/static/img/kuma.svg";
                    logoIcon = <img src={image} className={classes.icon} />;
                    break;
                  case "nginx service mesh":
                    image = "/static/img/nginx-sm.svg";
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
      <div>
        {showPrometheus}
        {showGrafana}
      </div>
    );

    /**
     * getMesheryVersionText returs a well formatted version text
     * @param {string} type type of version could be "latest" or "current"
     */
    const getMesheryVersionText = (type) => {
      const {build, latest, outdated} = this.state.versionDetail
      if (type === "current") {
        if (outdated) return `Running Meshery: stable-${build}`
        if (build === "Unknown") return "Unknown"

        return `Meshery is up to date: stable-${build}`
      }

      if (type === "latest") {
        if (outdated) return `Latest Available: stable-${latest}`
      }

      return ``
    };

    const showRelease = (
      <Grid container justify="space-between" spacing={1}>
        <Grid item xs={12} md={6}>
          <FormControl component="fieldset">
            <FormLabel component="legend" style={{ fontWeight: "bold" }}>
            Channel
            </FormLabel>
            <RadioGroup aria-label="release_channel_option" name="release_channel">
              <FormControlLabel value="stable_channel" disabled control={<Radio checked={true} />} label="Stable Channel" />
              <FormControlLabel value="edge_channel" disabled control={<Radio />} label="Edge Channel" />
            </RadioGroup>
          </FormControl>
        </Grid>
        <Grid item xs={12} md={6} style={{padding: "0"}}>
          <Typography style={{fontWeight: "bold", paddingBottom: "4px"}}>Version</Typography>
          <Typography style={{paddingTop: "2px", paddingBottom: "8px"}}>
            {getMesheryVersionText("current")}
          </Typography>
          <Typography style={{paddingTop: "8px"}}>
            {getMesheryVersionText("latest")}
          </Typography>
        </Grid>
      </ Grid>
    );

    return (
      <NoSsr>
        <div className={classes.root}>
          <Typography variant="h6" gutterBottom className={classes.chartTitle}>
            Connection Status
          </Typography>

          <Grid container spacing={1}>
            <Grid item xs={12} md={6}>
              {self.showCard("Kubernetes", showConfigured)}
            </Grid>
            <Grid item xs={12} md={6}>
              {self.showCard("Adapters", showAdapters)}
            </Grid>
            <Grid item xs={12} md={6}>
              {self.showCard("Metrics", showMetrics)}
            </Grid>
            <Grid item xs={12} md={6}>
              {self.showCard("Release", showRelease)}
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
