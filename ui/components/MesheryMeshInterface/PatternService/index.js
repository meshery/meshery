// @ts-check
import React from "react";
import Switch from "./Switch";
import RJSF from "./RJSF";

/**
 * componentType takes in json schema and returns the type
 * of the component that should be used for that schema
 *
 * This allows to use custom components along with the
 * react json form schema component
 * @param {Record<string, any>} jsonSchema
 * @return {"rjsf" | "switch"}
 */
function componentType(jsonSchema) {
  if (Object.keys(jsonSchema?.properties).length) return "rjsf";
  return "switch";
}

/**
 * PatternService returns a component for the given jsonSchema
 * @param {{
 *  jsonSchema: Record<string, any>;
 *  onChange: Function;
 *  onSubmit: Function;
 *  type: "trait" | "workload"
 * }} props
 * @returns
 */
function PatternService({ jsonSchema, onChange, type, onSubmit }) {
  const ctype = componentType(jsonSchema);

  if (ctype === "rjsf")
    return (
      <RJSF
        hideSubmit={type === "trait"}
        hideTitle={type === "workload"}
        jsonSchema={jsonSchema}
        onChange={onChange}
        onSubmit={onSubmit}
      />
    );
  if (ctype === "switch")
    return (
      <Switch
        jsonSchema={jsonSchema}
        onChange={onChange}
        onSubmit={onSubmit}
      />
    );

  return null;
}

export default PatternService;
