import { createTheme } from '@material-ui/core/styles';
import { blueGrey } from '@material-ui/core/colors';

const drawerWidth = 256;
export var darkTheme = createTheme({
  typography : {
    useNextVariants : true,
    h5 : {
      fontWeight : 'bolder',
      fontSize : 26,
      color : '#FFF',
      letterSpacing : 0.5,
    },
    p : {
      color : '#FFF',
    },
    h6 : {
      color : '#FFF',
    }
  },
  palette : {
    type : "dark",
    primary : blueGrey,
    secondary : {
      main : '#EE5351',
      dark : '#121212',
      mainBackground : '#202020',
    },
  },
  p : {
    color : '#FFF',
  },
  shape : { borderRadius : 8, },
  breakpoints : {
    values : {
      xs : 0,
      sm : 600,
      md : 960,
      lg : 1280,
      xl : 1920,
    },
  },
});
darkTheme = {
  ...darkTheme,
  overrides : {
    MuiDrawer : { paper : { backgroundColor : '#263238', }, },
    MuiButton : {
      label : { textTransform : 'initial', },
      contained : {
        boxShadow : 'none',
        '&:active' : { boxShadow : 'none', },
      },
      containedPrimary : {
        backgroundColor : "#00B39F",
        '&:hover' : { backgroundColor : "#00D3A9", },
      },
    },
    MuiToggleButton : {
      label : {
        textTransform : 'initial',
        color : '#00B39F',
      },
    },
    MuiTabs : {
      root : { marginLeft : darkTheme.spacing(1), },
      indicator : {
        height : 3,
        borderTopLeftRadius : 3,
        borderTopRightRadius : 3,
      },
    },
    MuiTab : {
      root : {
        textTransform : 'initial',
        margin : '0 16px',
        minWidth : 0,
        // [darkTheme.breakpoints.up('md')]: {
        //   minWidth: 0,
        // },
      },
      labelContainer : {
        padding : 0,
        // [darkTheme.breakpoints.up('md')]: {
        //   padding: 0,
        // },
      },
    },
    MuiPaper : { root : { backgroundColor : '#363636' }, elevation2 : { boxShadow : "0px 1px 0px -2px rgb(0 179 159 / 10%), 0px 1px 0px 0px rgb(0 179 159 / 10%), 0px 2px 0px 0px rgb(0 179 159 / 20%)" } },
    MuiIconButton : { root : { padding : darkTheme.spacing(1), }, colorPrimary : { color : "#FFF" }, },
    MuiTooltip : { tooltip : { borderRadius : 4, }, },
    MuiDivider : { root : { backgroundColor : '#404854', }, },
    MuiListItemText : { primary : { fontWeight : darkTheme.typography.fontWeightMedium, }, },
    MuiListItemIcon : {
      root : {
        color : 'inherit',
        marginRight : 0,
        '& svg' : { fontSize : 20, },
        justifyContent : 'center',
        minWidth : 0
      },
    },
    MuiAvatar : {
      root : {
        width : 32,
        height : 32,
      },
    },
    // Global scrollbar styles
    MuiCssBaseline : {
      "@global" : {
        body : {
          scrollbarColor : "#6b6b6b #263238",
          "&::-webkit-scrollbar, & *::-webkit-scrollbar" : {
            backgroundColor : "#263238",
            width : '0.7rem',
          },
          "&::-webkit-scrollbar-thumb, & *::-webkit-scrollbar-thumb" : {
            borderRadius : 8,
            backgroundColor : "#6b6b6b",
            minHeight : 24,
            border : "3px solid #263238",
          },
          "&::-webkit-scrollbar-thumb:focus, & *::-webkit-scrollbar-thumb:focus" : {
            backgroundColor : "#959595",
          },
          "&::-webkit-scrollbar-thumb:active, & *::-webkit-scrollbar-thumb:active" : {
            backgroundColor : "#959595",
          },
          "&::-webkit-scrollbar-thumb:hover, & *::-webkit-scrollbar-thumb:hover" : {
            backgroundColor : "#959595",
          },
          "&::-webkit-scrollbar-corner, & *::-webkit-scrollbar-corner" : {
            backgroundColor : "#263238",
          },
        },
      },
    },
  },
  props : { MuiTab : { disableRipple : true, }, },
  mixins : { ...darkTheme.mixins, },
};

let theme = createTheme({
  overrides : {
    MuiCssBaseline : {
      "@global" : {
        body : {
          backgroundColor : "red"
        }
      }
    }
  },
  typography : {
    useNextVariants : true,
    h5 : {
      fontWeight : 'bolder',
      fontSize : 26,
      letterSpacing : 0.5,
    },
  },
  palette : {
    type : "light",
    // primary: {
    //   light: '#cfd8dc',
    //   main: '#607d8b',
    //   dark: '#455a64',
    // },
    primary : blueGrey,
    secondary : {
      main : '#EE5351',
      dark : '#455a64',
      mainBackground : '#396679',
    },
  },
  shape : { borderRadius : 8, },
  breakpoints : {
    values : {
      xs : 0,
      sm : 600,
      md : 960,
      lg : 1280,
      xl : 1920,
    },
  },
});

theme = {
  ...theme,
  overrides : {
    MuiDrawer : { paper : { backgroundColor : '#263238', }, },
    MuiButton : {
      label : { textTransform : 'initial', },
      contained : {
        boxShadow : 'none',
        '&:active' : { boxShadow : 'none', },
      },
    },
    MuiToggleButton : {
      label : {
        textTransform : 'initial',
        color : '#607d8b',
      },
    },
    MuiTabs : {
      root : { marginLeft : theme.spacing(1), },
      indicator : {
        height : 3,
        borderTopLeftRadius : 3,
        borderTopRightRadius : 3,
      },
    },
    MuiTab : {
      root : {
        textTransform : 'initial',
        margin : '0 16px',
        minWidth : 0,
        // [theme.breakpoints.up('md')]: {
        //   minWidth: 0,
        // },
      },
      labelContainer : {
        padding : 0,
        // [theme.breakpoints.up('md')]: {
        //   padding: 0,
        // },
      },
    },
    MuiIconButton : { root : { padding : theme.spacing(1), }, colorPrimary : { color : "#607d8b" } },
    MuiTooltip : { tooltip : { borderRadius : 4, }, },
    MuiDivider : { root : { backgroundColor : '#404854', }, },
    MuiListItemText : { primary : { fontWeight : theme.typography.fontWeightMedium, }, },
    MuiListItemIcon : {
      root : {
        color : 'inherit',
        marginRight : 0,
        '& svg' : { fontSize : 20, },
        justifyContent : 'center',
        minWidth : 0
      },
    },
    MuiAvatar : {
      root : {
        width : 32,
        height : 32,
      },
    },

    // Global scrollbar styles
    MuiCssBaseline : {
      "@global" : {
        body : {
          scrollbarColor : "#DDDDDD #FFF",
          backgroundColor : "red",
          "&::-webkit-scrollbar, & *::-webkit-scrollbar" : {
            backgroundColor : "#DDDDDD",
            width : '0.7rem',
          },
          "&::-webkit-scrollbar-thumb, & *::-webkit-scrollbar-thumb" : {
            borderRadius : 8,
            backgroundColor : "#959595",
            minHeight : 24,
            border : "3px solid #DDDDDD",
          },
          "&::-webkit-scrollbar-thumb:focus, & *::-webkit-scrollbar-thumb:focus" : {
            backgroundColor : "#959595",
          },
          "&::-webkit-scrollbar-thumb:active, & *::-webkit-scrollbar-thumb:active" : {
            backgroundColor : "#959595",
          },
          "&::-webkit-scrollbar-thumb:hover, & *::-webkit-scrollbar-thumb:hover" : {
            backgroundColor : "#959595",
          },
          "&::-webkit-scrollbar-corner, & *::-webkit-scrollbar-corner" : {
            backgroundColor : "#263238",
          },
        },
      },
    },
  },
  props : { MuiTab : { disableRipple : true, }, },
  mixins : { ...theme.mixins, },
};

export default theme


export const notificationColors = {
  error : "#F91313",
  warning : "#F0A303",
  success : "#206D24",
  info : "#2196F3"
};

export const styles = (theme) => ({
  root : {
    display : 'flex',
    minHeight : '100vh',
  },
  drawer : {
    [theme.breakpoints.up('sm')] : {
      width : drawerWidth,
      flexShrink : 0,
    },
    transition : theme.transitions.create('width', {
      easing : theme.transitions.easing.sharp,
      duration : theme.transitions.duration.enteringScreen,
    }),
  },
  drawerCollapsed : {
    [theme.breakpoints.up('sm')] : { width : theme.spacing(8.4) + 1, },
    transition : theme.transitions.create('width', {
      easing : theme.transitions.easing.sharp,
      duration : theme.transitions.duration.leavingScreen,
    }),
    overflowX : 'hidden',
  },
  appContent : {
    flex : 1,
    display : 'flex',
    flexDirection : 'column',
  },
  mainContent : {
    flex : 1,
    padding : '48px 36px 24px',
    backgroundColor : theme.palette.type == 'dark' ? "rgb(20, 20, 20)" : "rgb(250, 250, 250)",
  },
  footer : {
    backgroundColor : "#212121",
    padding : theme.spacing(2),
  },
  footerText : {
    cursor : 'pointer',
    display : 'inline',
    verticalAlign : 'middle',
  },
  footerIcon : {
    display : 'inline',
    verticalAlign : 'top',
  },
  icon : { fontSize : 20, },
  notifSuccess : {
    backgroundColor : "rgba(248, 252, 248) !important",
    color : `${notificationColors.success} !important`, pointerEvents : "auto !important"
  },
  notifInfo : {
    backgroundColor : "rgba(250, 254, 255) !important",
    color : `${notificationColors.info} !important`, pointerEvents : "auto !important"
  },
  notifWarn : {
    backgroundColor : "rgba(240, 163, 3, 0.04) !important",
    color : `${notificationColors.warning} !important`, pointerEvents : "auto !important"
  },
  notifError : {
    backgroundColor : "rgba(255, 250, 250) !important",
    color : `${notificationColors.error} !important`, pointerEvents : "auto !important"
  },
  playgroundFooter : {
    backgroundColor : notificationColors.warning,
    padding : theme.spacing(2),
  }
});