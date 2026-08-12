import {
  mesheryApi,
  useDeleteUserCredentialMutation as useSchemasDeleteUserCredentialMutation,
  useGetCredentialByIdQuery as useSchemasGetCredentialByIdQuery,
  useGetUserCredentialsQuery as useSchemasGetUserCredentialsQuery,
  useSaveUserCredentialMutation as useSchemasSaveUserCredentialMutation,
  useUpdateUserCredentialMutation as useSchemasUpdateUserCredentialMutation,
} from '@meshery/schemas/mesheryApi';

// Every credential endpoint is now the schemas-generated one; the hooks below
// only adapt call shapes (bare id / bare credential) to the generated args.
//
// `deleteCredential` used to be re-declared here because the generated
// `deleteUserCredential` sends `?credentialId=` while the server read only
// `credential_id`, so the generated endpoint resolved a nil UUID and deleted
// nothing. The fix belonged in the server, not in a fork of the request:
// DeleteUserCredential now reads the canonical camelCase param (keeping
// `credential_id` as a legacy fallback) and rejects an unusable id outright.

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

export const useDeleteCredentialMutation = () => {
  const [trigger, result] = useSchemasDeleteUserCredentialMutation();
  const wrappedTrigger = (credentialId: string) => trigger({ credentialId });
  return [wrappedTrigger, result] as const;
};

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
