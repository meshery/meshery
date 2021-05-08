import React from "react";
import { makeStyles, Typography } from "@material-ui/core/";
import FiberManualRecordRoundedIcon from "@material-ui/icons/FiberManualRecordRounded";

const useStyles = makeStyles({
  infoContainer: {
    position: "relative",
    top: "2rem",
    width: "12rem",
    height: "8rem",
    padding: "1rem 1rem",
    margin: "0rem auto",
    boxShadow: "0px 1px 6px 1px rgba(0,0,0,0.75)",
  },
  infoContext: {
    display: "inline",
    padding: "0.25rem 0.5rem",
    fontSize: ".75rem",
    fontWeight: "300",
    background: "lightgray",
  },
  infoText: {
    marginTop: "1rem",
  },
  planeIcon: {
    marginBottom: "-0.4rem",
    color: "#00B39F",
  },
});

const InfoContainer = ({ controlPlane, dataPlane }) => {
  const classes = useStyles();
  return (
    <div className={classes.infoContainer}>
      <Typography className={classes.infoContext}>Context</Typography>
      <Typography className={classes.infoText}>
        Control Plane: {controlPlane} <FiberManualRecordRoundedIcon className={classes.planeIcon} />
      </Typography>
      <Typography>
        Data Plane: {dataPlane}
        <FiberManualRecordRoundedIcon className={classes.planeIcon} />
      </Typography>
    </div>
  );
};

export default InfoContainer;
