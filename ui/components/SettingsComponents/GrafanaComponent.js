import React, {useState} from 'react'
import { TextField, Grid } from '@mui/material';
import ReactSelectWrapper from "@/components/ReactSelectWrapper"
import {RightAlignButton} from "@/components/Button"
import {CustomWrapper} from "./PrometheusComponent"

function GrafanaComponent() {
  return (
    <CustomWrapper >
    <Grid container spacing={1}>
      <Grid item xs={12} md={6}>
          <ReactSelectWrapper
            label="Grafana Base URL"
          />
      </Grid>
      <Grid item xs={12} md={6}>
        <TextField
          id="grafanaAPIKey"
          name="grafanaAPIKey"
          label="API Key"
          fullWidth
          variant="outlined"
        />
      </Grid>
      
      <RightAlignButton title="Submit" />
      
    </Grid>
  </CustomWrapper>
  )
}

export default GrafanaComponent