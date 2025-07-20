import React, { useMemo } from 'react';
import { ListItemText, Modal, useTheme, Box, DownloadIcon } from '@sistent/sistent';
import KubernetesIcon from '@/assets/icons/technology/kubernetes';
import { useModal } from '@sistent/sistent';
import PatternIcon from '@/assets/icons/Pattern';
import { ListItem } from '@sistent/sistent';
import { ListItemIcon } from '@sistent/sistent';
import { ModalBody } from '@sistent/sistent';
import { Colors } from '@/themes/app';
import { InfoTooltip } from '@sistent/sistent';
import { IconButton } from '@sistent/sistent';
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
  
  const handleClick = (e) => {
    if (disabled) return;
    onClick(e);
  };

  const handleKeyPress = (e) => {
    if (disabled) return;
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      onClick(e);
    }
  };

  return (
    <ListItem
      role="button"
      tabIndex={disabled ? -1 : 0}
      onKeyDown={handleKeyPress}
      sx={{
        display: 'flex',
        flexDirection: 'column',
        gap: '0.5rem',
        borderRadius: '0.25rem',
        opacity: disabled ? 0.5 : 1,
        marginBottom: '1rem',
        border: `1px solid ${theme.palette.border.normal}`, // Fixed template literal
        boxShadow: theme.shadows[1],
        cursor: disabled ? 'not-allowed' : 'pointer',
        transition: theme.transitions.create(['border-color', 'box-shadow'], {
          duration: theme.transitions.duration.short,
        }),
        '&:hover': {
          borderColor: disabled ? theme.palette.border.normal : theme.palette.border.brand,
          boxShadow: disabled ? theme.shadows[1] : theme.shadows[2],
        },
        '&:focus-visible': {
          outline: `2px solid ${theme.palette.border.brand}`, // Fixed template literal
          outlineOffset: '2px',
        },
      }}
      onClick={handleClick}
    >
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          width: '100%',
          padding: theme.spacing(1),
          justifyContent: 'space-between',
        }}
      >
        <ListItemIcon 
          sx={{ 
            minWidth: '2rem', 
            marginRight: '1rem',
            color: disabled ? theme.palette.text.disabled : 'inherit',
          }}
        >
          {Icon}
        </ListItemIcon>
        <ListItemText 
          primary={title} 
          sx={{
            '& .MuiListItemText-primary': {
              color: disabled ? theme.palette.text.disabled : theme.palette.text.primary,
            }
          }}
        />
        <InfoTooltip 
          placement="top" 
          title={description} 
          content={description}
          sx={{ marginRight: '0.5rem' }}
        />
        <IconButton
          disabled={disabled}
          onClick={(e) => {
            e.stopPropagation();
            if (!disabled) onClick(e);
          }}
          aria-label={`Download ${title}`} // Fixed template literal
          sx={{
            marginLeft: '0.5rem',
            '&:hover': {
              backgroundColor: disabled ? 'transparent' : theme.palette.action.hover,
            },
          }}
        >
          <DownloadIcon fill={disabled ? theme.palette.text.disabled : theme.palette.icon.default} />
        </IconButton>
      </Box>
      {content && (
        <Box sx={{ padding: theme.spacing(0, 1, 1, 1), width: '100%' }}>
          {content}
        </Box>
      )}
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

  const theme = useTheme();

  // Memoize export options to prevent unnecessary re-renders
  const exportOptions = useMemo(() => [
    {
      title: 'Meshery Design (yaml)',
      icon: <PatternIcon width={'30'} height="30" fill={Colors.caribbeanGreen} />,
      onClick: (e) => handleDesignDownload(e, downloadModal.content),
      description:
        'Export your design as a complete, self-contained Meshery Design file (YAML). This file includes embedded images and all configuration details. It\'s the perfect format for creating backups, sharing with colleagues using Meshery, or transferring designs between Meshery environments without losing any information (lossless transfer).',
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
        <Box component="div">
          <Box component="p" sx={{ margin: '0 0 0.5rem 0' }}>
            Download your design as a standard Kubernetes Manifest file. This file contains the
            Kubernetes resource definitions from your design and can be directly applied to a
            cluster using tools like <code>kubectl</code>.
          </Box>
          <Box component="p" sx={{ margin: 0 }}>
            <Box component="strong">Lossy Export:</Box> This process strips out Meshery-specific information
            (e.g., visual arrangement, comments, and so on). The resulting manifest only includes
            the core Kubernetes resource definitions, not the extra context that might be present in
            your Meshery design.
          </Box>
        </Box>
      ),
    },
    {
      title: 'Helm Chart (tar.gz)',
      icon: <HelmIcon width={'30'} height="30" />,
      onClick: (e) => handleDesignDownload(e, downloadModal.content, null, 'export=helm-chart'),
      disabled: false,
      description: (
        <Box component="div">
          <Box component="p" sx={{ margin: '0 0 0.5rem 0' }}>
            Download your design as a Helm Chart.
          </Box>
          <Box component="p" sx={{ margin: 0 }}>
            <Box component="strong">Lossy Export:</Box> This process strips out Meshery-specific information
            (e.g., visual arrangement, comments, and so on). The resulting manifest only includes
            the core Kubernetes resource definitions, not the extra context that might be present in
            your Meshery design.
          </Box>
        </Box>
      ),
    },
    ...extensionExportOptions,
  ], [downloadModal.content, handleDesignDownload, extensionExportOptions]);

  const exportModal = useModal({
    headerIcon: <PatternIcon fill={'#fff'} height={'2rem'} width="2rem" />,
  });

  return (
    <Modal
      {...exportModal}
      open={downloadModal.open}
      onClose={handleDownloadDialogClose}
      closeModal={handleDownloadDialogClose}
      title="Export Design as..."
      aria-describedby="export-modal-description"
    >
      <ModalBody>
        <Box 
          id="export-modal-description" 
          sx={{ 
            marginBottom: theme.spacing(2),
            color: theme.palette.text.secondary,
            fontSize: theme.typography.body2.fontSize,
          }}
        >
          Choose your preferred export format below:
        </Box>
        <Box role="list" aria-label="Export options">
          {exportOptions.map((option, index) => (
            <Box role="listitem" key={option.title || index}>
              <ExportOption
                title={option.title}
                Icon={option.icon}
                content={option.content}
                disabled={option.disabled}
                description={option.description}
                onClick={(e) => option.onClick(e)}
              />
            </Box>
          ))}
        </Box>
      </ModalBody>
    </Modal>
  );
};

export default ExportModal;