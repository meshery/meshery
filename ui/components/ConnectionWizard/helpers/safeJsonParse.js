/**
 * Attempts to parse a given string as JSON. If the parsing fails,
 * it returns an empty object instead of throwing an error.
 *
 * @function
 * @param {string} str - The string to be parsed as JSON.
 * @returns {Object} The parsed JSON object, or an empty object if parsing fails.
 * @example

 */
const safeJsonParse = (str) => {
  try {
    return JSON.parse(str);
  } catch (e) {
    return {};
  }
}

export default safeJsonParse;
