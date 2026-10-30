import { compose } from 'lodash/fp';
import { WILDCARD_V } from './hooks/useMeshModelComponents';

/**
 * returns the API Version that should be used and is most stable
 * Priority Order v1 > v1beta > v1alpha
 *
 * @param {Array.<String>} versionList
 */
export default function getMostRecentVersion(versionList) {
  if (!versionList || !Array.isArray(versionList) || versionList.length === 0) return;

  const stableList: string[] = [];
  const alphaList: string[] = [];
  const betaList: string[] = [];

  versionList.forEach((apiVersion) => {
    if (!apiVersion || typeof apiVersion !== 'string') return;
    const isStable = /^v?[0-9]+(\.[0-9]+)*$/.test(apiVersion);
    const isAlpha = apiVersion.includes('alpha');
    const isBeta = apiVersion.includes('beta');

    if (isStable) {
      stableList.push(apiVersion);
    } else if (isBeta) {
      betaList.push(apiVersion);
    } else if (isAlpha) {
      alphaList.push(apiVersion);
    }
  });

  stableList.sort(versionSortComparatorFn).reverse();
  alphaList.sort(versionSortComparatorFn).reverse();
  betaList.sort(versionSortComparatorFn).reverse();

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
  if (versionA === undefined || versionB === undefined || versionA === null || versionB === null) {
    return;
  }

  if (versionA === versionB) {
    return 0;
  }

  if (versionA === WILDCARD_V) {
    return -1;
  }
  if (versionB === WILDCARD_V) {
    return 1;
  }

  const verA = String(versionA).split('.');
  const verB = String(versionB).split('.');
  const maxLen = Math.max(verA.length, verB.length);

  for (let i = 0; i < maxLen; i++) {
    let segA = verA[i];
    let segB = verB[i];

    if (segA === undefined) return -1;
    if (segB === undefined) return 1;

    // index 0 is the start of the version, remove v if present for proper sorting
    if (i === 0) {
      segA = removeVFromVersion(segA);
      segB = removeVFromVersion(segB);
    }

    const numA = parseInt(segA, 10);
    const numB = parseInt(segB, 10);

    if (!isNaN(numA) && !isNaN(numB)) {
      if (numA !== numB) {
        return numA - numB;
      }
      if (segA !== segB) {
        return segA.localeCompare(segB);
      }
    } else {
      const cmp = String(segA).localeCompare(String(segB));
      if (cmp !== 0) {
        return cmp;
      }
    }
  }

  return 0;
}

function removeVFromVersion(version) {
  if (!version) return;
  if (version.startsWith('v')) {
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
export const sortByVersionInDecreasingOrder = (versions) => {
  if (!versions || !Array.isArray(versions)) {
    return;
  }

  const uniqueVersions = [...new Set(versions.filter((v) => v !== WILDCARD_V))];

  // add wildcard only in the case of multiple distinct versions
  let wildCardV: string[] = [];
  if (uniqueVersions.length > 1) {
    wildCardV = [WILDCARD_V];
  }

  return [...wildCardV, ...[...uniqueVersions].sort(versionSortComparatorFn).reverse()];
};

/**
 * Get Greater version between the two versions passed
 *
 * @param {string} v1
 * @param {string} v2
 */
export function getGreaterVersion(v1, v2) {
  const comparatorResult = versionSortComparatorFn(v1, v2);

  if (comparatorResult === undefined) {
    return v1 || v2;
  }

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
  if (!models || !Array.isArray(models)) {
    return [];
  }

  const modelMap: Record<string, any> = {};
  models.forEach((model) => {
    if (!model || !model.name) return;
    let modelVersions: string[] = [];
    if (Array.isArray(model.version)) {
      modelVersions = model.version;
    } else if (model.version !== undefined) {
      modelVersions = [model.version];
    }
    const existing = modelMap[model.name];
    if (existing) {
      modelMap[model.name] = {
        ...existing,
        version: [...new Set([...existing.version, ...modelVersions])],
      };
    } else {
      modelMap[model.name] = {
        ...model,
        version: [...modelVersions],
      };
    }
  });

  return Object.values(modelMap);
}

function sortVersionsInModel(models) {
  return [...models].map((model) => ({
    ...model,
    version: sortByVersionInDecreasingOrder(model.version),
  }));
}

export const sortAndGroupVersionsInModel = compose(sortVersionsInModel, groupModelsByVersion);
