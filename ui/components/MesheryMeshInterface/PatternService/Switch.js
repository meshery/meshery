// @ts-check
import React, { useState } from "react";
import { FormControlLabel } from "@material-ui/core";
import MSwitch from "@material-ui/core/Switch";

function Switch({
  intialState, jsonSchema, onChange, onSubmit, onDelete
}) {
  const [isOn, setIsOn] = useState(intialState); // TODO: Hook with meshsync

  return (
    <FormControlLabel
      control={
        <MSwitch
          name={jsonSchema?.title}
          checked={isOn}
          onChange={() => {
            setIsOn((isOn) => {
              const newState = !isOn;

              if (!newState) onDelete?.(!newState); // Trigger this before actually updating the state
              onChange?.(newState, (state) => {
                if (state) onSubmit?.(state); // Trigger this after state update
              });

              return newState;
            });
          }}
        />
      }
      label={jsonSchema?.title}
    />
  );
}

export default Switch;
