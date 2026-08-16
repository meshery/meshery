import React, { useEffect, useRef, useState, type ReactNode } from 'react';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Alert,
  Box,
  Button,
  Checkbox,
  Chip,
  Divider,
  ExpandMoreIcon,
  FormControlLabel,
  FormGroup,
  InfoTooltip,
  MenuItem,
  Stack,
  TextField,
  Typography,
} from '@sistent/sistent';
import type { UpdateControllersDefaultConfigApiArg } from '@meshery/schemas/mesheryApi';
import { alpha } from '@/theme';
import { getPath, setPath, type FieldPath } from './fieldPath';
import {
  INHERIT,
  WATCH_EVENTS,
  WATCH_MODE_OPTIONS,
  dormantPathsIn,
  fitWidth,
  type ConfigSection,
  type DeploymentMode,
  type DeploymentModeGovernance,
  type WatchList,
} from './deploymentMode';
import { DeploymentModeBanner, SectionHeading, SectionNotice } from './DeploymentModeNotices';
import {
  ControllersConfigModePicker,
  OperatorVersionField,
  controlSx,
  createControllersConfigFields,
  fieldRowSx,
} from './ControllersConfigFields';

function ControllersConfigWatchList({
  label,
  value,
  disabled,
  onChange,
}: {
  label: ReactNode;
  value: WatchList | undefined;
  disabled: boolean;
  onChange: (next: WatchList | undefined) => void;
}) {
  const watchMode = !value ? INHERIT : value.whitelist ? 'whitelist' : 'blacklist';
  const whitelist = value?.whitelist ?? [];
  const blacklist = value?.blacklist ?? [];
  // Client-only row keys so React does not reuse the wrong row when one is
  // removed from the middle. Never written to the config document.
  const whitelistKeysRef = useRef<string[]>([]);
  if (whitelist.length === 0) {
    whitelistKeysRef.current = [];
  } else if (whitelistKeysRef.current.length < whitelist.length) {
    whitelistKeysRef.current = [
      ...whitelistKeysRef.current,
      ...Array.from({ length: whitelist.length - whitelistKeysRef.current.length }, () =>
        crypto.randomUUID(),
      ),
    ];
  }

  return (
    <Box
      sx={{ marginTop: '1rem', display: 'flex', flexDirection: 'column', alignItems: 'flex-start' }}
      data-testid="controllers-config-watch-list"
    >
      {label}
      <TextField
        select
        size="small"
        disabled={disabled}
        value={watchMode}
        aria-label="Watch mode"
        slotProps={{ htmlInput: { 'aria-label': 'Watch mode' } }}
        onChange={(e) => {
          const mode = e.target.value;
          whitelistKeysRef.current = [];
          if (mode === INHERIT) onChange(undefined);
          else if (mode === 'whitelist') onChange({ whitelist: [] });
          else onChange({ blacklist: [] });
        }}
        sx={{
          width: fitWidth(...WATCH_MODE_OPTIONS.map((option) => option.label)),
          maxWidth: '100%',
        }}
      >
        {WATCH_MODE_OPTIONS.map((option) => (
          <MenuItem key={option.value} value={option.value}>
            {option.label}
          </MenuItem>
        ))}
      </TextField>

      {watchMode === 'whitelist' && (
        <Box sx={{ marginTop: '1rem', width: '100%' }}>
          {whitelist.map((row, index) => (
            <Stack
              key={whitelistKeysRef.current[index]}
              direction="row"
              spacing={1.5}
              sx={{ alignItems: 'center', mb: 0.5 }}
            >
              <TextField
                size="small"
                disabled={disabled}
                value={row.resource}
                placeholder="pods.v1. or deployments.v1.apps"
                aria-label={`Resource ${index + 1}`}
                slotProps={{ htmlInput: { 'aria-label': `Resource ${index + 1}` } }}
                sx={{
                  width: fitWidth(row.resource, 'pods.v1. or deployments.v1.apps'),
                  maxWidth: '100%',
                }}
                onChange={(e) => {
                  const rows = [...whitelist];
                  rows[index] = { ...rows[index], resource: e.target.value };
                  onChange({ whitelist: rows });
                }}
              />
              <FormGroup row>
                {WATCH_EVENTS.map((eventType) => (
                  <FormControlLabel
                    key={eventType}
                    control={
                      <Checkbox
                        size="small"
                        disabled={disabled}
                        checked={(row.events ?? []).includes(eventType)}
                        onChange={(e) => {
                          const rows = [...whitelist];
                          const events = new Set(rows[index].events ?? []);
                          if (e.target.checked) events.add(eventType);
                          else events.delete(eventType);
                          rows[index] = { ...rows[index], events: Array.from(events) };
                          onChange({ whitelist: rows });
                        }}
                      />
                    }
                    label={eventType}
                  />
                ))}
              </FormGroup>
              <Button
                size="small"
                color="error"
                disabled={disabled}
                onClick={() => {
                  whitelistKeysRef.current = whitelistKeysRef.current.filter((_, i) => i !== index);
                  onChange({ whitelist: whitelist.filter((_, i) => i !== index) });
                }}
              >
                Remove
              </Button>
            </Stack>
          ))}
          <Button
            size="small"
            variant="outlined"
            color="primary"
            disabled={disabled}
            onClick={() => {
              whitelistKeysRef.current = [...whitelistKeysRef.current, crypto.randomUUID()];
              onChange({
                whitelist: [...whitelist, { resource: '', events: [...WATCH_EVENTS] }],
              });
            }}
          >
            Add resource
          </Button>
        </Box>
      )}

      {watchMode === 'blacklist' && (
        <TextField
          multiline
          minRows={3}
          size="small"
          disabled={disabled}
          sx={{ marginTop: '1rem', width: '100%' }}
          aria-label="Blacklist resources"
          slotProps={{ htmlInput: { 'aria-label': 'Blacklist resources' } }}
          value={blacklist.join('\n')}
          placeholder={'secrets.v1.\nevents.v1.'}
          onChange={(e) =>
            onChange({
              blacklist: e.target.value
                .split('\n')
                .map((s) => s.trim())
                .filter(Boolean),
            })
          }
        />
      )}
    </Box>
  );
}

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

export type ControllersConfigFormProps = {
  value: ControllersConfigDoc;
  onChange: (next: ControllersConfigDoc) => void;
  inheritedLayers?: (ControllersConfigDoc | null | undefined)[];
  inheritLabel?: string;
  showSourceIndicators?: boolean;
  deploymentMode?: DeploymentModeGovernance;
  disabled?: boolean;
};

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
  const liveGovernance = deploymentMode
    ? {
        ...deploymentMode,
        mode: modeForDisclosure,
        unsaved: modeForDisclosure !== deploymentMode.mode || deploymentMode.unsaved,
      }
    : undefined;
  const [meshsyncDeployOpen, setMeshsyncDeployOpen] = useState(operatorModeApplies);
  const [brokerOpen, setBrokerOpen] = useState(operatorModeApplies);

  useEffect(() => {
    setMeshsyncDeployOpen(operatorModeApplies);
    setBrokerOpen(operatorModeApplies);
  }, [operatorModeApplies]);

  const {
    inheritedValue,
    isInert,
    isDisabled,
    fieldLabel,
    triStateBoolean,
    textInput,
    listInput,
    enumSelect,
  } = createControllersConfigFields({
    value,
    onChange,
    inheritedLayers,
    inheritLabel,
    showSourceIndicators,
    disabled,
    liveGovernance,
  });

  const clearDormant = (section: ConfigSection) => {
    onChange(
      dormantPathsIn(liveGovernance, value, section).reduce(
        (doc, path) => setPath(doc, path, undefined),
        value,
      ),
    );
  };

  const notice = (section: ConfigSection) => (
    <SectionNotice
      section={section}
      governance={liveGovernance}
      value={value}
      onClearDormant={clearDormant}
      disabled={disabled}
    />
  );

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
      <DeploymentModeBanner governance={liveGovernance} />

      <SectionHeading
        title="Meshery Operator"
        section="operator"
        governance={liveGovernance}
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
        governance={liveGovernance}
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
        <Stack
          direction="row"
          spacing={0.5}
          sx={{ alignItems: 'center', flexWrap: 'wrap', mb: 1.5 }}
        >
          <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
            Discovery filters
          </Typography>
          <InfoTooltip
            helpText="Limits what MeshSync publishes into Meshery after discovery. Available for Embedded and Operator modes."
            placement="top"
          />
          <Chip size="small" label="Applies in both modes" variant="outlined" />
        </Stack>
        <Stack direction="row" spacing={2} useFlexGap sx={fieldRowSx}>
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
        </Stack>
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
          <Stack direction="row" spacing={2} useFlexGap sx={fieldRowSx}>
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
          </Stack>

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
        governance={liveGovernance}
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
          <Stack direction="row" spacing={2} useFlexGap sx={fieldRowSx}>
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
            <Box
              sx={{
                width: fitWidth(
                  annotationsText,
                  'service.beta.kubernetes.io/aws-load-balancer-internal=true',
                ),
                maxWidth: '100%',
              }}
            >
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
                sx={controlSx(false)}
                onChange={(e) => setAnnotationsFromText(e.target.value)}
              />
            </Box>
          </Stack>
        </AccordionDetails>
      </Accordion>
    </Box>
  );
}
