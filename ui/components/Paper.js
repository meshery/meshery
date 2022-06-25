/* eslint-disable react/prop-types */
import React from "react";
import { Grid, Paper, Typography } from "@mui/material";
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";

const CustomPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(2),
}));

export const PaperWithTitle = ({ title, children, titleVariant, containerProps }) => {
  const theme = useTheme();
  return (
    <CustomPaper>
      <Grid container spacing={1}>
        <Grid item xs={12} sx={{ mb: theme.spacing(2) }}>
          {typeof title === "string" ? (
            <Typography variant={titleVariant ? titleVariant : "subtitle1"}>{title}</Typography>
          ) : (
            title
          )}
        </Grid>
        <Grid item xs={12} container {...containerProps}>
          {children}
        </Grid>
      </Grid>
    </CustomPaper>
  );
          }

          export const PaperWithoutTitle = ({ children, containerProps }) => {
            return (
              <CustomPaper>
                <Grid container spacing={1}>
                  <Grid item xs={12} container {...containerProps}>
                    {children}
                  </Grid>
                </Grid>
              </CustomPaper>
            );
          };


export default CustomPaper;
