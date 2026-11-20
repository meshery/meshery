import _ from 'lodash';

function recursivelySearchObjKey(obj, arr, index) {
  if (index === arr.length) {
    return obj;
  }

  if (!obj || typeof obj !== 'object') {
    return false;
  }

  const keyName = arr[index];
  if (obj[keyName]) {
    return recursivelySearchObjKey(obj[keyName], arr, index + 1);
  }

  return false;
}

/**
 * Wraps the raw /api/provider/capabilities payload and answers provider-driven
 * UI access questions.
 *
 * "Registry" in Meshery UI refers to the Model/Relationship Registry under
 * ui/components/registry, so this helper deliberately avoids that name.
 */
export class ProviderUiAccessControl {
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

    if (!Array.isArray(navigatorWalker) || navigatorWalker.length === 0) {
      return false;
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

    if (!Array.isArray(headerWalker) || headerWalker.length === 0) {
      return false;
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

    if (!Array.isArray(walkerArray) || walkerArray.length === 0) {
      return false;
    }

    if (
      !Array.isArray(this.providerPayload?.extensions?.navigator) ||
      this.providerPayload.extensions.navigator.length === 0
    ) {
      return false;
    }

    const navigatorObj = this.providerPayload.extensions.navigator[0]?.allowedTo;
    if (!navigatorObj || typeof navigatorObj !== 'object') {
      return false;
    }

    const searchResult = recursivelySearchObjKey(navigatorObj, walkerArray, 0);
    if (_.isObject(searchResult) && _.isEmpty(searchResult)) {
      return false;
    }
    return searchResult;
  }
}
