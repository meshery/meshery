import { api, mesheryApiPath } from './index';

/**
 * RTK queries for schemas present in meshkit
 */
const schemasApi = api.injectEndpoints({
  overrideExisting: module.hot?.status() === 'apply',
  endpoints: (builder) => ({
    getSchema: builder.query({
      query: (queryArg) => mesheryApiPath(`schema/resource/${queryArg.schemaName}`),
    }),
  }),
});

export const { useGetSchemaQuery } = schemasApi;
