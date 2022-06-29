import makeStyles from "@mui/styles/makeStyles";
import { lightColor } from "styles/colors";

export const getStyles = (theme) => ({
  secondaryBar: { zIndex: 0 },
  menuButton: { marginLeft: -theme.spacing(1) },
  link: { textDecoration: "none", color: lightColor, "&:hover": { color: theme.palette.common.white } },
  button: { borderColor: lightColor },
  notifications: { paddingLeft: theme.spacing(4), paddingRight: theme.spacing(0), marginLeft: theme.spacing(4) },
  userContainer: { paddingLeft: 1, display: "flex", alignItems: "center" },
  pageTitleWrapper: { flexGrow: 1, marginRight: "auto" },
  pageTitle: {
    paddingLeft: theme.spacing(2),
    fontSize: "1.25rem",
    [theme.breakpoints.up("sm")]: { fontSize: "1.65rem" },
  },
  appBarOnDrawerOpen: {
    padding: theme.spacing(1.4),
    backgroundColor: "#396679",
    zIndex: theme.zIndex.drawer + 1,
    [theme.breakpoints.between(635, undefined)]: { padding: theme.spacing(0.75, 1.4) },
    [theme.breakpoints.between(600, undefined)]: { padding: theme.spacing(0.4, 1.4) },
  },
  appBarOnDrawerClosed: { padding: theme.spacing(1.4), backgroundColor: "#396679", zIndex: theme.zIndex.drawer + 1 },
  toolbarOnDrawerClosed: { minHeight: 59, paddingLeft: 24, paddingRight: 24 },
  toolbarOnDrawerOpen: {
    minHeight: 58,
    paddingLeft: 20,
    paddingRight: 20,
    [theme.breakpoints.between(620, undefined)]: { minHeight: 68, paddingLeft: 20, paddingRight: 20 },
  },
  itemActiveItem: { color: "#00B39F" },
  headerIcons: { fontSize: "1.5rem", height: "1.5rem", width: "1.5rem" },
});
export const useStyles = makeStyles(getStyles);
