import React from 'react';
import { FormControl, MenuItem, Select } from '@sistent/sistent';
import { RANGE_PRESETS } from './time';

interface TimeRangePickerProps {
  durationSec: number;
  onChange: (durationSec: number) => void;
}

/**
 * TimeRangePicker selects the relative time window (ending "now") applied to all
 * panels on the page.
 */
const TimeRangePicker: React.FC<TimeRangePickerProps> = ({ durationSec, onChange }) => (
  <FormControl size="small" sx={{ minWidth: 160 }}>
    <Select
      value={durationSec}
      onChange={(e) => onChange(Number(e.target.value))}
      aria-label="Time range"
    >
      {RANGE_PRESETS.map((preset) => (
        <MenuItem key={preset.durationSec} value={preset.durationSec}>
          {preset.label}
        </MenuItem>
      ))}
    </Select>
  </FormControl>
);

export default TimeRangePicker;
