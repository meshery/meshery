import React from "react";
import Link  from "next/link";
import { makeStyles, Button, Typography } from "@material-ui/core/";
import CheckCircleIcon from "@material-ui/icons/CheckCircle";

const useStyles = makeStyles((theme) => ({
  buttonContainer: {
    textAlign: "right",
    paddingBottom: "5rem",
    marginRight: "11.5rem",
  },
  button: {
    marginRight: theme.spacing(1),
    padding: "0.5rem 2rem",
    textDecoration: "none",
    background: "white",
    color: "#647881",
    border: "1.5px solid #647881",
    "&:hover": {
      backgroundColor: "#647881",
      color: "white",
    },
  },
  linkButton: {
    textDecoration: "none",
  },
  completed: {
    textAlign: "center",
  },
  checkCircleIcon: {
    color: "#00B39F",
    padding: "1rem",
    height: "auto",
    width: "4rem",
  },
}));

const ConfigurationDone = () => {
  const classes = useStyles();
  return (
    <div className={classes.completed}>
      <CheckCircleIcon className={classes.checkCircleIcon} />
      <Typography
        variant="h4"
        gutterBottom="true"
        className={classes.instructions}
      >
        Configuration done
      </Typography>
      <Typography
        variant="subtitle1"
        paragraph="true"
        gutterBottom="true"
        className={classes.instructions}
      >
        Your configuration was successful
      </Typography>
      <div className={classes.buttonContainer}>
        <Link href="/" className={classes.linkButton}>
          <Button className={classes.button}>Done</Button>
        </Link>
      </div>
    </div>
  );
};

export default ConfigurationDone;
