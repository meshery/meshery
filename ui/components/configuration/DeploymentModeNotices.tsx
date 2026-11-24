// The parts of the controllers editor that state what the effective deployment
// mode does to the document: which mode governs, where it came from, and why a
// section's settings are inert (or, on the server-wide layer, conditional).
//
// Kept out of ControllersConfigForm so the form stays a field renderer and this
// stays the one place the mode is explained to the user. What is said here is
// derived from ./deploymentMode, which mirrors the server's dependency
// structure - it is not a second opinion about it.

import React from 'react';
import { Alert, AlertTitle, Box, Button, Chip, Typography } from '@sistent/sistent';
import type { ControllersConfigDoc } from './ControllersConfigForm';
import {
  DEPLOYMENT_MODE_LABEL,
  EMBEDDED_INERT_REASON,
  SECTION_PATHS,
  SERVER_DEFAULT_SCOPE_NOTE,
  dormantPathsIn,
  isInertIn,
  type ConfigSection,
  type DeploymentModeGovernance,
} from './deploymentMode';

/**
 * What mode governs this editor and which layer it came from. The user has to
 * be able to read that before reading anything else on the form: it is what
 * decides whether the rest of the document reaches a cluster at all.
 */
export const DeploymentModeBanner: React.FC<{ governance?: DeploymentModeGovernance }> = ({
  governance,
}) => {
  if (!governance) return null;
  const { mode, sourceLabel, scope, unsaved } = governance;
  const isConnection = scope === 'connection';

  const consequence = isConnection
    ? mode === 'embedded'
      ? 'MeshSync runs inside Meshery Server for this connection and nothing is installed into the cluster, so the settings marked "Not applied" below cannot take effect until this connection uses Operator mode.'
      : 'Meshery Operator manages MeshSync and Meshery Broker on this cluster, so every setting below applies.'
    : mode === 'embedded'
      ? 'Connections that do not override the mode run MeshSync inside Meshery Server, where the in-cluster MeshSync and Meshery Broker settings below do not apply. They still reach connections that override the mode to Operator.'
      : 'Connections that do not override the mode run Meshery Operator, which applies every setting below.';

  return (
    <Alert
      severity={mode === 'embedded' ? 'info' : 'success'}
      variant="outlined"
      sx={{ marginBottom: '1.5rem' }}
      data-testid="controllers-config-mode-banner"
    >
      <AlertTitle>
        {isConnection ? 'Effective deployment mode' : 'Default deployment mode'}:{' '}
        {DEPLOYMENT_MODE_LABEL[mode]}
        {unsaved ? ' (unsaved)' : ''}
      </AlertTitle>
      {isConnection ? `Resolved from ${sourceLabel}. ` : `Set by ${sourceLabel}. `}
      {consequence}
      {unsaved ? ' This mode is applied when you save.' : ''}
    </Alert>
  );
};

/**
 * A section heading, carrying the mode marker when every setting in the section
 * is inert. Each such setting is chipped individually too; the heading chip is
 * what makes it legible while scanning.
 */
export const SectionHeading: React.FC<{
  title: string;
  section: ConfigSection;
  governance?: DeploymentModeGovernance;
}> = ({ title, section, governance }) => {
  const fullyInert = SECTION_PATHS[section].every((path) => isInertIn(governance, path));
  return (
    <Box
      sx={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '0.5rem' }}
      data-testid={`controllers-config-section-${section}`}
    >
      <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
        {title}
      </Typography>
      {fullyInert ? (
        <Chip size="small" label="Not applied in Embedded mode" variant="outlined" />
      ) : null}
    </Box>
  );
};

/**
 * The per-section statement of what the mode does to the section: why its
 * settings are inert (per-connection) or which connections they reach
 * (server-wide defaults), plus a way to clear values left dormant by the mode.
 * Stated in the form itself, not in a tooltip.
 */
export const SectionNotice: React.FC<{
  section: ConfigSection;
  governance?: DeploymentModeGovernance;
  value: ControllersConfigDoc;
  onClearDormant: (section: ConfigSection) => void;
  disabled?: boolean;
}> = ({ section, governance, value, onClearDormant, disabled = false }) => {
  if (!governance || governance.mode === 'operator') return null;

  if (governance.scope === 'serverDefault') {
    const note = SERVER_DEFAULT_SCOPE_NOTE[section];
    if (!note) return null;
    return (
      <Alert severity="info" variant="outlined" sx={{ marginBottom: '1rem' }}>
        {note}
      </Alert>
    );
  }

  const reason = EMBEDDED_INERT_REASON[section];
  if (!reason) return null;
  const dormant = dormantPathsIn(governance, value, section);

  return (
    <Alert
      severity="info"
      variant="outlined"
      sx={{ marginBottom: '1rem' }}
      data-testid={`controllers-config-inert-${section}`}
      action={
        dormant.length > 0 ? (
          <Button size="small" disabled={disabled} onClick={() => onClearDormant(section)}>
            Clear {dormant.length} dormant {dormant.length === 1 ? 'value' : 'values'}
          </Button>
        ) : undefined
      }
    >
      <AlertTitle>Not applied in Embedded mode</AlertTitle>
      {reason}
      {dormant.length > 0
        ? ' Values already stored here are kept and shown below; they become live again if this connection uses Operator mode.'
        : ''}
    </Alert>
  );
};
