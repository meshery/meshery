import React from 'react'
import { Button,TextField, MenuItem, Grid, Chip, } from '@mui/material';
import ReactSelectWrapper from "@/components/ReactSelectWrapper"
import {RightAlignButton} from "@/components/Button"

function MeshAdapterConfigComponent() {
  return (
    <div style={{padding: "3rem"}}>
       <Grid item xs={12}>
       <ReactSelectWrapper label="Mesh Adapter UR" />
            </Grid>
            <RightAlignButton title="Connect" />
    </div>
  )
}

export default MeshAdapterConfigComponent