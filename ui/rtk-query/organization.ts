import { useGetOrgsQuery as useSchemasGetOrgsQuery } from '@meshery/schemas/mesheryApi';

export const useGetOrgsQuery = (queryArgs, options) =>
  useSchemasGetOrgsQuery(
    {
      page: queryArgs?.page?.toString(),
      pagesize: queryArgs?.pagesize?.toString(),
      search: queryArgs?.search,
      order: queryArgs?.order,
      all: queryArgs?.all,
    },
    options,
  );
