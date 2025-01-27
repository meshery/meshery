// @ts-check
import React from 'react';
import { Fade } from '@mui/material';
import { Modal, Backdrop } from '@layer5/sistent';

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
export default function GenericModal({ open, Content, handleClose, container }) {
  return (
    <Modal
      style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}
      open={open}
      onClose={handleClose}
      closeAfterTransition
      BackdropComponent={Backdrop}
      BackdropProps={{ timeout: 200 }}
      container={container}
    >
      <Fade in={open} style={{ maxHeight: '90vh', overflow: 'auto' }}>
        {Content}
      </Fade>
    </Modal>
  );
}
