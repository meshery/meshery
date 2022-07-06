import React from 'react'
import { Button, Grid, Chip, } from '@mui/material';
import {MetricsButton} from "@/components/Button"

function MeshAdapterConfigComponent() {
  return (
    <div>
           <MetricsButton                  type="submit"
                variant="contained"
                color="primary"
                size="large" >
             Connect
            </MetricsButton>
    </div>
  )
}

export default MeshAdapterConfigComponent