import React, { useState } from "react";
import {
  Typography,
  Button,
  Dialog,
  DialogActions,
  DialogContentText,
  DialogContent,
  DialogTitle,
} from "@mui/material";
import { useStyles } from "./Prompt.styles.js";

const Prompt = () => {
  const classes = useStyles();
  const [showModal, setShowModal] = useState(false);
  const [title, setTitle] = useState("");
  const [subtitle, setSubtitle] = useState("");
  const [options, setOptions] = useState([]);
  let promiseInfo = {};

  const show = async (passed) => {
    return new Promise((resolve, reject) => {
      promiseInfo = { resolve, reject };
      setTitle(passed.title);
      setSubtitle(passed.subtitle);
      setOptions(passed.options);
      showModal(true);
    });
  };

  const hide = () => {
    setShowModal(false);
  };

  const { resolve } = promiseInfo;
  return (
    <div className={classes.root}>
      <Dialog
        open={show}
        onClose={hide}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
        className={classes.dialogBox}
      >
        {title !== "" && (
          <DialogTitle id="alert-dialog-title" className={classes.title}>
            <b>{title}</b>
          </DialogTitle>
        )}
        {subtitle !== "" && (
          <DialogContent>
            <DialogContentText id="alert-dialog-description" className={classes.subtitle}>
              <Typography variant="body1">{subtitle}</Typography>
            </DialogContentText>
          </DialogContent>
        )}
        <DialogActions className={classes.actions}>
          <Button
            onClick={() => {
              hide();
              resolve(options[1]);
            }}
            key={options[1]}
            className={classes.button1}
          >
            {options[1]}
          </Button>
          <Button
            onClick={() => {
              hide();
              resolve(options[0]);
            }}
            key={options[0]}
            className={classes.button0}
            type="submit"
            variant="contained"
            color="primary"
          >
            {options[0]}
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default Prompt;
