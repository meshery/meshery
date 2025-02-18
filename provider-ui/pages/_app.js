import * as React from "react";
import PropTypes from "prop-types";
import Head from "next/head";
import { CacheProvider } from "@emotion/react";
import createEmotionCache from "../lib/createEmotionCache";
import Footer from "../components/Footer";
import { CssBaseline, charcoal, Box, Paper, SistentThemeProvider, styled, useTheme } from "@layer5/sistent";

//styled-components:
const StyledBox = styled(Box)(() => ({
  display : "flex",
  flex : 1,
  flexDirection : "column",
  background : charcoal[10],
}));

// Client-side cache, shared for the whole session of the user in the browser.
const clientSideEmotionCache = createEmotionCache();

const MesheryThemeProvider = ({ children }) => {
  const theme = useTheme();
  const mode = theme.palette.mode;
  return (
    <SistentThemeProvider initialMode={mode}>{children}</SistentThemeProvider>
  );
};

export default function MyApp(props) {
  const { Component, emotionCache = clientSideEmotionCache, pageProps } = props;
  return (
    <CacheProvider value={emotionCache}>
      <Head>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
        <link rel="shortcut icon" href="/provider/favicon.ico" />
        <title>Provider | Meshery</title>
      </Head>

      <MesheryThemeProvider>
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
                  // background: theme.palette.background.brand.default
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
      </MesheryThemeProvider>
    </CacheProvider>
  );
}

MyApp.propTypes = {
  Component : PropTypes.elementType.isRequired,
  emotionCache : PropTypes.object,
  pageProps : PropTypes.object.isRequired,
};
