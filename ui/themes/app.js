import { createTheme } from '@material-ui/core/styles';
import { blueGrey } from '@material-ui/core/colors';

const drawerWidth = 256;

let theme = createTheme({
  typography : { useNextVariants : true,
    h5 : { fontWeight : 'bolder',
      fontSize : 26,
      letterSpacing : 0.5, }, },
  palette : {
    // primary: {
    //   light: '#cfd8dc',
    //   main: '#607d8b',
    //   dark: '#455a64',
    // },
    primary : blueGrey,
    secondary : { main : '#EE5351',
      dark : '#1E2117', }, },
  shape : { borderRadius : 8, },
  breakpoints : { values : {
    xs : 0,
    sm : 600,
    md : 960,
    lg : 1280,
    xl : 1920,
  }, },
});

theme = {
  ...theme,
  overrides : {
    MuiDrawer : { paper : { backgroundColor : '#263238', }, },
    MuiButton : { label : { textTransform : 'initial', },
      contained : { boxShadow : 'none',
        '&:active' : { boxShadow : 'none', }, }, },
    MuiToggleButton : { label : { textTransform : 'initial',
      color : '#607d8b', }, },
    MuiTabs : { root : { marginLeft : theme.spacing(1), },
      indicator : { height : 3,
        borderTopLeftRadius : 3,
        borderTopRightRadius : 3, }, },
    MuiTab : { root : { textTransform : 'initial',
      margin : '0 16px',
      minWidth : 0,
      // [theme.breakpoints.up('md')]: {
      //   minWidth: 0,
      // },
    },
    labelContainer : { padding : 0,
      // [theme.breakpoints.up('md')]: {
      //   padding: 0,
      // },
    }, },
    MuiIconButton : { root : { padding : theme.spacing(1), }, },
    MuiTooltip : { tooltip : { borderRadius : 4, }, },
    MuiDivider : { root : { backgroundColor : '#404854', }, },
    MuiListItemText : { primary : { fontWeight : theme.typography.fontWeightMedium, }, },
    MuiListItemIcon : { root : {
      color : 'inherit',
      marginRight : 0,
      '& svg' : { fontSize : 20, },
      justifyContent : 'center',
      minWidth : 0
    }, },
    MuiAvatar : { root : { width : 32,
      height : 32, }, },
    // Global scrollbar styles
    // MuiCssBaseline : {
    //   "@global" : {
    //     body : {
    //       scrollbarColor : "#6b6b6b #263238",
    //       "&::-webkit-scrollbar, & *::-webkit-scrollbar" : {
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
  mixins : { ...theme.mixins, },
};

export default theme


export const notificationColors = {
  error : "#F91313",
  warning : "#F0A303",
  success : "#206D24",
  info : "#2196F3"
};

export const styles = {
  root : { display : 'flex',
    minHeight : '100vh', },
  drawer : { [theme.breakpoints.up('sm')] : { width : drawerWidth,
    flexShrink : 0, },
  transition : theme.transitions.create('width', { easing : theme.transitions.easing.sharp,
    duration : theme.transitions.duration.enteringScreen, }), },
  drawerCollapsed : { [theme.breakpoints.up('sm')] : { width : theme.spacing(8.4) + 1, },
    transition : theme.transitions.create('width', { easing : theme.transitions.easing.sharp,
      duration : theme.transitions.duration.leavingScreen, }),
    overflowX : 'hidden', },
  appContent : { flex : 1,
    display : 'flex',
    flexDirection : 'column', },
  mainContent : { flex : 1,
    padding : '48px 36px 24px',
    background : '#eaeff1', },
  footer : { backgroundColor : theme.palette.background.paper,
    padding : theme.spacing(2),
    color : '#737373', },
  footerText : { cursor : 'pointer',
    display : 'inline',
    verticalAlign : 'middle', },
  footerIcon : { display : 'inline',
    verticalAlign : 'top', },
  icon : { fontSize : 20, },
  notifSuccess : { backgroundColor : "rgba(248, 252, 248) !important",
    color : `${notificationColors.success} !important`, pointerEvents : "auto !important" },
  notifInfo : { backgroundColor : "rgba(250, 254, 255) !important",
    color : `${notificationColors.info} !important`, pointerEvents : "auto !important" },
  notifWarn : { backgroundColor : "rgba(240, 163, 3, 0.04) !important",
    color : `${notificationColors.warning} !important`, pointerEvents : "auto !important" },
  notifError : { backgroundColor : "rgba(255, 250, 250) !important",
    color : `${notificationColors.error} !important`, pointerEvents : "auto !important" },
};