import React from 'react'
import { Button,TextField, MenuItem, Grid, Chip, } from '@mui/material';
import {MetricsButton} from "@/components/Button"
import ReactSelectWrapper from "@/components/ReactSelectWrapper"

function MeshAdapterConfigComponent() {
  return (
    <div style={{padding: "3rem"}}>
       <Grid item xs={12}>
       <ReactSelectWrapper label="Mesh Adapter UR" />
            </Grid>
           <MetricsButton type="submit"
                variant="contained"
                color="primary"
                size="large" >
             Connect
            </MetricsButton>
    </div>
  )
}

export default MeshAdapterConfigComponent