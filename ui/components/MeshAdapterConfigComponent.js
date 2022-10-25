import React from "react";
import PropTypes from "prop-types";
import Button from "@material-ui/core/Button";
import { withStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import { NoSsr, Chip, IconButton } from "@material-ui/core";
import blue from "@material-ui/core/colors/blue";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { withRouter } from "next/router";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import ReactSelectWrapper from "./ReactSelectWrapper";
import { updateAdaptersInfo, updateProgress } from "../lib/store";
import dataFetch from "../lib/data-fetch";

const styles = (theme) => ({
  wrapperClass : { padding : theme.spacing(5), },
  buttons : { display : "flex",
    justifyContent : "flex-end", },
  button : { marginTop : theme.spacing(3),
    marginLeft : theme.spacing(1), },
  margin : { margin : theme.spacing(1), },
  alreadyConfigured : { textAlign : "center",
    padding : theme.spacing(20), },
  colorSwitchBase : { color : blue[300],
    "&$colorChecked" : { color : blue[500],
      "& + $colorBar" : { backgroundColor : blue[500], }, }, },
  colorBar : {},
  colorChecked : {},
  fileLabel : { width : "100%", },
  fileLabelText : {},
  inClusterLabel : { paddingRight : theme.spacing(2), },
  alignCenter : { textAlign : "center", },
  alignRight : { textAlign : "right",
    marginBottom : theme.spacing(2), },
  fileInputStyle : { opacity : "0.01", },
  icon : { width : theme.spacing(2.5), },
  istioIcon : { width : theme.spacing(1.5), },
});

class MeshAdapterConfigComponent extends React.Component {
  constructor(props) {
    super(props);
    const { meshAdapters } = props;
    this.state = {
      meshAdapters,
      availableAdapters : [],
      ts : new Date(),
      meshLocationURLError : false,
    };
  }

  static getDerivedStateFromProps(props, state) {
    const { meshAdapters, meshAdaptersts } = props;
    // if(meshAdapters.sort().join(',') !== state.meshAdapters.sort().join(',')){
    if (meshAdaptersts > state.ts) {
      return { meshAdapters,
        ts : meshAdaptersts, };
    }
    return {};
  }

  componentDidMount = () => {
    this.fetchAvailableAdapters();
  }

  fetchAvailableAdapters = () => {
    const self = this;
    this.props.updateProgress({ showProgress : true });
    dataFetch(
      "/api/system/adapters",
      {
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
      self.handleError("Unable to fetch available adapters")
    );
  };

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
          self.props.enqueueSnackbar("Adapter was successfully configured!", {
            variant : "success",
            "data-cy" : "adapterSuccessSnackbar",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ),
          });
          self.props.updateAdaptersInfo({ meshAdapters : result });
          self.fetchAvailableAdapters();
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
          this.props.enqueueSnackbar("Adapter was successfully removed!", { variant : "success",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ), });
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
        credentials : "include", },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Adapter was successfully pinged!", { variant : "success",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ), });
        }
      },
      self.handleError("error")
    );
  };

  handleError = (msg) => (error) => {
    this.props.updateProgress({ showProgress : false });
    const self = this;
    this.props.enqueueSnackbar(`${msg}: ${error}`, { variant : "error",
      action : (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration : 8000, });
  };

  configureTemplate = () => {
    const { classes } = this.props;
    const {
      availableAdapters, meshAdapters, meshLocationURL, meshLocationURLError
    } = this.state;

    let showAdapters = "";
    const self = this;
    if (meshAdapters.length > 0) {
      showAdapters = (
        <div className={classes.alignRight}>
          {meshAdapters.map((adapter) => {
            let image = "/static/img/meshery-logo.png";
            let logoIcon = <img src={image} className={classes.icon} />;
            if (adapter.name) {
              image = "/static/img/" + adapter.name.toLowerCase() + ".svg";
              logoIcon = <img src={image} className={classes.icon} />;
            }

            return (
              <Chip
                key={adapter.uniqueID}
                label={adapter.adapter_location}
                onDelete={self.handleDelete(adapter.adapter_location)}
                onClick={self.handleClick(adapter.adapter_location)}
                icon={logoIcon}
                variant="outlined"
                data-cy="chipAdapterLocation"
              />
            );
          })}
        </div>
      );
    }

    return (
      <NoSsr>
        <div className={classes.wrapperClass}>
          {showAdapters}

          <Grid container spacing={1} alignItems="flex-end">
            <Grid item xs={12}>
              <ReactSelectWrapper
                onChange={this.handleMeshLocURLChange}
                options={availableAdapters}
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
                onClick={this.handleSubmit}
                className={classes.button}
                data-cy="btnSubmitMeshAdapter"
              >
                Connect
              </Button>
            </div>
          </React.Fragment>
        </div>
      </NoSsr>
    );
  };

  render() {
    return this.configureTemplate();
  }
}

MeshAdapterConfigComponent.propTypes = { classes : PropTypes.object.isRequired, };

const mapDispatchToProps = (dispatch) => ({ updateAdaptersInfo : bindActionCreators(updateAdaptersInfo, dispatch),
  updateProgress : bindActionCreators(updateProgress, dispatch), });
const mapStateToProps = (state) => {
  const meshAdapters = state.get("meshAdapters").toJS();
  const meshAdaptersts = state.get("meshAdaptersts");
  return { meshAdapters, meshAdaptersts };
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(MeshAdapterConfigComponent)))
);
