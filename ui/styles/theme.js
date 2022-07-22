import { createTheme } from "@mui/material/styles";
import { blueGrey } from "@mui/material/colors";

const theme = createTheme({
    typography: { useNextVariants: true, h5: { fontWeight: "bolder", fontSize: 26, letterSpacing: 0.5 } },
    palette: {
      // primary: {
      //   light: '#cfd8dc',
      //   main: '#607d8b',
      //   dark: '#455a64',
      // },
      primary: blueGrey,
      secondary: { main: "#EE5351", dark: "#1E2117" },
      white: "#FFFFFF"
    },
    shape: { borderRadius: 8 },
    breakpoints: {
      values: {
        xs: 0,
        sm: 600,
        md: 960,
        lg: 1280,
        xl: 1920,
      },
    },
    components: {
      MuiDrawer: { styleOverrides: { paper: { backgroundColor: "#263238" } } },
      MuiButton: {
        styleOverrides: {
        root: { textTransform: "none" },
        contained: { boxShadow: "none", "&:active": { boxShadow: "none" } },
      }
    },
      MuiToggleButton: { styleOverrides: { label: { textTransform: "initial", color: "#607d8b" } } },
      MuiTabs: {
        styleOverrides: {
        // root: { marginLeft: theme.spacing(1) },
        indicator: { height: 3, borderTopLeftRadius: 3, borderTopRightRadius: 3 },
      }
    },
      MuiTab: {
        styleOverrides: {
        root: {
          textTransform: "initial",
          margin: "0 16px",
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
      }
    },
      // MuiIconButton: { root: { padding: theme.spacing(1) } },
      MuiTooltip: { styleOverrides: { tooltip: { borderRadius: 4 } } },
      MuiDivider: { styleOverrides: { root: { backgroundColor: "#404854" } } },
      MuiDialogTitle: { 
        styleOverrides: {
        root: {
          background: '#607d8b',
          color: '#fff',
          textAlign: 'center',
        },
      }
    },
      // MuiListItemText: { primary: { fontWeight: theme.typography.fontWeightMedium } },
      MuiListItemIcon: {
        styleOverrides: {
        root: {
          color: "inherit",
          marginRight: 0,
          "& svg": { fontSize: 20 },
          justifyContent: "center",
          minWidth: 0,
        },
      }
    },
      MuiAvatar: { root: { width: 32, height: 32 } },
    },
    props: { MuiTab: { disableRipple: true } },
    // mixins: { ...theme.mixins },
  }
);

export default theme;
