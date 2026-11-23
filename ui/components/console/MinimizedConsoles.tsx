import React from 'react';
import { Badge, TerminalIcon, Tooltip, useTheme } from '@sistent/sistent';
import { restoreConsoles, useConsoleDock } from './dock-store';
import { NavConsolesButton } from './console-styles';

export interface MinimizedConsolesProps {
  isDrawerCollapsed: boolean;
}

/**
 * The way back to a minimized consoles panel, at the foot of the Navigator.
 *
 * It renders nothing at all unless a panel is actually minimized, so the
 * Navigator's footer is unchanged for the ordinary case of no consoles.
 *
 * State comes from the module-scope dock store rather than `useConsole()`: the
 * Navigator is mounted above `ConsoleProvider`, outside its context. See
 * `dock-store.ts` for why the provider stays where it is.
 */
export const MinimizedConsoles: React.FC<MinimizedConsolesProps> = ({ isDrawerCollapsed }) => {
  const { count, minimized } = useConsoleDock();
  const theme = useTheme();

  if (!count || !minimized) return null;

  const label = `Restore ${count} console${count > 1 ? 's' : ''}`;

  return (
    <Tooltip title={label} placement={isDrawerCollapsed ? 'right' : 'top'}>
      <NavConsolesButton
        $collapsed={isDrawerCollapsed}
        aria-label={label}
        onClick={restoreConsoles}
      >
        <Badge badgeContent={count} color="primary">
          <TerminalIcon width={20} height={20} fill={theme.palette.background.constant?.white} />
        </Badge>
        {!isDrawerCollapsed && <span>Consoles</span>}
      </NavConsolesButton>
    </Tooltip>
  );
};

export default MinimizedConsoles;
