import React from "react";
import { connect } from "react-redux";
import NoSsr from "@material-ui/core/NoSsr";
import { withStyles, Button, Divider, MenuItem, TextField, Grid } from "@material-ui/core";
import { blue } from "@material-ui/core/colors";
import PropTypes from "prop-types";
import { withRouter } from "next/router";
import SettingsIcon from "@material-ui/icons/Settings";
import MesheryAdapterPlayComponent from "./MesheryAdapterPlayComponent";

const styles = (theme) => ({
  icon: {
    fontSize: 20,
  },
  root: {
    padding: theme.spacing(0),
    marginBottom: theme.spacing(2),
  },
  buttons: {
    display: "flex",
    justifyContent: "flex-end",
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  margin: {
    margin: theme.spacing(1),
  },
  alreadyConfigured: {
    textAlign: "center",
    padding: theme.spacing(20),
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
  uploadButton: {
    margin: theme.spacing(1),
    marginTop: theme.spacing(3),
  },
  fileLabel: {
    width: "100%",
  },
  editorContainer: {
    width: "100%",
  },
  deleteLabel: {
    paddingRight: theme.spacing(2),
  },
  alignRight: {
    textAlign: "right",
  },
  expTitleIcon: {
    width: theme.spacing(3),
    display: "inline",
    verticalAlign: "middle",
  },
  expIstioTitleIcon: {
    width: theme.spacing(2),
    display: "inline",
    verticalAlign: "middle",
    marginLeft: theme.spacing(0.5),
    marginRight: theme.spacing(0.5),
  },
  expTitle: {
    display: "inline",
    verticalAlign: "middle",
    marginLeft: theme.spacing(1),
  },
  paneSection: {
    backgroundColor: "#fff",
    padding: theme.spacing(2.5),
    borderRadius: 4,
  },
});

class MesheryPlayComponent extends React.Component {
  constructor(props) {
    super(props);

    const { k8sconfig, meshAdapters } = props;
    let adapter = {};
    if (meshAdapters && meshAdapters.length > 0) {
      adapter = meshAdapters[0];
    }
    this.state = {
      k8sconfig,
      kts: new Date(),

      meshAdapters,
      mts: new Date(),

      adapter,
    };
  }

  static getDerivedStateFromProps(props, state) {
    let { meshAdapters, meshAdaptersts, k8sconfig } = props;
    const st = {};
    if (meshAdaptersts > state.mts) {
      st.meshAdapters = meshAdapters;
      st.mts = meshAdaptersts;
      if (meshAdapters && meshAdapters.length > 0) {
        st.adapter = meshAdapters[0];
      }
    }
    if (k8sconfig.ts > state.kts) {
      st.inClusterConfig = k8sconfig.inClusterConfig;
      st.k8sfile = k8sconfig.k8sfile;
      st.contextName = k8sconfig.contextName;
      st.clusterConfigured = k8sconfig.clusterConfigured;
      st.configuredServer = k8sconfig.configuredServer;
      st.kts = props.ts;
    }

    return st;
  }

  handleConfigure = () => {
    this.props.router.push("/settings#service-mesh");
  };

  pickImage(adapter) {
    const { classes } = this.props;
    let image = "/static/img/meshery-logo.png";
    let imageIcon = <img src={image} className={classes.expTitleIcon} />;
    if (adapter && adapter.name) {
      switch (adapter.name.toLowerCase()) {
        case "istio":
          image = "/static/img/istio.svg";
          imageIcon = <img src={image} className={classes.expIstioTitleIcon} />;
          break;
        case "linkerd":
          image = "/static/img/linkerd.svg";
          imageIcon = <img src={image} className={classes.expTitleIcon} />;
          break;
        case "consul":
          image = "/static/img/consul.svg";
          imageIcon = <img src={image} className={classes.expTitleIcon} />;
          break;
        case "network service mesh":
          image = "/static/img/nsm.svg";
          imageIcon = <img src={image} className={classes.expTitleIcon} />;
          break;
        case "octarine":
          image = "/static/img/octarine.svg";
          imageIcon = <img src={image} className={classes.expTitleIcon} />;
          break;
        case "citrix service mesh":
          image = "/static/img/citrix.svg";
          imageIcon = <img src={image} className={classes.expTitleIcon} />;
          break;
        case "osm":
          image = "/static/img/osm.svg";
          imageIcon = <img src={image} className={classes.expTitleIcon} />;
          break;
        case "kuma":
          image = "/static/img/kuma.svg";
          imageIcon = <img src={image} className={classes.expTitleIcon} />;
          break;
        // Disable support for NGINX SM
        // case "nginx service mesh":
        //   image = "/static/img/nginx-sm.svg";
        //   imageIcon = <img src={image} className={classes.expTitleIcon} />;
        //   break;
        case "traefik mesh":
          image = "/static/img/traefikmesh.svg";
          imageIcon = <img src={image} className={classes.expTitleIcon} />;
          break; 
      }
    }
    return imageIcon;
  }

  handleAdapterChange = () => {
    const self = this;
    return (event) => {
      const { meshAdapters } = self.state;
      if (event.target.value !== "") {
        const selectedAdapter = meshAdapters.filter(({ adapter_location }) => adapter_location === event.target.value);
        if (typeof selectedAdapter !== "undefined" && selectedAdapter.length === 1) {
          self.setState({ adapter: selectedAdapter[0] });
        }
      }
    };
  };

  renderIndividualAdapter() {
    const { meshAdapters } = this.props;
    let adapCount = 0;
    let adapter;
    meshAdapters.forEach((adap) => {
      if (adap.adapter_location === this.props.adapter) {
        adapter = adap;
        meshAdapters.forEach((ad) => {
          if (ad.name == adap.name) adapCount += 1;
        });
      }
    });
    if (adapter) {
      const imageIcon = this.pickImage(adapter);
      return (
        <React.Fragment>
          <MesheryAdapterPlayComponent adapter={adapter} adapCount={adapCount} adapter_icon={imageIcon} />
        </React.Fragment>
      );
    }
    return "";
  }

  render() {
    const { classes, k8sconfig, meshAdapters } = this.props;
    let { adapter } = this.state;

    if (k8sconfig.clusterConfigured === false || meshAdapters.length === 0) {
      return (
        <NoSsr>
          <React.Fragment>
            <div className={classes.alreadyConfigured}>
              <Button variant="contained" color="primary" size="large" onClick={this.handleConfigure}>
                <SettingsIcon className={classes.icon} />
                Configure Settings
              </Button>
            </div>
          </React.Fragment>
        </NoSsr>
      );
    }
    if (this.props.adapter && this.props.adapter !== "") {
      const indContent = this.renderIndividualAdapter();
      if (indContent !== "") {
        return indContent;
      } // else it will render all the available adapters
    }

    const self = this;
    const imageIcon = self.pickImage(adapter);
    let adapCount = 0;
    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.root}>
            <Grid container>
              <Grid item xs={12} className={classes.paneSection}>
                <TextField
                  select
                  id="adapter_id"
                  name="adapter_name"
                  label="Select Service Mesh Type"
                  fullWidth
                  value={adapter && adapter.adapter_location ? adapter.adapter_location : ""}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleAdapterChange()}
                >
                  {meshAdapters.map((ada) => (
                    <MenuItem key={`${ada.adapter_location}_${new Date().getTime()}`} value={ada.adapter_location}>
                      {/* <ListItemIcon> */}
                      {self.pickImage(ada)}
                      {/* </ListItemIcon> */}
                      <span className={classes.expTitle}>{ada.adapter_location}</span>
                    </MenuItem>
                  ))}
                </TextField>
              </Grid>
            </Grid>
          </div>
          <Divider variant="fullWidth" light />
          <Divider variant="fullWidth" light />
          {meshAdapters.forEach((adap) => {
            if (adap.adapter_location === this.props.adapter) {
              adapter = adap;
              meshAdapters.forEach((ad) => {
                if (ad.name == adap.name) adapCount += 1;
              });
            }
          })}
          {adapter && adapter.adapter_location && (
            <MesheryAdapterPlayComponent adapter={adapter} adapCount={adapCount} adapter_icon={imageIcon} />
          )}
        </React.Fragment>
      </NoSsr>
    );
  }
}

MesheryPlayComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = () => ({});

const mapStateToProps = (state) => {
  const k8sconfig = state.get("k8sConfig").toJS();
  const meshAdapters = state.get("meshAdapters").toJS();
  const meshAdaptersts = state.get("meshAdaptersts");
  return { k8sconfig, meshAdapters, meshAdaptersts };
};

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withRouter(MesheryPlayComponent)));
