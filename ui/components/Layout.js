/* eslint-disable react/prop-types */
import { createTheme, CssBaseline, MuiThemeProvider } from "@material-ui/core";
import { Head } from "next/document";
import React from "react";
import { MuiPickersUtilsProvider } from "@material-ui/pickers";
import MomentUtils from "@date-io/moment";

const theme = createTheme({
  palette: {
    newcolor: {
      main: "red",
    },
  },
});
export const Layout = ({ children }) => (
  <MuiThemeProvider theme={theme}>
    <CssBaseline />
    {children}
  </MuiThemeProvider>
);
