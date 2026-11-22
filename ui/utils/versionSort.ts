import { compose } from 'lodash/fp';
import { WILDCARD_V } from './hooks/useMeshModelComponents';

/**
 * returns the API Version that should be used and is most stable
 * Priority Order v1 > v1beta > v1alpha
 *
 * @param {Array.<String>} versionList
 */
export default function getMostRecentVersion(versionList: unknown): string | undefined {
  if (!Array.isArray(versionList) || versionList.length === 0) return undefined;

  const stableList: string[] = [];
  const alphaList: string[] = [];
  const betaList: string[] = [];

  versionList.forEach((apiVersion) => {
    if (typeof apiVersion !== 'string' || !apiVersion) return;

    const isAlpha = apiVersion.includes('alpha'); // returns true if matches alpha in string
    const isBeta = apiVersion.includes('beta'); // returns true if matches beta in string
    const isStable = !isAlpha && !isBeta && /^v?[0-9]+(\.[0-9]+)*$/i.test(apiVersion);

    isStable && stableList.push(apiVersion);
    isAlpha && alphaList.push(apiVersion);
    isBeta && betaList.push(apiVersion);
  });

  stableList.sort(versionSortComparatorFn).reverse();
  alphaList.sort().reverse();
  betaList.sort().reverse();

  // priority order: stable > beta > alpha
  const fallback = versionList.find((v): v is string => typeof v === 'string' && v.length > 0);
  return stableList?.[0] || betaList?.[0] || alphaList?.[0] || fallback;
}

/**
 * Sorts version in "INCREASING ORDER", apply reverse if wanted to sort in decreasing order
 *
 * @param {string} versionA
 * @param {string} versionB
 * @returns
 */
export function versionSortComparatorFn(versionA: unknown, versionB: unknown): number | undefined {
  if (typeof versionA !== 'string' || typeof versionB !== 'string' || !versionA || !versionB) {
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

  const verA = versionA.split('.');
  const verB = versionB.split('.');
  const maxLen = Math.max(verA.length, verB.length);

  for (let i = 0; i < maxLen; i++) {
    let segA = verA[i];
    let segB = verB[i];

    // index 0 is the start of the version, remove v if present for proper sorting
    if (i === 0) {
      if (segA !== undefined) segA = removeVFromVersion(segA);
      if (segB !== undefined) segB = removeVFromVersion(segB);
    }

    const numA = segA !== undefined ? Number(segA) : 0;
    const numB = segB !== undefined ? Number(segB) : 0;

    const diff = numA - numB;
    if (diff !== 0 && !Number.isNaN(diff)) {
      return diff;
    }
    if (Number.isNaN(diff)) {
      const cmp = String(segA ?? '').localeCompare(String(segB ?? ''));
      if (cmp !== 0) {
        return cmp;
      }
    }
  }

  if (verA.length !== verB.length) {
    return verA.length - verB.length;
  }

  return 0;
}

function removeVFromVersion(version: unknown): string | undefined {
  if (typeof version !== 'string' || !version) return;
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
export const sortByVersionInDecreasingOrder = (versions: unknown): string[] | undefined => {
  if (!Array.isArray(versions)) {
    return;
  }

  const validVersions = versions.filter((v): v is string => typeof v === 'string' && v.length > 0);

  // Count only concrete versions: a WILDCARD_V already in the input must not
  // trigger a second prepended copy, and it must stay first in the result.
  const concreteVersions = validVersions.filter((v) => v !== WILDCARD_V);
  const hasWildcard = concreteVersions.length !== validVersions.length;
  if (concreteVersions.length > 1 || hasWildcard) {
    return [WILDCARD_V, ...[...concreteVersions].sort(versionSortComparatorFn).reverse()];
  }

  return [...concreteVersions].sort(versionSortComparatorFn).reverse();
};

/**
 * Get Greater version between the two versions passed
 *
 * @param {string} v1
 * @param {string} v2
 */
export function getGreaterVersion(v1: unknown, v2: unknown): unknown {
  const comparatorResult = versionSortComparatorFn(v1, v2);

  if (comparatorResult === undefined) {
    if (typeof v1 === 'string' && v1.length > 0) return v1;
    if (typeof v2 === 'string' && v2.length > 0) return v2;
    return typeof v1 === 'string' ? v1 : v2;
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
function groupModelsByVersion(models: unknown): Array<any> {
  if (!Array.isArray(models)) {
    return [];
  }

  let modelMap: Record<string, any> = {};
  models.forEach((model) => {
    if (!model || typeof model !== 'object' || !('name' in model)) return;
    const modelName = (model as any).name;
    if (typeof modelName !== 'string' || !modelName) return;

    const rawVersion = (model as any).version;
    const versionList = Array.isArray(rawVersion)
      ? [...new Set(rawVersion)]
      : rawVersion !== undefined
        ? [rawVersion]
        : [];

    const modelMapCurr = modelMap[modelName];
    if (modelMapCurr) {
      modelMap[modelName] = {
        ...modelMapCurr,
        version: [...new Set([...modelMapCurr.version, ...versionList])],
      };
    } else {
      modelMap[modelName] = {
        ...model,
        version: versionList,
      };
    }
  });

  return Object.values(modelMap);
}

function sortVersionsInModel(models: unknown) {
  if (!Array.isArray(models)) {
    return [];
  }
  return [...models].map((model) => ({
    ...model,
    version: sortByVersionInDecreasingOrder((model as any)?.version),
  }));
}

export const sortAndGroupVersionsInModel = compose(sortVersionsInModel, groupModelsByVersion);
