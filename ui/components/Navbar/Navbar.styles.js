import makeStyles from "@mui/styles/makeStyles";

export const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    minHeight: "100vh",
  },
  drawer: {
    [theme.breakpoints.up("sm")]: {
      width: "256px",
      flexShrink: 0,
    },
    transition: theme.transitions.create("width", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  drawerCollapsed: {
    [theme.breakpoints.up("sm")]: {
      width: theme.spacing(8),
    },
    transition: theme.transitions.create("width", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    overflowX: "hidden",
  },
  appContent: {
    flex: 1,
    display: "flex",
    flexDirection: "column",
  },
  mainContent: {
    flex: 1,
    padding: "48px 36px 24px",
    background: "#eaeff1",
  },
  footer: {
    backgroundColor: theme.palette.background.paper,
    padding: theme.spacing(2),
    color: "#737373",
  },
  footerText: {
    cursor: "pointer",
    display: "inline",
    verticalAlign: "middle",
  },
  footerIcon: {
    display: "inline",
    verticalAlign: "top",
  },
  icon: {
    fontSize: 20,
  },

  // Navigator
  categoryHeader: {
    paddingTop: 16,
    paddingBottom: 16,
  },
  categoryHeaderPrimary: {
    color: theme.palette.common.white,
  },
  item: {
    paddingTop: 4,
    paddingBottom: 4,
    color: "rgba(255, 255, 255, 0.7)",
    fill: "#fff",
    "&:hover": {
      "& $expandMoreIcon": {
        opacity: 1,
        transition: "opacity 200ms ease-in",
      },
    },
  },
  itemCategory: {
    backgroundColor: "#263238",
    boxShadow: "0 -1px 0 #404854 inset",
    paddingTop: 16,
    paddingBottom: 16,
  },
  firebase: {
    top: 0,
    position: "sticky",
    zIndex: 5,
  },
  link: {
    display: "inline-flex",
    width: "100%",
    height: "30px",
    alignItems: "self-end",
  },
  itemActionable: {
    "&:hover": {
      backgroundColor: "rgba(255, 255, 255, 0.08)",
    },
  },
  itemActiveItem: {
    color: "#4fc3f7",
    fill: "#4fc3f7",
  },
  itemPrimary: {
    color: "inherit",
    fontSize: theme.typography.fontSize,
    "&$textDense": {
      fontSize: theme.typography.fontSize,
    },
  },
  textDense: {},
  divider: {
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
  mainLogo: {
    marginRight: theme.spacing(1),
    marginTop: theme.spacing(1),
    marginLeft: theme.spacing(-1),
    width: 40,
    height: 40,
    borderRadius: "unset",
  },
  mainLogoText: {
    marginLeft: theme.spacing(0.5),
    marginTop: theme.spacing(1),
    width: 170,
    height: "100%",
    borderRadius: "unset",
  },
  mainLogoCollapsed: {
    marginRight: theme.spacing(1),
    marginTop: theme.spacing(1),
    marginLeft: theme.spacing(-0.5),
    width: 40,
    height: 40,
    borderRadius: "unset",
  },
  mainLogoTextCollapsed: {
    marginLeft: theme.spacing(1),
    marginTop: theme.spacing(1),
    width: 170,
    height: "100%",
    borderRadius: "unset",
  },
  settingsIcon: {
    marginLeft: theme.spacing(2),
  },
  cursorPointer: {
    cursor: "pointer",
  },
  listIcon: {
    minWidth: theme.spacing(3.5),
    paddingTop: theme.spacing(0.5),
    textAlign: "center",
    display: "inline-table",
    paddingRight: theme.spacing(0.5),
    marginLeft: theme.spacing(0.3),
  },
  listIcon1: {
    minWidth: theme.spacing(3.5),
    paddingTop: theme.spacing(0.5),
    textAlign: "center",
    display: "inline-table",
    paddingRight: theme.spacing(0.5),
    opacity: 0.5,
  },
  listIconSlack: {
    minWidth: theme.spacing(3.5),
    paddingTop: theme.spacing(0.5),
    textAlign: "center",
    display: "inline-table",
    marginLeft: theme.spacing(-0.1),
    paddingRight: theme.spacing(0.5),
    opacity: 0.5,
  },
  nested1: {
    paddingLeft: theme.spacing(3),
  },
  nested2: {
    paddingLeft: theme.spacing(5),
  },
  istioIcon: {
    width: theme.spacing(1.8),
  },
  isHidden: {
    opacity: 0,
    transition: "opacity 200ms ease-in-out",
  },
  isDisplayed: {
    opacity: 1,
    transition: "opacity 200ms ease-in-out",
  },
  sidebarCollapsed: {
    transition: theme.transitions.create("width", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    overflowX: "hidden",
    width: theme.spacing(8),
  },
  sidebarExpanded: {
    width: "256px",
    overflowX: "hidden",
    transition: theme.transitions.create("width", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  fixedSidebarFooter: {
    display: "none",
    [theme.breakpoints.up("sm")]: {
      display: "flex",
      flexDirection: "column",
      "margin-top": "auto",
      "margin-bottom": "0.5rem",
    },
  },
  collapseButtonWrapper: {
    width: "auto",
    "margin-left": "auto",
    opacity: "0.7",
    transition: "opacity 200ms linear",
    "&:hover": {
      opacity: 1,
      background: "transparent",
    },
    "&:focus": {
      opacity: 1,
      background: "transparent",
    },
  },
  collapseButtonWrapperRotated: {
    width: "auto",
    "margin-left": "auto",
    marginRight: theme.spacing(1),
    opacity: "0.7",
    transition: "opacity 200ms linear",
    transform: "rotate(180deg)",
    justifyContent: "center",
    alignSelf: "baseline",
    marginLeft: "3px",
    "&:hover": {
      opacity: 1,
      background: "transparent",
    },
    "&:focus": {
      opacity: 1,
      background: "transparent",
    },
  },
  noPadding: {
    paddingLeft: "16px",
    paddingRight: "16px",
  },
  drawerIcons: {
    height: "1.21rem",
    width: "1.21rem",
    fontSize: "1.21rem",
  },
  avatarGroup: {
    "& .MuiAvatarGroup-avatar": {
      border: "none",
    },
  },
  marginLeft: {
    marginLeft: 8,
    "& .MuiListItem-gutters": {
      paddingLeft: 8,
      paddingRight: 8,
    },
  },
  rightMargin: {
    marginRight: 8,
  },
  btnGrpMarginRight: {
    marginRight: 4,
    alignItems: "center",
  },
  helpIcon: {
    color: "#fff",
    opacity: "0.7",
    transition: "opacity 200ms linear",
    "&:hover": {
      opacity: 1,
      background: "transparent",
    },
    "&:focus": {
      opacity: 1,
      background: "transparent",
    },
  },
  extraPadding: {
    paddingTop: 4,
    paddingBottom: 4,
  },
  restrictPointer: {
    pointerEvents: "none",
  },
  expandMoreIcon: {
    opacity: 0,
    cursor: "pointer",
    transform: "translateX(3px)",
    "&:hover": {
      color: "#4fc3f7",
    },
  },
  collapsed: {
    transform: "rotate(180deg) translateX(-3px)",
  },
  collapsedHelpButton: {
    height: "30px",
    marginTop: "-4px",
    transform: "translateX(-1px)",
  },
  rightTranslate: {
    transform: "translateX(0.5px)",
  },
}));

export const drawerIconsStyle = { height: "1.21rem", width: "1.21rem", fontSize: "1.21rem" };

export const externalLinkIconStyle = { width: "1.11rem", fontSize: "1.11rem" };
