import React, { useEffect, useState } from 'react';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Alert,
  Box,
  Chip,
  Divider,
  ExpandMoreIcon,
  InfoTooltip,
  MenuItem,
  TextField,
  Typography,
} from '@sistent/sistent';
import type { UpdateControllersDefaultConfigApiArg } from '@meshery/schemas/mesheryApi';
import { alpha } from '@/theme';
import { getPath, setPath, type FieldPath } from './fieldPath';
import {
  dormantPathsIn,
  isInertIn,
  type ConfigSection,
  type DeploymentMode,
  type DeploymentModeGovernance,
} from './deploymentMode';
import { DeploymentModeBanner, SectionHeading, SectionNotice } from './DeploymentModeNotices';
import ControllersConfigModePicker from './ControllersConfigModePicker';
import ControllersConfigWatchList from './ControllersConfigWatchList';
import OperatorVersionField from './OperatorVersionField';
import { INHERIT, fitWidth, type WatchList } from './controllersConfigForm.shared';

/**
 * The editable controllers configuration document: the generated PUT request
 * body, i.e. the wire document without the server-stamped `schemaVersion`.
 * Both layers (server-wide defaults and the per-connection override) accept
 * this exact shape.
 */
export type ControllersConfigDoc = UpdateControllersDefaultConfigApiArg['body'];

/** As returned by the server, which stamps the schema revision onto the document. */
export type VersionedControllersConfigDoc = ControllersConfigDoc & { schemaVersion?: string };

// Built-in defaults mirrored from Meshery Server (BuiltInControllersConfig):
// what applies when no layer sets a field.
export const BUILT_IN_CONTROLLERS_CONFIG: ControllersConfigDoc = {
  operator: { deploymentMode: 'embedded' },
  meshsync: { replicas: 1, redactSecrets: false, brokerContentDedup: false, debugLogging: false },
  broker: { replicas: 1, service: { type: 'ClusterIP' } },
};

type SourceInfo = { label: string; overridden: boolean };

export type ControllersConfigFormProps = {
  value: ControllersConfigDoc;
  onChange: (next: ControllersConfigDoc) => void;
  inheritedLayers?: (ControllersConfigDoc | null | undefined)[];
  inheritLabel?: string;
  showSourceIndicators?: boolean;
  deploymentMode?: DeploymentModeGovernance;
  disabled?: boolean;
};

const SubsectionTitle = ({
  title,
  helpText,
  chip,
}: {
  title: string;
  helpText?: string;
  chip?: string;
}) => (
  <Box
    sx={{
      display: 'flex',
      alignItems: 'center',
      gap: '0.5rem',
      flexWrap: 'wrap',
      marginBottom: '0.75rem',
    }}
  >
    <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
      {title}
    </Typography>
    {helpText ? <InfoTooltip helpText={helpText} placement="top" /> : null}
    {chip ? <Chip size="small" label={chip} variant="outlined" /> : null}
  </Box>
);

const FieldRow = ({ children }: { children: React.ReactNode }) => (
  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2, alignItems: 'flex-start' }}>{children}</Box>
);

/**
 * Layered editor for the Meshery Operator, MeshSync, and Broker configuration.
 * Every control is tri-state: Inherit (or empty) keeps the field absent so the
 * next layer applies. Layout is mode-first.
 */
export default function ControllersConfigForm({
  value,
  onChange,
  inheritedLayers = [BUILT_IN_CONTROLLERS_CONFIG],
  inheritLabel = 'Inherited',
  showSourceIndicators = false,
  deploymentMode,
  disabled = false,
}: ControllersConfigFormProps) {
  const modeForDisclosure: DeploymentMode =
    (getPath(value, ['operator', 'deploymentMode']) as DeploymentMode | undefined) ??
    deploymentMode?.mode ??
    'embedded';
  const operatorModeApplies = modeForDisclosure === 'operator';
  const [meshsyncDeployOpen, setMeshsyncDeployOpen] = useState(operatorModeApplies);
  const [brokerOpen, setBrokerOpen] = useState(operatorModeApplies);

  useEffect(() => {
    setMeshsyncDeployOpen(operatorModeApplies);
    setBrokerOpen(operatorModeApplies);
  }, [operatorModeApplies]);

  const inheritedValue = (path: FieldPath): unknown => {
    for (const layer of inheritedLayers) {
      const v = getPath(layer ?? undefined, path);
      if (v !== undefined) return v;
    }
    return undefined;
  };

  const sourceOf = (path: FieldPath): SourceInfo => {
    if (getPath(value, path) !== undefined) return { label: 'Override', overridden: true };
    if (getPath(inheritedLayers[0] ?? undefined, path) !== undefined)
      return { label: inheritLabel, overridden: false };
    return { label: 'Built-in default', overridden: false };
  };

  const sourceChip = (path: FieldPath) => {
    if (!showSourceIndicators) return null;
    const source = sourceOf(path);
    return (
      <Chip
        size="small"
        label={source.label}
        color={source.overridden ? 'primary' : 'default'}
        variant={source.overridden ? 'filled' : 'outlined'}
        sx={{ marginLeft: '0.5rem', height: '20px' }}
      />
    );
  };

  const isInert = (path: FieldPath): boolean => isInertIn(deploymentMode, path);
  const isDisabled = (path: FieldPath): boolean => disabled || isInert(path);

  const clearDormant = (section: ConfigSection) => {
    onChange(
      dormantPathsIn(deploymentMode, value, section).reduce(
        (doc, path) => setPath(doc, path, undefined),
        value,
      ),
    );
  };

  const fieldLabel = (text: string, path: FieldPath, helper?: string) => (
    <Box
      sx={{
        display: 'flex',
        alignItems: 'center',
        marginBottom: '0.25rem',
        gap: '0.25rem',
        flexWrap: 'wrap',
      }}
    >
      <Typography variant="body2" sx={{ fontWeight: 500 }}>
        {text}
      </Typography>
      {helper ? <InfoTooltip helpText={helper} placement="top" /> : null}
      {sourceChip(path)}
      {isInert(path) ? (
        <Chip
          size="small"
          label="Not applied"
          variant="outlined"
          sx={{ marginLeft: '0.25rem', height: '20px' }}
        />
      ) : null}
    </Box>
  );

  const notice = (section: ConfigSection) => (
    <SectionNotice
      section={section}
      governance={deploymentMode}
      value={value}
      onClearDormant={clearDormant}
      disabled={disabled}
    />
  );

  const triStateBoolean = (label: string, path: FieldPath, helper?: string) => {
    const current = getPath(value, path) as boolean | undefined;
    const inherited = inheritedValue(path) as boolean | undefined;
    const inheritOption = `Inherit (${inherited === undefined ? 'unset' : inherited ? 'Enabled' : 'Disabled'})`;
    return (
      <Box>
        {fieldLabel(label, path, helper)}
        <TextField
          select
          size="small"
          disabled={isDisabled(path)}
          value={current === undefined ? INHERIT : current ? 'true' : 'false'}
          sx={{ width: fitWidth(inheritOption, 'Enabled', 'Disabled') }}
          onChange={(e) => {
            const v = e.target.value;
            onChange(setPath(value, path, v === INHERIT ? undefined : v === 'true'));
          }}
        >
          <MenuItem value={INHERIT}>{inheritOption}</MenuItem>
          <MenuItem value="true">Enabled</MenuItem>
          <MenuItem value="false">Disabled</MenuItem>
        </TextField>
      </Box>
    );
  };

  const textInput = (
    label: string,
    path: FieldPath,
    helper?: string,
    opts?: { number?: boolean; min?: number; max?: number },
  ) => {
    const current = getPath(value, path) as string | number | undefined;
    const inherited = inheritedValue(path);
    const placeholder = inherited !== undefined ? `Inherit (${inherited})` : 'Inherit';
    return (
      <Box>
        {fieldLabel(label, path, helper)}
        <TextField
          size="small"
          type={opts?.number ? 'number' : 'text'}
          disabled={isDisabled(path)}
          value={current ?? ''}
          placeholder={placeholder}
          sx={{ width: fitWidth(current, placeholder) }}
          slotProps={opts?.number ? { htmlInput: { min: opts.min, max: opts.max } } : undefined}
          onChange={(e) => {
            const raw = e.target.value;
            if (raw === '') {
              onChange(setPath(value, path, undefined));
              return;
            }
            onChange(setPath(value, path, opts?.number ? Number(raw) : raw));
          }}
        />
      </Box>
    );
  };

  const listInput = (label: string, path: FieldPath, helper: string) => {
    const current = getPath(value, path) as string[] | undefined;
    const inherited = inheritedValue(path) as string[] | undefined;
    const joined = current?.join(', ') ?? '';
    const placeholder =
      inherited && inherited.length > 0 ? `Inherit (${inherited.join(', ')})` : 'Inherit (all)';
    return (
      <Box>
        {fieldLabel(label, path, helper)}
        <TextField
          size="small"
          disabled={isDisabled(path)}
          value={joined}
          placeholder={placeholder}
          sx={{ width: fitWidth(joined, placeholder) }}
          onChange={(e) => {
            const raw = e.target.value;
            if (raw.trim() === '') {
              onChange(setPath(value, path, undefined));
              return;
            }
            onChange(
              setPath(
                value,
                path,
                raw
                  .split(',')
                  .map((s) => s.trim())
                  .filter(Boolean),
              ),
            );
          }}
        />
      </Box>
    );
  };

  const enumSelect = (
    label: string,
    path: FieldPath,
    options: { value: string; label: string }[],
    helper?: string,
    postProcess?: (
      next: ControllersConfigDoc,
      selected: string | undefined,
    ) => ControllersConfigDoc,
  ) => {
    const current = getPath(value, path) as string | undefined;
    const inherited = inheritedValue(path) as string | undefined;
    const inheritOption = `Inherit (${inherited ?? 'unset'})`;
    return (
      <Box>
        {fieldLabel(label, path, helper)}
        <TextField
          select
          size="small"
          disabled={isDisabled(path)}
          value={current ?? INHERIT}
          sx={{ width: fitWidth(inheritOption, ...options.map((option) => option.label)) }}
          onChange={(e) => {
            const selected = e.target.value === INHERIT ? undefined : e.target.value;
            let next = setPath(value, path, selected);
            if (postProcess) next = postProcess(next, selected);
            onChange(next);
          }}
        >
          <MenuItem value={INHERIT}>{inheritOption}</MenuItem>
          {options.map((option) => (
            <MenuItem key={option.value} value={option.value}>
              {option.label}
            </MenuItem>
          ))}
        </TextField>
      </Box>
    );
  };

  const clearLoadBalancerFieldsUnlessLB = (
    next: ControllersConfigDoc,
    selected: string | undefined,
  ): ControllersConfigDoc => {
    const effectiveType =
      selected ?? (inheritedValue(['broker', 'service', 'type']) as string | undefined);
    if (effectiveType === 'LoadBalancer') return next;
    return setPath(
      setPath(next, ['broker', 'service', 'loadBalancerClass'], undefined),
      ['broker', 'service', 'loadBalancerSourceRanges'],
      undefined,
    );
  };

  const redactSecrets =
    getPath(value, ['meshsync', 'redactSecrets']) ?? inheritedValue(['meshsync', 'redactSecrets']);

  const annotations = getPath(value, ['broker', 'service', 'annotations']) as
    | Record<string, string>
    | undefined;
  const annotationsText = annotations
    ? Object.entries(annotations)
        .map(([k, v]) => `${k}=${v}`)
        .join('\n')
    : '';
  const setAnnotationsFromText = (raw: string) => {
    if (raw.trim() === '') {
      onChange(setPath(value, ['broker', 'service', 'annotations'], undefined));
      return;
    }
    const next: Record<string, string> = {};
    raw.split('\n').forEach((line) => {
      const idx = line.indexOf('=');
      if (idx > 0) next[line.slice(0, idx).trim()] = line.slice(idx + 1).trim();
    });
    onChange(setPath(value, ['broker', 'service', 'annotations'], next));
  };

  const serviceType = (getPath(value, ['broker', 'service', 'type']) ??
    inheritedValue(['broker', 'service', 'type'])) as string | undefined;
  const isLoadBalancer = serviceType === 'LoadBalancer';

  const modePath: FieldPath = ['operator', 'deploymentMode'];
  const modeInherited = inheritedValue(modePath) as string | undefined;
  const inheritModeLabel =
    modeInherited === 'operator'
      ? 'Operator'
      : modeInherited === 'embedded'
        ? 'Embedded'
        : (modeInherited ?? 'embedded');

  return (
    <Box>
      <DeploymentModeBanner governance={deploymentMode} />

      <SectionHeading
        title="Meshery Operator"
        section="operator"
        governance={deploymentMode}
        id="controllers-config-operator"
      />
      {notice('operator')}
      <ControllersConfigModePicker
        label={fieldLabel(
          'Deployment mode',
          modePath,
          'How Meshery discovers the cluster. Operator installs MeshSync and Broker in-cluster; Embedded runs MeshSync inside Meshery Server. Hover a card for details. Changing the mode redeploys controllers.',
        )}
        selected={(getPath(value, modePath) as string | undefined) ?? INHERIT}
        inheritModeLabel={inheritModeLabel}
        inheritedMode={modeInherited === 'operator' ? 'operator' : 'embedded'}
        scope={deploymentMode?.scope ?? 'serverDefault'}
        disabled={disabled || isInert(modePath)}
        onChange={(selected) => onChange(setPath(value, modePath, selected))}
      />
      <Box>
        {fieldLabel(
          'Operator version',
          ['operator', 'version'],
          'Helm chart version (operator mode). Applying upgrades the Meshery Operator release. Inherit tracks the Meshery Server release.',
        )}
        <OperatorVersionField
          value={getPath(value, ['operator', 'version']) as string | undefined}
          placeholder={
            inheritedValue(['operator', 'version']) !== undefined
              ? `Inherit (${inheritedValue(['operator', 'version'])})`
              : 'Inherit (server release)'
          }
          disabled={isDisabled(['operator', 'version'])}
          onChange={(next) => onChange(setPath(value, ['operator', 'version'], next))}
        />
      </Box>

      <Divider sx={{ margin: '1.5rem 0' }} />

      <SectionHeading
        title="MeshSync"
        section="meshsync"
        governance={deploymentMode}
        id="controllers-config-meshsync"
      />
      {notice('meshsync')}

      <Box
        sx={{
          padding: '1rem',
          border: (theme) => `1px solid ${theme.palette.divider}`,
          borderRadius: 1,
          marginBottom: '1rem',
          background: (theme) =>
            alpha(theme.palette.background.brand?.default || theme.palette.primary.main, 0.03),
        }}
        data-testid="controllers-config-meshsync-filters"
      >
        <SubsectionTitle
          title="Discovery filters"
          chip="Applies in both modes"
          helpText="Limits what MeshSync publishes into Meshery after discovery. Available for Embedded and Operator modes."
        />
        <FieldRow>
          {listInput(
            'Output namespaces',
            ['meshsync', 'outputNamespaces'],
            'Comma-separated. Only these namespaces are published; empty publishes all.',
          )}
          {listInput(
            'Output resources',
            ['meshsync', 'outputResources'],
            'Comma-separated lowercase kinds (e.g. pod, deployment); empty publishes all.',
          )}
        </FieldRow>
      </Box>

      <Accordion
        expanded={meshsyncDeployOpen}
        onChange={(_, open) => setMeshsyncDeployOpen(open)}
        disableGutters
        elevation={0}
        sx={{
          border: (theme) => `1px solid ${theme.palette.divider}`,
          borderRadius: 1,
          '&:before': { display: 'none' },
        }}
        data-testid="controllers-config-accordion-meshsync"
      >
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: '0.5rem', flexWrap: 'wrap' }}>
            <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
              In-cluster MeshSync deployment
            </Typography>
            <InfoTooltip
              helpText="Version, replicas, security toggles, and watch scope for the MeshSync Deployment. Applies in Operator mode only."
              placement="top"
            />
            <Chip size="small" label="Operator mode only" variant="outlined" />
          </Box>
        </AccordionSummary>
        <AccordionDetails>
          <FieldRow>
            {textInput(
              'MeshSync version',
              ['meshsync', 'version'],
              'Image tag (operator mode). Applying rolls MeshSync pods.',
            )}
            {textInput('Replicas', ['meshsync', 'replicas'], '1-10 (operator mode).', {
              number: true,
              min: 1,
              max: 10,
            })}
            {triStateBoolean(
              'Secret redaction',
              ['meshsync', 'redactSecrets'],
              'Redacts Secret values before publishing. Requires MeshSync v1.0.2+.',
            )}
            {triStateBoolean(
              'Broker content dedup',
              ['meshsync', 'brokerContentDedup'],
              'Suppresses byte-identical republishes. Requires MeshSync v1.0.2+.',
            )}
            {triStateBoolean(
              'Debug logging',
              ['meshsync', 'debugLogging'],
              'DEBUG env on MeshSync.',
            )}
          </FieldRow>

          {redactSecrets !== true && !isInert(['meshsync', 'redactSecrets']) && (
            <Alert severity="warning" sx={{ marginTop: '1rem' }}>
              Secret redaction is disabled: Kubernetes Secret values within the watch scope are
              published un-redacted. Enable secret redaction or exclude Secrets from the watch
              scope.
            </Alert>
          )}

          <ControllersConfigWatchList
            label={fieldLabel(
              'Watched resources (discovery scope)',
              ['meshsync', 'watchList'],
              'At most one of whitelist or blacklist. Applying a watch-scope change restarts MeshSync pods automatically. In blacklist mode, enter one resource key per line in "<plural>.<version>.<group>" form.',
            )}
            value={getPath(value, ['meshsync', 'watchList']) as WatchList | undefined}
            disabled={isDisabled(['meshsync', 'watchList'])}
            onChange={(next) => onChange(setPath(value, ['meshsync', 'watchList'], next))}
          />
        </AccordionDetails>
      </Accordion>

      <Divider sx={{ margin: '1.5rem 0' }} />

      <SectionHeading
        title="Meshery Broker"
        section="broker"
        governance={deploymentMode}
        id="controllers-config-broker"
      />
      {notice('broker')}

      <Accordion
        expanded={brokerOpen}
        onChange={(_, open) => setBrokerOpen(open)}
        disableGutters
        elevation={0}
        sx={{
          border: (theme) => `1px solid ${theme.palette.divider}`,
          borderRadius: 1,
          '&:before': { display: 'none' },
        }}
        data-testid="controllers-config-accordion-broker"
      >
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: '0.5rem', flexWrap: 'wrap' }}>
            <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
              Broker deployment &amp; exposure
            </Typography>
            <InfoTooltip
              helpText="NATS image, replicas, and how the broker Service is reached. Applies in Operator mode only."
              placement="top"
            />
            <Chip size="small" label="Operator mode only" variant="outlined" />
          </Box>
        </AccordionSummary>
        <AccordionDetails>
          <FieldRow>
            {textInput(
              'Broker version',
              ['broker', 'version'],
              'NATS image tag (operator mode). Applying rolls broker pods.',
            )}
            {textInput('Replicas', ['broker', 'replicas'], '1-10 (operator mode).', {
              number: true,
              min: 1,
              max: 10,
            })}
            {enumSelect(
              'Service type',
              ['broker', 'service', 'type'],
              [
                { value: 'ClusterIP', label: 'ClusterIP (cluster-internal)' },
                { value: 'NodePort', label: 'NodePort' },
                { value: 'LoadBalancer', label: 'LoadBalancer' },
              ],
              'How the broker is exposed. Reconciles in place without restarting broker pods.',
              clearLoadBalancerFieldsUnlessLB,
            )}
            {textInput(
              'External endpoint override',
              ['broker', 'service', 'externalEndpointOverride'],
              'host:port; pins the advertised endpoint (ingress, NAT, air-gapped).',
            )}
            {isLoadBalancer &&
              textInput(
                'Load balancer class',
                ['broker', 'service', 'loadBalancerClass'],
                'LoadBalancer type only.',
              )}
            {isLoadBalancer &&
              listInput(
                'Load balancer source ranges',
                ['broker', 'service', 'loadBalancerSourceRanges'],
                'Comma-separated CIDRs allowed to reach the broker.',
              )}
            <Box>
              {fieldLabel(
                'Service annotations',
                ['broker', 'service', 'annotations'],
                'One key=value per line. Merged onto the broker client Service.',
              )}
              <TextField
                multiline
                minRows={2}
                size="small"
                disabled={isDisabled(['broker', 'service', 'annotations'])}
                value={annotationsText}
                placeholder={
                  'key=value\nservice.beta.kubernetes.io/aws-load-balancer-internal=true'
                }
                sx={{
                  width: fitWidth(
                    annotationsText,
                    'service.beta.kubernetes.io/aws-load-balancer-internal=true',
                  ),
                }}
                onChange={(e) => setAnnotationsFromText(e.target.value)}
              />
            </Box>
          </FieldRow>
        </AccordionDetails>
      </Accordion>
    </Box>
  );
}
