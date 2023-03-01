import * as React from 'react';
import PropTypes from 'prop-types';
import Head from 'next/head';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { CacheProvider } from '@emotion/react';
import theme from '../styles/theme';
import createEmotionCache from '../lib/createEmotionCache';
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import Footer from "../components/Footer";

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
        {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
        <CssBaseline />
        <Box
          sx={{ display : "flex", flexDirection : "column", minHeight : "100vh" }}
        >
          <Box sx={{ display : "flex", flexGrow : 1 }}>
            <Box sx={{ display : "flex", flex : 1, flexDirection : "column" }}>
              <Box
                sx={{
                  flex : 1,
                  padding : "48px 36px 24px",
                  background : "#eaeff1",
                }}
              >
                <Paper>
                  <Component {...pageProps} />
                </Paper>
              </Box>
            </Box>
          </Box>
          <Footer />
        </Box>
      </ThemeProvider>
    </CacheProvider>
  );
}

MyApp.propTypes = {
  Component : PropTypes.elementType.isRequired,
  emotionCache : PropTypes.object,
  pageProps : PropTypes.object.isRequired,
};
