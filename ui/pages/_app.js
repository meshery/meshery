import React from 'react';
import App, { Container } from 'next/app';
import Head from 'next/head';
import { MuiThemeProvider, createMuiTheme, withStyles } from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';
import JssProvider from 'react-jss/lib/JssProvider';
import getPageContext from '../components/PageContext';
import Navigator from '../components/Navigator';
import Header from '../components/Header';
import PropTypes from 'prop-types';
import Hidden from '@material-ui/core/Hidden';
import Paper from '@material-ui/core/Paper';

let theme = createMuiTheme({
    typography: {
      useNextVariants: true,
      h5: {
        fontWeight: 500,
        fontSize: 26,
        letterSpacing: 0.5,
      },
    },
    palette: {
      primary: {
        light: '#63ccff',
        main: '#009be5',
        dark: '#006db3',
      },
    },
    shape: {
      borderRadius: 8,
    },
  });
  
  theme = {
    ...theme,
    overrides: {
      MuiDrawer: {
        paper: {
          backgroundColor: '#18202c',
        },
      },
      MuiButton: {
        label: {
          textTransform: 'initial',
        },
        contained: {
          boxShadow: 'none',
          '&:active': {
            boxShadow: 'none',
          },
        },
      },
      MuiTabs: {
        root: {
          marginLeft: theme.spacing(1),
        },
        indicator: {
          height: 3,
          borderTopLeftRadius: 3,
          borderTopRightRadius: 3,
          backgroundColor: theme.palette.common.white,
        },
      },
      MuiTab: {
        root: {
          textTransform: 'initial',
          margin: '0 16px',
          minWidth: 0,
          [theme.breakpoints.up('md')]: {
            minWidth: 0,
          },
        },
        labelContainer: {
          padding: 0,
          [theme.breakpoints.up('md')]: {
            padding: 0,
          },
        },
      },
      MuiIconButton: {
        root: {
          padding: theme.spacing(1),
        },
      },
      MuiTooltip: {
        tooltip: {
          borderRadius: 4,
        },
      },
      MuiDivider: {
        root: {
          backgroundColor: '#404854',
        },
      },
      MuiListItemText: {
        primary: {
          fontWeight: theme.typography.fontWeightMedium,
        },
      },
      MuiListItemIcon: {
        root: {
          color: 'inherit',
          marginRight: 0,
          '& svg': {
            fontSize: 20,
          },
        },
      },
      MuiAvatar: {
        root: {
          width: 32,
          height: 32,
        },
      },
    },
    props: {
      MuiTab: {
        disableRipple: true,
      },
    },
    mixins: {
      ...theme.mixins,
      toolbar: {
        minHeight: 48,
      },
    },
  };
  
  const drawerWidth = 256;
  console.log("theme.breakpoints: "+ JSON.stringify(theme.breakpoints));
  
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
    },
    appContent: {
      flex: 1,
      display: 'flex',
      flexDirection: 'column',
    },
    mainContent: {
      flex: 1,
      padding: '48px 36px 0',
      background: '#eaeff1',
    },
    paper: {
        maxWidth: 936,
        margin: 'auto',
        overflow: 'hidden',
      },
    //   searchBar: {
    //     borderBottom: '1px solid rgba(0, 0, 0, 0.12)',
    //   },
    //   searchInput: {
    //     fontSize: theme.typography.fontSize,
    //   },
    //   block: {
    //     display: 'block',
    //   },
    //   addUser: {
    //     marginRight: theme.spacing(1),
    //   },
    //   contentWrapper: {
    //     margin: '40px 16px',
    //   },
  };


class MesheryApp extends App {
  constructor() {
    super();
    this.pageContext = getPageContext();
  }

  state = {
    mobileOpen: false,
  };

  handleDrawerToggle = () => {
    this.setState(state => ({ mobileOpen: !state.mobileOpen }));
  };

  componentDidMount() {
    // Remove the server-side injected CSS.
    // const jssStyles = document.querySelector('#jss-server-side');
    // if (jssStyles && jssStyles.parentNode) {
    //   jssStyles.parentNode.removeChild(jssStyles);
    // }
  }

  render() {
    const { Component, pageProps, classes } = this.props;
    return (
      <Container>
        <Head>
          <title>Meshery</title>
        </Head>
        {/* Wrap every page in Jss and Theme providers */}
        {/* <JssProvider
          registry={this.pageContext.sheetsRegistry}
          generateClassName={this.pageContext.generateClassName}
        > */}
          {/* MuiThemeProvider makes the theme available down the React
              tree thanks to React context. */}
          <MuiThemeProvider theme={theme}>
                <div className={classes.root}>
                    {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
                    <CssBaseline />
                    <nav className={classes.drawer}>
                        <Hidden smUp implementation="js">
                        <Navigator
                            PaperProps={{ style: { width: drawerWidth } }}
                            variant="temporary"
                            open={this.state.mobileOpen}
                            onClose={this.handleDrawerToggle}
                        />
                        </Hidden>
                        <Hidden xsDown implementation="css">
                        <Navigator PaperProps={{ style: { width: drawerWidth } }} />
                        </Hidden>
                    </nav>
                    <div className={classes.appContent}>
                        <Header onDrawerToggle={this.handleDrawerToggle} />
                        <main className={classes.mainContent}>
                    {/* Pass pageContext to the _document though the renderPage enhancer
                        to render collected styles on server-side. */}
                            <Paper className={classes.paper}>
                                <Component pageContext={this.pageContext} {...pageProps} />
                            </Paper>
                        </main>
                    </div>
                </div>
          </MuiThemeProvider>
        {/* </JssProvider> */}
      </Container>
    );
  }
}

MesheryApp.propTypes = {
    classes: PropTypes.object.isRequired,
  };
  
export default withStyles(styles)(MesheryApp);