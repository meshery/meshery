import makeStyles from "@mui/styles/makeStyles";

export const getStyles = (theme) => ({
  secondaryBar: { zIndex: 0 },
  menuButton: { marginLeft: -theme.spacing(1) },
  iconButtonAvatar: { padding: 4 },
  link: {
    textDecoration: "none",
    color: "rgba(255, 255, 255, 0.7)",
    "&:hover": { color: theme.palette.common.white },
  },
  button: { borderColor: "rgba(255, 255, 255, 0.7)" },
  notifications: {
    paddingLeft: theme.spacing(4),
    paddingRight: theme.spacing(0),
    marginLeft: theme.spacing(4),
  },
  userContainer: {
    paddingLeft: 1,
    display: "flex",
    backgroundColor: "#396679",
  },
  userSpan: { marginLeft: theme.spacing(1) },
  pageTitleWrapper: {
    flexGrow: 1,
    marginRight: "auto",
  },
  betaBadge: { color: "#EEEEEE", fontWeight: "300", fontSize: "13px" },
  pageTitle: {
    paddingLeft: theme.spacing(2),
    fontSize: "1.25rem",
    [theme.breakpoints.up("sm")]: { fontSize: "1.65rem" },
  },
  appBarOnDrawerOpen: {
    backgroundColor: "#396679",
    shadowColor: " #808080",
    zIndex: theme.zIndex.drawer + 1,
    [theme.breakpoints.between(635, 732)]: { padding: theme.spacing(0.75, 1.4) },
    [theme.breakpoints.between(600, 635)]: { padding: theme.spacing(0.4, 1.4) },
  },
  appBarOnDrawerClosed: {
    backgroundColor: "#396679",
    zIndex: theme.zIndex.drawer + 1,
  },
  toolbarOnDrawerClosed: {
    minHeight: 59,
    padding: theme.spacing(2.4),
    paddingLeft: 34,
    paddingRight: 34,
    backgroundColor: "#396679",
  },
  toolbarOnDrawerOpen: {
    minHeight: 58,
    padding: theme.spacing(2.4),
    paddingLeft: 34,
    paddingRight: 34,
    backgroundColor: "#396679",
    [theme.breakpoints.between(620, 732)]: { minHeight: 68, paddingLeft: 20, paddingRight: 20 },
  },
  itemActiveItem: { color: "#00B39F" },
  headerIcons: { fontSize: "1.5rem", height: "1.5rem", width: "1.5rem" },
  cbadge: {
    fontSize: "0.65rem",
    backgroundColor: "white",
    borderRadius: "50%",
    color: "black",
    height: "1.30rem",
    width: "1.30rem",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    position: "absolute",
    zIndex: 1,
    right: "-0.75rem",
    top: "-0.29rem",
  },
  cbadgeContainer: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    position: "relative",
  },
  icon: {
    width: 24,
    height: 24,
  },
  Chip: {
    backgroundColor: "white",
    cursor: "pointer",
  },
  cMenuContainer: {
    backgroundColor: "revert",
    marginTop: "-0.7rem",
    borderRadius: "3px",
    padding: "1rem",
    zIndex: 1201,
    boxShadow: "20px #979797",
    transition: "linear .2s",
    transitionProperty: "height",
  },
  alertEnter: {
    opacity: "0",
    transform: "scale(0.9)",
  },
  alertEnterActive: {
    opacity: "1",
    transform: "translateX(0)",
    transition: "opacity 300ms, transform 300ms",
  },
  chip: {
    margin: "0.25rem 0",
  },
  AddIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  searchIcon: {
    width: theme.spacing(3.5),
  },
});
export const useStyles = makeStyles(getStyles);
