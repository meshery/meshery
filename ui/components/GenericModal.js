import React from "react";
import { Modal, Backdrop, Fade } from "@mui/material";

/**
 *
 * @param {{
 *  open?: boolean,
 *  Content?: JSX.Element,
 *  handleClose?: (event: {}, reason: "backdropClick" | "escapeKeyDown") => void,
 *  container?: React.ReactInstance | (() => React.ReactInstance)
 * }} props
 * @returns
 */
export default function GenericModal({ open, handleClose, Content, container }) {
  return (
    <Modal
      sx={{
        display : "flex",
        alignItems : "center",
        justifyContent : "center",
      }}
      open={open}
      onClose={handleClose}
      closeAfterTransition
      slots={Backdrop}
      slotProps={{ timeout : 200 }}
      container={container}
    >
      <Fade in={open} sx={{ maxHeight : "90vh", overflow : "auto" }}>
        {Content}
      </Fade>
    </Modal>
  );
}
