import {
  mesheryApi,
  useGetCredentialByIdQuery as useSchemasGetCredentialByIdQuery,
  useGetUserCredentialsQuery as useSchemasGetUserCredentialsQuery,
  useSaveUserCredentialMutation as useSchemasSaveUserCredentialMutation,
  useUpdateUserCredentialMutation as useSchemasUpdateUserCredentialMutation,
} from '@meshery/schemas/mesheryApi';
import { api, mesheryApiPath } from './index';

// The tag type registered on the shared `mesheryApi` (see
// @meshery/schemas/mesheryApi), which the schemas credential endpoints provide
// and invalidate. A bare 'credentials' is not a registered tag type, so
// invalidating it silently invalidates nothing and the list never refetches.
const TAGS = {
  CREDENTIALS: 'credential_credentials',
};

// `deleteCredential` is the one credential endpoint still declared here: the
// schemas-generated `deleteUserCredential` sends `?credentialId=`, but the
// server reads `credential_id` (DeleteUserCredential in
// server/handlers/credentials_handlers.go), so the generated endpoint would
// resolve a nil UUID and delete nothing. Everything else delegates to schemas.
const credentialsApi = api.injectEndpoints({
  endpoints: (builder) => ({
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

export const { useDeleteCredentialMutation } = credentialsApi;

// Backed by the schemas-generated `getUserCredentials` (GET
// /api/integrations/credentials). Callers pass no list args, so every schemas
// param stays undefined and the request matches what this module used to build.
// queryArg is forwarded as-is rather than defaulted to `{}`, because RTK derives
// the cache key from it: `{}` would key separately from a plain schemas call,
// and from `connection.ts`'s wrapper over the same endpoint.
export const useGetCredentialsQuery = (queryArg?: undefined, options?: object) =>
  useSchemasGetUserCredentialsQuery(queryArg, options);

// Callers pass a bare id; the schemas endpoints take `{ credentialId }`.
export const useGetCredentialByIdQuery = (credentialId: string, options?: object) =>
  useSchemasGetCredentialByIdQuery({ credentialId }, options);

export const useLazyGetCredentialByIdQuery = () => {
  const [trigger, ...rest] = mesheryApi.endpoints.getCredentialById.useLazyQuery();
  const wrappedTrigger = (credentialId: string, preferCacheValue?: boolean) =>
    trigger({ credentialId }, preferCacheValue);
  return [wrappedTrigger, ...rest] as const;
};

// Callers pass the credential itself; the schemas mutations take `{ body }`.
export const useCreateCredentialMutation = () => {
  const [trigger, result] = useSchemasSaveUserCredentialMutation();
  const wrappedTrigger = (credential: unknown) => trigger({ body: credential });
  return [wrappedTrigger, result] as const;
};

export const useUpdateCredentialMutation = () => {
  const [trigger, result] = useSchemasUpdateUserCredentialMutation();
  const wrappedTrigger = (credential: unknown) => trigger({ body: credential });
  return [wrappedTrigger, result] as const;
};
