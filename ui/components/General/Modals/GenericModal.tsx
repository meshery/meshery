import React from 'react';
import { Modal, Backdrop, Box, styled, Fade } from '@sistent/sistent';

const StyledModal = styled(Modal)({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
});

const StyledFade = styled(Fade)({
  maxHeight: '90vh',
  overflow: 'auto',
});

type GenericModalProps = {
  open: boolean;
  Content: React.ReactNode;
  handleClose: () => void;
  container?: HTMLElement | null;
};

export default function GenericModal({ open, Content, handleClose, container }: GenericModalProps) {
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
