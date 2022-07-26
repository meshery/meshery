import React from 'react'
import { Box, Button, Paper, Typography, } from "@mui/material";
import UploadImport from "@/components/UploadImport";
import {PaperWithoutTitle} from "@/components/Paper";
import { styled } from "@mui/material/styles";

const CustomWrapper = styled(Box)(({theme}) => ({
  padding : theme.spacing(5),
  display : "flex",
  justifyContent : "center",
  alignItems : "center",
  flexDirection : "column",
  margin: "auto",
  gap: theme.spacing(3),
}))
const CustomBox = styled(Box)(({theme}) => ({
  margin: "5rem 5rem"
}))

function EmptyConfigurationList({configuration, NewButton}) {

  return (
    <CustomBox>
    <PaperWithoutTitle >
    <CustomWrapper >
      <Typography  sx={{fontSize : "2rem"}} align="center" color="textSecondary">
        No {configuration} Found
      </Typography>
          <Box sx={{flexDirection: "row"}} >
          {NewButton}
          <UploadImport aria-label="URL upload button" configuration={configuration} /> 
          </Box>

        </CustomWrapper>
  </PaperWithoutTitle>
  </CustomBox>
  )
}

export default EmptyConfigurationList