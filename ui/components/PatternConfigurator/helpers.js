/**
 * getHumanReadablePatternServiceName takes in the pattern service metadata and returns
 * the readable name of the service
 *
 * @param {*} item pattern service component
 * @returns {string} service name
 */
 export function getHumanReadablePatternServiceName(item) {
    return item?.metadata?.["display.ui.meshery.io/name"];
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
    if (includeDisplayName)
      return (
        item?.metadata?.["display.ui.meshery.io/name"] ||
        item?.oam_definition?.metadata?.name ||
        getPatternAttributeName(item) ||
        "NA"
      );
  
    return item?.oam_definition?.metadata?.name || "NA";
  }