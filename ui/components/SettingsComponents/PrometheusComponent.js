import React, {useState} from 'react'
import ReactSelectWrapper from "@/components/ReactSelectWrapper"
import { Button,TextField, Grid } from '@mui/material';

function PrometheusComponent() {
  return (
    <div style={{padding:"3rem" }} >
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <ReactSelectWrapper
                    label="Prometheus Base URL"
                  />
                </Grid>
                <Grid item xs={4} >
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                >
                  Submit
                </Button>
                </Grid>
                </Grid>
            </div>
  )
}

export default PrometheusComponent