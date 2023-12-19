import React from 'react';
import { Backdrop, Box, Fade, Dialog, Typography, IconButton } from '@mui/material';
import {
  CreateTokenModalBody,
  CreateTokenModalHeader,
  ButtonContainer,
  GenerateButton,
  CancelButton,
  ModalCloseIcon,
  CustomTooltip,
  ModalWrapper,
} from './styles';
import HelpOutlineIcon from '../../../assets/icons/HelpOutlineIcon';

export default function GenericModal({
  open,
  handleClose,
  title,
  body,
  selector,
  action,
  buttonTitle,
  leftHeaderIcon,
  actionBtnIcon,
  hideFooter = false,
  disabled = false,
  helpText,
  maxWidth = 'xs',
}) {
  return (
    <Dialog
      fullWidth={true}
      maxWidth={maxWidth}
      aria-labelledby="transition-modal-title"
      aria-describedby="transition-modal-description"
      open={open}
      onClose={handleClose}
      closeAfterTransition
      slots={{ backdrop: Backdrop }}
      slotProps={{
        backdrop: {
          timeout: 500,
        },
      }}
      id="searchClick"
    >
      <Fade in={open}>
        <Box>
          <ModalWrapper>
            <CreateTokenModalHeader>
              {leftHeaderIcon && (
                <div style={{ display: 'flex', alignItems: 'center' }}>{leftHeaderIcon}</div>
              )}
              <div style={{ display: 'flex', alignItems: 'center' }}>
                <Typography variant="h6" component="h6">
                  {title}
                </Typography>
                {selector ? selector : null}
              </div>
              <div style={{ display: 'flex', alignItems: 'center' }}>
                {helpText && (
                  <CustomTooltip title={helpText}>
                    <IconButton>
                      <HelpOutlineIcon height={'2rem'} wdith={'2rem'} fill="#ffffff" />
                    </IconButton>
                  </CustomTooltip>
                )}
                <ModalCloseIcon onClick={handleClose} rotate={90} />
              </div>
            </CreateTokenModalHeader>

            <CreateTokenModalBody>{body}</CreateTokenModalBody>

            {!hideFooter && (
              <ButtonContainer>
                <CancelButton variant="contained" onClick={handleClose}>
                  Cancel
                </CancelButton>
                <GenerateButton disabled={disabled} variant="contained" onClick={action}>
                  {actionBtnIcon && (
                    <div
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        marginRight: '0.27rem',
                        color: '#ffffff',
                      }}
                    >
                      {actionBtnIcon}
                    </div>
                  )}
                  {buttonTitle ? buttonTitle : 'Continue'}
                </GenerateButton>
              </ButtonContainer>
            )}
          </ModalWrapper>
        </Box>
      </Fade>
    </Dialog>
  );
}
