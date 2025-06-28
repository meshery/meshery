import { useState, useCallback, useEffect } from 'react';
import { useLazyGetMeshModelsQuery } from '@/rtk-query/meshModel';
import { removeDuplicateVersions } from '@/components/Settings/Registry/helper';

export const useModelSelection = () => {
  const [models, setModels] = useState([]);
  const [searchText, setSearchText] = useState('');
  const [page, setPage] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [isSearching, setIsSearching] = useState(false);

  const [getMeshModelsData, { isLoading, isFetching }] = useLazyGetMeshModelsQuery();

  const fetchModels = useCallback(
    async (resetData = false) => {
      try {
        const currentPage = resetData ? 0 : page;
        const response = await getMeshModelsData({
          params: {
            page: currentPage,
            pagesize: searchText ? 'all' : 25,
            components: false,
            relationships: false,
            search: searchText || '',
          },
        });

        if (response.data?.models) {
          const newModels = removeDuplicateVersions(response.data.models);

          if (resetData || searchText) {
            setModels(newModels);

            const totalCount = response.data.total_count || 0;
            setHasMore(newModels.length < totalCount && !searchText);
          } else {
            setModels((prev) => {
              const updatedModels = [...prev, ...newModels];

              const totalCount = response.data.total_count || 0;
              setHasMore(updatedModels.length < totalCount);
              return updatedModels;
            });
          }
        }
      } catch (error) {
        console.error('Failed to fetch models:', error);
      }
    },
    [getMeshModelsData, page, searchText],
  );

  const handleSearch = useCallback((query) => {
    setSearchText(query);
    setPage(0);
    setIsSearching(true);
  }, []);

  const loadNextPage = useCallback(() => {
    if (!isLoading && !isFetching && hasMore && !searchText) {
      setPage((prev) => prev + 1);
    }
  }, [isLoading, isFetching, hasMore, searchText]);

  useEffect(() => {
    if (searchText || page === 0) {
      fetchModels(true);
      setIsSearching(false);
    } else if (page > 0) {
      fetchModels(false);
    }
  }, [searchText, page]);

  //-> initial
  useEffect(() => {
    fetchModels(true);
  }, []);

  return {
    models,
    searchText,
    isLoading,
    isFetching,
    hasMore,
    isSearching,
    handleSearch,
    loadNextPage,
  };
};
