import React, {useState} from 'react'
import { Button,TextField, Grid } from '@mui/material';
import ReactSelectWrapper from "@/components/ReactSelectWrapper"
import {RightAlignButton} from "@/components/Button"

function GrafanaComponent() {
  return (
    <div style={{padding: "3rem"}} >
    <Grid container spacing={1}>
      <Grid item xs={12} md={6}>
          <ReactSelectWrapper
            label="Grafana Base URL"
            // placeholder="Address of Grafana Server"
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
  </div>
  )
}

export default GrafanaComponent