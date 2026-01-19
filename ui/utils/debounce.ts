/**
 * The debouncing mechanism for calling the function
 * on the timeout reached.
 *
 * @param func - The function to debounce
 * @param timeout - Time in milliseconds
 * @returns A debounced version of the function
 */
// eslint-disable-next-line no-unused-vars
export default function debounce<T extends (...args: unknown[]) => void>(
  func: T,
  timeout = 500,
  // eslint-disable-next-line no-unused-vars
): (..._funcArgs: Parameters<T>) => void {
  let timer: ReturnType<typeof setTimeout>;
  return (...funcArgs: Parameters<T>) => {
    clearTimeout(timer);
    timer = setTimeout(() => {
      func(...funcArgs);
    }, timeout);
  };
}
