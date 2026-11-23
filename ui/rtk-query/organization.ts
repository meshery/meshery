import { useGetOrgsQuery as useSchemasGetOrgsQuery } from '@meshery/schemas/mesheryApi';

// Callers pass `pagesize` like the other list wrappers, but the generated
// getOrgs endpoint (schemas v1beta2) reads only `pageSize`. Any other key is
// dropped before the request is built, and the server falls back to its
// default page.
export const useGetOrgsQuery = (queryArgs, options) =>
  useSchemasGetOrgsQuery(
    {
      page: queryArgs?.page?.toString(),
      pageSize: queryArgs?.pagesize?.toString(),
      search: queryArgs?.search,
      order: queryArgs?.order,
      all: queryArgs?.all,
    },
    options,
  );
