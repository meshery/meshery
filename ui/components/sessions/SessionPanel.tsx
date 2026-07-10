import React, { useCallback, useState } from 'react';
import {
  Alert,
  Box,
  CircularProgress,
  FullScreenExitIcon,
  FullScreenIcon,
  IconButton,
  Tooltip,
  Typography,
} from '@sistent/sistent';
import { useGetSessionCapabilitiesQuery } from '@meshery/schemas/mesheryApi';
import type { SessionKind, SessionTarget } from 'lib/sessions/protocol';
import { describeSessionError } from 'lib/sessions/errors';
import TerminalSession from './TerminalSession';
import LogSession from './LogSession';
import { SessionHost } from './session-styles';

export interface SessionPanelProps {
  connectionId: string;
  target: SessionTarget;
  kind: SessionKind;
  active?: boolean;
  /** Hides the fullscreen control, e.g. where the host already owns that chrome. */
  allowFullScreen?: boolean;
}

/**
 * Renders one session of the requested kind, once the server confirms the
 * target admits it.
 *
 * Capabilities are resolved before the socket is opened so an unsupported
 * target produces a readable explanation instead of a WebSocket that opens and
 * immediately closes. The server checks again on the handshake — this is a UX
 * affordance, not the authorization boundary.
 */
const SessionPanel: React.FC<SessionPanelProps> = ({
  connectionId,
  target,
  kind,
  active = true,
  allowFullScreen = true,
}) => {
  const [fullScreen, setFullScreen] = useState(false);
  const toggleFullScreen = useCallback(() => setFullScreen((current) => !current), []);

  // The generated argument type is flat, so the target is spread rather than
  // nested. RTK serializes the argument for its cache key, so an object literal
  // here does not re-fetch on every render.
  const {
    data: capabilities,
    isLoading,
    error,
  } = useGetSessionCapabilitiesQuery(
    {
      connectionId,
      resource: target?.resource,
      name: target?.name,
      namespace: target?.namespace,
      container: target?.container,
    },
    { skip: !active || !connectionId || !target?.name || !target?.resource },
  );

  if (!active) return null;

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
        <CircularProgress size={24} />
      </Box>
    );
  }

  if (error) {
    // Surface what the server actually said. A stale MeshSync view routinely
    // offers sessions on resources the cluster has already deleted, and
    // "something went wrong" sends the user looking for a bug that isn't there.
    const described = describeSessionError(error, target?.name);
    return (
      <Alert severity="error" variant="outlined">
        <Typography variant="body2">{described.message}</Typography>
        {described.detail && (
          <Typography variant="caption" component="p">
            {described.detail}
          </Typography>
        )}
        {described.code && (
          <Typography variant="caption" component="p" sx={{ opacity: 0.7 }}>
            {described.code}
          </Typography>
        )}
      </Alert>
    );
  }

  const supported = kind === 'terminal' ? capabilities?.terminal : capabilities?.logs;
  if (!supported) {
    return (
      <Alert severity="info" variant="outlined">
        {capabilities?.reason ?? `This resource does not support a ${kind} session.`}
      </Alert>
    );
  }

  const resolved: SessionTarget = {
    ...target,
    container: target.container || capabilities?.defaultContainer,
  };

  // The fullscreen control lives at the right of the session's own toolbar. It
  // used to float over the top-right of the pane, where it sat on top of the log
  // toolbar's controls.
  const fullScreenControl = allowFullScreen ? (
    <Tooltip title={fullScreen ? 'Exit fullscreen' : 'Fullscreen'}>
      <IconButton
        size="small"
        aria-label={fullScreen ? 'Exit fullscreen' : 'Enter fullscreen'}
        onClick={toggleFullScreen}
      >
        {/*
         * `fill="currentColor"`: unlike sistent's other icons these two take no
         * `fill` prop, and their path carries no fill attribute — so SVG
         * defaults it to black, which vanishes against a dark terminal. Setting
         * it on the root svg is inherited by the path, and `currentColor`
         * tracks the IconButton's own colour in either theme.
         */}
        {fullScreen ? (
          <FullScreenExitIcon width={16} height={16} fill="currentColor" />
        ) : (
          <FullScreenIcon width={16} height={16} fill="currentColor" />
        )}
      </IconButton>
    </Tooltip>
  ) : null;

  const session =
    kind === 'terminal' ? (
      <TerminalSession
        connectionId={connectionId}
        target={resolved}
        containers={capabilities?.containers ?? []}
        toolbarEnd={fullScreenControl}
        active={active}
      />
    ) : (
      <LogSession
        connectionId={connectionId}
        target={resolved}
        containers={capabilities?.containers ?? []}
        toolbarEnd={fullScreenControl}
        active={active}
      />
    );

  // One host element in both states, toggled by a style prop. Swapping the
  // wrapper element instead would move the session in the element tree, and
  // React would unmount and remount it — killing a live shell every time the
  // user hit fullscreen.
  return <SessionHost $fullScreen={fullScreen}>{session}</SessionHost>;
};

export default SessionPanel;
