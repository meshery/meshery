import React, {useState} from 'react'
import ReactSelectWrapper from "@/components/ReactSelectWrapper"
import { Button,TextField, Grid } from '@mui/material';
import {RightAlignButton} from "@/components/Button"

function PrometheusComponent() {
  return (
    <div style={{padding:"3rem" }} >
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <ReactSelectWrapper  label="Prometheus Base URl" />
                </Grid>
                <RightAlignButton title="Submit" />
                </Grid>
            </div>
  )
}

export default PrometheusComponent