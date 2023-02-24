//@ts-check
import React, { useState } from "react";
import {
  Button, Grid, IconButton, Typography
} from "@material-ui/core";
import DeleteIcon from "@material-ui/icons/Delete";
import EditIcon from "@material-ui/icons/Edit";
import Moment from "react-moment";
import PerformanceResults from "./PerformanceResults";
import FlipCard from "../FlipCard";
import { makeStyles } from "@material-ui/core/styles";
import { iconMedium } from "../../css/icons.styles";

const useStyles= makeStyles((theme) => ({
  cardButtons : {
    display : "flex",
    justifyContent : "flex-end",
    alignItems : "center"
  },
  testsButton : {
    marginRight : "0.5rem"
  },
  perfResultsContainer : {
    marginTop : "0.5rem"
  },
  backGrid : {
    marginBottom : "0.25rem",
    minHeight : "6rem"
  },
  deleteEditButton : {
    width : "fit-content",
    margin : "0 0 0 auto"
  },
  noOfResultsContainer : {
    margin : "0 0 1rem",
    '& div' : {
      display : "flex",
      alignItems : "center"
    },
  },
  bottomPart : {
    display : "flex",
    justifyContent : "space-between",
    alignItems : "center",
  },
  lastRunText : {
    marginRight : "0.5rem"
  },
  resultText : {
    color : theme.palette.secondary.lightText,
  },
}))

function PerformanceCard({
  id,
  name,
  endpoints,
  loadGenerators,
  testRunDuration,
  lastRun,
  reqHeaders,
  results,
  handleDelete,
  handleEdit,
  handleRunTest,
  requestFullSize,
  requestSizeRestore,
  concurrentRequest,
  qps,
  serviceMesh,
  contentType,
  requestBody,
  requestCookies,
  requestHeaders,
}) {
  const [renderTable, setRenderTable] = useState(false);

  function genericClickHandler(ev, fn) {
    ev.stopPropagation();
    fn();
  }
  const classes=useStyles()

  return (
    <FlipCard
      onClick={() => {
        setRenderTable(false);
        requestSizeRestore();
      }}
      duration={600}
    >
      {/* FRONT PART */}
      <>
        <div>
          <Typography variant="h6" component="div">
            {name}
          </Typography>
          <div className={classes.noOfResultsContainer} >
            <div >
              <Typography variant="h2" component="div" color="primary" style={{ marginRight : "0.75rem" }}>
                {(results).toLocaleString('en')}
              </Typography>
              <Typography variant="body1" className={classes.resultText} component="div">
                Results
              </Typography>
            </div>
          </div>
        </div>
        <div className={classes.bottomPart} >
          <div className={classes.lastRunText} >
            <div>
              {lastRun
                ? (
                  <Typography color="primary" variant="caption" style={{ fontStyle : "italic" }}>
                  Last Run: <Moment format="LLL">{lastRun}</Moment>
                  </Typography>
                )
                : null}
            </div>
          </div>
          <div className={classes.cardButtons} >
            <Button
              variant="contained"
              onClick={(ev) =>
                genericClickHandler(ev, () => {
                  setRenderTable((renderTable) => {
                    if (renderTable) {
                      requestSizeRestore();
                      return false;
                    }

                    requestFullSize();
                    return true;
                  });
                })
              }
              className={classes.testsButton}
            >
              {renderTable
                ? "Hide"
                : "View"} Results
            </Button>
            <Button
              color="primary"
              variant="contained"
              onClick={(ev) => genericClickHandler(ev, handleRunTest)}
            >
              Run Test
            </Button>
          </div>
        </div>
        {renderTable ? (
          <div onClick={(ev) => ev.stopPropagation()} className={classes.perfResultsContainer} >
            <PerformanceResults
              // @ts-ignore
              CustomHeader={<Typography variant="h6">Test Results</Typography>}
              // @ts-ignore
              endpoint={`/api/user/performance/profiles/${id}/results`}
              // @ts-ignore
              elevation={0}
            />
          </div>
        ) : null}
      </>

      {/* BACK PART */}
      <>
        <Grid className={classes.backGrid}
          container
          spacing={1}
          alignContent="space-between"
          alignItems="center"
        >
          <Grid item xs={8}>
            <Typography variant="h6" gutterBottom>
              {name} Details
            </Typography>
          </Grid>
          <Grid item xs={4}>
            <div className={classes.deleteEditButton} >
              <IconButton onClick={(ev) => genericClickHandler(ev, handleEdit)}>
                <EditIcon style={iconMedium} />
              </IconButton>
              <IconButton onClick={(ev) => genericClickHandler(ev, handleDelete)}>
                <DeleteIcon style={iconMedium} />
              </IconButton>
            </div>
          </Grid>
        </Grid>
        {Array.isArray(endpoints)
          ? (
            <div>
              <b>Endpoints:</b> {endpoints.join(", ")}
            </div>
          )
          : null}
        {Array.isArray(loadGenerators)
          ? (
            <div>
              <b>Load Generators:</b> {loadGenerators.join(", ")}
            </div>
          )
          : null}
        {testRunDuration
          ? (
            <div>
              <b>Test Run Duration:</b> {testRunDuration}
            </div>
          )
          : null}
        {reqHeaders
          ? (
            <div>
              <b>Request Headers:</b> <code>{reqHeaders}</code>
            </div>
          )
          : null}
        {concurrentRequest
          ? (
            <div>
              <b>Concurrent Request:</b> <code>{concurrentRequest}</code>
            </div>
          )
          : null}
        {qps
          ? (
            <div>
              <b>Queries Per Second:</b> <code>{qps}</code>
            </div>
          )
          : null}
        {serviceMesh
          ? (
            <div>
              <b>Service Mesh:</b> <code>{serviceMesh}</code>
            </div>
          )
          : null}
        {
          contentType||requestBody||requestCookies||requestHeaders
            ?(<h4>Advanced Options</h4>)
            :(null)
        }
        {contentType
          ? (
            <div>
              <b>Content Type:</b> <code>{contentType}</code>
            </div>
          )
          : null}
        {requestBody
          ? (
            <div>
              <b>Request Body:</b> <code>{requestBody}</code>
            </div>
          )
          : null}
        {requestCookies
          ? (
            <div>
              <b>Request Cookies:</b> <code>{requestCookies}</code>
            </div>
          )
          : null}
        {requestHeaders
          ? (
            <div>
              <b>Request Headers:</b> <code>{requestHeaders}</code>
            </div>
          )
          : null}
      </>
    </FlipCard >
  );
}

// @ts-ignore
export default PerformanceCard;
