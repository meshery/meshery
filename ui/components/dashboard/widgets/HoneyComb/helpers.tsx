export const getGridColumnsCount = (hexagonSide: number, containerWidth: number): number => {
  if (!hexagonSide || !containerWidth || hexagonSide <= 0 || containerWidth <= 0) {
    return 1;
  }
  const hexagonWidth = Math.sqrt(3) * hexagonSide;
  const columns = Math.floor(containerWidth / hexagonWidth);
  return Math.max(1, columns);
};

export const getRowSize = (hexagonSide: number): number => Math.max(1, hexagonSide / 2);

export const getColumnSize = (hexagonSide: number): number =>
  Math.max(1, (Math.sqrt(3) * hexagonSide) / 4);
