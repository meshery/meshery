import React, { useEffect, useState } from "react";
import { styled } from "@mui/material/styles";
import Typography from "@mui/material/Typography";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import WarningIcon from "@mui/icons-material/Warning";

const SessionExpired = styled(DialogContentText)(() => ({
  minWidth : 400,
  overflowWrap : "anywhere",
  textAlign : "center",
  padding : 5,
  margin : 2,
  display : "flex",
  flexDirection : "column",
  height : "7rem",
  justifyContent : "space-evenly",
}));
const IconContainer = styled("div")(() => ({
  width : "24px",
  height : "24px",
  marginRight : "1px",
}));

function AlertUnauthenticatedSession() {
  const [open, setOpen] = useState(false);
  const [countDown, setCountDown] = useState(3);

  useEffect(() => {
    const timer = setTimeout(() => {
      if (countDown === 1) {
        handleClose();
        window.location = "/user/login";
      }
      setCountDown((countDown) => countDown - 1);
    }, 1000);
    return () => clearTimeout(timer);
  });

  useEffect(() => {
    setOpen(true);
  }, []);

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle
        id="alert-dialog-title"
        sx={{
          display : "flex",
          justifyContent : "center",
          alignItems : "center",
          textAlign : "center",
          minWidth : 400,
          padding : "10px",
          color : "#ebf1f5",
          backgroundColor : "#F0A303",
        }}
      >
        <IconContainer>
          <WarningIcon color="#F0D053" />
        </IconContainer>
        Session Expired
      </DialogTitle>
      <DialogContent>
        <SessionExpired id="alert-dialog-description">
          <Typography variant="body1">Your session has expired</Typography>
          <Typography>
            You will be redirected to login in {countDown}
          </Typography>
        </SessionExpired>
      </DialogContent>
    </Dialog>
  );
}

export default AlertUnauthenticatedSession;
