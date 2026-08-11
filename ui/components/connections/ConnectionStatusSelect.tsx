import { FormControl, MenuItem } from '@sistent/sistent';
import { ConnectionStyledSelect } from './styles';
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
  const options = getNextStates(transitionMap, status);
  options.push(status);

  return (
    <FormControl sx={{ width: 'fit-content' }}>
      <ConnectionStyledSelect
        labelId="connection-status-select-label"
        id="connection-status-select"
        disabled={disabled}
        value={status}
        defaultValue={status}
        onClick={(event) => event.stopPropagation()}
        onChange={(event) => onChange(event.target.value as string)}
        disableUnderline
        MenuProps={{
          anchorOrigin: { vertical: 'bottom', horizontal: 'left' },
          transformOrigin: { vertical: 'top', horizontal: 'left' },
          getContentAnchorEl: null,
          MenuListProps: { disablePadding: true },
          PaperProps: { square: true },
        }}
      >
        {options.length === 1 && <MenuItem disabled>No transitions Available</MenuItem>}
        {options.map((option) => (
          <MenuItem
            disabled={option === status}
            style={{
              padding: 0,
              display: option === status ? 'none' : 'flex',
              justifyContent: 'center',
            }}
            value={option}
            key={option}
          >
            <ConnectionStateChip status={option} actionable={option !== status} />
          </MenuItem>
        ))}
      </ConnectionStyledSelect>
    </FormControl>
  );
};

export default ConnectionStatusSelect;
