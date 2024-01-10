import {
  Dialog,
  DialogTitle,
  Typography,
  IconButton,
  DialogContent,
  Button,
  Box,
  // CircularProgress,
  // Tooltip,
} from '@material-ui/core';
import React from 'react';
import PatternIcon from '@/assets/icons/Pattern';
import { CloseIcon } from '@layer5/sistent-svg';
import { GetApp as GetAppIcon } from '@material-ui/icons';
import OriginalApplicationFileIcon from '@/assets/icons/OriginalApplicationIcon';
import ModifiedApplicationFileIcon from '@/assets/icons/ModifiedApplicationIcon';
// import { useNotification } from '@/utils/hooks/useNotification';
// import HelpOutlineIcon from '@/assets/icons/HelpOutlineIcon';
// import { iconSmall } from '@/utils/icon';
// import theme from '@/themes/app';

export const ExportModal = ({ downloadModal, handleDownloadDialogClose, handleDesignDownload }) => {
  // const { notify } = useNotification();

  const handleClose = () => {
    handleDownloadDialogClose();
  };

  const exportBtnStyles = {
    boxShadow: '0px 0px 6px 2px rgba(0, 0, 0, 0.25)',
    borderRadius: '20px',
  };

  const exportWrpStyles = {
    width: '8rem',
    padding: '1rem',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
  };
  return (
    <Dialog
      open={downloadModal.open}
      onClose={handleClose}
      aria-labelledby="download-design-dialog"
      aria-describedby="download-design-dialog-description"
      maxWidth="xl"
    >
      <DialogTitle
        textAlign="center"
        id="download-design-dialog-title"
        // className={dialogClasses.dialogTitle}
      >
        <PatternIcon width={30} height={30} style={{ filter: 'none', opacity: 1 }} fill="#FFF" />
        <Typography
          // className={dialogClasses.textHeader}
          variant="h6"
        >
          Export Design
        </Typography>
        <IconButton
          aria-label="close"
          onClick={handleClose}
          component="button"
          style={{
            color: '#FFFFFF',
          }}
        >
          <CloseIcon
            // className={dialogClasses.closing}
            fill={'#FFF'}
          />
        </IconButton>
      </DialogTitle>
      <DialogContent
        style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          gap: '3.5rem',
          maxWidth: '688px',
          padding: '5rem',
        }}
      >
        <Box
          sx={{
            display: 'flex',
            gap: '2.5rem',
            width: '100%',
            textAlign: 'center',
            alignItems: 'center',
          }}
        >
          {downloadModal?.content?.type?.String && (
            <div
            // className={classes.exportBtns}
            >
              <Typography
                component={'h4'}
                style={{ paddingBottom: '1.5rem' }}
                // className={dialogClasses.text}
              >
                Original ({downloadModal?.content?.type?.String})
              </Typography>
              <div style={exportBtnStyles}>
                <Button onClick={(e) => handleDesignDownload(e, downloadModal.content, false)}>
                  <div style={exportWrpStyles}>
                    <OriginalApplicationFileIcon width={75} height={75} />
                    <div style={{ display: 'flex', padding: '0.8rem' }}>
                      <Typography> EXPORT </Typography>
                      <GetAppIcon />
                    </div>
                  </div>
                </Button>
              </div>
            </div>
          )}
          <div
          // className={classes.exportBtns}
          >
            <Typography
              component="p"
              style={{ paddingBottom: '1.5rem' }}
              // className={dialogClasses.text}
            >
              Current
            </Typography>
            <div style={exportBtnStyles}>
              <Button onClick={(e) => handleDesignDownload(e, downloadModal.content, false)}>
                <div style={exportWrpStyles}>
                  <ModifiedApplicationFileIcon width={75} height={82} />
                  <div style={{ display: 'flex', padding: '0.4rem' }}>
                    <Typography> EXPORT </Typography>
                    <GetAppIcon />
                  </div>
                </div>
              </Button>
            </div>
          </div>
          <div
          // className={classes.exportBtns}
          >
            <Typography
              style={{ paddingBottom: '1.5rem' }}
              // className={dialogClasses.text}
            >
              OCI
            </Typography>

            <div style={exportBtnStyles}>
              <Button onClick={(e) => handleDesignDownload(e, downloadModal.content, true)}>
                <div style={exportWrpStyles}>
                  <ModifiedApplicationFileIcon width={75} height={82} />
                  <div style={{ display: 'flex', padding: '0.4rem' }}>
                    <Typography> EXPORT </Typography>
                    <GetAppIcon />
                  </div>
                </div>
              </Button>
            </div>
          </div>
        </Box>
      </DialogContent>
    </Dialog>
  );
};
