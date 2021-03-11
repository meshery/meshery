//@ts-check
import React, { useState, useEffect } from "react";
import { connect } from "react-redux";
import { withSnackbar } from "notistack";
import { updateProgress } from "../../lib/store";
import { bindActionCreators } from "redux";
import { Grid } from "@material-ui/core";
import IconButton from "@material-ui/core/IconButton";
import CloseIcon from "@material-ui/icons/Close";
import { withRouter } from "next/router";
import dataFetch from "../../lib/data-fetch";
import MesheryMetrics from "../MesheryMetrics";

const MESHERY_PERFORMANCE_URL = "/api/user/performance/profiles";
const MESHERY_PERFORMANCE_TEST_URL = "/api/user/performance/profiles/results";

function Dashboard({ updateProgress, enqueueSnackbar, closeSnackbar, grafana, router }) {
  const [profiles, setProfiles] = useState({
    count: 0,
    profiles: [],
  });
  const [tests, setTests] = useState({
    count: 0,
    tests: [],
  });

  /**
   * fetch performance profiles when the page loads
   */
  useEffect(() => {
    fetchTestProfiles();
    fetchTests();
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
          setProfiles({
            count: result.total_count || 0,
            profiles: result.profiles || [],
          });
        }
      },
      handleError("Failed to Fetch Profiles")
    );
  }

  function fetchTests() {
    updateProgress({ showProgress: true });

    dataFetch(
      `${MESHERY_PERFORMANCE_TEST_URL}`,
      {
        credentials: "include",
      },
      (result) => {
        updateProgress({ showProgress: false });
        if (result) {
          setTests({
            count: result.total_count || 0,
            tests: result.results || [],
          });
        }
      },
      handleError("Failed to Fetch Results")
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
          <b>Total Profiles Created:</b> {profiles.count}
        </div>
        <div>
          <b>Total Tests Run:</b> {tests.count}
        </div>
      </Grid>
      <Grid item md={6} xs={12}>
        <MesheryMetrics
          boardConfigs={grafana.selectedBoardsConfigs}
          grafanaURL={grafana.grafanaURL}
          grafanaAPIKey={grafana.grafanaAPIKey}
          handleGrafanaChartAddition={() => router.push("/settings/#metrics")}
        />
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

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(Dashboard)));
