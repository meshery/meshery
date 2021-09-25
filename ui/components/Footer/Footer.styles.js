import makeStyles from "@mui/styles/makeStyles";

export const getStyles = (theme) => ({
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
  ComLink: {
    textDecoration: "none",
  },
});
export const useStyles = makeStyles(getStyles);
