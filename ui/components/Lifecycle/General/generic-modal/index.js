import React from 'react';
import CloseIcon from '@material-ui/icons/Close';
import {
  Typography,
  IconButton,
  Button,
  Grid,
  DialogTitle,
  DialogContent,
  DialogActions,
  Dialog,
  Backdrop,
} from '@material-ui/core';
import useStyles from './styles';
import HelpOutlineIcon from '../../../../assets/icons/HelpOutlineIcon';
import theme from '../../../../themes/app';
import { CustomTextTooltip } from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';

const GenericModal = ({
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
}) => {
  const classes = useStyles();

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
    >
      <DialogTitle textAlign="center" id="form-dialog-title" className={classes.dialogTitle}>
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
            <CustomTextTooltip title={helpText}>
              <IconButton>
                <HelpOutlineIcon
                  height={'2rem'}
                  wdith={'2rem'}
                  fill={theme.palette.secondary.primaryModalText}
                />
              </IconButton>
            </CustomTextTooltip>
          )}
          <IconButton
            aria-label="close"
            onClick={handleClose}
            component="button"
            style={{
              color: theme.palette.secondary.primaryModalText,
            }}
          >
            <CloseIcon className={classes.closing} />
          </IconButton>
        </div>
      </DialogTitle>
      <DialogContent style={{ padding: '1.5rem' }}>{body}</DialogContent>
      {!hideFooter && (
        <DialogActions
          style={{
            justifyContent: 'space-evenly',
            marginBottom: '0.5rem',
          }}
        >
          <Grid className={classes.modalActions}>
            <Button variant="outlined" onClick={handleClose} className={classes.copyButton}>
              Cancel
            </Button>
            <Button
              title={buttonTitle ? buttonTitle : 'Submit'}
              variant="contained"
              color="primary"
              className={classes.submitButton}
              disabled={disabled}
              onClick={action}
            >
              {actionBtnIcon && (
                <div
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    marginRight: '0.27rem',
                    color: theme.palette.primaryModalText,
                  }}
                >
                  {actionBtnIcon}
                </div>
              )}
              <span className={classes.btnText}>{buttonTitle ? buttonTitle : 'Submit'}</span>
            </Button>
          </Grid>
        </DialogActions>
      )}
    </Dialog>
  );
};

export default GenericModal;
