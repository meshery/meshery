import makeStyles from "@mui/styles/makeStyles";
export const getStyles = () => ({
  title: {
    textAlign: "center",
    minWidth: 400,
    padding: "10px",
    color: "#fff",
    backgroundColor: "#607d8b",
  },
  subtitle: {
    minWidth: 400,
    overflowWrap: "anywhere",
    textAlign: "center",
    padding: "5px",
  },
  actions: {
    display: "flex",
    justifyContent: "center",
  },
  button0: {
    margin: "8px 0px",
    width: "100%",
  },
  button1: {
    margin: "8px 0px",
    width: "100%",
    backgroundColor: "#e0e0e0",
    color: "rgba(0, 0, 0, 0.87)",
    "&:hover": {
      backgroundColor: "#d5d5d5",
      boxShadow:
        "0px 2px 4px -1px rgb(0 0 0 / 20%), 0px 4px 5px 0px rgb(0 0 0 / 14%), 0px 1px 10px 0px rgb(0 0 0 / 12%)",
    },
  },
});
export const useStyles = makeStyles(getStyles);
