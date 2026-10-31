import { FormControl, MenuItem } from '@sistent/sistent';
import { ConnectionStyledSelect, ConnectionStyledMenuItem } from './styles';
import { ConnectionStateChip } from './ConnectionChip';
import { getNextStates, type ConnectionTransitionMap } from './ConnectionTable.constants';

type ConnectionStatusSelectProps = {
  /** The connection's current lifecycle status. */
  status: string;
  /** Per-kind transition map (from the connection definition). */
  transitionMap?: ConnectionTransitionMap;
  disabled?: boolean;
  /** Called with the chosen target status when the user picks a transition. */
  onChange: (nextStatus: string) => void;
};

/**
 * The connection lifecycle status dropdown, shared across the app (the
 * Connections table and the configure modal). It shows the current status as a
 * chip and, when opened, the states the connection can transition to (derived
 * from the per-kind transition map). Selecting a state invokes `onChange`.
 */
export const ConnectionStatusSelect = ({
  status,
  transitionMap,
  disabled,
  onChange,
}: ConnectionStatusSelectProps) => {
  // Only the reachable next states go into the menu; current status is
  // shown in the closed trigger via renderValue.
  const nextStates = getNextStates(transitionMap, status);

  return (
    <FormControl sx={{ width: 'fit-content' }}>
      <ConnectionStyledSelect
        labelId="connection-status-select-label"
        id="connection-status-select"
        disabled={disabled}
        value={status}
        // renderValue owns the closed-trigger chip; menu options are next-states only.
        renderValue={(value) => <ConnectionStateChip status={value as string} actionable={false} />}
        onClick={(event) => event.stopPropagation()}
        onChange={(event) => onChange(event.target.value as string)}
        disableUnderline
        MenuProps={{
          anchorOrigin: { vertical: 'bottom', horizontal: 'left' },
          transformOrigin: { vertical: 'top', horizontal: 'left' },
          // sx owns list padding; no disablePadding conflict.
          sx: { '& .MuiList-root': { padding: '4px' } },
          PaperProps: { square: true },
        }}
      >
        {nextStates.length === 0 && (
          <MenuItem disabled sx={{ padding: '4px 8px' }}>
            No transitions available
          </MenuItem>
        )}
        {nextStates.map((option) => (
          <ConnectionStyledMenuItem value={option} key={option}>
            <ConnectionStateChip status={option} actionable />
          </ConnectionStyledMenuItem>
        ))}
      </ConnectionStyledSelect>
    </FormControl>
  );
};

export default ConnectionStatusSelect;
