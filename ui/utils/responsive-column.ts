/**
 * updateVisibleColumns Function
 *
 * Purpose:
 * This function calculates which columns should be visible based on the provided screen sizes and the current screen width.
 *
 * Input:
 * @param {Array} colViews - An array of column views, where each element is an array containing the column name and the screen size till which it should be visible.
 * @param {number} width - The current screen width.
 *
 * Output:
 * @returns {Object} - An object representing column visibility, where keys are column names, and values are booleans indicating visibility.
 */

export const updateVisibleColumns = (colViews, width) => {
  // Create a mapping of screen sizes to their minimum widths
  const screenSizeMap = {
    na: 0,
    xs: 585,
    s: 690,
    m: 775,
    l: 915,
    xl: 1140,
  };

  // Determine the minimum width for visibility based on the provided screen size
  const getMinWidthForScreenSize = (screenSize) => {
    return screenSizeMap[screenSize] || 0;
  };

  // Initialize showCols as an empty object
  const showCols = {};

  colViews.forEach((col) => {
    const [columnName, screenSize] = col;
    const minScreenWidth = getMinWidthForScreenSize(screenSize);

    // Determine if the column should be visible based on screen width and rules
    if (screenSize === 'na') {
      showCols[columnName] = false; // "na" columns are not visible in any screen size
    } else if (screenSize === 'xs') {
      showCols[columnName] = true; // "xs" columns are visible in all screen sizes
    } else {
      // For other screen sizes, check if the screen width is greater than or equal to the minimum width
      showCols[columnName] = width >= minScreenWidth;
    }
  });

  return showCols;
};
