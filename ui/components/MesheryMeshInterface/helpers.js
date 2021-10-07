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