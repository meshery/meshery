import { api, mesheryApiPath } from './index';

/**
 * RTK queries for schemas present in meshkit
 */
const schemasApi = api.injectEndpoints({
  overrideExisting: true,
  endpoints: (builder) => ({
    getSchema: builder.query({
      query: (queryArg) => mesheryApiPath(`schema/resource/${queryArg.schemaName}`),
    }),
  }),
});

export const { useGetSchemaQuery } = schemasApi;
