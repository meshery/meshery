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
}

/**
 * ConnectionPicker selects which registered Grafana connection the telemetry
 * section is operating against.
 */
const ConnectionPicker: React.FC<ConnectionPickerProps> = ({ connections, value, onChange }) => (
  <FormControl size="small" sx={{ minWidth: 220 }}>
    <InputLabel id="telemetry-connection-label">Grafana connection</InputLabel>
    <Select
      labelId="telemetry-connection-label"
      label="Grafana connection"
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
