/**
 * updateVisibleColumns Function
 *
 * Calculates column visibility based on screen size and current screen width.
 */
export const updateVisibleColumns = (colViews?: unknown, width?: number | unknown) => {
  if (!Array.isArray(colViews)) {
    return {};
  }

  const screenSizeMap: Record<string, number> = Object.assign(Object.create(null), {
    na: 0,
    xs: 585,
    s: 690,
    m: 775,
    l: 915,
    xl: 1140,
  });

  const safeWidth = typeof width === 'number' && Number.isFinite(width) ? Math.max(0, width) : 0;
  const showCols: Record<string, boolean> = Object.create(null);

  colViews.forEach((col) => {
    if (!Array.isArray(col) || col.length !== 2) {
      return;
    }

    const [columnName, screenSize] = col;
    if (
      typeof columnName !== 'string' ||
      columnName.trim() === '' ||
      typeof screenSize !== 'string' ||
      screenSize.trim() === ''
    ) {
      return;
    }

    const normalizedScreenSize = screenSize.trim();
    const minScreenWidth = Object.prototype.hasOwnProperty.call(screenSizeMap, normalizedScreenSize)
      ? screenSizeMap[normalizedScreenSize]
      : 0;

    if (normalizedScreenSize === 'na') {
      showCols[columnName] = false;
    } else if (normalizedScreenSize === 'xs') {
      showCols[columnName] = true;
    } else {
      showCols[columnName] = safeWidth >= minScreenWidth;
    }
  });

  return showCols;
};

export const getResponsiveColumnVisibility = (
  columnNames?: unknown,
  colViews?: unknown,
  width?: number | unknown,
) => {
  if (!Array.isArray(columnNames)) {
    return {};
  }

  const showCols = updateVisibleColumns(colViews, width);

  return Object.fromEntries(
    columnNames
      .filter(
        (columnName): columnName is string =>
          typeof columnName === 'string' && columnName.trim() !== '',
      )
      .map((columnName) => [
        columnName,
        Object.prototype.hasOwnProperty.call(showCols, columnName)
          ? showCols[columnName]
          : undefined,
      ]),
  );
};
