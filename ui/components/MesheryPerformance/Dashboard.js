//@ts-check
import React, { useState, useEffect } from "react";
import { connect } from "react-redux";
import { withSnackbar } from "notistack";
import { updateProgress } from "../../lib/store";
import { bindActionCreators } from "redux";
import {
  Button, Grid, Paper, Typography
} from "@material-ui/core";
import IconButton from "@material-ui/core/IconButton";
import { makeStyles, useTheme } from "@material-ui/core/styles";
import useMediaQuery from '@material-ui/core/useMediaQuery';
import CloseIcon from "@material-ui/icons/Close";
import { withRouter } from "next/router";
import MesheryMetrics from "../MesheryMetrics";
import PerformanceCalendar from "./PerformanceCalendar";
import GenericModal from "../GenericModal";
import MesheryPerformanceComponent from "./index";
import fetchPerformanceProfiles from "../graphql/queries/PerformanceProfilesQuery";
import fetchAllResults from "../graphql/queries/FetchAllResultsQuery";

// const MESHERY_PERFORMANCE_URL = "/api/user/performance/profiles";
// const MESHERY_PERFORMANCE_TEST_URL = "/api/user/performance/profiles/results";

const useStyles = makeStyles(() => ({
  paper : { padding : "1rem", },
  resultContainer : {
    display : "flex",
    flexDirection : "row",
    justifyContent : "space-between",
    ["@media (max-width: 830px)"] : {
      flexDirection : "column"
    },
  },
  vSep : {
    height : "10.4rem",
    width : "1px",
    background : "black",
    marginTop : "1.1rem",
    bottom : "0" ,
    left : "36%",
    backgroundColor : "#36454f",
    opacity : "0.7",
    ["@media (max-width: 830px)"] : {
      display : "none",
    }
  },
  hSep : {
    display : "none",
    ["@media (max-width: 830px)"] : {
      display : "block",
      width : "100%",
      height : "1px",
      background : "black",
      marginTop : "1.1rem",
      bottom : "0" ,
      left : "36%",
      backgroundColor : "#36454f",
      opacity : "0.7",
    }
  }
}));

function Dashboard({ updateProgress, enqueueSnackbar, closeSnackbar, grafana, router }) {
  const [profiles, setProfiles] = useState({ count : 0, profiles : [] });
  const [tests, setTests] = useState({ count : 0, tests : [] });
  const [runTest, setRunTest] = useState(false);
  const classes = useStyles();

  const theme = useTheme();
  const matches = useMediaQuery(theme.breakpoints.up("xs"))

  if (matches) {
    console.log("HIT")
  }

  /**
   * fetch performance profiles when the page loads
   */
  useEffect(() => {
    fetchTestProfiles();
    fetchTests();
  }, []);

  function fetchTestProfiles() {
    updateProgress({ showProgress : true });

    fetchPerformanceProfiles({
      selector : {
        // default
        pageSize : `10`,
        page : `0`,
        search : ``,
        order : ``,
      },
    }).subscribe({
      next : (res) => {
        // @ts-ignore
        let result = res?.getPerformanceProfiles;
        updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          if (result) {
            setProfiles({ count : result.total_count || 0,
              profiles : result.profiles || [], });
          }
        }
      },
      error : handleError("Failed to Fetch Profiles"),
    });
  }

  function fetchTests() {
    updateProgress({ showProgress : true });

    fetchAllResults({
      selector : {
        // default
        pageSize : `10`,
        page : `0`,
        search : ``,
        order : ``,
        from : ``,
        to : ``,
      },
    }).subscribe({
      next : (res) => {
        // @ts-ignore
        let result = res?.fetchAllResults;
        updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          if (result) {
            setTests({ count : result.total_count || 0,
              tests : result.results || [], });
          }
        }
      },
      error : handleError("Failed to Fetch Results"),
    });
  }

  function handleError(msg) {
    return function (error) {
      updateProgress({ showProgress : false });

      enqueueSnackbar(`${msg}: ${error}`, { variant : "error",
        action : function Action(key) {
          return (
            <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
              <CloseIcon />
            </IconButton>
          );
        },
        autoHideDuration : 8000, });
    };
  }

  return (
    <>
      <Grid container spacing={2} style={{ padding : "0.5rem" }} alignContent="space-around">
        <Grid container item spacing={1} direction="column" lg xs={12}>
          <Grid item>
            <Paper className={classes.paper}>
              <div className={classes.resultContainer}>
                <div className={classes.paper}>
                  <div style={{ display : "flex", alignItems : "center" , height : "6.8rem" }}>
                    <Typography variant="h2" component="div" color="primary" style={{ marginRight : "0.75rem" }}>
                      {(tests.count).toLocaleString('en')}
                    </Typography>
                    <Typography variant="body1" style={{ color : "rgba(0, 0, 0, 0.54)" }} component="div">
                      Results
                    </Typography>
                  </div>
                  <div style={{ margin : "2rem 0 0 auto", width : "fit-content" }}>
                    <Button variant="contained" color="primary" onClick={() => setRunTest(true)}>
                      Run Test
                    </Button>
                  </div>
                </div>
                <div className={classes.vSep} />
                <div className={classes.hSep} />
                <div className={classes.paper}>
                  <div style={{ display : "flex", alignItems : "center", height : "6.8rem" }}>
                    <Typography variant="h2" component="div" color="primary" style={{ marginRight : "0.75rem" }}>
                      {profiles.count}
                    </Typography>
                    <Typography variant="body1" style={{ color : "rgba(0, 0, 0, 0.54)" }} component="div">
                      Profiles
                    </Typography>
                  </div>
                  <div style={{ margin : "2rem 0 0 auto", width : "fit-content" }}>
                    <Button variant="contained" color="primary" onClick={() => router.push("/performance/profiles")}>
                      Manage Profiles
                    </Button>
                  </div>
                </div>
              </div>
            </Paper>
          </Grid>
          <Grid item>
            <Paper className={classes.paper}>
              <PerformanceCalendar style={{ height : "40rem", margin : "2rem 0 0" }} />
            </Paper>
          </Grid>
        </Grid>
        <Grid item lg xs={12}>
          <Paper className={classes.paper} style={{ height : "100%" }}>
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
          <Paper style={{ margin : "auto", maxWidth : "90%", outline : "none" }}>
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
  return { grafana : { ...grafana, ts : new Date(grafana.ts) } };
};

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(Dashboard)));
