import React, { useEffect } from 'react';
import {
  Box,
  CachedIcon,
  FormControl,
  IconButton,
  MenuItem,
  Select,
  Tooltip,
} from '@sistent/sistent';
import { REFRESH_INTERVALS } from './time';

interface RefreshControlProps {
  intervalMs: number;
  onIntervalChange: (ms: number) => void;
  onRefresh: () => void;
}

/**
 * RefreshControl provides a manual refresh button plus an auto-refresh interval
 * selector. It owns the interval timer and calls onRefresh on each tick.
 */
const RefreshControl: React.FC<RefreshControlProps> = ({
  intervalMs,
  onIntervalChange,
  onRefresh,
}) => {
  useEffect(() => {
    if (!intervalMs) return undefined;
    const id = setInterval(onRefresh, intervalMs);
    return () => clearInterval(id);
  }, [intervalMs, onRefresh]);

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
      <Tooltip title="Refresh now">
        <IconButton onClick={onRefresh} aria-label="Refresh now">
          <CachedIcon />
        </IconButton>
      </Tooltip>
      <FormControl size="small" sx={{ minWidth: 88 }}>
        <Select
          value={intervalMs}
          onChange={(e) => onIntervalChange(Number(e.target.value))}
          aria-label="Auto refresh interval"
        >
          {REFRESH_INTERVALS.map((opt) => (
            <MenuItem key={opt.ms} value={opt.ms}>
              {opt.label}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </Box>
  );
};

export default RefreshControl;
