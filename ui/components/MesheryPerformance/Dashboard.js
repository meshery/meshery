//@ts-check
import React, { useState, useEffect } from "react";
import { connect } from "react-redux";
import { withSnackbar } from "notistack";
import { updateProgress } from "../../lib/store";
import { bindActionCreators } from "redux";
import { Grid, Typography, Button } from "@material-ui/core";
import IconButton from "@material-ui/core/IconButton";
import CloseIcon from "@material-ui/icons/Close";
import AddIcon from "@material-ui/icons/Add";
import dataFetch from "../../lib/data-fetch";
import GrafanaCustomCharts from "../GrafanaCustomCharts";

const MESHERY_PERFORMANCE_URL = "/api/user/performance/profiles";

function Dashboard({ updateProgress, enqueueSnackbar, closeSnackbar, grafana }) {
  const [count, setCount] = useState(0);
  const [testProfiles, setTestProfiles] = useState([]);
  console.log(testProfiles);

  /**
   * renderGrafanaCustomCharts takes in the configuration and renders
   * the grafana boards. If the configuration is empty then it renders
   * a note directing a user to install grafana and prometheus
   * @param {Array<{ board: any, panels: Array<any>, templateVars: Array<any>}>} boardConfigs grafana board configs
   * @param {string} grafanaURL grafana URL
   * @param {string} grafanaAPIKey grafana API keey
   */
  function renderGrafanaCustomCharts(boardConfigs, grafanaURL, grafanaAPIKey) {
    if (boardConfigs?.length)
      return (
        <>
          <Typography
            align="center"
            variant="h6"
            style={{
              margin: "0 0 2.5rem 0",
            }}
          >
            Service Mesh Metrics
          </Typography>
          <GrafanaCustomCharts
            // @ts-ignore
            enableGrafanaChip
            boardPanelConfigs={boardConfigs || []}
            grafanaURL={grafanaURL || ""}
            grafanaAPIKey={grafanaAPIKey || ""}
          />
        </>
      );

    return (
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
          No Service Mesh Metrics Configurations Found
        </Typography>
        <Button
          aria-label="Add Grafana Charts"
          variant="contained"
          color="primary"
          size="large"
          onClick={() => this.props.router.push("/settings/#metrics")}
        >
          <AddIcon />
          Configure Service Mesh Metrics
        </Button>
      </div>
    );
  }

  /**
   * fetch performance profiles when the page loads
   */
  useEffect(() => {
    fetchTestProfiles();
  }, []);

  function fetchTestProfiles() {
    updateProgress({ showProgress: true });

    dataFetch(
      `${MESHERY_PERFORMANCE_URL}`,
      {
        credentials: "include",
      },
      (result) => {
        updateProgress({ showProgress: false });
        if (result) {
          setTestProfiles(result.profiles || []);
          setCount(result.total_count || 0);
        }
      },
      handleError("Failed to Fetch Profiles")
    );
  }

  function handleError(msg) {
    return function (error) {
      updateProgress({ showProgress: false });

      enqueueSnackbar(`${msg}: ${error}`, {
        variant: "error",
        action: function Action(key) {
          return (
            <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
              <CloseIcon />
            </IconButton>
          );
        },
        autoHideDuration: 8000,
      });
    };
  }

  return (
    <Grid container spacing={1} style={{ padding: "0.5rem" }} alignItems="center" alignContent="space-around">
      <Grid item md={6} xs={12}>
        <div>
          <b>Total Profiles Created:</b> {count}
        </div>
      </Grid>
      <Grid item md={6} xs={12}>
        {renderGrafanaCustomCharts(grafana.selectedBoardsConfigs, grafana.grafanaURL, grafana.grafanaAPIKey)}
      </Grid>
    </Grid>
  );
}

const mapStateToProps = (st) => {
  const grafana = st.get("grafana").toJS();
  return { grafana: { ...grafana, ts: new Date(grafana.ts) } };
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(withSnackbar(Dashboard));
