/**
 * The debouncing mechanism for calling the function
 * on the timeout reached.
 *
 * @param {Function} func
 * @param {Number} timeout time in milliseconds
 * @returns
 */
export default function debounce(func, timeout = 500) {
  if (typeof func !== 'function') {
    return () => {};
  }
  const delay =
    typeof timeout === 'number' && !Number.isNaN(timeout) && timeout >= 0 ? timeout : 500;
  let timer;
  return function (...args) {
    clearTimeout(timer);
    timer = setTimeout(() => {
      func.apply(this, args);
    }, delay);
  };
}
