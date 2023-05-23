// @ts-check
// ********************************** TYPE DEFINTIONS **********************************

import { trueRandom } from "../../lib/trueRandom";
import { userPromptKeys } from "./PatternService/helper";

/**
 * getPatternAttributeName will take a json schema and will return a pattern
 * attribute name if it does exists, then it returns "NA"
 * @param {*} jsonSchema json schema of the pattern
 * @returns {string} pattern attribute name
 */
export function getPatternAttributeName(jsonSchema) {
  return jsonSchema?._internal?.patternAttributeName || "NA";
}

/**
 * recursiveCleanObject will take an object and will remove all
 * of the "falsy" objects
 * @param {*} obj object that needs to be cleaned
 */
export function recursiveCleanObject(obj) {
  for (const k in obj) {
    if (!obj[k] || typeof obj[k] !== "object") continue;

    recursiveCleanObject(obj[k]);

    if (Object.keys(obj[k]).length === 0) delete obj[k];
  }
  return obj
}

/**
 * recursiveCleanObjectExceptEmptyArray will take an object and will remove all
 * of the "falsy" objects except empty array
 * @param {*} obj object that needs to be cleaned
 */
export function recursiveCleanObjectExceptEmptyArray(obj) {
  for (const k in obj) {
    if (!obj[k] || typeof obj[k] !== "object" || Array.isArray(obj[k])) continue;

    recursiveCleanObjectExceptEmptyArray(obj[k]);

    if (Object.keys(obj[k]).length === 0) delete obj[k];
  }
}

/**
 * createPatternFromConfig will take in the form data
 * and will create a valid pattern from it
 *
 * It will/may also perform some sanitization on the
 * given inputs
 * @param {*} config
 */
export function createPatternFromConfig(config, namespace, partialClean = false) {
  const pattern = {
    name : `pattern-${trueRandom().toString(36).substr(2, 5)}`,
    services : {},
  };
  partialClean ? recursiveCleanObjectExceptEmptyArray(config) : recursiveCleanObject(config);

  Object.keys(config).forEach((key) => {
    // Add it only if the settings are non empty or "true"
    if (config[key].settings) {
      // const name = PascalCaseToKebab(key);
      pattern.services[key] = config[key];

      pattern.services[key].namespace = namespace;
    }
  });

  Object.keys(pattern.services).forEach((key) => {
    // Delete the settings attribute/field if it is set to "true"
    if (pattern.services[key].settings === true) delete pattern.services[key].settings;
  });

  return pattern;
}


/**
 * The rjsf json schema builder for the ui
 * inplace builds the obj recursively according
 * to the schema provided
 *
 * @param {Record<string, any>} schema The RJSF schema
 * @param {*} uiSchema uiSchema
 * @returns
 */
function jsonSchemaBuilder(schema, uiSchema) {
  if (!schema) return

  userPromptKeys.forEach((key) => {
    if (Object.prototype.hasOwnProperty.call(schema, key)) {
      schema[key]?.forEach((item) => {
        jsonSchemaBuilder(item, uiSchema);
      })
    }
  })

  if (schema.type === 'object' || Object.prototype.hasOwnProperty.call(schema, "properties")) { // to handle objects as well as oneof, anyof and allof fields
    for (let key in schema.properties) {
      uiSchema[key] = {};

      // handle percentage for range widget
      if ((schema.properties?.[key]["type"] === 'number' || schema.properties?.[key].type === 'integer')
        && key.toLowerCase().includes("percent")) {
        uiSchema[key]["ui:widget"] = "range"
      }

      jsonSchemaBuilder(schema.properties?.[key], uiSchema[key]);
    }
    return
  }

  if (schema.type === 'array') {
    uiSchema["items"] = {
      "ui:label" : false
    }
    jsonSchemaBuilder(schema.items, uiSchema["items"]);
    return
  }

  if (uiSchema["ui:widget"]) { // if widget is already assigned, don't go over
    return
  }

  if (schema.type === 'boolean') {
    uiSchema["ui:widget"] = "checkbox";
    uiSchema["ui:description"] = "";
  }


  if (schema.type === 'number' || schema.type === 'integer') {
    schema["maximum"] = 99999;
    schema["minimum"] = 0;
  }
}

/**
 * Builds ui schema and sets the required rjsf ui
 * properties
 *
 * @param {Record.<string, any>} schema RJSF json Schema
 * @returns
 */
export function buildUiSchema(schema) {
  const uiSchemaObj = {};
  // 1. Build ui schema
  jsonSchemaBuilder(schema, uiSchemaObj);

  // 2. Set the ordering of the components
  uiSchemaObj["ui:order"] = ["name", "namespace", "label", "annotation", "*"]

  //3. Return the final uiSchema Object
  return uiSchemaObj
}
