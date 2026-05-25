import { useMemo } from 'react';
import { getAllCustomResourceDefinitionsKinds, ResourceMenuConfig } from '../../resources/config';

export type ResourceKind = { Kind: string; Model?: string; Count?: number };

type ResourceGroupOption = { value: string; label: string };

export const DEFAULT_GROUP_BY = 'all';
export const SORT_DIRECTIONS = {
  ASC: 'asc',
  DESC: 'desc',
} as const;

export type SortDirection = (typeof SORT_DIRECTIONS)[keyof typeof SORT_DIRECTIONS];

export const useResourceOptions = (): ResourceGroupOption[] => {
  const groupOptions = useMemo(
    () => [
      { value: DEFAULT_GROUP_BY, label: 'All Resources' },
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

export const useResourceFiltering = (
  kinds: ResourceKind[] | undefined,
  groupBy: string,
  sortDirection: SortDirection | null,
) => {
  const filteredAndSortedKinds = useMemo(() => {
    if (!kinds) {
      return [];
    }

    const filteredKinds = filterKindsByGroup(kinds, groupBy);

    if (!sortDirection) {
      return filteredKinds;
    }

    return sortKindsByCount(filteredKinds, sortDirection);
  }, [kinds, groupBy, sortDirection]);

  return filteredAndSortedKinds;
};

const filterKindsByGroup = (kinds: ResourceKind[], groupBy: string): ResourceKind[] => {
  if (groupBy === DEFAULT_GROUP_BY) {
    return [...kinds];
  }

  if (groupBy === 'crds') {
    const crdKinds = new Set(
      getAllCustomResourceDefinitionsKinds(kinds).map(
        (customResource: ResourceKind) => customResource.Kind,
      ),
    );

    return kinds.filter((item) => crdKinds.has(item.Kind));
  }

  const categoryKey = groupBy.charAt(0).toUpperCase() + groupBy.slice(1);
  const categoryKinds = ResourceMenuConfig[categoryKey] || [];

  return kinds.filter((item) => categoryKinds.includes(item.Kind));
};

const sortKindsByCount = (kinds: ResourceKind[], direction: SortDirection): ResourceKind[] => {
  return [...kinds].sort((a, b) => {
    return direction === SORT_DIRECTIONS.ASC
      ? (a.Count ?? 0) - (b.Count ?? 0)
      : (b.Count ?? 0) - (a.Count ?? 0);
  });
};
