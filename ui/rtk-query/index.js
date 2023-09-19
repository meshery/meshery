import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

export const api = createApi({
  reducerPath: 'mesheryApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/' }),
  endpoints: () => ({
  }),
})