export const getGridColumnsCount = (hexagonSide, containerWidth) => {
  if (!hexagonSide || !containerWidth || hexagonSide <= 0 || containerWidth <= 0) {
    return 1; // Return default value for invalid inputs
  }
  const hexagonWidth = Math.sqrt(3) * hexagonSide;
  const columns = Math.floor(containerWidth / hexagonWidth);
  return Math.max(1, columns); // Ensure at least 1 column
};

export const getRowSize = (hexagonSide) => {
  return Math.max(1, hexagonSide / 2); // Ensure positive value
};

export const getColumnSize = (hexagonSide) => {
  return Math.max(1, (Math.sqrt(3) * hexagonSide) / 4); // Ensure positive value
};
