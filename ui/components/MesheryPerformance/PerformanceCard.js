//@ts-check
import React, { useState } from "react";
import { Button, Grid, IconButton, Typography } from "@material-ui/core";
import DeleteIcon from "@material-ui/icons/Delete";
import EditIcon from "@material-ui/icons/Edit";
import Moment from "react-moment";
import PerformanceResults from "./PerformanceResults";
import FlipCard from "../FlipCard";

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
}) {
  const [renderTable, setRenderTable] = useState(false);

  function genericClickHandler(ev, fn) {
    ev.stopPropagation();
    fn();
  }

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
        <Grid
          style={{ marginBottom: "0.25rem" }}
          container
          spacing={1}
          alignContent="space-between"
          alignItems="center"
        >
          <Grid item xs={8}>
            <Typography variant="h6" component="div">
              {name}
            </Typography>
          </Grid>
          <Grid item xs={4}>
            <div style={{ width: "fit-content", margin: "0 0 0 auto" }}>
              <IconButton onClick={(ev) => genericClickHandler(ev, handleEdit)}>
                <EditIcon />
              </IconButton>
              <IconButton onClick={(ev) => genericClickHandler(ev, handleDelete)}>
                <DeleteIcon />
              </IconButton>
            </div>
          </Grid>
        </Grid>
        <div style={{ margin: "0 0 1rem" }}>
          <div style={{ display: "flex", alignItems: "center" }}>
            <Typography variant="h2" component="div" color="primary" style={{ marginRight: "0.75rem" }}>
              {results}
            </Typography>
            <Typography variant="body1" style={{ color: "rgba(0, 0, 0, 0.54)" }} component="div">
              Results
            </Typography>
          </div>
        </div>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
          <div style={{ marginRight: "0.5rem" }}>
            <div>
              {lastRun ? (
                <Typography color="primary" variant="caption" style={{ fontStyle: "italic" }}>
                  Last Run: <Moment format="LLL">{lastRun}</Moment>
                </Typography>
              ) : null}
            </div>
          </div>
          <div style={{ display: "flex", justifyContent: "flex-end", alignItems: "center" }}>
            <Button
              variant="contained"
              onClick={() => {
                // Let this propagate to the flip card which will trigger its
                // onClick handlers resulting in the card flipping automatically
              }}
              style={{ marginRight: "0.5rem", width: "5.7rem" }}
            >
              Details
            </Button>
            <Button
              color="primary"
              variant="contained"
              onClick={(ev) => genericClickHandler(ev, handleRunTest)}
              style={{ width: "5.7rem" }}
            >
              Run Test
            </Button>
          </div>
        </div>
      </>

      {/* BACK PART */}
      <>
        <Typography variant="h6" gutterBottom>
          {name} Details
        </Typography>
        {Array.isArray(endpoints) ? (
          <div>
            <b>Endpoints:</b> {endpoints.join(", ")}
          </div>
        ) : null}
        {Array.isArray(loadGenerators) ? (
          <div>
            <b>Load Generators:</b> {loadGenerators.join(", ")}
          </div>
        ) : null}
        {testRunDuration ? (
          <div>
            <b>Test Run Duration:</b> {testRunDuration}
          </div>
        ) : null}
        {reqHeaders ? (
          <div>
            <b>Request Headers:</b> <code>{reqHeaders}</code>
          </div>
        ) : null}
        <Button
          variant="contained"
          style={{ marginTop: "1.5rem" }}
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
        >
          {renderTable ? "Hide" : "View"} Test Results
        </Button>
        {renderTable ? (
          <div onClick={(ev) => ev.stopPropagation()} style={{ marginTop: "0.5rem" }}>
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
    </FlipCard>
  );
}

// @ts-ignore
export default PerformanceCard;
