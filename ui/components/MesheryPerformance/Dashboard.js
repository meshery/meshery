//@ts-check
import React, { useState, useEffect } from "react";
import { connect } from "react-redux";
import { withSnackbar } from "notistack";
import { updateProgress } from "../../lib/store";
import { bindActionCreators } from "redux";
import { Button, Grid, Paper, Typography } from "@material-ui/core";
import IconButton from "@material-ui/core/IconButton";
import CloseIcon from "@material-ui/icons/Close";
import { withRouter } from "next/router";
import dataFetch from "../../lib/data-fetch";
import MesheryMetrics from "../MesheryMetrics";
import { withStyles } from "@material-ui/core/styles";
import PerformanceCalendar from "./PerformanceCalendar";
import GenericModal from "../GenericModal";
import MesheryPerformanceComponent from "./index";

const MESHERY_PERFORMANCE_URL = "/api/user/performance/profiles";
const MESHERY_PERFORMANCE_TEST_URL = "/api/user/performance/profiles/results";

const styles = () => ({
  paper: {
    padding: "1rem",
  },
});

function Dashboard({ updateProgress, enqueueSnackbar, closeSnackbar, grafana, router, classes }) {
  const [profiles, setProfiles] = useState({
    count: 0,
    profiles: [],
  });
  const [tests, setTests] = useState({
    count: 0,
    tests: [],
  });

  const [runTest, setRunTest] = useState(false);

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
    <>
      <Grid container spacing={2} style={{ padding: "0.5rem" }} alignItems="flex-start" alignContent="space-around">
        <Grid item lg={6} xs={12}>
          <Paper className={classes.paper}>
            <Grid container spacing={1}>
              <Grid item xs>
                <Paper className={classes.paper}>
                  <div style={{ display: "flex", alignItems: "center" }}>
                    <Typography variant="h2" component="div" color="primary" style={{ marginRight: "0.75rem" }}>
                      {tests.count}
                    </Typography>
                    <Typography variant="body1" style={{ color: "rgba(0, 0, 0, 0.54)" }} component="div">
                      Results
                    </Typography>
                  </div>
                  <div style={{ margin: "2rem 0 0 auto", width: "fit-content" }}>
                    <Button variant="contained" color="primary" onClick={() => setRunTest(true)}>
                      Run Test
                    </Button>
                  </div>
                </Paper>
              </Grid>
              <Grid item xs>
                <Paper className={classes.paper}>
                  <div style={{ display: "flex", alignItems: "center" }}>
                    <Typography variant="h2" component="div" color="primary" style={{ marginRight: "0.75rem" }}>
                      {profiles.count}
                    </Typography>
                    <Typography variant="body1" style={{ color: "rgba(0, 0, 0, 0.54)" }} component="div">
                      Profiles
                    </Typography>
                  </div>
                  <div style={{ margin: "2rem 0 0 auto", width: "fit-content" }}>
                    <Button variant="contained" color="primary" onClick={() => router.push("/performance/profiles")}>
                      Manage Profiles
                    </Button>
                  </div>
                </Paper>
              </Grid>
            </Grid>
            <PerformanceCalendar style={{ height: "40rem", margin: "2rem 0 0" }} />
          </Paper>
        </Grid>
        <Grid item lg={6} xs={12}>
          <Paper className={classes.paper}>
            <MesheryMetrics
              boardConfigs={grafana.selectedBoardsConfigs}
              grafanaURL={grafana.grafanaURL}
              grafanaAPIKey={grafana.grafanaAPIKey}
              handleGrafanaChartAddition={() => router.push("/settings/#metrics")}
            />
          </Paper>
        </Grid>
      </Grid>

      <GenericModal
        open={!!runTest}
        Content={
          <Paper style={{ margin: "auto", maxWidth: "90%", outline: "none" }}>
            <MesheryPerformanceComponent />
          </Paper>
        }
        handleClose={() => setRunTest(false)}
      />
    </>
  );
}

const mapStateToProps = (st) => {
  const grafana = st.get("grafana").toJS();
  return { grafana: { ...grafana, ts: new Date(grafana.ts) } };
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(Dashboard))));
