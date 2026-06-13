import { useMemo, useState } from 'react';
import {
  Alert,
  Box,
  Button,
  CheckCircleIcon,
  Chip,
  CircularProgress,
  Typography,
} from '@sistent/sistent';
import { alpha, styled } from '@/theme';
import { EVENT_TYPES } from 'lib/event-types';
import { CONNECTION_STATES } from '@/utils/Enum';
import { buildKubernetesImportSummary, getKubernetesContexts } from '../ConnectionWizard.helpers';
import { formatWizardError } from './errors';
import { KubernetesImportStep, StepHeader } from '../ConnectionWizardStepContent';
import type { ConnectionExtension, GenericRecord, WizardContext, WizardStep } from './types';

const ContextRow = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: theme.spacing(2),
  padding: theme.spacing(1.5, 2),
  borderRadius: theme.spacing(1),
  border: `1px solid ${theme.palette.divider}`,
  background: theme.palette.background.card,
  transition: 'border-color 0.15s ease',
  '&:hover': {
    borderColor: theme.palette.background.brand?.default,
  },
}));

const SuccessBadge = styled(Box)(({ theme }) => ({
  width: 88,
  height: 88,
  borderRadius: '50%',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: alpha(theme.palette.background.brand?.default, 0.12),
}));

const SuccessIcon = styled(CheckCircleIcon)(({ theme }) => ({
  width: 52,
  height: 52,
  fill: theme.palette.background.brand?.default,
}));

type DiscoveredContext = {
  connectionId: string;
  name: string;
  server: string;
  status: string;
  errored: boolean;
};

const STATUS_CHIP: Record<string, 'success' | 'info' | 'warning' | 'error'> = {
  [CONNECTION_STATES.CONNECTED]: 'success',
  [CONNECTION_STATES.REGISTERED]: 'info',
  [CONNECTION_STATES.IGNORED]: 'warning',
};

/** Flattens the kubeconfig upload buckets into a single labelled list. */
const collectContexts = (result: GenericRecord | null): DiscoveredContext[] => {
  const buckets: Array<{
    key: 'connected' | 'registered' | 'ignored' | 'errored';
    status: string;
    errored: boolean;
  }> = [
    { key: 'connected', status: CONNECTION_STATES.CONNECTED, errored: false },
    { key: 'registered', status: CONNECTION_STATES.REGISTERED, errored: false },
    { key: 'ignored', status: CONNECTION_STATES.IGNORED, errored: false },
    { key: 'errored', status: 'error', errored: true },
  ];

  return buckets.flatMap(({ key, status, errored }) =>
    getKubernetesContexts(result, key).map((ctx) => ({
      connectionId: String(ctx.connectionId ?? ctx.id ?? ''),
      name: String(ctx.name ?? ctx.connectionId ?? 'context'),
      server: String(ctx.server ?? ''),
      status,
      errored,
    })),
  );
};

// ---------------------------------------------------------------------------
// 2. Kubeconfig (overrides the generic details step)
// ---------------------------------------------------------------------------

const KubeconfigStepBody = ({ ctx }: { ctx: WizardContext }) => (
  <KubernetesImportStep
    kubeconfigFile={ctx.data.kubeconfigFile}
    onPickFile={(kubeconfigFile) =>
      ctx.patch({ kubeconfigFile, registrationResult: null, registrationError: null })
    }
  />
);

const kubernetesDetailsStep: WizardStep = {
  id: 'kubeconfig',
  label: 'Import Kubeconfig',
  Component: KubeconfigStepBody,
  canProceed: (ctx) => Boolean(ctx.data.kubeconfigFile),
};

// ---------------------------------------------------------------------------
// 4a. Upload (register): discover contexts from the kubeconfig
// ---------------------------------------------------------------------------

const UploadStepBody = ({ ctx }: { ctx: WizardContext }) => (
  <Box sx={{ display: 'grid', gap: 3 }}>
    <StepHeader
      title="Discover contexts"
      subtitle={
        <>
          Meshery will read <strong>{ctx.data.kubeconfigFile?.name || 'the kubeconfig'}</strong>,
          discover every reachable context, and register each one as a Kubernetes connection.
        </>
      }
    />
    {Boolean(ctx.data.registrationError) && (
      <Alert severity="error" variant="filled">
        {formatWizardError(ctx.data.registrationError)}
      </Alert>
    )}
  </Box>
);

const kubernetesRegisterStep: WizardStep = {
  id: 'kubernetes-upload',
  label: 'Discover Contexts',
  Component: UploadStepBody,
  nextLabel: () => 'Import Contexts',
  canProceed: (ctx) => Boolean(ctx.data.kubeconfigFile),
  onNext: async (ctx) => {
    if (!ctx.data.kubeconfigFile) {
      return false;
    }
    ctx.patch({ registrationError: null });
    try {
      const result = await ctx.services.uploadKubeconfig(ctx.data.kubeconfigFile);
      ctx.patch({ registrationResult: result ?? {} });
      const summary = buildKubernetesImportSummary(result);
      ctx.services.notify({
        message:
          summary.importedCount > 0
            ? `Discovered ${summary.importedCount} Kubernetes context${summary.importedCount === 1 ? '' : 's'}.`
            : 'Kubeconfig uploaded. No reachable contexts were found.',
        event_type: summary.importedCount > 0 ? EVENT_TYPES.SUCCESS : EVENT_TYPES.WARNING,
      });
      return true;
    } catch (error) {
      ctx.patch({ registrationError: error });
      ctx.services.notify({
        message: `Failed to import kubeconfig: ${formatWizardError(error)}`,
        event_type: EVENT_TYPES.ERROR,
      });
      return false;
    }
  },
};

// ---------------------------------------------------------------------------
// 4b. Post-config: act on discovered contexts (connect / disconnect)
// ---------------------------------------------------------------------------

const ContextsStepBody = ({ ctx }: { ctx: WizardContext }) => {
  const contexts = useMemo(
    () => collectContexts(ctx.data.registrationResult),
    [ctx.data.registrationResult],
  );
  // Per-context status overrides applied during this session.
  const overrides = (ctx.data.postConfig.contextStatuses as Record<string, string>) || {};
  const [busyId, setBusyId] = useState<string | null>(null);

  const applyStatus = async (context: DiscoveredContext, status: string) => {
    if (!context.connectionId) {
      return;
    }
    setBusyId(context.connectionId);
    try {
      await ctx.services.updateConnectionById(context.connectionId, { status });
      ctx.patchPostConfig({
        contextStatuses: { ...overrides, [context.connectionId]: status },
      });
    } catch (error) {
      ctx.services.notify({
        message: `Failed to update ${context.name}: ${formatWizardError(error)}`,
        event_type: EVENT_TYPES.ERROR,
      });
    } finally {
      setBusyId(null);
    }
  };

  if (contexts.length === 0) {
    return (
      <Typography variant="body2" color="text.secondary">
        No Kubernetes contexts were discovered. Go back and choose a different kubeconfig.
      </Typography>
    );
  }

  return (
    <Box sx={{ display: 'grid', gap: 3 }}>
      <StepHeader
        title="Select clusters"
        subtitle="Choose which discovered clusters Meshery should keep connected."
      />
      <Box sx={{ display: 'grid', gap: 1.5 }}>
        {contexts.map((context) => {
          const status = overrides[context.connectionId] || context.status;
          const isConnected = status === CONNECTION_STATES.CONNECTED;
          return (
            <ContextRow key={context.connectionId || context.name}>
              <Box sx={{ display: 'grid', gap: 0.25, minWidth: 0 }}>
                <Typography variant="body1" noWrap>
                  {context.name}
                </Typography>
                <Typography variant="caption" color="text.secondary" noWrap>
                  {context.server || 'unknown server'}
                </Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                <Chip label={status} color={STATUS_CHIP[status] || 'default'} size="small" />
                {!context.errored && context.connectionId && (
                  <Button
                    variant={isConnected ? 'outlined' : 'contained'}
                    size="small"
                    disabled={busyId === context.connectionId}
                    onClick={() =>
                      applyStatus(
                        context,
                        isConnected ? CONNECTION_STATES.DISCONNECTED : CONNECTION_STATES.CONNECTED,
                      )
                    }
                  >
                    {busyId === context.connectionId ? (
                      <CircularProgress size={16} />
                    ) : isConnected ? (
                      'Disconnect'
                    ) : (
                      'Connect'
                    )}
                  </Button>
                )}
              </Box>
            </ContextRow>
          );
        })}
      </Box>
    </Box>
  );
};

const kubernetesContextsStep: WizardStep = {
  id: 'kubernetes-contexts',
  label: 'Select Clusters',
  Component: ContextsStepBody,
};

// ---------------------------------------------------------------------------
// 5. Receipt
// ---------------------------------------------------------------------------

const KubernetesReceiptBody = ({ ctx }: { ctx: WizardContext }) => {
  const summary = buildKubernetesImportSummary(ctx.data.registrationResult);
  const overrides = (ctx.data.postConfig.contextStatuses as Record<string, string>) || {};
  const connectedNow = collectContexts(ctx.data.registrationResult).filter(
    (context) =>
      (overrides[context.connectionId] || context.status) === CONNECTION_STATES.CONNECTED,
  ).length;

  return (
    <Box
      sx={{ display: 'grid', gap: 2, justifyItems: 'center', textAlign: 'center', py: 5, px: 2 }}
    >
      <SuccessBadge>
        <SuccessIcon />
      </SuccessBadge>
      <Typography variant="h6" sx={{ fontWeight: 600 }}>
        Kubernetes import complete
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ maxWidth: 420 }}>
        {summary.importedCount} context{summary.importedCount === 1 ? '' : 's'} discovered,{' '}
        {connectedNow} connected
        {summary.erroredCount > 0 ? `, ${summary.erroredCount} failed to import` : ''}.
      </Typography>
    </Box>
  );
};

const kubernetesReceiptStep: WizardStep = {
  id: 'kubernetes-receipt',
  label: 'Done',
  Component: KubernetesReceiptBody,
  nextLabel: () => 'Finish',
};

export const kubernetesExtension: ConnectionExtension = {
  match: { kind: 'kubernetes' },
  detailsStep: kubernetesDetailsStep,
  credentialStep: null, // the kubeconfig is the credential
  registerStep: kubernetesRegisterStep,
  postConfigSteps: [kubernetesContextsStep],
  receiptStep: kubernetesReceiptStep,
};
