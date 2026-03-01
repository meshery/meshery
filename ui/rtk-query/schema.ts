import { mesheryApi } from '@meshery/schemas/dist/mesheryApi';

/**
 * RTK queries for schemas present in meshkit
 */
const schemasApi = mesheryApi.injectEndpoints({
  endpoints: (builder) => ({
    getSchema: builder.query({
      query: (queryArg) => `schema/resource/${queryArg.schemaName}`,
    }),
  }),
});

export const { useGetSchemaQuery } = schemasApi;
