// @ts-check
// ********************************** TYPE DEFINTIONS **********************************

import { promisifiedDataFetch } from "../../lib/data-fetch";
import PascalCaseToKebab from "../../utils/PascalCaseToKebab";

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
    const res = await promisifiedDataFetch("/api/oam/workload?trim=true");

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
    const item = { workload : w, traits : [], type : getPatternServiceType(w?.metadata) };

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
 * @param {*} item pattern service component
 * @param {boolean} includeDisplayName if set to true, display name is checked first
 * @returns {string} service name
 */
export function getPatternServiceName(item, includeDisplayName = true) {
  if (includeDisplayName) return item?.metadata?.["display.ui.meshery.io/name"] || item?.oam_definition?.metadata?.name || "NA";

  return item?.oam_definition?.metadata?.name || "NA";
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
  return item?.metadata?.["ui.meshery.io/category"]
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
    name : `pattern-${Math.random().toString(36).substr(2, 5)}`,
    services : {},
  };

  partialClean ? recursiveCleanObjectExceptEmptyArray(config) : recursiveCleanObject(config);

  Object.keys(config).forEach((key) => {
    // Add it only if the settings are non empty or "true"
    if (config[key].settings) {
      const name = PascalCaseToKebab(key);
      pattern.services[name] = config[key];

      pattern.services[name].type = key;
      pattern.services[name].namespace = namespace;
    }
  });

  Object.keys(pattern.services).forEach((key) => {
    // Delete the settings attribute/field if it is set to "true"
    if (pattern.services[key].settings === true) delete pattern.services[key].settings;
  });

  return pattern;
}