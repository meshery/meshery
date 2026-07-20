import { Alert, Box, CustomTooltip, Radio, Typography, SettingsIcon } from '@sistent/sistent';
import { alpha, styled } from '@/theme';
import { EVENT_TYPES } from 'lib/event-types';
import { MESHSYNC_DEPLOYMENT_TYPE } from '@/utils/Enum';
import { formatWizardError } from './errors';
import { StepHeader } from '../ConnectionWizardStepContent';
import type { GenericRecord, WizardContext, WizardStep } from './types';

// Matches the backend's MeshsyncDeploymentModeDefault (embedded) used when no
// mode is supplied. Keep in sync with server/models/connections.
export const DEFAULT_MESHSYNC_DEPLOYMENT_MODE = MESHSYNC_DEPLOYMENT_TYPE.EMBEDDED;

/** Canonical docs for MeshSync + deployment modes (anchor deep links). */
export const MESHSYNC_DOCS_URL = 'https://docs.meshery.io/concepts/architecture/meshsync';
export const MESHSYNC_MODES_DOCS_URL =
  'https://docs.meshery.io/concepts/architecture/meshsync#meshsync-deployment-mode';
export const MESHSYNC_OPERATOR_MODE_DOCS_URL =
  'https://docs.meshery.io/concepts/architecture/meshsync#operator-mode';
export const MESHSYNC_EMBEDDED_MODE_DOCS_URL =
  'https://docs.meshery.io/concepts/architecture/meshsync#embedded-mode-default';

const ModeCard = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'selected',
})<{ selected?: boolean }>(({ theme, selected }) => ({
  display: 'flex',
  alignItems: 'flex-start',
  gap: theme.spacing(1.5),
  padding: theme.spacing(2),
  borderRadius: theme.spacing(1),
  cursor: 'pointer',
  border: `1px solid ${selected ? theme.palette.background.brand?.default : theme.palette.divider}`,
  background: selected
    ? alpha(theme.palette.background.brand?.default, 0.08)
    : theme.palette.background.card,
  transition: 'border-color 0.15s ease, background 0.15s ease',
  '&:hover': {
    borderColor: theme.palette.background.brand?.default,
  },
}));

export const MESHSYNC_DEPLOYMENT_MODE_OPTIONS = [
  {
    value: MESHSYNC_DEPLOYMENT_TYPE.OPERATOR,
    label: 'Operator',
    description:
      'Installs Meshery Operator in the cluster, along with the MeshSync controller, which subscribes to and streams updates on "interesting" resources in real-time. You can configure and refine which resources MeshSync subscribes to.',
    docsUrl: MESHSYNC_OPERATOR_MODE_DOCS_URL,
  },
  {
    value: MESHSYNC_DEPLOYMENT_TYPE.EMBEDDED,
    label: 'Embedded',
    description:
      'Runs MeshSync inside Meshery Server (default). Nothing is installed into the cluster.',
    docsUrl: MESHSYNC_EMBEDDED_MODE_DOCS_URL,
  },
] as const;

/** Hover copy for a mode option, with a docs deep link. */
export const getMeshsyncModeTooltip = (option: { description: string; docsUrl: string }): string =>
  `${option.description} [Learn more](${option.docsUrl})`;

/**
 * Presentational picker for MeshSync deployment mode. Mode details appear on
 * hover (CustomTooltip), not as always-on body copy.
 */
export const MeshsyncDeploymentModePicker = ({
  value,
  currentValue,
  onChange,
}: {
  value: string;
  currentValue?: string;
  onChange: (mode: string) => void;
}) => (
  <Box sx={{ display: 'grid', gap: 1.5 }}>
    {MESHSYNC_DEPLOYMENT_MODE_OPTIONS.map((option) => {
      const selected = value === option.value;
      return (
        <CustomTooltip
          key={option.value}
          interactive
          title={getMeshsyncModeTooltip(option)}
          placement="right"
        >
          <ModeCard
            selected={selected}
            role="radio"
            aria-checked={selected}
            tabIndex={0}
            onClick={() => onChange(option.value)}
            onKeyDown={(event) => {
              if (event.key === 'Enter' || event.key === ' ') {
                event.preventDefault();
                onChange(option.value);
              }
            }}
          >
            <Radio
              checked={selected}
              tabIndex={-1}
              onChange={() => onChange(option.value)}
              sx={{ p: 0, mt: 0.25 }}
            />
            <Box sx={{ display: 'grid', gap: 0.25, minWidth: 0 }}>
              <Typography variant="body1" sx={{ fontWeight: 600 }}>
                {option.label}
                {option.value === currentValue && (
                  <Typography
                    component="span"
                    variant="caption"
                    color="text.secondary"
                    sx={{ ml: 1 }}
                  >
                    Current
                  </Typography>
                )}
              </Typography>
            </Box>
          </ModeCard>
        </CustomTooltip>
      );
    })}
  </Box>
);

export const getConfiguredConnection = (ctx: WizardContext): GenericRecord =>
  (ctx.data.registrationResult as GenericRecord) || {};

// The persisted mode, accepting either the camelCase or snake_case metadata key
// (mirrors how the connections table reads it).
export const getCurrentDeploymentMode = (ctx: WizardContext): string => {
  const metadata = (getConfiguredConnection(ctx).metadata as GenericRecord) || {};
  const mode = metadata.meshsyncDeploymentMode ?? metadata.meshsync_deployment_mode;
  return typeof mode === 'string' && mode ? mode : DEFAULT_MESHSYNC_DEPLOYMENT_MODE;
};

// The mode currently chosen in the step, falling back to the persisted one.
export const getSelectedDeploymentMode = (ctx: WizardContext): string => {
  const pending = ctx.data.postConfig.meshsyncDeploymentMode;
  return typeof pending === 'string' && pending ? pending : getCurrentDeploymentMode(ctx);
};

const DeploymentModeStepBody = ({ ctx }: { ctx: WizardContext }) => (
  <Box sx={{ display: 'grid', gap: 3 }}>
    <StepHeader
      title="MeshSync deployment mode"
      subtitle="Choose how MeshSync runs for this cluster. MeshSync keeps Meshery's view of the cluster's resources in sync; changing the mode redeploys it."
    />
    {Boolean(ctx.data.registrationError) && (
      <Alert severity="error" variant="filled">
        {formatWizardError(ctx.data.registrationError)}
      </Alert>
    )}
    <MeshsyncDeploymentModePicker
      value={getSelectedDeploymentMode(ctx)}
      currentValue={getCurrentDeploymentMode(ctx)}
      onChange={(mode) => ctx.patchPostConfig({ meshsyncDeploymentMode: mode })}
    />
  </Box>
);

// Post-config (configure mode only): choose how MeshSync runs for this
// kubernetes connection. Switching the mode makes the backend undeploy and
// redeploy MeshSync (operator vs embedded) for the connection's cluster.
export const kubernetesDeploymentModeStep: WizardStep = {
  id: 'kubernetes-meshsync-mode',
  label: 'MeshSync Mode',
  icon: SettingsIcon,
  Component: DeploymentModeStepBody,
  // Operates on a single already-registered connection; only meaningful when
  // (re)configuring an existing kubernetes connection.
  hidden: (ctx) => ctx.mode === 'create' || !getConfiguredConnection(ctx).id,
  nextLabel: (ctx) =>
    getSelectedDeploymentMode(ctx) === getCurrentDeploymentMode(ctx) ? 'Next' : 'Apply',
  onNext: async (ctx) => {
    const connection = getConfiguredConnection(ctx);
    const connectionId = String(connection.id ?? '');
    const selectedMode = getSelectedDeploymentMode(ctx);

    // Nothing to persist when the mode is unchanged.
    if (!connectionId || selectedMode === getCurrentDeploymentMode(ctx)) {
      return true;
    }

    ctx.patch({ registrationError: null });
    try {
      // Dedicated action endpoint: the server owns the metadata merge and the
      // MeshSync redeploy, keyed on the connection id.
      await ctx.services.setMeshsyncMode(connectionId, selectedMode as 'operator' | 'embedded');
      // Keep the local connection in sync so the step shows the new mode as
      // current if the user steps back into it.
      const metadata = (connection.metadata as GenericRecord) || {};
      const nextMetadata = { ...metadata, meshsync_deployment_mode: selectedMode };
      ctx.patch({ registrationResult: { ...connection, metadata: nextMetadata } });
      ctx.services.notify({
        message: `MeshSync deployment mode set to ${selectedMode}.`,
        event_type: EVENT_TYPES.SUCCESS,
      });
      return true;
    } catch (error) {
      ctx.patch({ registrationError: error });
      ctx.services.notify({
        message: `Failed to change MeshSync deployment mode: ${formatWizardError(error)}`,
        event_type: EVENT_TYPES.ERROR,
      });
      return false;
    }
  },
};
