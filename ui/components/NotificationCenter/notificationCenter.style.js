import { makeStyles } from "@material-ui/core"

export const useStyles = makeStyles((theme) => ({
  sidelist: {
    width: "45rem",
    maxWidth: "90vw",
  },
  notificationButton: { height: "100%" },
  notificationDrawer: {
    backgroundColor: theme.palette.secondary.drawer,
    display: "flex",
    flexDirection: "column",
    justifyContent: "space-between",
  },
  drawerButton: {
    padding: "0.45rem",
    margin: "0.2rem",
    backgroundColor: theme.palette.secondary.dark,
    color: "#FFFFFF",
    "&:hover": {
      backgroundColor: "#FFFFFF",
      color: theme.palette.secondary.dark,
    },
  },
  fullView: {
    right: 0,
    transition: "0.3s ease-in-out !important",
  },
  peekView: {
    right: "-42.1rem",
    transition: "0.3s ease-in-out !important",
  },

  container: {
    padding: "20px"
  },
  header: {
    display: "flex",
    gap: "0.5rem",
    justifyContent: "space-between",
    alignItems: "center",
    background: theme.palette.secondary.headerColor,
  },
  title: {
    display: "flex",
    alignItems: "center",
    gap: "0.5rem",
  },
  titleBellIcon: {
    width: "36px",
    height: "36px",
    borderRadius: "100%",
    backgroundColor: "black",
    display: "flex",
    padding: "0.2rem",
    justifyContent: "center",
    alignItems: "center"
  },
  severityChip: {
    borderRadius: "4px",
    display: "flex",
    gap: "4px",
    justifyContent: "start",
    alignItems: "center",
    fontSize: "16px",
    cursor: "pointer",
  },

  severityChips: {
    display: "flex",
    gap: "12px",
    alignItems: "center",

  },
  notification: {
    margin: theme.spacing(0.5, 1),
  },
}));

