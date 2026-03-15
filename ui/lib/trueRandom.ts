export function trueRandom() {
  const cryptoObj = globalThis.crypto;
  const getRandomValues = cryptoObj?.getRandomValues?.bind(cryptoObj);
  if (getRandomValues) {
    return getRandomValues(new Uint32Array(1))[0] / 2 ** 32;
  }

  // Fallback for SSR / older runtimes
  return Math.random();
}
