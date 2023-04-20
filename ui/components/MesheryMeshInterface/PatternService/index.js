// @ts-check
import React from "react";
import Switch from "./Switch";
import RJSFWrapper from "./RJSF_wrapper";
import { isEmptyObj } from "../../../utils/utils";
import { sortProperties } from "./helper";

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
  if (jsonSchema?.properties) {
    if (Object.keys(jsonSchema?.properties).length)
      return "rjsf";

    return "switch";
  }
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
 *  RJSFWrapperComponent?: any;
 *  RJSFFormChildComponent?: any;
 * }} props
 *
 * @returns
 */
function PatternService({ formData, jsonSchema, onChange, type, onSubmit, onDelete, RJSFWrapperComponent, RJSFFormChildComponent }) {
  const ctype = componentType(jsonSchema);

  const sortedProperties = sortProperties(jsonSchema.properties);
  jsonSchema.properties = sortedProperties;
  if (ctype === "rjsf")
    return (
      <RJSFWrapper
        formData={formData}
        hideSubmit={type === "trait"}
        hideTitle={type === "workload"}
        jsonSchema={jsonSchema}
        onChange={onChange}
        onSubmit={onSubmit}
        onDelete={onDelete}
        RJSFWrapperComponent={RJSFWrapperComponent}
        RJSFFormChildComponent={RJSFFormChildComponent}
      />
    );
  if (ctype === "switch")
    return (
      <Switch
        intialState={!isEmptyObj(formData)}
        jsonSchema={jsonSchema}
        onChange={onChange}
        onSubmit={onSubmit}
        onDelete={onDelete}
      />
    );

  return null;
}

export default PatternService;
