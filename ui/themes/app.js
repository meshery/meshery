import { createTheme } from '@material-ui/core/styles';
import { blueGrey } from '@material-ui/core/colors';
import { iconMedium } from '../css/icons.styles';

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
      primeColor : '#303030',
      dark : '#121212',
      titleText : '#FBFBFB',
      text : '#FFF',
      text2 : '#7494a1',
      text3 : '#FFF',
      textMain : "#F6F8F8",
      titleBackground : "#000",
      mainBackground : '#202020',
      mainBackground2 : '#303030',
      elevatedComponents : '#202020',
      elevatedComponents2 : '#303030',
      elevatedComponents3 : '#303030',
      lightText : 'rgba(255, 255, 255, 0.54)',
      icon : 'rgba(255, 255, 255, 0.54)',
      icon2 : '#E6E6E6',
      iconMain : '#F6F8F8',
      disabledIcon : 'rgba(255, 255, 255, 0.26)',
      chevron : 'rgb(255, 255, 255, 0.2)',
      link : 'rgba(255, 255, 255, 0.7)',
      link2 : "#05FFCD",
      headerColor : '#202020',
      sideBar : '#1A1A1A',
      drawer : '#252E31',
      drawerHover : '#202020',
      img : 'invert(0.8)',
      appBar : '#363636',
      number : '#eee',
      completeInvert : 'invert(1)',
      canvas : '#1A1A1A',
      brightness : 'brightness(200%)',
      switcherButtons : '#1e1e1e',
      honeyComb : '#202020',
      filterChipBackground : '#222222',
      tabs : '#202020',
      modalTabs : '#363636',
      tabHover : '#212121',
      confirmationModal : '#111111',
      focused : '#00b39f',
      primaryModalText : '#FFF',
      penColorPrimary : '#E6E6E6',
      penColorSecondary : '#E6E6E6'
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
    MuiSvgIcon : {
      root : {
        ...iconMedium
      }
    },
    MuiOutlinedInput : {
      root : {
        "&:hover $notchedOutline" : {
          borderColor : "#00B39F",
        },
        "&$focused $notchedOutline" : {
          borderColor : "#00B39F",
        },
      },
    },
    MuiCheckbox : {
      colorPrimary : {
        "&$checked" : {
          color : "rgba(255, 255, 255, 0.7)",
        }
      },
    },
    MuiDrawer : { paper : { backgroundColor : '#263238', }, },
    MuiFormLabel : {
      root : {
        "&$focused" : {
          color : "#00B39F",
        },
      }
    },
    MuiButton : {
      label : { textTransform : 'initial', },
      contained : {
        boxShadow : 'none',
        color : "rgba(255, 255, 255, 0.87)",
        backgroundColor : "#3C494F",
        '&:hover' : { backgroundColor : "#505b61", },
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
        backgroundColor : "#00B39F",
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
      selected : {
        color : "#00B39F !important"
      }
    },
    MuiPaper : { root : { backgroundColor : '#363636' }, elevation2 : { boxShadow : "0px 4px 0px -2px rgb(0 179 159 / 10%), 0px 4px 0px 0px rgb(0 179 159 / 10%), 0px 2px 0px 0px rgb(0 179 159 / 20%)" } },
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
    // Global scrollbar and body styles
    // MuiCssBaseline : {
    //   "@global" : {
    //     body : {
    //       backgroundColor : "#303030",
    //       scrollbarColor : "#6b6b6b #263238",
    //       "&::-webkit-scrollbar, & *::-webkit-scrollbar" : {
    //         backgroundColor : "#263238",
    //         width : '0.7rem',
    //       },
    //       "&::-webkit-scrollbar-track, & *::-webkit-scrollbar-track" : {
    //         backgroundColor : "#263238",
    //         width : '0.7rem',
    //       },
    //       "&::-webkit-scrollbar-thumb, & *::-webkit-scrollbar-thumb" : {
    //         borderRadius : 8,
    //         backgroundColor : "#6b6b6b",
    //         minHeight : 24,
    //         border : "3px solid #263238",
    //       },
    //       "&::-webkit-scrollbar-thumb:focus, & *::-webkit-scrollbar-thumb:focus" : {
    //         backgroundColor : "#959595",
    //       },
    //       "&::-webkit-scrollbar-thumb:active, & *::-webkit-scrollbar-thumb:active" : {
    //         backgroundColor : "#959595",
    //       },
    //       "&::-webkit-scrollbar-thumb:hover, & *::-webkit-scrollbar-thumb:hover" : {
    //         backgroundColor : "#959595",
    //       },
    //       "&::-webkit-scrollbar-corner, & *::-webkit-scrollbar-corner" : {
    //         backgroundColor : "#263238",
    //       },
    //     },
    //   },
    // },
  },
  props : { MuiTab : { disableRipple : true, }, },
  mixins : { ...darkTheme.mixins, },
};

let theme = createTheme({
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
      primeColor : '#ebeff1',
      dark : '#455a64',
      titleText : '#7494A1',
      text : '#000',
      text2 : 'rgba(57, 102, 121, .9)',
      text3 : '#333333',
      textMain : "#3C494F",
      titleBackground : "rgba(57, 102, 121, .1)",
      mainBackground : '#396679',
      mainBackground2 : '#FFF',
      elevatedComponents : '#FFF',
      elevatedComponents2 : "#eaeff1",
      elevatedComponents3 : '#FFF',
      lightText : 'rgba(0, 0, 0, 0.54)',
      icon : 'rgba(0, 0, 0, 0.54)',
      icon2 : 'gray',
      iconMain : '#3C494F',
      disabledIcon : 'rgba(0, 0, 0, 0.26)',
      chevron : '#FFF',
      link : '#000',
      link2 : "#00b39F",
      headerColor : '#eeeeee',
      sideBar : '#FFF',
      drawer : '#FFF',
      drawerHover : '#f2f5f7',
      img : 'none',
      appBar : '#FFF',
      number : '#607d8b',
      completeInvert : 'none',
      canvas : '#fff',
      brightness : 'none',
      switcherButtons : '#335c6d',
      honeyComb : '#F0F0F0',
      filterChipBackground : '#CCCCCC',
      tabs : '#eeeeee87',
      modalTabs : '#dddddd',
      tabHover : '#e3e3e3',
      confirmationModal : 'rgb(234, 235, 236)',
      focused : '#607d8b',
      primaryModalText : '#FFF',
      penColorPrimary : '#3C494F',
      penColorSecondary : '#677E88'
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
    MuiSvgIcon : {
      root : {
        ...iconMedium
      }
    },
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

    // global style for body throughout meshery-ui
    MuiCssBaseline : {
      "@global" : {
        body : {
          backgroundColor : "#eaeff1",
        },
      },
    },
    // Global scrollbar styles
    // MuiCssBaseline : {
    //   "@global" : {

    //     body : {
    //       backgroundColor : "#eaeff1",
    //       scrollbarColor : "#DDDDDD #FFF",
    //       "&::-webkit-scrollbar, & *::-webkit-scrollbar" : {
    //         backgroundColor : "#DDDDDD",
    //         width : '0.7rem',
    //       },
    //       "&::-webkit-scrollbar-thumb, & *::-webkit-scrollbar-thumb" : {
    //         borderRadius : 8,
    //         backgroundColor : "#959595",
    //         minHeight : 24,
    //         border : "3px solid #DDDDDD",
    //       },
    //       "&::-webkit-scrollbar-thumb:focus, & *::-webkit-scrollbar-thumb:focus" : {
    //         backgroundColor : "#959595",
    //       },
    //       "&::-webkit-scrollbar-thumb:active, & *::-webkit-scrollbar-thumb:active" : {
    //         backgroundColor : "#959595",
    //       },
    //       "&::-webkit-scrollbar-thumb:hover, & *::-webkit-scrollbar-thumb:hover" : {
    //         backgroundColor : "#959595",
    //       },
    //     },
    //   },
    // },
  },
  props : { MuiTab : { disableRipple : true, }, },
  mixins : { ...theme.mixins, },
};

export default theme


export const notificationColors = {
  error : "#F91313",
  warning : "#F0A303",
  success : "#206D24",
  info : "#2196F3",
  darkRed : "#B32700"
};

export const darkNotificationColors = {
  error : "#F91313",
  warning : "#F0D053",
  success : "#78C57C",
  info : "#5FD4FF"
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
  },
  footer : {
    backgroundColor : '#fff',
    padding : theme.spacing(2),
  },
  footerDark : {
    backgroundColor : '#202020',
    padding : theme.spacing(2),
  },
  footerText : {
    cursor : 'pointer',
    display : 'inline',
    verticalAlign : 'middle',
  },
  footerIcon : {
    display : 'inline',
    verticalAlign : 'bottom'
  },
  icon : { fontSize : 20, },
  notifSuccess : {
    backgroundColor : "rgb(248, 252, 248) !important",
    color : `${notificationColors.success} !important`, pointerEvents : "auto !important"
  },
  notifInfo : {
    backgroundColor : "rgb(248, 252, 248) !important",
    color : `${notificationColors.info} !important`, pointerEvents : "auto !important"
  },
  notifWarn : {
    backgroundColor : "rgba(240, 163, 3, 0.04) !important",
    color : `${notificationColors.warning} !important`, pointerEvents : "auto !important"
  },
  notifError : {
    backgroundColor : "rgb(248, 252, 248) !important",
    color : `${notificationColors.error} !important`,
    pointerEvents : "auto !important"
  },
  darknotifSuccess : {
    backgroundColor : "#323232 !important",
    color : `${darkNotificationColors.success} !important`,
    pointerEvents : "auto !important"
  },
  darknotifInfo : {
    backgroundColor : "#323232 !important",
    color : `${darkNotificationColors.info} !important`,
    pointerEvents : "auto !important"
  },
  darknotifWarn : {
    backgroundColor : "#323232 !important",
    color : `${darkNotificationColors.warning} !important`,
    pointerEvents : "auto !important"
  },
  darknotifError : {
    backgroundColor : "#323232 !important",
    color : `${darkNotificationColors.error} !important`,
    pointerEvents : "auto !important"
  },
  playgroundFooter : {
    backgroundColor : notificationColors.warning,
    padding : theme.spacing(2),
  }
});
