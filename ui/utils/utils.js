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