import React from 'react';
import { FormControl, InputLabel, MenuItem, Select } from '@sistent/sistent';

export interface TelemetryConnection {
  id: string;
  name: string;
  kind: string;
  metadata?: Record<string, any>;
}

interface ConnectionPickerProps {
  connections: TelemetryConnection[];
  value: string;
  onChange: (connectionID: string) => void;
  // Label for the picker, e.g. "Grafana connection" / "Prometheus connection".
  label?: string;
}

/**
 * ConnectionPicker selects which registered telemetry connection the section is
 * operating against. Product-agnostic — pass a label per module.
 */
const ConnectionPicker: React.FC<ConnectionPickerProps> = ({
  connections,
  value,
  onChange,
  label = 'Connection',
}) => (
  <FormControl size="small" sx={{ minWidth: 220 }}>
    <InputLabel id="telemetry-connection-label">{label}</InputLabel>
    <Select
      labelId="telemetry-connection-label"
      label={label}
      value={value}
      onChange={(e) => onChange(String(e.target.value))}
    >
      {connections.map((conn) => (
        <MenuItem key={conn.id} value={conn.id}>
          {conn.name || conn.metadata?.url || conn.id}
        </MenuItem>
      ))}
    </Select>
  </FormControl>
);

export default ConnectionPicker;
