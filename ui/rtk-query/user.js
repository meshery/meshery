import { api } from './index';

const userApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getLoggedInUser: builder.query({
      query: () => `user`,
    }),
    getUserById: builder.query({
      query: (id) => `user/profile/${id}`,
    }),
    getUserPrefs: builder.query({
      query: () => `user/prefs`,
    }),
  }),
});

export const { useGetLoggedInUserQuery, useGetUserByIdQuery, useGetUserPrefsQuery } = userApi;
