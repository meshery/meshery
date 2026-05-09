import { useMemo } from 'react';
import { getAllCustomResourceDefinitionsKinds, ResourceMenuConfig } from '../../resources/config';

type Kind = { Kind: string; Count?: number };

export const useResourceOptions = () => {
  const groupOptions = useMemo(
    () => [
      { value: 'all', label: 'All Resources' },
      ...Object.entries(ResourceMenuConfig)
        .filter(([, resources]) => resources.length > 0)
        .map(([category]) => ({
          value: category.toLowerCase(),
          label: category,
        })),
      { value: 'crds', label: 'Custom Resources' },
    ],
    [],
  );

  return groupOptions;
};

const SORT_DIRECTIONS = {
  ASC: 'asc',
  DESC: 'desc',
};

export const useResourceFiltering = (
  kinds: Kind[] | undefined,
  groupBy: string,
  sortDirection: string | null,
) => {
  const filteredAndSortedKinds = useMemo(() => {
    if (!kinds) return [];
    const filteredKinds = filterKindsByGroup(kinds, groupBy);

    if (!sortDirection) return filteredKinds;

    return sortKindsByCount(filteredKinds, sortDirection);
  }, [kinds, groupBy, sortDirection]);

  return filteredAndSortedKinds;
};

const filterKindsByGroup = (kinds: Kind[], groupBy: string): Kind[] => {
  if (groupBy === 'all') return [...kinds];

  if (groupBy === 'crds') {
    const crdKinds = getAllCustomResourceDefinitionsKinds(kinds).map((crd: Kind) => crd.Kind);
    return kinds.filter((item) => crdKinds.includes(item.Kind));
  }

  const categoryKinds =
    ResourceMenuConfig[groupBy.charAt(0).toUpperCase() + groupBy.slice(1)] || [];
  return kinds.filter((item) => categoryKinds.includes(item.Kind));
};

const sortKindsByCount = (kinds: Kind[], direction: string): Kind[] => {
  return [...kinds].sort((a, b) => {
    return direction === SORT_DIRECTIONS.ASC
      ? (a.Count ?? 0) - (b.Count ?? 0)
      : (b.Count ?? 0) - (a.Count ?? 0);
  });
};

export { SORT_DIRECTIONS };
