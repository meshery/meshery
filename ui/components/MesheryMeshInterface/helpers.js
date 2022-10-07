// @ts-check
// ********************************** TYPE DEFINTIONS **********************************

import { promisifiedDataFetch } from "../../lib/data-fetch";
import { trueRandom } from "../../lib/trueRandom";
import { CustomFieldTemplate } from "./PatternService/RJSFCustomComponents/FieldTemplate";

/**
 * @typedef {Object} OAMDefinition
 * @property {string} kind
 * @property {string} apiVersion
 * @property {Record<string, any>} metadata
 * @property {Record<string, any>} spec
 */

/**
 * @typedef {string} OAMRefSchema
 */

/**
 * @typedef {Object} OAMGenericResponse
 * @property {OAMDefinition} oam_definition
 * @property {OAMRefSchema} oam_ref_schema
 * @property {string} host
 * @property {Record<string, any>} metadata
 */

// ******************************************************************************************

/**
 * getWorkloadDefinitionsForAdapter will fetch workloads for the given
 * adapter from meshery server
 * @param {string} adapter
 * @returns {Promise<Array<OAMGenericResponse>>}
 */
export async function getWorkloadDefinitionsForAdapter(adapter) {
  try {
    const res = await promisifiedDataFetch("/api/oam/workload");
    if (adapter) return res?.filter((el) => el?.metadata?.["adapter.meshery.io/name"] === adapter);
    return res;
  } catch (error) {
    console.error(error);
    return [];
  }
}

/**
 * getTraitDefinitionsForAdapter will fetch tratis for the given
 * adapter from meshery server
 * @param {string} adapter
 * @returns {Promise<Array<OAMGenericResponse>>}
 */
export async function getTraitDefinitionsForAdapter(adapter) {
  try {
    const res = await promisifiedDataFetch("/api/oam/trait?trim=true");

    if (adapter) return res?.filter((el) => el?.metadata?.["adapter.meshery.io/name"] === adapter);
    return res;
  } catch (error) {
    console.error(error);
    return [];
  }
}

/**
 * createWorkloadTraitSets returns an array of workloads and traits object
 * which are interrelated
 * @param {string} adapter
 * @returns {Promise<Array<{
 *  workload: OAMGenericResponse;
 *  traits: Array<OAMGenericResponse>;
 *  type?: string;
 * }>>}
 */
export async function createWorkloadTraitSets(adapter) {
  const workloads = await getWorkloadDefinitionsForAdapter(adapter);
  const traits = await getTraitDefinitionsForAdapter(adapter);

  const sets = [];
  workloads?.forEach((w) => {
    const item = { workload : w, traits : [], type : getPatternServiceType(w) };

    item.traits = traits?.filter((t) => {
      if (Array.isArray(t?.oam_definition?.spec?.appliesToWorkloads))
        return t?.oam_definition?.spec?.appliesToWorkloads?.includes(w?.oam_definition?.metadata?.name);

      return false;
    });

    sets.push(item);
  });

  return sets;
}

/**
 * getPatternServiceName takes in the pattern service metadata and returns
 * the name of the service
 *
 * @param {*} item pattern service component
 * @param {boolean} includeDisplayName if set to true, display name is checked first
 * @returns {string} service name
 */
export function getPatternServiceName(item, includeDisplayName = true) {
  if (includeDisplayName) return item?.metadata?.["display.ui.meshery.io/name"] || item?.oam_definition?.metadata?.name || getPatternAttributeName(item) || "NA";

  return item?.oam_definition?.metadata?.name || "NA";
}

/**
 * getHumanReadablePatternServiceName takes in the pattern service metadata and returns
 * the readable name of the service
 *
 * @param {*} item pattern service component
 * @returns {string} service name
 */
export function getHumanReadablePatternServiceName(item) {
  return (
    item?.metadata?.["display.ui.meshery.io/name"]
  )
}

/**
 * getPatternServiceID takes in the pattern service metadata and returns
 * the ID of the service
 * @param {*} item pattern service component
 * @returns {string | undefined}
 */
export function getPatternServiceID(item) {
  return item?.id;
}

/**
 * getPatternServiceType takes in the pattern service metadata and returns
 * the category of the service
 * @param {*} item pattern service coponent
 * @returns {string | undefined} service name
 */
export function getPatternServiceType(item) {
  return item?.metadata?.["ui.meshery.io/category"];
}

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
 * Capitalises camelcase-string
 *
 * @param {String} text
 * @returns
 */
export function camelCaseToCapitalize(text) {
  if (!text) return null

  return text?.replaceAll(/([A-Z])/g, " $1")?.trim();
}

/**
 * Formats text for prettified view
 *
 * @param {String} text
 * @returns
 */
export function formatString(text) {
  if (!text) return null

  // format string for prettified camelCase
  // @ts-ignore
  let formattedText = text.replaceAll("IP", "Ip");
  formattedText = camelCaseToCapitalize(formattedText),
  formattedText = formattedText.replaceAll("Ip", "IP")
  return formattedText
}

/**
 * The rjsf json schema builder for the ui
 * inplace builds the obj recursively according
 * to the schema provided
 *
 * @param {Record<string, any>} schema The RJSF schema
 * @param {*} obj
 * @returns
 */
function jsonSchemaBuilder(schema, obj) {
  if (!schema) return

  const uiDesc = "ui:description"

  if (schema.type === 'object') {
    for (let key in schema.properties) {
      obj[key] = {};

      // handle percentage for range widget
      if ((schema.properties?.[key]["type"] === 'number' || schema.properties?.[key].type === 'integer')
        && key.toLowerCase().includes("percent")) {
        obj[key]["ui:widget"] = "range"
      }

      jsonSchemaBuilder(schema.properties?.[key], obj[key]);
    }
    return
  }

  if (schema.type === 'array') {
    obj["items"] = {}
    jsonSchemaBuilder(schema.items, obj["items"]);
    return
  }


  // Don't remove the description from additonal Fields Title
  if (!schema?.additionalProperties) {
    obj[uiDesc] = " ";
  }

  if (obj["ui:widget"]) { // if widget is already assigned, don't go over
    return
  }

  if (schema.type === 'boolean') {
    obj["ui:widget"] = "checkbox";
  }

  if (schema.type==='string'&&!schema?.enum) {
    obj["ui:FieldTemplate"] = CustomFieldTemplate;
  }

  if (schema.type === 'number' || schema.type === 'integer') {
    schema["maximum"] = 99999;
    schema["minimum"] = 0;
    obj["ui:widget"] = "updown";
    obj["ui:FieldTemplate"] = CustomFieldTemplate;
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
