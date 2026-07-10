import { useState, type ReactNode } from 'react';
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Divider,
  TextField,
  Typography,
  SettingsIcon,
} from '@sistent/sistent';
import { useSelector } from 'react-redux';
import { EVENT_TYPES } from 'lib/event-types';
import { CONNECTION_STATES, MESHSYNC_DEPLOYMENT_TYPE } from '@/utils/Enum';
import { useGetControllerDiagnosticsQuery } from '@/rtk-query/connection';
import { formatWizardError } from './errors';
import { StepHeader } from '../ConnectionWizardStepContent';
import { ConnectionStatusSelect } from '../ConnectionStatusSelect';
import type { ConnectionTransitionMap } from '../ConnectionTable.constants';
import {
  MeshsyncDeploymentModePicker,
  getConfiguredConnection,
  getCurrentDeploymentMode,
  getSelectedDeploymentMode,
} from './kubernetesDeploymentMode';
import type { GenericRecord, WizardContext, WizardStep } from './types';

const getMetadata = (ctx: WizardContext): GenericRecord =>
  (getConfiguredConnection(ctx).metadata as GenericRecord) || {};

// The connection's display name lives on metadata.name for kubernetes
// connections (that's what the table and detail view render).
const getCurrentName = (ctx: WizardContext): string =>
  String(getMetadata(ctx).name ?? getConfiguredConnection(ctx).name ?? '');

const getPendingName = (ctx: WizardContext): string => {
  const pending = ctx.data.postConfig.connectionName;
  return typeof pending === 'string' ? pending : getCurrentName(ctx);
};

const getStatus = (ctx: WizardContext): string =>
  String(getConfiguredConnection(ctx).status ?? getMetadata(ctx).status ?? '').toLowerCase();

const hasPendingChanges = (ctx: WizardContext): boolean =>
  getPendingName(ctx).trim() !== getCurrentName(ctx).trim() ||
  getSelectedDeploymentMode(ctx) !== getCurrentDeploymentMode(ctx);

// findBrokerNetworking pulls the informational broker_networking diagnostic (how
// Meshery reaches the broker) out of the diagnostics response, if present.
const findBrokerNetworking = (
  diagnostics?: { code?: string; summary?: string; endpoint?: string }[],
): { summary?: string; endpoint?: string } | undefined =>
  (diagnostics || []).find((d) => d.code === 'broker_networking');

// A small, muted, uppercase field/section label for a consistent, refined look.
const FieldLabel = ({ children }: { children: ReactNode }) => (
  <Typography
    variant="caption"
    sx={{
      color: 'text.secondary',
      textTransform: 'uppercase',
      letterSpacing: '0.06em',
      fontWeight: 600,
    }}
  >
    {children}
  </Typography>
);

const SettingsStepBody = ({ ctx }: { ctx: WizardContext }) => {
  const connection = getConfiguredConnection(ctx);
  const connectionId = String(connection.id ?? '');
  const status = getStatus(ctx);
  const [statusBusy, setStatusBusy] = useState(false);
  const [flushBusy, setFlushBusy] = useState(false);

  // The per-kind transition map that drives the status dropdown — the same one
  // the Connections table uses (populated on `state.ui.connectionMetadataState`).
  const transitionMap = useSelector(
    (state: {
      ui: {
        connectionMetadataState: Record<string, { transitionMap?: ConnectionTransitionMap }> | null;
      };
    }) =>
      state.ui.connectionMetadataState?.[String(connection.kind ?? 'kubernetes')]?.transitionMap,
  );

  // Surface how Meshery reaches the broker (managed port-forward / ClusterIP /
  // direct). Only meaningful in operator mode (embedded uses no in-cluster broker).
  const isOperatorMode = getCurrentDeploymentMode(ctx) === MESHSYNC_DEPLOYMENT_TYPE.OPERATOR;
  const { data: diagnosticsData } = useGetControllerDiagnosticsQuery(connectionId, {
    skip: !connectionId || !isOperatorMode,
  });
  const networking = findBrokerNetworking(diagnosticsData?.diagnostics);
  // The diagnostic summary is "Broker networking: <transport>"; strip the prefix
  // since the UI supplies its own label.
  const brokerTransport = (networking?.summary ?? '').replace(/^Broker networking:\s*/i, '');

  // Transition the connection to the picked lifecycle state (reuses the shared
  // status selector). Applied immediately, like the table.
  const changeStatus = async (nextStatus: string) => {
    if (!connectionId || nextStatus === status) {
      return;
    }
    setStatusBusy(true);
    try {
      await ctx.services.updateConnectionById(connectionId, { status: nextStatus });
      ctx.patch({ registrationResult: { ...connection, status: nextStatus } });
      ctx.services.notify({
        message: `Connection transitioned to ${nextStatus}.`,
        event_type: EVENT_TYPES.SUCCESS,
      });
    } catch (error) {
      ctx.services.notify({
        message: `Failed to update connection status: ${formatWizardError(error)}`,
        event_type: EVENT_TYPES.ERROR,
      });
    } finally {
      setStatusBusy(false);
    }
  };

  // Flush MeshSync data for this cluster and re-sync: clears the cached cluster
  // state (repopulated by a live MeshSync) via the connection actions endpoint,
  // keyed on the connection id (server resolves the cluster).
  const flushMeshSync = async () => {
    if (!connectionId) {
      return;
    }
    setFlushBusy(true);
    try {
      await ctx.services.flushMeshsync(connectionId);
      ctx.services.notify({
        message: 'MeshSync data flush requested; data will be repopulated from the cluster.',
        event_type: EVENT_TYPES.SUCCESS,
      });
    } catch (error) {
      ctx.services.notify({
        message: `Failed to flush MeshSync data: ${formatWizardError(error)}`,
        event_type: EVENT_TYPES.ERROR,
      });
    } finally {
      setFlushBusy(false);
    }
  };

  return (
    <Box sx={{ display: 'grid', gap: 2.5 }}>
      <StepHeader
        title="Settings"
        subtitle="Rename the connection, change its lifecycle status, choose how MeshSync runs, and flush MeshSync data for this cluster."
      />
      {Boolean(ctx.data.registrationError) && (
        <Alert severity="error" variant="filled">
          {formatWizardError(ctx.data.registrationError)}
        </Alert>
      )}

      {/* Name + Status on one row — both are compact, single-value fields. */}
      <Box sx={{ display: 'flex', gap: 2.5, flexWrap: 'wrap', alignItems: 'flex-start' }}>
        <Box sx={{ display: 'grid', gap: 0.75, flex: 1, minWidth: 220 }}>
          <FieldLabel>Name</FieldLabel>
          <TextField
            size="small"
            fullWidth
            value={getPendingName(ctx)}
            onChange={(event) => ctx.patchPostConfig({ connectionName: event.target.value })}
            inputProps={{ 'aria-label': 'Connection name' }}
          />
        </Box>
        <Box sx={{ display: 'grid', gap: 0.75 }}>
          <FieldLabel>Status</FieldLabel>
          <ConnectionStatusSelect
            status={status || CONNECTION_STATES.DISCOVERED}
            transitionMap={transitionMap}
            disabled={statusBusy}
            onChange={changeStatus}
          />
        </Box>
      </Box>

      <Divider />

      {/* MeshSync deployment mode = operator stack lifecycle */}
      <Box sx={{ display: 'grid', gap: 1.5 }}>
        <Box sx={{ display: 'grid', gap: 0.25 }}>
          <FieldLabel>MeshSync &amp; Operator</FieldLabel>
          <Typography variant="body2" color="text.secondary">
            Choose how MeshSync runs. Switching modes deploys or uninstalls the in-cluster Operator,
            Broker and MeshSync.
          </Typography>
        </Box>
        <MeshsyncDeploymentModePicker
          value={getSelectedDeploymentMode(ctx)}
          currentValue={getCurrentDeploymentMode(ctx)}
          onChange={(mode) => ctx.patchPostConfig({ meshsyncDeploymentMode: mode })}
        />
        {/* Read-only: how Meshery currently reaches the broker (operator mode only). */}
        {isOperatorMode && (
          <Box sx={{ display: 'flex', alignItems: 'baseline', gap: 1, flexWrap: 'wrap' }}>
            <FieldLabel>Broker networking</FieldLabel>
            <Typography variant="body2" color="text.secondary">
              {brokerTransport || 'Resolving how Meshery reaches the broker…'}
            </Typography>
            {networking?.endpoint && (
              <Box
                component="code"
                sx={{ fontFamily: 'monospace', fontSize: '0.8rem', color: 'text.secondary' }}
              >
                {networking.endpoint}
              </Box>
            )}
          </Box>
        )}
      </Box>

      <Divider />

      {/* MeshSync data maintenance. */}
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          gap: 2,
          flexWrap: 'wrap',
        }}
      >
        <Box sx={{ display: 'grid', gap: 0.25 }}>
          <FieldLabel>MeshSync data</FieldLabel>
          <Typography variant="body2" color="text.secondary">
            Clear this cluster&apos;s cached state; it is repopulated by a running MeshSync.
          </Typography>
        </Box>
        <Button variant="outlined" onClick={flushMeshSync} disabled={flushBusy}>
          {flushBusy ? <CircularProgress size={16} /> : 'Flush MeshSync'}
        </Button>
      </Box>
    </Box>
  );
};

/**
 * Post-config (configure mode only) settings for a single kubernetes connection:
 * rename, the lifecycle status transition dropdown (shared with the table), the
 * MeshSync deployment mode (which also installs/uninstalls the in-cluster
 * Operator + Broker + MeshSync stack), and a Flush MeshSync action. Name and mode
 * changes are applied on "Apply"; status changes and the flush act immediately.
 */
export const kubernetesSettingsStep: WizardStep = {
  id: 'kubernetes-settings',
  label: 'Settings',
  icon: SettingsIcon,
  Component: SettingsStepBody,
  hidden: (ctx) => ctx.mode === 'create' || !getConfiguredConnection(ctx).id,
  nextLabel: (ctx) => (hasPendingChanges(ctx) ? 'Apply' : 'Next'),
  onNext: async (ctx) => {
    const connection = getConfiguredConnection(ctx);
    const connectionId = String(connection.id ?? '');
    if (!connectionId || !hasPendingChanges(ctx)) {
      return true;
    }

    const newName = getPendingName(ctx).trim();
    const nameChanged = Boolean(newName) && newName !== getCurrentName(ctx).trim();
    const selectedMode = getSelectedDeploymentMode(ctx);
    const modeChanged = selectedMode !== getCurrentDeploymentMode(ctx);

    ctx.patch({ registrationError: null });
    try {
      let nextConnection = connection;

      if (nameChanged) {
        const metadata = { ...getMetadata(ctx), name: newName };
        // Send the current status alongside so the rename can't clear it.
        await ctx.services.updateConnectionById(connectionId, {
          status: connection.status,
          metadata,
        });
        nextConnection = { ...nextConnection, name: newName, metadata };
      }

      if (modeChanged) {
        // Dedicated action endpoint: the server owns the metadata merge and the
        // MeshSync (and operator stack) redeploy, keyed on the connection id.
        await ctx.services.setMeshsyncMode(connectionId, selectedMode as 'operator' | 'embedded');
        const metadata = {
          ...((nextConnection.metadata as GenericRecord) || {}),
          meshsync_deployment_mode: selectedMode,
        };
        nextConnection = { ...nextConnection, metadata };
      }

      ctx.patch({ registrationResult: nextConnection });
      ctx.services.notify({ message: 'Settings applied.', event_type: EVENT_TYPES.SUCCESS });
      return true;
    } catch (error) {
      ctx.patch({ registrationError: error });
      ctx.services.notify({
        message: `Failed to apply settings: ${formatWizardError(error)}`,
        event_type: EVENT_TYPES.ERROR,
      });
      return false;
    }
  },
};
