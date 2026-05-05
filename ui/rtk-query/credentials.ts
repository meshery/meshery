import { api, mesheryApiPath } from './index';

const TAGS = {
  CREDENTIALS: 'credentials',
};
const credentialsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getCredentials: builder.query({
      query: () => ({
        url: mesheryApiPath('integrations/credentials'),
        method: 'GET',
      }),
      providesTags: [TAGS.CREDENTIALS],
    }),
    getCredentialById: builder.query({
      query: (credentialId) => mesheryApiPath(`integrations/credentials/${credentialId}`),
      providesTags: (result, error, id) => [{ type: TAGS.CREDENTIALS, id }],
    }),
    createCredential: builder.mutation({
      query: (credential) => ({
        url: mesheryApiPath('integrations/credentials'),
        method: 'POST',
        body: credential,
      }),
      invalidatesTags: [TAGS.CREDENTIALS],
    }),
    updateCredential: builder.mutation({
      query: (credential) => ({
        url: mesheryApiPath('integrations/credentials'),
        method: 'PUT',
        body: credential,
      }),
      invalidatesTags: [TAGS.CREDENTIALS],
    }),
    deleteCredential: builder.mutation({
      query: (credentialId) => ({
        url: mesheryApiPath(`integrations/credentials?credential_id=${credentialId}`),
        method: 'DELETE',
      }),
      invalidatesTags: [TAGS.CREDENTIALS],
    }),
  }),
  overrideExisting: true,
});

export const {
  useGetCredentialsQuery,
  useGetCredentialByIdQuery,
  useLazyGetCredentialByIdQuery,
  useCreateCredentialMutation,
  useUpdateCredentialMutation,
  useDeleteCredentialMutation,
} = credentialsApi;
