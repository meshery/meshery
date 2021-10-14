/* eslint-disable react/prop-types */
import { pathToPageTitleMapper } from "@/utils/path";
import { StyledEngineProvider, CssBaseline } from "@mui/material";
import { ThemeProvider } from "@mui/material/styles";
import React from "react";
import theme from "../../styles/theme";
import Header from "../Header/Header";
import Navbar from "../Navbar/Navbar";
import Footer from "../Footer/Footer";
import { NavbarContainer } from "../Navbar/NavbarContainer";

export const Layout = ({ children }) => (
  <StyledEngineProvider injectFirst>
    <ThemeProvider theme={theme}>
      <div style={{ minHeight: "100vh", display: "flex" }}>
        <CssBaseline />
        <NavbarContainer render={(props) => <Navbar {...props} />} />
        <div style={{ flex: 1, display: "flex", flexDirection: "column" }}>
          <Header
            drawerOpen={false}
            onDrawerToggle={() => null}
            pageTitle={pathToPageTitleMapper[window.location.pathname]}
          />
          <main style={{ flex: 1, padding: "48px 36px 24px", background: "#eaeff1" }}>{children}</main>

          <Footer />
        </div>
      </div>
    </ThemeProvider>
  </StyledEngineProvider>
);
