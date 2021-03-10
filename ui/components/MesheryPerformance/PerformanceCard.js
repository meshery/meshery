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
  updatedAt,
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
        <Grid style={{ marginBottom: "3rem" }} container spacing={1} alignContent="space-between" alignItems="center">
          <Grid item xs={8}>
            <Typography variant="h6" component="div">
              {name}
            </Typography>
          </Grid>
          <Grid item xs={4}>
            <div style={{ width: "fit-content", margin: "0 0 0 auto" }}>
              <IconButton onClick={(ev) => genericClickHandler(ev, handleDelete)}>
                <DeleteIcon />
              </IconButton>
              <IconButton onClick={(ev) => genericClickHandler(ev, handleEdit)}>
                <EditIcon />
              </IconButton>
            </div>
          </Grid>
        </Grid>
        <div>
          <Button
            color="primary"
            variant="contained"
            onClick={(ev) => genericClickHandler(ev, handleRunTest)}
            style={{ display: "block", marginLeft: "auto" }}
          >
            Run Test
          </Button>
        </div>
      </>

      {/* BACK PART */}
      <>
        <Typography variant="h6" gutterBottom>
          Details
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
        {updatedAt ? (
          <div>
            <b>Last Updated At:</b> <Moment format="LLLL">{updatedAt}</Moment>
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
              CustomHeader={<Typography variant="h6">Performance Profile Results</Typography>}
              // @ts-ignore
              endpoint={`/api/user/performance/profiles/${id}/results`}
            />
          </div>
        ) : null}
      </>
    </FlipCard>
  );
}

// @ts-ignore
export default PerformanceCard;
