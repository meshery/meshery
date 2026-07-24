import React, { useEffect, useState } from 'react';
import { Typography } from '@sistent/sistent';
import { Modal } from '@/components/shared/Modal';
import { ModalButtonPrimary, ModalButtonSecondary } from '@sistent/sistent';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import {
  useGetConnectionControllersConfigQuery,
  useUpdateConnectionControllersConfigMutation,
  type ControllersConfigDoc,
} from '@/rtk-query/controllersConfig';
import ControllersConfigForm, {
  BUILT_IN_CONTROLLERS_CONFIG,
} from '@/components/configuration/ControllersConfigForm';

const stripSchemaVersion = (doc?: ControllersConfigDoc | null): ControllersConfigDoc => {
  if (!doc) return {};
  const rest = { ...doc };
  delete rest.schemaVersion;
  return rest;
};

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
  const { notify } = useNotification();
  const { data, isLoading } = useGetConnectionControllersConfigQuery(
    { connectionId },
    { skip: !isOpen || !connectionId },
  );
  const [updateOverride, { isLoading: isSaving }] = useUpdateConnectionControllersConfigMutation();
  const [draft, setDraft] = useState<ControllersConfigDoc>({});
  const [dirty, setDirty] = useState(false);

  useEffect(() => {
    if (data && !dirty) {
      setDraft(stripSchemaVersion(data.override));
    }
  }, [data, dirty]);

  const handleSave = async () => {
    try {
      await updateOverride({ connectionId, body: draft }).unwrap();
      setDirty(false);
      notify({
        message: `Controllers configuration applied to ${connectionName || 'connection'}.`,
        event_type: EVENT_TYPES.SUCCESS,
      });
      onClose();
    } catch (err) {
      notify({
        message: 'Failed to apply the controllers configuration override.',
        event_type: EVENT_TYPES.ERROR,
        details: String((err as { data?: unknown })?.data ?? err),
      });
    }
  };

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
          <ModalButtonPrimary
            onClick={handleSave}
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
        onChange={(next) => {
          setDraft(next);
          setDirty(true);
        }}
        inheritedLayers={[data?.default, BUILT_IN_CONTROLLERS_CONFIG]}
        inheritLabel="Server default"
        showSourceIndicators
        disabled={isLoading || isSaving}
      />
    </Modal>
  );
}
