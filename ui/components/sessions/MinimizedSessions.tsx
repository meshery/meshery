import React from 'react';
import { Badge, TerminalIcon, Tooltip, useTheme } from '@sistent/sistent';
import { restoreSessions, useSessionsDock } from './dock-store';
import { NavSessionsButton } from './session-styles';

export interface MinimizedSessionsProps {
  isDrawerCollapsed: boolean;
}

/**
 * The way back to a minimized sessions panel, at the foot of the Navigator.
 *
 * It renders nothing at all unless a panel is actually minimized, so the
 * Navigator's footer is unchanged for the ordinary case of no sessions.
 *
 * State comes from the module-scope dock store rather than `useSessions()`: the
 * Navigator is mounted above `SessionsProvider`, outside its context. See
 * `dock-store.ts` for why the provider stays where it is.
 */
export const MinimizedSessions: React.FC<MinimizedSessionsProps> = ({ isDrawerCollapsed }) => {
  const { count, minimized } = useSessionsDock();
  const theme = useTheme();

  if (!count || !minimized) return null;

  const label = `Restore ${count} session${count > 1 ? 's' : ''}`;

  return (
    <Tooltip title={label} placement={isDrawerCollapsed ? 'right' : 'top'}>
      <NavSessionsButton
        $collapsed={isDrawerCollapsed}
        aria-label={label}
        onClick={restoreSessions}
      >
        <Badge badgeContent={count} color="primary">
          <TerminalIcon width={20} height={20} fill={theme.palette.background.constant?.white} />
        </Badge>
        {!isDrawerCollapsed && <span>Sessions</span>}
      </NavSessionsButton>
    </Tooltip>
  );
};

export default MinimizedSessions;
