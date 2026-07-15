import { useMemo, useState } from 'react';
import {
  Alert,
  Box,
  Checkbox,
  CheckCircleIcon,
  CircularProgress,
  MenuItem,
  TextField,
  Typography,
  CloudUploadIcon,
  AssignmentTurnedInIcon,
  SettingsIcon,
} from '@sistent/sistent';
import { useSelector } from 'react-redux';
import { alpha, styled } from '@/theme';
import { EVENT_TYPES } from 'lib/event-types';
import { CONNECTION_STATES } from '@/utils/Enum';
import { getKubernetesContexts } from '../ConnectionWizard.helpers';
import { formatWizardError } from './errors';
import {
  DEFAULT_MESHSYNC_DEPLOYMENT_MODE,
  MESHSYNC_DEPLOYMENT_MODE_OPTIONS,
} from './kubernetesDeploymentMode';
import { kubernetesSettingsStep } from './kubernetesSettings';
import FormatConnectionMetadata from '../metadata';
import { ConnectionStateChip } from '../ConnectionChip';
import { KubernetesImportStep, StepHeader } from '../ConnectionWizardStepContent';
import type {
  ConnectionExtension,
  DiscoveredKubeContext,
  GenericRecord,
  WizardContext,
  WizardStep,
} from './types';

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

const ConnectNotice = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(1),
  borderRadius: theme.spacing(1),
  border: `1px solid ${theme.palette.divider}`,
  background: alpha(theme.palette.info.main, 0.06),
  padding: theme.spacing(1, 1.5),
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

const SummaryList = styled(Box)(({ theme }) => ({
  width: '100%',
  maxWidth: 460,
  borderRadius: theme.spacing(1.5),
  border: `1px solid ${theme.palette.divider}`,
  background: theme.palette.background.card,
  overflow: 'hidden',
  textAlign: 'left',
}));

const SummaryListRow = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: theme.spacing(2),
  padding: theme.spacing(1.25, 2),
  '& + &': {
    borderTop: `1px solid ${theme.palette.divider}`,
  },
}));

type ContextChoice = { selected: boolean; name: string; meshsyncDeploymentMode: string };

// ---------------------------------------------------------------------------
// postConfig accessors. The kubernetes flow keeps its working state on the
// wizard's free-form postConfig bag.
// ---------------------------------------------------------------------------

const getDiscovered = (ctx: WizardContext): DiscoveredKubeContext[] =>
  (ctx.data.postConfig.discoveredContexts as DiscoveredKubeContext[]) || [];

const getChoices = (ctx: WizardContext): Record<string, ContextChoice> =>
  (ctx.data.postConfig.contextChoices as Record<string, ContextChoice>) || {};

const getConnectOnImport = (ctx: WizardContext): boolean =>
  ctx.data.postConfig.connectOnImport !== false;

// A discovered context's choice, defaulted for contexts not yet touched.
const defaultChoice = (name: string): ContextChoice => ({
  selected: true,
  name,
  meshsyncDeploymentMode: DEFAULT_MESHSYNC_DEPLOYMENT_MODE,
});

// ---------------------------------------------------------------------------
// 2. Import kubeconfig (overrides the generic details step). Picks a file and,
// on Next, discovers its contexts WITHOUT persisting them.
// ---------------------------------------------------------------------------

const KubeconfigStepBody = ({ ctx }: { ctx: WizardContext }) => (
  <KubernetesImportStep
    kubeconfigFile={ctx.data.kubeconfigFile}
    onPickFile={(kubeconfigFile) => {
      ctx.patch({ kubeconfigFile, registrationResult: null, registrationError: null });
      // Drop any previously discovered contexts so a new file is re-discovered.
      ctx.patchPostConfig({ discoveredContexts: undefined, contextChoices: undefined });
    }}
  />
);

const kubernetesDetailsStep: WizardStep = {
  id: 'kubeconfig',
  label: 'Import Kubeconfig',
  icon: CloudUploadIcon,
  Component: KubeconfigStepBody,
  canProceed: (ctx) => Boolean(ctx.data.kubeconfigFile),
  nextLabel: () => 'Discover Contexts',
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

// ---------------------------------------------------------------------------
// 3. Review contexts (overrides the generic register step). Select which
// discovered contexts to import, rename them, and choose whether reachable
// clusters are connected on import. On Next the selected contexts are created
// and the reachable ones (when chosen) are transitioned to connected.
// ---------------------------------------------------------------------------

const ReviewContextsStepBody = ({ ctx }: { ctx: WizardContext }) => {
  const discovered = getDiscovered(ctx);
  const choices = getChoices(ctx);
  const connectOnImport = getConnectOnImport(ctx);
  const hasReachable = discovered.some((context) => context.reachable);

  const updateChoice = (id: string, partial: Partial<ContextChoice>) =>
    ctx.patchPostConfig({
      contextChoices: { ...choices, [id]: { ...choices[id], ...partial } },
    });

  return (
    <Box sx={{ display: 'grid', gap: 3 }}>
      <StepHeader
        title="Review contexts"
        subtitle="Choose which clusters to import, rename them if you like, and decide whether reachable clusters should be connected."
      />
      {Boolean(ctx.data.registrationError) && (
        <Alert severity="error" variant="filled">
          {formatWizardError(ctx.data.registrationError)}
        </Alert>
      )}
      <Box sx={{ display: 'grid', gap: 1.5 }}>
        {discovered.map((context) => {
          const choice = choices[context.id] || defaultChoice(context.name);
          return (
            <ContextRow key={context.id} muted={!choice.selected}>
              <Checkbox
                checked={choice.selected}
                onChange={(event) => updateChoice(context.id, { selected: event.target.checked })}
                sx={{ p: 0 }}
              />
              <Box sx={{ flex: 1, minWidth: 0 }}>
                <TextField
                  variant="standard"
                  fullWidth
                  value={choice.name}
                  disabled={!choice.selected}
                  onChange={(event) => updateChoice(context.id, { name: event.target.value })}
                  inputProps={{ 'aria-label': `Name for ${context.name}` }}
                />
                <Typography variant="caption" color="text.secondary" noWrap component="div">
                  {context.server || 'unknown server'}
                </Typography>
              </Box>
              {/* Per-context MeshSync deployment mode; persisted to the
                  connection's metadata when the context is imported. */}
              <TextField
                select
                variant="standard"
                label="MeshSync"
                value={choice.meshsyncDeploymentMode}
                disabled={!choice.selected}
                onChange={(event) =>
                  updateChoice(context.id, { meshsyncDeploymentMode: event.target.value })
                }
                sx={{ minWidth: 120 }}
                inputProps={{ 'aria-label': `MeshSync deployment mode for ${context.name}` }}
              >
                {MESHSYNC_DEPLOYMENT_MODE_OPTIONS.map((option) => (
                  <MenuItem key={option.value} value={option.value}>
                    {option.label}
                  </MenuItem>
                ))}
              </TextField>
              {/* Reachable contexts land in the Discovered state; unreachable
                  ones could not be contacted, mirroring the Not Found state. */}
              <ConnectionStateChip
                status={
                  context.reachable ? CONNECTION_STATES.DISCOVERED : CONNECTION_STATES.NOTFOUND
                }
              />
            </ContextRow>
          );
        })}
      </Box>
      <ConnectNotice
        onClick={() => hasReachable && ctx.patchPostConfig({ connectOnImport: !connectOnImport })}
        sx={{ cursor: hasReachable ? 'pointer' : 'default' }}
      >
        <Checkbox
          checked={connectOnImport}
          disabled={!hasReachable}
          onChange={(event) => ctx.patchPostConfig({ connectOnImport: event.target.checked })}
          onClick={(event) => event.stopPropagation()}
          sx={{ p: 0 }}
        />
        <Typography variant="body2">
          Connect reachable clusters after import. Unreachable clusters are imported as discovered
          and can be connected later.
        </Typography>
      </ConnectNotice>
    </Box>
  );
};

/** Connections returned by the upload endpoint, with how they should be acted on. */
type CreatedContext = {
  connectionId: string;
  name: string;
  reachable: boolean;
  alreadyConnected: boolean;
};

const collectCreated = (result: GenericRecord | null): CreatedContext[] => {
  const toCreated = (records: GenericRecord[], alreadyConnected: boolean): CreatedContext[] =>
    records.map((record) => ({
      connectionId: String(record.connectionId ?? record.id ?? ''),
      name: String(record.name ?? ''),
      reachable: Boolean(record.reachable),
      alreadyConnected,
    }));

  return [
    ...toCreated(getKubernetesContexts(result, 'registered'), false),
    ...toCreated(getKubernetesContexts(result, 'connected'), true),
  ];
};

const kubernetesReviewStep: WizardStep = {
  id: 'kubernetes-review',
  label: 'Review Contexts',
  icon: AssignmentTurnedInIcon,
  Component: ReviewContextsStepBody,
  nextLabel: () => 'Import',
  canProceed: (ctx) => getDiscovered(ctx).some((context) => getChoices(ctx)[context.id]?.selected),
  onNext: async (ctx) => {
    const discovered = getDiscovered(ctx);
    const choices = getChoices(ctx);
    const selected = discovered.filter((context) => choices[context.id]?.selected);
    if (!ctx.data.kubeconfigFile || selected.length === 0) {
      return false;
    }
    ctx.patch({ registrationError: null });
    try {
      const contexts = Object.fromEntries(
        selected.map((context) => {
          const choice = choices[context.id];
          return [
            context.id,
            {
              name: choice.name.trim() || context.name,
              meshsyncDeploymentMode: choice.meshsyncDeploymentMode,
            },
          ];
        }),
      );
      const result = await ctx.services.uploadKubeconfig(ctx.data.kubeconfigFile, {
        selectedContextIds: selected.map((context) => context.id),
        contexts,
      });
      ctx.patch({ registrationResult: result ?? {} });

      const created = collectCreated(result ?? {});
      const connectedIds = new Set(
        created
          .filter((context) => context.alreadyConnected)
          .map((context) => context.connectionId),
      );
      if (getConnectOnImport(ctx)) {
        const toConnect = created.filter(
          (context) => context.reachable && !context.alreadyConnected && context.connectionId,
        );
        const outcomes = await Promise.allSettled(
          toConnect.map((context) =>
            ctx.services.updateConnectionById(context.connectionId, {
              status: CONNECTION_STATES.CONNECTED,
            }),
          ),
        );
        toConnect.forEach((context, index) => {
          if (outcomes[index].status === 'fulfilled') {
            connectedIds.add(context.connectionId);
          }
        });
      }

      // Per-context final state, mapped to the canonical connection states so
      // the receipt can reuse the shared status chip.
      const importedContexts = created.map((context) => ({
        name: context.name,
        status: connectedIds.has(context.connectionId)
          ? CONNECTION_STATES.CONNECTED
          : context.reachable
            ? CONNECTION_STATES.DISCOVERED
            : CONNECTION_STATES.NOTFOUND,
      }));

      ctx.patchPostConfig({
        importedContexts,
        createdCount: created.length,
        connectedCount: connectedIds.size,
        unreachableCount: created.filter((context) => !context.reachable).length,
      });
      ctx.services.notify({
        message: `Imported ${created.length} Kubernetes connection${created.length === 1 ? '' : 's'}.`,
        event_type: created.length > 0 ? EVENT_TYPES.SUCCESS : EVENT_TYPES.WARNING,
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
// Post-config (configure mode only): act on an already-registered connection's
// contexts. Hidden during creation, where the review step covers connecting.
// ---------------------------------------------------------------------------

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
  // The creation flow handles selecting + connecting in the review step; this
  // step only adds value when (re)configuring an existing connection that
  // actually carries discovered contexts to manage.
  hidden: (ctx) =>
    ctx.mode === 'create' || collectContexts(ctx.data.registrationResult).length === 0,
};

// ---------------------------------------------------------------------------
// 5. Receipt
// ---------------------------------------------------------------------------

const KubernetesReceiptBody = ({ ctx }: { ctx: WizardContext }) => {
  const controllerState = useSelector(
    (state: { ui: { controllerState: unknown } }) => state.ui.controllerState,
  );

  // Configure mode operates on a single connection: show the same live detail
  // "receipt" that the Connections table renders when a row is expanded —
  // status chips, controller versions, and the Diagnostics section.
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
    unreachable > 0 ? `${unreachable} pending` : null,
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
          {summaryParts.join(' · ')}. Manage these connections anytime from the Connections table.
        </Typography>
      )}
      {imported.length > 0 && (
        <SummaryList>
          {imported.map((context) => (
            <SummaryListRow key={context.name}>
              <Typography variant="body2" noWrap sx={{ minWidth: 0 }}>
                {context.name}
              </Typography>
              <ConnectionStateChip status={context.status} />
            </SummaryListRow>
          ))}
        </SummaryList>
      )}
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
  registerStep: kubernetesReviewStep,
  postConfigSteps: [kubernetesSettingsStep, kubernetesContextsStep],
  receiptStep: kubernetesReceiptStep,
};
