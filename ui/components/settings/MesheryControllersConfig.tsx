import React from 'react';
import { Box, Button, Link, Paper, Typography } from '@sistent/sistent';
import NextLink from 'next/link';
import {
  useGetControllersDefaultConfigQuery,
  useUpdateControllersDefaultConfigMutation,
} from '@meshery/schemas/mesheryApi';
import ControllersConfigForm, {
  BUILT_IN_CONTROLLERS_CONFIG,
} from '@/components/configuration/ControllersConfigForm';
import { useControllersConfigDraft } from '@/components/configuration/useControllersConfigDraft';
import { serverDefaultDeploymentMode } from '@/components/configuration/deploymentMode';

/**
 * Settings tab: server-wide defaults for the Meshery Operator, MeshSync, and
 * Meshery Broker deployed to every managed Kubernetes cluster. Fields left on
 * Inherit fall back to the controllers' built-in defaults; per-connection
 * overrides (Connections page) take precedence over everything set here.
 */
export default function MesheryControllersConfig() {
  const { data, isLoading, error } = useGetControllersDefaultConfigQuery();
  const [updateDefaults, { isLoading: isSaving }] = useUpdateControllersDefaultConfigMutation();
  const { draft, dirty, onChange, discard, save } = useControllersConfigDraft({
    isLoaded: Boolean(data),
    source: data,
    loadError: error,
    save: (body) => updateDefaults({ body }).unwrap(),
    messages: {
      loadError: 'Failed to load controllers configuration defaults.',
      saveError: 'Failed to save controllers configuration defaults.',
      saveSuccess:
        'Server-wide controllers configuration defaults saved. Re-applying to connected clusters.',
    },
  });

  return (
    <Paper sx={{ padding: '1.5rem', marginTop: '1rem' }}>
      <Typography variant="h6" sx={{ fontWeight: 600 }}>
        Meshery Operator, MeshSync &amp; Broker
      </Typography>
      <Typography variant="body2" color="text.secondary">
        Server-wide defaults for every managed Kubernetes cluster. Fields left on Inherit use
        built-in defaults.
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ marginBottom: '1rem' }}>
        For per-cluster overrides, go to{' '}
        <Link
          className="keppel"
          component={NextLink}
          href="/management/connections"
          underline="hover"
        >
          Connection
        </Link>
        {' → Actions → Configure Controllers.'}
      </Typography>

      <ControllersConfigForm
        value={draft}
        onChange={onChange}
        inheritedLayers={[BUILT_IN_CONTROLLERS_CONFIG]}
        inheritLabel="Built-in default"
        deploymentMode={serverDefaultDeploymentMode(draft)}
        disabled={isLoading || isSaving}
      />

      <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1.5rem' }}>
        <Button variant="outlined" onClick={discard} disabled={!dirty || isSaving}>
          Discard changes
        </Button>
        <Button
          variant="contained"
          onClick={save}
          disabled={!dirty || isSaving}
          data-testid="controllers-config-save"
        >
          Save defaults
        </Button>
      </Box>
    </Paper>
  );
}
