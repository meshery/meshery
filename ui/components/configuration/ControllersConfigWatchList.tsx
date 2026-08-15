import type { ReactNode } from 'react';
import { Box, Button, Checkbox, FormControlLabel, MenuItem, TextField } from '@sistent/sistent';
import {
  INHERIT,
  WATCH_EVENTS,
  WATCH_MODE_OPTIONS,
  fitWidth,
  type WatchList,
} from './controllersConfigForm.shared';

type ControllersConfigWatchListProps = {
  label: ReactNode;
  value: WatchList | undefined;
  disabled: boolean;
  onChange: (next: WatchList | undefined) => void;
};

export default function ControllersConfigWatchList({
  label,
  value,
  disabled,
  onChange,
}: ControllersConfigWatchListProps) {
  const watchMode = !value ? INHERIT : value.whitelist ? 'whitelist' : 'blacklist';
  const whitelist = value?.whitelist ?? [];
  const blacklist = value?.blacklist ?? [];

  return (
    <Box
      sx={{ marginTop: '1rem', display: 'flex', flexDirection: 'column', alignItems: 'flex-start' }}
      data-testid="controllers-config-watch-list"
    >
      {label}
      <TextField
        select
        size="small"
        disabled={disabled}
        value={watchMode}
        aria-label="Watch mode"
        slotProps={{ htmlInput: { 'aria-label': 'Watch mode' } }}
        onChange={(e) => {
          const mode = e.target.value;
          if (mode === INHERIT) onChange(undefined);
          else if (mode === 'whitelist') onChange({ whitelist: [] });
          else onChange({ blacklist: [] });
        }}
        sx={{
          width: fitWidth(...WATCH_MODE_OPTIONS.map((option) => option.label)),
          maxWidth: '100%',
        }}
      >
        {WATCH_MODE_OPTIONS.map((option) => (
          <MenuItem key={option.value} value={option.value}>
            {option.label}
          </MenuItem>
        ))}
      </TextField>

      {watchMode === 'whitelist' && (
        <Box sx={{ marginTop: '1rem', width: '100%' }}>
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
                aria-label={`Resource ${index + 1}`}
                slotProps={{ htmlInput: { 'aria-label': `Resource ${index + 1}` } }}
                sx={{
                  width: fitWidth(row.resource, 'pods.v1. or deployments.v1.apps'),
                  maxWidth: '100%',
                }}
                onChange={(e) => {
                  const rows = [...whitelist];
                  rows[index] = { ...rows[index], resource: e.target.value };
                  onChange({ whitelist: rows });
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
                        if (e.target.checked) events.add(eventType);
                        else events.delete(eventType);
                        rows[index] = { ...rows[index], events: Array.from(events) };
                        onChange({ whitelist: rows });
                      }}
                    />
                  }
                  label={eventType}
                />
              ))}
              <Button
                size="small"
                color="error"
                disabled={disabled}
                onClick={() => onChange({ whitelist: whitelist.filter((_, i) => i !== index) })}
              >
                Remove
              </Button>
            </Box>
          ))}
          <Button
            size="small"
            variant="outlined"
            color="primary"
            disabled={disabled}
            onClick={() =>
              onChange({ whitelist: [...whitelist, { resource: '', events: [...WATCH_EVENTS] }] })
            }
          >
            Add resource
          </Button>
        </Box>
      )}

      {watchMode === 'blacklist' && (
        <TextField
          multiline
          minRows={3}
          size="small"
          disabled={disabled}
          sx={{ marginTop: '1rem', width: '100%' }}
          aria-label="Blacklist resources"
          slotProps={{ htmlInput: { 'aria-label': 'Blacklist resources' } }}
          value={blacklist.join('\n')}
          placeholder={'secrets.v1.\nevents.v1.'}
          onChange={(e) =>
            onChange({
              blacklist: e.target.value
                .split('\n')
                .map((s) => s.trim())
                .filter(Boolean),
            })
          }
        />
      )}
    </Box>
  );
}
