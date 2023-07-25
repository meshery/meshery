// @ts-check
import React from "react";
import RJSFWrapper from "./RJSF_wrapper";

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
  if (Object.keys(jsonSchema?.properties).length > 0 )
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
  return null;
}

export default PatternService;
