import React from 'react';
import App from 'next/app';
import Head from 'next/head';
import { MuiThemeProvider, createTheme, withStyles } from '@material-ui/core/styles';
// TODO: check if below CssBaseline import can be removed, its causing 'error  'CssBaseline' is defined but never used  no-unused-vars'
// during cypress e2e test run (i.e. `npm run test`)
// import CssBaseline from '@material-ui/core/CssBaseline';
import PropTypes from 'prop-types';
import Paper from '@material-ui/core/Paper';
import { NoSsr, Typography } from '@material-ui/core';
import FavoriteIcon from '@material-ui/icons/Favorite';

import { blueGrey } from '@material-ui/core/colors';
import getPageContext from '../components/PageContext';

let theme = createTheme({
  typography: {
    useNextVariants: true,
    h5: {
      fontWeight: 500,
      fontSize: 26,
      letterSpacing: 0.5,
    },
  },
  palette: {
    primary: blueGrey,
    secondary: {
      main: '#EE5351',
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
        backgroundColor: '#263238',
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
      },
    },
    MuiTab: {
      root: {
        textTransform: 'initial',
        margin: '0 16px',
        minWidth: 0,
        // [theme.breakpoints.up('md')]: {
        //   minWidth: 0,
        // },
      },
      labelContainer: {
        padding: 0,
        // [theme.breakpoints.up('md')]: {
        //   padding: 0,
        // },
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
          <Head>
            <title>Meshery</title>
          </Head>
          <MuiThemeProvider theme={theme}>
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
          </MuiThemeProvider>
      </NoSsr>
    );
  }
}

MesheryProviderApp.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(MesheryProviderApp);
