// @ts-check
import React, { useEffect, useState } from "react";
import { FormControlLabel } from "@material-ui/core";
import MSwitch from "@material-ui/core/Switch";

function Switch({ jsonSchema, onChange }) {
  const [isOn, setIsOn] = useState(false); // TODO: Hook with meshsync

  useEffect(() => {
    onChange?.(isOn)
  }, [isOn])

  return (
    <FormControlLabel
      control={<MSwitch name={jsonSchema?.title} checked={isOn} onChange={() => setIsOn((isOn) => !isOn)} />}
      label={jsonSchema?.title}
    />
  );
}

export default Switch;
