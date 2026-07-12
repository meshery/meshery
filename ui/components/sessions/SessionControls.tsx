import React from 'react';
import { createPortal } from 'react-dom';
import { InputAdornment, MenuItem, SearchIcon, Select, TextField, Tooltip } from '@sistent/sistent';
import { useSessionHeaderSlot } from './header-slot';

export interface SessionControlsProps {
  /** Containers the session may switch between; from the capabilities endpoint. */
  containers: string[];
  container: string;
  onContainerChange: (container: string) => void;
  query: string;
  onQueryChange: (query: string) => void;
  /** Runs the search. `forward` is false when the user shift-returns, to go back. */
  onSearch: (query: string, forward: boolean) => void;
  /**
   * Whether this session is the one on show. Only the focused session's controls
   * belong in the panel's single header; the others render nothing at all rather
   * than stacking three container selects into one bar.
   */
  focused: boolean;
}

/**
 * A session's container select and search box.
 *
 * They live in the panel's header rather than in the session's own toolbar, which
 * is where they were: a shell attached to one container, searching its own
 * scrollback, is still the panel's subject, and a second row of chrome under the
 * tab strip cost a line of terminal height on a panel the user can make short.
 *
 * With no header to project into — a session embedded on its own, without the
 * shell — they fall back to rendering inline, so an embedded session is not left
 * with no way to switch container.
 */
const SessionControls: React.FC<SessionControlsProps> = ({
  containers,
  container,
  onContainerChange,
  query,
  onQueryChange,
  onSearch,
  focused,
}) => {
  const slot = useSessionHeaderSlot();

  const onKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter') onSearch(query, !event.shiftKey);
  };

  const controls = (
    <>
      {containers.length > 1 && (
        <Tooltip title="Container">
          <Select
            size="small"
            value={container}
            displayEmpty
            onChange={(event) => onContainerChange(String(event.target.value))}
            // The header is a 2rem bar; a labelled form control does not fit in it,
            // and the tooltip carries what the label would have said.
            sx={{ height: '1.75rem', minWidth: '9rem', fontSize: '0.75rem' }}
            renderValue={(value) => (value as string) || 'Container'}
          >
            {containers.map((candidate) => (
              <MenuItem key={candidate} value={candidate} sx={{ fontSize: '0.75rem' }}>
                {candidate}
              </MenuItem>
            ))}
          </Select>
        </Tooltip>
      )}

      <TextField
        size="small"
        placeholder="Search"
        value={query}
        onChange={(event) => onQueryChange(event.target.value)}
        onKeyDown={onKeyDown}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon width={14} height={14} fill="currentColor" />
            </InputAdornment>
          ),
          sx: { height: '1.75rem', fontSize: '0.75rem' },
        }}
        sx={{ width: '12rem' }}
      />
    </>
  );

  if (!slot) return controls;
  if (!focused) return null;
  return createPortal(controls, slot);
};

export default SessionControls;
