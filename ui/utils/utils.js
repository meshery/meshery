import { trueRandom } from "../lib/trueRandom";
import jsYaml from "js-yaml";
import { findWorkloadByName } from "./workloadFilter";
import { EVENT_TYPES } from "./Enum";

/**
 * Check if an object is empty
 *
 * @param {Object} obj
 *
 * @returns {Boolean} if obj is empty
 */
export function isEmptyObj(obj) {
  return !obj
    || obj
    && Object.keys(obj).length === 0
    && Object.getPrototypeOf(obj) === Object.prototype;
}

/**
 * Check if array is empty
 *
 * @param {Array} arr
 * @returns {Boolean} if arr is empty
 */
export function isEmptyArr(arr) {
  return arr && arr.length === 0;
}

/**
 * Check if two arrays are equal
 *
 * @param {Array} arr1
 * @param {Array} arr2
 * @param {Boolean} orderMatters
 * @returns
 */
export function isEqualArr(arr1, arr2, orderMatters = true) {
  if (arr1 === arr2) return true;
  if (arr1 == null || arr2 == null) return false;
  if (arr1.length !== arr2.length) return false;

  if (!orderMatters) {
    arr1.sort()
    arr2.sort()
  }

  for (var i = 0; i < arr1.length; ++i) {
    if (arr1[i] !== arr2[i]) return false;
  }
  return true;
}

/**
 * ScrollToTop scrolls the window to top
 *
 * @param {(
 * "auto"
 * |"smooth"
 * |"inherit"
 * |"initial"
 * |"revert"
 * |"unset"
 * )} behavior : scroll-behaviour, see https://developer.mozilla.org/en-US/docs/Web/CSS/scroll-behavior
 */
export function scrollToTop(behavior = "smooth") {
  setTimeout(() => {
    window.scrollTo({
      top : 0,
      left : 0,
      behavior,
    })
  }, 0);
}

/**
 * Generates random Pattern Name with the prefix meshery_
 */
export function randomPatternNameGenerator() {
  return "meshery_" + Math.floor(trueRandom() * 100)
}

/**
 * Returns number of components in Pattern/Application/Filters file
 */
export function getComponentsinFile(file) {
  if (file) {
    try {
      let keys = Object.keys(jsYaml.load(file).services);
      return keys.length;
    } catch (e) {
      if (e.reason?.includes("expected a single document")) {
        return file.split("---").length
      }
    }
  }
  return 0;
}

export function generateValidatePayload(pattern_file, workloadTraitSet) {
  let pattern = jsYaml.loadAll(pattern_file)
  const services = pattern[0]?.services;
  if (!services) {
    return { err : "Services not found in the design" };
  }

  const validationPayloads = {};

  for (const serviceId in services) {
    let valueType;

    let { workload } = findWorkloadByName(services[serviceId].type, workloadTraitSet);

    if (!(workload && workload?.oam_ref_schema)) {
      continue;
    }
    const schema = workload.oam_ref_schema;
    const value = services[serviceId]?.settings;
    if (!value) {
      continue;
    }
    valueType = "JSON";
    const validationPayload = {
      schema,
      value : JSON.stringify(value),
      valueType
    };
    validationPayloads[serviceId] = validationPayload;
  }

  return validationPayloads;
}

export function updateURLs(urlsSet, newUrls, eventType) {
  switch (eventType) {
    case EVENT_TYPES.DELETED:
      newUrls.forEach(url => {
        urlsSet.delete(url);
      })
      break;
    case EVENT_TYPES.ADDED:
    case EVENT_TYPES.MODIFIED:
      newUrls.forEach(url => {
        urlsSet.add(url);
      })

  }
}
