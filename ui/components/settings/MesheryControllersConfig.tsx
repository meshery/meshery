import React, { useEffect, useState } from 'react';
import { Box, Button, Paper, Typography } from '@sistent/sistent';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import {
  useGetControllersDefaultConfigQuery,
  useUpdateControllersDefaultConfigMutation,
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

/**
 * Settings tab: server-wide defaults for the Meshery Operator, MeshSync, and
 * Meshery Broker deployed to every managed Kubernetes cluster. Fields left on
 * Inherit fall back to the controllers' built-in defaults; per-connection
 * overrides (Connections page) take precedence over everything set here.
 */
export default function MesheryControllersConfig() {
  const { notify } = useNotification();
  const { data, isLoading, error } = useGetControllersDefaultConfigQuery();
  const [updateDefaults, { isLoading: isSaving }] = useUpdateControllersDefaultConfigMutation();
  const [draft, setDraft] = useState<ControllersConfigDoc>({});
  const [dirty, setDirty] = useState(false);

  useEffect(() => {
    if (data && !dirty) {
      setDraft(stripSchemaVersion(data));
    }
  }, [data, dirty]);

  const handleSave = async () => {
    try {
      await updateDefaults({ body: draft }).unwrap();
      setDirty(false);
      notify({
        message:
          'Server-wide controllers configuration defaults saved. Re-applying to connected clusters.',
        event_type: EVENT_TYPES.SUCCESS,
      });
    } catch (err) {
      notify({
        message: 'Failed to save controllers configuration defaults.',
        event_type: EVENT_TYPES.ERROR,
        details: String((err as { data?: unknown })?.data ?? err),
      });
    }
  };

  const handleDiscard = () => {
    setDraft(stripSchemaVersion(data));
    setDirty(false);
  };

  // Notify once per load failure rather than on every render.
  useEffect(() => {
    if (error) {
      notify({
        message: 'Failed to load controllers configuration defaults.',
        event_type: EVENT_TYPES.ERROR,
      });
    }
  }, [error]);

  return (
    <Paper sx={{ padding: '1.5rem', marginTop: '1rem' }}>
      <Typography variant="h6" sx={{ fontWeight: 600 }}>
        Meshery Operator, MeshSync &amp; Broker
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ marginBottom: '1.5rem' }}>
        Server-wide defaults applied to every managed Kubernetes cluster. Individual connections can
        override any of these on the Connections page; fields left on Inherit use the
        controllers&apos; built-in defaults.
      </Typography>

      <ControllersConfigForm
        value={draft}
        onChange={(next) => {
          setDraft(next);
          setDirty(true);
        }}
        inheritedLayers={[BUILT_IN_CONTROLLERS_CONFIG]}
        inheritLabel="Built-in default"
        disabled={isLoading || isSaving}
      />

      <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1.5rem' }}>
        <Button variant="outlined" onClick={handleDiscard} disabled={!dirty || isSaving}>
          Discard changes
        </Button>
        <Button
          variant="contained"
          onClick={handleSave}
          disabled={!dirty || isSaving}
          data-testid="controllers-config-save"
        >
          Save defaults
        </Button>
      </Box>
    </Paper>
  );
}
