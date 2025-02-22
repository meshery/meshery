import React from 'react';
import { Modal, Backdrop, Box, styled, Fade } from '@layer5/sistent';

const StyledModal = styled(Modal)({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
});

const StyledFade = styled(Fade)({
  maxHeight: '90vh',
  overflow: 'auto',
});

export default function GenericModal({ open, Content, handleClose, container }) {
  return (
    <StyledModal
      open={open}
      onClose={handleClose}
      closeAfterTransition
      slots={{ backdrop: Backdrop }}
      slotProps={{ backdrop: { timeout: 200 } }}
      container={container}
      maxWidth="lg"
    >
      <StyledFade in={open}>
        <Box sx={{ outline: 'none', width: '100%' }}>{Content}</Box>
      </StyledFade>
    </StyledModal>
  );
}
