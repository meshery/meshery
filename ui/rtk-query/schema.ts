import { mesheryBaseApi } from '@meshery/schemas/dist/api';

/**
 * RTK queries for schemas present in meshkit
 */
const schemasApi = mesheryBaseApi.injectEndpoints({
  endpoints: (builder) => ({
    getSchema: builder.query({
      query: (queryArg) => `schema/resource/${queryArg.schemaName}`,
    }),
  }),
});

export const { useGetSchemaQuery } = schemasApi;
