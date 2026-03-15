import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

export const api = createApi({
  reducerPath: 'mesheryApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/' }),
  // Use a broad tag type to allow tag usage across injected APIs without
  // forcing every tag string to be declared up-front.
  tagTypes: [] as string[],
  endpoints: () => ({}),
});
