/**
 * function used to calculate the zIndex
 * @param {number} p - power of zIndex - directly proportional to `zIndex` (css property) value
 * @returns {string} zIndex
 */
export const ziCalc = (p = 1) => {
  if (p >= 1) {
    // 1 is the least power for zIndex
    let zIndex = "";
    for (let i = 0; i < p; i++) {
      zIndex = zIndex + "9";
    }
    return zIndex;
  }
  `0`;
};
