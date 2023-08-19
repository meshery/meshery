import React, { Component } from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import { NoSsr, Typography, IconButton, Box } from "@material-ui/core";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import CloseIcon from "@material-ui/icons/Close";
// import { withSnackbar } from "notistack";
import dataFetch from "../../../lib/data-fetch";
import GrafanaConfigComponent from "./GrafanaConfigComponent";
import GrafanaSelectionComponent from "./GrafanaSelectionComponent";
import GrafanaDisplaySelection from "./GrafanaDisplaySelection";
// import GrafanaCharts from './GrafanaCharts';
import { updateGrafanaConfig, updateProgress } from "../../../lib/store";
import GrafanaCustomCharts from "./GrafanaCustomCharts";
import fetchAvailableAddons from "../../graphql/queries/AddonsStatusQuery";
import { getK8sClusterIdsFromCtxId } from "../../../utils/multi-ctx";
import { withNotify } from "../../../utils/hooks/useNotification";

const grafanaStyles = (theme) => ({
  buttons : { display : "flex",
    //   justifyContent: 'flex-end',
  },
  button : { marginTop : theme.spacing(3),
    //   marginLeft: theme.spacing(1),
  },
  margin : { margin : theme.spacing(1), },
  icon : { width : theme.spacing(2.5), },
  alignRight : { textAlign : "right", },
  formControl : { margin : theme.spacing(1),
    minWidth : 180, },
  panelChips : { display : "flex",
    flexWrap : "wrap", },
  panelChip : { margin : theme.spacing(0.25), },
  chartTitle : { marginLeft : theme.spacing(3),
    marginTop : theme.spacing(2), textAlign : "center", },
  grafanChartsWrapper : { backgroundColor : theme.palette.secondary.elevatedComponents,
    padding : theme.spacing(1) , marginTop : theme.spacing(2) }
});

const getGrafanaBoards = (self, cb = () => {}) => {
  const {
    grafanaURL, grafanaAPIKey, grafanaBoardSearch, selectedBoardsConfigs
  } = self.state;
  if (typeof grafanaURL === "undefined" || grafanaURL === "") {
    return;
  }
  self.props.updateProgress({ showProgress : true });
  dataFetch(
    `/api/telemetry/metrics/grafana/boards?dashboardSearch=${grafanaBoardSearch}`,
    {
      method : "GET",
      credentials : "include", },
    (result) => {
      self.props.updateProgress({ showProgress : false });
      if (typeof result !== "undefined") {
        self.setState({ grafanaBoards : result });
        self.props.updateGrafanaConfig({ grafana : {
          grafanaURL,
          grafanaAPIKey,
          grafanaBoardSearch,
          grafanaBoards : result,
          selectedBoardsConfigs,
        }, });
        cb()
      }
    },
    self.handleError("There was an error communicating with Grafana")
  );
};

export const submitGrafanaConfigure = (self, cb) => {
  const {
    grafanaURL, grafanaAPIKey, grafanaBoards, grafanaBoardSearch, selectedBoardsConfigs
  } = self.state;
  if (grafanaURL === "") {
    return;
  }
  const data = { grafanaURL,
    grafanaAPIKey, };
  const params = Object.keys(data)
    .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
    .join("&");
  // console.log(`data to be submitted for load test: ${params}`);
  self.props.updateProgress({ showProgress : true });
  dataFetch(
    "/api/telemetry/metrics/grafana/config",
    {

      method : "POST",
      credentials : "include",
      headers : { "Content-Type" : "application/x-www-form-urlencoded;charset=UTF-8", },
      body : params,
    },
    (result) => {
      self.props.updateProgress({ showProgress : false });
      if (typeof result !== "undefined") {
        self.props.enqueueSnackbar("Grafana was configured!", { variant : "success",
          autoHideDuration : 2000,
          action : function Loader(key) {
            return (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            );
          }, });
        self.setState({ grafanaConfigSuccess : true });
        self.props.updateGrafanaConfig({ grafana : {
          grafanaURL,
          grafanaAPIKey,
          grafanaBoardSearch,
          grafanaBoards,
          selectedBoardsConfigs,
        }, });
        getGrafanaBoards(self, cb);
      }
    },
    self.handleError("There was an error communicating with Grafana")
  );
};

class GrafanaComponent extends Component {
  constructor(props) {
    super(props);

    this.state = {
      urlError : false,

      grafanaConfigSuccess : props.grafana.grafanaURL !== "",
      grafanaBoardSearch : "",
      grafanaURL : props.grafana.grafanaURL,
      grafanaAPIKey : props.grafana.grafanaAPIKey,
      grafanaBoards : props.grafana.grafanaBoards,
      selectedBoardsConfigs : props.grafana.selectedBoardsConfigs,
      ts : props.grafana.ts,
    };
  }

  UNSAFE_componentWillReceiveProps(nextProps) {
    const { grafanaURL, grafanaAPIKey, selectedBoardsConfigs } = nextProps.grafana;
    if (nextProps.grafana.ts > this.state.ts) {
      this.setState(
        {
          grafanaURL,
          grafanaAPIKey,
          selectedBoardsConfigs,
          grafanaConfigSuccess : grafanaURL !== "",
          ts : nextProps.ts,
        },
        () => getGrafanaBoards(this)
      );
    }
    return {};
  }

  componentDidMount() {
    const self = this;

    if (self.props.isMeshConfigured)
      dataFetch(
        "/api/telemetry/metrics/grafana/config",
        {
          method : "GET",
          credentials : "include",
          headers : {
            "Content-Type" : "application/x-www-form-urlencoded;charset=UTF-8",
          },
        },
        (result) => {
          self.props.updateProgress({ showProgress : false });
          if (!(typeof result !== "undefined" && result?.grafanaURL && result?.grafanaURL !="")) {
            let selector = {
              type : "ALL_MESH",
              k8sClusterIDs : this.getK8sClusterIds()
            };
            fetchAvailableAddons(selector).subscribe({
              next : (res) => {
                res?.addonsState?.forEach((addon) => {
                  if (addon.name === "grafana" && self.state.grafanaURL === "") {
                    self.setState({ grafanaURL : "http://" + addon.endpoint })
                    submitGrafanaConfigure(self, () => self.setState({ selectedBoardsConfigs : self.state.grafanaBoards?.[2]
                      ? [self.state.grafanaBoards[2]]
                      : [] }));
                  }
                });
              },
              error : (err) => console.log("error registering Grafana: " + err),
            });
          }
        },
        self.handleError("There was an error communicating with grafana config")
      )
  }

  getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(this.props.selectedK8sContexts, this.props.k8sconfig)
  }

  handleChange = (name) => (event) => {
    if (name === "grafanaURL" && event.target.value !== "") {
      this.setState({ urlError : false });
    }
    if (name === "grafanaBoardSearch") {
      if (this.boardSearchTimeout) clearTimeout(this.boardSearchTimeout);
      this.boardSearchTimeout = setTimeout(() => getGrafanaBoards(this), 500); // to delay the search by a few.
    }

    this.setState({ [name] : event.target.value });
  };

  handleChangeApiKey = (event) => {
    this.setState({ grafanaAPIKey : event.target.value });
  };

  handleGrafanaConfigure = () => {
    const { grafanaURL } = this.state;
    if (
      grafanaURL === "" ||
      !(grafanaURL.toLowerCase().startsWith("http://") || grafanaURL.toLowerCase().startsWith("https://"))
    ) {
      this.setState({ urlError : true });
      return;
    }
    submitGrafanaConfigure(this);
  };

  handleError = (msg) => () => {
    const self = this;
    // this.setState({timerDialogOpen: false });
    this.props.updateProgress({ showProgress : false });
    this.props.enqueueSnackbar(msg, { variant : "error",
      action : (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration : 8000, });
  };

  handleGrafanaChipDelete = () => {
    this.props.updateProgress({ showProgress : true });
    const self = this;
    dataFetch(
      "/api/telemetry/metrics/grafana/config",
      {
        method : "DELETE",
        credentials : "include", },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          self.setState({
            grafanaConfigSuccess : false,
            grafanaURL : "",
            grafanaAPIKey : "",
            grafanaBoardSearch : "",
            grafanaBoards : [],
            selectedBoardsConfigs : [],
          });
          self.props.updateGrafanaConfig({ grafana : {
            grafanaURL : "",
            grafanaAPIKey : "",
            grafanaBoardSearch : "",
            grafanaBoards : [],
            selectedBoardsConfigs : [],
          }, });
        }
      },
      self.handleError("There was an error communicating with Grafana")
    );
  };

  handleGrafanaClick = () => {
    this.props.updateProgress({ showProgress : true });
    const self = this;
    dataFetch(
      "/api/telemetry/metrics/grafana/ping",
      {
        credentials : "include", },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          this.props.enqueueSnackbar("Grafana pinged!", { variant : "success",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ), });
        }
      },
      self.handleError("Could not ping Grafana.")
    );
  };

  addSelectedBoardPanelConfig = (boardsSelection) => {
    const { selectedBoardsConfigs } = this.state;

    if (boardsSelection && boardsSelection.panels && boardsSelection.panels.length) {
      selectedBoardsConfigs.push(boardsSelection);
      this.persistBoardSelection(selectedBoardsConfigs);
    }
  };

  deleteSelectedBoardPanelConfig = (indexes) => {
    const { selectedBoardsConfigs } = this.state;
    indexes.sort();
    for (let i = indexes.length - 1; i >= 0; i--) {
      selectedBoardsConfigs.splice(indexes[i], 1);
    }
    this.persistBoardSelection(selectedBoardsConfigs);
  };

  persistBoardSelection(selectedBoardsConfigs) {
    const {
      grafanaURL, grafanaAPIKey, grafanaBoards, grafanaBoardSearch
    } = this.state;
    const self = this;
    dataFetch(
      "/api/telemetry/metrics/grafana/boards",
      {
        method : "POST",
        credentials : "include",
        headers : { "Content-Type" : "application/json;charset=UTF-8", },
        body : JSON.stringify(selectedBoardsConfigs),
      },
      (result) => {
        this.props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          self.setState({ selectedBoardsConfigs });
          self.props.updateGrafanaConfig({ grafana : {
            grafanaURL,
            grafanaAPIKey,
            grafanaBoardSearch,
            grafanaBoards,
            selectedBoardsConfigs,
          }, });

          self.props.enqueueSnackbar("Grafana board selection was saved!", { variant : "success",
            autoHideDuration : 2000,
            action : (key) => (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            ), });
        }
      },
      self.handleError("There was an error persisting the board selection")
    );
  }

  render() {
    const { classes } = this.props;
    const {
      urlError,
      grafanaURL,
      grafanaConfigSuccess,
      grafanaAPIKey,
      grafanaBoards,
      grafanaBoardSearch,
      selectedBoardsConfigs,
    } = this.state;
    if (grafanaConfigSuccess) {
      let displaySelec = "";
      if (selectedBoardsConfigs.length > 0) {
        displaySelec = (
          <React.Fragment>
            <GrafanaDisplaySelection
              boardPanelConfigs={selectedBoardsConfigs}
              deleteSelectedBoardPanelConfig={this.deleteSelectedBoardPanelConfig}
            />
            <Box className={classes.grafanChartsWrapper}>
              <Typography variant="h6" gutterBottom className={classes.chartTitle}>
                Grafana charts
              </Typography>
              {/* <GrafanaCharts
                    boardPanelConfigs={selectedBoardsConfigs}
                    grafanaURL={grafanaURL} /> */}
              <div style={{ padding : "0 1rem" }}>
                <GrafanaCustomCharts
                  boardPanelConfigs={selectedBoardsConfigs}
                  grafanaURL={grafanaURL}
                  grafanaAPIKey={grafanaAPIKey}
                />
              </div>
            </Box>
          </React.Fragment>
        );
      }

      return (
        <NoSsr>
          <React.Fragment>
            <GrafanaSelectionComponent
              grafanaURL={grafanaURL}
              grafanaBoards={grafanaBoards}
              grafanaBoardSearch={grafanaBoardSearch}
              handleGrafanaBoardSearchChange={this.handleChange}
              handleGrafanaChipDelete={this.handleGrafanaChipDelete}
              handleGrafanaClick={this.handleGrafanaClick}
              addSelectedBoardPanelConfig={this.addSelectedBoardPanelConfig}
              handleError={this.handleError("There was an error communicating with Grafana")}
            />
            {displaySelec}
          </React.Fragment>
        </NoSsr>
      );
    }
    return (
      <NoSsr>
        <GrafanaConfigComponent
          grafanaURL={grafanaURL && { label : grafanaURL, value : grafanaURL }}
          options={this.props.scannedGrafana.map((graf) => ({ label : graf, value : graf }))}
          grafanaAPIKey={grafanaAPIKey}
          urlError={urlError}
          handleChange={(name) => {
            // Simulating event.target.value
            return (value) => this.handleChange(name)({ target : { value } });
          }}
          handleChangeApiKey={this.handleChangeApiKey}
          handleGrafanaConfigure={this.handleGrafanaConfigure}
        />
      </NoSsr>
    );
  }
}

GrafanaComponent.propTypes = { classes : PropTypes.object.isRequired,
  scannedGrafana : PropTypes.array.isRequired, };

const mapDispatchToProps = (dispatch) => ({ updateGrafanaConfig : bindActionCreators(updateGrafanaConfig, dispatch),
  updateProgress : bindActionCreators(updateProgress, dispatch), });

const mapStateToProps = (st) => {
  const grafana = st.get("grafana").toJS();
  const selectedK8sContexts = st.get('selectedK8sContexts');
  const k8sconfig = st.get("k8sConfig");

  return { grafana : { ...grafana, ts : new Date(grafana.ts) }, selectedK8sContexts, k8sconfig };
};

export default withStyles(grafanaStyles)(connect(mapStateToProps, mapDispatchToProps)(withNotify(GrafanaComponent)));
