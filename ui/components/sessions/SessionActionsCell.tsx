import React, { useCallback } from 'react';
import { ArticleIcon, Box, IconButton, TerminalIcon, Tooltip } from '@sistent/sistent';
import { getK8sContextFromClusterId } from '@/utils/multi-ctx';
import type { SessionTarget } from 'lib/sessions/protocol';
import { useSessions } from './SessionsProvider';

export interface SessionActionsCellProps {
  /** A MeshSync resource row. */
  resource: Record<string, any> | undefined;
  k8sConfig: unknown;
}

/** The resource types that can host a session, keyed as MeshSync reports them. */
const SESSION_RESOURCE = 'pod';

/**
 * Row actions that open a log stream or a shell against the row's resource.
 *
 * Whether a target really admits a session is the server's call, resolved
 * against live state by the capabilities endpoint once the panel mounts. This
 * cell only decides whether to *offer* the action, from the cached MeshSync
 * row: a button that opens onto "the pod is Pending" is worse than no button,
 * but a stale row must not be able to hide a session the server would allow.
 */
const SessionActionsCell: React.FC<SessionActionsCellProps> = ({ resource, k8sConfig }) => {
  const { openSession } = useSessions();

  const connectionId = resource?.cluster_id
    ? getK8sContextFromClusterId(resource.cluster_id, k8sConfig)?.connectionId
    : undefined;

  const target: SessionTarget | null = resource?.metadata?.name
    ? {
        resource: SESSION_RESOURCE,
        namespace: resource.metadata.namespace,
        name: resource.metadata.name,
      }
    : null;

  const open = useCallback(
    (kind: 'terminal' | 'logs') => () => {
      if (connectionId && target) {
        openSession({ connectionId, target, kind });
      }
    },
    [connectionId, target, openSession],
  );

  if (!connectionId || !target) return null;

  // MeshSync stores the pod status as a JSON blob in `status.attribute`.
  let phase: string | undefined;
  try {
    phase = JSON.parse(resource.status?.attribute ?? '{}')?.phase;
  } catch {
    phase = undefined;
  }
  const running = phase === 'Running';

  return (
    <Box sx={{ display: 'flex', gap: '0.25rem' }}>
      <Tooltip title="View logs">
        <IconButton size="small" aria-label="View logs" onClick={open('logs')}>
          <ArticleIcon width={18} height={18} />
        </IconButton>
      </Tooltip>

      <Tooltip
        title={
          running
            ? 'Open a terminal'
            : `A terminal needs a running pod (this one is ${phase ?? 'not running'})`
        }
      >
        {/* A disabled button emits no events, so the tooltip needs a live wrapper. */}
        <Box component="span" sx={{ display: 'inline-flex' }}>
          <IconButton
            size="small"
            aria-label="Open a terminal"
            disabled={!running}
            onClick={open('terminal')}
          >
            <TerminalIcon width={18} height={18} />
          </IconButton>
        </Box>
      </Tooltip>
    </Box>
  );
};

export default SessionActionsCell;
