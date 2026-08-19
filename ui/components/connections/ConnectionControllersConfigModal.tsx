import React from 'react';
import { Typography } from '@sistent/sistent';
import { Modal } from '@/components/shared/Modal';
import { ModalButtonPrimary, ModalButtonSecondary } from '@sistent/sistent';
import {
  useGetConnectionControllersConfigQuery,
  useUpdateConnectionControllersConfigMutation,
} from '@meshery/schemas/mesheryApi';
import ControllersConfigForm, {
  BUILT_IN_CONTROLLERS_CONFIG,
} from '@/components/configuration/ControllersConfigForm';
import { useControllersConfigDraft } from '@/components/configuration/useControllersConfigDraft';
import { connectionDeploymentMode } from '@/components/configuration/deploymentMode';

type ConnectionControllersConfigModalProps = {
  isOpen: boolean;
  onClose: () => void;
  connectionId: string;
  connectionName?: string;
};

/**
 * Per-connection override editor for the Meshery Operator, MeshSync, and
 * Broker configuration. Shows, per field, whether the effective value is an
 * override, the server-wide default, or the built-in default; fields left on
 * Inherit follow the server-wide defaults from Settings.
 */
export default function ConnectionControllersConfigModal({
  isOpen,
  onClose,
  connectionId,
  connectionName,
}: ConnectionControllersConfigModalProps) {
  const { data, isLoading, error } = useGetConnectionControllersConfigQuery(
    { connectionId },
    { skip: !isOpen || !connectionId },
  );
  const [updateOverride, { isLoading: isSaving }] = useUpdateConnectionControllersConfigMutation();
  const { draft, dirty, onChange, discard, save } = useControllersConfigDraft({
    isLoaded: Boolean(data),
    source: data?.override,
    loadError: error,
    save: (body) => updateOverride({ connectionId, body }).unwrap(),
    messages: {
      loadError: 'Failed to load the controllers configuration for this connection.',
      saveError: 'Failed to apply the controllers configuration override.',
      saveSuccess: `Controllers configuration applied to ${connectionName || 'connection'}.`,
    },
    onSaved: onClose,
  });

  // Meshery Operator manages MeshSync and Meshery Broker, so the mode this
  // connection resolves to decides which settings can reach the cluster at all.
  // It follows the draft, not just the persisted state: switching the mode
  // select brings the dependent sections to life (or puts them to sleep)
  // before the save rather than after it.
  const deploymentMode = connectionDeploymentMode({
    draft,
    persistedOverride: data?.override,
    serverDefault: data?.default,
    serverEffective: data?.effective,
  });

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={`Operator, MeshSync & Broker Configuration${
        connectionName ? ` - ${connectionName}` : ''
      }`}
      size="lg"
      helpText="Absent fields inherit the server-wide defaults from Settings."
      actions={
        <>
          <ModalButtonSecondary onClick={onClose} disabled={isSaving}>
            Cancel
          </ModalButtonSecondary>
          <ModalButtonSecondary onClick={discard} disabled={!dirty || isSaving}>
            Discard changes
          </ModalButtonSecondary>
          <ModalButtonPrimary
            onClick={save}
            disabled={!dirty || isSaving}
            data-testid="connection-controllers-config-save"
          >
            Save &amp; Apply
          </ModalButtonPrimary>
        </>
      }
    >
      <Typography variant="body2" color="text.secondary" sx={{ marginBottom: '1.5rem' }}>
        Overrides apply to this connection only. Fields left on Inherit follow the server-wide
        defaults from Settings, then the controllers&apos; built-in defaults. Changes are applied to
        the cluster immediately; restart-required changes restart MeshSync automatically.
      </Typography>
      <ControllersConfigForm
        value={draft}
        onChange={onChange}
        inheritedLayers={[data?.default, BUILT_IN_CONTROLLERS_CONFIG]}
        inheritLabel="Server default"
        showSourceIndicators
        deploymentMode={deploymentMode}
        disabled={isLoading || isSaving}
      />
    </Modal>
  );
}
