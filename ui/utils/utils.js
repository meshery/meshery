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

/**
 * ScrollToTop scrolls the window to top
 *
 * @param {(
 * "auto"
 * |"smooth"
 * |"inherit"
 * |"initial"
 * |"revert"
 * |"unset"
 * )} behavior : scroll-behaviour, see https://developer.mozilla.org/en-US/docs/Web/CSS/scroll-behavior
 */
export function scrollToTop(behavior='smooth') {
  setTimeout(() => {
    window.scrollTo({
      top: 0,
      left: 0,
      behavior,
    });
  }, 0);
}