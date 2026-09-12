import React from 'react';
import { Modal, Backdrop, Box, styled, Fade, IconButton, CloseIcon } from '@sistent/sistent';

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
      closeModal={handleClose}
      closeAfterTransition
      slots={{ backdrop: Backdrop }}
      slotProps={{ backdrop: { timeout: 200 } }}
      container={container}
      maxWidth="lg"
    >
      <StyledFade in={open}>
        <Box sx={{ outline: 'none', width: '100%', position: 'relative' }}>
          <IconButton
            onClick={(e) => {
              e.stopPropagation();
              handleClose();
            }}
            sx={{
              position: 'absolute',
              top: 12,
              right: 12,
              color: (theme) => theme.palette.text.secondary,
              zIndex: 1300,
            }}
          >
            <CloseIcon />
          </IconButton>
          {Content}
        </Box>
      </StyledFade>
    </StyledModal>
  );
}
