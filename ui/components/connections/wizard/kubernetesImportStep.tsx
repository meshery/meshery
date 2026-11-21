import { useRef, useState } from 'react';
import { Alert, Box, CircularProgress, CloudUploadIcon, Typography } from '@sistent/sistent';
import { alpha, styled } from '@/theme';
import { EVENT_TYPES } from 'lib/event-types';
import { StepHeader, StepLayout } from '../ConnectionWizardStepContent';
import { formatWizardError } from './errors';
import { defaultChoice } from './kubernetesReviewStep';
import type { DiscoveredKubeContext, WizardContext, WizardStep } from './types';

const KUBECONFIG_DOCS_URL = 'https://docs.meshery.io/installation/kubernetes';

const UploadDropzone = styled('button')(({ theme }) => ({
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  gap: theme.spacing(1),
  padding: theme.spacing(4),
  borderRadius: theme.spacing(1.5),
  border: `1.5px dashed ${theme.palette.divider}`,
  background: theme.palette.background.card,
  color: theme.palette.text.primary,
  cursor: 'pointer',
  textAlign: 'center',
  transition: 'border-color 0.15s ease, background 0.15s ease',
  '&:hover': {
    borderColor: theme.palette.background.brand.default,
    background: alpha(theme.palette.background.brand.default, 0.04),
  },
}));

const UploadIconCircle = styled(Box)(({ theme }) => ({
  width: 56,
  height: 56,
  borderRadius: '50%',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: alpha(theme.palette.background.brand.default, 0.12),
}));

const UploadIcon = styled(CloudUploadIcon)(({ theme }) => ({
  width: 28,
  height: 28,
  fill: theme.palette.background.brand.default,
}));

const HiddenFileInput = styled('input')({
  display: 'none',
});

type KubernetesImportStepProps = {
  kubeconfigFile: File | null;
  previewContexts?: DiscoveredKubeContext[];
  previewLoading?: boolean;
  onPickFile: (file: File | null) => void;
};

const KubernetesImportStep = ({
  kubeconfigFile,
  previewContexts,
  previewLoading,
  onPickFile,
}: KubernetesImportStepProps) => (
  <StepLayout>
    <StepHeader
      title="Upload a kubeconfig"
      subtitle="Upload a kubeconfig file. Meshery will read the Kubernetes contexts inside and let you choose which ones to import."
    />
    <UploadDropzone
      type="button"
      onClick={() => document.getElementById('connection-wizard-kubeconfig-input')?.click()}
    >
      <UploadIconCircle>
        <UploadIcon />
      </UploadIconCircle>
      {kubeconfigFile ? (
        <Typography
          variant="body1"
          title={kubeconfigFile.name}
          sx={{
            fontWeight: 600,
            maxWidth: '100%',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
          }}
        >
          {kubeconfigFile.name}
        </Typography>
      ) : (
        <Typography variant="body1" sx={{ fontWeight: 600 }}>
          Click to choose a kubeconfig file
        </Typography>
      )}
      <Typography variant="caption" color="text.secondary">
        {kubeconfigFile
          ? 'Click to replace the selected file'
          : 'Accepts kubeconfigs with embedded certificates'}
      </Typography>
    </UploadDropzone>
    {previewLoading && (
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          gap: 1,
          mt: 2,
        }}
      >
        <CircularProgress size={18} />
        <Typography variant="body2" color="text.secondary">
          Discovering Kubernetes contexts...
        </Typography>
      </Box>
    )}
    {!previewLoading && previewContexts && previewContexts.length > 0 && (
      <Box
        sx={{
          mt: 2,
          display: 'grid',
          gap: 1,
        }}
      >
        <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
          Detected Kubernetes contexts
        </Typography>

        {previewContexts.map((context) => (
          <Box
            key={context.id}
            sx={{
              p: 1.5,
              border: (theme) => `1px solid ${theme.palette.divider}`,
              borderRadius: 1,
            }}
          >
            <Typography variant="body2" sx={{ fontWeight: 500 }}>
              {context.name}
            </Typography>

            <Typography variant="caption" color="text.secondary" component="div">
              {context.server || 'Unknown server'}
            </Typography>
          </Box>
        ))}
      </Box>
    )}
    {!previewLoading && previewContexts && previewContexts.length === 0 && (
      <Alert severity="warning">
        No Kubernetes contexts were found in the selected kubeconfig.
      </Alert>
    )}
    <HiddenFileInput
      id="connection-wizard-kubeconfig-input"
      type="file"
      onChange={(event) => onPickFile(event.target.files?.[0] || null)}
    />
  </StepLayout>
);

const KubeconfigStepBody = ({ ctx }: { ctx: WizardContext }) => {
  const [previewLoading, setPreviewLoading] = useState(false);
  const requestVersionRef = useRef(0);

  const handlePickFile = async (kubeconfigFile: File | null) => {
    const requestVersion = ++requestVersionRef.current;

    ctx.patch({
      kubeconfigFile,
      registrationResult: null,
      registrationError: null,
    });

    ctx.patchPostConfig({
      discoveredContexts: undefined,
      contextChoices: undefined,
      previewContexts: undefined,
    });

    if (!kubeconfigFile) {
      return;
    }

    setPreviewLoading(true);

    try {
      const previewContexts = await ctx.services.discoverKubeContexts(kubeconfigFile);

      if (requestVersion !== requestVersionRef.current) {
        return;
      }

      ctx.patchPostConfig({
        previewContexts,
      });
    } catch {
      // Preview discovery is best-effort. Swallow errors so they do not surface
      // as unhandled rejections; Continue performs authoritative discovery and error handling.
    } finally {
      if (requestVersion === requestVersionRef.current) {
        setPreviewLoading(false);
      }
    }
  };

  return (
    <KubernetesImportStep
      kubeconfigFile={ctx.data.kubeconfigFile}
      previewContexts={ctx.data.postConfig.previewContexts as DiscoveredKubeContext[] | undefined}
      previewLoading={previewLoading}
      onPickFile={handlePickFile}
    />
  );
};

export const kubernetesDetailsStep: WizardStep = {
  id: 'kubeconfig',
  label: 'Import Kubeconfig',
  icon: CloudUploadIcon,
  Component: KubeconfigStepBody,
  canProceed: (ctx) => Boolean(ctx.data.kubeconfigFile),
  nextLabel: () => 'Continue',
  helpText: `Upload a kubeconfig with embedded certificates. Meshery reads the file's contexts, then registers the ones you select as Kubernetes connections. [Learn more about connecting Kubernetes](${KUBECONFIG_DOCS_URL}).`,
  onNext: async (ctx) => {
    if (!ctx.data.kubeconfigFile) {
      return false;
    }
    ctx.patch({ registrationError: null });
    try {
      const discovered = await ctx.services.discoverKubeContexts(ctx.data.kubeconfigFile);
      if (discovered.length === 0) {
        ctx.services.notify({
          message: 'No Kubernetes contexts were found in the kubeconfig.',
          event_type: EVENT_TYPES.WARNING,
        });
        return false;
      }
      const contextChoices = Object.fromEntries(
        discovered.map((context) => [context.id, defaultChoice(context.name)]),
      );
      ctx.patchPostConfig({
        discoveredContexts: discovered,
        contextChoices,
        connectOnImport: true,
      });
      return true;
    } catch (error) {
      ctx.patch({ registrationError: error });
      ctx.services.notify({
        message: `Failed to read kubeconfig: ${formatWizardError(error)}`,
        event_type: EVENT_TYPES.ERROR,
      });
      return false;
    }
  },
};
