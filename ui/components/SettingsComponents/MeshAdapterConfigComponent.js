import React, {useState} from 'react'
import { Box,TextField, MenuItem, Grid, Chip, } from '@mui/material';
import { useTheme } from "@mui/system";
import { AdaptersChipList, AdaptersListContainer } from "@/features/mesheryComponents";
import ReactSelectWrapper from "@/components/ReactSelectWrapper"
import {RightAlignButton} from "@/components/Button"

function MeshAdapterConfigComponent() {
  const theme = useTheme();

  return (
    <Box sx={{p: theme.spacing(4)}}>
       <AdaptersListContainer>{(props) => <AdaptersChipList {...props} />}</AdaptersListContainer>
         <ReactSelectWrapper style={{ marginTop: theme.spacing(4)}} label="Mesh Adapter UR" />
    <RightAlignButton title="Connect" />
    </Box>
  )
}

export default MeshAdapterConfigComponent