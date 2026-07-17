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
  MenuItem,
  TextField,
  Typography,
} from '@sistent/sistent';
import type { ControllersConfigDoc } from '@/rtk-query/controllersConfig';

// Built-in defaults mirrored from Meshery Server (BuiltInControllersConfig):
// what applies when no layer sets a field.
export const BUILT_IN_CONTROLLERS_CONFIG: ControllersConfigDoc = {
  operator: { deploymentMode: 'embedded' },
  meshsync: { replicas: 1, redactSecrets: false, brokerContentDedup: false, debugLogging: false },
  broker: { replicas: 1, service: { type: 'ClusterIP' } },
};

const INHERIT = '__inherit__';
const WATCH_EVENTS = ['ADDED', 'MODIFIED', 'DELETED'] as const;

type FieldPath = (string | number)[];

const getPath = (doc: ControllersConfigDoc | null | undefined, path: FieldPath): unknown => {
  let node: unknown = doc;
  for (const key of path) {
    if (node == null || typeof node !== 'object') return undefined;
    node = (node as Record<string | number, unknown>)[key];
  }
  return node;
};

// deleteAtPath removes the leaf at path and prunes any parents left empty,
// so an all-inherit section disappears from the document entirely.
const deleteAtPath = (node: unknown, path: FieldPath): void => {
  if (node == null || typeof node !== 'object' || path.length === 0) return;
  const obj = node as Record<string | number, unknown>;
  const [head, ...rest] = path;
  if (rest.length === 0) {
    delete obj[head];
    return;
  }
  deleteAtPath(obj[head], rest);
  const child = obj[head];
  if (child && typeof child === 'object' && Object.keys(child as object).length === 0) {
    delete obj[head];
  }
};

const setPath = (
  doc: ControllersConfigDoc,
  path: FieldPath,
  value: unknown,
): ControllersConfigDoc => {
  const next: ControllersConfigDoc = JSON.parse(JSON.stringify(doc ?? {}));
  if (value === undefined) {
    deleteAtPath(next, path);
    return next;
  }
  let node: Record<string | number, unknown> = next as Record<string | number, unknown>;
  for (const key of path.slice(0, -1)) {
    if (node[key] == null || typeof node[key] !== 'object') {
      node[key] = {};
    }
    node = node[key] as Record<string | number, unknown>;
  }
  node[path[path.length - 1]] = value;
  return next;
};

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
  disabled?: boolean;
};

/**
 * Layered editor for the Meshery Operator, MeshSync, and Broker
 * configuration. Every control is tri-state: leaving a field on "Inherit"
 * (or empty) keeps it absent from the document so the next layer applies.
 */
export default function ControllersConfigForm({
  value,
  onChange,
  inheritedLayers = [BUILT_IN_CONTROLLERS_CONFIG],
  inheritLabel = 'Inherited',
  showSourceIndicators = false,
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

  const fieldLabel = (text: string, path: FieldPath) => (
    <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '0.25rem' }}>
      <Typography variant="body2" sx={{ fontWeight: 500 }}>
        {text}
      </Typography>
      {sourceChip(path)}
    </Box>
  );

  const triStateBoolean = (label: string, path: FieldPath, helper?: string) => {
    const current = getPath(value, path) as boolean | undefined;
    const inherited = inheritedValue(path) as boolean | undefined;
    return (
      <Grid2 size={{ xs: 12, md: 4 }}>
        {fieldLabel(label, path)}
        <TextField
          select
          fullWidth
          size="small"
          disabled={disabled}
          value={current === undefined ? INHERIT : current ? 'true' : 'false'}
          onChange={(e) => {
            const v = e.target.value;
            onChange(setPath(value, path, v === INHERIT ? undefined : v === 'true'));
          }}
          helperText={helper}
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
        {fieldLabel(label, path)}
        <TextField
          fullWidth
          size="small"
          type={opts?.number ? 'number' : 'text'}
          disabled={disabled}
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
          helperText={helper}
        />
      </Grid2>
    );
  };

  const listInput = (label: string, path: FieldPath, helper: string) => {
    const current = getPath(value, path) as string[] | undefined;
    const inherited = inheritedValue(path) as string[] | undefined;
    return (
      <Grid2 size={{ xs: 12, md: 6 }}>
        {fieldLabel(label, path)}
        <TextField
          fullWidth
          size="small"
          disabled={disabled}
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
          helperText={helper}
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
        {fieldLabel(label, path)}
        <TextField
          select
          fullWidth
          size="small"
          disabled={disabled}
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
          helperText={helper}
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
      {/* Meshery Operator */}
      <Typography variant="subtitle1" sx={{ fontWeight: 600, marginBottom: '0.5rem' }}>
        Meshery Operator
      </Typography>
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
          'Helm chart version. Inherit tracks the Meshery Server release.',
        )}
      </Grid2>

      <Divider sx={{ margin: '1.5rem 0' }} />

      {/* MeshSync */}
      <Typography variant="subtitle1" sx={{ fontWeight: 600, marginBottom: '0.5rem' }}>
        MeshSync
      </Typography>
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

      {redactSecrets !== true && (
        <Alert severity="warning" sx={{ marginTop: '1rem' }}>
          Secret redaction is disabled: Kubernetes Secret values within the watch scope are
          published un-redacted. Enable secret redaction or exclude Secrets from the watch scope.
        </Alert>
      )}

      <Box sx={{ marginTop: '1.5rem' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', marginBottom: '0.25rem' }}>
          <Typography variant="body2" sx={{ fontWeight: 500 }}>
            Watched resources (discovery scope)
          </Typography>
          {sourceChip(['meshsync', 'watchList'])}
        </Box>
        <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: 1 }}>
          At most one of whitelist or blacklist. Applying a watch-scope change restarts MeshSync
          pods automatically.
        </Typography>
        <TextField
          select
          size="small"
          disabled={disabled}
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
                  disabled={disabled}
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
                        disabled={disabled}
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
                  disabled={disabled}
                  onClick={() => setWhitelist(whitelist.filter((_, i) => i !== index))}
                >
                  Remove
                </Button>
              </Box>
            ))}
            <Button
              size="small"
              variant="outlined"
              disabled={disabled}
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
            disabled={disabled}
            sx={{ marginTop: '1rem' }}
            value={blacklist.join('\n')}
            placeholder={'secrets.v1.\nevents.v1.'}
            helperText='One resource key per line, in "<plural>.<version>.<group>" form.'
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
      <Typography variant="subtitle1" sx={{ fontWeight: 600, marginBottom: '0.5rem' }}>
        Meshery Broker
      </Typography>
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
          {fieldLabel('Service annotations', ['broker', 'service', 'annotations'])}
          <TextField
            fullWidth
            multiline
            minRows={2}
            size="small"
            disabled={disabled}
            value={annotationsText}
            placeholder={'key=value\nservice.beta.kubernetes.io/aws-load-balancer-internal=true'}
            helperText="One key=value per line. Merged onto the broker client Service."
            onChange={(e) => setAnnotationsFromText(e.target.value)}
          />
        </Grid2>
      </Grid2>
    </Box>
  );
}
