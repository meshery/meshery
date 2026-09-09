import React from 'react';
import {
  Alert,
  Box,
  Button,
  Checkbox,
  Chip,
  Divider,
  FormControlLabel,
  Grid2,
  InfoTooltip,
  MenuItem,
  TextField,
  Typography,
} from '@sistent/sistent';
import type { UpdateControllersDefaultConfigApiArg } from '@meshery/schemas/mesheryApi';
import { getPath, setPath, type FieldPath } from './fieldPath';
import {
  dormantPathsIn,
  isInertIn,
  type ConfigSection,
  type DeploymentModeGovernance,
} from './deploymentMode';
import { DeploymentModeBanner, SectionHeading, SectionNotice } from './DeploymentModeNotices';

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

const INHERIT = '__inherit__';
const WATCH_EVENTS = ['ADDED', 'MODIFIED', 'DELETED'] as const;

type SourceInfo = { label: string; overridden: boolean };

export type ControllersConfigFormProps = {
  /** The document being edited (only explicitly-set fields present). */
  value: ControllersConfigDoc;
  onChange: (next: ControllersConfigDoc) => void;
  /**
   * The layers this document inherits from, outermost first (for a
   * per-connection override: [server defaults, built-ins]; for the
   * server-wide defaults: [built-ins]).
   */
  inheritedLayers?: (ControllersConfigDoc | null | undefined)[];
  /** Label describing where inherited values come from, e.g. "Server default". */
  inheritLabel?: string;
  /** Show per-field source chips (used on the per-connection editor). */
  showSourceIndicators?: boolean;
  /**
   * The deployment mode governing this editor. Meshery Operator manages
   * MeshSync and Meshery Broker, so the mode decides which settings can reach
   * anything at all; without it the form treats every field as live.
   */
  deploymentMode?: DeploymentModeGovernance;
  disabled?: boolean;
};

/**
 * Layered editor for the Meshery Operator, MeshSync, and Broker
 * configuration. Every control is tri-state: leaving a field on "Inherit"
 * (or empty) keeps it absent from the document so the next layer applies.
 *
 * The form is governed by the effective deployment mode: on a connection
 * running MeshSync embedded in Meshery Server, the settings that configure
 * in-cluster objects are rendered inert and say why, and any value already
 * stored for them is shown as dormant with a way to clear it - rather than
 * being offered as a live control that Meshery Server will silently skip.
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

  // A field is inert when the connection's effective mode cannot apply it.
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
    return (
      <Grid2 size={{ xs: 12, md: 4 }}>
        {fieldLabel(label, path, helper)}
        <TextField
          select
          fullWidth
          size="small"
          disabled={isDisabled(path)}
          value={current === undefined ? INHERIT : current ? 'true' : 'false'}
          onChange={(e) => {
            const v = e.target.value;
            onChange(setPath(value, path, v === INHERIT ? undefined : v === 'true'));
          }}
        >
          <MenuItem value={INHERIT}>
            Inherit ({inherited === undefined ? 'unset' : inherited ? 'Enabled' : 'Disabled'})
          </MenuItem>
          <MenuItem value="true">Enabled</MenuItem>
          <MenuItem value="false">Disabled</MenuItem>
        </TextField>
      </Grid2>
    );
  };

  const textInput = (
    label: string,
    path: FieldPath,
    helper?: string,
    opts?: { number?: boolean; min?: number; max?: number; mdSize?: number },
  ) => {
    const current = getPath(value, path) as string | number | undefined;
    const inherited = inheritedValue(path);
    return (
      <Grid2 size={{ xs: 12, md: opts?.mdSize ?? 4 }}>
        {fieldLabel(label, path, helper)}
        <TextField
          fullWidth
          size="small"
          type={opts?.number ? 'number' : 'text'}
          disabled={isDisabled(path)}
          value={current ?? ''}
          placeholder={inherited !== undefined ? `Inherit (${inherited})` : 'Inherit'}
          slotProps={opts?.number ? { htmlInput: { min: opts?.min, max: opts?.max } } : undefined}
          onChange={(e) => {
            const raw = e.target.value;
            if (raw === '') {
              onChange(setPath(value, path, undefined));
              return;
            }
            onChange(setPath(value, path, opts?.number ? Number(raw) : raw));
          }}
        />
      </Grid2>
    );
  };

  const listInput = (label: string, path: FieldPath, helper: string) => {
    const current = getPath(value, path) as string[] | undefined;
    const inherited = inheritedValue(path) as string[] | undefined;
    return (
      <Grid2 size={{ xs: 12, md: 6 }}>
        {fieldLabel(label, path, helper)}
        <TextField
          fullWidth
          size="small"
          disabled={isDisabled(path)}
          value={current ? current.join(', ') : ''}
          placeholder={
            inherited && inherited.length > 0
              ? `Inherit (${inherited.join(', ')})`
              : 'Inherit (all)'
          }
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
      </Grid2>
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
    return (
      <Grid2 size={{ xs: 12, md: 4 }}>
        {fieldLabel(label, path, helper)}
        <TextField
          select
          fullWidth
          size="small"
          disabled={isDisabled(path)}
          value={current ?? INHERIT}
          onChange={(e) => {
            const v = e.target.value;
            const selected = v === INHERIT ? undefined : v;
            let next = setPath(value, path, selected);
            if (postProcess) {
              next = postProcess(next, selected);
            }
            onChange(next);
          }}
        >
          <MenuItem value={INHERIT}>Inherit ({inherited ?? 'unset'})</MenuItem>
          {options.map((option) => (
            <MenuItem key={option.value} value={option.value}>
              {option.label}
            </MenuItem>
          ))}
        </TextField>
      </Grid2>
    );
  };

  // clearLoadBalancerFieldsUnlessLB drops the LoadBalancer-only service
  // fields whenever the effective service type is not LoadBalancer: the
  // inputs are hidden then, and stale values would trip server-side
  // validation the user cannot see or clear from the form.
  const clearLoadBalancerFieldsUnlessLB = (
    next: ControllersConfigDoc,
    selected: string | undefined,
  ): ControllersConfigDoc => {
    const effectiveType =
      selected ?? (inheritedValue(['broker', 'service', 'type']) as string | undefined);
    if (effectiveType === 'LoadBalancer') {
      return next;
    }
    let cleared = setPath(next, ['broker', 'service', 'loadBalancerClass'], undefined);
    cleared = setPath(cleared, ['broker', 'service', 'loadBalancerSourceRanges'], undefined);
    return cleared;
  };

  // Watch list -------------------------------------------------------------
  const watchList = getPath(value, ['meshsync', 'watchList']) as
    | { whitelist?: { resource: string; events?: string[] }[]; blacklist?: string[] }
    | undefined;
  const watchMode = !watchList ? INHERIT : watchList.whitelist ? 'whitelist' : 'blacklist';

  const setWatchMode = (mode: string) => {
    if (mode === INHERIT) {
      onChange(setPath(value, ['meshsync', 'watchList'], undefined));
    } else if (mode === 'whitelist') {
      onChange(setPath(value, ['meshsync', 'watchList'], { whitelist: [] }));
    } else {
      onChange(setPath(value, ['meshsync', 'watchList'], { blacklist: [] }));
    }
  };

  const whitelist = watchList?.whitelist ?? [];
  const setWhitelist = (rows: { resource: string; events?: string[] }[]) =>
    onChange(setPath(value, ['meshsync', 'watchList'], { whitelist: rows }));

  const blacklist = watchList?.blacklist ?? [];
  const setBlacklist = (rows: string[]) =>
    onChange(setPath(value, ['meshsync', 'watchList'], { blacklist: rows }));

  const redactSecrets =
    getPath(value, ['meshsync', 'redactSecrets']) ?? inheritedValue(['meshsync', 'redactSecrets']);

  // Annotations ------------------------------------------------------------
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
      if (idx > 0) {
        next[line.slice(0, idx).trim()] = line.slice(idx + 1).trim();
      }
    });
    onChange(setPath(value, ['broker', 'service', 'annotations'], next));
  };

  const serviceType = (getPath(value, ['broker', 'service', 'type']) ??
    inheritedValue(['broker', 'service', 'type'])) as string | undefined;
  const isLoadBalancer = serviceType === 'LoadBalancer';

  return (
    <Box>
      <DeploymentModeBanner governance={deploymentMode} />

      {/* Meshery Operator */}
      <SectionHeading title="Meshery Operator" section="operator" governance={deploymentMode} />
      {notice('operator')}
      <Grid2 container spacing={2}>
        {enumSelect(
          'Deployment mode',
          ['operator', 'deploymentMode'],
          [
            { value: 'operator', label: 'Operator (in-cluster)' },
            { value: 'embedded', label: 'Embedded (in Meshery Server)' },
          ],
          'Operator installs MeshSync and Broker into the cluster; Embedded runs MeshSync inside Meshery Server. Changing the mode redeploys controllers.',
        )}
        {textInput(
          'Operator version',
          ['operator', 'version'],
          'Helm chart version (operator mode). Applying upgrades the Meshery Operator release. Inherit tracks the Meshery Server release.',
        )}
      </Grid2>

      <Divider sx={{ margin: '1.5rem 0' }} />

      {/* MeshSync */}
      <SectionHeading title="MeshSync" section="meshsync" governance={deploymentMode} />
      {notice('meshsync')}
      <Grid2 container spacing={2}>
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
        {triStateBoolean('Debug logging', ['meshsync', 'debugLogging'], 'DEBUG env on MeshSync.')}
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
      </Grid2>

      {/* In embedded mode this setting is not what decides redaction - the
          Meshery Server process environment is - so the warning would name the
          wrong control. The section notice states that instead. */}
      {redactSecrets !== true && !isInert(['meshsync', 'redactSecrets']) && (
        <Alert severity="warning" sx={{ marginTop: '1rem' }}>
          Secret redaction is disabled: Kubernetes Secret values within the watch scope are
          published un-redacted. Enable secret redaction or exclude Secrets from the watch scope.
        </Alert>
      )}

      <Box sx={{ marginTop: '1.5rem' }}>
        {fieldLabel(
          'Watched resources (discovery scope)',
          ['meshsync', 'watchList'],
          'At most one of whitelist or blacklist. Applying a watch-scope change restarts MeshSync pods automatically. In blacklist mode, enter one resource key per line in "<plural>.<version>.<group>" form.',
        )}
        <TextField
          select
          size="small"
          disabled={isDisabled(['meshsync', 'watchList'])}
          value={watchMode}
          onChange={(e) => setWatchMode(e.target.value)}
          sx={{ minWidth: '260px' }}
        >
          <MenuItem value={INHERIT}>Inherit</MenuItem>
          <MenuItem value="whitelist">Whitelist (watch only these)</MenuItem>
          <MenuItem value="blacklist">Blacklist (default scope minus these)</MenuItem>
        </TextField>

        {watchMode === 'whitelist' && (
          <Box sx={{ marginTop: '1rem' }}>
            {whitelist.map((row, index) => (
              <Box
                key={index}
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '0.75rem',
                  marginBottom: '0.5rem',
                }}
              >
                <TextField
                  size="small"
                  disabled={isDisabled(['meshsync', 'watchList'])}
                  value={row.resource}
                  placeholder="pods.v1. or deployments.v1.apps"
                  sx={{ minWidth: '280px' }}
                  onChange={(e) => {
                    const rows = [...whitelist];
                    rows[index] = { ...rows[index], resource: e.target.value };
                    setWhitelist(rows);
                  }}
                />
                {WATCH_EVENTS.map((eventType) => (
                  <FormControlLabel
                    key={eventType}
                    control={
                      <Checkbox
                        size="small"
                        disabled={isDisabled(['meshsync', 'watchList'])}
                        checked={(row.events ?? []).includes(eventType)}
                        onChange={(e) => {
                          const rows = [...whitelist];
                          const events = new Set(rows[index].events ?? []);
                          if (e.target.checked) {
                            events.add(eventType);
                          } else {
                            events.delete(eventType);
                          }
                          rows[index] = { ...rows[index], events: Array.from(events) };
                          setWhitelist(rows);
                        }}
                      />
                    }
                    label={eventType}
                  />
                ))}
                <Button
                  size="small"
                  disabled={isDisabled(['meshsync', 'watchList'])}
                  onClick={() => setWhitelist(whitelist.filter((_, i) => i !== index))}
                >
                  Remove
                </Button>
              </Box>
            ))}
            <Button
              size="small"
              variant="outlined"
              disabled={isDisabled(['meshsync', 'watchList'])}
              onClick={() =>
                setWhitelist([...whitelist, { resource: '', events: [...WATCH_EVENTS] }])
              }
            >
              Add resource
            </Button>
          </Box>
        )}

        {watchMode === 'blacklist' && (
          <TextField
            fullWidth
            multiline
            minRows={3}
            size="small"
            disabled={isDisabled(['meshsync', 'watchList'])}
            sx={{ marginTop: '1rem' }}
            value={blacklist.join('\n')}
            placeholder={'secrets.v1.\nevents.v1.'}
            onChange={(e) =>
              setBlacklist(
                e.target.value
                  .split('\n')
                  .map((s) => s.trim())
                  .filter(Boolean),
              )
            }
          />
        )}
      </Box>

      <Divider sx={{ margin: '1.5rem 0' }} />

      {/* Meshery Broker */}
      <SectionHeading title="Meshery Broker" section="broker" governance={deploymentMode} />
      {notice('broker')}
      <Grid2 container spacing={2}>
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
        <Grid2 size={{ xs: 12, md: 6 }}>
          {fieldLabel(
            'Service annotations',
            ['broker', 'service', 'annotations'],
            'One key=value per line. Merged onto the broker client Service.',
          )}
          <TextField
            fullWidth
            multiline
            minRows={2}
            size="small"
            disabled={isDisabled(['broker', 'service', 'annotations'])}
            value={annotationsText}
            placeholder={'key=value\nservice.beta.kubernetes.io/aws-load-balancer-internal=true'}
            onChange={(e) => setAnnotationsFromText(e.target.value)}
          />
        </Grid2>
      </Grid2>
    </Box>
  );
}
