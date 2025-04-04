import React from 'react';
import { ListItemText, Modal, useTheme, Box, DownloadIcon } from '@layer5/sistent';
import KubernetesIcon from '@/assets/icons/technology/kubernetes';
import { useModal } from '@layer5/sistent';
import PatternIcon from '@/assets/icons/Pattern';
import { ListItem } from '@layer5/sistent';
import { ListItemIcon } from '@layer5/sistent';
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
      title: 'Meshery Design (yaml)',
      icon: <PatternIcon width={'30'} height="30" fill={Colors.caribbeanGreen} />,
      onClick: (e) => handleDesignDownload(e, downloadModal.content),
      description:
        'Export your design as a complete, self-contained Meshery Design file (YAML). This file includes embedded images and all configuration details. It&apos;s the perfect format for creating backups, sharing with colleagues using Meshery, or transferring designs between Meshery environments without losing any information (lossless transfer).',
    },
    {
      title: 'Meshery Design (OCI image)',
      icon: <OCIImageIcon width={'30'} height="30" />,
      onClick: (e) => handleDesignDownload(e, downloadModal.content, null, 'oci=true'),
      description:
        'Download your design as an OCI compatible container image, which can be pushed to and pulled from container registries like Docker Hub, AWS ECR, and so on.',
    },
    {
      title: 'Kubernetes Manifest (yaml)',
      icon: <KubernetesIcon width={'30'} height="30" />,
      onClick: (e) =>
        handleDesignDownload(e, downloadModal.content, null, 'export=Kubernetes Manifest'),
      description: (
        <div>
          <p>
            Download your design as a standard Kubernetes Manifest file. This file contains the
            Kubernetes resource definitions from your design and can be directly applied to a
            cluster using tools `kubectl`.
          </p>
          <p>
            <strong>Lossy Export:</strong> This process strips out Meshery-specific information
            (e.g., visual arrangement, comments, and so on). The resulting manifest only includes
            the core Kubernetes resource definitions, not the extra context that might be present in
            your Meshery design.
          </p>
        </div>
      ),
    },
    {
      title: 'Helm Chart (tar.gz) (Coming Soon)',
      icon: <HelmIcon width={'30'} height="30" />,

      onClick: (e) => handleDesignDownload(e, downloadModal.content, null, 'export=Helm Chart'),
      disabled: false,
      description: (
        <div>
          <p>Download your design as a Helm Chart.</p>
          <p>
            <strong>Lossy Export:</strong> This process strips out Meshery-specific information
            (e.g., visual arrangement, comments, and so on). The resulting manifest only includes
            the core Kubernetes resource definitions, not the extra context that might be present in
            your Meshery design.
          </p>
        </div>
      ),
    },
    ...extensionExportOptions,
  ];

  const exportModal = useModal({
    headerIcon: <PatternIcon fill={'#fff'} height={'2rem'} width="2rem" />,
  });
  return (
    <>
      <Modal
        {...exportModal}
        open={downloadModal.open}
        onClose={handleDownloadDialogClose}
        closeModal={handleDownloadDialogClose}
        title="Export Design as..."
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
    </>
  );
};

export default ExportModal;
