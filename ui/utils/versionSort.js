import { compose } from "lodash/fp";
import { WILDCARD_V } from "./hooks/useMeshModelComponents";

/**
 * returns the API Version that should be used and is most stable
 * Priority Order v1 > v1beta > v1alpha
 *
 * @param {Array.<String>} versionList
 */
export default function getMostRecentVersion(versionList) {
  if (!versionList) return;

  const stableList = [];
  const alphaList = [];
  const betaList = [];

  versionList.forEach(apiVersion => {
    const isStable = /^v[0-9]$/.test(apiVersion); // returns true if matches v1-v9
    const isAlpha = apiVersion.includes("alpha"); // returns true if matches alpha in string
    const isBeta = apiVersion.includes("beta"); // returns true if matches beta in string

    isStable && stableList.push(apiVersion);
    isAlpha && alphaList.push(apiVersion);
    isBeta && betaList.push(apiVersion);
  });

  stableList.sort().reverse();
  alphaList.sort().reverse();
  betaList.sort().reverse();

  // priority order: stable > beta > alpha
  return stableList?.[0] || betaList?.[0] || alphaList?.[0] || versionList?.[0];
}

/**
 * Sorts version in "INCREASING ORDER", apply reverse if wanted to sort in descreasing order
 *
 * @param {string} versionA
 * @param {string} versionB
 * @returns
 */
export function versionSortComparatorFn(versionA, versionB) {
  if (!versionA || !versionB) {
    return;
  }

  if (versionA === WILDCARD_V || versionB === WILDCARD_V) { // wildcard support
    return -1;
  }

  const verA = versionA.split(".");
  const verB = versionB.split(".");


  for (let i = 0; i < verA.length && i < verB.length; i++) {
    let vA = verA[i];
    let vB = verB[i];
    // index 0 is the start of the version, remove v if present for proper sorting
    if (i == 0) {
      vA = removeVFromVersion(vA);
      vB = removeVFromVersion(vB);
    }
    // move to next comparison
    if (vA - vB === 0) {
      continue;
    }
    return vA - vB;
  }
}

function removeVFromVersion(version) {
  if (!version) return;
  if (version.startsWith("v")) {
    return version.substring(1);
  }
  return version;
}

/**
 * usually when you sort the version by string, the version 10.0.0 < 2.0.0, because of string sort,
 * ideally, it should be 10.0.0 > 2.0.0
 *
 * [ "2.2.1", "2.10.11", "10.1.2", "10.1.1" ] using this comparator function returns [ '10.1.2', '10.1.1', '2.10.11', '2.2.1' ]
 * @param {Array.<string>} versions
 * @returns Versions sorted in decreasing order
 */
export const sortByVersionInDecreasingOrder = versions => {
  if (!versions) {
    return;
  }

  // add wildcard only in the case of multiple distinct versions
  let wildCardV = [];
  if (versions.length > 1) {
    wildCardV = [WILDCARD_V]
  }

  return [...wildCardV, ...[...versions].sort(versionSortComparatorFn).reverse()];
};

/**
 * Get Greater version between the two versions passed
 *
 * @param {string} v1
 * @param {string} v2
 */
export function getGreaterVersion(v1, v2) {
  const comparatorResult = versionSortComparatorFn(v1, v2);

  if (comparatorResult >= 0) {
    return v1;
  }

  return v2;
}


/*
    MeshModel Specific
*/

/**
 *
 * @param {Array} models
 * @returns {Array} the de-duplicated models array with all the versions available
 */
function groupModelsByVersion(models) {
  if (!models) {
    return [];
  }

  let modelMap = {};
  models.forEach(model => {
    const modelMapCurr = modelMap[model.name];
    if (modelMapCurr) {
      modelMap[model.name].version = [
        ...new Set([...modelMapCurr.version, model.version])
      ]; // remove duplicate entries for version
    } else {
      model.version = [model.version]; // the version is a string, to derive consistency the version is changed to array
      modelMap[model.name] = model;
    }
  });

  return Object.values(modelMap);
}

function sortVersionsInModel(models) {
  return [...models].map(model => ({
    ...model,
    version : sortByVersionInDecreasingOrder(model.version)
  }));
}

export const sortAndGroupVersionsInModel = compose(
  sortVersionsInModel,
  groupModelsByVersion
);
