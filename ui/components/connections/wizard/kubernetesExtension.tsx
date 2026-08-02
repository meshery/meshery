import { useMemo, useState } from 'react';
import Link from 'next/link';
import {
  Box,
  CheckCircleIcon,
  Chip,
  CircularProgress,
  Typography,
  SettingsIcon,
} from '@sistent/sistent';
import { useSelector } from 'react-redux';
import { alpha, styled } from '@/theme';
import { CONNECTION_STATES } from '@/utils/Enum';
import { CONNECTIONS_PATH, getKubernetesContexts } from '../ConnectionWizard.helpers';
import { kubernetesSettingsStep } from './kubernetesSettings';
import FormatConnectionMetadata from '../metadata';
import { ConnectionStateChip } from '../ConnectionChip';
import { StepHeader } from '../ConnectionWizardStepContent';
import { kubernetesDetailsStep } from './kubernetesImportStep';
import { kubernetesReviewStep } from './kubernetesReviewStep';
import type { ConnectionExtension, GenericRecord, WizardContext, WizardStep } from './types';

const KUBERNETES_CONNECTION_DOCS_URL =
  'https://docs.meshery.io/guides/infrastructure-management/kubernetes-connection-lifecycle';

const ContextRow = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'muted',
})<{ muted?: boolean }>(({ theme, muted }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(1.5),
  padding: theme.spacing(1.5, 2),
  borderRadius: theme.spacing(1),
  border: `1px solid ${theme.palette.divider}`,
  background: theme.palette.background.card,
  opacity: muted ? 0.6 : 1,
  transition: 'border-color 0.15s ease, opacity 0.15s ease',
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

const ConnectionsLink = styled(Link)(({ theme }) => ({
  color: theme.palette.background.brand?.default || theme.palette.primary.main,
  fontWeight: 600,
  textDecoration: 'underline',
  textUnderlineOffset: 2,
}));

type DiscoveredContext = {
  connectionId: string;
  name: string;
  server: string;
  status: string;
  errored: boolean;
};

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

const ContextsStepBody = ({ ctx }: { ctx: WizardContext }) => {
  const contexts = useMemo(
    () => collectContexts(ctx.data.registrationResult),
    [ctx.data.registrationResult],
  );
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
        message: `Failed to update ${context.name}: ${error}`,
        event_type: 'error',
      });
    } finally {
      setBusyId(null);
    }
  };

  if (contexts.length === 0) {
    return (
      <Typography variant="body2" color="text.secondary">
        No Kubernetes contexts are associated with this connection.
      </Typography>
    );
  }

  return (
    <Box sx={{ display: 'grid', gap: 3 }}>
      <StepHeader
        title="Manage clusters"
        subtitle="Choose which discovered clusters Meshery should keep connected."
      />
      <Box sx={{ display: 'grid', gap: 1.5 }}>
        {contexts.map((context) => {
          const status = overrides[context.connectionId] || context.status;
          const isConnected = status === CONNECTION_STATES.CONNECTED;
          return (
            <ContextRow key={context.connectionId || context.name}>
              <Box sx={{ display: 'grid', gap: 0.25, minWidth: 0, flex: 1 }}>
                <Typography variant="body1" noWrap>
                  {context.name}
                </Typography>
                <Typography variant="caption" color="text.secondary" noWrap>
                  {context.server || 'unknown server'}
                </Typography>
              </Box>
              <ConnectionStateChip status={status} />
              {!context.errored && context.connectionId && (
                <Box
                  component="button"
                  type="button"
                  disabled={busyId === context.connectionId}
                  onClick={() =>
                    applyStatus(
                      context,
                      isConnected ? CONNECTION_STATES.DISCONNECTED : CONNECTION_STATES.CONNECTED,
                    )
                  }
                  sx={{
                    minWidth: 96,
                    border: (theme) => `1px solid ${theme.palette.background.brand?.default}`,
                    borderRadius: 1,
                    py: 0.5,
                    px: 1.5,
                    cursor: 'pointer',
                    background: (theme) =>
                      isConnected ? 'transparent' : theme.palette.background.brand?.default,
                    color: (theme) =>
                      isConnected ? theme.palette.text.primary : theme.palette.common.white,
                  }}
                >
                  {busyId === context.connectionId ? (
                    <CircularProgress size={16} />
                  ) : isConnected ? (
                    'Disconnect'
                  ) : (
                    'Connect'
                  )}
                </Box>
              )}
            </ContextRow>
          );
        })}
      </Box>
    </Box>
  );
};

const kubernetesContextsStep: WizardStep = {
  id: 'kubernetes-contexts',
  label: 'Manage Clusters',
  icon: SettingsIcon,
  Component: ContextsStepBody,
  hidden: (ctx) =>
    ctx.mode === 'create' || collectContexts(ctx.data.registrationResult).length === 0,
};

const KubernetesReceiptBody = ({ ctx }: { ctx: WizardContext }) => {
  const controllerState = useSelector(
    (state: { ui: { controllerState: unknown } }) => state.ui.controllerState,
  );

  const configuredConnection = (ctx.data.registrationResult as GenericRecord) || {};
  if (
    ctx.mode === 'configure' &&
    configuredConnection.kind === 'kubernetes' &&
    configuredConnection.id
  ) {
    return (
      <Box sx={{ display: 'grid', gap: 2 }}>
        <StepHeader
          title="Connection details"
          subtitle="Live status and diagnostics for this Kubernetes connection."
        />
        <FormatConnectionMetadata
          connection={configuredConnection}
          meshsyncControllerState={controllerState}
        />
      </Box>
    );
  }

  const imported =
    (ctx.data.postConfig.importedContexts as { name: string; status: string }[]) || [];
  const created = Number(ctx.data.postConfig.createdCount ?? imported.length);
  const connected = Number(ctx.data.postConfig.connectedCount ?? 0);
  const unreachable = Number(ctx.data.postConfig.unreachableCount ?? 0);
  const isImport = ctx.data.postConfig.createdCount !== undefined;

  const summaryParts = [
    `${created} imported`,
    connected > 0 ? `${connected} connected` : null,
    unreachable > 0 ? `${unreachable} not found` : null,
  ].filter(Boolean);

  return (
    <Box
      sx={{ display: 'grid', gap: 2, justifyItems: 'center', textAlign: 'center', py: 4, px: 2 }}
    >
      <SuccessBadge>
        <SuccessIcon />
      </SuccessBadge>
      <Typography variant="h6" sx={{ fontWeight: 600 }}>
        {isImport ? 'Kubernetes import complete' : 'Configuration saved'}
      </Typography>
      {isImport && (
        <Typography variant="body2" color="text.secondary" sx={{ maxWidth: 440 }}>
          {summaryParts.join(' · ')}. Manage these anytime from{' '}
          <ConnectionsLink href={CONNECTIONS_PATH}>connections</ConnectionsLink>.
        </Typography>
      )}
      {imported.length > 0 && (
        <Box sx={{ display: 'grid', gap: 1, width: '100%', maxWidth: 460, textAlign: 'left' }}>
          {imported.map((context, index) => (
            <Box
              key={`${context.name}-${index}`}
              sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                gap: 1.5,
                minWidth: 0,
              }}
            >
              <Chip
                label={context.name || 'Unnamed context'}
                size="small"
                title={context.name}
                sx={{
                  minWidth: 0,
                  flex: '1 1 auto',
                  '& .MuiChip-label': {
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                  },
                }}
              />
              <ConnectionStateChip status={context.status} />
            </Box>
          ))}
        </Box>
      )}
    </Box>
  );
};

const kubernetesReceiptStep: WizardStep = {
  id: 'kubernetes-receipt',
  label: 'Done',
  Component: KubernetesReceiptBody,
  nextLabel: () => 'Finish',
  helpText: `Import complete. Manage connections anytime from the Connections page. [Learn more about Kubernetes connection lifecycle](${KUBERNETES_CONNECTION_DOCS_URL}).`,
};

export const kubernetesExtension: ConnectionExtension = {
  match: { kind: 'kubernetes' },
  detailsStep: kubernetesDetailsStep,
  credentialStep: null,
  registerStep: kubernetesReviewStep,
  postConfigSteps: [kubernetesSettingsStep, kubernetesContextsStep],
  receiptStep: kubernetesReceiptStep,
};
