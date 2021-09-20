/* eslint-disable react/prop-types */
import React from "react";
import { Grid, Paper, Typography } from "@mui/material";
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";

const CustomPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(2),
}));

export const PaperWithTitle = ({ title, children }) => {
  const theme = useTheme();
  return (
    <CustomPaper>
      <Grid container spacing={1}>
        <Grid item xs={12} sx={{ mb: theme.spacing(2) }}>
          <Typography variant="subtitle1">{title}</Typography>
        </Grid>
        <Grid item xs={12} container>
          {children}
        </Grid>
      </Grid>
    </CustomPaper>
  );
};
export default CustomPaper;
