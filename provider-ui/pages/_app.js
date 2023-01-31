<<<<<<< HEAD
import React from 'react';
import App from 'next/app';
import Head from 'next/head';
import { createTheme, withStyles } from '@material-ui/core/styles';
// TODO: check if below CssBaseline import can be removed, its causing 'error  'CssBaseline' is defined but never used  no-unused-vars'
// during cypress e2e test run (i.e. `npm run test`)
// import CssBaseline from '@material-ui/core/CssBaseline';
=======
import * as React from 'react';
>>>>>>> 940f63251e002f2be03174a21b9774ead55848c8
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

<<<<<<< HEAD
const drawerWidth = 256;

const styles = {
  root: {
    display: 'flex',
    minHeight: '100vh',
  },
  drawer: {
    [theme.breakpoints.up('sm')]: {
      width: drawerWidth,
      flexShrink: 0,
    },
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  drawerCollapsed: {
    [theme.breakpoints.up('sm')]: {
      width: theme.spacing(7) + 1,
    },
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    overflowX: 'hidden',
  },
  appContent: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
  },
  mainContent: {
    flex: 1,
    padding: '48px 36px 24px',
    background: '#eaeff1',
  },
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',

  },
  footer: {
    backgroundColor: theme.palette.background.paper,
    padding: theme.spacing(2),
    color: '#737373',
  },
  footerText: {
    cursor: 'pointer',
    display: 'inline',
    verticalAlign: 'middle',
  },
  footerIcon: {
    display: 'inline',
    verticalAlign: 'top',
  },
  icon: {
    fontSize: 20,
  },
};


class MesheryProviderApp extends App {
  constructor() {
    super();
    this.pageContext = getPageContext();

    this.state = {
      mobileOpen: false,
    };
  }

  handleDrawerToggle = () => {
    this.setState((state) => ({ mobileOpen: !state.mobileOpen }));
  };

  handleL5CommunityClick = () => {
    if (typeof window !== 'undefined') {
      const w = window.open('https://layer5.io', '_blank');
      w.focus();
    }
  }

  static async getInitialProps({ Component, ctx }) {
    const pageProps = Component.getInitialProps ? await Component.getInitialProps(ctx) : {}
    return { pageProps };
  }

  componentDidMount() {

  }

  render() {
    const {
      Component, pageProps, classes,
    } = this.props;
    return (
      <NoSsr>
        <Global
          styles={{
            body: {
              margin: 0,
              padding: 0
            }
          }}
        />
        <Head>
          <title>Meshery</title>
        </Head>

        <div className={classes.root}>

          <div className={classes.appContent}>
            <main className={classes.mainContent}>
              <Paper className={classes.paper}>
                <Component pageContext={this.pageContext} {...pageProps} />
              </Paper>
            </main>
            <footer className={classes.footer}>
              <Typography variant="body2" align="center" color="textSecondary" component="p">
                <span onClick={this.handleL5CommunityClick} className={classes.footerText}>
                  Built with
                  {' '}
                  <FavoriteIcon className={classes.footerIcon} />
                  {' '}
                  by the Layer5 Community
                </span>
              </Typography>
            </footer>
          </div>
        </div>

      </NoSsr>
    );
  }
=======
  return (
    <CacheProvider value={emotionCache}>
      <Head>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
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
>>>>>>> 940f63251e002f2be03174a21b9774ead55848c8
}

MyApp.propTypes = {
  Component : PropTypes.elementType.isRequired,
  emotionCache : PropTypes.object,
  pageProps : PropTypes.object.isRequired,
};
