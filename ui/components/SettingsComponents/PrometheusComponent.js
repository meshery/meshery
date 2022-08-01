import React, {useState} from 'react'
import ReactSelectWrapper from "@/components/ReactSelectWrapper"
import { Button,TextField, Grid } from '@mui/material';
import {RightAlignButton} from "@/components/Button"
import { styled } from "@mui/material/styles";


export const CustomWrapper = styled("div")(({ theme }) => ({
  padding: theme.spacing(6)
}));

function PrometheusComponent() {

  return (
    <CustomWrapper >
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <ReactSelectWrapper  label="Prometheus Base URl" />
                </Grid>
                <RightAlignButton title="Submit" />
                </Grid>
            </CustomWrapper>
  )
}

export default PrometheusComponent