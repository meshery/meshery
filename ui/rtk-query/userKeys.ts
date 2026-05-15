import {
  mesheryApi,
  useGetUserKeysQuery as useSchemasGetUserKeysQuery,
} from '@meshery/schemas/mesheryApi';

export const useGetUserKeysQuery = (queryArgs, options) =>
  useSchemasGetUserKeysQuery(
    {
      orgId: queryArgs?.orgId,
    },
    options,
  );

export const useLazyGetUserKeysQuery = () => {
  const [trigger, result, lastPromiseInfo] = mesheryApi.endpoints.getUserKeys.useLazyQuery();

  const wrappedTrigger = (queryArgs, preferCacheValue) =>
    trigger(
      {
        orgId: queryArgs?.orgId,
      },
      preferCacheValue,
    );

  return [wrappedTrigger, result, lastPromiseInfo] as const;
};
