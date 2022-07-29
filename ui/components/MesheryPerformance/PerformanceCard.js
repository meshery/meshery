import React, { useState } from 'react'
import { Button, Grid, IconButton, Typography } from "@mui/material";
import FlipCard from "@/components/FlipCard";
import Moment from "react-moment";
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

function PerformanceCard({
    name,
    endpoints,
    loadGenerators,
    testRunDuration,
    lastRun,
    reqHeaders,
    results,
    concurrentRequest,
    qps,
    serviceMesh,
    contentType,
    requestBody,
    requestCookies,
    requestHeaders,
}) {
    const [renderTable, setRenderTable] = useState(false);
  return (
    <FlipCard>
          {/* FRONT PART */}
          <>
        <div>
          <Typography variant="h6" component="div">
            {name}
          </Typography>
          <div>
              <Typography variant="h2" component="div" color="primary" style={{ marginRight : "0.75rem" }}>
                {(results).toLocaleString('en')}
              </Typography>
              <Typography variant="body1" style={{ color : "rgba(0, 0, 0, 0.54)" }} component="div">
                Results
              </Typography>
          </div>
        </div>
        <div>
          <div>
              {lastRun
                ? (
                  <Typography color="primary" variant="caption" style={{ fontStyle : "italic" }}>
                  Last Run: <Moment format="LLL">{lastRun}</Moment>
                  </Typography>
                )
                : null}
          </div>
          <div>
            <Button variant="contained" >
              {renderTable
                ? "Hide"
                : "View"} Results
            </Button>
            <Button color="primary" variant="contained" >
              Run Test
            </Button>
          </div>
        </div>
      </>

      {/* BACK PART */}
      <>
        <Grid
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
            <div>
              <IconButton>
                <EditIcon color="primary" />
              </IconButton>
              <IconButton>
                <DeleteIcon color="primary" />
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

    </FlipCard>
  )
}

export default PerformanceCard