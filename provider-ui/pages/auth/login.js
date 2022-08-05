import { Typography } from "@material-ui/core";
import React, { useEffect, useState } from "react";
import { withStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle'
import WarningIcon from '@material-ui/icons/Warning';

const styles = (() => ({
  title : {
    textAlign : 'center',
    minWidth : 400,
    padding : '10px',
    color : '#ebf1f5',
    backgroundColor : '#F0A303',
  },
  subtitle : {
    minWidth : 400,
    overflowWrap : 'anywhere',
    textAlign : 'center',
    padding : '5px',
    margin: "2px",
    display: "flex",
    flexDirection: "column",
    height: "7rem",
    justifyContent: "space-evenly"
  },
  icon: {
    padding: "inherit",
    position: "relative",
    right: "120px",
    top: "5px" 
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
            <WarningIcon color="#F0D053"/>
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
