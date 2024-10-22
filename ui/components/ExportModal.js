import React from 'react';
import { ListItemText, Modal, useTheme, Box, DownloadIcon } from '@layer5/sistent';
import KubernetesIcon from '@/assets/icons/technology/kubernetes';
import { useModal } from '@layer5/sistent';
import PatternIcon from '@/assets/icons/Pattern';
import { ListItem } from '@layer5/sistent';
import { ListItemIcon } from '@layer5/sistent';
import { UsesSistent } from './SistentWrapper';
import { ModalBody } from '@layer5/sistent';
import { Colors } from '@/themes/app';
import { InfoTooltip } from '@layer5/sistent';
import { IconButton } from '@layer5/sistent';
import { OCIImageIcon } from '@/assets/icons/OciImage';
import HelmIcon from '@/assets/icons/technology/HelmIcon';

const ExportOption = ({
  title,
  Icon,
  onClick,
  content,
  disabled = false,
  description = 'Download the design in the selected format',
}) => {
  const theme = useTheme();
  return (
    <ListItem
      sx={{
        display: 'flex',
        flexDirection: 'column',
        gap: '0.5rem',
        borderRadius: '0.25rem',
        opacity: disabled ? 0.5 : 1,
        marginBottom: '1rem',
        border: `1px solid ${theme.palette.border.normal}`,
        boxShadow: theme.shadows[1],
        cursor: disabled ? 'not-allowed' : 'default',
        '&:hover': {
          borderColor: theme.palette.border.brand,
        },
      }}
    >
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          width: '100%',
          padding: theme.spacing(1),
          // borderBottom: `1px solid ${theme.palette.border.normal}`,
          justifyContent: 'space-between',
        }}
      >
        <ListItemIcon style={{ minWidth: '2rem', marginRight: '1rem' }}>{Icon}</ListItemIcon>
        <ListItemText primary={title} />
        <InfoTooltip placement="top" title={description} content={description} />
        <IconButton
          disabled={disabled}
          onClick={onClick}
          sx={{
            marginLeft: '0.5rem',
          }}
        >
          <DownloadIcon fill={theme.palette.icon.default} />
        </IconButton>
      </Box>
      {content}
    </ListItem>
  );
};
const ExportModal = (props) => {
  const {
    downloadModal,
    handleDownloadDialogClose,
    handleDesignDownload,
    extensionExportOptions = [],
  } = props;

  const ExportOptions = [
    {
      title: 'Design File',
      icon: <PatternIcon width={'30'} height="30" fill={Colors.caribbeanGreen} />,
      onClick: (e) => handleDesignDownload(e, downloadModal.content),
      description:
        'Download the design as a Meshery Design file , that can be imported later without any loss.',
    },
    {
      title: 'Kubernetes Manifest',
      icon: <KubernetesIcon width={'30'} height="30" />,
      onClick: (e) =>
        handleDesignDownload(e, downloadModal.content, null, 'export=Kubernetes Manifest'),
      description:
        'Download the design as a Kubernetes Manifest file, that can be applied to a Kubernetes cluster. This is a lossy export as some meshery specific  metadata gets stripped off.',
    },
    {
      title: 'OCI Image',
      icon: <OCIImageIcon width={'30'} height="30" />,
      onClick: (e) => handleDesignDownload(e, downloadModal.content, null, 'export=OCI Image'),
      description:
        'Download the design as an OCI Image, that can be pushed to a container registry. This is a lossless export.',
    },
    {
      title: 'Helm Chart (Coming Soon)',
      icon: <HelmIcon width={'30'} height="30" />,
      onClick: () => {},
      disabled: true,
      description:
        'Download the design as a Helm Chart. This is a lossy export as some meshery specific metadata gets stripped off.',
    },
    ...extensionExportOptions,
  ];

  const exportModal = useModal({
    headerIcon: <PatternIcon fill={'#fff'} height={'2rem'} width="2rem" />,
  });
  return (
    <UsesSistent>
      <Modal
        {...exportModal}
        open={downloadModal.open}
        onClose={handleDownloadDialogClose}
        closeModal={handleDownloadDialogClose}
        title="Export Design"
      >
        <ModalBody>
          {ExportOptions.map((option) => (
            <ExportOption
              key={option.title}
              title={option.title}
              Icon={option.icon}
              content={option.content}
              disabled={option.disabled}
              description={option.description}
              onClick={(e) => option.onClick(e)}
            />
          ))}
        </ModalBody>
      </Modal>
    </UsesSistent>
  );
};

export default ExportModal;
// <div>
//        <Typography component="p" style={{ paddingBottom: '1.5rem' }}>
//          Design File
//        </Typography>
//        <div style={exportBtnStyles}>
//          <Button
//            data-testid="export-current"
//            onClick={(e) => handleDesignDownload(e, downloadModal.content)}
//          >
//            <div style={exportWrpStyles}>
//              <ModifiedApplicationFileIcon width={75} height={82} />
//              <div style={{ display: 'flex', padding: '0.4rem' }}>
//                <Typography> EXPORT </Typography>
//                <GetAppIcon />
//              </div>
//            </div>
//          </Button>
//        </div>
//      </div>
//      <div>
//        <Typography component={'h4'} style={{ paddingBottom: '1.5rem', maxWidth: '9rem' }}>
//          Kubernetes Manifest
//        </Typography>
//        <div style={exportBtnStyles}>
//          <Button
//            data-testid="export-original"
//            onClick={(e) =>
//              handleDesignDownload(e, downloadModal.content, null, 'export=Kubernetes Manifest')
//            }
//          >
//            <div style={exportWrpStyles}>
//              <KubernetesIcon width={75} height={75} />
//              <div style={{ display: 'flex', padding: '0.8rem' }}>
//                <Typography> EXPORT </Typography>
//                <GetAppIcon />
//              </div>
//            </div>
//          </Button>
//        </div>
//      </div>
//      <div>
//        <Typography style={{ paddingBottom: '1.5rem' }}>OCI</Typography>

//        <div style={exportBtnStyles}>
//          <Button
//            data-testid="export-oci"
//            onClick={(e) => handleDesignDownload(e, downloadModal.content, null, 'oci=true')}
//          >
//            <div style={exportWrpStyles}>
//              <ModifiedApplicationFileIcon width={75} height={82} />
//              <div style={{ display: 'flex', padding: '0.4rem' }}>
//                <Typography> EXPORT </Typography>
//                <GetAppIcon />
//              </div>
//            </div>
//          </Button>
//        </div>
//      </div>
