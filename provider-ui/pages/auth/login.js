import { Typography } from "@material-ui/core";
import React, { useEffect, useState } from "react";
import { withStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle'
import InfoIcon from '@material-ui/icons/Info';

const styles = (() => ({
  title : {
    textAlign : 'center',
    minWidth : 400,
    padding : '10px',
    color : '#2198f3',
    backgroundColor : '#e6f2fb',
  },
  subtitle : {
    minWidth : 400,
    overflowWrap : 'anywhere',
    textAlign : 'center',
    padding : '5px',
    margin: "2px",
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    height: "7rem",
    justifyContent: "space-evenly"
  },
  icon: {
    padding: "inherit",
    position: "absolute",
    left: "25px",
    top: "20px" 
  }
}));


const AlertUnauthenticatedSession = ({ classes }) => {

  const [open, setOpen] = useState(false);
  const [countDown, setCountDown] = useState(3);

  useEffect(() => {
    const timer = setTimeout(() => {
      if (countDown === 1) {
        handleClose();
        window.location = "/user/login";
      }
      setCountDown(countDown => countDown - 1)
    }, 1000)
    return () => clearTimeout(timer);
  })
 
  useEffect(() => {
    setOpen(true);
  },[]);

  const handleClose = () => {
    setOpen(false);
  };

  return (
      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
        className={classes.dialogBox}
      >
        <DialogTitle id="alert-dialog-title" className={classes.title}>
          <span className={classes.icon}>
            <InfoIcon color="#2196f3"/>
          </span>
              Session Expired 
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description" className={classes.subtitle}>
            <Typography variant="body1">
              Your session has expired
            </Typography>
            <Typography>
              You will be redirected to login in {countDown}
            </Typography>
          </DialogContentText>
        </DialogContent>
      </Dialog>
  );
}

export default withStyles(styles)(AlertUnauthenticatedSession);
