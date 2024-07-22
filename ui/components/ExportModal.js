import { Typography, Button, Box } from '@material-ui/core';
import React from 'react';
import PatternIcon from '@/assets/icons/Pattern';
import { ModalFooter } from '@layer5/sistent';
import { GetApp as GetAppIcon } from '@material-ui/icons';
import OriginalApplicationFileIcon from '@/assets/icons/OriginalApplicationIcon';
import ModifiedApplicationFileIcon from '@/assets/icons/ModifiedApplicationIcon';
import { withStyles } from '@material-ui/core/styles';
import { Modal, ModalBody } from '@layer5/sistent';
import theme, { Colors } from '@/themes/app';
import { iconMedium } from 'css/icons.styles';
import { UsesSistent } from './SistentWrapper';

const styles = (theme) => ({
  dialogTitle: {
    backgroundColor: theme.palette.secondary.mainBackground,
    justifyContent: 'space-between',
    flexDirection: 'row',
    alignItems: 'center',
    padding: '12px 20px',
    gap: '146px',
    color: '#FFFFFF',
    textAlign: 'center',
    textOverflow: 'ellipsis',
    '& h2': {
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'center',
    },
  },
  text: {
    fontFamily: 'Qanelas Soft, sans-serif',
    '&.MuiTypography-root': {
      fontFamily: 'Qanelas Soft, sans-serif',
    },
  },
  textHeader: {
    fontFamily: 'Qanelas Soft, sans-serif',
    textAlign: 'center',
  },
  closing: {
    transform: 'rotate(-90deg)',
    '&:hover': {
      transform: 'rotate(90deg)',
      transition: 'all .3s ease-in',
      cursor: 'pointer',
    },
  },
  exportBtns: {
    '& p': {
      margin: '0 0 0.75rem',
    },

    '& button': {
      height: '8rem',
      border: 0,
      boxShadow: '0rem 0rem 0.375rem rgba(0, 0, 0, 0.25)',

      '& span': {
        flexWrap: 'wrap',

        '& .MuiButton-startIcon': {
          flex: '0 0 100%',
          justifyContent: 'center',
        },

        '& .MuiButton-endIcon': {
          margin: 0,
        },
      },
    },
  },
  infoIconButton: {
    color: theme.palette.secondary.focused,
  },
  infoIcon: {
    color: theme.palette.secondary.focused,
    fill: 'currentColor',
    '&:hover': {
      color: Colors.keppelGreen,
    },
  },
  toolTip: {
    textDecoration: 'underline',
    color: theme.palette.secondary.link2,
  },
  dialogAction: {
    padding: '0.5rem 1rem',
    '&.MuiDialogActions-root': {
      padding: '0.5rem 1rem',
    },
  },
});

const ExportModal = (props) => {
  const {
    downloadModal,
    handleDownloadDialogClose,
    handleDesignDownload,
    classes,
    ExtensibleButton,
  } = props;

  const [closed, setClosed] = React.useState(false);

  const handleClose = () => {
    setClosed(true);
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
    <UsesSistent>
      <Modal
        open={downloadModal.open}
        closeModal={handleClose}
        title="Export Design"
        headerIcon={<PatternIcon style={iconMedium} fill={theme.palette.secondary.whiteIcon} />}
        aria-labelledby="download-design-dialog"
        aria-describedby="download-design-dialog-description"
        maxWidth="sm"
      >
        <ModalBody>
          <Box
            sx={{
              display: 'flex',
              gap: '2.5rem',
              width: 'auto',
              textAlign: 'center',
              alignItems: 'center',
              justifyContent: 'center',
              flexWrap: 'wrap',
            }}
          >
            {downloadModal?.content?.type?.String && (
              <div>
                <Typography
                  component={'h4'}
                  style={{ paddingBottom: '1.5rem', maxWidth: '9rem' }}
                  className={classes.text}
                >
                  Original ({downloadModal?.content?.type?.String})
                </Typography>
                <div style={exportBtnStyles}>
                  <Button
                    onClick={(e) =>
                      handleDesignDownload(
                        e,
                        downloadModal.content,
                        downloadModal?.content.type?.String,
                      )
                    }
                  >
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
            <div>
              <Typography
                component="p"
                style={{ paddingBottom: '1.5rem' }}
                className={classes.text}
              >
                Current
              </Typography>
              <div style={exportBtnStyles}>
                <Button onClick={(e) => handleDesignDownload(e, downloadModal.content)}>
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
            <div>
              <Typography style={{ paddingBottom: '1.5rem' }} className={classes.text}>
                OCI
              </Typography>

              <div style={exportBtnStyles}>
                <Button
                  onClick={(e) => handleDesignDownload(e, downloadModal.content, null, 'oci=true')}
                >
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
            {ExtensibleButton && <ExtensibleButton {...props} closed={closed} />}
          </Box>
        </ModalBody>
        <ModalFooter
          variant="filled"
          helpText="MeshMap Designer offers multiple export options, allowing you to choose the format that suits your needs. [Learn more](https://docs.layer5.io/meshmap/designer/export-designs/)"
        ></ModalFooter>
      </Modal>
    </UsesSistent>
  );
};

export default withStyles(styles)(ExportModal);
