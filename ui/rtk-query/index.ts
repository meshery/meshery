import { mesheryApi } from '@meshery/schemas/mesheryApi';

const PASSTHROUGH_PREFIXES = [
  'http://',
  'https://',
  '//',
  '/api',
  'api/',
  '/evaluate',
  'evaluate',
  '/extensions/api',
  'extensions/api',
];

export const normalizeMesheryRequest = (request: string | Record<string, any>) => {
  if (typeof request === 'string') {
    return normalizeMesheryURL(request);
  }

  if (!request || typeof request !== 'object' || !('url' in request)) {
    return request;
  }

  return {
    ...request,
    url: normalizeMesheryURL(String(request.url)),
  };
};

const normalizeMesheryURL = (url: string) => {
  if (!url || PASSTHROUGH_PREFIXES.some((prefix) => url.startsWith(prefix))) {
    return url;
  }

  return url.startsWith('/') ? `/api${url}` : `/api/${url}`;
};

const wrapQuery =
  (query: (...args: any[]) => any) =>
  (...args: any[]) =>
    normalizeMesheryRequest(query(...args));

const wrapQueryFn =
  (queryFn: (...args: any[]) => any) => (arg: any, api: any, extraOptions: any, baseQuery: any) =>
    queryFn(arg, api, extraOptions, (baseQueryArg: any) =>
      baseQuery(normalizeMesheryRequest(baseQueryArg)),
    );

const wrapEndpointDefinition = (definition: Record<string, any>) => {
  if (!definition || typeof definition !== 'object') {
    return definition;
  }

  const wrappedDefinition = { ...definition };

  if (typeof definition.query === 'function') {
    wrappedDefinition.query = wrapQuery(definition.query);
  }

  if (typeof definition.queryFn === 'function') {
    wrappedDefinition.queryFn = wrapQueryFn(definition.queryFn);
  }

  return wrappedDefinition;
};

const wrapBuilder = (builder: Record<string, any>) => ({
  ...builder,
  query: (definition: Record<string, any>) => builder.query(wrapEndpointDefinition(definition)),
  mutation: (definition: Record<string, any>) =>
    builder.mutation(wrapEndpointDefinition(definition)),
  ...(typeof builder.infiniteQuery === 'function'
    ? {
        infiniteQuery: (definition: Record<string, any>) =>
          builder.infiniteQuery(wrapEndpointDefinition(definition)),
      }
    : {}),
});

const wrapApi = (baseApi: typeof mesheryApi): typeof mesheryApi =>
  Object.assign(Object.create(baseApi), baseApi, {
    injectEndpoints: (options: Parameters<typeof baseApi.injectEndpoints>[0]) =>
      wrapApi(
        baseApi.injectEndpoints({
          ...options,
          endpoints: (builder) => options.endpoints(wrapBuilder(builder)),
        }),
      ),
    enhanceEndpoints: (options: Parameters<typeof baseApi.enhanceEndpoints>[0]) =>
      wrapApi(baseApi.enhanceEndpoints(options)),
  });

export const api = wrapApi(mesheryApi);
