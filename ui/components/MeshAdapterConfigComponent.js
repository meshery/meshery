import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import { NoSsr, Chip, Button, TextField, Tooltip, Avatar } from "@material-ui/core";
import blue from "@material-ui/core/colors/blue";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { withRouter } from "next/router";
import ReactSelectWrapper from "./ReactSelectWrapper";
import { updateAdaptersInfo, updateProgress } from "../lib/store";
import dataFetch from "../lib/data-fetch";
import changeAdapterState from './graphql/mutations/AdapterStatusMutation';
import { withNotify } from "../utils/hooks/useNotification";
import { EVENT_TYPES } from "../lib/event-types";
import BadgeAvatars from './CustomAvatar';

const styles = (theme) => ({
  wrapperClass : {
    padding : theme.spacing(5),
    backgroundColor : theme.palette.secondary.elevatedComponents,
    borderBottomLeftRadius : theme.spacing(1),
    borderBottomRightRadius : theme.spacing(1),
    marginTop : theme.spacing(2),
  },
  buttons : {
    display : "flex",
    justifyContent : "flex-end", paddingTop : "2rem"
  },
  button : {
    marginLeft : theme.spacing(1),
  },
  margin : { margin : theme.spacing(1), },
  alreadyConfigured : {
    textAlign : "center",
    padding : theme.spacing(20),
  },
  colorSwitchBase : {
    color : blue[300],
    "&$colorChecked" : {
      color : blue[500],
      "& + $colorBar" : { backgroundColor : blue[500], },
    },
  },
  colorBar : {},
  colorChecked : {},
  fileLabel : { width : "100%", },
  fileLabelText : {},
  inClusterLabel : { paddingRight : theme.spacing(2), },
  alignCenter : { textAlign : "center", },
  alignRight : {
    textAlign : "right",
    marginBottom : theme.spacing(2),
  },
  fileInputStyle : { opacity : "0.01", },
  // icon : { width : theme.spacing(2.5), },
  icon : {
    width : 20,
    height : 20
  },
  istioIcon : { width : theme.spacing(1.5), },
  chip : {
    marginRight : theme.spacing(1),
    marginBottom : theme.spacing(1),
  }
});

const STATUS = {
  DEPLOYED : "DEPLOYED",
  UNDEPLOYED : "UNDEPLOYED",
  DEPLOYING : "DEPLOYING",
  UNDEPLOYING : "UNDEPLOYING",
}

class  MeshAdapterConfigComponent extends React.Component {
  constructor(props) {
    super(props);
    const { meshAdapters } = props;
    this.labelRef = React.createRef();
    this.state = {
      meshAdapters,
      setAdapterURLs : [],
      availableAdapters : [],
      ts : new Date(),
      meshLocationURLError : false,
      selectedAvailableAdapterError : false,
      adapterStates : {},
    };
  }

  static getDerivedStateFromProps(props, state) {
    const { meshAdapters, meshAdaptersts } = props;
    // if(meshAdapters.sort().join(',') !== state.meshAdapters.sort().join(',')){
    if (meshAdaptersts > state.ts) {
      return {
        meshAdapters,
        ts : meshAdaptersts
      };
    }
    return {};
  }

  componentDidMount = () => {
    this.fetchSetAdapterURLs();
    this.fetchAvailableAdapters();
    this.setAdapterStates();
  }

  fetchSetAdapterURLs = () => {
    const self = this;
    this.props.updateProgress({ showProgress : true });
    dataFetch(
      "/api/system/adapters",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          const options = result.map((res) => ({
            value : res.adapter_location,
            label : res.adapter_location,
          }));
          this.setState({ setAdapterURLs : options });
        }
      },
      self.handleError("Unable to fetch available adapters")
    );
  };

  fetchAvailableAdapters = () => {
    const self = this;
    this.props.updateProgress({ showProgress : true });
    dataFetch(
      "/api/system/availableAdapters",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          const options = result.map((res) => ({
            value : res.adapter_location,
            label : res.name,
          }));
          this.setState({ availableAdapters : options });
        }
      },
      self.handleError("Unable to fetch available adapters")
    );
  };

  setAdapterStates = () => {
    const { meshAdapters } = this.state;
    const initialAdapterStates = {};

    meshAdapters.forEach((adapter) => {
      const lable = adapter.name.toUpperCase()
      initialAdapterStates[lable] = STATUS.UNDEPLOYED;
    });

    this.setState({
      adapterStates : initialAdapterStates,
    });
  }

  getStatusColor = (status) => {
    if (status === STATUS.DEPLOYED) {
      return "#00B39F";
    } else if (status === STATUS.UNDEPLOYED) {
      return "#808080";
    } else if (status === STATUS.DEPLOYING) {
      return "#EBC017"
    } else if (status === STATUS.UNDEPLOYING) {
      return "#E75225"
    }
  }

  handleChange = (name) => (event) => {
    if (name === "meshLocationURL" && event.target.value !== "") {
      this.setState({ meshLocationURLError : false });
    }
    this.setState({ [name] : event.target.value });
  };

  handleMeshLocURLChange = (newValue) => {
    // console.log(newValue);
    // console.log(`action: ${actionMeta.action}`);
    // console.groupEnd();
    if (typeof newValue !== "undefined") {
      this.setState({ meshLocationURL : newValue, meshLocationURLError : false });
    }
  };

  handleDeployPortChange = (newValue) => {
    if (typeof newValue !== "undefined") {
      console.log("port change to " + (newValue.value))
      this.setState({ meshDeployURL : newValue.value, meshDeployURLError : false });
    }
  }

  handleAvailableAdapterChange = (newValue) => {
    if (typeof newValue !== "undefined") {
      // Trigger label animation manually
      this.labelRef.current.querySelector('label').classList.add('MuiInputLabel-shrink');
      this.setState({ selectedAvailableAdapter : newValue, selectedAvailableAdapterError : false });
      if (newValue !== null) {
        this.setState({ meshDeployURL : newValue.value, meshDeployURLError : false });
      }
    }
  };

  handleSubmit = () => {
    const { meshLocationURL } = this.state;

    if (!meshLocationURL || !meshLocationURL.value || meshLocationURL.value === "") {
      this.setState({ meshLocationURLError : true });
      return;
    }

    this.submitConfig();
  };

  submitConfig = () => {
    const { meshLocationURL } = this.state;

    const data = { meshLocationURL : meshLocationURL.value };

    const params = Object.keys(data)
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
      .join("&");

    this.props.updateProgress({ showProgress : true });
    const self = this;
    dataFetch(
      "/api/system/adapter/manage",
      {

        method : "POST",
        credentials : "include",
        headers : { "Content-Type" : "application/x-www-form-urlencoded;charset=UTF-8", },
        body : params,
      },
      (result) => {
        self.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          self.setState({ meshAdapters : result, meshLocationURL : "" });
          const notify = self.props.notify;
          notify({ message : "Adapter was configured!", event_type : EVENT_TYPES.SUCCESS })
          self.props.updateAdaptersInfo({ meshAdapters : result });
          self.fetchSetAdapterURLs();
        }
      },
      self.handleError("Adapter was not configured due to an error")
    );
  };

  handleDelete = (adapterLoc) => () => {
    // const { meshAdapters } = this.state;
    this.props.updateProgress({ showProgress : true });
    const self = this;
    dataFetch(
      `/api/system/adapter/manage?adapter=${encodeURIComponent(adapterLoc)}`,
      {

        method : "DELETE",
        credentials : "include",
      },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          this.setState({ meshAdapters : result });
          const notify = self.props.notify;
          notify({ message : "Adapter was removed!", event_type : EVENT_TYPES.SUCCESS })
          this.props.updateAdaptersInfo({ meshAdapters : result });
        }
      },
      self.handleError("Adapter was not removed due to an error")
    );
  };

  handleClick = (adapterLoc) => () => {
    // const { meshAdapters } = this.state;
    this.props.updateProgress({ showProgress : true });
    const self = this;
    dataFetch(
      `/api/system/adapters?adapter=${encodeURIComponent(adapterLoc)}`,
      {
        credentials : "include",
      },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          const notify = self.props.notify;
          notify({ message : "Adapter was pinged!", event_type : EVENT_TYPES.SUCCESS })
        }
      },
      self.handleError("error")
    );
  };

  handleAdapterDeploy = () => {
    const { selectedAvailableAdapter, meshDeployURL } = this.state;
    if (!selectedAvailableAdapter || !selectedAvailableAdapter.value || selectedAvailableAdapter.value === "") {
      this.setState({ selectedAvailableAdapterError : true });
      return;
    }

    const adapterLabel = selectedAvailableAdapter.label.replace(/^meshery-/, '').toUpperCase();
    this.setState(prevState => ({
      adapterStates : {
        ...prevState.adapterStates,
        [adapterLabel] : STATUS.DEPLOYING,
      },
    }));

    if (!meshDeployURL || meshDeployURL === "") {
      console.log(meshDeployURL)
      this.setState({ meshDeployURLError : true });
      return;
    }

    this.props.updateProgress({ showProgress : true });

    const variables = {
      status : "ENABLED",
      adapter : selectedAvailableAdapter.label,
      targetPort : meshDeployURL
    };

    changeAdapterState((response, errors) => {
      this.props.updateProgress({ showProgress : false });

      if (errors !== undefined) {
        this.handleError("Unable to Deploy adapter");
        this.setState(prevState => ({
          adapterStates : {
            ...prevState.adapterStates,
            [adapterLabel] : STATUS.UNDEPLOYED,
          },
        }))
      }

      this.setState(prevState => ({
        adapterStates : {
          ...prevState.adapterStates,
          [adapterLabel] : STATUS.DEPLOYED,
        },
      }))
      const notify = this.props.notify;
      notify({ message : "Adapter " + response.adapterStatus.toLowerCase(), event_type : EVENT_TYPES.SUCCESS })
    }, variables);
  };

  handleAdapterUndeploy = () => {
    const { meshLocationURL, availableAdapters } = this.state;

    if (!meshLocationURL || !meshLocationURL.value || meshLocationURL.value === "") {
      this.setState({ meshLocationURLError : true });
      return;
    }

    this.props.updateProgress({ showProgress : true });

    const targetPort = function getTargetPort(location) {
      if (!location.value) {
        return null
      }

      if (location.value.includes(":")) {
        return location.value.split(":")[1]
      }

      return location.value;
    }(meshLocationURL)

    const adapterName = function getAdapterName(location) {
      if (!location.value) {
        return null
      }

      if (location.value.includes(":")) {
        return location.value.split(":")[0]
      }

      return location.value;
    }(meshLocationURL)

    const adapterLabel = (availableAdapters.find(adapter => adapter.value === targetPort)?.label || "").replace(/^meshery-/, '').toUpperCase();

    this.setState(prevState => ({
      adapterStates : {
        ...prevState.adapterStates,
        [adapterLabel] : STATUS.UNDEPLOYING,
      },
    }));

    const variables = {
      status : "DISABLED",
      adapter : adapterName,
      targetPort : targetPort,
    };

    changeAdapterState((response, errors) => {
      this.props.updateProgress({ showProgress : false });

      if (errors !== undefined) {
        console.error(errors)
        this.handleError("Unable to Deploy adapter");
        this.setState(prevState => ({
          adapterStates : {
            ...prevState.adapterStates,
            [adapterLabel] : STATUS.DEPLOYED,
          },
        }))
      }
      const notify = this.props.notify;
      notify({ message : "Adapter " + response.adapterStatus.toLowerCase(), event_type : EVENT_TYPES.SUCCESS })
      this.setState(prevState => ({
        adapterStates : {
          ...prevState.adapterStates,
          [adapterLabel] : STATUS.UNDEPLOYED,
        },
      }));
    }, variables);
  };

  handleError = (msg) => (error) => {
    this.props.updateProgress({ showProgress : false });
    const notify = this.props.notify;
    notify({ message : msg, event_type : EVENT_TYPES.ERROR, details : error.toString() })
  };

  configureTemplate = () => {
    const { classes } = this.props;
    const {
      availableAdapters, setAdapterURLs, meshAdapters, meshLocationURL, meshLocationURLError, meshDeployURLError, selectedAvailableAdapter, selectedAvailableAdapterError, meshDeployURL, adapterStates
    } = this.state;

    let showAdapters = "";
    const self = this;
    if (meshAdapters.length > 0) {
      showAdapters = (
        <div className={classes.alignRight}>

          {meshAdapters.map((adapter) => {
            let image = "/static/img/meshery-logo.png";
            // let logoIcon = <img src={image} className={classes.icon} />;
            if (adapter.name) {
              image = "/static/img/" + adapter.name.toLowerCase() + ".svg";
              // logoIcon = <img src={image} className={classes.icon} />;
            }

            return (
              <Tooltip
                key={adapter.uniqueID}
                title={
                  `Meshery Adapter for
                        ${adapter.name
                    .toLowerCase()
                    .split(" ")
                    .map((s) => s.charAt(0).toUpperCase() + s.substring(1))
                    .join(" ")} (${adapter.version})`}>
                <Chip
                  className={classes.chip}
                  label={adapter.adapter_location}
                  onDelete={self.handleDelete(adapter.adapter_location)}
                  onClick={self.handleClick(adapter.adapter_location)}
                  icon={
                    // logoIcon
                    <BadgeAvatars color={this.getStatusColor(adapterStates[adapter.name])}>
                      <Avatar alt={adapter.name} src={image} className={classes.icon} />
                    </BadgeAvatars>
                  }
                  variant="outlined"
                  data-cy="chipAdapterLocation"
                />
              </Tooltip>
            );
          })}
        </div>
      );
    }

    return (
      <NoSsr>
        <div className={classes.wrapperClass} data-cy="mesh-adapter-connections" >
          {showAdapters}

          <Grid container spacing={1} alignItems="flex-end">
            <Grid item xs={12} data-cy="mesh-adapter-url">
              <ReactSelectWrapper
                onChange={this.handleMeshLocURLChange}
                options={setAdapterURLs}
                value={meshLocationURL}
                // placeholder={'Mesh Adapter URL'}
                label="Mesh Adapter URL"
                error={meshLocationURLError}
              />
            </Grid>
          </Grid>
          <React.Fragment>
            <div className={classes.buttons}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                size="large"
                onClick={this.handleAdapterUndeploy}
                className={classes.button}
              >
                Undeploy
              </Button>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                size="large"
                onClick={this.handleSubmit}
                className={classes.button}
                data-cy="btnSubmitMeshAdapter"
              >
                Connect
              </Button>
            </div>
          </React.Fragment>
          <Grid container spacing={1} alignItems="flex-end" style={{ marginTop : '50px' }}>
            <Grid item xs={12}>
              <ReactSelectWrapper
                onChange={this.handleAvailableAdapterChange}
                options={availableAdapters}
                value={selectedAvailableAdapter}
                // placeholder={'Mesh Adapter URL'}
                label="Available Mesh Adapter"
                error={selectedAvailableAdapterError}
              />
            </Grid>
          </Grid>
          <Grid container spacing={1} alignItems="flex-end" justifyContent="flex-end">
            <div ref={this.labelRef}>
              <TextField
                id="deployPort"
                type="text"
                label="Enter Port"
                onChange={(e) => this.handleDeployPortChange(e.target)}
                value={meshDeployURL}
                error={meshDeployURLError}
              />
            </div>
            <React.Fragment>
              <div className={classes.buttons}>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  onClick={this.handleAdapterDeploy}
                  className={classes.button}
                >
                  Deploy
                </Button>
              </div>
            </React.Fragment>
          </Grid>
        </div>
      </NoSsr>
    );
  };

  render() {
    return this.configureTemplate();
  }
}

MeshAdapterConfigComponent.propTypes = { classes : PropTypes.object.isRequired, };

const mapDispatchToProps = (dispatch) => ({
  updateAdaptersInfo : bindActionCreators(updateAdaptersInfo, dispatch),
  updateProgress : bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const meshAdapters = state.get("meshAdapters").toJS();
  const meshAdaptersts = state.get("meshAdaptersts");

  return { meshAdapters, meshAdaptersts };
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withNotify(MeshAdapterConfigComponent)))
);