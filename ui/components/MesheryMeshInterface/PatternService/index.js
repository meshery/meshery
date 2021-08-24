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
 *  onSubmit?: Function;
 *  onDelete?: Function;
 *  type: "trait" | "workload"
 *  formData: Record<string, any>;
 *  renderAsTooltip: boolean;
 * }} props
 *
 * @returns
 */
function PatternService({ formData, jsonSchema, onChange, type, onSubmit, onDelete, renderAsTooltip, ...rest }) {
  const ctype = componentType(jsonSchema);

  if (ctype === "rjsf")
    return (
      <RJSF
        formData={formData}
        hideSubmit={type === "trait"}
        hideTitle={type === "workload"}
        jsonSchema={jsonSchema}
        onChange={onChange}
        onSubmit={onSubmit}
        onDelete={onDelete}
        renderAsTooltip={renderAsTooltip}
        {...rest}
      />
    );
  if (ctype === "switch")
    return (
      <Switch
        jsonSchema={jsonSchema}
        onChange={onChange}
        onSubmit={onSubmit}
        onDelete={onDelete}
      />
    );

  return null;
}

export default PatternService;
