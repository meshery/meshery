import { Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Slide } from "@mui/material";
import { v4 } from "uuid";
import React from "react";
import { PromptButton } from "../Button";

const Transition = React.forwardRef(function Transition(props, ref) {
  return <Slide direction="up" ref={ref} {...props} />;
});

/**
 * @typedef {{displayName: string, handler: () => void}} Option
 */

/**
 * @param {{ title: string, content: string, options: Array.<Option>, open: boolean, handleClose: () => void }} Props
 */
// eslint-disable-next-line react/prop-types
export const Prompt = ({ title, content, options, open, handleClose }) => {
  return (
    <Dialog
      open={open}
      TransitionComponent={Transition}
      keepMounted
      onClose={handleClose}
      aria-describedby="alert-dialog-slide-description"
    >
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <DialogContentText>{content}</DialogContentText>
      </DialogContent>
      <DialogActions>
        {/* eslint-disable-next-line react/prop-types */}
        {options.map((opt) => (
          <PromptButton key={v4()} onClick={opt.handler}>
            {opt.displayName}
          </PromptButton>
        ))}
      </DialogActions>
    </Dialog>
  );
};
