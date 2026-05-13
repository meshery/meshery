import _ from 'lodash';

function recursivelySearchObjKey(obj, arr, index) {
  if (index === arr.length) {
    return obj;
  }

  const keyName = arr[index];
  if (obj[keyName]) {
    return recursivelySearchObjKey(obj[keyName], arr, index + 1);
  }

  return false;
}

/**
 * CapabilitiesRegistry is a UI component-access guard that wraps the raw provider
 * capabilities payload returned by /api/provider/capabilities.
 *
 * Distinguish between two related but separate concepts:
 *
 *   • Provider capabilities  – the `{ feature, endpoint }[]` list inside the
 *     payload that describes which features the remote provider supports (e.g.
 *     "persist-meshery-patterns").  Access them via `providerCapabilities`.
 *
 *   • Component-access registry – the methods on this class
 *     (`isNavigatorComponentEnabled`, `isHeaderComponentEnabled`, …) that consult
 *     the provider's `restrictedAccess` block to decide whether individual UI
 *     components should be rendered (playground / restricted-mode guards).
 */
export class CapabilitiesRegistry {
  /** Raw payload from /api/provider/capabilities. */
  providerPayload;
  isPlaygroundEnv = false;

  constructor(providerPayload) {
    this.providerPayload = providerPayload;
    this.isPlaygroundEnv = providerPayload?.restrictedAccess?.isMesheryUIRestricted || false;
  }

  /**
   * The list of features the current provider declares it supports.
   * Each entry is a `{ feature: string, endpoint: string }` object.
   * Returns an empty array when the payload is absent or malformed.
   */
  get providerCapabilities() {
    return Array.isArray(this.providerPayload?.capabilities)
      ? this.providerPayload.capabilities
      : [];
  }

  isNavigatorComponentEnabled(navigatorWalker) {
    if (!this.isPlaygroundEnv) {
      return true;
    }

    let walkerArray = ['restrictedAccess', 'allowedComponents', 'navigator', ...navigatorWalker];

    const searchResult = recursivelySearchObjKey(this.providerPayload, walkerArray, 0);
    if (_.isObject(searchResult) && _.isEmpty(searchResult)) {
      return false;
    }
    return searchResult;
  }

  isHeaderComponentEnabled(headerWalker) {
    if (!this.isPlaygroundEnv) {
      return true;
    }

    let walkerArray = ['restrictedAccess', 'allowedComponents', 'header', ...headerWalker];

    const searchResult = recursivelySearchObjKey(this.providerPayload, walkerArray, 0);
    if (_.isObject(searchResult) && _.isEmpty(searchResult)) {
      return false;
    }
    return searchResult;
  }

  isExtensionComponentEnabled(walkerArray) {
    if (!this.isPlaygroundEnv) {
      return true;
    }

    if (!this.providerPayload.extensions?.navigator) {
      return false;
    }

    const navigatorObj = this.providerPayload.extensions.navigator[0]?.allowedTo;

    const searchResult = recursivelySearchObjKey(navigatorObj, walkerArray, 0);
    if (_.isObject(searchResult) && _.isEmpty(searchResult)) {
      return false;
    }
    return searchResult;
  }
}
