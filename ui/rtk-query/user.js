import { api } from './index';

const userApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getLoggedInUser: builder.query({
      query: () => `user`,
    }),
    getUserById: builder.query({
      query: (id) => `user/profile/${id}`,
    }),
  }),
});

export const { useGetLoggedInUserQuery, useGetUserByIdQuery } = userApi;
