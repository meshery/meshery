/* eslint-disable react/prop-types */
import { CssBaseline, MuiThemeProvider } from "@material-ui/core";
import React from "react";
import { theme } from "..";

export const Layout = ({ children }) => (
  <MuiThemeProvider theme={theme}>
    <CssBaseline />
    {children}
  </MuiThemeProvider>
);
