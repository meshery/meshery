// @ts-check
import React from "react";
import { Modal, Backdrop, Fade } from "@material-ui/core";

export default function GenericModal({ open, Content, handleClose }) {
  return (
    <Modal
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
      open={open}
      onClose={handleClose}
      closeAfterTransition
      BackdropComponent={Backdrop}
      BackdropProps={{
        timeout: 200,
      }}
    >
      <Fade in={open} style={{ maxHeight: "90vh", overflow: "auto" }} >{Content}</Fade>
    </Modal>
  );
}
