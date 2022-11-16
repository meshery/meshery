import _ from "lodash"

function recursivelySearchObjKey(obj, arr, index) {
  if (index == arr.length) {
    return obj;
  }

  const objectKeys = Object.keys(obj);

  // eslint-disable-next-line no-unused-vars
  for (let _key in objectKeys) {
    const keyName = arr[index];
    if (obj[keyName]) {
      return recursivelySearchObjKey(obj[keyName], arr, index + 1);
    }

    return false;
  }
}


export class CapabilitiesRegistry {
  capabilitiesRegistry;
  isPlaygroundEnv = false;

  constructor(capabilitiesRegistry){
    this.capabilitiesRegistry = capabilitiesRegistry;
    this.isPlaygroundEnv = capabilitiesRegistry?.restrictedAccess?.isMesheryUiRestricted || false;
  }

  isNavigatorComponentEnabled(navigatorWalker){
    if (!this.isPlaygroundEnv) {
      return true;
    }

    let walkerArray = ["restrictedAccess", "allowedComponents", "navigator", ...navigatorWalker]

    const searchResult = recursivelySearchObjKey(this.capabilitiesRegistry, walkerArray, 0);
    if (_.isObject(searchResult) && _.isEmpty(searchResult)) {
      return false;
    }
    return searchResult
  }

  isHeaderComponentEnabled(headerWalker){
    if (!this.isPlaygroundEnv) {
      return true;
    }

    let walkerArray = ["restrictedAccess", "allowedComponents", "header", ...headerWalker]

    const searchResult = recursivelySearchObjKey(this.capabilitiesRegistry, walkerArray, 0);
    if (_.isObject(searchResult) && _.isEmpty(searchResult)) {
      return false;
    }
    return searchResult
  }

  isExtensionComponentEnabled(walkerArray) {
    if (!this.isPlaygroundEnv) {
      return true;
    }

    if (!this.capabilitiesRegistry.extensions?.navigator) {
      return false;
    }

    const navigatorObj = this.capabilitiesRegistry.extensions.navigator[0]?.allowedTo;

    const searchResult = recursivelySearchObjKey(navigatorObj, walkerArray, 0);
    if (_.isObject(searchResult) && _.isEmpty(searchResult)) {
      return false;
    }
    return searchResult
  }
}
