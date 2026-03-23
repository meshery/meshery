import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

const mesheryEndpointPrefix = process.env.RTK_MESHERY_ENDPOINT_PREFIX ?? '';

const normalizeMesheryRequest = (request) => {
  if (mesheryEndpointPrefix || !request) {
    return request;
  }

  const normalizedRequest = typeof request === 'string' ? { url: request } : request;
  const requestUrl = normalizedRequest?.url;

  if (typeof requestUrl !== 'string' || !requestUrl) {
    return request;
  }

  if (
    /^https?:\/\//.test(requestUrl) ||
    requestUrl === 'api' ||
    requestUrl.startsWith('api/') ||
    requestUrl === '/api' ||
    requestUrl.startsWith('/api/') ||
    requestUrl.startsWith('extensions/api/') ||
    requestUrl.startsWith('/extensions/api/') ||
    requestUrl.startsWith('/user/login') ||
    requestUrl.startsWith('/user/logout') ||
    requestUrl.startsWith('/provider')
  ) {
    return request;
  }

  const absoluteUrl = requestUrl.startsWith('/') ? requestUrl : `/${requestUrl}`;
  const rewrittenUrl =
    absoluteUrl === '/evaluate' ? '/api/meshmodels/relationships/evaluate' : `/api${absoluteUrl}`;

  if (typeof request === 'string') {
    return rewrittenUrl;
  }

  return {
    ...normalizedRequest,
    url: rewrittenUrl,
  };
};

const baseQuery = fetchBaseQuery({
  baseUrl: mesheryEndpointPrefix,
  credentials: 'include',
});

export const api = createApi({
  reducerPath: 'mesheryRtkSchemasApi',
  baseQuery: (request, apiContext, extraOptions) =>
    baseQuery(normalizeMesheryRequest(request), apiContext, extraOptions),
  tagTypes: [],
  endpoints: () => ({}),
});
