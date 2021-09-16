/* eslint-disable react/prop-types */
import { CssBaseline, MuiThemeProvider } from "@material-ui/core";
import React from "react";
import { theme } from "../../styles/theme";
import Header from "../Header/Header";

export const Layout = ({ children }) => (
  <MuiThemeProvider theme={theme}>
    <CssBaseline />
    <Header drawerOpen={false} onDrawerToggle={() => null} pageTitle={"Connection Wizard"} />
    {/* Navbar */}

    {children}
  </MuiThemeProvider>
);
