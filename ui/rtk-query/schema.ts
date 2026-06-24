import { api, mesheryApiPath } from './index';
import { shouldOverrideExisting } from './utils';

/**
 * RTK queries for schemas present in meshkit
 */
const schemasApi = api.injectEndpoints({
  overrideExisting: shouldOverrideExisting,
  endpoints: (builder) => ({
    getSchema: builder.query({
      query: (queryArg) => mesheryApiPath(`schema/resource/${queryArg.schemaName}`),
    }),
  }),
});

export const { useGetSchemaQuery } = schemasApi;
