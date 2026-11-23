import React, { useCallback } from 'react';
import { ArticleIcon, Box, IconButton, TerminalIcon, Tooltip } from '@sistent/sistent';
import { getK8sContextFromClusterId } from '@/utils/multi-ctx';
import type { ConsoleTarget } from 'lib/console/protocol';
import { useConsole } from './ConsoleProvider';

export interface ConsoleActionsCellProps {
  /** A MeshSync resource row. */
  resource: Record<string, any> | undefined;
  k8sConfig: unknown;
}

/** The resource types that can host a console, keyed as MeshSync reports them. */
const CONSOLE_RESOURCE = 'pod';

/**
 * The actions that open a log stream or a shell against a resource — on a table
 * row, and in the header of that resource's detail view.
 *
 * Whether a target really admits a console is the server's call, resolved against
 * live state by the capabilities endpoint once the console mounts, which answers
 * with a reason the panel can show. So the buttons are offered unconditionally
 * rather than gated on the phase in the cached MeshSync row: that row goes stale,
 * and a control greyed out against a pod that is in fact running — with the
 * explanation hidden behind a tooltip — reads as a broken button, not a busy pod.
 */
const ConsoleActionsCell: React.FC<ConsoleActionsCellProps> = ({ resource, k8sConfig }) => {
  const { openConsole } = useConsole();

  const connectionId = resource?.cluster_id
    ? getK8sContextFromClusterId(resource.cluster_id, k8sConfig)?.connectionId
    : undefined;

  const target: ConsoleTarget | null = resource?.metadata?.name
    ? {
        resource: CONSOLE_RESOURCE,
        namespace: resource.metadata.namespace,
        name: resource.metadata.name,
      }
    : null;

  const open = useCallback(
    (kind: 'terminal' | 'logs') => () => {
      if (connectionId && target) {
        openConsole({ connectionId, target, kind });
      }
    },
    [connectionId, target, openConsole],
  );

  if (!connectionId || !target) return null;

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
