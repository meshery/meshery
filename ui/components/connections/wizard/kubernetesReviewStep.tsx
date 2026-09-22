import {
  Alert,
  Box,
  Checkbox,
  CustomTooltip,
  MenuItem,
  TextField,
  Typography,
  AssignmentTurnedInIcon,
} from '@sistent/sistent';
import { alpha, styled } from '@/theme';
import { EVENT_TYPES } from 'lib/event-types';
import { CONNECTION_STATES } from '@/utils/Enum';
import { getKubernetesContexts, kubernetesImportedNotify } from '../ConnectionWizard.helpers';
import { formatWizardError } from './errors';
import {
  DEFAULT_MESHSYNC_DEPLOYMENT_MODE,
  getMeshsyncModeTooltip,
  MESHSYNC_DEPLOYMENT_MODE_OPTIONS,
  MESHSYNC_MODES_DOCS_URL,
} from './kubernetesDeploymentMode';
import { ConnectionStateChip } from '../ConnectionChip';
import { StepHeader } from '../ConnectionWizardStepContent';
import type { DiscoveredKubeContext, GenericRecord, WizardContext, WizardStep } from './types';

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

const ConnectNotice = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(1),
  borderRadius: theme.spacing(1),
  border: `1px solid ${theme.palette.divider}`,
  background: alpha(theme.palette.info.main, 0.06),
  padding: theme.spacing(1, 1.5),
}));

type ContextChoice = { selected: boolean; name: string; meshsyncDeploymentMode: string };

const getDiscovered = (ctx: WizardContext): DiscoveredKubeContext[] =>
  (ctx.data.postConfig.discoveredContexts as DiscoveredKubeContext[]) || [];

const getChoices = (ctx: WizardContext): Record<string, ContextChoice> =>
  (ctx.data.postConfig.contextChoices as Record<string, ContextChoice>) || {};

const getConnectOnImport = (ctx: WizardContext): boolean =>
  ctx.data.postConfig.connectOnImport !== false;

export const defaultChoice = (name: string): ContextChoice => ({
  selected: true,
  name,
  meshsyncDeploymentMode: DEFAULT_MESHSYNC_DEPLOYMENT_MODE,
});

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
                  inputProps={{
                    'aria-label': `Name for ${context.name}`,
                    title: choice.name,
                  }}
                  sx={{
                    maxWidth: '100%',
                    '& .MuiInput-input': {
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
                    },
                    '& .MuiInput-underline:before, & .MuiInput-underline:after': {
                      maxWidth: '100%',
                    },
                  }}
                />
                <Typography variant="caption" color="text.secondary" noWrap component="div">
                  {context.server || 'unknown server'}
                </Typography>
              </Box>
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
                    <CustomTooltip
                      interactive
                      title={getMeshsyncModeTooltip(option)}
                      placement="left"
                    >
                      <Box component="span" sx={{ display: 'block', width: '100%' }}>
                        {option.label}
                      </Box>
                    </CustomTooltip>
                  </MenuItem>
                ))}
              </TextField>
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

export const kubernetesReviewStep: WizardStep = {
  id: 'kubernetes-review',
  label: 'Review Contexts',
  icon: AssignmentTurnedInIcon,
  Component: ReviewContextsStepBody,
  nextLabel: () => 'Import',
  helpText: `Select which contexts to import, choose a MeshSync mode (embedded or operator), and decide whether reachable clusters should connect immediately. Hover a MeshSync option for details. [MeshSync modes](${MESHSYNC_MODES_DOCS_URL}). [Kubernetes connection lifecycle](${KUBERNETES_CONNECTION_DOCS_URL}).`,
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
      ctx.services.notify(kubernetesImportedNotify(created.length));
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
