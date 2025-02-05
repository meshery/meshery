import * as React from "react";
import PropTypes from "prop-types";
import Head from "next/head";
import { ThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import { CacheProvider } from "@emotion/react";
import theme from "../styles/theme";
import createEmotionCache from "../lib/createEmotionCache";
import Footer from "../components/Footer";
import { UsesSistent } from "../components/SistentWrapper";
import { Box, Paper, styled } from "@layer5/sistent";

//styled-components:
const StyledBox = styled(Box)(({ theme }) => ({
  display : "flex",
  flex : 1,
  flexDirection : "column",
  background : theme.palette.mode === "dark" ? "#202020" : "#fff",
}));

// Client-side cache, shared for the whole session of the user in the browser.
const clientSideEmotionCache = createEmotionCache();

export default function MyApp(props) {
  const { Component, emotionCache = clientSideEmotionCache, pageProps } = props;

  return (
    <CacheProvider value={emotionCache}>
      <Head>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
        <link rel="shortcut icon" href="/provider/favicon.ico" />
        <title>Provider | Meshery</title>
      </Head>

      <ThemeProvider theme={theme}>
        <UsesSistent>
          {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
          <CssBaseline />
          <Box
            sx={{
              display : "flex",
              flexDirection : "column",
              minHeight : "100vh",
            }}
          >
            <Box sx={{ display : "flex", flexGrow : 1 }}>
              <StyledBox>
                <Box
                  sx={{
                    flex : 1,
                    padding : "48px 36px 24px",
                  }}
                >
                  <Paper>
                    <Component {...pageProps} />
                  </Paper>
                </Box>
              </StyledBox>
            </Box>
            <Footer />
          </Box>
        </UsesSistent>
      </ThemeProvider>
    </CacheProvider>
  );
}

MyApp.propTypes = {
  Component : PropTypes.elementType.isRequired,
  emotionCache : PropTypes.object,
  pageProps : PropTypes.object.isRequired,
};
