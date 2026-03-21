import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

const baseQuery = fetchBaseQuery({
  baseUrl: '/api/',
  credentials: 'include', // for pushing client-cookies in all requests to server
  prepareHeaders: (headers) => {
    return headers;
  },
});

const baseQueryWithReauth = async (args, api, extraOptions) => {
  let result = await baseQuery(args, api, extraOptions);

  // Handle redirect to login page
  if (result.error && result.error.status === 401) {
    if (window.location.host.endsWith('3000')) {
      window.location.href = '/user/login'; // for local dev thru node server
    } else {
      window.location.reload(); // for use with Go server
    }
  }

  return result;
};

export const api = createApi({
  reducerPath: 'mesheryApi',
  baseQuery: baseQueryWithReauth,
  endpoints: () => ({}),
});
