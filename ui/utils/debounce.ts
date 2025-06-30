/**
 * The debouncing mechanism for calling the function
 * on the timeout reached.
 *
 * @param func - Function to debounce
 * @param timeout - Time in milliseconds
 * @returns Debounced function
 */
export default function debounce<T extends (...args: any[]) => void>(
  func: T,
  timeout = 500
): (...args: Parameters<T>) => void {
  let timer: ReturnType<typeof setTimeout>;
  return (...args: Parameters<T>) => {
    clearTimeout(timer);
    timer = setTimeout(() => {
      func.apply(this, args);
    }, timeout);
  };
}