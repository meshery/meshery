/**
 * The debouncing mechanism for calling the function
 * on the timeout reached.
 *
 * @param {Function} func
 * @param {Number} timeout time in miliseconds
 * @returns
 */
export default function debounce(func, timeout = 500) {
  let timer;
  return (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => {
      func.apply(this, args);
    }, timeout);
  };
}