import React from 'react';
import { FormControl, FormControlLabel, FormGroup, FormLabel, Switch } from '@sistent/sistent';

interface AdapterOperation {
  key: string;
  value: string;
  category?: number;
}

interface AdapterAddonSwitchesProps {
  selectedAdapterOps: AdapterOperation[];
  addonSwitchGroup: Record<string, boolean>;
  onSwitchChange: (name: string, checked: boolean, ops: AdapterOperation) => void;
}

/**
 * Renders the per-adapter add-on toggle list (configuration category
 * 2). Extracted from MesheryAdapterPlayComponent's
 * generateAddonSwitches() — same FormControl markup, same per-switch
 * change semantics.
 */
const AdapterAddonSwitches: React.FC<AdapterAddonSwitchesProps> = ({
  selectedAdapterOps,
  addonSwitchGroup,
  onSwitchChange,
}) => {
  if (!selectedAdapterOps.length) return null;

  return (
    <FormControl component="fieldset" style={{ padding: '1rem' }}>
      <FormLabel component="legend">Customize Addons</FormLabel>
      <FormGroup>
        {selectedAdapterOps
          .map((ops) => ({ ...ops, value: ops.value.replace('Add-on:', '') }))
          .sort((ops1, ops2) => ops1.value.localeCompare(ops2.value))
          .map((ops) => (
            <FormControlLabel
              control={
                <Switch
                  color="primary"
                  checked={!!addonSwitchGroup[ops.key]}
                  onChange={(ev) => onSwitchChange(ev.target.name, ev.target.checked, ops)}
                  name={ops.key}
                />
              }
              label={ops.value}
              key={ops.key}
            />
          ))}
      </FormGroup>
    </FormControl>
  );
};

export default AdapterAddonSwitches;
