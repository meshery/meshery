import React from 'react'
import { Box, Button, Paper, Typography, } from "@mui/material";
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
  margin: theme.spacing(5)

}))

function EmptyState({configuration, Button1, Button2}) {

  return (
    <CustomBox>
    <PaperWithoutTitle >
    <CustomWrapper >
      <Typography  sx={{fontSize : "2rem"}} align="center" color="textSecondary">
        No {configuration} Found
      </Typography>
          <Box sx={{flexDirection: "row"}} >
          {Button1}
          {Button2}
          </Box>

        </CustomWrapper>
  </PaperWithoutTitle>
  </CustomBox>
  )
}

export default EmptyState