import React from 'react';
import { ArticleIcon, Box, IconButton, TerminalIcon, Tooltip } from '@sistent/sistent';
import { getK8sContextFromClusterId } from '@/utils/multi-ctx';
import type { ConsoleTarget } from 'lib/console/protocol';
import { useConsole } from './ConsoleProvider';

export interface ConsoleActionsCellProps {
  /** A MeshSync resource row. */
  resource: Record<string, any> | undefined;
  k8sConfig: unknown;
}

/**
 * The resource kinds a console can attach to, in the form MeshSync reports them.
 * Every one of them runs pods, and the console driver resolves an owner workload
 * onto the pod it runs, so a Deployment row opens a console on a real pod rather
 * than asking the API server for a pod named after the deployment.
 */
const CONSOLE_RESOURCES = new Set([
  'pod',
  'deployment',
  'statefulset',
  'daemonset',
  'job',
  'replicaset',
  'replicationcontroller',
]);

/** The console resource kind of a MeshSync row, or undefined if it hosts none. */
const consoleResourceOf = (resource: Record<string, any> | undefined): string | undefined => {
  const kind = typeof resource?.kind === 'string' ? resource.kind.trim().toLowerCase() : '';
  return CONSOLE_RESOURCES.has(kind) ? kind : undefined;
};

/**
 * Where a row's console buttons would attach: the connection it belongs to and
 * the resource to address. Returns null for a row that names nothing a console can
 * attach to — no connection, no name, or a kind that runs no containers.
 */
export const consoleAddressFor = (
  resource: Record<string, any> | undefined,
  k8sConfig: unknown,
): { connectionId: string; target: ConsoleTarget } | null => {
  const resourceKind = consoleResourceOf(resource);
  const name = resource?.metadata?.name;
  const connectionId = resource?.cluster_id
    ? getK8sContextFromClusterId(resource.cluster_id, k8sConfig)?.connectionId
    : undefined;

  if (!connectionId || !resourceKind || !name) return null;

  return {
    connectionId,
    target: { resource: resourceKind, namespace: resource.metadata.namespace, name },
  };
};

/**
 * The actions that open a log stream or a shell against a resource — on a table
 * row, and in the header of that resource's detail view.
 *
 * Which kind of resource this is comes from the row itself; whether that target
 * really admits a console is the server's call, resolved against live state by the
 * capabilities endpoint once the console mounts, which answers with a reason the
 * panel can show. So the buttons are offered for every console-capable kind
 * whatever the cached row says about its state: a greyed-out control against a pod
 * that is in fact running reads as a broken button, not a busy pod.
 */
const ConsoleActionsCell: React.FC<ConsoleActionsCellProps> = ({ resource, k8sConfig }) => {
  const { openConsole } = useConsole();
  const address = consoleAddressFor(resource, k8sConfig);

  if (!address) return null;

  const open = (kind: 'terminal' | 'logs') => () => openConsole({ ...address, kind });

  return (
    <Box sx={{ display: 'flex', gap: '0.25rem' }}>
      <Tooltip title="View logs">
        <IconButton size="small" aria-label="View logs" onClick={open('logs')}>
          {/*
           * `fill="currentColor"`: sistent's icons carry no fill of their own, so
           * SVG defaults them to black — invisible against the dark table. Setting
           * it on the root svg is inherited by the paths, and `currentColor` tracks
           * the IconButton's own colour in either theme.
           */}
          <ArticleIcon width={18} height={18} fill="currentColor" />
        </IconButton>
      </Tooltip>

      <Tooltip title="Open a terminal">
        <IconButton size="small" aria-label="Open a terminal" onClick={open('terminal')}>
          <TerminalIcon width={18} height={18} fill="currentColor" />
        </IconButton>
      </Tooltip>
    </Box>
  );
};

export default ConsoleActionsCell;
