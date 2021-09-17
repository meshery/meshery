/* eslint-disable react/prop-types */
import { CssBaseline, MuiThemeProvider } from "@material-ui/core";
import React from "react";
import theme from "../../styles/theme";
import Header from "../Header/Header";
import Navbar from "../Navbar/Navbar";
import { NavbarContainer } from "../Navbar/NavbarContainer";

export const Layout = ({ children }) => (
  <MuiThemeProvider theme={theme}>
    <div style={{ minHeight: "100vh", display: "flex" }}>
      <CssBaseline />
      <NavbarContainer render={(props) => <Navbar {...props} />} />
      <div style={{ flex: 1, display: "flex", flexDirection: "column" }}>
        <Header drawerOpen={false} onDrawerToggle={() => null} pageTitle={"Connection Wizard"} />
        <main style={{ flex: 1, padding: "48px 36px 24px", background: "#eaeff1" }}>{children}</main>

        <footer>Hello</footer>
      </div>
    </div>
  </MuiThemeProvider>
);
