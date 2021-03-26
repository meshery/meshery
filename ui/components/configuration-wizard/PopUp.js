import React from "react";
import Link from "next/link";
import { makeStyles } from "@material-ui/core/styles";
import {
  Typography,
  Modal,
  Button,
  Checkbox,
  Zoom,
} from "@material-ui/core/";

const useStyles = makeStyles((theme) => ({
  modal: {
    position: "absolute",
    display: "flex",
    flexDirection: "column",
    justifyContent: "space-evenly",
    height: "20rem",
    padding: "1rem",
    width: "500px",
    marginLeft: "-250px",
    left: "50%",
    bottom: "35%",
    backgroundColor: "#017374",
    border: "none",
    boxShadow: theme.shadows[5],
  },
  buttonContainer: {
    alignSelf: "flex-end",
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
  },
  checkbox: {
    display: "flex",
    flexDirection: "row",
    alignItems: "center",
  },
  label: {
    //fontSize: "0.75rem",
    color: "lightgray",
    marginRight: "7rem",
  },
  startButton: {
    padding: "0.75rem 1rem",
    width: "8rem",
    background: "#455A64",
    color: "white",
    "&:hover": {
      backgroundColor: "#607D8B",
    },
  },
  link: {
    textDecoration: "none",
    color: "white",
  },
  text: {
    color: "white",
  },
}));

const PopUp = ({ open, handleClose }) => {
  const classes = useStyles();
  return (
    <Modal
      open={open}
      onClose={handleClose}
      aria-labelledby="simple-modal-title"
      aria-describedby="simple-modal-description"
    >
      <Zoom in={open}>
        <div className={classes.modal}>
          <Typography
            variant="h4"
            gutterBottom="true"
            paragraph="true"
            className={classes.text}
            id="simple-modal-title"
          >
            Configuration Wizard
          </Typography>
          <Typography
            variant="body2"
            gutterBottom="true"
            paragraph="true"
            className={classes.text}
            id="simple-modal-description"
          >
            Welcome to configuration wizard, Mesherys unique feature. The
            configuration wizard represents an assembly of existing
            functionality and settings that you can configure.
          </Typography>
          <div className={classes.buttonContainer}>
            <div className={classes.checkbox}>
              <Checkbox color="default" />
              <Typography className={classes.label}>
                Dont show this again
              </Typography>
            </div>
            <Button onClick={handleClose} className={classes.text}>
              Skip
            </Button>
            <Link href="/wizard" className={classes.link}>
              <Button onClick={handleClose} className={classes.startButton}>
                Get Started
              </Button>
            </Link>
          </div>
        </div>
      </Zoom>
    </Modal>
  );
};

export default PopUp;
