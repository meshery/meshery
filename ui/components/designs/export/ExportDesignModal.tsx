/**
 * ExportDesignModal — design-domain export flow.
 *
 * Replaces the legacy `components/shared/Modal/ExportModal.tsx`. Rendered as a
 * list of export options inside the shared `Modal` primitive from
 * `@/components/shared/Modal` so it inherits the standard header/body styling,
 * size tokens, and a11y wiring. Each option fires its own download handler;
 * there is no aggregate submit, so the base `Modal` (no footer) is the right
 * shared primitive rather than `FormModal`.
 *
 * User-facing export choices are preserved verbatim: Meshery YAML, OCI image,
 * Kubernetes manifest, Helm chart, plus any extension-supplied options.
 */
import { FC, MouseEvent, ReactNode, useMemo, memo } from 'react';
import {
  Box,
  DownloadIcon,
  IconButton,
  InfoTooltip,
  ListItem,
  ListItemIcon,
  ListItemText,
} from '@sistent/sistent';
import { styled, useTheme } from '@/theme';
import { Modal } from '@/components/shared/Modal';
import PatternIcon from '@/assets/icons/Pattern';
import KubernetesIcon from '@/assets/icons/technology/kubernetes';
import HelmIcon from '@/assets/icons/technology/HelmIcon';
import { OCIImageIcon } from '@/assets/icons/OciImage';

export interface ExportDesignOption {
  title: string;
  icon: ReactNode;
  onClick: (e: MouseEvent) => void;
  description?: string | ReactNode;
  disabled?: boolean;
  content?: ReactNode;
}

interface ExportDesignOptionRowProps {
  title: string;
  Icon: ReactNode;
  onClick: (e: MouseEvent) => void;
  content?: ReactNode;
  disabled?: boolean;
  description?: string | ReactNode;
}

const OptionListItem = styled(ListItem, {
  shouldForwardProp: (prop) => prop !== 'isDisabled',
})<{ isDisabled?: boolean }>(({ theme, isDisabled }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: '0.5rem',
  borderRadius: '0.25rem',
  opacity: isDisabled ? 0.5 : 1,
  marginBottom: '1rem',
  border: `1px solid ${theme.palette.border.normal}`,
  boxShadow: theme.shadows[1],
  cursor: isDisabled ? 'not-allowed' : 'default',
  '&:hover': {
    borderColor: theme.palette.border.brand,
  },
}));

const OptionRow = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  width: '100%',
  padding: theme.spacing(1),
  justifyContent: 'space-between',
}));

const OptionIconSlot = styled(ListItemIcon)({
  minWidth: '2rem',
  marginRight: '1rem',
});

const OptionDownloadButton = styled(IconButton)({
  marginLeft: '0.5rem',
});

const ExportDesignOptionRow: FC<ExportDesignOptionRowProps> = ({
  title,
  Icon,
  onClick,
  content,
  disabled = false,
  description = 'Download the design in the selected format',
}) => {
  const theme = useTheme();
  return (
    <OptionListItem isDisabled={disabled}>
      <OptionRow>
        <OptionIconSlot>{Icon}</OptionIconSlot>
        <ListItemText primary={title} />
        <InfoTooltip placement="top" title={description} content={description} />
        <OptionDownloadButton disabled={disabled} onClick={onClick}>
          <DownloadIcon fill={theme.palette.icon.default} />
        </OptionDownloadButton>
      </OptionRow>
      {content}
    </OptionListItem>
  );
};

export interface ExportDesignModalProps {
  /** Controls visibility plus the design content the user is exporting. */
  downloadModal: { open: boolean; content: unknown };
  handleDownloadDialogClose: () => void;
  handleDesignDownload: (e: MouseEvent, content: unknown, arg?: null, params?: string) => void;
  /** Extra export options supplied by extension hosts (e.g. Meshery Cloud). */
  extensionExportOptions?: ExportDesignOption[];
}

const ExportDesignModalComponent: FC<ExportDesignModalProps> = ({
  downloadModal,
  handleDownloadDialogClose,
  handleDesignDownload,
  extensionExportOptions = [],
}) => {
  const theme = useTheme();

  const baseOptions: ExportDesignOption[] = useMemo(
    () => [
      {
        title: 'Meshery Design (yaml)',
        icon: <PatternIcon width="30" height="30" fill={theme.palette.primary.main} />,
        onClick: (e) => handleDesignDownload(e, downloadModal.content),
        description:
          'Export your design as a complete, self-contained Meshery Design file (YAML). This file includes embedded images and all configuration details. It&apos;s the perfect format for creating backups, sharing with colleagues using Meshery, or transferring designs between Meshery environments without losing any information (lossless transfer).',
      },
      {
        title: 'Meshery Design (OCI image)',
        icon: <OCIImageIcon width={30} height={30} />,
        onClick: (e) =>
          handleDesignDownload(
            e,
            downloadModal.content,
            null,
            new URLSearchParams({ oci: 'true' }).toString(),
          ),
        description:
          'Download your design as an OCI compatible container image, which can be pushed to and pulled from container registries like Docker Hub, AWS ECR, and so on.',
      },
      {
        title: 'Kubernetes Manifest (yaml)',
        icon: <KubernetesIcon width="30" height="30" />,
        onClick: (e) =>
          handleDesignDownload(
            e,
            downloadModal.content,
            null,
            new URLSearchParams({ export: 'Kubernetes Manifest' }).toString(),
          ),
        description: (
          <div>
            <p>
              Download your design as a standard Kubernetes Manifest file. This file contains the
              Kubernetes resource definitions from your design and can be directly applied to a
              cluster using tools like `kubectl`.
            </p>
            <p>
              <strong>Lossy Export:</strong> This process strips out Meshery-specific information
              (e.g., visual arrangement, comments, and so on). The resulting manifest only includes
              the core Kubernetes resource definitions, not the extra context that might be present
              in your Meshery design.
            </p>
          </div>
        ),
      },
      {
        title: 'Helm Chart (tar.gz)',
        icon: <HelmIcon width="30" height="30" />,
        onClick: (e) =>
          handleDesignDownload(
            e,
            downloadModal.content,
            null,
            new URLSearchParams({ export: 'helm-chart' }).toString(),
          ),
        disabled: false,
        description: (
          <div>
            <p>Download your design as a Helm Chart.</p>
            <p>
              <strong>Lossy Export:</strong> This process strips out Meshery-specific information
              (e.g., visual arrangement, comments, and so on). The resulting manifest only includes
              the core Kubernetes resource definitions, not the extra context that might be present
              in your Meshery design.
            </p>
          </div>
        ),
      },
      ...extensionExportOptions,
    ],
    [theme, handleDesignDownload, downloadModal.content, extensionExportOptions],
  );

  return (
    <Modal
      isOpen={downloadModal.open}
      onClose={handleDownloadDialogClose}
      title="Export Design as..."
      headerIcon={<PatternIcon fill={theme.palette.common.white} height="2rem" width="2rem" />}
    >
      {baseOptions.map((option) => (
        <ExportDesignOptionRow
          key={option.title}
          title={option.title}
          Icon={option.icon}
          content={option.content}
          disabled={option.disabled}
          description={option.description}
          onClick={option.onClick}
        />
      ))}
    </Modal>
  );
};

ExportDesignModalComponent.displayName = 'ExportDesignModal';

export const ExportDesignModal = memo(ExportDesignModalComponent);

/**
 * Unwrapped (non-memo) variant for use in remote extension `injectProps`.
 * Remote bundles (e.g. Kanvas) may call the component directly as a function;
 * `memo()` returns a plain object without `.apply`, which causes
 * "i.apply is not a function" at runtime in the remote bundle.
 */
export { ExportDesignModalComponent };

export default ExportDesignModal;
